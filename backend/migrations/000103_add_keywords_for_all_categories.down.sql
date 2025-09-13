-- Откат миграции: удаление добавленных ключевых слов

DELETE FROM category_keywords WHERE source = 'manual';