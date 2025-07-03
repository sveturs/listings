# Session 2025-07-01 - Payment System Fixed

## Проблема
Баланс показывал 0 после успешной оплаты на странице `/balance/deposit/success`

## Решение
1. **Миграция 000064**: добавлен mock_payment в payment_methods
2. **Frontend fixes**:
   - Mock payment button теперь вызывает API `/balance/mock/complete`
   - Токен берется из sessionStorage (svetu_access_token)
   - BalanceWidget обновляется 4 раза (500ms, 1.5s, 3s, 5s)

## Ключевые файлы
- `/backend/migrations/000064_add_mock_payment_method.up.sql`
- `/frontend/svetu/src/app/[locale]/payment/mock/page.tsx` (строки 257-290)
- `/frontend/svetu/src/app/[locale]/balance/deposit/success/page.tsx` (строки 31-42)

## Тестовый пользователь
- Email: test@user.rs
- Password: testuser
- Баланс после теста: 5000 RSD

## Команды
```bash
# Frontend: yarn dev -p 3001
# Backend: go run ./cmd/api/main.go
# Миграция уже применена в БД
```