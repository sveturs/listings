# Проблемы с миграциями и их решение

## Описание проблемы

При попытке создания продукта для витрины возникла серия ошибок, связанных с отсутствием автоинкремента для полей `id` в нескольких таблицах БД.

### Последовательность ошибок:

1. **Первая ошибка**: `ERROR: null value in column "id" of relation "storefront_products" violates not-null constraint`
   - Причина: Отсутствие автоинкремента для поля `id` в таблице `storefront_products`
   - Решение: Создана миграция 000121

2. **Вторая ошибка**: `ERROR: function calculate_blurred_location(double precision, double precision, location_privacy_level) does not exist`
   - Причина: Триггер `auto_geocode_storefront_product` вызывал несуществующую функцию
   - Решение: Обнаружено, что функция существует, но с другой сигнатурой (numeric вместо double precision)
   - Создана миграция 000124 для исправления триггера

3. **Третья ошибка**: `ERROR: null value in column "id" of relation "unified_geo" violates not-null constraint`
   - Причина: Отсутствие автоинкремента для поля `id` в таблице `unified_geo`
   - Решение: Создана миграция 000125

## Проблемы с процессом миграции

### 1. Проблемы с мигратором
- Мигратор периодически впадал в "dirty" состояние
- Команда `./migrator force` не поддерживается используемым мигратором
- Конфликты с уже существующими типами данных при повторном применении миграций

### 2. Непоследовательность применения
- Миграции применялись напрямую через `psql` вместо мигратора из-за его ошибок
- Это привело к рассинхронизации состояния БД и таблицы `schema_migrations`

## Правильное решение

### 1. Единый подход к миграциям

```bash
# Всегда использовать мигратор для применения миграций
cd /data/hostel-booking-system/backend
./migrator up

# Если возникла ошибка dirty state:
psql $DATABASE_URL -c "UPDATE schema_migrations SET dirty = false WHERE version = (SELECT MAX(version) FROM schema_migrations);"

# Если нужно откатить конкретную версию:
./migrator down [version]
```

### 2. Структура миграций для автоинкремента

Для каждой таблицы с проблемой автоинкремента создавать миграцию по шаблону:

**UP миграция:**
```sql
-- Создаем последовательность
CREATE SEQUENCE IF NOT EXISTS table_name_id_seq;

-- Устанавливаем владельца последовательности
ALTER SEQUENCE table_name_id_seq OWNED BY table_name.id;

-- Устанавливаем значение по умолчанию
ALTER TABLE table_name 
    ALTER COLUMN id SET DEFAULT nextval('table_name_id_seq');

-- Устанавливаем текущее значение последовательности
SELECT setval('table_name_id_seq', COALESCE((SELECT MAX(id) FROM table_name), 0) + 1, false);
```

**DOWN миграция:**
```sql
-- Удаляем значение по умолчанию
ALTER TABLE table_name 
    ALTER COLUMN id DROP DEFAULT;

-- Удаляем последовательность
DROP SEQUENCE IF EXISTS table_name_id_seq;
```

### 3. Проверка перед применением

Перед созданием миграции всегда проверять:

```bash
# Структуру таблицы
psql $DATABASE_URL -c "\d table_name"

# Наличие последовательностей
psql $DATABASE_URL -c "\ds"

# Текущее состояние миграций
psql $DATABASE_URL -c "SELECT * FROM schema_migrations ORDER BY version DESC LIMIT 10;"
```

### 4. Обработка ошибок мигратора

Если мигратор не может применить миграцию из-за конфликтов:

1. Проверить, не была ли часть миграции уже применена
2. При необходимости создать условные конструкции в миграции:
   ```sql
   -- Пример условного создания
   DO $$ 
   BEGIN
       IF NOT EXISTS (SELECT 1 FROM pg_sequences WHERE sequencename = 'table_name_id_seq') THEN
           CREATE SEQUENCE table_name_id_seq;
       END IF;
   END $$;
   ```

## Рекомендации

1. **Всегда использовать мигратор** для применения миграций
2. **Тестировать миграции** на копии БД перед применением в продакшене
3. **Документировать** все прямые SQL команды, если они были применены в обход мигратора
4. **Проверять состояние** БД после каждой миграции
5. **Создавать резервные копии** БД перед массовыми изменениями

## Итоговое состояние

После всех исправлений:
- ✅ Таблица `storefront_products` имеет автоинкремент для `id`
- ✅ Таблица `unified_geo` имеет автоинкремент для `id`
- ✅ Функция `auto_geocode_storefront_product` использует правильную сигнатуру `calculate_blurred_location`
- ✅ Создание продуктов витрины работает корректно