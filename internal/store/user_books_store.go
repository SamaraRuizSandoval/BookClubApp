package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

type UserBook struct {
	ID                int64     `json:"id"`
	UserID            int64     `json:"user_id"`
	BookID            int64     `json:"book_id"`
	Status            string    `json:"status"` // "wishlist", "reading", "completed"
	StartedAt         *JSONDate `json:"started_at,omitempty"`
	CompletedAt       *JSONDate `json:"completed_at,omitempty"`
	PagesRead         *int      `json:"pages_read,omitempty"`
	PercentageRead    *float64  `json:"percentage_read,omitempty"`
	ProgressUpdatedAt *JSONDate `json:"progress_updated_at,omitempty"`
	UpdatedAt         JSONDate  `json:"updated_at"`
	Book              *Book     `json:"book,omitempty"`
}

type UpdateUserBookRequest struct {
	Status         *string    `json:"status"`
	PagesRead      *int       `json:"pages_read"`
	PercentageRead *float64   `json:"percentage_read"`
	CompletedAt    **JSONDate `json:"completed_at"` // pointer-to-pointer allows null explicitly
}

type BasicUserBook struct {
	ID        int64    `json:"id"`
	UserID    int64    `json:"user_id"`
	Status    string   `json:"status"`
	UpdatedAt JSONDate `json:"updated_at"`
	Book      *Book    `json:"book,omitempty"`
}

type PostgresUserBooksStore struct {
	db *sql.DB
}

func NewUserBooksStore(db *sql.DB) *PostgresUserBooksStore {
	return &PostgresUserBooksStore{db: db}
}

type UserBooksStore interface {
	GetUserBooksByUserID(userID int64, status *string, page, limit int) ([]*BasicUserBook, error)
	AddUserBook(userid, bookid int64, status string) (*UserBook, error)
	UpdateUserBook(userID, userBookID int64, req UpdateUserBookRequest) (*UserBook, error)
	DeleteUserBook(userID, userBookID int64) error
}

func (pub *PostgresUserBooksStore) GetUserBooksByUserID(userID int64, status *string, page, limit int) ([]*BasicUserBook, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}

	offset := (page - 1) * limit

	rows, err := pub.db.Query(`
        SELECT ub.id, ub.user_id, ub.status, ub.updated_at,

		jsonb_build_object(
			'id', b.id,
			'title', b.title,
			'published_date', b.published_date,
			'description', b.description,
			'page_count', b.page_count,
			'isbn_13', b.isbn_13,
			'isbn_10', b.isbn_10,
			'publisher', p.name,

			'authors',
			COALESCE(
				json_agg(DISTINCT a.name) FILTER (WHERE a.id IS NOT NULL),
				'[]'::json
			),

			'images',
			COALESCE(
				jsonb_build_object(
					'thumbnail_url', bi.thumbnail_url,
					'small_url',     bi.small_url,
					'medium_url',    bi.medium_url,
					'large_url',     bi.large_url
				),
				'{}'::jsonb
			),

			'chapters',
			COALESCE(
				json_agg(
					jsonb_build_object(
						'id',     c.id,
						'number', c.number,
						'title',  c.title
					) ORDER BY c.number
				) FILTER (WHERE c.id IS NOT NULL),
				'[]'::json
			)
		) AS book

	FROM user_books ub
	JOIN books b ON ub.book_id = b.id
	LEFT JOIN publishers p ON b.publisher_id = p.id
	LEFT JOIN book_authors ba ON b.id = ba.book_id
	LEFT JOIN authors a ON ba.author_id = a.id
	LEFT JOIN book_images bi ON b.id = bi.book_id
	LEFT JOIN chapters c ON b.id = c.book_id

	WHERE ub.user_id = $1
	AND ($4::user_book_status IS NULL OR ub.status = $4::user_book_status)

	GROUP BY 
		ub.id,
		ub.user_id,
		ub.status,
		ub.updated_at,
		b.id,
		p.name,
		bi.thumbnail_url, bi.small_url, bi.medium_url, bi.large_url

	ORDER BY ub.updated_at DESC
	LIMIT $2 OFFSET $3;
		`, userID, limit, offset, status)
	if err != nil {
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("failed to close transaction: %v", closeErr)
		}
	}()

	userBooks := []*BasicUserBook{}
	for rows.Next() {
		var ub BasicUserBook
		var bookJson []byte

		err := rows.Scan(&ub.ID, &ub.UserID, &ub.Status, &ub.UpdatedAt, &bookJson)
		if err != nil {
			return nil, err
		}

		var book Book
		if err := json.Unmarshal(bookJson, &book); err != nil {
			return nil, err
		}
		ub.Book = &book

		userBooks = append(userBooks, &ub)
	}

	return userBooks, nil
}

func (pub *PostgresUserBooksStore) AddUserBook(userid, bookid int64, status string) (*UserBook, error) {
	userBook := &UserBook{}
	err := pub.db.QueryRow(`
		INSERT INTO user_books (user_id, book_id, status)
		VALUES ($1, $2, $3)
		RETURNING id, updated_at`,
		userid, bookid, status,
	).Scan(&userBook.ID, &userBook.UpdatedAt)
	if err != nil {
		return nil, err
	}

	userBook.UserID = userid
	userBook.BookID = bookid
	userBook.Status = status
	return userBook, nil
}

func (pub *PostgresUserBooksStore) UpdateUserBook(userID, userBookID int64, req UpdateUserBookRequest) (*UserBook, error) {
	// "setClauses" collects the SQL pieces for columns that actually change.
	setClauses := []string{}
	args := []interface{}{}

	// If user sent a new status
	if req.Status != nil {
		setClauses = append(setClauses,
			fmt.Sprintf("status = $%d", len(args)+1),
		)
		args = append(args, *req.Status)
	}

	// If user sent a completed_at field
	if req.CompletedAt != nil {
		if *req.CompletedAt == nil {
			// User explicitly sent: "completed_at": null
			setClauses = append(setClauses, "completed_at = NULL")
		} else {
			// User sent a real timestamp
			setClauses = append(setClauses,
				fmt.Sprintf("completed_at = $%d", len(args)+1),
			)
			args = append(args, **req.CompletedAt)
		}
	}

	if req.Status != nil && *req.Status == "completed" {
		// Only set completed_at automatically if the user didn't set it manually
		if req.CompletedAt == nil {
			setClauses = append(setClauses, "completed_at = NOW()")
		}
	}

	// If user updated pages_read
	if req.PagesRead != nil {
		setClauses = append(setClauses,
			fmt.Sprintf("pages_read = $%d", len(args)+1),
		)
		args = append(args, *req.PagesRead)
		// We also update progress timestamp
		setClauses = append(setClauses, "progress_updated_at = NOW()")
	}

	// If user updated percentage_read
	if req.PercentageRead != nil {
		setClauses = append(setClauses,
			fmt.Sprintf("percentage_read = $%d", len(args)+1),
		)
		args = append(args, *req.PercentageRead)
		// Same: progress was updated
		setClauses = append(setClauses, "progress_updated_at = NOW()")
	}

	// User sent an empty payload (nothing to update)
	if len(setClauses) == 0 {
		return nil, fmt.Errorf("no fields to update")
	}

	// Always update "updated_at"
	setClauses = append(setClauses, "updated_at = NOW()")

	// Build the SQL dynamically
	query := fmt.Sprintf(`
        UPDATE user_books
        SET %s
        WHERE id = $%d AND user_id = $%d
        RETURNING id, user_id, book_id, status, updated_at,
                  started_at, completed_at, pages_read, percentage_read,
                  progress_updated_at
    `,
		strings.Join(setClauses, ", "),
		len(args)+1, // ID placeholder
		len(args)+2, // userID placeholder
	)

	// Add WHERE clause parameters
	args = append(args, userBookID, userID)

	// Execute query and scan result
	userBook := &UserBook{}
	err := pub.db.QueryRow(query, args...).Scan(
		&userBook.ID,
		&userBook.UserID,
		&userBook.BookID,
		&userBook.Status,
		&userBook.UpdatedAt,
		&userBook.StartedAt,
		&userBook.CompletedAt,
		&userBook.PagesRead,
		&userBook.PercentageRead,
		&userBook.ProgressUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return userBook, nil
}

func (pub *PostgresUserBooksStore) DeleteUserBook(userID, userBookID int64) error {
	// Delete only the row belonging to this user
	res, err := pub.db.Exec(`
        DELETE FROM user_books
        WHERE id = $1 AND user_id = $2
    `, userBookID, userID)
	if err != nil {
		return err
	}

	// Check if any row was actually deleted
	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
