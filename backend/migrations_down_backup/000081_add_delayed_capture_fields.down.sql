-- Откат миграции
DROP TRIGGER IF EXISTS update_marketplace_orders_updated_at ON marketplace_orders;
DROP FUNCTION IF EXISTS update_updated_at_column();

DROP TABLE IF EXISTS order_messages;
DROP TABLE IF EXISTS order_status_history;
DROP TABLE IF EXISTS marketplace_orders;

DROP INDEX IF EXISTS idx_payment_transactions_auto_capture;

ALTER TABLE payment_transactions 
DROP COLUMN IF EXISTS capture_mode,
DROP COLUMN IF EXISTS auto_capture_at,
DROP COLUMN IF EXISTS capture_deadline_at,
DROP COLUMN IF EXISTS capture_attempted_at,
DROP COLUMN IF EXISTS capture_attempts;