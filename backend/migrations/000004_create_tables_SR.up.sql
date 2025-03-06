-- Српски преводи категорија за табелу translations
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated, is_verified, created_at, updated_at) VALUES

('category', 1000, 'sr', 'name', 'Некретнине', true, true, NOW(), NOW()),
    -- Подкатегорије некретнина 1000
    ('category', 1100, 'sr', 'name', 'Стан', true, true, NOW(), NOW()),
    ('category', 1200, 'sr', 'name', 'Соба', true, true, NOW(), NOW()),
    ('category', 1300, 'sr', 'name', 'Кућа, викендица, вила', true, true, NOW(), NOW()),
    ('category', 1400, 'sr', 'name', 'Земљишна парцела', true, true, NOW(), NOW()),
    ('category', 1500, 'sr', 'name', 'Гаража и паркинг место', true, true, NOW(), NOW()),
    ('category', 1600, 'sr', 'name', 'Пословни простор', true, true, NOW(), NOW()),
    ('category', 1700, 'sr', 'name', 'Некретнине у иностранству', true, true, NOW(), NOW()),
    ('category', 1800, 'sr', 'name', 'Хотел', true, true, NOW(), NOW()),
    ('category', 1900, 'sr', 'name', 'Апартмани', true, true, NOW(), NOW()),

('category', 2000, 'sr', 'name', 'Возила', true, true, NOW(), NOW()),

    -- Ауто категорије 2000
    ('category', 2100, 'sr', 'name', 'Путнички аутомобили', true, true, NOW(), NOW()),
    ('category', 2200, 'sr', 'name', 'Комерцијална возила', true, true, NOW(), NOW()),
        -- Подкатегорије комерцијалних возила 2200
        ('category', 2210, 'sr', 'name', 'Камиони', true, true, NOW(), NOW()),
        ('category', 2220, 'sr', 'name', 'Полуприколице', true, true, NOW(), NOW()),
        ('category', 2230, 'sr', 'name', 'Лака комерцијална возила', true, true, NOW(), NOW()),
        ('category', 2240, 'sr', 'name', 'Аутобуси', true, true, NOW(), NOW()),
        
    ('category', 2300, 'sr', 'name', 'Специјална механизација', true, true, NOW(), NOW()),
        -- Подкатегорије специјалне механизације 2300
        ('category', 2310, 'sr', 'name', 'Багери', true, true, NOW(), NOW()),
        ('category', 2315, 'sr', 'name', 'Утоваривачи', true, true, NOW(), NOW()),
        ('category', 2320, 'sr', 'name', 'Комбиновани багери', true, true, NOW(), NOW()),
        ('category', 2325, 'sr', 'name', 'Ауто-дизалице', true, true, NOW(), NOW()),
        ('category', 2330, 'sr', 'name', 'Аутомиксери', true, true, NOW(), NOW()),
        ('category', 2335, 'sr', 'name', 'Друмски ваљци', true, true, NOW(), NOW()),
        ('category', 2340, 'sr', 'name', 'Чистилице', true, true, NOW(), NOW()),
        ('category', 2345, 'sr', 'name', 'Камиони за смеће', true, true, NOW(), NOW()),
        ('category', 2350, 'sr', 'name', 'Аутокорпе', true, true, NOW(), NOW()),
        ('category', 2355, 'sr', 'name', 'Булдожери', true, true, NOW(), NOW()),
        ('category', 2360, 'sr', 'name', 'Грејдери', true, true, NOW(), NOW()),
        ('category', 2365, 'sr', 'name', 'Бушилице', true, true, NOW(), NOW()),

    ('category', 2400, 'sr', 'name', 'Пољопривредна механизација', true, true, NOW(), NOW()),
        -- Подкатегорије пољопривредне механизације 2400
        ('category', 2410, 'sr', 'name', 'Трактори', true, true, NOW(), NOW()),
        ('category', 2415, 'sr', 'name', 'Мини трактори', true, true, NOW(), NOW()),
        ('category', 2420, 'sr', 'name', 'Балирке', true, true, NOW(), NOW()),
        ('category', 2425, 'sr', 'name', 'Дрљаче', true, true, NOW(), NOW()),
        ('category', 2430, 'sr', 'name', 'Косачице', true, true, NOW(), NOW()),
        ('category', 2435, 'sr', 'name', 'Комбајни', true, true, NOW(), NOW()),
        ('category', 2440, 'sr', 'name', 'Телескопски утоваривачи', true, true, NOW(), NOW()),
        ('category', 2445, 'sr', 'name', 'Сејачице', true, true, NOW(), NOW()),
        ('category', 2450, 'sr', 'name', 'Култиватори', true, true, NOW(), NOW()),
        ('category', 2455, 'sr', 'name', 'Плугови', true, true, NOW(), NOW()),
        ('category', 2460, 'sr', 'name', 'Прскалице', true, true, NOW(), NOW()),

    ('category', 2500, 'sr', 'name', 'Изнајмљивање возила и механизације', true, true, NOW(), NOW()),
            -- Изнајмљивање возила и механизације (надређена id=2500)
        ('category', 2510, 'sr', 'name', 'Аутомобили', true, true, NOW(), NOW()),
        ('category', 2515, 'sr', 'name', 'Опрема за подизање', true, true, NOW(), NOW()),
        ('category', 2520, 'sr', 'name', 'Опрема за земљане радове', true, true, NOW(), NOW()),
        ('category', 2525, 'sr', 'name', 'Комунална механизација', true, true, NOW(), NOW()),
        ('category', 2530, 'sr', 'name', 'Опрема за изградњу путева', true, true, NOW(), NOW()),
        ('category', 2535, 'sr', 'name', 'Теретни транспорт', true, true, NOW(), NOW()),
        ('category', 2540, 'sr', 'name', 'Опрема за утовар', true, true, NOW(), NOW()),
        ('category', 2545, 'sr', 'name', 'Прикључци', true, true, NOW(), NOW()),
        ('category', 2550, 'sr', 'name', 'Приколице', true, true, NOW(), NOW()),
        ('category', 2555, 'sr', 'name', 'Пољопривредна опрема', true, true, NOW(), NOW()),
        ('category', 2560, 'sr', 'name', 'Камп кућице', true, true, NOW(), NOW()),

    ('category', 2600, 'sr', 'name', 'Мотоцикли и моторна возила', true, true, NOW(), NOW()),
        -- Мотоцикли и моторна возила (надређена id=2600)
        ('category', 2610, 'sr', 'name', 'Теренска возила', true, true, NOW(), NOW()),
        ('category', 2615, 'sr', 'name', 'Картинг', true, true, NOW(), NOW()),
        ('category', 2620, 'sr', 'name', 'Квадови и багији', true, true, NOW(), NOW()),
        ('category', 2625, 'sr', 'name', 'Мопеди', true, true, NOW(), NOW()),
        ('category', 2630, 'sr', 'name', 'Скутери', true, true, NOW(), NOW()),
        ('category', 2635, 'sr', 'name', 'Мотоцикли', true, true, NOW(), NOW()),
        ('category', 2640, 'sr', 'name', 'Моторне санке', true, true, NOW(), NOW()),

    ('category', 2700, 'sr', 'name', 'Пловила', true, true, NOW(), NOW()),
        -- Пловила (надређена id=2700)
        ('category', 2710, 'sr', 'name', 'Чамци на весла', true, true, NOW(), NOW()),
        ('category', 2720, 'sr', 'name', 'Кајаци', true, true, NOW(), NOW()),
        ('category', 2730, 'sr', 'name', 'Џет-ски', true, true, NOW(), NOW()),
        ('category', 2740, 'sr', 'name', 'Чамци и јахте', true, true, NOW(), NOW()),
        ('category', 2750, 'sr', 'name', 'Моторни чамци и мотори', true, true, NOW(), NOW()),

    ('category', 2800, 'sr', 'name', 'Делови и додаци', true, true, NOW(), NOW()),
        -- Делови и додаци (надређена id=2800)
        ('category', 2810, 'sr', 'name', 'Резервни делови', true, true, NOW(), NOW()),
        ('category', 2815, 'sr', 'name', 'Гуме, фелне и точкови', true, true, NOW(), NOW()),
        ('category', 2820, 'sr', 'name', 'Аудио и видео опрема', true, true, NOW(), NOW()),
        ('category', 2825, 'sr', 'name', 'Додаци', true, true, NOW(), NOW()),
        ('category', 2830, 'sr', 'name', 'Уља и ауто-хемија', true, true, NOW(), NOW()),
        ('category', 2835, 'sr', 'name', 'Алати', true, true, NOW(), NOW()),
        ('category', 2840, 'sr', 'name', 'Кровни носачи и куке', true, true, NOW(), NOW()),
        ('category', 2845, 'sr', 'name', 'Приколице', true, true, NOW(), NOW()),
        ('category', 2850, 'sr', 'name', 'Опрема', true, true, NOW(), NOW()),
        ('category', 2855, 'sr', 'name', 'Уређаји против крађе', true, true, NOW(), NOW()),
        ('category', 2860, 'sr', 'name', 'ГПС навигација', true, true, NOW(), NOW()),

('category', 3000, 'sr', 'name', 'Електроника', true, true, NOW(), NOW()),
    -- Електроника (надређена id=3000)
    ('category', 3100, 'sr', 'name', 'Телефони', true, true, NOW(), NOW()),
            -- Подкатегорије телефона (надређена id=3100)
        ('category', 3110, 'sr', 'name', 'Мобилни телефони', true, true, NOW(), NOW()),
        ('category', 3120, 'sr', 'name', 'Додаци', true, true, NOW(), NOW()),
                    -- Подкатегорије додатака за телефоне (надређена id=3120)
            ('category', 3121, 'sr', 'name', 'Батерије', true, true, NOW(), NOW()),
            ('category', 3122, 'sr', 'name', 'Слушалице', true, true, NOW(), NOW()),
            ('category', 3123, 'sr', 'name', 'Пуњачи', true, true, NOW(), NOW()),
            ('category', 3124, 'sr', 'name', 'Каблови и адаптери', true, true, NOW(), NOW()),
            ('category', 3125, 'sr', 'name', 'Модеми и рутери', true, true, NOW(), NOW()),
            ('category', 3126, 'sr', 'name', 'Маске и заштитне фолије', true, true, NOW(), NOW()),
            ('category', 3127, 'sr', 'name', 'Резервни делови', true, true, NOW(), NOW()),

        ('category', 3130, 'sr', 'name', 'Радио станице', true, true, NOW(), NOW()),
        ('category', 3140, 'sr', 'name', 'Фиксни телефони', true, true, NOW(), NOW()),

    ('category', 3200, 'sr', 'name', 'Аудио и видео', true, true, NOW(), NOW()),
            -- Подкатегорије аудио и видео (надређена id=3200)
        ('category', 3210, 'sr', 'name', 'Телевизори и пројектори', true, true, NOW(), NOW()),
        ('category', 3215, 'sr', 'name', 'Слушалице', true, true, NOW(), NOW()),
        ('category', 3220, 'sr', 'name', 'Звучници, системи, сабвуфери', true, true, NOW(), NOW()),
        ('category', 3225, 'sr', 'name', 'Додаци', true, true, NOW(), NOW()),
        ('category', 3230, 'sr', 'name', 'Музички центри, касетофони', true, true, NOW(), NOW()),
        ('category', 3235, 'sr', 'name', 'Појачала и ресивери', true, true, NOW(), NOW()),
        ('category', 3240, 'sr', 'name', 'Видео камере', true, true, NOW(), NOW()),
        ('category', 3245, 'sr', 'name', 'Видео, ДВД и Блу-реј плејери', true, true, NOW(), NOW()),
        ('category', 3250, 'sr', 'name', 'Каблови и адаптери', true, true, NOW(), NOW()),
        ('category', 3255, 'sr', 'name', 'Музика и филмови', true, true, NOW(), NOW()),
        ('category', 3260, 'sr', 'name', 'Микрофони', true, true, NOW(), NOW()),
        ('category', 3265, 'sr', 'name', 'МП3 плејери', true, true, NOW(), NOW()),

    ('category', 3300, 'sr', 'name', 'Рачунарска опрема', true, true, NOW(), NOW()),
    -- Подкатегорије рачунарске опреме (надређена id=3300)
        ('category', 3310, 'sr', 'name', 'Десктоп рачунари', true, true, NOW(), NOW()),
        ('category', 3320, 'sr', 'name', 'Монитори', true, true, NOW(), NOW()),
        ('category', 3330, 'sr', 'name', 'Компоненте', true, true, NOW(), NOW()),
        -- Подкатегорије компоненти (надређена id=3330)
            ('category', 3331, 'sr', 'name', 'ЦД, ДВД и Блу-реј уређаји', true, true, NOW(), NOW()),
            ('category', 3332, 'sr', 'name', 'Напајања', true, true, NOW(), NOW()),
            ('category', 3333, 'sr', 'name', 'Графичке карте', true, true, NOW(), NOW()),
            ('category', 3334, 'sr', 'name', 'Хард дискови', true, true, NOW(), NOW()),
            ('category', 3335, 'sr', 'name', 'Звучне карте', true, true, NOW(), NOW()),
            ('category', 3336, 'sr', 'name', 'Контролери', true, true, NOW(), NOW()),
            ('category', 3337, 'sr', 'name', 'Кућишта', true, true, NOW(), NOW()),
            ('category', 3338, 'sr', 'name', 'Матичне плоче', true, true, NOW(), NOW()),
            ('category', 3339, 'sr', 'name', 'РАМ меморија', true, true, NOW(), NOW()),
            ('category', 3340, 'sr', 'name', 'Процесори', true, true, NOW(), NOW()),
            ('category', 3341, 'sr', 'name', 'Системи за хлађење', true, true, NOW(), NOW()),

        ('category', 3360, 'sr', 'name', 'Монитори и делови', true, true, NOW(), NOW()),
        ('category', 3365, 'sr', 'name', 'Мрежна опрема', true, true, NOW(), NOW()),
        ('category', 3370, 'sr', 'name', 'Тастатуре и мишеви', true, true, NOW(), NOW()),
        ('category', 3375, 'sr', 'name', 'Џојстици и волани', true, true, NOW(), NOW()),
        ('category', 3380, 'sr', 'name', 'Флеш меморије и меморијске картице', true, true, NOW(), NOW()),
        ('category', 3385, 'sr', 'name', 'Веб камере', true, true, NOW(), NOW()),
        ('category', 3390, 'sr', 'name', 'ТВ тјунери', true, true, NOW(), NOW()),

    ('category', 3500, 'sr', 'name', 'Игре, конзоле и програми', true, true, NOW(), NOW()),
            -- Подкатегорије игара, конзола и програма (надређена id=3500)
        ('category', 3510, 'sr', 'name', 'Конзоле за игре', true, true, NOW(), NOW()),
        ('category', 3520, 'sr', 'name', 'Игре за конзоле', true, true, NOW(), NOW()),
        ('category', 3530, 'sr', 'name', 'Рачунарске игре', true, true, NOW(), NOW()),
        ('category', 3540, 'sr', 'name', 'Програми', true, true, NOW(), NOW()),

    ('category', 3600, 'sr', 'name', 'Лаптопови', true, true, NOW(), NOW()),
    ('category', 3700, 'sr', 'name', 'Фото опрема', true, true, NOW(), NOW()),
            -- Подкатегорије фото опреме (надређена id=3700)
        ('category', 3710, 'sr', 'name', 'Опрема и додаци', true, true, NOW(), NOW()),
        ('category', 3720, 'sr', 'name', 'Објективи', true, true, NOW(), NOW()),
        ('category', 3730, 'sr', 'name', 'Компактни фотоапарати', true, true, NOW(), NOW()),
        ('category', 3740, 'sr', 'name', 'Филмски фотоапарати', true, true, NOW(), NOW()),
        ('category', 3750, 'sr', 'name', 'ДСЛР фотоапарати', true, true, NOW(), NOW()),

    ('category', 3800, 'sr', 'name', 'Таблети и е-читачи', true, true, NOW(), NOW()),
            -- Подкатегорије таблета и е-читача (надређена id=3800)
        ('category', 3810, 'sr', 'name', 'Таблети', true, true, NOW(), NOW()),
        ('category', 3820, 'sr', 'name', 'Е-читачи', true, true, NOW(), NOW()),
        ('category', 3830, 'sr', 'name', 'Додаци', true, true, NOW(), NOW()),

    ('category', 3900, 'sr', 'name', 'Канцеларијска опрема и потрошни материјал', true, true, NOW(), NOW()),
            -- Подкатегорије канцеларијске опреме (надређена id=3900)
        ('category', 3910, 'sr', 'name', 'МФП, копир и скенер уређаји', true, true, NOW(), NOW()),
        ('category', 3920, 'sr', 'name', 'Штампачи', true, true, NOW(), NOW()),
        ('category', 3930, 'sr', 'name', 'Канцеларијски материјал', true, true, NOW(), NOW()),
        ('category', 3940, 'sr', 'name', 'УПС, филтери напајања', true, true, NOW(), NOW()),
        ('category', 3950, 'sr', 'name', 'Телефонија', true, true, NOW(), NOW()),
        ('category', 3960, 'sr', 'name', 'Уништавачи папира', true, true, NOW(), NOW()),
        ('category', 3970, 'sr', 'name', 'Потрошни материјал', true, true, NOW(), NOW()),

    ('category', 4100, 'sr', 'name', 'Кућни апарати', true, true, NOW(), NOW()),
        -- Подкатегорије кућних апарата (надређена id=4100)
        ('category', 4110, 'sr', 'name', 'Кухињски апарати', true, true, NOW(), NOW()),
                    -- Подкатегорије кухињских апарата (надређена id=4110)
            ('category', 4111, 'sr', 'name', 'Аспиратори', true, true, NOW(), NOW()),
            ('category', 4112, 'sr', 'name', 'Мали кухињски апарати', true, true, NOW(), NOW()),
            ('category', 4113, 'sr', 'name', 'Микроталасне пећнице', true, true, NOW(), NOW()),
            ('category', 4114, 'sr', 'name', 'Шпорети и рерне', true, true, NOW(), NOW()),
            ('category', 4115, 'sr', 'name', 'Машине за прање судова', true, true, NOW(), NOW()),
            ('category', 4116, 'sr', 'name', 'Фрижидери и замрзивачи', true, true, NOW(), NOW()),

        ('category', 4120, 'sr', 'name', 'Кућни апарати', true, true, NOW(), NOW()),
                    -- Подкатегорије кућних апарата (надређена id=4120)
            ('category', 4121, 'sr', 'name', 'Усисивачи и делови', true, true, NOW(), NOW()),
            ('category', 4122, 'sr', 'name', 'Машине за прање и сушење веша', true, true, NOW(), NOW()),
            ('category', 4123, 'sr', 'name', 'Пегле', true, true, NOW(), NOW()),
            ('category', 4124, 'sr', 'name', 'Опрема за шивење', true, true, NOW(), NOW()),

        ('category', 4130, 'sr', 'name', 'Климатизациона опрема', true, true, NOW(), NOW()),
                    -- Подкатегорије климатизационе опреме (надређена id=4130)
            ('category', 4131, 'sr', 'name', 'Вентилатори', true, true, NOW(), NOW()),
            ('category', 4132, 'sr', 'name', 'Клима уређаји и делови', true, true, NOW(), NOW()),
            ('category', 4133, 'sr', 'name', 'Грејалице', true, true, NOW(), NOW()),
            ('category', 4134, 'sr', 'name', 'Пречишћивачи ваздуха', true, true, NOW(), NOW()),
            ('category', 4135, 'sr', 'name', 'Термометри и метео станице', true, true, NOW(), NOW()),

        ('category', 4140, 'sr', 'name', 'Апарати за личну негу', true, true, NOW(), NOW()),
            -- Подкатегорије апарата за личну негу (надређена id=4140)
            ('category', 4141, 'sr', 'name', 'Бријачи и тримери', true, true, NOW(), NOW()),
            ('category', 4142, 'sr', 'name', 'Машинице за шишање', true, true, NOW(), NOW()),
            ('category', 4143, 'sr', 'name', 'Фенови и уређаји за обликовање', true, true, NOW(), NOW()),
            ('category', 4144, 'sr', 'name', 'Епилатори', true, true, NOW(), NOW()),
            
('category', 5000, 'sr', 'name', 'Све за кућу и стан', true, true, NOW(), NOW()),
    -- Све за кућу и стан (надређена id=5000)
    -- Реновирање и изградња
    ('category', 5100, 'sr', 'name', 'Реновирање и изградња', true, true, NOW(), NOW()),
        -- Подкатегорије реновирања и изградње (надређена id=5100)
        ('category', 5110, 'sr', 'name', 'Врата', true, true, NOW(), NOW()),
        ('category', 5115, 'sr', 'name', 'Алати', true, true, NOW(), NOW()),
        ('category', 5120, 'sr', 'name', 'Камини и грејалице', true, true, NOW(), NOW()),
        ('category', 5125, 'sr', 'name', 'Прозори и балкони', true, true, NOW(), NOW()),
        ('category', 5130, 'sr', 'name', 'Плафони', true, true, NOW(), NOW()),
        ('category', 5135, 'sr', 'name', 'За башту и викендицу', true, true, NOW(), NOW()),
        ('category', 5140, 'sr', 'name', 'Водовод, водоснабдевање и сауна', true, true, NOW(), NOW()),
        ('category', 5145, 'sr', 'name', 'Готове конструкције и брвнаре', true, true, NOW(), NOW()),
        ('category', 5150, 'sr', 'name', 'Капије, ограде и баријере', true, true, NOW(), NOW()),
        ('category', 5155, 'sr', 'name', 'Сигурносни и алармни системи', true, true, NOW(), NOW()),

    ('category', 5200, 'sr', 'name', 'Намештај и ентеријер', true, true, NOW(), NOW()),
        -- Подкатегорије намештаја и ентеријера (надређена id=5200)
        ('category', 5210, 'sr', 'name', 'Кревети, каучи и фотеље', true, true, NOW(), NOW()),
        ('category', 5215, 'sr', 'name', 'Текстил и теписи', true, true, NOW(), NOW()),
        ('category', 5220, 'sr', 'name', 'Осветљење', true, true, NOW(), NOW()),
        ('category', 5225, 'sr', 'name', 'Компјутерски столови и столице', true, true, NOW(), NOW()),
        ('category', 5230, 'sr', 'name', 'Ормари, комоде и полице', true, true, NOW(), NOW()),
        ('category', 5235, 'sr', 'name', 'Кухињске гарнитуре', true, true, NOW(), NOW()),
        ('category', 5240, 'sr', 'name', 'Столови и столице', true, true, NOW(), NOW()),
        ('category', 5250, 'sr', 'name', 'Собне биљке', true, true, NOW(), NOW()),
            -- Подкатегорије собних биљака (надређена id=5250)
            ('category', 5251, 'sr', 'name', 'Декоративно лишће биљака', true, true, NOW(), NOW()),
            ('category', 5252, 'sr', 'name', 'Цветне биљке', true, true, NOW(), NOW()),
            ('category', 5253, 'sr', 'name', 'Палме и фикуси', true, true, NOW(), NOW()),
            ('category', 5254, 'sr', 'name', 'Кактуси и сукуленти', true, true, NOW(), NOW()),

('category', 5300, 'sr', 'name', 'Прехрамбени производи', true, true, NOW(), NOW()),
            -- Подкатегорије прехрамбених производа (надређена id=5300)
        ('category', 5310, 'sr', 'name', 'Чај, кафа, какао', true, true, NOW(), NOW()),
        ('category', 5315, 'sr', 'name', 'Пића', true, true, NOW(), NOW()),
        ('category', 5320, 'sr', 'name', 'Риба, морски плодови, кавијар', true, true, NOW(), NOW()),
        ('category', 5325, 'sr', 'name', 'Месо, живина, изнутрице', true, true, NOW(), NOW()),
        ('category', 5330, 'sr', 'name', 'Слаткиши', true, true, NOW(), NOW()),
        ('category', 5340, 'sr', 'name', 'Ракија и вино', true, true, NOW(), NOW()),
                    -- Подкатегорије ракије и вина (надређена id=5340)
            ('category', 5341, 'sr', 'name', 'Шљивовица', true, true, NOW(), NOW()),
            ('category', 5342, 'sr', 'name', 'Лозовача', true, true, NOW(), NOW()),
            ('category', 5343, 'sr', 'name', 'Воћна ракија', true, true, NOW(), NOW()),
            ('category', 5344, 'sr', 'name', 'Домаће вино', true, true, NOW(), NOW()),

        ('category', 5350, 'sr', 'name', 'Домаћи сиреви', true, true, NOW(), NOW()),
        ('category', 5360, 'sr', 'name', 'Кајмак', true, true, NOW(), NOW()),
        ('category', 5370, 'sr', 'name', 'Ајвар', true, true, NOW(), NOW()),

    ('category', 5400, 'sr', 'name', 'Посуђе и кухињски прибор', true, true, NOW(), NOW()),
        -- Подкатегорије посуђа и кухињског прибора (надређена id=5400)
        ('category', 5405, 'sr', 'name', 'Посуђе', true, true, NOW(), NOW()),
        ('category', 5410, 'sr', 'name', 'Кухињски прибор', true, true, NOW(), NOW()),
        ('category', 5415, 'sr', 'name', 'Прибор за постављање стола', true, true, NOW(), NOW()),
        ('category', 5420, 'sr', 'name', 'Прибор за кување', true, true, NOW(), NOW()),
        ('category', 5425, 'sr', 'name', 'Чување хране', true, true, NOW(), NOW()),
        ('category', 5430, 'sr', 'name', 'Припрема напитака', true, true, NOW(), NOW()),
        ('category', 5435, 'sr', 'name', 'Кућна хемија', true, true, NOW(), NOW()),

('category', 6000, 'sr', 'name', 'Све за башту', true, true, NOW(), NOW()),
    -- Све за башту (надређена id=6000, подкатегорије 51-60)
    ('category', 6050, 'sr', 'name', 'Баштенски намештај', true, true, NOW(), NOW()),
    ('category', 6100, 'sr', 'name', 'Баштенски алати', true, true, NOW(), NOW()),
    ('category', 6200, 'sr', 'name', 'Семе и расад', true, true, NOW(), NOW()),
    ('category', 6250, 'sr', 'name', 'Роштиљ и додаци', true, true, NOW(), NOW()),
    ('category', 6300, 'sr', 'name', 'Базени и опрема', true, true, NOW(), NOW()),
    ('category', 6350, 'sr', 'name', 'Системи за наводњавање', true, true, NOW(), NOW()),
    ('category', 6400, 'sr', 'name', 'Компостирање', true, true, NOW(), NOW()),
    ('category', 6450, 'sr', 'name', 'Стакленици и расадници', true, true, NOW(), NOW()),
    ('category', 6500, 'sr', 'name', 'Ђубрива и земља', true, true, NOW(), NOW()),
    ('category', 6550, 'sr', 'name', 'Осветљење', true, true, NOW(), NOW()),
    ('category', 6600, 'sr', 'name', 'Уређење ентеријера', true, true, NOW(), NOW()),
    ('category', 6650, 'sr', 'name', 'Биљке и семена', true, true, NOW(), NOW()),
    ('category', 6700, 'sr', 'name', 'Башта и повртњак', true, true, NOW(), NOW()),
    ('category', 6750, 'sr', 'name', 'Баштенске биљке', true, true, NOW(), NOW()),
            -- Подкатегорије баштенских биљака (надређена id=6750)
        ('category', 6751, 'sr', 'name', 'Декоративно грмље и дрвеће', true, true, NOW(), NOW()),
        ('category', 6752, 'sr', 'name', 'Четинари', true, true, NOW(), NOW()),
        ('category', 6753, 'sr', 'name', 'Вишегодишње биљке', true, true, NOW(), NOW()),
        ('category', 6754, 'sr', 'name', 'Воћке', true, true, NOW(), NOW()),
        ('category', 6755, 'sr', 'name', 'Травњак', true, true, NOW(), NOW()),
        ('category', 6756, 'sr', 'name', 'Зелениш и зачини', true, true, NOW(), NOW()),

    ('category', 6850, 'sr', 'name', 'Семена, луковице, гомољи', true, true, NOW(), NOW()),
    ('category', 6900, 'sr', 'name', 'Средства за негу биљака', true, true, NOW(), NOW()),
            -- Подкатегорије средстава за негу биљака (надређена id=6900)
        ('category', 6901, 'sr', 'name', 'Земље и супстрати', true, true, NOW(), NOW()),
        ('category', 6902, 'sr', 'name', 'Ђубрива', true, true, NOW(), NOW()),
        ('category', 6903, 'sr', 'name', 'Средства против штеточина и корова', true, true, NOW(), NOW()),
        ('category', 6904, 'sr', 'name', 'Саксије и жардињере', true, true, NOW(), NOW()),
        ('category', 6905, 'sr', 'name', 'Фито лампе', true, true, NOW(), NOW()),
        ('category', 6906, 'sr', 'name', 'Мерачи влаге', true, true, NOW(), NOW()),

    ('category', 6950, 'sr', 'name', 'Панои и вештачке биљке', true, true, NOW(), NOW()),

('category', 7000, 'sr', 'name', 'Хоби и рекреација', true, true, NOW(), NOW()),
        -- Хоби и рекреација (надређена id=7000)
    ('category', 7050, 'sr', 'name', 'Музички инструменти', true, true, NOW(), NOW()),
        -- Музички инструменти (надређена id=7050)
        ('category', 7055, 'sr', 'name', 'Жичани инструменти', true, true, NOW(), NOW()),
        ('category', 7060, 'sr', 'name', 'Клавири и клавијатуре', true, true, NOW(), NOW()),
        ('category', 7065, 'sr', 'name', 'Удараљке', true, true, NOW(), NOW()),
        ('category', 7070, 'sr', 'name', 'Дувачки инструменти', true, true, NOW(), NOW()),
        ('category', 7075, 'sr', 'name', 'Хармонике', true, true, NOW(), NOW()),
        ('category', 7080, 'sr', 'name', 'Аудио опрема', true, true, NOW(), NOW()),
        ('category', 7085, 'sr', 'name', 'Додаци за инструменте', true, true, NOW(), NOW()),

    ('category', 7100, 'sr', 'name', 'Књиге и часописи', true, true, NOW(), NOW()),
            -- Подкатегорије књига и часописа (надређена id=7100)
        ('category', 7105, 'sr', 'name', 'Часописи, новине, брошуре', true, true, NOW(), NOW()),
        ('category', 7115, 'sr', 'name', 'Књиге', true, true, NOW(), NOW()),
        ('category', 7130, 'sr', 'name', 'Уџбеници', true, true, NOW(), NOW()),
    ('category', 7150, 'sr', 'name', 'Спортски реквизити', true, true, NOW(), NOW()),
    ('category', 7250, 'sr', 'name', 'Колекционарство', true, true, NOW(), NOW()),
            -- Подкатегорије колекционарства (надређена id=7250)
        ('category', 7251, 'sr', 'name', 'Новчанице', true, true, NOW(), NOW()),
        ('category', 7252, 'sr', 'name', 'Карте', true, true, NOW(), NOW()),
        ('category', 7253, 'sr', 'name', 'Ствари познатих, аутограми', true, true, NOW(), NOW()),
        ('category', 7254, 'sr', 'name', 'Војне ствари', true, true, NOW(), NOW()),
        ('category', 7255, 'sr', 'name', 'Грамофонске плоче', true, true, NOW(), NOW()),
        ('category', 7256, 'sr', 'name', 'Документа', true, true, NOW(), NOW()),
        ('category', 7257, 'sr', 'name', 'Жетони, медаље, значке', true, true, NOW(), NOW()),
        ('category', 7258, 'sr', 'name', 'Игре', true, true, NOW(), NOW()),
        ('category', 7259, 'sr', 'name', 'Календари', true, true, NOW(), NOW()),
        ('category', 7261, 'sr', 'name', 'Слике', true, true, NOW(), NOW()),
        ('category', 7262, 'sr', 'name', 'Поштанске марке', true, true, NOW(), NOW()),
        ('category', 7263, 'sr', 'name', 'Модели', true, true, NOW(), NOW()),
        ('category', 7264, 'sr', 'name', 'Новчићи', true, true, NOW(), NOW()),

    ('category', 7300, 'sr', 'name', 'Уметнички предмети', true, true, NOW(), NOW()),
    ('category', 7400, 'sr', 'name', 'Бицикли', true, true, NOW(), NOW()),
            -- Подкатегорије бицикала (надређена id=7400)
        ('category', 7410, 'sr', 'name', 'БМХ', true, true, NOW(), NOW()),
        ('category', 7415, 'sr', 'name', 'Градски бицикли', true, true, NOW(), NOW()),
        ('category', 7420, 'sr', 'name', 'Бицикли за друм', true, true, NOW(), NOW()),
        ('category', 7425, 'sr', 'name', 'Дечији бицикли', true, true, NOW(), NOW()),
        ('category', 7430, 'sr', 'name', 'Брдски бицикли', true, true, NOW(), NOW()),
        ('category', 7435, 'sr', 'name', 'Делови и додаци', true, true, NOW(), NOW()),

    ('category', 7500, 'sr', 'name', 'Лов и риболов', true, true, NOW(), NOW()),
        -- Подкатегорије лова и риболова (надређена id=7500)
        ('category', 7510, 'sr', 'name', 'Ножеви, мултиалати, секире', true, true, NOW(), NOW()),
        ('category', 7520, 'sr', 'name', 'Лов', true, true, NOW(), NOW()),
                    -- Подкатегорије лова (надређена id=7520)
            ('category', 7521, 'sr', 'name', 'Оптички нишани', true, true, NOW(), NOW()),
            ('category', 7522, 'sr', 'name', 'Додаци за оптичке нишане', true, true, NOW(), NOW()),
            ('category', 7523, 'sr', 'name', 'Монокулари, двогледи, даљиномери', true, true, NOW(), NOW()),

        ('category', 7530, 'sr', 'name', 'Риболов', true, true, NOW(), NOW()),
            -- Подкатегорије риболова (надређена id=7530)
            ('category', 7531, 'sr', 'name', 'Штапови, машинице', true, true, NOW(), NOW()),
            ('category', 7551, 'sr', 'name', 'Мамци и прибор', true, true, NOW(), NOW()),
            ('category', 7571, 'sr', 'name', 'Сонари и опрема', true, true, NOW(), NOW()),

    ('category', 7650, 'sr', 'name', 'Камповање', true, true, NOW(), NOW()),
    ('category', 7700, 'sr', 'name', 'Антиквитети', true, true, NOW(), NOW()),
    ('category', 7750, 'sr', 'name', 'Карте, догађаји и путовања', true, true, NOW(), NOW()),
        -- Подкатегорије карата, догађаја и путовања (надређена id=7750)
        ('category', 7751, 'sr', 'name', 'Карте, купони', true, true, NOW(), NOW()),
        ('category', 7752, 'sr', 'name', 'Концерти', true, true, NOW(), NOW()),
        ('category', 7753, 'sr', 'name', 'Путовања', true, true, NOW(), NOW()),
        ('category', 7754, 'sr', 'name', 'Спорт', true, true, NOW(), NOW()),
        ('category', 7755, 'sr', 'name', 'Позориште, опера, балет', true, true, NOW(), NOW()),
        ('category', 7756, 'sr', 'name', 'Циркус, биоскоп', true, true, NOW(), NOW()),
        ('category', 7758, 'sr', 'name', 'Шоу, мјузикл', true, true, NOW(), NOW()),

    ('category', 7800, 'sr', 'name', 'Спорт', true, true, NOW(), NOW()),
            -- Подкатегорије спорта (надређена id=7800)
        ('category', 7805, 'sr', 'name', 'Билијар и куглање', true, true, NOW(), NOW()),
        ('category', 7810, 'sr', 'name', 'Роњење и водени спортови', true, true, NOW(), NOW()),
        ('category', 7815, 'sr', 'name', 'Борилачки спортови', true, true, NOW(), NOW()),
        ('category', 7820, 'sr', 'name', 'Зимски спортови', true, true, NOW(), NOW()),
        ('category', 7825, 'sr', 'name', 'Игре са лоптом', true, true, NOW(), NOW()),
        ('category', 7830, 'sr', 'name', 'Друштвене игре', true, true, NOW(), NOW()),
        ('category', 7835, 'sr', 'name', 'Пејнтбол и ерсофт', true, true, NOW(), NOW()),
        ('category', 7840, 'sr', 'name', 'Ролери и скејтборд', true, true, NOW(), NOW()),
        ('category', 7845, 'sr', 'name', 'Тенис, бадминтон, пинг-понг', true, true, NOW(), NOW()),
        ('category', 7850, 'sr', 'name', 'Туризам и боравак у природи', true, true, NOW(), NOW()),
        ('category', 7855, 'sr', 'name', 'Фитнес и тренажери', true, true, NOW(), NOW()),
        ('category', 7860, 'sr', 'name', 'Спортска исхрана', true, true, NOW(), NOW()),

    -- Традиционални занати и сувенири
    ('category', 7865, 'sr', 'name', 'Народни занати и рукотворине', true, true, NOW(), NOW()),
        -- Подкатегорије народних заната (надређена id=7865)
        ('category', 7866, 'sr', 'name', 'Опанци', true, true, NOW(), NOW()),
        ('category', 7867, 'sr', 'name', 'Керамика', true, true, NOW(), NOW()),
        ('category', 7868, 'sr', 'name', 'Вез', true, true, NOW(), NOW()),
        ('category', 7869, 'sr', 'name', 'Ткање', true, true, NOW(), NOW()),
        ('category', 7870, 'sr', 'name', 'Народни инструменти', true, true, NOW(), NOW()),
        ('category', 7871, 'sr', 'name', 'Дрводеља', true, true, NOW(), NOW()),

    -- Пољопривредне категорије
    ('category', 7900, 'sr', 'name', 'Пчеларство', true, true, NOW(), NOW()),
        -- Подкатегорије пчеларства (надређена id=7900)
        ('category', 7910, 'sr', 'name', 'Мед', true, true, NOW(), NOW()),
        ('category', 7920, 'sr', 'name', 'Пчелињи восак', true, true, NOW(), NOW()),
        ('category', 7930, 'sr', 'name', 'Прополис', true, true, NOW(), NOW()),
        ('category', 7935, 'sr', 'name', 'Пчеларски прибор', true, true, NOW(), NOW()),
        ('category', 7945, 'sr', 'name', 'Пчеле', true, true, NOW(), NOW()),

    -- Туристичке услуге
    ('category', 7950, 'sr', 'name', 'Сеоски туризам', true, true, NOW(), NOW()),
            -- Подкатегорије сеоског туризма (надређена id=7950)
        ('category', 7951, 'sr', 'name', 'Етно-села', true, true, NOW(), NOW()),
        ('category', 7952, 'sr', 'name', 'Винске туре', true, true, NOW(), NOW()),
        ('category', 7953, 'sr', 'name', 'Агро туризам', true, true, NOW(), NOW()),
        ('category', 7954, 'sr', 'name', 'Планински туризам', true, true, NOW(), NOW()),

('category', 8000, 'sr', 'name', 'Животиње', true, true, NOW(), NOW()),
    -- Категорије животиња (надређена id=8000)
    ('category', 8050, 'sr', 'name', 'Пси', true, true, NOW(), NOW()),
    ('category', 8100, 'sr', 'name', 'Мачке', true, true, NOW(), NOW()),
    ('category', 8150, 'sr', 'name', 'Птице', true, true, NOW(), NOW()),
    ('category', 8200, 'sr', 'name', 'Акваријум', true, true, NOW(), NOW()),
            -- Подкатегорије акваријума (надређена id=8200)
        ('category', 8205, 'sr', 'name', 'Акваријуми', true, true, NOW(), NOW()),
        ('category', 8210, 'sr', 'name', 'Рибе', true, true, NOW(), NOW()),
        ('category', 8215, 'sr', 'name', 'Друге акваријумске животиње', true, true, NOW(), NOW()),
        ('category', 8220, 'sr', 'name', 'Опрема', true, true, NOW(), NOW()),
        ('category', 8225, 'sr', 'name', 'Биљке', true, true, NOW(), NOW()),
        ('category', 8230, 'sr', 'name', 'Акваријумски намештај', true, true, NOW(), NOW()),
        ('category', 8235, 'sr', 'name', 'Морска акваристика', true, true, NOW(), NOW()),

    ('category', 8250, 'sr', 'name', 'Друге животиње', true, true, NOW(), NOW()),
        -- Подкатегорије других животиња (надређена id=8250)
        ('category', 8251, 'sr', 'name', 'Водоземци', true, true, NOW(), NOW()),
        ('category', 8252, 'sr', 'name', 'Глодари', true, true, NOW(), NOW()),
        ('category', 8253, 'sr', 'name', 'Зечеви', true, true, NOW(), NOW()),
        ('category', 8254, 'sr', 'name', 'Коњи', true, true, NOW(), NOW()),
        ('category', 8255, 'sr', 'name', 'Гмизавци', true, true, NOW(), NOW()),
        ('category', 8256, 'sr', 'name', 'Домаће животиње', true, true, NOW(), NOW()),
        ('category', 8257, 'sr', 'name', 'Живина', true, true, NOW(), NOW()),
        ('category', 8258, 'sr', 'name', 'Производи за животиње', true, true, NOW(), NOW()),

('category', 8500, 'sr', 'name', 'Готов бизнис и опрема', true, true, NOW(), NOW()),


('category', 9000, 'sr', 'name', 'Посао', true, true, NOW(), NOW()),
    ('category', 9050, 'sr', 'name', 'Слободна радна места', true, true, NOW(), NOW()),
    ('category', 9100, 'sr', 'name', 'Биографије', true, true, NOW(), NOW()),
    ('category', 9150, 'sr', 'name', 'Рад на даљину', true, true, NOW(), NOW()),
    ('category', 9200, 'sr', 'name', 'Партнерство и сарадња', true, true, NOW(), NOW()),
    ('category', 9250, 'sr', 'name', 'Обука и пракса', true, true, NOW(), NOW()),
    ('category', 9300, 'sr', 'name', 'Сезонски посао', true, true, NOW(), NOW()),
        -- Подкатегорије сезонског посла (надређена id=9300)
        ('category', 9310, 'sr', 'name', 'Берба', true, true, NOW(), NOW()),
        ('category', 9315, 'sr', 'name', 'Рад у винограду', true, true, NOW(), NOW()),
        ('category', 9320, 'sr', 'name', 'Сезонски грађевински радови', true, true, NOW(), NOW()),

('category', 9500, 'sr', 'name', 'Одећа, обућа, додаци', true, true, NOW(), NOW()),

-- Дечије ствари и играчке (нова категорија)
('category', 9700, 'sr', 'name', 'Дечије ствари и играчке', true, true, NOW(), NOW()),
    -- Подкатегорије дечијих ствари и играчака (надређена id=9700)
    ('category', 9705, 'sr', 'name', 'Колица за бебе', true, true, NOW(), NOW()),
    ('category', 9710, 'sr', 'name', 'Дечији намештај', true, true, NOW(), NOW()),
    ('category', 9715, 'sr', 'name', 'Бицикли и тротинети', true, true, NOW(), NOW()),
    ('category', 9720, 'sr', 'name', 'Производи за храњење', true, true, NOW(), NOW()),
    ('category', 9725, 'sr', 'name', 'Ауто-седишта', true, true, NOW(), NOW()),
    ('category', 9730, 'sr', 'name', 'Играчке', true, true, NOW(), NOW()),
    ('category', 9735, 'sr', 'name', 'Постељина', true, true, NOW(), NOW()),
    ('category', 9740, 'sr', 'name', 'Производи за купање', true, true, NOW(), NOW()),
    ('category', 9745, 'sr', 'name', 'Школски прибор', true, true, NOW(), NOW()),
    ('category', 9750, 'sr', 'name', 'Дечија одећа и обућа, додаци', true, true, NOW(), NOW()),

('category', 9999, 'sr', 'name', 'Остало', true, true, NOW(), NOW());

-- Ажурирање sequence-а за translations
SELECT setval('translations_id_seq', (SELECT MAX(id) FROM translations), true);


-- Insert marketplace listings
INSERT INTO marketplace_listings (id, user_id, category_id, title, description, price, condition, status, location, latitude, longitude, address_city, address_country, views_count, created_at, updated_at, show_on_map, original_language) VALUES
(8, 2, 2100, 'Toyota Corolla 2018', 'Продајем Toyota Corolla 2018 годиште, 80.000 км, одлично стање. Први власник, редовно одржавање, сва документација доступна.', 1150000.00, 'used', 'active', 'Нови Сад, Србија', 45.26710000, 19.83350000, 'Нови Сад', 'Србија', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'sr'),
(9, 3, 3100, 'mobile Samsung Galaxy S21', 'Selling Samsung Galaxy S21, 256GB, Deep Purple. Perfect condition, complete set with original box and accessories. AppleCare+ until 2024.', 120000.00, 'used', 'active', 'Novi Sad, Serbia', 45.25510000, 19.84520000, 'Novi Sad', 'Serbia', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'en'),
(10, 4, 3320, 'Игровой компьютер RTX 4080', 'Продаю мощный игровой ПК: Intel Core i9-13900K, RTX 4080, 64GB RAM, 2TB NVMe SSD. Идеален для любых игр и тяжелых задач.', 350000.00, 'used', 'active', 'Нови-Сад, Сербия', 45.25410000, 19.84010000, 'Нови-Сад', 'Сербия', 0, '2025-02-07 07:13:52.973909', '2025-02-07 07:13:52.973909', true, 'ru'),
(12, 2, 2100, 'автомобиль Toyota Corolla 2018', 'Продаю Toyota Corolla 2018 года, 80.000 км, отличное состояние. Первый владелец, регулярное обслуживание, вся документация в наличии.', 1475000.00, 'used', 'active', 'Косте Мајинског 4, Ветерник, Сербия', 45.24755670, 19.76878366, 'Ветерник', 'Сербия', 0, '2025-02-07 17:33:27.680035', '2025-02-07 17:40:23.957971', true, 'ru');

SELECT setval('marketplace_listings_id_seq', 12, true);

-- Insert marketplace images
INSERT INTO marketplace_images (id, listing_id, file_path, file_name, file_size, content_type, is_main, created_at) VALUES
(15, 8, 'toyota_1.jpg', 'toyota_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(16, 8, 'toyota_2.jpg', 'toyota_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(17, 9, 'galaxy_s21_1.jpg', 'galaxy_s21_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(18, 9, 'galaxy_s21_2.jpg', 'galaxy_s21_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(19, 10, 'gaming_pc_1.jpg', 'gaming_pc_1.jpg', 1024, 'image/jpeg', true, '2025-02-07 07:13:52.973909'),
(20, 10, 'gaming_pc_2.jpg', 'gaming_pc_2.jpg', 1024, 'image/jpeg', false, '2025-02-07 07:13:52.973909'),
(21, 12, 'toyota_1.jpg', 'toyota_1.jpg', 454842, 'image/jpeg', true, '2025-02-07 17:35:09.579393'),
(22, 12, 'toyota_2.jpg', 'toyota_2.jpg', 398035, 'image/jpeg', true, '2025-02-07 17:40:24.397595');

SELECT setval('marketplace_images_id_seq', 22, true);

-- Insert reviews and related data
INSERT INTO reviews (id, user_id, entity_type, entity_id, rating, comment, pros, cons, photos, likes_count, is_verified_purchase, status, created_at, updated_at, helpful_votes, not_helpful_votes, original_language) VALUES
(1, 2, 'listing', 8, 5, 'норм', NULL, NULL, NULL, 0, true, 'published', '2025-02-07 07:47:17.001726', '2025-02-07 14:25:23.586871', 0, 1, 'ru');

SELECT setval('reviews_id_seq', 1, true);

INSERT INTO review_responses (id, review_id, user_id, response, created_at, updated_at) VALUES
(1, 1, 3, 'ok', '2025-02-07 07:49:14.935308', '2025-02-07 07:49:14.935308');

SELECT setval('review_responses_id_seq', 1, true);

INSERT INTO review_votes (review_id, user_id, vote_type, created_at) VALUES
(1, 3, 'not_helpful', '2025-02-07 07:48:11.709016');
