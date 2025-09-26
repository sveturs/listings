-- Откат миграции: удаляем добавленные маппинги

DELETE FROM category_ai_mappings
WHERE (ai_domain, product_type, category_id) IN (
    -- Автомобильные категории
    ('automotive', 'car', 1301),
    ('automotive', 'vehicle', 1301),
    ('automotive', 'automobile', 1301),
    ('automotive', 'auto', 1301),
    ('automotive', 'sedan', 1301),
    ('automotive', 'suv', 1301),
    ('automotive', 'minivan', 1301),
    ('automotive', 'hatchback', 1301),
    ('automotive', 'coupe', 1301),
    ('automotive', 'wagon', 1301),
    ('automotive', 'pickup', 1301),
    ('automotive', 'crossover', 1301),
    ('automotive', 'convertible', 1301),
    ('automotive', 'van', 1301),
    ('automotive', 'electric-car', 10170),
    ('automotive', 'hybrid', 10171),
    ('automotive', 'luxury-car', 10172),
    ('automotive', 'sports-car', 10173),
    ('automotive', 'city-car', 10176),
    ('automotive', 'motorcycle', 1302),
    ('automotive', 'scooter', 1302),
    ('automotive', 'moped', 1302),
    ('automotive', 'auto-parts', 1303),
    ('automotive', 'spare-parts', 1303),
    ('automotive', 'car-parts', 1303),
    ('automotive', 'tires', 1303),
    ('automotive', 'wheels', 1303),
    ('automotive', 'tires-wheels', 1303),

    -- Другие категории
    ('electronics', 'smartphone', 1001),
    ('electronics', 'laptop', 1001),
    ('electronics', 'computer', 1001),
    ('electronics', 'tablet', 1001),
    ('electronics', 'tv', 1001),
    ('electronics', 'camera', 1001),
    ('electronics', 'headphones', 1001),
    ('electronics', 'smartwatch', 1001),
    ('electronics', 'gaming-console', 1001),
    ('fashion', 'clothing', 1002),
    ('fashion', 'shoes', 1002),
    ('fashion', 'accessories', 1002),
    ('fashion', 'bags', 1002),
    ('fashion', 'jewelry', 1002),
    ('fashion', 'watches', 1002),
    ('real-estate', 'apartment', 1004),
    ('real-estate', 'house', 1004),
    ('real-estate', 'land', 1004),
    ('real-estate', 'commercial', 1004),
    ('real-estate', 'garage', 1004),
    ('home-garden', 'furniture', 1005),
    ('home-garden', 'appliances', 1005),
    ('home-garden', 'tools', 1005),
    ('home-garden', 'garden', 1005),
    ('home-garden', 'decor', 1005),
    ('sports', 'equipment', 1010),
    ('sports', 'fitness', 1010),
    ('sports', 'outdoor', 1010),
    ('sports', 'bicycle', 1010),
    ('sports', 'camping', 1010)
);

-- Восстанавливаем неправильные маппинги (для полного отката)
INSERT INTO category_ai_mappings (ai_domain, product_type, category_id, weight, success_count, failure_count, is_active)
VALUES
    ('automotive', 'car', 1002, 0.95, 1, 0, true),
    ('automotive', 'vehicle', 1002, 0.90, 0, 0, true)
ON CONFLICT (ai_domain, product_type, category_id) DO NOTHING;