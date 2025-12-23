# Category Detection gRPC Integration Example

## gRPC Handler пример

```go
// internal/grpc/handlers/category_detection_handler.go
package handlers

import (
    "context"

    "github.com/google/uuid"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    
    pb "github.com/vondi-global/listings/api/proto/categories/v2"
    "github.com/vondi-global/listings/internal/domain"
    "github.com/vondi-global/listings/internal/service"
)

type CategoryDetectionHandler struct {
    pb.UnimplementedCategoryDetectionServiceServer
    detectionService *service.CategoryDetectionService
}

func NewCategoryDetectionHandler(svc *service.CategoryDetectionService) *CategoryDetectionHandler {
    return &CategoryDetectionHandler{
        detectionService: svc,
    }
}

// DetectFromText реализует gRPC метод
func (h *CategoryDetectionHandler) DetectFromText(
    ctx context.Context,
    req *pb.DetectFromTextRequest,
) (*pb.DetectFromTextResponse, error) {
    // Валидация
    if req.Title == "" && req.Description == "" {
        return nil, status.Error(codes.InvalidArgument, "title or description required")
    }

    if req.Language == "" {
        req.Language = "en" // default
    }

    // Подготовка input
    input := domain.DetectFromTextInput{
        Title:       req.Title,
        Description: req.Description,
        Language:    req.Language,
    }

    // Hints (опционально)
    if req.Hints != nil {
        input.Hints = &domain.CategoryHints{
            Domain:      req.Hints.Domain,
            ProductType: req.Hints.ProductType,
            Keywords:    req.Hints.Keywords,
        }
    }

    // Детекция
    detection, err := h.detectionService.DetectFromText(ctx, input)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "detection failed: %v", err)
    }

    // Конвертация в proto
    return h.convertToProto(detection), nil
}

// DetectBatch реализует batch детекцию
func (h *CategoryDetectionHandler) DetectBatch(
    ctx context.Context,
    req *pb.DetectBatchRequest,
) (*pb.DetectBatchResponse, error) {
    if len(req.Items) == 0 {
        return nil, status.Error(codes.InvalidArgument, "items required")
    }

    // Конвертация proto → domain
    input := domain.DetectBatchInput{
        Items: make([]domain.DetectFromTextInput, len(req.Items)),
    }

    for i, item := range req.Items {
        input.Items[i] = domain.DetectFromTextInput{
            Title:       item.Title,
            Description: item.Description,
            Language:    item.Language,
        }
        if item.Hints != nil {
            input.Items[i].Hints = &domain.CategoryHints{
                Domain:      item.Hints.Domain,
                ProductType: item.Hints.ProductType,
                Keywords:    item.Hints.Keywords,
            }
        }
    }

    // Batch детекция
    result, err := h.detectionService.DetectBatch(ctx, input)
    if err != nil {
        return nil, status.Errorf(codes.Internal, "batch detection failed: %v", err)
    }

    // Конвертация в proto
    resp := &pb.DetectBatchResponse{
        Results:             make([]*pb.CategoryDetection, len(result.Results)),
        TotalProcessingTime: result.TotalProcessingTime,
    }

    for i, detection := range result.Results {
        resp.Results[i] = h.convertToProto(&detection)
    }

    return resp, nil
}

// ConfirmSelection подтверждает выбор пользователя
func (h *CategoryDetectionHandler) ConfirmSelection(
    ctx context.Context,
    req *pb.ConfirmSelectionRequest,
) (*pb.ConfirmSelectionResponse, error) {
    detectionID, err := uuid.Parse(req.DetectionId)
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, "invalid detection_id")
    }

    categoryID, err := uuid.Parse(req.SelectedCategoryId)
    if err != nil {
        return nil, status.Error(codes.InvalidArgument, "invalid selected_category_id")
    }

    if err := h.detectionService.ConfirmSelection(ctx, detectionID, categoryID); err != nil {
        return nil, status.Errorf(codes.Internal, "confirm failed: %v", err)
    }

    return &pb.ConfirmSelectionResponse{Success: true}, nil
}

// convertToProto конвертирует domain → proto
func (h *CategoryDetectionHandler) convertToProto(d *domain.CategoryDetection) *pb.CategoryDetection {
    resp := &pb.CategoryDetection{
        Id:               d.ID.String(),
        ProcessingTimeMs: d.ProcessingTimeMs,
        InputTitle:       d.InputTitle,
        InputDescription: d.InputDescription,
        InputLanguage:    d.InputLanguage,
    }

    if d.Primary != nil {
        resp.Primary = h.convertMatchToProto(d.Primary)
    }

    resp.Alternatives = make([]*pb.CategoryMatch, len(d.Alternatives))
    for i, alt := range d.Alternatives {
        resp.Alternatives[i] = h.convertMatchToProto(&alt)
    }

    return resp
}

func (h *CategoryDetectionHandler) convertMatchToProto(m *domain.CategoryMatch) *pb.CategoryMatch {
    return &pb.CategoryMatch{
        CategoryId:      m.CategoryID.String(),
        CategoryName:    m.CategoryName,
        CategorySlug:    m.CategorySlug,
        CategoryPath:    m.CategoryPath,
        ConfidenceScore: float32(m.ConfidenceScore),
        DetectionMethod: string(m.DetectionMethod),
        MatchedKeywords: m.MatchedKeywords,
    }
}
```

## Регистрация в gRPC сервере

```go
// cmd/server/main.go или internal/grpc/server.go

import (
    "github.com/vondi-global/listings/internal/grpc/handlers"
    pb "github.com/vondi-global/listings/api/proto/categories/v2"
)

// В функции инициализации gRPC сервера
func setupGRPC(
    grpcServer *grpc.Server,
    detectionService *service.CategoryDetectionService,
) {
    // Регистрируем CategoryDetection handler
    detectionHandler := handlers.NewCategoryDetectionHandler(detectionService)
    pb.RegisterCategoryDetectionServiceServer(grpcServer, detectionHandler)
}
```

## Пример вызова из монолита

```go
// В монолите (vondi/backend)
import (
    pb "github.com/vondi-global/listings/api/proto/categories/v2"
    "google.golang.org/grpc"
)

// Подключение к Listings микросервису
conn, err := grpc.Dial(
    "localhost:50053",
    grpc.WithInsecure(),
    grpc.WithTimeout(10*time.Second),
)
if err != nil {
    log.Fatal("Failed to connect to listings service", err)
}
defer conn.Close()

client := pb.NewCategoryDetectionServiceClient(conn)

// Детекция категории
resp, err := client.DetectFromText(context.Background(), &pb.DetectFromTextRequest{
    Title:       "iPhone 15 Pro 256GB",
    Description: "Новый смартфон Apple iPhone 15 Pro",
    Language:    "ru",
})

if err != nil {
    log.Error("Detection failed", err)
    return
}

log.Info("Detected category",
    "category", resp.Primary.CategoryName,
    "confidence", resp.Primary.ConfidenceScore,
    "method", resp.Primary.DetectionMethod,
)

// Альтернативы
for _, alt := range resp.Alternatives {
    log.Debug("Alternative",
        "category", alt.CategoryName,
        "confidence", alt.ConfidenceScore,
    )
}
```

## Proto Definition

```protobuf
// api/proto/categories/v2/category_detection.proto

syntax = "proto3";

package categories.v2;

option go_package = "github.com/vondi-global/listings/api/proto/categories/v2;categoriesv2";

service CategoryDetectionService {
  rpc DetectFromText(DetectFromTextRequest) returns (DetectFromTextResponse);
  rpc DetectBatch(DetectBatchRequest) returns (DetectBatchResponse);
  rpc ConfirmSelection(ConfirmSelectionRequest) returns (ConfirmSelectionResponse);
}

message DetectFromTextRequest {
  string title = 1;
  string description = 2;
  string language = 3;
  CategoryHints hints = 4;
}

message DetectFromTextResponse {
  CategoryDetection detection = 1;
}

message CategoryHints {
  string domain = 1;
  string product_type = 2;
  repeated string keywords = 3;
}

message CategoryDetection {
  string id = 1;
  CategoryMatch primary = 2;
  repeated CategoryMatch alternatives = 3;
  int32 processing_time_ms = 4;
  string input_title = 5;
  string input_description = 6;
  string input_language = 7;
}

message CategoryMatch {
  string category_id = 1;
  string category_name = 2;
  string category_slug = 3;
  string category_path = 4;
  float confidence_score = 5;
  string detection_method = 6;
  repeated string matched_keywords = 7;
}

message DetectBatchRequest {
  repeated DetectFromTextRequest items = 1;
}

message DetectBatchResponse {
  repeated CategoryDetection results = 1;
  int32 total_processing_time = 2;
}

message ConfirmSelectionRequest {
  string detection_id = 1;
  string selected_category_id = 2;
}

message ConfirmSelectionResponse {
  bool success = 1;
}
```

## Testing

```go
// internal/grpc/handlers/category_detection_handler_test.go

func TestDetectFromText(t *testing.T) {
    // Mock service
    mockService := &mockCategoryDetectionService{
        detectFromTextFn: func(ctx context.Context, input domain.DetectFromTextInput) (*domain.CategoryDetection, error) {
            return &domain.CategoryDetection{
                ID: uuid.New(),
                Primary: &domain.CategoryMatch{
                    CategorySlug:    "electronics",
                    CategoryName:    "Electronics",
                    ConfidenceScore: 0.95,
                    DetectionMethod: domain.MethodAIClaude,
                },
                ProcessingTimeMs: 3500,
            }, nil
        },
    }

    handler := NewCategoryDetectionHandler(mockService)

    resp, err := handler.DetectFromText(context.Background(), &pb.DetectFromTextRequest{
        Title:    "iPhone 15 Pro",
        Language: "en",
    })

    require.NoError(t, err)
    require.NotNil(t, resp.Detection.Primary)
    assert.Equal(t, "electronics", resp.Detection.Primary.CategorySlug)
    assert.Greater(t, resp.Detection.Primary.ConfidenceScore, float32(0.9))
}
```
