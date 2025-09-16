package store

import (
	"database/sql"
	"time"
)

type Book struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Authors       []string   `json:"authors"`
	Publisher     string     `json:"publisher"`
	PublishedDate time.Time  `json:"published_date"`
	Description   *string    `json:"description"`
	PageCount     *int       `json:"page_count"`
	Isbn13        string     `json:"isbn_13"`
	Isbn10        *string    `json:"isbn_10"`
	Images        BookImages `json:"book_images"`
}

type BookImages struct {
	ThumbnailUrl *string `json:"thumbnail_url"`
	SmallUrl     *string `json:"small_url"`
	MediumUrl    *string `json:"medium_url"`
	LargeUrl     *string `json:"large_url"`
}

type PostgresBookStore struct {
	db *sql.DB
}

func NewPostgresBookStore(db *sql.DB) *PostgresBookStore {
	return &PostgresBookStore{db: db}
}

type BookStore interface {
	AddBook(*Book) (*Book, error)
	GetBookByID(id int64) (*Book, error)
}

func (pg *PostgresBookStore) AddBook(book *Book) (*Book, error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var publisherID int
	err = tx.QueryRow(`
        INSERT INTO publishers (name) 
        VALUES ($1) 
        ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name 
        RETURNING id`,
		book.Publisher,
	).Scan(&publisherID)
	if err != nil {
		return nil, err
	}

	var bookID int
	err = tx.QueryRow(`
        INSERT INTO books (title, publisher_id, published_date, description, page_count, isbn_13, isbn_10) 
        VALUES ($1,$2,$3,$4,$5,$6,$7) 
        RETURNING id`,
		book.Title, publisherID, book.PublishedDate, book.Description, book.PageCount, book.Isbn13, book.Isbn10,
	).Scan(&bookID)
	if err != nil {
		return nil, err
	}

	for _, author := range book.Authors {
		var authorID int
		err = tx.QueryRow(`
            INSERT INTO authors (name) 
            VALUES ($1) 
            ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name 
            RETURNING id`,
			author,
		).Scan(&authorID)
		if err != nil {
			return nil, err
		}

		_, err = tx.Exec(`INSERT INTO book_authors (book_id, author_id) VALUES ($1,$2) ON CONFLICT DO NOTHING`, bookID, authorID)
		if err != nil {
			return nil, err
		}
	}

	_, err = tx.Exec(`
    INSERT INTO book_images (book_id, thumbnail_url, small_url, medium_url, large_url)
    VALUES ($1, $2, $3, $4, $5)`,
		bookID, book.Images.ThumbnailUrl, book.Images.SmallUrl, book.Images.MediumUrl,
		book.Images.LargeUrl)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (pg *PostgresBookStore) GetBookByID(id int64) (*Book, error) {
	book := &Book{}
	return book, nil
}
