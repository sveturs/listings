-- /backend/migrations/0030_add_review_votes_columns.down.sql
ALTER TABLE reviews 
DROP COLUMN helpful_votes,
DROP COLUMN not_helpful_votes;