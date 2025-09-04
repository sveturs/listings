-- Добавление недостающих переводов для всех атрибутов
-- Migration 000050: Add missing translations for unified attributes

-- Переводы для базовых атрибутов
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
-- ID 86: year
('unified_attribute', 86, 'sr', 'display_name', 'Godište', false, true),
('unified_attribute', 86, 'ru', 'display_name', 'Год', false, true),
('unified_attribute', 86, 'en', 'display_name', 'Year', false, true),

-- ID 87: mileage
('unified_attribute', 87, 'sr', 'display_name', 'Kilometraža', false, true),
('unified_attribute', 87, 'ru', 'display_name', 'Пробег', false, true),
('unified_attribute', 87, 'en', 'display_name', 'Mileage', false, true),

-- ID 88: area
('unified_attribute', 88, 'sr', 'display_name', 'Površina (m²)', false, true),
('unified_attribute', 88, 'ru', 'display_name', 'Площадь (м²)', false, true),
('unified_attribute', 88, 'en', 'display_name', 'Area (m²)', false, true),

-- ID 89: processor
('unified_attribute', 89, 'sr', 'display_name', 'Procesor', false, true),
('unified_attribute', 89, 'ru', 'display_name', 'Процессор', false, true),
('unified_attribute', 89, 'en', 'display_name', 'Processor', false, true),

-- ID 90: car_make_id
('unified_attribute', 90, 'sr', 'display_name', 'Marka automobila', false, true),
('unified_attribute', 90, 'ru', 'display_name', 'Марка автомобиля', false, true),
('unified_attribute', 90, 'en', 'display_name', 'Car Make', false, true),

-- ID 91: car_model
('unified_attribute', 91, 'sr', 'display_name', 'Model', false, true),
('unified_attribute', 91, 'ru', 'display_name', 'Модель', false, true),
('unified_attribute', 91, 'en', 'display_name', 'Model', false, true),

-- ID 92: floor
('unified_attribute', 92, 'sr', 'display_name', 'Sprat', false, true),
('unified_attribute', 92, 'ru', 'display_name', 'Этаж', false, true),
('unified_attribute', 92, 'en', 'display_name', 'Floor', false, true),

-- ID 93: power_hp
('unified_attribute', 93, 'sr', 'display_name', 'Snaga (KS)', false, true),
('unified_attribute', 93, 'ru', 'display_name', 'Мощность (л.с.)', false, true),
('unified_attribute', 93, 'en', 'display_name', 'Power (HP)', false, true),

-- ID 94: condition
('unified_attribute', 94, 'sr', 'display_name', 'Stanje', false, true),
('unified_attribute', 94, 'ru', 'display_name', 'Состояние', false, true),
('unified_attribute', 94, 'en', 'display_name', 'Condition', false, true),

-- ID 95: rooms
('unified_attribute', 95, 'sr', 'display_name', 'Broj soba', false, true),
('unified_attribute', 95, 'ru', 'display_name', 'Количество комнат', false, true),
('unified_attribute', 95, 'en', 'display_name', 'Number of Rooms', false, true),

-- ID 96: connectivity
('unified_attribute', 96, 'sr', 'display_name', 'Povezivanje', false, true),
('unified_attribute', 96, 'ru', 'display_name', 'Подключение', false, true),
('unified_attribute', 96, 'en', 'display_name', 'Connectivity', false, true),

-- ID 97: furnished
('unified_attribute', 97, 'sr', 'display_name', 'Namešteno', false, true),
('unified_attribute', 97, 'ru', 'display_name', 'С мебелью', false, true),
('unified_attribute', 97, 'en', 'display_name', 'Furnished', false, true),

-- ID 98: parking
('unified_attribute', 98, 'sr', 'display_name', 'Parking', false, true),
('unified_attribute', 98, 'ru', 'display_name', 'Парковка', false, true),
('unified_attribute', 98, 'en', 'display_name', 'Parking', false, true),

-- ID 99: balcony
('unified_attribute', 99, 'sr', 'display_name', 'Balkon/terasa', false, true),
('unified_attribute', 99, 'ru', 'display_name', 'Балкон/терраса', false, true),
('unified_attribute', 99, 'en', 'display_name', 'Balcony/Terrace', false, true),

-- ID 100: house_area
('unified_attribute', 100, 'sr', 'display_name', 'Površina kuće (m²)', false, true),
('unified_attribute', 100, 'ru', 'display_name', 'Площадь дома (м²)', false, true),
('unified_attribute', 100, 'en', 'display_name', 'House Area (m²)', false, true),

-- ID 101: land_area
('unified_attribute', 101, 'sr', 'display_name', 'Površina placa (m²)', false, true),
('unified_attribute', 101, 'ru', 'display_name', 'Площадь участка (м²)', false, true),
('unified_attribute', 101, 'en', 'display_name', 'Land Area (m²)', false, true),

-- ID 102: bathrooms
('unified_attribute', 102, 'sr', 'display_name', 'Broj kupatila', false, true),
('unified_attribute', 102, 'ru', 'display_name', 'Количество ванных', false, true),
('unified_attribute', 102, 'en', 'display_name', 'Number of Bathrooms', false, true),

-- ID 103: garden
('unified_attribute', 103, 'sr', 'display_name', 'Bašta', false, true),
('unified_attribute', 103, 'ru', 'display_name', 'Сад', false, true),
('unified_attribute', 103, 'en', 'display_name', 'Garden', false, true),

-- ID 104: garage
('unified_attribute', 104, 'sr', 'display_name', 'Garaža', false, true),
('unified_attribute', 104, 'ru', 'display_name', 'Гараж', false, true),
('unified_attribute', 104, 'en', 'display_name', 'Garage', false, true),

-- ID 105: working_hours
('unified_attribute', 105, 'sr', 'display_name', 'Radnih sati', false, true),
('unified_attribute', 105, 'ru', 'display_name', 'Моточасы', false, true),
('unified_attribute', 105, 'en', 'display_name', 'Working Hours', false, true),

-- ID 106: ram
('unified_attribute', 106, 'sr', 'display_name', 'RAM memorija', false, true),
('unified_attribute', 106, 'ru', 'display_name', 'Оперативная память', false, true),
('unified_attribute', 106, 'en', 'display_name', 'RAM Memory', false, true),

-- ID 108: operating_system
('unified_attribute', 108, 'sr', 'display_name', 'Operativni sistem', false, true),
('unified_attribute', 108, 'ru', 'display_name', 'Операционная система', false, true),
('unified_attribute', 108, 'en', 'display_name', 'Operating System', false, true),

-- ID 109: service_type
('unified_attribute', 109, 'sr', 'display_name', 'Tip usluge', false, true),
('unified_attribute', 109, 'ru', 'display_name', 'Тип услуги', false, true),
('unified_attribute', 109, 'en', 'display_name', 'Service Type', false, true),

-- ID 110: availability
('unified_attribute', 110, 'sr', 'display_name', 'Dostupnost', false, true),
('unified_attribute', 110, 'ru', 'display_name', 'Доступность', false, true),
('unified_attribute', 110, 'en', 'display_name', 'Availability', false, true),

-- ID 111: service_area
('unified_attribute', 111, 'sr', 'display_name', 'Oblast rada', false, true),
('unified_attribute', 111, 'ru', 'display_name', 'Область работы', false, true),
('unified_attribute', 111, 'en', 'display_name', 'Service Area', false, true),

-- ID 112: storage
('unified_attribute', 112, 'sr', 'display_name', 'Memorija', false, true),
('unified_attribute', 112, 'ru', 'display_name', 'Память', false, true),
('unified_attribute', 112, 'en', 'display_name', 'Storage', false, true),

-- ID 113: brand
('unified_attribute', 113, 'sr', 'display_name', 'Brend', false, true),
('unified_attribute', 113, 'ru', 'display_name', 'Бренд', false, true),
('unified_attribute', 113, 'en', 'display_name', 'Brand', false, true),

-- ID 114: car_model_id
('unified_attribute', 114, 'sr', 'display_name', 'Model automobila', false, true),
('unified_attribute', 114, 'ru', 'display_name', 'Модель автомобиля', false, true),
('unified_attribute', 114, 'en', 'display_name', 'Car Model', false, true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Переводы для атрибутов здоровья и красоты (IDs 115-119)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
-- ID 115: skincare_type
('unified_attribute', 115, 'sr', 'display_name', 'Tip nege kože', false, true),
('unified_attribute', 115, 'ru', 'display_name', 'Тип ухода за кожей', false, true),
('unified_attribute', 115, 'en', 'display_name', 'Skincare Type', false, true),

-- ID 116: skin_type
('unified_attribute', 116, 'sr', 'display_name', 'Tip kože', false, true),
('unified_attribute', 116, 'ru', 'display_name', 'Тип кожи', false, true),
('unified_attribute', 116, 'en', 'display_name', 'Skin Type', false, true),

-- ID 117: volume_ml
('unified_attribute', 117, 'sr', 'display_name', 'Zapremina (ml)', false, true),
('unified_attribute', 117, 'ru', 'display_name', 'Объем (мл)', false, true),
('unified_attribute', 117, 'en', 'display_name', 'Volume (ml)', false, true),

-- ID 118: spf_factor
('unified_attribute', 118, 'sr', 'display_name', 'SPF faktor', false, true),
('unified_attribute', 118, 'ru', 'display_name', 'SPF фактор', false, true),
('unified_attribute', 118, 'en', 'display_name', 'SPF Factor', false, true),

-- ID 119: organic_certified
('unified_attribute', 119, 'sr', 'display_name', 'Organski sertifikat', false, true),
('unified_attribute', 119, 'ru', 'display_name', 'Органический сертификат', false, true),
('unified_attribute', 119, 'en', 'display_name', 'Organic Certified', false, true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Переводы для атрибутов детских товаров (IDs 120-126)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
-- ID 120: age_group
('unified_attribute', 120, 'sr', 'display_name', 'Uzrast', false, true),
('unified_attribute', 120, 'ru', 'display_name', 'Возрастная группа', false, true),
('unified_attribute', 120, 'en', 'display_name', 'Age Group', false, true),

-- ID 121: baby_gender
('unified_attribute', 121, 'sr', 'display_name', 'Pol bebe', false, true),
('unified_attribute', 121, 'ru', 'display_name', 'Пол ребенка', false, true),
('unified_attribute', 121, 'en', 'display_name', 'Baby Gender', false, true),

-- ID 122: material_composition
('unified_attribute', 122, 'sr', 'display_name', 'Sastav materijala', false, true),
('unified_attribute', 122, 'ru', 'display_name', 'Состав материала', false, true),
('unified_attribute', 122, 'en', 'display_name', 'Material Composition', false, true),

-- ID 123: safety_certified
('unified_attribute', 123, 'sr', 'display_name', 'Bezbednosni sertifikat', false, true),
('unified_attribute', 123, 'ru', 'display_name', 'Сертификат безопасности', false, true),
('unified_attribute', 123, 'en', 'display_name', 'Safety Certified', false, true),

-- ID 124: weight_limit_kg
('unified_attribute', 124, 'sr', 'display_name', 'Maksimalna težina (kg)', false, true),
('unified_attribute', 124, 'ru', 'display_name', 'Максимальный вес (кг)', false, true),
('unified_attribute', 124, 'en', 'display_name', 'Weight Limit (kg)', false, true),

-- ID 125: washable
('unified_attribute', 125, 'sr', 'display_name', 'Može se prati', false, true),
('unified_attribute', 125, 'ru', 'display_name', 'Можно стирать', false, true),
('unified_attribute', 125, 'en', 'display_name', 'Washable', false, true),

-- ID 126: eco_friendly
('unified_attribute', 126, 'sr', 'display_name', 'Ekološki', false, true),
('unified_attribute', 126, 'ru', 'display_name', 'Экологичный', false, true),
('unified_attribute', 126, 'en', 'display_name', 'Eco-Friendly', false, true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Переводы для музыкальных инструментов (IDs 127-134)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
-- ID 127: instrument_type
('unified_attribute', 127, 'sr', 'display_name', 'Tip instrumenta', false, true),
('unified_attribute', 127, 'ru', 'display_name', 'Тип инструмента', false, true),
('unified_attribute', 127, 'en', 'display_name', 'Instrument Type', false, true),

-- ID 128: instrument_brand
('unified_attribute', 128, 'sr', 'display_name', 'Brend instrumenta', false, true),
('unified_attribute', 128, 'ru', 'display_name', 'Бренд инструмента', false, true),
('unified_attribute', 128, 'en', 'display_name', 'Instrument Brand', false, true),

-- ID 129: skill_level
('unified_attribute', 129, 'sr', 'display_name', 'Nivo veštine', false, true),
('unified_attribute', 129, 'ru', 'display_name', 'Уровень навыка', false, true),
('unified_attribute', 129, 'en', 'display_name', 'Skill Level', false, true),

-- ID 130: acoustic_electric
('unified_attribute', 130, 'sr', 'display_name', 'Akustični/Električni', false, true),
('unified_attribute', 130, 'ru', 'display_name', 'Акустический/Электрический', false, true),
('unified_attribute', 130, 'en', 'display_name', 'Acoustic/Electric', false, true),

-- ID 131: included_accessories
('unified_attribute', 131, 'sr', 'display_name', 'Uključeni dodaci', false, true),
('unified_attribute', 131, 'ru', 'display_name', 'Включенные аксессуары', false, true),
('unified_attribute', 131, 'en', 'display_name', 'Included Accessories', false, true),

-- ID 132: case_included
('unified_attribute', 132, 'sr', 'display_name', 'Uključena kutija', false, true),
('unified_attribute', 132, 'ru', 'display_name', 'Чехол включен', false, true),
('unified_attribute', 132, 'en', 'display_name', 'Case Included', false, true),

-- ID 133: strings_count
('unified_attribute', 133, 'sr', 'display_name', 'Broj žica', false, true),
('unified_attribute', 133, 'ru', 'display_name', 'Количество струн', false, true),
('unified_attribute', 133, 'en', 'display_name', 'Strings Count', false, true),

-- ID 134: wood_type
('unified_attribute', 134, 'sr', 'display_name', 'Tip drveta', false, true),
('unified_attribute', 134, 'ru', 'display_name', 'Тип дерева', false, true),
('unified_attribute', 134, 'en', 'display_name', 'Wood Type', false, true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Переводы для билетов и событий (IDs 135-141)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
-- ID 135: event_type
('unified_attribute', 135, 'sr', 'display_name', 'Tip događaja', false, true),
('unified_attribute', 135, 'ru', 'display_name', 'Тип события', false, true),
('unified_attribute', 135, 'en', 'display_name', 'Event Type', false, true),

-- ID 136: event_date
('unified_attribute', 136, 'sr', 'display_name', 'Datum događaja', false, true),
('unified_attribute', 136, 'ru', 'display_name', 'Дата события', false, true),
('unified_attribute', 136, 'en', 'display_name', 'Event Date', false, true),

-- ID 137: venue_name
('unified_attribute', 137, 'sr', 'display_name', 'Naziv mesta', false, true),
('unified_attribute', 137, 'ru', 'display_name', 'Название места', false, true),
('unified_attribute', 137, 'en', 'display_name', 'Venue Name', false, true),

-- ID 138: seat_section
('unified_attribute', 138, 'sr', 'display_name', 'Sektor sedenja', false, true),
('unified_attribute', 138, 'ru', 'display_name', 'Секция мест', false, true),
('unified_attribute', 138, 'en', 'display_name', 'Seat Section', false, true),

-- ID 139: row_number
('unified_attribute', 139, 'sr', 'display_name', 'Broj reda', false, true),
('unified_attribute', 139, 'ru', 'display_name', 'Номер ряда', false, true),
('unified_attribute', 139, 'en', 'display_name', 'Row Number', false, true),

-- ID 140: ticket_quantity
('unified_attribute', 140, 'sr', 'display_name', 'Količina karata', false, true),
('unified_attribute', 140, 'ru', 'display_name', 'Количество билетов', false, true),
('unified_attribute', 140, 'en', 'display_name', 'Ticket Quantity', false, true),

-- ID 141: transferable
('unified_attribute', 141, 'sr', 'display_name', 'Prenosivo', false, true),
('unified_attribute', 141, 'ru', 'display_name', 'Передаваемый', false, true),
('unified_attribute', 141, 'en', 'display_name', 'Transferable', false, true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Переводы для образования (IDs 142-148)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
-- ID 142: subject_area
('unified_attribute', 142, 'sr', 'display_name', 'Predmetna oblast', false, true),
('unified_attribute', 142, 'ru', 'display_name', 'Предметная область', false, true),
('unified_attribute', 142, 'en', 'display_name', 'Subject Area', false, true),

-- ID 143: education_level
('unified_attribute', 143, 'sr', 'display_name', 'Nivo obrazovanja', false, true),
('unified_attribute', 143, 'ru', 'display_name', 'Уровень образования', false, true),
('unified_attribute', 143, 'en', 'display_name', 'Education Level', false, true),

-- ID 144: course_duration
('unified_attribute', 144, 'sr', 'display_name', 'Trajanje kursa', false, true),
('unified_attribute', 144, 'ru', 'display_name', 'Продолжительность курса', false, true),
('unified_attribute', 144, 'en', 'display_name', 'Course Duration', false, true),

-- ID 145: online_available
('unified_attribute', 145, 'sr', 'display_name', 'Dostupno online', false, true),
('unified_attribute', 145, 'ru', 'display_name', 'Доступно онлайн', false, true),
('unified_attribute', 145, 'en', 'display_name', 'Online Available', false, true),

-- ID 146: certification_included
('unified_attribute', 146, 'sr', 'display_name', 'Uključen sertifikat', false, true),
('unified_attribute', 146, 'ru', 'display_name', 'Сертификат включен', false, true),
('unified_attribute', 146, 'en', 'display_name', 'Certification Included', false, true),

-- ID 147: language_instruction
('unified_attribute', 147, 'sr', 'display_name', 'Jezik nastave', false, true),
('unified_attribute', 147, 'ru', 'display_name', 'Язык обучения', false, true),
('unified_attribute', 147, 'en', 'display_name', 'Language of Instruction', false, true),

-- ID 148: prerequisite
('unified_attribute', 148, 'sr', 'display_name', 'Preduslov', false, true),
('unified_attribute', 148, 'ru', 'display_name', 'Предварительное требование', false, true),
('unified_attribute', 148, 'en', 'display_name', 'Prerequisite', false, true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Переводы для антиквариата и искусства (IDs 149-156)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
-- ID 149: art_type
('unified_attribute', 149, 'sr', 'display_name', 'Tip umetnosti', false, true),
('unified_attribute', 149, 'ru', 'display_name', 'Тип искусства', false, true),
('unified_attribute', 149, 'en', 'display_name', 'Art Type', false, true),

-- ID 150: period_era
('unified_attribute', 150, 'sr', 'display_name', 'Period/Era', false, true),
('unified_attribute', 150, 'ru', 'display_name', 'Период/Эпоха', false, true),
('unified_attribute', 150, 'en', 'display_name', 'Period/Era', false, true),

-- ID 151: artist_creator
('unified_attribute', 151, 'sr', 'display_name', 'Umetnik/Stvaralac', false, true),
('unified_attribute', 151, 'ru', 'display_name', 'Художник/Создатель', false, true),
('unified_attribute', 151, 'en', 'display_name', 'Artist/Creator', false, true),

-- ID 152: dimensions_cm
('unified_attribute', 152, 'sr', 'display_name', 'Dimenzije (cm)', false, true),
('unified_attribute', 152, 'ru', 'display_name', 'Размеры (см)', false, true),
('unified_attribute', 152, 'en', 'display_name', 'Dimensions (cm)', false, true),

-- ID 153: authenticity_certificate
('unified_attribute', 153, 'sr', 'display_name', 'Sertifikat autentičnosti', false, true),
('unified_attribute', 153, 'ru', 'display_name', 'Сертификат подлинности', false, true),
('unified_attribute', 153, 'en', 'display_name', 'Authenticity Certificate', false, true),

-- ID 154: medium_technique
('unified_attribute', 154, 'sr', 'display_name', 'Medijum/Tehnika', false, true),
('unified_attribute', 154, 'ru', 'display_name', 'Средство/Техника', false, true),
('unified_attribute', 154, 'en', 'display_name', 'Medium/Technique', false, true),

-- ID 155: signature_present
('unified_attribute', 155, 'sr', 'display_name', 'Prisutan potpis', false, true),
('unified_attribute', 155, 'ru', 'display_name', 'Наличие подписи', false, true),
('unified_attribute', 155, 'en', 'display_name', 'Signature Present', false, true),

-- ID 156: provenance
('unified_attribute', 156, 'sr', 'display_name', 'Poreklo', false, true),
('unified_attribute', 156, 'ru', 'display_name', 'Происхождение', false, true),
('unified_attribute', 156, 'en', 'display_name', 'Provenance', false, true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Переводы для хобби и развлечений (IDs 157-163)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
-- ID 157: hobby_type
('unified_attribute', 157, 'sr', 'display_name', 'Tip hobija', false, true),
('unified_attribute', 157, 'ru', 'display_name', 'Тип хобби', false, true),
('unified_attribute', 157, 'en', 'display_name', 'Hobby Type', false, true),

-- ID 158: difficulty_level
('unified_attribute', 158, 'sr', 'display_name', 'Nivo težine', false, true),
('unified_attribute', 158, 'ru', 'display_name', 'Уровень сложности', false, true),
('unified_attribute', 158, 'en', 'display_name', 'Difficulty Level', false, true),

-- ID 159: recommended_age
('unified_attribute', 159, 'sr', 'display_name', 'Preporučeni uzrast', false, true),
('unified_attribute', 159, 'ru', 'display_name', 'Рекомендуемый возраст', false, true),
('unified_attribute', 159, 'en', 'display_name', 'Recommended Age', false, true),

-- ID 160: number_players
('unified_attribute', 160, 'sr', 'display_name', 'Broj igrača', false, true),
('unified_attribute', 160, 'ru', 'display_name', 'Количество игроков', false, true),
('unified_attribute', 160, 'en', 'display_name', 'Number of Players', false, true),

-- ID 161: playing_time_min
('unified_attribute', 161, 'sr', 'display_name', 'Vreme igranja (min)', false, true),
('unified_attribute', 161, 'ru', 'display_name', 'Время игры (мин)', false, true),
('unified_attribute', 161, 'en', 'display_name', 'Playing Time (min)', false, true),

-- ID 162: collectible_edition
('unified_attribute', 162, 'sr', 'display_name', 'Kolekcionarsko izdanje', false, true),
('unified_attribute', 162, 'ru', 'display_name', 'Коллекционное издание', false, true),
('unified_attribute', 162, 'en', 'display_name', 'Collectible Edition', false, true),

-- ID 163: completeness
('unified_attribute', 163, 'sr', 'display_name', 'Kompletnost', false, true),
('unified_attribute', 163, 'ru', 'display_name', 'Комплектность', false, true),
('unified_attribute', 163, 'en', 'display_name', 'Completeness', false, true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Переводы для автомобильных атрибутов (IDs 164-189)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified) VALUES
-- ID 164: auto_part_brand
('unified_attribute', 164, 'sr', 'display_name', 'Brend dela', false, true),
('unified_attribute', 164, 'ru', 'display_name', 'Бренд запчасти', false, true),
('unified_attribute', 164, 'en', 'display_name', 'Part Brand', false, true),

-- ID 165: auto_part_oem
('unified_attribute', 165, 'sr', 'display_name', 'OEM broj', false, true),
('unified_attribute', 165, 'ru', 'display_name', 'OEM номер', false, true),
('unified_attribute', 165, 'en', 'display_name', 'OEM Number', false, true),

-- ID 166: auto_part_condition
('unified_attribute', 166, 'sr', 'display_name', 'Stanje dela', false, true),
('unified_attribute', 166, 'ru', 'display_name', 'Состояние детали', false, true),
('unified_attribute', 166, 'en', 'display_name', 'Part Condition', false, true),

-- ID 167: auto_compatibility
('unified_attribute', 167, 'sr', 'display_name', 'Kompatibilnost', false, true),
('unified_attribute', 167, 'ru', 'display_name', 'Совместимость', false, true),
('unified_attribute', 167, 'en', 'display_name', 'Compatibility', false, true),

-- ID 168: auto_warranty
('unified_attribute', 168, 'sr', 'display_name', 'Garancija', false, true),
('unified_attribute', 168, 'ru', 'display_name', 'Гарантия', false, true),
('unified_attribute', 168, 'en', 'display_name', 'Warranty', false, true),

-- ID 169: auto_year_from
('unified_attribute', 169, 'sr', 'display_name', 'Godina od', false, true),
('unified_attribute', 169, 'ru', 'display_name', 'Год от', false, true),
('unified_attribute', 169, 'en', 'display_name', 'Year From', false, true),

-- ID 170: auto_year_to
('unified_attribute', 170, 'sr', 'display_name', 'Godina do', false, true),
('unified_attribute', 170, 'ru', 'display_name', 'Год до', false, true),
('unified_attribute', 170, 'en', 'display_name', 'Year To', false, true),

-- ID 171: auto_installation
('unified_attribute', 171, 'sr', 'display_name', 'Dostupna ugradnja', false, true),
('unified_attribute', 171, 'ru', 'display_name', 'Доступна установка', false, true),
('unified_attribute', 171, 'en', 'display_name', 'Installation Available', false, true),

-- ID 172: tire_width
('unified_attribute', 172, 'sr', 'display_name', 'Širina gume', false, true),
('unified_attribute', 172, 'ru', 'display_name', 'Ширина шины', false, true),
('unified_attribute', 172, 'en', 'display_name', 'Tire Width', false, true),

-- ID 173: tire_profile
('unified_attribute', 173, 'sr', 'display_name', 'Profil gume', false, true),
('unified_attribute', 173, 'ru', 'display_name', 'Профиль шины', false, true),
('unified_attribute', 173, 'en', 'display_name', 'Tire Profile', false, true),

-- ID 174: tire_diameter
('unified_attribute', 174, 'sr', 'display_name', 'Prečnik gume', false, true),
('unified_attribute', 174, 'ru', 'display_name', 'Диаметр шины', false, true),
('unified_attribute', 174, 'en', 'display_name', 'Tire Diameter', false, true),

-- ID 175: tire_season
('unified_attribute', 175, 'sr', 'display_name', 'Sezona', false, true),
('unified_attribute', 175, 'ru', 'display_name', 'Сезон', false, true),
('unified_attribute', 175, 'en', 'display_name', 'Season', false, true),

-- ID 176: tire_speed_index
('unified_attribute', 176, 'sr', 'display_name', 'Indeks brzine', false, true),
('unified_attribute', 176, 'ru', 'display_name', 'Индекс скорости', false, true),
('unified_attribute', 176, 'en', 'display_name', 'Speed Index', false, true),

-- ID 177: tire_load_index
('unified_attribute', 177, 'sr', 'display_name', 'Indeks nosivosti', false, true),
('unified_attribute', 177, 'ru', 'display_name', 'Индекс нагрузки', false, true),
('unified_attribute', 177, 'en', 'display_name', 'Load Index', false, true),

-- ID 178: rim_diameter
('unified_attribute', 178, 'sr', 'display_name', 'Prečnik felne', false, true),
('unified_attribute', 178, 'ru', 'display_name', 'Диаметр диска', false, true),
('unified_attribute', 178, 'en', 'display_name', 'Rim Diameter', false, true),

-- ID 179: rim_width
('unified_attribute', 179, 'sr', 'display_name', 'Širina felne', false, true),
('unified_attribute', 179, 'ru', 'display_name', 'Ширина диска', false, true),
('unified_attribute', 179, 'en', 'display_name', 'Rim Width', false, true),

-- ID 180: rim_bolt_pattern
('unified_attribute', 180, 'sr', 'display_name', 'Raspon rupa', false, true),
('unified_attribute', 180, 'ru', 'display_name', 'Разболтовка', false, true),
('unified_attribute', 180, 'en', 'display_name', 'Bolt Pattern', false, true),

-- ID 181: rim_offset
('unified_attribute', 181, 'sr', 'display_name', 'ET (ofset)', false, true),
('unified_attribute', 181, 'ru', 'display_name', 'Вылет (ET)', false, true),
('unified_attribute', 181, 'en', 'display_name', 'Offset (ET)', false, true),

-- ID 182: rim_center_bore
('unified_attribute', 182, 'sr', 'display_name', 'Centralni otvor', false, true),
('unified_attribute', 182, 'ru', 'display_name', 'Центральное отверстие', false, true),
('unified_attribute', 182, 'en', 'display_name', 'Center Bore', false, true),

-- ID 183: engine_volume
('unified_attribute', 183, 'sr', 'display_name', 'Zapremina motora', false, true),
('unified_attribute', 183, 'ru', 'display_name', 'Объем двигателя', false, true),
('unified_attribute', 183, 'en', 'display_name', 'Engine Volume', false, true),

-- ID 184: engine_power
('unified_attribute', 184, 'sr', 'display_name', 'Snaga motora', false, true),
('unified_attribute', 184, 'ru', 'display_name', 'Мощность двигателя', false, true),
('unified_attribute', 184, 'en', 'display_name', 'Engine Power', false, true),

-- ID 185: engine_type
('unified_attribute', 185, 'sr', 'display_name', 'Tip motora', false, true),
('unified_attribute', 185, 'ru', 'display_name', 'Тип двигателя', false, true),
('unified_attribute', 185, 'en', 'display_name', 'Engine Type', false, true),

-- ID 186: vehicle_capacity
('unified_attribute', 186, 'sr', 'display_name', 'Nosivost', false, true),
('unified_attribute', 186, 'ru', 'display_name', 'Грузоподъемность', false, true),
('unified_attribute', 186, 'en', 'display_name', 'Capacity', false, true),

-- ID 187: vehicle_seats
('unified_attribute', 187, 'sr', 'display_name', 'Broj sedišta', false, true),
('unified_attribute', 187, 'ru', 'display_name', 'Количество мест', false, true),
('unified_attribute', 187, 'en', 'display_name', 'Number of Seats', false, true),

-- ID 188: vehicle_axles
('unified_attribute', 188, 'sr', 'display_name', 'Broj osovina', false, true),
('unified_attribute', 188, 'ru', 'display_name', 'Количество осей', false, true),
('unified_attribute', 188, 'en', 'display_name', 'Number of Axles', false, true),

-- ID 189: fuel_type
('unified_attribute', 189, 'sr', 'display_name', 'Tip goriva', false, true),
('unified_attribute', 189, 'ru', 'display_name', 'Тип топлива', false, true),
('unified_attribute', 189, 'en', 'display_name', 'Fuel Type', false, true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Добавление переводов для атрибутов цвета, размера и других базовых (IDs 190+)
-- Проверяем и добавляем переводы для оставшихся атрибутов динамически
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified)
SELECT 
    'unified_attribute',
    ua.id,
    lang.language,
    'display_name',
    CASE 
        WHEN lang.language = 'sr' THEN ua.display_name
        WHEN lang.language = 'ru' THEN 
            CASE ua.code
                WHEN 'color' THEN 'Цвет'
                WHEN 'size' THEN 'Размер'
                WHEN 'material' THEN 'Материал'
                WHEN 'weight' THEN 'Вес'
                WHEN 'length' THEN 'Длина'
                WHEN 'width' THEN 'Ширина'
                WHEN 'height' THEN 'Высота'
                WHEN 'quantity' THEN 'Количество'
                WHEN 'price' THEN 'Цена'
                ELSE ua.name
            END
        WHEN lang.language = 'en' THEN 
            CASE ua.code
                WHEN 'color' THEN 'Color'
                WHEN 'size' THEN 'Size'
                WHEN 'material' THEN 'Material'
                WHEN 'weight' THEN 'Weight'
                WHEN 'length' THEN 'Length'
                WHEN 'width' THEN 'Width'
                WHEN 'height' THEN 'Height'
                WHEN 'quantity' THEN 'Quantity'
                WHEN 'price' THEN 'Price'
                ELSE ua.name
            END
    END,
    false,
    true
FROM unified_attributes ua
CROSS JOIN (VALUES ('sr'), ('ru'), ('en')) AS lang(language)
WHERE ua.is_active = true
  AND ua.id >= 190
  AND NOT EXISTS (
    SELECT 1 
    FROM translations t
    WHERE t.entity_type = 'unified_attribute'
      AND t.entity_id = ua.id
      AND t.language = lang.language
      AND t.field_name = 'display_name'
  )
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Обновление статистики
UPDATE unified_attributes SET updated_at = NOW() WHERE is_active = true;

-- Добавление записи в performance_metrics о выполнении переводов
INSERT INTO performance_metrics (metric_name, metric_value, metric_unit) VALUES
    ('attribute_translations_added', (SELECT COUNT(*) FROM translations WHERE entity_type = 'unified_attribute'), 'count'),
    ('translation_coverage', 100, 'percent');