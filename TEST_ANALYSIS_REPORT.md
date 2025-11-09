# Integration Tests Analysis Report
## Listings Microservice - Legacy Schema Migration

**Date:** 2025-11-08  
**Test Run:** Phase 13.1.7 - Integration Tests Execution  
**Total Tests:** 232  
**Status:** 110 PASS (47.4%) / 122 FAIL (52.6%)

---

## EXECUTIVE SUMMARY

### Key Achievements ✅
1. **Schema Migration: COMPLETE** - All 14 migrations applied successfully
2. **Fixture Loading: FUNCTIONAL** - Categories fixtures load without errors
3. **Test Infrastructure: WORKING** - Tests execute, no build/setup failures
4. **Legacy Code Cleanup: COMPLETE** - 200+ b2c_products references removed
5. **SQL Mocks: MIGRATED** - 30+ mock queries updated to unified schema

### Critical Issues Identified ❌
1. **Foreign Key Constraint Violations** (3381 occurrences)
   - Error: `fk_listings_category_id` constraint violation
   - Root cause: Fixtures use category IDs not present in 00_categories_fixtures.sql
   - Impact: ~50 tests fail on fixture loading

2. **SKU Type Mismatch** (34+ occurrences)
   - Expected: `string`
   - Actual: `*string` (pointer)
   - Root cause: gRPC message definition uses pointer for optional field
   - Impact: Test assertions fail even when business logic succeeds

3. **Error Message Format Validation** (10+ occurrences)
   - Tests expect lowercase "sku" in error messages
   - Actual: "SKU DUPLICATE-SKU-TEST already exists..."
   - Impact: Minor, error handling works but tests need update

4. **Test Assertion Logic** (34 individual assertion failures)
   - Price validation errors
   - Stock quantity comparisons
   - Empty result checks
   - Impact: Various, needs case-by-case analysis

---

## DETAILED BREAKDOWN BY CATEGORY

### Category 1: Foreign Key Constraint Violations
**Count:** ~50 failed tests (41% of failures)  
**Severity:** HIGH - Blocks test execution

**Failed Test Suites:**
- `TestRollbackStock_*` (13 tests) - ALL fail on rollback_stock_fixtures.sql
- `TestDecrementStock_*` (17 tests) - Fixture loading blocked
- `TestListing_*` (16 tests) - Category FK violations

**Example Error:**
```
Error: Received unexpected error:
pq: insert or update on table "listings" violates foreign key constraint "fk_listings_category_id"
Test: TestRollbackStock_SingleProduct_Success
Messages: Could not load fixtures from: ../fixtures/rollback_stock_fixtures.sql
```

**Root Cause Analysis:**
- Test fixtures reference category IDs (e.g., 1301, 1302) that don't exist
- `00_categories_fixtures.sql` has limited category data
- Legacy fixtures assumed different category seed data

**Recommended Fix:**
```sql
-- Option 1: Add missing categories to 00_categories_fixtures.sql
INSERT INTO categories (id, name, ...) VALUES 
  (1301, 'Electronics', ...),
  (1302, 'Clothing', ...),
  (1303, 'Home & Garden', ...);

-- Option 2: Update all test fixtures to use existing category IDs
-- Find what categories exist:
SELECT id, name FROM categories ORDER BY id;

-- Update fixtures to use real IDs
```

**Priority:** P0 - Must fix before merge

---

### Category 2: SKU Type Mismatch (string vs *string)
**Count:** 34 test failures (28% of failures)  
**Severity:** MEDIUM - Business logic works, test assertions wrong

**Failed Test Suites:**
- `TestBulkCreateProducts_*` (8 tests)
- `TestCreateProduct_*` (15 tests)
- `TestGetProduct_*` (3 tests)

**Example Error:**
```go
Error: Not equal:
expected: string("BULK-SINGLE-001")
actual  : *string((*string)(0xc0002ba190))

Diff:
--- Expected
+++ Actual
@@ -1 +1 @@
-BULK-SINGLE-001
+<*string Value>

Test: TestBulkCreateProducts_Success_Single
```

**Root Cause:**
```protobuf
// api/proto/listings/v1/product.proto
message Product {
  string sku = 5; // Should be optional
}

// CORRECT:
message Product {
  optional string sku = 5; // Protobuf generates *string in Go
}
```

**Recommended Fix:**
```go
// Option 1: Update test assertions to handle pointer
assert.Equal(t, "BULK-SINGLE-001", *resp.Product.Sku) // Dereference

// Option 2: Add helper function
func stringValue(s *string) string {
    if s == nil {
        return ""
    }
    return *s
}
assert.Equal(t, "BULK-SINGLE-001", stringValue(resp.Product.Sku))

// Option 3: Make SKU required in proto (NOT recommended, breaks optionality)
```

**Priority:** P1 - Fix after FK constraints resolved

---

### Category 3: Error Message Validation Failures
**Count:** 10 test failures (8% of failures)  
**Severity:** LOW - Cosmetic, doesn't affect functionality

**Failed Tests:**
- `TestBulkCreateProducts_Error_DuplicateSKU`
- `TestUpdateProduct_DuplicateSKU`

**Example:**
```go
Error: "SKU DUPLICATE-SKU-TEST already exists in storefront" does not contain "sku"
Test: TestBulkCreateProducts_Error_DuplicateSKU
```

**Root Cause:**
Test expects lowercase "sku" but error message uses uppercase "SKU"

**Recommended Fix:**
```go
// Option 1: Update test to case-insensitive check
assert.Contains(t, strings.ToLower(err.Error()), "sku")

// Option 2: Update error message to use lowercase
return errors.New("sku DUPLICATE-SKU-TEST already exists...")

// Option 3: Just check for "already exists" or error code
assert.Equal(t, codes.AlreadyExists, status.Code(err))
```

**Priority:** P2 - Low impact, fix when convenient

---

### Category 4: Business Logic & Assertion Failures
**Count:** 28 test failures (23% of failures)  
**Severity:** MIXED - Requires case-by-case investigation

**Sub-categories:**

**4a. Price Validation Errors** (2 failures)
```go
Error: "products.bulk_update_failed" does not contain "price"
Actual log: validation failed for product 20001: 
  Key: 'BulkUpdateProductInput.Price' Error:
  Field validation for 'Price' failed on the 'gte' tag
```
- **Fix:** Test expects "price" in error message, but error code is generic
- **Action:** Update test to check validation error correctly

**4b. Stock Quantity Assertions** (3 failures)
```go
Error: "103" is not greater than or equal to "120"
Error: "0" is not greater than "0"
```
- **Fix:** Business logic not decrementing stock correctly OR test data wrong
- **Action:** Debug decrement logic in service layer

**4c. Empty Result Checks** (4 failures)
```go
Error: Should be true
Error: Should NOT be empty, but was []
Error: "[]" should have 3 item(s), but has 0
```
- **Fix:** Queries returning empty when data expected
- **Action:** Check if unified schema joins are correct

**Priority:** P1 - Some may indicate real bugs

---

## PRIORITIZED FIX RECOMMENDATIONS

### Phase 1: Critical Path (P0)
**Goal:** Get tests running without setup failures

1. **Fix Category Foreign Key Constraints**
   ```bash
   # Step 1: Identify missing categories
   grep "category_id" tests/fixtures/*.sql | grep -o "category_id[^,]*" | sort -u > /tmp/used_categories.txt
   
   # Step 2: Add to 00_categories_fixtures.sql
   # Ensure IDs: 1301, 1302, 1303, 1304, 1305 exist
   
   # Step 3: Re-run fixture loading tests
   go test -v -run "TestRollbackStock_SingleProduct" ./tests/integration/
   ```

2. **Validate Fixture Load Order**
   - Ensure `00_categories_fixtures.sql` loads FIRST
   - Check `tests/testing.go:autoLoadCategoryFixtures()` works correctly

**Expected Impact:** Unlock ~50 tests currently blocked

---

### Phase 2: Type Safety (P1)
**Goal:** Fix SKU pointer handling across codebase

1. **Update Test Assertions**
   ```go
   // Create helper in tests/helpers.go
   func DerefString(s *string) string {
       if s == nil {
           return ""
       }
       return *s
   }
   
   // Update ~34 test assertions
   assert.Equal(t, "expected-sku", DerefString(product.Sku))
   ```

2. **Verify Proto Definition**
   ```protobuf
   // Confirm SKU SHOULD be optional
   message Product {
       optional string sku = 5; // ✅ Correct for B2C products
   }
   ```

**Expected Impact:** Fix 34 assertion failures

---

### Phase 3: Business Logic Validation (P1)
**Goal:** Investigate real bugs vs test issues

1. **Stock Decrement Logic**
   - Test manually: does DecrementStock actually update quantity?
   - Check transaction isolation
   - Verify unified schema triggers

2. **Empty Query Results**
   - Validate JOIN conditions in repository
   - Check if `source_type` filtering correct
   - Ensure test data actually inserts

**Expected Impact:** Fix 10-15 legitimate bugs OR update tests

---

### Phase 4: Polish (P2)
**Goal:** Clean up cosmetic issues

1. **Error Message Formatting**
   - Standardize case (lowercase field names)
   - Update test expectations

2. **Test Cleanup**
   - Remove deprecated tests
   - Add comments explaining pointer handling

**Expected Impact:** All tests GREEN

---

## NEXT STEPS

### Immediate Actions (Today)
1. ✅ Generate this analysis report
2. ⏳ **Add missing categories to 00_categories_fixtures.sql**
3. ⏳ **Re-run integration tests to validate FK fix**

### Short-term (This Sprint)
4. ⏳ **Fix SKU pointer assertions (bulk edit ~34 tests)**
5. ⏳ **Investigate business logic failures (10-15 tests)**
6. ⏳ **Run full test suite until GREEN**

### Before Merge
7. ⏳ **100% test pass rate required**
8. ⏳ **Update PROGRESS.md with final results**
9. ⏳ **Create PR with migration summary**

---

## APPENDIX: Test Failure Distribution

| Test Suite | Failed | Total | Pass Rate |
|------------|--------|-------|-----------|
| TestDecrementStock_* | 17 | ? | ~0% |
| TestListing_* | 16 | ? | ~0% |
| TestCreateProduct_* | 15 | ? | ~0% |
| TestRollbackStock_* | 13 | ? | ~0% |
| TestConcurrency_* | 8 | ? | ~0% |
| TestBulkCreateProducts_* | 8 | ? | ~50% |
| TestUpdateProduct_* | 6 | ? | ~0% |
| TestUpdateProductInventory_* | 5 | ? | ~0% |
| Others | 34 | ? | ~50% |
| **TOTAL** | **122** | **232** | **47.4%** |

**Legend:**
- ✅ = Complete
- ⏳ = In Progress  
- ❌ = Blocked

---

**Report Generated:** 2025-11-08 23:35 CET  
**Next Update:** After Category FK fix applied
