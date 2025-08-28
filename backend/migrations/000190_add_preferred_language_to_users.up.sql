-- Добавляем поле preferred_language в таблицу users для поддержки мультиязычности
ALTER TABLE users
ADD COLUMN preferred_language VARCHAR(10) DEFAULT 'ru';

-- Добавляем проверку на допустимые значения языка
ALTER TABLE users
ADD CONSTRAINT users_preferred_language_check CHECK (preferred_language IN ('ru', 'sr', 'en'));

-- Создаем индекс для быстрого поиска по языку (для аналитики)
CREATE INDEX idx_users_preferred_language ON users(preferred_language);

-- Обновляем существующих пользователей на основе их текущих настроек или страны
-- Пользователи с сербскими телефонами получают 'sr', остальные - 'ru'
UPDATE users 
SET preferred_language = CASE 
    WHEN phone LIKE '+381%' THEN 'sr'
    WHEN phone LIKE '+1%' OR phone LIKE '+44%' THEN 'en' 
    ELSE 'ru'
END
WHERE preferred_language IS NULL;