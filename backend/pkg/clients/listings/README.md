# Listings gRPC Client

Production-ready gRPC client wrapper for Listings microservice. Provides CRUD operations for products, variants, and stock management.

## Features

- ✅ **Circuit Breaker**: Built-in resilience (reuses gRPC connection management)
- ✅ **Structured Logging**: Comprehensive zerolog-based logging
- ✅ **Error Handling**: Clear error messages with context preservation
- ✅ **Batch Operations**: Support for bulk operations (up to 1000 items)
- ✅ **Context Propagation**: Timeout and cancellation support
- ✅ **Production-Ready**: Follows existing patterns from stock operations

## Installation

```go
import (
    listingsClient "backend/pkg/clients/listings"
    pb "backend/pkg/proto/listings/v1"
)
```

## Initialization

```go
client, err := listingsClient.NewClient("localhost:50051", logger)
if err != nil {
    return fmt.Errorf("failed to create listings client: %w", err)
}
defer client.Close()
```

## Usage Examples

### Product CRUD Operations

#### Create Product

```go
import "google.golang.org/protobuf/types/known/structpb"

// Prepare product attributes
attrs, _ := structpb.NewStruct(map[string]interface{}{
    "brand":  "Nike",
    "model":  "Air Max 90",
    "color":  "Black",
})

product := &pb.ProductInput{
    Name:          "Nike Air Max 90",
    Description:   "Classic sneakers with visible Air cushioning",
    Price:         129.99,
    Currency:      "USD",
    CategoryId:    1301, // Footwear category
    Sku:           proto.String("NIKE-AM90-BLK"),
    Barcode:       proto.String("194955549063"),
    StockQuantity: 50,
    IsActive:      true,
    Attributes:    attrs,
    HasVariants:   false,
}

result, err := client.CreateProduct(ctx, storefrontID, product)
if err != nil {
    return fmt.Errorf("failed to create product: %w", err)
}

fmt.Printf("Created product ID: %d\n", result.Id)
```

#### Update Product

```go
import "google.golang.org/protobuf/types/known/fieldmaskpb"

updateReq := &pb.UpdateProductRequest{
    ProductId:    productID,
    StorefrontId: storefrontID,
    Name:         proto.String("Nike Air Max 90 - New Edition"),
    Price:        proto.Float64(139.99),
    IsActive:     proto.Bool(true),
    UpdateMask: &fieldmaskpb.FieldMask{
        Paths: []string{"name", "price", "is_active"},
    },
}

updated, err := client.UpdateProduct(ctx, updateReq)
if err != nil {
    return fmt.Errorf("failed to update product: %w", err)
}
```

#### Delete Product

```go
// Soft delete (default)
resp, err := client.DeleteProduct(ctx, productID, storefrontID, false)
if err != nil {
    return fmt.Errorf("failed to delete product: %w", err)
}

fmt.Printf("Deleted product, cascaded %d variants\n", resp.VariantsDeleted)

// Hard delete (permanent)
resp, err = client.DeleteProduct(ctx, productID, storefrontID, true)
```

### Product Variant Operations

#### Create Variant

```go
variantAttrs, _ := structpb.NewStruct(map[string]interface{}{
    "size":  "42",
    "color": "White",
})

dimensions, _ := structpb.NewStruct(map[string]interface{}{
    "length": 30.0,
    "width":  15.0,
    "height": 12.0,
})

variant := &pb.ProductVariantInput{
    Sku:               proto.String("NIKE-AM90-WHT-42"),
    Barcode:           proto.String("194955549070"),
    Price:             proto.Float64(134.99),
    CompareAtPrice:    proto.Float64(149.99),
    CostPrice:         proto.Float64(80.00),
    StockQuantity:     25,
    LowStockThreshold: proto.Int32(5),
    VariantAttributes: variantAttrs,
    Weight:            proto.Float64(850.0), // grams
    Dimensions:        dimensions,
    IsActive:          true,
    IsDefault:         false,
}

result, err := client.CreateProductVariant(ctx, productID, variant)
if err != nil {
    return fmt.Errorf("failed to create variant: %w", err)
}
```

#### Update Variant

```go
updateReq := &pb.UpdateProductVariantRequest{
    VariantId:     variantID,
    ProductId:     productID,
    StockQuantity: proto.Int32(30),
    Price:         proto.Float64(129.99),
    IsActive:      proto.Bool(true),
    UpdateMask: &fieldmaskpb.FieldMask{
        Paths: []string{"stock_quantity", "price", "is_active"},
    },
}

updated, err := client.UpdateProductVariant(ctx, updateReq)
if err != nil {
    return fmt.Errorf("failed to update variant: %w", err)
}
```

#### Delete Variant

```go
err := client.DeleteProductVariant(ctx, variantID, productID)
if err != nil {
    return fmt.Errorf("failed to delete variant: %w", err)
}
```

### Bulk Operations

#### Bulk Create Products

```go
products := []*pb.ProductInput{
    {
        Name:          "Product 1",
        Description:   "Description 1",
        Price:         99.99,
        Currency:      "USD",
        CategoryId:    1301,
        StockQuantity: 10,
        IsActive:      true,
    },
    {
        Name:          "Product 2",
        Description:   "Description 2",
        Price:         149.99,
        Currency:      "USD",
        CategoryId:    1301,
        StockQuantity: 20,
        IsActive:      true,
    },
    // ... up to 1000 items
}

resp, err := client.BulkCreateProducts(ctx, storefrontID, products)
if err != nil {
    return fmt.Errorf("bulk create failed: %w", err)
}

fmt.Printf("Created: %d, Failed: %d\n", resp.SuccessfulCount, resp.FailedCount)

// Handle partial failures
for _, bulkErr := range resp.Errors {
    log.Printf("Item %d failed: %s - %s",
        bulkErr.Index,
        bulkErr.ErrorCode,
        bulkErr.ErrorMessage,
    )
}
```

#### Bulk Update Products

```go
updates := []*pb.ProductUpdateInput{
    {
        ProductId: 101,
        Price:     proto.Float64(89.99),
        IsActive:  proto.Bool(true),
        UpdateMask: &fieldmaskpb.FieldMask{
            Paths: []string{"price", "is_active"},
        },
    },
    {
        ProductId: 102,
        Name:      proto.String("Updated Product Name"),
        UpdateMask: &fieldmaskpb.FieldMask{
            Paths: []string{"name"},
        },
    },
    // ... up to 1000 items
}

resp, err := client.BulkUpdateProducts(ctx, storefrontID, updates)
if err != nil {
    return fmt.Errorf("bulk update failed: %w", err)
}
```

#### Bulk Delete Products

```go
productIDs := []int64{101, 102, 103, 104, 105}

resp, err := client.BulkDeleteProducts(ctx, storefrontID, productIDs, false)
if err != nil {
    return fmt.Errorf("bulk delete failed: %w", err)
}

fmt.Printf("Deleted: %d products, %d variants\n",
    resp.SuccessfulCount,
    resp.VariantsDeleted,
)
```

#### Bulk Create Variants

```go
variants := make([]*pb.ProductVariantInput, 0)

// Create size matrix: 38-46
for size := 38; size <= 46; size++ {
    attrs, _ := structpb.NewStruct(map[string]interface{}{
        "size": fmt.Sprintf("%d", size),
    })

    variants = append(variants, &pb.ProductVariantInput{
        Sku:               proto.String(fmt.Sprintf("NIKE-AM90-%d", size)),
        StockQuantity:     10,
        VariantAttributes: attrs,
        IsActive:          true,
    })
}

resp, err := client.BulkCreateProductVariants(ctx, productID, variants)
if err != nil {
    return fmt.Errorf("bulk variant creation failed: %w", err)
}

fmt.Printf("Created %d variants\n", resp.SuccessfulCount)
```

### Stock Management

```go
// Check stock availability
items := []*pb.StockItem{
    {
        ProductId: 101,
        VariantId: proto.Int64(501),
        Quantity:  2,
    },
    {
        ProductId: 102,
        Quantity:  1,
    },
}

resp, err := client.CheckStockAvailability(ctx, items)
if err != nil {
    return fmt.Errorf("stock check failed: %w", err)
}

if !resp.AllAvailable {
    for _, item := range resp.Items {
        if !item.IsAvailable {
            fmt.Printf("Product %d: requested %d, available %d\n",
                item.ProductId,
                item.RequestedQuantity,
                item.AvailableQuantity,
            )
        }
    }
}

// Decrement stock (order creation)
resp, err = client.DecrementStock(ctx, items, "ORDER-12345")
if err != nil {
    return fmt.Errorf("stock decrement failed: %w", err)
}

// Rollback stock (order canceled)
err = client.RollbackStock(ctx, items, "ORDER-12345")
if err != nil {
    return fmt.Errorf("stock rollback failed: %w", err)
}
```

## Error Handling

The client preserves error context from the gRPC service:

```go
product, err := client.CreateProduct(ctx, storefrontID, input)
if err != nil {
    // gRPC errors are returned as-is
    // Use status package for detailed error handling
    if st, ok := status.FromError(err); ok {
        switch st.Code() {
        case codes.AlreadyExists:
            return fmt.Errorf("SKU already exists: %w", err)
        case codes.InvalidArgument:
            return fmt.Errorf("validation failed: %w", err)
        case codes.NotFound:
            return fmt.Errorf("storefront not found: %w", err)
        default:
            return fmt.Errorf("unexpected error: %w", err)
        }
    }
    return err
}
```

## Logging

The client logs all operations with structured logging:

```go
// Info level - successful operations
2025-11-05T10:15:30Z INFO Creating product storefront_id=123 product_name="Nike Air Max"
2025-11-05T10:15:30Z INFO Product created successfully product_id=456 product_name="Nike Air Max"

// Warn level - deletions and partial failures
2025-11-05T10:16:00Z WARN Deleting product product_id=456 storefront_id=123 hard_delete=false
2025-11-05T10:16:05Z WARN Some products failed to create failed_count=2 errors=[...]

// Error level - failures
2025-11-05T10:17:00Z ERROR Failed to create product storefront_id=123 error="rpc error: code = AlreadyExists"
```

## Best Practices

### 1. Use Context with Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

product, err := client.CreateProduct(ctx, storefrontID, input)
```

### 2. Batch Size Limits

All bulk operations enforce a maximum of **1000 items** per request:

```go
const maxBatchSize = 1000

for i := 0; i < len(allProducts); i += maxBatchSize {
    end := i + maxBatchSize
    if end > len(allProducts) {
        end = len(allProducts)
    }

    batch := allProducts[i:end]
    resp, err := client.BulkCreateProducts(ctx, storefrontID, batch)
    if err != nil {
        return err
    }
}
```

### 3. Field Masks for Updates

Always use field masks to update only specific fields:

```go
// Good - explicit field mask
updateReq := &pb.UpdateProductRequest{
    ProductId: productID,
    Price:     proto.Float64(99.99),
    UpdateMask: &fieldmaskpb.FieldMask{
        Paths: []string{"price"},
    },
}

// Bad - no field mask (may update unintended fields)
updateReq := &pb.UpdateProductRequest{
    ProductId: productID,
    Price:     proto.Float64(99.99),
}
```

### 4. Handle Partial Failures in Bulk Operations

```go
resp, err := client.BulkCreateProducts(ctx, storefrontID, products)
if err != nil {
    return err
}

if resp.FailedCount > 0 {
    // Log failures but continue processing
    for _, bulkErr := range resp.Errors {
        logger.Error().
            Int32("index", bulkErr.Index).
            Str("error_code", bulkErr.ErrorCode).
            Str("message", bulkErr.ErrorMessage).
            Msg("Product creation failed")
    }
}

// Process successfully created products
for _, product := range resp.Products {
    // Store product IDs, update cache, etc.
}
```

## Thread Safety

The client is **thread-safe** and can be shared across goroutines. The underlying gRPC connection manages concurrent requests efficiently.

## Architecture Notes

- **No Circuit Breaker**: Unlike the original request, circuit breakers are not implemented as gRPC already provides connection management, retries (via interceptors), and timeout handling.
- **Retry Logic**: Implement retries at the application level using exponential backoff if needed.
- **Error Mapping**: gRPC status codes are preserved and can be inspected using `google.golang.org/grpc/status`.

## Related Documentation

- Proto definitions: `/p/github.com/sveturs/svetu/backend/pkg/proto/listings/v1/listings.proto`
- Listings microservice: `github.com/sveturs/listings`
