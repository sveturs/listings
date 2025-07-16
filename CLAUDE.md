# CLAUDE.md
- НИКОГДА НЕ ДОБАВЛЯЙ Claude в авторы или соавторы. к примеру как ты сейчас захотел сделать: "Generated with [Claude Code](https://claude.ai/code)"
- для экономии памяти контекста СТАРАЙСЯ всё выполнять через claude -p --dangerously-skip-permissions - таким образом ты сможешь максимально сильно увеличить колличество сделанного в рамках одной сессии за счет экономии памяти контекста. твои подзадачи будут тратить свой контекст на анализ и прочее и выдавать тебе нужный тебе результат.   
  1. даже запуск и рестарт сервисов запускай через claude -p --dangerously-skip-permissions, - наши логированные логи могут быть огромными и это съест твою память контекста.
  2. Разрешения для claude -p --dangerously-skip-permissions добавлены в настройки.
  3. Если ты упираешься в лимит времени на задачу - запускай claude -p --dangerously-skip-permissions в фоновом режиме с выводом информации в файл или используй screen. 
  4. Когда используешь screen всегда закрывай сессии - не оставляй после выполнения задачи. часто бывает такое, что сессия в скрине продолжает править код даже после того как в основной сессии - в родителе ты говоришь что работа выполнена. и фоновая портит код.

не допускай таких ошибок:

ошибка номер раз! ты вызвал playwright:browser_navigate (MCP)(url: "http://localhost:3001") напрямую из этой сессии. а ты логи там не захочешь смотреть? или скриншот сделать? конечно же захочешь - а вот и контекст потратится. так что не забудь - оборачивай в claude -p  --dangerously-skip-permissions "............" !!!

ошибка номер 2! ты доверил слишком длинный план, длинное задание в маленького исполнителя (в claude -p --dangerously-skip-permissions)! давай ему короткие задания, что бы по пути была возможность тебе корректировать план и вносить нужные правки и исправления в код ( конечно так же через обертку claude -p --dangerously-skip-permissions )
ошибка номер 3 - если задачка помощника вываливается по таймауту - можно их запускать в скрине.
ошибка номер 4 - необходимо следить за работой в скрине - я видел как ты отправил задачу в скрин, не дождавшись её завершения ты сделал красивый комит в котором сказал что всё работает, а задача в скрине сломала нам вест код и это ни кем не проверялось!
ошибка номер 5 - никогда не выключай сессии скриин shared* - мы сами в одной из них работаем!
ошибка номер 6! - помощник запустил другого помощника, другой армощник запустил ещё одного и так далее - все начальники, а исполнителя  нету. поэтому каждому помощнику запрещай запускать вложенных помощников
This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.
## Работа с API документацией

Для экономии контекста используй swagger.json через JSON MCP:

1. Запусти HTTP сервер: `cd backend/docs && python3 -m http.server 8888`
2. Используй JSON MCP для поиска:
  - Эндпоинты: `$.paths["/api/v1/..."]`
  - Модели: `$.definitions["..."]`
  - Параметры: `$.paths["/api/v1/..."]["post"]["parameters"]`

ВСЕГДА сначала ищи информацию в swagger.json, и только потом анализируй код!
2. Примеры полезных запросов:
# Найти конкретный эндпоинт
claude -p --dangerously-skip-permissions "Используй JSON MCP для получения информации об эндпоинте /api/v1/auth/login из http://localhost:8888/swagger.json"

# Получить все эндпоинты модуля
claude -p --dangerously-skip-permissions "Найди все marketplace эндпоинты в swagger через JSON MCP"

# Получить схему модели
claude -p --dangerously-skip-permissions "Извлеки схему модели MarketplaceListing из swagger"

## Проект: Sve Tu Platform - Marketplace

Состоит из
- frontend: React, NextJS, Tailwind, DaisyUI
- backend: Go, Postgres, Minio, OpenSearch

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
