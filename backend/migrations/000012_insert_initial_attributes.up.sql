-- backend/migrations/000012_insert_initial_attributes.up.sql
-- Автомобили
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required) VALUES
('make', 'Марка', 'select', '{"values": ["Audi", "BMW", "Mercedes", "Toyota", "Honda", "Ford", "Chevrolet", "Volkswagen", "Nissan", "Other"]}', true, true, true),
('model', 'Модель', 'text', NULL, true, true, true),
('year', 'Год выпуска', 'number', '{"min": 1900, "max": 2025}', true, true, true),
('mileage', 'Пробег (км)', 'number', '{"min": 0}', true, true, false),
('engine_capacity', 'Объем двигателя (л)', 'number', '{"min": 0.1, "max": 10, "step": 0.1}', true, true, false),
('fuel_type', 'Тип топлива', 'select', '{"values": ["Бензин", "Дизель", "Гибрид", "Электро", "Газ"]}', true, true, false),
('transmission', 'Коробка передач', 'select', '{"values": ["Механика", "Автомат", "Робот", "Вариатор"]}', true, true, false),
('body_type', 'Тип кузова', 'select', '{"values": ["Седан", "Хэтчбек", "Универсал", "Внедорожник", "Купе", "Кабриолет", "Минивэн", "Пикап"]}', true, true, false),
('color', 'Цвет', 'select', '{"values": ["Белый", "Черный", "Серый", "Серебристый", "Красный", "Синий", "Зеленый", "Желтый", "Коричневый", "Другой"]}', true, true, false);

-- Недвижимость
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required) VALUES
('property_type', 'Тип недвижимости', 'select', '{"values": ["Квартира", "Дом", "Комната", "Земельный участок", "Коммерческая недвижимость", "Гараж"]}', true, true, true),
('rooms', 'Количество комнат', 'select', '{"values": ["Студия", "1", "2", "3", "4", "5+"]}', true, true, false),
('floor', 'Этаж', 'number', '{"min": 0, "max": 100}', true, true, false),
('total_floors', 'Этажей в доме', 'number', '{"min": 1, "max": 100}', true, true, false),
('area', 'Площадь (м²)', 'number', '{"min": 1}', true, true, true),
('land_area', 'Площадь участка (сот.)', 'number', '{"min": 0}', true, true, false),
('building_type', 'Тип дома', 'select', '{"values": ["Панельный", "Кирпичный", "Монолитный", "Деревянный", "Блочный", "Другой"]}', true, true, false),
('has_balcony', 'Балкон', 'boolean', NULL, true, true, false),
('has_elevator', 'Лифт', 'boolean', NULL, true, true, false),
('has_parking', 'Парковка', 'boolean', NULL, true, true, false);

-- Электроника - Телефоны
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required) VALUES
('brand', 'Бренд', 'select', '{"values": ["Apple", "Samsung", "Xiaomi", "Huawei", "Google", "OnePlus", "Sony", "Nokia", "LG", "Other"]}', true, true, true),
('model_phone', 'Модель', 'text', NULL, true, true, true),
('memory', 'Память (ГБ)', 'select', '{"values": ["8", "16", "32", "64", "128", "256", "512", "1024"]}', true, true, false),
('ram', 'ОЗУ (ГБ)', 'select', '{"values": ["1", "2", "3", "4", "6", "8", "12", "16"]}', true, true, false),
('os', 'Операционная система', 'select', '{"values": ["iOS", "Android", "Windows", "Другая"]}', true, true, false),
('screen_size', 'Размер экрана (дюймы)', 'number', '{"min": 1, "max": 15, "step": 0.1}', true, true, false),
('camera', 'Камера (МП)', 'number', '{"min": 1}', true, true, false),
('has_5g', '5G', 'boolean', NULL, true, true, false);

-- Компьютеры
INSERT INTO category_attributes (name, display_name, attribute_type, options, is_searchable, is_filterable, is_required) VALUES
('pc_brand', 'Бренд', 'select', '{"values": ["Apple", "Dell", "HP", "Lenovo", "Asus", "Acer", "MSI", "Gigabyte", "Сборка"]}', true, true, true),
('pc_type', 'Тип', 'select', '{"values": ["Ноутбук", "Настольный ПК", "Моноблок", "Сервер", "Другое"]}', true, true, true),
('cpu', 'Процессор', 'text', NULL, true, true, false),
('gpu', 'Видеокарта', 'text', NULL, true, true, false),
('ram_pc', 'ОЗУ (ГБ)', 'select', '{"values": ["2", "4", "8", "16", "32", "64", "128"]}', true, true, false),
('storage_type', 'Тип накопителя', 'select', '{"values": ["HDD", "SSD", "HDD+SSD"]}', true, true, false),
('storage_capacity', 'Объем накопителя (ГБ)', 'number', '{"min": 1}', true, true, false),
('os_pc', 'Операционная система', 'select', '{"values": ["Windows", "macOS", "Linux", "Без ОС"]}', true, true, false);