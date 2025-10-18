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

## Notes

1. **Focus on Critical Paths**: Don't aim for 100% coverage - focus on critical user journeys
2. **E2E over Unit**: Functional tests should test complete workflows, not individual endpoints
3. **Data Setup**: Each test should create its own data (listings, users, orders) and clean up
4. **Idempotency**: Tests should be runnable multiple times without side effects
5. **Error Scenarios**: Test both happy path AND error cases (validation, not found, unauthorized)

---

**Created**: 2025-10-18
**Status**: Initial draft - ready for implementation
**Next Steps**: Implement Phase 1 tests (Marketplace Listing Creation Flow)
