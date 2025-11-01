# Foreign Keys Integration Tests

## Overview

Этот набор тестов проверяет корректность работы Foreign Key constraints после применения миграции `000194_add_foreign_keys_c2c_b2c.sql`.

## Test Files

1. **test_foreign_keys_cascade.sql** - тесты CASCADE DELETE поведения
2. **test_foreign_keys_restrict.sql** - тесты RESTRICT поведения
3. **run_fk_tests.sh** - bash скрипт для запуска всех тестов
4. **foreign_keys_test.go** - Go интеграционные тесты (в `internal/storage/postgres/`)

## Important Notes

### Database Schema Considerations

**⚠️ ВАЖНО:** Эти тесты требуют определённой схемы базы данных:

1. **Auth Service Integration**: Проект использует внешний Auth Service, поэтому:
   - Таблицы `users` НЕТ в локальной БД
   - Тесты которые требуют `users` будут **SKIP** или **FAIL** в текущей БД
   - Это **ОЖИДАЕМОЕ ПОВЕДЕНИЕ**

2. **Naming Conventions**: Некоторые таблицы имеют другие названия:
   - `storefronts` → `b2c_stores`
   - `users` → отсутствует (Auth Service)

3. **Migration Status**: Тесты ожидают что миграция `000194` уже применена:
   - Если миграция не применена → тесты покажут 0 FK constraints
   - Это **НЕ ОШИБКА**, нужно сначала применить миграцию

## Running Tests

### Prerequisites

1. **Apply FK Migration First:**
   ```bash
   cd /p/github.com/sveturs/svetu/backend
   ./migrator up
   ```

2. **Verify Migration Applied:**
   ```bash
   psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" -c "
   SELECT COUNT(*) as fk_count
   FROM information_schema.table_constraints
   WHERE constraint_type = 'FOREIGN KEY'
   AND table_schema = 'public';
   "
   ```

   Expected: `fk_count > 0`

### Option 1: Run All Tests (Recommended)

```bash
cd /p/github.com/sveturs/svetu/backend/migrations/tests
chmod +x run_fk_tests.sh
./run_fk_tests.sh
```

### Option 2: Run Individual SQL Tests

```bash
# CASCADE DELETE tests
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" \
  -f test_foreign_keys_cascade.sql

# RESTRICT tests
psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable" \
  -f test_foreign_keys_restrict.sql
```

### Option 3: Run Go Tests

```bash
cd /p/github.com/sveturs/svetu/backend
go test -v ./internal/storage/postgres -run TestForeignKeyConstraints
```

## Expected Results

### If Migration NOT Applied

```
⚠ WARNING: No FK constraints found - migration may not be applied yet
Total FK constraints: 0
```

**Solution:** Apply migration first using `./migrator up`

### If Migration IS Applied

```
✅ All CASCADE DELETE tests pass
✅ All RESTRICT tests pass
📊 FK Constraints Summary:
   Total FK constraints: 17+
   CASCADE DELETE: 9+
   RESTRICT/NO ACTION: 7+
```

## Test Coverage

### CASCADE DELETE Tests (7 test cases):
1. `c2c_images.listing_id` → listing deletion cascades to images
2. `c2c_attributes.listing_id` → listing deletion cascades to attributes
3. `c2c_favorites.listing_id` → listing deletion cascades to favorites
4. `b2c_product_images.product_id` → product deletion cascades to images
5. `b2c_product_variants.product_id` → product deletion cascades to variants
6. Multi-layer CASCADE (listing + images + attributes + favorites)
7. User deletion CASCADE (if users table exists)

### RESTRICT Tests (7 test cases):
1. Cannot delete category with existing listings
2. Cannot delete user with existing storefronts (SKIP if no users table)
3. Cannot delete attribute_meta with existing attribute values
4. Cannot delete B2C category with existing products
5. Cannot delete storefront with existing products
6. RESTRICT vs CASCADE comparison
7. FK metadata validation

### Go Integration Tests (9 test cases):
- All CASCADE DELETE scenarios
- All RESTRICT scenarios
- FK metadata validation
- Performance benchmarks

## Troubleshooting

### Issue: "relation users does not exist"

**Reason:** Auth Service architecture - users managed externally

**Solution:** Tests will SKIP or FAIL gracefully. This is expected.

### Issue: "No FK constraints found"

**Reason:** Migration `000194` not applied yet

**Solution:**
```bash
cd /p/github.com/sveturs/svetu/backend
./migrator up
```

### Issue: "Cannot connect to database"

**Reason:** Wrong port or credentials

**Solution:**
- Check port: Should be `5433` (not 5432)
- Verify connection: `psql "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable"`

### Issue: Tests timeout

**Reason:** Large dataset or slow queries

**Solution:** Increase timeout in `run_fk_tests.sh` or run tests individually

## Maintenance

### Adding New Tests

When adding new FK constraints:

1. Add CASCADE test in `test_foreign_keys_cascade.sql`
2. Add RESTRICT test in `test_foreign_keys_restrict.sql` (if applicable)
3. Add Go test in `foreign_keys_test.go`
4. Update coverage table in phase-1-p0.md
5. Run all tests to verify

### Updating Tests

When schema changes:
1. Update table/column names in SQL tests
2. Update Go test expectations
3. Re-run full test suite
4. Update documentation

## Performance

- **SQL tests runtime:** ~5-10 seconds (with transactions)
- **Go tests runtime:** ~3-5 seconds
- **Total test suite:** <20 seconds

All tests use `BEGIN/ROLLBACK` transactions to avoid data pollution.

## References

- **Migration file:** `backend/migrations/000194_add_foreign_keys_c2c_b2c.up.sql`
- **Documentation:** `docs/migration/phases/phase-1-p0.md`
- **Auth Service:** External microservice (no local users table)

## Support

For issues or questions:
1. Check CLAUDE.md troubleshooting section
2. Verify migration applied correctly
3. Check database connection (port 5433)
4. Review test output for specific error messages
