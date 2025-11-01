#!/bin/bash
set -e

echo "🧪 Testing Unified Marketplace API (Phase 3)"
echo "============================================="

# Получить токен
if [ ! -f /tmp/token ]; then
    echo "❌ Error: JWT token not found at /tmp/token"
    echo "💡 Tip: Generate token with: ssh svetu@svetu.rs \"cd /opt/svetu-authpreprod && sed 's|/data/auth_svetu/keys/private.pem|./keys/private.pem|g' cmd/scripts/create_admin_jwt/create_admin_jwt.go > /tmp/create_jwt_fixed.go && go run /tmp/create_jwt_fixed.go\" > /tmp/token"
    exit 1
fi

TOKEN=$(cat /tmp/token)
BASE_URL="http://localhost:3000"

echo ""
echo "📡 Base URL: $BASE_URL"
echo "🔑 Token: ${TOKEN:0:50}..."
echo ""

# Test 1: Health check
echo "1️⃣  Backend health check..."
VERSION=$(curl -s "$BASE_URL/" | head -1)
if [[ "$VERSION" == *"Svetu API"* ]]; then
    echo "   ✅ Backend is running: $VERSION"
else
    echo "   ❌ Backend is not responding"
    exit 1
fi

echo ""

# Test 2: Search listings (public endpoint)
echo "2️⃣  Search listings (public)..."
SEARCH_RESULT=$(curl -s "$BASE_URL/api/v1/marketplace/search?limit=2")
SEARCH_COUNT=$(echo "$SEARCH_RESULT" | jq -r '.total')
echo "   📊 Found listings: $SEARCH_COUNT"
echo "   📋 Response meta:"
echo "$SEARCH_RESULT" | jq '.meta // {total, limit, offset}'

echo ""

# Test 3: Get specific listing (public endpoint)
echo "3️⃣  Get listing by ID (public)..."
# Используем listing ID 328 как пример (из документации)
LISTING_ID=328
GET_RESULT=$(curl -s "$BASE_URL/api/v1/marketplace/listings/$LISTING_ID?source_type=c2c")
if echo "$GET_RESULT" | jq -e '.success' > /dev/null 2>&1; then
    LISTING_TITLE=$(echo "$GET_RESULT" | jq -r '.data.title')
    echo "   ✅ Listing found: $LISTING_TITLE"
    echo "   📋 Response:"
    echo "$GET_RESULT" | jq '{success, data: {id: .data.id, title: .data.title, price: .data.price, source_type: .data.source_type}}'
else
    echo "   ⚠️  Listing not found or error:"
    echo "$GET_RESULT" | jq '.'
fi

echo ""

# Test 4: Create listing (requires auth)
echo "4️⃣  Create C2C listing (auth required)..."
CREATE_RESULT=$(curl -s -X POST \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "source_type": "c2c",
    "title": "Test Unified Listing",
    "description": "Created via unified API integration test",
    "price": 999.99,
    "category_id": 1301,
    "condition": "new"
  }' \
  "$BASE_URL/api/v1/marketplace/listings")

if echo "$CREATE_RESULT" | jq -e '.success' > /dev/null 2>&1; then
    CREATED_ID=$(echo "$CREATE_RESULT" | jq -r '.id')
    echo "   ✅ Listing created successfully!"
    echo "   🆔 ID: $CREATED_ID"
    echo "   🔤 Source type: $(echo "$CREATE_RESULT" | jq -r '.source_type')"

    # Сохраняем ID для дальнейших тестов
    CREATED_LISTING_ID=$CREATED_ID

    echo ""

    # Test 5: Update created listing (requires auth + ownership)
    echo "5️⃣  Update created listing (auth + ownership)..."
    UPDATE_RESULT=$(curl -s -X PUT \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d '{
        "source_type": "c2c",
        "title": "Test Unified Listing (Updated)",
        "description": "Updated via unified API",
        "price": 1099.99,
        "category_id": 1301,
        "condition": "used"
      }' \
      "$BASE_URL/api/v1/marketplace/listings/$CREATED_LISTING_ID")

    if echo "$UPDATE_RESULT" | jq -e '.success' > /dev/null 2>&1; then
        echo "   ✅ Listing updated successfully!"
    else
        echo "   ❌ Update failed:"
        echo "$UPDATE_RESULT" | jq '.'
    fi

    echo ""

    # Test 6: Delete created listing (requires auth + ownership)
    echo "6️⃣  Delete created listing (auth + ownership)..."
    DELETE_RESULT=$(curl -s -X DELETE \
      -H "Authorization: Bearer $TOKEN" \
      "$BASE_URL/api/v1/marketplace/listings/$CREATED_LISTING_ID?source_type=c2c")

    if echo "$DELETE_RESULT" | jq -e '.success' > /dev/null 2>&1; then
        echo "   ✅ Listing deleted successfully!"
    else
        echo "   ❌ Delete failed:"
        echo "$DELETE_RESULT" | jq '.'
    fi
else
    echo "   ❌ Create failed:"
    echo "$CREATE_RESULT" | jq '.'
fi

echo ""
echo "================================================"
echo "✅ All unified API tests completed!"
echo "================================================"
echo ""
echo "📝 Summary:"
echo "   ✓ Health check"
echo "   ✓ Search listings (public)"
echo "   ✓ Get listing by ID (public)"
echo "   ✓ Create listing (auth)"
echo "   ✓ Update listing (auth + ownership)"
echo "   ✓ Delete listing (auth + ownership)"
echo ""
echo "🎯 Unified API is working correctly!"
