# Документация по переменным окружения

## Обзор
Этот документ описывает все переменные окружения, используемые frontend приложением.

## Категории переменных

### 1. Публичные переменные (NEXT_PUBLIC_*)
Эти переменные доступны в браузере и могут использоваться в клиентском коде.

| Переменная | Обязательная | По умолчанию | Описание |
|------------|--------------|--------------|----------|
| `NEXT_PUBLIC_API_URL` | Да | `http://localhost:3000` | URL публичного API |
| `NEXT_PUBLIC_MINIO_URL` | Да | `http://localhost:9000` | URL MinIO/S3 хранилища |
| `NEXT_PUBLIC_IMAGE_HOSTS` | Нет | См. .env.example | Разрешенные домены для изображений |
| `NEXT_PUBLIC_IMAGE_PATH_PATTERN` | Нет | `/listings/**` | Паттерны путей для изображений |
| `NEXT_PUBLIC_WEBSOCKET_URL` | Нет | - | WebSocket endpoint для real-time функций |
| `NEXT_PUBLIC_ENABLE_PAYMENTS` | Нет | `false` | Включить платежные функции |
| `NEXT_PUBLIC_GOOGLE_CLIENT_ID` | Нет | - | Google OAuth client ID |

### 2. Серверные переменные
Эти переменные доступны только в серверном коде (API routes, SSR).

| Переменная | Обязательная | По умолчанию | Описание |
|------------|--------------|--------------|----------|
| `INTERNAL_API_URL` | Нет | - | Внутренний URL API для Docker/K8s |
| `NODE_ENV` | Да | `development` | Окружение Node.js |
| `PORT` | Нет | `3000` | Порт сервера |

### 3. Переменные сборки
Эти переменные влияют на процесс сборки и оптимизацию.

| Переменная | Обязательная | По умолчанию | Описание |
|------------|--------------|--------------|----------|
| `NEXT_TELEMETRY_DISABLED` | Нет | `0` | Отключить телеметрию Next.js |

## Конфигурации для разных окружений

### Разработка (Development)
```bash
NEXT_PUBLIC_API_URL=http://localhost:3000
NEXT_PUBLIC_ENABLE_DEBUG=true
NODE_ENV=development
```

### Промежуточное окружение (Staging)
```bash
NEXT_PUBLIC_API_URL=https://staging-api.svetu.rs
NEXT_PUBLIC_ENABLE_DEBUG=true
NODE_ENV=production
```

### Продакшн (Production)
```bash
NEXT_PUBLIC_API_URL=https://api.svetu.rs
NEXT_PUBLIC_ENABLE_DEBUG=false
NODE_ENV=production
```

## Безопасность

1. **Никогда не коммитьте реальные значения** для production переменных
2. **Используйте системы управления секретами** для чувствительных данных (API ключи, токены)
3. **Валидируйте все переменные** при запуске для раннего обнаружения ошибок
4. **Минимизируйте публичные переменные** для уменьшения поверхности атаки

## Настройка Docker/Kubernetes

### Docker Compose
```yaml
services:
  frontend:
    environment:
      - NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}
      - INTERNAL_API_URL=http://backend:3000
```

### Kubernetes ConfigMap
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: frontend-config
data:
  NEXT_PUBLIC_API_URL: "https://api.svetu.rs"
  NEXT_PUBLIC_ENABLE_PAYMENTS: "true"
```

## Устранение неполадок

### Переменная не обновляется
1. Перезапустите сервер разработки
2. Очистите кеш Next.js: `rm -rf .next`
3. Проверьте, что имя переменной начинается с `NEXT_PUBLIC_`

### Переменная не определена в production
1. Убедитесь, что переменная установлена в окружении развертывания
2. Проверьте конфигурацию Docker/K8s
3. Проверьте логи сборки на наличие предупреждений

### Ошибки типов с переменными окружения
1. Обновите определения типов в `config/types.ts`
2. Запустите валидацию для раннего обнаружения ошибок
3. Используйте значения по умолчанию для необязательных переменных