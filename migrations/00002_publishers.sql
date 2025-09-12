-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS publishers (
    id BIGSERIAL PRIMARY KEY, 
    name VARCHAR(100) NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS publishers;
-- +goose StatementEnd