# Delivery Microservice - Быстрый старт

## 🎯 Что нужно

Реализовать 3 RPC метода в отдельном gRPC микросервисе:
1. **GetSettlements** - поиск населенных пунктов (TX 3)
2. **GetStreets** - поиск улиц (TX 4)
3. **GetParcelLockers** - список паккетоматов (TX 10)

---

## 📂 Где создавать

```bash
mkdir -p /data/hostel-booking-system/services/delivery-service
cd /data/hostel-booking-system/services/delivery-service
```

---

## ⚡ Минимальная реализация

### 1. Инициализация проекта

```bash
# Init go module
go mod init services/delivery-service

# Install dependencies
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf@latest
go get github.com/rs/zerolog@latest
go get github.com/joho/godotenv@latest

go mod tidy
```

### 2. Скопировать proto и сгенерировать

```bash
# Copy proto
mkdir -p proto/delivery/v1
cp /data/hostel-booking-system/backend/proto/delivery/v1/delivery.proto proto/delivery/v1/

# Generate Go code
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/delivery/v1/delivery.proto
```

### 3. Создать .env

```bash
SERVICE_NAME=delivery-service
GRPC_PORT=50051
LOG_LEVEL=debug

WSP_ENDPOINT=https://wsp.posta.rs/api
WSP_USERNAME=your_username
WSP_PASSWORD=your_password
WSP_PARTNER_ID=10109
```

### 4. Структура файлов

```
services/delivery-service/
├── cmd/server/main.go          # Entry point
├── internal/
│   ├── config/config.go        # Load .env
│   ├── service/
│   │   ├── service.go          # DeliveryService struct
│   │   ├── settlements.go      # GetSettlements RPC
│   │   ├── streets.go          # GetStreets RPC
│   │   └── parcel_lockers.go   # GetParcelLockers RPC
│   ├── wspapi/
│   │   ├── client.go           # HTTP client
│   │   ├── transactions.go     # TX 3, 4, 10
│   │   └── types.go            # WSP structs
│   └── mapper/
│       └── mappers.go          # WSP -> Proto conversion
└── .env
```

---

## 🚀 Запуск

```bash
# Run
go run cmd/server/main.go

# Test
grpcurl -plaintext -d '{"provider":1,"country":"RS","search_query":"Београд"}' \
  localhost:50051 delivery.v1.DeliveryService/GetSettlements
```

---

## 📖 Полная документация

См. `/data/hostel-booking-system/docs/DELIVERY_MICROSERVICE_IMPLEMENTATION_GUIDE.md`

В полном руководстве:
- ✅ Подробная структура проекта
- ✅ Примеры кода для каждого метода
- ✅ WSP API integration details
- ✅ Dockerfile и docker-compose
- ✅ Troubleshooting guide

---

## ✅ После реализации

Main backend endpoints начнут работать автоматически:

```bash
# Вместо "Unimplemented" вернут реальные данные
curl 'http://localhost:3000/api/public/delivery/test/settlements?country=RS'
curl 'http://localhost:3000/api/public/delivery/test/streets?settlement_name=Beograd'
curl 'http://localhost:3000/api/public/delivery/test/parcel-lockers?city=Beograd'
```

---

## 🔗 Справка

- **Proto:** `/data/hostel-booking-system/backend/proto/delivery/v1/delivery.proto`
- **WSP Client example:** `/data/hostel-booking-system/backend/internal/proj/postexpress/service/client.go`
- **WSP Types:** `/data/hostel-booking-system/backend/internal/proj/postexpress/types.go`
