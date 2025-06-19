# Проблема с запуском yarn dev - Статус отладки

## Проблема
Frontend сервер (yarn dev -p 3001) не запускается - процесс Next.js завершается сразу после старта с сообщением "Done in X.XXs"

## Что обнаружено

### 1. Найдена основная причина
В файле `next.config.ts` был параметр `output: 'standalone'`, который предназначен только для production сборок и вызывает некорректное поведение в dev режиме.

**Статус**: ✅ Уже исправлено (закомментировано)

### 2. Debug информация
При запуске с `DEBUG=*` видно, что процесс завершается с:
```
next:start-server start-server process cleanup
next:start-server start-server process cleanup finished
```

### 3. Проверены внешние факторы
- В `/home/dim/.local/bin/` добавлены скрипты для управления портами, но они НЕ убивают процессы автоматически
- Нет cron задач или фоновых процессов, убивающих порт 3001
- Порт 3001 свободен

## Что уже сделано

1. ✅ Закомментирован `output: 'standalone'` в `next.config.ts`
2. ✅ Очищены кеши (.next, node_modules/.cache)
3. ✅ Попробован запуск без turbopack - та же проблема
4. ✅ Проверены различные способы запуска (yarn, npx, node напрямую)
5. ✅ Убран `--turbopack` из package.json для теста

## Возможные причины оставшейся проблемы

1. Проблема с контекстом выполнения (не поддерживает интерактивные процессы)
2. Конфликт с какими-то системными настройками после добавления скриптов
3. Проблема с конфигурацией Next.js 15.3.2
4. Возможно, нужна полная перезагрузка системы для применения изменений

## План действий после перезагрузки

### 1. Первый тест - обычный запуск
```bash
cd /data/hostel-booking-system/frontend/svetu
yarn dev -p 3001
```

### 2. Если не работает - альтернативные варианты
```bash
# Другой порт
yarn dev -p 3002

# С turbopack (вернуть в package.json)
npx next dev --turbopack -p 3001

# Через screen/tmux
screen -S frontend
yarn dev -p 3001

# С явным NODE_ENV
NODE_ENV=development yarn dev -p 3001
```

### 3. Дополнительная диагностика
```bash
# Проверить процессы
ps aux | grep node

# Проверить порты
lsof -i :3001

# Системные логи
journalctl --user --since "5 minutes ago" | grep next
```

## Файлы для проверки

1. `/data/hostel-booking-system/frontend/svetu/next.config.ts` - убедиться, что `output: 'standalone'` закомментировано
2. `/data/hostel-booking-system/frontend/svetu/package.json` - при необходимости вернуть `--turbopack`

## Версии
- Node.js: v20.19.2
- npm: 10.8.2
- Next.js: 15.3.2
- yarn: 1.22.22