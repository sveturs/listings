-- Индексы для улучшения производительности рекомендаций

-- Индекс для быстрого поиска trending (популярных) товаров
CREATE INDEX IF NOT EXISTS idx_listings_trending ON marketplace_listings(status, created_at DESC, views_count DESC)
WHERE status = 'active';

-- Составной индекс для поиска похожих товаров по категории и цене
CREATE INDEX IF NOT EXISTS idx_listings_similar ON marketplace_listings(category_id, status, price)
WHERE status = 'active';

-- Индексы для истории просмотров (если таблица существует)
-- CREATE INDEX IF NOT EXISTS idx_view_history_user_date ON universal_view_history(user_id, created_at DESC);
-- CREATE INDEX IF NOT EXISTS idx_view_history_listing ON universal_view_history(listing_id, interaction_type);

-- Частичный индекс для активных товаров с высокими просмотрами
CREATE INDEX IF NOT EXISTS idx_listings_popular ON marketplace_listings(views_count DESC)
WHERE status = 'active' AND views_count > 10;

-- Индекс для поиска по городу для локальных рекомендаций
CREATE INDEX IF NOT EXISTS idx_listings_city ON marketplace_listings(address_city, status)
WHERE status = 'active';

-- Индекс для featured товаров
CREATE INDEX IF NOT EXISTS idx_listings_featured ON marketplace_listings(created_at DESC)
WHERE status = 'active' AND metadata->>'is_featured' = 'true';

-- Составной индекс для collaborative filtering
-- CREATE INDEX IF NOT EXISTS idx_view_history_collaborative ON universal_view_history(category_id, user_id, created_at DESC)
-- WHERE created_at > NOW() - INTERVAL '30 days';