-- Русские переводы категорий для таблицы translations
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES

('category', 1, 'ru', 'name', 'Недвижимость', true, true, NOW(), NOW()),
('category', 1000, 'Недвижимость', 'real-estate', NULL, 'home', '2025-02-07 07:13:52.823283'),
    -- Недвижимость подкатегории 1000
    ('category', 1100, 'ru', 'name', 'Квартира', true, true, NOW(), NOW()),
    ('category', 1200, 'ru', 'name', 'Комната', true, true, NOW(), NOW()),
    ('category', 1300, 'ru', 'name', 'Дом, дача, коттедж', true, true, NOW(), NOW()),
    ('category', 1400, 'ru', 'name', 'Земельный участок', true, true, NOW(), NOW()),
    ('category', 1500, 'ru', 'name', 'Гараж и машиноместо', true, true, NOW(), NOW()),
    ('category', 1600, 'ru', 'name', 'Коммерческая недвижимость', true, true, NOW(), NOW()),
    ('category', 1700, 'ru', 'name', 'Недвижимость за рубежом', true, true, NOW(), NOW()),
    ('category', 1800, 'ru', 'name', 'Отель', true, true, NOW(), NOW()),
    ('category', 1900, 'ru', 'name', 'Апартаменты', true, true, NOW(), NOW()),

('category', 2000, 'ru', 'name', 'Автомобили', true, true, NOW(), NOW()),

    -- Auto categories 2000
    ('category', 2100, 'ru', 'name', 'Легковые автомобили', true, true, NOW(), NOW()),
    ('category', 2200, 'ru', 'name', 'Грузовые автомобили', true, true, NOW(), NOW()),
        -- Commercial vehicles subcategories 2200
        ('category', 2210, 'ru', 'name', 'Грузовики', true, true, NOW(), NOW()),
        ('category', 2220, 'ru', 'name', 'Полуприцепы', true, true, NOW(), NOW()),
        ('category', 2230, 'ru', 'name', 'Лёгкий коммерческий транспорт', true, true, NOW(), NOW()),
        ('category', 2240, 'ru', 'name', 'Автобусы', true, true, NOW(), NOW()),
        
    ('category', 2300, 'ru', 'name', 'Спецтехника', true, true, NOW(), NOW()),
        -- Special equipment subcategories 2300
        ('category', 2310, 'ru', 'name', 'Экскаваторы', true, true, NOW(), NOW()),
        ('category', 2315, 'ru', 'name', 'Погрузчики', true, true, NOW(), NOW()),
        ('category', 2320, 'ru', 'name', 'Экскаваторы-погрузчики', true, true, NOW(), NOW()),
        ('category', 2325, 'ru', 'name', 'Автокраны', true, true, NOW(), NOW()),
        ('category', 2330, 'ru', 'name', 'Автобетоносмесители', true, true, NOW(), NOW()),
        ('category', 2335, 'ru', 'name', 'Дорожные катки', true, true, NOW(), NOW()),
        ('category', 2340, 'ru', 'name', 'Поливомоечные машины', true, true, NOW(), NOW()),
        ('category', 2345, 'ru', 'name', 'Мусоровозы', true, true, NOW(), NOW()),
        ('category', 2350, 'ru', 'name', 'Автовышки', true, true, NOW(), NOW()),
        ('category', 2355, 'ru', 'name', 'Бульдозеры', true, true, NOW(), NOW()),
        ('category', 2360, 'ru', 'name', 'Автогрейдеры', true, true, NOW(), NOW()),
        ('category', 2365, 'ru', 'name', 'Буровые установки', true, true, NOW(), NOW()),

    ('category', 2400, 'ru', 'name', 'Сельхозтехника', true, true, NOW(), NOW()),
        -- Agricultural machinery subcategories 2400
        ('category', 2410, 'ru', 'name', 'Тракторы', true, true, NOW(), NOW()),
        ('category', 2415, 'ru', 'name', 'Мини-тракторы', true, true, NOW(), NOW()),
        ('category', 2420, 'ru', 'name', 'Пресс-подборщики', true, true, NOW(), NOW()),
        ('category', 2425, 'ru', 'name', 'Бороны', true, true, NOW(), NOW()),
        ('category', 2430, 'ru', 'name', 'Косилки', true, true, NOW(), NOW()),
        ('category', 2435, 'ru', 'name', 'Комбайны', true, true, NOW(), NOW()),
        ('category', 2440, 'ru', 'name', 'Телескопические погрузчики', true, true, NOW(), NOW()),
        ('category', 2445, 'ru', 'name', 'Сеялки', true, true, NOW(), NOW()),
        ('category', 2450, 'ru', 'name', 'Культиваторы', true, true, NOW(), NOW()),
        ('category', 2455, 'ru', 'name', 'Плуги', true, true, NOW(), NOW()),
        ('category', 2460, 'ru', 'name', 'Опрыскиватели', true, true, NOW(), NOW()),

    ('category', 2500, 'ru', 'name', 'Аренда авто и спецтехники', true, true, NOW(), NOW()),
            -- Аренда авто и спецтехники (родитель id=2500)
        ('category', 2510, 'ru', 'name', 'Авто', true, true, NOW(), NOW()),
        ('category', 2515, 'ru', 'name', 'Подъёмная техника', true, true, NOW(), NOW()),
        ('category', 2520, 'ru', 'name', 'Землеройная техника', true, true, NOW(), NOW()),
        ('category', 2525, 'ru', 'name', 'Коммунальная техника', true, true, NOW(), NOW()),
        ('category', 2530, 'ru', 'name', 'Дорожно-строительная техника', true, true, NOW(), NOW()),
        ('category', 2535, 'ru', 'name', 'Грузовой транспорт', true, true, NOW(), NOW()),
        ('category', 2540, 'ru', 'name', 'Погрузочная техника', true, true, NOW(), NOW()),
        ('category', 2545, 'ru', 'name', 'Навесное оборудование', true, true, NOW(), NOW()),
        ('category', 2550, 'ru', 'name', 'Прицепы', true, true, NOW(), NOW()),
        ('category', 2555, 'ru', 'name', 'Сельхозтехника', true, true, NOW(), NOW()),
        ('category', 2560, 'ru', 'name', 'Автодома', true, true, NOW(), NOW()),

    ('category', 2600, 'ru', 'name', 'Мотоциклы и мототехника', true, true, NOW(), NOW()),
        -- Мотоциклы и мототехника (родитель id=2600)
        ('category', 2610, 'ru', 'name', 'Вездеходы', true, true, NOW(), NOW()),
        ('category', 2615, 'ru', 'name', 'Картинг', true, true, NOW(), NOW()),
        ('category', 2620, 'ru', 'name', 'Квадроциклы и багги', true, true, NOW(), NOW()),
        ('category', 2625, 'ru', 'name', 'Мопеды', true, true, NOW(), NOW()),
        ('category', 2630, 'ru', 'name', 'Скутеры', true, true, NOW(), NOW()),
        ('category', 2635, 'ru', 'name', 'Мотоциклы', true, true, NOW(), NOW()),
        ('category', 2640, 'ru', 'name', 'Снегоходы', true, true, NOW(), NOW()),

    ('category', 2700, 'ru', 'name', 'Водный транспорт', true, true, NOW(), NOW()),
        -- Водный транспорт (родитель id=2700)
        ('category', 2710, 'ru', 'name', 'Вёсельные лодки', true, true, NOW(), NOW()),
        ('category', 2720, 'ru', 'name', 'Каяки', true, true, NOW(), NOW()),
        ('category', 2730, 'ru', 'name', 'Гидроциклы', true, true, NOW(), NOW()),
        ('category', 2740, 'ru', 'name', 'Катера и яхты', true, true, NOW(), NOW()),
        ('category', 2750, 'ru', 'name', 'Моторные лодки и моторы', true, true, NOW(), NOW()),

    ('category', 2800, 'ru', 'name', 'Запчасти и аксессуары', true, true, NOW(), NOW()),
        -- Запчасти и аксессуары (родитель id=2800)
        ('category', 2810, 'ru', 'name', 'Запчасти', true, true, NOW(), NOW()),
        ('category', 2815, 'ru', 'name', 'Шины, диски и колёса', true, true, NOW(), NOW()),
        ('category', 2820, 'ru', 'name', 'Аудио- и видеотехника', true, true, NOW(), NOW()),
        ('category', 2825, 'ru', 'name', 'Аксессуары', true, true, NOW(), NOW()),
        ('category', 2830, 'ru', 'name', 'Масла и автохимия', true, true, NOW(), NOW()),
        ('category', 2835, 'ru', 'name', 'Инструменты', true, true, NOW(), NOW()),
        ('category', 2840, 'ru', 'name', 'Багажники и фаркопы', true, true, NOW(), NOW()),
        ('category', 2845, 'ru', 'name', 'Прицепы', true, true, NOW(), NOW()),
        ('category', 2850, 'ru', 'name', 'Экипировка', true, true, NOW(), NOW()),
        ('category', 2855, 'ru', 'name', 'Противоугонные устройства', true, true, NOW(), NOW()),
        ('category', 2860, 'ru', 'name', 'GPS-навигаторы', true, true, NOW(), NOW()),

('category', 3000, 'ru', 'name', 'Электроника', true, true, NOW(), NOW()),
    -- Электроника (родитель id=3)
    ('category', 3100, 'ru', 'name', 'Телефоны', true, true, NOW(), NOW()),
            -- Подкатегории Телефонов (родитель id=3100)
        ('category', 3110, 'ru', 'name', 'Мобильные телефоны', true, true, NOW(), NOW()),
        ('category', 3120, 'ru', 'name', 'Аксессуары', true, true, NOW(), NOW()),
                    -- Подкатегории аксессуаров для телефонов (родитель id=3120)
            ('category', 3121, 'ru', 'name', 'ru', 'name', 'Аккумуляторы', 'batteries', 3120, 'battery', '2025-02-07 07:13:52.823283'),
            ('category', 3122, 'ru', 'name', 'Гарнитуры и наушники', true, true, NOW(), NOW()),
            ('category', 3123, 'ru', 'name', 'Зарядные устройства', true, true, NOW(), NOW()),
            ('category', 3124, 'ru', 'name', 'Кабели и адаптеры', true, true, NOW(), NOW()),
            ('category', 3125, 'ru', 'name', 'Модемы и роутеры', true, true, NOW(), NOW()),
            ('category', 3126, 'ru', 'name', 'Чехлы и плёнки', true, true, NOW(), NOW()),
            ('category', 3127, 'ru', 'name', 'Запчасти', true, true, NOW(), NOW()),

        ('category', 3130, 'ru', 'name', 'Рации', true, true, NOW(), NOW()),
        ('category', 3140, 'ru', 'name', 'Стационарные телефоны', true, true, NOW(), NOW()),

    ('category', 3200, 'ru', 'name', 'Аудио и видео', true, true, NOW(), NOW()),
            -- Подкатегории Аудио и видео (родитель id=240)
        ('category', 3210, 'ru', 'name', 'Телевизоры и проекторы', true, true, NOW(), NOW()),
        ('category', 3215, 'ru', 'name', 'Наушники', true, true, NOW(), NOW()),
        ('category', 3220, 'ru', 'name', 'Акустика, колонки, сабвуферы', true, true, NOW(), NOW()),
        ('category', 3225, 'ru', 'name', 'Аксессуары', true, true, NOW(), NOW()),
        ('category', 3230, 'ru', 'name', 'Музыкальные центры, магнитолы', true, true, NOW(), NOW()),
        ('category', 3235, 'ru', 'name', 'Усилители и ресиверы', true, true, NOW(), NOW()),
        ('category', 3240, 'ru', 'name', 'Видеокамеры', true, true, NOW(), NOW()),
        ('category', 3245, 'ru', 'name', 'Видео, DVD и Blu-ray плееры', true, true, NOW(), NOW()),
        ('category', 3250, 'ru', 'name', 'Кабели и адаптеры', true, true, NOW(), NOW()),
        ('category', 3255, 'ru', 'name', 'Музыка и фильмы', true, true, NOW(), NOW()),
        ('category', 3260, 'ru', 'name', 'Микрофоны', true, true, NOW(), NOW()),
        ('category', 3265, 'ru', 'name', 'MP3-плееры', true, true, NOW(), NOW()),

    ('category', 3300, 'ru', 'name', 'Товары для компьютера', true, true, NOW(), NOW()),
    -- Подкатегории Товары для компьютера (родитель id=3300)
        ('category', 3310, 'ru', 'name', 'Системные блоки', true, true, NOW(), NOW()),
        ('category', 3320, 'ru', 'name', 'Моноблоки', true, true, NOW(), NOW()),
        ('category', 3330, 'ru', 'name', 'Комплектующие', true, true, NOW(), NOW()),
        -- Подкатегории Комплектующих (родитель id=3330)
            ('category', 3331, 'ru', 'name', 'CD, DVD и Blu-ray приводы', true, true, NOW(), NOW()),
            ('category', 3332, 'ru', 'name', 'Блоки питания', true, true, NOW(), NOW()),
            ('category', 3333, 'ru', 'name', 'Видеокарты', true, true, NOW(), NOW()),
            ('category', 3334, 'ru', 'name', 'Жёсткие диски', true, true, NOW(), NOW()),
            ('category', 3335, 'ru', 'name', 'Звуковые карты', true, true, NOW(), NOW()),
            ('category', 3336, 'ru', 'name', 'Контроллеры', true, true, NOW(), NOW()),
            ('category', 3337, 'ru', 'name', 'Корпуса', true, true, NOW(), NOW()),
            ('category', 3338, 'ru', 'name', 'Материнские платы', true, true, NOW(), NOW()),
            ('category', 3339, 'ru', 'name', 'Оперативная память', true, true, NOW(), NOW()),
            ('category', 3340, 'ru', 'name', 'Процессоры', true, true, NOW(), NOW()),
            ('category', 3341, 'ru', 'name', 'Системы охлаждения', true, true, NOW(), NOW()),

        ('category', 3360, 'ru', 'name', 'Мониторы и запчасти', true, true, NOW(), NOW()),
        ('category', 3365, 'ru', 'name', 'Сетевое оборудование', true, true, NOW(), NOW()),
        ('category', 3370, 'ru', 'name', 'Клавиатуры и мыши', true, true, NOW(), NOW()),
        ('category', 3375, 'ru', 'name', 'Джойстики и рули', true, true, NOW(), NOW()),
        ('category', 3380, 'ru', 'name', 'Флэшки и карты памяти', true, true, NOW(), NOW()),
        ('category', 3385, 'ru', 'name', 'Веб-камеры', true, true, NOW(), NOW()),
        ('category', 3390, 'ru', 'name', 'ТВ-тюнеры', true, true, NOW(), NOW()),

    ('category', 3500, 'ru', 'name', 'Игры, приставки и программы', true, true, NOW(), NOW()),
            -- Подкатегории Игры, приставки и программы (родитель id=3500)
        ('category', 3510, 'ru', 'name', 'Игровые приставки', true, true, NOW(), NOW()),
        ('category', 3520, 'ru', 'name', 'Игры для приставок', true, true, NOW(), NOW()),
        ('category', 3530, 'ru', 'name', 'Компьютерные игры', true, true, NOW(), NOW()),
        ('category', 3540, 'ru', 'name', 'Программы', true, true, NOW(), NOW()),

    ('category', 3600, 'ru', 'name', 'Ноутбуки', true, true, NOW(), NOW()),
    ('category', 3700, 'ru', 'name', 'Фототехника', true, true, NOW(), NOW()),
            -- Подкатегории Фототехника (родитель id=3700)
        ('category', 3710, 'ru', 'name', 'Оборудование и аксессуары', true, true, NOW(), NOW()),
        ('category', 3720, 'ru', 'name', 'Объективы', true, true, NOW(), NOW()),
        ('category', 3730, 'ru', 'name', 'Компактные фотоаппараты', true, true, NOW(), NOW()),
        ('category', 3740, 'ru', 'name', 'Плёночные фотоаппараты', true, true, NOW(), NOW()),
        ('category', 3750, 'ru', 'name', 'Зеркальные фотоаппараты', true, true, NOW(), NOW()),

    ('category', 3800, 'ru', 'name', 'Планшеты и электронные книги', true, true, NOW(), NOW()),
            -- Подкатегории Планшеты и электронные книги (родитель id=3800)
        ('category', 3810, 'ru', 'name', 'Планшеты', true, true, NOW(), NOW()),
        ('category', 3820, 'ru', 'name', 'Электронные книги', true, true, NOW(), NOW()),
        ('category', 3830, 'ru', 'name', 'Аксессуары', true, true, NOW(), NOW()),

    ('category', 3900, 'ru', 'name', 'Оргтехника и расходники', true, true, NOW(), NOW()),
            -- Подкатегории Оргтехника и расходники (родитель id=3900)
        ('category', 3910, 'ru', 'name', 'МФУ, копиры и сканеры', true, true, NOW(), NOW()),
        ('category', 3920, 'ru', 'name', 'Принтеры', true, true, NOW(), NOW()),
        ('category', 3930, 'ru', 'name', 'Канцелярия', true, true, NOW(), NOW()),
        ('category', 3940, 'ru', 'name', 'ИБП, сетевые фильтры', true, true, NOW(), NOW()),
        ('category', 3950, 'ru', 'name', 'Телефония', true, true, NOW(), NOW()),
        ('category', 3960, 'ru', 'name', 'Уничтожители бумаг', true, true, NOW(), NOW()),
        ('category', 3970, 'ru', 'name', 'Расходные материалы', true, true, NOW(), NOW()),

    ('category', 4100, 'ru', 'name', 'Бытовая техника', true, true, NOW(), NOW()),
        -- Подкатегории Бытовая техника (родитель id=4100)
        ('category', 4110, 'ru', 'name', 'Для кухни', true, true, NOW(), NOW()),
                    -- Подкатегории техники для кухни (родитель id=4110)
            ('category', 4111, 'ru', 'name', 'Вытяжки', true, true, NOW(), NOW()),
            ('category', 4112, 'ru', 'name', 'Мелкая кухонная техника', true, true, NOW(), NOW()),
            ('category', 4113, 'ru', 'name', 'Микроволновые печи', true, true, NOW(), NOW()),
            ('category', 4114, 'ru', 'name', 'Плиты и духовые шкафы', true, true, NOW(), NOW()),
            ('category', 4115, 'ru', 'name', 'Посудомоечные машины', true, true, NOW(), NOW()),
            ('category', 4116, 'ru', 'name', 'Холодильники и морозильные камеры', true, true, NOW(), NOW()),

        ('category', 4120, 'ru', 'name', 'Для дома', true, true, NOW(), NOW()),
                    -- Подкатегории техники для дома (родитель id=4120)
            ('category', 4121, 'ru', 'name', 'Пылесосы и запчасти', true, true, NOW(), NOW()),
            ('category', 4122, 'ru', 'name', 'Стиральные и сушильные машины', true, true, NOW(), NOW()),
            ('category', 4123, 'ru', 'name', 'Утюги', true, true, NOW(), NOW()),
            ('category', 4124, 'ru', 'name', 'Швейное оборудование', true, true, NOW(), NOW()),

        ('category', 4130, 'ru', 'name', 'Климатическое оборудование', true, true, NOW(), NOW()),
                    -- Подкатегории климатического оборудования (родитель id=4130)
            ('category', 4131, 'ru', 'name', 'Вентиляторы', true, true, NOW(), NOW()),
            ('category', 4132, 'ru', 'name', 'Кондиционеры и запчасти', true, true, NOW(), NOW()),
            ('category', 4133, 'ru', 'name', 'Обогреватели', true, true, NOW(), NOW()),
            ('category', 4134, 'ru', 'name', 'Очистители воздуха', true, true, NOW(), NOW()),
            ('category', 4135, 'ru', 'name', 'Термометры и метеостанции', true, true, NOW(), NOW()),

        ('category', 4140, 'ru', 'name', 'Для индивидуального ухода', true, true, NOW(), NOW()),
            -- Подкатегории техники для индивидуального ухода (родитель id=4140)
            ('category', 4141, 'ru', 'name', 'Бритвы и триммеры', true, true, NOW(), NOW()),
            ('category', 4142, 'ru', 'name', 'Машинки для стрижки', true, true, NOW(), NOW()),
            ('category', 4143, 'ru', 'name', 'Фены и приборы для укладки', true, true, NOW(), NOW()),
            ('category', 4144, 'ru', 'name', 'Эпиляторы', true, true, NOW(), NOW()),
            
('category', 5000, 'ru', 'name', 'Все для дома и квартиры', true, true, NOW(), NOW()),
    -- Все для дома и квартиры (родитель id=5000)
    -- Ремонт и строительство
    ('category', 5100, 'ru', 'name', 'Ремонт и строительство', true, true, NOW(), NOW()),
        -- Подкатегории Ремонт и строительство (родитель id=5100)
        ('category', 5110, 'ru', 'name', 'Двери', true, true, NOW(), NOW()),
        ('category', 5115, 'ru', 'name', 'Инструменты', true, true, NOW(), NOW()),
        ('category', 5120, 'ru', 'name', 'Камины и обогреватели', true, true, NOW(), NOW()),
        ('category', 5125, 'ru', 'name', 'Окна и балконы', true, true, NOW(), NOW()),
        ('category', 5130, 'ru', 'name', 'Потолки', true, true, NOW(), NOW()),
        ('category', 5135, 'ru', 'name', 'Для сада и дачи', true, true, NOW(), NOW()),
        ('category', 5140, 'ru', 'name', 'Сантехника, водоснабжение и сауна', true, true, NOW(), NOW()),
        ('category', 5145, 'ru', 'name', 'Готовые строения и срубы', true, true, NOW(), NOW()),
        ('category', 5150, 'ru', 'name', 'Ворота, заборы и ограждения', true, true, NOW(), NOW()),
        ('category', 5155, 'ru', 'name', 'Охрана и сигнализации', true, true, NOW(), NOW()),

    ('category', 5200, 'ru', 'name', 'Мебель и интерьер', true, true, NOW(), NOW()),
        -- Подкатегории Мебель и интерьер (родитель id=5200)
        ('category', 5210, 'ru', 'name', 'Кровати, диваны и кресла', true, true, NOW(), NOW()),
        ('category', 5215, 'ru', 'name', 'Текстиль и ковры', true, true, NOW(), NOW()),
        ('category', 5220, 'ru', 'name', 'Освещение', true, true, NOW(), NOW()),
        ('category', 5225, 'ru', 'name', 'Компьютерные столы и кресла', true, true, NOW(), NOW()),
        ('category', 5230, 'ru', 'name', 'Шкафы, комоды и стеллажи', true, true, NOW(), NOW()),
        ('category', 5235, 'ru', 'name', 'Кухонные гарнитуры', true, true, NOW(), NOW()),
        ('category', 5240, 'ru', 'name', 'Столы и стулья', true, true, NOW(), NOW()),
        ('category', 5250, 'ru', 'name', 'Комнатные растения', true, true, NOW(), NOW()),
            -- Подкатегории Комнатных растений (родитель id=5250)
            ('category', 5251, 'ru', 'name', 'Декоративно-лиственные растения', true, true, NOW(), NOW()),
            ('category', 5252, 'ru', 'name', 'Цветущие растения', true, true, NOW(), NOW()),
            ('category', 5253, 'ru', 'name', 'Пальмы и фикусы', true, true, NOW(), NOW()),
            ('category', 5254, 'ru', 'name', 'Кактусы и суккуленты', true, true, NOW(), NOW()),

    ('category', 5300, 'ru', 'name', 'Продукты питания', true, true, NOW(), NOW()),
            -- Подкатегории Продукты питания (родитель id=5300)
        ('category', 5310, 'ru', 'name', 'Чай, кофе, какао', true, true, NOW(), NOW()),
        ('category', 5315, 'ru', 'name', 'Напитки', true, true, NOW(), NOW()),
        ('category', 5320, 'ru', 'name', 'Рыба, морепродукты, икра', true, true, NOW(), NOW()),
        ('category', 5325, 'ru', 'name', 'Мясо, птица, субпродукты', true, true, NOW(), NOW()),
        ('category', 5330, 'ru', 'name', 'Кондитерские изделия', true, true, NOW(), NOW()),
        ('category', 5340, 'ru', 'name', 'Ракия и вино', true, true, NOW(), NOW()),
                    -- Подкатегории Ракии и вина (родитель id=5340)
            ('category', 5341, 'ru', 'name', 'Сливовая ракия', true, true, NOW(), NOW()),
            ('category', 5342, 'ru', 'name', 'Виноградная ракия', true, true, NOW(), NOW()),
            ('category', 5343, 'ru', 'name', 'Фруктовая ракия', true, true, NOW(), NOW()),
            ('category', 5344, 'ru', 'name', 'Домашнее вино', true, true, NOW(), NOW()),

        ('category', 5350, 'ru', 'name', 'Домашние сыры', true, true, NOW(), NOW()),
        ('category', 5360, 'ru', 'name', 'Каймак', true, true, NOW(), NOW()),
        ('category', 5370, 'ru', 'name', 'Айвар', true, true, NOW(), NOW()),

    ('category', 5400, 'ru', 'name', 'Посуда и товары для кухни', true, true, NOW(), NOW()),
        -- Подкатегории Посуда и товары для кухни (родитель id=5400)
        ('category', 5405, 'ru', 'name', 'Посуда', true, true, NOW(), NOW()),
        ('category', 5410, 'ru', 'name', 'Товары для кухни', true, true, NOW(), NOW()),
        ('category', 5415, 'ru', 'name', 'Сервировка стола', true, true, NOW(), NOW()),
        ('category', 5420, 'ru', 'name', 'Приготовление пищи', true, true, NOW(), NOW()),
        ('category', 5425, 'ru', 'name', 'Хранение продуктов', true, true, NOW(), NOW()),
        ('category', 5430, 'ru', 'name', 'Приготовление напитков', true, true, NOW(), NOW()),
        ('category', 5435, 'ru', 'name', 'Хозяйственные товары', true, true, NOW(), NOW()),

('category', 6000, 'ru', 'name', 'Все для сада', true, true, NOW(), NOW()),
    -- Все для сада ('category', подкатегории 51-60)
    ('category', 6050, 'ru', 'name', 'Садовая мебель', true, true, NOW(), NOW()),
    ('category', 6100, 'ru', 'name', 'Садовые инструменты', true, true, NOW(), NOW()),
    ('category', 6150, 'ru', 'name', 'Растения для сада', true, true, NOW(), NOW()),
    ('category', 6200, 'ru', 'name', 'Семена и рассада', true, true, NOW(), NOW()),
    ('category', 6250, 'ru', 'name', 'Барбекю и аксессуары', true, true, NOW(), NOW()),
    ('category', 6300, 'ru', 'name', 'Бассейны и оборудование', true, true, NOW(), NOW()),
    ('category', 6350, 'ru', 'name', 'Системы полива', true, true, NOW(), NOW()),
    ('category', 6400, 'ru', 'name', 'Компостирование', true, true, NOW(), NOW()),
    ('category', 6450, 'ru', 'name', 'Теплицы и парники', true, true, NOW(), NOW()),
    ('category', 6500, 'ru', 'name', 'Удобрения и грунты', true, true, NOW(), NOW()),
    ('category', 6550, 'ru', 'name', 'Освещение', true, true, NOW(), NOW()),
    ('category', 6600, 'ru', 'name', 'Оформление интерьера', true, true, NOW(), NOW()),
    ('category', 6650, 'ru', 'name', 'Растения и семена', true, true, NOW(), NOW()),
    ('category', 6700, 'ru', 'name', 'Сад и огород', true, true, NOW(), NOW()),
    ('category', 6750, 'ru', 'name', 'Садовые растения', true, true, NOW(), NOW()),
            -- Подкатегории Садовые растения (родитель id=6750)
        ('category', 6751, 'ru', 'name', 'Декоративные кустарники и деревья', true, true, NOW(), NOW()),
        ('category', 6752, 'ru', 'name', 'Хвойные растения', true, true, NOW(), NOW()),
        ('category', 6753, 'ru', 'name', 'Многолетние растения', true, true, NOW(), NOW()),
        ('category', 6754, 'ru', 'name', 'Плодовые растения', true, true, NOW(), NOW()),
        ('category', 6755, 'ru', 'name', 'Газон', true, true, NOW(), NOW()),
        ('category', 6756, 'ru', 'name', 'Зелень и пряные травы', true, true, NOW(), NOW()),

    ('category', 6850, 'ru', 'name', 'Семена, луковицы, клубни', true, true, NOW(), NOW()),
    ('category', 6900, 'ru', 'name', 'Товары для ухода за растениями', true, true, NOW(), NOW()),
            -- Подкатегории Товары для ухода за растениями (родитель id=6900)
        ('category', 6901, 'ru', 'name', 'Грунты и субстраты', true, true, NOW(), NOW()),
        ('category', 6902, 'ru', 'name', 'Удобрения', true, true, NOW(), NOW()),
        ('category', 6903, 'ru', 'name', 'Средства от вредителей и сорняков', true, true, NOW(), NOW()),
        ('category', 6904, 'ru', 'name', 'Горшки и кашпо', true, true, NOW(), NOW()),
        ('category', 6905, 'ru', 'name', 'Фитолампы', true, true, NOW(), NOW()),
        ('category', 6906, 'ru', 'name', 'Измерители влаги', true, true, NOW(), NOW()),
        ('category', 6907, 'ru', 'name', 'Теплицы, грядки, клумбы', true, true, NOW(), NOW()),

    ('category', 6950, 'ru', 'name', 'Панно и искусственные растения', true, true, NOW(), NOW()),

('category', 7000, 'ru', 'name', 'Хобби и отдых', true, true, NOW(), NOW()),
        -- Хобби и отдых (родитель id=7000)
    ('category', 7050, 'ru', 'name', 'Музыкальные инструменты', true, true, NOW(), NOW()),
        -- Музыкальные инструменты (родитель id=7050)
        ('category', 7055, 'ru', 'name', 'Струнные инструменты', true, true, NOW(), NOW()),
        ('category', 7060, 'ru', 'name', 'Фортепиано и клавишные', true, true, NOW(), NOW()),
        ('category', 7065, 'ru', 'name', 'Ударные инструменты', true, true, NOW(), NOW()),
        ('category', 7070, 'ru', 'name', 'Духовые инструменты', true, true, NOW(), NOW()),
        ('category', 7075, 'ru', 'name', 'Аккордеоны и гармони', true, true, NOW(), NOW()),
        ('category', 7080, 'ru', 'name', 'Аудиооборудование', true, true, NOW(), NOW()),
        ('category', 7085, 'ru', 'name', 'Аксессуары для инструментов', true, true, NOW(), NOW()),

    ('category', 7100, 'ru', 'name', 'Книги и журналы', true, true, NOW(), NOW()),
            -- Подкатегории Книги и журналы (родитель id=7100)
        ('category', 7105, 'ru', 'name', 'Журналы, газеты, брошюры', true, true, NOW(), NOW()),
        ('category', 7115, 'ru', 'name', 'Книги', true, true, NOW(), NOW()),
        ('category', 7130, 'ru', 'name', 'Учебная литература', true, true, NOW(), NOW()),
    ('category', 7150, 'ru', 'name', 'Спортивный инвентарь', true, true, NOW(), NOW()),
    ('category', 7200, 'ru', 'name', 'Путешествия', true, true, NOW(), NOW()),
    ('category', 7250, 'ru', 'name', 'Коллекционирование', true, true, NOW(), NOW()),
            -- Подкатегории Коллекционирование (родитель id=7250)
        ('category', 7251, 'ru', 'name', 'Банкноты', true, true, NOW(), NOW()),
        ('category', 7252, 'ru', 'name', 'Билеты', true, true, NOW(), NOW()),
        ('category', 7253, 'ru', 'name', 'Вещи знаменитостей, автографы', true, true, NOW(), NOW()),
        ('category', 7254, 'ru', 'name', 'Военные вещи', true, true, NOW(), NOW()),
        ('category', 7255, 'ru', 'name',  'Грампластинки', 'vinyl', 7250, 'vinyl', '2025-02-07 07:13:52.823283'),
        ('category', 7256, 'ru', 'name', 'Документы', true, true, NOW(), NOW()),
        ('category', 7257, 'ru', 'name', 'Жетоны, медали, значки', true, true, NOW(), NOW()),
        ('category', 7258, 'ru', 'name', 'Игры', true, true, NOW(), NOW()),
        ('category', 7259, 'ru', 'name', 'Календари', true, true, NOW(), NOW()),
        ('category', 7261, 'ru', 'name', 'Картины', true, true, NOW(), NOW()),
        ('category', 7262, 'ru', 'name', 'Марки', true, true, NOW(), NOW()),
        ('category', 7263, 'ru', 'name', 'Модели', true, true, NOW(), NOW()),
        ('category', 7264, 'ru', 'name', 'Монеты', true, true, NOW(), NOW()),

    ('category', 7300, 'ru', 'name', 'Предметы искусства', true, true, NOW(), NOW()),
    ('category', 7350, 'ru', 'name', 'Игрушки', true, true, NOW(), NOW()),
    ('category', 7400, 'ru', 'name', 'Велосипеды', true, true, NOW(), NOW()),
            -- Подкатегории Велосипеды (родитель id=7400)
        ('category', 7410, 'ru', 'name', 'ВМХ', true, true, NOW(), NOW()),
        ('category', 7415, 'ru', 'name', 'Городские', true, true, NOW(), NOW()),
        ('category', 7420, 'ru', 'name', 'Шоссейные', true, true, NOW(), NOW()),
        ('category', 7425, 'ru', 'name', 'Детские', true, true, NOW(), NOW()),
        ('category', 7430, 'ru', 'name', 'Горные', true, true, NOW(), NOW()),
        ('category', 7435, 'ru', 'name', 'Запчасти и аксессуары', true, true, NOW(), NOW()),

    ('category', 7500, 'ru', 'name', 'Охота и рыбалка', true, true, NOW(), NOW()),
        -- Подкатегории Охота и рыбалка (родитель id=7500)
        ('category', 7510, 'ru', 'name', 'Ножи, мультитулы, топоры', true, true, NOW(), NOW()),
        ('category', 7520, 'ru', 'name', 'Охота', true, true, NOW(), NOW()),
                    -- Подкатегории Охота (родитель id=7520)
            ('category', 7521, 'ru', 'name', 'Прицелы', true, true, NOW(), NOW()),
            ('category', 7522, 'ru', 'name', 'Аксессуары для прицелов', true, true, NOW(), NOW()),
            ('category', 7523, 'ru', 'name', 'Монокуляры, бинокли, дальномеры', true, true, NOW(), NOW()),

        ('category', 7530, 'ru', 'name', 'Рыбалка', true, true, NOW(), NOW()),
            -- Подкатегории Рыбалка (родитель id=7530)
            ('category', 7531, 'ru', 'name', 'Удочки, спиннинги и катушки', true, true, NOW(), NOW()),
            ('category', 7551, 'ru', 'name', 'Приманки и снасти', true, true, NOW(), NOW()),
            ('category', 7571, 'ru', 'name', 'Эхолоты и снаряжение', true, true, NOW(), NOW()),

    ('category', 7650, 'ru', 'name', 'Кемпинг', true, true, NOW(), NOW()),
    ('category', 7700, 'ru', 'name', 'Антиквариат', true, true, NOW(), NOW()),
    ('category', 7750, 'ru', 'name', 'Билеты, мероприятия и путешествия', true, true, NOW(), NOW()),
        -- Подкатегории Билеты и путешествия (родитель id=7750)
        ('category', 7751, 'ru', 'name', 'Карты, купоны', true, true, NOW(), NOW()),
        ('category', 7752, 'ru', 'name', 'Концерты', true, true, NOW(), NOW()),
        ('category', 7753, 'ru', 'name', 'Путешествия', true, true, NOW(), NOW()),
        ('category', 7754, 'ru', 'name', 'Спорт', true, true, NOW(), NOW()),
        ('category', 7755, 'ru', 'name', 'Театр, опера, балет', true, true, NOW(), NOW()),
        ('category', 7756, 'ru', 'name', 'Цирк, кино', true, true, NOW(), NOW()),
        ('category', 7758, 'ru', 'name', 'Шоу, мюзикл', true, true, NOW(), NOW()),

    ('category', 7800, 'ru', 'name', 'Спорт', true, true, NOW(), NOW()),
            -- Подкатегории Спорт и отдых (родитель id=7800)
        ('category', 7805, 'ru', 'name', 'Бильярд и боулинг', true, true, NOW(), NOW()),
        ('category', 7810, 'ru', 'name', 'Дайвинг и водный спорт', true, true, NOW(), NOW()),
        ('category', 7815, 'ru', 'name', 'Единоборства', true, true, NOW(), NOW()),
        ('category', 7820, 'ru', 'name', 'Зимний спорт', true, true, NOW(), NOW()),
        ('category', 7825, 'ru', 'name', 'Игры с мячом', true, true, NOW(), NOW()),
        ('category', 7830, 'ru', 'name', 'Настольные игры', true, true, NOW(), NOW()),
        ('category', 7835, 'ru', 'name', 'Пейнтбол и страйкбол', true, true, NOW(), NOW()),
        ('category', 7840, 'ru', 'name', 'Ролики и скейтбординг', true, true, NOW(), NOW()),
        ('category', 7845, 'ru', 'name', 'Теннис, бадминтон, пинг-понг', true, true, NOW(), NOW()),
        ('category', 7850, 'ru', 'name', 'Туризм и отдых на природе', true, true, NOW(), NOW()),
        ('category', 7855, 'ru', 'name', 'Фитнес и тренажёры', true, true, NOW(), NOW()),
        ('category', 7860, 'ru', 'name', 'Спортивное питание', true, true, NOW(), NOW()),

    -- Традиционные ремесла и сувениры
    ('category', 7850, 'ru', 'name', 'Народное ремесло и рукоделие', true, true, NOW(), NOW()),
        -- Подкатегории народных ремесел (родитель id=7850)
        ('category', 7851, 'ru', 'name', 'Опанци', true, true, NOW(), NOW()),
        ('category', 7852, 'ru', 'name', 'Керамика', true, true, NOW(), NOW()),
        ('category', 7853, 'ru', 'name', 'Вышивка', true, true, NOW(), NOW()),
        ('category', 7854, 'ru', 'name', 'Ткачество', true, true, NOW(), NOW()),
        ('category', 7855, 'ru', 'name', 'Народные инструменты', true, true, NOW(), NOW()),
        ('category', 7856, 'ru', 'name', 'Деревообработка', true, true, NOW(), NOW()),

    -- Сельскохозяйственные категории
    ('category', 7900, 'ru', 'name', 'Пчеловодство', true, true, NOW(), NOW()),
        -- Подкатегории пчеловодства (родитель id=7900)
        ('category', 7910, 'ru', 'name', 'Мёд', true, true, NOW(), NOW()),
        ('category', 7920, 'ru', 'name', 'Пчелиный воск', true, true, NOW(), NOW()),
        ('category', 7930, 'ru', 'name', 'Прополис', true, true, NOW(), NOW()),
        ('category', 7935, 'ru', 'name', 'Пчеловодный инвентарь', true, true, NOW(), NOW()),
        ('category', 7945, 'ru', 'name', 'Пчелы', true, true, NOW(), NOW()),

    -- Туристические услуги
    ('category', 7950, 'ru', 'name', 'Сельский туризм', true, true, NOW(), NOW()),
            -- Подкатегории сельского туризма (родитель id=7950)
        ('category', 7951, 'ru', 'name', 'Этно-деревни', true, true, NOW(), NOW()),
        ('category', 7952, 'ru', 'name', 'Винные туры', true, true, NOW(), NOW()),
        ('category', 7953, 'ru', 'name', 'Агротуризм', true, true, NOW(), NOW()),
        ('category', 7954, 'ru', 'name', 'Горный туризм', true, true, NOW(), NOW()),

('category', 8000, 'ru', 'name', 'Животные', true, true, NOW(), NOW()),
    -- Animals categories (родитель id=8000)
    ('category', 8050, 'ru', 'name', 'Собаки', true, true, NOW(), NOW()),
    ('category', 8100, 'ru', 'name', 'Кошки', true, true, NOW(), NOW()),
    ('category', 8150, 'ru', 'name', 'Птицы', true, true, NOW(), NOW()),
    ('category', 8200, 'ru', 'name', 'Аквариум', true, true, NOW(), NOW()),
            -- Подкатегории Аквариум (родитель id=8200)
        ('category', 8205, 'ru', 'name', 'Аквариумы', true, true, NOW(), NOW()),
        ('category', 8210, 'ru', 'name', 'Рыбы', true, true, NOW(), NOW()),
        ('category', 8215, 'ru', 'name', 'Другие аквариумные животные', true, true, NOW(), NOW()),
        ('category', 8220, 'ru', 'name', 'Оборудование', true, true, NOW(), NOW()),
        ('category', 8225, 'ru', 'name', 'Растения', true, true, NOW(), NOW()),
        ('category', 8230, 'ru', 'name', 'Аквариумная мебель', true, true, NOW(), NOW()),
        ('category', 8235, 'ru', 'name', 'Морская аквариумистика', true, true, NOW(), NOW()),

    ('category', 8250, 'ru', 'name', 'Другие животные', true, true, NOW(), NOW()),
        -- Подкатегории Другие животные (родитель id=8250)
        ('category', 8251, 'ru', 'name', 'Амфибии', true, true, NOW(), NOW()),
        ('category', 8252, 'ru', 'name', 'Грызуны', true, true, NOW(), NOW()),
        ('category', 8253, 'ru', 'name', 'Кролики', true, true, NOW(), NOW()),
        ('category', 8254, 'ru', 'name', 'Лошади', true, true, NOW(), NOW()),
        ('category', 8255, 'ru', 'name', 'Рептилии', true, true, NOW(), NOW()),
        ('category', 8256, 'ru', 'name', 'Сельхоз животные', true, true, NOW(), NOW()),
        ('category', 8257, 'ru', 'name', 'Домашние птицы', true, true, NOW(), NOW()),
        ('category', 8258, 'ru', 'name', 'Товары для животных', true, true, NOW(), NOW()),

('category', 8500, 'ru', 'name', 'Готовый бизнес и оборудование', true, true, NOW(), NOW()),


('category', 9000, 'ru', 'name', 'Работа', true, true, NOW(), NOW()),
    ('category', 9050, 'ru', 'name', 'Вакансии', true, true, NOW(), NOW()),
    ('category', 9100, 'ru', 'name', 'Резюме', true, true, NOW(), NOW()),
    ('category', 9150, 'ru', 'name', 'Удаленная работа', true, true, NOW(), NOW()),
    ('category', 9200, 'ru', 'name', 'Партнерство и сотрудничество', true, true, NOW(), NOW()),
    ('category', 9250, 'ru', 'name', 'Обучение и стажировка', true, true, NOW(), NOW()),
    ('category', 9300, 'ru', 'name', 'Сезонные работы', true, true, NOW(), NOW()),
        -- Подкатегории сезонных работ (родитель id=9300)
        ('category', 9310, 'ru', 'name', 'Сбор урожая', true, true, NOW(), NOW()),
        ('category', 9315, 'ru', 'name', 'Работа на винограднике', true, true, NOW(), NOW()),
        ('category', 9320, 'ru', 'name', 'Сезонные строительные работы', true, true, NOW(), NOW()),

('category', 9500, 'ru', 'name', 'Одежда, обувь, аксессуары', true, true, NOW(), NOW()),

-- Товары для детей и игрушки ('category', новая категория)
('category', 9700, 'ru', 'name', 'Товары для детей и игрушки', true, true, NOW(), NOW()),
    -- Подкатегории Товары для детей и игрушки (родитель id=9700)
    ('category', 9705, 'ru', 'name', 'Детские коляски', true, true, NOW(), NOW()),
    ('category', 9710, 'ru', 'name', 'Детская мебель', true, true, NOW(), NOW()),
    ('category', 9715, 'ru', 'name', 'Велосипеды и самокаты', true, true, NOW(), NOW()),
    ('category', 9720, 'ru', 'name', 'Товары для кормления', true, true, NOW(), NOW()),
    ('category', 9725, 'ru', 'name', 'Автомобильные кресла', true, true, NOW(), NOW()),
    ('category', 9730, 'ru', 'name', 'Игрушки', true, true, NOW(), NOW()),
    ('category', 9735, 'ru', 'name', 'Постельные принадлежности', true, true, NOW(), NOW()),
    ('category', 9740, 'ru', 'name', 'Товары для купания', true, true, NOW(), NOW()),
    ('category', 9745, 'ru', 'name', 'Товары для школы', true, true, NOW(), NOW()),
    ('category', 9750, 'ru', 'name', 'Детская одежда и обувь, аксессуары', true, true, NOW(), NOW()),

('category', 9999, 'ru', 'name', 'Прочее', true, true, NOW(), NOW()),

-- Обновляем sequence для translations
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);