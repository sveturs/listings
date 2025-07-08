-- Удаляем mock payment method
DELETE FROM payment_methods WHERE code = 'mock_payment';

-- Удаляем комментарий
COMMENT ON COLUMN payment_methods.code IS NULL;