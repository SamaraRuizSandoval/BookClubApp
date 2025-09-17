-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS books (
    id BIGSERIAL PRIMARY KEY, 
    title VARCHAR(150) NOT NULL,
    publisher_id BIGINT NOT NULL REFERENCES publishers(id) ON DELETE CASCADE,
    published_date DATE NOT NULL,
    description TEXT, 
    page_count INT,
    isbn_13 VARCHAR(20) UNIQUE NOT NULL,
    isbn_10 VARCHAR(20)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS books;
-- +goose StatementEnd