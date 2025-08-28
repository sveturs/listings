-- Добавление переводов для опций атрибутов
-- Примечание: en_translation не используется в текущей структуре таблицы

-- Переводы для condition
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('condition', 'new', 'Новый', 'Novo'),
('condition', 'used', 'Б/у', 'Polovno'),
('condition', 'refurbished', 'Восстановленный', 'Obnovljeno'),
('condition', 'damaged', 'Поврежденный', 'Oštećeno')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для color
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('color', 'black', 'Черный', 'Crna'),
('color', 'white', 'Белый', 'Bela'),
('color', 'silver', 'Серебристый', 'Srebrna'),
('color', 'gold', 'Золотой', 'Zlatna'),
('color', 'blue', 'Синий', 'Plava'),
('color', 'red', 'Красный', 'Crvena'),
('color', 'green', 'Зеленый', 'Zelena'),
('color', 'yellow', 'Желтый', 'Žuta'),
('color', 'purple', 'Фиолетовый', 'Ljubičasta'),
('color', 'grey', 'Серый', 'Siva'),
('color', 'other', 'Другой', 'Druga')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для transmission
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('transmission', 'manual', 'Механическая', 'Manuelni'),
('transmission', 'automatic', 'Автоматическая', 'Automatik'),
('transmission', 'semi-automatic', 'Полуавтоматическая', 'Poluautomatik'),
('transmission', 'cvt', 'Вариатор', 'CVT')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для fuel_type
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('fuel_type', 'petrol', 'Бензин', 'Benzin'),
('fuel_type', 'diesel', 'Дизель', 'Dizel'),
('fuel_type', 'electric', 'Электро', 'Električni'),
('fuel_type', 'hybrid', 'Гибрид', 'Hibrid'),
('fuel_type', 'lpg', 'Газ', 'Plin'),
('fuel_type', 'cng', 'Метан', 'Metan')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для body_type
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('body_type', 'sedan', 'Седан', 'Limuzina'),
('body_type', 'hatchback', 'Хэтчбек', 'Hečbek'),
('body_type', 'suv', 'Внедорожник', 'Džip'),
('body_type', 'coupe', 'Купе', 'Kupe'),
('body_type', 'wagon', 'Универсал', 'Karavan'),
('body_type', 'minivan', 'Минивэн', 'Monovolumen'),
('body_type', 'pickup', 'Пикап', 'Pikap'),
('body_type', 'convertible', 'Кабриолет', 'Kabriolet')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для storage
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('storage', '16GB', '16 ГБ', '16GB'),
('storage', '32GB', '32 ГБ', '32GB'),
('storage', '64GB', '64 ГБ', '64GB'),
('storage', '128GB', '128 ГБ', '128GB'),
('storage', '256GB', '256 ГБ', '256GB'),
('storage', '512GB', '512 ГБ', '512GB'),
('storage', '1TB', '1 ТБ', '1TB'),
('storage', '2TB', '2 ТБ', '2TB')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для operating_system
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('operating_system', 'iOS', 'iOS', 'iOS'),
('operating_system', 'Android', 'Android', 'Android'),
('operating_system', 'Windows', 'Windows', 'Windows'),
('operating_system', 'macOS', 'macOS', 'macOS'),
('operating_system', 'Linux', 'Linux', 'Linux'),
('operating_system', 'Other', 'Другая', 'Druga')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для ram
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('ram', '2GB', '2 ГБ', '2GB'),
('ram', '4GB', '4 ГБ', '4GB'),
('ram', '8GB', '8 ГБ', '8GB'),
('ram', '16GB', '16 ГБ', '16GB'),
('ram', '32GB', '32 ГБ', '32GB'),
('ram', '64GB', '64 ГБ', '64GB'),
('ram', '128GB', '128 ГБ', '128GB')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для storage_type
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('storage_type', 'HDD', 'HDD', 'HDD'),
('storage_type', 'SSD', 'SSD', 'SSD'),
('storage_type', 'NVMe', 'NVMe', 'NVMe'),
('storage_type', 'eMMC', 'eMMC', 'eMMC'),
('storage_type', 'Hybrid', 'Гибридный', 'Hibridni'),
('storage_type', 'Memory Card', 'Карта памяти', 'Memorijska kartica')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для rooms
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('rooms', 'studio', 'Студия', 'Garsonjera'),
('rooms', '1', '1 комната', '1 soba'),
('rooms', '2', '2 комнаты', '2 sobe'),
('rooms', '3', '3 комнаты', '3 sobe'),
('rooms', '4', '4 комнаты', '4 sobe'),
('rooms', '5', '5 комнат', '5 soba'),
('rooms', '6+', '6+ комнат', '6+ soba')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для warranty
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('warranty', 'no_warranty', 'Без гарантии', 'Bez garancije'),
('warranty', 'manufacturer', 'Гарантия производителя', 'Garancija proizvođača'),
('warranty', 'store', 'Гарантия магазина', 'Garancija prodavnice'),
('warranty', 'extended', 'Расширенная гарантия', 'Produžena garancija')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для return_policy
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('return_policy', 'no_returns', 'Возврат не принимается', 'Bez povrata'),
('return_policy', '7_days', '7 дней', '7 dana'),
('return_policy', '14_days', '14 дней', '14 dana'),
('return_policy', '30_days', '30 дней', '30 dana')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для service_type
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('service_type', 'consulting', 'Консультации', 'Konsultacije'),
('service_type', 'repair', 'Ремонт', 'Popravka'),
('service_type', 'installation', 'Установка', 'Instalacija'),
('service_type', 'maintenance', 'Обслуживание', 'Održavanje'),
('service_type', 'cleaning', 'Уборка', 'Čišćenje'),
('service_type', 'transport', 'Транспорт', 'Transport'),
('service_type', 'design', 'Дизайн', 'Dizajn'),
('service_type', 'development', 'Разработка', 'Razvoj'),
('service_type', 'education', 'Обучение', 'Edukacija'),
('service_type', 'other', 'Другое', 'Drugo')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для availability
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('availability', 'immediate', 'Немедленно', 'Odmah'),
('availability', 'within_24h', 'В течение 24 часов', 'U roku od 24h'),
('availability', 'within_week', 'В течение недели', 'U roku od nedelje'),
('availability', 'by_appointment', 'По записи', 'Po dogovoru'),
('availability', 'weekdays', 'Будние дни', 'Radnim danima'),
('availability', 'weekends', 'Выходные', 'Vikendom'),
('availability', '24_7', 'Круглосуточно', '24/7')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;

-- Переводы для service_area
INSERT INTO attribute_option_translations (attribute_name, option_value, ru_translation, sr_translation) VALUES
('service_area', 'local', 'Местный', 'Lokalno'),
('service_area', 'city', 'Город', 'Grad'),
('service_area', 'region', 'Регион', 'Region'),
('service_area', 'country', 'Страна', 'Država'),
('service_area', 'international', 'Международный', 'Međunarodno'),
('service_area', 'online', 'Онлайн', 'Online')
ON CONFLICT (attribute_name, option_value) DO UPDATE SET
    ru_translation = EXCLUDED.ru_translation,
    sr_translation = EXCLUDED.sr_translation;