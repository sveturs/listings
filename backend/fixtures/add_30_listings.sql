-- Скрипт для создания 30 тестовых объявлений для проверки подгрузки при скролле

-- Создаем объявления в существующих категориях

-- Телефоны (категория 2) - еще 10 штук
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, address_city, address_country, latitude, longitude)
VALUES 
(11, 2, 'Samsung Galaxy S21', 'Отличный смартфон, 128GB', 599.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(12, 2, 'Google Pixel 6', 'В идеальном состоянии', 499.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(13, 2, 'OnePlus 9 Pro', 'Флагманский смартфон', 649.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(11, 2, 'Xiaomi Mi 11', 'Мощный процессор', 449.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(12, 2, 'Sony Xperia 5 III', 'Компактный флагман', 799.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(13, 2, 'Motorola Edge 20', 'Легкий и тонкий', 399.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(11, 2, 'Nokia 8.3 5G', 'Чистый Android', 349.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(12, 2, 'ASUS ROG Phone 5', 'Игровой смартфон', 699.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(13, 2, 'Realme GT', 'Быстрая зарядка', 399.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(11, 2, 'Oppo Find X3', 'Отличная камера', 599.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489);

-- Ноутбуки (категория 3) - 10 штук
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, address_city, address_country, latitude, longitude)
VALUES 
(13, 3, 'MacBook Pro 13"', 'M1 чип, 256GB SSD', 1299.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(11, 3, 'Dell XPS 13', 'Intel i7, 16GB RAM', 1099.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(12, 3, 'Lenovo ThinkPad X1', 'Бизнес ноутбук', 899.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(13, 3, 'HP Spectre x360', 'Трансформер 2-в-1', 999.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(11, 3, 'ASUS ZenBook 14', 'Ультрабук с OLED', 849.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(12, 3, 'MSI Prestige 14', 'Для креативщиков', 1149.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(13, 3, 'Acer Swift 5', 'Легкий и мощный', 749.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(11, 3, 'Surface Laptop 4', 'Стильный дизайн', 1199.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(12, 3, 'Razer Blade 15', 'Игровой ноутбук', 1999.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(13, 3, 'LG Gram 17', 'Большой и легкий', 1299.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489);

-- Мебель (категория 20) - 10 штук
INSERT INTO marketplace_listings (user_id, category_id, title, description, price, condition, status, location, address_city, address_country, latitude, longitude)
VALUES 
(11, 20, 'Кресло офисное', 'Эргономичное кресло', 199.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(12, 20, 'Стеллаж IKEA', '5 полок, белый', 89.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(13, 20, 'Шкаф-купе', '2 метра ширина', 399.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(11, 20, 'Комод деревянный', '6 ящиков', 249.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(12, 20, 'Кровать двуспальная', 'С матрасом', 599.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(13, 20, 'Тумба под ТВ', 'Современный дизайн', 149.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(11, 20, 'Барная стойка', 'С двумя стульями', 299.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(12, 20, 'Кухонный уголок', 'Мягкий диван', 449.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(13, 20, 'Письменный стол', 'С полками', 179.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489),
(11, 20, 'Журнальный столик', 'Стеклянный', 129.00, 'used', 'active', 'Belgrade, Serbia', 'Belgrade', 'Serbia', 44.7866, 20.4489);

-- Выводим количество записей
SELECT COUNT(*) as total_listings FROM marketplace_listings;