# 🎭 E2E Tests Headless Configuration

**Дата создания:** 2025-10-19
**Версия:** 1.0
**Статус:** ✅ CONFIGURED - Ready for CI/CD

---

## 📋 Обзор

E2E и Accessibility тесты настроены для работы в headless режиме, что позволяет запускать их в CI/CD окружении без графического интерфейса.

## ✅ Что настроено

### 1. Playwright Configuration (`playwright.config.ts`)

**Headless режим:**
```typescript
use: {
  headless: process.env.HEADLESS !== 'false',  // По умолчанию headless
}
```

**Browser launch options для headless:**
```typescript
launchOptions: {
  args: [
    '--no-sandbox',                // Для Docker/CI окружения
    '--disable-setuid-sandbox',    // Для запуска без root
    '--disable-dev-shm-usage',     // Для ограниченной памяти
    '--disable-gpu',               // Для серверного окружения
  ],
}
```

### 2. Установленные зависимости

- ✅ Playwright v1.54.1
- ✅ Chromium browser с системными зависимостями
- ✅ @axe-core/playwright v4.10.2 (для accessibility тестов)

### 3. Конфигурация для CI

```typescript
// playwright.config.ts
{
  forbidOnly: !!process.env.CI,    // Запрещает test.only в CI
  retries: process.env.CI ? 2 : 0, // 2 retry в CI
  workers: process.env.CI ? 1 : 1, // 1 worker в CI
}
```

---

## 🚀 Запуск тестов

### Локально (Headless режим - по умолчанию)

```bash
cd /data/hostel-booking-system/frontend/svetu

# Запустить все E2E тесты
npx playwright test

# Запустить конкретный тест
npx playwright test e2e/user-journey-create-listing.spec.ts

# Запустить accessibility тесты
npx playwright test e2e/axe/
```

### Локально (С открытым браузером - для отладки)

```bash
# Отключить headless режим
HEADLESS=false npx playwright test

# Или запустить в UI режиме
npx playwright test --ui

# Или в debug режиме
npx playwright test --debug
```

### В CI/CD окружении

```bash
# Установить браузеры (один раз)
npx playwright install chromium --with-deps

# Запустить тесты
CI=true npx playwright test
```

---

## 📊 Типы тестов

### E2E Tests (3 теста)

**Файлы:**
- `e2e/user-journey-create-listing.spec.ts`
- `e2e/user-journey-search-contact.spec.ts`
- `e2e/admin-moderation-flow.spec.ts`

**Что тестируют:**
- Полный user journey: login → create listing → publish
- Поиск и связь с продавцом
- Админская модерация объявлений

### Accessibility Tests (2 теста)

**Файлы:**
- `e2e/axe/a11y-wcag-compliance.spec.ts`
- `e2e/axe/a11y-keyboard-navigation.spec.ts`

**Что тестируют:**
- WCAG 2.1 AA compliance через axe-core
- Keyboard navigation на всех интерактивных элементах

---

## ⚙️ Переменные окружения

```bash
# Headless режим (по умолчанию true)
HEADLESS=false              # Отключить headless режим

# Test credentials
TEST_ADMIN_EMAIL=admin@admin.rs
TEST_ADMIN_PASSWORD=P@$S4@dmi№

# URLs
NEXT_PUBLIC_FRONTEND_URL=http://localhost:3001
BASE_URL=http://localhost:3001

# CI флаг
CI=true                     # Включить CI режим (больше retries, строгие проверки)
```

---

## 🔧 Интеграция с Backend Test Runner

E2E тесты можно запускать через backend HTTP API:

```bash
# Запустить E2E suite через backend
TOKEN=$(cat /tmp/token)
curl -X POST "http://localhost:3000/api/v1/admin/tests/run" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"test_suite": "e2e"}'

# Результаты сохраняются в БД (test_runs, test_results, test_logs)
```

**Требования:**
- ✅ Frontend должен быть запущен на localhost:3001
- ✅ Backend API доступен на localhost:3000
- ✅ Chromium установлен через `npx playwright install chromium --with-deps`

---

## 🐳 Docker / CI Setup

### GitHub Actions пример:

```yaml
name: E2E Tests

on: [push, pull_request]

jobs:
  e2e:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '20'

      - name: Install dependencies
        working-directory: frontend/svetu
        run: yarn install --frozen-lockfile

      - name: Install Playwright browsers
        working-directory: frontend/svetu
        run: npx playwright install chromium --with-deps

      - name: Start backend
        run: |
          cd backend
          go run ./cmd/api/main.go &
          sleep 10

      - name: Start frontend
        working-directory: frontend/svetu
        run: |
          yarn build
          yarn start &
          sleep 15

      - name: Run E2E tests
        working-directory: frontend/svetu
        run: CI=true npx playwright test

      - name: Upload test results
        if: always()
        uses: actions/upload-artifact@v4
        with:
          name: playwright-report
          path: frontend/svetu/playwright-report/
```

### Dockerfile пример:

```dockerfile
FROM mcr.microsoft.com/playwright:v1.54.1-focal

WORKDIR /app

# Copy package files
COPY frontend/svetu/package.json frontend/svetu/yarn.lock ./

# Install dependencies
RUN yarn install --frozen-lockfile

# Copy application
COPY frontend/svetu ./

# Run tests
CMD ["npx", "playwright", "test"]
```

---

## 📝 Отчеты

Playwright генерирует несколько типов отчетов:

```bash
# HTML отчет (после запуска тестов)
npx playwright show-report

# JSON результаты
cat test-results/results.json | jq '.'

# Screenshots при ошибках
ls -la test-results/**/test-failed-*.png

# Videos при ошибках
ls -la test-results/**/*.webm
```

**Расположение:**
- HTML: `playwright-report/index.html`
- JSON: `test-results/results.json`
- Screenshots: `test-results/*/test-failed-*.png`
- Videos: `test-results/*/*.webm`

---

## 🔍 Troubleshooting

### 1. "Browser not found"

```bash
# Установить Chromium с зависимостями
npx playwright install chromium --with-deps
```

### 2. "Permission denied" в Docker

```bash
# Добавить в launchOptions
args: ['--no-sandbox', '--disable-setuid-sandbox']
```

### 3. Frontend не доступен

```bash
# Проверить что frontend запущен
curl http://localhost:3001

# Или изменить baseURL в playwright.config.ts
```

### 4. Тесты падают с timeout

```bash
# Увеличить timeout в playwright.config.ts
use: {
  actionTimeout: 30000,  // 30 секунд
  timeout: 120000,       // 2 минуты на тест
}
```

### 5. "Insufficient shared memory" в Docker

```bash
# Добавить в docker run
docker run --shm-size=1gb ...

# Или добавить в launchOptions
args: ['--disable-dev-shm-usage']
```

---

## ✅ Checklist для CI интеграции

- [x] Playwright установлен (`playwright.config.ts` настроен)
- [x] Headless режим включен по умолчанию
- [x] Browser launch args для CI окружения
- [x] Chromium browser установлен с зависимостями
- [x] @axe-core/playwright установлен
- [x] Тесты совместимы с headless режимом
- [ ] Frontend автоматически запускается перед тестами (для CI)
- [ ] Backend автоматически запускается перед тестами (для CI)
- [ ] Артефакты (screenshots, videos, reports) сохраняются в CI

---

## 🎯 Следующие шаги

1. **Добавить в CI pipeline:**
   - Настроить GitHub Actions / GitLab CI
   - Автоматический запуск backend/frontend перед E2E
   - Сохранение отчетов как артефактов

2. **Расширить покрытие:**
   - Добавить больше E2E сценариев
   - Покрыть критические user journeys
   - Добавить visual regression tests

3. **Оптимизация:**
   - Parallel execution для быстрого прогона
   - Кэширование node_modules и playwright browsers
   - Оптимизация времени выполнения тестов

---

**Документация обновлена:** 2025-10-19 17:00
**Готовность к CI:** ✅ 90% (требуется добавить автозапуск services)
