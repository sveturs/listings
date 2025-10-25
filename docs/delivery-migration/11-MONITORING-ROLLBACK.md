# Мониторинг и Rollback
          periodSeconds: 5
        resources:
          requests:
            memory: "256Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
```

---

## 📊 Мониторинг

### Prometheus Metrics

**Файл**: `internal/server/grpc/metrics.go`

```go
var (
    grpcRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "delivery_grpc_requests_total",
            Help: "Total number of gRPC requests",
        },
        []string{"method", "status"},
    )

    grpcRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "delivery_grpc_request_duration_seconds",
            Help:    "Duration of gRPC requests",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method"},
    )

    shipmentsCreatedTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "delivery_shipments_created_total",
            Help: "Total number of shipments created",
        },
        []string{"provider"},
    )
)
```

### Grafana Dashboard

**Панели**:
- Request rate (RPS)
- Request latency (p50, p95, p99)
- Error rate
- Shipments created by provider
- Active shipments by status

---

## 🔄 Rollback Plan

Если что-то пойдет не так:

```bash
# 1. Остановка микросервиса
docker-compose -f docker-compose.dev.yml stop delivery-service

# 2. Откат монолита к предыдущей версии
git checkout HEAD~1
