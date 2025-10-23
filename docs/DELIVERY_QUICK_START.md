# 🚀 Delivery System - Quick Start Guide

**Для разработчиков, которым нужно быстро начать работать с новой системой доставки**

---

## ⚡ За 5 минут

### Что изменилось?

**ДО:** Backend содержал всю логику доставки + провайдеров
**ПОСЛЕ:** Backend → gRPC вызов → Delivery Microservice → Провайдеры

### Архитектура одной строкой

```
Frontend → Backend (localhost:3000) → gRPC → Delivery Service (svetu.rs:30051) → 5 провайдеров
```

### Что нужно знать?

1. **Backend больше НЕ содержит логику провайдеров**
2. **Все операции доставки идут через gRPC**
3. **Deprecated эндпоинты возвращают HTTP 501**
4. **Локальная БД backend используется только для кеширования**

---

## 🔌 Как использовать?

### Вариант 1: HTTP API (рекомендуется для Frontend)

```bash
# 1. Получить JWT токен
TOKEN=$(cat /tmp/token)

# 2. Получить список провайдеров
curl -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/delivery/providers | jq .

# 3. Создать отправление
curl -X POST -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "provider_id": 1,
    "provider_code": "post_express",
    "order_id": 123,
    "from_address": {
      "name": "John Doe",
      "phone": "+381611234567",
      "street": "Kneza Milosa 10",
      "city": "Belgrade",
      "postal_code": "11000",
      "country": "RS"
    },
    "to_address": {
      "name": "Jane Smith",
      "phone": "+381621234567",
      "street": "Bulevar Oslobodjenja 1",
      "city": "Novi Sad",
      "postal_code": "21000",
      "country": "RS"
    },
    "packages": [{
      "weight": 1.5,
      "dimensions": {"length": 30, "width": 20, "height": 10},
      "value": 5000,
      "description": "Electronics"
    }]
  }' http://localhost:3000/api/v1/delivery/shipments | jq .

# 4. Отследить отправление
curl -H "Authorization: Bearer $TOKEN" \
  "http://localhost:3000/api/v1/delivery/shipments/track/post_express-1761215005-6768" | jq .
```

### Вариант 2: gRPC (для backend разработчиков)

```go
import (
    "backend/internal/proj/delivery/grpcclient"
    pb "backend/pkg/grpc/delivery/v1"
)

// Создать клиент
client, err := grpcclient.NewClient("svetu.rs:30051", logger)
if err != nil {
    return err
}
defer client.Close()

// Создать отправление
resp, err := client.CreateShipment(ctx, &pb.CreateShipmentRequest{
    Provider: pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS,
    FromAddress: &pb.Address{
        ContactName:  "John Doe",
        ContactPhone: "+381611234567",
        Street:       "Kneza Milosa 10",
        City:         "Belgrade",
        PostalCode:   "11000",
        Country:      "RS",
    },
    // ... остальные поля
})
```

---

## 🛠️ Настройка окружения

### Backend (.env)

```bash
# Единственная новая переменная
DELIVERY_GRPC_URL=svetu.rs:30051

# Если не задано, используется svetu.rs:30051 по умолчанию
```

### Проверка подключения

```bash
# 1. Запустить backend
cd /data/hostel-booking-system/backend
go run ./cmd/api/main.go

# 2. Проверить логи (должны увидеть):
# ✅ "Successfully connected to delivery gRPC service" url=svetu.rs:30051

# 3. Проверить gRPC (опционально)
grpcurl -plaintext svetu.rs:30051 list
# Должно вывести: delivery.v1.DeliveryService
```

---

## 📋 Доступные эндпоинты

### ✅ Работают (требуют JWT)

| Метод | Endpoint | Описание |
|-------|----------|----------|
| GET | `/api/v1/delivery/providers` | Список провайдеров |
| POST | `/api/v1/delivery/shipments` | Создать отправление |
| GET | `/api/v1/delivery/shipments/:id` | Получить отправление |
| GET | `/api/v1/delivery/shipments/track/:tracking` | Отследить отправление |
| DELETE | `/api/v1/delivery/shipments/:id` | Отменить отправление |
| GET | `/api/v1/products/:id/delivery-attributes` | Атрибуты доставки товара |
| PUT | `/api/v1/products/:id/delivery-attributes` | Обновить атрибуты |

### ❌ НЕ работают (deprecated)

| Метод | Endpoint | Возвращает |
|-------|----------|------------|
| POST | `/api/v1/delivery/calculate-universal` | HTTP 501 |
| POST | `/api/v1/delivery/calculate-cart` | HTTP 501 |

**Миграция:** Используйте gRPC метод `CalculateRate` вместо этих эндпоинтов.

---

## 🐛 Частые проблемы

### Backend не подключается к микросервису

```
ERROR: Failed to connect to delivery gRPC service at svetu.rs:30051
```

**Решение:**
```bash
# Проверить, что микросервис работает
ssh svetu@svetu.rs 'docker ps | grep delivery-service'

# Проверить порт
ssh svetu@svetu.rs 'netstat -tlnp | grep 30051'

# Проверить файрвол (если на production)
ssh svetu@svetu.rs 'sudo ufw status | grep 30051'
```

### Эндпоинт возвращает 501

**Это нормально!** Эндпоинты `calculate-universal` и `calculate-cart` больше не поддерживаются.

**Используйте вместо них:**
- gRPC метод `CalculateRate`
- Или создайте новый HTTP эндпоинт в backend

### Провайдер в mock режиме

```
WARN: Post Express provider not available
INFO: Using mock provider for post_express
```

**Это ожидаемо** - Post Express требует production credentials.

**Для тестирования:** Используйте mock провайдер (работает как настоящий, но не создает реальные отправления)

---

## 📊 Мониторинг

### Быстрая проверка здоровья

```bash
# Backend
curl http://localhost:3000/
# Ожидаем: Svetu API 0.2.4

# Microservice (gRPC)
grpcurl -plaintext svetu.rs:30051 list
# Ожидаем: delivery.v1.DeliveryService

# Microservice (Docker)
ssh svetu@svetu.rs 'docker ps | grep delivery-service'
# Ожидаем: Up X hours
```

### Логи

```bash
# Backend
tail -f /tmp/backend.log | grep delivery

# Microservice
ssh svetu@svetu.rs 'docker logs -f delivery-service'

# Ошибки микросервиса
ssh svetu@svetu.rs 'docker logs delivery-service 2>&1 | grep ERROR'
```

---

## 🎯 Следующие шаги

1. ✅ Прочитал Quick Start
2. ✅ Настроил окружение
3. ✅ Протестировал основные эндпоинты
4. 📚 Читаю [полную документацию](DELIVERY_MICROSERVICE_MIGRATION_COMPLETE.md)
5. 🔍 Изучаю [proto схему](../backend/proto/delivery/v1/delivery.proto)
6. 💻 Смотрю [примеры кода](../backend/internal/proj/delivery/)

---

## 📚 Дополнительно

- **Полная документация:** [DELIVERY_MICROSERVICE_MIGRATION_COMPLETE.md](DELIVERY_MICROSERVICE_MIGRATION_COMPLETE.md)
- **Proto схема:** `backend/proto/delivery/v1/delivery.proto`
- **Код backend модуля:** `backend/internal/proj/delivery/`
- **Microservice repo:** `github.com/sveturs/delivery`

---

## 🆘 Нужна помощь?

1. Проверь [Troubleshooting раздел](DELIVERY_MICROSERVICE_MIGRATION_COMPLETE.md#troubleshooting) в полной документации
2. Посмотри логи (см. раздел Мониторинг выше)
3. Спроси в команде backend
4. Создай issue в репозитории

---

**Обновлено:** 2025-10-23 | **Версия:** 1.0
