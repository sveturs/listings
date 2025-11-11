#!/bin/bash
#
# Rollback Script for Data Migration
# Removes migrated data from new database
#

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Database credentials
NEW_DB="postgres://listings_user:listings_secret@localhost:35434/listings_dev_db?sslmode=disable"

echo "================================================================================"
echo "                     MIGRATION ROLLBACK SCRIPT"
echo "                     WARNING: This will delete migrated data!"
echo "================================================================================"
echo ""

# Function to run SQL on new DB
run_new_db() {
    psql "$NEW_DB" -c "$1" 2>&1
}

# Safety check - ask for confirmation
echo -e "${RED}⚠️  WARNING: This will DELETE the following data from listings_dev_db:${NC}"
echo "  - All C2C listings (source_type='c2c')"
echo "  - All related listing images"
echo "  - All related listing attributes"
echo "  - All related listing locations"
echo ""
echo -e "${YELLOW}Note: B2C storefronts will NOT be deleted (manual cleanup required if needed)${NC}"
echo ""

# Show current counts
echo "Current data counts:"
echo "-------------------"
c2c_count=$(psql "$NEW_DB" -t -c "SELECT COUNT(*) FROM listings WHERE source_type='c2c'" | xargs)
images_count=$(psql "$NEW_DB" -t -c "SELECT COUNT(*) FROM listing_images li JOIN listings l ON l.id = li.listing_id WHERE l.source_type='c2c'" | xargs)
echo "  - C2C Listings: $c2c_count"
echo "  - Related Images: $images_count"
echo ""

read -p "Are you sure you want to rollback? Type 'YES' to confirm: " confirmation

if [ "$confirmation" != "YES" ]; then
    echo -e "${GREEN}Rollback cancelled.${NC}"
    exit 0
fi

echo ""
echo -e "${YELLOW}Starting rollback...${NC}"
echo ""

# Create backup before rollback
timestamp=$(date +%Y%m%d_%H%M%S)
backup_file="/tmp/listings_dev_db_before_rollback_${timestamp}.sql"

echo "Creating backup before rollback..."
PGPASSWORD=listings_secret pg_dump -h localhost -p 35434 -U listings_user -d listings_dev_db \
    --no-owner --no-acl -f "$backup_file"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓${NC} Backup created: $backup_file"
else
    echo -e "${RED}✗${NC} Backup failed! Aborting rollback."
    exit 1
fi

echo ""
echo "Deleting migrated data..."
echo ""

# Start transaction
psql "$NEW_DB" <<EOF
BEGIN;

-- Delete images for C2C listings
DELETE FROM listing_images
WHERE listing_id IN (SELECT id FROM listings WHERE source_type = 'c2c');

-- Delete attributes for C2C listings (if table exists)
DELETE FROM listing_attributes
WHERE listing_id IN (SELECT id FROM listings WHERE source_type = 'c2c');

-- Delete locations for C2C listings (if table exists)
DELETE FROM listing_locations
WHERE listing_id IN (SELECT id FROM listings WHERE source_type = 'c2c');

-- Delete tags for C2C listings (if table exists)
DELETE FROM listing_tags
WHERE listing_id IN (SELECT id FROM listings WHERE source_type = 'c2c');

-- Delete favorites for C2C listings (if table exists)
DELETE FROM c2c_favorites
WHERE listing_id IN (SELECT id FROM listings WHERE source_type = 'c2c');

-- Delete stats for C2C listings (if table exists)
DELETE FROM listing_stats
WHERE listing_id IN (SELECT id FROM listings WHERE source_type = 'c2c');

-- Finally, delete the C2C listings
DELETE FROM listings WHERE source_type = 'c2c';

COMMIT;
EOF

if [ $? -eq 0 ]; then
    echo ""
    echo -e "${GREEN}✓✓✓ Rollback completed successfully!${NC}"
    echo ""
    echo "Final counts:"
    echo "-------------"
    c2c_count_after=$(psql "$NEW_DB" -t -c "SELECT COUNT(*) FROM listings WHERE source_type='c2c'" | xargs)
    images_count_after=$(psql "$NEW_DB" -t -c "SELECT COUNT(*) FROM listing_images li JOIN listings l ON l.id = li.listing_id WHERE l.source_type='c2c'" | xargs)
    echo "  - C2C Listings: $c2c_count_after"
    echo "  - Related Images: $images_count_after"
    echo ""
    echo "Backup saved at: $backup_file"
    echo ""
else
    echo ""
    echo -e "${RED}✗✗✗ Rollback failed!${NC}"
    echo "Database state preserved. Check logs above for errors."
    echo "You can restore from backup: $backup_file"
    exit 1
fi

echo "================================================================================"
echo "                     ROLLBACK COMPLETED"
echo "================================================================================"
echo ""
echo "The migration has been rolled back."
echo "You can now re-run the migration if needed:"
echo "  python3 /p/github.com/sveturs/listings/scripts/migrate_data.py"
echo ""
