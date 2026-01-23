# Migration Summary: Category ID bigint → UUID

**Date:** 2025-12-17
**Status:** ✅ COMPLETED
**Migration:** `20251217150000_fix_category_id_to_uuid`

## Problem Statement

The `listings.category_id` column was defined as `bigint` while `categories.id` is `UUID`, causing:
1. ❌ Foreign key constraint impossible
2. ❌ No referential integrity
3. ❌ Empty listings database (no test data)
4. ❌ Filters couldn't be tested

## Solution Implemented

### 1. Database Schema Migration

**Migration files created:**
- `/p/github.com/vondi-global/listings/migrations/20251217150000_fix_category_id_to_uuid.up.sql`
- `/p/github.com/vondi-global/listings/migrations/20251217150000_fix_category_id_to_uuid.down.sql`

**Changes:**
```sql
-- Before
category_id BIGINT

-- After
category_id UUID REFERENCES categories(id) ON DELETE SET NULL ON UPDATE CASCADE
```

**Migration applied successfully:**
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" \
  -f migrations/20251217150000_fix_category_id_to_uuid.up.sql
```

### 2. Test Data Insertion

**15 test listings created:**

| Category | UUID | Count | Products |
|----------|------|-------|----------|
| Apple iPhone | `dbd48c81-ab21-42f0-b74d-45e2a408dc5e` | 5 | iPhone 15 Pro, 15 Pro Max, 14 Pro, 13, 12 |
| Samsung telefoni | `8b75d427-1bf9-4b7d-9446-94b40f085546` | 5 | Galaxy S24 Ultra, S24+, S23 FE, S23, A54 |
| Pametni telefoni | `ba829402-0e4e-467b-b117-e1da1e2b51ce` | 5 | Xiaomi 14 Pro, Google Pixel 8 Pro, OnePlus 12, Huawei Pura 70 Pro, Nothing Phone (2) |

**Price range:** 45,000 RSD - 135,000 RSD
**Total quantity:** 113 units

**Attributes included:**
- `brand` (apple, samsung, xiaomi, google, oneplus, huawei, nothing)
- `color` (Black Titanium, Blue Titanium, Deep Purple, etc.)
- `storage` (128GB, 256GB, 512GB)
- `ram` (4GB, 6GB, 8GB, 12GB, 16GB)
- `condition` (new, used)

### 3. OpenSearch Reindexing

**Command:**
```bash
python3 scripts/reindex_listings.py \
  --target-password listings_secret \
  --target-db listings_dev_db \
  --target-port 35434 \
  --batch-size 100
```

**Results:**
- ✅ Total Listings: 15
- ✅ Indexed: 15
- ❌ Failed: 0
- ⚡ Duration: 0.27s
- 📈 Rate: 54.83 docs/sec

### 4. Verification

#### Database Schema ✅
```sql
\d listings

category_id | uuid | | |
```

#### Foreign Key Constraint ✅
```sql
listings_category_id_fkey
  FOREIGN KEY (category_id) REFERENCES categories(id)
  ON UPDATE CASCADE ON DELETE SET NULL
```

#### OpenSearch Mapping ✅
```json
{
  "category_id": {
    "type": "keyword"
  }
}
```

#### Filter Queries Work ✅

**Test 1: Category filter**
```bash
curl "http://localhost:9200/listings_microservice/_search" \
  -d '{"query": {"term": {"category_id": "dbd48c81-ab21-42f0-b74d-45e2a408dc5e"}}}'

# Result: 5 iPhone listings found
```

**Test 2: Category + Price range filter**
```bash
curl "http://localhost:9200/listings_microservice/_search" \
  -d '{
    "query": {
      "bool": {
        "must": [
          {"term": {"category_id": "dbd48c81-ab21-42f0-b74d-45e2a408dc5e"}},
          {"range": {"price": {"gte": 50000, "lte": 100000}}}
        ]
      }
    }
  }'

# Result: 3 iPhones (14 Pro, 13, 12) in price range 55k-95k RSD
```

**Test 3: Attribute filter (brand)**
```bash
curl "http://localhost:9200/listings_microservice/_search" \
  -d '{"query": {"term": {"attributes.code": "brand"}}, "size": 10}'

# Result: All listings have brand attribute
```

## Known Issues

### Search API Returns 500 Error

**Error:** `rpc error: code = Unimplemented desc = unknown service search.v1.SearchService`

**Root Cause:**
The backend monolith expects a separate **Search microservice** (gRPC service) but:
- ❌ Search microservice is NOT running
- ❌ No fallback to direct OpenSearch queries
- ✅ OpenSearch itself works perfectly

**Workaround:**
Use OpenSearch API directly for now:
```bash
# Direct OpenSearch query (works perfectly)
curl "http://localhost:9200/listings_microservice/_search?q=title:iPhone"
```

**TODO:**
Either:
1. Start Search microservice (if it exists)
2. Add fallback to direct OpenSearch queries in backend
3. Disable `USE_SEARCH_MICROSERVICE` feature flag

## Summary

| Task | Status | Notes |
|------|--------|-------|
| Create migration files | ✅ DONE | Up/down migrations |
| Apply migration | ✅ DONE | category_id is now UUID |
| Verify schema | ✅ DONE | FK constraint added |
| Insert test data | ✅ DONE | 15 listings, 3 categories |
| Reindex OpenSearch | ✅ DONE | 15/15 indexed successfully |
| Test filters (direct) | ✅ DONE | Category, price, attribute filters work |
| Test API endpoints | ⚠️ BLOCKED | Search microservice not running |

## Files Created

1. `/p/github.com/vondi-global/listings/migrations/20251217150000_fix_category_id_to_uuid.up.sql`
2. `/p/github.com/vondi-global/listings/migrations/20251217150000_fix_category_id_to_uuid.down.sql`
3. `/tmp/insert_test_listings.sql` (test data)
4. `/p/github.com/vondi-global/listings/MIGRATION_20251217_CATEGORY_UUID_SUMMARY.md` (this file)

## Next Steps

1. ✅ Migration complete - category_id is now UUID
2. ✅ Test data available - 15 listings across 3 categories
3. ✅ OpenSearch working - filters operational
4. ⏭️ Fix Search microservice or add fallback mechanism
5. ⏭️ Test frontend filters once Search API is working

---

**Completed by:** Claude Code
**Duration:** ~30 minutes
**Lines of code changed:** ~150 (migration + test data)
