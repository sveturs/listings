#!/bin/bash

# Скрипт для обновления маппинга OpenSearch для поддержки нечеткого поиска

OPENSEARCH_URL="${OPENSEARCH_URL:-http://localhost:9200}"
INDEX_NAME="${INDEX_NAME:-marketplace_listings}"

echo "Updating OpenSearch mapping for index: $INDEX_NAME"

# Закрываем индекс для обновления настроек
echo "Closing index..."
curl -X POST "$OPENSEARCH_URL/$INDEX_NAME/_close" \
  -H "Content-Type: application/json"

# Обновляем настройки анализаторов
echo "Updating analyzer settings..."
curl -X PUT "$OPENSEARCH_URL/$INDEX_NAME/_settings" \
  -H "Content-Type: application/json" \
  -d '{
    "analysis": {
      "filter": {
        "ngram_filter": {
          "type": "ngram",
          "min_gram": 2,
          "max_gram": 3
        },
        "edge_ngram_filter": {
          "type": "edge_ngram",
          "min_gram": 2,
          "max_gram": 20
        }
      },
      "analyzer": {
        "ngram_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "ngram_filter"]
        },
        "edge_ngram_analyzer": {
          "type": "custom",
          "tokenizer": "standard",
          "filter": ["lowercase", "edge_ngram_filter"]
        }
      }
    }
  }'

# Открываем индекс
echo "Opening index..."
curl -X POST "$OPENSEARCH_URL/$INDEX_NAME/_open" \
  -H "Content-Type: application/json"

# Обновляем маппинг для новых полей
echo "Updating mapping..."
curl -X PUT "$OPENSEARCH_URL/$INDEX_NAME/_mapping" \
  -H "Content-Type: application/json" \
  -d '{
    "properties": {
      "title": {
        "type": "text",
        "fields": {
          "ngram": {
            "type": "text",
            "analyzer": "ngram_analyzer",
            "search_analyzer": "standard"
          },
          "edge_ngram": {
            "type": "text",
            "analyzer": "edge_ngram_analyzer",
            "search_analyzer": "standard"
          }
        }
      },
      "description": {
        "type": "text",
        "fields": {
          "ngram": {
            "type": "text",
            "analyzer": "ngram_analyzer",
            "search_analyzer": "standard"
          }
        }
      }
    }
  }'

echo "Mapping update completed!"

# Переиндексируем данные для применения новых анализаторов
echo "Consider running reindex to apply new analyzers to existing data"