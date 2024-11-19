-- 0005_add_status_and_indexes.down.sql
DROP INDEX IF EXISTS idx_bookings_status;
DROP INDEX IF EXISTS idx_bookings_dates;
ALTER TABLE bookings DROP COLUMN status;