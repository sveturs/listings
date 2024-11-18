
-- Remove the NOT NULL constraint.
ALTER TABLE bookings
ALTER COLUMN user_id DROP NOT NULL,
ALTER COLUMN room_id DROP NOT NULL;

-- Drop the added indexes.
DROP INDEX IF EXISTS idx_bookings_user_id;
DROP INDEX IF EXISTS idx_bookings_room_id;
