# Quick Start - Fix Attribute Options

## TL;DR

```bash
# 1. Проверить что будет изменено (безопасно)
cd /p/github.com/sveturs/listings
go run ./cmd/fix_attribute_options/main.go

# 2. Применить изменения (ВНИМАНИЕ: изменяет БД!)
go run ./cmd/fix_attribute_options/main.go --dry-run=false
```

## Что исправляет

**До:**
```json
{"value": "manual", "label": "Manual"}
```

**После:**
```json
{"value": "manual", "label": {"en": "Manual", "ru": "Manual", "sr": "Manual"}}
```

## Флаги

- `--dry-run=true` (default) - Только показать, не применять
- `--dry-run=false` - Применить изменения в БД
- `-v` - Подробный вывод

## Безопасность

✅ Dry-run по умолчанию
✅ Транзакции
✅ Валидация JSON
✅ Rollback при ошибках

## Статистика (актуальная)

По состоянию на 2025-11-17:
- **Всего атрибутов:** 65
- **Требуют обновления:** 13
- **Опций конвертируется:** 91
- **Уже правильный формат:** 52

## Примеры атрибутов для обновления

1. `transmission` - 5 опций (Manual, Automatic, Semi-Automatic, CVT, DSG)
2. `fuel_type` - 7 опций (Petrol, Diesel, Electric, Hybrid, etc.)
3. `book_language` - 8 опций (ru, en, sr, de, fr, es, it, other)
4. `book_genre` - 10 опций (fiction, non-fiction, science, etc.)
5. `book_format` - 4 опции (hardcover, paperback, ebook, audiobook)

И ещё 8 атрибутов...

## Подключение

```
postgresql://listings_user:listings_secret@localhost:35434/listings_dev_db
```

## Полная документация

См. [README.md](./README.md)
