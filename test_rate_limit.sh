#!/bin/bash
# Rate Limit Load Testing Script
# Tests the rate limiter by making rapid requests to the listings service

set -e

GRPC_HOST="localhost:50051"
LIMIT=200  # GetListing limit (from config)
EXTRA_REQUESTS=50  # Additional requests to test blocking

echo "================================================"
echo "Rate Limit Load Test"
echo "================================================"
echo "Target: $GRPC_HOST"
echo "Expected limit: $LIMIT requests/min"
echo "Test plan: Send $(($LIMIT + $EXTRA_REQUESTS)) requests rapidly"
echo ""

# Check if grpcurl is available
if ! command -v grpcurl &> /dev/null; then
    echo "ERROR: grpcurl not found. Please install it:"
    echo "  go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest"
    exit 1
fi

# Check if the service is running
echo "Checking if service is running..."
if ! grpcurl -plaintext $GRPC_HOST list &> /dev/null; then
    echo "ERROR: Cannot connect to gRPC service at $GRPC_HOST"
    echo "Please start the listings service first."
    exit 1
fi
echo "✓ Service is running"
echo ""

# Prepare test data
LISTING_ID=328

# Counters
SUCCESS_COUNT=0
RATE_LIMIT_COUNT=0
ERROR_COUNT=0

echo "Starting load test..."
echo "Sending $(($LIMIT + $EXTRA_REQUESTS)) requests..."
echo ""

# Make rapid requests
START_TIME=$(date +%s)
for i in $(seq 1 $(($LIMIT + $EXTRA_REQUESTS))); do
    RESPONSE=$(grpcurl -plaintext \
        -d "{\"id\": $LISTING_ID}" \
        $GRPC_HOST \
        listings.v1.ListingsService/GetListing 2>&1)

    if echo "$RESPONSE" | grep -q "ResourceExhausted"; then
        RATE_LIMIT_COUNT=$((RATE_LIMIT_COUNT + 1))
        if [ $i -eq $(($LIMIT + 1)) ]; then
            echo "✓ Rate limit triggered at request $i (expected at ~$LIMIT)"
        fi
    elif echo "$RESPONSE" | grep -q "\"listing\""; then
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    else
        ERROR_COUNT=$((ERROR_COUNT + 1))
    fi

    # Show progress every 50 requests
    if [ $((i % 50)) -eq 0 ]; then
        echo "Progress: $i/$(($LIMIT + $EXTRA_REQUESTS)) requests sent..."
    fi
done
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo ""
echo "================================================"
echo "Test Results"
echo "================================================"
echo "Total requests:       $(($LIMIT + $EXTRA_REQUESTS))"
echo "Successful:           $SUCCESS_COUNT"
echo "Rate limited:         $RATE_LIMIT_COUNT"
echo "Errors:               $ERROR_COUNT"
echo "Duration:             ${DURATION}s"
echo ""

# Validation
echo "================================================"
echo "Validation"
echo "================================================"

if [ $RATE_LIMIT_COUNT -gt 0 ]; then
    echo "✓ Rate limiting is ACTIVE"

    # Check if roughly the right number were allowed
    ALLOWED_RANGE_MIN=$((LIMIT - 10))
    ALLOWED_RANGE_MAX=$((LIMIT + 10))

    if [ $SUCCESS_COUNT -ge $ALLOWED_RANGE_MIN ] && [ $SUCCESS_COUNT -le $ALLOWED_RANGE_MAX ]; then
        echo "✓ Rate limit working correctly (~$SUCCESS_COUNT allowed out of $LIMIT limit)"
    else
        echo "⚠ Warning: Expected ~$LIMIT successful requests, got $SUCCESS_COUNT"
    fi

    EXPECTED_BLOCKED=$EXTRA_REQUESTS
    if [ $RATE_LIMIT_COUNT -ge $((EXPECTED_BLOCKED - 10)) ] && [ $RATE_LIMIT_COUNT -le $((EXPECTED_BLOCKED + 10)) ]; then
        echo "✓ Correct number of requests blocked (~$RATE_LIMIT_COUNT blocked)"
    else
        echo "⚠ Warning: Expected ~$EXPECTED_BLOCKED blocked requests, got $RATE_LIMIT_COUNT"
    fi
else
    echo "✗ FAIL: No rate limiting detected!"
    echo "   All requests succeeded, rate limiter may not be active"
fi

if [ $ERROR_COUNT -gt 0 ]; then
    echo "⚠ Warning: $ERROR_COUNT requests failed with errors (not rate limits)"
fi

echo ""
echo "Rate limit test completed!"
echo ""
echo "Next steps:"
echo "1. Check service logs for rate limit messages"
echo "2. Check Prometheus metrics at http://localhost:9090/metrics"
echo "   - listings_rate_limit_hits_total"
echo "   - listings_rate_limit_allowed_total"
echo "   - listings_rate_limit_rejected_total"
echo "3. Check Redis keys: redis-cli KEYS 'rate_limit:*'"
