
-- Add NOT NULL constraint to ensure data integrity.
ALTER TABLE bookings
ALTER COLUMN user_id SET NOT NULL,
ALTER COLUMN room_id SET NOT NULL;

-- Add an index to optimize queries on bookings table.
CREATE INDEX idx_bookings_user_id ON bookings(user_id);
CREATE INDEX idx_bookings_room_id ON bookings(room_id);
