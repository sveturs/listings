-- 0013_update_bed_bookings.down.sql
DROP INDEX IF EXISTS idx_bed_bookings_dates_status;
ALTER TABLE bed_bookings DROP COLUMN IF EXISTS status;