# Category Translations Migration Guide

## Overview

This directory contains SQL scripts to add multi-language support (English, Russian, Serbian) to the `c2c_categories` table.

**Total categories to translate:** 77

## Scripts Execution Order

### 1. Add Translation Columns
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" -f 01_add_translation_columns.sql
```

**What it does:**
- Adds `title_en`, `title_ru`, `title_sr` columns
- Creates indexes for better search performance
- Adds column comments

### 2. Create Backup (Optional but Recommended)
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" -f 02_backup_categories_before_translation.sql
```

**What it does:**
- Creates backup table `c2c_categories_backup_20251110`
- Shows verification counts
- Reports categories without translations

### 3. Apply Translations
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" -f 03_add_category_translations.sql
```

**What it does:**
- Updates all 77 categories with translations
- Organized by category hierarchy
- Includes verification queries
- Shows final statistics

**Expected output:**
```
total_categories | with_english | with_russian | with_serbian | missing_translations
----------------+--------------+--------------+--------------+---------------------
              77 |           77 |           77 |           77 |                    0
```

### 4. Rollback (If Needed)
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" -f 04_rollback_category_translations.sql
```

**What it does:**
- Provides 3 rollback options (all commented by default)
- Option 1: Restore from backup
- Option 2: Clear translations (set NULL)
- Option 3: Drop columns completely

## Quick Start (All-in-One)

```bash
# Navigate to scripts directory
cd /p/github.com/sveturs/listings/scripts

# Run all scripts in sequence
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f 01_add_translation_columns.sql \
  -f 02_backup_categories_before_translation.sql \
  -f 03_add_category_translations.sql
```

## Verification Queries

### Check translation coverage:
```sql
SELECT
  COUNT(*) as total,
  COUNT(title_en) as en,
  COUNT(title_ru) as ru,
  COUNT(title_sr) as sr
FROM c2c_categories;
```

### Find missing translations:
```sql
SELECT id, name, slug, title_en, title_ru, title_sr
FROM c2c_categories
WHERE title_en IS NULL OR title_ru IS NULL OR title_sr IS NULL;
```

### View sample translations:
```sql
SELECT id, name, title_en, title_ru, title_sr
FROM c2c_categories
WHERE parent_id IS NULL
ORDER BY id;
```

## Translation Strategy

### English (title_en)
- Proper English translation of the category name
- Uses standard marketplace terminology

### Russian (title_ru)
- Russian translation
- Preserves original meaning from `name` field where applicable

### Serbian (title_sr)
- Serbian Latin translation
- Preserves original Serbian names from `name` field where already in Serbian

## Category Hierarchy

```
Root Categories (20)
├── Electronics (1001) - 10 subcategories
├── Fashion (1002) - 2 subcategories
├── Automotive (1003) - 24 subcategories
│   ├── Cars (1301) - 8 specialized subcategories
│   ├── Motorcycles (1302) - 1 subcategory
│   ├── Auto Parts (1303) - 7 subcategories
│   ├── Domestic Production (10100) - 3 subcategories
│   └── Imported Vehicles (10110) - 3 subcategories
├── Real Estate (1004) - 4 subcategories
├── Home & Garden (1005) - 2 subcategories
├── Agriculture (1006) - 4 subcategories
├── Industrial (1007) - 2 subcategories
├── Hobbies & Entertainment (1015) - 3 subcategories
└── ... other root categories
```

## Safety Features

1. **Transaction support**: All updates wrapped in BEGIN/COMMIT
2. **Backup script**: Creates snapshot before changes
3. **Verification queries**: Automatic checks after updates
4. **Rollback options**: Multiple ways to undo changes
5. **Idempotent operations**: Safe to run multiple times

## Testing

### Test on development first:
```bash
# Development database
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
```

### After verification, apply to production:
```bash
# Production database (update connection string)
psql "postgres://[PROD_CONNECTION_STRING]"
```

## Troubleshooting

### If columns already exist:
- Script uses `ADD COLUMN IF NOT EXISTS` - safe to re-run

### If backup table exists:
- Manually drop: `DROP TABLE c2c_categories_backup_20251110;`
- Or modify backup script with new date

### If some translations fail:
- Check transaction status: `SELECT txid_current();`
- Rollback if needed: `ROLLBACK;`
- Review error messages and fix specific UPDATEs

## Notes

- All scripts use transaction blocks for atomicity
- Serbian translations use Latin alphabet (not Cyrillic)
- English translations follow standard marketplace terminology
- Script execution time: ~1-2 seconds for all 77 categories

## Author

System
Date: 2025-11-10

## Related Documentation

- `/p/github.com/sveturs/listings/docs/CATEGORY_STRUCTURE.md` (if exists)
- Database schema: `c2c_categories` table definition
