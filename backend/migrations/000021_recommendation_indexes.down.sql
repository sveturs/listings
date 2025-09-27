-- Удаление индексов для рекомендаций

DROP INDEX IF EXISTS idx_listings_trending;
DROP INDEX IF EXISTS idx_listings_similar;
DROP INDEX IF EXISTS idx_view_history_user_date;
DROP INDEX IF EXISTS idx_view_history_listing;
DROP INDEX IF EXISTS idx_listings_popular;
DROP INDEX IF EXISTS idx_listings_city;
DROP INDEX IF EXISTS idx_listings_featured;
DROP INDEX IF EXISTS idx_view_history_collaborative;