# 🔍 OpenSearch Reindexing Instructions

**Phase:** 13.1.7 Post-Deployment  
**Purpose:** Index new schema fields (attributes, stock_status, view_count, etc.)

---

## 📋 When to Reindex

Reindex OpenSearch **AFTER** deploying Phase 13.1.7:

1. ✅ Migrations applied (000012, 000013, 000014)
2. ✅ Application deployed with new code
3. ⏳ OpenSearch needs reindexing

---

## 🚀 Reindexing Options

### Option 1: Using Production Script (Recommended)

```bash
cd /p/github.com/sveturs/listings

python3 scripts/reindex_listings.py \
  --target-password YOUR_DB_PASSWORD \
  --target-host localhost \
  --target-port 5432 \
  --target-db vondi_db \
  --opensearch-url http://localhost:9200 \
  --opensearch-index marketplace_listings \
  --batch-size 1000 \
  --verbose
```

**Parameters:**
- `--target-password` - PostgreSQL password (**required**)
- `--target-host` - Database host (default: localhost)
- `--target-port` - Database port (default: 5432)
- `--target-db` - Database name (default: listings_db)
- `--opensearch-url` - OpenSearch URL (default: http://localhost:9200)
- `--opensearch-index` - Index name (default: marketplace_listings)
- `--batch-size` - Documents per batch (default: 1000)
- `--verbose` - Detailed logging
- `--dry-run` - Test without actually indexing

---

### Option 2: Direct Database Query (Quick Test)

```bash
# Connect to database
psql "postgres://postgres:YOUR_PASSWORD@localhost:5432/vondi_db?sslmode=disable"

# Check new columns exist
\d+ listings

# Sample query to verify data
SELECT 
  id, 
  title, 
  view_count,  -- was views_count
  sold_count,   -- NEW
  stock_status, -- NEW
  attributes    -- NEW (JSONB)
FROM listings 
WHERE source_type = 'b2c' 
LIMIT 5;
```

---

### Option 3: Docker-based Reindex

```bash
# If using Docker
python3 scripts/reindex_via_docker.py
```

---

## 🔍 Verification

### Check Index Status

```bash
# Count documents
curl -X GET "http://localhost:9200/marketplace_listings/_count" | jq '.'

# Check mapping includes new fields
curl -X GET "http://localhost:9200/marketplace_listings/_mapping" | jq '.marketplace_listings.mappings.properties' | grep -E 'stock_status|attributes|view_count|sold_count'

# Sample document with new fields
curl -X GET "http://localhost:9200/marketplace_listings/_search?size=1" | jq '.hits.hits[0]._source | {id, stock_status, attributes, view_count, sold_count}'
```

### Expected New Fields in Index

```json
{
  "_source": {
    "id": 328,
    "title": "Product Name",
    "view_count": 150,        // ✅ Renamed from views_count
    "sold_count": 25,         // ✅ NEW
    "stock_status": "in_stock", // ✅ NEW
    "attributes": {           // ✅ NEW (JSONB)
      "color": "blue",
      "size": "L"
    },
    "source_type": "b2c"      // ✅ Used for filtering
  }
}
```

---

## ⏱️ Performance

**Estimated Duration:**
- Small dataset (<10k documents): 5-10 minutes
- Medium dataset (10k-100k): 15-30 minutes
- Large dataset (>100k): 30-60 minutes

**Monitor Progress:**
```bash
# Watch reindex log
tail -f /tmp/opensearch_reindex.log

# Check indexed count during process
watch -n 5 'curl -s "http://localhost:9200/marketplace_listings/_count" | jq ".count"'
```

---

## ⚠️ Important Notes

### Before Reindexing:

1. **Backup current index:**
   ```bash
   # Create snapshot (if configured)
   curl -X PUT "http://localhost:9200/_snapshot/my_backup/snapshot_$(date +%Y%m%d)"
   ```

2. **Check disk space:**
   ```bash
   df -h
   # OpenSearch needs ~2x current index size free
   ```

3. **Plan for downtime (if needed):**
   - Read-only mode: Search still works
   - Full reindex: Brief search unavailability

### During Reindexing:

- ✅ Application continues to work
- ✅ New writes are handled
- ⚠️ Search may be incomplete until done
- ⚠️ Don't restart OpenSearch

### After Reindexing:

1. **Verify counts match:**
   ```sql
   -- Database count
   SELECT COUNT(*) FROM listings WHERE source_type = 'b2c';
   
   -- OpenSearch count (should match)
   curl "http://localhost:9200/marketplace_listings/_count?q=source_type:b2c"
   ```

2. **Test search with new fields:**
   ```bash
   # Search by stock_status
   curl -X GET "http://localhost:9200/marketplace_listings/_search" -H 'Content-Type: application/json' -d'
   {
     "query": {
       "term": { "stock_status": "in_stock" }
     }
   }'
   
   # Search in attributes (JSONB)
   curl -X GET "http://localhost:9200/marketplace_listings/_search" -H 'Content-Type: application/json' -d'
   {
     "query": {
       "match": { "attributes.color": "blue" }
     }
   }'
   ```

---

## 🐛 Troubleshooting

### Issue: Reindex script fails with connection error

**Solution:**
```bash
# Check PostgreSQL is running
systemctl status postgresql
# Check OpenSearch is running  
systemctl status opensearch
# Verify connectivity
psql "postgres://postgres:password@localhost:5432/vondi_db" -c "SELECT 1;"
curl http://localhost:9200
```

### Issue: Some documents missing new fields

**Cause:** Documents indexed before migration

**Solution:** Full reindex required
```bash
python3 scripts/reindex_listings.py --target-password PASSWORD
```

### Issue: Search returns unexpected results

**Cause:** Mapping might be cached

**Solution:** Refresh index
```bash
curl -X POST "http://localhost:9200/marketplace_listings/_refresh"
```

---

## 📚 Related Documentation

- [DEPLOYMENT_GUIDE_13_1_7.md](./DEPLOYMENT_GUIDE_13_1_7.md) - Full deployment guide
- [PHASE_13_1_7_FINAL_REPORT.md](./PHASE_13_1_7_FINAL_REPORT.md) - Technical details
- Migration 000014 - Complete list of new columns

---

## ✅ Success Criteria

Reindexing is successful when:

- [x] Document count matches database
- [x] All new fields present in mapping
- [x] Sample queries return expected data
- [x] No errors in reindex log
- [x] Application search works correctly

---

**Document Created:** 2025-11-08  
**Last Updated:** 2025-11-08  
**Status:** Ready for Production Use

---

**END OF REINDEX INSTRUCTIONS**
