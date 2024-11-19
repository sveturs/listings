-- 0012_add_accommodation_types.down.sql

-- Удаляем поле status из таблицы bookings
ALTER TABLE bookings
    DROP COLUMN IF EXISTS status;

-- Удаляем таблицу bed_bookings
DROP TABLE IF EXISTS bed_bookings;

-- Удаляем таблицу beds
DROP TABLE IF EXISTS beds;

-- Удаляем добавленные колонки из таблицы rooms
ALTER TABLE rooms
    DROP COLUMN IF EXISTS accommodation_type,
    DROP COLUMN IF EXISTS is_shared,
    DROP COLUMN IF EXISTS total_beds,
    DROP COLUMN IF EXISTS available_beds,
    DROP COLUMN IF EXISTS has_private_bathroom;