-- 0011_extend_room_types.down.sql
DROP TABLE IF EXISTS bed_bookings;
DROP TABLE IF EXISTS beds;
ALTER TABLE rooms
    DROP COLUMN accommodation_type,
    DROP COLUMN is_shared,
    DROP COLUMN total_beds,
    DROP COLUMN available_beds,
    DROP COLUMN has_private_bathroom;