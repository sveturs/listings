-- Создание постоянной таблицы для хранения значений атрибутов и их переводов
CREATE TABLE IF NOT EXISTS attribute_option_translations (
    id SERIAL PRIMARY KEY,
    attribute_name VARCHAR(100) NOT NULL,
    option_value TEXT NOT NULL,
    en_translation TEXT NOT NULL,
    sr_translation TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(attribute_name, option_value)
);

-- Заполнение таблицы значениями атрибутов и их переводами
-- Автомобили: Марки
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('make', 'Audi', 'Audi', 'Audi'),
('make', 'BMW', 'BMW', 'BMW'),
('make', 'Mercedes', 'Mercedes', 'Mercedes'),
('make', 'Toyota', 'Toyota', 'Toyota'),
('make', 'Honda', 'Honda', 'Honda'),
('make', 'Ford', 'Ford', 'Ford'),
('make', 'Chevrolet', 'Chevrolet', 'Chevrolet'),
('make', 'Volkswagen', 'Volkswagen', 'Volkswagen'),
('make', 'Nissan', 'Nissan', 'Nissan'),
('make', 'Other', 'Other', 'Drugo')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Автомобили: Тип топлива
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('fuel_type', 'Бензин', 'Petrol', 'Benzin'),
('fuel_type', 'Дизель', 'Diesel', 'Dizel'),
('fuel_type', 'Гибрид', 'Hybrid', 'Hibrid'),
('fuel_type', 'Электро', 'Electric', 'Električni'),
('fuel_type', 'Газ', 'Gas', 'Gas')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Автомобили: Коробка передач
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('transmission', 'Механика', 'Manual', 'Manuelni'),
('transmission', 'Автомат', 'Automatic', 'Automatski'),
('transmission', 'Робот', 'Automated Manual', 'Robotizovani'),
('transmission', 'Вариатор', 'CVT', 'Varijator')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Автомобили: Тип кузова
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('body_type', 'Седан', 'Sedan', 'Sedan'),
('body_type', 'Хэтчбек', 'Hatchback', 'Hečbek'),
('body_type', 'Универсал', 'Station Wagon', 'Karavan'),
('body_type', 'Внедорожник', 'SUV', 'Džip/SUV'),
('body_type', 'Купе', 'Coupe', 'Kupe'),
('body_type', 'Кабриолет', 'Convertible', 'Kabriolet'),
('body_type', 'Минивэн', 'Minivan', 'Miniven'),
('body_type', 'Пикап', 'Pickup', 'Pikap')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Автомобили: Цвет
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('color', 'Белый', 'White', 'Bela'),
('color', 'Черный', 'Black', 'Crna'),
('color', 'Серый', 'Gray', 'Siva'),
('color', 'Серебристый', 'Silver', 'Srebrna'),
('color', 'Красный', 'Red', 'Crvena'),
('color', 'Синий', 'Blue', 'Plava'),
('color', 'Зеленый', 'Green', 'Zelena'),
('color', 'Желтый', 'Yellow', 'Žuta'),
('color', 'Коричневый', 'Brown', 'Braon'),
('color', 'Другой', 'Other', 'Druga')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Автомобили: Привод
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('drive_type', 'Передний', 'Front-wheel drive', 'Prednji pogon'),
('drive_type', 'Задний', 'Rear-wheel drive', 'Zadnji pogon'),
('drive_type', 'Полный', 'All-wheel drive', 'Pogon na sve točkove'),
('drive_type', 'Другой', 'Other', 'Drugi')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Автомобили: Количество дверей
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('number_of_doors', '2', '2', '2'),
('number_of_doors', '3', '3', '3'),
('number_of_doors', '4', '4', '4'),
('number_of_doors', '5', '5', '5'),
('number_of_doors', '6+', '6+', '6+')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Автомобили: Количество мест
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('number_of_seats', '1', '1', '1'),
('number_of_seats', '2', '2', '2'),
('number_of_seats', '3', '3', '3'),
('number_of_seats', '4', '4', '4'),
('number_of_seats', '5', '5', '5'),
('number_of_seats', '6', '6', '6'),
('number_of_seats', '7', '7', '7'),
('number_of_seats', '8+', '8+', '8+')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Недвижимость: Тип недвижимости
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('property_type', 'Квартира', 'Apartment', 'Stan'),
('property_type', 'Дом', 'House', 'Kuća'),
('property_type', 'Комната', 'Room', 'Soba'),
('property_type', 'Земельный участок', 'Land', 'Zemljište'),
('property_type', 'Коммерческая недвижимость', 'Commercial property', 'Poslovna nekretnina'),
('property_type', 'Гараж', 'Garage', 'Garaža')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Недвижимость: Количество комнат
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('rooms', 'Студия', 'Studio', 'Garsonjera'),
('rooms', '1', '1', '1'),
('rooms', '2', '2', '2'),
('rooms', '3', '3', '3'),
('rooms', '4', '4', '4'),
('rooms', '5+', '5+', '5+')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Недвижимость: Тип дома
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('building_type', 'Панельный', 'Panel', 'Montažni'),
('building_type', 'Кирпичный', 'Brick', 'Ciglani'),
('building_type', 'Монолитный', 'Monolithic', 'Monolitni'),
('building_type', 'Деревянный', 'Wooden', 'Drveni'),
('building_type', 'Блочный', 'Block', 'Blok'),
('building_type', 'Другой', 'Other', 'Drugi')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Телефоны: Бренд
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('brand', 'Apple', 'Apple', 'Apple'),
('brand', 'Samsung', 'Samsung', 'Samsung'),
('brand', 'Xiaomi', 'Xiaomi', 'Xiaomi'),
('brand', 'Huawei', 'Huawei', 'Huawei'),
('brand', 'Google', 'Google', 'Google'),
('brand', 'OnePlus', 'OnePlus', 'OnePlus'),
('brand', 'Sony', 'Sony', 'Sony'),
('brand', 'Nokia', 'Nokia', 'Nokia'),
('brand', 'LG', 'LG', 'LG'),
('brand', 'Other', 'Other', 'Drugi')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Телефоны: Память
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('memory', '8', '8', '8'),
('memory', '16', '16', '16'),
('memory', '32', '32', '32'),
('memory', '64', '64', '64'),
('memory', '128', '128', '128'),
('memory', '256', '256', '256'),
('memory', '512', '512', '512'),
('memory', '1024', '1024', '1024')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Телефоны: ОЗУ
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('ram', '1', '1', '1'),
('ram', '2', '2', '2'),
('ram', '3', '3', '3'),
('ram', '4', '4', '4'),
('ram', '6', '6', '6'),
('ram', '8', '8', '8'),
('ram', '12', '12', '12'),
('ram', '16', '16', '16')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Телефоны: ОС
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('os', 'iOS', 'iOS', 'iOS'),
('os', 'Android', 'Android', 'Android'),
('os', 'Windows', 'Windows', 'Windows'),
('os', 'Другая', 'Other', 'Drugi')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Компьютеры: Бренд
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('pc_brand', 'Apple', 'Apple', 'Apple'),
('pc_brand', 'Dell', 'Dell', 'Dell'),
('pc_brand', 'HP', 'HP', 'HP'),
('pc_brand', 'Lenovo', 'Lenovo', 'Lenovo'),
('pc_brand', 'Asus', 'Asus', 'Asus'),
('pc_brand', 'Acer', 'Acer', 'Acer'),
('pc_brand', 'MSI', 'MSI', 'MSI'),
('pc_brand', 'Gigabyte', 'Gigabyte', 'Gigabyte'),
('pc_brand', 'Сборка', 'Custom build', 'Sastavljeno')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Компьютеры: Тип
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('pc_type', 'Ноутбук', 'Laptop', 'Laptop'),
('pc_type', 'Настольный ПК', 'Desktop PC', 'Desktop računar'),
('pc_type', 'Моноблок', 'All-in-One', 'All-in-One'),
('pc_type', 'Сервер', 'Server', 'Server'),
('pc_type', 'Другое', 'Other', 'Drugo')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Компьютеры: ОЗУ
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('ram_pc', '2', '2', '2'),
('ram_pc', '4', '4', '4'),
('ram_pc', '8', '8', '8'),
('ram_pc', '16', '16', '16'),
('ram_pc', '32', '32', '32'),
('ram_pc', '64', '64', '64'),
('ram_pc', '128', '128', '128')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Компьютеры: Тип накопителя
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('storage_type', 'HDD', 'HDD', 'HDD'),
('storage_type', 'SSD', 'SSD', 'SSD'),
('storage_type', 'HDD+SSD', 'HDD+SSD', 'HDD+SSD')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Компьютеры: ОС
INSERT INTO attribute_option_translations (attribute_name, option_value, en_translation, sr_translation) VALUES
('os_pc', 'Windows', 'Windows', 'Windows'),
('os_pc', 'macOS', 'macOS', 'macOS'),
('os_pc', 'Linux', 'Linux', 'Linux'),
('os_pc', 'Без ОС', 'No OS', 'Bez OS-a')
ON CONFLICT (attribute_name, option_value) DO UPDATE
SET en_translation = EXCLUDED.en_translation,
    sr_translation = EXCLUDED.sr_translation,
    updated_at = CURRENT_TIMESTAMP;

-- Функция для добавления переводов атрибутов из постоянной таблицы
CREATE OR REPLACE FUNCTION sync_attribute_option_translations() RETURNS void AS $$
DECLARE
    attr RECORD;
    opt RECORD;
    opt_index INTEGER;
    options_json JSONB;
    attr_id INTEGER;
BEGIN
    FOR attr IN SELECT DISTINCT attribute_name FROM attribute_option_translations LOOP
        -- Получение ID атрибута
        SELECT id, options INTO attr_id, options_json FROM category_attributes WHERE name = attr.attribute_name;
        
        IF attr_id IS NOT NULL AND options_json IS NOT NULL THEN
            -- Получение массива значений из JSON
            IF options_json ? 'values' THEN
                FOR opt_index IN 0..jsonb_array_length(options_json->'values')-1 LOOP
                    -- Получение оригинального значения
                    DECLARE
                        option_value TEXT := options_json->'values'->opt_index;
                        option_text TEXT;
                    BEGIN
                        -- Удаление кавычек из JSON-значения
                        option_text := trim(both '"' from option_value::text);
                        
                        -- Добавление перевода на английский
                        INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at)
                        SELECT 'attribute_option', attr_id, 'en', 'option_' || option_text, aot.en_translation, false, true, NOW(), NOW()
                        FROM attribute_option_translations aot
                        WHERE aot.attribute_name = attr.attribute_name AND aot.option_value = option_text
                        ON CONFLICT DO NOTHING;
                        
                        -- Добавление перевода на сербский
                        INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at)
                        SELECT 'attribute_option', attr_id, 'sr', 'option_' || option_text, aot.sr_translation, false, true, NOW(), NOW()
                        FROM attribute_option_translations aot
                        WHERE aot.attribute_name = attr.attribute_name AND aot.option_value = option_text
                        ON CONFLICT DO NOTHING;
                    END;
                END LOOP;
            END IF;
        END IF;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- Вызов функции для добавления переводов
SELECT sync_attribute_option_translations();

-- Создадим индекс для ускорения поиска по атрибутам и их значениям
CREATE INDEX IF NOT EXISTS idx_attribute_option_translations ON attribute_option_translations(attribute_name, option_value);

-- Комментарии к таблице и полям
COMMENT ON TABLE attribute_option_translations IS 'Содержит переводы значений атрибутов на различные языки';
COMMENT ON COLUMN attribute_option_translations.attribute_name IS 'Название атрибута в системе';
COMMENT ON COLUMN attribute_option_translations.option_value IS 'Значение опции атрибута на русском языке';
COMMENT ON COLUMN attribute_option_translations.en_translation IS 'Перевод на английский язык';
COMMENT ON COLUMN attribute_option_translations.sr_translation IS 'Перевод на сербский язык';