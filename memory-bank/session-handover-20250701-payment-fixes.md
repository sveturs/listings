# Session Handover - Payment System Fixes
Date: 2025-07-01
Status: Completed

## Summary
Исправлены критические проблемы в платежной системе: удален хардкод localhost, исправлена поддержка локализации, исправлен WebSocket URL.

## Completed Tasks

### 1. Исправлен хардкод localhost в платежной системе
**Проблема**: В MockPaymentService был захардкожен URL `http://localhost:3001`
**Решение**: 
- Модифицирован `MockPaymentService` для приема `frontendURL` из конфигурации
- Обновлен `service.go` для передачи `cfg.FrontendURL` в конструктор
- Теперь URL генерируется динамически на основе переменной окружения `FRONTEND_URL`

### 2. Добавлена поддержка динамических локалей
**Проблема**: Payment URL всегда генерировался с `/en/`, игнорируя локаль пользователя
**Решение**:
- Добавлено поле `ReturnURL` в `DepositRequest` структуру
- Реализован метод `extractLocaleFromURL` для извлечения локали из return_url
- Обновлены методы `CreatePaymentSession` и `CreateOrderPayment` для использования правильной локали
- Теперь система поддерживает как `/en/`, так и `/ru/` интерфейсы

### 3. Исправлен WebSocket URL
**Проблема**: WebSocket пытался подключаться к `ws://localhost:3001` (фронтенд) вместо бэкенда
**Решение**:
- Обновлен `chat.ts` для использования `config.api.websocketUrl` из конфигурации
- Удалена логика с `window.location.host` для development окружения
- Теперь WebSocket правильно подключается к `ws://localhost:3000`

### 4. Удален фоллбэк на английскую локаль
**Проблема**: В `balanceService.ts` был хардкод фоллбэка на `/en/balance/deposit/success`
**Решение**:
- Удален автоматический фоллбэк
- Теперь `return_url` всегда передается явно из компонента

## Technical Changes

### Backend
- `/backend/internal/proj/payments/service/mock_service.go`:
  - Добавлено поле `frontendURL` в структуру
  - Добавлен метод `extractLocaleFromURL`
  - Обновлены методы генерации payment URL

- `/backend/internal/proj/balance/handler/balance.go`:
  - Добавлено поле `ReturnURL` в `DepositRequest`
  - Добавлена передача `return_url` через контекст

- `/backend/internal/proj/global/service/service.go`:
  - Обновлен вызов `NewMockPaymentService` с передачей `cfg.FrontendURL`

### Frontend
- `/frontend/svetu/src/services/chat.ts`:
  - Исправлена логика подключения WebSocket
  - Теперь используется `config.api.websocketUrl`

- `/frontend/svetu/src/services/balance.ts`:
  - Удален хардкод фоллбэка на `/en/`

## Testing Notes
1. Депозит теперь работает корректно для обеих локалей (en/ru)
2. WebSocket успешно подключается к правильному серверу
3. В логах видно правильное извлечение локали:
   - `MockPaymentService: Extracted locale 'ru' from return_url`
   - `MockPaymentService: Extracted locale 'en' from return_url`

## Environment Variables
- `FRONTEND_URL` - используется для генерации payment URLs (например: `http://localhost:3001` или `https://svetu.rs`)
- `NEXT_PUBLIC_WEBSOCKET_URL` - WebSocket URL для чата (например: `ws://localhost:3000`)

## Future Considerations
1. В production окружении нужно правильно настроить `FRONTEND_URL`
2. Рассмотреть возможность передачи локали через API параметры вместо парсинга URL
3. Rate limiting для WebSocket может требовать настройки в production

## Current State
- ✅ Платежная система полностью функциональна
- ✅ Поддержка мультиязычности работает корректно
- ✅ WebSocket подключение стабильно
- ✅ Конфигурация через переменные окружения