-- Исправление неправильно закешированных категорий для автомобилей
-- Категория 1002 (Fashion/Moda) ошибочно использовалась вместо 1301 (Lični automobili/Cars)

-- Обновляем все закешированные решения для автомобилей
UPDATE ai_category_decisions
SET category_id = 1301,
    reasoning = COALESCE(reasoning, '') || ' [FIXED: Corrected from Fashion(1002) to Cars(1301)]'
WHERE category_id = 1002
  AND (
    -- Проверяем по названию
    lower(title) LIKE '%car%'
    OR lower(title) LIKE '%auto%'
    OR lower(title) LIKE '%vehicle%'
    OR lower(title) LIKE '%volkswagen%'
    OR lower(title) LIKE '%touran%'
    OR lower(title) LIKE '%bmw%'
    OR lower(title) LIKE '%mercedes%'
    OR lower(title) LIKE '%audi%'
    OR lower(title) LIKE '%toyota%'
    OR lower(title) LIKE '%honda%'
    OR lower(title) LIKE '%ford%'
    OR lower(title) LIKE '%opel%'
    OR lower(title) LIKE '%peugeot%'
    OR lower(title) LIKE '%renault%'
    OR lower(title) LIKE '%citroen%'
    OR lower(title) LIKE '%fiat%'
    OR lower(title) LIKE '%škoda%'
    OR lower(title) LIKE '%skoda%'
    OR lower(title) LIKE '%mazda%'
    OR lower(title) LIKE '%nissan%'
    OR lower(title) LIKE '%hyundai%'
    OR lower(title) LIKE '%kia%'
    -- Проверяем по домену
    OR ai_domain = 'automotive'
    OR ai_product_type IN ('car', 'suv', 'sedan', 'minivan', 'hatchback', 'wagon', 'coupe', 'convertible', 'truck', 'van')
  );

-- Также обновляем альтернативные категории если они содержат 1002
UPDATE ai_category_decisions
SET alternative_category_ids = array_replace(alternative_category_ids, 1002, 1301)
WHERE 1002 = ANY(alternative_category_ids)
  AND (
    lower(title) LIKE '%car%'
    OR lower(title) LIKE '%auto%'
    OR lower(title) LIKE '%vehicle%'
    OR ai_domain = 'automotive'
  );

-- Удаляем совсем неправильные записи где Fashion категория для автомобилей
DELETE FROM ai_category_decisions
WHERE category_id = 1002
  AND ai_domain = 'automotive'
  AND created_at < NOW() - INTERVAL '1 hour'; -- Удаляем только старые записи