#!/bin/bash
# ============================================================================
# Script: 002_import_to_listings.sh
# Description: Import attributes data to listings microservice PostgreSQL
# Source: /tmp/attribute_migration/*.csv
# Target: postgres://listings_user:listings_secret@localhost:35434/listings_dev_db
# ============================================================================

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
LISTINGS_DB="postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"
INPUT_DIR="/tmp/attribute_migration"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo -e "${GREEN}============================================================================${NC}"
echo -e "${GREEN}Attributes Migration: Import to Listings Microservice${NC}"
echo -e "${GREEN}============================================================================${NC}"
echo ""
echo -e "${YELLOW}Source Dir:${NC} $INPUT_DIR"
echo -e "${YELLOW}Target DB:${NC} $LISTINGS_DB"
echo ""

# Check if export files exist
if [ ! -f "$INPUT_DIR/attributes.csv" ]; then
    echo -e "${RED}✗${NC} Error: Export files not found in $INPUT_DIR"
    echo -e "${YELLOW}→${NC} Run ./001_export_monolith_attributes.sh first"
    exit 1
fi

echo -e "${GREEN}✓${NC} Export files found"

# Check database connection
if ! psql "$LISTINGS_DB" -c "SELECT 1" > /dev/null 2>&1; then
    echo -e "${RED}✗${NC} Error: Cannot connect to listings database"
    echo -e "${YELLOW}→${NC} Check if database is running and credentials are correct"
    exit 1
fi

echo -e "${GREEN}✓${NC} Database connection OK"

# ============================================================================
# IMPORT 1: attributes table
# ============================================================================
echo ""
echo -e "${YELLOW}[1/4] Importing attributes...${NC}"

psql "$LISTINGS_DB" <<EOF
-- Disable triggers temporarily for faster import
SET session_replication_role = replica;

-- Import attributes
COPY attributes (
    id,
    code,
    name,
    display_name,
    attribute_type,
    purpose,
    options,
    validation_rules,
    ui_settings,
    is_searchable,
    is_filterable,
    is_required,
    is_variant_compatible,
    affects_stock,
    affects_price,
    show_in_card,
    is_active,
    sort_order,
    legacy_category_attribute_id,
    legacy_product_variant_attribute_id,
    icon,
    created_at,
    updated_at
)
FROM STDIN WITH CSV HEADER DELIMITER '|' NULL '\\N';
$(cat "$INPUT_DIR/attributes.csv")
\.

-- Re-enable triggers
SET session_replication_role = DEFAULT;

-- Update sequence
SELECT setval('attributes_id_seq', COALESCE((SELECT MAX(id) FROM attributes), 1), true);

-- Show count
SELECT COUNT(*) as imported_attributes FROM attributes;
EOF

echo -e "${GREEN}✓${NC} Attributes imported successfully"

# ============================================================================
# IMPORT 2: category_attributes table
# ============================================================================
echo ""
echo -e "${YELLOW}[2/4] Importing category_attributes...${NC}"

# First, check if categories exist in listings DB
CATEGORY_COUNT=$(psql "$LISTINGS_DB" -t -c "SELECT COUNT(*) FROM categories" | tr -d ' ')

if [ "$CATEGORY_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}⚠ WARNING: No categories found in listings database${NC}"
    echo -e "${YELLOW}→${NC} Category attributes will be imported, but foreign keys may fail"
    echo -e "${YELLOW}→${NC} You may need to sync categories first"
fi

psql "$LISTINGS_DB" <<EOF
-- Disable triggers
SET session_replication_role = replica;

-- Import category_attributes
COPY category_attributes (
    id,
    category_id,
    attribute_id,
    is_enabled,
    is_required,
    sort_order,
    category_specific_options,
    created_at,
    updated_at
)
FROM STDIN WITH CSV HEADER DELIMITER '|' NULL '\\N';
$(cat "$INPUT_DIR/category_attributes.csv")
\.

-- Re-enable triggers
SET session_replication_role = DEFAULT;

-- Update sequence
SELECT setval('category_attributes_id_seq', COALESCE((SELECT MAX(id) FROM category_attributes), 1), true);

-- Show count
SELECT COUNT(*) as imported_category_attributes FROM category_attributes;
EOF

echo -e "${GREEN}✓${NC} Category attributes imported successfully"

# ============================================================================
# IMPORT 3: listing_attribute_values table
# ============================================================================
echo ""
echo -e "${YELLOW}[3/4] Importing listing_attribute_values...${NC}"

# Check if listings exist
LISTING_COUNT=$(psql "$LISTINGS_DB" -t -c "SELECT COUNT(*) FROM listings" | tr -d ' ')

if [ "$LISTING_COUNT" -eq 0 ]; then
    echo -e "${YELLOW}⚠ WARNING: No listings found in listings database${NC}"
    echo -e "${YELLOW}→${NC} Listing attribute values will be imported, but foreign keys may fail"
fi

psql "$LISTINGS_DB" <<EOF
-- Disable triggers
SET session_replication_role = replica;

-- Import listing_attribute_values (map entity_id to listing_id)
CREATE TEMP TABLE temp_listing_attr_values (
    id INTEGER,
    entity_type VARCHAR(50),
    entity_id INTEGER,
    attribute_id INTEGER,
    text_value TEXT,
    numeric_value DECIMAL(20, 4),
    boolean_value BOOLEAN,
    date_value DATE,
    json_value JSONB,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

COPY temp_listing_attr_values
FROM STDIN WITH CSV HEADER DELIMITER '|' NULL '\\N';
$(cat "$INPUT_DIR/attribute_values.csv")
\.

-- Insert into listing_attribute_values (rename entity_id to listing_id)
INSERT INTO listing_attribute_values (
    id,
    listing_id,
    attribute_id,
    value_text,
    value_number,
    value_boolean,
    value_date,
    value_json,
    created_at,
    updated_at
)
SELECT
    id,
    entity_id as listing_id,
    attribute_id,
    text_value,
    numeric_value,
    boolean_value,
    date_value,
    json_value,
    created_at,
    updated_at
FROM temp_listing_attr_values
WHERE entity_type = 'listing';

DROP TABLE temp_listing_attr_values;

-- Re-enable triggers
SET session_replication_role = DEFAULT;

-- Update sequence
SELECT setval('listing_attribute_values_id_seq', COALESCE((SELECT MAX(id) FROM listing_attribute_values), 1), true);

-- Show count
SELECT COUNT(*) as imported_listing_attr_values FROM listing_attribute_values;
EOF

echo -e "${GREEN}✓${NC} Listing attribute values imported successfully"

# ============================================================================
# IMPORT 4: category_variant_attributes (if exists)
# ============================================================================
echo ""
echo -e "${YELLOW}[4/4] Importing category_variant_attributes...${NC}"

VARIANT_LINE_COUNT=$(wc -l < "$INPUT_DIR/variant_attribute_mappings.csv")

if [ "$VARIANT_LINE_COUNT" -le 1 ]; then
    echo -e "${YELLOW}⊘${NC} No variant attribute mappings to import, skipping"
else
    psql "$LISTINGS_DB" <<EOF
-- Disable triggers
SET session_replication_role = replica;

-- Import category_variant_attributes
CREATE TEMP TABLE temp_variant_mappings (
    id INTEGER,
    category_id INTEGER,
    variant_attribute_id INTEGER,
    is_required BOOLEAN,
    affects_price BOOLEAN,
    affects_stock BOOLEAN,
    sort_order INTEGER,
    display_as VARCHAR(50),
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

COPY temp_variant_mappings
FROM STDIN WITH CSV HEADER DELIMITER '|' NULL '\\N';
$(cat "$INPUT_DIR/variant_attribute_mappings.csv")
\.

-- Insert into category_variant_attributes (rename variant_attribute_id to attribute_id)
INSERT INTO category_variant_attributes (
    id,
    category_id,
    attribute_id,
    is_required,
    affects_price,
    affects_stock,
    sort_order,
    display_as,
    created_at,
    updated_at
)
SELECT
    id,
    category_id,
    variant_attribute_id,
    is_required,
    affects_price,
    affects_stock,
    sort_order,
    display_as,
    created_at,
    updated_at
FROM temp_variant_mappings;

DROP TABLE temp_variant_mappings;

-- Re-enable triggers
SET session_replication_role = DEFAULT;

-- Update sequence
SELECT setval('category_variant_attributes_id_seq', COALESCE((SELECT MAX(id) FROM category_variant_attributes), 1), true);

-- Show count
SELECT COUNT(*) as imported_variant_attrs FROM category_variant_attributes;
EOF

    echo -e "${GREEN}✓${NC} Category variant attributes imported successfully"
fi

# ============================================================================
# VERIFY FOREIGN KEY INTEGRITY
# ============================================================================
echo ""
echo -e "${YELLOW}Verifying foreign key integrity...${NC}"

psql "$LISTINGS_DB" <<EOF
-- Check for orphaned category_attributes (category_id not in categories)
SELECT
    COUNT(*) as orphaned_category_attrs,
    'category_attributes with invalid category_id' as description
FROM category_attributes ca
LEFT JOIN categories c ON ca.category_id = c.id
WHERE c.id IS NULL;

-- Check for orphaned listing_attribute_values (listing_id not in listings)
SELECT
    COUNT(*) as orphaned_listing_attrs,
    'listing_attribute_values with invalid listing_id' as description
FROM listing_attribute_values lav
LEFT JOIN listings l ON lav.listing_id = l.id
WHERE l.id IS NULL;
EOF

# ============================================================================
# SUMMARY
# ============================================================================
echo ""
echo -e "${GREEN}============================================================================${NC}"
echo -e "${GREEN}Import Complete!${NC}"
echo -e "${GREEN}============================================================================${NC}"
echo ""

# Get final counts
psql "$LISTINGS_DB" -t <<EOF
\echo 'Final Counts:'
\echo '  Attributes:               ' || (SELECT COUNT(*) FROM attributes)
\echo '  Category Attributes:      ' || (SELECT COUNT(*) FROM category_attributes)
\echo '  Listing Attribute Values: ' || (SELECT COUNT(*) FROM listing_attribute_values)
\echo '  Variant Attributes:       ' || (SELECT COUNT(*) FROM category_variant_attributes)
\echo ''
EOF

echo -e "${GREEN}✓${NC} Next step: Run ./003_validate_migration.sql to verify data integrity"
echo ""
