#!/usr/bin/env python3
"""
validate_data_consistency.py

Validates data consistency between monolith PostgreSQL and marketplace microservice.
- Compares listing counts
- Validates random sample (10 listings)
- Checks image URLs integrity
- Reports: 100% consistency or FAIL with details

Usage:
    python3 validate_data_consistency.py

Exit codes:
    0 - 100% consistency (PASS)
    1 - Data inconsistency detected (FAIL)
    2 - Error (connection issues, etc.)
"""

import sys
import os
import random
import psycopg2
import grpc
from datetime import datetime

# Import generated protobuf files
sys.path.insert(0, os.path.join(os.path.dirname(__file__), '..', 'internal', 'proj', 'marketplace_listings'))

try:
    import listings_pb2
    import listings_pb2_grpc
except ImportError:
    print("❌ ERROR: Failed to import protobuf files")
    print("   Make sure to generate protobuf files first:")
    print("   cd /p/github.com/sveturs/svetu/backend && make generate-proto")
    sys.exit(2)

# Configuration
DB_HOST = os.getenv('DB_HOST', 'localhost')
DB_PORT = os.getenv('DB_PORT', '5433')  # dev.svetu.rs PostgreSQL
DB_NAME = os.getenv('DB_NAME', 'svetu_dev_db')
DB_USER = os.getenv('DB_USER', 'svetu_dev_user')
DB_PASSWORD = os.getenv('DB_PASSWORD', 'svetu_dev_password')

GRPC_HOST = os.getenv('GRPC_HOST', 'localhost')
GRPC_PORT = os.getenv('GRPC_PORT', '50053')

SAMPLE_SIZE = 10  # Number of listings to validate in detail


def connect_to_postgres():
    """Connect to PostgreSQL monolith database."""
    try:
        conn = psycopg2.connect(
            host=DB_HOST,
            port=DB_PORT,
            database=DB_NAME,
            user=DB_USER,
            password=DB_PASSWORD,
            connect_timeout=5
        )
        return conn
    except Exception as e:
        print(f"❌ ERROR: Failed to connect to PostgreSQL: {e}")
        sys.exit(2)


def connect_to_grpc():
    """Connect to gRPC microservice."""
    try:
        channel = grpc.insecure_channel(f'{GRPC_HOST}:{GRPC_PORT}')
        stub = listings_pb2_grpc.MarketplaceListingsServiceStub(channel)
        return stub
    except Exception as e:
        print(f"❌ ERROR: Failed to connect to gRPC service: {e}")
        sys.exit(2)


def get_postgres_listing_count(conn):
    """Get total listing count from PostgreSQL."""
    try:
        cursor = conn.cursor()
        cursor.execute("SELECT COUNT(*) FROM marketplace_listings WHERE deleted_at IS NULL")
        count = cursor.fetchone()[0]
        cursor.close()
        return count
    except Exception as e:
        print(f"❌ ERROR: Failed to query PostgreSQL listing count: {e}")
        sys.exit(2)


def get_grpc_listing_count(stub):
    """Get total listing count from gRPC microservice."""
    try:
        request = listings_pb2.SearchListingsRequest(
            limit=1,  # We only need the total count
            offset=0
        )
        response = stub.SearchListings(request)
        return response.total
    except Exception as e:
        print(f"❌ ERROR: Failed to query gRPC listing count: {e}")
        sys.exit(2)


def get_random_listing_ids(conn, count):
    """Get random listing IDs from PostgreSQL."""
    try:
        cursor = conn.cursor()
        cursor.execute("""
            SELECT id FROM marketplace_listings
            WHERE deleted_at IS NULL
            ORDER BY RANDOM()
            LIMIT %s
        """, (count,))
        ids = [row[0] for row in cursor.fetchall()]
        cursor.close()
        return ids
    except Exception as e:
        print(f"❌ ERROR: Failed to get random listing IDs: {e}")
        sys.exit(2)


def get_postgres_listing(conn, listing_id):
    """Get listing details from PostgreSQL."""
    try:
        cursor = conn.cursor()
        cursor.execute("""
            SELECT
                id, title, description, price, currency, category_id,
                user_id, status, created_at, updated_at
            FROM marketplace_listings
            WHERE id = %s AND deleted_at IS NULL
        """, (listing_id,))
        row = cursor.fetchone()
        cursor.close()

        if not row:
            return None

        listing = {
            'id': row[0],
            'title': row[1],
            'description': row[2],
            'price': float(row[3]) if row[3] else 0.0,
            'currency': row[4],
            'category_id': row[5],
            'user_id': row[6],
            'status': row[7],
            'created_at': row[8],
            'updated_at': row[9]
        }

        # Get images
        cursor = conn.cursor()
        cursor.execute("""
            SELECT image_url, position
            FROM marketplace_listing_images
            WHERE listing_id = %s AND deleted_at IS NULL
            ORDER BY position
        """, (listing_id,))
        listing['images'] = [{'url': row[0], 'position': row[1]} for row in cursor.fetchall()]
        cursor.close()

        return listing
    except Exception as e:
        print(f"❌ ERROR: Failed to get PostgreSQL listing {listing_id}: {e}")
        return None


def get_grpc_listing(stub, listing_id):
    """Get listing details from gRPC microservice."""
    try:
        request = listings_pb2.GetListingRequest(id=listing_id)
        response = stub.GetListing(request)

        listing = {
            'id': response.id,
            'title': response.title,
            'description': response.description,
            'price': response.price,
            'currency': response.currency,
            'category_id': response.category_id,
            'user_id': response.user_id,
            'status': response.status,
            'created_at': response.created_at,
            'updated_at': response.updated_at,
            'images': [{'url': img.url, 'position': img.position} for img in response.images]
        }
        return listing
    except grpc.RpcError as e:
        if e.code() == grpc.StatusCode.NOT_FOUND:
            return None
        print(f"❌ ERROR: Failed to get gRPC listing {listing_id}: {e}")
        return None
    except Exception as e:
        print(f"❌ ERROR: Failed to get gRPC listing {listing_id}: {e}")
        return None


def compare_listings(pg_listing, grpc_listing):
    """Compare two listings and return list of differences."""
    differences = []

    if pg_listing is None and grpc_listing is None:
        return differences

    if pg_listing is None:
        differences.append("Listing exists in microservice but not in PostgreSQL")
        return differences

    if grpc_listing is None:
        differences.append("Listing exists in PostgreSQL but not in microservice")
        return differences

    # Compare fields
    if pg_listing['title'] != grpc_listing['title']:
        differences.append(f"Title mismatch: '{pg_listing['title']}' != '{grpc_listing['title']}'")

    if pg_listing['description'] != grpc_listing['description']:
        differences.append(f"Description mismatch (lengths: {len(pg_listing['description'])} != {len(grpc_listing['description'])})")

    # Price comparison with tolerance (0.01 for floating point)
    if abs(pg_listing['price'] - grpc_listing['price']) > 0.01:
        differences.append(f"Price mismatch: {pg_listing['price']} != {grpc_listing['price']}")

    if pg_listing['currency'] != grpc_listing['currency']:
        differences.append(f"Currency mismatch: '{pg_listing['currency']}' != '{grpc_listing['currency']}'")

    if pg_listing['category_id'] != grpc_listing['category_id']:
        differences.append(f"Category ID mismatch: {pg_listing['category_id']} != {grpc_listing['category_id']}")

    if pg_listing['user_id'] != grpc_listing['user_id']:
        differences.append(f"User ID mismatch: {pg_listing['user_id']} != {grpc_listing['user_id']}")

    if pg_listing['status'] != grpc_listing['status']:
        differences.append(f"Status mismatch: '{pg_listing['status']}' != '{grpc_listing['status']}'")

    # Compare images
    if len(pg_listing['images']) != len(grpc_listing['images']):
        differences.append(f"Images count mismatch: {len(pg_listing['images'])} != {len(grpc_listing['images'])}")
    else:
        for i, (pg_img, grpc_img) in enumerate(zip(pg_listing['images'], grpc_listing['images'])):
            if pg_img['url'] != grpc_img['url']:
                differences.append(f"Image {i} URL mismatch: '{pg_img['url']}' != '{grpc_img['url']}'")
            if pg_img['position'] != grpc_img['position']:
                differences.append(f"Image {i} position mismatch: {pg_img['position']} != {grpc_img['position']}")

    return differences


def main():
    """Main validation logic."""
    print("=" * 70)
    print("Data Consistency Validation: PostgreSQL ↔ Microservice")
    print("=" * 70)
    print(f"Timestamp: {datetime.now().strftime('%Y-%m-%d %H:%M:%S')}")
    print(f"Sample size: {SAMPLE_SIZE} listings")
    print("")

    # Connect to databases
    print("Connecting to PostgreSQL...")
    pg_conn = connect_to_postgres()
    print("✅ Connected to PostgreSQL")

    print("Connecting to gRPC microservice...")
    grpc_stub = connect_to_grpc()
    print("✅ Connected to gRPC microservice")
    print("")

    # Phase 1: Compare total counts
    print("-" * 70)
    print("Phase 1: Total Listing Count Comparison")
    print("-" * 70)

    pg_count = get_postgres_listing_count(pg_conn)
    grpc_count = get_grpc_listing_count(grpc_stub)

    print(f"PostgreSQL count: {pg_count}")
    print(f"Microservice count: {grpc_count}")

    count_match = pg_count == grpc_count
    if count_match:
        print("✅ PASS - Counts match")
    else:
        print(f"❌ FAIL - Count mismatch (delta: {abs(pg_count - grpc_count)})")

    print("")

    # Phase 2: Validate random sample
    print("-" * 70)
    print(f"Phase 2: Random Sample Validation ({SAMPLE_SIZE} listings)")
    print("-" * 70)

    if pg_count == 0:
        print("⚠ WARNING: No listings found in PostgreSQL - skipping sample validation")
        print("")
        pg_conn.close()

        if count_match:
            print("=" * 70)
            print("✅ VALIDATION PASSED (0 listings, counts match)")
            print("=" * 70)
            sys.exit(0)
        else:
            print("=" * 70)
            print("❌ VALIDATION FAILED (count mismatch with 0 listings)")
            print("=" * 70)
            sys.exit(1)

    sample_ids = get_random_listing_ids(pg_conn, min(SAMPLE_SIZE, pg_count))
    print(f"Selected {len(sample_ids)} random listings: {sample_ids}")
    print("")

    inconsistencies = []
    consistent_count = 0

    for listing_id in sample_ids:
        print(f"Validating listing ID {listing_id}...")

        pg_listing = get_postgres_listing(pg_conn, listing_id)
        grpc_listing = get_grpc_listing(grpc_stub, listing_id)

        differences = compare_listings(pg_listing, grpc_listing)

        if differences:
            print(f"  ❌ INCONSISTENT ({len(differences)} differences)")
            for diff in differences:
                print(f"     - {diff}")
            inconsistencies.append({
                'listing_id': listing_id,
                'differences': differences
            })
        else:
            print(f"  ✅ CONSISTENT")
            consistent_count += 1

        print("")

    # Close connections
    pg_conn.close()

    # Phase 3: Summary
    print("=" * 70)
    print("Validation Summary")
    print("=" * 70)
    print(f"Total listings checked: {len(sample_ids)}")
    print(f"Consistent: {consistent_count}")
    print(f"Inconsistent: {len(inconsistencies)}")
    print(f"Consistency rate: {(consistent_count / len(sample_ids) * 100):.1f}%")
    print("")

    # Final verdict
    all_passed = count_match and len(inconsistencies) == 0

    if all_passed:
        print("=" * 70)
        print("✅ VALIDATION PASSED - 100% DATA CONSISTENCY")
        print("=" * 70)
        print("")
        print("- Total count: ✅ MATCH")
        print(f"- Sample validation: ✅ {consistent_count}/{len(sample_ids)} consistent")
        print("")
        sys.exit(0)
    else:
        print("=" * 70)
        print("❌ VALIDATION FAILED - DATA INCONSISTENCY DETECTED")
        print("=" * 70)
        print("")
        print(f"- Total count: {'✅ MATCH' if count_match else '❌ MISMATCH'}")
        print(f"- Sample validation: ❌ {len(inconsistencies)}/{len(sample_ids)} inconsistent")
        print("")

        if inconsistencies:
            print("Inconsistent listings:")
            for inc in inconsistencies:
                print(f"  - Listing ID {inc['listing_id']}: {len(inc['differences'])} differences")

        print("")
        print("Action required: Investigate and fix data inconsistencies before proceeding")
        print("")
        sys.exit(1)


if __name__ == '__main__':
    try:
        main()
    except KeyboardInterrupt:
        print("\n⚠ Validation interrupted by user")
        sys.exit(2)
    except Exception as e:
        print(f"\n❌ UNEXPECTED ERROR: {e}")
        import traceback
        traceback.print_exc()
        sys.exit(2)
