# Quick Start Guide - Playwright Tests

## Варианты запуска тестов

### 1. Простая проверка работы Playwright

```bash
# Запуск простого теста для проверки установки
npx playwright test tests/example.spec.ts --config=playwright.config.simple.ts

# Просмотр отчета
npx playwright show-report
```

### 2. Локальная разработка (только frontend)

**Требования:**
- Frontend запущен на http://localhost:3001

```bash
# Запуск frontend (из другого терминала)
cd ../frontend/svetu
yarn dev -p 3001

# Запуск тестов
cd playwright
npx playwright test --config=playwright.config.local.ts

# Запуск в UI режиме для отладки
npx playwright test --ui --config=playwright.config.local.ts
```

### 3. Полная интеграция (backend + frontend)

**Требования:**
- PostgreSQL на порту 5433
- OpenSearch на порту 9201
- Backend на порту 3000
- Frontend на порту 3001

```bash
# Запуск всех тестов (автоматически запустит backend и frontend)
npx playwright test

# Запуск с отладкой
npx playwright test --debug
```

## Полезные команды

```bash
# Установка браузеров
yarn install:browsers

# Запуск конкретного теста
npx playwright test tests/auth/login.spec.ts

# Запуск в headed режиме (видимый браузер)
npx playwright test --headed

# Генерация кода тестов
npx playwright codegen http://localhost:3001

# Просмотр трейсов
npx playwright show-trace test-results/trace.zip
```

## Отладка тестов

1. **VS Code Extension**: Установите Playwright Test for VSCode
2. **Debug режим**: `npx playwright test --debug`
3. **UI режим**: `npx playwright test --ui`
4. **Трейсы**: Автоматически записываются при падении тестов

## Структура тестов

- `tests/auth/` - Тесты аутентификации
- `tests/marketplace/` - Тесты маркетплейса
- `tests/storefronts/` - Тесты витрин
- `helpers/` - Вспомогательные функции
- `fixtures/` - Тестовые файлы и данные

## Переменные окружения

Скопируйте `.env.test` и настройте под ваше окружение:

```bash
TEST_FRONTEND_URL=http://localhost:3001
TEST_BACKEND_URL=http://localhost:3000
TEST_DB_URL=postgres://test:test@localhost:5433/svetu_test
```