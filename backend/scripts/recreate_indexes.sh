#!/bin/bash

# Создание индекса marketplace с правильной структурой
echo "Creating marketplace index..."
curl -X PUT "http://localhost:9200/marketplace" \
  -H "Content-Type: application/json" \
  -d '{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
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
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {"type": "integer"},
      "user_id": {"type": "integer"},
      "category_id": {"type": "integer"},
      "category_slug": {"type": "keyword"},
      "status": {"type": "keyword"},
      "title": {
        "type": "text",
        "analyzer": "default_analyzer",
        "fields": {
          "keyword": {"type": "keyword"},
          "serbian": {"type": "text", "analyzer": "serbian_analyzer"},
          "russian": {"type": "text", "analyzer": "russian_analyzer"},
          "english": {"type": "text", "analyzer": "english_analyzer"}
        }
      },
      "description": {
        "type": "text",
        "analyzer": "default_analyzer",
        "fields": {
          "serbian": {"type": "text", "analyzer": "serbian_analyzer"},
          "russian": {"type": "text", "analyzer": "russian_analyzer"},
          "english": {"type": "text", "analyzer": "english_analyzer"}
        }
      },
      "price": {"type": "float"},
      "currency": {"type": "keyword"},
      "condition": {"type": "keyword"},
      "location": {"type": "geo_point"},
      "city": {"type": "keyword"},
      "country": {"type": "keyword"},
      "created_at": {"type": "date"},
      "updated_at": {"type": "date"},
      "view_count": {"type": "integer"},
      "is_promoted": {"type": "boolean"},
      "has_delivery": {"type": "boolean"},
      "seller_type": {"type": "keyword"},
      "storefront_id": {"type": "integer"},
      "images": {
        "type": "nested",
        "properties": {
          "id": {"type": "integer"},
          "file_path": {"type": "keyword"},
          "display_order": {"type": "integer"}
        }
      },
      "attributes": {
        "type": "nested",
        "properties": {
          "attribute_id": {"type": "integer"},
          "attribute_name": {"type": "keyword"},
          "value": {"type": "text"},
          "value_id": {"type": "integer"}
        }
      },
      "db_translations": {
        "type": "nested",
        "properties": {
          "field_name": {"type": "keyword"},
          "language": {"type": "keyword"},
          "translation": {"type": "text"}
        }
      }
    }
  }
}'

echo ""
echo "Creating storefront_products index..."
curl -X PUT "http://localhost:9200/storefront_products" \
  -H "Content-Type: application/json" \
  -d '{
  "settings": {
    "number_of_shards": 1,
    "number_of_replicas": 0,
    "analysis": {
      "analyzer": {
        "default_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase"],
          "char_filter": ["html_strip"]
        }
      }
    }
  },
  "mappings": {
    "properties": {
      "id": {"type": "integer"},
      "storefront_id": {"type": "integer"},
      "name": {
        "type": "text",
        "analyzer": "default_analyzer",
        "fields": {
          "keyword": {"type": "keyword"}
        }
      },
      "description": {
        "type": "text",
        "analyzer": "default_analyzer"
      },
      "base_price": {"type": "float"},
      "currency": {"type": "keyword"},
      "category_id": {"type": "integer"},
      "is_active": {"type": "boolean"},
      "has_variants": {"type": "boolean"},
      "stock_quantity": {"type": "integer"},
      "created_at": {"type": "date"},
      "updated_at": {"type": "date"},
      "view_count": {"type": "integer"},
      "variants": {
        "type": "nested",
        "properties": {
          "id": {"type": "integer"},
          "name": {"type": "text"},
          "sku": {"type": "keyword"},
          "price": {"type": "float"},
          "stock_quantity": {"type": "integer"},
          "is_active": {"type": "boolean"},
          "attributes": {
            "type": "nested",
            "properties": {
              "id": {"type": "integer"},
              "name": {"type": "keyword"},
              "value": {"type": "text"}
            }
          }
        }
      },
      "images": {
        "type": "nested",
        "properties": {
          "id": {"type": "integer"},
          "url": {"type": "keyword"},
          "display_order": {"type": "integer"}
        }
      },
      "attributes": {
        "type": "nested",
        "properties": {
          "id": {"type": "integer"},
          "name": {"type": "keyword"},
          "value": {"type": "text"}
        }
      },
      "popularity_score": {"type": "float"},
      "quality_score": {"type": "float"}
    }
  }
}'

echo ""
echo "Indexes created successfully!"