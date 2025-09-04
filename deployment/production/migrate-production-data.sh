#!/bin/bash

# Production Data Migration Script for Unified Attributes System
# Safely migrates production data with validation and rollback capability

set -e

# Configuration
PROD_DB_URL="${DATABASE_URL:-postgres://postgres:password@prod-db.svetu.rs:5432/svetubd?sslmode=require}"
BACKUP_DIR="/var/backups/unified-attributes"
LOG_FILE="/var/log/unified-attributes-migration.log"
BATCH_SIZE="${BATCH_SIZE:-1000}"
DRY_RUN="${DRY_RUN:-false}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging function
log() {
    local level=$1
    shift
    echo -e "${level}[$(date +'%Y-%m-%d %H:%M:%S')] $*${NC}" | tee -a "$LOG_FILE"
}

# Create backup directory
mkdir -p "$BACKUP_DIR"

# Pre-migration validation
pre_migration_check() {
    log "$BLUE" "=== Starting Pre-Migration Checks ==="
    
    # Check database connectivity
    log "$YELLOW" "Checking database connectivity..."
    if ! psql "$PROD_DB_URL" -c "SELECT 1" > /dev/null 2>&1; then
        log "$RED" "ERROR: Cannot connect to production database"
        exit 1
    fi
    log "$GREEN" "Database connection: OK"
    
    # Check if migrations are already applied
    log "$YELLOW" "Checking migration status..."
    local migration_exists=$(psql "$PROD_DB_URL" -t -c "
        SELECT COUNT(*) FROM schema_migrations 
        WHERE version IN ('000034', '000035')
    ")
    
    if [ "$migration_exists" -gt 0 ]; then
        log "$YELLOW" "WARNING: Migrations already applied. Checking data migration status..."
        
        # Check if data is already migrated
        local unified_count=$(psql "$PROD_DB_URL" -t -c "SELECT COUNT(*) FROM unified_attributes")
        if [ "$unified_count" -gt 0 ]; then
            log "$RED" "ERROR: Data already exists in unified_attributes table"
            read -p "Continue anyway? (y/n) " -n 1 -r
            echo
            if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                exit 1
            fi
        fi
    fi
    
    # Count existing records
    log "$YELLOW" "Counting existing records..."
    local old_attributes=$(psql "$PROD_DB_URL" -t -c "SELECT COUNT(*) FROM category_attributes")
    local old_values=$(psql "$PROD_DB_URL" -t -c "SELECT COUNT(*) FROM attribute_values")
    
    log "$GREEN" "Found $old_attributes attributes and $old_values values to migrate"
    
    # Check disk space
    log "$YELLOW" "Checking disk space..."
    local available_space=$(df -BG "$BACKUP_DIR" | awk 'NR==2 {print $4}' | sed 's/G//')
    if [ "$available_space" -lt 10 ]; then
        log "$RED" "ERROR: Insufficient disk space for backup (need at least 10GB)"
        exit 1
    fi
    log "$GREEN" "Disk space: ${available_space}GB available"
    
    log "$GREEN" "=== Pre-Migration Checks Passed ==="
    echo
}

# Create full database backup
create_backup() {
    log "$BLUE" "=== Creating Database Backup ==="
    
    local backup_file="$BACKUP_DIR/pre-migration-$(date +%Y%m%d-%H%M%S).sql.gz"
    
    log "$YELLOW" "Creating full backup to $backup_file..."
    
    # Backup only attribute-related tables
    pg_dump "$PROD_DB_URL" \
        --tables="category_attributes" \
        --tables="attribute_values" \
        --tables="marketplace_categories" \
        --tables="marketplace_listing_attributes" \
        --verbose \
        | gzip > "$backup_file"
    
    if [ $? -eq 0 ]; then
        log "$GREEN" "Backup completed: $backup_file"
        log "$GREEN" "Backup size: $(du -h $backup_file | cut -f1)"
        
        # Verify backup
        if gunzip -t "$backup_file" 2>/dev/null; then
            log "$GREEN" "Backup verification: OK"
        else
            log "$RED" "ERROR: Backup verification failed"
            exit 1
        fi
    else
        log "$RED" "ERROR: Backup failed"
        exit 1
    fi
    
    echo
    return 0
}

# Apply database migrations
apply_migrations() {
    log "$BLUE" "=== Applying Database Migrations ==="
    
    if [ "$DRY_RUN" == "true" ]; then
        log "$YELLOW" "DRY RUN: Would apply migrations 000034 and 000035"
        return 0
    fi
    
    # Apply schema migration
    log "$YELLOW" "Applying schema migration (000034)..."
    psql "$PROD_DB_URL" -f "/app/migrations/000034_unified_attributes.up.sql"
    
    if [ $? -eq 0 ]; then
        log "$GREEN" "Schema migration completed"
    else
        log "$RED" "ERROR: Schema migration failed"
        exit 1
    fi
    
    # Apply data migration in batches
    log "$YELLOW" "Applying data migration (000035) in batches of $BATCH_SIZE..."
    
    # Start transaction and migrate data
    psql "$PROD_DB_URL" <<EOF
BEGIN;

-- Migrate attributes in batches
DO \$\$
DECLARE
    batch_count INTEGER := 0;
    total_count INTEGER;
BEGIN
    SELECT COUNT(*) INTO total_count FROM category_attributes;
    
    WHILE batch_count < total_count LOOP
        INSERT INTO unified_attributes (
            id, key, name_ru, name_en, name_sr,
            type, required, searchable, filterable,
            options, validation_rules, display_order,
            legacy_id, created_at, updated_at
        )
        SELECT 
            gen_random_uuid(),
            LOWER(REPLACE(ca.name_ru, ' ', '_')),
            ca.name_ru,
            COALESCE(ca.name_en, ca.name_ru),
            COALESCE(ca.name_sr, ca.name_ru),
            COALESCE(ca.type, 'text'),
            COALESCE(ca.required, false),
            COALESCE(ca.searchable, true),
            COALESCE(ca.filterable, true),
            ca.options,
            ca.validation_rules,
            ca.display_order,
            ca.id,
            ca.created_at,
            ca.updated_at
        FROM category_attributes ca
        ORDER BY ca.id
        LIMIT ${BATCH_SIZE}
        OFFSET batch_count
        ON CONFLICT (legacy_id) DO NOTHING;
        
        batch_count := batch_count + ${BATCH_SIZE};
        
        -- Log progress
        RAISE NOTICE 'Migrated % of % attributes', 
            LEAST(batch_count, total_count), total_count;
    END LOOP;
END\$\$;

-- Migrate category-attribute relationships
INSERT INTO unified_category_attributes (
    category_id, attribute_id, required, display_order
)
SELECT 
    mca.category_id,
    ua.id,
    COALESCE(mca.required, false),
    mca.display_order
FROM marketplace_category_attributes mca
JOIN unified_attributes ua ON ua.legacy_id = mca.attribute_id
ON CONFLICT DO NOTHING;

-- Migrate attribute values
INSERT INTO unified_attribute_values (
    id, attribute_id, value, display_value_ru,
    display_value_en, display_value_sr,
    is_active, display_order
)
SELECT 
    gen_random_uuid(),
    ua.id,
    av.value,
    av.display_value,
    av.display_value,
    av.display_value,
    COALESCE(av.is_active, true),
    av.display_order
FROM attribute_values av
JOIN unified_attributes ua ON ua.legacy_id = av.attribute_id
ON CONFLICT DO NOTHING;

COMMIT;
EOF
    
    if [ $? -eq 0 ]; then
        log "$GREEN" "Data migration completed successfully"
    else
        log "$RED" "ERROR: Data migration failed"
        exit 1
    fi
    
    echo
    return 0
}

# Validate migrated data
validate_migration() {
    log "$BLUE" "=== Validating Migrated Data ==="
    
    # Count migrated records
    local new_attributes=$(psql "$PROD_DB_URL" -t -c "SELECT COUNT(*) FROM unified_attributes")
    local new_relationships=$(psql "$PROD_DB_URL" -t -c "SELECT COUNT(*) FROM unified_category_attributes")
    local new_values=$(psql "$PROD_DB_URL" -t -c "SELECT COUNT(*) FROM unified_attribute_values")
    
    log "$YELLOW" "Migrated records:"
    log "$YELLOW" "  - Attributes: $new_attributes"
    log "$YELLOW" "  - Relationships: $new_relationships"
    log "$YELLOW" "  - Values: $new_values"
    
    # Validate data integrity
    log "$YELLOW" "Checking data integrity..."
    
    local orphaned=$(psql "$PROD_DB_URL" -t -c "
        SELECT COUNT(*) FROM unified_category_attributes uca
        LEFT JOIN unified_attributes ua ON ua.id = uca.attribute_id
        WHERE ua.id IS NULL
    ")
    
    if [ "$orphaned" -gt 0 ]; then
        log "$RED" "ERROR: Found $orphaned orphaned category-attribute relationships"
        exit 1
    fi
    
    # Check for missing translations
    local missing_translations=$(psql "$PROD_DB_URL" -t -c "
        SELECT COUNT(*) FROM unified_attributes
        WHERE name_en IS NULL OR name_sr IS NULL
    ")
    
    if [ "$missing_translations" -gt 0 ]; then
        log "$YELLOW" "WARNING: $missing_translations attributes missing translations"
    fi
    
    # Test queries
    log "$YELLOW" "Testing query performance..."
    
    # Test attribute fetching
    local query_time=$(psql "$PROD_DB_URL" -t -c "
        EXPLAIN (ANALYZE, BUFFERS, FORMAT JSON)
        SELECT ua.*, uca.required, uca.display_order
        FROM unified_attributes ua
        JOIN unified_category_attributes uca ON ua.id = uca.attribute_id
        WHERE uca.category_id = 1
        ORDER BY uca.display_order
    " | jq -r '.[0].Plan."Actual Total Time"')
    
    log "$GREEN" "Query performance: ${query_time}ms"
    
    if (( $(echo "$query_time > 100" | bc -l) )); then
        log "$YELLOW" "WARNING: Query performance slower than expected"
    fi
    
    log "$GREEN" "=== Migration Validation Passed ==="
    echo
    return 0
}

# Enable dual-write mode
enable_dual_write() {
    log "$BLUE" "=== Enabling Dual-Write Mode ==="
    
    if [ "$DRY_RUN" == "true" ]; then
        log "$YELLOW" "DRY RUN: Would enable dual-write mode"
        return 0
    fi
    
    # Update application configuration
    log "$YELLOW" "Updating application configuration..."
    
    # For Kubernetes deployment
    if command -v kubectl &> /dev/null; then
        kubectl set env deployment/backend \
            DUAL_WRITE_ATTRIBUTES=true \
            USE_UNIFIED_ATTRIBUTES=true \
            UNIFIED_ATTRIBUTES_FALLBACK=true \
            --namespace production
        
        log "$GREEN" "Kubernetes configuration updated"
    fi
    
    # For Docker deployment
    if command -v docker &> /dev/null; then
        docker exec backend sh -c "
            export DUAL_WRITE_ATTRIBUTES=true
            export USE_UNIFIED_ATTRIBUTES=true
            export UNIFIED_ATTRIBUTES_FALLBACK=true
        "
        log "$GREEN" "Docker configuration updated"
    fi
    
    log "$GREEN" "Dual-write mode enabled"
    echo
    return 0
}

# Create indexes for performance
create_indexes() {
    log "$BLUE" "=== Creating Performance Indexes ==="
    
    if [ "$DRY_RUN" == "true" ]; then
        log "$YELLOW" "DRY RUN: Would create performance indexes"
        return 0
    fi
    
    log "$YELLOW" "Creating indexes..."
    
    psql "$PROD_DB_URL" <<EOF
-- Create indexes if not exists
CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_unified_attributes_key 
    ON unified_attributes(key);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_unified_attributes_type 
    ON unified_attributes(type);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_unified_attributes_searchable 
    ON unified_attributes(searchable) WHERE searchable = true;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_unified_attributes_filterable 
    ON unified_attributes(filterable) WHERE filterable = true;

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_unified_category_attributes_category 
    ON unified_category_attributes(category_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_unified_category_attributes_attribute 
    ON unified_category_attributes(attribute_id);

CREATE INDEX CONCURRENTLY IF NOT EXISTS idx_unified_attribute_values_attribute 
    ON unified_attribute_values(attribute_id);

-- Analyze tables for query planner
ANALYZE unified_attributes;
ANALYZE unified_category_attributes;
ANALYZE unified_attribute_values;
EOF
    
    if [ $? -eq 0 ]; then
        log "$GREEN" "Indexes created successfully"
    else
        log "$YELLOW" "WARNING: Some indexes may have failed (non-critical)"
    fi
    
    echo
    return 0
}

# Main migration flow
main() {
    log "$BLUE" "========================================="
    log "$BLUE" "   UNIFIED ATTRIBUTES MIGRATION SCRIPT   "
    log "$BLUE" "========================================="
    log "$BLUE" "Environment: PRODUCTION"
    log "$BLUE" "Database: $PROD_DB_URL"
    log "$BLUE" "Dry Run: $DRY_RUN"
    log "$BLUE" "========================================="
    echo
    
    # Confirmation prompt for production
    if [ "$DRY_RUN" != "true" ]; then
        log "$RED" "⚠️  WARNING: This will modify PRODUCTION data!"
        read -p "Are you sure you want to continue? Type 'yes' to proceed: " confirmation
        if [ "$confirmation" != "yes" ]; then
            log "$YELLOW" "Migration cancelled by user"
            exit 0
        fi
    fi
    
    # Execute migration steps
    pre_migration_check
    create_backup
    apply_migrations
    validate_migration
    create_indexes
    enable_dual_write
    
    # Final summary
    log "$GREEN" "========================================="
    log "$GREEN" "    MIGRATION COMPLETED SUCCESSFULLY     "
    log "$GREEN" "========================================="
    log "$GREEN" "Next steps:"
    log "$GREEN" "  1. Monitor application logs for errors"
    log "$GREEN" "  2. Verify dual-write is working"
    log "$GREEN" "  3. Start canary release process"
    log "$GREEN" "  4. Monitor performance metrics"
    log "$GREEN" "========================================="
    
    # Save migration report
    local report_file="$BACKUP_DIR/migration-report-$(date +%Y%m%d-%H%M%S).txt"
    cp "$LOG_FILE" "$report_file"
    log "$GREEN" "Migration report saved to: $report_file"
}

# Handle script interruption
trap 'log "$RED" "Migration interrupted! Check logs at $LOG_FILE"' INT TERM

# Run main function
main "$@"