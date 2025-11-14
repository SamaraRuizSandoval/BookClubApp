package store

import (
	"database/sql"
	"fmt"
	"time"
)

type ChapterComment struct {
	ID        int64     `json:"id"`
	Body      string    `json:"body"`
	UserID    int64     `json:"user_id"`
	ChapterID int64     `json:"chapter_id"`
	User      *User     `json:"user,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostgresChapterCommentStore struct {
	db *sql.DB
}

func NewPostgresChapterCommentStore(db *sql.DB) *PostgresChapterCommentStore {
	return &PostgresChapterCommentStore{db: db}
}

type ChapterCommentStore interface {
	AddComment(comment *ChapterComment, chapterID int64, userID int64) (*ChapterComment, error)
	UpdateComment(comment *ChapterComment) error
	GetCommentByID(id int64) (*ChapterComment, error)
	DeleteCommentByID(id int64) error
}

func (cs *PostgresChapterCommentStore) AddComment(comment *ChapterComment, chapterID int64, userID int64) (*ChapterComment, error) {
	err := cs.db.QueryRow(`
		INSERT INTO comments (body, user_id, chapter_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`,
		comment.Body, userID, chapterID,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)
	if err != nil {
		return nil, err
	}

	comment.UserID = userID
	comment.ChapterID = chapterID

	return comment, nil
}

func (cs *PostgresChapterCommentStore) UpdateComment(comment *ChapterComment) error {
	err := cs.db.QueryRow(`
		UPDATE comments
		SET body = $1,
		    updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
		RETURNING updated_at`,
		comment.Body, comment.ID,
	).Scan(&comment.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return sql.ErrNoRows
		}
		return err
	}

	return nil
}

func (cs *PostgresChapterCommentStore) GetCommentByID(id int64) (*ChapterComment, error) {
	comment := &ChapterComment{
		User: &User{},
	}
	err := cs.db.QueryRow(`
        SELECT c.id, c.body, c.user_id, c.chapter_id, c.created_at, c.updated_at, u.id, u.username, u.email, u.role
        FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.id = $1`, id).Scan(
		&comment.ID,
		&comment.Body,
		&comment.UserID,
		&comment.ChapterID,
		&comment.CreatedAt,
		&comment.UpdatedAt,
		&comment.User.ID,
		&comment.User.Username,
		&comment.User.Email,
		&comment.User.Role,
	)
	if err != nil {
		return nil, err
	}

	return comment, nil
}

func (cs *PostgresChapterCommentStore) DeleteCommentByID(id int64) error {
	res, err := cs.db.Exec(`
		DELETE FROM comments
		WHERE id = $1`,
		id,
	)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rows == 0 {
		return sql.ErrNoRows
	}

	return nil
}
