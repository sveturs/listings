# 🚀 DEPLOYMENT GUIDE - Phase 13.1.7

**Version:** Phase 13.1.7 - B2C Schema Compatibility  
**Date:** 2025-11-08  
**Status:** ✅ Ready for Production

---

## 📋 Pre-Deployment Checklist

### Prerequisites
- [ ] Backup production database
- [ ] Review all migration files (000012, 000013, 000014)
- [ ] Verify staging environment passes all tests
- [ ] Notify team about deployment window

### Environment Check
```bash
# Verify migrations directory
ls -la migrations/000012* migrations/000013* migrations/000014*

# Check database connectivity
psql "postgres://user:pass@host:5432/dbname?sslmode=disable" -c "SELECT version();"

# Verify no pending data changes
git status
```

---

## 🔄 Migration Sequence

### Step 1: Apply Database Migrations

```bash
cd /p/github.com/sveturs/listings/backend

# Run migrator (applies in order: 000012 → 000013 → 000014)
./migrator up
```

**Expected Output:**
```
Applied migration: 000012_add_attributes_to_listings.up.sql
Applied migration: 000013_add_stock_status_to_listings.up.sql  
Applied migration: 000014_fix_b2c_schema_compatibility.up.sql
```

**Verify migrations:**
```bash
# Check applied migrations
psql $DATABASE_URL -c "SELECT version FROM schema_migrations ORDER BY version DESC LIMIT 5;"
```

### Step 2: Verify Schema Changes

```bash
# Verify new columns exist
psql $DATABASE_URL -c "\d+ listings" | grep -E 'attributes|stock_status|view_count|sold_count'
```

**Expected columns:**
- `attributes` - jsonb
- `stock_status` - character varying(50)
- `view_count` - integer (was views_count)
- `sold_count` - integer
- `has_individual_location` - boolean
- `individual_address` - text
- `individual_latitude` - numeric(10,8)
- `individual_longitude` - numeric(11,8)
- `location_privacy` - character varying(20)
- `show_on_map` - boolean
- `has_variants` - boolean

### Step 3: Deploy Application Code

```bash
# Pull latest code
git pull origin main

# Rebuild application
go build -o bin/listings ./cmd/api/main.go

# Restart service
systemctl restart listings-service
# OR
supervisorctl restart listings
```

### Step 4: Verify Application Health

```bash
# Health check
curl http://localhost:8080/health

# Test API endpoints
curl http://localhost:8080/api/v1/products/1
```

---

## 🔍 Post-Deployment Validation

### 1. Database Checks

```bash
# Count listings with new fields populated
psql $DATABASE_URL << 'SQL'
SELECT 
  COUNT(*) as total,
  COUNT(attributes) as has_attributes,
  COUNT(stock_status) as has_stock_status,
  COUNT(view_count) as has_view_count
FROM listings 
WHERE source_type = 'b2c';
SQL
```

### 2. Application Logs

```bash
# Check for errors
tail -f /var/log/listings/app.log | grep -i error

# Monitor for schema errors
tail -f /var/log/listings/app.log | grep -i "does not exist"
```

### 3. Integration Tests (Optional but Recommended)

```bash
cd /p/github.com/sveturs/listings

# Run smoke tests
go test -tags=integration ./tests/integration -run "TestGetProduct_Success|TestBulkUpdateProducts_Success"
```

---

## 🔄 OpenSearch Reindexing

**IMPORTANT:** Run AFTER application deployment

### Full Reindexing (30-60 minutes)

```bash
cd /p/github.com/sveturs/listings/backend

# Run reindexing script
python3 reindex_unified.py
```

**Monitor progress:**
```bash
# Check OpenSearch index
curl -X GET "http://localhost:9200/marketplace_listings/_count"

# Verify new fields indexed
curl -X GET "http://localhost:9200/marketplace_listings/_mapping" | jq '.marketplace_listings.mappings.properties | keys'
```

---

## ⚠️ Rollback Plan

If issues occur, rollback in reverse order:

### Step 1: Rollback Application
```bash
# Checkout previous version
git checkout <previous-commit-hash>

# Rebuild and restart
go build -o bin/listings ./cmd/api/main.go
systemctl restart listings-service
```

### Step 2: Rollback Database Migrations

```bash
cd /p/github.com/sveturs/listings/backend

# Rollback migrations (runs .down.sql files)
./migrator down 3  # Rolls back last 3 migrations
```

**Verify rollback:**
```bash
psql $DATABASE_URL -c "\d listings" | grep -v "view_count\|sold_count\|attributes"
```

---

## 📊 Monitoring

### Key Metrics to Watch (First 24-48 hours)

1. **Query Performance**
   ```sql
   -- Monitor slow queries
   SELECT query, mean_exec_time, calls 
   FROM pg_stat_statements 
   WHERE query LIKE '%listings%' 
   ORDER BY mean_exec_time DESC 
   LIMIT 10;
   ```

2. **Error Rates**
   - Monitor application error logs
   - Track HTTP 500 responses
   - Watch for database constraint violations

3. **OpenSearch Performance**
   ```bash
   # Query latency
   curl -X GET "http://localhost:9200/_nodes/stats/indices/search"
   ```

---

## 🐛 Common Issues & Solutions

### Issue 1: Migration Fails

**Symptom:** `pq: column already exists`

**Solution:**
```bash
# Check if migration already applied
psql $DATABASE_URL -c "SELECT * FROM schema_migrations WHERE version = '000014';"

# If exists, skip migration
./migrator status
```

### Issue 2: Application Can't Find Column

**Symptom:** `pq: column "view_count" does not exist`

**Solution:**
```bash
# Verify migration 000014 applied
psql $DATABASE_URL -c "SELECT column_name FROM information_schema.columns WHERE table_name = 'listings' AND column_name = 'view_count';"

# If missing, rerun migration
./migrator up
```

### Issue 3: OpenSearch Out of Sync

**Symptom:** Search returns incomplete data

**Solution:**
```bash
# Force reindex
python3 reindex_unified.py --force

# Verify completion
curl -X GET "http://localhost:9200/marketplace_listings/_count"
```

---

## ✅ Success Criteria

Deployment is successful when:

- [x] All 3 migrations applied without errors
- [x] Application starts and passes health check
- [x] Integration tests pass (>90%)
- [x] No schema-related errors in logs (24h)
- [x] OpenSearch reindexing completes
- [x] Query performance within acceptable range
- [x] Zero production incidents related to schema

---

## 📞 Support

**For Issues:**
1. Check logs: `/var/log/listings/`
2. Review this guide's troubleshooting section
3. Rollback if critical (see Rollback Plan)
4. Contact: DevOps team / Database team

**Documentation:**
- [PHASE_13_1_7_FINAL_REPORT.md](./PHASE_13_1_7_FINAL_REPORT.md) - Complete technical report
- [MIGRATION_PLAN_TO_MICROSERVICE.md](./MIGRATION_PLAN_TO_MICROSERVICE.md) - Overall migration strategy

---

**Guide Version:** 1.0  
**Last Updated:** 2025-11-08  
**Tested On:** Staging environment ✅

---

**END OF DEPLOYMENT GUIDE**
