#!/bin/bash
echo "📦 Создание еще 2 тестовых отправлений..."

for i in 2 3; do
  echo "=== Отправление #$i ==="
  sed "s/TEST-ORDER-001/TEST-ORDER-00$i/g; s/TEST-REF-001/TEST-REF-00$i/g; s/Test paket za SVETU/Test paket #$i za SVETU/g" test_post_express_working.go > /tmp/test_temp_$i.go
  go run /tmp/test_temp_$i.go 2>&1 | tail -8
  echo ""
  sleep 1
done

rm -f /tmp/test_temp_*.go
