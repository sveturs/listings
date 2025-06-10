# Анализ поиска, фильтрации и сортировки по атрибутам

## Обзор реализации

### 1. Frontend компоненты

#### Созданные компоненты:
- **MarketplaceFilters.tsx** - компонент для отображения фильтров
  - Динамическая загрузка атрибутов категории
  - Поддержка разных типов атрибутов (select, number, boolean, multiselect)
  - Мобильная адаптивность
  - Передача фильтров через колбэк

- **Обновленный MarketplaceList.tsx**
  - Интеграция с фильтрами
  - Поиск по тексту
  - Сортировка (дата, цена)
  - Синхронизация с URL параметрами
  - Динамическая подгрузка (infinite scroll)

#### Формирование запросов:
```typescript
// Добавление атрибутов к параметрам поиска
if (filters.attributes) {
  Object.entries(filters.attributes).forEach(([key, value]) => {
    (params as any)[`attr_${key}`] = value;
  });
}
```

### 2. Backend обработка

#### API endpoint `/api/v1/marketplace/search`:
```go
// Обработка фильтров атрибутов в handler/search.go
attributeFilters := make(map[string]string)
c.Context().QueryArgs().VisitAll(func(key, value []byte) {
    keyStr := string(key)
    if len(keyStr) > 5 && keyStr[:5] == "attr_" {
        attrName := keyStr[5:]
        attributeFilters[attrName] = string(value)
    }
})
```

#### Передача фильтров в сервисный слой:
```go
params.AttributeFilters = attributeFilters
```

### 3. OpenSearch интеграция

#### Построение запроса (buildSearchQuery):
```go
if params.AttributeFilters != nil && len(params.AttributeFilters) > 0 {
    for attrName, attrValue := range params.AttributeFilters {
        // Обработка разных типов атрибутов
        
        // Числовые диапазоны (min,max)
        if strings.Contains(attrValue, ",") {
            parts := strings.Split(attrValue, ",")
            // Создание range запроса
        }
        
        // Nested запросы для атрибутов
        nestedQuery := map[string]interface{}{
            "nested": map[string]interface{}{
                "path": "attributes",
                "query": map[string]interface{}{
                    "bool": map[string]interface{}{
                        "must": []map[string]interface{}{
                            {
                                "term": map[string]interface{}{
                                    "attributes.attribute_name": attrName,
                                },
                            },
                            // Дополнительные условия по значению
                        },
                    },
                },
            },
        }
    }
}
```

#### Маппинг атрибутов в OpenSearch:
```json
"attributes": {
    "type": "nested",
    "properties": {
        "attribute_id": {"type": "integer"},
        "attribute_name": {"type": "keyword"},
        "text_value": {
            "type": "text",
            "fields": {
                "keyword": {"type": "keyword"}
            }
        },
        "numeric_value": {"type": "double"},
        "boolean_value": {"type": "boolean"}
    }
}
```

### 4. PostgreSQL fallback

Создан пример SQL запроса для фильтрации в `search_with_attributes_example.sql`:
- CTE для фильтрации основных полей
- CTE для фильтрации по атрибутам
- INTERSECT для объединения результатов
- Индексы для оптимизации

### 5. Индексация атрибутов

В методе `listingToDoc`:
- Атрибуты индексируются как nested документы
- Важные атрибуты (make, model, brand) дублируются в корень документа
- Создаются текстовые представления для числовых и булевых значений
- Формируется поле `all_attributes_text` для полнотекстового поиска

## Производительность

### Оптимизации:
1. **Индексы OpenSearch**:
   - Nested mapping для эффективной фильтрации
   - Keyword поля для точного сопоставления
   - Text поля с анализаторами для поиска

2. **PostgreSQL индексы**:
   ```sql
   CREATE INDEX idx_listing_attribute_values_listing_id ON listing_attribute_values(listing_id);
   CREATE INDEX idx_listing_attribute_values_attribute_id ON listing_attribute_values(attribute_id);
   CREATE INDEX idx_lav_attr_id_text_value ON listing_attribute_values(attribute_id, text_value);
   ```

3. **Кэширование**:
   - Атрибуты категории кэшируются на frontend
   - OpenSearch кэширует частые запросы

## UX решения

1. **Адаптивный дизайн**:
   - Мобильная версия с выездной панелью фильтров
   - Desktop версия с постоянно видимыми фильтрами

2. **Динамические фильтры**:
   - Фильтры загружаются в зависимости от выбранной категории
   - Только filterable атрибуты отображаются

3. **Типы фильтров**:
   - Select - выпадающий список
   - Number - поля "от" и "до" 
   - Boolean - чекбокс
   - Multiselect - кнопки множественного выбора

4. **Обратная связь**:
   - Индикация количества активных фильтров
   - Кнопка "Очистить все"
   - Синхронизация с URL для возможности поделиться

## Проблемы и рекомендации

### Обнаруженные проблемы:
1. **Отсутствие агрегаций** - нет подсчета количества товаров для каждого значения фильтра
2. **Нет сохранения фильтров** - при перезагрузке теряются настройки (кроме URL параметров)
3. **Производительность** - при большом количестве атрибутов может быть медленно

### Рекомендации:
1. Добавить faceted search с подсчетом товаров
2. Реализовать сохранение последних поисков
3. Добавить предложения по поиску (suggestions)
4. Оптимизировать запросы с большим количеством фильтров
5. Добавить визуализацию диапазонов для числовых атрибутов (слайдеры)

## Код для тестирования

```bash
# Поиск с фильтрами через curl
curl "http://localhost:3000/api/v1/marketplace/search?query=BMW&attr_year=2020,2024&attr_fuel_type=diesel&sort_by=price_asc"

# Переиндексация для обновления атрибутов
cd backend && ./reindex
```