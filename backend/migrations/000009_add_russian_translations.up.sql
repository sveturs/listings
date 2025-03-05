-- Русские переводы категорий для таблицы translations
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES
-- Основные категории
('category', 1, 'ru', 'name', 'Недвижимость', true, true, NOW(), NOW()),
('category', 2, 'ru', 'name', 'Автомобили', true, true, NOW(), NOW()),
('category', 3, 'ru', 'name', 'Электроника', true, true, NOW(), NOW()),
('category', 4, 'ru', 'name', 'Все для дома и квартиры', true, true, NOW(), NOW()),
('category', 5, 'ru', 'name', 'Все для сада', true, true, NOW(), NOW()),
('category', 6, 'ru', 'name', 'Хобби и отдых', true, true, NOW(), NOW()),
('category', 7, 'ru', 'name', 'Животные', true, true, NOW(), NOW()),
('category', 8, 'ru', 'name', 'Готовый бизнес и оборудование', true, true, NOW(), NOW()),
('category', 9, 'ru', 'name', 'Прочее', true, true, NOW(), NOW()),
('category', 10, 'ru', 'name', 'Работа', true, true, NOW(), NOW()),
('category', 11, 'ru', 'name', 'Одежда, обувь, аксессуары', true, true, NOW(), NOW()),
('category', 12, 'ru', 'name', 'Детская одежда и обувь, аксессуары', true, true, NOW(), NOW()),

-- Подкатегории недвижимость
('category', 13, 'ru', 'name', 'Квартира', true, true, NOW(), NOW()),
('category', 14, 'ru', 'name', 'Комната', true, true, NOW(), NOW()),
('category', 15, 'ru', 'name', 'Дом, дача, коттедж', true, true, NOW(), NOW()),
('category', 16, 'ru', 'name', 'Земельный участок', true, true, NOW(), NOW()),
('category', 17, 'ru', 'name', 'Гараж и машиноместо', true, true, NOW(), NOW()),
('category', 18, 'ru', 'name', 'Коммерческая недвижимость', true, true, NOW(), NOW()),
('category', 19, 'ru', 'name', 'Недвижимость за рубежом', true, true, NOW(), NOW()),
('category', 20, 'ru', 'name', 'Отель', 'hotel-property', true, true, NOW(), NOW()),
('category', 21, 'ru', 'name', 'Апартаменты', true, true, NOW(), NOW()),

-- Подкатегории автотранспорта
('category', 22, 'ru', 'name', 'Легковые автомобили', true, true, NOW(), NOW()),
('category', 23, 'ru', 'name', 'Грузовые автомобили', true, true, NOW(), NOW()),
('category', 24, 'ru', 'name', 'Спецтехника', true, true, NOW(), NOW()),
('category', 25, 'ru', 'name', 'Сельхозтехника', true, true, NOW(), NOW()),
('category', 26, 'ru', 'name', 'Аренда авто и спецтехники', true, true, NOW(), NOW()),
('category', 27, 'ru', 'name', 'Мотоциклы и мототехника', true, true, NOW(), NOW()),
('category', 28, 'ru', 'name', 'Водный транспорт', true, true, NOW(), NOW()),
('category', 29, 'ru', 'name', 'Запчасти и аксессуары', true, true, NOW(), NOW()),

-- Грузовые автомобили
('category', 30, 'ru', 'name', 'Грузовики', true, true, NOW(), NOW()),
('category', 31, 'ru', 'name', 'Полуприцепы', true, true, NOW(), NOW()),
('category', 32, 'ru', 'name', 'Лёгкий коммерческий транспорт', true, true, NOW(), NOW()),
('category', 33, 'ru', 'name', 'Автобусы', true, true, NOW(), NOW()),

-- Спецтехника
('category', 34, 'ru', 'name', 'Экскаваторы', true, true, NOW(), NOW()),
('category', 35, 'ru', 'name', 'Погрузчики', true, true, NOW(), NOW()),
('category', 36, 'ru', 'name', 'Экскаваторы-погрузчики', true, true, NOW(), NOW()),
('category', 37, 'ru', 'name', 'Автокраны', true, true, NOW(), NOW()),
('category', 38, 'ru', 'name', 'Автобетоносмесители', true, true, NOW(), NOW()),
('category', 39, 'ru', 'name', 'Дорожные катки', true, true, NOW(), NOW()),
('category', 40, 'ru', 'name', 'Поливомоечные машины', true, true, NOW(), NOW()),
('category', 41, 'ru', 'name', 'Мусоровозы', true, true, NOW(), NOW()),
('category', 42, 'ru', 'name', 'Автовышки', true, true, NOW(), NOW()),
('category', 43, 'ru', 'name', 'Бульдозеры', true, true, NOW(), NOW()),
('category', 44, 'ru', 'name', 'Автогрейдеры', true, true, NOW(), NOW()),
('category', 45, 'ru', 'name', 'Буровые установки', true, true, NOW(), NOW()),

-- Сельхозтехника
('category', 46, 'ru', 'name', 'Тракторы', true, true, NOW(), NOW()),
('category', 47, 'ru', 'name', 'Мини-тракторы', true, true, NOW(), NOW()),
('category', 48, 'ru', 'name', 'Пресс-подборщики', true, true, NOW(), NOW()),
('category', 49, 'ru', 'name', 'Бороны', true, true, NOW(), NOW()),
('category', 50, 'ru', 'name', 'Косилки', true, true, NOW(), NOW()),

-- Животные
('category', 200, 'ru', 'name', 'Собаки', true, true, NOW(), NOW()),
('category', 201, 'ru', 'name', 'Кошки', true, true, NOW(), NOW()),
('category', 202, 'ru', 'name', 'Птицы', true, true, NOW(), NOW()),
('category', 203, 'ru', 'name', 'Аквариум', true, true, NOW(), NOW()),
('category', 204, 'ru', 'name', 'Другие животные', true, true, NOW(), NOW()),

-- Аренда авто и спецтехники
('category', 205, 'ru', 'name', 'Авто', true, true, NOW(), NOW()),
('category', 206, 'ru', 'name', 'Подъёмная техника', true, true, NOW(), NOW()),
('category', 207, 'ru', 'name', 'Землеройная техника', true, true, NOW(), NOW()),
('category', 208, 'ru', 'name', 'Коммунальная техника', true, true, NOW(), NOW()),
('category', 209, 'ru', 'name', 'Дорожно-строительная техника', true, true, NOW(), NOW()),
('category', 210, 'ru', 'name', 'Грузовой транспорт', true, true, NOW(), NOW()),
('category', 211, 'ru', 'name', 'Погрузочная техника', true, true, NOW(), NOW()),
('category', 212, 'ru', 'name', 'Навесное оборудование', true, true, NOW(), NOW()),
('category', 213, 'ru', 'name', 'Прицепы', true, true, NOW(), NOW()),
('category', 214, 'ru', 'name', 'Сельхозтехника', true, true, NOW(), NOW()),
('category', 215, 'ru', 'name', 'Автодома', true, true, NOW(), NOW()),

-- Мотоциклы и мототехника
('category', 216, 'ru', 'name', 'Вездеходы', true, true, NOW(), NOW()),
('category', 217, 'ru', 'name', 'Картинг', true, true, NOW(), NOW()),
('category', 218, 'ru', 'name', 'Квадроциклы и багги', true, true, NOW(), NOW()),
('category', 219, 'ru', 'name', 'Мопеды', true, true, NOW(), NOW()),
('category', 220, 'ru', 'name', 'Скутеры', true, true, NOW(), NOW()),
('category', 221, 'ru', 'name', 'Мотоциклы', true, true, NOW(), NOW()),
('category', 222, 'ru', 'name', 'Снегоходы', true, true, NOW(), NOW()),

-- Водный транспорт
('category', 223, 'ru', 'name', 'Вёсельные лодки', true, true, NOW(), NOW()),
('category', 224, 'ru', 'name', 'Каяки', true, true, NOW(), NOW()),
('category', 225, 'ru', 'name', 'Гидроциклы', true, true, NOW(), NOW()),
('category', 226, 'ru', 'name', 'Катера и яхты', true, true, NOW(), NOW()),
('category', 227, 'ru', 'name', 'Моторные лодки и моторы', true, true, NOW(), NOW()),

-- Запчасти и аксессуары
('category', 228, 'ru', 'name', 'Запчасти', true, true, NOW(), NOW()),
('category', 229, 'ru', 'name', 'Шины, диски и колёса', true, true, NOW(), NOW()),
('category', 230, 'ru', 'name', 'Аудио- и видеотехника', true, true, NOW(), NOW()),
('category', 231, 'ru', 'name', 'Аксессуары', true, true, NOW(), NOW()),
('category', 232, 'ru', 'name', 'Масла и автохимия', true, true, NOW(), NOW()),
('category', 233, 'ru', 'name', 'Инструменты', true, true, NOW(), NOW()),
('category', 234, 'ru', 'name', 'Багажники и фаркопы', true, true, NOW(), NOW()),
('category', 235, 'ru', 'name', 'Прицепы', true, true, NOW(), NOW()),
('category', 236, 'ru', 'name', 'Экипировка', true, true, NOW(), NOW()),
('category', 237, 'ru', 'name', 'Противоугонные устройства', true, true, NOW(), NOW()),
('category', 238, 'ru', 'name', 'GPS-навигаторы', true, true, NOW(), NOW()),

-- Электроника
('category', 239, 'ru', 'name', 'Телефоны', true, true, NOW(), NOW()),
('category', 240, 'ru', 'name', 'Аудио и видео', true, true, NOW(), NOW()),
('category', 241, 'ru', 'name', 'Товары для компьютера', true, true, NOW(), NOW()),
('category', 242, 'ru', 'name', 'Игры, приставки и программы', true, true, NOW(), NOW()),
('category', 243, 'ru', 'name', 'Ноутбуки', true, true, NOW(), NOW()),
('category', 244, 'ru', 'name', 'Фототехника', true, true, NOW(), NOW()),
('category', 245, 'ru', 'name', 'Планшеты и электронные книги', true, true, NOW(), NOW()),
('category', 246, 'ru', 'name', 'Оргтехника и расходники', true, true, NOW(), NOW()),
('category', 247, 'ru', 'name', 'Бытовая техника', true, true, NOW(), NOW()),

-- Подкатегории Телефонов
('category', 248, 'ru', 'name', 'Мобильные телефоны', true, true, NOW(), NOW()),
('category', 249, 'ru', 'name', 'Аксессуары', true, true, NOW(), NOW()),
('category', 250, 'ru', 'name', 'Рации', true, true, NOW(), NOW()),
('category', 251, 'ru', 'name', 'Стационарные телефоны', true, true, NOW(), NOW()),

-- Аксессуары для телефонов
('category', 252, 'ru', 'name', 'Аккумуляторы', true, true, NOW(), NOW()),
('category', 253, 'ru', 'name', 'Гарнитуры и наушники', true, true, NOW(), NOW()),
('category', 254, 'ru', 'name', 'Зарядные устройства', true, true, NOW(), NOW()),
('category', 255, 'ru', 'name', 'Кабели и адаптеры', true, true, NOW(), NOW()),
('category', 256, 'ru', 'name', 'Модемы и роутеры', true, true, NOW(), NOW()),
('category', 257, 'ru', 'name', 'Чехлы и плёнки', true, true, NOW(), NOW()),
('category', 258, 'ru', 'name', 'Запчасти', true, true, NOW(), NOW()),

-- Аудио и видео
('category', 259, 'ru', 'name', 'Телевизоры и проекторы', true, true, NOW(), NOW()),
('category', 260, 'ru', 'name', 'Наушники', true, true, NOW(), NOW()),
('category', 261, 'ru', 'name', 'Акустика, колонки, сабвуферы', true, true, NOW(), NOW()),
('category', 262, 'ru', 'name', 'Аксессуары', true, true, NOW(), NOW()),
('category', 263, 'ru', 'name', 'Музыкальные центры, магнитолы', true, true, NOW(), NOW()),
('category', 264, 'ru', 'name', 'Усилители и ресиверы', true, true, NOW(), NOW()),
('category', 265, 'ru', 'name', 'Видеокамеры', true, true, NOW(), NOW()),
('category', 266, 'ru', 'name', 'Видео, DVD и Blu-ray плееры', true, true, NOW(), NOW()),
('category', 267, 'ru', 'name', 'Кабели и адаптеры', true, true, NOW(), NOW()),
('category', 268, 'ru', 'name', 'Музыка и фильмы', true, true, NOW(), NOW()),
('category', 269, 'ru', 'name', 'Микрофоны', true, true, NOW(), NOW()),
('category', 270, 'ru', 'name', 'MP3-плееры', true, true, NOW(), NOW()),

-- Товары для компьютера
('category', 271, 'ru', 'name', 'Системные блоки', true, true, NOW(), NOW()),
('category', 272, 'ru', 'name', 'Моноблоки', true, true, NOW(), NOW()),
('category', 273, 'ru', 'name', 'Комплектующие', true, true, NOW(), NOW()),
('category', 274, 'ru', 'name', 'Мониторы и запчасти', true, true, NOW(), NOW()),
('category', 275, 'ru', 'name', 'Сетевое оборудование', true, true, NOW(), NOW()),
('category', 276, 'ru', 'name', 'Клавиатуры и мыши', true, true, NOW(), NOW()),
('category', 277, 'ru', 'name', 'Джойстики и рули', true, true, NOW(), NOW()),
('category', 278, 'ru', 'name', 'Флэшки и карты памяти', true, true, NOW(), NOW()),
('category', 279, 'ru', 'name', 'Веб-камеры', true, true, NOW(), NOW()),
('category', 280, 'ru', 'name', 'ТВ-тюнеры', true, true, NOW(), NOW()),

-- Комплектующие для компьютера
('category', 281, 'ru', 'name', 'CD, DVD и Blu-ray приводы', true, true, NOW(), NOW()),
('category', 282, 'ru', 'name', 'Блоки питания', true, true, NOW(), NOW()),
('category', 283, 'ru', 'name', 'Видеокарты', true, true, NOW(), NOW()),
('category', 284, 'ru', 'name', 'Жёсткие диски', true, true, NOW(), NOW()),
('category', 285, 'ru', 'name', 'Звуковые карты', true, true, NOW(), NOW()),
('category', 286, 'ru', 'name', 'Контроллеры', true, true, NOW(), NOW()),
('category', 287, 'ru', 'name', 'Корпуса', true, true, NOW(), NOW()),
('category', 288, 'ru', 'name', 'Материнские платы', true, true, NOW(), NOW()),
('category', 289, 'ru', 'name', 'Оперативная память', true, true, NOW(), NOW()),
('category', 290, 'ru', 'name', 'Процессоры', true, true, NOW(), NOW()),
('category', 291, 'ru', 'name', 'Системы охлаждения', true, true, NOW(), NOW()),

-- Игры, приставки и программы
('category', 292, 'ru', 'name', 'Игровые приставки', true, true, NOW(), NOW()),
('category', 293, 'ru', 'name', 'Игры для приставок', true, true, NOW(), NOW()),
('category', 294, 'ru', 'name', 'Компьютерные игры', true, true, NOW(), NOW()),
('category', 295, 'ru', 'name', 'Программы', true, true, NOW(), NOW()),

-- Фототехника
('category', 296, 'ru', 'name', 'Оборудование и аксессуары', true, true, NOW(), NOW()),
('category', 297, 'ru', 'name', 'Объективы', true, true, NOW(), NOW()),
('category', 298, 'ru', 'name', 'Компактные фотоаппараты', true, true, NOW(), NOW()),
('category', 299, 'ru', 'name', 'Плёночные фотоаппараты', true, true, NOW(), NOW()),
('category', 300, 'ru', 'name', 'Зеркальные фотоаппараты', true, true, NOW(), NOW()),

-- Планшеты и электронные книги
('category', 301, 'ru', 'name', 'Планшеты', true, true, NOW(), NOW()),
('category', 302, 'ru', 'name', 'Электронные книги', true, true, NOW(), NOW()),
('category', 303, 'ru', 'name', 'Аксессуары', true, true, NOW(), NOW()),

-- Оргтехника и расходники
('category', 304, 'ru', 'name', 'МФУ, копиры и сканеры', true, true, NOW(), NOW()),
('category', 305, 'ru', 'name', 'Принтеры', true, true, NOW(), NOW()),
('category', 306, 'ru', 'name', 'Канцелярия', true, true, NOW(), NOW()),
('category', 307, 'ru', 'name', 'ИБП, сетевые фильтры', true, true, NOW(), NOW()),
('category', 308, 'ru', 'name', 'Телефония', true, true, NOW(), NOW()),
('category', 309, 'ru', 'name', 'Уничтожители бумаг', true, true, NOW(), NOW()),
('category', 310, 'ru', 'name', 'Расходные материалы', true, true, NOW(), NOW()),

-- Бытовая техника
('category', 311, 'ru', 'name', 'Для кухни', true, true, NOW(), NOW()),
('category', 312, 'ru', 'name', 'Для дома', true, true, NOW(), NOW()),
('category', 313, 'ru', 'name', 'Климатическое оборудование', true, true, NOW(), NOW()),
('category', 314, 'ru', 'name', 'Для индивидуального ухода', true, true, NOW(), NOW()),

-- Техника для кухни
('category', 315, 'ru', 'name', 'Вытяжки', true, true, NOW(), NOW()),
('category', 316, 'ru', 'name', 'Мелкая кухонная техника', true, true, NOW(), NOW()),
('category', 317, 'ru', 'name', 'Микроволновые печи', true, true, NOW(), NOW()),
('category', 318, 'ru', 'name', 'Плиты и духовые шкафы', true, true, NOW(), NOW()),
('category', 319, 'ru', 'name', 'Посудомоечные машины', true, true, NOW(), NOW()),
('category', 320, 'ru', 'name', 'Холодильники и морозильные камеры', true, true, NOW(), NOW()),

-- Техника для дома
('category', 321, 'ru', 'name', 'Пылесосы и запчасти', true, true, NOW(), NOW()),
('category', 322, 'ru', 'name', 'Стиральные и сушильные машины', true, true, NOW(), NOW()),
('category', 323, 'ru', 'name', 'Утюги', true, true, NOW(), NOW()),
('category', 324, 'ru', 'name', 'Швейное оборудование', true, true, NOW(), NOW()),

-- Климатическое оборудование
('category', 325, 'ru', 'name', 'Вентиляторы', true, true, NOW(), NOW()),
('category', 326, 'ru', 'name', 'Кондиционеры и запчасти', true, true, NOW(), NOW()),
('category', 327, 'ru', 'name', 'Обогреватели', true, true, NOW(), NOW()),
('category', 328, 'ru', 'name', 'Очистители воздуха', true, true, NOW(), NOW()),
('category', 329, 'ru', 'name', 'Термометры и метеостанции', true, true, NOW(), NOW()),

-- Техника для индивидуального ухода
('category', 330, 'ru', 'name', 'Бритвы и триммеры', true, true, NOW(), NOW()),
('category', 331, 'ru', 'name', 'Машинки для стрижки', true, true, NOW(), NOW()),
('category', 332, 'ru', 'name', 'Фены и приборы для укладки', true, true, NOW(), NOW()),
('category', 333, 'ru', 'name', 'Эпиляторы', true, true, NOW(), NOW()),

-- Для дома и квартиры
('category', 334, 'ru', 'name', 'Ремонт и строительство', true, true, NOW(), NOW()),
('category', 335, 'ru', 'name', 'Мебель и интерьер', true, true, NOW(), NOW()),
('category', 336, 'ru', 'name', 'Продукты питания', true, true, NOW(), NOW()),
('category', 337, 'ru', 'name', 'Посуда и товары для кухни', true, true, NOW(), NOW()),

-- Основные категории и другие важные для демонстрации
('category', 338, 'ru', 'name', 'Двери', true, true, NOW(), NOW()),
('category', 339, 'ru', 'name', 'Инструменты', true, true, NOW(), NOW()),
('category', 340, 'ru', 'name', 'Камины и обогреватели', true, true, NOW(), NOW()),
('category', 341, 'ru', 'name', 'Окна и балконы', true, true, NOW(), NOW()),
('category', 342, 'ru', 'name', 'Потолки', true, true, NOW(), NOW()),
('category', 343, 'ru', 'name', 'Для сада и дачи', true, true, NOW(), NOW()),
('category', 344, 'ru', 'name', 'Сантехника, водоснабжение и сауна', true, true, NOW(), NOW()),
('category', 345, 'ru', 'name', 'Готовые строения и срубы', true, true, NOW(), NOW()),
('category', 346, 'ru', 'name', 'Ворота, заборы и ограждения', true, true, NOW(), NOW()),
('category', 347, 'ru', 'name', 'Охрана и сигнализации', true, true, NOW(), NOW()),
('category', 348, 'ru', 'name', 'Кровати, диваны и кресла', true, true, NOW(), NOW()),
('category', 349, 'ru', 'name', 'Текстиль и ковры', true, true, NOW(), NOW()),
('category', 350, 'ru', 'name', 'Освещение', true, true, NOW(), NOW()),
('category', 351, 'ru', 'name', 'Компьютерные столы и кресла', true, true, NOW(), NOW()),
('category', 352, 'ru', 'name', 'Шкафы, комоды и стеллажи', true, true, NOW(), NOW()),
('category', 353, 'ru', 'name', 'Кухонные гарнитуры', true, true, NOW(), NOW()),
('category', 354, 'ru', 'name', 'Столы и стулья', true, true, NOW(), NOW()),
('category', 355, 'ru', 'name', 'Комнатные растения', true, true, NOW(), NOW()),
('category', 356, 'ru', 'name', 'Декоративно-лиственные растения', true, true, NOW(), NOW()),
('category', 357, 'ru', 'name', 'Цветущие растения', true, true, NOW(), NOW()),
('category', 358, 'ru', 'name', 'Пальмы и фикусы', true, true, NOW(), NOW()),
('category', 359, 'ru', 'name', 'Кактусы и суккуленты', true, true, NOW(), NOW()),
('category', 360, 'ru', 'name', 'Чай, кофе, какао', true, true, NOW(), NOW()),
('category', 361, 'ru', 'name', 'Напитки', true, true, NOW(), NOW()),
('category', 362, 'ru', 'name', 'Рыба, морепродукты, икра', true, true, NOW(), NOW()),
('category', 363, 'ru', 'name', 'Мясо, птица, субпродукты', true, true, NOW(), NOW()),
('category', 364, 'ru', 'name', 'Кондитерские изделия', true, true, NOW(), NOW()),
('category', 366, 'ru', 'name', 'Посуда', true, true, NOW(), NOW()),
('category', 367, 'ru', 'name', 'Товары для кухни', true, true, NOW(), NOW()),
('category', 368, 'ru', 'name', 'Сервировка стола', true, true, NOW(), NOW()),
('category', 369, 'ru', 'name', 'Приготовление пищи', true, true, NOW(), NOW()),
('category', 370, 'ru', 'name', 'Хранение продуктов', true, true, NOW(), NOW()),
('category', 371, 'ru', 'name', 'Приготовление напитков', true, true, NOW(), NOW()),
('category', 372, 'ru', 'name', 'Хозяйственные товары', true, true, NOW(), NOW()),

-- Все для сада (родитель id=5)
('category', 373, 'ru', 'name', 'Садовая мебель', true, true, NOW(), NOW()),
('category', 374, 'ru', 'name', 'Освещение', true, true, NOW(), NOW()),
('category', 375, 'ru', 'name', 'Оформление интерьера', true, true, NOW(), NOW()),
('category', 376, 'ru', 'name', 'Растения и семена', true, true, NOW(), NOW()),
('category', 377, 'ru', 'name', 'Сад и огород', true, true, NOW(), NOW()),
('category', 378, 'ru', 'name', 'Садовые растения', true, true, NOW(), NOW()),
('category', 379, 'ru', 'name', 'Семена, луковицы, клубни', true, true, NOW(), NOW()),
('category', 380, 'ru', 'name', 'Товары для ухода за растениями', true, true, NOW(), NOW()),
('category', 381, 'ru', 'name', 'Панно и искусственные растения', true, true, NOW(), NOW()),

-- Подкатегории Садовые растения (родитель id=378)
('category', 382, 'ru', 'name', 'Декоративные кустарники и деревья', true, true, NOW(), NOW()),
('category', 383, 'ru', 'name', 'Хвойные растения', true, true, NOW(), NOW()),
('category', 384, 'ru', 'name', 'Многолетние растения', true, true, NOW(), NOW()),
('category', 385, 'ru', 'name', 'Плодовые растения', true, true, NOW(), NOW()),
('category', 386, 'ru', 'name', 'Газон', true, true, NOW(), NOW()),
('category', 387, 'ru', 'name', 'Зелень и пряные травы', true, true, NOW(), NOW()),

-- Подкатегории Товары для ухода за растениями (родитель id=380)
('category', 388, 'ru', 'name', 'Грунты и субстраты', true, true, NOW(), NOW()),
('category', 389, 'ru', 'name', 'Удобрения', true, true, NOW(), NOW()),
('category', 390, 'ru', 'name', 'Средства от вредителей и сорняков', true, true, NOW(), NOW()),
('category', 391, 'ru', 'name', 'Горшки и кашпо', true, true, NOW(), NOW()),
('category', 392, 'ru', 'name', 'Фитолампы', true, true, NOW(), NOW()),
('category', 393, 'ru', 'name', 'Измерители влаги', true, true, NOW(), NOW()),
('category', 394, 'ru', 'name', 'Теплицы, грядки, клумбы', true, true, NOW(), NOW()),

-- Прочие категории (добавьте остальные категории при необходимости)
('category', 473, 'ru', 'name', 'Ракия и вино', true, true, NOW(), NOW()),
('category', 474, 'ru', 'name', 'Домашние сыры', true, true, NOW(), NOW()),
('category', 475, 'ru', 'name', 'Каймак', true, true, NOW(), NOW()),
('category', 476, 'ru', 'name', 'Айвар', true, true, NOW(), NOW()),
('category', 477, 'ru', 'name', 'Сливовая ракия', true, true, NOW(), NOW()),
('category', 478, 'ru', 'name', 'Виноградная ракия', true, true, NOW(), NOW()),
('category', 479, 'ru', 'name', 'Фруктовая ракия', true, true, NOW(), NOW()),
('category', 480, 'ru', 'name', 'Домашнее вино', true, true, NOW(), NOW()),
('category', 481, 'ru', 'name', 'Народные ремесла', true, true, NOW(), NOW()),
('category', 482, 'ru', 'name', 'Опанци', true, true, NOW(), NOW()),
('category', 483, 'ru', 'name', 'Керамика', true, true, NOW(), NOW()),
('category', 484, 'ru', 'name', 'Вышивка', true, true, NOW(), NOW()),
('category', 485, 'ru', 'name', 'Ткачество', true, true, NOW(), NOW()),
('category', 486, 'ru', 'name', 'Народные инструменты', true, true, NOW(), NOW()),

-- Сельскохозяйственные категории
('category', 487, 'ru', 'name', 'Пчеловодство', true, true, NOW(), NOW()),

-- Подкатегории пчеловодства (родитель id=487)
('category', 488, 'ru', 'name', 'Мёд', true, true, NOW(), NOW()),
('category', 489, 'ru', 'name', 'Пчелиный воск', true, true, NOW(), NOW()),
('category', 490, 'ru', 'name', 'Прополис', true, true, NOW(), NOW()),
('category', 491, 'ru', 'name', 'Пчеловодный инвентарь', true, true, NOW(), NOW()),

-- Категории для сезонных работ
('category', 492, 'ru', 'name', 'Сезонные работы', true, true, NOW(), NOW()),
-- Подкатегории сезонных работ (родитель id=492)
('category', 493, 'ru', 'name', 'Сбор урожая', true, true, NOW(), NOW()),
('category', 494, 'ru', 'name', 'Работа на винограднике', true, true, NOW(), NOW()),
('category', 495, 'ru', 'name', 'Сезонные строительные работы', true, true, NOW(), NOW()),

-- Туристические услуги
('category', 496, 'ru', 'name', 'Сельский туризм', true, true, NOW(), NOW()),

-- Подкатегории сельского туризма (родитель id=496)

('category', 497, 'ru', 'name', 'Этно-деревни', true, true, NOW(), NOW()),
('category', 498, 'ru', 'name', 'Винные туры', true, true, NOW(), NOW()),
('category', 499, 'ru', 'name', 'Агротуризм', true, true, NOW(), NOW()),
('category', 500, 'ru', 'name', 'Горный туризм', true, true, NOW(), NOW());




-- Обновляем sequence для translations
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);