# Инструкция по исправлению ошибок 403 в админских разделах

## Проблема
Ошибка 403 (Forbidden) при доступе к админским API endpoints, даже если пользователь авторизован.

## Причины
1. **Неправильные параметры запроса** - API ожидает одни параметры, а фронтенд отправляет другие
2. **Токен не передается в заголовках** - tokenManager не инициализирован
3. **Пользователь не в списке админов** - user_id отсутствует в hardcoded списке

## Решение

### 1. Проверка параметров запроса
```bash
# Проверить какие параметры ожидает backend через swagger
cd /data/hostel-booking-system/backend/docs && python3 -m http.server 8888
# Затем через JSON MCP проверить параметры endpoint
```

Пример исправления:
- Было: `period_start`, `period_end`
- Стало: `date_from`, `date_to`

### 2. Исправление авторизации на фронтенде

В сервисе добавить асинхронный импорт tokenManager:

```typescript
// Было
import { tokenManager } from '@/utils/tokenManager';
const token = tokenManager.getAccessToken();

// Стало
if (typeof window !== 'undefined') {
  const { tokenManager } = await import('@/utils/tokenManager');
  tokenManager.initializeFromStorage(); // Важно!
  const accessToken = tokenManager.getAccessToken();
  if (accessToken) {
    headers['Authorization'] = `Bearer ${accessToken}`;
  }
}
```

### 3. Добавление пользователя в список админов

Найти handler в backend и добавить user_id:

```go
// Файл: /backend/internal/proj/{module}/handler/{handler}.go

// Было
adminUsers := []int{1, 2, 3, 4, 5}

// Стало
adminUsers := []int{1, 2, 3, 4, 5, 6}
```

### 4. Перезапуск backend
```bash
/home/dim/.local/bin/kill-port-3000.sh && \
screen -dmS backend-3000 bash -c 'cd /data/hostel-booking-system/backend && go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'
```

## Список файлов с админскими проверками

Для быстрого поиска где нужно добавить user_id:
```bash
grep -r "adminUsers.*\[\]int" /data/hostel-booking-system/backend
```

Основные модули с админской проверкой:
- `/backend/internal/proj/analytics/handler/`
- `/backend/internal/proj/search_admin/handler/`
- `/backend/internal/proj/admin/handler/`
- `/backend/internal/proj/logistics/handler/`