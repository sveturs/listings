-- Добавляем популярные модели SEAT
INSERT INTO car_models (make_id, name, slug, is_active, created_at, updated_at)
VALUES
    (22, 'Arona', 'arona', true, NOW(), NOW()),
    (22, 'Ateca', 'ateca', true, NOW(), NOW()),
    (22, 'Ibiza', 'ibiza', true, NOW(), NOW()),
    (22, 'Leon', 'leon', true, NOW(), NOW()),
    (22, 'Formentor', 'formentor', true, NOW(), NOW()),
    (22, 'Tarraco', 'tarraco', true, NOW(), NOW()),
    (22, 'Altea', 'altea', true, NOW(), NOW()),
    (22, 'Toledo', 'toledo', true, NOW(), NOW()),
    (22, 'Alhambra', 'alhambra', true, NOW(), NOW()),
    (22, 'Mii', 'mii', true, NOW(), NOW()),
    (22, 'Exeo', 'exeo', true, NOW(), NOW()),
    (22, 'Cordoba', 'cordoba', true, NOW(), NOW())
ON CONFLICT (make_id, slug) DO NOTHING;