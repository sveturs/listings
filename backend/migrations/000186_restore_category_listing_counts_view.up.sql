-- Восстановление материализованного представления category_listing_counts
-- Это view используется для отображения количества объявлений в каждой категории

CREATE MATERIALIZED VIEW IF NOT EXISTS category_listing_counts AS
SELECT
    c.id AS category_id,
    COUNT(l.id) AS listing_count
FROM
    c2c_categories c
    LEFT JOIN c2c_listings l ON l.category_id = c.id AND l.status = 'active'
GROUP BY
    c.id;

-- Создаём индекс для быстрого доступа
CREATE UNIQUE INDEX IF NOT EXISTS idx_category_listing_counts_category_id
ON category_listing_counts(category_id);

-- Комментарий для документации
COMMENT ON MATERIALIZED VIEW category_listing_counts IS
'Материализованное представление с кол-вом активных объявлений по категориям. Обновляется после создания/изменения объявления.';
