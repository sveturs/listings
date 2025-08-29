-- Добавляем переводы для объявлений которые еще не имеют переводов

-- iPhone 14 Pro Max 256GB Deep Purple (ID: 172)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 172, 'ru', 'title', 'iPhone 14 Pro Max 256GB Deep Purple', true),
    ('listing', 172, 'ru', 'description', 'Новый iPhone 14 Pro Max в заводской пленке. Цвет Deep Purple. 256GB памяти. Dynamic Island, A16 Bionic, камера 48MP. Полный комплект с зарядкой.', true),
    ('listing', 172, 'en', 'title', 'iPhone 14 Pro Max 256GB Deep Purple', true),
    ('listing', 172, 'en', 'description', 'New iPhone 14 Pro Max in factory film. Deep Purple color. 256GB storage. Dynamic Island, A16 Bionic, 48MP camera. Complete set with charger.', true),
    ('listing', 172, 'sr', 'title', 'iPhone 14 Pro Max 256GB Deep Purple', true),
    ('listing', 172, 'sr', 'description', 'Novi iPhone 14 Pro Max u fabričkoj foliji. Deep Purple boja. 256GB memorije. Dynamic Island, A16 Bionic, 48MP kamera. Kompletan set sa punjačem.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Samsung Galaxy S23 Ultra 512GB (ID: 173)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 173, 'ru', 'title', 'Samsung Galaxy S23 Ultra 512GB', true),
    ('listing', 173, 'ru', 'description', 'Samsung Galaxy S23 Ultra со стилусом S Pen. 512GB памяти, Snapdragon 8 Gen 2, камера 200MP. AMOLED 2X дисплей 120Hz. Гарантия.', true),
    ('listing', 173, 'en', 'title', 'Samsung Galaxy S23 Ultra 512GB', true),
    ('listing', 173, 'en', 'description', 'Samsung Galaxy S23 Ultra with S Pen stylus. 512GB storage, Snapdragon 8 Gen 2, 200MP camera. AMOLED 2X display 120Hz. Warranty included.', true),
    ('listing', 173, 'sr', 'title', 'Samsung Galaxy S23 Ultra 512GB', true),
    ('listing', 173, 'sr', 'description', 'Samsung Galaxy S23 Ultra sa S Pen olovkom. 512GB memorije, Snapdragon 8 Gen 2, 200MP kamera. AMOLED 2X ekran 120Hz. Garancija.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Samsung QLED TV 65" 4K Smart TV (ID: 174)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 174, 'ru', 'title', 'Samsung QLED TV 65" 4K Smart TV', true),
    ('listing', 174, 'ru', 'description', 'Samsung QLED телевизор 65 дюймов с 4K разрешением. Quantum Processor 4K, Object Tracking Sound+, Gaming Hub. Smart TV с Tizen OS.', true),
    ('listing', 174, 'en', 'title', 'Samsung QLED TV 65" 4K Smart TV', true),
    ('listing', 174, 'en', 'description', 'Samsung QLED TV 65 inches with 4K resolution. Quantum Processor 4K, Object Tracking Sound+, Gaming Hub. Smart TV with Tizen OS.', true),
    ('listing', 174, 'sr', 'title', 'Samsung QLED TV 65" 4K Smart TV', true),
    ('listing', 174, 'sr', 'description', 'Samsung QLED televizor 65 inča sa 4K rezolucijom. Quantum Processor 4K, Object Tracking Sound+, Gaming Hub. Smart TV sa Tizen OS.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Hugo Boss odelo - tamno plavo (ID: 175)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 175, 'ru', 'title', 'Костюм Hugo Boss - темно-синий', true),
    ('listing', 175, 'ru', 'description', 'Элегантный костюм Hugo Boss, темно-синий цвет. Размер 52. 100% шерсть. Slim fit крой. Состояние отличное, носился 2 раза.', true),
    ('listing', 175, 'en', 'title', 'Hugo Boss suit - dark blue', true),
    ('listing', 175, 'en', 'description', 'Elegant Hugo Boss suit, dark blue color. Size 52. 100% wool. Slim fit cut. Excellent condition, worn twice.', true),
    ('listing', 175, 'sr', 'title', 'Hugo Boss odelo - tamno plavo', true),
    ('listing', 175, 'sr', 'description', 'Elegantno Hugo Boss odelo, tamno plava boja. Veličina 52. 100% vuna. Slim fit kroj. Odlično stanje, nošeno 2 puta.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Nike Air Max 90 - original patike (ID: 176)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 176, 'ru', 'title', 'Nike Air Max 90 - оригинальные кроссовки', true),
    ('listing', 176, 'ru', 'description', 'Оригинальные кроссовки Nike Air Max 90, размер 43. Куплены в официальном магазине Nike. Классическая модель, белый/черный цвет. С коробкой.', true),
    ('listing', 176, 'en', 'title', 'Nike Air Max 90 - original sneakers', true),
    ('listing', 176, 'en', 'description', 'Original Nike Air Max 90 sneakers, size 43. Purchased at official Nike store. Classic model, white/black color. With box.', true),
    ('listing', 176, 'sr', 'title', 'Nike Air Max 90 - original patike', true),
    ('listing', 176, 'sr', 'description', 'Original Nike Air Max 90 patike, veličina 43. Kupljene u zvaničnoj Nike prodavnici. Klasičan model, bela/crna boja. Sa kutijom.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Michael Kors torba original (ID: 177)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 177, 'ru', 'title', 'Сумка Michael Kors оригинал', true),
    ('listing', 177, 'ru', 'description', 'Оригинальная сумка Michael Kors, модель Jet Set. Натуральная кожа, черный цвет. Вместительная, с внутренними карманами. Состояние идеальное.', true),
    ('listing', 177, 'en', 'title', 'Michael Kors bag original', true),
    ('listing', 177, 'en', 'description', 'Original Michael Kors bag, Jet Set model. Genuine leather, black color. Spacious with interior pockets. Perfect condition.', true),
    ('listing', 177, 'sr', 'title', 'Michael Kors torba original', true),
    ('listing', 177, 'sr', 'description', 'Originalna Michael Kors torba, model Jet Set. Prirodna koža, crna boja. Prostrana sa unutrašnjim džepovima. Idealno stanje.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Bosch kosilica za travu + trimer (ID: 178)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 178, 'ru', 'title', 'Газонокосилка Bosch + триммер', true),
    ('listing', 178, 'ru', 'description', 'Электрическая газонокосилка Bosch ARM 37 и триммер ART 26. Мощность 1400W и 450W. Ширина кошения 37см. В отличном состоянии.', true),
    ('listing', 178, 'en', 'title', 'Bosch lawn mower + trimmer', true),
    ('listing', 178, 'en', 'description', 'Bosch electric lawn mower ARM 37 and trimmer ART 26. Power 1400W and 450W. Cutting width 37cm. Excellent condition.', true),
    ('listing', 178, 'sr', 'title', 'Bosch kosilica za travu + trimer', true),
    ('listing', 178, 'sr', 'description', 'Bosch električna kosilica ARM 37 i trimer ART 26. Snaga 1400W i 450W. Širina košenja 37cm. Odlično stanje.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Persijski tepih 3x4m (ID: 179)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 179, 'ru', 'title', 'Персидский ковер 3x4м', true),
    ('listing', 179, 'ru', 'description', 'Прекрасный персидский ковер ручной работы. Размеры 3x4 метра. Натуральная шерсть, традиционный узор. Возраст около 20 лет, состояние отличное.', true),
    ('listing', 179, 'en', 'title', 'Persian carpet 3x4m', true),
    ('listing', 179, 'en', 'description', 'Beautiful hand-woven Persian carpet. Dimensions 3x4 meters. Natural wool, traditional pattern. About 20 years old, excellent condition.', true),
    ('listing', 179, 'sr', 'title', 'Persijski tepih 3x4m', true),
    ('listing', 179, 'sr', 'description', 'Prekrasan ručno tkan persijski tepih. Dimenzije 3x4 metra. Prirodna vuna, tradicionalni dezen. Star oko 20 godina, odlično stanje.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Traktor IMT 539 sa prikljuccima (ID: 180)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 180, 'ru', 'title', 'Трактор IMT 539 с навесным оборудованием', true),
    ('listing', 180, 'ru', 'description', 'Трактор IMT 539 в отличном состоянии. Регулярное обслуживание. В комплекте плуг, культиватор, прицеп. Мощность 39 л.с. Год выпуска 1989.', true),
    ('listing', 180, 'en', 'title', 'Tractor IMT 539 with attachments', true),
    ('listing', 180, 'en', 'description', 'IMT 539 tractor in excellent condition. Regular maintenance. Includes plow, cultivator, trailer. Power 39 HP. Year 1989.', true),
    ('listing', 180, 'sr', 'title', 'Traktor IMT 539 sa priključcima', true),
    ('listing', 180, 'sr', 'description', 'IMT 539 traktor u odličnom stanju. Redovno održavan. U kompletu plug, kultivator, prikolica. Snaga 39 KS. Godište 1989.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Seme kukuruza - hibrid NS (ID: 181)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('listing', 181, 'ru', 'title', 'Семена кукурузы - гибрид NS', true),
    ('listing', 181, 'ru', 'description', 'Качественные семена кукурузы, гибрид Новосадского института. FAO 600. Высокая урожайность, устойчивость к засухе. Упаковка 25кг.', true),
    ('listing', 181, 'en', 'title', 'Corn seeds - NS hybrid', true),
    ('listing', 181, 'en', 'description', 'Quality corn seeds, Novi Sad Institute hybrid. FAO 600. High yield, drought resistant. Package 25kg.', true),
    ('listing', 181, 'sr', 'title', 'Seme kukuruza - hibrid NS', true),
    ('listing', 181, 'sr', 'description', 'Kvalitetno seme kukuruza, hibrid Novosadskog instituta. FAO 600. Visok prinos, otporno na sušu. Pakovanje 25kg.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;


-- Добавляем переводы для витрин (user_storefronts)

-- agenstvonedvizimosty (ID: 1)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront', 1, 'ru', 'name', 'Агентство недвижимости', true),
    ('storefront', 1, 'ru', 'description', 'Мы размещаем квартиры по карте и помогаем с их продажей', true),
    ('storefront', 1, 'en', 'name', 'Real Estate Agency', true),
    ('storefront', 1, 'en', 'description', 'We place apartments on the map and help with their sale', true),
    ('storefront', 1, 'sr', 'name', 'Agencija za nekretnine', true),
    ('storefront', 1, 'sr', 'description', 'Postavljamo stanove na mapu i pomažemo u njihovoj prodaji', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Test storefront (ID: 18)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront', 18, 'ru', 'name', 'Тестовый магазин', true),
    ('storefront', 18, 'ru', 'description', 'Тестовое описание магазина', true),
    ('storefront', 18, 'en', 'name', 'Test Store', true),
    ('storefront', 18, 'en', 'description', 'Test store description', true),
    ('storefront', 18, 'sr', 'name', 'Test prodavnica', true),
    ('storefront', 18, 'sr', 'description', 'Test opis prodavnice', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Demo storefront (ID: 19)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront', 19, 'ru', 'name', 'Демо магазин', true),
    ('storefront', 19, 'ru', 'description', 'Демонстрационное описание магазина', true),
    ('storefront', 19, 'en', 'name', 'Demo Store', true),
    ('storefront', 19, 'en', 'description', 'Demo store description', true),
    ('storefront', 19, 'sr', 'name', 'Demo prodavnica', true),
    ('storefront', 19, 'sr', 'description', 'Demo opis prodavnice', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Example storefront (ID: 20)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront', 20, 'ru', 'name', 'Примерный магазин', true),
    ('storefront', 20, 'ru', 'description', 'Пример описания магазина для тестирования', true),
    ('storefront', 20, 'en', 'name', 'Example Store', true),
    ('storefront', 20, 'en', 'description', 'Example store description for testing', true),
    ('storefront', 20, 'sr', 'name', 'Primer prodavnica', true),
    ('storefront', 20, 'sr', 'description', 'Primer opisa prodavnice za testiranje', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;


-- Добавляем переводы для товаров (storefront_products)

-- iPhone 15 Pro Max 256GB (ID: 215)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 215, 'ru', 'name', 'iPhone 15 Pro Max 256GB', true),
    ('storefront_product', 215, 'ru', 'description', 'Новейший iPhone с титановым корпусом и чипом A17 Pro. 256GB памяти, тройная камера с оптическим зумом 5x.', true),
    ('storefront_product', 215, 'en', 'name', 'iPhone 15 Pro Max 256GB', true),
    ('storefront_product', 215, 'en', 'description', 'Latest iPhone with titanium frame and A17 Pro chip. 256GB storage, triple camera with 5x optical zoom.', true),
    ('storefront_product', 215, 'sr', 'name', 'iPhone 15 Pro Max 256GB', true),
    ('storefront_product', 215, 'sr', 'description', 'Najnoviji iPhone sa titanijumskim kućištem i A17 Pro čipom. 256GB memorije, trostruka kamera sa 5x optičkim zumom.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Samsung Galaxy S24 Ultra (ID: 216)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 216, 'ru', 'name', 'Samsung Galaxy S24 Ultra', true),
    ('storefront_product', 216, 'ru', 'description', 'Флагманский Samsung с S Pen, камерой 200MP и AI функциями. Snapdragon 8 Gen 3, 12GB RAM, 512GB памяти.', true),
    ('storefront_product', 216, 'en', 'name', 'Samsung Galaxy S24 Ultra', true),
    ('storefront_product', 216, 'en', 'description', 'Flagship Samsung phone with S Pen, 200MP camera and AI features. Snapdragon 8 Gen 3, 12GB RAM, 512GB storage.', true),
    ('storefront_product', 216, 'sr', 'name', 'Samsung Galaxy S24 Ultra', true),
    ('storefront_product', 216, 'sr', 'description', 'Flagship Samsung telefon sa S Pen, 200MP kamerom i AI funkcijama. Snapdragon 8 Gen 3, 12GB RAM, 512GB memorije.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Google Pixel 8 Pro (ID: 217)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 217, 'ru', 'name', 'Google Pixel 8 Pro', true),
    ('storefront_product', 217, 'ru', 'description', 'Google телефон с лучшей камерой и чистым Android. Tensor G3, 12GB RAM, 256GB. 7 лет обновлений.', true),
    ('storefront_product', 217, 'en', 'name', 'Google Pixel 8 Pro', true),
    ('storefront_product', 217, 'en', 'description', 'Google phone with best camera and pure Android. Tensor G3, 12GB RAM, 256GB. 7 years of updates.', true),
    ('storefront_product', 217, 'sr', 'name', 'Google Pixel 8 Pro', true),
    ('storefront_product', 217, 'sr', 'description', 'Google telefon sa najboljom kamerom i čistim Android-om. Tensor G3, 12GB RAM, 256GB. 7 godina ažuriranja.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- MacBook Pro 14" M3 (ID: 218)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 218, 'ru', 'name', 'MacBook Pro 14" M3', true),
    ('storefront_product', 218, 'ru', 'description', 'Apple ноутбук с чипом M3, 16GB RAM, 512GB SSD. Идеально для профессионалов. Liquid Retina XDR дисплей.', true),
    ('storefront_product', 218, 'en', 'name', 'MacBook Pro 14" M3', true),
    ('storefront_product', 218, 'en', 'description', 'Apple laptop with M3 chip, 16GB RAM, 512GB SSD. Ideal for professionals. Liquid Retina XDR display.', true),
    ('storefront_product', 218, 'sr', 'name', 'MacBook Pro 14" M3', true),
    ('storefront_product', 218, 'sr', 'description', 'Apple laptop sa M3 čipom, 16GB RAM, 512GB SSD. Idealno za profesionalce. Liquid Retina XDR ekran.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Dell XPS 15 (ID: 219)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 219, 'ru', 'name', 'Dell XPS 15', true),
    ('storefront_product', 219, 'ru', 'description', 'Премиум Windows ноутбук с 4K OLED экраном и Intel Core i9. 32GB RAM, 1TB SSD, NVIDIA RTX 4060.', true),
    ('storefront_product', 219, 'en', 'name', 'Dell XPS 15', true),
    ('storefront_product', 219, 'en', 'description', 'Premium Windows laptop with 4K OLED screen and Intel Core i9. 32GB RAM, 1TB SSD, NVIDIA RTX 4060.', true),
    ('storefront_product', 219, 'sr', 'name', 'Dell XPS 15', true),
    ('storefront_product', 219, 'sr', 'description', 'Premium Windows laptop sa 4K OLED ekranom i Intel Core i9. 32GB RAM, 1TB SSD, NVIDIA RTX 4060.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- LG OLED 65" C3 (ID: 220)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 220, 'ru', 'name', 'LG OLED 65" C3', true),
    ('storefront_product', 220, 'ru', 'description', 'OLED TV с 4K разрешением, 120Hz, HDMI 2.1 для игр. Dolby Vision IQ, Dolby Atmos. WebOS Smart TV.', true),
    ('storefront_product', 220, 'en', 'name', 'LG OLED 65" C3', true),
    ('storefront_product', 220, 'en', 'description', 'OLED TV with 4K resolution, 120Hz, HDMI 2.1 for gaming. Dolby Vision IQ, Dolby Atmos. WebOS Smart TV.', true),
    ('storefront_product', 220, 'sr', 'name', 'LG OLED 65" C3', true),
    ('storefront_product', 220, 'sr', 'description', 'OLED TV sa 4K rezolucijom, 120Hz, HDMI 2.1 za gaming. Dolby Vision IQ, Dolby Atmos. WebOS Smart TV.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Sony Bravia XR 55" (ID: 221)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 221, 'ru', 'name', 'Sony Bravia XR 55"', true),
    ('storefront_product', 221, 'ru', 'description', 'Sony TV с Cognitive Processor XR и Google TV. Dolby Vision, HDMI 2.1, Perfect for PlayStation 5.', true),
    ('storefront_product', 221, 'en', 'name', 'Sony Bravia XR 55"', true),
    ('storefront_product', 221, 'en', 'description', 'Sony TV with Cognitive Processor XR and Google TV. Dolby Vision, HDMI 2.1, Perfect for PlayStation 5.', true),
    ('storefront_product', 221, 'sr', 'name', 'Sony Bravia XR 55"', true),
    ('storefront_product', 221, 'sr', 'description', 'Sony TV sa Cognitive Processor XR i Google TV. Dolby Vision, HDMI 2.1, savršen za PlayStation 5.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Tommy Hilfiger Polo majica (ID: 222)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 222, 'ru', 'name', 'Поло Tommy Hilfiger', true),
    ('storefront_product', 222, 'ru', 'description', 'Оригинальная поло Tommy Hilfiger, 100% хлопок. Классический крой, размер L. Синий цвет с фирменным логотипом.', true),
    ('storefront_product', 222, 'en', 'name', 'Tommy Hilfiger Polo shirt', true),
    ('storefront_product', 222, 'en', 'description', 'Original Tommy Hilfiger polo shirt, 100% cotton. Classic fit, size L. Navy blue with signature logo.', true),
    ('storefront_product', 222, 'sr', 'name', 'Tommy Hilfiger Polo majica', true),
    ('storefront_product', 222, 'sr', 'description', 'Originalna Tommy Hilfiger polo majica, 100% pamuk. Klasičan kroj, veličina L. Teget boja sa prepoznatljivim logom.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Calvin Klein džins (ID: 223)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 223, 'ru', 'name', 'Джинсы Calvin Klein', true),
    ('storefront_product', 223, 'ru', 'description', 'CK Jeans модель slim fit. Премиум деним, удобные и современные. Размер 32/32. Темно-синий цвет.', true),
    ('storefront_product', 223, 'en', 'name', 'Calvin Klein jeans', true),
    ('storefront_product', 223, 'en', 'description', 'CK Jeans slim fit model. Premium denim, comfortable and modern. Size 32/32. Dark blue color.', true),
    ('storefront_product', 223, 'sr', 'name', 'Calvin Klein džins', true),
    ('storefront_product', 223, 'sr', 'description', 'CK Jeans slim fit model. Premium denim, udoban i moderan. Veličina 32/32. Tamno plava boja.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;

-- Zara haljina (ID: 224)
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES 
    ('storefront_product', 224, 'ru', 'name', 'Платье Zara', true),
    ('storefront_product', 224, 'ru', 'description', 'Элегантное вечернее платье, черный цвет. Идеально для особых случаев. Размер M. Материал: полиэстер с добавлением эластана.', true),
    ('storefront_product', 224, 'en', 'name', 'Zara dress', true),
    ('storefront_product', 224, 'en', 'description', 'Elegant evening dress, black color. Perfect for special occasions. Size M. Material: polyester with elastane.', true),
    ('storefront_product', 224, 'sr', 'name', 'Zara haljina', true),
    ('storefront_product', 224, 'sr', 'description', 'Elegantna večernja haljina, crna boja. Idealna za posebne prilike. Veličina M. Materijal: poliester sa elastinom.', true)
ON CONFLICT (entity_type, entity_id, language, field_name) DO NOTHING;