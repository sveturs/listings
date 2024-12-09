-- Создаем таблицу отзывов
CREATE TABLE room_reviews (
    id SERIAL PRIMARY KEY,
    room_id INT REFERENCES rooms(id) ON DELETE CASCADE,
    user_id INT REFERENCES users(id) ON DELETE CASCADE,
    rating INT CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(room_id, user_id) -- один отзыв от пользователя на комнату
);

-- Индексы для оптимизации запросов
CREATE INDEX idx_room_reviews_room_id ON room_reviews(room_id);
CREATE INDEX idx_room_reviews_user_id ON room_reviews(user_id);
CREATE INDEX idx_room_reviews_rating ON room_reviews(rating);

-- Индексы для улучшения производительности фильтрации и сортировки
CREATE INDEX idx_rooms_price_created ON rooms(price_per_night, created_at);
CREATE INDEX idx_rooms_type_price ON rooms(accommodation_type, price_per_night);
CREATE INDEX idx_rooms_city_country ON rooms(address_city, address_country);