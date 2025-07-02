# Паспорт таблицы `translations`

## Назначение
Система мультиязычных переводов для всех сущностей проекта. Позволяет хранить переводы любых полей любых таблиц на разные языки.

## Полная структура таблицы

```sql
CREATE TABLE translations (
    id SERIAL PRIMARY KEY,
    entity_type VARCHAR(50) NOT NULL,
    entity_id INTEGER NOT NULL,
    language VARCHAR(10) NOT NULL,
    field_name VARCHAR(50) NOT NULL,
    translated_text TEXT NOT NULL,
    is_machine_translated BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    metadata JSONB DEFAULT '{}',
    UNIQUE(entity_type, entity_id, language, field_name)
);
```

## Описание полей

| Поле | Тип | Обязательность | Описание |
|------|-----|---------------|----------|
| `id` | SERIAL | PRIMARY KEY | Уникальный идентификатор перевода |
| `entity_type` | VARCHAR(50) | NOT NULL | Тип сущности (например, 'category', 'listing') |
| `entity_id` | INTEGER | NOT NULL | ID сущности в исходной таблице |
| `language` | VARCHAR(10) | NOT NULL | Код языка (ru, en, sr и т.д.) |
| `field_name` | VARCHAR(50) | NOT NULL | Название поля для перевода |
| `translated_text` | TEXT | NOT NULL | Переведенный текст |
| `is_machine_translated` | BOOLEAN | DEFAULT true | Флаг машинного перевода |
| `is_verified` | BOOLEAN | DEFAULT false | Флаг верификации перевода человеком |
| `created_at` | TIMESTAMP | DEFAULT NOW | Время создания записи |
| `updated_at` | TIMESTAMP | AUTO UPDATE | Время последнего обновления (триггер) |
| `metadata` | JSONB | DEFAULT '{}' | Дополнительные метаданные в JSON |

## Индексы

```sql
-- Основной поиск переводов
CREATE INDEX idx_translations_lookup ON translations(entity_type, entity_id, language);

-- Поиск по метаданным
CREATE INDEX idx_translations_metadata ON translations USING GIN (metadata);
```

## Ограничения

- **UNIQUE KEY**: (`entity_type`, `entity_id`, `language`, `field_name`) - один перевод одного поля на один язык для одной сущности

## Триггеры

```sql
-- Автоматическое обновление updated_at при изменении
CREATE TRIGGER update_translations_timestamp
    BEFORE UPDATE ON translations
    FOR EACH ROW
    EXECUTE FUNCTION update_translations_updated_at();
```

## Связи с другими таблицами

Эта таблица имеет **динамические связи** - не классические FK, а логические:
- `entity_type` + `entity_id` ссылается на любую таблицу системы
- Например: `entity_type='category'` + `entity_id=5` → запись в `marketplace_categories` с `id=5`

## Бизнес-правила

1. **Уникальность**: Только один перевод для каждой комбинации (сущность + поле + язык)
2. **Типы сущностей**: Поддерживаются все модели системы
3. **Языки**: Поддерживаются ru, en, sr и другие ISO коды
4. **Машинный перевод**: По умолчанию считается машинным до верификации
5. **Верификация**: Требует ручного подтверждения качества

## Примеры использования

### Добавление перевода категории
```sql
INSERT INTO translations (entity_type, entity_id, language, field_name, translated_text, is_machine_translated)
VALUES ('category', 1, 'en', 'name', 'Electronics', false);
```

### Получение всех переводов объявления
```sql
SELECT field_name, language, translated_text 
FROM translations 
WHERE entity_type = 'listing' AND entity_id = 123;
```

### Поиск непроверенных переводов
```sql
SELECT * FROM translations 
WHERE is_machine_translated = true AND is_verified = false;
```

## Известные особенности

1. **Гибкая схема**: Таблица поддерживает перевод любых полей любых сущностей
2. **Отсутствие FK**: Связи логические, не физические (для гибкости)
3. **Метаданные**: JSONB поле для хранения дополнительной информации о переводе
4. **Машинный перевод**: Система различает автоматические и ручные переводы

## Использование в коде

**Backend**: 
- Модуль в `internal/proj/translations/`
- Handler для управления переводами
- Middleware для автоматической подстановки переводов

**Frontend**:
- Файлы переводов в `src/messages/{en,ru}.json`
- Использование next-intl для интернационализации
- Плейсхолдеры типа "notifications.getError"

## Связанные компоненты

- **Middleware**: Автоматическая подстановка переводов в API ответы
- **Admin панель**: Управление переводами
- **Import/Export**: Возможность массового импорта переводов
- **Translation API**: REST API для работы с переводами