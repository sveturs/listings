-- backend/migrations/0020_add_demo_data.up.sql

-- Сначала фиксим ограничение для изображений
ALTER TABLE room_images DROP CONSTRAINT IF EXISTS unique_main_image_per_room;
CREATE UNIQUE INDEX unique_main_image_per_room ON room_images (room_id) WHERE is_main = true;

-- Очищаем существующие данные
TRUNCATE bed_images, beds, room_images, rooms, users CASCADE;

-- Сбрасываем последовательности
ALTER SEQUENCE rooms_id_seq RESTART WITH 1;
ALTER SEQUENCE beds_id_seq RESTART WITH 1;
ALTER SEQUENCE room_images_id_seq RESTART WITH 1;
ALTER SEQUENCE bed_images_id_seq RESTART WITH 1;
ALTER SEQUENCE users_id_seq RESTART WITH 1;

-- Создаем тестового пользователя
INSERT INTO users (name, email) VALUES
('Demo User', 'demo@example.com');

-- Создаем демонстрационные объекты

-- 1. Апартаменты в центре
INSERT INTO rooms (
    name, capacity, price_per_night,
    address_street, address_city, address_state, address_country, address_postal_code,
    accommodation_type, is_shared, has_private_bathroom,
    latitude, longitude, formatted_address
) VALUES (
    'Уютные апартаменты в центре', 4, 80,
    'Dunavska 35', 'Novi Sad', 'Vojvodina', 'Serbia', '21000',
    'apartment', false, true,
    45.255421, 19.845241,
    'Dunavska 35, Novi Sad, Serbia'
);

-- Добавляем изображения для первых апартаментов
INSERT INTO room_images (room_id, file_path, file_name, file_size, content_type, is_main) VALUES
((SELECT MAX(id) FROM rooms), '1.jpg', 'apartment1.jpg', 1024, 'image/jpeg', true),
((SELECT MAX(id) FROM rooms), '2.jpg', 'apartment2.jpg', 1024, 'image/jpeg', false),
((SELECT MAX(id) FROM rooms), '3.jpg', 'apartment3.jpg', 1024, 'image/jpeg', false);

-- 2. Апартаменты возле парка
INSERT INTO rooms (
    name, capacity, price_per_night,
    address_street, address_city, address_state, address_country, address_postal_code,
    accommodation_type, is_shared, has_private_bathroom,
    latitude, longitude, formatted_address
) VALUES (
    'Просторные апартаменты у парка', 6, 120,
    'Futoška 12', 'Novi Sad', 'Vojvodina', 'Serbia', '21000',
    'apartment', false, true,
    45.249877, 19.833657,
    'Futoška 12, Novi Sad, Serbia'
);

-- Добавляем изображения для вторых апартаментов
INSERT INTO room_images (room_id, file_path, file_name, file_size, content_type, is_main) VALUES
((SELECT MAX(id) FROM rooms), '4.jpg', 'apartment4.jpg', 1024, 'image/jpeg', true),
((SELECT MAX(id) FROM rooms), '5.jpg', 'apartment5.jpg', 1024, 'image/jpeg', false),
((SELECT MAX(id) FROM rooms), '6.jpg', 'apartment6.jpg', 1024, 'image/jpeg', false);

-- 3. Приватная комната в историческом центре
INSERT INTO rooms (
    name, capacity, price_per_night,
    address_street, address_city, address_state, address_country, address_postal_code,
    accommodation_type, is_shared, has_private_bathroom,
    latitude, longitude, formatted_address
) VALUES (
    'Уютная комната в историческом центре', 2, 35,
    'Zmaj Jovina 4', 'Novi Sad', 'Vojvodina', 'Serbia', '21000',
    'room', false, true,
    45.254663, 19.844966,
    'Zmaj Jovina 4, Novi Sad, Serbia'
);

-- Добавляем изображения для первой приватной комнаты
INSERT INTO room_images (room_id, file_path, file_name, file_size, content_type, is_main) VALUES
((SELECT MAX(id) FROM rooms), '7.jpg', 'room1.jpg', 1024, 'image/jpeg', true),
((SELECT MAX(id) FROM rooms), '8.jpg', 'room2.jpg', 1024, 'image/jpeg', false),
((SELECT MAX(id) FROM rooms), '1.jpg', 'room3.jpg', 1024, 'image/jpeg', false);

-- 4. Приватная комната возле набережной
INSERT INTO rooms (
    name, capacity, price_per_night,
    address_street, address_city, address_state, address_country, address_postal_code,
    accommodation_type, is_shared, has_private_bathroom,
    latitude, longitude, formatted_address
) VALUES (
    'Комната с видом на Дунай', 2, 40,
    'Beogradski kej 31', 'Novi Sad', 'Vojvodina', 'Serbia', '21000',
    'room', false, true,
    45.256893, 19.861559,
    'Beogradski kej 31, Novi Sad, Serbia'
);

-- Добавляем изображения для второй приватной комнаты
INSERT INTO room_images (room_id, file_path, file_name, file_size, content_type, is_main) VALUES
((SELECT MAX(id) FROM rooms), '2.jpg', 'room4.jpg', 1024, 'image/jpeg', true),
((SELECT MAX(id) FROM rooms), '3.jpg', 'room5.jpg', 1024, 'image/jpeg', false),
((SELECT MAX(id) FROM rooms), '4.jpg', 'room6.jpg', 1024, 'image/jpeg', false);

-- 5. Хостел в центре (койко-места)
INSERT INTO rooms (
    name, capacity, price_per_night,
    address_street, address_city, address_state, address_country, address_postal_code,
    accommodation_type, is_shared, total_beds, available_beds, has_private_bathroom,
    latitude, longitude, formatted_address
) VALUES (
    'Центральный хостел', 6, 12,
    'Miletićeva 15', 'Novi Sad', 'Vojvodina', 'Serbia', '21000',
    'bed', true, 6, 6, true,
    45.252558, 19.842895,
    'Miletićeva 15, Novi Sad, Serbia'
);

-- Добавляем изображения для первого хостела
INSERT INTO room_images (room_id, file_path, file_name, file_size, content_type, is_main) VALUES
((SELECT MAX(id) FROM rooms), '5.jpg', 'hostel1.jpg', 1024, 'image/jpeg', true),
((SELECT MAX(id) FROM rooms), '6.jpg', 'hostel2.jpg', 1024, 'image/jpeg', false),
((SELECT MAX(id) FROM rooms), '7.jpg', 'hostel3.jpg', 1024, 'image/jpeg', false);

-- Добавляем кровати для первого хостела
INSERT INTO beds (room_id, bed_number, price_per_night, is_available) VALUES
((SELECT MAX(id) FROM rooms), 11, 12, true),
((SELECT MAX(id) FROM rooms), 12, 12, true),
((SELECT MAX(id) FROM rooms), 21, 12, true),
((SELECT MAX(id) FROM rooms), 22, 12, true),
((SELECT MAX(id) FROM rooms), 31, 12, true),
((SELECT MAX(id) FROM rooms), 32, 12, true);

-- Добавляем изображения для кроватей первого хостела
INSERT INTO bed_images (bed_id, file_path, file_name, file_size, content_type) VALUES
((SELECT MAX(id) FROM beds), '8.jpg', 'bed1.jpg', 1024, 'image/jpeg'),
((SELECT MAX(id) FROM beds), '1.jpg', 'bed2.jpg', 1024, 'image/jpeg'),
((SELECT MAX(id) - 1 FROM beds), '2.jpg', 'bed3.jpg', 1024, 'image/jpeg'),
((SELECT MAX(id) - 2 FROM beds), '3.jpg', 'bed4.jpg', 1024, 'image/jpeg'),
((SELECT MAX(id) - 3 FROM beds), '4.jpg', 'bed5.jpg', 1024, 'image/jpeg'),
((SELECT MAX(id) - 4 FROM beds), '5.jpg', 'bed6.jpg', 1024, 'image/jpeg');

-- 6. Хостел возле вокзала (койко-места)
INSERT INTO rooms (
    name, capacity, price_per_night,
    address_street, address_city, address_state, address_country, address_postal_code,
    accommodation_type, is_shared, total_beds, available_beds, has_private_bathroom,
    latitude, longitude, formatted_address
) VALUES (
    'Хостел у вокзала', 4, 10,
    'Bulevar Jaše Tomića 5', 'Novi Sad', 'Vojvodina', 'Serbia', '21000',
    'bed', true, 4, 4, true,
    45.260721, 19.831572,
    'Bulevar Jaše Tomića 5, Novi Sad, Serbia'
);

-- Добавляем изображения для второго хостела
INSERT INTO room_images (room_id, file_path, file_name, file_size, content_type, is_main) VALUES
((SELECT MAX(id) FROM rooms), '6.jpg', 'hostel4.jpg', 1024, 'image/jpeg', true),
((SELECT MAX(id) FROM rooms), '7.jpg', 'hostel5.jpg', 1024, 'image/jpeg', false),
((SELECT MAX(id) FROM rooms), '8.jpg', 'hostel6.jpg', 1024, 'image/jpeg', false);

-- Добавляем кровати для второго хостела
INSERT INTO beds (room_id, bed_number, price_per_night, is_available) VALUES
((SELECT MAX(id) FROM rooms), 1, 10, true),
((SELECT MAX(id) FROM rooms), 2, 10, true),
((SELECT MAX(id) FROM rooms), 3, 10, true),
((SELECT MAX(id) FROM rooms), 4, 10, true);

-- Добавляем изображения для кроватей второго хостела
INSERT INTO bed_images (bed_id, file_path, file_name, file_size, content_type) VALUES
((SELECT MAX(id) FROM beds), '1.jpg', 'bed7.jpg', 1024, 'image/jpeg'),
((SELECT MAX(id) - 1 FROM beds), '2.jpg', 'bed8.jpg', 1024, 'image/jpeg'),
((SELECT MAX(id) - 2 FROM beds), '3.jpg', 'bed9.jpg', 1024, 'image/jpeg'),
((SELECT MAX(id) - 3 FROM beds), '4.jpg', 'bed10.jpg', 1024, 'image/jpeg');

