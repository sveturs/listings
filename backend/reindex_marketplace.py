#!/usr/bin/env python3
"""
Marketplace Listings Reindexing Script - Fixed Version
"""
import psycopg2
from opensearchpy import OpenSearch

# PostgreSQL configuration
PG_CONN = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5433/svetubd?sslmode=disable"

# OpenSearch configuration
OS_INDEX = "marketplace_listings"

def main():
    # Connect to PostgreSQL
    conn = psycopg2.connect(PG_CONN)
    cursor = conn.cursor()

    # Connect to OpenSearch
    os_client = OpenSearch([{"host": "localhost", "port": 9200}],
                          http_compress=True, use_ssl=False, verify_certs=False)

    # Delete old index if exists
    if os_client.indices.exists(index=OS_INDEX):
        os_client.indices.delete(index=OS_INDEX)
        print(f"✅ Deleted old index: {OS_INDEX}")

    # Create new index with proper mapping
    os_client.indices.create(index=OS_INDEX, body={
        "settings": {
            "number_of_shards": 1,
            "number_of_replicas": 0,
            "analysis": {
                "analyzer": {
                    "serbian_analyzer": {
                        "type": "custom",
                        "tokenizer": "standard",
                        "filter": ["lowercase", "asciifolding"]
                    }
                }
            }
        },
        "mappings": {
            "properties": {
                "id": {"type": "integer"},
                "user_id": {"type": "integer"},
                "category_id": {"type": "integer"},
                "title": {"type": "text", "analyzer": "serbian_analyzer"},
                "description": {"type": "text", "analyzer": "serbian_analyzer"},
                "price": {"type": "float"},
                "condition": {"type": "keyword"},
                "status": {"type": "keyword"},
                "location": {"type": "text"},
                "city": {"type": "text"},
                "country": {"type": "text"},
                "created_at": {"type": "date"},
                "updated_at": {"type": "date"},
                "views_count": {"type": "integer"},
                "storefront_id": {"type": "integer"},
                "document_type": {"type": "keyword"},
                "images": {
                    "type": "nested",
                    "properties": {
                        "public_url": {"type": "keyword"},
                        "is_main": {"type": "boolean"}
                    }
                },
                "category": {
                    "properties": {
                        "id": {"type": "integer"},
                        "name": {"type": "text"},
                        "slug": {"type": "keyword"}
                    }
                }
            }
        }
    })
    print(f"✅ Created new index: {OS_INDEX}")

    # Get active listings
    cursor.execute('''
        SELECT l.id, l.user_id, l.category_id, l.title, l.description, l.price,
               l.condition, l.address_city, l.address_country,
               l.latitude, l.longitude, l.created_at, l.updated_at,
               l.views_count, l.storefront_id,
               c.name as category_name, c.slug as category_slug
        FROM c2c_listings l
        LEFT JOIN c2c_categories c ON l.category_id = c.id
        WHERE l.status = 'active'
    ''')

    listings = cursor.fetchall()
    print(f"📦 Found {len(listings)} active listings")

    # Index each listing
    for listing in listings:
        (listing_id, user_id, category_id, title, description, price, condition,
         city, country, lat, lng, created_at, updated_at, views_count, storefront_id,
         category_name, category_slug) = listing

        # Get images
        cursor.execute('''
            SELECT public_url, is_main
            FROM c2c_images
            WHERE listing_id = %s
            ORDER BY display_order
        ''', (listing_id,))
        images = cursor.fetchall()

        # Build document with all required fields
        doc = {
            "id": listing_id,
            "user_id": user_id,
            "category_id": category_id,
            "title": title,
            "description": description or "",
            "price": float(price) if price else 0,
            "condition": condition if condition else "",  # ✅ Empty string instead of null
            "location": "",  # ✅ Always include location field
            "city": city or "",
            "country": country or "",
            "created_at": created_at.isoformat() + 'Z' if created_at else None,
            "updated_at": updated_at.isoformat() + 'Z' if updated_at else None,
            "views_count": views_count or 0,
            "storefront_id": storefront_id,
            "status": "active",
            "document_type": "listing",
            "images": [{"public_url": img[0], "is_main": img[1]} for img in images if img[0]]
        }

        # Add category info
        if category_id:
            doc["category"] = {
                "id": category_id,
                "name": category_name,
                "slug": category_slug
            }

        # Index document
        os_client.index(index=OS_INDEX, id=listing_id, body=doc)
        print(f"✅ Indexed listing #{listing_id}: {title}")

    # Refresh index
    os_client.indices.refresh(index=OS_INDEX)
    print(f"\n🎉 Reindexing complete! Total: {len(listings)} listings")

    conn.close()

if __name__ == "__main__":
    main()
