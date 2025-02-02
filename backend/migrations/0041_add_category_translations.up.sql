-- Сначала проверяем существование таблицы
DO $$ 
BEGIN
    -- Проверяем существование таблицы translations
    IF NOT EXISTS (
        SELECT FROM pg_tables 
        WHERE schemaname = 'public' 
        AND tablename = 'translations'
    ) THEN
        -- Если таблица не существует, создаем её
        CREATE TABLE translations (
            id SERIAL PRIMARY KEY,
            entity_type VARCHAR(50) NOT NULL,
            entity_id INTEGER NOT NULL,
            language VARCHAR(10) NOT NULL,
            field_name VARCHAR(50) NOT NULL,
            translated_text TEXT NOT NULL,
            is_machine_translated BOOLEAN DEFAULT true,
            is_verified BOOLEAN DEFAULT false,
            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            UNIQUE(entity_type, entity_id, language, field_name)
        );
    END IF;
END $$;

-- Очищаем старые переводы категорий если они есть
DELETE FROM translations WHERE entity_type = 'category';

-- Добавляем переводы для основных категорий
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'category',
    id,
    'en',
    'name',
    CASE 
        WHEN name = 'Превоз' THEN 'Transport'
        WHEN name = 'Некретнине' THEN 'Real Estate'
        WHEN name = 'Електроника' THEN 'Electronics'
        WHEN name = 'Одећа и обућа' THEN 'Clothing and Shoes'
        WHEN name = 'Кућа и башта' THEN 'Home and Garden'
        WHEN name = 'Пољопривреда' THEN 'Agriculture'
        WHEN name = 'Послови' THEN 'Jobs'
        WHEN name = 'Лични предмети' THEN 'Personal Items'
        WHEN name = 'Хоби и разонода' THEN 'Hobbies and Leisure'
        WHEN name = 'Кућни љубимци' THEN 'Pets'
        WHEN name = 'Услуге' THEN 'Services'
        WHEN name = 'Бизнис и индустрија' THEN 'Business and Industry'
        -- Добавляем переводы для всех подкатегорий
        WHEN name = 'Аутомобили' THEN 'Cars'
        WHEN name = 'Мотоцикли' THEN 'Motorcycles'
        WHEN name = 'Електрична возила' THEN 'Electric Vehicles'
        WHEN name = 'Теретна возила' THEN 'Trucks'
        WHEN name = 'Делови и опрема' THEN 'Parts and Accessories'
        WHEN name = 'Електрични аутомобили' THEN 'Electric Cars'
        WHEN name = 'Електрични тротинети' THEN 'Electric Scooters'
        WHEN name = 'Електрични бицикли' THEN 'Electric Bikes'
        WHEN name = 'Издавање' THEN 'Rent'
        WHEN name = 'Продаја' THEN 'Sale'
        WHEN name = 'Гараже и паркинг' THEN 'Garages and Parking'
        -- Добавляем остальные подкатегории...
        ELSE name -- Для остальных оставляем оригинальное имя
    END,
    true,
    true
FROM marketplace_categories;

-- Аналогично добавляем русские переводы
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'category',
    id,
    'ru',
    'name',
    CASE 
        WHEN name = 'Превоз' THEN 'Транспорт'
        WHEN name = 'Некретнине' THEN 'Недвижимость'
        WHEN name = 'Електроника' THEN 'Электроника'
        WHEN name = 'Одећа и обућа' THEN 'Одежда и обувь'
        WHEN name = 'Кућа и башта' THEN 'Дом и сад'
        WHEN name = 'Пољопривреда' THEN 'Сельское хозяйство'
        WHEN name = 'Послови' THEN 'Работа'
        WHEN name = 'Лични предмети' THEN 'Личные вещи'
        WHEN name = 'Хоби и разонода' THEN 'Хобби и развлечения'
        WHEN name = 'Кућни љубимци' THEN 'Домашние животные'
        WHEN name = 'Услуге' THEN 'Услуги'
        WHEN name = 'Бизнис и индустрија' THEN 'Бизнес и промышленность'
        -- Добавляем переводы для всех подкатегорий
        WHEN name = 'Аутомобили' THEN 'Автомобили'
        WHEN name = 'Мотоцикли' THEN 'Мотоциклы'
        WHEN name = 'Електрична возила' THEN 'Электротранспорт'
        WHEN name = 'Теретна возила' THEN 'Грузовой транспорт'
        WHEN name = 'Делови и опрема' THEN 'Запчасти и оборудование'
        WHEN name = 'Електрични аутомобили' THEN 'Электромобили'
        WHEN name = 'Електрични тротинети' THEN 'Электросамокаты'
        WHEN name = 'Електрични бицикли' THEN 'Электровелосипеды'
        WHEN name = 'Издавање' THEN 'Аренда'
        WHEN name = 'Продаја' THEN 'Продажа'
        WHEN name = 'Гараже и паркинг' THEN 'Гаражи и парковки'
        -- Добавляем остальные подкатегории...
        ELSE name -- Для остальных оставляем оригинальное имя
    END,
    true,
    true
FROM marketplace_categories;

-- Создаем индекс если его еще нет
CREATE INDEX IF NOT EXISTS idx_translations_lookup ON translations(entity_type, entity_id, language);