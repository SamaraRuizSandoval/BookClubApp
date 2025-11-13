package store

import (
	"database/sql"
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
	DeleteCommentByID(id int64) error
}

func (pcs *PostgresChapterCommentStore) AddComment(comment *ChapterComment, chapterID int64, userID int64) (*ChapterComment, error) {
	err := pcs.db.QueryRow(`
		INSERT INTO comments (body, user_id, chapter_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at`,
		comment.Body, userID, chapterID,
	).Scan(&comment.ID, &comment.CreatedAt, &comment.UpdatedAt)

	if err != nil {
		//TODO: Chapter doesn't exist error
		return nil, err
	}

	comment.UserID = userID
	comment.ChapterID = chapterID

	return comment, nil
}

func (pcs *PostgresChapterCommentStore) UpdateComment(comment *ChapterComment) error {
	//TODO
	return nil
}

func (pcs *PostgresChapterCommentStore) DeleteCommentByID(id int64) error {
	//TODO
	return nil
}
