# Mock Payment Method Migration

## Описание
Миграция 000064 добавляет mock payment method для тестирования системы платежей без реальных транзакций.

## Что добавляется
- Mock payment method с кодом `mock_payment`
- Настройки: без комиссий, лимиты от 100 до 1,000,000 RSD
- Тип: card (для совместимости с UI)

## Использование

### Backend
При создании депозита указывайте:
```json
{
  "amount": 5000,
  "payment_method": "mock"
}
```

### Frontend
Mock платежи обрабатываются на странице `/payment/mock` с симуляцией процесса оплаты.

### Тестирование через API
```bash
# 1. Создать депозит
POST /api/v1/balance/deposit
{
  "amount": 5000,
  "payment_method": "mock"
}

# 2. Завершить mock платеж
POST /api/v1/balance/mock/complete
{
  "session_id": "mock_session_XXX",
  "amount": 5000
}
```

## Важно
- Mock payment method активен только в development окружении
- Не использовать на production!
- Все mock транзакции помечаются как "mock_payment" в базе данных