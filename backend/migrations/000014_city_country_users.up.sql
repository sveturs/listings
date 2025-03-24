-- Добавление полей city и country в таблицу users, если их еще нет

DO $$
BEGIN
    -- Проверяем наличие колонки city
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'city'
    ) THEN
        -- Если колонки city нет, добавляем ее
        ALTER TABLE users ADD COLUMN city VARCHAR(100);
    END IF;
    
    -- Проверяем наличие колонки country
    IF NOT EXISTS (
        SELECT 1 
        FROM information_schema.columns 
        WHERE table_name = 'users' AND column_name = 'country'
    ) THEN
        -- Если колонки country нет, добавляем ее
        ALTER TABLE users ADD COLUMN country VARCHAR(100);
    END IF;
END
$$;