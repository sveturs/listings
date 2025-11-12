# ReindexAll gRPC Method - Testing Guide

## Overview
The `ReindexAll` method performs full reindexing of all products to OpenSearch. This is an administrative operation used for rebuilding the search index after schema changes or data migration.

## gRPC Method Signature

```protobuf
rpc ReindexAll(ReindexAllRequest) returns (ReindexAllResponse);
```

## Request Parameters

### ReindexAllRequest
```protobuf
message ReindexAllRequest {
  optional string source_type = 1; // Filter: "c2c", "b2c", or empty for all
  optional int32 batch_size = 2;   // Batch size for processing (default: 1000)
}
```

- **source_type** (optional): Filter products by source type
  - `"b2c"` - Index only B2C products
  - `"c2c"` - Index only C2C listings
  - Empty/null - Index all products (both B2C and C2C)

- **batch_size** (optional): Number of products to process in each batch
  - Default: 1000
  - Range: 1-10000
  - Invalid values will be clamped to default

## Response Format

### ReindexAllResponse
```protobuf
message ReindexAllResponse {
  int32 total_indexed = 1;     // Total successfully indexed products
  int32 total_failed = 2;      // Total failed products
  int32 duration_seconds = 3;  // Total operation duration in seconds
  repeated string errors = 4;  // Sample error messages (max 10)
}
```

## Example Usage

### 1. Using grpcurl (Command Line)

#### Reindex ALL products (B2C + C2C)
```bash
grpcurl -plaintext \
  -d '{}' \
  localhost:50051 \
  listings.v1.ListingsService/ReindexAll
```

#### Reindex only B2C products
```bash
grpcurl -plaintext \
  -d '{"source_type": "b2c"}' \
  localhost:50051 \
  listings.v1.ListingsService/ReindexAll
```

#### Reindex with custom batch size
```bash
grpcurl -plaintext \
  -d '{"source_type": "b2c", "batch_size": 500}' \
  localhost:50051 \
  listings.v1.ListingsService/ReindexAll
```

### 2. Using Go Client

```go
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sveturs/listings/api/proto/listings/v1"
)

func main() {
	// Connect to gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewListingsServiceClient(conn)

	// Create context with timeout (reindexing can take a while)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	// Example 1: Reindex all products
	req := &pb.ReindexAllRequest{}

	resp, err := client.ReindexAll(ctx, req)
	if err != nil {
		log.Fatalf("ReindexAll failed: %v", err)
	}

	fmt.Printf("Reindexing completed:\n")
	fmt.Printf("  Indexed: %d\n", resp.TotalIndexed)
	fmt.Printf("  Failed: %d\n", resp.TotalFailed)
	fmt.Printf("  Duration: %d seconds\n", resp.DurationSeconds)
	if len(resp.Errors) > 0 {
		fmt.Printf("  Sample errors:\n")
		for _, errMsg := range resp.Errors {
			fmt.Printf("    - %s\n", errMsg)
		}
	}
}
```

### 3. Using BloomRPC or Postman

#### Request JSON
```json
{
  "source_type": "b2c",
  "batch_size": 1000
}
```

## Error Handling

### Possible Error Codes

| gRPC Code | Description | Solution |
|-----------|-------------|----------|
| `FAILED_PRECONDITION` | OpenSearch indexer is not configured | Check that OpenSearch is properly configured in the service |
| `INVALID_ARGUMENT` | Invalid source_type parameter | Use "b2c", "c2c", or leave empty |
| `CANCELED` | Operation was cancelled by client | Retry or increase timeout |
| `DEADLINE_EXCEEDED` | Operation timed out | Increase context timeout or reduce batch size |
| `INTERNAL` | Generic internal error | Check logs for details |

### Example Error Response
```json
{
  "code": 9,
  "message": "OpenSearch indexer is not configured",
  "details": []
}
```

## Performance Considerations

### Recommended Settings

| Database Size | Batch Size | Expected Duration | Recommended Timeout |
|---------------|------------|-------------------|---------------------|
| < 10,000 products | 1000 | 1-2 minutes | 5 minutes |
| 10,000 - 100,000 | 1000 | 5-15 minutes | 20 minutes |
| > 100,000 | 500-1000 | 15+ minutes | 30+ minutes |

### Tips for Large Datasets

1. **Use smaller batch sizes** (500) to reduce memory pressure
2. **Increase gRPC timeout** to account for large datasets
3. **Monitor logs** to track progress during reindexing
4. **Run during off-peak hours** to minimize impact on production

## Monitoring

### Log Messages

The operation logs progress at INFO level:
```
INFO starting full reindexing source_type=b2c batch_size=1000
DEBUG fetching batch offset=0 batch_size=1000
DEBUG fetched batch count=1000 total=5000
INFO batch indexed indexed_so_far=1000 failed_so_far=0 batch_size=1000
INFO reindexing completed total_indexed=5000 total_failed=0 duration_seconds=45
```

### Metrics to Monitor

- **total_indexed**: Number of successfully indexed products
- **total_failed**: Number of failed products (should be 0 in normal operation)
- **duration_seconds**: Total time taken
- **errors**: Sample error messages (first 10 errors)

## Troubleshooting

### Issue: "OpenSearch indexer is not configured"
**Solution**: Ensure OpenSearch client is properly initialized in the service constructor.

### Issue: Operation times out
**Solutions**:
1. Increase context timeout in client
2. Reduce batch_size parameter (e.g., 500 instead of 1000)
3. Check OpenSearch cluster health and performance

### Issue: High failure rate (total_failed > 0)
**Solutions**:
1. Check error messages in response
2. Verify OpenSearch index mapping is correct
3. Check for data integrity issues in PostgreSQL
4. Review logs for detailed error messages

## Production Usage

### Pre-Reindexing Checklist
- [ ] Backup OpenSearch index
- [ ] Check OpenSearch cluster health (yellow/green)
- [ ] Verify sufficient disk space
- [ ] Schedule during maintenance window
- [ ] Set up monitoring/alerting

### Post-Reindexing Verification
```bash
# Check index document count
curl -X GET "localhost:9200/marketplace_listings/_count" | jq '.'

# Verify search functionality
grpcurl -plaintext \
  -d '{"query": "test", "limit": 10, "offset": 0}' \
  localhost:50051 \
  listings.v1.ListingsService/SearchListings
```

## See Also

- [OpenSearch Indexing Architecture](../docs/opensearch-integration.md)
- [gRPC API Documentation](../api/proto/listings/v1/listings.proto)
- [Service Layer Implementation](../internal/service/listings/service.go)
