-- 0006_add_room_address.up.sql
ALTER TABLE rooms
    ADD COLUMN address_street VARCHAR(255) NOT NULL DEFAULT '',
    ADD COLUMN address_city VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN address_state VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN address_country VARCHAR(100) NOT NULL DEFAULT '',
    ADD COLUMN address_postal_code VARCHAR(20) NOT NULL DEFAULT '';

CREATE INDEX idx_rooms_city ON rooms(address_city);
CREATE INDEX idx_rooms_country ON rooms(address_country);