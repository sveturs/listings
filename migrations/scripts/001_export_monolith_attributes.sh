#!/bin/bash
# ============================================================================
# Script: 001_export_monolith_attributes.sh
# Description: Export attributes data from monolith PostgreSQL database
# Source: postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/vondi_db
# Output: /tmp/attribute_migration/
# ============================================================================

set -e  # Exit on error

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Configuration
MONOLITH_DB="postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/vondi_db?sslmode=disable"
OUTPUT_DIR="/tmp/attribute_migration"
TIMESTAMP=$(date +%Y%m%d_%H%M%S)

echo -e "${GREEN}============================================================================${NC}"
echo -e "${GREEN}Attributes Migration: Export from Monolith${NC}"
echo -e "${GREEN}============================================================================${NC}"
echo ""
echo -e "${YELLOW}Source DB:${NC} $MONOLITH_DB"
echo -e "${YELLOW}Output Dir:${NC} $OUTPUT_DIR"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"
echo -e "${GREEN}✓${NC} Created output directory: $OUTPUT_DIR"

# ============================================================================
# EXPORT 1: unified_attributes (203 records expected)
# ============================================================================
echo ""
echo -e "${YELLOW}[1/3] Exporting unified_attributes...${NC}"

psql "$MONOLITH_DB" -c "COPY (
    SELECT
        id,
        code,
        -- Convert name to JSONB format (currently VARCHAR)
        jsonb_build_object('en', name, 'ru', name, 'sr', name) as name,
        -- Convert display_name to JSONB format
        jsonb_build_object('en', display_name, 'ru', display_name, 'sr', display_name) as display_name,
        attribute_type,
        purpose,
        options::text as options,
        validation_rules::text as validation_rules,
        ui_settings::text as ui_settings,
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
    FROM unified_attributes
    ORDER BY id
) TO STDOUT WITH CSV HEADER DELIMITER '|' NULL '\\N'" > "$OUTPUT_DIR/attributes.csv"

ATTR_COUNT=$(wc -l < "$OUTPUT_DIR/attributes.csv")
ATTR_COUNT=$((ATTR_COUNT - 1))  # Subtract header
echo -e "${GREEN}✓${NC} Exported $ATTR_COUNT attributes to attributes.csv"

if [ $ATTR_COUNT -ne 203 ]; then
    echo -e "${RED}⚠ WARNING: Expected 203 attributes, but got $ATTR_COUNT${NC}"
fi

# ============================================================================
# EXPORT 2: unified_category_attributes (category relationships)
# ============================================================================
echo ""
echo -e "${YELLOW}[2/3] Exporting unified_category_attributes...${NC}"

psql "$MONOLITH_DB" -c "COPY (
    SELECT
        id,
        category_id,
        attribute_id,
        is_enabled,
        is_required,
        sort_order,
        category_specific_options::text as category_specific_options,
        created_at,
        updated_at
    FROM unified_category_attributes
    ORDER BY category_id, sort_order, attribute_id
) TO STDOUT WITH CSV HEADER DELIMITER '|' NULL '\\N'" > "$OUTPUT_DIR/category_attributes.csv"

CAT_ATTR_COUNT=$(wc -l < "$OUTPUT_DIR/category_attributes.csv")
CAT_ATTR_COUNT=$((CAT_ATTR_COUNT - 1))
echo -e "${GREEN}✓${NC} Exported $CAT_ATTR_COUNT category-attribute relationships to category_attributes.csv"

# ============================================================================
# EXPORT 3: unified_attribute_values (listing values)
# ============================================================================
echo ""
echo -e "${YELLOW}[3/3] Exporting unified_attribute_values...${NC}"

psql "$MONOLITH_DB" -c "COPY (
    SELECT
        id,
        entity_type,
        entity_id,
        attribute_id,
        text_value,
        numeric_value,
        boolean_value,
        date_value,
        json_value::text as json_value,
        created_at,
        updated_at
    FROM unified_attribute_values
    WHERE entity_type = 'listing'
    ORDER BY entity_id, attribute_id
) TO STDOUT WITH CSV HEADER DELIMITER '|' NULL '\\N'" > "$OUTPUT_DIR/attribute_values.csv"

ATTR_VAL_COUNT=$(wc -l < "$OUTPUT_DIR/attribute_values.csv")
ATTR_VAL_COUNT=$((ATTR_VAL_COUNT - 1))
echo -e "${GREEN}✓${NC} Exported $ATTR_VAL_COUNT listing attribute values to attribute_values.csv"

# ============================================================================
# EXPORT 4: variant_attribute_mappings (if exists)
# ============================================================================
echo ""
echo -e "${YELLOW}[OPTIONAL] Checking for variant_attribute_mappings...${NC}"

if psql "$MONOLITH_DB" -t -c "SELECT EXISTS (SELECT 1 FROM information_schema.tables WHERE table_name='variant_attribute_mappings')" | grep -q t; then
    psql "$MONOLITH_DB" -c "COPY (
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
        FROM variant_attribute_mappings
        ORDER BY category_id, sort_order
    ) TO STDOUT WITH CSV HEADER DELIMITER '|' NULL '\\N'" > "$OUTPUT_DIR/variant_attribute_mappings.csv"

    VARIANT_COUNT=$(wc -l < "$OUTPUT_DIR/variant_attribute_mappings.csv")
    VARIANT_COUNT=$((VARIANT_COUNT - 1))
    echo -e "${GREEN}✓${NC} Exported $VARIANT_COUNT variant attribute mappings to variant_attribute_mappings.csv"
else
    echo -e "${YELLOW}⊘${NC} Table variant_attribute_mappings does not exist, skipping"
    touch "$OUTPUT_DIR/variant_attribute_mappings.csv"
    echo "id|category_id|variant_attribute_id|is_required|affects_price|affects_stock|sort_order|display_as|created_at|updated_at" > "$OUTPUT_DIR/variant_attribute_mappings.csv"
fi

# ============================================================================
# SUMMARY
# ============================================================================
echo ""
echo -e "${GREEN}============================================================================${NC}"
echo -e "${GREEN}Export Complete!${NC}"
echo -e "${GREEN}============================================================================${NC}"
echo ""
echo -e "${YELLOW}Summary:${NC}"
echo -e "  Attributes:                 $ATTR_COUNT"
echo -e "  Category Relationships:     $CAT_ATTR_COUNT"
echo -e "  Listing Attribute Values:   $ATTR_VAL_COUNT"
echo -e "  Variant Mappings:           ${VARIANT_COUNT:-0}"
echo ""
echo -e "${YELLOW}Files created in $OUTPUT_DIR:${NC}"
ls -lh "$OUTPUT_DIR"/*.csv
echo ""
echo -e "${GREEN}✓${NC} Next step: Run ./002_import_to_listings.sh"
echo ""
