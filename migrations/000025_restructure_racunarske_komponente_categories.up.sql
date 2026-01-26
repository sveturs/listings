-- ==========================================
-- ФАЗА 0: Переструктурирование категории "Računarske komponente"
-- Дата: 2026-01-26
-- Описание: Удаление 20 конкретных подкатегорий, создание 9 общих категорий
-- КРИТИЧНО: Эта миграция должна быть применена ДО миграции 000023 (привязка атрибутов)
-- ==========================================

BEGIN;

-- ==========================================
-- ШАГ 1: ПРОВЕРКА - Убедиться что нет листингов в старых категориях
-- ==========================================

DO $$
DECLARE
    listings_count INTEGER;
    parent_uuid uuid := 'fe64130f-deea-4767-a937-6f9d584a4395'::uuid; -- racunarske-komponente
BEGIN
    -- Подсчитать листинги в старых подкатегориях
    SELECT COUNT(*) INTO listings_count
    FROM listings l
    JOIN categories c ON l.category_id = c.id
    WHERE c.parent_id = parent_uuid;

    IF listings_count > 0 THEN
        RAISE WARNING '⚠️  Found % listings in old subcategories!', listings_count;
        RAISE WARNING '⚠️  These listings will be orphaned after category deletion.';
        RAISE WARNING '⚠️  Please migrate listings to new categories first, or set category_id = parent category.';
        -- НЕ останавливаем миграцию - просто предупреждаем
        -- Листинги автоматически станут прямыми детьми "Računarske komponente"
    ELSE
        RAISE NOTICE '✅ No listings found in old subcategories. Safe to proceed.';
    END IF;
END $$;

-- ==========================================
-- ШАГ 2: УДАЛЕНИЕ - Старые 20 конкретных подкатегорий
-- ==========================================

DO $$
DECLARE
    parent_uuid uuid := 'fe64130f-deea-4767-a937-6f9d584a4395';
    deleted_attrs INTEGER;
    deleted_cats INTEGER;
BEGIN
    RAISE NOTICE '';
    RAISE NOTICE '🗑️  Deleting old concrete subcategories...';
    RAISE NOTICE '';

    -- Удалить связи атрибутов (category_attributes)
    DELETE FROM category_attributes WHERE category_id IN (
        SELECT id FROM categories WHERE parent_id = parent_uuid
    );
    GET DIAGNOSTICS deleted_attrs = ROW_COUNT;

    -- Удалить связи вариативных атрибутов (category_variant_attributes)
    -- ВАЖНО: category_variant_attributes.category_id - varchar(36), а categories.id - uuid
    DELETE FROM category_variant_attributes WHERE category_id IN (
        SELECT id::varchar FROM categories WHERE parent_id = parent_uuid
    );

    -- Удалить старые подкатегории
    DELETE FROM categories WHERE parent_id = parent_uuid;
    GET DIAGNOSTICS deleted_cats = ROW_COUNT;

    RAISE NOTICE '✅ Deleted % old subcategories', deleted_cats;
    RAISE NOTICE '✅ Deleted % attribute links', deleted_attrs;
    RAISE NOTICE '';
END $$;

-- ==========================================
-- ШАГ 3: СОЗДАНИЕ - 9 новых общих подкатегорий
-- ==========================================

DO $$
DECLARE
    parent_uuid uuid := 'fe64130f-deea-4767-a937-6f9d584a4395'::uuid; -- racunarske-komponente
BEGIN
    RAISE NOTICE '📦 Creating 9 new general subcategories...';
    RAISE NOTICE '';

    -- 1. Grafičke kartice (Видеокарты)
    INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, meta_keywords, is_active, icon)
    VALUES (
        'graficke-kartice',
        parent_uuid,
        3,
        'elektronika/racunarske-komponente/graficke-kartice',
        10,
        '{"en": "Graphics Cards", "ru": "Видеокарты", "sr": "Grafičke kartice"}',
        '{"en": "NVIDIA RTX, AMD Radeon graphics cards for gaming and professional use", "ru": "Видеокарты NVIDIA RTX, AMD Radeon для игр и профессиональной работы", "sr": "Sve grafičke kartice NVIDIA RTX, AMD Radeon za gaming i profesionalnu upotrebu"}',
        '{"en": "Graphics Cards | Vondi", "ru": "Видеокарты | Vondi", "sr": "Grafičke kartice | Vondi"}',
        '{"en": "Buy graphics cards online - RTX 4090, RTX 4080, RX 7900 XTX", "ru": "Купить видеокарты онлайн - RTX 4090, RTX 4080, RX 7900 XTX", "sr": "Kupite grafičke kartice online - RTX 4090, RTX 4080, RX 7900 XTX"}',
        '{"en": "graphics cards, gpu, rtx 4090, rtx 4080, rtx 4070, rtx 3090, rtx 3080, rx 7900 xtx, rx 7900 xt, rx 6900 xt, nvidia graphics, amd graphics, gaming gpu, video card, graphics card price, buy graphics card, gpu deals", "ru": "видеокарты, gpu, ртх 4090, ртх 4080, ртх 4070, ртх 3090, видеокарта nvidia, видеокарта amd, игровая видеокарта, купить видеокарту, цена видеокарты", "sr": "grafičke kartice, gpu, rtx 4090, rtx 4080, rtx 4070, rtx 3090, amd radeon, nvidia grafika, gejming gpu, video kartica, cena grafičke, kupiti grafičku"}',
        true,
        '🎮'
    );
    RAISE NOTICE '  ✅ Created: graficke-kartice';

    -- 2. Procesori (Процессоры)
    INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, meta_keywords, is_active, icon)
    VALUES (
        'procesori',
        parent_uuid,
        3,
        'elektronika/racunarske-komponente/procesori',
        20,
        '{"en": "Processors (CPU)", "ru": "Процессоры", "sr": "Procesori"}',
        '{"en": "Intel Core, AMD Ryzen processors for desktop PCs", "ru": "Процессоры Intel Core, AMD Ryzen для настольных компьютеров", "sr": "Intel Core, AMD Ryzen procesori za desktop računare"}',
        '{"en": "Processors (CPU) | Vondi", "ru": "Процессоры | Vondi", "sr": "Procesori | Vondi"}',
        '{"en": "Buy CPU online - Intel Core i9, i7, AMD Ryzen 9, Ryzen 7", "ru": "Купить процессор онлайн - Intel Core i9, i7, AMD Ryzen 9, Ryzen 7", "sr": "Kupite procesor online - Intel Core i9, i7, AMD Ryzen 9, Ryzen 7"}',
        '{"en": "cpu, processor, intel core i9, intel core i7, intel core i5, amd ryzen 9, amd ryzen 7, amd ryzen 5, desktop cpu, gaming cpu, cpu price, buy cpu, processor deals", "ru": "процессор, cpu, intel core i9, intel core i7, amd ryzen 9, amd ryzen 7, купить процессор, цена процессора, игровой процессор", "sr": "procesor, cpu, intel core i9, intel core i7, amd ryzen 9, amd ryzen 7, desktop procesor, gejming procesor, cena procesora, kupiti procesor"}',
        true,
        '⚙️'
    );
    RAISE NOTICE '  ✅ Created: procesori';

    -- 3. RAM memorija (Оперативная память)
    INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, meta_keywords, is_active, icon)
    VALUES (
        'ram-memorija',
        parent_uuid,
        3,
        'elektronika/racunarske-komponente/ram-memorija',
        30,
        '{"en": "RAM Memory", "ru": "Оперативная память", "sr": "RAM memorija"}',
        '{"en": "DDR4, DDR5 RAM modules for desktop computers", "ru": "Модули оперативной памяти DDR4, DDR5 для настольных компьютеров", "sr": "DDR4, DDR5 RAM memorija za desktop računare"}',
        '{"en": "RAM Memory | Vondi", "ru": "Оперативная память | Vondi", "sr": "RAM memorija | Vondi"}',
        '{"en": "Buy RAM online - DDR5, DDR4 8GB, 16GB, 32GB, 64GB", "ru": "Купить оперативную память онлайн - DDR5, DDR4 8ГБ, 16ГБ, 32ГБ, 64ГБ", "sr": "Kupite RAM online - DDR5, DDR4 8GB, 16GB, 32GB, 64GB"}',
        '{"en": "ram memory, ddr5, ddr4, 8gb ram, 16gb ram, 32gb ram, 64gb ram, corsair ram, gskill ram, kingston ram, desktop memory, gaming ram, ram price, buy ram", "ru": "оперативная память, озу, ddr5, ddr4, 8гб памяти, 16гб памяти, 32гб памяти, купить память, цена памяти", "sr": "ram memorija, ddr5, ddr4, 8gb ram, 16gb ram, 32gb ram, corsair ram, gskill ram, desktop memorija, gejming ram, cena ram, kupiti ram"}',
        true,
        '🧠'
    );
    RAISE NOTICE '  ✅ Created: ram-memorija';

    -- 4. SSD nakopitelji (SSD накопители)
    INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, meta_keywords, is_active, icon)
    VALUES (
        'ssd-nakopitelji',
        parent_uuid,
        3,
        'elektronika/racunarske-komponente/ssd-nakopitelji',
        40,
        '{"en": "SSD Storage", "ru": "SSD накопители", "sr": "SSD nakopitelji"}',
        '{"en": "NVMe, SATA SSD drives - fast storage for your PC", "ru": "NVMe, SATA SSD диски - быстрая память для вашего ПК", "sr": "NVMe, SATA SSD diskovi - brzo skladište za vaš računar"}',
        '{"en": "SSD Storage | Vondi", "ru": "SSD накопители | Vondi", "sr": "SSD nakopitelji | Vondi"}',
        '{"en": "Buy SSD online - NVMe, SATA 500GB, 1TB, 2TB, 4TB", "ru": "Купить SSD онлайн - NVMe, SATA 500ГБ, 1ТБ, 2ТБ, 4ТБ", "sr": "Kupite SSD online - NVMe, SATA 500GB, 1TB, 2TB, 4TB"}',
        '{"en": "ssd, nvme ssd, sata ssd, m.2 ssd, 500gb ssd, 1tb ssd, 2tb ssd, samsung ssd, crucial ssd, wd ssd, fast storage, ssd price, buy ssd", "ru": "ссд, nvme ssd, sata ssd, м.2 ssd, 500гб ссд, 1тб ссд, 2тб ссд, купить ssd, цена ssd", "sr": "ssd, nvme ssd, sata ssd, m.2 ssd, 500gb ssd, 1tb ssd, samsung ssd, crucial ssd, brzo skladište, cena ssd, kupiti ssd"}',
        true,
        '💿'
    );
    RAISE NOTICE '  ✅ Created: ssd-nakopitelji';

    -- 5. HDD nakopitelji (HDD накопители)
    INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, meta_keywords, is_active, icon)
    VALUES (
        'hdd-nakopitelji',
        parent_uuid,
        3,
        'elektronika/racunarske-komponente/hdd-nakopitelji',
        50,
        '{"en": "HDD Storage", "ru": "HDD накопители", "sr": "HDD nakopitelji"}',
        '{"en": "Hard disk drives - large capacity storage", "ru": "Жёсткие диски - большая ёмкость для хранения данных", "sr": "Hard diskovi - veliko skladište podataka"}',
        '{"en": "HDD Storage | Vondi", "ru": "HDD накопители | Vondi", "sr": "HDD nakopitelji | Vondi"}',
        '{"en": "Buy HDD online - 1TB, 2TB, 4TB, 6TB, 8TB hard drives", "ru": "Купить жёсткие диски онлайн - 1ТБ, 2ТБ, 4ТБ, 6ТБ, 8ТБ", "sr": "Kupite HDD online - hard diskovi 1TB, 2TB, 4TB, 6TB, 8TB"}',
        '{"en": "hdd, hard drive, 1tb hdd, 2tb hdd, 4tb hdd, 8tb hdd, wd hdd, seagate hdd, toshiba hdd, storage, hdd price, buy hdd", "ru": "жёсткий диск, хдд, 1тб хдд, 2тб хдд, 4тб хдд, купить hdd, цена hdd", "sr": "hdd, hard disk, 1tb hdd, 2tb hdd, 4tb hdd, wd hdd, seagate hdd, skladište, cena hdd, kupiti hdd"}',
        true,
        '💾'
    );
    RAISE NOTICE '  ✅ Created: hdd-nakopitelji';

    -- 6. Matične ploče (Материнские платы)
    INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, meta_keywords, is_active, icon)
    VALUES (
        'maticne-ploce',
        parent_uuid,
        3,
        'elektronika/racunarske-komponente/maticne-ploce',
        60,
        '{"en": "Motherboards", "ru": "Материнские платы", "sr": "Matične ploče"}',
        '{"en": "Intel, AMD motherboards for desktop PCs", "ru": "Материнские платы Intel, AMD для настольных ПК", "sr": "Intel, AMD matične ploče za desktop računare"}',
        '{"en": "Motherboards | Vondi", "ru": "Материнские платы | Vondi", "sr": "Matične ploče | Vondi"}',
        '{"en": "Buy motherboards online - Z790, B760, X670, B650", "ru": "Купить материнские платы онлайн - Z790, B760, X670, B650", "sr": "Kupite matične ploče online - Z790, B760, X670, B650"}',
        '{"en": "motherboard, mainboard, z790, b760, x670, b650, intel motherboard, amd motherboard, atx motherboard, micro atx, gaming motherboard, motherboard price, buy motherboard", "ru": "материнская плата, z790, b760, x670, b650, материнка intel, материнка amd, atx плата, купить материнку, цена материнки", "sr": "matična ploča, z790, b760, x670, b650, intel ploča, amd ploča, atx ploča, gejming ploča, cena ploče, kupiti ploču"}',
        true,
        '🔌'
    );
    RAISE NOTICE '  ✅ Created: maticne-ploce';

    -- 7. Napajanja (Блоки питания)
    INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, meta_keywords, is_active, icon)
    VALUES (
        'napajanja',
        parent_uuid,
        3,
        'elektronika/racunarske-komponente/napajanja',
        70,
        '{"en": "Power Supplies (PSU)", "ru": "Блоки питания", "sr": "Napajanja (PSU)"}',
        '{"en": "PC power supplies - modular, semi-modular, 80+ certified", "ru": "Блоки питания для ПК - модульные, полумодульные, сертификация 80+", "sr": "PC napajanja - modularna, polu-modularna, 80+ sertifikat"}',
        '{"en": "Power Supplies | Vondi", "ru": "Блоки питания | Vondi", "sr": "Napajanja | Vondi"}',
        '{"en": "Buy PSU online - 650W, 750W, 850W, 1000W power supplies", "ru": "Купить блоки питания - 650Вт, 750Вт, 850Вт, 1000Вт", "sr": "Kupite napajanja - 650W, 750W, 850W, 1000W"}',
        '{"en": "psu, power supply, 650w psu, 750w psu, 850w psu, modular psu, 80 plus gold, 80 plus platinum, corsair psu, seasonic psu, evga psu, psu price, buy psu", "ru": "блок питания, бп, 650вт бп, 750вт бп, 850вт бп, модульный бп, 80 plus gold, купить бп, цена бп", "sr": "napajanje, psu, 650w napajanje, 750w napajanje, 850w napajanje, modularno, 80 plus gold, corsair psu, cena napajanja, kupiti napajanje"}',
        true,
        '⚡'
    );
    RAISE NOTICE '  ✅ Created: napajanja';

    -- 8. Kućišta (Корпуса)
    INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, meta_keywords, is_active, icon)
    VALUES (
        'kucista',
        parent_uuid,
        3,
        'elektronika/racunarske-komponente/kucista',
        80,
        '{"en": "PC Cases", "ru": "Корпуса", "sr": "Kućišta"}',
        '{"en": "PC cases - ATX, Micro-ATX, Mini-ITX, tower cases", "ru": "Корпуса для ПК - ATX, Micro-ATX, Mini-ITX, башенные корпуса", "sr": "PC kućišta - ATX, Micro-ATX, Mini-ITX, tower kućišta"}',
        '{"en": "PC Cases | Vondi", "ru": "Корпуса | Vondi", "sr": "Kućišta | Vondi"}',
        '{"en": "Buy PC cases online - ATX, Mid Tower, Full Tower, Mini-ITX", "ru": "Купить корпуса для ПК онлайн - ATX, Mid Tower, Full Tower, Mini-ITX", "sr": "Kupite PC kućišta online - ATX, Mid Tower, Full Tower, Mini-ITX"}',
        '{"en": "pc case, computer case, atx case, mid tower, full tower, mini itx case, tempered glass case, rgb case, fractal design, nzxt case, corsair case, case price, buy case", "ru": "корпус пк, корпус компьютера, atx корпус, mid tower, full tower, mini itx, корпус со стеклом, rgb корпус, купить корпус, цена корпуса", "sr": "pc kućište, atx kućište, mid tower, full tower, mini itx, staklo kućište, rgb kućište, fractal design, nzxt kućište, cena kućišta, kupiti kućište"}',
        true,
        '📦'
    );
    RAISE NOTICE '  ✅ Created: kucista';

    -- 9. Hlađenje (Охлаждение)
    INSERT INTO categories (slug, parent_id, level, path, sort_order, name, description, meta_title, meta_description, meta_keywords, is_active, icon)
    VALUES (
        'hladjenje',
        parent_uuid,
        3,
        'elektronika/racunarske-komponente/hladjenje',
        90,
        '{"en": "Cooling", "ru": "Охлаждение", "sr": "Hlađenje"}',
        '{"en": "CPU coolers - AIO liquid cooling, air coolers, case fans", "ru": "Системы охлаждения процессора - AIO СВО, воздушные кулеры, вентиляторы", "sr": "CPU hladnjaci - AIO vodeno hlađenje, vazdušno hlađenje, ventilatori"}',
        '{"en": "Cooling | Vondi", "ru": "Охлаждение | Vondi", "sr": "Hlađenje | Vondi"}',
        '{"en": "Buy CPU coolers online - AIO 240mm, 360mm, tower coolers", "ru": "Купить системы охлаждения онлайн - AIO 240мм, 360мм, башенные кулеры", "sr": "Kupite sisteme hlađenja online - AIO 240mm, 360mm, kula hlađenje"}',
        '{"en": "cpu cooler, aio cooler, liquid cooling, air cooler, tower cooler, 240mm aio, 360mm aio, case fans, rgb fans, noctua, corsair cooling, cooler price, buy cooler", "ru": "кулер процессора, aio кулер, водяное охлаждение, воздушный кулер, 240мм aio, 360мм aio, купить кулер, цена кулера", "sr": "cpu hladnjak, aio hladnjak, vodeno hlađenje, vazdušni hladnjak, 240mm aio, 360mm aio, ventilatori, rgb ventilatori, cena hladnjaka, kupiti hladnjak"}',
        true,
        '❄️'
    );
    RAISE NOTICE '  ✅ Created: hladjenje';

    RAISE NOTICE '';
    RAISE NOTICE '✅ Phase 0 completed! Created 9 new general subcategories.';
    RAISE NOTICE '';
    RAISE NOTICE 'Next step: Run migration 000023 to link attributes to these categories.';
    RAISE NOTICE '';
END $$;

COMMIT;
