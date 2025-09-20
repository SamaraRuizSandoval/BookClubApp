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
	Chapters      []Chapter  `json:"chapters"`
}

type BookImages struct {
	ThumbnailUrl *string `json:"thumbnail_url"`
	SmallUrl     *string `json:"small_url"`
	MediumUrl    *string `json:"medium_url"`
	LargeUrl     *string `json:"large_url"`
}

type Chapter struct {
	Number int    `json:"number"`
	Title  string `json:"title"`
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
	UpdateBook(*Book) error
}

func (pg *PostgresBookStore) AddBook(book *Book) (_ *Book, err error) {
	tx, err := pg.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		// Rollback will be a no-op if the tx is already committed.
		if rbErr := tx.Rollback(); rbErr != nil && rbErr != sql.ErrTxDone {
			log.Printf("failed to rollback transaction: %v", rbErr)
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

	for _, ch := range book.Chapters {
		_, err := tx.Exec(`
			INSERT INTO chapters (book_id, number, title) 
			VALUES ($1, $2, $3)`,
			bookID, ch.Number, ch.Title,
		)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (pg *PostgresBookStore) GetBookByID(id int64) (*Book, error) {
	book := &Book{}

	err := pg.db.QueryRow(`
        SELECT b.id, b.title, b.published_date, b.description, b.page_count, b.isbn_13, b.isbn_10, p.name
        FROM books b
        JOIN publishers p ON b.publisher_id = p.id
        WHERE b.id = $1`, id).Scan(
		&book.ID,
		&book.Title,
		&book.PublishedDate,
		&book.Description,
		&book.PageCount,
		&book.ISBN13,
		&book.ISBN10,
		&book.Publisher,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	rows, err := pg.db.Query(`
        SELECT a.name
        FROM authors a
        JOIN book_authors ba ON a.id = ba.author_id
        WHERE ba.book_id = $1
    `, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var authors []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		authors = append(authors, name)
	}
	book.Authors = authors

	var images BookImages
	err = pg.db.QueryRow(`
        SELECT thumbnail_url, small_url, medium_url, large_url
        FROM book_images
        WHERE book_id = $1
	`, id).Scan(&images.ThumbnailUrl, &images.SmallUrl, &images.MediumUrl, &images.LargeUrl)

	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}

	chapterRows, err := pg.db.Query(`
        SELECT number, title
        FROM chapters
        WHERE book_id = $1
        ORDER BY number
    `, id)
	if err != nil {
		if err != sql.ErrNoRows {
			return nil, err
		}
	}
	defer chapterRows.Close()

	var chapters []Chapter
	for chapterRows.Next() {
		var ch Chapter
		if err := chapterRows.Scan(&ch.Number, &ch.Title); err != nil {
			return nil, err
		}
		chapters = append(chapters, ch)
	}
	book.Chapters = chapters

	return book, nil
}

func (pg *PostgresBookStore) UpdateBook(book *Book) error {
	return nil
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
