-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS authors (
    id BIGSERIAL PRIMARY KEY, 
    name VARCHAR(150) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS authors;
-- +goose StatementEnd