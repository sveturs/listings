-- Добавляем mock payment method для тестирования и разработки
INSERT INTO payment_methods (
    code, 
    name, 
    type, 
    is_active, 
    minimum_amount, 
    maximum_amount, 
    fee_percentage, 
    fixed_fee, 
    created_at
)
VALUES (
    'mock_payment', 
    'Mock Payment (для тестирования)', 
    'card', 
    true, 
    100, 
    1000000, 
    0, 
    0, 
    NOW()
)
ON CONFLICT (code) DO UPDATE 
SET 
    name = EXCLUDED.name,
    type = EXCLUDED.type,
    is_active = EXCLUDED.is_active,
    minimum_amount = EXCLUDED.minimum_amount,
    maximum_amount = EXCLUDED.maximum_amount,
    fee_percentage = EXCLUDED.fee_percentage,
    fixed_fee = EXCLUDED.fixed_fee;

-- Добавляем комментарий для документации
COMMENT ON COLUMN payment_methods.code IS 'Уникальный код метода оплаты. mock_payment используется для тестирования';