-- 0013_update_bed_bookings.up.sql
ALTER TABLE bed_bookings
ADD COLUMN IF NOT EXISTS status VARCHAR(20) NOT NULL DEFAULT 'pending'
    CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed'));

-- Создаем индекс для оптимизации поиска
CREATE INDEX IF NOT EXISTS idx_bed_bookings_dates_status 
ON bed_bookings(start_date, end_date, status);