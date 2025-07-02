-- Добавляем недостающие колонки в таблицу storefront_orders

-- Платежные поля
ALTER TABLE storefront_orders 
ADD COLUMN IF NOT EXISTS payment_method VARCHAR(50) NOT NULL DEFAULT 'allsecure',
ADD COLUMN IF NOT EXISTS payment_status VARCHAR(50) NOT NULL DEFAULT 'pending',
ADD COLUMN IF NOT EXISTS notes TEXT,
ADD COLUMN IF NOT EXISTS metadata JSONB DEFAULT '{}';

-- Добавляем колонку discount
ALTER TABLE storefront_orders
ADD COLUMN IF NOT EXISTS discount DECIMAL(12,2) DEFAULT 0;

-- Создаем алиасы для совместимости с моделью
-- Используем views для прозрачного маппинга
CREATE OR REPLACE VIEW storefront_orders_view AS
SELECT 
    id,
    order_number,
    storefront_id,
    customer_id,
    payment_transaction_id,
    subtotal_amount as subtotal,
    shipping_amount as shipping,
    tax_amount as tax,
    discount,
    total_amount as total,
    commission_amount,
    seller_amount,
    currency,
    status,
    escrow_release_date,
    escrow_days,
    shipping_address,
    billing_address,
    shipping_method,
    shipping_provider,
    tracking_number,
    customer_notes,
    seller_notes,
    payment_method,
    payment_status,
    notes,
    metadata,
    confirmed_at,
    shipped_at,
    delivered_at,
    cancelled_at,
    created_at,
    updated_at
FROM storefront_orders;

-- Комментарии к новым колонкам
COMMENT ON COLUMN storefront_orders.payment_method IS 'Метод оплаты: allsecure, mock_payment и т.д.';
COMMENT ON COLUMN storefront_orders.payment_status IS 'Статус платежа: pending, processing, completed, failed, refunded';
COMMENT ON COLUMN storefront_orders.notes IS 'Общие заметки к заказу';
COMMENT ON COLUMN storefront_orders.metadata IS 'Дополнительные метаданные заказа в формате JSON';
COMMENT ON COLUMN storefront_orders.discount IS 'Сумма скидки на заказ';