# Orders Data Migration Tool

Migration script for transferring inventory reservations from monolith database to Orders microservice.

## 🚀 Quick Start

```bash
# Dry-run (safe, no writes)
go run main.go --verbose

# Execute migration (writes data)
go run main.go --dry-run=false --verbose

# Build standalone binary
go build -o migrate_orders main.go
./migrate_orders --dry-run=false
```

## 📚 Full Documentation

See comprehensive guide: [/p/github.com/sveturs/listings/docs/PHASE_17_DATA_MIGRATION_GUIDE.md](../../docs/PHASE_17_DATA_MIGRATION_GUIDE.md)

## ⚡ Features

- ✅ Dry-run mode by default
- ✅ Idempotent (safe to run multiple times)
- ✅ Transaction-based (atomic)
- ✅ FK constraint validation
- ✅ Order ID nullification if order doesn't exist
- ✅ Detailed logging

## 🔧 Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--monolith-dsn` | postgres://postgres:...@localhost:5433/svetubd | Monolith DB connection string |
| `--microservice-dsn` | postgres://listings_user:...@localhost:35434/listings_dev_db | Microservice DB connection string |
| `--dry-run` | `true` | Dry-run mode (no writes) |
| `--verbose` | `false` | Verbose logging |

## ✅ Verification

```bash
# Check migrated count
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT COUNT(*) FROM inventory_reservations;"

# View migrated records
psql "postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable" \
  -c "SELECT id, listing_id, quantity, status, expires_at FROM inventory_reservations;"
```

## 📝 What Gets Migrated

- **Inventory Reservations**: Only active (`status = 'active'`) and non-expired (`expires_at > NOW()`)
- **Shopping Carts**: N/A (don't exist in monolith)
- **Orders**: N/A (don't exist in monolith)
- **Order Items**: N/A (don't exist in monolith)

## 🔍 Schema Differences

| Monolith | Microservice | Handling |
|----------|--------------|----------|
| `product_id` | `listing_id` | Direct 1:1 mapping |
| enum `reservation_status` | varchar(20) `status` | Cast to text |
| `order_id` NOT NULL | `order_id` NULL | Set NULL if order not found |

## 🐛 Troubleshooting

See [Troubleshooting section](../../docs/PHASE_17_DATA_MIGRATION_GUIDE.md#-troubleshooting) in full documentation.

---

**Phase:** 17 Days 18-19
**Status:** ✅ Production Ready
**Last Updated:** 2025-11-14
