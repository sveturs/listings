-- Увеличиваем размер поля token для refresh токенов
ALTER TABLE refresh_tokens ALTER COLUMN token TYPE TEXT;