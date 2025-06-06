-- Возвращаем размер поля token обратно к VARCHAR(255)
ALTER TABLE refresh_tokens ALTER COLUMN token TYPE VARCHAR(255);