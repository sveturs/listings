#!/bin/bash
# ============================================================================
# Apply All Category Translations
# Description: Applies all translation scripts in correct order
# Author: System
# Date: 2025-11-10
# ============================================================================

set -e  # Exit on error

# Database connection string
DB_CONN="postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"

# Script directory
SCRIPT_DIR="/p/github.com/sveturs/listings/scripts"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}============================================================================${NC}"
echo -e "${BLUE}Category Translations Migration${NC}"
echo -e "${BLUE}============================================================================${NC}"
echo ""

# Step 1: Add translation columns
echo -e "${YELLOW}Step 1/3: Adding translation columns...${NC}"
psql "$DB_CONN" -f "$SCRIPT_DIR/01_add_translation_columns.sql" > /dev/null
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Translation columns added successfully${NC}"
else
    echo -e "${RED}✗ Failed to add translation columns${NC}"
    exit 1
fi
echo ""

# Step 2: Create backup
echo -e "${YELLOW}Step 2/3: Creating backup...${NC}"
psql "$DB_CONN" -f "$SCRIPT_DIR/02_backup_categories_before_translation.sql" > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Backup created successfully${NC}"
else
    echo -e "${RED}✗ Failed to create backup${NC}"
    exit 1
fi
echo ""

# Step 3: Apply translations
echo -e "${YELLOW}Step 3/3: Applying translations...${NC}"
psql "$DB_CONN" -f "$SCRIPT_DIR/03_add_category_translations.sql" | tail -5
if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Translations applied successfully${NC}"
else
    echo -e "${RED}✗ Failed to apply translations${NC}"
    exit 1
fi
echo ""

# Verification
echo -e "${BLUE}============================================================================${NC}"
echo -e "${BLUE}Verification Results${NC}"
echo -e "${BLUE}============================================================================${NC}"
echo ""

echo -e "${YELLOW}Translation Coverage:${NC}"
psql "$DB_CONN" -c "
SELECT
  COUNT(*) as total_categories,
  COUNT(title_en) as with_english,
  COUNT(title_ru) as with_russian,
  COUNT(title_sr) as with_serbian,
  COUNT(CASE WHEN title_en IS NULL OR title_ru IS NULL OR title_sr IS NULL THEN 1 END) as missing
FROM c2c_categories;
"

echo ""
echo -e "${YELLOW}Sample Root Categories:${NC}"
psql "$DB_CONN" -c "
SELECT id, name, title_en, title_ru, title_sr
FROM c2c_categories
WHERE parent_id IS NULL
ORDER BY id
LIMIT 5;
"

echo ""
echo -e "${GREEN}============================================================================${NC}"
echo -e "${GREEN}Migration completed successfully!${NC}"
echo -e "${GREEN}============================================================================${NC}"
echo ""
echo -e "${YELLOW}Backup table created: ${NC}c2c_categories_backup_20251110"
echo -e "${YELLOW}To rollback: ${NC}psql \"$DB_CONN\" -f $SCRIPT_DIR/04_rollback_category_translations.sql"
echo ""
