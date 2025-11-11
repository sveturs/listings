# Category Translations - File Index

> Quick reference guide to all translation-related files

## 📋 Quick Start

**Want to apply translations immediately?**
```bash
cd /p/github.com/sveturs/listings/scripts
./apply_all_translations.sh
```

**Need help?** See [TRANSLATIONS_QUICKSTART.md](TRANSLATIONS_QUICKSTART.md)

---

## 📁 Files Overview

### 🔧 SQL Scripts (Execution Order)

1. **`01_add_translation_columns.sql`** (1.1KB)
   - Adds `title_en`, `title_ru`, `title_sr` columns
   - Creates performance indexes
   - Safe to re-run (uses IF NOT EXISTS)

2. **`02_backup_categories_before_translation.sql`** (907B)
   - Creates backup table: `c2c_categories_backup_20251110`
   - Shows verification statistics
   - Always run before applying translations

3. **`03_add_category_translations.sql`** (18KB)
   - Translates all 77 categories
   - Organized by hierarchy
   - Includes verification queries
   - Transaction-safe

4. **`04_rollback_category_translations.sql`** (2.7KB)
   - 3 rollback options (all commented by default)
   - Restore from backup / Clear translations / Drop columns
   - Edit file to uncomment desired option

### 🚀 Shell Scripts

- **`apply_all_translations.sh`** (3.2KB, executable)
  - One-click solution
  - Runs all scripts in correct order
  - Color-coded output
  - Automatic verification

### 📚 Documentation

- **`TRANSLATIONS_INDEX.md`** (This file)
  - Quick navigation to all files
  - File purposes and sizes
  - Usage examples

- **`TRANSLATIONS_QUICKSTART.md`** (3.7KB)
  - Quick reference guide
  - Common commands
  - Connection strings
  - Troubleshooting

- **`README_CATEGORY_TRANSLATIONS.md`** (5.4KB)
  - Complete migration guide
  - Detailed instructions
  - Category hierarchy
  - Testing procedures

- **`TRANSLATIONS_SUMMARY.md`** (8.0KB)
  - Implementation summary
  - Translation breakdown
  - Verification results
  - Next steps for integration

---

## 🎯 Common Tasks

### Apply All Translations
```bash
./apply_all_translations.sh
```

### Apply Individual Scripts
```bash
DB="postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"

# Step by step
psql "$DB" -f 01_add_translation_columns.sql
psql "$DB" -f 02_backup_categories_before_translation.sql
psql "$DB" -f 03_add_category_translations.sql
```

### Verify Translations
```bash
DB="postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"

psql "$DB" -c "
  SELECT COUNT(*) as total,
         COUNT(title_en) as en,
         COUNT(title_ru) as ru,
         COUNT(title_sr) as sr
  FROM c2c_categories;
"
```

### Rollback Changes
```bash
# First, edit 04_rollback_category_translations.sql
# Uncomment your preferred rollback option
nano 04_rollback_category_translations.sql

# Then run it
psql "$DB" -f 04_rollback_category_translations.sql
```

---

## 📊 Statistics

**Database:** `listings_dev_db`
**Table:** `c2c_categories`
**Total Categories:** 77
**Translation Coverage:** 100%

**Breakdown:**
- Root categories (level 0): 23
- First-level subcategories (level 1): 29
- Second-level subcategories (level 2): 25

**Languages:**
- English (title_en): 77/77 ✅
- Russian (title_ru): 77/77 ✅
- Serbian (title_sr): 77/77 ✅

---

## 🔗 Quick Links

| Document | Purpose |
|----------|---------|
| [QUICKSTART](TRANSLATIONS_QUICKSTART.md) | Fast implementation guide |
| [README](README_CATEGORY_TRANSLATIONS.md) | Complete documentation |
| [SUMMARY](TRANSLATIONS_SUMMARY.md) | Implementation details |
| [INDEX](TRANSLATIONS_INDEX.md) | This file |

---

## 🔍 File Locations

All files located in:
```
/p/github.com/sveturs/listings/scripts/
```

### SQL Scripts
```
01_add_translation_columns.sql
02_backup_categories_before_translation.sql
03_add_category_translations.sql
04_rollback_category_translations.sql
```

### Shell Scripts
```
apply_all_translations.sh
```

### Documentation
```
TRANSLATIONS_INDEX.md
TRANSLATIONS_QUICKSTART.md
README_CATEGORY_TRANSLATIONS.md
TRANSLATIONS_SUMMARY.md
```

---

## 💡 Tips

1. **Always create backup first** - Run script 02 before 03
2. **Use the all-in-one script** - `apply_all_translations.sh` handles everything
3. **Verify after applying** - Check the output for any missing translations
4. **Keep backup table** - Don't drop `c2c_categories_backup_20251110` until verified in production
5. **Test on dev first** - Always test on development database before production

---

## ⚠️ Important Notes

- **Transaction Safety:** All operations are wrapped in BEGIN/COMMIT
- **Idempotency:** Safe to re-run scripts multiple times
- **Backup Required:** Always create backup before applying translations
- **Rollback Available:** Multiple rollback options provided
- **Production Ready:** All scripts tested on development database

---

## 📞 Support

For issues or questions:
1. Check [TRANSLATIONS_QUICKSTART.md](TRANSLATIONS_QUICKSTART.md#troubleshooting)
2. Review [README_CATEGORY_TRANSLATIONS.md](README_CATEGORY_TRANSLATIONS.md)
3. Verify database connection
4. Check backup table exists

---

**Last Updated:** 2025-11-10
**Status:** ✅ Production Ready
**Version:** 1.0
