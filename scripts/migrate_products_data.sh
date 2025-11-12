#!/bin/bash
set -e  # Exit on error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

BACKUP_FILE="/tmp/backup_before_drop_legacy_20251103_215749.sql"
MICROSERVICE_DB="postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
TEMP_DIR="/tmp/products_migration_$$"

echo -e "${YELLOW}Starting Products data migration to microservice...${NC}"

# Check if backup file exists
if [ ! -f "$BACKUP_FILE" ]; then
    echo -e "${RED}Error: Backup file not found: $BACKUP_FILE${NC}"
    exit 1
fi

# Test DB connectivity
echo -e "${YELLOW}Testing database connectivity...${NC}"
if ! psql "$MICROSERVICE_DB" -c "SELECT 1;" > /dev/null 2>&1; then
    echo -e "${RED}Error: Cannot connect to microservice database${NC}"
    exit 1
fi
echo -e "${GREEN}✓ Database connection OK${NC}"

# Create temp directory
mkdir -p "$TEMP_DIR"

# Step 1: Extract CREATE TABLE statements
echo -e "${YELLOW}Step 1: Extracting table schemas...${NC}"

# Extract b2c_products table
awk '/CREATE TABLE public\.b2c_products/,/;/' "$BACKUP_FILE" > "$TEMP_DIR/create_products.sql"
if [ ! -s "$TEMP_DIR/create_products.sql" ]; then
    echo -e "${RED}Error: Could not extract b2c_products schema${NC}"
    exit 1
fi

# Extract b2c_product_variants table
awk '/CREATE TABLE public\.b2c_product_variants/,/;/' "$BACKUP_FILE" > "$TEMP_DIR/create_variants.sql"
if [ ! -s "$TEMP_DIR/create_variants.sql" ]; then
    echo -e "${RED}Error: Could not extract b2c_product_variants schema${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Schemas extracted${NC}"

# Step 2: Extract INSERT statements
echo -e "${YELLOW}Step 2: Extracting data...${NC}"

grep "^INSERT INTO public\.b2c_products" "$BACKUP_FILE" > "$TEMP_DIR/insert_products.sql" || echo "No products data found"
grep "^INSERT INTO public\.b2c_product_variants" "$BACKUP_FILE" > "$TEMP_DIR/insert_variants.sql" || echo "No variants data found"

PRODUCTS_INSERTS=$([ -f "$TEMP_DIR/insert_products.sql" ] && wc -l < "$TEMP_DIR/insert_products.sql" || echo "0")
VARIANTS_INSERTS=$([ -f "$TEMP_DIR/insert_variants.sql" ] && wc -l < "$TEMP_DIR/insert_variants.sql" || echo "0")

echo -e "${GREEN}✓ Found $PRODUCTS_INSERTS products INSERT statements${NC}"
echo -e "${GREEN}✓ Found $VARIANTS_INSERTS variants INSERT statements${NC}"

# Step 3: Drop existing tables if they exist
echo -e "${YELLOW}Step 3: Dropping existing tables if they exist...${NC}"
psql "$MICROSERVICE_DB" -c "DROP TABLE IF EXISTS b2c_product_variants CASCADE;" > /dev/null
psql "$MICROSERVICE_DB" -c "DROP TABLE IF EXISTS b2c_products CASCADE;" > /dev/null
echo -e "${GREEN}✓ Tables dropped${NC}"

# Step 4: Create required ENUM types and sequences
echo -e "${YELLOW}Step 4: Creating required ENUM types and sequences...${NC}"
psql "$MICROSERVICE_DB" -c "DROP TYPE IF EXISTS location_privacy_level CASCADE;" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE TYPE location_privacy_level AS ENUM ('exact', 'street', 'district', 'city');" > /dev/null

# Create sequence for product IDs
psql "$MICROSERVICE_DB" -c "DROP SEQUENCE IF EXISTS global_product_id_seq CASCADE;" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE SEQUENCE global_product_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;" > /dev/null

# Create sequence for variant IDs
psql "$MICROSERVICE_DB" -c "DROP SEQUENCE IF EXISTS b2c_product_variants_id_seq CASCADE;" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE SEQUENCE b2c_product_variants_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1;" > /dev/null

echo -e "${GREEN}✓ ENUM types and sequences created${NC}"

# Step 5: Create tables
echo -e "${YELLOW}Step 5: Creating tables...${NC}"
# Remove public schema prefix from CREATE TABLE statements
sed -e 's/public\.location_privacy_level/location_privacy_level/g' \
    -e 's/public\.global_product_id_seq/global_product_id_seq/g' \
    -e 's/public\.b2c_product_variants_id_seq/b2c_product_variants_id_seq/g' \
    "$TEMP_DIR/create_products.sql" > "$TEMP_DIR/create_products_clean.sql"

sed -e 's/public\.location_privacy_level/location_privacy_level/g' \
    -e 's/public\.global_product_id_seq/global_product_id_seq/g' \
    -e 's/public\.b2c_product_variants_id_seq/b2c_product_variants_id_seq/g' \
    "$TEMP_DIR/create_variants.sql" > "$TEMP_DIR/create_variants_clean.sql"

psql "$MICROSERVICE_DB" < "$TEMP_DIR/create_products_clean.sql" > /dev/null
psql "$MICROSERVICE_DB" < "$TEMP_DIR/create_variants_clean.sql" > /dev/null
echo -e "${GREEN}✓ Tables created${NC}"

# Step 5.1: Add PRIMARY KEYs
echo -e "${YELLOW}Step 5.1: Adding PRIMARY KEYs...${NC}"
psql "$MICROSERVICE_DB" -c "ALTER TABLE b2c_products ADD CONSTRAINT b2c_products_pkey PRIMARY KEY (id);" > /dev/null
psql "$MICROSERVICE_DB" -c "ALTER TABLE b2c_product_variants ADD CONSTRAINT b2c_product_variants_pkey PRIMARY KEY (id);" > /dev/null
echo -e "${GREEN}✓ PRIMARY KEYs added${NC}"

# Step 5.2: Add UNIQUE constraints
echo -e "${YELLOW}Step 5.2: Adding UNIQUE constraints...${NC}"
psql "$MICROSERVICE_DB" -c "ALTER TABLE b2c_product_variants ADD CONSTRAINT b2c_product_variants_sku_key UNIQUE (sku);" > /dev/null 2>&1 || echo "  (UNIQUE constraint on sku skipped - may already exist)"
echo -e "${GREEN}✓ UNIQUE constraints added${NC}"

# Step 5.3: Create indexes
echo -e "${YELLOW}Step 5.3: Creating indexes...${NC}"

# Products indexes
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_barcode_idx ON b2c_products USING btree (barcode) WHERE (barcode IS NOT NULL);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_category_id_idx ON b2c_products USING btree (category_id);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_has_individual_location_idx ON b2c_products USING btree (has_individual_location);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_has_variants_idx ON b2c_products USING btree (has_variants);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_is_active_idx ON b2c_products USING btree (is_active);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_location_privacy_idx ON b2c_products USING btree (location_privacy);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_show_on_map_idx ON b2c_products USING btree (show_on_map);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_sku_idx ON b2c_products USING btree (sku) WHERE (sku IS NOT NULL);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_stock_status_idx ON b2c_products USING btree (stock_status);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE UNIQUE INDEX b2c_products_storefront_id_barcode_idx ON b2c_products USING btree (storefront_id, barcode) WHERE (barcode IS NOT NULL);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE UNIQUE INDEX b2c_products_storefront_id_sku_idx ON b2c_products USING btree (storefront_id, sku) WHERE (sku IS NOT NULL);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX b2c_products_storefront_id_view_count_idx ON b2c_products USING btree (storefront_id, view_count);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX idx_b2c_products_active_created ON b2c_products USING btree (is_active, created_at DESC) WHERE (is_active = true);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX idx_b2c_products_category_active ON b2c_products USING btree (category_id, is_active) WHERE (is_active = true);" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX idx_b2c_products_price ON b2c_products USING btree (price) WHERE ((is_active = true) AND (price IS NOT NULL));" > /dev/null
psql "$MICROSERVICE_DB" -c "CREATE INDEX idx_b2c_products_storefront ON b2c_products USING btree (storefront_id, is_active) WHERE (is_active = true);" > /dev/null

echo -e "${GREEN}✓ Indexes created (16 indexes)${NC}"

# Step 6: Insert data
echo -e "${YELLOW}Step 6: Inserting data...${NC}"
if [ -f "$TEMP_DIR/insert_products.sql" ] && [ -s "$TEMP_DIR/insert_products.sql" ]; then
    # Replace public schema references with default schema
    sed 's/public\.b2c_products/b2c_products/g' "$TEMP_DIR/insert_products.sql" > "$TEMP_DIR/insert_products_clean.sql"
    psql "$MICROSERVICE_DB" < "$TEMP_DIR/insert_products_clean.sql" > /dev/null
    echo -e "${GREEN}✓ Products data inserted${NC}"
else
    echo -e "${YELLOW}⚠ No products data to insert${NC}"
fi

if [ -f "$TEMP_DIR/insert_variants.sql" ] && [ -s "$TEMP_DIR/insert_variants.sql" ]; then
    # Replace public schema references with default schema
    sed 's/public\.b2c_product_variants/b2c_product_variants/g' "$TEMP_DIR/insert_variants.sql" > "$TEMP_DIR/insert_variants_clean.sql"
    psql "$MICROSERVICE_DB" < "$TEMP_DIR/insert_variants_clean.sql" > /dev/null
    echo -e "${GREEN}✓ Variants data inserted${NC}"
else
    echo -e "${YELLOW}⚠ No variants data to insert${NC}"
fi

# Step 7: Update sequence value
echo -e "${YELLOW}Step 7: Updating sequence value...${NC}"
MAX_ID=$(psql "$MICROSERVICE_DB" -t -c "SELECT COALESCE(MAX(id), 0) FROM b2c_products;" | tr -d ' ')
if [ "$MAX_ID" -gt 0 ]; then
    psql "$MICROSERVICE_DB" -c "SELECT pg_catalog.setval('global_product_id_seq', $MAX_ID, true);" > /dev/null
    echo -e "${GREEN}✓ Sequence updated to $MAX_ID${NC}"
else
    echo -e "${YELLOW}⚠ No products found, sequence not updated${NC}"
fi

# Step 8: Verify data integrity
echo -e "${YELLOW}Step 8: Verifying data integrity...${NC}"
PRODUCTS_COUNT=$(psql "$MICROSERVICE_DB" -t -c "SELECT COUNT(*) FROM b2c_products;" | tr -d ' ')
VARIANTS_COUNT=$(psql "$MICROSERVICE_DB" -t -c "SELECT COUNT(*) FROM b2c_product_variants;" | tr -d ' ')

echo -e ""
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}Migration completed successfully!${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e "${GREEN}✓ Products migrated: ${PRODUCTS_COUNT}${NC}"
echo -e "${GREEN}✓ Variants migrated: ${VARIANTS_COUNT}${NC}"
echo -e "${GREEN}========================================${NC}"
echo -e ""

# Sample data check
echo -e "${YELLOW}Sample data check:${NC}"
psql "$MICROSERVICE_DB" -c "SELECT id, name, price, created_at FROM b2c_products ORDER BY created_at DESC LIMIT 3;"

# Cleanup
rm -rf "$TEMP_DIR"
echo -e "${GREEN}✓ Temporary files cleaned up${NC}"
