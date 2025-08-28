-- Создаем функцию для синхронизации товаров витрин с marketplace_listings
CREATE OR REPLACE FUNCTION sync_storefront_to_marketplace()
RETURNS TRIGGER AS $$
DECLARE
    storefront_user_id INTEGER;
    storefront_location_data JSONB;
    formatted_address TEXT;
    lat FLOAT;
    lon FLOAT;
BEGIN
    -- Получаем данные о владельце витрины и местоположении
    SELECT 
        us.user_id,
        jsonb_build_object(
            'country', l.country,
            'city', l.city,
            'address', l.address,
            'latitude', l.latitude,
            'longitude', l.longitude
        )
    INTO storefront_user_id, storefront_location_data
    FROM user_storefronts us
    LEFT JOIN locations l ON us.location_id = l.id
    WHERE us.id = NEW.storefront_id;

    -- Определяем адрес и координаты
    IF NEW.has_individual_location AND NEW.individual_address IS NOT NULL THEN
        -- Используем индивидуальный адрес товара
        formatted_address := NEW.individual_address;
        lat := NEW.individual_latitude;
        lon := NEW.individual_longitude;
    ELSE
        -- Используем адрес витрины
        formatted_address := COALESCE(
            storefront_location_data->>'address',
            storefront_location_data->>'city',
            ''
        );
        lat := (storefront_location_data->>'latitude')::FLOAT;
        lon := (storefront_location_data->>'longitude')::FLOAT;
    END IF;

    IF TG_OP = 'INSERT' THEN
        -- Создаем запись в marketplace_listings
        INSERT INTO marketplace_listings (
            id,
            user_id,
            category_id,
            title,
            description,
            price,
            location,
            latitude,
            longitude,
            status,
            condition,
            storefront_id,
            attributes,
            show_on_map,
            created_at,
            updated_at
        ) VALUES (
            NEW.id,
            storefront_user_id,
            NEW.category_id,
            NEW.name,
            NEW.description,
            NEW.price,
            formatted_address,
            lat,
            lon,
            CASE WHEN NEW.is_active THEN 'active' ELSE 'disabled' END,
            'new',  -- Товары витрин всегда новые
            NEW.storefront_id,
            NEW.attributes,
            COALESCE(NEW.show_on_map, true),
            NEW.created_at,
            NEW.updated_at
        );
        
    ELSIF TG_OP = 'UPDATE' THEN
        -- Обновляем запись в marketplace_listings
        UPDATE marketplace_listings
        SET
            title = NEW.name,
            description = NEW.description,
            price = NEW.price,
            category_id = NEW.category_id,
            location = formatted_address,
            latitude = lat,
            longitude = lon,
            status = CASE WHEN NEW.is_active THEN 'active' ELSE 'disabled' END,
            attributes = NEW.attributes,
            show_on_map = COALESCE(NEW.show_on_map, true),
            updated_at = NEW.updated_at
        WHERE id = NEW.id AND storefront_id = NEW.storefront_id;
        
    ELSIF TG_OP = 'DELETE' THEN
        -- При удалении товара витрины удаляем запись из marketplace_listings
        DELETE FROM marketplace_listings
        WHERE id = OLD.id AND storefront_id = OLD.storefront_id;
        
        RETURN OLD;
    END IF;
    
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Создаем триггер
CREATE TRIGGER sync_storefront_products_to_marketplace
    AFTER INSERT OR UPDATE OR DELETE ON storefront_products
    FOR EACH ROW
    EXECUTE FUNCTION sync_storefront_to_marketplace();

-- Синхронизируем существующие товары витрин
INSERT INTO marketplace_listings (
    id,
    user_id,
    category_id,
    title,
    description,
    price,
    location,
    latitude,
    longitude,
    status,
    condition,
    storefront_id,
    attributes,
    show_on_map,
    created_at,
    updated_at
)
SELECT 
    sp.id,
    us.user_id,
    sp.category_id,
    sp.name,
    sp.description,
    sp.price,
    CASE 
        WHEN sp.has_individual_location AND sp.individual_address IS NOT NULL 
        THEN sp.individual_address
        ELSE COALESCE(l.address, l.city, '')
    END as location,
    CASE 
        WHEN sp.has_individual_location AND sp.individual_latitude IS NOT NULL 
        THEN sp.individual_latitude
        ELSE l.latitude
    END as latitude,
    CASE 
        WHEN sp.has_individual_location AND sp.individual_longitude IS NOT NULL 
        THEN sp.individual_longitude
        ELSE l.longitude
    END as longitude,
    CASE WHEN sp.is_active THEN 'active' ELSE 'disabled' END as status,
    'new' as condition,
    sp.storefront_id,
    sp.attributes,
    COALESCE(sp.show_on_map, true),
    sp.created_at,
    sp.updated_at
FROM storefront_products sp
JOIN user_storefronts us ON sp.storefront_id = us.id
LEFT JOIN locations l ON us.location_id = l.id
WHERE NOT EXISTS (
    SELECT 1 FROM marketplace_listings ml 
    WHERE ml.id = sp.id AND ml.storefront_id = sp.storefront_id
);

-- Логирование
DO $$
DECLARE
    synced_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO synced_count
    FROM storefront_products sp
    WHERE EXISTS (
        SELECT 1 FROM marketplace_listings ml 
        WHERE ml.id = sp.id AND ml.storefront_id = sp.storefront_id
    );
    
    RAISE NOTICE 'Created trigger for automatic sync between storefront_products and marketplace_listings. Synced % existing products.', synced_count;
END $$;