-- Проверка данных для endpoint /api/v1/analytics/metrics/items

-- 1. Общее количество событий
SELECT COUNT(*) as total_events FROM user_behavior_events;

-- 2. Количество событий по типам
SELECT event_type, COUNT(*) as count 
FROM user_behavior_events 
GROUP BY event_type 
ORDER BY count DESC;

-- 3. Количество событий с item_id
SELECT COUNT(*) as events_with_items 
FROM user_behavior_events 
WHERE item_id IS NOT NULL;

-- 4. События по item_type
SELECT item_type, COUNT(*) as count 
FROM user_behavior_events 
WHERE item_id IS NOT NULL 
GROUP BY item_type;

-- 5. Примеры событий с товарами
SELECT event_type, item_id, item_type, position, created_at 
FROM user_behavior_events 
WHERE item_id IS NOT NULL 
ORDER BY created_at DESC 
LIMIT 10;

-- 6. Проверка метрик для конкретных товаров
WITH item_events AS (
    SELECT 
        item_id,
        item_type,
        event_type,
        position,
        created_at
    FROM user_behavior_events
    WHERE item_id IS NOT NULL
        AND created_at >= NOW() - INTERVAL '7 days'
)
SELECT 
    item_id,
    item_type,
    COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END) as views,
    COUNT(CASE WHEN event_type = 'result_clicked' THEN 1 END) as clicks,
    COUNT(CASE WHEN event_type = 'item_purchased' THEN 1 END) as purchases,
    CASE 
        WHEN COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END) > 0 
        THEN COUNT(CASE WHEN event_type = 'result_clicked' THEN 1 END)::float / 
             COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END)
        ELSE 0 
    END as ctr,
    CASE 
        WHEN COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END) > 0 
        THEN COUNT(CASE WHEN event_type = 'item_purchased' THEN 1 END)::float / 
             COUNT(CASE WHEN event_type = 'item_viewed' THEN 1 END)
        ELSE 0 
    END as conversion_rate
FROM item_events
GROUP BY item_id, item_type
ORDER BY views DESC
LIMIT 10;