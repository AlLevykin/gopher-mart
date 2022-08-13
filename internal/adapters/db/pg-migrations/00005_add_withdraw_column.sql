-- +goose Up
ALTER TABLE withdraw ADD COLUMN processed timestamp DEFAULT NOW();
-- +goose Down
ALTER TABLE withdraw DROP COLUMN processed;