-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS book_authors (
    id BIGSERIAL PRIMARY KEY, 
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    author_id BIGINT NOT NULL REFERENCES authors(id) ON DELETE CASCADE,
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS book_authors;
-- +goose StatementEnd