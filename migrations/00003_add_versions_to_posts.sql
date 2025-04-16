-- +goose Up
-- SQL in this section is executed when the migration is applied.
ALTER TABLE posts ADD COLUMN version INT DEFAULT 0;

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
ALTER TABLE posts DROP COLUMN version;
