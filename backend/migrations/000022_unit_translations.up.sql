CREATE TABLE IF NOT EXISTS unit_translations (
    unit VARCHAR(20) NOT NULL,
    language VARCHAR(10) NOT NULL,
    translated_unit VARCHAR(20) NOT NULL,
    display_format VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (unit, language)
);

-- Заполняем переводы для основных единиц измерения
INSERT INTO unit_translations (unit, language, translated_unit, display_format) VALUES
-- Площадь
('m²', 'en', 'sq.m.', '%g sq.m.'),
('m²', 'ru', 'м²', '%g м²'),
('m²', 'sr', 'm²', '%g m²'),

-- Соток
('ar', 'en', 'acres', '%g acres'),
('ar', 'ru', 'сот.', '%g сот.'),
('ar', 'sr', 'ar', '%g ar'),

-- Километры
('km', 'en', 'km', '%g km'),
('km', 'ru', 'км', '%g км'),
('km', 'sr', 'km', '%g km'),

-- Литры
('l', 'en', 'l', '%g l'),
('l', 'ru', 'л', '%g л'),
('l', 'sr', 'l', '%g l'),

-- Комнаты
('soba', 'en', 'room', '%g rooms'),
('soba', 'ru', 'комн.', '%g комн.'),
('soba', 'sr', 'soba', '%g soba'),

-- Этажи
('sprat', 'en', 'floor', '%g floor'),
('sprat', 'ru', 'эт.', '%g эт.'),
('sprat', 'sr', 'sprat', '%g sprat'),

-- Лошадиные силы
('ks', 'en', 'hp', '%g hp'),
('ks', 'ru', 'л.с.', '%g л.с.'),
('ks', 'sr', 'ks', '%g ks'),

-- Дюймы
('inč', 'en', 'inch', '%g"'),
('inč', 'ru', 'дюйм', '%g"'),
('inč', 'sr', 'inč', '%g"')

ON CONFLICT (unit, language) DO UPDATE 
SET translated_unit = EXCLUDED.translated_unit,
    display_format = EXCLUDED.display_format;