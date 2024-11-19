-- 0005_add_status_and_indexes.up.sql
ALTER TABLE bookings 
    ADD COLUMN status VARCHAR(20) NOT NULL 
    DEFAULT 'pending' 
    CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed'));

CREATE INDEX idx_bookings_status ON bookings(status);
CREATE INDEX idx_bookings_dates ON bookings(start_date, end_date);