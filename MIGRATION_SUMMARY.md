# Migration Summary & Examples

## Successful Migration Report

**Date:** 2025-11-10
**Duration:** ~1 second
**Success Rate:** 100%

### Records Migrated

| Table | Old DB | New DB | Status |
|-------|--------|--------|--------|
| C2C Listings | 5 | 5 | ✅ |
| B2C Stores | 1 | 1 | ✅ |
| Images | 37 | 37 | ✅ |
| **Total** | **43** | **43** | **✅** |

### Files Created

1. **Migration Script:** `/p/github.com/sveturs/listings/scripts/migrate_data.py`
2. **Verification Script:** `/p/github.com/sveturs/listings/scripts/verify_migration.sh`
3. **Rollback Script:** `/p/github.com/sveturs/listings/scripts/rollback_migration.sh`
4. **Full Documentation:** `/p/github.com/sveturs/listings/scripts/MIGRATION_GUIDE.md`
5. **Quick Start:** `/p/github.com/sveturs/listings/DATA_MIGRATION.md`

## Example Output

### Migration Log

```
╔══════════════════════════════════════════════════════════════════════════════╗
║                    DATA MIGRATION SCRIPT                                     ║
║                vondi_db → listings_dev_db                                     ║
╚══════════════════════════════════════════════════════════════════════════════╝

Started at: 2025-11-10 15:00:16

================================================================================
PHASE 1: Verifying Database Connections
================================================================================
✓ Old DB connection successful
✓ New DB connection successful

================================================================================
PHASE 2: Verifying Database Schema
================================================================================
Checking OLD database tables:
✓ Table 'c2c_listings' exists
✓ Table 'b2c_stores' exists
✓ Table 'c2c_images' exists

Checking NEW database tables:
✓ Table 'listings' exists
✓ Table 'storefronts' exists
✓ Table 'listing_images' exists

================================================================================
PHASE 3: Checking for ID Conflicts
================================================================================
✓ No conflicts in listings (old: 5, new: 10)
✓ No conflicts in storefronts (old: 1, new: 2)

================================================================================
PHASE 4: Creating Backup
================================================================================
✓ Backup created: /tmp/listings_dev_db_backup_20251110_150016.sql

================================================================================
PHASE 5: Migrating C2C Listings
================================================================================
Found 5 C2C listings to migrate
  ✓ Migrated listing 1080 → 62: Пылесос Miele
  ✓ Migrated listing 1081 → 63: Test Listing for Image Upload
  ✓ Migrated listing 1082 → 64: Цветной струйный принтер Canon
  ✓ Migrated listing 1083 → 65: Sony PS5
  ✓ Migrated listing 1084 → 66: PS5 Slim

✓ C2C Listings migration completed: 5/5 successful

================================================================================
PHASE 6: Migrating B2C Stores
================================================================================
Found 1 B2C stores to migrate
  ✓ Migrated store 43 → 3: shop

✓ B2C Stores migration completed: 1/1 successful

================================================================================
PHASE 7: Migrating C2C Images
================================================================================
Found 37 images to migrate

✓ Images migration completed: 37/37 successful

================================================================================
MIGRATION SUMMARY
================================================================================

c2c_listings:
  Total:    5
  Migrated: 5 (100.0%)
  Failed:   0

b2c_stores:
  Total:    1
  Migrated: 1 (100.0%)
  Failed:   0

c2c_images:
  Total:    37
  Migrated: 37 (100.0%)
  Failed:   0

OVERALL:
  Total records:    43
  Total migrated:   43
  Total failed:     0
  Success rate:     100.0%

ID Mappings:
  Listings:    5 mappings
  Storefronts: 1 mappings

================================================================================

✓✓✓ Migration completed successfully at 2025-11-10 15:00:16
Backup file: /tmp/listings_dev_db_backup_20251110_150016.sql
Log file: /tmp/migrate_data.log
```

## Command Examples

### Run Migration

```bash
# Basic migration
python3 /p/github.com/sveturs/listings/scripts/migrate_data.py

# View logs in real-time
tail -f /tmp/migrate_data.log
```

### Verify Results

```bash
# Run verification script
bash /p/github.com/sveturs/listings/scripts/verify_migration.sh

# Manual verification - check counts
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -c "
SELECT
  'C2C Listings' as type, COUNT(*) as count FROM listings WHERE source_type='c2c'
UNION ALL
SELECT 'Storefronts', COUNT(*) FROM storefronts
UNION ALL
SELECT 'Images', COUNT(*) FROM listing_images;
"

# Check sample data
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -c "
SELECT id, title, price, status, source_type, created_at
FROM listings
WHERE source_type = 'c2c'
ORDER BY id
LIMIT 5;
"
```

### Check Images

```bash
# Images per listing
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -c "
SELECT
  l.id,
  l.title,
  COUNT(li.id) as image_count,
  STRING_AGG(li.url, ', ' ORDER BY li.display_order) as image_urls
FROM listings l
LEFT JOIN listing_images li ON l.id = li.listing_id
WHERE l.source_type = 'c2c'
GROUP BY l.id, l.title
ORDER BY l.id;
"
```

### Rollback

```bash
# Rollback migration (requires typing 'YES')
bash /p/github.com/sveturs/listings/scripts/rollback_migration.sh

# Verify rollback
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -c "
SELECT COUNT(*) FROM listings WHERE source_type='c2c';
"
```

## Data Mapping Examples

### Listing Attributes (Old → New)

**Old DB (c2c_listings):**
```json
{
  "id": 1080,
  "title": "Пылесос Miele",
  "condition": "used",
  "address_city": "Belgrade",
  "address_country": "RS",
  "original_language": "sr",
  "show_on_map": true,
  "latitude": 44.787197,
  "longitude": 20.457273
}
```

**New DB (listings):**
```json
{
  "id": 62,
  "title": "Пылесос Miele",
  "source_type": "c2c",
  "attributes": {
    "condition": "used",
    "city": "Belgrade",
    "country": "RS",
    "original_language": "sr"
  },
  "has_individual_location": true,
  "individual_latitude": 44.787197,
  "individual_longitude": 20.457273,
  "show_on_map": true,
  "location_privacy": "exact",
  "currency": "RSD",
  "visibility": "public",
  "quantity": 1
}
```

### Image Transformation

**Old DB (c2c_images):**
```json
{
  "id": 123,
  "listing_id": 1080,
  "file_path": "/uploads/images/vacuum.jpg",
  "public_url": "https://vondi.rs/images/vacuum.jpg",
  "is_main": true,
  "display_order": 0,
  "content_type": "image/jpeg",
  "file_size": 245678
}
```

**New DB (listing_images):**
```json
{
  "id": 456,
  "listing_id": 62,
  "url": "https://vondi.rs/images/vacuum.jpg",
  "storage_path": "/uploads/images/vacuum.jpg",
  "is_primary": true,
  "display_order": 0,
  "mime_type": "image/jpeg",
  "file_size": 245678
}
```

## Post-Migration Checklist

- [x] Migration completed successfully
- [x] All records migrated (43/43)
- [x] Backup created
- [x] Verification script passed
- [ ] OpenSearch reindexed
- [ ] Application tested
- [ ] Frontend verified
- [ ] API endpoints tested
- [ ] Image URLs accessible
- [ ] Search functionality working

## Next Steps

1. **Reindex OpenSearch:**
   ```bash
   python3 /p/github.com/sveturs/listings/scripts/reindex_listings.py
   ```

2. **Test API endpoints:**
   ```bash
   # Test listing retrieval
   curl http://localhost:35433/v1/listings/62 | jq

   # Test search
   curl http://localhost:35433/v1/listings/search?q=Miele | jq
   ```

3. **Verify frontend:**
   - Browse to marketplace
   - Check listing details
   - Verify images load
   - Test search functionality

4. **Monitor logs:**
   ```bash
   # Application logs
   docker logs -f listings-service

   # Database logs
   docker logs -f listings-postgres
   ```

## Troubleshooting

### Issue: Connection refused

```bash
# Check database is running
docker ps | grep listings

# Check port
netstat -tlnp | grep 35434

# Restart if needed
docker-compose restart db
```

### Issue: Images not loading

```bash
# Check image URLs
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db" -c "
SELECT id, listing_id, url, storage_path FROM listing_images LIMIT 5;
"

# Verify MinIO/S3 storage
curl -I https://vondi.rs/images/vacuum.jpg
```

### Issue: Search not working

```bash
# Reindex OpenSearch
python3 scripts/reindex_listings.py

# Check OpenSearch indices
curl -X GET "http://localhost:9200/_cat/indices?v"

# Check document count
curl -X GET "http://localhost:9200/listings/_count"
```

## Support

- **Log file:** `/tmp/migrate_data.log`
- **Backup location:** `/tmp/listings_dev_db_backup_*.sql`
- **Full documentation:** `/p/github.com/sveturs/listings/scripts/MIGRATION_GUIDE.md`

---

**Status:** ✅ Migration Complete
**Date:** 2025-11-10
**Success Rate:** 100% (43/43 records)
