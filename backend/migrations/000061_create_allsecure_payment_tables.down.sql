-- Удаление таблиц для интеграции с AllSecure Payment Gateway

-- Удаляем добавленные колонки из существующей таблицы если она существует
DO $$
BEGIN
    IF EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'user_transactions') THEN
        ALTER TABLE user_transactions DROP COLUMN IF EXISTS payment_gateway;
        ALTER TABLE user_transactions DROP COLUMN IF EXISTS gateway_transaction_id;
        ALTER TABLE user_transactions DROP COLUMN IF EXISTS gateway_reference;
    END IF;
END $$;

-- Удаляем новые таблицы в обратном порядке зависимостей
DROP TABLE IF EXISTS merchant_payouts;
DROP TABLE IF EXISTS escrow_payments;
DROP TABLE IF EXISTS payment_transactions;
DROP TABLE IF EXISTS payment_gateways;