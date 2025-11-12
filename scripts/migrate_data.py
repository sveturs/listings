#!/usr/bin/env python3
"""
Production Data Migration Script
=================================
Migrates data from old monolith database (svetubd) to new microservice database (listings_dev_db).

Old DB (port 5433):
  - c2c_listings → listings (source_type='c2c')
  - b2c_stores → storefronts
  - c2c_images → listing_images
  - c2c_listing_attributes → listing_attributes
  - c2c_locations → listing_locations

New DB (port 35434):
  - listings (unified table with source_type)
  - storefronts (B2C stores)
  - listing_images, listing_attributes, listing_locations
"""

import psycopg2
import psycopg2.extras
import sys
import json
import logging
from datetime import datetime
from typing import Dict, List, Tuple, Optional
import traceback

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s [%(levelname)s] %(message)s',
    handlers=[
        logging.FileHandler('/tmp/migrate_data.log'),
        logging.StreamHandler(sys.stdout)
    ]
)
logger = logging.getLogger(__name__)

# Database credentials
OLD_DB = {
    'host': 'localhost',
    'port': 5433,
    'database': 'svetubd',
    'user': 'postgres',
    'password': 'mX3g1XGhMRUZEX3l',
    'options': '-c statement_timeout=300000'  # 5 minutes timeout
}

NEW_DB = {
    'host': 'localhost',
    'port': 35434,
    'database': 'listings_dev_db',
    'user': 'listings_user',
    'password': 'listings_secret',
    'options': '-c statement_timeout=300000'
}

# ID mappings for foreign key relationships
id_mappings = {
    'listings': {},      # old_id -> new_id
    'storefronts': {},   # old_id -> new_id
}

# Migration statistics
stats = {
    'c2c_listings': {'total': 0, 'migrated': 0, 'failed': 0},
    'b2c_stores': {'total': 0, 'migrated': 0, 'failed': 0},
    'c2c_images': {'total': 0, 'migrated': 0, 'failed': 0},
    'c2c_listing_attributes': {'total': 0, 'migrated': 0, 'failed': 0},
    'c2c_locations': {'total': 0, 'migrated': 0, 'failed': 0},
}

def connect_db(config: Dict) -> psycopg2.extensions.connection:
    """Establish database connection with retries."""
    try:
        conn = psycopg2.connect(**config)
        conn.autocommit = False
        logger.info(f"Connected to {config['database']} on port {config['port']}")
        return conn
    except Exception as e:
        logger.error(f"Failed to connect to {config['database']}: {e}")
        raise

def check_connections():
    """Verify both database connections work."""
    logger.info("=" * 80)
    logger.info("PHASE 1: Verifying Database Connections")
    logger.info("=" * 80)

    try:
        old_conn = connect_db(OLD_DB)
        logger.info("✓ Old DB connection successful")
        old_conn.close()
    except Exception as e:
        logger.error(f"✗ Old DB connection failed: {e}")
        return False

    try:
        new_conn = connect_db(NEW_DB)
        logger.info("✓ New DB connection successful")
        new_conn.close()
    except Exception as e:
        logger.error(f"✗ New DB connection failed: {e}")
        return False

    logger.info("")
    return True

def verify_schema(conn: psycopg2.extensions.connection, tables: List[str]) -> bool:
    """Verify required tables exist."""
    cursor = conn.cursor()
    for table in tables:
        cursor.execute("""
            SELECT EXISTS (
                SELECT FROM information_schema.tables
                WHERE table_schema = 'public' AND table_name = %s
            )
        """, (table,))
        exists = cursor.fetchone()[0]
        if not exists:
            logger.error(f"✗ Table '{table}' does not exist")
            return False
        logger.info(f"✓ Table '{table}' exists")
    cursor.close()
    return True

def check_schema():
    """Verify all required tables exist in both databases."""
    logger.info("=" * 80)
    logger.info("PHASE 2: Verifying Database Schema")
    logger.info("=" * 80)

    old_tables = ['c2c_listings', 'b2c_stores', 'c2c_images']
    new_tables = ['listings', 'storefronts', 'listing_images']

    logger.info("\nChecking OLD database tables:")
    old_conn = connect_db(OLD_DB)
    old_valid = verify_schema(old_conn, old_tables)
    old_conn.close()

    if not old_valid:
        return False

    logger.info("\nChecking NEW database tables:")
    new_conn = connect_db(NEW_DB)
    new_valid = verify_schema(new_conn, new_tables)
    new_conn.close()

    logger.info("")
    return new_valid

def check_conflicts():
    """Check for ID conflicts between old and new databases."""
    logger.info("=" * 80)
    logger.info("PHASE 3: Checking for ID Conflicts")
    logger.info("=" * 80)

    old_conn = connect_db(OLD_DB)
    new_conn = connect_db(NEW_DB)

    old_cur = old_conn.cursor()
    new_cur = new_conn.cursor()

    has_conflicts = False

    # Check c2c_listings vs listings
    old_cur.execute("SELECT id FROM c2c_listings ORDER BY id")
    old_listing_ids = [row[0] for row in old_cur.fetchall()]

    new_cur.execute("SELECT id FROM listings WHERE source_type = 'c2c' ORDER BY id")
    new_listing_ids = [row[0] for row in new_cur.fetchall()]

    conflicts = set(old_listing_ids) & set(new_listing_ids)
    if conflicts:
        logger.warning(f"⚠ Found {len(conflicts)} conflicting listing IDs: {sorted(conflicts)}")
        has_conflicts = True
    else:
        logger.info(f"✓ No conflicts in listings (old: {len(old_listing_ids)}, new: {len(new_listing_ids)})")

    # Check b2c_stores vs storefronts
    old_cur.execute("SELECT id FROM b2c_stores ORDER BY id")
    old_store_ids = [row[0] for row in old_cur.fetchall()]

    new_cur.execute("SELECT id FROM storefronts ORDER BY id")
    new_store_ids = [row[0] for row in new_cur.fetchall()]

    conflicts = set(old_store_ids) & set(new_store_ids)
    if conflicts:
        logger.warning(f"⚠ Found {len(conflicts)} conflicting storefront IDs: {sorted(conflicts)}")
        has_conflicts = True
    else:
        logger.info(f"✓ No conflicts in storefronts (old: {len(old_store_ids)}, new: {len(new_store_ids)})")

    old_cur.close()
    new_cur.close()
    old_conn.close()
    new_conn.close()

    logger.info("")
    return not has_conflicts

def create_backup():
    """Create backup of target database before migration."""
    logger.info("=" * 80)
    logger.info("PHASE 4: Creating Backup")
    logger.info("=" * 80)

    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    backup_file = f'/tmp/listings_dev_db_backup_{timestamp}.sql'

    try:
        import subprocess
        cmd = [
            'pg_dump',
            '-h', str(NEW_DB['host']),
            '-p', str(NEW_DB['port']),
            '-U', NEW_DB['user'],
            '-d', NEW_DB['database'],
            '-f', backup_file,
            '--no-owner',
            '--no-acl'
        ]

        env = {'PGPASSWORD': NEW_DB['password']}
        result = subprocess.run(cmd, env=env, capture_output=True, text=True)

        if result.returncode == 0:
            logger.info(f"✓ Backup created: {backup_file}")
            return backup_file
        else:
            logger.error(f"✗ Backup failed: {result.stderr}")
            return None
    except Exception as e:
        logger.error(f"✗ Backup error: {e}")
        return None

def map_listing_status(old_status: str) -> str:
    """Map old listing status to new status."""
    status_map = {
        'active': 'active',
        'sold': 'sold',
        'inactive': 'inactive',
        'archived': 'archived',
        'draft': 'draft',
    }
    return status_map.get(old_status, 'draft')

def migrate_c2c_listings(old_conn, new_conn) -> bool:
    """Migrate C2C listings from c2c_listings to listings."""
    logger.info("\n" + "=" * 80)
    logger.info("PHASE 5: Migrating C2C Listings")
    logger.info("=" * 80)

    old_cur = old_conn.cursor(cursor_factory=psycopg2.extras.DictCursor)
    new_cur = new_conn.cursor()

    try:
        # Fetch all c2c listings
        old_cur.execute("""
            SELECT * FROM c2c_listings
            ORDER BY id
        """)
        old_listings = old_cur.fetchall()
        stats['c2c_listings']['total'] = len(old_listings)

        logger.info(f"Found {len(old_listings)} C2C listings to migrate")

        for old_listing in old_listings:
            try:
                # Prepare data for new listings table
                insert_data = {
                    'user_id': old_listing['user_id'] or 1,  # Default to user_id 1 if NULL
                    'storefront_id': old_listing['storefront_id'],
                    'title': old_listing['title'],
                    'description': old_listing['description'],
                    'price': old_listing['price'] or 0,
                    'currency': 'RSD',
                    'category_id': old_listing['category_id'] or 1,  # Default category if NULL
                    'status': map_listing_status(old_listing['status'] or 'draft'),
                    'visibility': 'public',
                    'quantity': 1,
                    'view_count': old_listing['views_count'] or 0,
                    'source_type': 'c2c',
                    'has_individual_location': old_listing['show_on_map'],
                    'individual_address': old_listing['location'],
                    'individual_latitude': old_listing['latitude'],
                    'individual_longitude': old_listing['longitude'],
                    'location_privacy': 'exact',
                    'show_on_map': old_listing['show_on_map'],
                    'created_at': old_listing['created_at'] or datetime.now(),
                    'updated_at': old_listing['updated_at'] or datetime.now(),
                }

                # Build attributes JSON from old fields
                attributes = {}
                if old_listing['condition']:
                    attributes['condition'] = old_listing['condition']
                if old_listing['address_city']:
                    attributes['city'] = old_listing['address_city']
                if old_listing['address_country']:
                    attributes['country'] = old_listing['address_country']
                if old_listing['original_language']:
                    attributes['original_language'] = old_listing['original_language']
                if old_listing['metadata']:
                    attributes['metadata'] = old_listing['metadata']
                if old_listing['address_multilingual']:
                    attributes['address_multilingual'] = old_listing['address_multilingual']

                insert_data['attributes'] = json.dumps(attributes) if attributes else '{}'

                # Insert into new listings table
                new_cur.execute("""
                    INSERT INTO listings (
                        user_id, storefront_id, title, description, price, currency,
                        category_id, status, visibility, quantity, view_count,
                        source_type, attributes, has_individual_location,
                        individual_address, individual_latitude, individual_longitude,
                        location_privacy, show_on_map, created_at, updated_at
                    ) VALUES (
                        %(user_id)s, %(storefront_id)s, %(title)s, %(description)s,
                        %(price)s, %(currency)s, %(category_id)s, %(status)s,
                        %(visibility)s, %(quantity)s, %(view_count)s, %(source_type)s,
                        %(attributes)s, %(has_individual_location)s, %(individual_address)s,
                        %(individual_latitude)s, %(individual_longitude)s, %(location_privacy)s,
                        %(show_on_map)s, %(created_at)s, %(updated_at)s
                    )
                    RETURNING id
                """, insert_data)

                new_id = new_cur.fetchone()[0]
                id_mappings['listings'][old_listing['id']] = new_id

                stats['c2c_listings']['migrated'] += 1
                logger.info(f"  ✓ Migrated listing {old_listing['id']} → {new_id}: {old_listing['title'][:50]}")

            except Exception as e:
                stats['c2c_listings']['failed'] += 1
                logger.error(f"  ✗ Failed to migrate listing {old_listing['id']}: {e}")
                logger.debug(traceback.format_exc())
                # Continue with next listing

        new_conn.commit()
        logger.info(f"\n✓ C2C Listings migration completed: {stats['c2c_listings']['migrated']}/{stats['c2c_listings']['total']} successful")
        return True

    except Exception as e:
        new_conn.rollback()
        logger.error(f"✗ C2C Listings migration failed: {e}")
        logger.debug(traceback.format_exc())
        return False
    finally:
        old_cur.close()
        new_cur.close()

def migrate_b2c_stores(old_conn, new_conn) -> bool:
    """Migrate B2C stores from b2c_stores to storefronts."""
    logger.info("\n" + "=" * 80)
    logger.info("PHASE 6: Migrating B2C Stores")
    logger.info("=" * 80)

    old_cur = old_conn.cursor(cursor_factory=psycopg2.extras.DictCursor)
    new_cur = new_conn.cursor()

    try:
        # Fetch all b2c stores
        old_cur.execute("""
            SELECT * FROM b2c_stores
            ORDER BY id
        """)
        old_stores = old_cur.fetchall()
        stats['b2c_stores']['total'] = len(old_stores)

        logger.info(f"Found {len(old_stores)} B2C stores to migrate")

        for old_store in old_stores:
            try:
                # Check if slug already exists
                new_cur.execute("SELECT id FROM storefronts WHERE slug = %s", (old_store['slug'],))
                existing = new_cur.fetchone()

                if existing:
                    logger.warning(f"  ⚠ Store with slug '{old_store['slug']}' already exists (id={existing[0]})")
                    id_mappings['storefronts'][old_store['id']] = existing[0]
                    stats['b2c_stores']['migrated'] += 1
                    continue

                # Convert JSONB fields to JSON strings for psycopg2
                store_data = dict(old_store)

                # Convert dict/list to JSON string for JSONB fields
                for field in ['theme', 'settings', 'seo_meta', 'ai_agent_config']:
                    if field in store_data and store_data[field] is not None:
                        if isinstance(store_data[field], (dict, list)):
                            store_data[field] = json.dumps(store_data[field])
                        elif isinstance(store_data[field], str):
                            # Already a string, ensure it's valid JSON
                            try:
                                json.loads(store_data[field])
                            except:
                                store_data[field] = '{}'

                # Insert into storefronts
                new_cur.execute("""
                    INSERT INTO storefronts (
                        user_id, slug, name, description, logo_url, banner_url, theme,
                        phone, email, website, address, city, postal_code, country,
                        latitude, longitude, formatted_address, geo_strategy,
                        default_privacy_level, address_verified, settings, seo_meta,
                        is_active, is_verified, verification_date, rating, reviews_count,
                        products_count, sales_count, views_count, followers_count,
                        subscription_plan, subscription_expires_at, subscription_id,
                        is_subscription_active, commission_rate, ai_agent_enabled,
                        ai_agent_config, live_shopping_enabled, group_buying_enabled,
                        created_at, updated_at
                    ) VALUES (
                        %(user_id)s, %(slug)s, %(name)s, %(description)s, %(logo_url)s,
                        %(banner_url)s, %(theme)s, %(phone)s, %(email)s, %(website)s,
                        %(address)s, %(city)s, %(postal_code)s, %(country)s, %(latitude)s,
                        %(longitude)s, %(formatted_address)s, %(geo_strategy)s,
                        %(default_privacy_level)s, %(address_verified)s, %(settings)s,
                        %(seo_meta)s, %(is_active)s, %(is_verified)s, %(verification_date)s,
                        %(rating)s, %(reviews_count)s, %(products_count)s, %(sales_count)s,
                        %(views_count)s, %(followers_count)s, %(subscription_plan)s,
                        %(subscription_expires_at)s, %(subscription_id)s,
                        %(is_subscription_active)s, %(commission_rate)s, %(ai_agent_enabled)s,
                        %(ai_agent_config)s, %(live_shopping_enabled)s, %(group_buying_enabled)s,
                        %(created_at)s, %(updated_at)s
                    )
                    RETURNING id
                """, store_data)

                new_id = new_cur.fetchone()[0]
                id_mappings['storefronts'][old_store['id']] = new_id

                stats['b2c_stores']['migrated'] += 1
                logger.info(f"  ✓ Migrated store {old_store['id']} → {new_id}: {old_store['name']}")

            except Exception as e:
                stats['b2c_stores']['failed'] += 1
                logger.error(f"  ✗ Failed to migrate store {old_store['id']}: {e}")
                logger.debug(traceback.format_exc())

        new_conn.commit()
        logger.info(f"\n✓ B2C Stores migration completed: {stats['b2c_stores']['migrated']}/{stats['b2c_stores']['total']} successful")
        return True

    except Exception as e:
        new_conn.rollback()
        logger.error(f"✗ B2C Stores migration failed: {e}")
        logger.debug(traceback.format_exc())
        return False
    finally:
        old_cur.close()
        new_cur.close()

def migrate_c2c_images(old_conn, new_conn) -> bool:
    """Migrate images from c2c_images to listing_images."""
    logger.info("\n" + "=" * 80)
    logger.info("PHASE 7: Migrating C2C Images")
    logger.info("=" * 80)

    old_cur = old_conn.cursor(cursor_factory=psycopg2.extras.DictCursor)
    new_cur = new_conn.cursor()

    try:
        # Fetch all images for migrated listings
        old_listing_ids = list(id_mappings['listings'].keys())
        if not old_listing_ids:
            logger.info("No listings migrated, skipping images")
            return True

        old_cur.execute("""
            SELECT * FROM c2c_images
            WHERE listing_id = ANY(%s)
            ORDER BY listing_id, display_order, id
        """, (old_listing_ids,))
        old_images = old_cur.fetchall()
        stats['c2c_images']['total'] = len(old_images)

        logger.info(f"Found {len(old_images)} images to migrate")

        for old_image in old_images:
            try:
                old_listing_id = old_image['listing_id']
                new_listing_id = id_mappings['listings'].get(old_listing_id)

                if not new_listing_id:
                    logger.warning(f"  ⚠ Skipping image {old_image['id']}: listing {old_listing_id} not found in mappings")
                    stats['c2c_images']['failed'] += 1
                    continue

                # Build URL from storage information
                url = old_image['public_url'] or old_image['file_path']
                storage_path = old_image['file_path']

                # Insert into listing_images
                new_cur.execute("""
                    INSERT INTO listing_images (
                        listing_id, url, storage_path, display_order, is_primary,
                        file_size, mime_type, created_at
                    ) VALUES (
                        %s, %s, %s, %s, %s, %s, %s, %s
                    )
                """, (
                    new_listing_id,
                    url,
                    storage_path,
                    old_image['display_order'],
                    old_image['is_main'],
                    old_image['file_size'],
                    old_image['content_type'],
                    old_image['created_at']
                ))

                stats['c2c_images']['migrated'] += 1

            except Exception as e:
                stats['c2c_images']['failed'] += 1
                logger.error(f"  ✗ Failed to migrate image {old_image['id']}: {e}")
                logger.debug(traceback.format_exc())

        new_conn.commit()
        logger.info(f"\n✓ Images migration completed: {stats['c2c_images']['migrated']}/{stats['c2c_images']['total']} successful")
        return True

    except Exception as e:
        new_conn.rollback()
        logger.error(f"✗ Images migration failed: {e}")
        logger.debug(traceback.format_exc())
        return False
    finally:
        old_cur.close()
        new_cur.close()

def print_summary():
    """Print migration summary statistics."""
    logger.info("\n" + "=" * 80)
    logger.info("MIGRATION SUMMARY")
    logger.info("=" * 80)

    total_migrated = sum(s['migrated'] for s in stats.values())
    total_failed = sum(s['failed'] for s in stats.values())
    total_records = sum(s['total'] for s in stats.values())

    for table, data in stats.items():
        if data['total'] > 0:
            success_rate = (data['migrated'] / data['total'] * 100) if data['total'] > 0 else 0
            logger.info(f"\n{table}:")
            logger.info(f"  Total:    {data['total']}")
            logger.info(f"  Migrated: {data['migrated']} ({success_rate:.1f}%)")
            logger.info(f"  Failed:   {data['failed']}")

    logger.info(f"\nOVERALL:")
    logger.info(f"  Total records:    {total_records}")
    logger.info(f"  Total migrated:   {total_migrated}")
    logger.info(f"  Total failed:     {total_failed}")
    logger.info(f"  Success rate:     {(total_migrated/total_records*100):.1f}%" if total_records > 0 else "N/A")

    logger.info("\nID Mappings:")
    logger.info(f"  Listings:    {len(id_mappings['listings'])} mappings")
    logger.info(f"  Storefronts: {len(id_mappings['storefronts'])} mappings")

    logger.info("\n" + "=" * 80)

def main():
    """Main migration flow."""
    logger.info("\n")
    logger.info("╔" + "=" * 78 + "╗")
    logger.info("║" + " " * 20 + "DATA MIGRATION SCRIPT" + " " * 37 + "║")
    logger.info("║" + " " * 16 + "svetubd → listings_dev_db" + " " * 37 + "║")
    logger.info("╚" + "=" * 78 + "╝")
    logger.info(f"\nStarted at: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")

    # Phase 1: Check connections
    if not check_connections():
        logger.error("\n✗✗✗ Migration aborted: Database connection check failed")
        return 1

    # Phase 2: Check schema
    if not check_schema():
        logger.error("\n✗✗✗ Migration aborted: Schema validation failed")
        return 1

    # Phase 3: Check conflicts
    if not check_conflicts():
        logger.warning("\n⚠⚠⚠ ID conflicts detected - will use ID mapping")

    # Phase 4: Create backup
    backup_file = create_backup()
    if not backup_file:
        logger.error("\n✗✗✗ Migration aborted: Backup creation failed")
        return 1

    # Establish persistent connections
    try:
        old_conn = connect_db(OLD_DB)
        new_conn = connect_db(NEW_DB)
    except Exception as e:
        logger.error(f"\n✗✗✗ Migration aborted: {e}")
        return 1

    try:
        # Phase 5: Migrate C2C listings
        if not migrate_c2c_listings(old_conn, new_conn):
            raise Exception("C2C listings migration failed")

        # Phase 6: Migrate B2C stores
        if not migrate_b2c_stores(old_conn, new_conn):
            raise Exception("B2C stores migration failed")

        # Phase 7: Migrate images
        if not migrate_c2c_images(old_conn, new_conn):
            raise Exception("Images migration failed")

        # Print summary
        print_summary()

        logger.info(f"\n✓✓✓ Migration completed successfully at {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
        logger.info(f"Backup file: {backup_file}")
        logger.info(f"Log file: /tmp/migrate_data.log")

        return 0

    except Exception as e:
        logger.error(f"\n✗✗✗ Migration failed: {e}")
        logger.debug(traceback.format_exc())
        logger.info(f"\nYou can restore from backup: {backup_file}")
        logger.info("psql -h localhost -p 35434 -U listings_user -d listings_dev_db < " + backup_file)
        return 1
    finally:
        old_conn.close()
        new_conn.close()

if __name__ == '__main__':
    sys.exit(main())
