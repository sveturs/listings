-- 0006_add_room_address.down.sql
DROP INDEX IF EXISTS idx_rooms_city;
DROP INDEX IF EXISTS idx_rooms_country;

ALTER TABLE rooms
    DROP COLUMN address_street,
    DROP COLUMN address_city,
    DROP COLUMN address_state,
    DROP COLUMN address_country,
    DROP COLUMN address_postal_code;