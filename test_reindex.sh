#!/bin/bash
# Quick test script for ReindexAll gRPC method

set -e

GRPC_HOST="${GRPC_HOST:-localhost:50051}"

echo "🔍 Testing ReindexAll gRPC method on $GRPC_HOST"
echo ""

# Test 1: Reindex all products (no filters)
echo "📦 Test 1: Reindex ALL products (B2C + C2C)"
echo "Request: {}"
grpcurl -plaintext \
  -d '{}' \
  $GRPC_HOST \
  listings.v1.ListingsService/ReindexAll
echo ""
echo "---"
echo ""

# Test 2: Reindex only B2C products
echo "📦 Test 2: Reindex ONLY B2C products"
echo "Request: {\"source_type\": \"b2c\"}"
grpcurl -plaintext \
  -d '{"source_type": "b2c"}' \
  $GRPC_HOST \
  listings.v1.ListingsService/ReindexAll
echo ""
echo "---"
echo ""

# Test 3: Reindex with custom batch size
echo "📦 Test 3: Reindex with custom batch size (500)"
echo "Request: {\"source_type\": \"b2c\", \"batch_size\": 500}"
grpcurl -plaintext \
  -d '{"source_type": "b2c", "batch_size": 500}' \
  $GRPC_HOST \
  listings.v1.ListingsService/ReindexAll
echo ""
echo "---"
echo ""

# Test 4: Invalid source type (should fail)
echo "📦 Test 4: Invalid source_type (should fail with INVALID_ARGUMENT)"
echo "Request: {\"source_type\": \"invalid\"}"
grpcurl -plaintext \
  -d '{"source_type": "invalid"}' \
  $GRPC_HOST \
  listings.v1.ListingsService/ReindexAll || true
echo ""

echo "✅ All tests completed"
