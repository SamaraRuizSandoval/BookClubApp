-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS book_images (
    id BIGSERIAL PRIMARY KEY, 
    book_id BIGINT NOT NULL UNIQUE REFERENCES books(id) ON DELETE CASCADE,
    thumbnail_url TEXT,
    small_url TEXT,
    medium_url TEXT,
    large_url TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS book_images;
-- +goose StatementEnd