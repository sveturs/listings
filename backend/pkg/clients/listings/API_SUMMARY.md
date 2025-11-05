# Listings Client API Summary

Quick reference for all available methods in the gRPC client.

## Client Initialization

```go
client, err := listings.NewClient("localhost:50051", logger)
defer client.Close()
```

## Product CRUD

### Create
```go
CreateProduct(ctx, storefrontID int64, product *pb.ProductInput) (*pb.Product, error)
```

### Update
```go
UpdateProduct(ctx, req *pb.UpdateProductRequest) (*pb.Product, error)
```

### Delete
```go
DeleteProduct(ctx, productID, storefrontID int64, hardDelete bool) (*pb.DeleteProductResponse, error)
```

## Product Variant CRUD

### Create
```go
CreateProductVariant(ctx, productID int64, variant *pb.ProductVariantInput) (*pb.ProductVariant, error)
```

### Update
```go
UpdateProductVariant(ctx, req *pb.UpdateProductVariantRequest) (*pb.ProductVariant, error)
```

### Delete
```go
DeleteProductVariant(ctx, variantID, productID int64) error
```

## Bulk Operations

### Bulk Create Products (max 1000)
```go
BulkCreateProducts(ctx, storefrontID int64, products []*pb.ProductInput) (*pb.BulkCreateProductsResponse, error)
```

### Bulk Update Products (max 1000)
```go
BulkUpdateProducts(ctx, storefrontID int64, updates []*pb.ProductUpdateInput) (*pb.BulkUpdateProductsResponse, error)
```

### Bulk Delete Products (max 1000)
```go
BulkDeleteProducts(ctx, storefrontID int64, productIDs []int64, hardDelete bool) (*pb.BulkDeleteProductsResponse, error)
```

### Bulk Create Variants (max 1000)
```go
BulkCreateProductVariants(ctx, productID int64, variants []*pb.ProductVariantInput) (*pb.BulkCreateProductVariantsResponse, error)
```

## Stock Operations (Existing)

### Check Availability
```go
CheckStockAvailability(ctx, items []*pb.StockItem) (*pb.CheckStockAvailabilityResponse, error)
```

### Decrement Stock
```go
DecrementStock(ctx, items []*pb.StockItem, orderID string) (*pb.DecrementStockResponse, error)
```

### Rollback Stock
```go
RollbackStock(ctx, items []*pb.StockItem, orderID string) error
```

## Key Features

- ✅ **Structured logging** with zerolog
- ✅ **Error context preservation** via gRPC status codes
- ✅ **Batch validation** (max 1000 items)
- ✅ **Soft/hard delete** support
- ✅ **Field mask** support for partial updates
- ✅ **Thread-safe** gRPC connection

## Common Patterns

### Single Product Creation
```go
product := &pb.ProductInput{
    Name:        "Product Name",
    Description: "Description",
    Price:       99.99,
    Currency:    "USD",
    CategoryId:  1301,
}

result, err := client.CreateProduct(ctx, storefrontID, product)
```

### Bulk Import with Error Handling
```go
resp, err := client.BulkCreateProducts(ctx, storefrontID, products)
if err != nil {
    return err
}

// Handle partial failures
for _, bulkErr := range resp.Errors {
    log.Printf("Item %d failed: %s", bulkErr.Index, bulkErr.ErrorMessage)
}
```

### Stock Management Workflow
```go
// 1. Check availability
available, _ := client.CheckStockAvailability(ctx, items)

// 2. Decrement on order
if available.AllAvailable {
    resp, _ := client.DecrementStock(ctx, items, orderID)
}

// 3. Rollback on failure
if orderFailed {
    client.RollbackStock(ctx, items, orderID)
}
```

## Error Handling

```go
import "google.golang.org/grpc/codes"
import "google.golang.org/grpc/status"

_, err := client.CreateProduct(ctx, storefrontID, product)
if err != nil {
    if st, ok := status.FromError(err); ok {
        switch st.Code() {
        case codes.AlreadyExists:
            // Handle duplicate SKU
        case codes.InvalidArgument:
            // Handle validation error
        case codes.NotFound:
            // Handle not found
        }
    }
}
```

## Best Practices

1. **Always use context with timeout**
   ```go
   ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
   defer cancel()
   ```

2. **Use field masks for updates**
   ```go
   UpdateMask: &fieldmaskpb.FieldMask{
       Paths: []string{"price", "is_active"},
   }
   ```

3. **Batch operations in chunks of 1000**
   ```go
   for i := 0; i < len(all); i += 1000 {
       batch := all[i:min(i+1000, len(all))]
       client.BulkCreateProducts(ctx, storefrontID, batch)
   }
   ```

4. **Handle partial failures gracefully**
   ```go
   if resp.FailedCount > 0 {
       for _, err := range resp.Errors {
           log.Error().Int32("index", err.Index).Msg(err.ErrorMessage)
       }
   }
   ```

## Method Summary Table

| Operation | Single | Bulk | Max Items |
|-----------|--------|------|-----------|
| Create Product | ✅ | ✅ | 1000 |
| Update Product | ✅ | ✅ | 1000 |
| Delete Product | ✅ | ✅ | 1000 |
| Create Variant | ✅ | ✅ | 1000 |
| Update Variant | ✅ | ❌ | - |
| Delete Variant | ✅ | ❌ | - |
| Stock Check | ✅ | - | - |
| Stock Decrement | ✅ | - | - |
| Stock Rollback | ✅ | - | - |

## Response Types

### Single Operations
- `*pb.Product` - created/updated product
- `*pb.ProductVariant` - created/updated variant
- `*pb.DeleteProductResponse` - deletion confirmation with cascade info

### Bulk Operations
- `*pb.BulkCreateProductsResponse` - products + error details
- `*pb.BulkUpdateProductsResponse` - updated products + error details
- `*pb.BulkDeleteProductsResponse` - counts + error details
- `*pb.BulkCreateProductVariantsResponse` - variants + error details

All bulk responses include:
- `SuccessfulCount int32`
- `FailedCount int32`
- `Errors []*BulkOperationError` (index, error_code, error_message)

## See Also

- [Full Documentation](README.md)
- [Proto Definitions](../../proto/listings/v1/listings.proto)
