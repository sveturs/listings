
-- Add additional constraints and indexes to fix issues
ALTER TABLE bookings
    ADD CONSTRAINT valid_date_range 
    CHECK (start_date < end_date AND start_date >= CURRENT_DATE);

-- Indexes to optimize queries
CREATE INDEX idx_rooms_price ON rooms (price_per_night);
CREATE INDEX idx_bookings_room_date ON bookings (room_id, start_date, end_date);
