#!/bin/bash
# Script to migrate test fixtures from old schema (b2c_products, b2c_storefronts)
# to unified schema (listings, storefronts)

set -e

FIXTURES_DIR="/p/github.com/sveturs/listings/tests/fixtures"

echo "=== Migrating fixtures to unified schema ==="

# Backup original files first
echo "Creating backups..."
for file in "$FIXTURES_DIR"/*.sql; do
    if [[ -f "$file" && ! "$file" =~ _backup ]]; then
        cp "$file" "${file}.backup_$(date +%Y%m%d_%H%M%S)"
    fi
done

# List of files to migrate
FILES=(
    "b2c_inventory_fixtures.sql"
    "bulk_operations_fixtures.sql"
    "create_product_fixtures.sql"
    "decrement_stock_fixtures.sql"
    "rollback_stock_fixtures.sql"
    "update_product_fixtures.sql"
)

for file in "${FILES[@]}"; do
    filepath="$FIXTURES_DIR/$file"

    if [[ ! -f "$filepath" ]]; then
        echo "⚠️  File not found: $file - skipping"
        continue
    fi

    echo "Processing: $file"

    # Apply transformations using sed
    sed -i \
        -e 's/b2c_products/listings/g' \
        -e 's/b2c_storefronts/storefronts/g' \
        "$filepath"

    echo "✅ Migrated: $file"
done

echo ""
echo "=== Migration complete ==="
echo "Backups created with timestamp suffix"
echo "Next steps:"
echo "1. Manually review and update field mappings:"
echo "   - name → title"
echo "   - stock_quantity → quantity"
echo "   - stock_status → status"
echo "   - is_active → status (active/inactive)"
echo "   - Add source_type = 'b2c' for all listings"
echo "2. Remove incompatible fields (barcode, attributes for listings)"
echo "3. Re-run tests"
