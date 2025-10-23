# 📦 PostExpress → Delivery Microservice Migration Plan

> **Стратегия:** Вариант 1 - Постепенная миграция (Gradual Migration)
> **Создан:** 2025-10-23
> **Статус:** 🚧 В разработке
> **Срок:** 5-7 рабочих дней

---

## 🎯 Цель миграции

Перевести все тестовые эндпоинты PostExpress с **прямого WSP API** на **delivery gRPC микросервис**, обеспечив:

- ✅ Централизованную логику доставки
- ✅ Унификацию работы со всеми провайдерами
- ✅ Упрощение backend кода
- ✅ Обратную совместимость в переходный период

---

## 📊 Текущее состояние

### PostExpress модуль (7,763 строк):

```
backend/internal/proj/postexpress/
├── client.go           (316 lines) - WSP API клиент
├── handler/            (2,413 lines)
│   ├── handler.go      - основные эндпоинты
│   └── test_handler.go - тестовые эндпоинты ⚠️
├── service/            (2,074 lines)
├── storage/            (1,545 lines)
└── models/             (394 lines)
```

### Тестовые эндпоинты (требуют миграции):

| Endpoint | Метод | Описание | WSP TX |
|----------|-------|----------|--------|
| `/api/v1/postexpress/test/shipment` | POST | Создать тестовое отправление | TX 3 |
| `/api/v1/postexpress/test/tracking/:number` | GET | Отследить отправление | TX 15 |
| `/api/v1/postexpress/test/cancel/:id` | POST | Отменить отправление | TX 73 |
| `/api/v1/postexpress/test/settlements` | GET | Получить список населенных пунктов | TX 4 |
| `/api/v1/postexpress/test/streets/:settlement` | GET | Получить улицы по городу | TX 6 |
| `/api/v1/postexpress/test/calculate` | POST | Рассчитать стоимость | TX 9 |
| `/api/v1/postexpress/test/validate-address` | POST | Валидировать адрес | TX 11 |
| `/api/v1/postexpress/test/parcel-lockers` | GET | Список паккетоматов | TX 20 |
| `/api/v1/postexpress/test/delivery-services` | GET | Услуги доставки | TX 25 |

### Frontend интеграция:

**Страница:** `frontend/svetu/src/app/[locale]/examples/postexpress-api/page.tsx`

**Проблема:** Использует старые эндпоинты напрямую:
```typescript
const response = await apiClient.post('/postexpress/test/shipment', data);
```

---

## 🚀 План миграции (4 фазы)

### Фаза 1: Создание новых эндпоинтов (2 дня) ⏱️

**Задачи:**

1. **Добавить тестовые эндпоинты в delivery handler** (backend/internal/proj/delivery/handler/test_handler.go):
   - Создать новый файл `test_handler.go`
   - Реализовать 9 тестовых методов через gRPC клиент
   - Зарегистрировать роуты в `module.go`

2. **Маппинг эндпоинтов:**

| Старый endpoint | Новый endpoint | gRPC метод |
|----------------|----------------|------------|
| `POST /postexpress/test/shipment` | `POST /delivery/test/shipment` | `CreateShipment` |
| `GET /postexpress/test/tracking/:number` | `GET /delivery/test/tracking/:number` | `TrackShipment` |
| `POST /postexpress/test/cancel/:id` | `POST /delivery/test/cancel/:id` | `CancelShipment` |
| `GET /postexpress/test/settlements` | `GET /delivery/test/settlements` | `GetSettlements` |
| `GET /postexpress/test/streets/:settlement` | `GET /delivery/test/streets/:settlement` | `GetStreets` |
| `POST /postexpress/test/calculate` | `POST /delivery/test/calculate` | `CalculateRate` |
| `POST /postexpress/test/validate-address` | `POST /delivery/test/validate-address` | `ValidateAddress` |
| `GET /postexpress/test/parcel-lockers` | `GET /delivery/test/parcel-lockers` | `GetParcelLockers` |
| `GET /postexpress/test/delivery-services` | `GET /delivery/test/delivery-services` | `GetDeliveryServices` |

**Примеры реализации:**

```go
// backend/internal/proj/delivery/handler/test_handler.go
func (h *Handler) CreateTestShipment(c *fiber.Ctx) error {
    var req TestShipmentRequest
    if err := c.BodyParser(&req); err != nil {
        return utils.SendErrorResponse(c, fiber.StatusBadRequest, "delivery.invalid_request", nil)
    }

    // Преобразуем в gRPC запрос
    grpcReq := &pb.CreateShipmentRequest{
        Provider: pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS,
        FromAddress: &pb.Address{
            ContactName:  req.SenderName,
            ContactPhone: req.SenderPhone,
            Street:       req.SenderAddress,
            City:         req.SenderCity,
            PostalCode:   req.SenderZip,
            Country:      "RS",
        },
        ToAddress: &pb.Address{
            ContactName:  req.RecipientName,
            ContactPhone: req.RecipientPhone,
            Street:       req.RecipientAddress,
            City:         req.RecipientCity,
            PostalCode:   req.RecipientZip,
            Country:      "RS",
        },
        Packages: []*pb.Package{{
            Weight:       float32(req.Weight) / 1000.0, // граммы → кг
            Value:        float32(req.InsuredValue),
            Description:  req.Content,
        }},
        CodAmount: float32(req.CODAmount),
    }

    // Вызываем микросервис
    resp, err := h.service.grpcClient.CreateShipment(c.Context(), grpcReq)
    if err != nil {
        return utils.SendErrorResponse(c, fiber.StatusInternalServerError, "delivery.shipment_creation_failed", fiber.Map{
            "error": err.Error(),
        })
    }

    return utils.SendSuccessResponse(c, fiber.Map{
        "success":         true,
        "tracking_number": resp.TrackingNumber,
        "shipment_id":     resp.ShipmentId,
        "cost":            resp.Cost,
    })
}
```

**Чеклист:**
- [ ] Создать `/backend/internal/proj/delivery/handler/test_handler.go`
- [ ] Реализовать 9 тестовых методов
- [ ] Добавить роуты в `module.go` (группа `/api/v1/delivery/test/*`)
- [ ] Протестировать через curl/Postman

---

### Фаза 2: DEPRECATED маркеры (1 день) ⏱️

**Задачи:**

1. **Пометить старые эндпоинты как DEPRECATED:**

```go
// backend/internal/proj/postexpress/handler/test_handler.go

func (h *Handler) CreateTestShipment(c *fiber.Ctx) error {
    // DEPRECATED: Используйте /api/v1/delivery/test/shipment
    h.logger.Warn("DEPRECATED endpoint called",
        "endpoint", "/api/v1/postexpress/test/shipment",
        "new_endpoint", "/api/v1/delivery/test/shipment",
    )

    // Можно вернуть HTTP 410 Gone или проксировать на новый эндпоинт
    return utils.SendErrorResponse(c, fiber.StatusGone, "postexpress.endpoint_deprecated", fiber.Map{
        "message": "This endpoint is deprecated. Use /api/v1/delivery/test/shipment instead",
        "new_endpoint": "/api/v1/delivery/test/shipment",
        "sunset_date": "2025-11-23", // через месяц
    })
}
```

2. **Добавить warning в логи:**
   - Все вызовы старых эндпоинтов логируются как DEPRECATED
   - Мониторинг использования (для определения, когда можно удалить)

**Чеклист:**
- [ ] Обновить все test handlers с DEPRECATED маркерами
- [ ] Добавить HTTP 410 Gone ответы
- [ ] Настроить логирование deprecated вызовов

---

### Фаза 3: Миграция Frontend (1-2 дня) ⏱️

**Задачи:**

1. **Обновить страницу `/examples/postexpress-api`:**

```typescript
// frontend/svetu/src/app/[locale]/examples/postexpress-api/page.tsx

// СТАРЫЙ КОД (удалить):
const response = await apiClient.post('/postexpress/test/shipment', {
  recipient_name: 'John Doe',
  // ...
});

// НОВЫЙ КОД:
const response = await apiClient.post('/delivery/test/shipment', {
  recipient_name: 'John Doe',
  // ...
});
```

2. **Обновить типы запросов/ответов:**
   - Адаптировать под новый формат delivery API
   - Обновить error handling

3. **Переименовать страницу (опционально):**
   - `/examples/postexpress-api` → `/examples/delivery-api`
   - Обновить навигацию и переводы

**Чеклист:**
- [ ] Обновить все вызовы API в page.tsx
- [ ] Адаптировать типы и интерфейсы
- [ ] Обновить переводы (en/ru/sr)
- [ ] Протестировать в браузере

---

### Фаза 4: Удаление legacy кода (1-2 дня) ⏱️

**⚠️ Только после того, как:**
- ✅ Все frontend мигрировал на новые эндпоинты
- ✅ Логи показывают 0 вызовов deprecated эндпоинтов (за неделю)
- ✅ Production тесты прошли успешно

**Задачи:**

1. **Удалить PostExpress тестовые эндпоинты:**
   - Удалить `backend/internal/proj/postexpress/handler/test_handler.go` (1,171 строк)
   - Удалить роуты из `handler.go`

2. **Опционально: Удалить весь PostExpress модуль** (если он больше не используется):
   - Проверить, есть ли другие места, использующие `postexpress` модуль
   - Если нет - удалить всю директорию `backend/internal/proj/postexpress/` (7,763 строк)

3. **Обновить документацию:**
   - Удалить упоминания PostExpress тестовых эндпоинтов
   - Обновить API документацию (swagger)

**Чеклист:**
- [ ] Проверить логи на отсутствие deprecated вызовов
- [ ] Удалить test_handler.go
- [ ] Удалить тестовые роуты из handler.go
- [ ] (Опционально) Удалить весь postexpress модуль
- [ ] Обновить swagger.json

---

## 🧪 Тестирование

### Unit тесты:

```bash
# Backend тесты для новых delivery test endpoints
cd /data/hostel-booking-system/backend
go test ./internal/proj/delivery/handler -v -run TestCreateTestShipment
go test ./internal/proj/delivery/handler -v -run TestTrackTestShipment
```

### Функциональные тесты:

```bash
# Получить токен
TOKEN="$(cat /tmp/token)"

# Тест 1: Создание отправления через новый эндпоинт
curl -X POST -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "recipient_name": "Test User",
    "recipient_phone": "0641234567",
    "recipient_city": "Beograd",
    "recipient_address": "Takovska 2",
    "recipient_zip": "11000",
    "sender_name": "Sve Tu d.o.o.",
    "sender_phone": "0641234567",
    "sender_city": "Beograd",
    "sender_address": "Bulevar kralja Aleksandra 73",
    "sender_zip": "11000",
    "weight": 500,
    "content": "Test paket"
  }' \
  http://localhost:3000/api/v1/delivery/test/shipment | jq '.'

# Тест 2: Получение населенных пунктов
curl -H "Authorization: Bearer ${TOKEN}" \
  http://localhost:3000/api/v1/delivery/test/settlements | jq '.'

# Тест 3: Расчет стоимости
curl -X POST -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "from_city": "Beograd",
    "to_city": "Novi Sad",
    "weight": 1000
  }' \
  http://localhost:3000/api/v1/delivery/test/calculate | jq '.'
```

### Frontend тесты:

```bash
# Запустить frontend
cd /data/hostel-booking-system/frontend/svetu
yarn dev -p 3001

# Открыть в браузере:
# http://localhost:3001/ru/examples/postexpress-api

# Проверить:
# 1. Страница загружается без ошибок
# 2. Можно создать тестовое отправление
# 3. Отслеживание работает
# 4. Все поля корректно отображаются
```

---

## 🚨 Риски и митигация

| Риск | Вероятность | Влияние | Митигация |
|------|-------------|---------|-----------|
| gRPC микросервис недоступен | Средняя | Высокое | Добавить fallback на старые эндпоинты (временно) |
| Несовпадение форматов данных | Высокая | Среднее | Детальный маппинг в Phase 1, unit тесты |
| Frontend ломается после миграции | Средняя | Среднее | Постепенное развертывание, feature flags |
| Production проблемы Post Express | Низкая | Высокое | Mock режим для тестирования, staging окружение |

---

## 📈 Метрики успеха

**После завершения миграции:**

- ✅ **Код:** -7,763 строк (PostExpress модуль удален)
- ✅ **Эндпоинты:** 9 новых delivery test endpoints работают
- ✅ **Frontend:** Страница `/examples/delivery-api` работает через gRPC
- ✅ **Тесты:** 100% coverage для новых endpoints
- ✅ **Production:** 0 ошибок, 0 deprecated вызовов

---

## 📚 Связанная документация

- [Delivery Microservice Complete](../DELIVERY_MICROSERVICE_MIGRATION_COMPLETE.md)
- [Delivery Quick Start](../DELIVERY_QUICK_START.md)
- [Delivery Module README](../../backend/internal/proj/delivery/README.md)
- [Proto Schema](../../backend/proto/delivery/v1/delivery.proto)

---

## 🎯 Следующие шаги

1. ✅ Прочитать этот план
2. 🚧 Запустить Phase 1 (создание новых эндпоинтов)
3. ⏸️ Запустить Phase 2 (DEPRECATED маркеры)
4. ⏸️ Запустить Phase 3 (миграция frontend)
5. ⏸️ Запустить Phase 4 (удаление legacy кода)

---

**Версия:** 1.0 | **Дата:** 2025-10-23 | **Статус:** 🚧 В разработке | **Срок:** 5-7 дней
