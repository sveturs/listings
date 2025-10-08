#!/bin/bash
echo "📦 Создание еще 2 простых тестовых отправлений..."
echo ""

# Модифицируем working script для 2 отправлений
for i in 2 3; do
  echo "=== Отправление #$i ==="
  go run test_post_express_working.go 2>&1 | grep -E "(✅|📦|🆔|Tracking)"
  echo ""
  sleep 2
done
