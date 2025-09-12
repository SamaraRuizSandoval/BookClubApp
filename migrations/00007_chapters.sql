-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chapters (
    id BIGSERIAL PRIMARY KEY, 
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    number INT NOT NULL,
    title VARCHAR(150) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS chapters;
-- +goose StatementEnd