-- Откат системы уведомлений

DROP TRIGGER IF EXISTS update_user_notification_preferences_updated_at ON user_notification_preferences;
DROP TRIGGER IF EXISTS update_notification_templates_updated_at ON notification_templates;
DROP TRIGGER IF EXISTS update_delivery_notifications_updated_at ON delivery_notifications;

DROP TABLE IF EXISTS user_notification_contacts CASCADE;
DROP TABLE IF EXISTS user_notification_preferences CASCADE;
DROP TABLE IF EXISTS notification_templates CASCADE;
DROP TABLE IF EXISTS delivery_notifications CASCADE;

DROP FUNCTION IF EXISTS update_updated_at_column() CASCADE;