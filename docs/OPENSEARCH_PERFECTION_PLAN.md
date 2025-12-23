# OpenSearch Perfection Plan

**Дата:** 2025-12-19
**Статус:** ПЛАН УТВЕРЖДЁН (v2.0 - с улучшениями)
**Цель:** Довести поисковик Vondi до совершенства

---

## Условия выполнения

> **КРИТИЧЕСКИ ВАЖНО:**
>
> 1. **Обратная совместимость НЕ требуется** — ломаем всё что нужно
> 2. **Даунтайм разрешён** — этап разработки, сервисы могут быть недоступны
> 3. **Рудименты ЗАПРЕЩЕНЫ** — удаляем всё устаревшее, не оставляем мёртвый код
> 4. **Единственный источник истины** — только Listings Microservice

---

## Обзор плана

```
┌─────────────────────────────────────────────────────────────────┐
│                    OPENSEARCH PERFECTION v2.0                   │
├─────────────────────────────────────────────────────────────────┤
│                                                                  │
│  ФАЗА 1: ЗАЧИСТКА РУДИМЕНТОВ              [~2 часа]             │
│  ├─ Удалить Python скрипты                                      │
│  ├─ Удалить индексацию из monolith                              │
│  └─ Удалить старые индексы в OpenSearch                         │
│                                                                  │
│  ФАЗА 2: НОВЫЕ МАППИНГИ                   [~5 часов] (+2ч)      │
│  ├─ Языковые анализаторы (sr/ru/en)                             │
│  ├─ ★ Синонимы (synonym filter)                                 │
│  ├─ ★ Транслитерация (сербский кир↔лат)                         │
│  ├─ ★ Trigram для Did You Mean                                  │
│  ├─ Variants (nested structure)                                  │
│  ├─ Недостающие поля                                            │
│  └─ Пересоздание индекса                                        │
│                                                                  │
│  ФАЗА 3: ИНДЕКСАТОР                       [~4 часа]             │
│  ├─ Загрузка variants                                           │
│  ├─ Все переводы (location, city, country)                      │
│  ├─ Новые поля (brand, rating, condition, etc)                  │
│  └─ Batch processing                                            │
│                                                                  │
│  ФАЗА 4: ПОИСКОВЫЕ ЗАПРОСЫ                [~6 часов] (+3ч)      │
│  ├─ Multi-language search                                       │
│  ├─ ★ Fuzzy search в основном поиске                            │
│  ├─ ★ Function Score (умное ранжирование)                       │
│  ├─ ★ Did You Mean (коррекция опечаток)                         │
│  ├─ ★ Highlighting (подсветка результатов)                      │
│  ├─ ★ Zero Results Fallback                                     │
│  ├─ ★ More Like This (похожие товары)                           │
│  ├─ ★ Query Rewriting (нормализация)                            │
│  ├─ Variant search                                              │
│  ├─ Advanced filters                                            │
│  ├─ ★ Расширенные фасеты                                        │
│  └─ Autocomplete                                                │
│                                                                  │
│  ФАЗА 5: ПЕРЕИНДЕКСАЦИЯ                   [~2 часа]             │
│  ├─ Blue-green strategy                                         │
│  ├─ Full reindex                                                │
│  └─ Verification                                                │
│                                                                  │
│  ФАЗА 6: МОНИТОРИНГ                       [~2 часа]             │
│  ├─ Health checks                                               │
│  ├─ Queue monitoring                                            │
│  └─ Alerts                                                      │
│                                                                  │
│  ★ ФАЗА 7: АНАЛИТИКА ПОИСКА               [~4 часа] (НОВАЯ)     │
│  ├─ Search Analytics индекс                                     │
│  ├─ Click-through tracking                                      │
│  ├─ Zero Results отчёты                                         │
│  └─ Trending searches                                           │
│                                                                  │
└─────────────────────────────────────────────────────────────────┘

ОБЩЕЕ ВРЕМЯ: ~25 часов (было 16)
★ = новые улучшения v2.0
```

---

## ФАЗА 1: Зачистка рудиментов

**Время:** ~2 часа
**Цель:** Удалить ВСЁ устаревшее, оставить единственный источник индексации

### 1.1 Удалить Python скрипты

| Файл | Действие | Причина |
|------|----------|---------|
| `vondi/backend/reindex_unified.py` | **УДАЛИТЬ** | Использует phantom индекс `unified_listings` |
| `vondi/backend/scripts/reindex_*.py` | **УДАЛИТЬ** | Дубликаты функционала microservice |

```bash
# Команды удаления
rm -f /p/github.com/vondi-global/vondi/backend/reindex_unified.py
rm -f /p/github.com/vondi-global/vondi/backend/scripts/reindex_*.py
```

### 1.2 Удалить индексацию из Monolith

| Файл | Действие |
|------|----------|
| `vondi/backend/internal/storage/opensearch/client.go` | Удалить методы IndexDocument, UpdateDocument, DeleteDocument |
| `vondi/backend/internal/storage/opensearch/mappings.go` | Удалить полностью (рудимент) |
| `vondi/backend/internal/storage/opensearch/` | Оставить только Search методы |

**Что оставить в monolith:**
- `Search()` - поиск (проксирует в microservice или напрямую в OS)
- `GetClient()` - получение клиента (для поиска)

**Что УДАЛИТЬ из monolith:**
- `CreateIndex()` - создание индекса
- `IndexDocument()` - индексация документа
- `UpdateDocument()` - обновление документа
- `DeleteDocument()` - удаление документа
- `BulkIndex()` - массовая индексация
- Все маппинги

### 1.3 Удалить старые индексы в OpenSearch

```bash
# Подключиться к OpenSearch
curl -X GET "http://localhost:9200/_cat/indices?v"

# Удалить phantom индексы
curl -X DELETE "http://localhost:9200/unified_listings"
curl -X DELETE "http://localhost:9200/listings_microservice"  # если есть
curl -X DELETE "http://localhost:9200/c2c_listings"           # если есть
curl -X DELETE "http://localhost:9200/b2c_listings"           # если есть

# Оставить ТОЛЬКО marketplace_listings (будет пересоздан в Фазе 2)
```

### 1.4 Очистить конфигурацию

**Monolith .env - УДАЛИТЬ:**
```bash
# УДАЛИТЬ эти переменные (индексация теперь только в microservice)
# OPENSEARCH_INDEX=...
# USE_OPENSEARCH_INDEXING=...
```

**Listings Microservice .env - ОСТАВИТЬ:**
```bash
VONDILISTINGS_OPENSEARCH_URL=http://localhost:9200
VONDILISTINGS_OPENSEARCH_INDEX=marketplace_listings
VONDILISTINGS_WORKER_ENABLED=true
VONDILISTINGS_WORKER_CONCURRENCY=5
```

---

## ФАЗА 2: Новые маппинги

**Время:** ~5 часов (было 3)
**Цель:** Создать идеальную схему индекса с синонимами и транслитерацией

### 2.1 Языковые анализаторы + Синонимы + Транслитерация

**Файл:** `listings/internal/repository/opensearch/mappings.go`

```go
var IndexSettings = map[string]interface{}{
    "settings": map[string]interface{}{
        "number_of_shards":   1,
        "number_of_replicas": 0,
        "analysis": map[string]interface{}{
            "char_filter": map[string]interface{}{
                // ★ НОВОЕ: Транслитерация сербский кириллица → латиница
                "serbian_cyrillic_to_latin": map[string]interface{}{
                    "type": "mapping",
                    "mappings": []string{
                        "а=>a", "б=>b", "в=>v", "г=>g", "д=>d",
                        "ђ=>đ", "е=>e", "ж=>ž", "з=>z", "и=>i",
                        "ј=>j", "к=>k", "л=>l", "љ=>lj", "м=>m",
                        "н=>n", "њ=>nj", "о=>o", "п=>p", "р=>r",
                        "с=>s", "т=>t", "ћ=>ć", "у=>u", "ф=>f",
                        "х=>h", "ц=>c", "ч=>č", "џ=>dž", "ш=>š",
                        "А=>A", "Б=>B", "В=>V", "Г=>G", "Д=>D",
                        "Ђ=>Đ", "Е=>E", "Ж=>Ž", "З=>Z", "И=>I",
                        "Ј=>J", "К=>K", "Л=>L", "Љ=>Lj", "М=>M",
                        "Н=>N", "Њ=>Nj", "О=>O", "П=>P", "Р=>R",
                        "С=>S", "Т=>T", "Ћ=>Ć", "У=>U", "Ф=>F",
                        "Х=>H", "Ц=>C", "Ч=>Č", "Џ=>Dž", "Ш=>Š",
                    },
                },
            },
            "analyzer": map[string]interface{}{
                // Сербский (латиница + кириллица с транслитерацией)
                "serbian_analyzer": map[string]interface{}{
                    "type":        "custom",
                    "char_filter": []string{"serbian_cyrillic_to_latin"},
                    "tokenizer":   "standard",
                    "filter":      []string{"lowercase", "serbian_normalization", "serbian_stemmer", "serbian_synonyms"},
                },
                // Русский
                "russian_analyzer": map[string]interface{}{
                    "type":      "custom",
                    "tokenizer": "standard",
                    "filter":    []string{"lowercase", "russian_stemmer", "russian_synonyms"},
                },
                // Английский
                "english_analyzer": map[string]interface{}{
                    "type":      "custom",
                    "tokenizer": "standard",
                    "filter":    []string{"lowercase", "english_stemmer", "english_possessive_stemmer", "english_synonyms"},
                },
                // ★ НОВОЕ: Универсальный анализатор с синонимами
                "universal_analyzer": map[string]interface{}{
                    "type":        "custom",
                    "char_filter": []string{"serbian_cyrillic_to_latin"},
                    "tokenizer":   "standard",
                    "filter":      []string{"lowercase", "icu_folding", "universal_synonyms"},
                },
                // Autocomplete
                "autocomplete_analyzer": map[string]interface{}{
                    "type":        "custom",
                    "char_filter": []string{"serbian_cyrillic_to_latin"},
                    "tokenizer":   "standard",
                    "filter":      []string{"lowercase", "autocomplete_filter"},
                },
                "autocomplete_search": map[string]interface{}{
                    "type":        "custom",
                    "char_filter": []string{"serbian_cyrillic_to_latin"},
                    "tokenizer":   "standard",
                    "filter":      []string{"lowercase"},
                },
                // ★ НОВОЕ: Trigram для Did You Mean
                "trigram_analyzer": map[string]interface{}{
                    "type":      "custom",
                    "tokenizer": "standard",
                    "filter":    []string{"lowercase", "shingle_filter"},
                },
            },
            "filter": map[string]interface{}{
                "serbian_normalization": map[string]interface{}{
                    "type": "icu_folding", // Нормализация диакритиков
                },
                "serbian_stemmer": map[string]interface{}{
                    "type":     "stemmer",
                    "language": "serbian",
                },
                "russian_stemmer": map[string]interface{}{
                    "type":     "stemmer",
                    "language": "russian",
                },
                "english_stemmer": map[string]interface{}{
                    "type":     "stemmer",
                    "language": "english",
                },
                "english_possessive_stemmer": map[string]interface{}{
                    "type":     "stemmer",
                    "language": "possessive_english",
                },
                "autocomplete_filter": map[string]interface{}{
                    "type":     "edge_ngram",
                    "min_gram": 2,
                    "max_gram": 20,
                },
                // ★ НОВОЕ: Shingle для trigram (Did You Mean)
                "shingle_filter": map[string]interface{}{
                    "type":             "shingle",
                    "min_shingle_size": 2,
                    "max_shingle_size": 3,
                },
                // ★ НОВОЕ: Синонимы для сербского
                "serbian_synonyms": map[string]interface{}{
                    "type": "synonym",
                    "synonyms": []string{
                        // Электроника
                        "telefon, mobilni, smartphone, pametni telefon",
                        "laptop, prenosni računar, notebook",
                        "televizor, tv, smart tv",
                        "slušalice, bubice, headphones, earphones",
                        "punjač, charger, adapter",
                        "tastatura, keyboard",
                        "miš, mouse",
                        // Одежда
                        "patike, sportska obuća, sneakers",
                        "cipele, obuća, shoes",
                        "majica, t-shirt, tshirt",
                        "haljina, dress",
                        "jakna, jacket, kaput",
                        "farmerke, jeans, džins",
                        "suknja, skirt",
                        // Дом
                        "krevet, bed",
                        "sto, table, stol",
                        "stolica, chair",
                        "ormar, wardrobe, plakar",
                        "frižider, hladnjak, refrigerator",
                        // Авто
                        "auto, automobil, kola, car, vozilo",
                        "motor, motocikl, motorcycle",
                        "gume, pneumatici, tires",
                    },
                },
                // ★ НОВОЕ: Синонимы для русского
                "russian_synonyms": map[string]interface{}{
                    "type": "synonym",
                    "synonyms": []string{
                        "телефон, смартфон, мобильный, мобилка",
                        "ноутбук, лэптоп, laptop",
                        "наушники, наушники, earphones, headphones",
                        "обувь, кроссовки, кеды, туфли",
                        "платье, сарафан",
                        "машина, авто, автомобиль, тачка",
                    },
                },
                // ★ НОВОЕ: Синонимы для английского
                "english_synonyms": map[string]interface{}{
                    "type": "synonym",
                    "synonyms": []string{
                        "phone, smartphone, mobile, cell",
                        "laptop, notebook, computer",
                        "headphones, earphones, earbuds, airpods",
                        "shoes, sneakers, trainers, footwear",
                        "dress, gown",
                        "car, auto, automobile, vehicle",
                    },
                },
                // ★ НОВОЕ: Универсальные синонимы (кросс-языковые бренды)
                "universal_synonyms": map[string]interface{}{
                    "type": "synonym",
                    "synonyms": []string{
                        // Бренды
                        "iphone, apple phone, ajfon",
                        "samsung, самсунг, sumsung",
                        "xiaomi, сяоми, mi",
                        "adidas, адидас",
                        "nike, найк, najk",
                        "zara, зара",
                        // Размеры
                        "xs, extra small",
                        "s, small, mali",
                        "m, medium, srednji",
                        "l, large, veliki",
                        "xl, extra large",
                        "xxl, 2xl",
                    },
                },
            },
        },
    },
}
```

### 2.2 Полная схема маппингов

```go
var IndexMappings = map[string]interface{}{
    "mappings": map[string]interface{}{
        "properties": map[string]interface{}{
            // ═══════════════════════════════════════════
            // ИДЕНТИФИКАТОРЫ
            // ═══════════════════════════════════════════
            "id":             map[string]interface{}{"type": "long"},
            "uuid":           map[string]interface{}{"type": "keyword"},
            "sku":            map[string]interface{}{"type": "keyword"},
            "storefront_id":  map[string]interface{}{"type": "long"},
            "category_id":    map[string]interface{}{"type": "keyword"},
            "category_slug":  map[string]interface{}{"type": "keyword"},
            "source_type":    map[string]interface{}{"type": "keyword"}, // c2c, b2c

            // ═══════════════════════════════════════════
            // МУЛЬТИЯЗЫЧНЫЕ ПОЛЯ (с анализаторами)
            // ═══════════════════════════════════════════

            // Title - основное поле с multiple sub-fields
            "title": map[string]interface{}{
                "type":     "text",
                "analyzer": "universal_analyzer",
                "fields": map[string]interface{}{
                    "keyword": map[string]interface{}{"type": "keyword"},
                    "autocomplete": map[string]interface{}{
                        "type":            "text",
                        "analyzer":        "autocomplete_analyzer",
                        "search_analyzer": "autocomplete_search",
                    },
                    // ★ НОВОЕ: Trigram для Did You Mean
                    "trigram": map[string]interface{}{
                        "type":     "text",
                        "analyzer": "trigram_analyzer",
                    },
                },
            },
            "title_sr": map[string]interface{}{
                "type":     "text",
                "analyzer": "serbian_analyzer",
            },
            "title_en": map[string]interface{}{
                "type":     "text",
                "analyzer": "english_analyzer",
            },
            "title_ru": map[string]interface{}{
                "type":     "text",
                "analyzer": "russian_analyzer",
            },

            // Description
            "description": map[string]interface{}{
                "type":     "text",
                "analyzer": "universal_analyzer",
            },
            "description_sr": map[string]interface{}{
                "type":     "text",
                "analyzer": "serbian_analyzer",
            },
            "description_en": map[string]interface{}{
                "type":     "text",
                "analyzer": "english_analyzer",
            },
            "description_ru": map[string]interface{}{
                "type":     "text",
                "analyzer": "russian_analyzer",
            },

            "original_language": map[string]interface{}{"type": "keyword"},

            // ═══════════════════════════════════════════
            // ЦЕНА И СКИДКИ
            // ═══════════════════════════════════════════
            "price":            map[string]interface{}{"type": "double"},
            "old_price":        map[string]interface{}{"type": "double"},
            "discount_percent": map[string]interface{}{"type": "integer"},
            "has_discount":     map[string]interface{}{"type": "boolean"},
            "currency":         map[string]interface{}{"type": "keyword"},

            // ═══════════════════════════════════════════
            // СТАТУС И НАЛИЧИЕ
            // ═══════════════════════════════════════════
            "status":       map[string]interface{}{"type": "keyword"},
            "visibility":   map[string]interface{}{"type": "keyword"},
            "quantity":     map[string]interface{}{"type": "integer"},
            "stock_status": map[string]interface{}{"type": "keyword"},
            "condition":    map[string]interface{}{"type": "keyword"}, // new, used, refurbished

            // ═══════════════════════════════════════════
            // БРЕНД И ТЕГИ
            // ═══════════════════════════════════════════
            "brand": map[string]interface{}{
                "type": "keyword",
                "fields": map[string]interface{}{
                    "text": map[string]interface{}{
                        "type":     "text",
                        "analyzer": "universal_analyzer",
                    },
                },
            },
            "tags": map[string]interface{}{
                "type": "keyword",
            },

            // ═══════════════════════════════════════════
            // РЕЙТИНГ И ОТЗЫВЫ
            // ═══════════════════════════════════════════
            "rating":       map[string]interface{}{"type": "float"},
            "review_count": map[string]interface{}{"type": "integer"},

            // ═══════════════════════════════════════════
            // ПРОДАВЕЦ (STOREFRONT)
            // ═══════════════════════════════════════════
            "storefront_name": map[string]interface{}{
                "type": "text",
                "fields": map[string]interface{}{
                    "keyword": map[string]interface{}{"type": "keyword"},
                },
            },
            "storefront_slug":   map[string]interface{}{"type": "keyword"},
            "storefront_rating": map[string]interface{}{"type": "float"},
            "seller_verified":   map[string]interface{}{"type": "boolean"},

            // ═══════════════════════════════════════════
            // ГЕОЛОКАЦИЯ (с переводами)
            // ═══════════════════════════════════════════
            "location": map[string]interface{}{"type": "geo_point"},
            "has_individual_location": map[string]interface{}{"type": "boolean"},
            "individual_latitude":     map[string]interface{}{"type": "double"},
            "individual_longitude":    map[string]interface{}{"type": "double"},

            "country":    map[string]interface{}{"type": "keyword"},
            "country_sr": map[string]interface{}{"type": "keyword"},
            "country_en": map[string]interface{}{"type": "keyword"},
            "country_ru": map[string]interface{}{"type": "keyword"},

            "city":    map[string]interface{}{"type": "keyword"},
            "city_sr": map[string]interface{}{"type": "keyword"},
            "city_en": map[string]interface{}{"type": "keyword"},
            "city_ru": map[string]interface{}{"type": "keyword"},

            "address":    map[string]interface{}{"type": "text"},
            "address_sr": map[string]interface{}{"type": "text", "analyzer": "serbian_analyzer"},
            "address_en": map[string]interface{}{"type": "text", "analyzer": "english_analyzer"},
            "address_ru": map[string]interface{}{"type": "text", "analyzer": "russian_analyzer"},

            // ═══════════════════════════════════════════
            // ДАТЫ
            // ═══════════════════════════════════════════
            "created_at":  map[string]interface{}{"type": "date"},
            "updated_at":  map[string]interface{}{"type": "date"},
            "expires_at":  map[string]interface{}{"type": "date"},
            "promoted_at": map[string]interface{}{"type": "date"},

            // ═══════════════════════════════════════════
            // СТАТИСТИКА
            // ═══════════════════════════════════════════
            "views_count":     map[string]interface{}{"type": "integer"},
            "favorites_count": map[string]interface{}{"type": "integer"},
            "orders_count":    map[string]interface{}{"type": "integer"},
            "popularity_score": map[string]interface{}{"type": "float"},

            // ═══════════════════════════════════════════
            // ПРОМО И ФЛАГИ
            // ═══════════════════════════════════════════
            "is_promoted":    map[string]interface{}{"type": "boolean"},
            "is_featured":    map[string]interface{}{"type": "boolean"},
            "is_new_arrival": map[string]interface{}{"type": "boolean"},
            "accepts_offers": map[string]interface{}{"type": "boolean"},

            // ═══════════════════════════════════════════
            // ДОСТАВКА
            // ═══════════════════════════════════════════
            "shipping_available": map[string]interface{}{"type": "boolean"},
            "shipping_free":      map[string]interface{}{"type": "boolean"},
            "shipping_price":     map[string]interface{}{"type": "double"},
            "delivery_days_min":  map[string]interface{}{"type": "integer"},
            "delivery_days_max":  map[string]interface{}{"type": "integer"},

            // ═══════════════════════════════════════════
            // ИЗОБРАЖЕНИЯ (nested)
            // ═══════════════════════════════════════════
            "images": map[string]interface{}{
                "type": "nested",
                "properties": map[string]interface{}{
                    "id":         map[string]interface{}{"type": "long"},
                    "url":        map[string]interface{}{"type": "keyword"},
                    "alt":        map[string]interface{}{"type": "text"},
                    "is_main":    map[string]interface{}{"type": "boolean"},
                    "sort_order": map[string]interface{}{"type": "integer"},
                },
            },
            "main_image_url": map[string]interface{}{"type": "keyword"},
            "images_count":   map[string]interface{}{"type": "integer"},

            // ═══════════════════════════════════════════
            // АТРИБУТЫ (nested)
            // ═══════════════════════════════════════════
            "attributes": map[string]interface{}{
                "type": "nested",
                "properties": map[string]interface{}{
                    "id":             map[string]interface{}{"type": "long"},
                    "code":           map[string]interface{}{"type": "keyword"},
                    "name":           map[string]interface{}{"type": "keyword"},
                    "name_sr":        map[string]interface{}{"type": "keyword"},
                    "name_en":        map[string]interface{}{"type": "keyword"},
                    "name_ru":        map[string]interface{}{"type": "keyword"},
                    "value_text":     map[string]interface{}{"type": "text"},
                    "value_number":   map[string]interface{}{"type": "double"},
                    "value_boolean":  map[string]interface{}{"type": "boolean"},
                    "value_keyword":  map[string]interface{}{"type": "keyword"},
                    "is_searchable":  map[string]interface{}{"type": "boolean"},
                    "is_filterable":  map[string]interface{}{"type": "boolean"},
                },
            },

            // Денормализованный текст для full-text поиска
            "all_attributes_text": map[string]interface{}{"type": "text"},

            // ═══════════════════════════════════════════
            // ВАРИАНТЫ ТОВАРА (nested)
            // ═══════════════════════════════════════════
            "variants": map[string]interface{}{
                "type": "nested",
                "properties": map[string]interface{}{
                    "id":             map[string]interface{}{"type": "long"},
                    "sku":            map[string]interface{}{"type": "keyword"},
                    "price":          map[string]interface{}{"type": "double"},
                    "old_price":      map[string]interface{}{"type": "double"},
                    "stock_quantity": map[string]interface{}{"type": "integer"},
                    "status":         map[string]interface{}{"type": "keyword"},
                    "is_default":     map[string]interface{}{"type": "boolean"},
                    "weight":         map[string]interface{}{"type": "double"},
                    "dimensions": map[string]interface{}{
                        "properties": map[string]interface{}{
                            "length": map[string]interface{}{"type": "double"},
                            "width":  map[string]interface{}{"type": "double"},
                            "height": map[string]interface{}{"type": "double"},
                        },
                    },
                    // Атрибуты варианта (цвет, размер)
                    "attributes": map[string]interface{}{
                        "type": "nested",
                        "properties": map[string]interface{}{
                            "code":  map[string]interface{}{"type": "keyword"},
                            "name":  map[string]interface{}{"type": "keyword"},
                            "value": map[string]interface{}{"type": "keyword"},
                        },
                    },
                    // Изображение варианта
                    "image_url": map[string]interface{}{"type": "keyword"},
                },
            },

            // Денормализованные поля для быстрой фильтрации
            "variant_skus":   map[string]interface{}{"type": "keyword"},     // все SKU вариантов
            "variant_colors": map[string]interface{}{"type": "keyword"},     // все цвета
            "variant_sizes":  map[string]interface{}{"type": "keyword"},     // все размеры
            "has_variants":   map[string]interface{}{"type": "boolean"},
            "variants_count": map[string]interface{}{"type": "integer"},
            "min_price":      map[string]interface{}{"type": "double"},      // мин. цена среди вариантов
            "max_price":      map[string]interface{}{"type": "double"},      // макс. цена среди вариантов
            "total_stock":    map[string]interface{}{"type": "integer"},     // общий сток всех вариантов

            // ═══════════════════════════════════════════
            // КАТЕГОРИИ (путь)
            // ═══════════════════════════════════════════
            "category_path":     map[string]interface{}{"type": "keyword"},  // ["electronics", "phones", "smartphones"]
            "category_path_ids": map[string]interface{}{"type": "long"},     // [1, 5, 23]
            "category_names": map[string]interface{}{
                "properties": map[string]interface{}{
                    "sr": map[string]interface{}{"type": "keyword"},
                    "en": map[string]interface{}{"type": "keyword"},
                    "ru": map[string]interface{}{"type": "keyword"},
                },
            },
        },
    },
}
```

### 2.3 Пересоздание индекса

```go
// Новый метод в client.go
func (c *Client) RecreateIndex(ctx context.Context) error {
    indexName := c.config.Index

    // 1. Удалить старый индекс
    _, err := c.client.Indices.Delete([]string{indexName})
    if err != nil {
        // Игнорируем если не существует
    }

    // 2. Создать с новыми настройками
    settings := mergeSettingsAndMappings(IndexSettings, IndexMappings)

    res, err := c.client.Indices.Create(
        indexName,
        c.client.Indices.Create.WithBody(esutil.NewJSONReader(settings)),
    )
    if err != nil {
        return fmt.Errorf("create index: %w", err)
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("create index response: %s", res.String())
    }

    return nil
}
```

---

## ФАЗА 3: Индексатор

**Время:** ~4 часа
**Цель:** Полностью переписать логику индексации

### 3.1 Загрузка вариантов

**Файл:** `listings/internal/indexer/listing_indexer.go`

```go
// Новый метод загрузки вариантов
func (idx *ListingIndexer) loadVariantsForListing(ctx context.Context, listingID int64) ([]*domain.ProductVariant, error) {
    variants, err := idx.variantRepo.GetByListingID(ctx, listingID)
    if err != nil {
        return nil, fmt.Errorf("load variants: %w", err)
    }
    return variants, nil
}

// Обновлённый buildListingDocument с вариантами
func (idx *ListingIndexer) buildListingDocument(
    ctx context.Context,
    listing *domain.Listing,
) (map[string]interface{}, error) {

    doc := map[string]interface{}{
        // ... основные поля ...
    }

    // ═══════════════════════════════════════════
    // ВАРИАНТЫ
    // ═══════════════════════════════════════════
    variants, err := idx.loadVariantsForListing(ctx, listing.ID)
    if err != nil {
        idx.logger.Warn().Err(err).Int64("listing_id", listing.ID).Msg("failed to load variants")
    }

    if len(variants) > 0 {
        doc["has_variants"] = true
        doc["variants_count"] = len(variants)

        variantDocs := make([]map[string]interface{}, 0, len(variants))
        variantSKUs := make([]string, 0, len(variants))
        variantColors := make([]string, 0)
        variantSizes := make([]string, 0)

        var minPrice, maxPrice float64 = math.MaxFloat64, 0
        var totalStock int32 = 0

        for _, v := range variants {
            variantDoc := map[string]interface{}{
                "id":             v.ID,
                "sku":            v.SKU,
                "price":          v.Price,
                "old_price":      v.OldPrice,
                "stock_quantity": v.StockQuantity,
                "status":         v.Status,
                "is_default":     v.IsDefault,
                "weight":         v.Weight,
                "image_url":      v.ImageURL,
            }

            // Dimensions
            if v.Length > 0 || v.Width > 0 || v.Height > 0 {
                variantDoc["dimensions"] = map[string]interface{}{
                    "length": v.Length,
                    "width":  v.Width,
                    "height": v.Height,
                }
            }

            // Атрибуты варианта
            if len(v.Attributes) > 0 {
                attrs := make([]map[string]interface{}, 0, len(v.Attributes))
                for _, attr := range v.Attributes {
                    attrs = append(attrs, map[string]interface{}{
                        "code":  attr.Code,
                        "name":  attr.Name,
                        "value": attr.Value,
                    })

                    // Собираем цвета и размеры для денормализации
                    if attr.Code == "color" {
                        variantColors = append(variantColors, attr.Value)
                    }
                    if attr.Code == "size" {
                        variantSizes = append(variantSizes, attr.Value)
                    }
                }
                variantDoc["attributes"] = attrs
            }

            variantDocs = append(variantDocs, variantDoc)
            variantSKUs = append(variantSKUs, v.SKU)

            // Min/Max price
            if v.Price < minPrice {
                minPrice = v.Price
            }
            if v.Price > maxPrice {
                maxPrice = v.Price
            }

            totalStock += v.StockQuantity
        }

        doc["variants"] = variantDocs
        doc["variant_skus"] = unique(variantSKUs)
        doc["variant_colors"] = unique(variantColors)
        doc["variant_sizes"] = unique(variantSizes)
        doc["min_price"] = minPrice
        doc["max_price"] = maxPrice
        doc["total_stock"] = totalStock
    } else {
        doc["has_variants"] = false
        doc["variants_count"] = 0
        doc["min_price"] = listing.Price
        doc["max_price"] = listing.Price
        doc["total_stock"] = listing.Quantity
    }

    return doc, nil
}
```

### 3.2 Индексация переводов локации

```go
// В buildListingDocument добавить:

// ═══════════════════════════════════════════
// ПЕРЕВОДЫ ЛОКАЦИИ
// ═══════════════════════════════════════════

// Country translations
if len(listing.CountryTranslations) > 0 {
    for lang, translation := range listing.CountryTranslations {
        if translation != "" {
            doc["country_"+lang] = translation
        }
    }
}

// City translations
if len(listing.CityTranslations) > 0 {
    for lang, translation := range listing.CityTranslations {
        if translation != "" {
            doc["city_"+lang] = translation
        }
    }
}

// Address translations
if len(listing.AddressTranslations) > 0 {
    for lang, translation := range listing.AddressTranslations {
        if translation != "" {
            doc["address_"+lang] = translation
        }
    }
}
```

### 3.3 Новые поля

```go
// В buildListingDocument добавить:

// ═══════════════════════════════════════════
// ЦЕНА И СКИДКИ
// ═══════════════════════════════════════════
doc["price"] = listing.Price
doc["old_price"] = listing.OldPrice
if listing.OldPrice != nil && *listing.OldPrice > listing.Price {
    discount := int(((*listing.OldPrice - listing.Price) / *listing.OldPrice) * 100)
    doc["discount_percent"] = discount
    doc["has_discount"] = true
} else {
    doc["discount_percent"] = 0
    doc["has_discount"] = false
}

// ═══════════════════════════════════════════
// БРЕНД
// ═══════════════════════════════════════════
doc["brand"] = extractBrand(listing.Attributes)

// ═══════════════════════════════════════════
// СОСТОЯНИЕ ТОВАРА
// ═══════════════════════════════════════════
doc["condition"] = listing.Condition // new, used, refurbished

// ═══════════════════════════════════════════
// РЕЙТИНГ (если есть)
// ═══════════════════════════════════════════
if listing.Rating != nil {
    doc["rating"] = *listing.Rating
}
doc["review_count"] = listing.ReviewCount

// ═══════════════════════════════════════════
// ПРОДАВЕЦ
// ═══════════════════════════════════════════
if listing.Storefront != nil {
    doc["storefront_name"] = listing.Storefront.Name
    doc["storefront_slug"] = listing.Storefront.Slug
    doc["storefront_rating"] = listing.Storefront.Rating
    doc["seller_verified"] = listing.Storefront.IsVerified
}

// ═══════════════════════════════════════════
// ДАТЫ
// ═══════════════════════════════════════════
doc["expires_at"] = listing.ExpiresAt
doc["promoted_at"] = listing.PromotedAt

// ═══════════════════════════════════════════
// ТЕГИ
// ═══════════════════════════════════════════
if len(listing.Tags) > 0 {
    doc["tags"] = listing.Tags
}

// ═══════════════════════════════════════════
// ДОСТАВКА
// ═══════════════════════════════════════════
doc["shipping_available"] = listing.ShippingAvailable
doc["shipping_free"] = listing.ShippingFree
doc["shipping_price"] = listing.ShippingPrice
doc["delivery_days_min"] = listing.DeliveryDaysMin
doc["delivery_days_max"] = listing.DeliveryDaysMax

// ═══════════════════════════════════════════
// ФЛАГИ
// ═══════════════════════════════════════════
doc["is_promoted"] = listing.IsPromoted
doc["is_featured"] = listing.IsFeatured
doc["is_new_arrival"] = isNewArrival(listing.CreatedAt) // < 7 дней
doc["accepts_offers"] = listing.AcceptsOffers

// ═══════════════════════════════════════════
// КАТЕГОРИИ (путь)
// ═══════════════════════════════════════════
if listing.CategoryPath != nil {
    doc["category_path"] = listing.CategoryPath.Slugs
    doc["category_path_ids"] = listing.CategoryPath.IDs
    doc["category_names"] = map[string]interface{}{
        "sr": listing.CategoryPath.Names["sr"],
        "en": listing.CategoryPath.Names["en"],
        "ru": listing.CategoryPath.Names["ru"],
    }
}

// ═══════════════════════════════════════════
// POPULARITY SCORE (рассчитывается)
// ═══════════════════════════════════════════
doc["popularity_score"] = calculatePopularityScore(listing)
```

### 3.4 Batch processing с правильным размером

```go
func (idx *ListingIndexer) BulkIndexListings(ctx context.Context, listings []*domain.Listing) error {
    const batchSize = 500 // Не более 500 за раз

    for i := 0; i < len(listings); i += batchSize {
        end := i + batchSize
        if end > len(listings) {
            end = len(listings)
        }

        batch := listings[i:end]

        if err := idx.indexBatch(ctx, batch); err != nil {
            return fmt.Errorf("batch %d-%d: %w", i, end, err)
        }

        idx.logger.Info().
            Int("from", i).
            Int("to", end).
            Int("total", len(listings)).
            Msg("batch indexed")
    }

    return nil
}

func (idx *ListingIndexer) indexBatch(ctx context.Context, listings []*domain.Listing) error {
    var buf bytes.Buffer

    for _, listing := range listings {
        doc, err := idx.buildListingDocument(ctx, listing)
        if err != nil {
            idx.logger.Error().Err(err).Int64("id", listing.ID).Msg("skip listing")
            continue
        }

        // Action line
        meta := map[string]interface{}{
            "index": map[string]interface{}{
                "_index": idx.config.Index,
                "_id":    strconv.FormatInt(listing.ID, 10),
            },
        }

        if err := json.NewEncoder(&buf).Encode(meta); err != nil {
            return err
        }
        if err := json.NewEncoder(&buf).Encode(doc); err != nil {
            return err
        }
    }

    res, err := idx.client.Bulk(
        bytes.NewReader(buf.Bytes()),
        idx.client.Bulk.WithContext(ctx),
        idx.client.Bulk.WithRefresh("false"), // Не refresh каждый раз
    )
    if err != nil {
        return fmt.Errorf("bulk request: %w", err)
    }
    defer res.Body.Close()

    if res.IsError() {
        return fmt.Errorf("bulk response error: %s", res.String())
    }

    return nil
}
```

---

## ФАЗА 4: Поисковые запросы

**Время:** ~6 часов (было 3)
**Цель:** Реализовать умный поиск с учётом всех новых возможностей

### 4.1 Multi-language search с Fuzzy (★ УЛУЧШЕНО)

```go
func (c *Client) BuildSearchQuery(params SearchParams) map[string]interface{} {
    must := make([]map[string]interface{}, 0)
    filter := make([]map[string]interface{}, 0)
    should := make([]map[string]interface{}, 0)

    // ═══════════════════════════════════════════
    // ★ QUERY REWRITING (нормализация запроса)
    // ═══════════════════════════════════════════
    normalizedQuery := normalizeQuery(params.Query)

    // ═══════════════════════════════════════════
    // МУЛЬТИЯЗЫЧНЫЙ ПОИСК ПО ТЕКСТУ
    // ═══════════════════════════════════════════
    if normalizedQuery != "" {
        // Определяем язык запроса
        queryLang := detectLanguage(normalizedQuery) // "sr", "ru", "en"

        // Multi-match по всем языковым полям с boost
        textQuery := map[string]interface{}{
            "bool": map[string]interface{}{
                "should": []map[string]interface{}{
                    // ★ Основное поле title с FUZZY
                    {
                        "multi_match": map[string]interface{}{
                            "query":  normalizedQuery,
                            "fields": []string{
                                "title^3",
                                "title_" + queryLang + "^5", // Язык запроса - максимальный boost
                                "title_sr^2",
                                "title_en^2",
                                "title_ru^2",
                            },
                            "type":           "best_fields",
                            "operator":       "or",
                            "fuzziness":      "AUTO",      // ★ FUZZY ДОБАВЛЕН!
                            "prefix_length":  2,           // Первые 2 символа без fuzzy
                            "max_expansions": 50,
                        },
                    },
                    // Description с fuzzy
                    {
                        "multi_match": map[string]interface{}{
                            "query":  normalizedQuery,
                            "fields": []string{
                                "description",
                                "description_" + queryLang + "^2",
                                "description_sr",
                                "description_en",
                                "description_ru",
                            },
                            "type":      "best_fields",
                            "fuzziness": "AUTO",
                        },
                    },
                    // SKU (exact match - высший приоритет)
                    {
                        "term": map[string]interface{}{
                            "sku": map[string]interface{}{
                                "value": normalizedQuery,
                                "boost": 10,
                            },
                        },
                    },
                    // Variant SKUs
                    {
                        "term": map[string]interface{}{
                            "variant_skus": map[string]interface{}{
                                "value": normalizedQuery,
                                "boost": 10,
                            },
                        },
                    },
                    // Brand с fuzzy
                    {
                        "match": map[string]interface{}{
                            "brand.text": map[string]interface{}{
                                "query":     normalizedQuery,
                                "boost":     3,
                                "fuzziness": "AUTO",
                            },
                        },
                    },
                    // Tags
                    {
                        "term": map[string]interface{}{
                            "tags": map[string]interface{}{
                                "value": strings.ToLower(normalizedQuery),
                                "boost": 2,
                            },
                        },
                    },
                    // Storefront name
                    {
                        "match": map[string]interface{}{
                            "storefront_name": map[string]interface{}{
                                "query":     normalizedQuery,
                                "boost":     2,
                                "fuzziness": "AUTO",
                            },
                        },
                    },
                    // All attributes text
                    {
                        "match": map[string]interface{}{
                            "all_attributes_text": map[string]interface{}{
                                "query":     normalizedQuery,
                                "fuzziness": "AUTO",
                            },
                        },
                    },
                },
                "minimum_should_match": 1,
            },
        }
        must = append(must, textQuery)
    }

    // Статус = active (обязательный фильтр)
    filter = append(filter, map[string]interface{}{
        "term": map[string]interface{}{"status": "active"},
    })

    // ... остальные фильтры ...

    return map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must":   must,
                "filter": filter,
                "should": should,
            },
        },
    }
}

// ★ НОВОЕ: Нормализация запроса
func normalizeQuery(query string) string {
    // 1. Trim и lowercase
    query = strings.TrimSpace(strings.ToLower(query))

    // 2. Убрать множественные пробелы
    query = regexp.MustCompile(`\s+`).ReplaceAllString(query, " ")

    // 3. Транслитерация сербской кириллицы
    query = transliterateSerbianCyrillic(query)

    // 4. Удалить специальные символы (кроме букв, цифр, пробелов, дефисов)
    query = regexp.MustCompile(`[^\p{L}\p{N}\s-]`).ReplaceAllString(query, "")

    return query
}

func transliterateSerbianCyrillic(s string) string {
    replacer := strings.NewReplacer(
        "а", "a", "б", "b", "в", "v", "г", "g", "д", "d",
        "ђ", "đ", "е", "e", "ж", "ž", "з", "z", "и", "i",
        "ј", "j", "к", "k", "л", "l", "љ", "lj", "м", "m",
        "н", "n", "њ", "nj", "о", "o", "п", "p", "р", "r",
        "с", "s", "т", "t", "ћ", "ć", "у", "u", "ф", "f",
        "х", "h", "ц", "c", "ч", "č", "џ", "dž", "ш", "š",
    )
    return replacer.Replace(s)
}
```

### 4.2 ★ Function Score (умное ранжирование)

```go
func (c *Client) BuildRankedSearchQuery(params SearchParams) map[string]interface{} {
    baseQuery := c.BuildSearchQuery(params)

    return map[string]interface{}{
        "query": map[string]interface{}{
            "function_score": map[string]interface{}{
                "query": baseQuery["query"],
                "functions": []map[string]interface{}{
                    // ★ Boost для продвигаемых товаров (x10)
                    {
                        "filter": map[string]interface{}{
                            "term": map[string]interface{}{"is_promoted": true},
                        },
                        "weight": 10,
                    },
                    // ★ Boost для featured товаров (x5)
                    {
                        "filter": map[string]interface{}{
                            "term": map[string]interface{}{"is_featured": true},
                        },
                        "weight": 5,
                    },
                    // ★ Boost по популярности (logarithm)
                    {
                        "field_value_factor": map[string]interface{}{
                            "field":    "popularity_score",
                            "factor":   1.2,
                            "modifier": "log1p",
                            "missing":  1,
                        },
                    },
                    // ★ Decay по свежести (товары старше 30 дней теряют score)
                    {
                        "gauss": map[string]interface{}{
                            "created_at": map[string]interface{}{
                                "origin": "now",
                                "scale":  "30d",
                                "offset": "7d",
                                "decay":  0.5,
                            },
                        },
                    },
                    // ★ Boost для товаров с высоким рейтингом (4+)
                    {
                        "filter": map[string]interface{}{
                            "range": map[string]interface{}{
                                "rating": map[string]interface{}{"gte": 4.0},
                            },
                        },
                        "weight": 1.5,
                    },
                    // ★ Boost для верифицированных продавцов
                    {
                        "filter": map[string]interface{}{
                            "term": map[string]interface{}{"seller_verified": true},
                        },
                        "weight": 1.3,
                    },
                    // ★ Boost для товаров с изображениями (3+)
                    {
                        "filter": map[string]interface{}{
                            "range": map[string]interface{}{
                                "images_count": map[string]interface{}{"gte": 3},
                            },
                        },
                        "weight": 1.2,
                    },
                    // ★ Boost для товаров со скидкой
                    {
                        "filter": map[string]interface{}{
                            "term": map[string]interface{}{"has_discount": true},
                        },
                        "weight": 1.1,
                    },
                },
                "score_mode": "sum",
                "boost_mode": "multiply",
            },
        },
        "size": params.Limit,
        "from": params.Offset,
    }
}
```

### 4.3 ★ Did You Mean (коррекция опечаток)

```go
func (c *Client) BuildDidYouMeanQuery(query string) map[string]interface{} {
    return map[string]interface{}{
        "suggest": map[string]interface{}{
            "text": query,
            "phrase_suggestion": map[string]interface{}{
                "phrase": map[string]interface{}{
                    "field":          "title.trigram",
                    "size":           3,
                    "gram_size":      3,
                    "confidence":     1.0,
                    "max_errors":     2.0,
                    "direct_generator": []map[string]interface{}{
                        {
                            "field":           "title.trigram",
                            "suggest_mode":    "always",
                            "min_word_length": 3,
                        },
                    },
                    "collate": map[string]interface{}{
                        "query": map[string]interface{}{
                            "source": map[string]interface{}{
                                "match": map[string]interface{}{
                                    "title": "{{suggestion}}",
                                },
                            },
                        },
                        "prune": true,
                    },
                    "highlight": map[string]interface{}{
                        "pre_tag":  "<em>",
                        "post_tag": "</em>",
                    },
                },
            },
        },
    }
}

// Использование:
func (s *SearchService) SearchWithDidYouMean(ctx context.Context, params SearchParams) (*SearchResponse, error) {
    // 1. Основной поиск
    result, err := s.Search(ctx, params)
    if err != nil {
        return nil, err
    }

    // 2. Если мало результатов - получить suggestions
    if result.Total < 5 && params.Query != "" {
        suggestions, err := s.client.Suggest(ctx, params.Query)
        if err == nil && len(suggestions) > 0 {
            result.DidYouMean = suggestions[0]
        }
    }

    return result, nil
}
```

### 4.4 ★ Highlighting (подсветка результатов)

```go
func (c *Client) AddHighlighting(query map[string]interface{}) map[string]interface{} {
    query["highlight"] = map[string]interface{}{
        "pre_tags":  []string{"<mark>"},
        "post_tags": []string{"</mark>"},
        "fields": map[string]interface{}{
            "title": map[string]interface{}{
                "number_of_fragments": 0, // Весь title
            },
            "description": map[string]interface{}{
                "fragment_size":       150,
                "number_of_fragments": 3,
            },
            "brand.text": map[string]interface{}{
                "number_of_fragments": 0,
            },
            "all_attributes_text": map[string]interface{}{
                "fragment_size":       100,
                "number_of_fragments": 2,
            },
        },
    }
    return query
}
```

### 4.5 ★ Zero Results Fallback

```go
func (s *SearchService) SearchWithFallback(ctx context.Context, params SearchParams) (*SearchResponse, error) {
    // 1. Основной поиск (точный)
    result, err := s.Search(ctx, params)
    if err != nil {
        return nil, err
    }

    if result.Total > 0 {
        result.SearchType = "exact"
        return result, nil
    }

    // 2. Если 0 результатов - попробовать fuzzy
    if params.Query != "" {
        params.Fuzziness = "AUTO"
        result, err = s.Search(ctx, params)
        if err != nil {
            return nil, err
        }

        if result.Total > 0 {
            result.SearchType = "fuzzy"
            return result, nil
        }
    }

    // 3. Если всё ещё 0 - показать популярные в категории
    if params.CategoryID > 0 {
        result, err = s.GetPopularInCategory(ctx, params.CategoryID, 12)
        if err != nil {
            return nil, err
        }

        if result.Total > 0 {
            result.SearchType = "popular_in_category"
            result.FallbackReason = "no_results_showing_popular"
            return result, nil
        }
    }

    // 4. Fallback - показать новые товары
    result, err = s.GetNewArrivals(ctx, 12)
    if err != nil {
        return nil, err
    }
    result.SearchType = "new_arrivals"
    result.FallbackReason = "no_results_showing_new"

    return result, nil
}
```

### 4.6 ★ More Like This (похожие товары)

```go
func (c *Client) BuildMoreLikeThisQuery(listingID int64, categoryID string, limit int) map[string]interface{} {
    return map[string]interface{}{
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{
                    {
                        "more_like_this": map[string]interface{}{
                            "fields": []string{
                                "title", "title_sr", "title_en", "title_ru",
                                "description", "tags", "brand",
                            },
                            "like": []map[string]interface{}{
                                {
                                    "_index": c.config.Index,
                                    "_id":    strconv.FormatInt(listingID, 10),
                                },
                            },
                            "min_term_freq":        1,
                            "min_doc_freq":         2,
                            "max_query_terms":      25,
                            "minimum_should_match": "30%",
                        },
                    },
                },
                "filter": []map[string]interface{}{
                    {"term": map[string]interface{}{"status": "active"}},
                    {"term": map[string]interface{}{"category_id": categoryID}},
                },
                "must_not": []map[string]interface{}{
                    {"term": map[string]interface{}{"id": listingID}}, // Исключить текущий
                },
            },
        },
        "size": limit,
    }
}
```

### 4.7 Variant search (nested query)

```go
// Поиск по атрибутам вариантов
func (c *Client) buildVariantFilter(params SearchParams) map[string]interface{} {
    variantFilters := make([]map[string]interface{}, 0)

    // Фильтр по цвету
    if len(params.Colors) > 0 {
        variantFilters = append(variantFilters, map[string]interface{}{
            "terms": map[string]interface{}{
                "variants.attributes.value": params.Colors,
            },
        })
    }

    // Фильтр по размеру
    if len(params.Sizes) > 0 {
        variantFilters = append(variantFilters, map[string]interface{}{
            "terms": map[string]interface{}{
                "variants.attributes.value": params.Sizes,
            },
        })
    }

    // Фильтр по наличию варианта
    if params.InStockOnly {
        variantFilters = append(variantFilters, map[string]interface{}{
            "range": map[string]interface{}{
                "variants.stock_quantity": map[string]interface{}{
                    "gt": 0,
                },
            },
        })
    }

    if len(variantFilters) == 0 {
        return nil
    }

    return map[string]interface{}{
        "nested": map[string]interface{}{
            "path": "variants",
            "query": map[string]interface{}{
                "bool": map[string]interface{}{
                    "must": variantFilters,
                },
            },
        },
    }
}
```

### 4.8 ★ Расширенные фасеты (Aggregations)

```go
func buildExtendedAggregations() map[string]interface{} {
    return map[string]interface{}{
        // ═══════════════════════════════════════════
        // СУЩЕСТВУЮЩИЕ
        // ═══════════════════════════════════════════

        // Категории
        "categories": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "category_id",
                "size":  50,
            },
        },

        // Цена (гистограмма)
        "price_ranges": map[string]interface{}{
            "histogram": map[string]interface{}{
                "field":    "min_price",
                "interval": 1000,
            },
        },

        // Типы источников
        "source_types": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "source_type",
            },
        },

        // Статусы стока
        "stock_statuses": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "stock_status",
            },
        },

        // Атрибуты (nested)
        "attributes": map[string]interface{}{
            "nested": map[string]interface{}{
                "path": "attributes",
            },
            "aggs": map[string]interface{}{
                "attribute_keys": map[string]interface{}{
                    "terms": map[string]interface{}{
                        "field": "attributes.code",
                        "size":  50,
                    },
                    "aggs": map[string]interface{}{
                        "attribute_values": map[string]interface{}{
                            "terms": map[string]interface{}{
                                "field": "attributes.value_keyword",
                                "size":  20,
                            },
                        },
                    },
                },
            },
        },

        // ═══════════════════════════════════════════
        // ★ НОВЫЕ ФАСЕТЫ
        // ═══════════════════════════════════════════

        // Бренды
        "brands": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "brand",
                "size":  30,
            },
        },

        // Скидки (диапазоны)
        "discounts": map[string]interface{}{
            "range": map[string]interface{}{
                "field": "discount_percent",
                "ranges": []map[string]interface{}{
                    {"key": "10-20%", "from": 10, "to": 20},
                    {"key": "20-30%", "from": 20, "to": 30},
                    {"key": "30-50%", "from": 30, "to": 50},
                    {"key": "50%+", "from": 50},
                },
            },
        },

        // Рейтинг
        "ratings": map[string]interface{}{
            "range": map[string]interface{}{
                "field": "rating",
                "ranges": []map[string]interface{}{
                    {"key": "4+ stars", "from": 4},
                    {"key": "3+ stars", "from": 3},
                    {"key": "2+ stars", "from": 2},
                },
            },
        },

        // Состояние товара
        "conditions": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "condition",
            },
        },

        // Города
        "cities": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "city",
                "size":  50,
            },
        },

        // Доставка
        "shipping": map[string]interface{}{
            "filters": map[string]interface{}{
                "filters": map[string]interface{}{
                    "free_shipping": map[string]interface{}{
                        "term": map[string]interface{}{"shipping_free": true},
                    },
                    "has_shipping": map[string]interface{}{
                        "term": map[string]interface{}{"shipping_available": true},
                    },
                },
            },
        },

        // Цвета вариантов
        "variant_colors": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "variant_colors",
                "size":  20,
            },
        },

        // Размеры вариантов
        "variant_sizes": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "variant_sizes",
                "size":  30,
            },
        },

        // Верифицированные продавцы
        "seller_verified": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "seller_verified",
            },
        },

        // С вариантами / без
        "has_variants": map[string]interface{}{
            "terms": map[string]interface{}{
                "field": "has_variants",
            },
        },

        // Статистика цен
        "price_stats": map[string]interface{}{
            "stats": map[string]interface{}{
                "field": "min_price",
            },
        },
    }
}
```

### 4.9 Autocomplete

```go
func (c *Client) Autocomplete(ctx context.Context, query string, limit int) ([]string, error) {
    searchQuery := map[string]interface{}{
        "size": 0,
        "query": map[string]interface{}{
            "match": map[string]interface{}{
                "title.autocomplete": map[string]interface{}{
                    "query":    query,
                    "operator": "and",
                },
            },
        },
        "aggs": map[string]interface{}{
            "suggestions": map[string]interface{}{
                "terms": map[string]interface{}{
                    "field": "title.keyword",
                    "size":  limit,
                },
            },
        },
    }

    res, err := c.client.Search(
        c.client.Search.WithContext(ctx),
        c.client.Search.WithIndex(c.config.Index),
        c.client.Search.WithBody(esutil.NewJSONReader(searchQuery)),
    )
    if err != nil {
        return nil, err
    }
    defer res.Body.Close()

    // Parse response and extract suggestions
    // ...

    return suggestions, nil
}
```

---

## ФАЗА 5: Переиндексация

**Время:** ~2 часа
**Цель:** Безопасная полная переиндексация

### 5.1 Blue-Green Strategy

```go
func (idx *ListingIndexer) SafeReindexAll(ctx context.Context) error {
    // 1. Создать новый индекс с версией
    newIndexName := fmt.Sprintf("%s_v%d", idx.config.Index, time.Now().Unix())

    if err := idx.client.CreateIndexWithName(ctx, newIndexName); err != nil {
        return fmt.Errorf("create new index: %w", err)
    }

    idx.logger.Info().Str("index", newIndexName).Msg("new index created")

    // 2. Полная индексация в новый индекс
    tempConfig := idx.config
    tempConfig.Index = newIndexName

    tempIndexer := NewListingIndexer(idx.client, tempConfig, idx.logger, idx.repos)

    if err := tempIndexer.ReindexAllWithAttributes(ctx, 500); err != nil {
        // Удалить неудачный индекс
        _ = idx.client.DeleteIndex(ctx, newIndexName)
        return fmt.Errorf("reindex: %w", err)
    }

    // 3. Проверить количество документов
    newCount, err := idx.client.CountDocuments(ctx, newIndexName)
    if err != nil {
        return fmt.Errorf("count new: %w", err)
    }

    oldCount, _ := idx.client.CountDocuments(ctx, idx.config.Index)

    if newCount < oldCount*90/100 { // Меньше 90% - что-то не так
        _ = idx.client.DeleteIndex(ctx, newIndexName)
        return fmt.Errorf("new index has fewer documents: %d vs %d", newCount, oldCount)
    }

    // 4. Переключить alias
    aliasName := idx.config.Index // marketplace_listings

    if err := idx.client.SwitchAlias(ctx, aliasName, newIndexName); err != nil {
        return fmt.Errorf("switch alias: %w", err)
    }

    idx.logger.Info().
        Str("alias", aliasName).
        Str("new_index", newIndexName).
        Int("documents", newCount).
        Msg("alias switched")

    // 5. Удалить старые индексы (оставить последние 2)
    if err := idx.cleanupOldIndices(ctx); err != nil {
        idx.logger.Warn().Err(err).Msg("cleanup old indices failed")
    }

    return nil
}

func (idx *ListingIndexer) cleanupOldIndices(ctx context.Context) error {
    indices, err := idx.client.ListIndices(ctx, idx.config.Index+"_v*")
    if err != nil {
        return err
    }

    // Сортируем по дате создания, удаляем все кроме последних 2
    sort.Slice(indices, func(i, j int) bool {
        return indices[i].CreatedAt.After(indices[j].CreatedAt)
    })

    for i := 2; i < len(indices); i++ {
        if err := idx.client.DeleteIndex(ctx, indices[i].Name); err != nil {
            idx.logger.Warn().Err(err).Str("index", indices[i].Name).Msg("delete old index failed")
        } else {
            idx.logger.Info().Str("index", indices[i].Name).Msg("old index deleted")
        }
    }

    return nil
}
```

### 5.2 Verification

```go
func (idx *ListingIndexer) VerifyIndex(ctx context.Context) (*VerificationResult, error) {
    result := &VerificationResult{}

    // 1. Количество документов в индексе
    indexCount, err := idx.client.CountDocuments(ctx, idx.config.Index)
    if err != nil {
        return nil, err
    }
    result.IndexCount = indexCount

    // 2. Количество активных листингов в БД
    dbCount, err := idx.repo.CountActiveListings(ctx)
    if err != nil {
        return nil, err
    }
    result.DBCount = dbCount

    // 3. Сравнение
    result.Diff = dbCount - indexCount
    result.DiffPercent = float64(result.Diff) / float64(dbCount) * 100

    // 4. Проверка случайных документов
    sampleIDs, err := idx.repo.GetRandomListingIDs(ctx, 10)
    if err != nil {
        return nil, err
    }

    for _, id := range sampleIDs {
        exists, err := idx.client.DocumentExists(ctx, idx.config.Index, id)
        if err != nil || !exists {
            result.MissingIDs = append(result.MissingIDs, id)
        }
    }

    // 5. Проверка маппингов
    mappings, err := idx.client.GetMappings(ctx, idx.config.Index)
    if err != nil {
        return nil, err
    }

    // Проверяем наличие критических полей
    requiredFields := []string{
        "variants", "brand", "condition", "rating",
        "storefront_name", "old_price", "discount_percent",
        "title.trigram", "variant_colors", "variant_sizes",
    }
    for _, field := range requiredFields {
        if !hasField(mappings, field) {
            result.MissingFields = append(result.MissingFields, field)
        }
    }

    result.IsHealthy = result.Diff < 100 &&
                       len(result.MissingIDs) == 0 &&
                       len(result.MissingFields) == 0

    return result, nil
}
```

---

## ФАЗА 6: Мониторинг

**Время:** ~2 часа
**Цель:** Настроить алерты и health checks

### 6.1 Health checks endpoint

```go
// В HTTP handler добавить
func (h *Handler) SearchHealthCheck(c *fiber.Ctx) error {
    ctx := c.Context()

    health := &SearchHealth{
        Status:    "ok",
        Timestamp: time.Now(),
    }

    // 1. OpenSearch connection
    if err := h.searchClient.Ping(ctx); err != nil {
        health.Status = "error"
        health.Errors = append(health.Errors, "opensearch connection failed")
    }

    // 2. Index exists
    exists, err := h.searchClient.IndexExists(ctx)
    if err != nil || !exists {
        health.Status = "error"
        health.Errors = append(health.Errors, "index not found")
    }

    // 3. Documents count
    count, err := h.searchClient.CountDocuments(ctx)
    if err != nil {
        health.Status = "warning"
        health.Warnings = append(health.Warnings, "cannot count documents")
    } else {
        health.DocumentCount = count
        if count == 0 {
            health.Status = "warning"
            health.Warnings = append(health.Warnings, "index is empty")
        }
    }

    // 4. Indexing queue status
    queueStats, err := h.queueRepo.GetStats(ctx)
    if err == nil {
        health.QueuePending = queueStats.Pending
        health.QueueFailed = queueStats.Failed

        if queueStats.Pending > 10000 {
            health.Status = "warning"
            health.Warnings = append(health.Warnings, "queue backlog > 10000")
        }
        if queueStats.Failed > 100 {
            health.Status = "warning"
            health.Warnings = append(health.Warnings, "too many failed jobs")
        }
    }

    // 5. Last index time
    lastIndexed, err := h.queueRepo.GetLastCompletedTime(ctx)
    if err == nil {
        health.LastIndexedAt = lastIndexed
        if time.Since(lastIndexed) > 10*time.Minute {
            health.Status = "warning"
            health.Warnings = append(health.Warnings, "no indexing for 10+ minutes")
        }
    }

    return c.JSON(health)
}
```

### 6.2 Prometheus metrics

```go
var (
    indexingQueuePending = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "vondi_indexing_queue_pending",
        Help: "Number of pending indexing jobs",
    })

    indexingQueueFailed = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "vondi_indexing_queue_failed",
        Help: "Number of failed indexing jobs",
    })

    indexingQueueCompleted = promauto.NewCounter(prometheus.CounterOpts{
        Name: "vondi_indexing_queue_completed_total",
        Help: "Total number of completed indexing jobs",
    })

    indexingDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "vondi_indexing_duration_seconds",
        Help:    "Duration of indexing operations",
        Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30},
    })

    searchDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "vondi_search_duration_seconds",
        Help:    "Duration of search queries",
        Buckets: []float64{0.01, 0.05, 0.1, 0.25, 0.5, 1},
    })

    searchResultsCount = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "vondi_search_results_count",
        Help:    "Number of search results returned",
        Buckets: []float64{0, 1, 10, 50, 100, 500, 1000},
    })

    // ★ НОВЫЕ метрики
    searchZeroResults = promauto.NewCounter(prometheus.CounterOpts{
        Name: "vondi_search_zero_results_total",
        Help: "Total searches with zero results",
    })

    searchFuzzyUsed = promauto.NewCounter(prometheus.CounterOpts{
        Name: "vondi_search_fuzzy_used_total",
        Help: "Total searches where fuzzy matching was used",
    })

    searchFallbackUsed = promauto.NewCounter(prometheus.CounterOpts{
        Name: "vondi_search_fallback_used_total",
        Help: "Total searches where fallback was used",
    })
)
```

---

## ★ ФАЗА 7: Аналитика поиска (НОВАЯ)

**Время:** ~4 часа
**Цель:** Понять что ищут пользователи и улучшать качество поиска

### 7.1 Search Analytics индекс

```go
var SearchAnalyticsSettings = map[string]interface{}{
    "settings": map[string]interface{}{
        "number_of_shards":   1,
        "number_of_replicas": 0,
        "index": map[string]interface{}{
            "refresh_interval": "30s",
        },
    },
    "mappings": map[string]interface{}{
        "properties": map[string]interface{}{
            // Поисковый запрос
            "query":            map[string]interface{}{"type": "keyword"},
            "query_normalized": map[string]interface{}{"type": "keyword"},
            "query_length":     map[string]interface{}{"type": "integer"},

            // Результаты
            "results_count":    map[string]interface{}{"type": "integer"},
            "response_time_ms": map[string]interface{}{"type": "integer"},
            "search_type":      map[string]interface{}{"type": "keyword"}, // exact, fuzzy, fallback

            // Клик
            "clicked_id":       map[string]interface{}{"type": "long"},
            "clicked_position": map[string]interface{}{"type": "integer"},
            "clicked_at":       map[string]interface{}{"type": "date"},

            // Контекст
            "category_id":      map[string]interface{}{"type": "keyword"},
            "filters_used":     map[string]interface{}{"type": "keyword"},
            "sort_used":        map[string]interface{}{"type": "keyword"},

            // Пользователь
            "user_id":          map[string]interface{}{"type": "long"},
            "session_id":       map[string]interface{}{"type": "keyword"},
            "locale":           map[string]interface{}{"type": "keyword"},

            // Время
            "timestamp":        map[string]interface{}{"type": "date"},
            "hour_of_day":      map[string]interface{}{"type": "integer"},
            "day_of_week":      map[string]interface{}{"type": "integer"},
        },
    },
}
```

### 7.2 Логирование событий поиска

```go
type SearchEvent struct {
    Query           string    `json:"query"`
    QueryNormalized string    `json:"query_normalized"`
    QueryLength     int       `json:"query_length"`
    ResultsCount    int       `json:"results_count"`
    ResponseTimeMs  int       `json:"response_time_ms"`
    SearchType      string    `json:"search_type"`
    CategoryID      string    `json:"category_id,omitempty"`
    FiltersUsed     []string  `json:"filters_used,omitempty"`
    SortUsed        string    `json:"sort_used,omitempty"`
    UserID          int64     `json:"user_id,omitempty"`
    SessionID       string    `json:"session_id"`
    Locale          string    `json:"locale"`
    Timestamp       time.Time `json:"timestamp"`
    HourOfDay       int       `json:"hour_of_day"`
    DayOfWeek       int       `json:"day_of_week"`
}

type ClickEvent struct {
    Query           string    `json:"query"`
    ClickedID       int64     `json:"clicked_id"`
    ClickedPosition int       `json:"clicked_position"`
    UserID          int64     `json:"user_id,omitempty"`
    SessionID       string    `json:"session_id"`
    ClickedAt       time.Time `json:"clicked_at"`
}

func (s *SearchService) LogSearchEvent(ctx context.Context, event SearchEvent) error {
    event.Timestamp = time.Now()
    event.HourOfDay = event.Timestamp.Hour()
    event.DayOfWeek = int(event.Timestamp.Weekday())
    event.QueryNormalized = normalizeQuery(event.Query)
    event.QueryLength = len(event.Query)

    return s.analyticsClient.Index(ctx, "search_analytics", event)
}

func (s *SearchService) LogClickEvent(ctx context.Context, event ClickEvent) error {
    event.ClickedAt = time.Now()

    return s.analyticsClient.Index(ctx, "search_analytics", event)
}
```

### 7.3 Отчёты Zero Results

```go
func (s *SearchService) GetZeroResultsQueries(ctx context.Context, days int, limit int) ([]ZeroResultQuery, error) {
    query := map[string]interface{}{
        "size": 0,
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{
                    {
                        "term": map[string]interface{}{"results_count": 0},
                    },
                    {
                        "range": map[string]interface{}{
                            "timestamp": map[string]interface{}{
                                "gte": fmt.Sprintf("now-%dd", days),
                            },
                        },
                    },
                },
            },
        },
        "aggs": map[string]interface{}{
            "zero_result_queries": map[string]interface{}{
                "terms": map[string]interface{}{
                    "field": "query_normalized",
                    "size":  limit,
                    "order": map[string]interface{}{
                        "_count": "desc",
                    },
                },
            },
        },
    }

    // Execute and parse...
    return results, nil
}
```

### 7.4 Trending Searches (реальные данные)

```go
func (s *SearchService) GetTrendingSearches(ctx context.Context, hours int, limit int) ([]TrendingSearch, error) {
    query := map[string]interface{}{
        "size": 0,
        "query": map[string]interface{}{
            "bool": map[string]interface{}{
                "must": []map[string]interface{}{
                    {
                        "range": map[string]interface{}{
                            "timestamp": map[string]interface{}{
                                "gte": fmt.Sprintf("now-%dh", hours),
                            },
                        },
                    },
                    {
                        "range": map[string]interface{}{
                            "results_count": map[string]interface{}{
                                "gt": 0, // Только успешные поиски
                            },
                        },
                    },
                },
            },
        },
        "aggs": map[string]interface{}{
            "trending": map[string]interface{}{
                "terms": map[string]interface{}{
                    "field": "query_normalized",
                    "size":  limit,
                    "order": map[string]interface{}{
                        "_count": "desc",
                    },
                },
                "aggs": map[string]interface{}{
                    "avg_results": map[string]interface{}{
                        "avg": map[string]interface{}{
                            "field": "results_count",
                        },
                    },
                    "click_through": map[string]interface{}{
                        "filter": map[string]interface{}{
                            "exists": map[string]interface{}{
                                "field": "clicked_id",
                            },
                        },
                    },
                },
            },
        },
    }

    // Execute and parse...
    return results, nil
}
```

---

## Чеклист выполнения

### ФАЗА 1: Зачистка рудиментов
- [ ] Удалить `reindex_unified.py`
- [ ] Удалить все `reindex_*.py` из monolith
- [ ] Удалить методы индексации из `storage/opensearch/client.go`
- [ ] Удалить `storage/opensearch/mappings.go` из monolith
- [ ] Удалить старые индексы в OpenSearch
- [ ] Почистить конфигурацию

### ФАЗА 2: Новые маппинги
- [ ] Добавить char_filter для транслитерации сербского
- [ ] Добавить языковые анализаторы (sr, ru, en)
- [ ] ★ Добавить synonym filters для всех языков
- [ ] ★ Добавить universal_analyzer с синонимами
- [ ] Добавить autocomplete analyzer
- [ ] ★ Добавить trigram analyzer для Did You Mean
- [ ] Добавить все новые поля
- [ ] Добавить nested структуру variants
- [ ] Добавить денормализованные поля вариантов
- [ ] Реализовать RecreateIndex()

### ФАЗА 3: Индексатор
- [ ] Добавить loadVariantsForListing()
- [ ] Обновить buildListingDocument() с вариантами
- [ ] Добавить индексацию переводов локации
- [ ] Добавить индексацию новых полей
- [ ] Реализовать batch processing с лимитом 500
- [ ] Добавить loadStorefrontForListing()

### ФАЗА 4: Поисковые запросы
- [ ] ★ Реализовать normalizeQuery() с транслитерацией
- [ ] ★ Добавить fuzziness в multi_match
- [ ] Реализовать multi-language search
- [ ] ★ Реализовать BuildRankedSearchQuery() с function_score
- [ ] ★ Реализовать BuildDidYouMeanQuery()
- [ ] ★ Добавить highlighting
- [ ] ★ Реализовать SearchWithFallback()
- [ ] ★ Реализовать BuildMoreLikeThisQuery()
- [ ] Реализовать variant nested query
- [ ] Добавить все новые фильтры
- [ ] ★ Добавить расширенные фасеты
- [ ] Реализовать Autocomplete()
- [ ] Добавить sorting options

### ФАЗА 5: Переиндексация
- [ ] Реализовать SafeReindexAll() с blue-green
- [ ] Реализовать SwitchAlias()
- [ ] Реализовать cleanupOldIndices()
- [ ] Реализовать VerifyIndex()
- [ ] Выполнить полную переиндексацию

### ФАЗА 6: Мониторинг
- [ ] Добавить /health/search endpoint
- [ ] Добавить Prometheus metrics
- [ ] ★ Добавить метрики zero results и fallback
- [ ] Настроить queue monitoring
- [ ] Добавить алерты

### ★ ФАЗА 7: Аналитика поиска
- [ ] Создать индекс search_analytics
- [ ] Реализовать LogSearchEvent()
- [ ] Реализовать LogClickEvent()
- [ ] Реализовать GetZeroResultsQueries()
- [ ] Реализовать GetTrendingSearches()
- [ ] Добавить dashboard для анализа

---

## Критерии успеха

| Критерий | Целевое значение |
|----------|------------------|
| Время поиска (p95) | < 100ms |
| Время индексации одного документа | < 50ms |
| Размер очереди индексации | < 100 (норма) |
| Точность мультиязычного поиска | > 95% |
| ★ Fuzzy matching accuracy | > 90% |
| ★ Zero results rate | < 5% |
| ★ Click-through rate | > 30% |
| Покрытие полей | 100% |
| Покрытие вариантов | 100% |
| Uptime после релиза | > 99.9% |

---

## Риски и митигация

| Риск | Вероятность | Митигация |
|------|-------------|-----------|
| OOM при переиндексации | Средняя | Batch size 500, мониторинг памяти |
| Рассинхронизация после миграции | Низкая | Verification после каждой фазы |
| Потеря данных при blue-green | Низкая | Сохраняем 2 последних индекса |
| Неполные данные в варинатах | Средняя | Логирование ошибок, retry |
| ★ Синонимы замедляют поиск | Низкая | Тестирование производительности |
| ★ Did You Mean даёт плохие suggestions | Средняя | Настройка confidence и min_word_length |

---

**Автор:** Claude Code
**Дата создания:** 2025-12-19
**Версия:** 2.0 (с улучшениями)
**Статус:** Готов к выполнению
