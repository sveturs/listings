-- backend/migrations/0030_add_review_votes_columns.up.sql
ALTER TABLE reviews 
ADD COLUMN helpful_votes INT DEFAULT 0,
ADD COLUMN not_helpful_votes INT DEFAULT 0;