# Настройка Augment AI для GoLand в проекте Hostel Booking System

## Текущее состояние конфигурации

### ✅ Что уже настроено:

1. **augment.config.json** - файл присутствует в корне проекта с базовой конфигурацией:
   - Указаны языки: TypeScript, Go, JavaScript
   - Фреймворки: Next.js, React, Fiber
   - Настройки контекстного окна (8000 токенов)
   - Включена индексация с интервалом обновления 300 секунд
   - Настроены паттерны исключения для больших/временных файлов

2. **.augmentignore** - файл настроен корректно:
   - Игнорируются node_modules, vendor, логи
   - Исключены большие директории (data/, uploads/)
   - Пропускаются сгенерированные файлы
   - Оптимизирована индексация для больших директорий backend

### ❌ Что отсутствует:

1. Директория **.idea/** не найдена в проекте
2. Нет специфичных настроек для GoLand в augment.config.json

## Рекомендации по настройке Augment в GoLand

### 1. Установка плагина Augment

1. Откройте GoLand
2. Перейдите в **File → Settings** (или **GoLand → Settings** на macOS)
3. Выберите **Plugins** в боковой панели
4. Найдите "Augment" в маркетплейсе
5. Нажмите **Install**
6. Перезапустите GoLand

### 2. Первоначальная настройка

После установки:
1. Используйте сочетание клавиш **Cmd/Ctrl + L** или кликните на иконку Augment в боковой панели
2. Войдите в свой аккаунт Augment
3. Дождитесь индексации проекта (может занять несколько минут для большого проекта)

### 3. Оптимизация для Go-разработки

Добавьте в **augment.config.json** специфичные настройки для Go:

```json
{
  "projectName": "Hostel Booking System - Sve Tu Platform",
  "description": "Marketplace platform with React frontend and Go backend",
  "language": ["typescript", "go", "javascript"],
  "framework": ["nextjs", "react", "fiber"],
  "contextWindow": {
    "maxTokens": 8000,
    "includeImports": true,
    "includeComments": true
  },
  "indexing": {
    "enabled": true,
    "refreshInterval": 300,
    "maxFileSize": "500KB",
    "excludePatterns": [
      "**/node_modules/**",
      "**/vendor/**",
      "**/*.log",
      "**/dist/**",
      "**/build/**",
      "**/.next/**",
      "**/uploads/**",
      "**/data/**"
    ]
  },
  "codeCompletion": {
    "enabled": true,
    "contextLines": 50,
    "debounceMs": 300
  },
  "chat": {
    "enabled": true,
    "defaultContext": [
      "CLAUDE.md",
      ".ai/*.md",
      "README.md"
    ]
  },
  "goSpecific": {
    "goModPath": "backend/go.mod",
    "testPatterns": ["**/*_test.go"],
    "buildTags": [],
    "includeVendor": false
  }
}
```

### 4. Создание .augment-guidelines

Создайте файл **.augment-guidelines** в корне проекта для командных соглашений:

```
# Augment Guidelines for Hostel Booking System

## Go Development
- Follow Go idioms and best practices
- Use meaningful variable and function names
- Add comments for exported functions
- Use error handling pattern: if err != nil
- Prefer table-driven tests

## Project Structure
- Backend code in /backend directory
- Frontend code in /frontend/svetu directory
- Database migrations in /backend/migrations
- API handlers in /backend/internal/proj/

## Code Style
- Use gofumpt for Go formatting
- Use Prettier for TypeScript/JavaScript
- Follow existing patterns in the codebase
```

### 5. Настройка MCP серверов (опционально)

Если используете MCP (Model Context Protocol):

1. Откройте настройки Augment в GoLand (иконка шестеренки в панели Augment)
2. Добавьте конфигурацию MCP серверов при необходимости

### 6. Проверка прав доступа

Убедитесь, что GoLand имеет права на чтение всех файлов проекта:

```bash
# Проверка прав доступа (выполните в терминале)
find /data/hostel-booking-system -type f ! -readable 2>&1 | grep -v "Permission denied"
```

### 7. Оптимизация производительности

Для больших проектов рекомендуется:

1. Увеличить память для GoLand:
   - Help → Edit Custom VM Options
   - Установите `-Xmx4096m` или больше

2. Настроить индексацию GoLand:
   - File → Settings → Appearance & Behavior → System Settings
   - Снимите галочку "Synchronize files on frame or editor tab activation"

3. Исключить ненужные директории из индексации GoLand:
   - Правый клик на директории → Mark Directory as → Excluded

## Проблемы и решения

### Проблема: Augment не видит изменения в коде
**Решение**: Проверьте интервал обновления индексации в augment.config.json (сейчас 300 секунд)

### Проблема: Слишком медленная работа
**Решение**: 
- Добавьте больше паттернов в excludePatterns
- Уменьшите maxTokens в contextWindow
- Проверьте .augmentignore

### Проблема: Не работают подсказки для Go
**Решение**: 
- Убедитесь, что go.mod правильно определен
- Проверьте, что GOPATH и GOROOT настроены в GoLand

## Дополнительные ресурсы

- [Документация Augment](https://docs.augmentcode.com/)
- [Плагин Augment для JetBrains](https://plugins.jetbrains.com/plugin/24072-augment)
- [Настройка GoLand](https://www.jetbrains.com/help/go/configuring-project-and-ide-settings.html)