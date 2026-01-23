# OpenSearch Filters & Facets - Quick Start Guide

**Цель:** Быстрое подключение фильтрации и фасетов для фронтенд разработчиков.

---

## 🚀 Quick Start (3 шага)

### 1. Получить фасеты для категории

**gRPC:**
```bash
grpcurl -plaintext -d '{
  "category_id": 123,
  "include_facets": true
}' localhost:50053 search.v1.SearchService/GetSearchFacets
```

**Ответ:**
```json
{
  "categories": [
    {"category_id": 10, "count": 150},
    {"category_id": 20, "count": 200}
  ],
  "price_ranges": [
    {"min": 0, "max": 50, "count": 30},
    {"min": 50, "max": 100, "count": 45}
  ],
  "attributes": {
    "brand": {
      "key": "brand",
      "values": [
        {"value": "apple", "count": 100},
        {"value": "samsung", "count": 150}
      ]
    },
    "color": {
      "key": "color",
      "values": [
        {"value": "black", "count": 80},
        {"value": "white", "count": 70}
      ]
    }
  },
  "source_types": [
    {"key": "c2c", "count": 120},
    {"key": "b2c", "count": 80}
  ],
  "stock_statuses": [
    {"key": "in_stock", "count": 180},
    {"key": "out_of_stock", "count": 20}
  ]
}
```

---

### 2. Поиск с фильтрами

**gRPC:**
```bash
grpcurl -plaintext -d '{
  "query": "smartphone",
  "category_id": 123,
  "filters": {
    "attributes": {
      "brand": {"values": ["apple"]},
      "color": {"values": ["black", "white"]}
    },
    "price": {
      "min": 100.0,
      "max": 500.0
    },
    "source_type": "b2c",
    "stock_status": "in_stock"
  },
  "sort": {
    "field": "price",
    "order": "asc"
  },
  "limit": 20,
  "offset": 0,
  "include_facets": true
}' localhost:50053 search.v1.SearchService/SearchWithFilters
```

**Ответ:**
```json
{
  "listings": [...],
  "total": 42,
  "took_ms": 15,
  "facets": {
    "attributes": {
      "brand": {
        "values": [
          {"value": "apple", "count": 42}
        ]
      },
      "color": {
        "values": [
          {"value": "black", "count": 25},
          {"value": "white", "count": 17}
        ]
      }
    },
    "price_ranges": [...]
  }
}
```

---

### 3. Автодополнение

**gRPC:**
```bash
grpcurl -plaintext -d '{
  "prefix": "ipho",
  "category_id": 123,
  "limit": 10
}' localhost:50053 search.v1.SearchService/GetSuggestions
```

**Ответ:**
```json
{
  "suggestions": [
    {"text": "iphone 15", "score": 0.95},
    {"text": "iphone 14", "score": 0.87},
    {"text": "iphone 13 pro", "score": 0.82}
  ]
}
```

---

## 📋 Поддерживаемые фильтры

| Фильтр | Тип proto | Пример значения | Описание |
|--------|-----------|-----------------|----------|
| **Категория** | `optional int64 category_id` | `123` | ID категории |
| **Текстовый поиск** | `string query` | `"iphone 15"` | Поиск в title + description |
| **Атрибуты** | `map<string, AttributeValues> attributes` | `{"brand": ["apple"], "color": ["black"]}` | Фильтры по атрибутам |
| **Цена** | `PriceRange price` | `{"min": 100, "max": 500}` | Диапазон цен |
| **Тип источника** | `optional string source_type` | `"b2c"` | c2c или b2c |
| **Наличие** | `optional string stock_status` | `"in_stock"` | in_stock, out_of_stock, low_stock |
| **Геолокация** | `LocationFilter location` | `{"lat": 44.78, "lon": 20.44, "radius_km": 10}` | Поиск по радиусу |
| **Сортировка** | `SortConfig sort` | `{"field": "price", "order": "asc"}` | Сортировка результатов |

---

## 🎨 UI Компоненты (рекомендации)

### Фасеты → UI

```typescript
// 1. Бренды (чекбоксы)
facets.attributes["brand"].values.map(v => (
  <Checkbox key={v.value} label={`${v.value} (${v.count})`} />
))

// 2. Цвета (цветовые чипы)
facets.attributes["color"].values.map(v => (
  <ColorChip color={v.value} count={v.count} />
))

// 3. Цена (range slider)
const priceRange = facets.price_ranges
<RangeSlider min={0} max={1000} />

// 4. Категории (древо)
facets.categories.map(c => (
  <CategoryLink id={c.category_id} count={c.count} />
))

// 5. Тип источника (radio buttons)
facets.source_types.map(st => (
  <Radio key={st.key} label={`${st.key} (${st.count})`} />
))
```

### Применение фильтров

```typescript
const [filters, setFilters] = useState({
  attributes: {},
  price: null,
  source_type: null,
  stock_status: null,
})

// Добавить фильтр
const addFilter = (code: string, value: string) => {
  setFilters(prev => ({
    ...prev,
    attributes: {
      ...prev.attributes,
      [code]: [...(prev.attributes[code] || []), value]
    }
  }))
}

// Удалить фильтр
const removeFilter = (code: string, value: string) => {
  setFilters(prev => ({
    ...prev,
    attributes: {
      ...prev.attributes,
      [code]: prev.attributes[code].filter(v => v !== value)
    }
  }))
}

// Очистить все фильтры
const clearFilters = () => {
  setFilters({attributes: {}, price: null, source_type: null, stock_status: null})
}
```

---

## 📊 Сортировка

| Поле | Значение `sort.field` | Описание |
|------|----------------------|----------|
| **Релевантность** | `relevance` | По _score (для текстового поиска) |
| **Цена** | `price` | По цене |
| **Дата создания** | `created_at` | По дате (новые/старые) |
| **Популярность** | `views_count` | По количеству просмотров |
| **Избранное** | `favorites_count` | По количеству добавлений в избранное |

**Порядок:**
- `"asc"` - по возрастанию
- `"desc"` - по убыванию (по умолчанию)

---

## 🔍 Примеры использования

### 1. Фильтр по бренду (один)

```json
{
  "category_id": 123,
  "filters": {
    "attributes": {
      "brand": {"values": ["apple"]}
    }
  }
}
```

### 2. Фильтр по цвету (несколько - OR)

```json
{
  "category_id": 123,
  "filters": {
    "attributes": {
      "color": {"values": ["black", "white", "silver"]}
    }
  }
}
```

### 3. Комбинированные фильтры (AND между атрибутами)

```json
{
  "category_id": 123,
  "filters": {
    "attributes": {
      "brand": {"values": ["apple"]},
      "color": {"values": ["black"]},
      "ram": {"values": ["8GB", "16GB"]}
    },
    "price": {
      "min": 500.0,
      "max": 1500.0
    },
    "stock_status": "in_stock"
  }
}
```

### 4. Геопоиск + фильтры

```json
{
  "category_id": 123,
  "filters": {
    "location": {
      "lat": 44.7866,
      "lon": 20.4489,
      "radius_km": 10.0
    },
    "source_type": "c2c"
  }
}
```

### 5. Поиск с сортировкой

```json
{
  "query": "laptop",
  "category_id": 123,
  "sort": {
    "field": "price",
    "order": "asc"
  }
}
```

---

## 🧪 Тестирование через curl

### Создать тестовый документ

```bash
curl -X POST "localhost:9200/listings_microservice/_doc/test1" \
  -H 'Content-Type: application/json' -d'
{
  "id": 1,
  "uuid": "test-uuid-1",
  "title": "iPhone 15 Pro",
  "description": "Latest iPhone with Pro features",
  "price": 1299.99,
  "currency": "EUR",
  "category_id": 123,
  "status": "active",
  "source_type": "b2c",
  "stock_status": "in_stock",
  "quantity": 10,
  "attributes": [
    {"code": "brand", "type": "select", "value_select": "apple"},
    {"code": "color", "type": "select", "value_select": "black"},
    {"code": "ram", "type": "select", "value_select": "8GB"},
    {"code": "storage", "type": "select", "value_select": "256GB"}
  ],
  "images": [
    {"id": 1, "url": "https://example.com/img1.jpg", "is_primary": true, "display_order": 1}
  ],
  "created_at": "2025-12-17T10:00:00Z"
}'
```

### Поиск по атрибутам

```bash
curl -X POST "localhost:9200/listings_microservice/_search" \
  -H 'Content-Type: application/json' -d'
{
  "query": {
    "bool": {
      "must": [
        {"term": {"category_id": 123}},
        {
          "nested": {
            "path": "attributes",
            "query": {
              "bool": {
                "must": [
                  {"term": {"attributes.code": "brand"}},
                  {"term": {"attributes.value_select": "apple"}}
                ]
              }
            }
          }
        }
      ]
    }
  }
}'
```

### Фасеты (агрегации)

```bash
curl -X POST "localhost:9200/listings_microservice/_search" \
  -H 'Content-Type: application/json' -d'
{
  "size": 0,
  "query": {"term": {"category_id": 123}},
  "aggs": {
    "facet_brand": {
      "nested": {"path": "attributes"},
      "aggs": {
        "filter_by_code": {
          "filter": {"term": {"attributes.code": "brand"}},
          "aggs": {
            "values": {
              "terms": {"field": "attributes.value_select", "size": 100}
            }
          }
        }
      }
    }
  }
}'
```

---

## ⚠️ Важные замечания

1. **Nested queries обязательны** для фильтрации по атрибутам
2. **Type mapping автоматический:** `select` → `value_select`, `number` → `value_number`
3. **OR внутри атрибута:** `{"color": ["black", "white"]}` → OR
4. **AND между атрибутами:** `{"brand": ["apple"], "color": ["black"]}` → AND
5. **Индекс:** `listings_microservice` (не `listings`)
6. **gRPC порт:** 50053
7. **Кэширование:** Включено по умолчанию через Redis

---

## 📚 Дополнительная документация

- **Полный отчет:** `PROGRESS_PHASE2_OPENSEARCH.md`
- **Proto definitions:** `/api/proto/search/v1/*.proto`
- **Query builder:** `/internal/opensearch/query_builder.go`
- **Facets:** `/internal/opensearch/facets.go`
- **Service:** `/internal/service/search/service.go`
- **Tests:** `/internal/opensearch/*_test.go`

---

**Готово! 🚀 Теперь можно интегрировать с фронтендом.**
