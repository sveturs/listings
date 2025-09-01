-- Создание витрины агентства недвижимости
INSERT INTO storefronts (
    user_id, slug, name, description,
    logo_url, banner_url,
    phone, email, website,
    address, city, postal_code, country,
    latitude, longitude,
    theme, settings, seo_meta,
    is_active, is_verified, verification_date,
    rating, reviews_count, products_count, sales_count, views_count,
    created_at, updated_at
) VALUES (
    2, -- user_id (Dmitry Voroshilov)
    'belgrade-real-estate', 
    'Belgrade Real Estate Agency',
    'Ведущее агентство недвижимости в Белграде. Продажа и аренда квартир, домов и коммерческой недвижимости. Полное сопровождение сделок.',
    '/storefronts/real-estate/logo.jpg',
    '/storefronts/real-estate/banner.jpg',
    '+381 11 3456789',
    'info@belgrade-realestate.rs',
    'https://belgrade-realestate.rs',
    'Кнез Михаилова 42',
    'Београд',
    '11000',
    'RS',
    44.8176, -- Belgrade center coordinates
    20.4633,
    '{"layout": "grid", "primaryColor": "#2c5530", "secondaryColor": "#8b7355"}'::jsonb,
    '{"showMap": true, "enableChat": true, "workingHours": "09:00-18:00"}'::jsonb,
    '{"title": "Belgrade Real Estate - Недвижимость в Белграде", "description": "Лучшие предложения недвижимости в Белграде", "keywords": "недвижимость белград, квартиры белград, дома белград"}'::jsonb,
    true, -- is_active
    true, -- is_verified
    NOW(), -- verification_date
    4.8, -- rating
    45, -- reviews_count
    0, -- products_count (will be updated)
    156, -- sales_count
    8750, -- views_count
    NOW(),
    NOW()
);

-- Получаем ID созданной витрины
DO $$
DECLARE
    storefront_id INTEGER;
BEGIN
    SELECT id INTO storefront_id FROM storefronts WHERE slug = 'belgrade-real-estate';
    
    -- Добавляем товары-квартиры с атрибутами недвижимости
    
    -- Квартира 1: Студия в центре
    INSERT INTO storefront_products (
        storefront_id, name, description, price, currency,
        category_id, sku, barcode,
        stock_quantity, stock_status, is_active,
        attributes,
        has_individual_location, individual_address,
        individual_latitude, individual_longitude,
        location_privacy, show_on_map,
        created_at, updated_at
    ) VALUES (
        storefront_id,
        'Студия 28м² - Стари Град',
        'Современная студия в историческом центре Белграда. Полностью меблирована, новый ремонт. Идеально для инвестиций или проживания. 3 этаж, есть лифт.',
        65000,
        'EUR',
        1004, -- Real Estate category
        'BRE-APT-001',
        NULL,
        1,
        'available',
        true,
        '{
            "property_type": "apartment",
            "rooms": "studio",
            "area_m2": 28,
            "floor": 3,
            "total_floors": 5,
            "year_built": 2019,
            "heating": "central",
            "elevator": true,
            "parking": false,
            "balcony": false,
            "furnished": true,
            "energy_class": "B",
            "ownership": "full",
            "view": "courtyard"
        }'::jsonb,
        true,
        'Карађорђева 15, Стари Град',
        44.8176,
        20.4590,
        'exact',
        true,
        NOW(),
        NOW()
    );
    
    -- Квартира 2: Двухкомнатная на Новом Белграде
    INSERT INTO storefront_products (
        storefront_id, name, description, price, currency,
        category_id, sku, barcode,
        stock_quantity, stock_status, is_active,
        attributes,
        has_individual_location, individual_address,
        individual_latitude, individual_longitude,
        location_privacy, show_on_map,
        created_at, updated_at
    ) VALUES (
        storefront_id,
        'Двухкомнатная квартира 65м² - Нови Београд',
        'Просторная квартира в Новом Белграде, рядом с Дельта Сити. 2 спальни, большая гостиная, балкон с видом на реку. Гараж включен в цену.',
        125000,
        'EUR',
        1004,
        'BRE-APT-002',
        NULL,
        1,
        'available',
        true,
        '{
            "property_type": "apartment",
            "rooms": "2",
            "area_m2": 65,
            "floor": 8,
            "total_floors": 12,
            "year_built": 2021,
            "heating": "heat_pump",
            "elevator": true,
            "parking": true,
            "parking_type": "garage",
            "balcony": true,
            "balcony_area": 8,
            "furnished": false,
            "energy_class": "A",
            "ownership": "full",
            "view": "river"
        }'::jsonb,
        true,
        'Булевар Михајла Пупина 10, Нови Београд',
        44.8176,
        20.4190,
        'exact',
        true,
        NOW(),
        NOW()
    );
    
    -- Квартира 3: Пентхаус в Врачаре
    INSERT INTO storefront_products (
        storefront_id, name, description, price, currency,
        category_id, sku, barcode,
        stock_quantity, stock_status, is_active,
        attributes,
        has_individual_location, individual_address,
        individual_latitude, individual_longitude,
        location_privacy, show_on_map,
        created_at, updated_at
    ) VALUES (
        storefront_id,
        'Роскошный пентхаус 150м² - Врачар',
        'Эксклюзивный пентхаус в престижном районе Врачар. 3 спальни, 2 ванные, терраса 50м² с панорамным видом. Умный дом, дизайнерский ремонт.',
        380000,
        'EUR',
        1004,
        'BRE-APT-003',
        NULL,
        1,
        'available',
        true,
        '{
            "property_type": "penthouse",
            "rooms": "4",
            "area_m2": 150,
            "floor": 10,
            "total_floors": 10,
            "year_built": 2023,
            "heating": "floor_heating",
            "elevator": true,
            "parking": true,
            "parking_spots": 2,
            "parking_type": "underground",
            "balcony": false,
            "terrace": true,
            "terrace_area": 50,
            "furnished": true,
            "energy_class": "A+",
            "ownership": "full",
            "view": "panoramic",
            "smart_home": true,
            "security": true
        }'::jsonb,
        true,
        'Светог Саве 25, Врачар',
        44.7980,
        20.4680,
        'exact',
        true,
        NOW(),
        NOW()
    );
    
    -- Квартира 4: Трехкомнатная в Земуне
    INSERT INTO storefront_products (
        storefront_id, name, description, price, currency,
        category_id, sku, barcode,
        stock_quantity, stock_status, is_active,
        attributes,
        has_individual_location, individual_address,
        individual_latitude, individual_longitude,
        location_privacy, show_on_map,
        created_at, updated_at
    ) VALUES (
        storefront_id,
        'Семейная квартира 85м² - Земун',
        'Идеальная квартира для семьи в тихом районе Земуна. 3 спальни, 2 балкона, рядом школа и детский сад. Недавний ремонт.',
        165000,
        'EUR',
        1004,
        'BRE-APT-004',
        NULL,
        1,
        'available',
        true,
        '{
            "property_type": "apartment",
            "rooms": "3",
            "area_m2": 85,
            "floor": 2,
            "total_floors": 4,
            "year_built": 2015,
            "heating": "gas",
            "elevator": false,
            "parking": true,
            "parking_type": "street",
            "balcony": true,
            "balcony_count": 2,
            "furnished": false,
            "energy_class": "B",
            "ownership": "full",
            "view": "street",
            "near_school": true,
            "near_kindergarten": true
        }'::jsonb,
        true,
        'Главна 78, Земун',
        44.8433,
        20.4011,
        'exact',
        true,
        NOW(),
        NOW()
    );
    
    -- Квартира 5: Дом в пригороде
    INSERT INTO storefront_products (
        storefront_id, name, description, price, currency,
        category_id, sku, barcode,
        stock_quantity, stock_status, is_active,
        attributes,
        has_individual_location, individual_address,
        individual_latitude, individual_longitude,
        location_privacy, show_on_map,
        created_at, updated_at
    ) VALUES (
        storefront_id,
        'Современный дом 220м² - Дедиње',
        'Новый дом в элитном районе Дедиње. 4 спальни, 3 ванные, бассейн, сад 500м². Закрытая территория, видеонаблюдение.',
        650000,
        'EUR',
        1004,
        'BRE-HSE-001',
        NULL,
        1,
        'available',
        true,
        '{
            "property_type": "house",
            "rooms": "5+",
            "area_m2": 220,
            "land_area_m2": 500,
            "floors": 2,
            "year_built": 2024,
            "heating": "heat_pump",
            "garage": true,
            "garage_spots": 2,
            "pool": true,
            "garden": true,
            "furnished": false,
            "energy_class": "A++",
            "ownership": "full",
            "view": "garden",
            "security": true,
            "gated_community": true,
            "basement": true,
            "attic": false
        }'::jsonb,
        true,
        'Ужичка 35, Дедиње',
        44.7650,
        20.4400,
        'exact',
        true,
        NOW(),
        NOW()
    );
    
    -- Обновляем количество товаров в витрине
    UPDATE storefronts 
    SET products_count = 5 
    WHERE id = storefront_id;
    
END $$;

-- Добавляем переводы для витрины
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text)
SELECT 'storefront', id, 'name', 'en', 'Belgrade Real Estate Agency'
FROM storefronts WHERE slug = 'belgrade-real-estate'
UNION ALL
SELECT 'storefront', id, 'name', 'ru', 'Агентство недвижимости Белград'
FROM storefronts WHERE slug = 'belgrade-real-estate'
UNION ALL
SELECT 'storefront', id, 'description', 'en', 
'Leading real estate agency in Belgrade. Sale and rent of apartments, houses and commercial properties. Full transaction support.'
FROM storefronts WHERE slug = 'belgrade-real-estate'
UNION ALL
SELECT 'storefront', id, 'description', 'ru', 
'Ведущее агентство недвижимости в Белграде. Продажа и аренда квартир, домов и коммерческой недвижимости. Полное сопровождение сделок.'
FROM storefronts WHERE slug = 'belgrade-real-estate';

-- Добавляем переводы для товаров
INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text)
SELECT 'storefront_product', sp.id, 'name', 'en',
    CASE 
        WHEN sp.sku = 'BRE-APT-001' THEN 'Studio 28m² - Old Town'
        WHEN sp.sku = 'BRE-APT-002' THEN 'Two bedroom apartment 65m² - New Belgrade'
        WHEN sp.sku = 'BRE-APT-003' THEN 'Luxury penthouse 150m² - Vračar'
        WHEN sp.sku = 'BRE-APT-004' THEN 'Family apartment 85m² - Zemun'
        WHEN sp.sku = 'BRE-HSE-001' THEN 'Modern house 220m² - Dedinje'
    END
FROM storefront_products sp
JOIN storefronts s ON sp.storefront_id = s.id
WHERE s.slug = 'belgrade-real-estate';

INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text)
SELECT 'storefront_product', sp.id, 'name', 'ru',
    CASE 
        WHEN sp.sku = 'BRE-APT-001' THEN 'Студия 28м² - Стари Град'
        WHEN sp.sku = 'BRE-APT-002' THEN 'Двухкомнатная квартира 65м² - Новый Белград'
        WHEN sp.sku = 'BRE-APT-003' THEN 'Роскошный пентхаус 150м² - Врачар'
        WHEN sp.sku = 'BRE-APT-004' THEN 'Семейная квартира 85м² - Земун'
        WHEN sp.sku = 'BRE-HSE-001' THEN 'Современный дом 220м² - Дединье'
    END
FROM storefront_products sp
JOIN storefronts s ON sp.storefront_id = s.id
WHERE s.slug = 'belgrade-real-estate';

INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text)
SELECT 'storefront_product', sp.id, 'description', 'en',
    CASE 
        WHEN sp.sku = 'BRE-APT-001' THEN 'Modern studio in the historic center of Belgrade. Fully furnished, newly renovated. Perfect for investment or living. 3rd floor with elevator.'
        WHEN sp.sku = 'BRE-APT-002' THEN 'Spacious apartment in New Belgrade, near Delta City. 2 bedrooms, large living room, balcony with river view. Garage included.'
        WHEN sp.sku = 'BRE-APT-003' THEN 'Exclusive penthouse in prestigious Vračar district. 3 bedrooms, 2 bathrooms, 50m² terrace with panoramic view. Smart home, designer renovation.'
        WHEN sp.sku = 'BRE-APT-004' THEN 'Perfect family apartment in quiet Zemun area. 3 bedrooms, 2 balconies, near school and kindergarten. Recent renovation.'
        WHEN sp.sku = 'BRE-HSE-001' THEN 'New house in elite Dedinje area. 4 bedrooms, 3 bathrooms, pool, 500m² garden. Gated community, video surveillance.'
    END
FROM storefront_products sp
JOIN storefronts s ON sp.storefront_id = s.id
WHERE s.slug = 'belgrade-real-estate';

INSERT INTO translations (entity_type, entity_id, field_name, language, translated_text)
SELECT 'storefront_product', sp.id, 'description', 'ru',
    CASE 
        WHEN sp.sku = 'BRE-APT-001' THEN 'Современная студия в историческом центре Белграда. Полностью меблирована, новый ремонт. Идеально для инвестиций или проживания. 3 этаж, есть лифт.'
        WHEN sp.sku = 'BRE-APT-002' THEN 'Просторная квартира в Новом Белграде, рядом с Дельта Сити. 2 спальни, большая гостиная, балкон с видом на реку. Гараж включен в цену.'
        WHEN sp.sku = 'BRE-APT-003' THEN 'Эксклюзивный пентхаус в престижном районе Врачар. 3 спальни, 2 ванные, терраса 50м² с панорамным видом. Умный дом, дизайнерский ремонт.'
        WHEN sp.sku = 'BRE-APT-004' THEN 'Идеальная квартира для семьи в тихом районе Земуна. 3 спальни, 2 балкона, рядом школа и детский сад. Недавний ремонт.'
        WHEN sp.sku = 'BRE-HSE-001' THEN 'Новый дом в элитном районе Дединье. 4 спальни, 3 ванные, бассейн, сад 500м². Закрытая территория, видеонаблюдение.'
    END
FROM storefront_products sp
JOIN storefronts s ON sp.storefront_id = s.id
WHERE s.slug = 'belgrade-real-estate';