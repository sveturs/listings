# Functional Tests Coverage Plan

## Current Coverage (6 tests)

### ✅ Already Implemented:
1. **Authentication Flow** (`api-auth-flow`)
   - GET `/api/v1/auth/me` - verify user authentication

2. **Marketplace CRUD** (`api-marketplace-crud`)
   - GET `/api/v1/unified/listings?limit=5` - fetch listings

3. **Categories API** (`api-categories-fetch`)
   - GET `/api/v1/admin/categories` - fetch all categories

4. **Search Functionality** (`api-search-functionality`)
   - GET `/api/v1/search?query=test&limit=5` - unified search

5. **Admin Operations** (`api-admin-operations`)
   - GET `/api/v1/admin/admins` - fetch admin list

6. **Review Creation** (`api-review-creation`)
   - GET `/api/v1/auth/me` - get current user
   - GET `/api/v1/unified/listings?limit=10` - find listing
   - GET `/api/v1/reviews?entity_type=listing&entity_id={id}&user_id={uid}` - check existing reviews
   - POST `/api/v1/reviews/draft` - create draft
   - POST `/api/v1/reviews/{id}/publish` - publish review

---

## API Coverage Analysis

### Total Endpoints: **538**

#### Module Distribution:
- **Admin** (138 endpoints) - 0.7% covered (1/138)
- **Marketplace** (67 endpoints) - 1.5% covered (1/67)
- **B2C Stores** (58 endpoints) - 0% covered
- **Post Express** (45 endpoints) - 0% covered
- **GIS** (27 endpoints) - 0% covered
- **Reviews** (18 endpoints) - 11% covered (2/18)
- **Storefront** (15 endpoints) - 0% covered
- **Search** (12 endpoints) - 8% covered (1/12)
- **Auth** (9 endpoints) - 11% covered (1/9)
- **Others** - 0% covered

**Overall Coverage: 1.1% (6/538 endpoints tested)**

---

## Priority 1: Critical User Flows (High Priority)

### 1. Marketplace Listing Creation Flow 🔴
**Why Critical**: Core business functionality - users create listings

**Endpoints to test**:
- POST `/api/v1/marketplace/listings` - create listing (draft)
- POST `/api/v1/marketplace/listings/{id}/publish` - publish listing
- GET `/api/v1/marketplace/listings/{id}` - verify created listing
- PUT `/api/v1/marketplace/listings/{id}` - update listing
- DELETE `/api/v1/marketplace/listings/{id}` - delete listing

**Test Flow**:
1. Create draft listing with minimal data
2. Publish listing
3. Verify listing appears in search/unified listings
4. Update listing (title, price)
5. Delete listing
6. Verify listing is deleted

---

### 2. B2C Storefront Product Management 🔴
**Why Critical**: Store owners manage products

**Endpoints to test**:
- POST `/api/v1/b2c_stores` - create storefront
- POST `/api/v1/b2c_stores/{id}/products` - add product
- GET `/api/v1/b2c_stores/{id}/products` - list products
- PUT `/api/v1/b2c_stores/{id}/products/{product_id}` - update product
- DELETE `/api/v1/b2c_stores/{id}/products/{product_id}` - delete product

**Test Flow**:
1. Create new storefront
2. Add product to storefront
3. List products and verify
4. Update product (price, inventory)
5. Delete product
6. Delete storefront

---

### 3. Order Creation and Management 🔴
**Why Critical**: Revenue generation - users place orders

**Endpoints to test**:
- POST `/api/v1/orders` - create order
- GET `/api/v1/orders/{id}` - get order details
- PUT `/api/v1/orders/{id}/status` - update order status
- POST `/api/v1/orders/{id}/cancel` - cancel order
- GET `/api/v1/users/orders` - list user orders

**Test Flow**:
1. Add item to cart
2. Create order from cart
3. Verify order created with correct status
4. Update order status (processing → shipped)
5. Cancel order
6. Verify refund/cancellation

---

### 4. Payment Processing 🔴
**Why Critical**: Payment flow must work perfectly

**Endpoints to test**:
- POST `/api/v1/payments/intent` - create payment intent
- POST `/api/v1/payments/confirm` - confirm payment
- GET `/api/v1/payments/{id}/status` - check payment status
- POST `/api/v1/payments/{id}/refund` - process refund

**Test Flow**:
1. Create payment intent for order
2. Confirm payment (mock payment gateway)
3. Verify payment status = success
4. Process refund
5. Verify refund status

---

### 5. Cart Management 🟡
**Why Important**: Users add items before ordering

**Endpoints to test**:
- POST `/api/v1/user/carts` - create cart
- POST `/api/v1/user/carts/{id}/items` - add item
- PUT `/api/v1/user/carts/{id}/items/{item_id}` - update quantity
- DELETE `/api/v1/user/carts/{id}/items/{item_id}` - remove item
- GET `/api/v1/user/carts` - get user carts
- DELETE `/api/v1/user/carts/{id}` - clear cart

**Test Flow**:
1. Create new cart
2. Add 3 different items
3. Update item quantity
4. Remove one item
5. Verify cart total calculation
6. Clear cart

---

### 6. User Profile Management 🟡
**Why Important**: Users manage their data

**Endpoints to test**:
- GET `/api/v1/users/profile` - get profile
- PUT `/api/v1/users/profile` - update profile
- GET `/api/v1/users/privacy-settings` - get privacy settings
- PUT `/api/v1/users/privacy-settings` - update privacy settings
- PUT `/api/v1/users/chat-settings` - update chat settings

**Test Flow**:
1. Get current profile
2. Update profile (name, phone)
3. Update privacy settings
4. Update chat settings
5. Verify all changes persisted

---

### 7. Post Express Delivery Integration 🟡
**Why Important**: Critical for logistics

**Endpoints to test**:
- POST `/api/v1/postexpress/calculate` - calculate shipping cost
- POST `/api/v1/postexpress/shipment` - create shipment
- GET `/api/v1/postexpress/shipment/{id}/track` - track shipment
- POST `/api/v1/postexpress/shipment/{id}/cancel` - cancel shipment

**Test Flow**:
1. Calculate shipping cost for order
2. Create shipment with tracking
3. Track shipment status
4. Cancel shipment
5. Verify cancellation

---

### 8. Category and Attribute Management 🟡
**Why Important**: Admin manages product categories

**Endpoints to test**:
- POST `/api/v1/admin/categories` - create category
- POST `/api/v1/admin/attributes` - create attribute
- POST `/api/v1/admin/categories/{id}/attributes/{attr_id}` - link attribute to category
- GET `/api/v1/admin/categories/{id}/attributes` - get category attributes
- DELETE `/api/v1/admin/categories/{id}/attributes/{attr_id}` - unlink attribute

**Test Flow**:
1. Create new category
2. Create attribute (e.g., "Color")
3. Link attribute to category
4. Verify attribute appears in category
5. Unlink attribute
6. Delete category

---

## Priority 2: Supporting Features (Medium Priority)

### 9. Notifications System 🟢
- POST `/api/v1/notifications/send` - send notification
- GET `/api/v1/notifications` - list user notifications
- PUT `/api/v1/notifications/{id}/read` - mark as read
- DELETE `/api/v1/notifications/{id}` - delete notification

### 10. Subscriptions Management 🟢
- POST `/api/v1/subscriptions` - create subscription
- GET `/api/v1/subscriptions` - list subscriptions
- PUT `/api/v1/subscriptions/{id}/cancel` - cancel subscription
- POST `/api/v1/subscriptions/{id}/renew` - renew subscription

### 11. Analytics and Reporting 🟢
- GET `/api/v1/analytics/sales` - sales analytics
- GET `/api/v1/analytics/users` - user analytics
- GET `/api/v1/analytics/listings` - listing analytics

### 12. GIS and Location Services 🟢
- POST `/api/v1/gis/geocode` - geocode address
- GET `/api/v1/gis/nearby` - find nearby listings
- POST `/api/v1/gis/route` - calculate route

---

## Priority 3: Admin Features (Lower Priority)

### 13. Admin User Management
- GET `/api/v1/admin/users` - list all users
- PUT `/api/v1/admin/users/{id}/status` - ban/unban user
- PUT `/api/v1/admin/users/{id}/role` - change user role

### 14. Admin Content Moderation
- GET `/api/v1/admin/listings` - list all listings
- PUT `/api/v1/admin/listings/{id}/approve` - approve listing
- PUT `/api/v1/admin/listings/{id}/reject` - reject listing

### 15. Admin Delivery Management
- GET `/api/v1/admin/delivery/problems` - list delivery problems
- POST `/api/v1/admin/delivery/problems/{id}/resolve` - resolve problem
- GET `/api/v1/admin/delivery/analytics` - delivery analytics

---

## Recommended Implementation Order

### Phase 1: Core Business (Week 1-2)
1. ✅ Authentication Flow (DONE)
2. ✅ Marketplace CRUD basics (DONE)
3. ✅ Review Creation (DONE)
4. 🔴 **Marketplace Listing Creation Flow** (HIGH PRIORITY)
5. 🔴 **Cart Management** (HIGH PRIORITY)
6. 🔴 **Order Creation and Management** (HIGH PRIORITY)

### Phase 2: Payment and Delivery (Week 3)
7. 🔴 **Payment Processing** (HIGH PRIORITY)
8. 🟡 **Post Express Delivery Integration**
9. 🟡 **B2C Storefront Product Management**

### Phase 3: User Management (Week 4)
10. 🟡 **User Profile Management**
11. 🟢 **Notifications System**
12. 🟢 **Subscriptions Management**

### Phase 4: Admin Features (Week 5)
13. 🟡 **Category and Attribute Management**
14. 🟢 **Admin User Management**
15. 🟢 **Admin Content Moderation**

### Phase 5: Advanced Features (Week 6+)
16. 🟢 **Analytics and Reporting**
17. 🟢 **GIS and Location Services**
18. 🟢 **Admin Delivery Management**

---

## Test Coverage Goals

### Short-term (1 month):
- **Target**: 20-30 functional tests
- **Coverage**: 5-6% of total endpoints
- **Focus**: Core business flows (Phases 1-2)

### Medium-term (3 months):
- **Target**: 50-70 functional tests
- **Coverage**: 10-12% of total endpoints
- **Focus**: All critical user flows (Phases 1-4)

### Long-term (6 months):
- **Target**: 100-150 functional tests
- **Coverage**: 20-30% of total endpoints
- **Focus**: All major features + admin panels (Phases 1-5)

---

## How to Write Functional Tests

### File Structure

All functional tests are in: `backend/internal/proj/admin/testing/service/functional_tests.go`

### Test Template

```go
// testYourFeature tests your feature workflow
func testYourFeature(ctx context.Context, baseURL, token string) *domain.TestResult {
    result := &domain.TestResult{
        TestName:  "api-your-feature",
        TestSuite: "api",
        Status:    domain.TestResultStatusPassed,
        StartedAt: time.Now().UTC(),
    }

    client := &http.Client{Timeout: 10 * time.Second}

    // Step 1: First API call
    req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/endpoint", nil)
    if err != nil {
        return failTest(result, "Failed to create request", err)
    }
    req.Header.Set("Authorization", "Bearer "+token)

    resp, err := client.Do(req)
    if err != nil {
        return failTest(result, "Failed to execute request", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
    }

    // Step 2: Verify response
    var response map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
        return failTest(result, "Failed to decode response", err)
    }

    // Step 3: Validate data
    data, ok := response["data"]
    if !ok {
        return failTest(result, "Missing data field", nil)
    }

    // Mark test as passed
    result.CompletedAt = time.Now().UTC()
    result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
    return result
}
```

### Step-by-Step Guide

#### 1. Add Test to Registry

In `APIEndpointTests` array:

```go
var APIEndpointTests = []FunctionalTest{
    // ... existing tests
    {
        Name:        "api-your-feature",
        Category:    domain.TestCategoryAPI,
        Description: "Test your feature description",
        RunFunc:     testYourFeature,
    },
}
```

#### 2. Create Test Function

Follow this structure:

**A. Initialize test result:**
```go
result := &domain.TestResult{
    TestName:  "api-your-feature",  // Must match Name in registry
    TestSuite: "api",               // Test suite category
    Status:    domain.TestResultStatusPassed,  // Default to passed
    StartedAt: time.Now().UTC(),   // ALWAYS use .UTC()
}
```

**B. Create HTTP client:**
```go
client := &http.Client{Timeout: 10 * time.Second}
```

**C. Make HTTP requests:**
```go
// GET request
req, err := http.NewRequestWithContext(ctx, "GET", baseURL+"/api/v1/endpoint", nil)
if err != nil {
    return failTest(result, "Failed to create request", err)
}
req.Header.Set("Authorization", "Bearer "+token)

// POST request with JSON body
payload := map[string]interface{}{
    "field1": "value1",
    "field2": 123,
}
payloadBytes, err := json.Marshal(payload)
if err != nil {
    return failTest(result, "Failed to marshal payload", err)
}
req, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/v1/endpoint", bytes.NewReader(payloadBytes))
req.Header.Set("Authorization", "Bearer "+token)
req.Header.Set("Content-Type", "application/json")
```

**D. Execute request:**
```go
resp, err := client.Do(req)
if err != nil {
    return failTest(result, "Failed to execute request", err)
}
defer resp.Body.Close()
```

**E. Check status code:**
```go
if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    return failTest(result, fmt.Sprintf("Expected status 200, got %d", resp.StatusCode), fmt.Errorf("response: %s", string(body)))
}
```

**F. Parse and validate response:**
```go
var response map[string]interface{}
if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
    return failTest(result, "Failed to decode response", err)
}

// Validate response structure
data, ok := response["data"]
if !ok {
    return failTest(result, "Missing data field in response", nil)
}

// For nested objects
dataMap, ok := data.(map[string]interface{})
if !ok {
    return failTest(result, "Data field is not an object", nil)
}

// For arrays
dataArray, ok := data.([]interface{})
if !ok {
    return failTest(result, "Data field is not an array", nil)
}
```

**G. Complete test successfully:**
```go
result.CompletedAt = time.Now().UTC()  // ALWAYS use .UTC()
result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
return result
```

#### 3. Add to Frontend UI

In `frontend/svetu/src/app/[locale]/admin/quality-tests/QualityTestsClient.tsx`:

```typescript
{
  id: 'api-your-feature',
  name: 'Your Feature Name',
  description: 'Description of what this test does',
  category: 'functional',
  icon: '🎯',  // Choose appropriate emoji
},
```

### Best Practices

#### ✅ DO:

1. **Use UTC timestamps everywhere:**
   ```go
   time.Now().UTC()  // ✅ Correct
   time.Now()        // ❌ Wrong - timezone issues
   ```

2. **Always close response bodies:**
   ```go
   defer resp.Body.Close()
   ```

3. **Use descriptive error messages:**
   ```go
   return failTest(result, "Failed to create draft review, status 400", err)
   ```

4. **Test complete workflows, not single endpoints:**
   ```go
   // ✅ Good - tests full flow
   // 1. Get user
   // 2. Find listing
   // 3. Create review
   // 4. Publish review

   // ❌ Bad - tests only one endpoint
   // 1. GET /api/v1/reviews
   ```

5. **Handle both success and error cases:**
   ```go
   // Check for null or empty data
   if dataField == nil {
       // Handle null case
   }

   if len(dataArray) == 0 {
       // Handle empty case
   }
   ```

6. **Use context for cancellation:**
   ```go
   http.NewRequestWithContext(ctx, "GET", url, nil)
   ```

#### ❌ DON'T:

1. **Don't hardcode user/listing IDs:**
   ```go
   // ❌ Bad
   listingID := 1073

   // ✅ Good - fetch dynamically
   listings := fetchListings()
   listingID := listings[0].ID
   ```

2. **Don't skip error handling:**
   ```go
   // ❌ Bad
   resp, _ := client.Do(req)

   // ✅ Good
   resp, err := client.Do(req)
   if err != nil {
       return failTest(result, "Request failed", err)
   }
   ```

3. **Don't leak response bodies:**
   ```go
   // ❌ Bad
   resp, err := client.Do(req)
   // ... forgot defer close

   // ✅ Good
   resp, err := client.Do(req)
   defer resp.Body.Close()
   ```

4. **Don't use inconsistent naming:**
   ```go
   // Test name in registry must match TestName in result
   Name: "api-your-feature",     // in APIEndpointTests
   TestName: "api-your-feature", // in result struct
   ```

### Example: Multi-Step Test

```go
func testMarketplaceListingFlow(ctx context.Context, baseURL, token string) *domain.TestResult {
    result := &domain.TestResult{
        TestName:  "api-marketplace-listing-flow",
        TestSuite: "api",
        Status:    domain.TestResultStatusPassed,
        StartedAt: time.Now().UTC(),
    }

    client := &http.Client{Timeout: 10 * time.Second}

    // Step 1: Create draft listing
    listingPayload := map[string]interface{}{
        "title":       "Test Listing",
        "price":       100,
        "category_id": 1001,
        "description": "Test description",
    }

    payloadBytes, err := json.Marshal(listingPayload)
    if err != nil {
        return failTest(result, "Failed to marshal listing payload", err)
    }

    reqCreate, err := http.NewRequestWithContext(ctx, "POST", baseURL+"/api/v1/marketplace/listings", bytes.NewReader(payloadBytes))
    if err != nil {
        return failTest(result, "Failed to create request", err)
    }
    reqCreate.Header.Set("Authorization", "Bearer "+token)
    reqCreate.Header.Set("Content-Type", "application/json")

    respCreate, err := client.Do(reqCreate)
    if err != nil {
        return failTest(result, "Failed to create listing", err)
    }
    defer respCreate.Body.Close()

    if respCreate.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(respCreate.Body)
        return failTest(result, fmt.Sprintf("Create listing failed, status %d", respCreate.StatusCode), fmt.Errorf("response: %s", string(body)))
    }

    var createResp map[string]interface{}
    if err := json.NewDecoder(respCreate.Body).Decode(&createResp); err != nil {
        return failTest(result, "Failed to decode create response", err)
    }

    listingData, ok := createResp["data"].(map[string]interface{})
    if !ok {
        return failTest(result, "Missing data in create response", nil)
    }

    listingID := int(listingData["id"].(float64))

    // Step 2: Publish listing
    reqPublish, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/api/v1/marketplace/listings/%d/publish", baseURL, listingID), nil)
    if err != nil {
        return failTest(result, "Failed to create publish request", err)
    }
    reqPublish.Header.Set("Authorization", "Bearer "+token)

    respPublish, err := client.Do(reqPublish)
    if err != nil {
        return failTest(result, "Failed to publish listing", err)
    }
    defer respPublish.Body.Close()

    if respPublish.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(respPublish.Body)
        return failTest(result, fmt.Sprintf("Publish listing failed, status %d", respPublish.StatusCode), fmt.Errorf("response: %s", string(body)))
    }

    // Step 3: Verify listing is published
    reqGet, err := http.NewRequestWithContext(ctx, "GET", fmt.Sprintf("%s/api/v1/marketplace/listings/%d", baseURL, listingID), nil)
    if err != nil {
        return failTest(result, "Failed to create get request", err)
    }
    reqGet.Header.Set("Authorization", "Bearer "+token)

    respGet, err := client.Do(reqGet)
    if err != nil {
        return failTest(result, "Failed to get listing", err)
    }
    defer respGet.Body.Close()

    if respGet.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(respGet.Body)
        return failTest(result, fmt.Sprintf("Get listing failed, status %d", respGet.StatusCode), fmt.Errorf("response: %s", string(body)))
    }

    var getResp map[string]interface{}
    if err := json.NewDecoder(respGet.Body).Decode(&getResp); err != nil {
        return failTest(result, "Failed to decode get response", err)
    }

    listing, ok := getResp["data"].(map[string]interface{})
    if !ok {
        return failTest(result, "Missing listing data", nil)
    }

    status, ok := listing["status"].(string)
    if !ok || status != "active" {
        return failTest(result, fmt.Sprintf("Listing status is %v, expected 'active'", listing["status"]), nil)
    }

    result.CompletedAt = time.Now().UTC()
    result.DurationMs = int(result.CompletedAt.Sub(result.StartedAt).Milliseconds())
    return result
}
```

### Debugging Tips

1. **Check backend logs:**
   ```bash
   tail -f /tmp/backend.log | grep "test_runner\|reviews\|marketplace"
   ```

2. **Test individual endpoints with curl:**
   ```bash
   TOKEN="$(cat /tmp/token)"
   curl -H "Authorization: Bearer ${TOKEN}" 'http://localhost:3000/api/v1/endpoint' | jq '.'
   ```

3. **Check response structure:**
   ```go
   // Print response for debugging (remove in production)
   body, _ := io.ReadAll(resp.Body)
   fmt.Printf("Response: %s\n", string(body))
   ```

---

## Notes

1. **Focus on Critical Paths**: Don't aim for 100% coverage - focus on critical user journeys
2. **E2E over Unit**: Functional tests should test complete workflows, not individual endpoints
3. **Data Setup**: Each test should create its own data (listings, users, orders) and clean up
4. **Idempotency**: Tests should be runnable multiple times without side effects
5. **Error Scenarios**: Test both happy path AND error cases (validation, not found, unauthorized)

---

**Created**: 2025-10-18
**Updated**: 2025-10-18 (Added "How to Write Functional Tests" section)
**Status**: Ready for implementation
**Next Steps**: Implement Phase 1 tests (Marketplace Listing Creation Flow)
