# Category Translations - Quick Start Guide

## TL;DR - Quick Apply

```bash
cd /p/github.com/sveturs/listings/scripts
./apply_all_translations.sh
```

This will:
1. Add translation columns (title_en, title_ru, title_sr)
2. Create backup table
3. Apply translations to all 77 categories
4. Show verification results

## Individual Scripts

### 1. Add Columns
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f 01_add_translation_columns.sql
```

### 2. Create Backup
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f 02_backup_categories_before_translation.sql
```

### 3. Apply Translations
```bash
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f 03_add_category_translations.sql
```

### 4. Rollback (if needed)
```bash
# Edit the file first to uncomment desired rollback option
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f 04_rollback_category_translations.sql
```

## Quick Verification

```bash
# Check if all categories have translations
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) as total, COUNT(title_en) as en, COUNT(title_ru) as ru, COUNT(title_sr) as sr FROM c2c_categories;"
```

Expected result:
```
 total | en | ru | sr
-------+----+----+----
    77 | 77 | 77 | 77
```

## View Sample Translations

```bash
# Root categories
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT id, name, title_en, title_ru, title_sr FROM c2c_categories WHERE parent_id IS NULL LIMIT 5;"
```

## Files Created

- `01_add_translation_columns.sql` - Adds columns and indexes
- `02_backup_categories_before_translation.sql` - Creates backup
- `03_add_category_translations.sql` - All 77 translations
- `04_rollback_category_translations.sql` - Rollback options
- `apply_all_translations.sh` - One-click apply script
- `README_CATEGORY_TRANSLATIONS.md` - Full documentation

## Connection Strings

**Development:**
```
postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable
```

**Production:**
```
# Update this with production credentials
postgres://[USER]:[PASSWORD]@[HOST]:[PORT]/[DATABASE]?sslmode=disable
```

## Safety Features

- All scripts use transactions (BEGIN/COMMIT)
- Backup table created before changes
- Idempotent operations (safe to re-run)
- Verification queries included
- Multiple rollback options

## Translation Examples

| Category | English | Russian | Serbian |
|----------|---------|---------|---------|
| Electronics | Electronics | Электроника | Elektronika |
| Automotive | Automotive | Автомобили | Automobili |
| Real Estate | Real Estate | Недвижимость | Nekretnine |
| Electric Cars | Electric Cars | Электромобили | Električni automobili |
| SUV Vehicles | SUV Vehicles | Внедорожники | SUV vozila |

## Troubleshooting

**Columns already exist?**
- Safe to re-run, uses `IF NOT EXISTS`

**Backup table exists?**
- Drop manually: `DROP TABLE c2c_categories_backup_20251110;`

**Need to revert?**
- Edit `04_rollback_category_translations.sql`
- Uncomment desired option
- Run the script

## Next Steps

After applying translations:
1. Update application code to use new columns
2. Add i18n support in frontend
3. Update API to return localized titles
4. Test with different locales

## Support

For detailed documentation see:
- `README_CATEGORY_TRANSLATIONS.md` - Full guide
- `/p/github.com/sveturs/listings/docs/` - Project docs
