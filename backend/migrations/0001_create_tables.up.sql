-- backend/migrations/0001_create_tables.up.sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Создаем демо-пользователя с ID=1
INSERT INTO users (id, name, email, created_at) VALUES 
(1, 'Demo User', 'test@example.com', CURRENT_TIMESTAMP);

-- Настраиваем sequence чтобы следующий ID был 2
SELECT setval('users_id_seq', 1, true);