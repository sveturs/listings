# Текущая задача: Исправление behavior tracking

## Статус: ✅ ЗАВЕРШЕНА

### Проблема
Frontend отправлял batch событий на `/api/v1/analytics/track`, но backend:
1. Ожидал одиночное событие вместо batch
2. Требовал авторизацию для публичного endpoint

### Решение
1. ✅ Добавлена структура `TrackEventBatch` в `backend/internal/domain/behavior/types.go`
2. ✅ Обновлен handler для поддержки как batch, так и одиночных событий
3. ✅ Исправлен middleware для пропуска авторизации на `/api/v1/analytics/track`
4. ✅ Обновлена swagger документация
5. ✅ Протестирована работа API - события успешно принимаются

### Результат
- API endpoint `/api/v1/analytics/track` теперь корректно принимает:
  - Batch события от frontend (с полями `events`, `batch_id`, `created_at`)
  - Одиночные события для обратной совместимости
- Авторизация не требуется для записи событий
- События успешно сохраняются в базу данных

### Тестирование
```bash
# Batch события
curl -X POST http://localhost:3000/api/v1/analytics/track \
  -H "Content-Type: application/json" \
  -d '{"events":[{"event_type":"search_performed","session_id":"test123","search_query":"test"}],"batch_id":"batch123","created_at":"2025-07-08T15:00:00Z"}'

# Результат: {"data":{"batch_id":"batch123","failed_count":0,"message":"Events batch processed","processed_count":1},"success":true}
```