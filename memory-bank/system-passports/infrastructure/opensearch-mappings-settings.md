# 📋 Паспорт OpenSearch Mappings и настройки

## 🏷️ Метаданные
- **Назначение:** Схемы данных и настройки анализаторов для OpenSearch индексов
- **Тип компонента:** Инфраструктура / Search Schema
- **Статус:** Активный, используется в production
- **Версия OpenSearch:** 2.x
- **Файлы:** `backend/internal/storage/opensearch/mappings.go`

## 🎯 Назначение
OpenSearch mappings и настройки определяют структуру данных, типы полей, анализаторы текста и правила индексирования для всех поисковых индексов системы Sve Tu Platform.

## 🗂️ Структура индексов

### 1. Marketplace Listings Index (`marketplace`)

#### 📊 Настройки индекса (Settings)
```json
{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 1,
    "analysis": {
      "analyzer": {
        "serbian_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "serbian_latin_stemmer"]
        },
        "russian_analyzer": {
          "type": "custom", 
          "tokenizer": "standard",
          "filter": ["lowercase", "russian_stemmer"]
        },
        "english_analyzer": {
          "type": "custom",
          "tokenizer": "standard", 
          "filter": ["lowercase", "english_stemmer"]
        },
        "default_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase"],
          "char_filter": ["html_strip"]
        },
        "autocomplete": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "autocomplete_filter"]
        },
        "shingle_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "shingle_filter"]
        }
      },
      "filter": {
        "serbian_latin_stemmer": {
          "type": "stemmer",
          "language": "serbian"
        },
        "russian_stemmer": {
          "type": "stemmer", 
          "language": "russian"
        },
        "english_stemmer": {
          "type": "stemmer",
          "language": "english"
        },
        "autocomplete_filter": {
          "type": "edge_ngram",
          "min_gram": 1,
          "max_gram": 20
        },
        "shingle_filter": {
          "type": "shingle",
          "min_shingle_size": 2,
          "max_shingle_size": 3
        }
      }
    }
  }
}\n```\n\n#### 🗺️ Схема полей (Mappings)\n\n##### Основные поля объявления\n```json\n{\n  \"properties\": {\n    \"id\": {\"type\": \"integer\"},\n    \"title\": {\n      \"type\": \"text\",\n      \"analyzer\": \"default_analyzer\",\n      \"fields\": {\n        \"keyword\": {\"type\": \"keyword\"},\n        \"serbian\": {\"type\": \"text\", \"analyzer\": \"serbian_analyzer\"},\n        \"russian\": {\"type\": \"text\", \"analyzer\": \"russian_analyzer\"},\n        \"english\": {\"type\": \"text\", \"analyzer\": \"english_analyzer\"},\n        \"autocomplete\": {\"type\": \"text\", \"analyzer\": \"autocomplete\"}\n      }\n    },\n    \"description\": {\n      \"type\": \"text\",\n      \"analyzer\": \"default_analyzer\",\n      \"fields\": {\n        \"serbian\": {\"type\": \"text\", \"analyzer\": \"serbian_analyzer\"},\n        \"russian\": {\"type\": \"text\", \"analyzer\": \"russian_analyzer\"},\n        \"english\": {\"type\": \"text\", \"analyzer\": \"english_analyzer\"}\n      }\n    },\n    \"price\": {\"type\": \"double\"},\n    \"old_price\": {\"type\": \"double\"},\n    \"has_discount\": {\"type\": \"boolean\"},\n    \"status\": {\"type\": \"keyword\"},\n    \"condition\": {\n      \"type\": \"text\",\n      \"fields\": {\"keyword\": {\"type\": \"keyword\"}}\n    }\n  }\n}\n```\n\n##### Географические поля\n```json\n{\n  \"location\": {\n    \"type\": \"text\",\n    \"fields\": {\"keyword\": {\"type\": \"keyword\"}}\n  },\n  \"coordinates\": {\"type\": \"geo_point\"},\n  \"city\": {\n    \"type\": \"text\",\n    \"fields\": {\"keyword\": {\"type\": \"keyword\"}}\n  },\n  \"country\": {\n    \"type\": \"text\",\n    \"fields\": {\"keyword\": {\"type\": \"keyword\"}}\n  }\n}\n```\n\n##### Мультиязычные переводы (nested)\n```json\n{\n  \"translations\": {\n    \"type\": \"object\",\n    \"properties\": {\n      \"sr\": {\n        \"properties\": {\n          \"title\": {\"type\": \"text\", \"analyzer\": \"serbian_analyzer\"},\n          \"description\": {\"type\": \"text\", \"analyzer\": \"serbian_analyzer\"}\n        }\n      },\n      \"ru\": {\n        \"properties\": {\n          \"title\": {\"type\": \"text\", \"analyzer\": \"russian_analyzer\"},\n          \"description\": {\"type\": \"text\", \"analyzer\": \"russian_analyzer\"}\n        }\n      },\n      \"en\": {\n        \"properties\": {\n          \"title\": {\"type\": \"text\", \"analyzer\": \"english_analyzer\"},\n          \"description\": {\"type\": \"text\", \"analyzer\": \"english_analyzer\"}\n        }\n      }\n    }\n  }\n}\n```\n\n##### Атрибуты категорий (nested)\n```json\n{\n  \"attributes\": {\n    \"type\": \"nested\",\n    \"properties\": {\n      \"attribute_id\": {\"type\": \"integer\"},\n      \"attribute_name\": {\"type\": \"keyword\"},\n      \"display_name\": {\"type\": \"text\"},\n      \"attribute_type\": {\"type\": \"keyword\"},\n      \"text_value\": {\n        \"type\": \"text\",\n        \"analyzer\": \"default_analyzer\",\n        \"fields\": {\n          \"keyword\": {\"type\": \"keyword\"},\n          \"lowercase\": {\"type\": \"keyword\", \"normalizer\": \"lowercase\"},\n          \"serbian\": {\"type\": \"text\", \"analyzer\": \"serbian_analyzer\"},\n          \"russian\": {\"type\": \"text\", \"analyzer\": \"russian_analyzer\"},\n          \"english\": {\"type\": \"text\", \"analyzer\": \"english_analyzer\"}\n        }\n      },\n      \"numeric_value\": {\"type\": \"double\"},\n      \"boolean_value\": {\"type\": \"boolean\"},\n      \"json_value\": {\"type\": \"text\"},\n      \"display_value\": {\"type\": \"text\"},\n      \"translations\": {\n        \"type\": \"object\",\n        \"properties\": {\n          \"en\": {\"type\": \"text\"},\n          \"sr\": {\"type\": \"text\"},\n          \"ru\": {\"type\": \"text\"}\n        }\n      }\n    }\n  }\n}\n```\n\n##### Изображения (nested)\n```json\n{\n  \"images\": {\n    \"type\": \"nested\",\n    \"properties\": {\n      \"id\": {\"type\": \"integer\"},\n      \"file_path\": {\"type\": \"keyword\"},\n      \"is_main\": {\"type\": \"boolean\"},\n      \"alt_text\": {\"type\": \"text\"}\n    }\n  }\n}\n```\n\n##### Поля автодополнения\n```json\n{\n  \"title_suggest\": {\n    \"type\": \"completion\",\n    \"analyzer\": \"default_analyzer\",\n    \"search_analyzer\": \"default_analyzer\",\n    \"contexts\": [\n      {\n        \"name\": \"category\",\n        \"type\": \"category\"\n      }\n    ]\n  },\n  \"all_attributes_text\": {\n    \"type\": \"text\",\n    \"analyzer\": \"default_analyzer\"\n  }\n}\n```\n\n### 2. Storefront Products Index (`storefront_products`)\n\n#### 📊 Настройки для товаров витрин\n```json\n{\n  \"settings\": {\n    \"analysis\": {\n      \"analyzer\": {\n        \"russian_analyzer\": {\n          \"tokenizer\": \"standard\",\n          \"filter\": [\"lowercase\", \"russian_stop\", \"russian_stemmer\"]\n        }\n      }\n    }\n  }\n}\n```\n\n#### 🗺️ Схема товаров\n```json\n{\n  \"properties\": {\n    \"product_id\": {\"type\": \"integer\"},\n    \"storefront_id\": {\"type\": \"integer\"},\n    \"category_id\": {\"type\": \"integer\"},\n    \"name\": {\n      \"type\": \"search_as_you_type\",\n      \"analyzer\": \"russian_analyzer\"\n    },\n    \"description\": {\n      \"type\": \"text\",\n      \"analyzer\": \"russian_analyzer\"\n    },\n    \"price\": {\"type\": \"float\"},\n    \"price_min\": {\"type\": \"float\"},\n    \"price_max\": {\"type\": \"float\"},\n    \"sku\": {\"type\": \"keyword\"},\n    \"barcode\": {\"type\": \"keyword\"},\n    \"brand\": {\n      \"type\": \"text\",\n      \"fields\": {\n        \"keyword\": {\"type\": \"keyword\"},\n        \"lowercase\": {\"type\": \"keyword\", \"normalizer\": \"lowercase\"}\n      }\n    },\n    \"model\": {\n      \"type\": \"text\",\n      \"fields\": {\n        \"keyword\": {\"type\": \"keyword\"},\n        \"lowercase\": {\"type\": \"keyword\", \"normalizer\": \"lowercase\"}\n      }\n    }\n  }\n}\n```\n\n##### Инвентаризация\n```json\n{\n  \"inventory\": {\n    \"properties\": {\n      \"track\": {\"type\": \"boolean\"},\n      \"count\": {\"type\": \"integer\"},\n      \"reserved\": {\"type\": \"integer\"},\n      \"available\": {\"type\": \"integer\"},\n      \"in_stock\": {\"type\": \"boolean\"},\n      \"low_stock\": {\"type\": \"boolean\"}\n    }\n  }\n}\n```\n\n##### Варианты товаров\n```json\n{\n  \"variants\": {\n    \"properties\": {\n      \"id\": {\"type\": \"integer\"},\n      \"name\": {\"type\": \"text\"},\n      \"sku\": {\"type\": \"keyword\"},\n      \"price\": {\"type\": \"float\"},\n      \"attributes\": {\"type\": \"object\"},\n      \"inventory\": {\"type\": \"object\"}\n    }\n  }\n}\n```\n\n### 3. Storefronts Index (`storefronts`)\n\n#### 🗺️ Схема витрин\n```json\n{\n  \"properties\": {\n    \"id\": {\"type\": \"integer\"},\n    \"user_id\": {\"type\": \"integer\"},\n    \"slug\": {\"type\": \"keyword\"},\n    \"name\": {\"type\": \"text\"},\n    \"description\": {\"type\": \"text\"},\n    \"address\": {\"type\": \"text\"},\n    \"city\": {\"type\": \"keyword\"},\n    \"postal_code\": {\"type\": \"keyword\"},\n    \"country\": {\"type\": \"keyword\"},\n    \"location\": {\"type\": \"geo_point\"},\n    \"phone\": {\"type\": \"keyword\"},\n    \"email\": {\"type\": \"keyword\"},\n    \"website\": {\"type\": \"keyword\"},\n    \"rating\": {\"type\": \"float\"},\n    \"reviews_count\": {\"type\": \"integer\"},\n    \"products_count\": {\"type\": \"integer\"},\n    \"sales_count\": {\"type\": \"integer\"},\n    \"views_count\": {\"type\": \"integer\"},\n    \"subscription_plan\": {\"type\": \"keyword\"},\n    \"is_active\": {\"type\": \"boolean\"},\n    \"is_verified\": {\"type\": \"boolean\"},\n    \"is_open_now\": {\"type\": \"boolean\"},\n    \"payment_methods\": {\"type\": \"keyword\"},\n    \"delivery_providers\": {\"type\": \"keyword\"}\n  }\n}\n```\n\n## 🔧 Анализаторы и фильтры\n\n### Языковые анализаторы\n\n#### Serbian Analyzer\n- **Tokenizer:** standard\n- **Filters:** lowercase, serbian_latin_stemmer\n- **Назначение:** Обработка сербского текста с латинской транслитерацией\n\n#### Russian Analyzer\n- **Tokenizer:** standard\n- **Filters:** lowercase, russian_stemmer, russian_stop\n- **Назначение:** Обработка русского текста с учетом морфологии\n\n#### English Analyzer\n- **Tokenizer:** standard\n- **Filters:** lowercase, english_stemmer\n- **Назначение:** Обработка английского текста\n\n### Специальные анализаторы\n\n#### Default Analyzer\n- **Tokenizer:** standard\n- **Filters:** lowercase\n- **Char Filters:** html_strip\n- **Назначение:** Универсальный анализатор с очисткой HTML\n\n#### Autocomplete Analyzer\n- **Tokenizer:** standard\n- **Filters:** lowercase, edge_ngram (1-20 символов)\n- **Назначение:** Автодополнение при вводе\n\n#### Shingle Analyzer\n- **Tokenizer:** standard\n- **Filters:** lowercase, shingle (2-3 слова)\n- **Назначение:** Фразовый поиск\n\n## 🎯 Особенности схем\n\n### Мультиязычная поддержка\n1. **Multi-field подход** - каждое текстовое поле имеет версии для разных языков\n2. **Translations объект** - структурированные переводы\n3. **Контекстные анализаторы** - выбор анализатора по языку запроса\n\n### Nested структуры\n1. **Атрибуты** - сложная типизация с переводами\n2. **Изображения** - метаданные и флаги\n3. **Варианты товаров** - структурированные опции\n\n### Автодополнение\n1. **Completion suggester** - быстрые подсказки\n2. **Edge n-gram** - поиск по частичному вводу\n3. **Search-as-you-type** - мгновенный поиск\n\n### Геопоиск\n1. **geo_point поля** - координаты для радиусного поиска\n2. **Keyword адреса** - точные совпадения локаций\n3. **Интеграция с картами** - поддержка картографических сервисов\n\n## ⚡ Производительность\n\n### Оптимизация индексов\n- **Single shard** - для небольших объемов данных\n- **HTML strip** - очистка контента при индексации\n- **Keyword поля** - быстрая фильтрация без анализа\n\n### Нормализаторы\n- **Lowercase normalizer** - приведение к нижнему регистру без токенизации\n- **Поддержка точных совпадений** через keyword поля\n\n## 🔗 Связи с компонентами\n\n### Зависимости\n- **PostgreSQL** - источник данных для индексации\n- **Backend API** - использование схем в запросах\n- **MinIO** - ссылки на изображения в mappings\n\n### Использование\n- **Search Service** - выполнение поисковых запросов\n- **Indexing Service** - создание и обновление документов\n- **Frontend Search** - отображение результатов\n\n---\n**Паспорт создан:** 2025-06-29  \n**Компонент:** OpenSearch Mappings и настройки  \n**Статус:** Активный в production