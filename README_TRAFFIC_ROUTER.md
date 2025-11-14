# Orders Microservice Traffic Router

## Обзор

Traffic Router для Orders микросервиса обеспечивает постепенный rollout с поддержкой:
- **Percentage-based routing** - процентное распределение трафика
- **Canary users** - тестирование на выбранных пользователях
- **Graceful fallback** - автоматический откат на монолит при ошибках

## Архитектура

```
┌─────────────────────────────────────────────────────────────┐
│                     Request                                  │
│  (с user_id из JWT token)                                   │
└────────────────────┬───────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────────┐
│             OrdersTrafficRouter                              │
│                                                              │
│  1. Check: USE_ORDERS_MICROSERVICE enabled?                 │
│  2. Check: User in canary list?                             │
│  3. Check: Random % < ORDERS_ROLLOUT_PERCENT?               │
└────────────┬───────────────────┬───────────────────────────┘
             │                   │
         YES │                   │ NO
             │                   │
             ▼                   ▼
   ┌─────────────────┐   ┌──────────────┐
   │ Orders Service  │   │   Monolith   │
   │  (gRPC:50052)   │   │  (Postgres)  │
   └─────────────────┘   └──────────────┘
```

## Конфигурация

### Environment Variables

```bash
# 1. Основной feature flag (обязательный)
USE_ORDERS_MICROSERVICE=true|false
# Default: false

# 2. Процент трафика на микросервис (0-100)
ORDERS_ROLLOUT_PERCENT=0-100
# Default: 0
# Examples:
#   0   - Весь трафик на монолит
#   25  - 25% на микросервис, 75% на монолит
#   50  - 50/50 распределение
#   100 - Весь трафик на микросервис

# 3. Canary users (comma-separated user IDs)
ORDERS_CANARY_USER_IDS="1,2,3,100"
# Default: ""
# Эти пользователи ВСЕГДА используют микросервис

# 4. gRPC endpoint
ORDERS_GRPC_URL=localhost:50052
# Default: localhost:50052

# 5. Timeout для gRPC запросов
ORDERS_GRPC_TIMEOUT=5s
# Default: 5s

# 6. Fallback при ошибках
ORDERS_FALLBACK_TO_MONOLITH=true|false
# Default: true
```

## Сценарии Rollout

### Сценарий 1: Начальный rollout (0% + canary users)

**Цель:** Протестировать микросервис на ограниченной группе пользователей

```bash
# .env
USE_ORDERS_MICROSERVICE=true
ORDERS_ROLLOUT_PERCENT=0
ORDERS_CANARY_USER_IDS="1,2,3,100,999"
ORDERS_FALLBACK_TO_MONOLITH=true
```

**Результат:**
- ✅ Canary users (1, 2, 3, 100, 999) → microservice
- ❌ Все остальные → monolith
- 🔄 Ошибки microservice → fallback to monolith

**Кто использует:**
- Dev/QA команда для тестирования
- Ранние adopters
- Внутренние пользователи

---

### Сценарий 2: Постепенный rollout (25%)

**Цель:** Расширить тестирование на 25% пользователей

```bash
# .env
USE_ORDERS_MICROSERVICE=true
ORDERS_ROLLOUT_PERCENT=25
ORDERS_CANARY_USER_IDS="1,2,3"
ORDERS_FALLBACK_TO_MONOLITH=true
```

**Результат:**
- ✅ Canary users (1, 2, 3) → microservice (всегда)
- 🎲 25% остальных пользователей → microservice
- ❌ 75% остальных пользователей → monolith

**Когда использовать:**
- После успешного canary testing
- Микросервис стабилен в production
- Нет критичных багов

---

### Сценарий 3: Массовый rollout (50%-75%)

**Цель:** Основная миграция на микросервис

```bash
# .env
USE_ORDERS_MICROSERVICE=true
ORDERS_ROLLOUT_PERCENT=50  # или 75
ORDERS_CANARY_USER_IDS=""
ORDERS_FALLBACK_TO_MONOLITH=true
```

**Результат:**
- 🎲 50% (или 75%) → microservice
- ❌ 50% (или 25%) → monolith

**Когда использовать:**
- Микросервис доказал стабильность на 25%
- Metrics показывают хорошую производительность
- Готовы к полной миграции

---

### Сценарий 4: Полный rollout (100%)

**Цель:** Завершение миграции

```bash
# .env
USE_ORDERS_MICROSERVICE=true
ORDERS_ROLLOUT_PERCENT=100
ORDERS_CANARY_USER_IDS=""
ORDERS_FALLBACK_TO_MONOLITH=true  # Оставляем для безопасности
```

**Результат:**
- ✅ Все пользователи → microservice
- 🔄 Fallback to monolith при ошибках (для безопасности)

**Когда использовать:**
- После успешного 75% rollout
- Metrics показывают лучшую производительность
- Готовы отключить монолит

---

### Сценарий 5: Emergency rollback

**Цель:** Быстрый откат на монолит

```bash
# .env
USE_ORDERS_MICROSERVICE=false
# Или:
ORDERS_ROLLOUT_PERCENT=0
```

**Результат:**
- ❌ ВСЕ пользователи → monolith
- Даже canary users игнорируются

**Когда использовать:**
- Критичные баги в микросервисе
- Performance проблемы
- Необходим hotfix в монолите

---

## Мониторинг и Логирование

### Логи Routing Decisions

Traffic Router логирует каждое решение о routing:

```json
{
  "level": "info",
  "component": "orders_traffic_router",
  "user_id": 123,
  "reason": "canary_user",
  "message": "Routing to microservice: canary user"
}

{
  "level": "debug",
  "component": "orders_traffic_router",
  "user_id": 456,
  "rollout_percent": 50,
  "random_value": 32,
  "reason": "percentage_match",
  "message": "Routing to microservice: percentage-based"
}

{
  "level": "debug",
  "component": "orders_traffic_router",
  "user_id": 789,
  "reason": "zero_percent",
  "message": "Routing to monolith: 0% rollout"
}
```

### HTTP Response Headers

Каждый ответ содержит header `X-Served-By`:

```bash
# Если запрос обработан микросервисом
X-Served-By: microservice

# Если запрос обработан монолитом
X-Served-By: monolith
```

**Пример проверки:**
```bash
TOKEN=$(cat /tmp/token)
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/orders/cart \
  -v 2>&1 | grep "X-Served-By"
```

---

## Тестирование Traffic Router

### Unit Tests

```bash
cd /p/github.com/sveturs/svetu/backend/internal/proj/orders
go test -v -run TestOrdersTrafficRouter
```

**Покрытие тестами:**
- ✅ Zero percent routing (0% → monolith)
- ✅ Full rollout (100% → microservice)
- ✅ Canary users (всегда → microservice)
- ✅ Percentage-based routing (10%, 25%, 50%, 75%)
- ✅ Microservice disabled (всегда → monolith)
- ✅ Canary user ID parsing
- ✅ Edge cases

### Интеграционные тесты

```bash
# 1. Запустить Orders microservice
/home/dim/.local/bin/start-listings-microservice.sh

# 2. Настроить canary user
export ORDERS_CANARY_USER_IDS="1"
export ORDERS_ROLLOUT_PERCENT=0
export USE_ORDERS_MICROSERVICE=true

# 3. Перезапустить монолит
/home/dim/.local/bin/kill-port-3000.sh
screen -dmS backend-3000 bash -c 'cd /p/github.com/sveturs/svetu/backend && go run ./cmd/api/main.go'

# 4. Тестировать routing
TOKEN=$(cat /tmp/token)  # Admin user (ID=1, canary)

# Должен вернуть X-Served-By: microservice
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/orders/cart \
  -v 2>&1 | grep "X-Served-By"
```

---

## Best Practices

### 1. Постепенный Rollout

**НЕ делай:**
```bash
# ❌ 0% → 100% сразу (слишком рискованно)
ORDERS_ROLLOUT_PERCENT=100
```

**Делай:**
```bash
# ✅ Постепенно: 0% → 10% → 25% → 50% → 75% → 100%
ORDERS_ROLLOUT_PERCENT=10   # Day 1
ORDERS_ROLLOUT_PERCENT=25   # Day 3
ORDERS_ROLLOUT_PERCENT=50   # Day 7
ORDERS_ROLLOUT_PERCENT=100  # Day 14
```

### 2. Мониторинг Metrics

**Что отслеживать:**
- Request latency (p50, p95, p99)
- Error rate
- Fallback rate (как часто используется monolith fallback)
- Traffic distribution (microservice vs monolith)

### 3. Canary Users

**Кого выбирать:**
- Dev/QA team
- Early adopters
- Внутренние пользователи
- Power users (с высокой активностью)

### 4. Fallback Strategy

**Всегда включай fallback:**
```bash
# ✅ Оставляй fallback включенным даже при 100%
ORDERS_FALLBACK_TO_MONOLITH=true
```

**Отключай fallback только когда:**
- Монолит полностью выведен из эксплуатации
- Микросервис доказал 100% reliability
- Нет планов возвращаться к монолиту

---

## Troubleshooting

### Проблема: Все запросы идут на монолит (даже canary users)

**Причина:** `USE_ORDERS_MICROSERVICE=false`

**Решение:**
```bash
export USE_ORDERS_MICROSERVICE=true
# Перезапустить backend
```

---

### Проблема: Canary users не используют микросервис

**Проверка:**
```bash
# 1. Проверить конфигурацию
echo $ORDERS_CANARY_USER_IDS

# 2. Проверить user_id в JWT token
TOKEN=$(cat /tmp/token)
curl -s -H "Authorization: Bearer $TOKEN" \
  http://localhost:3000/api/v1/auth/me | jq '.id'

# 3. Убедиться что user_id в списке canary users
```

**Решение:**
```bash
# Добавить user_id в canary list
export ORDERS_CANARY_USER_IDS="1,2,3,YOUR_USER_ID"
```

---

### Проблема: Процентное распределение не работает

**Причина:** Возможно `ORDERS_ROLLOUT_PERCENT=0`

**Проверка:**
```bash
echo $ORDERS_ROLLOUT_PERCENT
```

**Решение:**
```bash
export ORDERS_ROLLOUT_PERCENT=25  # Или любое значение 1-100
```

---

## Миграция с Legacy Feature Flag

### До (старый подход):

```bash
# Все или ничего
USE_ORDERS_MICROSERVICE=true  # 100% microservice
USE_ORDERS_MICROSERVICE=false # 100% monolith
```

### После (Traffic Router):

```bash
# Постепенный rollout
USE_ORDERS_MICROSERVICE=true
ORDERS_ROLLOUT_PERCENT=25      # 25% microservice
ORDERS_CANARY_USER_IDS="1,2,3"  # + canary users
```

**Обратная совместимость:**
- Если `ORDERS_ROLLOUT_PERCENT` не указан → используется legacy behavior (all-or-nothing)
- Traffic Router имеет приоритет над legacy flag

---

## Примеры конфигураций для разных сред

### Development (локальная разработка)

```bash
# .env.development
USE_ORDERS_MICROSERVICE=true
ORDERS_ROLLOUT_PERCENT=100
ORDERS_CANARY_USER_IDS=""
ORDERS_GRPC_URL=localhost:50052
ORDERS_FALLBACK_TO_MONOLITH=true
```

### Staging (тестирование)

```bash
# .env.staging
USE_ORDERS_MICROSERVICE=true
ORDERS_ROLLOUT_PERCENT=50
ORDERS_CANARY_USER_IDS="1,2,3,100"
ORDERS_GRPC_URL=orders-staging.internal:50052
ORDERS_FALLBACK_TO_MONOLITH=true
```

### Production (prod rollout)

```bash
# .env.production - Phase 1 (canary)
USE_ORDERS_MICROSERVICE=true
ORDERS_ROLLOUT_PERCENT=0
ORDERS_CANARY_USER_IDS="1,2,3,10,20"
ORDERS_GRPC_URL=orders-prod.internal:50052
ORDERS_FALLBACK_TO_MONOLITH=true

# .env.production - Phase 2 (10%)
ORDERS_ROLLOUT_PERCENT=10

# .env.production - Phase 3 (25%)
ORDERS_ROLLOUT_PERCENT=25

# .env.production - Phase 4 (50%)
ORDERS_ROLLOUT_PERCENT=50

# .env.production - Phase 5 (100%)
ORDERS_ROLLOUT_PERCENT=100
```

---

## Дополнительные ресурсы

- **Integration Tests:** `/p/github.com/sveturs/svetu/backend/tests/integration/traffic_router_integration_test.go`
- **Unit Tests:** `/p/github.com/sveturs/svetu/backend/internal/proj/orders/traffic_router_test.go`
- **Config Reference:** `/p/github.com/sveturs/svetu/backend/internal/config/config.go` (OrdersConfig)
- **Traffic Router Implementation:** `/p/github.com/sveturs/svetu/backend/internal/proj/orders/traffic_router.go`

---

**Автор:** Orders Microservice Team
**Последнее обновление:** 2025-01-14
**Версия:** 1.0.0
