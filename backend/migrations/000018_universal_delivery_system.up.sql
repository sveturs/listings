-- ================================================================================
-- УНИВЕРСАЛЬНАЯ СИСТЕМА ДОСТАВКИ
-- Миграция 000018: Создание структуры для поддержки множественных провайдеров
-- ================================================================================

-- 1. ДОБАВЛЕНИЕ АТРИБУТОВ ДОСТАВКИ В СУЩЕСТВУЮЩИЕ ТАБЛИЦЫ
-- ================================================================================

-- Добавляем атрибуты доставки в metadata для marketplace_listings
UPDATE marketplace_listings
SET metadata = jsonb_set(
  COALESCE(metadata, '{}'),
  '{delivery_attributes}',
  '{}'::jsonb
)
WHERE metadata IS NULL OR metadata->>'delivery_attributes' IS NULL;

-- Добавляем атрибуты доставки в attributes для storefront_products
UPDATE storefront_products
SET attributes = jsonb_set(
  COALESCE(attributes, '{}'),
  '{delivery_attributes}',
  '{}'::jsonb
)
WHERE attributes IS NULL OR attributes->>'delivery_attributes' IS NULL;

-- Структура delivery_attributes:
-- {
--   "weight_kg": 0.5,
--   "dimensions": {
--     "length_cm": 30,
--     "width_cm": 20,
--     "height_cm": 10
--   },
--   "volume_m3": 0.006,
--   "is_fragile": false,
--   "requires_special_handling": false,
--   "stackable": true,
--   "max_stack_weight_kg": 50,
--   "packaging_type": "box",
--   "hazmat_class": null
-- }

-- 2. ТАБЛИЦА ПРОВАЙДЕРОВ ДОСТАВКИ
-- ================================================================================

CREATE TABLE IF NOT EXISTS delivery_providers (
    id SERIAL PRIMARY KEY,
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    logo_url VARCHAR(500),
    is_active BOOLEAN DEFAULT false,
    supports_cod BOOLEAN DEFAULT false,
    supports_insurance BOOLEAN DEFAULT false,
    supports_tracking BOOLEAN DEFAULT true,
    api_config JSONB,
    capabilities JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE delivery_providers IS 'Универсальная таблица провайдеров доставки';
COMMENT ON COLUMN delivery_providers.code IS 'Уникальный код провайдера (post_express, bex_express, etc)';
COMMENT ON COLUMN delivery_providers.api_config IS 'Конфигурация API провайдера';
COMMENT ON COLUMN delivery_providers.capabilities IS 'Возможности провайдера (max_weight, delivery_zones, etc)';

-- 3. УНИВЕРСАЛЬНЫЕ ОТПРАВЛЕНИЯ
-- ================================================================================

CREATE TABLE IF NOT EXISTS delivery_shipments (
    id SERIAL PRIMARY KEY,
    provider_id INTEGER REFERENCES delivery_providers(id),
    order_id INTEGER REFERENCES marketplace_orders(id),
    external_id VARCHAR(255),
    tracking_number VARCHAR(255) UNIQUE,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',

    -- Адреса
    sender_info JSONB NOT NULL,
    recipient_info JSONB NOT NULL,

    -- Параметры посылки
    package_info JSONB NOT NULL,

    -- Финансы
    delivery_cost DECIMAL(10,2),
    insurance_cost DECIMAL(10,2),
    cod_amount DECIMAL(10,2),

    -- Детализация стоимости
    cost_breakdown JSONB,

    -- Даты
    pickup_date DATE,
    estimated_delivery DATE,
    actual_delivery_date TIMESTAMPTZ,

    -- Метаданные
    provider_response JSONB,
    labels JSONB,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE delivery_shipments IS 'Универсальная таблица отправлений';
COMMENT ON COLUMN delivery_shipments.package_info IS 'Включает weight_kg, dimensions, is_fragile и другие атрибуты';
COMMENT ON COLUMN delivery_shipments.cost_breakdown IS 'base_price, weight_surcharge, fragile_surcharge, etc';

-- 4. СОБЫТИЯ ОТСЛЕЖИВАНИЯ
-- ================================================================================

CREATE TABLE IF NOT EXISTS delivery_tracking_events (
    id SERIAL PRIMARY KEY,
    shipment_id INTEGER REFERENCES delivery_shipments(id),
    provider_id INTEGER REFERENCES delivery_providers(id),
    event_time TIMESTAMPTZ NOT NULL,
    status VARCHAR(100) NOT NULL,
    location VARCHAR(500),
    description TEXT,
    raw_data JSONB,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE delivery_tracking_events IS 'События отслеживания для всех провайдеров';

-- 5. ДЕФОЛТНЫЕ АТРИБУТЫ ДОСТАВКИ ПО КАТЕГОРИЯМ
-- ================================================================================

CREATE TABLE IF NOT EXISTS delivery_category_defaults (
    id SERIAL PRIMARY KEY,
    category_id INTEGER REFERENCES marketplace_categories(id) UNIQUE,
    default_weight_kg DECIMAL(10,3),
    default_length_cm DECIMAL(10,2),
    default_width_cm DECIMAL(10,2),
    default_height_cm DECIMAL(10,2),
    default_packaging_type VARCHAR(50),
    is_typically_fragile BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE delivery_category_defaults IS 'Дефолтные атрибуты доставки для категорий товаров';

-- 6. ПРАВИЛА РАСЧЕТА СТОИМОСТИ ДОСТАВКИ
-- ================================================================================

CREATE TABLE IF NOT EXISTS delivery_pricing_rules (
    id SERIAL PRIMARY KEY,
    provider_id INTEGER REFERENCES delivery_providers(id),
    rule_type VARCHAR(50) NOT NULL,

    -- Весовые правила
    weight_ranges JSONB,

    -- Объемные правила
    volume_ranges JSONB,

    -- Зональные правила
    zone_multipliers JSONB,

    -- Дополнительные сборы
    fragile_surcharge DECIMAL(10,2) DEFAULT 0,
    oversized_surcharge DECIMAL(10,2) DEFAULT 0,
    special_handling_surcharge DECIMAL(10,2) DEFAULT 0,

    -- Минимальная и максимальная стоимость
    min_price DECIMAL(10,2),
    max_price DECIMAL(10,2),

    -- Формула расчета
    custom_formula TEXT,

    priority INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE delivery_pricing_rules IS 'Правила расчета стоимости доставки';
COMMENT ON COLUMN delivery_pricing_rules.rule_type IS 'weight_based, volume_based, zone_based, combined';
COMMENT ON COLUMN delivery_pricing_rules.weight_ranges IS '[{"from": 0, "to": 1, "price_per_kg": 5}, ...]';
COMMENT ON COLUMN delivery_pricing_rules.volume_ranges IS '[{"from": 0, "to": 0.01, "price_per_m3": 100}, ...]';

-- 7. ЗОНЫ ДОСТАВКИ
-- ================================================================================

-- Добавляем расширение PostGIS если не установлено
CREATE EXTENSION IF NOT EXISTS postgis;

CREATE TABLE IF NOT EXISTS delivery_zones (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL,

    -- Географические границы
    countries TEXT[],
    regions TEXT[],
    cities TEXT[],
    postal_codes TEXT[],

    -- Полигон для точного определения (GIS)
    boundary GEOMETRY(POLYGON, 4326),

    -- Расстояние от центра (для радиусных зон)
    center_point GEOMETRY(POINT, 4326),
    radius_km DECIMAL(10,2),

    created_at TIMESTAMPTZ DEFAULT NOW()
);

COMMENT ON TABLE delivery_zones IS 'Зоны доставки';
COMMENT ON COLUMN delivery_zones.type IS 'local, regional, national, international';

-- 8. СВЯЗЬ С ЗАКАЗАМИ
-- ================================================================================

-- Добавляем связь с универсальными отправлениями
ALTER TABLE marketplace_orders
ADD COLUMN IF NOT EXISTS delivery_shipment_id INTEGER REFERENCES delivery_shipments(id);

-- 9. ИНДЕКСЫ ДЛЯ ПРОИЗВОДИТЕЛЬНОСТИ
-- ================================================================================

CREATE INDEX IF NOT EXISTS idx_delivery_shipments_tracking ON delivery_shipments(tracking_number);
CREATE INDEX IF NOT EXISTS idx_delivery_shipments_status ON delivery_shipments(status);
CREATE INDEX IF NOT EXISTS idx_delivery_shipments_order ON delivery_shipments(order_id);
CREATE INDEX IF NOT EXISTS idx_delivery_tracking_events_shipment ON delivery_tracking_events(shipment_id);
CREATE INDEX IF NOT EXISTS idx_delivery_tracking_events_time ON delivery_tracking_events(event_time);
CREATE INDEX IF NOT EXISTS idx_delivery_category_defaults_category ON delivery_category_defaults(category_id);
CREATE INDEX IF NOT EXISTS idx_delivery_pricing_rules_provider ON delivery_pricing_rules(provider_id);
CREATE INDEX IF NOT EXISTS idx_delivery_pricing_rules_active ON delivery_pricing_rules(is_active) WHERE is_active = true;
CREATE INDEX IF NOT EXISTS idx_delivery_zones_boundary ON delivery_zones USING GIST(boundary);
CREATE INDEX IF NOT EXISTS idx_delivery_zones_center ON delivery_zones USING GIST(center_point);

-- 10. ЗАПОЛНЕНИЕ НАЧАЛЬНЫХ ДАННЫХ
-- ================================================================================

-- Добавляем провайдеров
INSERT INTO delivery_providers (code, name, logo_url, is_active, supports_cod, supports_insurance, supports_tracking, capabilities) VALUES
('post_express', 'Post Express', '/images/providers/post-express.png', true, true, true, true,
 '{"max_weight_kg": 30, "max_volume_m3": 0.5, "delivery_zones": ["serbia", "montenegro", "bosnia"], "delivery_types": ["standard", "express", "office_pickup"]}'),
('bex_express', 'BEX Express', '/images/providers/bex-express.png', false, true, true, true,
 '{"max_weight_kg": 50, "max_volume_m3": 1.0, "delivery_zones": ["serbia"], "delivery_types": ["standard", "express"]}'),
('aks_express', 'AKS Express', '/images/providers/aks-express.png', false, true, true, true,
 '{"max_weight_kg": 40, "max_volume_m3": 0.8, "delivery_zones": ["serbia", "region"], "delivery_types": ["standard"]}'),
('d_express', 'D Express', '/images/providers/d-express.png', false, true, false, true,
 '{"max_weight_kg": 25, "max_volume_m3": 0.4, "delivery_zones": ["serbia"], "delivery_types": ["standard", "same_day"]}'),
('city_express', 'City Express', '/images/providers/city-express.png', false, true, true, true,
 '{"max_weight_kg": 35, "max_volume_m3": 0.6, "delivery_zones": ["serbia", "montenegro"], "delivery_types": ["standard", "express"]}'),
('dhl_express', 'DHL Express', '/images/providers/dhl-express.png', false, false, true, true,
 '{"max_weight_kg": 70, "max_volume_m3": 2.0, "delivery_zones": ["worldwide"], "delivery_types": ["express", "international"]}')
ON CONFLICT (code) DO UPDATE SET
    name = EXCLUDED.name,
    capabilities = EXCLUDED.capabilities;

-- Добавляем базовые правила расчета для Post Express (активного провайдера)
INSERT INTO delivery_pricing_rules (provider_id, rule_type, weight_ranges, fragile_surcharge, min_price, max_price, is_active)
SELECT
    id,
    'weight_based',
    '[
        {"from": 0, "to": 1, "base_price": 300, "price_per_kg": 0},
        {"from": 1, "to": 5, "base_price": 300, "price_per_kg": 50},
        {"from": 5, "to": 10, "base_price": 500, "price_per_kg": 40},
        {"from": 10, "to": 20, "base_price": 700, "price_per_kg": 35},
        {"from": 20, "to": 30, "base_price": 1000, "price_per_kg": 30}
    ]'::jsonb,
    50, -- доплата за хрупкий товар
    250, -- минимальная стоимость
    5000, -- максимальная стоимость
    true
FROM delivery_providers WHERE code = 'post_express';

-- Добавляем зоны доставки для Сербии
INSERT INTO delivery_zones (name, type, countries, cities) VALUES
('Белград и окрестности', 'local', ARRAY['RS'], ARRAY['Белград', 'Земун', 'Нови Београд', 'Палилула', 'Звездара']),
('Центральная Сербия', 'regional', ARRAY['RS'], ARRAY['Крагујевац', 'Ниш', 'Нови Сад', 'Суботица']),
('Вся Сербия', 'national', ARRAY['RS'], NULL),
('Балканы', 'international', ARRAY['RS', 'ME', 'BA', 'HR', 'MK'], NULL);

-- 11. ФУНКЦИИ ДЛЯ РАБОТЫ С АТРИБУТАМИ
-- ================================================================================

-- Функция получения атрибутов доставки товара
CREATE OR REPLACE FUNCTION get_delivery_attributes(
    p_product_id INTEGER,
    p_product_type VARCHAR DEFAULT 'listing'
) RETURNS JSONB AS $$
DECLARE
    v_attributes JSONB;
    v_category_id INTEGER;
    v_defaults RECORD;
BEGIN
    IF p_product_type = 'listing' THEN
        SELECT
            ml.metadata->'delivery_attributes',
            ml.category_id
        INTO v_attributes, v_category_id
        FROM marketplace_listings ml
        WHERE ml.id = p_product_id;
    ELSE
        SELECT
            sp.attributes->'delivery_attributes',
            sp.category_id
        INTO v_attributes, v_category_id
        FROM storefront_products sp
        WHERE sp.id = p_product_id;
    END IF;

    -- Если атрибуты пустые, используем дефолтные из категории
    IF v_attributes IS NULL OR v_attributes = '{}'::jsonb THEN
        SELECT * INTO v_defaults
        FROM delivery_category_defaults
        WHERE category_id = v_category_id;

        IF FOUND THEN
            v_attributes = jsonb_build_object(
                'weight_kg', v_defaults.default_weight_kg,
                'dimensions', jsonb_build_object(
                    'length_cm', v_defaults.default_length_cm,
                    'width_cm', v_defaults.default_width_cm,
                    'height_cm', v_defaults.default_height_cm
                ),
                'packaging_type', v_defaults.default_packaging_type,
                'is_fragile', v_defaults.is_typically_fragile
            );
        END IF;
    END IF;

    RETURN COALESCE(v_attributes, '{}'::jsonb);
END;
$$ LANGUAGE plpgsql;

COMMENT ON FUNCTION get_delivery_attributes IS 'Получение атрибутов доставки товара с учетом дефолтных значений категории';

-- Функция расчета объемного веса
CREATE OR REPLACE FUNCTION calculate_volumetric_weight(
    p_length_cm NUMERIC,
    p_width_cm NUMERIC,
    p_height_cm NUMERIC,
    p_divisor INTEGER DEFAULT 5000
) RETURNS NUMERIC AS $$
BEGIN
    RETURN (p_length_cm * p_width_cm * p_height_cm) / p_divisor;
END;
$$ LANGUAGE plpgsql IMMUTABLE;

COMMENT ON FUNCTION calculate_volumetric_weight IS 'Расчет объемного веса по габаритам (стандартный делитель 5000)';

-- 12. ТРИГГЕРЫ ДЛЯ АВТОМАТИЧЕСКОГО ОБНОВЛЕНИЯ
-- ================================================================================

-- Триггер для обновления updated_at
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_delivery_providers_updated_at BEFORE UPDATE ON delivery_providers
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_delivery_category_defaults_updated_at BEFORE UPDATE ON delivery_category_defaults
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_delivery_pricing_rules_updated_at BEFORE UPDATE ON delivery_pricing_rules
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_delivery_shipments_updated_at BEFORE UPDATE ON delivery_shipments
FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ================================================================================
-- КОНЕЦ МИГРАЦИИ
-- ================================================================================