# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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

- **Проверка качества кода**

    Перед завершением любой задачи на Frontend **ОБЯЗАТЕЛЬНО** выполните:

    ```bash
    yarn format && yarn lint && yarn build
    ```
    
    Задача считается выполненной только если все три команды выполнились успешно!
