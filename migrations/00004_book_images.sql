-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS book_images (
    id BIGSERIAL PRIMARY KEY, 
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    thumbnail TEXT
    small TEXT
    medium TEXT
    large TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS book_images;
-- +goose StatementEnd