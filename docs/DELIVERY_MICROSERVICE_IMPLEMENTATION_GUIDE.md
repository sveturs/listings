# DELIVERY MICROSERVICE - Техническое задание для реализации

> **Статус:** Ready to implement
> **Дата создания:** 2025-10-23
> **Версия:** 1.0.0

---

## 🎯 Что нужно сделать

Создать отдельный gRPC микросервис `delivery-service`, который реализует 3 RPC метода для работы с Post Express WSP API:

1. **GetSettlements** - поиск населенных пунктов (вызывает TX 3)
2. **GetStreets** - поиск улиц по населенному пункту (вызывает TX 4)
3. **GetParcelLockers** - получение списка паккетоматов (вызывает TX 10 с фильтрацией)

---

## 📍 Где создавать микросервис

```bash
# Создать новую директорию для микросервиса
mkdir -p /data/hostel-booking-system/services/delivery-service
cd /data/hostel-booking-system/services/delivery-service
```

**ВАЖНО:** Это отдельный проект, НЕ внутри backend!

---

## 📂 Структура проекта

```
services/delivery-service/
├── cmd/
│   └── server/
│       └── main.go                 # Точка входа gRPC сервера
├── internal/
│   ├── config/
│   │   └── config.go               # Конфигурация из ENV
│   ├── server/
│   │   └── grpc.go                 # gRPC server setup
│   ├── service/
│   │   ├── delivery_service.go     # Реализация DeliveryService
│   │   ├── settlements.go          # GetSettlements метод
│   │   ├── streets.go              # GetStreets метод
│   │   └── parcel_lockers.go       # GetParcelLockers метод
│   ├── wspapi/
│   │   ├── client.go               # WSP API HTTP клиент
│   │   ├── transactions.go         # TX 3, TX 4, TX 10 методы
│   │   └── types.go                # WSP типы данных
│   └── mapper/
│       ├── settlements.go          # Proto <-> WSP маппинг
│       ├── streets.go              # Proto <-> WSP маппинг
│       └── parcel_lockers.go       # Proto <-> WSP маппинг
├── pkg/
│   └── logger/
│       └── logger.go               # Логирование (zerolog)
├── proto/
│   └── delivery/
│       └── v1/
│           └── delivery.proto      # Копия из backend/proto
├── Dockerfile                      # Multi-stage build
├── docker-compose.yml              # Локальная разработка
├── Makefile                        # Build команды
├── go.mod                          # Go зависимости
├── .env.example                    # Пример конфигурации
└── README.md                       # Документация
```

---

## 📦 Зависимости (go.mod)

```bash
cd /data/hostel-booking-system/services/delivery-service
go mod init services/delivery-service

# Установить необходимые пакеты
go get google.golang.org/grpc@latest
go get google.golang.org/protobuf@latest
go get github.com/rs/zerolog@latest
go get github.com/joho/godotenv@latest
go get github.com/go-playground/validator/v10@latest

go mod tidy
```

---

## 🔧 Конфигурация (.env)

```bash
# Service Configuration
SERVICE_NAME=delivery-service
SERVICE_VERSION=1.0.0
GRPC_PORT=50051
LOG_LEVEL=debug

# Post Express WSP API
WSP_ENDPOINT=https://wsp.posta.rs/api
WSP_USERNAME=your_username
WSP_PASSWORD=your_password
WSP_LANGUAGE=sr-Latn
WSP_DEVICE_TYPE=2
WSP_PARTNER_ID=10109
WSP_TIMEOUT_SECONDS=30
WSP_MAX_RETRIES=3
```

---

## 🚀 Реализация RPC методов

### 1. GetSettlements (TX 3)

**Что делает:**
- Принимает `search_query` (название города, например "Београд")
- Вызывает Post Express WSP API TX 3
- Возвращает список населенных пунктов

**Пример WSP запроса:**
```json
{
  "Servis": 101,
  "StrKlijent": "{\"Username\":\"user\",\"Password\":\"pass\"}",
  "Transakcija": 3,
  "Ulazni": "{\"Naziv\":\"Београд\"}"
}
```

**Пример WSP ответа:**
```json
{
  "Success": true,
  "OutputData": {
    "Rezultat": 0,
    "Naselja": [
      {
        "Id": 123,
        "Naziv": "Београд",
        "PostanskiBroj": "11000"
      }
    ]
  }
}
```

**Код реализации:**

`internal/service/settlements.go`:
```go
func (s *DeliveryService) GetSettlements(ctx context.Context, req *pb.GetSettlementsRequest) (*pb.GetSettlementsResponse, error) {
    // 1. Валидация provider
    if req.Provider != pb.DeliveryProvider_DELIVERY_PROVIDER_POST_EXPRESS {
        return nil, fmt.Errorf("unsupported provider: %s", req.Provider)
    }

    // 2. Вызов WSP API TX 3
    wspResp, err := s.wspClient.GetSettlements(ctx, req.SearchQuery)
    if err != nil {
        return nil, fmt.Errorf("failed to get settlements: %w", err)
    }

    // 3. Проверка результата
    if wspResp.Rezultat != 0 {
        return nil, fmt.Errorf("Post Express error: %s", wspResp.Poruka)
    }

    // 4. Маппинг WSP -> Proto
    settlements := mapper.MapSettlementsWSPToProto(wspResp.Naselja, req.Country)

    return &pb.GetSettlementsResponse{
        Settlements: settlements,
    }, nil
}
```

---

### 2. GetStreets (TX 4)

**Что делает:**
- Принимает `settlement_name` (название города, например "Београд")
- Сначала находит ID населенного пункта через TX 3
- Затем вызывает TX 4 с этим ID
- Возвращает список улиц

**ВАЖНО:** TX 4 требует `IdNaselje` (ID населенного пункта), а не название!

**Пример WSP запроса:**
```json
{
  "Servis": 101,
  "Transakcija": 4,
  "Ulazni": "{\"IdNaselje\":123,\"Naziv\":\"Кнеза\"}"
}
```

**Код реализации:**

`internal/service/streets.go`:
```go
func (s *DeliveryService) GetStreets(ctx context.Context, req *pb.GetStreetsRequest) (*pb.GetStreetsResponse, error) {
    // 1. Найти ID населенного пункта по имени
    settlementsResp, err := s.wspClient.GetSettlements(ctx, req.SettlementName)
    if err != nil {
        return nil, fmt.Errorf("failed to find settlement: %w", err)
    }

    if len(settlementsResp.Naselja) == 0 {
        return nil, fmt.Errorf("settlement not found: %s", req.SettlementName)
    }

    settlementID := settlementsResp.Naselja[0].Id

    // 2. Вызов TX 4 с ID населенного пункта
    wspResp, err := s.wspClient.GetStreets(ctx, settlementID, req.SearchQuery)
    if err != nil {
        return nil, fmt.Errorf("failed to get streets: %w", err)
    }

    // 3. Маппинг WSP -> Proto
    streets := mapper.MapStreetsWSPToProto(wspResp.Ulice, req.SettlementName)

    return &pb.GetStreetsResponse{
        Streets: streets,
    }, nil
}
```

---

### 3. GetParcelLockers (TX 10)

**Что делает:**
- Принимает `city` (опционально) и `search_query`
- Находит ID города через TX 3 (если указан)
- Вызывает TX 10 (GetOffices) - получает ВСЕ отделения
- Фильтрует только паккетоматы (`TipPoste == "PL"`)
- Возвращает список паккетоматов

**ВАЖНО:** У Post Express нет отдельной транзакции для паккетоматов! Используем TX 10 с фильтрацией.

**Пример WSP ответа TX 10:**
```json
{
  "Success": true,
  "OutputData": {
    "Rezultat": 0,
    "PostanskeJedinice": [
      {
        "IdPoste": 456,
        "SifraPoste": "BG001",
        "Naziv": "Пакетомат Немањина",
        "TipPoste": "PL",
        "Adresa": "Немањина 2",
        "Mesto": "Београд",
        "PostanskiBroj": "11000",
        "Latitude": 44.816,
        "Longitude": 20.456
      }
    ]
  }
}
```

**Код реализации:**

`internal/service/parcel_lockers.go`:
```go
func (s *DeliveryService) GetParcelLockers(ctx context.Context, req *pb.GetParcelLockersRequest) (*pb.GetParcelLockersResponse, error) {
    // 1. Найти ID города
    var settlementID int
    if req.City != "" {
        settlementsResp, err := s.wspClient.GetSettlements(ctx, req.City)
        if err != nil {
            return nil, fmt.Errorf("failed to find city: %w", err)
        }
        if len(settlementsResp.Naselja) == 0 {
            return nil, fmt.Errorf("city not found: %s", req.City)
        }
        settlementID = settlementsResp.Naselja[0].Id
    } else {
        // Default: Белград
        settlementsResp, _ := s.wspClient.GetSettlements(ctx, "Beograd")
        settlementID = settlementsResp.Naselja[0].Id
    }

    // 2. Получить все отделения через TX 10
    wspResp, err := s.wspClient.GetOffices(ctx, settlementID)
    if err != nil {
        return nil, fmt.Errorf("failed to get offices: %w", err)
    }

    // 3. Фильтровать только паккетоматы
    var parcelLockers []wspapi.Office
    for _, office := range wspResp.PostanskeJedinice {
        if office.TipPoste == "PL" {
            // Дополнительная фильтрация по search_query
            if req.SearchQuery == "" ||
               strings.Contains(strings.ToLower(office.Naziv), strings.ToLower(req.SearchQuery)) {
                parcelLockers = append(parcelLockers, office)
            }
        }
    }

    // 4. Маппинг WSP -> Proto
    lockers := mapper.MapParcelLockersWSPToProto(parcelLockers)

    return &pb.GetParcelLockersResponse{
        ParcelLockers: lockers,
    }, nil
}
```

---

## 🏃 Запуск микросервиса

### 1. Генерация proto файлов

```bash
# Скопировать proto из backend
cp /data/hostel-booking-system/backend/proto/delivery/v1/delivery.proto \
   proto/delivery/v1/

# Сгенерировать Go код
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/delivery/v1/delivery.proto
```

### 2. Создать main.go

`cmd/server/main.go`:
```go
package main

import (
    "fmt"
    "net"
    "os"

    "google.golang.org/grpc"
    "google.golang.org/grpc/reflection"

    "services/delivery-service/internal/config"
    "services/delivery-service/internal/service"
    "services/delivery-service/internal/wspapi"
    "services/delivery-service/pkg/logger"
    pb "services/delivery-service/proto/delivery/v1"
)

func main() {
    // Загрузить конфигурацию
    cfg, err := config.Load()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
        os.Exit(1)
    }

    // Инициализировать логгер
    log := logger.New(cfg.LogLevel)

    // Создать WSP клиент
    wspClient := wspapi.NewClient(cfg.WSP, log)

    // Создать gRPC сервис
    deliveryService := service.NewDeliveryService(wspClient, log)

    // Запустить gRPC сервер
    listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
    if err != nil {
        log.Fatal("Failed to listen: %v", err)
    }

    grpcServer := grpc.NewServer()
    pb.RegisterDeliveryServiceServer(grpcServer, deliveryService)

    // Включить reflection для grpcurl
    reflection.Register(grpcServer)

    log.Info("Starting delivery gRPC service on port %s", cfg.GRPCPort)
    if err := grpcServer.Serve(listener); err != nil {
        log.Fatal("Failed to serve: %v", err)
    }
}
```

### 3. Запуск

```bash
# Создать .env
cp .env.example .env
# Добавить реальные credentials WSP_USERNAME и WSP_PASSWORD

# Запустить
go run cmd/server/main.go

# Вывод:
# 2025/10/23 23:00:00 INFO Starting delivery gRPC service on port 50051
```

---

## 🧪 Тестирование

### Проверка что сервис работает

```bash
# Установить grpcurl (если нет)
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# Проверить доступные методы
grpcurl -plaintext localhost:50051 list delivery.v1.DeliveryService

# Вывод:
# delivery.v1.DeliveryService.GetSettlements
# delivery.v1.DeliveryService.GetStreets
# delivery.v1.DeliveryService.GetParcelLockers
```

### Тест GetSettlements

```bash
grpcurl -plaintext -d '{
  "provider": 1,
  "country": "RS",
  "search_query": "Београд"
}' localhost:50051 delivery.v1.DeliveryService/GetSettlements

# Ожидаемый ответ:
# {
#   "settlements": [
#     {
#       "id": 123,
#       "name": "Београд",
#       "zip_code": "11000",
#       "country": "RS"
#     }
#   ]
# }
```

### Тест GetStreets

```bash
grpcurl -plaintext -d '{
  "provider": 1,
  "settlement_name": "Београд",
  "search_query": "Кнеза"
}' localhost:50051 delivery.v1.DeliveryService/GetStreets
```

### Тест GetParcelLockers

```bash
grpcurl -plaintext -d '{
  "provider": 1,
  "city": "Београд"
}' localhost:50051 delivery.v1.DeliveryService/GetParcelLockers
```

---

## 🔗 Интеграция с main backend

После того как микросервис запущен и работает, main backend автоматически начнёт использовать его:

```bash
# Main backend уже настроен!
# Эндпоинты начнут возвращать данные вместо "Unimplemented"

curl -H "Authorization: Bearer $(cat /tmp/token)" \
  'http://localhost:3000/api/public/delivery/test/settlements?country=RS&search_query=Београд'

# Вместо ошибки "Unimplemented" получишь:
# {
#   "success": true,
#   "data": {
#     "settlements": [...],
#     "count": 5
#   }
# }
```

---

## 📚 Справочная информация

### Post Express WSP API транзакции

| TX  | Название              | Описание                          | Используется в |
|-----|-----------------------|-----------------------------------|----------------|
| 3   | GetNaselje            | Поиск населенных пунктов          | GetSettlements |
| 4   | GetUlica              | Поиск улиц по населенному пункту  | GetStreets     |
| 10  | GetPostanskaJedinica  | Получение отделений почты         | GetParcelLockers (с фильтром `TipPoste == "PL"`) |
| 73  | B2BManifest           | Создание отправления              | CreateShipment (already implemented) |

### Где искать примеры кода

1. **WSP Client реализация:** `/data/hostel-booking-system/backend/internal/proj/postexpress/service/client.go`
2. **WSP типы данных:** `/data/hostel-booking-system/backend/internal/proj/postexpress/types.go`
3. **Proto контракт:** `/data/hostel-booking-system/backend/proto/delivery/v1/delivery.proto`
4. **gRPC клиент (main backend):** `/data/hostel-booking-system/backend/internal/proj/delivery/grpcclient/client.go`

### Архитектура взаимодействия

```
Browser
   ↓ HTTP
Main Backend (Fiber)
http://localhost:3000
   ↓ gRPC (internal)
Delivery Microservice
grpc://localhost:50051
   ↓ HTTPS
Post Express WSP API
https://wsp.posta.rs/api
```

---

## ✅ Чеклист реализации

### Базовая инфраструктура
- [ ] Создать структуру папок
- [ ] Инициализировать go.mod
- [ ] Скопировать и сгенерировать proto
- [ ] Создать config.go
- [ ] Создать logger.go

### WSP Client
- [ ] Реализовать базовый HTTP клиент
- [ ] Реализовать метод Transaction (базовый)
- [ ] Добавить TX 3 (GetSettlements)
- [ ] Добавить TX 4 (GetStreets)
- [ ] Добавить TX 10 (GetOffices)

### gRPC Service
- [ ] Реализовать GetSettlements RPC
- [ ] Реализовать GetStreets RPC
- [ ] Реализовать GetParcelLockers RPC
- [ ] Создать mappers (WSP -> Proto)

### Тестирование
- [ ] Тест GetSettlements через grpcurl
- [ ] Тест GetStreets через grpcurl
- [ ] Тест GetParcelLockers через grpcurl
- [ ] Интеграционный тест через main backend

### Production
- [ ] Dockerfile
- [ ] docker-compose.yml
- [ ] README.md с документацией
- [ ] Деплой на dev сервер

---

## 🎯 Итог

После реализации микросервиса по этому ТЗ:

✅ Main backend endpoints начнут работать:
- `GET /api/public/delivery/test/settlements`
- `GET /api/public/delivery/test/streets`
- `GET /api/public/delivery/test/parcel-lockers`

✅ Они вернут реальные данные вместо "Unimplemented"

✅ Микросервисная архитектура полностью готова к продакшену

---

**Подробное техническое задание с примерами кода создано агентом и сохранено в этом документе.**
