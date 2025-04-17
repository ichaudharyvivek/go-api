-- +goose Up
-- SQL in this section is executed when the migration is applied.
-- Add a tsvector column for full-text search
ALTER TABLE posts ADD COLUMN tsv tsvector;
-- Create a GIN index on the tsvector column
CREATE INDEX idx_posts_tsv ON posts USING GIN(tsv);
-- Update the tsvector column with a combination of title and content for full-text search
UPDATE posts SET tsv = to_tsvector('english', title || ' ' || content);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
-- Drop the GIN index
DROP INDEX IF EXISTS idx_posts_tsv;
-- Drop the tsvector column
ALTER TABLE posts DROP COLUMN IF EXISTS tsv;
