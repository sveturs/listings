
-- Remove constraints and indexes added in the up migration
ALTER TABLE bookings DROP CONSTRAINT valid_date_range;

DROP INDEX IF EXISTS idx_rooms_price;
DROP INDEX IF EXISTS idx_bookings_room_date;
