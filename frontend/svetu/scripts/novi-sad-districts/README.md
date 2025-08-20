# Управление районами Нови-Сада

Набор скриптов для управления интерактивной картой районов города Нови-Сад.

## Структура

```
scripts/novi-sad-districts/
├── auto_add_district.js      # Добавление района из OpenStreetMap
├── safe_remove_districts.js  # Безопасное удаление районов
├── fix_syntax.js             # Исправление синтаксических ошибок
└── README.md                 # Эта документация
```

## Использование

### Добавление района

```bash
# Из корня frontend/svetu:
node scripts/novi-sad-districts/auto_add_district.js "Название" "Novi Sad" "Serbia"

# Примеры:
node scripts/novi-sad-districts/auto_add_district.js "Veternik" "Novi Sad" "Serbia"
node scripts/novi-sad-districts/auto_add_district.js "Slana bara" "Novi Sad" "Serbia"
```

### Удаление районов

```bash
# Из корня frontend/svetu:
node scripts/novi-sad-districts/safe_remove_districts.js district-id-1 district-id-2

# Пример:
node scripts/novi-sad-districts/safe_remove_districts.js veternik futog petrovaradin
```

### Исправление синтаксиса

```bash
# Из корня frontend/svetu:
node scripts/novi-sad-districts/fix_syntax.js
```

## Веб-интерфейс

Доступен по адресу: http://localhost:3001/ru/examples/novi-sad-districts/manage

Позволяет:

- Искать районы через OpenStreetMap
- Выбирать и удалять существующие районы
- Генерировать команды для выполнения в терминале

## Особенности

1. **Автоматическое добавление** - скрипт получает точные границы района из OpenStreetMap
2. **Безопасное удаление** - создается резервная копия перед удалением
3. **Проверка синтаксиса** - автоматическая проверка и исправление после изменений
4. **31 район** - текущее количество районов на карте

## Файлы данных

Районы хранятся в файле:

```
src/app/[locale]/examples/novi-sad-districts/page.tsx
```

## Резервные копии

При удалении районов автоматически создаются резервные копии с временной меткой:

```
page.tsx.backup_1234567890
```

Для восстановления:

```bash
cp page.tsx.backup_1234567890 page.tsx
```
