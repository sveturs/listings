-- Добавляем ключевые слова для роутеров
INSERT INTO category_keywords (category_id, keyword, language, weight, keyword_type, source) VALUES
-- Русские ключевые слова
(2007, 'роутер', 'ru', 10.0, 'main', 'manual'),
(2007, 'маршрутизатор', 'ru', 9.0, 'synonym', 'manual'),
(2007, 'вайфай', 'ru', 8.0, 'attribute', 'manual'),
(2007, 'wifi', 'ru', 8.0, 'attribute', 'manual'),
(2007, 'модем', 'ru', 7.0, 'synonym', 'manual'),
(2007, 'интернет', 'ru', 6.0, 'context', 'manual'),
(2007, 'сетевое', 'ru', 5.0, 'attribute', 'manual'),
(2007, 'tp-link', 'ru', 7.0, 'brand', 'manual'),
(2007, 'asus', 'ru', 7.0, 'brand', 'manual'),
(2007, 'd-link', 'ru', 7.0, 'brand', 'manual'),
(2007, 'mikrotik', 'ru', 7.0, 'brand', 'manual'),

-- Английские ключевые слова
(2007, 'router', 'en', 10.0, 'main', 'manual'),
(2007, 'wifi', 'en', 8.0, 'attribute', 'manual'),
(2007, 'wireless', 'en', 8.0, 'attribute', 'manual'),
(2007, 'modem', 'en', 7.0, 'synonym', 'manual'),
(2007, 'network', 'en', 6.0, 'context', 'manual'),
(2007, 'internet', 'en', 6.0, 'context', 'manual'),
(2007, 'tp-link', 'en', 7.0, 'brand', 'manual'),
(2007, 'asus', 'en', 7.0, 'brand', 'manual'),
(2007, 'd-link', 'en', 7.0, 'brand', 'manual'),
(2007, 'netgear', 'en', 7.0, 'brand', 'manual'),

-- Сербские ключевые слова
(2007, 'ruter', 'sr', 10.0, 'main', 'manual'),
(2007, 'router', 'sr', 10.0, 'main', 'manual'),
(2007, 'wifi', 'sr', 8.0, 'attribute', 'manual'),
(2007, 'vajfaj', 'sr', 8.0, 'attribute', 'manual'),
(2007, 'modem', 'sr', 7.0, 'synonym', 'manual'),
(2007, 'internet', 'sr', 6.0, 'context', 'manual'),
(2007, 'mreža', 'sr', 5.0, 'context', 'manual')

ON CONFLICT (category_id, keyword, language) DO UPDATE SET
  weight = EXCLUDED.weight,
  keyword_type = EXCLUDED.keyword_type,
  updated_at = CURRENT_TIMESTAMP;