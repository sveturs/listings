# Listings Microservice Scripts

Comprehensive automation tools for the Listings microservice, including production deployment, monitoring, and profiling.

---

## Table of Contents

1. [Data Migration](#data-migration)
2. [Production Deployment](#production-deployment)
3. [Memory Profiling](#memory-profiling)
4. [Monitoring & Testing](#monitoring--testing)

---

# Data Migration

Production-ready data migration from old monolith database (svetubd) to new microservice database (listings_dev_db).

## Quick Start - Migration

```bash
# 1. Run migration
python3 scripts/migrate_data.py

# 2. Verify migration
bash scripts/verify_migration.sh

# 3. Rollback if needed
bash scripts/rollback_migration.sh
```

## Migration Scripts

### 🔄 `migrate_data.py`

**Purpose:** Migrate production data from old database to new microservice database

**What it migrates:**
- C2C Listings (5 records) → `listings` table with `source_type='c2c'`
- B2C Stores (1 record) → `storefronts` table
- Images (37 records) → `listing_images` table

**Features:**
- ✅ Pre-flight checks (connectivity, schema, conflicts)
- ✅ Automatic backup before migration
- ✅ Transaction support (rollback on error)
- ✅ ID mapping for foreign keys
- ✅ Detailed logging to `/tmp/migrate_data.log`
- ✅ Error resilience (continues on individual failures)
- ✅ Progress tracking in real-time

**Usage:**
```bash
python3 scripts/migrate_data.py
```

**Output:**
- Migration log: `/tmp/migrate_data.log`
- Backup: `/tmp/listings_dev_db_backup_YYYYMMDD_HHMMSS.sql`

### ✅ `verify_migration.sh`

**Purpose:** Verify migration completed successfully

**What it checks:**
- Database connectivity
- Record counts (old vs new)
- Sample data integrity
- Foreign key relationships
- Data quality (NULL values, price ranges)
- Image distribution
- Status mapping correctness
- Attributes migration
- Location data

**Usage:**
```bash
bash scripts/verify_migration.sh
```

**Example output:**
```
✓ C2C Listings: 5 (old) = 5 (new)
✓ B2C Stores: 1 (old) = 1 (new)
✓ Images: 37 (old) = 37 (new)
✓ All images reference existing listings
✓ All listings have valid titles
```

### ⏪ `rollback_migration.sh`

**Purpose:** Rollback migration (delete migrated data)

**What it does:**
1. Creates backup before rollback
2. Deletes all C2C listings (`source_type='c2c'`)
3. Deletes related images
4. Deletes related attributes/locations/tags
5. Shows before/after counts

**Safety:**
- ⚠️ Requires explicit confirmation (type 'YES')
- ✅ Creates backup before deletion
- ⚠️ Does NOT delete B2C storefronts (manual cleanup)

**Usage:**
```bash
bash scripts/rollback_migration.sh
```

## Data Transformations

### C2C Listings → Listings Table

**Direct mappings:**
- `user_id`, `title`, `description`, `price`, `category_id`
- `status` (validated: active/sold/inactive/archived/draft)
- `views_count` → `view_count`

**New fields:**
- `source_type` = `'c2c'` (identifies source)
- `currency` = `'RSD'` (default)
- `visibility` = `'public'` (default)
- `quantity` = `1` (default)

**Location fields:**
- `show_on_map` → `has_individual_location` + `show_on_map`
- `location` → `individual_address`
- `latitude`, `longitude` → `individual_latitude`, `individual_longitude`
- `location_privacy` = `'exact'` (default)

**Attributes (JSONB):**
- Old fields → JSON: `condition`, `address_city`, `address_country`, `original_language`, `metadata`, `address_multilingual`

### B2C Stores → Storefronts Table

**Direct 1:1 mapping** of all fields with JSONB conversion:
- `theme`, `settings`, `seo_meta`, `ai_agent_config` (dict → JSON string)

### C2C Images → Listing Images

**Mappings:**
- `file_path` → `storage_path`
- `public_url` or `file_path` → `url`
- `is_main` → `is_primary`
- `content_type` → `mime_type`
- `display_order` (preserved)

## Database Credentials

**Old DB (port 5433):**
```
postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd
```

**New DB (port 35434):**
```
postgres://listings_user:listings_secret@localhost:35434/listings_dev_db
```

## Post-Migration Steps

1. **Verify migration:**
   ```bash
   bash scripts/verify_migration.sh
   ```

2. **Update OpenSearch indices:**
   ```bash
   python3 scripts/reindex_listings.py
   ```

3. **Test application:**
   - Start services
   - Test API endpoints
   - Verify frontend displays data correctly

## Troubleshooting

**Connection errors:**
```bash
# Check database running
docker ps | grep listings

# Test connection
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"
```

**Python dependencies:**
```bash
pip3 install psycopg2-binary
```

**View logs:**
```bash
tail -f /tmp/migrate_data.log
```

## Documentation

- **Quick Start:** `/p/github.com/sveturs/listings/DATA_MIGRATION.md`
- **Full Guide:** `/p/github.com/sveturs/listings/scripts/MIGRATION_GUIDE.md`

---

## Attributes Migration (unified_attributes → attributes)

Migration of attribute definitions from monolith to microservice with i18n JSONB conversion.

### Quick Start - Attributes

```bash
# Recommended: Use Go migration tool
cd /p/github.com/sveturs/listings

# Step 1: Preview migration (dry-run)
go run ./cmd/migrate_attributes/main.go --dry-run

# Step 2: Execute migration
go run ./cmd/migrate_attributes/main.go -v

# Step 3: Validate results
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f ./scripts/validate_attributes.sql
```

### Key Features

✅ **Idempotent:** Safe to run multiple times - skips existing records
✅ **JSONB Conversion:** VARCHAR → `{"en", "ru", "sr"}` i18n format
✅ **Transaction Support:** All-or-nothing migration
✅ **Built-in Validation:** Automatic post-migration checks
✅ **Dry-run Mode:** Preview before executing

### Migration Scripts

| File | Purpose |
|------|---------|
| `../cmd/migrate_attributes/main.go` | Automated Go migration tool ⭐ (recommended) |
| `migrate_attributes.sql` | Manual SQL migration (requires CSV export) |
| `validate_attributes.sql` | Comprehensive validation queries |
| `rollback_attributes.sql` | Emergency rollback (deletes all attributes) |

### Migration Results (2025-11-17)

**Status:** ✅ Completed successfully
- **Records:** 203 attributes migrated
- **Time:** ~1 second
- **Idempotency:** Verified (safe to re-run)

### Validation Summary

- ✅ All 203 records migrated
- ✅ JSONB structure correct (`en`, `ru`, `sr` keys)
- ✅ No NULL values in required fields
- ✅ No duplicate codes
- ✅ Sequence updated correctly (current: 549)
- ✅ Search vectors generated for all records
- ✅ Attribute type distribution: select(83), number(45), text(34), boolean(19), multiselect(16), date(5), textarea(1)

### Data Transformation

**Source:** `svetubd.unified_attributes` (port 5433)
**Target:** `listings_dev_db.attributes` (port 35434)

**Key Changes:**
- `name` (VARCHAR) → `name` (JSONB) with `{en, ru, sr}` keys
- `display_name` (VARCHAR) → `display_name` (JSONB) with `{en, ru, sr}` keys

**Example:**
```sql
-- Before (monolith)
name = 'year'
display_name = 'Godište'

-- After (microservice)
name = {"en": "year", "ru": "year", "sr": "year"}
display_name = {"en": "Godište", "ru": "Godište", "sr": "Godište"}
```

### Troubleshooting

**Issue:** "duplicate key violation"
```bash
# Migration is idempotent - this is expected if re-running
# Script will skip existing records automatically
go run ./cmd/migrate_attributes/main.go  # Will show "Already migrated: 203"
```

**Issue:** Need to rollback
```bash
# Use interactive rollback script
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -f ./scripts/rollback_attributes.sql
```

**Issue:** Verify specific attribute
```sql
-- Check JSONB structure
SELECT id, code, name, display_name
FROM attributes
WHERE code = 'year';
```

### Documentation

📚 **Full guide:** [../docs/ATTRIBUTES_MIGRATION_GUIDE.md](../docs/ATTRIBUTES_MIGRATION_GUIDE.md)

---

# Production Deployment

Production-ready deployment automation using Blue-Green strategy with zero downtime.

## Quick Start - Deployment

```bash
# 1. Configure deployment
cp .env.deploy.example .env.deploy
nano .env.deploy

# 2. Validate environment
./validate-deployment.sh

# 3. Deploy to production
./deploy-to-prod.sh
```

## Deployment Scripts

### 📊 `profile_memory.sh`

**Purpose:** Automated memory leak detection with before/after comparison

**What it does:**
1. Captures baseline heap and goroutine profiles
2. Generates load for 60 seconds (configurable)
3. Captures post-load profiles
4. Analyzes memory growth
5. Generates reports and flamegraphs

**Usage:**
```bash
./scripts/profile_memory.sh
```

**Output:**
```
/tmp/memory_profiles_YYYYMMDD_HHMMSS/
├── heap_baseline.pprof
├── heap_postload.pprof
├── heap_diff.txt
├── goroutine_baseline.pprof
├── goroutine_postload.pprof
├── goroutine_diff.txt
└── heap_flamegraph.svg (if go-torch available)
```

**Customize:**
```bash
# Edit DURATION variable in script for longer tests
DURATION="300s"  # 5 minutes
```

---

### 📈 `monitor_memory.sh`

**Purpose:** Continuous real-time memory monitoring

**What it does:**
- Monitors heap, goroutines, DB connections every 10 seconds
- Saves timestamped data to CSV
- Displays real-time metrics in terminal

**Usage:**
```bash
./scripts/monitor_memory.sh
```

**Output:**
```
/tmp/memory_monitoring_YYYYMMDD_HHMMSS.csv
```

**CSV Columns:**
- `timestamp` - Unix timestamp
- `heap_alloc_mb` - Heap allocation (MB)
- `heap_sys_mb` - Heap system memory (MB)
- `num_gc` - Total GC runs
- `goroutines` - Active goroutine count
- `db_connections_open` - Open DB connections
- `db_connections_in_use` - In-use DB connections

**Example:**
```bash
# Start monitoring in background
./scripts/monitor_memory.sh &
MONITOR_PID=$!

# Generate load for 30 minutes
# ... your load tests ...

# Stop monitoring
kill $MONITOR_PID
```

---

### 🔍 `detect_leaks.py`

**Purpose:** Automated leak detection from monitoring data

**What it does:**
- Analyzes CSV data from `monitor_memory.sh`
- Detects heap growth, goroutine leaks, connection leaks
- Calculates growth rates and trends
- Generates detailed report

**Usage:**
```bash
python3 ./scripts/detect_leaks.py <csv_file> [threshold_mb]
```

**Examples:**
```bash
# Default threshold (50 MB)
python3 ./scripts/detect_leaks.py /tmp/memory_monitoring_20250104_120000.csv

# Custom threshold (100 MB)
python3 ./scripts/detect_leaks.py /tmp/memory_monitoring_20250104_120000.csv 100
```

**Exit codes:**
- `0` - No leaks detected
- `1` - Leaks detected

**Detection criteria:**
- Heap growth > 50 MB (or custom threshold)
- Heap growth rate > 1 MB/min sustained
- Goroutine growth > 100
- DB connection growth > 10

---

## Common Workflows

### Workflow 1: Quick Check (5 minutes)

```bash
# 1. Run automated profiling
./scripts/profile_memory.sh

# 2. Review output
cat /tmp/memory_profiles_*/heap_diff.txt
cat /tmp/memory_profiles_*/goroutine_diff.txt
```

### Workflow 2: Soak Test (30+ minutes)

```bash
# 1. Start continuous monitoring
./scripts/monitor_memory.sh &
MONITOR_PID=$!

# 2. Generate sustained load
ghz --insecure \
    --proto="api/proto/listings/v1/listings.proto" \
    --call=listings.v1.ListingsService.GetListing \
    -d '{"id": 328}' \
    -c 100 \
    -z 30m \
    --rps 5000 \
    localhost:50051

# 3. Stop monitoring
kill $MONITOR_PID

# 4. Analyze results
CSV_FILE=$(ls -t /tmp/memory_monitoring_*.csv | head -1)
python3 ./scripts/detect_leaks.py $CSV_FILE
```

### Workflow 3: Production Investigation

```bash
# If production shows high memory usage:

# 1. Capture snapshot
curl http://localhost:6060/debug/pprof/heap > production_leak.pprof

# 2. Analyze top allocations
go tool pprof -top production_leak.pprof

# 3. Interactive analysis
go tool pprof -http=:8081 production_leak.pprof
# Open http://localhost:8081 in browser

# 4. Check goroutines
curl http://localhost:6060/debug/pprof/goroutine > goroutines.pprof
go tool pprof -top goroutines.pprof
```

---

## Manual pprof Commands

### Heap Analysis

```bash
# Capture current heap
curl http://localhost:6060/debug/pprof/heap > heap.pprof

# View top memory consumers
go tool pprof -top heap.pprof

# Show allocation call graph
go tool pprof -web heap.pprof

# Interactive mode
go tool pprof heap.pprof
# Commands: top, list <function>, web, pdf
```

### Goroutine Analysis

```bash
# Capture goroutine profile
curl http://localhost:6060/debug/pprof/goroutine > goroutine.pprof

# View goroutine stacks
go tool pprof -top goroutine.pprof

# Show goroutine graph
go tool pprof -web goroutine.pprof
```

### Compare Profiles

```bash
# Capture before load
curl http://localhost:6060/debug/pprof/heap > before.pprof

# ... generate load ...

# Capture after load
curl http://localhost:6060/debug/pprof/heap > after.pprof

# Compare
go tool pprof -base before.pprof after.pprof
```

---

## Load Testing

### Using ghz (gRPC)

```bash
# Install ghz
go install github.com/bojand/ghz/cmd/ghz@latest

# Load test GetListing endpoint
ghz --insecure \
    --proto="api/proto/listings/v1/listings.proto" \
    --call=listings.v1.ListingsService.GetListing \
    -d '{"id": 328}' \
    -c 50 \          # 50 concurrent connections
    -z 5m \          # 5 minutes duration
    --rps 1000 \     # 1000 requests per second
    localhost:50051
```

### Using curl (HTTP)

```bash
# Simple load test
for i in {1..10000}; do
    curl -s http://localhost:8080/health > /dev/null &
done
wait

# Sustained load
while true; do
    curl -s http://localhost:8080/health > /dev/null &
    sleep 0.01
done
```

---

## Troubleshooting

### Script doesn't run

```bash
# Make sure scripts are executable
chmod +x scripts/*.sh scripts/*.py

# Check if pprof is accessible
curl http://localhost:6060/debug/pprof/
```

### No ghz available

The script falls back to curl-based load testing if `ghz` is not installed.

To install ghz:
```bash
go install github.com/bojand/ghz/cmd/ghz@latest
```

### pprof server not accessible

Check if service is running with pprof enabled:
```bash
# Service should log on startup:
# INFO Starting pprof server addr=:6060

# Check if port is open
nc -z localhost 6060
```

### Python script errors

```bash
# Make sure Python 3 is installed
python3 --version

# No additional dependencies required
# Script uses only standard library
```

---

## Best Practices

### ✅ DO

- Run profiling during realistic load
- Monitor for at least 30 minutes under sustained load
- Capture profiles before and after load
- Use automated scripts for consistent results
- Save all profile outputs for comparison
- Run regular leak detection in CI/CD

### ❌ DON'T

- Profile in production without planning (can add overhead)
- Ignore small leaks (they compound over time)
- Skip goroutine profiling (goroutines are cheap but can leak)
- Run profiling with no load (won't show real-world behavior)
- Delete profile files immediately (keep for comparison)

---

## References

- [Go pprof Documentation](https://pkg.go.dev/net/http/pprof)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [Memory Leak Report](../docs/MEMORY_LEAK_REPORT.md)
- [ghz - gRPC Load Testing](https://github.com/bojand/ghz)

---

## Support

For questions or issues with profiling:

1. Check the [Memory Leak Report](../docs/MEMORY_LEAK_REPORT.md)
2. Review script output for errors
3. Verify pprof endpoint is accessible
4. Check service logs for errors
