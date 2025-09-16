-- Синхронизация ID пользователей с центральным auth микросервисом
-- Auth сервис имеет:
-- ID=6 для voroshilovdo@gmail.com (был ID=2 в локальной БД)
-- ID=8 для 4hash92@gmail.com (был ID=6 в локальной БД)
-- ID=7 для www.svetu.rs@gmail.com (уже синхронизирован)

BEGIN;

-- Для безопасности сначала проверяем, что эти пользователи существуют
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = 2 AND email = 'voroshilovdo@gmail.com') THEN
        RAISE EXCEPTION 'User with ID=2 and email=voroshilovdo@gmail.com not found';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE id = 6 AND email = '4hash92@gmail.com') THEN
        RAISE EXCEPTION 'User with ID=6 and email=4hash92@gmail.com not found';
    END IF;
END $$;

-- Временно отключаем ограничения внешних ключей для обновления
SET session_replication_role = 'replica';

-- 0. Сначала освобождаем ID=8, который занят другим пользователем
-- Перемещаем EmailEmail@EmailEmail.ru (ID=8) на временный ID=208
UPDATE users SET id = 208 WHERE id = 8 AND email = 'EmailEmail@EmailEmail.ru';
UPDATE marketplace_listings SET user_id = 208 WHERE user_id = 8;
UPDATE marketplace_chats SET seller_id = 208 WHERE seller_id = 8;
UPDATE marketplace_chats SET buyer_id = 208 WHERE buyer_id = 8;
UPDATE marketplace_messages SET sender_id = 208 WHERE sender_id = 8;
UPDATE marketplace_messages SET receiver_id = 208 WHERE receiver_id = 8;
UPDATE marketplace_favorites SET user_id = 208 WHERE user_id = 8;
UPDATE marketplace_orders SET seller_id = 208 WHERE seller_id = 8;
UPDATE marketplace_orders SET buyer_id = 208 WHERE buyer_id = 8;
UPDATE storefronts SET user_id = 208 WHERE user_id = 8;
UPDATE notifications SET user_id = 208 WHERE user_id = 8;
UPDATE notification_settings SET user_id = 208 WHERE user_id = 8;
UPDATE reviews SET user_id = 208 WHERE user_id = 8;
UPDATE review_responses SET user_id = 208 WHERE user_id = 8;
UPDATE review_votes SET user_id = 208 WHERE user_id = 8;
UPDATE user_roles SET user_id = 208 WHERE user_id = 8;
UPDATE user_contacts SET user_id = 208 WHERE user_id = 8;
UPDATE user_contacts SET contact_user_id = 208 WHERE contact_user_id = 8;
UPDATE user_balances SET user_id = 208 WHERE user_id = 8;
UPDATE user_privacy_settings SET user_id = 208 WHERE user_id = 8;
UPDATE user_behavior_events SET user_id = 208 WHERE user_id = 8;
UPDATE user_telegram_connections SET user_id = 208 WHERE user_id = 8;
UPDATE subscription_history SET user_id = 208 WHERE user_id = 8;
UPDATE subscription_payments SET user_id = 208 WHERE user_id = 8;
UPDATE user_subscriptions SET user_id = 208 WHERE user_id = 8;
UPDATE balance_transactions SET user_id = 208 WHERE user_id = 8;
UPDATE payment_transactions SET user_id = 208 WHERE user_id = 8;
UPDATE escrow_payments SET seller_id = 208 WHERE seller_id = 8;
UPDATE escrow_payments SET buyer_id = 208 WHERE buyer_id = 8;
UPDATE merchant_payouts SET seller_id = 208 WHERE seller_id = 8;
UPDATE shopping_carts SET user_id = 208 WHERE user_id = 8;
UPDATE listing_views SET user_id = 208 WHERE user_id = 8;
UPDATE search_statistics SET user_id = 208 WHERE user_id = 8;
UPDATE gis_filter_analytics SET user_id = 208 WHERE user_id = 8;
UPDATE address_change_log SET user_id = 208 WHERE user_id = 8;
UPDATE translation_audit_log SET user_id = 208 WHERE user_id = 8;
UPDATE role_audit_log SET user_id = 208 WHERE user_id = 8;
UPDATE role_audit_log SET target_user_id = 208 WHERE target_user_id = 8;
UPDATE storefront_staff SET user_id = 208 WHERE user_id = 8;
UPDATE storefront_inventory_movements SET user_id = 208 WHERE user_id = 8;
UPDATE user_storefronts SET user_id = 208 WHERE user_id = 8;

-- 1. Создаем временных пользователей с новыми ID
INSERT INTO users (id, email, name, provider, created_at, updated_at, old_email)
SELECT
    CASE
        WHEN id = 2 THEN 106  -- временный ID для voroshilovdo
        WHEN id = 6 THEN 108  -- временный ID для 4hash92
    END as new_id,
    CONCAT('temp_', email) as temp_email,
    name,
    provider,
    created_at,
    updated_at,
    email as old_email
FROM users
WHERE id IN (2, 6);

-- 2. Обновляем все связанные таблицы для voroshilovdo@gmail.com (2 -> 106)
UPDATE marketplace_listings SET user_id = 106 WHERE user_id = 2;
UPDATE marketplace_chats SET seller_id = 106 WHERE seller_id = 2;
UPDATE marketplace_chats SET buyer_id = 106 WHERE buyer_id = 2;
UPDATE marketplace_messages SET sender_id = 106 WHERE sender_id = 2;
UPDATE marketplace_messages SET receiver_id = 106 WHERE receiver_id = 2;
UPDATE marketplace_favorites SET user_id = 106 WHERE user_id = 2;
UPDATE marketplace_orders SET seller_id = 106 WHERE seller_id = 2;
UPDATE marketplace_orders SET buyer_id = 106 WHERE buyer_id = 2;
UPDATE storefronts SET user_id = 106 WHERE user_id = 2;
UPDATE storefront_staff SET user_id = 106 WHERE user_id = 2;
UPDATE storefront_inventory_movements SET user_id = 106 WHERE user_id = 2;
UPDATE user_storefronts SET user_id = 106 WHERE user_id = 2;
UPDATE notifications SET user_id = 106 WHERE user_id = 2;
UPDATE notification_settings SET user_id = 106 WHERE user_id = 2;
UPDATE reviews SET user_id = 106 WHERE user_id = 2;
UPDATE review_responses SET user_id = 106 WHERE user_id = 2;
UPDATE review_votes SET user_id = 106 WHERE user_id = 2;
UPDATE user_roles SET user_id = 106 WHERE user_id = 2;
UPDATE user_contacts SET user_id = 106 WHERE user_id = 2;
UPDATE user_contacts SET contact_user_id = 106 WHERE contact_user_id = 2;
UPDATE user_balances SET user_id = 106 WHERE user_id = 2;
UPDATE user_privacy_settings SET user_id = 106 WHERE user_id = 2;
UPDATE user_behavior_events SET user_id = 106 WHERE user_id = 2;
UPDATE user_telegram_connections SET user_id = 106 WHERE user_id = 2;
UPDATE subscription_history SET user_id = 106 WHERE user_id = 2;
UPDATE subscription_payments SET user_id = 106 WHERE user_id = 2;
UPDATE user_subscriptions SET user_id = 106 WHERE user_id = 2;
UPDATE balance_transactions SET user_id = 106 WHERE user_id = 2;
UPDATE payment_transactions SET user_id = 106 WHERE user_id = 2;
UPDATE escrow_payments SET seller_id = 106 WHERE seller_id = 2;
UPDATE escrow_payments SET buyer_id = 106 WHERE buyer_id = 2;
UPDATE merchant_payouts SET seller_id = 106 WHERE seller_id = 2;
UPDATE shopping_carts SET user_id = 106 WHERE user_id = 2;
UPDATE listing_views SET user_id = 106 WHERE user_id = 2;
UPDATE search_statistics SET user_id = 106 WHERE user_id = 2;
UPDATE gis_filter_analytics SET user_id = 106 WHERE user_id = 2;
UPDATE address_change_log SET user_id = 106 WHERE user_id = 2;
UPDATE translation_audit_log SET user_id = 106 WHERE user_id = 2;
UPDATE role_audit_log SET user_id = 106 WHERE user_id = 2;
UPDATE role_audit_log SET target_user_id = 106 WHERE target_user_id = 2;

-- 3. Обновляем все связанные таблицы для 4hash92@gmail.com (6 -> 108)
UPDATE marketplace_listings SET user_id = 108 WHERE user_id = 6;
UPDATE marketplace_chats SET seller_id = 108 WHERE seller_id = 6;
UPDATE marketplace_chats SET buyer_id = 108 WHERE buyer_id = 6;
UPDATE marketplace_messages SET sender_id = 108 WHERE sender_id = 6;
UPDATE marketplace_messages SET receiver_id = 108 WHERE receiver_id = 6;
UPDATE marketplace_favorites SET user_id = 108 WHERE user_id = 6;
UPDATE marketplace_orders SET seller_id = 108 WHERE seller_id = 6;
UPDATE marketplace_orders SET buyer_id = 108 WHERE buyer_id = 6;
UPDATE storefronts SET user_id = 108 WHERE user_id = 6;
UPDATE storefront_staff SET user_id = 108 WHERE user_id = 6;
UPDATE storefront_inventory_movements SET user_id = 108 WHERE user_id = 6;
UPDATE user_storefronts SET user_id = 108 WHERE user_id = 6;
UPDATE notifications SET user_id = 108 WHERE user_id = 6;
UPDATE notification_settings SET user_id = 108 WHERE user_id = 6;
UPDATE reviews SET user_id = 108 WHERE user_id = 6;
UPDATE review_responses SET user_id = 108 WHERE user_id = 6;
UPDATE review_votes SET user_id = 108 WHERE user_id = 6;
UPDATE user_roles SET user_id = 108 WHERE user_id = 6;
UPDATE user_contacts SET user_id = 108 WHERE user_id = 6;
UPDATE user_contacts SET contact_user_id = 108 WHERE contact_user_id = 6;
UPDATE user_balances SET user_id = 108 WHERE user_id = 6;
UPDATE user_privacy_settings SET user_id = 108 WHERE user_id = 6;
UPDATE user_behavior_events SET user_id = 108 WHERE user_id = 6;
UPDATE user_telegram_connections SET user_id = 108 WHERE user_id = 6;
UPDATE subscription_history SET user_id = 108 WHERE user_id = 6;
UPDATE subscription_payments SET user_id = 108 WHERE user_id = 6;
UPDATE user_subscriptions SET user_id = 108 WHERE user_id = 6;
UPDATE balance_transactions SET user_id = 108 WHERE user_id = 6;
UPDATE payment_transactions SET user_id = 108 WHERE user_id = 6;
UPDATE escrow_payments SET seller_id = 108 WHERE seller_id = 6;
UPDATE escrow_payments SET buyer_id = 108 WHERE buyer_id = 6;
UPDATE merchant_payouts SET seller_id = 108 WHERE seller_id = 6;
UPDATE shopping_carts SET user_id = 108 WHERE user_id = 6;
UPDATE listing_views SET user_id = 108 WHERE user_id = 6;
UPDATE search_statistics SET user_id = 108 WHERE user_id = 6;
UPDATE gis_filter_analytics SET user_id = 108 WHERE user_id = 6;
UPDATE address_change_log SET user_id = 108 WHERE user_id = 6;
UPDATE translation_audit_log SET user_id = 108 WHERE user_id = 6;
UPDATE role_audit_log SET user_id = 108 WHERE user_id = 6;
UPDATE role_audit_log SET target_user_id = 108 WHERE target_user_id = 6;

-- 4. Удаляем старых пользователей
DELETE FROM users WHERE id IN (2, 6);

-- 5. Обновляем временные ID на финальные
UPDATE users SET id = 6, email = REPLACE(email, 'temp_', '') WHERE id = 106;
UPDATE users SET id = 8, email = REPLACE(email, 'temp_', '') WHERE id = 108;

-- 6. Обновляем все связанные таблицы на финальные ID
-- Для voroshilovdo@gmail.com (106 -> 6)
UPDATE marketplace_listings SET user_id = 6 WHERE user_id = 106;
UPDATE marketplace_chats SET seller_id = 6 WHERE seller_id = 106;
UPDATE marketplace_chats SET buyer_id = 6 WHERE buyer_id = 106;
UPDATE marketplace_messages SET sender_id = 6 WHERE sender_id = 106;
UPDATE marketplace_messages SET receiver_id = 6 WHERE receiver_id = 106;
UPDATE marketplace_favorites SET user_id = 6 WHERE user_id = 106;
UPDATE marketplace_orders SET seller_id = 6 WHERE seller_id = 106;
UPDATE marketplace_orders SET buyer_id = 6 WHERE buyer_id = 106;
UPDATE storefronts SET user_id = 6 WHERE user_id = 106;
UPDATE storefront_staff SET user_id = 6 WHERE user_id = 106;
UPDATE storefront_inventory_movements SET user_id = 6 WHERE user_id = 106;
UPDATE user_storefronts SET user_id = 6 WHERE user_id = 106;
UPDATE notifications SET user_id = 6 WHERE user_id = 106;
UPDATE notification_settings SET user_id = 6 WHERE user_id = 106;
UPDATE reviews SET user_id = 6 WHERE user_id = 106;
UPDATE review_responses SET user_id = 6 WHERE user_id = 106;
UPDATE review_votes SET user_id = 6 WHERE user_id = 106;
UPDATE user_roles SET user_id = 6 WHERE user_id = 106;
UPDATE user_contacts SET user_id = 6 WHERE user_id = 106;
UPDATE user_contacts SET contact_user_id = 6 WHERE contact_user_id = 106;
UPDATE user_balances SET user_id = 6 WHERE user_id = 106;
UPDATE user_privacy_settings SET user_id = 6 WHERE user_id = 106;
UPDATE user_behavior_events SET user_id = 6 WHERE user_id = 106;
UPDATE user_telegram_connections SET user_id = 6 WHERE user_id = 106;
UPDATE subscription_history SET user_id = 6 WHERE user_id = 106;
UPDATE subscription_payments SET user_id = 6 WHERE user_id = 106;
UPDATE user_subscriptions SET user_id = 6 WHERE user_id = 106;
UPDATE balance_transactions SET user_id = 6 WHERE user_id = 106;
UPDATE payment_transactions SET user_id = 6 WHERE user_id = 106;
UPDATE escrow_payments SET seller_id = 6 WHERE seller_id = 106;
UPDATE escrow_payments SET buyer_id = 6 WHERE buyer_id = 106;
UPDATE merchant_payouts SET seller_id = 6 WHERE seller_id = 106;
UPDATE shopping_carts SET user_id = 6 WHERE user_id = 106;
UPDATE listing_views SET user_id = 6 WHERE user_id = 106;
UPDATE search_statistics SET user_id = 6 WHERE user_id = 106;
UPDATE gis_filter_analytics SET user_id = 6 WHERE user_id = 106;
UPDATE address_change_log SET user_id = 6 WHERE user_id = 106;
UPDATE translation_audit_log SET user_id = 6 WHERE user_id = 106;
UPDATE role_audit_log SET user_id = 6 WHERE user_id = 106;
UPDATE role_audit_log SET target_user_id = 6 WHERE target_user_id = 106;

-- Для 4hash92@gmail.com (108 -> 8)
UPDATE marketplace_listings SET user_id = 8 WHERE user_id = 108;
UPDATE marketplace_chats SET seller_id = 8 WHERE seller_id = 108;
UPDATE marketplace_chats SET buyer_id = 8 WHERE buyer_id = 108;
UPDATE marketplace_messages SET sender_id = 8 WHERE sender_id = 108;
UPDATE marketplace_messages SET receiver_id = 8 WHERE receiver_id = 108;
UPDATE marketplace_favorites SET user_id = 8 WHERE user_id = 108;
UPDATE marketplace_orders SET seller_id = 8 WHERE seller_id = 108;
UPDATE marketplace_orders SET buyer_id = 8 WHERE buyer_id = 108;
UPDATE storefronts SET user_id = 8 WHERE user_id = 108;
UPDATE storefront_staff SET user_id = 8 WHERE user_id = 108;
UPDATE storefront_inventory_movements SET user_id = 8 WHERE user_id = 108;
UPDATE user_storefronts SET user_id = 8 WHERE user_id = 108;
UPDATE notifications SET user_id = 8 WHERE user_id = 108;
UPDATE notification_settings SET user_id = 8 WHERE user_id = 108;
UPDATE reviews SET user_id = 8 WHERE user_id = 108;
UPDATE review_responses SET user_id = 8 WHERE user_id = 108;
UPDATE review_votes SET user_id = 8 WHERE user_id = 108;
UPDATE user_roles SET user_id = 8 WHERE user_id = 108;
UPDATE user_contacts SET user_id = 8 WHERE user_id = 108;
UPDATE user_contacts SET contact_user_id = 8 WHERE contact_user_id = 108;
UPDATE user_balances SET user_id = 8 WHERE user_id = 108;
UPDATE user_privacy_settings SET user_id = 8 WHERE user_id = 108;
UPDATE user_behavior_events SET user_id = 8 WHERE user_id = 108;
UPDATE user_telegram_connections SET user_id = 8 WHERE user_id = 108;
UPDATE subscription_history SET user_id = 8 WHERE user_id = 108;
UPDATE subscription_payments SET user_id = 8 WHERE user_id = 108;
UPDATE user_subscriptions SET user_id = 8 WHERE user_id = 108;
UPDATE balance_transactions SET user_id = 8 WHERE user_id = 108;
UPDATE payment_transactions SET user_id = 8 WHERE user_id = 108;
UPDATE escrow_payments SET seller_id = 8 WHERE seller_id = 108;
UPDATE escrow_payments SET buyer_id = 8 WHERE buyer_id = 108;
UPDATE merchant_payouts SET seller_id = 8 WHERE seller_id = 108;
UPDATE shopping_carts SET user_id = 8 WHERE user_id = 108;
UPDATE listing_views SET user_id = 8 WHERE user_id = 108;
UPDATE search_statistics SET user_id = 8 WHERE user_id = 108;
UPDATE gis_filter_analytics SET user_id = 8 WHERE user_id = 108;
UPDATE address_change_log SET user_id = 8 WHERE user_id = 108;
UPDATE translation_audit_log SET user_id = 8 WHERE user_id = 108;
UPDATE role_audit_log SET user_id = 8 WHERE user_id = 108;
UPDATE role_audit_log SET target_user_id = 8 WHERE target_user_id = 108;

-- 7. Возвращаем EmailEmail@EmailEmail.ru с временного ID=208 обратно на ID=8
-- Теперь ID=8 свободен, так как 4hash92@gmail.com уже на нем
-- Поэтому перемещаем EmailEmail на новый свободный ID
UPDATE users SET id = 12 WHERE id = 208 AND email = 'EmailEmail@EmailEmail.ru';
UPDATE marketplace_listings SET user_id = 12 WHERE user_id = 208;
UPDATE marketplace_chats SET seller_id = 12 WHERE seller_id = 208;
UPDATE marketplace_chats SET buyer_id = 12 WHERE buyer_id = 208;
UPDATE marketplace_messages SET sender_id = 12 WHERE sender_id = 208;
UPDATE marketplace_messages SET receiver_id = 12 WHERE receiver_id = 208;
UPDATE marketplace_favorites SET user_id = 12 WHERE user_id = 208;
UPDATE marketplace_orders SET seller_id = 12 WHERE seller_id = 208;
UPDATE marketplace_orders SET buyer_id = 12 WHERE buyer_id = 208;
UPDATE storefronts SET user_id = 12 WHERE user_id = 208;
UPDATE notifications SET user_id = 12 WHERE user_id = 208;
UPDATE notification_settings SET user_id = 12 WHERE user_id = 208;
UPDATE reviews SET user_id = 12 WHERE user_id = 208;
UPDATE review_responses SET user_id = 12 WHERE user_id = 208;
UPDATE review_votes SET user_id = 12 WHERE user_id = 208;
UPDATE user_roles SET user_id = 12 WHERE user_id = 208;
UPDATE user_contacts SET user_id = 12 WHERE user_id = 208;
UPDATE user_contacts SET contact_user_id = 12 WHERE contact_user_id = 208;
UPDATE user_balances SET user_id = 12 WHERE user_id = 208;
UPDATE user_privacy_settings SET user_id = 12 WHERE user_id = 208;
UPDATE user_behavior_events SET user_id = 12 WHERE user_id = 208;
UPDATE user_telegram_connections SET user_id = 12 WHERE user_id = 208;
UPDATE subscription_history SET user_id = 12 WHERE user_id = 208;
UPDATE subscription_payments SET user_id = 12 WHERE user_id = 208;
UPDATE user_subscriptions SET user_id = 12 WHERE user_id = 208;
UPDATE balance_transactions SET user_id = 12 WHERE user_id = 208;
UPDATE payment_transactions SET user_id = 12 WHERE user_id = 208;
UPDATE escrow_payments SET seller_id = 12 WHERE seller_id = 208;
UPDATE escrow_payments SET buyer_id = 12 WHERE buyer_id = 208;
UPDATE merchant_payouts SET seller_id = 12 WHERE seller_id = 208;
UPDATE shopping_carts SET user_id = 12 WHERE user_id = 208;
UPDATE listing_views SET user_id = 12 WHERE user_id = 208;
UPDATE search_statistics SET user_id = 12 WHERE user_id = 208;
UPDATE gis_filter_analytics SET user_id = 12 WHERE user_id = 208;
UPDATE address_change_log SET user_id = 12 WHERE user_id = 208;
UPDATE translation_audit_log SET user_id = 12 WHERE user_id = 208;
UPDATE role_audit_log SET user_id = 12 WHERE user_id = 208;
UPDATE role_audit_log SET target_user_id = 12 WHERE target_user_id = 208;
UPDATE storefront_staff SET user_id = 12 WHERE user_id = 208;
UPDATE storefront_inventory_movements SET user_id = 12 WHERE user_id = 208;
UPDATE user_storefronts SET user_id = 12 WHERE user_id = 208;

-- 8. Сбрасываем последовательность для правильного следующего ID
SELECT setval('users_id_seq', GREATEST(
    (SELECT MAX(id) FROM users),
    12
));

-- Включаем обратно ограничения внешних ключей
SET session_replication_role = 'origin';

-- Проверяем результат
DO $$
BEGIN
    -- Проверяем что пользователи имеют правильные ID
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = 6 AND email = 'voroshilovdo@gmail.com') THEN
        RAISE EXCEPTION 'Failed to update voroshilovdo@gmail.com to ID=6';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE id = 8 AND email = '4hash92@gmail.com') THEN
        RAISE EXCEPTION 'Failed to update 4hash92@gmail.com to ID=8';
    END IF;

    RAISE NOTICE 'User IDs successfully synchronized with auth service';
END $$;

COMMIT;