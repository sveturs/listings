-- Удаляем демо данные
DELETE FROM bed_images;
DELETE FROM beds;
DELETE FROM room_images;
DELETE FROM rooms;
DELETE FROM users WHERE email = 'demo@example.com';

-- Восстанавливаем старое ограничение
DROP INDEX IF EXISTS unique_main_image_per_room;
ALTER TABLE room_images ADD CONSTRAINT unique_main_image_per_room UNIQUE (room_id, is_main);