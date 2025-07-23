# CLAUDE.md
- НИКОГДА НЕ ДОБАВЛЯЙ Claude в авторы или соавторы. к примеру как ты сейчас захотел сделать: "Generated with [Claude Code](https://claude.ai/code)"

## ⚠️ КРИТИЧЕСКИ ВАЖНОЕ ПРАВИЛО: Изменения в базе данных

**ВСЕ изменения структуры и данных в базе данных должны производиться ТОЛЬКО через миграции!**

- ❌ **ЗАПРЕЩЕНО** выполнять прямые SQL запросы для изменения схемы или данных
- ✅ **ОБЯЗАТЕЛЬНО** создавать миграции в директории `backend/migrations/`
- ✅ **ОБЯЗАТЕЛЬНО** создавать как up, так и down миграции
- ✅ **ОБЯЗАТЕЛЬНО** тестировать откат миграций

Пример: Если нужно изменить данные в таблице, создай миграцию:
```bash
# НЕ ДЕЛАЙ ТАК:
psql -c "UPDATE table SET ..."

# ДЕЛАЙ ТАК:
# 1. Создай файлы миграции
backend/migrations/000XXX_description.up.sql
backend/migrations/000XXX_description.down.sql

# 2. Примени миграцию через мигратор
```

## Работа с API документацией

Для экономии контекста используй swagger.json через JSON MCP:

1. Запусти HTTP сервер: `cd /data/hostel-booking-system/backend/docs && python3 -m http.server 8888`
2. Используй JSON MCP для поиска:
  - Эндпоинты: `$.paths["/api/v1/..."]`
  - Модели: `$.definitions["..."]`
  - Параметры: `$.paths["/api/v1/..."]["post"]["parameters"]`

ВСЕГДА сначала ищи информацию в swagger.json, и только потом анализируй код!
2. Примеры полезных запросов:
# Найти конкретный эндпоинт
Используй JSON MCP для получения информации об эндпоинте /api/v1/auth/login из http://localhost:8888/swagger.json. 

# Получить все эндпоинты модуля
Найди все marketplace эндпоинты в swagger через JSON MCP. 

# Получить схему модели
Извлеки схему модели MarketplaceListing из swagger.

## ⚠️ КРИТИЧЕСКИ ВАЖНОЕ ПРАВИЛО: Управление процессами и screen сессиями

**ВСЕГДА закрывайте старые процессы и screen сессии перед запуском новых!**

PostgreSQL имеет лимит в 100 подключений, а каждый экземпляр backend создаёт пул до 50 подключений. 
Множественные экземпляры быстро исчерпают лимит и вызовут ошибку "too many clients already".

### Правильная последовательность запуска backend:
```bash
# 1. СНАЧАЛА останови все процессы на порту 3000
/home/dim/.local/bin/kill-port-3000.sh

# 2. Закрой ВСЕ старые screen сессии backend
screen -ls | grep backend-3000 | awk '{print $1}' | xargs -I {} screen -S {} -X quit
screen -wipe  # очистить мёртвые сессии

# 3. ТОЛЬКО ПОТОМ запускай новый экземпляр
screen -dmS backend-3000 bash -c 'cd /data/hostel-booking-system/backend && go run ./cmd/api/main.go 2>&1 | tee /tmp/backend.log'
```

### Мониторинг подключений к БД:
```bash
# Проверить количество текущих подключений
psql "postgres://postgres:password@localhost:5432/svetubd?sslmode=disable" -c "SELECT COUNT(*) FROM pg_stat_activity;"

# Если подключений больше 90 - это критично!
# Перезапусти PostgreSQL для сброса всех подключений:
sudo systemctl restart postgresql
```

## Проект: Sve Tu Platform - Marketplace

Состоит из
- frontend: React, NextJS, Tailwind, DaisyUI
- backend: Go, Postgres, Minio, OpenSearch

## Документация по категориям и фильтрам

- **IMPLEMENTATION_CATEGORY_SELECTOR.md** - общая инструкция по реализации выбора категорий во всей системе
- **DISPLAY_CATEGORIES_INSTRUCTION.md** - инструкция по правильному отображению категорий
- **CATEGORY_ATTRIBUTES_STATUS.md** - статус системы атрибутов категорий и их интеграции с фильтрами
- **MAP_CATEGORIES_FILTERS_INSTRUCTION.md** - инструкция по работе с категориями и фильтрами на карте

## Документация по витринам (Storefronts)

- **STOREFRONTS_STATUS.md** - актуальный статус функционала витрин и последние обновления
- **STOREFRONT_PRODUCT_CATEGORY_SELECTION.md** - реализация выбора категорий при создании товара в витрине
- **PRODUCT_LOCATION_IMPLEMENTATION_PLAN.md** - управление местоположением и приватностью адресов товаров
- **LOCATION_PICKER_INSTRUCTION.md** - использование компонента LocationPicker для выбора адреса

## IMPORTANT INSTRUCTIONS

- When working with frontend - use instructions from [@.ai/frontend.md](.ai/frontend.md)
- When working with backend - use instructions from [@.ai/backend.md](.ai/backend.md)
- When working with migrations - use instructions from [@.ai/migrations.md](.ai/migrations.md)
- When working with GIT - use instructions from [@.ai/git.md](.ai/git.md)
- When working with github - use instructions from [@.ai/github.md](.ai/github.md)
- Follow instruction by [@.ai/openapi.md](.ai/openapi.md:3)
- For manager frontend and backend running use instructions from [@.ai/server-management.md]

## IMPORTANT WORKFLOW RULES

- **Язык общения**: Use Russian in your answers

- **Управление задачами**: ВСЕГДА Информацию о прогрессе текущей задачи храни в файле `/memory-bank/currentTask.md`. При начале новой сессии получай статус из этого файла. Прогресс сохраняй в `/memory-bank/tasks/<next_number_as_\d\d\d>-<task_name>.md`. После завершения задачи очищай `/memory-bank/current-task.md`

- **Коммиты**: Подбирай описание для коммита на основании внесенных изменений. НЕ ИСПОЛЬЗУЙ УКАЗАНИЕ АВТОРСТВА CLAUDE

- **Разработка**:
   - Запускай сервер командой `yarn dev -p 3001`. Если надо его остановить, ищи процесс для остановке ИМЕННО ПО ПОРТУ 3001. НЕ ИСПОЛЬЗУЙ `killall node 2>/dev/null || true`
   - Используй mcp playwright для работы с браузером Google Chrome
   - ИСПОЛЬЗУЙ ДЕФОЛТНЫЕ КЛАССЫ DAISYUI. Используй mcp context7 для поиска по документации daisyui при работе с версткой
   - **Галерея изображений**: Для отображения фотографий ВСЕГДА используй компонент `ImageGallery` из `src/components/reviews/ImageGallery.tsx`. Этот компонент обеспечивает единообразный просмотр изображений с полупрозрачным фоном, навигацией и миниатюрами

- **Качество кода**: Перед завершением задачи ОБЯЗАТЕЛЬНО выполни `yarn format && yarn lint && yarn build`. Задача считается выполненной только если все команды прошли успешно

- **Зависимости**: При добавлении ключевых зависимостей указывай их в разделе "Key Dependencies"

- **Переводы**: Backend возвращает только placeholder'ы (например, "notifications.getError"), а фронт переводит их через файлы `frontend/svetu/src/messages/{en,ru}.json`

- **Переменные окружения**: При добавлении новых переменных:
   1. Обновляй `.env` и `.env.example` (с примерами и комментариями)
   2. Обновляй интерфейсы в `src/config/types.ts` (`EnvVariables` и `Config`)
   3. Добавляй обработку в `src/config/index.ts` в методе `loadConfig()`

- **Проверка качества кода и форматирование**

    ### Frontend (в папке `frontend/svetu/`)
    Перед завершением любой задачи на Frontend **ОБЯЗАТЕЛЬНО** выполните:

    ```bash
    yarn format && yarn lint && yarn build
    ```
    
    ### Backend (в папке `backend/`)
    Перед завершением любой задачи на Backend **ОБЯЗАТЕЛЬНО** выполните:

    ```bash
    make format && make lint
    ```
    
    Задача считается выполненной только если все команды выполнились успешно!

- **Pre-commit hooks и автоформатирование**

    В проекте настроена система автоматического форматирования кода:
    
    ### Установка pre-commit (один раз на машине)
    ```bash
    # macOS: brew install pre-commit
    # Ubuntu: sudo apt install pre-commit
    # или: pip install pre-commit
    
    # Активация в проекте
    cd /path/to/project
    pre-commit install
    ```
    
    ### Команды форматирования
    
    **Backend (Go):**
    - `make format` - автоформатирование с gofumpt + goimports
    - `make format-check` - проверка форматирования без изменений
    - `make lint` - запуск golangci-lint
    - `make pre-commit` - все проверки перед коммитом
    
    **Frontend:**
    - `yarn format` - форматирование с Prettier
    - `yarn format:check` - проверка форматирования
    - `yarn lint` - ESLint проверки
    - `yarn lint:fix` - исправление ESLint ошибок
    
    ### Важно
    - Pre-commit hooks автоматически запускаются при каждом коммите
    - GitHub Actions проверяют форматирование в PR
    - EditorConfig (.editorconfig) обеспечивает единые настройки IDE
    - Больше НЕ БУДЕТ проблем с trailing spaces и line endings!
