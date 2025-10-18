package store

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserRole string

type password struct {
	plainText *string
	hash      []byte
}

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           int       `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash password  `json:"-"`
	Role         UserRole  `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

var AnonymusUser = &User{}

func (u *User) IsAnonymus() bool {
	return u == nil || u == AnonymusUser || u.ID == 0
}

type PostgresUserStore struct {
	db *sql.DB
}

func NewPostgresUserStore(db *sql.DB) *PostgresUserStore {
	return &PostgresUserStore{
		db: db,
	}
}

type UserStore interface {
	CreateUser(*User) (*User, error)
	GetUserByUsername(username string) (*User, error)
	UpdateUser(*User) error
	GetUserToken(scope, plainTextPassword string) (*User, error)
}

func (p *password) Set(plainTextPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainTextPassword), 12)
	if err != nil {
		return err
	}

	p.plainText = &plainTextPassword
	p.hash = hash

	return nil
}

func (p *password) Matches(plainTextPassword string) (bool, error) {
	err := bcrypt.CompareHashAndPassword(p.hash, []byte(plainTextPassword))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			return false, nil
		default:
			return false, err

		}
	}

	return true, nil
}

func (us *PostgresUserStore) CreateUser(user *User) (*User, error) {
	err := us.db.QueryRow(`
        INSERT INTO users (username, email, password_hash, role) 
        VALUES ($1, $2, $3, $4)  
        RETURNING id, created_at`,
		user.Username, user.Email, user.PasswordHash.hash, user.Role,
	).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *PostgresUserStore) GetUserByUsername(username string) (*User, error) {
	user := &User{
		PasswordHash: password{},
	}

	err := us.db.QueryRow(`
        SELECT id, username, email, password_hash, role, created_at
        FROM users
		WHERE username = $1`,
		username,
	).Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.Role, &user.CreatedAt)

	if err == sql.ErrNoRows {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (us *PostgresUserStore) UpdateUser(user *User) error {
	result, err := us.db.Exec(`
        UPDATE users
		SET username = $1, email = $2
		WHERE id = $3`,
		user.Username, user.Email, user.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (us *PostgresUserStore) GetUserToken(scope, plainTextToken string) (*User, error) {
	fmt.Println("scope:", scope, "plainTextToken:", plainTextToken)

	// Hash the token string exactly as created
	tokenHash := sha256.Sum256([]byte(plainTextToken))

	query := `
	SELECT u.id, u.username, u.email, u.password_hash, u.role, u.created_at
	FROM users u
	INNER JOIN tokens t ON t.user_id = u.id
	WHERE t.hash = $1 AND t.scope = $2 AND t.expiry > $3
	`

	user := &User{PasswordHash: password{}}

	err := us.db.QueryRow(query, tokenHash[:], scope, time.Now()).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash.hash, &user.Role, &user.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}
