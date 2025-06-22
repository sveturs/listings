# Migration Rules

- Директория для миграций - [@backend/migrations](backend/migrations)
- Каждая миграция должна иметь up и down файл
- Файлы именуем как `<migration-number>_<migration-short-description>.<up/down>.sql`
- Старайся не использовать SQL триггеры в миграциях
