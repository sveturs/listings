# Конфигурация URL для Viber Bot интеграции

## Обзор

Viber Bot интеграция теперь поддерживает гибкую настройку URL через переменные окружения, что позволяет легко переключаться между локальной разработкой, dev, staging и production окружениями.

## Переменные окружения

### VIBER_PUBLIC_URL
Специальная переменная для URL, которые будут отображаться в Viber сообщениях. Может отличаться от `FRONTEND_URL` для удобства тестирования.

```bash
# Для локальной разработки (по умолчанию использует FRONTEND_URL)
# Не устанавливайте, будет использован FRONTEND_URL

# Для dev окружения
VIBER_PUBLIC_URL=https://dev.svetu.rs

# Для staging
VIBER_PUBLIC_URL=https://staging.svetu.rs

# Для production
VIBER_PUBLIC_URL=https://svetu.rs
```

### FRONTEND_URL
Основной URL фронтенда. Используется как fallback для `VIBER_PUBLIC_URL`.

```bash
# Локальная разработка
FRONTEND_URL=http://localhost:3001

# Dev окружение
FRONTEND_URL=https://dev.svetu.rs

# Production
FRONTEND_URL=https://svetu.rs
```

## Приоритет конфигурации

1. Если установлен `VIBER_PUBLIC_URL` - используется он
2. Если не установлен - используется `FRONTEND_URL`
3. Если оба не установлены - используется дефолт `https://svetu.rs`

## Примеры конфигурации

### Локальная разработка
```bash
# .env
FRONTEND_URL=http://localhost:3001
# VIBER_PUBLIC_URL не устанавливаем
# Viber будет отправлять ссылки на localhost:3001
```

### Dev сервер с тестированием через Viber
```bash
# .env
FRONTEND_URL=http://localhost:3001  # Для локальной разработки
VIBER_PUBLIC_URL=https://dev.svetu.rs  # Для ссылок в Viber
```

### Production
```bash
# .env
FRONTEND_URL=https://svetu.rs
VIBER_PUBLIC_URL=https://svetu.rs
```

## Использование в коде

Все URL в Viber сообщениях теперь генерируются динамически:

```go
// Вместо хардкода
msg := "Посетите: https://svetu.rs/product/123"

// Используем конфиг
msg := fmt.Sprintf("Посетите: %s/product/123", config.FrontendURL)
```

## Проверка конфигурации

После изменения переменных окружения:

1. Перезапустите backend сервер
2. Отправьте тестовое сообщение боту
3. Проверьте, что ссылки ведут на правильный домен

## Важно

- Всегда используйте `VIBER_PUBLIC_URL` для production и staging
- Не коммитьте локальные значения в репозиторий
- При деплое убедитесь, что переменные настроены правильно