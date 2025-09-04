-- Добавление недостающих сербских переводов для категорий
-- Migration 000046: Complete Serbian translations for marketplace categories

-- Автомобили - дополнительные подкатегории
INSERT INTO translations (key_name, language, translation) VALUES
    ('category.automotive.additional_equipment', 'sr', 'Dodatna oprema'),
    ('category.automotive.electrical_electronic', 'sr', 'Električni i elektronski delovi'),
    ('category.automotive.tires_wheels', 'sr', 'Gume i točkovi'),
    ('category.automotive.body_parts', 'sr', 'Karoserija i delovi'),
    ('category.automotive.engine_parts', 'sr', 'Motor i delovi motora'),
    ('category.automotive.cooling_system', 'sr', 'Sistem hlađenja'),
    ('category.automotive.suspension_steering', 'sr', 'Oslanjanje i upravljanje'),
    ('category.automotive.exhaust_system', 'sr', 'Izduvni sistem'),
    ('category.automotive.fuel_system', 'sr', 'Sistem za gorivo'),
    ('category.automotive.lighting', 'sr', 'Osvetljenje'),
    ('category.automotive.interior_accessories', 'sr', 'Unutrašnji pribor'),
    ('category.automotive.exterior_accessories', 'sr', 'Spoljašnji pribor'),
    ('category.automotive.car_care', 'sr', 'Nega automobila'),
    ('category.automotive.tools_equipment', 'sr', 'Alati i oprema'),
    ('category.automotive.safety_security', 'sr', 'Bezbednost i sigurnost'),
    ('category.automotive.navigation_electronics', 'sr', 'Navigacija i elektronika'),
    ('category.automotive.car_audio', 'sr', 'Auto hi-fi'),
    ('category.automotive.performance_tuning', 'sr', 'Performanse i tuning'),
    ('category.automotive.vintage_classic', 'sr', 'Oldtajmer i klasični automobili'),
    ('category.automotive.motorcycles', 'sr', 'Motocikli'),
    ('category.automotive.bicycles', 'sr', 'Bicikli'),
    ('category.automotive.boats_marine', 'sr', 'Čamci i pomorstvo'),
    ('category.automotive.trucks_commercial', 'sr', 'Kamioni i komercijalna vozila'),
    ('category.automotive.trailers_parts', 'sr', 'Prikolice i delovi'),
    ('category.automotive.racing_sports', 'sr', 'Trka i sport'),
    ('category.automotive.maintenance_repair', 'sr', 'Održavanje i popravke'),
    ('category.automotive.documentation', 'sr', 'Dokumentacija'),
    ('category.automotive.insurance_financing', 'sr', 'Osiguranje i finansiranje')

-- Резервная обработка потенциальных дубликатов
ON CONFLICT (key_name, language) DO UPDATE SET
    translation = EXCLUDED.translation,
    updated_at = NOW();

-- Обновление счетчика переводов для отчетности
INSERT INTO translation_stats (category, total_translations, last_updated) 
VALUES ('automotive_serbian_additions', 28, NOW()) 
ON CONFLICT (category) DO UPDATE SET 
    total_translations = EXCLUDED.total_translations,
    last_updated = NOW();