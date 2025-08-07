# IT архитектура складской системы Sve Tu с собственной WMS

## 1. Обновленная общая архитектура

```
┌─────────────────────────────────────────────────────────────────────┐
│                        Sve Tu Platform                               │
├─────────────────┬────────────────┬──────────────┬──────────────────┤
│  Frontend       │  Backend API   │  PostgreSQL  │  OpenSearch      │
│  (Next.js)      │  (Go/Fiber)    │  Database    │  Search Engine   │
└────────┬────────┴───────┬────────┴──────┬───────┴──────────────────┘
         │                │               │
         │                ▼               │
     ┌───┴─────────────────────────────────────────────┐
     │         Warehouse Management Service            │
     │              (Go Microservice)                  │
     ├─────────────────────┬───────────────────────────┤
     │   API Gateway       │   Business Logic          │
     │   - REST API        │   - Order Processing      │
     │   - WebSocket       │   - Inventory Management  │
     │   - gRPC            │   - Route Optimization    │
     └─────────┬───────────┴──────────┬────────────────┘
               │                      │
               ▼                      ▼
     ┌─────────────────────┐ ┌──────────────────────┐
     │   Собственная WMS   │ │  External Services   │
     ├─────────────────────┤ ├──────────────────────┤
     │ • Web Dashboard     │ │ • Carrier APIs       │
     │ • Mobile PWA        │ │ • Payment Gateway    │
     │ • Scanner Support   │ │ • SMS/Email Service  │
     │ • Real-time Sync    │ │ • Analytics Platform │
     └─────────────────────┘ └──────────────────────┘
```

## 2. Компоненты собственной WMS

### 2.1 Backend архитектура WMS

```
wms-backend/
├── cmd/
│   ├── api/            # HTTP API сервер
│   ├── worker/         # Background задачи
│   └── migrator/       # Миграции БД
├── internal/
│   ├── api/
│   │   ├── rest/       # REST endpoints
│   │   ├── grpc/       # gRPC сервисы
│   │   └── websocket/  # Real-time обновления
│   ├── core/
│   │   ├── inventory/  # Управление остатками
│   │   ├── orders/     # Обработка заказов
│   │   ├── locations/  # Управление локациями
│   │   └── reports/    # Отчетность
│   ├── services/
│   │   ├── picking/    # Сервис сборки
│   │   ├── packing/    # Сервис упаковки
│   │   ├── shipping/   # Сервис отправки
│   │   └── returns/    # Обработка возвратов
│   └── infrastructure/
│       ├── database/   # PostgreSQL
│       ├── cache/      # Redis
│       ├── queue/      # RabbitMQ
│       └── storage/    # MinIO для документов
```

### 2.2 Frontend архитектура WMS

```
wms-frontend/
├── apps/
│   ├── dashboard/      # Веб-панель управления
│   └── mobile/         # PWA для сборщиков
├── packages/
│   ├── ui/             # Общие UI компоненты
│   ├── api-client/     # API клиент
│   ├── scanner/        # Библиотека сканирования
│   └── types/          # TypeScript типы
└── shared/
    ├── hooks/          # React hooks
    ├── utils/          # Утилиты
    └── constants/      # Константы
```

## 3. Интеграция с основной платформой

### 3.1 API Gateway паттерн

```go
// Единая точка входа для всех складских операций
type WarehouseGateway struct {
    wmsClient      WMSClient
    svetuAPI       SvetuAPIClient
    carrierManager CarrierManager
    eventBus       EventBus
}

// Обработка заказа через gateway
func (gw *WarehouseGateway) ProcessOrder(ctx context.Context, order Order) error {
    // 1. Проверяем доступность в WMS
    availability, err := gw.wmsClient.CheckAvailability(order.Items)
    if err != nil {
        return err
    }
    
    // 2. Резервируем товары
    reservation, err := gw.wmsClient.ReserveItems(order.Items)
    if err != nil {
        return err
    }
    
    // 3. Создаем задание на fulfillment
    fulfillmentOrder := gw.wmsClient.CreateFulfillmentOrder(order, reservation)
    
    // 4. Уведомляем основную платформу
    gw.eventBus.Publish("order.fulfillment.created", fulfillmentOrder)
    
    return nil
}
```

### 3.2 Event-driven архитектура

```yaml
# События между системами
events:
  # От платформы к WMS
  - order.created
  - order.cancelled
  - product.updated
  - storefront.inventory.updated
  
  # От WMS к платформе
  - inventory.level.low
  - order.picked
  - order.packed
  - order.shipped
  - return.received
```

## 4. Технологический стек

### 4.1 Backend технологии

```yaml
Core:
  - Language: Go 1.21+
  - Framework: Fiber v2
  - Database: PostgreSQL 15
  - Cache: Redis 7
  - Queue: RabbitMQ
  - Search: Elasticsearch (для логов)

Monitoring:
  - Metrics: Prometheus + Grafana
  - Tracing: Jaeger
  - Logging: ELK Stack
  - Alerting: AlertManager

DevOps:
  - Containers: Docker
  - Orchestration: Kubernetes
  - CI/CD: GitLab CI
  - IaC: Terraform
```

### 4.2 Frontend технологии

```yaml
Core:
  - Framework: React 18
  - Language: TypeScript 5
  - State: Zustand
  - Routing: React Router 6
  - UI: Tailwind CSS + Radix UI

Mobile PWA:
  - Framework: Next.js 14
  - Offline: Service Workers
  - Camera: WebRTC API
  - Scanner: QuaggaJS

Build Tools:
  - Bundler: Vite
  - Testing: Vitest + React Testing Library
  - Linting: ESLint + Prettier
```

## 5. Безопасность и масштабируемость

### 5.1 Безопасность

```go
// JWT авторизация с ролями
type Claims struct {
    UserID   string   `json:"user_id"`
    Role     string   `json:"role"`
    Scopes   []string `json:"scopes"`
    TenantID string   `json:"tenant_id"`
    jwt.StandardClaims
}

// API Rate Limiting
rateLimiter := limiter.New(
    limiter.WithTrustedForwardHeader(true),
    limiter.WithRateLimit(rate.NewLimiter(rate.Every(time.Second), 100)),
)

// Шифрование чувствительных данных
type EncryptedField struct {
    Value string `json:"value" encrypt:"aes"`
    IV    string `json:"iv"`
}
```

### 5.2 Масштабируемость

```yaml
# Kubernetes deployment
apiVersion: apps/v1
kind: Deployment
metadata:
  name: wms-api
spec:
  replicas: 3
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1
      maxUnavailable: 0
  template:
    spec:
      containers:
      - name: wms-api
        image: svetu/wms-api:latest
        resources:
          requests:
            memory: "256Mi"
            cpu: "250m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
        readinessProbe:
          httpGet:
            path: /ready
            port: 8080
```

## 6. Интеграции с внешними сервисами

### 6.1 Унифицированный API курьеров

```go
// Интерфейс для всех курьерских служб
type CarrierAPI interface {
    GetRates(shipment Shipment) ([]Rate, error)
    CreateShipment(shipment Shipment) (*Label, error)
    TrackShipment(trackingNumber string) (*TrackingInfo, error)
    CancelShipment(shipmentID string) error
    GetPickupSlots(date time.Time) ([]TimeSlot, error)
    SchedulePickup(pickup PickupRequest) (*Confirmation, error)
}

// Реализации для каждого партнера
carriers := map[string]CarrierAPI{
    "posta_srbije": &PostaSrbijeAPI{},
    "post_express": &PostExpressAPI{},
    "dexpress":     &DExpressAPI{},
    "bex":          &BexAPI{},
}

// Автоматический выбор оптимального курьера
func SelectBestCarrier(shipment Shipment) (string, Rate) {
    var bestRate Rate
    var bestCarrier string
    
    for name, api := range carriers {
        rates, _ := api.GetRates(shipment)
        for _, rate := range rates {
            if bestRate.Price == 0 || rate.Price < bestRate.Price {
                bestRate = rate
                bestCarrier = name
            }
        }
    }
    
    return bestCarrier, bestRate
}
```

### 6.2 Hardware интеграции

```typescript
// Поддержка различных типов сканеров
class ScannerManager {
  private scanners: Map<string, Scanner> = new Map();
  
  constructor() {
    // USB сканеры через WebUSB
    this.scanners.set('usb', new USBScanner());
    
    // Bluetooth сканеры
    this.scanners.set('bluetooth', new BluetoothScanner());
    
    // Камера устройства
    this.scanners.set('camera', new CameraScanner());
  }
  
  async scan(preferredType?: string): Promise<string> {
    const scanner = preferredType 
      ? this.scanners.get(preferredType) 
      : this.getAvailableScanner();
      
    return scanner.scan();
  }
}

// Интеграция с принтерами этикеток
class LabelPrinter {
  async print(label: Label) {
    if ('usb' in navigator) {
      // Прямая печать на USB принтер
      const device = await navigator.usb.requestDevice({
        filters: [{ vendorId: 0x0a5f }] // Zebra
      });
      await device.open();
      await device.transferOut(1, label.toZPL());
    } else {
      // Fallback на системную печать
      window.print();
    }
  }
}
```

## 7. Миграция данных и развертывание

### 7.1 План миграции

```sql
-- Новые таблицы для WMS
CREATE SCHEMA IF NOT EXISTS wms;

-- Перенос существующих данных
INSERT INTO wms.products 
SELECT 
    p.id,
    p.sku,
    p.name,
    p.barcode,
    jsonb_build_object(
        'length', p.length,
        'width', p.width,
        'height', p.height,
        'weight', p.weight
    ) as dimensions
FROM products p
WHERE p.storefront_id IN (
    SELECT id FROM storefronts WHERE fulfillment_enabled = true
);

-- Создание начальных локаций
INSERT INTO wms.locations (warehouse_id, code, type, capacity)
SELECT 
    w.id,
    z.code || '-' || r || '-' || s || '-' || b as code,
    'shelf' as type,
    100 as capacity
FROM warehouses w
CROSS JOIN LATERAL (VALUES ('A'), ('B'), ('C')) as z(code)
CROSS JOIN LATERAL generate_series(1, 10) as r
CROSS JOIN LATERAL generate_series(1, 5) as s
CROSS JOIN LATERAL generate_series(1, 4) as b;
```

### 7.2 Развертывание

```bash
#!/bin/bash
# deploy-wms.sh

# 1. Build Docker образы
docker build -t svetu/wms-api:latest ./wms-backend
docker build -t svetu/wms-dashboard:latest ./wms-frontend/apps/dashboard
docker build -t svetu/wms-mobile:latest ./wms-frontend/apps/mobile

# 2. Push в registry
docker push svetu/wms-api:latest
docker push svetu/wms-dashboard:latest
docker push svetu/wms-mobile:latest

# 3. Deploy в Kubernetes
kubectl apply -f k8s/wms-namespace.yaml
kubectl apply -f k8s/wms-configmap.yaml
kubectl apply -f k8s/wms-secrets.yaml
kubectl apply -f k8s/wms-deployment.yaml
kubectl apply -f k8s/wms-service.yaml
kubectl apply -f k8s/wms-ingress.yaml

# 4. Run migrations
kubectl exec -it deployment/wms-api -- ./migrator up
```

## 8. Мониторинг и поддержка

### 8.1 Дашборды Grafana

```json
{
  "dashboard": {
    "title": "WMS Operations",
    "panels": [
      {
        "title": "Orders Processing Rate",
        "targets": [{
          "expr": "rate(wms_orders_processed_total[5m])"
        }]
      },
      {
        "title": "Inventory Accuracy",
        "targets": [{
          "expr": "wms_inventory_accuracy_percent"
        }]
      },
      {
        "title": "Picking Performance",
        "targets": [{
          "expr": "histogram_quantile(0.95, wms_picking_duration_seconds_bucket)"
        }]
      },
      {
        "title": "Warehouse Utilization",
        "targets": [{
          "expr": "wms_location_utilization_percent"
        }]
      }
    ]
  }
}
```

### 8.2 SLA и поддержка

| Метрика | Целевое значение |
|---------|------------------|
| Uptime | 99.9% |
| Response time (API) | < 200ms |
| Order processing time | < 2 часа |
| Inventory accuracy | > 99.5% |
| Picking accuracy | > 99.9% |

## 9. Roadmap интеграции

### Q1 2025 - MVP
- Базовая WMS функциональность
- Интеграция с основной платформой
- Пилот с 10 витринами

### Q2 2025 - Расширение
- Мобильное приложение для сборщиков
- Интеграция всех курьеров
- Автоматизация процессов

### Q3 2025 - Оптимизация
- AI для прогнозирования спроса
- Мультискладская логика
- B2B функциональность

### Q4 2025 - Масштабирование
- Международные отправки
- White-label WMS как продукт
- Интеграция с ERP системами

## Заключение

Собственная WMS система позволит Sve Tu:
1. Полностью контролировать складские процессы
2. Быстро адаптироваться под потребности бизнеса
3. Экономить на лицензиях (€3,600/год)
4. Создать конкурентное преимущество
5. Масштабироваться без ограничений

Инвестиции в €15,000 окупятся за 12 месяцев и обеспечат технологическое лидерство на рынке.