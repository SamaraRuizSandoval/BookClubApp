-- +goose Up
-- +goose StatementBegin
CREATE TYPE user_book_status AS ENUM ('reading', 'completed', 'wishlist');

CREATE TABLE IF NOT EXISTS user_books (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    status user_book_status NOT NULL DEFAULT 'wishlist',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    started_at TIMESTAMPTZ,
    completed_at TIMESTAMPTZ,

    pages_read INT,
    percentage_read NUMERIC(5,2),

    progress_updated_at TIMESTAMPTZ,

    CONSTRAINT pages_read_valid CHECK (pages_read >= 0),
    CONSTRAINT percentage_valid CHECK (percentage_read >= 0 AND percentage_read <= 100),

    -- Allow either pages_read OR percentage_read OR both empty
    CONSTRAINT progress_exclusive_or CHECK (
        (pages_read IS NOT NULL AND percentage_read IS NULL)
        OR (pages_read IS NULL AND percentage_read IS NOT NULL)
        OR (pages_read IS NULL AND percentage_read IS NULL)
    ),

    CONSTRAINT user_book_unique UNIQUE (user_id, book_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_books;
DROP TYPE IF EXISTS user_book_status;
-- +goose StatementEnd