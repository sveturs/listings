-- Синхронизация пользователей с auth сервисом
-- Auth сервис говорит, что www.svetu.rs@gmail.com имеет ID=7
-- Поэтому приводим локальную БД в соответствие

BEGIN;

-- 1. Сохраняем старый email пользователя ID=7 в временное поле для истории
ALTER TABLE users ADD COLUMN IF NOT EXISTS old_email TEXT;
UPDATE users SET old_email = email WHERE id = 7;

-- 2. Сначала архивируем пользователя ID=10 чтобы освободить email
UPDATE users
SET
    email = CONCAT('archived_', email),
    name = CONCAT('[ARCHIVED] ', name)
WHERE id = 10;

-- 3. Переносим все связанные данные с user_id=10 на user_id=7
-- (так как по факту это один и тот же пользователь согласно auth сервису)

-- Обновляем объявления
UPDATE marketplace_listings SET user_id = 7 WHERE user_id = 10;

-- Обновляем чаты где пользователь продавец
UPDATE marketplace_chats SET seller_id = 7 WHERE seller_id = 10;

-- Обновляем чаты где пользователь покупатель
UPDATE marketplace_chats SET buyer_id = 7 WHERE buyer_id = 10;

-- Обновляем сообщения где пользователь отправитель
UPDATE marketplace_messages SET sender_id = 7 WHERE sender_id = 10;

-- Обновляем сообщения где пользователь получатель
UPDATE marketplace_messages SET receiver_id = 7 WHERE receiver_id = 10;

-- Обновляем витрины магазинов
UPDATE storefronts SET user_id = 7 WHERE user_id = 10;

-- Обновляем избранное (если таблица существует)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'marketplace_favorites') THEN
        UPDATE marketplace_favorites SET user_id = 7 WHERE user_id = 10;
    END IF;
END $$;

-- Обновляем отзывы (если таблица существует)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'marketplace_reviews') THEN
        UPDATE marketplace_reviews SET reviewer_id = 7 WHERE reviewer_id = 10;
        UPDATE marketplace_reviews SET reviewee_id = 7 WHERE reviewee_id = 10;
    END IF;
END $$;

-- Обновляем жалобы (если таблица существует)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'marketplace_reports') THEN
        UPDATE marketplace_reports SET reporter_id = 7 WHERE reporter_id = 10;
    END IF;
END $$;

-- Обновляем просмотры (если таблица существует)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'marketplace_views') THEN
        UPDATE marketplace_views SET user_id = 7 WHERE user_id = 10;
    END IF;
END $$;

-- Обновляем балансы (если таблица существует)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'user_balances') THEN
        -- Сначала удаляем существующий баланс для user_id=7, если он есть
        DELETE FROM user_balances WHERE user_id = 7;
        -- Затем переносим баланс от user_id=10
        UPDATE user_balances SET user_id = 7 WHERE user_id = 10;
    END IF;
END $$;

-- Обновляем контакты (если таблица существует)
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name = 'user_contacts') THEN
        UPDATE user_contacts SET user_id = 7 WHERE user_id = 10;
        UPDATE user_contacts SET contact_user_id = 7 WHERE contact_user_id = 10;
    END IF;
END $$;

-- 4. Теперь обновляем данные пользователя ID=7 на актуальные из auth сервиса
UPDATE users
SET
    email = 'www.svetu.rs@gmail.com',
    name = 'sveturs',
    provider = 'google'
WHERE id = 7;

-- 5. Удаляем временную колонку (опционально, можно оставить для истории)
-- ALTER TABLE users DROP COLUMN IF EXISTS old_email;

COMMIT;