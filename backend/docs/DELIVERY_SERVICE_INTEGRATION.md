# Delivery Service Integration Guide

## 📋 Table of Contents

- [Overview](#overview)
- [Architecture](#architecture)
- [Integration Points](#integration-points)
- [API Endpoints](#api-endpoints)
- [Error Handling](#error-handling)
- [Environment Configuration](#environment-configuration)
- [Troubleshooting](#troubleshooting)

---

## Overview

Svetu Marketplace интегрирован с внешним delivery microservice для управления доставкой заказов. Интеграция построена на gRPC для backend-to-backend коммуникации и REST API через BFF proxy для frontend.

### Key Features

- ✅ **Multi-provider support**: Post Express, BEX Express, AKS Express, D Express, City Express
- ✅ **Rate calculation**: Сравнение стоимости доставки от всех провайдеров
- ✅ **Automatic shipment creation**: При подтверждении заказа автоматически создается shipment
- ✅ **Real-time tracking**: Отслеживание статуса доставки
- ✅ **Graceful degradation**: Система работает даже если microservice недоступен
- ✅ **Circuit breaker**: Защита от каскадных ошибок
- ✅ **Retry logic**: 3 попытки с exponential backoff

---

## Architecture

### System Overview

```
┌─────────────┐      ┌──────────────┐      ┌─────────────────┐      ┌──────────────────┐
│   Browser   │─────▶│  Next.js BFF │─────▶│ Backend (Fiber) │─────▶│ Delivery Service │
│  (React)    │      │  /api/v2/*   │      │   /api/v1/*     │      │   (gRPC :50053)  │
└─────────────┘      └──────────────┘      └─────────────────┘      └──────────────────┘
     │                     │                      │                         │
     │                     │                      │                         │
     └─ httpOnly cookies   └─ JWT Bearer          └─ gRPC Protocol          └─ External APIs
        (access_token)         Authorization         (protobuf)                 (Post Express, etc)
```

### Component Layers

#### 1. Frontend Layer (React/Next.js)

**Компоненты:**
- `UnifiedDeliverySelector`: Универсальный селектор для выбора доставки
- `CartDeliveryCalculator`: Расчет доставки в корзине
- `DeliveryInfo`: Отображение tracking информации

**State Management:**
- `deliverySlice` (Redux Toolkit): Кэширование providers, calculations, tracking

**API Client:**
- `deliveryService`: Wrapper для всех delivery API calls через BFF proxy

**Технологии:**
- React 19
- Next.js 15 (App Router)
- Redux Toolkit
- TypeScript

#### 2. BFF Proxy Layer (Next.js API Routes)

**Файл:** `frontend/svetu/src/app/api/v2/[...path]/route.ts`

**Функции:**
- Маппинг `/api/v2/*` → `/api/v1/*`
- Автоматическое добавление JWT из httpOnly cookies
- Централизованная обработка ошибок
- Логирование запросов

**Преимущества:**
- Нет CORS проблем (все на одном домене)
- JWT недоступны для JavaScript (защита от XSS)
- Простая инфраструктура

#### 3. Backend Layer (Go/Fiber)

**Модули:**

**`internal/proj/delivery/`** - Delivery module
- `handler/`: HTTP handlers для REST API
- `service/`: Бизнес-логика
- `grpcclient/`: gRPC клиент для микросервиса
- `repository/`: Database access (если нужно)

**`internal/proj/orders/service/`** - Orders integration
- `createShipmentForOrder()`: Автоматическое создание shipment при подтверждении заказа
- `enrichOrderWithTracking()`: Обогащение заказа tracking информацией

**gRPC Client:**
- Retry logic: 3 попытки с exponential backoff
- Circuit breaker: Защита от перегрузки
- Timeout: 30 секунд на запрос

**Технологии:**
- Go 1.22+
- Fiber v2 (HTTP framework)
- gRPC
- PostgreSQL

#### 4. Delivery Microservice (gRPC)

**Протокол:** gRPC (Protocol Buffers)

**Методы:**
- `CalculateRate`: Расчет стоимости доставки
- `CreateShipment`: Создание отправления
- `TrackShipment`: Отслеживание статуса
- `CancelShipment`: Отмена отправления
- `GetShipment`: Получение информации о shipment

**Провайдеры:**
- Post Express
- BEX Express
- AKS Express
- D Express
- City Express

---

## Integration Points

### 1. Order Checkout Flow

**Шаги интеграции:**

```
1. User adds items to cart
   └─▶ CartDeliveryCalculator запрашивает расчет стоимости

2. Calculate delivery rates
   ├─▶ Frontend: deliveryService.calculateRate()
   ├─▶ BFF Proxy: POST /api/v2/delivery/calculate-universal
   ├─▶ Backend: POST /api/v1/delivery/calculate-universal
   └─▶ Delivery Service: CalculateRate (gRPC)

3. User selects delivery option
   └─▶ Redux: selectQuote({ storefrontId, quote })

4. User confirms order
   ├─▶ POST /api/v2/orders (BFF proxy)
   ├─▶ Backend: CreateOrder() → status = "pending"
   └─▶ Redis: Save quote in cache

5. Payment successful
   └─▶ Backend: ConfirmOrder()

6. Automatic shipment creation
   ├─▶ Backend: createShipmentForOrder()
   ├─▶ Delivery Service: CreateShipment (gRPC)
   ├─▶ Response: shipment_id, tracking_number
   └─▶ Database: Update order.tracking_number

7. Order status updated
   └─▶ Frontend: Shows tracking link
```

**Код примера (Backend):**

```go
// backend/internal/proj/orders/service/order_service.go

func (s *OrderService) ConfirmOrder(ctx context.Context, orderID int64) error {
    // 1. Validate order
    order, err := s.orderRepo.GetByID(ctx, orderID)
    if err != nil {
        return fmt.Errorf("failed to get order: %w", err)
    }

    // 2. Commit inventory reservations
    if err := s.inventoryMgr.CommitOrderReservations(ctx, orderID); err != nil {
        return fmt.Errorf("failed to commit reservations: %w", err)
    }

    // 3. Update order status
    now := time.Now()
    order.Status = models.OrderStatusConfirmed
    order.ConfirmedAt = &now

    if err := s.orderRepo.Update(ctx, order); err != nil {
        return fmt.Errorf("failed to update order status: %w", err)
    }

    // 4. Create shipment in delivery service (GRACEFUL)
    if err := s.createShipmentForOrder(ctx, order); err != nil {
        s.logger.Error("Failed to create shipment: %v (order_id: %d)", err, orderID)
        // НЕ фейлим весь order - shipment можно создать позже вручную
    }

    return nil
}
```

### 2. Tracking Integration

**Flow:**

```
1. User opens order details page
   └─▶ GET /api/v2/orders/:id

2. Backend enriches order with tracking
   ├─▶ enrichOrderWithTracking()
   ├─▶ Delivery Service: TrackShipment (gRPC)
   └─▶ Response includes tracking events

3. Frontend displays tracking
   └─▶ DeliveryInfo component shows status, events, ETA
```

**Код примера (Backend):**

```go
// backend/internal/proj/orders/service/order_service.go

func (s *OrderService) enrichOrderWithTracking(ctx context.Context, order *models.StorefrontOrder) error {
    if order.TrackingNumber == nil || *order.TrackingNumber == "" {
        return nil // No tracking number yet
    }

    // Call delivery microservice
    resp, err := s.deliveryClient.TrackShipment(ctx, &deliveryv1.TrackShipmentRequest{
        TrackingNumber: *order.TrackingNumber,
    })
    if err != nil {
        s.logger.Error("Failed to track shipment: %v", err)
        return err // Non-critical, just log
    }

    // Parse tracking info
    trackingInfo := map[string]interface{}{
        "status":             resp.Status.String(),
        "current_location":   resp.CurrentLocation,
        "estimated_delivery": resp.EstimatedDelivery.AsTime(),
        "events":             resp.Events,
    }

    // Store as JSON in order.tracking_info field
    trackingJSON, _ := json.Marshal(trackingInfo)
    trackingStr := string(trackingJSON)
    order.TrackingInfo = &trackingStr

    return nil
}
```

### 3. Admin Management

**Функции:**

- **View shipments:** Admin может просматривать все shipments
- **Track shipments:** Просмотр детальной tracking информации
- **Retry creation:** Если shipment не создался - можно создать вручную
- **Cancel shipments:** Отмена shipment при отмене заказа

**Эндпоинты:**

```
GET    /api/v1/admin/delivery/shipments       - Список всех shipments
GET    /api/v1/admin/delivery/shipments/:id   - Детали shipment
POST   /api/v1/admin/delivery/shipments       - Создать shipment вручную
DELETE /api/v1/admin/delivery/shipments/:id   - Отменить shipment
```

---

## API Endpoints

### Calculate Rate

**Endpoint:** `POST /api/v2/delivery/calculate-universal`

**Request:**
```json
{
  "from_location": {
    "city": "Belgrade",
    "postal_code": "11000"
  },
  "to_location": {
    "city": "Novi Sad",
    "postal_code": "21000"
  },
  "items": [
    {
      "weight": 2.5,
      "length": 30,
      "width": 20,
      "height": 10,
      "quantity": 1
    }
  ],
  "provider_id": "post_express",
  "insurance_value": 10000,
  "cod_amount": 5000
}
```

**Response:**
```json
{
  "providers": [
    {
      "provider_id": "post_express",
      "provider_name": "Post Express",
      "base_price": 350.0,
      "insurance": 50.0,
      "cod_fee": 100.0,
      "weight_fee": 0.0,
      "distance_fee": 0.0,
      "total_cost": 500.0,
      "estimated_delivery_days": 2,
      "currency": "RSD"
    },
    {
      "provider_id": "bex_express",
      "provider_name": "BEX Express",
      "total_cost": 450.0,
      "estimated_delivery_days": 1,
      "currency": "RSD"
    }
  ],
  "recommended": {
    "provider_id": "post_express",
    "total_cost": 500.0
  },
  "cheapest": {
    "provider_id": "bex_express",
    "total_cost": 450.0
  },
  "fastest": {
    "provider_id": "bex_express",
    "estimated_delivery_days": 1
  }
}
```

### Create Shipment

**Endpoint:** `POST /api/v2/delivery/shipments`

**Note:** Обычно вызывается автоматически при подтверждении заказа. Для ручного создания (admin):

**Request:**
```json
{
  "order_id": 123,
  "provider_code": "post_express",
  "from_address": {
    "contact_name": "Store Name",
    "contact_phone": "+381601234567",
    "street": "Main Street 1",
    "city": "Belgrade",
    "postal_code": "11000",
    "country": "RS"
  },
  "to_address": {
    "contact_name": "Customer Name",
    "contact_phone": "+381607654321",
    "street": "Customer Street 5",
    "city": "Novi Sad",
    "postal_code": "21000",
    "country": "RS"
  },
  "packages": [
    {
      "weight": 2.5,
      "length": 30,
      "width": 20,
      "height": 10,
      "description": "Order #123"
    }
  ]
}
```

**Response:**
```json
{
  "shipment_id": 456,
  "tracking_number": "PE1234567890RS",
  "provider_code": "post_express",
  "status": "PENDING",
  "label_url": "https://delivery-service.com/labels/PE1234567890RS.pdf"
}
```

### Track Shipment

**Endpoint:** `GET /api/v2/delivery/track/:trackingToken`

**Example:** `GET /api/v2/delivery/track/PE1234567890RS`

**Response:**
```json
{
  "shipment_id": 456,
  "tracking_number": "PE1234567890RS",
  "status": "IN_TRANSIT",
  "current_location": "Postal center Belgrade",
  "estimated_delivery": "2025-10-31T15:00:00Z",
  "events": [
    {
      "timestamp": "2025-10-29T10:00:00Z",
      "location": "Belgrade depot",
      "status": "CONFIRMED",
      "description": "Shipment confirmed and accepted"
    },
    {
      "timestamp": "2025-10-29T14:30:00Z",
      "location": "Postal center Belgrade",
      "status": "IN_TRANSIT",
      "description": "Package in transit to Novi Sad"
    }
  ]
}
```

### Cancel Shipment

**Endpoint:** `DELETE /api/v2/delivery/shipments/:shipmentId`

**Response:**
```json
{
  "success": true,
  "message": "Shipment cancelled successfully"
}
```

---

## Error Handling

### Graceful Degradation Strategy

**Principle:** Система должна работать даже если delivery microservice недоступен.

#### Scenario 1: Microservice недоступен при checkout

**Поведение:**
1. User видит ошибку при расчете доставки
2. Может продолжить checkout БЕЗ доставки
3. Доставку можно добавить позже через admin

**Code:**
```go
// Backend: Graceful handling
if err := s.createShipmentForOrder(ctx, order); err != nil {
    s.logger.Error("Failed to create shipment: %v (order_id: %d)", err, orderID)
    // НЕ возвращаем ошибку - заказ подтверждается без shipment
    // TODO: Добавить в admin UI кнопку "Retry Create Shipment"
}
```

#### Scenario 2: Microservice недоступен при tracking

**Поведение:**
1. Order details показываются
2. Tracking section показывает "Tracking information unavailable"
3. Retry через 1 минуту

**Code:**
```typescript
// Frontend: Fallback UI
if (trackingError) {
  return (
    <div className="alert alert-warning">
      <InformationCircleIcon className="w-5 h-5" />
      <div>
        <div>Tracking information temporarily unavailable</div>
        <button onClick={retry}>Retry</button>
      </div>
    </div>
  );
}
```

### Circuit Breaker

**Implementation:** `backend/internal/proj/delivery/grpcclient/client.go`

**Parameters:**
- **Threshold:** 5 consecutive failures
- **Open duration:** 30 seconds
- **Half-open:** После 30 сек пробует 1 запрос

**Behavior:**
```go
func (c *Client) isCircuitBreakerOpen() bool {
    if c.failureCount < circuitBreakerOpen {
        return false
    }

    // If circuit opened recently, keep it open
    if time.Since(c.lastFailureTime) < 30*time.Second {
        return true
    }

    // Half-open: allow one request to test
    c.failureCount = 0
    return false
}
```

### Retry Logic

**Parameters:**
- **Max retries:** 3
- **Initial backoff:** 100ms
- **Max backoff:** 2 seconds
- **Multiplier:** 2.0 (exponential)

**Retry on:**
- `UNAVAILABLE` - Сервис временно недоступен
- `DEADLINE_EXCEEDED` - Timeout
- Network errors

**Do NOT retry on:**
- `INVALID_ARGUMENT` - Неверные параметры
- `NOT_FOUND` - Shipment не найден
- `PERMISSION_DENIED` - Нет прав
- `ALREADY_EXISTS` - Дубликат

**Code:**
```go
func (c *Client) shouldRetry(err error) bool {
    st, ok := status.FromError(err)
    if !ok {
        return true // Network error - retry
    }

    switch st.Code() {
    case codes.Unavailable, codes.DeadlineExceeded:
        return true
    case codes.InvalidArgument, codes.NotFound, codes.PermissionDenied, codes.AlreadyExists:
        return false
    default:
        return true
    }
}
```

### Error Response Format

**Standard error response:**
```json
{
  "error": {
    "code": "DELIVERY_SERVICE_UNAVAILABLE",
    "message": "Delivery service temporarily unavailable. Order confirmed, shipment will be created later.",
    "details": {
      "order_id": 123,
      "can_retry": true
    }
  }
}
```

**Error codes:**
- `DELIVERY_SERVICE_UNAVAILABLE` - Microservice недоступен
- `INVALID_DELIVERY_ADDRESS` - Неверный адрес
- `PROVIDER_NOT_AVAILABLE` - Провайдер не работает в этом регионе
- `RATE_CALCULATION_FAILED` - Ошибка расчета стоимости
- `SHIPMENT_CREATION_FAILED` - Ошибка создания shipment
- `TRACKING_NOT_FOUND` - Tracking информация не найдена

---

## Environment Configuration

### Backend Environment Variables

**Required:**

```bash
# Delivery microservice gRPC URL
DELIVERY_SERVICE_URL=localhost:50053

# Production:
# DELIVERY_SERVICE_URL=delivery-service.internal:50053
```

**Optional:**

```bash
# Timeout для gRPC запросов (default: 30s)
DELIVERY_GRPC_TIMEOUT=30s

# Enable circuit breaker (default: true)
DELIVERY_CIRCUIT_BREAKER_ENABLED=true

# Max retries (default: 3)
DELIVERY_MAX_RETRIES=3
```

### Frontend Environment Variables

**BFF Proxy:**

```bash
# Backend URL для BFF proxy (server-side)
BACKEND_INTERNAL_URL=http://localhost:3000

# Production:
# BACKEND_INTERNAL_URL=http://backend-internal:3000
```

**Note:** Frontend НЕ обращается напрямую к delivery microservice!

### Docker Compose

```yaml
services:
  backend:
    environment:
      - DELIVERY_SERVICE_URL=delivery-service:50053
    depends_on:
      - delivery-service

  delivery-service:
    image: sveturs/delivery-service:latest
    ports:
      - "50053:50053"
    environment:
      - POST_EXPRESS_API_KEY=${POST_EXPRESS_API_KEY}
      - BEX_EXPRESS_API_KEY=${BEX_EXPRESS_API_KEY}
```

---

## Troubleshooting

### Issue 1: "Delivery service unavailable"

**Symptoms:**
- Checkout fails на этапе расчета доставки
- Logs: `Failed to create shipment: rpc error: code = Unavailable`

**Solutions:**
1. Check delivery service is running:
   ```bash
   curl http://localhost:50053/health  # HTTP health check
   ```

2. Check network connectivity:
   ```bash
   telnet localhost 50053
   ```

3. Check environment variable:
   ```bash
   echo $DELIVERY_SERVICE_URL
   ```

4. Check logs:
   ```bash
   docker logs delivery-service
   ```

### Issue 2: "Circuit breaker is open"

**Symptoms:**
- Все delivery requests fail
- Logs: `Circuit breaker is open, rejecting CreateShipment request`

**Cause:** 5+ consecutive failures

**Solution:**
1. Wait 30 seconds for circuit to half-open
2. Fix underlying issue (service down, network)
3. Restart backend to reset circuit breaker:
   ```bash
   systemctl restart backend
   ```

### Issue 3: Tracking не работает

**Symptoms:**
- Order details не показывают tracking
- Response: `tracking_number: null`

**Причины:**

1. **Shipment не создался:**
   - Check logs: `Failed to create shipment`
   - Solution: Создать shipment вручную через admin

2. **Delivery service не вернул tracking_number:**
   - Check delivery service logs
   - Check provider integration

3. **Tracking_number не сохранился:**
   - Check database: `SELECT tracking_number FROM storefront_orders WHERE id = ?`
   - Check order update logic

**Fix:**
```sql
-- Manually set tracking number if shipment created externally
UPDATE storefront_orders
SET tracking_number = 'PE1234567890RS'
WHERE id = 123;
```

### Issue 4: Duplicate shipments

**Symptoms:**
- Provider shows 2+ shipments for same order
- Extra charges

**Причина:** Retry logic создал дубликат

**Prevention:**
- Используй idempotency key: `order_id`
- Delivery service должен проверять дубликаты

**Fix:**
```bash
# Cancel duplicate shipment
curl -X DELETE http://localhost:3000/api/v1/admin/delivery/shipments/456
```

### Issue 5: Rate calculation очень медленный

**Symptoms:**
- Checkout долго загружается
- Timeout errors

**Причины:**

1. **Multiple sequential calls:**
   - Solution: Используй `/calculate-universal` вместо отдельных calls

2. **No caching:**
   - Solution: deliverySlice кэширует на 5 минут
   - Check Redux DevTools: `state.delivery.calculations`

3. **Slow provider API:**
   - Check delivery service logs для provider response times
   - Consider increasing timeout

**Optimization:**
```typescript
// Frontend: Use cached calculation
const cachedCalc = useAppSelector(selectCalculation(request));

if (cachedCalc) {
  console.log('Using cached calculation');
  return cachedCalc;
}

// Otherwise fetch
dispatch(calculateRate({ request }));
```

---

## Performance Metrics

### Expected Response Times

| Endpoint | Expected | Acceptable | Slow |
|----------|----------|------------|------|
| Calculate Rate | < 500ms | < 1s | > 2s |
| Create Shipment | < 1s | < 3s | > 5s |
| Track Shipment | < 200ms | < 500ms | > 1s |
| Cancel Shipment | < 500ms | < 1s | > 2s |

### Monitoring

**Backend logs:**
```
[INFO] Shipment created successfully: PE1234567890RS (duration: 450ms)
[WARN] Slow delivery service response: 2500ms (endpoint: CreateShipment)
[ERROR] Circuit breaker opened after 5 failures
```

**Frontend Redux DevTools:**
- Check `delivery.calculationsLoading` для active requests
- Check `delivery.calculations` для cache hits/misses
- Monitor TTL (5 минут)

---

## Best Practices

### 1. Always use BFF proxy

```typescript
// ✅ CORRECT
import { deliveryService } from '@/services/delivery';
const response = await deliveryService.calculateRate(request);

// ❌ WRONG
const response = await fetch('http://localhost:3000/api/v1/delivery/...');
```

### 2. Handle errors gracefully

```typescript
// Frontend
const { data, error } = await deliveryService.calculateRate(request);

if (error) {
  // Show fallback UI, allow user to continue
  return <FallbackDeliveryUI />;
}
```

### 3. Use caching

```typescript
// Check cache first
const cached = useAppSelector(selectCalculation(request));
if (cached) return cached;

// Otherwise fetch
dispatch(calculateRate({ request }));
```

### 4. Monitor performance

```go
// Backend
start := time.Now()
resp, err := s.deliveryClient.CreateShipment(ctx, req)
duration := time.Since(start)

if duration > 2*time.Second {
    s.logger.Warn("Slow delivery service response: %v", duration)
}
```

### 5. Log all integration points

```go
s.logger.Info("Creating shipment for order_id=%d, provider=%s", orderID, providerCode)
// ... API call ...
s.logger.Info("Shipment created: tracking_number=%s", trackingNumber)
```

---

## Related Documentation

- [Delivery Microservice API Reference](./DELIVERY_MICROSERVICE_API.md) - gRPC методы и примеры
- [Frontend Delivery Guide](./DELIVERY_FRONTEND_GUIDE.md) - Components, Redux, best practices
- [Orders Integration](./ORDERS_INTEGRATION.md) - Полный checkout flow

---

**Last updated:** 2025-10-29
**Version:** 1.0.0
