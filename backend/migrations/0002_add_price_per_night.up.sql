ALTER TABLE rooms 
    ADD COLUMN price_per_night NUMERIC(10, 2) NOT NULL DEFAULT 0 
    CHECK (price_per_night >= 0);
CREATE INDEX idx_users_email ON users(email);

-- Добавить составной индекс для фильтрации комнат
CREATE INDEX idx_rooms_capacity_price ON rooms(capacity, price_per_night);