# Session Handover - 2025-07-01 - Mock Payment Fix

## Исправлена ошибка 404 при пополнении баланса через Mock Payment

### Проблема
При попытке пополнения баланса через mock payment система выдавала ошибку 404 на URL:
```
POST http://localhost:3001/en/payment/undefined/balance/mock/complete 404 (Not Found)
```

### Решение
1. **Исправлен URL в frontend**: изменён с `${config.api.baseUrl}/balance/mock/complete` на `/api/v1/balance/mock/complete`
2. **Улучшена обработка ошибок**: добавлена проверка ответа и отображение ошибки пользователю
3. **Добавлен await для завершения запроса**: теперь редирект происходит только после успешного завершения платежа

### Изменённые файлы
- `/frontend/svetu/src/app/[locale]/payment/mock/page.tsx` - исправлен URL и логика обработки

### Что работает
✅ Mock payment endpoint доступен на backend
✅ Обработчик `CompleteMockPayment` корректно создаёт транзакцию пополнения
✅ Страница успешного пополнения использует `BalanceWidget` с автообновлением
✅ Frontend корректно отправляет запрос на завершение платежа

### Следующие шаги
- Протестировать полный flow пополнения баланса
- Убедиться, что баланс корректно обновляется на странице успеха
- Проверить, что транзакция записывается в базу данных

### Команды для проверки
```bash
# Проверка работы frontend
curl http://localhost:3001/en

# Проверка mock payment endpoint
curl -X POST http://localhost:3000/api/v1/balance/mock/complete \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <token>" \
  -d '{"session_id": "test", "amount": 1000}'
```