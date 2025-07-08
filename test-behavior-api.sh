#!/bin/bash

# Тест behavior tracking API

echo "=== Тестирование Behavior Tracking API ==="
echo

# Тест 1: Отправка batch событий
echo "1. Отправка batch событий:"
curl -X POST http://localhost:3000/api/v1/analytics/track \
  -H "Content-Type: application/json" \
  -d '{
    "events": [
      {
        "event_type": "search_performed",
        "session_id": "test_session_123",
        "search_query": "test query",
        "metadata": {
          "results_count": 10
        }
      },
      {
        "event_type": "result_clicked",
        "session_id": "test_session_123",
        "search_query": "test query",
        "item_id": "123",
        "position": 1,
        "metadata": {
          "click_time": 1234
        }
      }
    ],
    "batch_id": "batch_test_123",
    "created_at": "2025-07-08T15:00:00Z"
  }' | jq

echo
echo "2. Отправка одиночного события:"
curl -X POST http://localhost:3000/api/v1/analytics/track \
  -H "Content-Type: application/json" \
  -d '{
    "event_type": "item_viewed",
    "session_id": "test_session_456",
    "item_id": "456",
    "item_type": "marketplace",
    "metadata": {
      "view_duration": 5000
    }
  }' | jq