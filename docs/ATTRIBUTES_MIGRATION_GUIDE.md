# Attributes Migration Guide

## Overview

Migrate `unified_attributes` table from monolith database to `attributes` table in microservice with i18n JSONB conversion.

**Source:** `svetubd.unified_attributes` (localhost:5433)
**Target:** `listings_dev_db.attributes` (localhost:35434)

**Key Transformation:** VARCHAR fields → JSONB i18n format
- `name` (VARCHAR) → `name` (JSONB) with `{en, ru, sr}` keys
- `display_name` (VARCHAR) → `display_name` (JSONB) with `{en, ru, sr}` keys

---

## Migration Methods

### ✅ Recommended: Go Migration Tool (Fully Automated)

**Advantages:**
- ✅ Fully automated (one command)
- ✅ Direct database-to-database transfer
- ✅ Built-in validation
- ✅ Idempotent (safe to re-run)
- ✅ Transaction support
- ✅ Dry-run mode

**Usage:**

```bash
# 1. Dry run (preview only)
cd /p/github.com/sveturs/listings
go run ./cmd/migrate_attributes/main.go --dry-run

# 2. Live migration (with verbose output)
go run ./cmd/migrate_attributes/main.go -v

# 3. Silent migration (production mode)
go run ./cmd/migrate_attributes/main.go
```

**Expected Output:**
```
=== Attributes Migration Tool ===
Mode: LIVE MIGRATION

Connecting to source database (monolith)...
Connecting to target database (microservice)...
Fetching attributes from source database...
✓ Found 203 attributes in source database

Checking existing attributes in target database...
✓ Found 0 existing attributes in target database

Migration plan:
  - Total attributes: 203
  - Already migrated: 0
  - To migrate: 203

Starting migration...
Executing batch insert of 203 attributes...
  Inserted 50/203 attributes...
  Inserted 100/203 attributes...
  Inserted 150/203 attributes...
  Inserted 200/203 attributes...

✓ Migration complete!
  - Successfully migrated: 203 attributes

Validating migration...
  - Total attributes in target: 203
  - Sample migrated attributes:
    [86] year: name=year, display_name=Godište
    [87] mileage: name=mileage, display_name=Kilometraža
    [88] area: name=area, display_name=Površina (m²)
    [89] processor: name=processor, display_name=Procesor
    [90] car_make_id: name=car_make_id, display_name=Car Make ID
✓ All validations passed!
```

---

### Alternative: SQL Script (Manual)

**Advantages:**
- Simple SQL script
- Easy to review

**Disadvantages:**
- Requires manual data export/import
- Two-step process

**Usage:**

```bash
# Step 1: Export from source DB
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" \
  -c "COPY (
    SELECT id, code, name, display_name, attribute_type, purpose,
           options, validation_rules, ui_settings,
           is_searchable, is_filterable, is_required,
           affects_stock, affects_price, sort_order, is_active,
           created_at, updated_at,
           legacy_category_attribute_id, legacy_product_variant_attribute_id,
           is_variant_compatible, icon, show_in_card
    FROM unified_attributes
    ORDER BY id
  ) TO '/tmp/attributes_export.csv' WITH CSV HEADER"

# Step 2: Import into target DB with JSONB conversion
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f /p/github.com/sveturs/listings/scripts/migrate_attributes.sql
```

---

## Post-Migration Validation

### Quick Check
```bash
# Run validation script
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f /p/github.com/sveturs/listings/scripts/validate_attributes.sql
```

### Manual Validation Queries

**1. Count Check:**
```sql
-- Source DB
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" \
  -c "SELECT COUNT(*) FROM unified_attributes"

-- Target DB
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM attributes"
```

**2. JSONB Structure Check:**
```sql
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT id, code, name, display_name FROM attributes LIMIT 5"
```

Expected output:
```
 id |    code     |                      name                       |                    display_name
----+-------------+-------------------------------------------------+-----------------------------------------------------
 86 | year        | {"en": "year", "ru": "year", "sr": "year"}      | {"en": "Godište", "ru": "Godište", "sr": "Godište"}
 87 | mileage     | {"en": "mileage", "ru": "mileage", "sr": ...}   | {"en": "Kilometraža", "ru": ...}
```

**3. Check for Missing Data:**
```sql
-- Should return 0
SELECT COUNT(*) FROM attributes
WHERE name IS NULL
   OR name->>'en' IS NULL
   OR display_name IS NULL
   OR display_name->>'en' IS NULL;
```

**4. Verify Sequence:**
```sql
SELECT
    last_value as seq_value,
    (SELECT MAX(id) FROM attributes) as max_id
FROM attributes_id_seq;
```
(seq_value should be >= max_id)

---

## Rollback

If something goes wrong:

```bash
# Run rollback script (interactive prompts for safety)
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f /p/github.com/sveturs/listings/scripts/rollback_attributes.sql
```

**OR** manual rollback:
```sql
BEGIN;
DELETE FROM attributes;
SELECT setval('attributes_id_seq', 1, false);
-- Review changes, then:
COMMIT;  -- or ROLLBACK to undo
```

---

## Troubleshooting

### Issue: "duplicate key value violates unique constraint"

**Cause:** Attributes with same code already exist in target DB

**Solution:** Migration is idempotent - it will skip existing records
```bash
# Check which codes already exist
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT code FROM attributes ORDER BY code"
```

### Issue: "invalid input syntax for type json"

**Cause:** Malformed JSONB conversion

**Solution:** Check source data for special characters
```sql
-- Find problematic records
SELECT id, code, name, display_name
FROM unified_attributes
WHERE name ~ '[\\"]'  -- contains backslash or quotes
   OR display_name ~ '[\\"]';
```

### Issue: "sequence is not updated"

**Cause:** Manual INSERT without sequence update

**Solution:**
```sql
SELECT setval('attributes_id_seq', (SELECT MAX(id) FROM attributes));
```

---

## Files Reference

**Go Migration Tool:**
- `/p/github.com/sveturs/listings/cmd/migrate_attributes/main.go`

**SQL Scripts:**
- Migration: `/p/github.com/sveturs/listings/scripts/migrate_attributes.sql`
- Validation: `/p/github.com/sveturs/listings/scripts/validate_attributes.sql`
- Rollback: `/p/github.com/sveturs/listings/scripts/rollback_attributes.sql`

**Documentation:**
- This guide: `/p/github.com/sveturs/listings/docs/ATTRIBUTES_MIGRATION_GUIDE.md`

---

## Migration Checklist

- [ ] Review source data: `SELECT COUNT(*) FROM unified_attributes`
- [ ] Check target is empty: `SELECT COUNT(*) FROM attributes`
- [ ] Run dry-run: `go run ./cmd/migrate_attributes/main.go --dry-run`
- [ ] Run migration: `go run ./cmd/migrate_attributes/main.go -v`
- [ ] Validate count matches
- [ ] Validate JSONB structure
- [ ] Verify sequence updated
- [ ] Test application with migrated data
- [ ] Document migration timestamp

---

## Post-Migration Notes

**Date Migrated:** _______________
**Records Migrated:** _______________
**Migration Time:** _______________
**Validated By:** _______________

**Issues Encountered:**
-

**Resolution:**
-
