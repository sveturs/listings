# Delivery Microservice API Reference

## 📋 Table of Contents

- [Overview](#overview)
- [Protocol & Connection](#protocol--connection)
- [Service Definition](#service-definition)
- [gRPC Methods](#grpc-methods)
- [Message Types](#message-types)
- [Enumerations](#enumerations)
- [Error Codes](#error-codes)
- [Examples](#examples)
- [Testing](#testing)

---

## Overview

Delivery Microservice предоставляет gRPC API для управления доставкой заказов через различных провайдеров (Post Express, BEX Express, AKS Express, D Express, City Express).

### Key Features

- ✅ **Multi-provider aggregation**: Единый API для всех провайдеров
- ✅ **Rate calculation**: Сравнение цен от всех провайдеров
- ✅ **Shipment management**: Создание, отслеживание, отмена
- ✅ **Real-time tracking**: Актуальная информация о статусе доставки
- ✅ **Idempotency**: Защита от дубликатов через order_id
- ✅ **Error handling**: Детальные коды ошибок

### Technology Stack

- **Protocol:** gRPC (Protocol Buffers v3)
- **Language:** Go 1.22+
- **Port:** 50053
- **Format:** Binary (protobuf)

---

## Protocol & Connection

### Connection String

```
localhost:50053                    # Development
delivery-service.internal:50053    # Production (internal network)
```

### gRPC Client Setup (Go)

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "backend/pkg/grpc/delivery/v1"
)

// Create connection
conn, err := grpc.NewClient(
    "localhost:50053",
    grpc.WithTransportCredentials(insecure.NewCredentials()),
)
if err != nil {
    log.Fatal(err)
}
defer conn.Close()

// Create client
client := pb.NewDeliveryServiceClient(conn)
```

### Health Check

```bash
# gRPC health check (requires grpc_health_probe)
grpc_health_probe -addr=localhost:50053

# HTTP health endpoint (if available)
curl http://localhost:50053/health
```

---

## Service Definition

### Proto File Location

```
proto/delivery/v1/delivery.proto
```

### Generated Code

```
backend/pkg/grpc/delivery/v1/
├── delivery.pb.go          # Message types
└── delivery_grpc.pb.go     # Service client/server
```

### Service Interface

```protobuf
service DeliveryService {
  // Calculate delivery rate for given parameters
  rpc CalculateRate(CalculateRateRequest) returns (CalculateRateResponse);

  // Create a new shipment
  rpc CreateShipment(CreateShipmentRequest) returns (CreateShipmentResponse);

  // Get shipment details
  rpc GetShipment(GetShipmentRequest) returns (GetShipmentResponse);

  // Track shipment status
  rpc TrackShipment(TrackShipmentRequest) returns (TrackShipmentResponse);

  // Cancel a shipment
  rpc CancelShipment(CancelShipmentRequest) returns (CancelShipmentResponse);

  // List all available providers
  rpc ListProviders(ListProvidersRequest) returns (ListProvidersResponse);
}
```

---

## gRPC Methods

### 1. CalculateRate

Рассчитывает стоимость доставки для указанных параметров от одного или нескольких провайдеров.

#### Request

```protobuf
message CalculateRateRequest {
  // Provider to calculate rate for (empty = all providers)
  string provider_code = 1;  // "post_express", "bex_express", etc.

  // From/To addresses
  Address from_address = 2;
  Address to_address = 3;

  // Package details
  repeated Package packages = 4;

  // Additional options
  bool insurance = 5;
  double insurance_value = 6;  // RSD
  bool cash_on_delivery = 7;
  double cod_amount = 8;  // RSD
}
```

#### Response

```protobuf
message CalculateRateResponse {
  // If single provider requested
  RateQuote quote = 1;

  // If multiple providers requested
  repeated RateQuote quotes = 2;

  // Recommended provider (best price/time balance)
  string recommended_provider = 3;
}

message RateQuote {
  string provider_code = 1;
  string provider_name = 2;

  // Cost breakdown
  double base_price = 3;
  double insurance_fee = 4;
  double cod_fee = 5;
  double weight_fee = 6;
  double distance_fee = 7;
  double total_cost = 8;
  string currency = 9;  // "RSD"

  // Delivery estimate
  int32 estimated_delivery_days = 10;
  google.protobuf.Timestamp estimated_delivery_date = 11;
}
```

#### Example (Go)

```go
resp, err := client.CalculateRate(ctx, &pb.CalculateRateRequest{
    ProviderCode: "post_express",
    FromAddress: &pb.Address{
        City:       "Belgrade",
        PostalCode: "11000",
        Country:    "RS",
    },
    ToAddress: &pb.Address{
        City:       "Novi Sad",
        PostalCode: "21000",
        Country:    "RS",
    },
    Packages: []*pb.Package{
        {
            Weight: 2.5,   // kg
            Length: 30,    // cm
            Width:  20,
            Height: 10,
        },
    },
    Insurance:      true,
    InsuranceValue: 10000,  // RSD
    CashOnDelivery: true,
    CodAmount:      5000,   // RSD
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Total cost: %.2f %s\n", resp.Quote.TotalCost, resp.Quote.Currency)
fmt.Printf("Estimated delivery: %d days\n", resp.Quote.EstimatedDeliveryDays)
```

#### Error Codes

- `INVALID_ARGUMENT` - Missing required fields or invalid values
- `NOT_FOUND` - Provider not found or not available in region
- `UNAVAILABLE` - Provider API temporarily unavailable
- `DEADLINE_EXCEEDED` - Provider API timeout

---

### 2. CreateShipment

Создает новое отправление у выбранного провайдера.

#### Request

```protobuf
message CreateShipmentRequest {
  // Provider
  string provider_code = 1;  // Required

  // Order reference (for idempotency)
  string order_id = 2;  // Optional but recommended

  // Addresses
  Address from_address = 3;  // Required
  Address to_address = 4;    // Required

  // Packages
  repeated Package packages = 5;  // At least 1 required

  // Additional services
  bool insurance = 6;
  double insurance_value = 7;
  bool cash_on_delivery = 8;
  double cod_amount = 9;
  bool sms_notification = 10;
  bool email_notification = 11;

  // Pickup/Delivery preferences
  string pickup_date = 12;  // YYYY-MM-DD (optional)
  string delivery_instructions = 13;
}

message Address {
  string street = 1;
  string city = 2;
  string state = 3;
  string postal_code = 4;
  string country = 5;  // ISO 3166-1 alpha-2 (e.g. "RS")
  string contact_name = 6;
  string contact_phone = 7;
}

message Package {
  double weight = 1;   // kg
  double length = 2;   // cm
  double width = 3;
  double height = 4;
  string description = 5;
  double declared_value = 6;  // RSD
}
```

#### Response

```protobuf
message CreateShipmentResponse {
  Shipment shipment = 1;
  string label_url = 2;  // URL to download shipping label (PDF)
}

message Shipment {
  string id = 1;                    // Internal ID
  string tracking_number = 2;       // Provider tracking number
  string provider_code = 3;
  ShipmentStatus status = 4;
  google.protobuf.Timestamp created_at = 5;
  google.protobuf.Timestamp updated_at = 6;
}
```

#### Example (Go)

```go
resp, err := client.CreateShipment(ctx, &pb.CreateShipmentRequest{
    ProviderCode: "post_express",
    OrderId:      "12345",  // For idempotency
    FromAddress: &pb.Address{
        Street:       "Main Street 1",
        City:         "Belgrade",
        PostalCode:   "11000",
        Country:      "RS",
        ContactName:  "Store Name",
        ContactPhone: "+381601234567",
    },
    ToAddress: &pb.Address{
        Street:       "Customer Street 5",
        City:         "Novi Sad",
        PostalCode:   "21000",
        Country:      "RS",
        ContactName:  "Customer Name",
        ContactPhone: "+381607654321",
    },
    Packages: []*pb.Package{
        {
            Weight:        2.5,
            Length:        30,
            Width:         20,
            Height:        10,
            Description:   "Order #12345",
            DeclaredValue: 10000,
        },
    },
    Insurance:      true,
    InsuranceValue: 10000,
    CashOnDelivery: true,
    CodAmount:      5000,
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Shipment created: %s\n", resp.Shipment.TrackingNumber)
fmt.Printf("Label URL: %s\n", resp.LabelUrl)
```

#### Idempotency

**Important:** Если указан `order_id`, delivery service проверяет дубликаты:

```go
// First call - creates shipment
resp1, _ := client.CreateShipment(ctx, &pb.CreateShipmentRequest{
    OrderId: "12345",
    // ... other fields
})
// Returns: shipment_id=1, tracking_number=PE1234567890RS

// Second call with same order_id - returns existing shipment
resp2, _ := client.CreateShipment(ctx, &pb.CreateShipmentRequest{
    OrderId: "12345",  // Same order_id!
    // ... other fields
})
// Returns: shipment_id=1, tracking_number=PE1234567890RS (same!)
```

#### Error Codes

- `INVALID_ARGUMENT` - Missing required fields, invalid address, invalid package dimensions
- `ALREADY_EXISTS` - Shipment already exists for this order_id
- `NOT_FOUND` - Provider not found
- `UNAVAILABLE` - Provider API unavailable
- `RESOURCE_EXHAUSTED` - Provider rate limit exceeded

---

### 3. GetShipment

Получает детальную информацию о shipment.

#### Request

```protobuf
message GetShipmentRequest {
  string shipment_id = 1;  // Internal ID or tracking number
}
```

#### Response

```protobuf
message GetShipmentResponse {
  Shipment shipment = 1;
  Address from_address = 2;
  Address to_address = 3;
  repeated Package packages = 4;
  double total_cost = 5;
  string currency = 6;
  string label_url = 7;
}
```

#### Example (Go)

```go
resp, err := client.GetShipment(ctx, &pb.GetShipmentRequest{
    ShipmentId: "PE1234567890RS",
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", resp.Shipment.Status)
fmt.Printf("From: %s\n", resp.FromAddress.City)
fmt.Printf("To: %s\n", resp.ToAddress.City)
```

#### Error Codes

- `NOT_FOUND` - Shipment not found
- `INVALID_ARGUMENT` - Empty shipment_id

---

### 4. TrackShipment

Отслеживает текущий статус и историю доставки.

#### Request

```protobuf
message TrackShipmentRequest {
  string tracking_number = 1;  // Required
}
```

#### Response

```protobuf
message TrackShipmentResponse {
  string tracking_number = 1;
  ShipmentStatus status = 2;
  string current_location = 3;

  // Delivery estimates
  google.protobuf.Timestamp estimated_delivery = 4;
  google.protobuf.Timestamp actual_delivery = 5;

  // Tracking events (sorted by timestamp DESC)
  repeated TrackingEvent events = 6;
}

message TrackingEvent {
  google.protobuf.Timestamp timestamp = 1;
  string location = 2;
  ShipmentStatus status = 3;
  string description = 4;
  string notes = 5;
}
```

#### Example (Go)

```go
resp, err := client.TrackShipment(ctx, &pb.TrackShipmentRequest{
    TrackingNumber: "PE1234567890RS",
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Status: %s\n", resp.Status)
fmt.Printf("Current location: %s\n", resp.CurrentLocation)
fmt.Printf("Events:\n")
for _, event := range resp.Events {
    fmt.Printf("  %s - %s: %s\n",
        event.Timestamp.AsTime().Format("2006-01-02 15:04"),
        event.Location,
        event.Description,
    )
}
```

#### Example Response

```
Status: IN_TRANSIT
Current location: Postal center Belgrade
Events:
  2025-10-29 10:00 - Belgrade depot: Shipment confirmed and accepted
  2025-10-29 14:30 - Postal center Belgrade: Package in transit to Novi Sad
  2025-10-30 08:15 - Novi Sad depot: Package arrived at destination city
```

#### Error Codes

- `NOT_FOUND` - Tracking number not found
- `INVALID_ARGUMENT` - Empty tracking_number
- `UNAVAILABLE` - Provider tracking API unavailable

---

### 5. CancelShipment

Отменяет shipment (если еще не доставлен).

#### Request

```protobuf
message CancelShipmentRequest {
  string shipment_id = 1;      // Internal ID or tracking number
  string cancellation_reason = 2;  // Optional
}
```

#### Response

```protobuf
message CancelShipmentResponse {
  bool success = 1;
  string message = 2;
  ShipmentStatus new_status = 3;  // Should be CANCELLED
}
```

#### Example (Go)

```go
resp, err := client.CancelShipment(ctx, &pb.CancelShipmentRequest{
    ShipmentId:         "PE1234567890RS",
    CancellationReason: "Customer cancelled order",
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Cancelled: %v\n", resp.Success)
fmt.Printf("Message: %s\n", resp.Message)
```

#### Error Codes

- `NOT_FOUND` - Shipment not found
- `FAILED_PRECONDITION` - Shipment already delivered or cannot be cancelled
- `UNAVAILABLE` - Provider API unavailable

---

### 6. ListProviders

Получает список доступных провайдеров доставки.

#### Request

```protobuf
message ListProvidersRequest {
  string city = 1;  // Filter by service area (optional)
}
```

#### Response

```protobuf
message ListProvidersResponse {
  repeated Provider providers = 1;
}

message Provider {
  string code = 1;           // "post_express"
  string name = 2;           // "Post Express"
  string description = 3;
  bool enabled = 4;
  string logo_url = 5;
  repeated string service_areas = 6;  // Cities where available
}
```

#### Example (Go)

```go
resp, err := client.ListProviders(ctx, &pb.ListProvidersRequest{})

if err != nil {
    log.Fatal(err)
}

for _, provider := range resp.Providers {
    fmt.Printf("%s (%s) - Enabled: %v\n",
        provider.Name,
        provider.Code,
        provider.Enabled,
    )
}
```

---

## Message Types

### Address

```protobuf
message Address {
  string street = 1;           // Required
  string city = 2;             // Required
  string state = 3;            // Optional
  string postal_code = 4;      // Required
  string country = 5;          // Required (ISO 3166-1 alpha-2)
  string contact_name = 6;     // Required
  string contact_phone = 7;    // Required (E.164 format preferred)
}
```

**Validation:**
- `street`: 1-200 characters
- `city`: 1-100 characters
- `postal_code`: 1-20 characters
- `country`: Exactly 2 characters (ISO code)
- `contact_name`: 1-100 characters
- `contact_phone`: 1-20 characters

**Example (Serbia):**
```json
{
  "street": "Kralja Milana 10",
  "city": "Belgrade",
  "postal_code": "11000",
  "country": "RS",
  "contact_name": "Marko Markovic",
  "contact_phone": "+381601234567"
}
```

### Package

```protobuf
message Package {
  double weight = 1;         // kg, Required, > 0
  double length = 2;         // cm, Optional
  double width = 3;          // cm, Optional
  double height = 4;         // cm, Optional
  string description = 5;    // Optional, max 200 chars
  double declared_value = 6; // RSD, Optional
}
```

**Validation:**
- `weight`: > 0, max 50kg (provider-specific)
- `length, width, height`: If specified, must be > 0
- `declared_value`: If specified, must be >= 0

**Example:**
```json
{
  "weight": 2.5,
  "length": 30,
  "width": 20,
  "height": 10,
  "description": "Electronics - Smartphone",
  "declared_value": 50000
}
```

---

## Enumerations

### ShipmentStatus

```protobuf
enum ShipmentStatus {
  SHIPMENT_STATUS_UNSPECIFIED = 0;
  SHIPMENT_STATUS_PENDING = 1;           // Created, awaiting pickup
  SHIPMENT_STATUS_CONFIRMED = 2;         // Confirmed by provider
  SHIPMENT_STATUS_IN_TRANSIT = 3;        // Package in transit
  SHIPMENT_STATUS_OUT_FOR_DELIVERY = 4;  // Out for final delivery
  SHIPMENT_STATUS_DELIVERED = 5;         // Successfully delivered
  SHIPMENT_STATUS_FAILED = 6;            // Delivery failed
  SHIPMENT_STATUS_CANCELLED = 7;         // Cancelled by sender
  SHIPMENT_STATUS_RETURNED = 8;          // Returned to sender
}
```

**State transitions:**

```
PENDING → CONFIRMED → IN_TRANSIT → OUT_FOR_DELIVERY → DELIVERED
                ↓           ↓              ↓
            CANCELLED   FAILED         RETURNED
```

### DeliveryProvider

```protobuf
enum DeliveryProvider {
  DELIVERY_PROVIDER_UNSPECIFIED = 0;
  DELIVERY_PROVIDER_POST_EXPRESS = 1;  // post_express
  DELIVERY_PROVIDER_BEX_EXPRESS = 2;   // bex_express
  DELIVERY_PROVIDER_AKS_EXPRESS = 3;   // aks_express
  DELIVERY_PROVIDER_D_EXPRESS = 4;     // d_express
  DELIVERY_PROVIDER_CITY_EXPRESS = 5;  // city_express
}
```

**Mapping:**

| Enum | Code | Name |
|------|------|------|
| 1 | `post_express` | Post Express |
| 2 | `bex_express` | BEX Express |
| 3 | `aks_express` | AKS Express |
| 4 | `d_express` | D Express |
| 5 | `city_express` | City Express |

---

## Error Codes

### gRPC Status Codes

| Code | Name | Description | Retry? |
|------|------|-------------|--------|
| 0 | OK | Success | N/A |
| 3 | INVALID_ARGUMENT | Invalid request parameters | ❌ No |
| 5 | NOT_FOUND | Resource not found | ❌ No |
| 6 | ALREADY_EXISTS | Duplicate resource (idempotency) | ❌ No |
| 7 | PERMISSION_DENIED | Access denied | ❌ No |
| 8 | RESOURCE_EXHAUSTED | Rate limit exceeded | ✅ Yes (after delay) |
| 9 | FAILED_PRECONDITION | Operation not allowed in current state | ❌ No |
| 13 | INTERNAL | Internal server error | ✅ Yes |
| 14 | UNAVAILABLE | Service temporarily unavailable | ✅ Yes |
| 4 | DEADLINE_EXCEEDED | Request timeout | ✅ Yes |

### Error Details

Errors include detailed messages in metadata:

```go
st, ok := status.FromError(err)
if ok {
    fmt.Println("Code:", st.Code())
    fmt.Println("Message:", st.Message())

    // Get details
    for _, detail := range st.Details() {
        // Process error details
    }
}
```

### Common Error Scenarios

#### 1. Invalid Address

```
Code: INVALID_ARGUMENT
Message: Invalid address: postal_code is required for country RS
```

#### 2. Provider Unavailable

```
Code: UNAVAILABLE
Message: Post Express API is temporarily unavailable
Retry-After: 30s
```

#### 3. Shipment Not Found

```
Code: NOT_FOUND
Message: Shipment not found: PE1234567890RS
```

#### 4. Duplicate Shipment

```
Code: ALREADY_EXISTS
Message: Shipment already exists for order_id: 12345
Existing shipment: PE1234567890RS
```

#### 5. Cannot Cancel

```
Code: FAILED_PRECONDITION
Message: Cannot cancel shipment in status: DELIVERED
```

---

## Examples

### Complete Workflow Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    pb "backend/pkg/grpc/delivery/v1"
)

func main() {
    // 1. Connect to delivery service
    conn, err := grpc.NewClient(
        "localhost:50053",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    client := pb.NewDeliveryServiceClient(conn)
    ctx := context.Background()

    // 2. Calculate rate
    fmt.Println("=== Calculating Rate ===")
    rateResp, err := client.CalculateRate(ctx, &pb.CalculateRateRequest{
        ProviderCode: "", // Empty = all providers
        FromAddress: &pb.Address{
            City:       "Belgrade",
            PostalCode: "11000",
            Country:    "RS",
        },
        ToAddress: &pb.Address{
            City:       "Novi Sad",
            PostalCode: "21000",
            Country:    "RS",
        },
        Packages: []*pb.Package{
            {Weight: 2.5, Length: 30, Width: 20, Height: 10},
        },
        Insurance:      true,
        InsuranceValue: 10000,
    })
    if err != nil {
        log.Fatal(err)
    }

    // Print all quotes
    for _, quote := range rateResp.Quotes {
        fmt.Printf("%s: %.2f %s (%d days)\n",
            quote.ProviderName,
            quote.TotalCost,
            quote.Currency,
            quote.EstimatedDeliveryDays,
        )
    }

    // Select cheapest provider
    selectedProvider := rateResp.Quotes[0].ProviderCode
    fmt.Printf("\nSelected provider: %s\n", selectedProvider)

    // 3. Create shipment
    fmt.Println("\n=== Creating Shipment ===")
    shipmentResp, err := client.CreateShipment(ctx, &pb.CreateShipmentRequest{
        ProviderCode: selectedProvider,
        OrderId:      "TEST-12345",
        FromAddress: &pb.Address{
            Street:       "Main Street 1",
            City:         "Belgrade",
            PostalCode:   "11000",
            Country:      "RS",
            ContactName:  "Test Store",
            ContactPhone: "+381601234567",
        },
        ToAddress: &pb.Address{
            Street:       "Customer Street 5",
            City:         "Novi Sad",
            PostalCode:   "21000",
            Country:      "RS",
            ContactName:  "Test Customer",
            ContactPhone: "+381607654321",
        },
        Packages: []*pb.Package{
            {
                Weight:        2.5,
                Length:        30,
                Width:         20,
                Height:        10,
                Description:   "Test Order",
                DeclaredValue: 10000,
            },
        },
        Insurance:      true,
        InsuranceValue: 10000,
    })
    if err != nil {
        log.Fatal(err)
    }

    trackingNumber := shipmentResp.Shipment.TrackingNumber
    fmt.Printf("Shipment created: %s\n", trackingNumber)
    fmt.Printf("Label URL: %s\n", shipmentResp.LabelUrl)

    // 4. Track shipment
    fmt.Println("\n=== Tracking Shipment ===")
    time.Sleep(2 * time.Second) // Wait a bit

    trackResp, err := client.TrackShipment(ctx, &pb.TrackShipmentRequest{
        TrackingNumber: trackingNumber,
    })
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Status: %s\n", trackResp.Status)
    fmt.Printf("Current location: %s\n", trackResp.CurrentLocation)
    fmt.Println("Events:")
    for _, event := range trackResp.Events {
        fmt.Printf("  %s - %s: %s\n",
            event.Timestamp.AsTime().Format("2006-01-02 15:04"),
            event.Location,
            event.Description,
        )
    }

    // 5. Cancel shipment (optional)
    // cancelResp, err := client.CancelShipment(ctx, &pb.CancelShipmentRequest{
    //     ShipmentId:         trackingNumber,
    //     CancellationReason: "Test cancellation",
    // })
}
```

---

## Testing

### Unit Testing (Go)

```go
package grpcclient_test

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "google.golang.org/grpc"
    pb "backend/pkg/grpc/delivery/v1"
)

func TestCalculateRate(t *testing.T) {
    // Setup
    conn, err := grpc.NewClient(
        "localhost:50053",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
    )
    require.NoError(t, err)
    defer conn.Close()

    client := pb.NewDeliveryServiceClient(conn)
    ctx := context.Background()

    // Test
    resp, err := client.CalculateRate(ctx, &pb.CalculateRateRequest{
        ProviderCode: "post_express",
        FromAddress: &pb.Address{
            City:       "Belgrade",
            PostalCode: "11000",
            Country:    "RS",
        },
        ToAddress: &pb.Address{
            City:       "Novi Sad",
            PostalCode: "21000",
            Country:    "RS",
        },
        Packages: []*pb.Package{
            {Weight: 2.5},
        },
    })

    // Assert
    require.NoError(t, err)
    assert.NotNil(t, resp.Quote)
    assert.Greater(t, resp.Quote.TotalCost, 0.0)
    assert.Equal(t, "RSD", resp.Quote.Currency)
}
```

### Manual Testing (grpcurl)

```bash
# Install grpcurl
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

# List services
grpcurl -plaintext localhost:50053 list

# List methods
grpcurl -plaintext localhost:50053 list delivery.v1.DeliveryService

# Call CalculateRate
grpcurl -plaintext -d '{
  "provider_code": "post_express",
  "from_address": {
    "city": "Belgrade",
    "postal_code": "11000",
    "country": "RS"
  },
  "to_address": {
    "city": "Novi Sad",
    "postal_code": "21000",
    "country": "RS"
  },
  "packages": [
    {"weight": 2.5}
  ]
}' localhost:50053 delivery.v1.DeliveryService/CalculateRate

# Call CreateShipment
grpcurl -plaintext -d '{
  "provider_code": "post_express",
  "order_id": "TEST-123",
  "from_address": {
    "street": "Main St 1",
    "city": "Belgrade",
    "postal_code": "11000",
    "country": "RS",
    "contact_name": "Store",
    "contact_phone": "+381601234567"
  },
  "to_address": {
    "street": "Customer St 5",
    "city": "Novi Sad",
    "postal_code": "21000",
    "country": "RS",
    "contact_name": "Customer",
    "contact_phone": "+381607654321"
  },
  "packages": [
    {"weight": 2.5, "description": "Test"}
  ]
}' localhost:50053 delivery.v1.DeliveryService/CreateShipment
```

### Performance Testing

```bash
# Install ghz
go install github.com/bojand/ghz/cmd/ghz@latest

# Load test CalculateRate (100 requests, 10 concurrent)
ghz --insecure \
  --proto proto/delivery/v1/delivery.proto \
  --call delivery.v1.DeliveryService/CalculateRate \
  -d '{
    "provider_code": "post_express",
    "from_address": {"city": "Belgrade", "postal_code": "11000", "country": "RS"},
    "to_address": {"city": "Novi Sad", "postal_code": "21000", "country": "RS"},
    "packages": [{"weight": 2.5}]
  }' \
  -n 100 \
  -c 10 \
  localhost:50053
```

---

## Best Practices

### 1. Always use context with timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

resp, err := client.CalculateRate(ctx, req)
```

### 2. Implement retry logic

```go
var resp *pb.CalculateRateResponse
var err error

for attempt := 0; attempt < 3; attempt++ {
    resp, err = client.CalculateRate(ctx, req)
    if err == nil {
        break
    }

    // Check if retryable
    if st, ok := status.FromError(err); ok {
        if st.Code() == codes.Unavailable || st.Code() == codes.DeadlineExceeded {
            time.Sleep(time.Duration(attempt+1) * 100 * time.Millisecond)
            continue
        }
    }

    break // Non-retryable error
}
```

### 3. Use idempotency keys

```go
// Always include order_id for CreateShipment
resp, err := client.CreateShipment(ctx, &pb.CreateShipmentRequest{
    OrderId: fmt.Sprintf("ORDER-%d", orderID), // Prevents duplicates
    // ... other fields
})
```

### 4. Handle errors gracefully

```go
resp, err := client.CreateShipment(ctx, req)
if err != nil {
    st, ok := status.FromError(err)
    if ok {
        switch st.Code() {
        case codes.AlreadyExists:
            log.Info("Shipment already exists (idempotency)")
            // Extract existing tracking number from error message
        case codes.Unavailable:
            log.Warn("Delivery service unavailable, will retry later")
            // Queue for retry
        default:
            log.Error("Failed to create shipment: %v", err)
        }
    }
}
```

### 5. Log all API calls

```go
start := time.Now()
resp, err := client.CreateShipment(ctx, req)
duration := time.Since(start)

log.Info("CreateShipment: provider=%s, duration=%v, error=%v",
    req.ProviderCode, duration, err)
```

---

## Related Documentation

- [Delivery Service Integration Guide](./DELIVERY_SERVICE_INTEGRATION.md) - Architecture and integration points
- [Frontend Delivery Guide](./DELIVERY_FRONTEND_GUIDE.md) - Frontend components and Redux usage

---

**Last updated:** 2025-10-29
**Version:** 1.0.0
**Proto version:** v1
