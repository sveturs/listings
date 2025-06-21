# План улучшения поисковой системы

## 1. Индексация товаров витрин

### Проблема
Товары витрин (storefront_products) не индексируются в OpenSearch, что делает их недоступными для поиска.

### Решение
1. Создать единый индекс `all_products` для всех товаров:
   - Товары маркетплейса (marketplace_listings)
   - Товары витрин (storefront_products)

2. Добавить поле `product_type` для различения:
   - `marketplace` - товары маркетплейса
   - `storefront` - товары витрин

3. Структура документа:
```json
{
  "id": "product_123",
  "product_type": "storefront",
  "storefront_id": 456,
  "title": "iPhone 15 Pro",
  "description": "...",
  "price": 999.99,
  "category_id": 10,
  "attributes": [...],
  "location": { "lat": 45.0, "lon": 20.0 },
  "search_keywords": ["iphone", "apple", "смартфон"],
  // Специфичные поля для витрин
  "storefront_name": "Tech Store",
  "storefront_rating": 4.8,
  "inventory_count": 5
}
```

## 2. Улучшенный маппинг с анализаторами

### Русский анализатор
```json
{
  "settings": {
    "analysis": {
      "analyzer": {
        "russian_analyzer": {
          "tokenizer": "standard",
          "filter": [
            "lowercase",
            "russian_stop",
            "russian_keywords",
            "russian_stemmer"
          ]
        }
      },
      "filter": {
        "russian_stop": {
          "type": "stop",
          "stopwords": "_russian_"
        },
        "russian_keywords": {
          "type": "keyword_marker",
          "keywords": ["iphone", "samsung", "xiaomi"]
        },
        "russian_stemmer": {
          "type": "stemmer",
          "language": "russian"
        }
      }
    }
  }
}
```

### Edge n-gram для автодополнения
```json
{
  "analysis": {
    "tokenizer": {
      "edge_ngram_tokenizer": {
        "type": "edge_ngram",
        "min_gram": 2,
        "max_gram": 20,
        "token_chars": ["letter", "digit"]
      }
    },
    "analyzer": {
      "autocomplete_analyzer": {
        "tokenizer": "edge_ngram_tokenizer",
        "filter": ["lowercase"]
      }
    }
  }
}
```

## 3. Расширенные поисковые возможности

### Multi-match с boost
```json
{
  "query": {
    "multi_match": {
      "query": "iphone 15",
      "fields": [
        "title^3",
        "title.autocomplete^2",
        "description",
        "attributes.brand^2",
        "attributes.model^2",
        "search_keywords"
      ],
      "type": "best_fields",
      "fuzziness": "AUTO"
    }
  }
}
```

### Функция скоринга
```json
{
  "query": {
    "function_score": {
      "query": { /* основной запрос */ },
      "functions": [
        {
          "filter": { "term": { "is_verified": true }},
          "weight": 1.5
        },
        {
          "filter": { "term": { "has_discount": true }},
          "weight": 1.2
        },
        {
          "gauss": {
            "created_at": {
              "scale": "7d",
              "offset": "1d",
              "decay": 0.5
            }
          }
        }
      ],
      "score_mode": "sum"
    }
  }
}
```

## 4. Агрегации и фасеты

### Динамические фасеты по атрибутам
```json
{
  "aggs": {
    "attributes": {
      "nested": {
        "path": "attributes"
      },
      "aggs": {
        "by_name": {
          "terms": {
            "field": "attributes.name.keyword"
          },
          "aggs": {
            "values": {
              "terms": {
                "field": "attributes.value.keyword",
                "size": 20
              }
            }
          }
        }
      }
    },
    "price_ranges": {
      "range": {
        "field": "price",
        "ranges": [
          { "to": 100 },
          { "from": 100, "to": 500 },
          { "from": 500, "to": 1000 },
          { "from": 1000 }
        ]
      }
    }
  }
}
```

## 5. Поиск по геолокации

### Сортировка по расстоянию
```json
{
  "sort": [
    {
      "_geo_distance": {
        "location": {
          "lat": 45.0,
          "lon": 20.0
        },
        "order": "asc",
        "unit": "km",
        "distance_type": "arc"
      }
    }
  ]
}
```

### Фильтр по радиусу с учетом доставки
```json
{
  "query": {
    "bool": {
      "should": [
        {
          "geo_distance": {
            "distance": "10km",
            "location": { "lat": 45.0, "lon": 20.0 }
          }
        },
        {
          "term": { "has_delivery": true }
        }
      ],
      "minimum_should_match": 1
    }
  }
}
```

## 6. Персонализация результатов

### Boost на основе истории пользователя
```json
{
  "query": {
    "bool": {
      "must": { /* основной запрос */ },
      "should": [
        {
          "terms": {
            "category_id": [/* любимые категории пользователя */],
            "boost": 1.5
          }
        },
        {
          "terms": {
            "storefront_id": [/* витрины, где покупал */],
            "boost": 1.3
          }
        }
      ]
    }
  }
}
```

## 7. Реализация в коде

### Сервис индексации товаров витрин
```go
// internal/proj/storefronts/storage/opensearch/product_indexer.go

type ProductIndexer struct {
    client    *opensearch.Client
    indexName string
}

func (pi *ProductIndexer) IndexProduct(product *models.StorefrontProduct) error {
    doc := map[string]interface{}{
        "id":            fmt.Sprintf("sp_%d", product.ID),
        "product_type":  "storefront",
        "storefront_id": product.StorefrontID,
        "title":         product.Name,
        "description":   product.Description,
        "price":         product.Price,
        "category_id":   product.CategoryID,
        "sku":          product.SKU,
        "barcode":      product.Barcode,
        "inventory": map[string]interface{}{
            "count":        product.InventoryCount,
            "reserved":     product.ReservedCount,
            "available":    product.InventoryCount - product.ReservedCount,
            "track":        product.TrackInventory,
        },
        "attributes":    pi.processAttributes(product.Attributes),
        "images":        pi.processImages(product.Images),
        "created_at":    product.CreatedAt,
        "updated_at":    product.UpdatedAt,
    }
    
    // Добавляем данные витрины
    if product.Storefront != nil {
        doc["storefront"] = map[string]interface{}{
            "name":       product.Storefront.Name,
            "slug":       product.Storefront.Slug,
            "rating":     product.Storefront.Rating,
            "verified":   product.Storefront.IsVerified,
        }
        
        if product.Storefront.Latitude != nil && product.Storefront.Longitude != nil {
            doc["location"] = map[string]interface{}{
                "lat": *product.Storefront.Latitude,
                "lon": *product.Storefront.Longitude,
            }
        }
    }
    
    return pi.client.Index(pi.indexName, doc)
}
```

### Unified Search API
```go
// internal/proj/search/service/unified_search.go

type UnifiedSearchService struct {
    opensearch *opensearch.Client
}

func (s *UnifiedSearchService) Search(params SearchParams) (*SearchResult, error) {
    query := s.buildQuery(params)
    
    // Добавляем персонализацию
    if params.UserID > 0 {
        query = s.addPersonalization(query, params.UserID)
    }
    
    // Выполняем поиск
    response, err := s.opensearch.Search("all_products", query)
    if err != nil {
        return nil, err
    }
    
    return s.parseResponse(response, params)
}
```

## 8. Миграция и переиндексация

### Скрипт миграции
```bash
#!/bin/bash
# reindex-all-products.sh

# 1. Создаем новый индекс с правильным маппингом
curl -X PUT "localhost:9200/all_products" \
  -H 'Content-Type: application/json' \
  -d @mappings/all_products.json

# 2. Переиндексируем товары маркетплейса
./reindex-marketplace-products

# 3. Индексируем товары витрин
./index-storefront-products

# 4. Создаем алиас
curl -X POST "localhost:9200/_aliases" \
  -H 'Content-Type: application/json' \
  -d '{
    "actions": [
      { "add": { "index": "all_products", "alias": "products" }}
    ]
  }'
```

## 9. Мониторинг и оптимизация

### Метрики для отслеживания
- Время ответа поиска (p50, p95, p99)
- Количество результатов по запросу
- Click-through rate (CTR) результатов
- Процент "пустых" поисков
- Популярные поисковые запросы

### Оптимизация производительности
- Использование кэша запросов
- Настройка refresh_interval
- Оптимизация размера шардов
- Использование search templates

## 10. A/B тестирование

### Варианты для тестирования
- Различные веса полей (title boost)
- Fuzzy search vs exact match
- Включение/выключение персонализации
- Различные алгоритмы ранжирования