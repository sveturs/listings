-- Добавляем больше телефонов для демонстрации похожих товаров

-- iPhone 14 (похож на iPhone 13 Pro Max по бренду)
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, address_city, address_country, latitude, longitude)
VALUES (11, 2, 'iPhone 14', 'Отличное состояние, 128GB', 1099.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489);

-- iPhone 12 (еще один Apple)
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, address_city, address_country, latitude, longitude)
VALUES (12, 2, 'iPhone 12', 'Как новый, 64GB', 699.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489);

-- Samsung Galaxy S22 (похож на S23)
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, address_city, address_country, latitude, longitude)
VALUES (13, 2, 'Samsung Galaxy S22', 'Флагман прошлого года', 699.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489);

-- Получаем ID новых объявлений
SELECT id, title FROM marketplace_listings WHERE title IN ('iPhone 14', 'iPhone 12', 'Samsung Galaxy S22');

-- Добавляем атрибуты для новых телефонов
-- iPhone 14
INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 1, 'Apple' FROM marketplace_listings ml WHERE ml.title = 'iPhone 14';

INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 2, 'used' FROM marketplace_listings ml WHERE ml.title = 'iPhone 14';

INSERT INTO listing_attribute_values (listing_id, attribute_id, numeric_value) 
SELECT ml.id, 4, 6.1 FROM marketplace_listings ml WHERE ml.title = 'iPhone 14';

INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 5, '128GB' FROM marketplace_listings ml WHERE ml.title = 'iPhone 14';

INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 7, 'purple' FROM marketplace_listings ml WHERE ml.title = 'iPhone 14';

-- iPhone 12 
INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 1, 'Apple' FROM marketplace_listings ml WHERE ml.title = 'iPhone 12';

INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 2, 'used' FROM marketplace_listings ml WHERE ml.title = 'iPhone 12';

INSERT INTO listing_attribute_values (listing_id, attribute_id, numeric_value) 
SELECT ml.id, 4, 6.1 FROM marketplace_listings ml WHERE ml.title = 'iPhone 12';

INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 5, '64GB' FROM marketplace_listings ml WHERE ml.title = 'iPhone 12';

INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 7, 'blue' FROM marketplace_listings ml WHERE ml.title = 'iPhone 12';

-- Samsung Galaxy S22
INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 1, 'Samsung' FROM marketplace_listings ml WHERE ml.title = 'Samsung Galaxy S22';

INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 2, 'used' FROM marketplace_listings ml WHERE ml.title = 'Samsung Galaxy S22';

INSERT INTO listing_attribute_values (listing_id, attribute_id, numeric_value) 
SELECT ml.id, 4, 6.1 FROM marketplace_listings ml WHERE ml.title = 'Samsung Galaxy S22';

INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 5, '128GB' FROM marketplace_listings ml WHERE ml.title = 'Samsung Galaxy S22';

INSERT INTO listing_attribute_values (listing_id, attribute_id, text_value) 
SELECT ml.id, 7, 'green' FROM marketplace_listings ml WHERE ml.title = 'Samsung Galaxy S22';

-- Проверяем результат
SELECT ml.id, ml.title, ml.price, ca.name, lav.text_value
FROM marketplace_listings ml
LEFT JOIN listing_attribute_values lav ON ml.id = lav.listing_id AND lav.attribute_id = 1
LEFT JOIN category_attributes ca ON lav.attribute_id = ca.id
WHERE ml.category_id = 2
ORDER BY ml.id;