#!/bin/bash
# Verify gRPC client initialization in server.go

echo "🔍 Checking gRPC client initialization..."

# Check if import is added
if grep -q "listingsClient.*internal/clients/listings" /p/github.com/sveturs/svetu/backend/internal/server/server.go; then
    echo "✅ Import statement found"
else
    echo "❌ Import statement MISSING"
    exit 1
fi

# Check if client is initialized
if grep -q "listingsClient.NewClient" /p/github.com/sveturs/svetu/backend/internal/server/server.go; then
    echo "✅ gRPC client initialization found"
else
    echo "❌ gRPC client initialization MISSING"
    exit 1
fi

# Check if adapter is created
if grep -q "listingsClient.NewClientAdapter" /p/github.com/sveturs/svetu/backend/internal/server/server.go; then
    echo "✅ gRPC client adapter creation found"
else
    echo "❌ gRPC client adapter creation MISSING"
    exit 1
fi

# Check if client is set in service
if grep -q "SetListingsGRPCClient" /p/github.com/sveturs/svetu/backend/internal/server/server.go; then
    echo "✅ SetListingsGRPCClient() call found"
else
    echo "❌ SetListingsGRPCClient() call MISSING"
    exit 1
fi

# Check if ClientAdapter file exists
if [ -f "/p/github.com/sveturs/svetu/backend/internal/clients/listings/client_adapter.go" ]; then
    echo "✅ ClientAdapter file exists"
else
    echo "❌ ClientAdapter file MISSING"
    exit 1
fi

# Check if code compiles
echo "🔨 Checking compilation..."
cd /p/github.com/sveturs/svetu/backend && go build -o /tmp/backend-verify-test ./cmd/api/main.go 2>&1 > /tmp/compile_check.log
if [ $? -eq 0 ]; then
    echo "✅ Code compiles successfully"
    rm -f /tmp/backend-verify-test
else
    echo "❌ Compilation FAILED"
    echo "Error log:"
    cat /tmp/compile_check.log
    exit 1
fi

echo ""
echo "✅ All gRPC client initialization checks passed!"
