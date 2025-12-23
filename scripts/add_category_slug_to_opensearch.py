#!/usr/bin/env python3
"""
Add category_slug to existing OpenSearch documents
Phase 4 - PF-4.4: OpenSearch mapping optimization
Date: 2025-12-18
"""

import psycopg2
from opensearchpy import OpenSearch
from opensearchpy.helpers import scan, bulk
import sys

# PostgreSQL configuration
PG_HOST = "localhost"
PG_PORT = 35434
PG_USER = "listings_user"
PG_PASSWORD = "listings_secret"
PG_DATABASE = "listings_dev_db"

# OpenSearch configuration
OS_HOST = "localhost"
OS_PORT = 9200
OS_INDEX = "listings_microservice"

def get_db_connection():
    """Connect to PostgreSQL"""
    return psycopg2.connect(
        host=PG_HOST,
        port=PG_PORT,
        user=PG_USER,
        password=PG_PASSWORD,
        database=PG_DATABASE
    )

def get_opensearch_client():
    """Create OpenSearch client"""
    return OpenSearch(
        hosts=[{"host": OS_HOST, "port": OS_PORT}],
        http_compress=True,
        use_ssl=False,
        verify_certs=False,
    )

def get_category_slugs(conn):
    """Get mapping of category UUID -> slug"""
    cursor = conn.cursor()
    cursor.execute("SELECT id::text, slug FROM categories")

    mapping = {}
    for category_id, slug in cursor.fetchall():
        mapping[category_id] = slug

    cursor.close()
    print(f"✅ Loaded {len(mapping)} category slugs from database")
    return mapping

def update_documents(os_client, category_mapping):
    """Update all documents in OpenSearch with category_slug"""

    # Get all documents
    query = {"query": {"match_all": {}}}

    print(f"📊 Fetching documents from {OS_INDEX}...")
    docs = list(scan(os_client, index=OS_INDEX, query=query))
    total_docs = len(docs)
    print(f"✅ Found {total_docs} documents")

    if total_docs == 0:
        print("⚠️  No documents to update")
        return 0, 0

    # Prepare bulk update actions
    actions = []
    updated_count = 0
    skipped_count = 0

    for doc in docs:
        doc_id = doc['_id']
        source = doc['_source']

        # Get category_id (might be string UUID or not exist)
        category_id = source.get('category_id')

        # Skip if already has category_slug
        if 'category_slug' in source and source['category_slug']:
            skipped_count += 1
            continue

        # Skip if no category_id
        if not category_id:
            skipped_count += 1
            continue

        # Lookup category_slug
        category_id_str = str(category_id)
        category_slug = category_mapping.get(category_id_str)

        if not category_slug:
            print(f"⚠️  No category slug found for category_id: {category_id_str} (doc {doc_id})")
            skipped_count += 1
            continue

        # Prepare update action
        action = {
            "_op_type": "update",
            "_index": OS_INDEX,
            "_id": doc_id,
            "doc": {
                "category_slug": category_slug
            }
        }
        actions.append(action)
        updated_count += 1

    if updated_count == 0:
        print("✅ All documents already have category_slug (or no category_id)")
        return updated_count, skipped_count

    # Execute bulk update
    print(f"📝 Updating {updated_count} documents...")
    success, failed = bulk(os_client, actions, raise_on_error=False, raise_on_exception=False)

    if failed:
        print(f"⚠️  Some updates failed: {len(failed)} errors")
        for item in failed[:5]:  # Show first 5 errors
            print(f"   - {item}")

    print(f"✅ Successfully updated {success} documents")
    return success, skipped_count

def main():
    print("="*60)
    print("OpenSearch: Add category_slug to existing documents")
    print("="*60)
    print("")

    # Step 1: Connect to PostgreSQL
    print("[1/3] Connecting to PostgreSQL...")
    try:
        conn = get_db_connection()
        print(f"✅ Connected to {PG_DATABASE}")
    except Exception as e:
        print(f"❌ Failed to connect to PostgreSQL: {e}")
        sys.exit(1)

    # Step 2: Load category slugs
    print("")
    print("[2/3] Loading category slugs...")
    try:
        category_mapping = get_category_slugs(conn)
    except Exception as e:
        print(f"❌ Failed to load category slugs: {e}")
        conn.close()
        sys.exit(1)

    # Step 3: Update OpenSearch documents
    print("")
    print("[3/3] Updating OpenSearch documents...")
    try:
        os_client = get_opensearch_client()
        updated, skipped = update_documents(os_client, category_mapping)
    except Exception as e:
        print(f"❌ Failed to update documents: {e}")
        conn.close()
        sys.exit(1)

    # Cleanup
    conn.close()

    # Summary
    print("")
    print("="*60)
    print("Migration Complete!")
    print("="*60)
    print(f"Updated: {updated}")
    print(f"Skipped: {skipped}")
    print("")
    print("Next steps:")
    print("1. Verify documents:")
    print(f"   curl -s 'http://localhost:9200/{OS_INDEX}/_search?size=5' | jq '.hits.hits[]._source | {{title, category_slug}}'")
    print("")
    print("2. Test category filter:")
    print(f"   curl -s 'http://localhost:9200/{OS_INDEX}/_search' -H 'Content-Type: application/json' -d '{{\"query\": {{\"term\": {{\"category_slug\": \"elektronika\"}}}}}}' | jq '.hits.total.value'")
    print("")

if __name__ == "__main__":
    main()
