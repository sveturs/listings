# Phase 17 Days 18-19 - Data Migration from Monolith

## 📋 Overview

This guide describes the data migration process from the monolith database (`svetubd`) to the Orders microservice database (`listings_dev_db`).

**Current status:**
- Monolith has **only** `inventory_reservations` table with orders data (2 active records)
- Shopping carts, orders, and order items **do not exist** in monolith - they are new features in the microservice
- Migration focuses on transferring active inventory reservations

---

## 🗄️ Database Schemas

### Monolith (svetubd)
**Connection:**
```bash
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable"
```

**Tables:**
- ✅ `inventory_reservations` (2 records, all active)
  - Uses `product_id` (FK to products)
  - Uses enum `reservation_status` (active, committed, released, expired)
  - Has trigger `trigger_update_inventory_reservations_updated_at`

### Microservice (listings_dev_db)
**Connection:**
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
```

**Tables:**
- ✅ `shopping_carts` - NEW (0 records)
- ✅ `cart_items` - NEW (0 records)
- ✅ `orders` - NEW (0 records)
- ✅ `order_items` - NEW (0 records)
- ✅ `inventory_reservations` (0 records)
  - Uses `listing_id` (FK to listings) - **DIFFERENT from monolith!**
  - Uses varchar(20) `status` with CHECK constraint
  - Has updated_at trigger

---

## 🔍 Key Differences

| Field | Monolith | Microservice | Migration Strategy |
|-------|----------|--------------|-------------------|
| Product reference | `product_id` (bigint) | `listing_id` (bigint) | Direct mapping: `product_id` → `listing_id` |
| Status type | enum `reservation_status` | varchar(20) with CHECK | Cast enum to text: `status::text` |
| Order FK | `order_id` bigint NOT NULL | `order_id` bigint NULL | Keep original value (NULL if no order) |

---

## 🚀 Migration Script

**Location:** `/p/github.com/sveturs/listings/cmd/migrate_orders_from_monolith/main.go`

**Features:**
- ✅ Dry-run mode by default (no writes)
- ✅ Idempotent (safe to run multiple times)
- ✅ Transaction-based (atomic commits)
- ✅ FK constraint validation
- ✅ Detailed progress logging
- ✅ Rollback support on errors

**What it migrates:**
- **Inventory Reservations:** Only active reservations (`expires_at > NOW()` AND `status = 'active'`)

**What it skips:**
- Reservations for non-existent listings
- Already migrated records (by ID)
- Expired or released reservations

---

## 📖 Prerequisites

1. **Both databases must be running:**
   ```bash
   # Check monolith PostgreSQL
   psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" -c "SELECT version();"

   # Check microservice PostgreSQL (Docker container)
   docker ps | grep listings-postgres
   psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" -c "SELECT version();"
   ```

2. **Microservice migrations must be applied:**
   ```bash
   cd /p/github.com/sveturs/listings
   make migrate-up
   # Should show: 5 migrations applied (shopping_carts, cart_items, orders, order_items, inventory_reservations)
   ```

3. **Go 1.21+ installed:**
   ```bash
   go version  # Should be >= 1.21
   ```

---

## 🏃 Usage

### Step 1: Dry-Run (Recommended First)

**Check what will be migrated without making changes:**

```bash
cd /p/github.com/sveturs/listings/cmd/migrate_orders_from_monolith

# Simple dry-run
go run main.go

# Verbose dry-run (see each record)
go run main.go --verbose
```

**Expected output:**
```
🚀 Starting Orders data migration from monolith to microservice
   Mode: DRY-RUN (no writes)

📡 Connecting to databases...
✅ Connected to both databases

📦 Migrating inventory_reservations...
   Found 2 active reservations to migrate
   [DRY-RUN] Would insert reservation 1 (listing=123, quantity=1, expires=2025-11-15T10:00:00Z)
   [DRY-RUN] Would insert reservation 2 (listing=456, quantity=2, expires=2025-11-15T11:00:00Z)

═══════════════════════════════════════════════════════════
📊 MIGRATION SUMMARY
═══════════════════════════════════════════════════════════
Mode:                    DRY-RUN (no writes)
Duration:                152ms

Inventory Reservations:
  Total found:           2
  Successfully migrated: 2
  Skipped:               0
  Failed:                0
═══════════════════════════════════════════════════════════

ℹ️  This was a DRY RUN. No data was written.
To execute migration, run with --dry-run=false
```

### Step 2: Execute Migration

**⚠️ WARNING: This will write data to microservice database!**

```bash
cd /p/github.com/sveturs/listings/cmd/migrate_orders_from_monolith

# Execute with 5-second safety delay
go run main.go --dry-run=false

# Execute with verbose logging
go run main.go --dry-run=false --verbose
```

**Expected output:**
```
⚠️  WARNING: Running in EXECUTE mode. Data will be written to microservice database.
Press Ctrl+C within 5 seconds to cancel...

🚀 Starting Orders data migration from monolith to microservice
   Mode: EXECUTE (writes enabled)

📡 Connecting to databases...
✅ Connected to both databases

📦 Migrating inventory_reservations...
   Found 2 active reservations to migrate
   ✅ Inserted reservation 1 (listing=123, quantity=1)
   ✅ Inserted reservation 2 (listing=456, quantity=2)
   ✅ Transaction committed successfully

🔍 Verifying foreign key constraints...
   ✅ All FK constraints valid

═══════════════════════════════════════════════════════════
📊 MIGRATION SUMMARY
═══════════════════════════════════════════════════════════
Mode:                    EXECUTE (writes enabled)
Duration:                234ms

Inventory Reservations:
  Total found:           2
  Successfully migrated: 2
  Skipped:               0
  Failed:                0
═══════════════════════════════════════════════════════════

✅ Migration completed successfully!
```

### Step 3: Verify Results

**Check migrated data:**

```bash
# Count migrated reservations
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM inventory_reservations;"

# View migrated records
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT id, listing_id, quantity, status, expires_at FROM inventory_reservations ORDER BY id;"

# Check FK integrity (should be 0)
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM inventory_reservations ir LEFT JOIN listings l ON ir.listing_id = l.id WHERE l.id IS NULL;"
```

---

## 🔧 Advanced Usage

### Custom Connection Strings

```bash
go run main.go \
  --monolith-dsn="postgres://user:pass@host:port/dbname?sslmode=disable" \
  --microservice-dsn="postgres://user:pass@host:port/dbname?sslmode=disable" \
  --dry-run=false
```

### Build Standalone Binary

```bash
cd /p/github.com/sveturs/listings/cmd/migrate_orders_from_monolith
go build -o migrate_orders main.go

# Run binary
./migrate_orders --dry-run=false --verbose
```

---

## ✅ Safety Features

### 1. Dry-Run by Default
- **Default mode:** `--dry-run=true` (no writes)
- Must explicitly set `--dry-run=false` to write data

### 2. 5-Second Safety Delay
- When running in execute mode, script waits 5 seconds
- Allows cancellation with Ctrl+C before writes

### 3. Idempotency
- Checks if record already exists by ID
- Skips duplicates automatically
- Safe to run multiple times

### 4. Transaction Safety
- All writes in single transaction
- Automatic rollback on errors
- Either all succeed or none

### 5. FK Validation
- Checks if `listing_id` exists before insert
- Verifies FK constraints after migration
- Reports orphaned records

### 6. Detailed Logging
- Progress updates for each step
- Error messages with context
- Final summary with statistics

---

## 🐛 Troubleshooting

### Error: "monolith database unreachable"

**Cause:** PostgreSQL on port 5433 not running

**Solution:**
```bash
sudo systemctl status postgresql
sudo systemctl start postgresql
```

### Error: "microservice database unreachable"

**Cause:** Docker container not running

**Solution:**
```bash
cd /p/github.com/sveturs/listings
docker-compose up -d
docker ps | grep listings-postgres
```

### Error: "relation 'shopping_carts' does not exist"

**Cause:** Migrations not applied

**Solution:**
```bash
cd /p/github.com/sveturs/listings
make migrate-up
```

### Warning: "Found N orphaned reservations"

**Cause:** Reservations reference non-existent listings

**Impact:** FK constraints violated, queries may fail

**Solution:**
```bash
# Identify orphaned records
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT ir.id, ir.listing_id FROM inventory_reservations ir LEFT JOIN listings l ON ir.listing_id = l.id WHERE l.id IS NULL;"

# Delete orphaned records (CAREFUL!)
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "DELETE FROM inventory_reservations WHERE id IN (SELECT ir.id FROM inventory_reservations ir LEFT JOIN listings l ON ir.listing_id = l.id WHERE l.id IS NULL);"
```

### Skipped reservations due to missing listings

**Cause:** Product from monolith doesn't exist as listing in microservice

**Solution:**
1. Run listings migration first (if needed)
2. Or accept that reservations without listings will be skipped
3. Check logs for specific product IDs that were skipped

---

## 📊 Verification Queries

### Compare counts between databases

```bash
# Monolith: Count active reservations
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" \
  -c "SELECT COUNT(*) FROM inventory_reservations WHERE expires_at > NOW() AND status = 'active';"

# Microservice: Count migrated reservations
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM inventory_reservations;"
```

### Verify data integrity

```bash
# Check status values are valid
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT DISTINCT status FROM inventory_reservations;"
# Expected: 'active' (only active were migrated)

# Check all reservations have valid listing_id
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM inventory_reservations WHERE listing_id NOT IN (SELECT id FROM listings);"
# Expected: 0 (all should reference existing listings)

# Check expires_at is in the future
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM inventory_reservations WHERE expires_at <= NOW();"
# Expected: 0 (only future expiries were migrated)
```

### Sample data comparison

```bash
# Monolith: Show sample reservation
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" \
  -c "SELECT id, product_id, quantity, status::text, expires_at FROM inventory_reservations LIMIT 1;"

# Microservice: Show same reservation (by ID)
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT id, listing_id, quantity, status, expires_at FROM inventory_reservations WHERE id = <ID>;"

# Values should match (product_id = listing_id, status same, etc.)
```

---

## 🔄 Rollback

If migration needs to be reverted:

```bash
# Delete all migrated reservations
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "DELETE FROM inventory_reservations;"

# Verify deletion
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM inventory_reservations;"
# Expected: 0
```

**⚠️ WARNING:** This deletes ALL reservations, including any created after migration!

For selective rollback, note the IDs migrated from the script output.

---

## 📈 Performance Considerations

**Current scale:** 2 records (negligible performance impact)

**Future scale considerations:**
- Migration script processes records one-by-one (safe but slower)
- For large datasets (10k+ records), consider batch inserts
- Transaction size may need tuning for very large migrations
- Add progress bar for migrations taking >10 seconds

**Recommended batch size for future:** 1000 records per transaction

---

## ✅ Migration Checklist

Before migration:
- [ ] Monolith database accessible (port 5433)
- [ ] Microservice database accessible (port 35434)
- [ ] Microservice migrations applied (5 tables exist)
- [ ] Listings data migrated (so inventory_reservations can reference them)
- [ ] Backup of both databases (optional but recommended)

During migration:
- [ ] Run dry-run first and review output
- [ ] Check for skipped/failed records in dry-run
- [ ] Verify counts match expectations
- [ ] Run execute mode
- [ ] Monitor logs for errors

After migration:
- [ ] Verify record counts match
- [ ] Check FK integrity (no orphaned records)
- [ ] Spot-check sample records for data accuracy
- [ ] Test microservice operations (create order, reserve inventory)
- [ ] Monitor microservice logs for errors

---

## 📝 Notes

1. **No data deletion from monolith:** Migration only reads from monolith, never writes or deletes
2. **ID preservation:** Original IDs are preserved to minimize disruption
3. **Timestamp preservation:** `created_at` and `updated_at` from monolith are kept
4. **Active only:** Only active, non-expired reservations are migrated (historical data stays in monolith)
5. **FK mapping:** `product_id` → `listing_id` (direct 1:1 mapping assumed)

---

## 📞 Support

For issues or questions:
1. Check logs for specific error messages
2. Consult Troubleshooting section above
3. Verify Prerequisites are met
4. Run verification queries to diagnose state

---

**Last updated:** 2025-11-14
**Phase:** 17 Days 18-19
**Status:** ✅ Ready for testing
