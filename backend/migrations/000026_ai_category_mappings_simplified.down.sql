-- Удаляем функции
DROP FUNCTION IF EXISTS get_category_accuracy_report(INTEGER);
DROP FUNCTION IF EXISTS cleanup_detection_cache();

-- Удаляем view
DROP VIEW IF EXISTS category_detection_accuracy;

-- Удаляем эксперимент
DELETE FROM category_detection_experiments WHERE experiment_name = 'enhanced_ai_mapping_v2';

-- Удаляем ключевые слова (для всех категорий которые добавили)
DELETE FROM category_keyword_weights WHERE
    (keyword IN ('пазл', 'puzzle', 'slagalica', 'игрушка', 'toy', 'igračka', 'игра', 'game') AND category_id = 1015) OR
    (keyword IN ('телефон', 'phone', 'telefon') AND category_id = 1101) OR
    (keyword IN ('компьютер', 'computer', 'laptop') AND category_id = 1102) OR
    (keyword IN ('автомобиль', 'car', 'automobil') AND category_id = 1301) OR
    (keyword IN ('мотоцикл', 'motorcycle') AND category_id = 1302) OR
    (keyword IN ('запчасти', 'parts') AND category_id = 1303) OR
    (keyword IN ('квартира', 'apartment', 'stan') AND category_id = 1401) OR
    (keyword IN ('дом', 'house', 'kuća') AND category_id = 1402) OR
    (keyword IN ('песок', 'цемент', 'кирпич', 'плитка') AND category_id = 1504) OR
    (keyword IN ('желудь', 'acorn', 'природный', 'natural') AND category_id = 1005);

-- Удаляем AI маппинги
DELETE FROM category_ai_mappings WHERE
    ai_domain IN ('electronics', 'entertainment', 'automotive', 'real-estate', 'fashion',
                  'home-garden', 'agriculture', 'industrial', 'construction', 'services',
                  'sports-recreation', 'pets', 'jobs', 'education', 'events', 'antiques',
                  'nature', 'other');

-- Удаляем добавленные категории
DELETE FROM marketplace_categories WHERE slug IN (
    'toys', 'puzzles', 'board-games', 'collectibles',
    'construction-materials', 'natural-materials'
);