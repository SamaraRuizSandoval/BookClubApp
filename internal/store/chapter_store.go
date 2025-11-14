package store

import "database/sql"

type PostgresChapter struct {
	db *sql.DB
}

func NewPostgresChapterStore(db *sql.DB) *PostgresChapter {
	return &PostgresChapter{db: db}
}

type ChapterStore interface {
	GetChapterByID(id int64) (*Chapter, error)
}

func (cs *PostgresChapter) GetChapterByID(id int64) (*Chapter, error) {
	chapterInfo := &Chapter{}

	err := cs.db.QueryRow(`
		SELECT id, number, title
		FROM chapters
		WHERE id = $1`, id).Scan(
		&chapterInfo.ID,
		&chapterInfo.Number,
		&chapterInfo.Title,
	)
	if err != nil {
		return nil, err
	}

	return chapterInfo, nil
}
