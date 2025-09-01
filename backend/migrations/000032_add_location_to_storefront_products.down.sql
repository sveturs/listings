-- Откат миграции - возвращаем товары витрин к состоянию без индивидуального местоположения
-- если они были добавлены этой миграцией

-- Сбрасываем индивидуальное местоположение для товаров, 
-- которые получили его от витрин в этой миграции
UPDATE storefront_products sp
SET 
    individual_address = NULL,
    individual_latitude = NULL,
    individual_longitude = NULL,
    has_individual_location = false
FROM storefronts s
WHERE sp.storefront_id = s.id
    AND sp.has_individual_location = true
    AND sp.individual_latitude = s.latitude
    AND sp.individual_longitude = s.longitude;

-- Для остальных товаров, которые получили случайные координаты
-- сбрасываем их только если адрес содержит известные города
UPDATE storefront_products
SET 
    individual_address = NULL,
    individual_latitude = NULL,
    individual_longitude = NULL,
    has_individual_location = false
WHERE has_individual_location = true
    AND individual_address IN (
        'Белград, Стари Град',
        'Белград, Нови Београд',
        'Белград, Земун',
        'Белград, Палилула',
        'Белград, Звездара',
        'Белград, Вождовац',
        'Белград, Чукарица',
        'Белград, Раковица',
        'Нови Сад, Центр',
        'Нови Сад, Лиман',
        'Нови Сад, Ново Насеље',
        'Нови Сад, Детелинара',
        'Ниш, Центр',
        'Крагујевац, Центр',
        'Суботица, Центр',
        'Зрењанин, Центр',
        'Панчево, Центр',
        'Смедерево, Центр',
        'Вршац, Центр',
        'Шабац, Центр'
    );

-- Обратные изменения для marketplace_listings не делаем,
-- так как добавление адресов в поле location не критично