#!/bin/bash

# Скрипт для принудительной переиндексации объявлений в OpenSearch

echo "Запуск переиндексации объявлений..."

# Запрос к административному API для переиндексации
curl -X POST "http://localhost:8080/api/v1/admin/reindex-listings" \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer YOUR_TOKEN_HERE"

echo "Запрос на переиндексацию отправлен."
echo "Проверка статуса индекса..."

# Проверяем статус индекса
sleep 2
curl -s "http://localhost:9200/_cat/indices" | grep marketplace_listings

echo "Готово. Проверьте логи для подтверждения успешной переиндексации."