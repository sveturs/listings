-- 0012_add_accommodation_types.up.sql
ALTER TABLE rooms
    ADD COLUMN accommodation_type VARCHAR(50) NOT NULL DEFAULT 'room'
        CHECK (accommodation_type IN ('bed', 'room', 'apartment')),
    ADD COLUMN is_shared BOOLEAN NOT NULL DEFAULT false,
    ADD COLUMN total_beds INT,
    ADD COLUMN available_beds INT,
    ADD COLUMN has_private_bathroom BOOLEAN NOT NULL DEFAULT true;

-- Создаем таблицу для отдельных кроватей
CREATE TABLE beds (
    id SERIAL PRIMARY KEY,
    room_id INT REFERENCES rooms(id) ON DELETE CASCADE,
    bed_number VARCHAR(50) NOT NULL,
    is_available BOOLEAN NOT NULL DEFAULT true,
    price_per_night DECIMAL(10, 2) NOT NULL CHECK (price_per_night >= 0),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(room_id, bed_number)
);

-- Создаем таблицу для бронирования кроватей
CREATE TABLE bed_bookings (
    id SERIAL PRIMARY KEY,
    bed_id INT REFERENCES beds(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL CHECK (end_date > start_date),
    status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Добавляем поле status в существующую таблицу bookings
ALTER TABLE bookings
    ADD COLUMN status VARCHAR(20) NOT NULL DEFAULT 'pending'
        CHECK (status IN ('pending', 'confirmed', 'cancelled', 'completed'));