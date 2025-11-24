# Listings Microservice - Migration Cheatsheet

Quick reference for all data migrations from monolith to microservice.

---

## 🎯 Migration Order

**IMPORTANT:** Migrations must be done in this exact order due to foreign key dependencies!

1. ✅ **Categories** → Foundational data
2. ✅ **Attributes** → Referenced by category_attributes
3. ✅ **Category Attributes** → Links categories to attributes
4. ⏳ **Listing Values** → Actual listing attribute data
5. ⏳ **Orders** → Order data
6. ⏳ **Cart Items** → Shopping cart data

---

## 📊 Current Status

| Migration           | Status | Records | Command                                          |
|---------------------|--------|---------|--------------------------------------------------|
| Categories          | ✅ Done | N/A     | Manual (seed data)                               |
| Attributes          | ✅ Done | 157     | `go run ./cmd/migrate_attributes/main.go`        |
| Category Attributes | ✅ Done | 479     | `go run ./cmd/migrate_category_attributes/main.go` |
| Listing Values      | ⏳ TODO | ~50k?   | TBD                                              |
| Orders              | ⏳ TODO | TBD     | TBD                                              |
| Cart Items          | ⏳ TODO | TBD     | TBD                                              |

---

## 🚀 Quick Commands

### 1. Attributes Migration

```bash
# Dry-run
cd /p/github.com/sveturs/listings
go run ./cmd/migrate_attributes/main.go --dry-run

# Real migration
go run ./cmd/migrate_attributes/main.go

# Validate
./scripts/validate_attributes_migration.sh
```

### 2. Category Attributes Migration

```bash
# Dry-run
cd /p/github.com/sveturs/listings
go run ./cmd/migrate_category_attributes/main.go --dry-run

# Real migration
go run ./cmd/migrate_category_attributes/main.go

# Validate
./scripts/validate_category_attributes_migration.sh
```

---

## 🔍 Validation Quick Checks

### Check Record Counts

```bash
# Source (Monolith) - port 5433
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" \
  -c "SELECT COUNT(*) FROM unified_attributes;"

# Destination (Microservice) - port 35434
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM attributes;"
```

### Compare Data

```bash
# Check specific attribute
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT * FROM attributes WHERE id = 90;"

# Check category attributes
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM category_attributes WHERE category_id = 1001;"
```

---

## 🗄️ Database Connections

### Monolith (Source)

```bash
# Connection string
postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable

# Quick connect
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable"

# Tables
unified_attributes
unified_category_attributes
unified_categories
```

### Microservice (Destination)

```bash
# Connection string
postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable

# Quick connect
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"

# Tables
attributes
category_attributes
categories
```

---

## 📂 File Locations

### Migration Tools

```
/p/github.com/sveturs/listings/
├── cmd/
│   ├── migrate_attributes/
│   │   ├── main.go
│   │   └── README.md
│   └── migrate_category_attributes/
│       ├── main.go
│       └── README.md
├── scripts/
│   ├── validate_attributes_migration.sh
│   └── validate_category_attributes_migration.sh
└── docs/
    ├── ATTRIBUTES_MIGRATION.md
    ├── ATTRIBUTES_MIGRATION_SUMMARY.md
    ├── CATEGORY_ATTRIBUTES_MIGRATION.md
    └── CATEGORY_ATTRIBUTES_MIGRATION_SUMMARY.md
```

---

## 🛠️ Common Tasks

### Rollback Migration

```sql
-- Rollback category_attributes
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "TRUNCATE TABLE category_attributes;"

-- Rollback attributes
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "TRUNCATE TABLE attributes CASCADE;"
```

### Re-run Migration

```bash
# Clear table
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "TRUNCATE TABLE category_attributes;"

# Run migration again
go run ./cmd/migrate_category_attributes/main.go
```

### Check Foreign Keys

```sql
-- Check category_id references
SELECT COUNT(*)
FROM category_attributes ca
LEFT JOIN categories c ON ca.category_id = c.id
WHERE c.id IS NULL;

-- Check attribute_id references
SELECT COUNT(*)
FROM category_attributes ca
LEFT JOIN attributes a ON ca.attribute_id = a.id
WHERE a.id IS NULL;
```

---

## 📊 Statistics

### Attributes Migration
- **Total Records:** 157
- **Success Rate:** 100%
- **Duration:** ~120ms
- **Date:** 2025-11-17

### Category Attributes Migration
- **Total Records:** 479
- **Success Rate:** 100%
- **Duration:** 83ms
- **Date:** 2025-11-17

---

## 🐛 Troubleshooting

### Port Connection Issues

```bash
# Check if PostgreSQL is running
docker ps | grep postgres

# Restart if needed
docker restart listings-postgres
```

### Password Authentication Failed

```bash
# Verify credentials in .env
cat /p/github.com/sveturs/listings/.env | grep DB_

# Correct credentials:
# User: listings_user
# Password: listings_secret
# Port: 35434
# DB: listings_dev_db
```

### Foreign Key Violations

```bash
# Ensure dependencies are migrated first
# 1. Categories must exist
# 2. Attributes must exist
# 3. Then category_attributes can be migrated

# Check if categories exist
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM categories;"
```

---

## 📚 Documentation Links

- [Attributes Migration Guide](./docs/ATTRIBUTES_MIGRATION.md)
- [Attributes Migration Summary](./docs/ATTRIBUTES_MIGRATION_SUMMARY.md)
- [Category Attributes Migration Guide](./docs/CATEGORY_ATTRIBUTES_MIGRATION.md)
- [Category Attributes Migration Summary](./docs/CATEGORY_ATTRIBUTES_MIGRATION_SUMMARY.md)

---

## 🔐 Environment Variables

```bash
# Monolith (rarely needed, defaults are correct)
export SOURCE_DB="postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable"

# Microservice (rarely needed, defaults are correct)
export DEST_DB="postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
```

---

## ⚡ Performance Tips

1. **Batch Size**: Default 100 is optimal for most cases
2. **Dry-run First**: Always test with `--dry-run` before real migration
3. **Validate After**: Always run validation scripts after migration
4. **Local Network**: Migrations are fast on localhost
5. **Transaction Safety**: All migrations use single transaction (atomic)

---

## 📝 Notes

- All migrations are **idempotent** (can be run multiple times safely)
- Use **UPSERT** strategy: `ON CONFLICT DO UPDATE`
- Always **validate** after migration
- **Rollback** is simple: `TRUNCATE TABLE`
- Foreign key **validation** prevents orphaned records

---

**Last Updated:** 2025-11-17
**Status:** 2/6 migrations completed (33%)
