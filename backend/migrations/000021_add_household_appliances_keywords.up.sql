-- Добавляем ключевые слова для бытовой техники в категорию "Dom i bašta" (ID: 1005)
-- Используем ON CONFLICT для избежания дубликатов

-- Английские ключевые слова
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, usage_count, success_rate) VALUES
-- Пылесосы
(1005, 'vacuum', 'en', 10, 'main', false, 0, 0.85),
(1005, 'cleaner', 'en', 9, 'main', false, 0, 0.85),
(1005, 'vacuum cleaner', 'en', 10, 'main', false, 0, 0.90),
(1005, 'hoover', 'en', 8, 'synonym', false, 0, 0.80),
(1005, 'dust', 'en', 6, 'context', false, 0, 0.70),
(1005, 'suction', 'en', 7, 'context', false, 0, 0.75),
(1005, 'miele', 'en', 8, 'brand', false, 0, 0.85),
(1005, 'dyson', 'en', 8, 'brand', false, 0, 0.85),
(1005, 'bosch', 'en', 8, 'brand', false, 0, 0.85),
(1005, 'samsung', 'en', 8, 'brand', false, 0, 0.85),
(1005, 'electrolux', 'en', 8, 'brand', false, 0, 0.85),
(1005, 'philips', 'en', 8, 'brand', false, 0, 0.85),
-- Общая бытовая техника
(1005, 'appliance', 'en', 9, 'main', false, 0, 0.85),
(1005, 'domestic', 'en', 8, 'synonym', false, 0, 0.80),
(1005, 'iron', 'en', 8, 'context', false, 0, 0.80),
(1005, 'washing', 'en', 8, 'context', false, 0, 0.80),
(1005, 'dishwasher', 'en', 9, 'main', false, 0, 0.85),
(1005, 'microwave', 'en', 9, 'main', false, 0, 0.85),
(1005, 'refrigerator', 'en', 9, 'main', false, 0, 0.85),
(1005, 'fridge', 'en', 9, 'synonym', false, 0, 0.85),
(1005, 'oven', 'en', 8, 'context', false, 0, 0.80),
(1005, 'stove', 'en', 8, 'context', false, 0, 0.80),
(1005, 'cooker', 'en', 8, 'context', false, 0, 0.80),
(1005, 'blender', 'en', 8, 'context', false, 0, 0.80),
(1005, 'mixer', 'en', 8, 'context', false, 0, 0.80),
(1005, 'toaster', 'en', 8, 'context', false, 0, 0.80),
(1005, 'kettle', 'en', 8, 'context', false, 0, 0.80),
(1005, 'coffee', 'en', 7, 'context', false, 0, 0.75),
(1005, 'kitchen', 'en', 8, 'context', false, 0, 0.80)
ON CONFLICT (category_id, keyword, language)
DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type,
    usage_count = category_keywords.usage_count + 1;

-- Русские ключевые слова
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, usage_count, success_rate) VALUES
-- Пылесосы
(1005, 'пылесос', 'ru', 10, 'main', false, 0, 0.90),
(1005, 'уборка', 'ru', 7, 'context', false, 0, 0.75),
(1005, 'пыль', 'ru', 6, 'context', false, 0, 0.70),
(1005, 'чистка', 'ru', 7, 'context', false, 0, 0.75),
(1005, 'мешок', 'ru', 5, 'context', false, 0, 0.65),
(1005, 'фильтр', 'ru', 5, 'context', false, 0, 0.65),
(1005, 'щетка', 'ru', 6, 'context', false, 0, 0.70),
-- Общая бытовая техника
(1005, 'техника', 'ru', 9, 'main', false, 0, 0.85),
(1005, 'бытовая', 'ru', 9, 'main', false, 0, 0.85),
(1005, 'домашняя', 'ru', 8, 'synonym', false, 0, 0.80),
(1005, 'утюг', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'стиральная', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'посудомоечная', 'ru', 9, 'main', false, 0, 0.85),
(1005, 'микроволновка', 'ru', 9, 'main', false, 0, 0.85),
(1005, 'холодильник', 'ru', 9, 'main', false, 0, 0.85),
(1005, 'плита', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'духовка', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'варочная', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'блендер', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'миксер', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'тостер', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'чайник', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'кофеварка', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'кофемашина', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'кухня', 'ru', 8, 'context', false, 0, 0.80),
(1005, 'кухонная', 'ru', 8, 'context', false, 0, 0.80)
ON CONFLICT (category_id, keyword, language)
DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type,
    usage_count = category_keywords.usage_count + 1;

-- Сербские ключевые слова
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, is_negative, usage_count, success_rate) VALUES
-- Усисивачи
(1005, 'usisivač', 'sr', 10, 'main', false, 0, 0.90),
(1005, 'čišćenje', 'sr', 7, 'context', false, 0, 0.75),
(1005, 'prašina', 'sr', 6, 'context', false, 0, 0.70),
(1005, 'vrećica', 'sr', 5, 'context', false, 0, 0.65),
(1005, 'četka', 'sr', 6, 'context', false, 0, 0.70),
-- Кућни апарати
(1005, 'aparat', 'sr', 9, 'main', false, 0, 0.85),
(1005, 'kućni', 'sr', 9, 'main', false, 0, 0.85),
(1005, 'pegla', 'sr', 8, 'context', false, 0, 0.80),
(1005, 'veš', 'sr', 8, 'context', false, 0, 0.80),
(1005, 'mašina', 'sr', 8, 'context', false, 0, 0.80),
(1005, 'mikrotalasna', 'sr', 9, 'main', false, 0, 0.85),
(1005, 'frižider', 'sr', 9, 'main', false, 0, 0.85),
(1005, 'šporet', 'sr', 8, 'context', false, 0, 0.80),
(1005, 'rerna', 'sr', 8, 'context', false, 0, 0.80),
(1005, 'blender', 'sr', 8, 'context', false, 0, 0.80),
(1005, 'mikser', 'sr', 8, 'context', false, 0, 0.80),
(1005, 'toster', 'sr', 8, 'context', false, 0, 0.80),
(1005, 'kuvalo', 'sr', 8, 'context', false, 0, 0.80),
(1005, 'kuhinja', 'sr', 8, 'context', false, 0, 0.80)
ON CONFLICT (category_id, keyword, language)
DO UPDATE SET
    weight = EXCLUDED.weight,
    keyword_type = EXCLUDED.keyword_type,
    usage_count = category_keywords.usage_count + 1;