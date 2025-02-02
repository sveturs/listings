-- Сначала проверяем существование таблицы
DO $$ 
BEGIN
    IF NOT EXISTS (
        SELECT FROM pg_tables 
        WHERE schemaname = 'public' 
        AND tablename = 'translations'
    ) THEN
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

-- Очищаем старые переводы категорий
DELETE FROM translations WHERE entity_type = 'category';

-- Создаем временную таблицу для хранения всех переводов
CREATE TEMP TABLE category_translations AS
SELECT 
    id,
    name,
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
        WHEN name = 'Пољопривредне машине' THEN 'Agricultural Machinery'
        WHEN name = 'Домаће животиње' THEN 'Farm Animals'
        WHEN name = 'Пољопривредни производи' THEN 'Agricultural Products'
        WHEN name = 'Трактори' THEN 'Tractors'
        WHEN name = 'Комбајни' THEN 'Harvesters'
        WHEN name = 'Плугови и дрљаче' THEN 'Plows and Harrows'
        WHEN name = 'Сејалице' THEN 'Seeders'
        WHEN name = 'Опрема за наводњавање' THEN 'Irrigation Equipment'
        WHEN name = 'Краве' THEN 'Cows'
        WHEN name = 'Свиње' THEN 'Pigs'
        WHEN name = 'Козе и овце' THEN 'Goats and Sheep'
        WHEN name = 'Живина' THEN 'Poultry'
        WHEN name = 'Сточна храна' THEN 'Animal Feed'
        WHEN name = 'Поврће' THEN 'Vegetables'
        WHEN name = 'Воће' THEN 'Fruits'
        WHEN name = 'Житарице' THEN 'Grains'
        WHEN name = 'Млечни производи' THEN 'Dairy Products'
        WHEN name = 'Месо и месни производи' THEN 'Meat Products'
        WHEN name = 'Мед и пчеларски производи' THEN 'Honey and Beekeeping Products'
        ELSE name
    END as en_name,
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
        WHEN name = 'Пољопривредне машине' THEN 'Сельскохозяйственная техника'
        WHEN name = 'Домаће животиње' THEN 'Домашние животные'
        WHEN name = 'Пољопривредни производи' THEN 'Сельхозпродукция'
        WHEN name = 'Трактори' THEN 'Тракторы'
        WHEN name = 'Комбајни' THEN 'Комбайны'
        WHEN name = 'Плугови и дрљаче' THEN 'Плуги и бороны'
        WHEN name = 'Сејалице' THEN 'Сеялки'
        WHEN name = 'Опрема за наводњавање' THEN 'Оборудование для полива'
        WHEN name = 'Краве' THEN 'Коровы'
        WHEN name = 'Свиње' THEN 'Свиньи'
        WHEN name = 'Козе и овце' THEN 'Козы и овцы'
        WHEN name = 'Живина' THEN 'Птица'
        WHEN name = 'Сточна храна' THEN 'Корм для животных'
        WHEN name = 'Поврће' THEN 'Овощи'
        WHEN name = 'Воће' THEN 'Фрукты'
        WHEN name = 'Житарице' THEN 'Зерновые'
        WHEN name = 'Млечни производи' THEN 'Молочные продукты'
        WHEN name = 'Месо и месни производи' THEN 'Мясо и мясные продукты'
        WHEN name = 'Мед и пчеларски производи' THEN 'Мёд и продукты пчеловодства'
        ELSE name
    END as ru_name
FROM marketplace_categories;

-- Добавляем английские переводы
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 'category', id, 'en', 'name', en_name, true, true
FROM category_translations
WHERE en_name != name;

-- Добавляем русские переводы
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 'category', id, 'ru', 'name', ru_name, true, true
FROM category_translations
WHERE ru_name != name;

-- Удаляем временную таблицу
DROP TABLE category_translations;

-- Создаем индекс для ускорения поиска переводов
CREATE INDEX IF NOT EXISTS idx_translations_lookup ON translations(entity_type, entity_id, language);