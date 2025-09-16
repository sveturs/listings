-- Откат синхронизации ID пользователей
-- Возвращаем:
-- ID=2 для voroshilovdo@gmail.com (был ID=6 после синхронизации)
-- ID=6 для 4hash92@gmail.com (был ID=8 после синхронизации)

BEGIN;

-- Временно отключаем ограничения внешних ключей
SET session_replication_role = 'replica';

-- 1. Создаем временных пользователей
INSERT INTO users (id, email, name, provider, created_at, updated_at, old_email)
SELECT
    CASE
        WHEN id = 6 AND email = 'voroshilovdo@gmail.com' THEN 102  -- временный для возврата к ID=2
        WHEN id = 8 AND email = '4hash92@gmail.com' THEN 106  -- временный для возврата к ID=6
    END as temp_id,
    CONCAT('rollback_', email) as temp_email,
    name,
    provider,
    created_at,
    updated_at,
    email as old_email
FROM users
WHERE (id = 6 AND email = 'voroshilovdo@gmail.com')
   OR (id = 8 AND email = '4hash92@gmail.com');

-- 2. Обновляем все связи для voroshilovdo@gmail.com (6 -> 102)
UPDATE marketplace_listings SET user_id = 102 WHERE user_id = 6;
UPDATE marketplace_chats SET seller_id = 102 WHERE seller_id = 6;
UPDATE marketplace_chats SET buyer_id = 102 WHERE buyer_id = 6;
UPDATE marketplace_messages SET sender_id = 102 WHERE sender_id = 6;
UPDATE marketplace_messages SET receiver_id = 102 WHERE receiver_id = 6;
UPDATE marketplace_favorites SET user_id = 102 WHERE user_id = 6;
UPDATE marketplace_orders SET seller_id = 102 WHERE seller_id = 6;
UPDATE marketplace_orders SET buyer_id = 102 WHERE buyer_id = 6;
UPDATE storefronts SET user_id = 102 WHERE user_id = 6;
UPDATE storefront_staff SET user_id = 102 WHERE user_id = 6;
UPDATE storefront_inventory_movements SET user_id = 102 WHERE user_id = 6;
UPDATE user_storefronts SET user_id = 102 WHERE user_id = 6;
UPDATE notifications SET user_id = 102 WHERE user_id = 6;
UPDATE notification_settings SET user_id = 102 WHERE user_id = 6;
UPDATE reviews SET user_id = 102 WHERE user_id = 6;
UPDATE review_responses SET user_id = 102 WHERE user_id = 6;
UPDATE review_votes SET user_id = 102 WHERE user_id = 6;
UPDATE user_roles SET user_id = 102 WHERE user_id = 6;
UPDATE user_contacts SET user_id = 102 WHERE user_id = 6;
UPDATE user_contacts SET contact_user_id = 102 WHERE contact_user_id = 6;
UPDATE user_balances SET user_id = 102 WHERE user_id = 6;
UPDATE user_privacy_settings SET user_id = 102 WHERE user_id = 6;
UPDATE user_behavior_events SET user_id = 102 WHERE user_id = 6;
UPDATE user_telegram_connections SET user_id = 102 WHERE user_id = 6;
UPDATE subscription_history SET user_id = 102 WHERE user_id = 6;
UPDATE subscription_payments SET user_id = 102 WHERE user_id = 6;
UPDATE user_subscriptions SET user_id = 102 WHERE user_id = 6;
UPDATE balance_transactions SET user_id = 102 WHERE user_id = 6;
UPDATE payment_transactions SET user_id = 102 WHERE user_id = 6;
UPDATE escrow_payments SET seller_id = 102 WHERE seller_id = 6;
UPDATE escrow_payments SET buyer_id = 102 WHERE buyer_id = 6;
UPDATE merchant_payouts SET seller_id = 102 WHERE seller_id = 6;
UPDATE shopping_carts SET user_id = 102 WHERE user_id = 6;
UPDATE listing_views SET user_id = 102 WHERE user_id = 6;
UPDATE search_statistics SET user_id = 102 WHERE user_id = 6;
UPDATE gis_filter_analytics SET user_id = 102 WHERE user_id = 6;
UPDATE address_change_log SET user_id = 102 WHERE user_id = 6;
UPDATE translation_audit_log SET user_id = 102 WHERE user_id = 6;
UPDATE role_audit_log SET user_id = 102 WHERE user_id = 6;
UPDATE role_audit_log SET target_user_id = 102 WHERE target_user_id = 6;

-- 3. Обновляем все связи для 4hash92@gmail.com (8 -> 106)
UPDATE marketplace_listings SET user_id = 106 WHERE user_id = 8;
UPDATE marketplace_chats SET seller_id = 106 WHERE seller_id = 8;
UPDATE marketplace_chats SET buyer_id = 106 WHERE buyer_id = 8;
UPDATE marketplace_messages SET sender_id = 106 WHERE sender_id = 8;
UPDATE marketplace_messages SET receiver_id = 106 WHERE receiver_id = 8;
UPDATE marketplace_favorites SET user_id = 106 WHERE user_id = 8;
UPDATE marketplace_orders SET seller_id = 106 WHERE seller_id = 8;
UPDATE marketplace_orders SET buyer_id = 106 WHERE buyer_id = 8;
UPDATE storefronts SET user_id = 106 WHERE user_id = 8;
UPDATE storefront_staff SET user_id = 106 WHERE user_id = 8;
UPDATE storefront_inventory_movements SET user_id = 106 WHERE user_id = 8;
UPDATE user_storefronts SET user_id = 106 WHERE user_id = 8;
UPDATE notifications SET user_id = 106 WHERE user_id = 8;
UPDATE notification_settings SET user_id = 106 WHERE user_id = 8;
UPDATE reviews SET user_id = 106 WHERE user_id = 8;
UPDATE review_responses SET user_id = 106 WHERE user_id = 8;
UPDATE review_votes SET user_id = 106 WHERE user_id = 8;
UPDATE user_roles SET user_id = 106 WHERE user_id = 8;
UPDATE user_contacts SET user_id = 106 WHERE user_id = 8;
UPDATE user_contacts SET contact_user_id = 106 WHERE contact_user_id = 8;
UPDATE user_balances SET user_id = 106 WHERE user_id = 8;
UPDATE user_privacy_settings SET user_id = 106 WHERE user_id = 8;
UPDATE user_behavior_events SET user_id = 106 WHERE user_id = 8;
UPDATE user_telegram_connections SET user_id = 106 WHERE user_id = 8;
UPDATE subscription_history SET user_id = 106 WHERE user_id = 8;
UPDATE subscription_payments SET user_id = 106 WHERE user_id = 8;
UPDATE user_subscriptions SET user_id = 106 WHERE user_id = 8;
UPDATE balance_transactions SET user_id = 106 WHERE user_id = 8;
UPDATE payment_transactions SET user_id = 106 WHERE user_id = 8;
UPDATE escrow_payments SET seller_id = 106 WHERE seller_id = 8;
UPDATE escrow_payments SET buyer_id = 106 WHERE buyer_id = 8;
UPDATE merchant_payouts SET seller_id = 106 WHERE seller_id = 8;
UPDATE shopping_carts SET user_id = 106 WHERE user_id = 8;
UPDATE listing_views SET user_id = 106 WHERE user_id = 8;
UPDATE search_statistics SET user_id = 106 WHERE user_id = 8;
UPDATE gis_filter_analytics SET user_id = 106 WHERE user_id = 8;
UPDATE address_change_log SET user_id = 106 WHERE user_id = 8;
UPDATE translation_audit_log SET user_id = 106 WHERE user_id = 8;
UPDATE role_audit_log SET user_id = 106 WHERE user_id = 8;
UPDATE role_audit_log SET target_user_id = 106 WHERE target_user_id = 8;

-- 4. Удаляем старых пользователей
DELETE FROM users WHERE id = 6 AND email = 'voroshilovdo@gmail.com';
DELETE FROM users WHERE id = 8 AND email = '4hash92@gmail.com';

-- 5. Восстанавливаем пользователей с оригинальными ID
UPDATE users SET id = 2, email = REPLACE(email, 'rollback_', '') WHERE id = 102;
UPDATE users SET id = 6, email = REPLACE(email, 'rollback_', '') WHERE id = 106;

-- 6. Финальное обновление всех связей
-- Для voroshilovdo@gmail.com (102 -> 2)
UPDATE marketplace_listings SET user_id = 2 WHERE user_id = 102;
UPDATE marketplace_chats SET seller_id = 2 WHERE seller_id = 102;
UPDATE marketplace_chats SET buyer_id = 2 WHERE buyer_id = 102;
UPDATE marketplace_messages SET sender_id = 2 WHERE sender_id = 102;
UPDATE marketplace_messages SET receiver_id = 2 WHERE receiver_id = 102;
UPDATE marketplace_favorites SET user_id = 2 WHERE user_id = 102;
UPDATE marketplace_orders SET seller_id = 2 WHERE seller_id = 102;
UPDATE marketplace_orders SET buyer_id = 2 WHERE buyer_id = 102;
UPDATE storefronts SET user_id = 2 WHERE user_id = 102;
UPDATE storefront_staff SET user_id = 2 WHERE user_id = 102;
UPDATE storefront_inventory_movements SET user_id = 2 WHERE user_id = 102;
UPDATE user_storefronts SET user_id = 2 WHERE user_id = 102;
UPDATE notifications SET user_id = 2 WHERE user_id = 102;
UPDATE notification_settings SET user_id = 2 WHERE user_id = 102;
UPDATE reviews SET user_id = 2 WHERE user_id = 102;
UPDATE review_responses SET user_id = 2 WHERE user_id = 102;
UPDATE review_votes SET user_id = 2 WHERE user_id = 102;
UPDATE user_roles SET user_id = 2 WHERE user_id = 102;
UPDATE user_contacts SET user_id = 2 WHERE user_id = 102;
UPDATE user_contacts SET contact_user_id = 2 WHERE contact_user_id = 102;
UPDATE user_balances SET user_id = 2 WHERE user_id = 102;
UPDATE user_privacy_settings SET user_id = 2 WHERE user_id = 102;
UPDATE user_behavior_events SET user_id = 2 WHERE user_id = 102;
UPDATE user_telegram_connections SET user_id = 2 WHERE user_id = 102;
UPDATE subscription_history SET user_id = 2 WHERE user_id = 102;
UPDATE subscription_payments SET user_id = 2 WHERE user_id = 102;
UPDATE user_subscriptions SET user_id = 2 WHERE user_id = 102;
UPDATE balance_transactions SET user_id = 2 WHERE user_id = 102;
UPDATE payment_transactions SET user_id = 2 WHERE user_id = 102;
UPDATE escrow_payments SET seller_id = 2 WHERE seller_id = 102;
UPDATE escrow_payments SET buyer_id = 2 WHERE buyer_id = 102;
UPDATE merchant_payouts SET seller_id = 2 WHERE seller_id = 102;
UPDATE shopping_carts SET user_id = 2 WHERE user_id = 102;
UPDATE listing_views SET user_id = 2 WHERE user_id = 102;
UPDATE search_statistics SET user_id = 2 WHERE user_id = 102;
UPDATE gis_filter_analytics SET user_id = 2 WHERE user_id = 102;
UPDATE address_change_log SET user_id = 2 WHERE user_id = 102;
UPDATE translation_audit_log SET user_id = 2 WHERE user_id = 102;
UPDATE role_audit_log SET user_id = 2 WHERE user_id = 102;
UPDATE role_audit_log SET target_user_id = 2 WHERE target_user_id = 102;

-- Для 4hash92@gmail.com (106 -> 6)
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

-- Включаем обратно ограничения
SET session_replication_role = 'origin';

-- Проверяем результат отката
DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM users WHERE id = 2 AND email = 'voroshilovdo@gmail.com') THEN
        RAISE EXCEPTION 'Failed to rollback voroshilovdo@gmail.com to ID=2';
    END IF;

    IF NOT EXISTS (SELECT 1 FROM users WHERE id = 6 AND email = '4hash92@gmail.com') THEN
        RAISE EXCEPTION 'Failed to rollback 4hash92@gmail.com to ID=6';
    END IF;

    RAISE NOTICE 'User IDs successfully rolled back to original values';
END $$;

COMMIT;