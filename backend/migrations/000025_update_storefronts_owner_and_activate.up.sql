-- Миграция для изменения владельца всех витрин на voroshilovdo@gmail.com
-- и активации витрины с изображениями

-- 1. Изменяем владельца всех витрин на user_id=6 (voroshilovdo@gmail.com)
UPDATE storefronts
SET user_id = 6
WHERE user_id != 6;

-- 2. Активируем витрину 123455 (id=39), у которой есть логотип и баннер
UPDATE storefronts
SET is_active = true
WHERE id = 39;

-- 3. Также активируем другие витрины с изображениями
UPDATE storefronts
SET is_active = true
WHERE id IN (36, 38) AND (logo_url IS NOT NULL OR banner_url IS NOT NULL);

-- 4. Товары витрин не имеют user_id, они связаны через storefront_id
-- Поэтому пропускаем этот шаг

-- 5. Добавляем недостающие поля для витрин если их нет
ALTER TABLE storefronts
ADD COLUMN IF NOT EXISTS followers_count INTEGER DEFAULT 0;

ALTER TABLE storefronts
ADD COLUMN IF NOT EXISTS rating DECIMAL(3,2) DEFAULT 0.00;

-- 6. Устанавливаем начальные значения для активных витрин
UPDATE storefronts
SET followers_count = CASE
    WHEN id = 39 THEN 125  -- 123455
    WHEN id = 36 THEN 890  -- 3G Store
    WHEN id = 38 THEN 456  -- sportmag
    WHEN id = 34 THEN 234  -- AgroShop
    WHEN id = 31 THEN 567  -- TechNova
    WHEN id = 32 THEN 789  -- Fashion House
    WHEN id = 33 THEN 345  -- Home & Garden
    ELSE 0
END,
rating = CASE
    WHEN id = 39 THEN 4.5  -- 123455
    WHEN id = 36 THEN 4.8  -- 3G Store
    WHEN id = 38 THEN 4.2  -- sportmag
    WHEN id = 34 THEN 4.6  -- AgroShop
    WHEN id = 31 THEN 4.7  -- TechNova
    WHEN id = 32 THEN 4.9  -- Fashion House
    WHEN id = 33 THEN 4.3  -- Home & Garden
    ELSE 0.00
END
WHERE is_active = true;