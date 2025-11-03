#!/bin/bash
set -e  # Exit on error

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Database connection strings
SOURCE_DB="postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd"
TARGET_DB="postgres://listings_user:listings_secret@localhost:35434/listings_dev_db"

# Tables to migrate in correct order (respecting FK dependencies)
TABLES=(
    "c2c_categories"
    "c2c_listings"
    "c2c_chats"
    "c2c_favorites"
    "c2c_images"
    "c2c_listing_variants"
    "c2c_orders"
    "c2c_messages"
)

echo -e "${BLUE}=====================================${NC}"
echo -e "${BLUE}C2C Data Migration Script${NC}"
echo -e "${BLUE}Source: svetubd (5433)${NC}"
echo -e "${BLUE}Target: listings_dev_db (35434)${NC}"
echo -e "${BLUE}=====================================${NC}"
echo ""

# Function to check database connectivity
check_db_connection() {
    local db_name=$1
    local db_string=$2

    echo -ne "${YELLOW}Checking connection to ${db_name}...${NC}"
    if psql "$db_string" -c "SELECT 1;" > /dev/null 2>&1; then
        echo -e " ${GREEN}OK${NC}"
        return 0
    else
        echo -e " ${RED}FAILED${NC}"
        echo -e "${RED}Cannot connect to ${db_name}${NC}"
        exit 1
    fi
}

# Function to get row count
get_row_count() {
    local db_string=$1
    local table=$2

    psql "$db_string" -t -c "SELECT COUNT(*) FROM $table;" 2>/dev/null | xargs || echo "0"
}

# Function to check if table is empty
check_table_empty() {
    local table=$1
    local count=$(get_row_count "$TARGET_DB" "$table")

    if [ "$count" -gt 0 ]; then
        echo -e "${RED}Error: Table ${table} in target DB is not empty (${count} rows)${NC}"
        echo -e "${RED}Please truncate tables before migration or use --force flag${NC}"
        return 1
    fi
    return 0
}

# Function to migrate single table
migrate_table() {
    local table=$1

    echo -e "${YELLOW}Migrating ${table}...${NC}"

    # Get source row count
    local source_count=$(get_row_count "$SOURCE_DB" "$table")
    echo -e "  Source rows: ${source_count}"

    if [ "$source_count" -eq 0 ]; then
        echo -e "  ${YELLOW}Skipping (no data)${NC}"
        return 0
    fi

    # Create temporary dump file
    local dump_file="/tmp/${table}_dump_$$.sql"

    # Dump table data from source
    echo -ne "  Dumping data..."
    PGPASSWORD=mX3g1XGhMRUZEX3l pg_dump \
        -h localhost \
        -p 5433 \
        -U postgres \
        -d svetubd \
        --table="$table" \
        --data-only \
        --column-inserts \
        --no-owner \
        --no-acl \
        -f "$dump_file" 2>/dev/null

    if [ $? -eq 0 ]; then
        echo -e " ${GREEN}OK${NC}"
    else
        echo -e " ${RED}FAILED${NC}"
        rm -f "$dump_file"
        return 1
    fi

    # Load data into target
    echo -ne "  Loading data..."
    PGPASSWORD=listings_secret psql \
        -h localhost \
        -p 35434 \
        -U listings_user \
        -d listings_dev_db \
        -f "$dump_file" > /dev/null 2>&1

    if [ $? -eq 0 ]; then
        echo -e " ${GREEN}OK${NC}"
    else
        echo -e " ${RED}FAILED${NC}"
        rm -f "$dump_file"
        return 1
    fi

    # Cleanup
    rm -f "$dump_file"

    # Verify target row count
    local target_count=$(get_row_count "$TARGET_DB" "$table")
    echo -e "  Target rows: ${target_count}"

    if [ "$source_count" -eq "$target_count" ]; then
        echo -e "  ${GREEN}Verified: ${target_count}/${source_count} rows migrated${NC}"
    else
        echo -e "  ${RED}Warning: Row count mismatch (source: ${source_count}, target: ${target_count})${NC}"
    fi

    echo ""
}

# Function to fix sequences after migration
fix_sequences() {
    echo -e "${YELLOW}Fixing sequences...${NC}"

    # Fix sequences for tables with SERIAL/BIGSERIAL columns
    local sequences=(
        "c2c_categories:id:c2c_categories_id_seq"
        "c2c_listings:id:c2c_listings_id_seq"
        "c2c_chats:id:c2c_chats_id_seq"
        "c2c_favorites:id:c2c_favorites_id_seq"
        "c2c_images:id:c2c_images_id_seq"
        "c2c_listing_variants:id:c2c_listing_variants_id_seq"
        "c2c_orders:id:c2c_orders_id_seq"
        "c2c_messages:id:c2c_messages_id_seq"
    )

    for seq_info in "${sequences[@]}"; do
        IFS=':' read -r table column sequence <<< "$seq_info"

        echo -ne "  Fixing ${sequence}..."
        psql "$TARGET_DB" -c "SELECT setval('${sequence}', COALESCE((SELECT MAX(${column}) FROM ${table}), 1), true);" > /dev/null 2>&1

        if [ $? -eq 0 ]; then
            echo -e " ${GREEN}OK${NC}"
        else
            echo -e " ${YELLOW}SKIPPED${NC}"
        fi
    done

    echo ""
}

# Main migration flow
main() {
    # Step 1: Check connections
    echo -e "${BLUE}Step 1: Checking database connections${NC}"
    check_db_connection "source DB (svetubd)" "$SOURCE_DB"
    check_db_connection "target DB (listings_dev_db)" "$TARGET_DB"
    echo ""

    # Step 2: Check if tables are empty (unless --force flag is used)
    if [ "$1" != "--force" ]; then
        echo -e "${BLUE}Step 2: Verifying target tables are empty${NC}"
        all_empty=true
        for table in "${TABLES[@]}"; do
            if ! check_table_empty "$table"; then
                all_empty=false
            fi
        done

        if [ "$all_empty" = false ]; then
            echo ""
            echo -e "${RED}Migration aborted. Use --force to skip this check.${NC}"
            exit 1
        fi
        echo -e "${GREEN}All target tables are empty${NC}"
        echo ""
    else
        echo -e "${YELLOW}Skipping empty table check (--force flag)${NC}"
        echo ""
    fi

    # Step 3: Disable triggers
    echo -e "${BLUE}Step 3: Disabling triggers${NC}"
    for table in "${TABLES[@]}"; do
        echo -ne "  Disabling triggers on ${table}..."
        psql "$TARGET_DB" -c "ALTER TABLE $table DISABLE TRIGGER ALL;" > /dev/null 2>&1
        echo -e " ${GREEN}OK${NC}"
    done
    echo ""

    # Step 4: Migrate data
    echo -e "${BLUE}Step 4: Migrating data${NC}"
    total_migrated=0
    for table in "${TABLES[@]}"; do
        if migrate_table "$table"; then
            count=$(get_row_count "$TARGET_DB" "$table")
            total_migrated=$((total_migrated + count))
        else
            echo -e "${RED}Migration failed for ${table}${NC}"
            exit 1
        fi
    done

    # Step 5: Re-enable triggers
    echo -e "${BLUE}Step 5: Re-enabling triggers${NC}"
    for table in "${TABLES[@]}"; do
        echo -ne "  Enabling triggers on ${table}..."
        psql "$TARGET_DB" -c "ALTER TABLE $table ENABLE TRIGGER ALL;" > /dev/null 2>&1
        echo -e " ${GREEN}OK${NC}"
    done
    echo ""

    # Step 6: Fix sequences
    echo -e "${BLUE}Step 6: Fixing sequences${NC}"
    fix_sequences

    # Summary
    echo -e "${GREEN}=====================================${NC}"
    echo -e "${GREEN}Migration completed successfully!${NC}"
    echo -e "${GREEN}=====================================${NC}"
    echo -e "Total rows migrated: ${GREEN}${total_migrated}${NC}"
    echo ""
    echo "Row counts by table:"
    for table in "${TABLES[@]}"; do
        count=$(get_row_count "$TARGET_DB" "$table")
        printf "  %-25s %s\n" "$table:" "$count"
    done
    echo ""
}

# Run main function
main "$@"
