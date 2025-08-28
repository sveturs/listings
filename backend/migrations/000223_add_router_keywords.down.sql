-- Удаляем добавленные ключевые слова для роутеров
DELETE FROM category_keywords 
WHERE category_id = 2007 
  AND source = 'manual'
  AND keyword IN (
    -- Русские
    'роутер', 'маршрутизатор', 'вайфай', 'wifi', 'модем', 'интернет', 'сетевое',
    'tp-link', 'asus', 'd-link', 'mikrotik',
    -- Английские
    'router', 'wireless', 'modem', 'network', 'internet', 'netgear',
    -- Сербские
    'ruter', 'vajfaj', 'mreža'
  );