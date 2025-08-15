# Playwright Integration Tests

Интеграционные тесты для полной проверки работы приложения Svetu.

## Что тестируется

Этот тест проверяет работу всего стека приложения:
- PostgreSQL база данных
- Redis для кеширования  
- OpenSearch для поиска
- MinIO для хранения файлов
- Go Backend API
- Next.js Frontend
- Загрузка главной страницы

## Локальный запуск

### Предварительные требования
1. Запустите все сервисы через docker-compose:
```bash
cd ..
docker-compose up -d
```

2. Запустите backend:
```bash
cd ../backend
go run ./cmd/api/main.go
```

3. Запустите frontend:
```bash
cd ../frontend/svetu
yarn dev -p 3001
```

### Запуск тестов
```bash
yarn install
yarn install:browsers
yarn test
```

## CI/CD в GitHub Actions

Тесты автоматически запускаются в GitHub Actions при push в ветки main, develop или tests.

GitHub Actions автоматически:
1. Поднимает все необходимые сервисы (PostgreSQL, Redis, OpenSearch, MinIO)
2. Настраивает базы данных и бакеты
3. Запускает миграции БД
4. Собирает и запускает backend
5. Собирает и запускает frontend  
6. Запускает Playwright тесты
7. Сохраняет отчеты и логи при ошибках

## Структура

- `tests/example.spec.ts` - интеграционный тест загрузки главной страницы
- `playwright.config.ts` - конфигурация Playwright
- `package.json` - зависимости и скрипты

## Текущие тесты

1. **example.spec.ts** - интеграционный тест приложения Svetu
   - Навигация на главную страницу (http://localhost:3001)
   - Проверка загрузки страницы
   - Проверка наличия основных элементов (header, main content)
   - Проверка что карта или контент маркетплейса загрузились