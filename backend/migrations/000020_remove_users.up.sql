-- Удаление materialized views зависящих от users
DROP MATERIALIZED VIEW IF EXISTS public.user_rating_distribution CASCADE;
DROP MATERIALIZED VIEW IF EXISTS public.user_rating_summary CASCADE;
DROP MATERIALIZED VIEW IF EXISTS public.user_ratings CASCADE;

-- Удаление всех foreign key constraints на таблицу users
ALTER TABLE IF EXISTS public.couriers DROP CONSTRAINT IF EXISTS couriers_user_id_fkey;
ALTER TABLE IF EXISTS public.delivery_notifications DROP CONSTRAINT IF EXISTS delivery_notifications_user_id_fkey;
ALTER TABLE IF EXISTS public.marketplace_favorites DROP CONSTRAINT IF EXISTS marketplace_favorites_user_id_fkey;
ALTER TABLE IF EXISTS public.marketplace_listings DROP CONSTRAINT IF EXISTS marketplace_listings_user_id_fkey;
ALTER TABLE IF EXISTS public.notification_settings DROP CONSTRAINT IF EXISTS notification_settings_user_id_fkey;
ALTER TABLE IF EXISTS public.notifications DROP CONSTRAINT IF EXISTS notifications_user_id_fkey;
ALTER TABLE IF EXISTS public.review_responses DROP CONSTRAINT IF EXISTS review_responses_user_id_fkey;
ALTER TABLE IF EXISTS public.review_votes DROP CONSTRAINT IF EXISTS review_votes_user_id_fkey;
ALTER TABLE IF EXISTS public.reviews DROP CONSTRAINT IF EXISTS reviews_user_id_fkey;
ALTER TABLE IF EXISTS public.role_audit_log DROP CONSTRAINT IF EXISTS role_audit_log_target_user_id_fkey;
ALTER TABLE IF EXISTS public.role_audit_log DROP CONSTRAINT IF EXISTS role_audit_log_user_id_fkey;
ALTER TABLE IF EXISTS public.saved_searches DROP CONSTRAINT IF EXISTS saved_searches_user_id_fkey;
ALTER TABLE IF EXISTS public.storefront_favorites DROP CONSTRAINT IF EXISTS storefront_favorites_user_id_fkey;
ALTER TABLE IF EXISTS public.subscription_history DROP CONSTRAINT IF EXISTS subscription_history_user_id_fkey;
ALTER TABLE IF EXISTS public.subscription_payments DROP CONSTRAINT IF EXISTS subscription_payments_user_id_fkey;
ALTER TABLE IF EXISTS public.tracking_websocket_connections DROP CONSTRAINT IF EXISTS tracking_websocket_connections_user_id_fkey;
ALTER TABLE IF EXISTS public.translation_audit_log DROP CONSTRAINT IF EXISTS translation_audit_log_user_id_fkey;
ALTER TABLE IF EXISTS public.user_car_view_history DROP CONSTRAINT IF EXISTS user_car_view_history_user_id_fkey;
ALTER TABLE IF EXISTS public.user_notification_contacts DROP CONSTRAINT IF EXISTS user_notification_contacts_user_id_fkey;
ALTER TABLE IF EXISTS public.user_notification_preferences DROP CONSTRAINT IF EXISTS user_notification_preferences_user_id_fkey;
ALTER TABLE IF EXISTS public.user_roles DROP CONSTRAINT IF EXISTS user_roles_user_id_fkey;
ALTER TABLE IF EXISTS public.user_subscriptions DROP CONSTRAINT IF EXISTS user_subscriptions_user_id_fkey;
ALTER TABLE IF EXISTS public.user_telegram_connections DROP CONSTRAINT IF EXISTS user_telegram_connections_user_id_fkey;
ALTER TABLE IF EXISTS public.user_view_history DROP CONSTRAINT IF EXISTS user_view_history_user_id_fkey;
ALTER TABLE IF EXISTS public.viber_users DROP CONSTRAINT IF EXISTS viber_users_user_id_fkey;
ALTER TABLE IF EXISTS public.vin_check_history DROP CONSTRAINT IF EXISTS vin_check_history_user_id_fkey;

-- Удаление таблиц
DROP TABLE IF EXISTS public.user_roles CASCADE;
DROP TABLE IF EXISTS public.users CASCADE;