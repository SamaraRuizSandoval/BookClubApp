package store

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"
)

type JSONDate time.Time

type Book struct {
	ID            int        `json:"id"`
	Title         string     `json:"title"`
	Authors       []string   `json:"authors"`
	Publisher     string     `json:"publisher"`
	PublishedDate JSONDate   `json:"published_date"`
	Description   *string    `json:"description"`
	PageCount     *int       `json:"page_count"`
	ISBN13        string     `json:"isbn_13"`
	ISBN10        *string    `json:"isbn_10"`
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
	defer func() {
		if err := tx.Rollback(); err != nil && err != sql.ErrTxDone {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

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
		book.Title, publisherID, book.PublishedDate.ToTime(), book.Description, book.PageCount, book.ISBN13, book.ISBN10,
	).Scan(&bookID)
	if err != nil {
		return nil, err
	}
	book.ID = bookID

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

func (d *JSONDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	if s == "" {
		return nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	*d = JSONDate(t)
	return nil
}

func (d JSONDate) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, time.Time(d).Format("2006-01-02"))), nil
}

func (d JSONDate) ToTime() time.Time {
	return time.Time(d)
}
