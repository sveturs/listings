-- Down миграция:
-- /backend/migrations/0029_create_reviews.down.sql
DROP TRIGGER IF EXISTS update_reviews_updated_at ON reviews;
DROP TRIGGER IF EXISTS update_review_responses_updated_at ON review_responses;
DROP FUNCTION IF EXISTS update_updated_at_column;
DROP FUNCTION IF EXISTS calculate_entity_rating;
DROP TABLE IF EXISTS review_votes;
DROP TABLE IF EXISTS review_responses;
DROP TABLE IF EXISTS reviews;