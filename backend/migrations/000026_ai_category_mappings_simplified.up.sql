-- AI маппинги для существующих категорий

-- Добавляем подкатегории для Hobbies & Entertainment (1015)
INSERT INTO marketplace_categories (name, slug, parent_id, icon, level, created_at) VALUES
('Игрушки', 'toys', 1015, 'toy-brick', 1, NOW()),
('Пазлы', 'puzzles', 1015, 'puzzle-piece', 1, NOW()),
('Настольные игры', 'board-games', 1015, 'dice', 1, NOW()),
('Коллекционирование', 'collectibles', 1015, 'star', 1, NOW())
ON CONFLICT (slug) DO NOTHING;

-- Добавляем категории для лучшей детекции
INSERT INTO marketplace_categories (name, slug, parent_id, icon, level, created_at) VALUES
('Строительные материалы', 'construction-materials', 1007, 'hammer', 1, NOW()),
('Природные материалы', 'natural-materials', NULL, 'tree', 0, NOW())
ON CONFLICT (slug) DO NOTHING;

-- AI маппинги для доменов и типов продуктов (используем только существующие категории)
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, success_count, failure_count, is_active) VALUES
-- Electronics mappings (1001 и подкатегории)
('electronics', 'laptop', 1102, 1.00, 100, 5, TRUE),
('electronics', 'smartphone', 1101, 1.00, 150, 3, TRUE),
('electronics', 'computer', 1102, 0.98, 120, 6, TRUE),
('electronics', 'tv', 1103, 0.98, 85, 2, TRUE),
('electronics', 'audio', 1103, 0.95, 70, 8, TRUE),
('electronics', 'appliance', 1104, 0.96, 90, 7, TRUE),
('electronics', 'gaming', 1105, 0.97, 110, 4, TRUE),
('electronics', 'camera', 1106, 0.95, 60, 5, TRUE),
('electronics', 'smarthome', 1107, 0.94, 45, 3, TRUE),
('electronics', 'accessories', 1108, 0.92, 200, 18, TRUE),

-- Entertainment mappings (1015)
('entertainment', 'puzzle', 1015, 0.95, 150, 8, TRUE),
('entertainment', 'toy', 1015, 0.94, 200, 12, TRUE),
('entertainment', 'game', 1015, 0.93, 180, 14, TRUE),
('entertainment', 'hobby', 1015, 0.92, 120, 10, TRUE),
('entertainment', 'collectible', 1015, 0.91, 80, 8, TRUE),
('entertainment', 'book', 1012, 0.95, 300, 15, TRUE),
('entertainment', 'music', 1016, 0.96, 100, 4, TRUE),

-- Automotive mappings (1003, 1301-1303)
('automotive', 'car', 1301, 0.98, 500, 10, TRUE),
('automotive', 'motorcycle', 1302, 0.97, 200, 6, TRUE),
('automotive', 'parts', 1303, 0.95, 400, 20, TRUE),
('automotive', 'accessories', 1303, 0.90, 150, 15, TRUE),

-- Real estate mappings (1004, 1401-1404)
('real-estate', 'apartment', 1401, 0.99, 800, 8, TRUE),
('real-estate', 'house', 1402, 0.99, 600, 6, TRUE),
('real-estate', 'land', 1403, 0.98, 300, 6, TRUE),
('real-estate', 'commercial', 1404, 0.97, 150, 5, TRUE),

-- Fashion mappings (1002)
('fashion', 'clothing', 1002, 0.96, 700, 28, TRUE),
('fashion', 'shoes', 1002, 0.95, 550, 28, TRUE),
('fashion', 'accessories', 1002, 0.94, 400, 24, TRUE),
('fashion', 'jewelry', 1002, 0.93, 250, 18, TRUE),

-- Home & Garden mappings (1005, 1501, 1504)
('home-garden', 'furniture', 1501, 0.97, 450, 14, TRUE),
('home-garden', 'appliance', 1104, 0.96, 380, 15, TRUE),
('home-garden', 'decoration', 1005, 0.93, 280, 20, TRUE),
('home-garden', 'garden', 1005, 0.95, 200, 10, TRUE),
('home-garden', 'materials', 1504, 0.94, 150, 9, TRUE),

-- Agriculture mappings (1006, 1601-1604)
('agriculture', 'machinery', 1601, 0.97, 120, 4, TRUE),
('agriculture', 'seeds', 1602, 0.96, 100, 4, TRUE),
('agriculture', 'livestock', 1603, 0.98, 80, 2, TRUE),
('agriculture', 'produce', 1604, 0.95, 150, 8, TRUE),

-- Industry mappings (1007)
('industrial', 'machinery', 1007, 0.95, 90, 5, TRUE),
('industrial', 'equipment', 1007, 0.94, 110, 7, TRUE),
('industrial', 'materials', 1007, 0.93, 130, 10, TRUE),
('construction', 'materials', 1504, 0.96, 180, 7, TRUE),
('construction', 'tools', 1007, 0.95, 150, 8, TRUE),

-- Services mappings (1009)
('services', 'professional', 1009, 0.94, 200, 12, TRUE),
('services', 'transport', 1009, 0.93, 150, 10, TRUE),
('services', 'repair', 1009, 0.92, 180, 15, TRUE),

-- Sports & Recreation mappings (1010)
('sports-recreation', 'fitness', 1010, 0.96, 140, 6, TRUE),
('sports-recreation', 'outdoor', 1010, 0.95, 120, 6, TRUE),
('sports-recreation', 'sports', 1010, 0.97, 180, 5, TRUE),

-- Pets (1011)
('pets', 'animal', 1011, 0.98, 200, 4, TRUE),
('pets', 'supplies', 1011, 0.96, 150, 6, TRUE),

-- Jobs (1018)
('jobs', 'employment', 1018, 0.97, 100, 3, TRUE),
('jobs', 'vacancy', 1018, 0.96, 80, 3, TRUE),

-- Education (1019)
('education', 'courses', 1019, 0.96, 90, 4, TRUE),
('education', 'training', 1019, 0.95, 70, 4, TRUE),

-- Events (1020)
('events', 'tickets', 1020, 0.97, 120, 4, TRUE),
('events', 'concert', 1020, 0.96, 80, 3, TRUE),

-- Special mappings
('antiques', 'art', 1017, 0.95, 60, 3, TRUE),
('antiques', 'vintage', 1017, 0.94, 50, 3, TRUE),

-- Nature (добавим после создания категории)
('nature', 'natural', 1005, 0.90, 40, 4, TRUE),
('nature', 'plant', 1005, 0.91, 60, 5, TRUE),

-- Общие маппинги для fallback
('other', 'unknown', 1001, 0.10, 10, 100, TRUE),
('other', 'general', 1001, 0.10, 10, 100, TRUE)
ON CONFLICT (ai_domain, product_type, category_id) DO UPDATE SET
    weight = EXCLUDED.weight,
    success_count = EXCLUDED.success_count,
    failure_count = EXCLUDED.failure_count,
    is_active = EXCLUDED.is_active;

-- Добавляем ключевые слова с весами для улучшения детекции
INSERT INTO category_keyword_weights (keyword, category_id, weight, occurrence_count, success_rate, language) VALUES
-- Пазлы и игрушки (1015)
('пазл', 1015, 1.8, 150, 0.98, 'ru'),
('puzzle', 1015, 1.8, 120, 0.97, 'en'),
('slagalica', 1015, 1.8, 100, 0.96, 'sr'),
('игрушка', 1015, 1.7, 200, 0.96, 'ru'),
('toy', 1015, 1.7, 180, 0.97, 'en'),
('igračka', 1015, 1.7, 160, 0.95, 'sr'),
('игра', 1015, 1.6, 150, 0.93, 'ru'),
('game', 1015, 1.6, 140, 0.94, 'en'),

-- Электроника (1001)
('телефон', 1101, 1.9, 300, 0.98, 'ru'),
('phone', 1101, 1.9, 280, 0.98, 'en'),
('telefon', 1101, 1.9, 260, 0.97, 'sr'),
('компьютер', 1102, 1.9, 250, 0.97, 'ru'),
('computer', 1102, 1.9, 240, 0.97, 'en'),
('laptop', 1102, 1.8, 200, 0.96, 'en'),

-- Автомобили (1301-1303)
('автомобиль', 1301, 1.9, 400, 0.99, 'ru'),
('car', 1301, 1.9, 380, 0.99, 'en'),
('automobil', 1301, 1.9, 360, 0.98, 'sr'),
('мотоцикл', 1302, 1.8, 150, 0.97, 'ru'),
('motorcycle', 1302, 1.8, 140, 0.97, 'en'),
('запчасти', 1303, 1.7, 300, 0.95, 'ru'),
('parts', 1303, 1.7, 280, 0.94, 'en'),

-- Недвижимость (1401-1404)
('квартира', 1401, 1.9, 500, 0.99, 'ru'),
('apartment', 1401, 1.9, 480, 0.99, 'en'),
('stan', 1401, 1.9, 460, 0.98, 'sr'),
('дом', 1402, 1.9, 400, 0.98, 'ru'),
('house', 1402, 1.9, 380, 0.98, 'en'),
('kuća', 1402, 1.9, 360, 0.97, 'sr'),

-- Строительные материалы (1504, 1007)
('песок', 1504, 1.8, 100, 0.96, 'ru'),
('цемент', 1504, 1.8, 90, 0.95, 'ru'),
('кирпич', 1504, 1.7, 80, 0.94, 'ru'),
('плитка', 1504, 1.6, 120, 0.93, 'ru'),

-- Природные материалы
('желудь', 1005, 1.9, 30, 0.95, 'ru'),
('acorn', 1005, 1.9, 25, 0.95, 'en'),
('природный', 1005, 1.4, 60, 0.90, 'ru'),
('natural', 1005, 1.4, 55, 0.89, 'en')
ON CONFLICT (keyword, category_id, language) DO UPDATE SET
    weight = GREATEST(category_keyword_weights.weight, EXCLUDED.weight),
    occurrence_count = category_keyword_weights.occurrence_count + EXCLUDED.occurrence_count,
    success_rate = (category_keyword_weights.success_rate * category_keyword_weights.occurrence_count +
                   EXCLUDED.success_rate * EXCLUDED.occurrence_count) /
                   (category_keyword_weights.occurrence_count + EXCLUDED.occurrence_count);

-- Создаем активный эксперимент для A/B тестирования
INSERT INTO category_detection_experiments (
    experiment_name,
    algorithm_a,
    algorithm_b,
    traffic_split,
    is_active
) VALUES (
    'enhanced_ai_mapping_v2',
    'stable_v1',
    'experimental_v2',
    0.10,
    TRUE
) ON CONFLICT (experiment_name) DO NOTHING;

-- Добавляем view для мониторинга точности
CREATE OR REPLACE VIEW category_detection_accuracy AS
SELECT
    DATE(created_at) as date,
    algorithm_version,
    COUNT(*) as total_detections,
    SUM(CASE WHEN user_confirmed THEN 1 ELSE 0 END) as confirmed,
    ROUND(100.0 * SUM(CASE WHEN user_confirmed THEN 1 ELSE 0 END) / NULLIF(COUNT(*), 0), 2) as accuracy_percent,
    AVG(confidence_score) as avg_confidence,
    PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY processing_time_ms) as median_time_ms,
    PERCENTILE_CONT(0.95) WITHIN GROUP (ORDER BY processing_time_ms) as p95_time_ms,
    PERCENTILE_CONT(0.99) WITHIN GROUP (ORDER BY processing_time_ms) as p99_time_ms
FROM category_detection_feedback
WHERE created_at > NOW() - INTERVAL '30 days'
GROUP BY DATE(created_at), algorithm_version
ORDER BY date DESC, algorithm_version;

-- Создаем функцию для автоматической очистки старого кэша
CREATE OR REPLACE FUNCTION cleanup_detection_cache() RETURNS void AS $$
BEGIN
    DELETE FROM category_detection_cache WHERE expires_at < NOW();
END;
$$ LANGUAGE plpgsql;

-- Создаем функцию для анализа точности по категориям
CREATE OR REPLACE FUNCTION get_category_accuracy_report(days_back INTEGER DEFAULT 7)
RETURNS TABLE(
    category_id INTEGER,
    category_name VARCHAR,
    total_detections BIGINT,
    correct_detections BIGINT,
    accuracy_percent NUMERIC
) AS $$
BEGIN
    RETURN QUERY
    SELECT
        c.id,
        c.name,
        COUNT(f.id) as total_detections,
        SUM(CASE WHEN f.detected_category_id = f.correct_category_id THEN 1 ELSE 0 END) as correct_detections,
        ROUND(100.0 * SUM(CASE WHEN f.detected_category_id = f.correct_category_id THEN 1 ELSE 0 END) /
              NULLIF(COUNT(f.id), 0), 2) as accuracy_percent
    FROM marketplace_categories c
    LEFT JOIN category_detection_feedback f ON f.detected_category_id = c.id
    WHERE f.created_at > NOW() - (days_back || ' days')::INTERVAL
        AND f.user_confirmed = TRUE
    GROUP BY c.id, c.name
    HAVING COUNT(f.id) > 0
    ORDER BY accuracy_percent DESC, total_detections DESC;
END;
$$ LANGUAGE plpgsql;