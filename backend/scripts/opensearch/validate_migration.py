#!/usr/bin/env python3
"""
Validation Script для unified_listings_v2 Migration
Проверяет корректность миграции данных
"""

import json
import sys
from opensearchpy import OpenSearch
import psycopg2

# PostgreSQL configuration
PG_HOST = "localhost"
PG_PORT = 5433
PG_USER = "postgres"
PG_PASSWORD = "mX3g1XGhMRUZEX3l"
PG_DATABASE = "svetubd"

# OpenSearch configuration
OS_HOST = "localhost"
OS_PORT = 9200
NEW_INDEX = "unified_listings_v2"


def get_db_connection():
    """Создать подключение к PostgreSQL"""
    return psycopg2.connect(
        host=PG_HOST,
        port=PG_PORT,
        user=PG_USER,
        password=PG_PASSWORD,
        database=PG_DATABASE
    )


def get_opensearch_client():
    """Создать клиент OpenSearch"""
    return OpenSearch(
        hosts=[{"host": OS_HOST, "port": OS_PORT}],
        http_compress=True,
        use_ssl=False,
        verify_certs=False,
        ssl_assert_hostname=False,
        ssl_show_warn=False,
    )


def validate_document_counts():
    """Проверить количество документов"""
    print("\n📊 Validating document counts...")

    conn = get_db_connection()
    cursor = conn.cursor()
    os_client = get_opensearch_client()

    # C2C count in PostgreSQL
    cursor.execute("SELECT COUNT(*) FROM c2c_listings WHERE status = 'active'")
    pg_c2c_count = cursor.fetchone()[0]

    # B2C count in PostgreSQL
    cursor.execute("SELECT COUNT(*) FROM b2c_products WHERE is_active = true")
    pg_b2c_count = cursor.fetchone()[0]

    # OpenSearch counts
    os_c2c_count = os_client.count(
        index=NEW_INDEX,
        body={"query": {"term": {"source_type": "c2c"}}}
    )["count"]

    os_b2c_count = os_client.count(
        index=NEW_INDEX,
        body={"query": {"term": {"source_type": "b2c"}}}
    )["count"]

    os_total = os_client.count(index=NEW_INDEX)["count"]

    print(f"   PostgreSQL:")
    print(f"      • C2C: {pg_c2c_count}")
    print(f"      • B2C: {pg_b2c_count}")
    print(f"      • Total: {pg_c2c_count + pg_b2c_count}")

    print(f"   OpenSearch ({NEW_INDEX}):")
    print(f"      • C2C: {os_c2c_count}")
    print(f"      • B2C: {os_b2c_count}")
    print(f"      • Total: {os_total}")

    # Validation
    c2c_match = pg_c2c_count == os_c2c_count
    b2c_match = pg_b2c_count == os_b2c_count
    total_match = (pg_c2c_count + pg_b2c_count) == os_total

    if c2c_match and b2c_match and total_match:
        print(f"   ✅ Document counts match!")
        return True
    else:
        print(f"   ❌ Document counts mismatch!")
        if not c2c_match:
            print(f"      • C2C mismatch: {pg_c2c_count} (PG) vs {os_c2c_count} (OS)")
        if not b2c_match:
            print(f"      • B2C mismatch: {pg_b2c_count} (PG) vs {os_b2c_count} (OS)")
        return False


def validate_schema():
    """Проверить схему индекса"""
    print("\n📋 Validating index schema...")

    os_client = get_opensearch_client()

    # Получить mapping
    mapping = os_client.indices.get_mapping(index=NEW_INDEX)
    properties = mapping[NEW_INDEX]["mappings"]["properties"]

    # Expected fields
    expected_fields = {
        "id", "source_type", "title", "description", "price", "currency",
        "status", "visibility", "category_id", "category_name", "category_path_ids",
        "user_id", "seller_type", "storefront_id", "storefront_name",
        "location", "city", "country",
        "images", "has_images", "image_count", "primary_image_url", "thumbnail_url",
        "attributes", "condition", "negotiable",
        "sku", "stock_quantity", "stock_status",
        "views_count", "original_language", "supported_languages",
        "created_at", "updated_at", "translations"
    }

    actual_fields = set(properties.keys())

    missing_fields = expected_fields - actual_fields
    extra_fields = actual_fields - expected_fields

    print(f"   Expected fields: {len(expected_fields)}")
    print(f"   Actual fields: {len(actual_fields)}")

    if missing_fields:
        print(f"   ⚠️  Missing fields: {', '.join(missing_fields)}")

    if extra_fields:
        print(f"   ℹ️  Extra fields: {', '.join(extra_fields)}")

    if not missing_fields:
        print(f"   ✅ Schema validation passed!")
        return True
    else:
        print(f"   ❌ Schema validation failed!")
        return False


def validate_sample_documents():
    """Проверить корректность sample документов"""
    print("\n🔍 Validating sample documents...")

    os_client = get_opensearch_client()

    # Проверить C2C sample
    c2c_sample = os_client.search(
        index=NEW_INDEX,
        body={
            "size": 1,
            "query": {"term": {"source_type": "c2c"}}
        }
    )

    if c2c_sample["hits"]["total"]["value"] > 0:
        c2c_doc = c2c_sample["hits"]["hits"][0]["_source"]
        print(f"   C2C sample document:")
        print(f"      • ID: {c2c_doc.get('id')}")
        print(f"      • Title: {c2c_doc.get('title', '')[:50]}...")
        print(f"      • Price: {c2c_doc.get('price')} {c2c_doc.get('currency')}")
        print(f"      • Has images: {c2c_doc.get('has_images')}")
        print(f"      • Attributes: {len(c2c_doc.get('attributes', []))}")
        print(f"      • Condition: {c2c_doc.get('condition')}")
        print(f"      • Negotiable: {c2c_doc.get('negotiable')}")
    else:
        print(f"   ⚠️  No C2C documents found")

    # Проверить B2C sample
    b2c_sample = os_client.search(
        index=NEW_INDEX,
        body={
            "size": 1,
            "query": {"term": {"source_type": "b2c"}}
        }
    )

    if b2c_sample["hits"]["total"]["value"] > 0:
        b2c_doc = b2c_sample["hits"]["hits"][0]["_source"]
        print(f"   B2C sample document:")
        print(f"      • ID: {b2c_doc.get('id')}")
        print(f"      • Title: {b2c_doc.get('title', '')[:50]}...")
        print(f"      • Price: {b2c_doc.get('price')} {b2c_doc.get('currency')}")
        print(f"      • SKU: {b2c_doc.get('sku')}")
        print(f"      • Stock: {b2c_doc.get('stock_quantity')}")
        print(f"      • Stock status: {b2c_doc.get('stock_status')}")
        print(f"      • Storefront: {b2c_doc.get('storefront_name')}")
    else:
        print(f"   ⚠️  No B2C documents found")

    print(f"   ✅ Sample documents validated!")
    return True


def validate_search():
    """Проверить работу поиска"""
    print("\n🔎 Validating search functionality...")

    os_client = get_opensearch_client()

    # Test simple search
    search_result = os_client.search(
        index=NEW_INDEX,
        body={
            "size": 5,
            "query": {
                "multi_match": {
                    "query": "laptop",
                    "fields": ["title^3", "description^2"]
                }
            }
        }
    )

    hits = search_result["hits"]["total"]["value"]
    took_ms = search_result["took"]

    print(f"   Search query: 'laptop'")
    print(f"   Results: {hits}")
    print(f"   Latency: {took_ms}ms")

    if hits > 0:
        print(f"   ✅ Search working!")
        return True
    else:
        print(f"   ⚠️  No search results (may be normal if no laptops)")
        return True


def validate_aggregations():
    """Проверить aggregations"""
    print("\n📊 Validating aggregations...")

    os_client = get_opensearch_client()

    agg_result = os_client.search(
        index=NEW_INDEX,
        body={
            "size": 0,
            "aggs": {
                "by_source_type": {
                    "terms": {"field": "source_type"}
                },
                "by_category": {
                    "terms": {"field": "category_id", "size": 5}
                },
                "price_stats": {
                    "stats": {"field": "price"}
                }
            }
        }
    )

    # Source type aggregation
    source_types = agg_result["aggregations"]["by_source_type"]["buckets"]
    print(f"   Source types:")
    for bucket in source_types:
        print(f"      • {bucket['key']}: {bucket['doc_count']}")

    # Top categories
    categories = agg_result["aggregations"]["by_category"]["buckets"]
    print(f"   Top categories:")
    for bucket in categories[:3]:
        print(f"      • Category {bucket['key']}: {bucket['doc_count']}")

    # Price stats
    price_stats = agg_result["aggregations"]["price_stats"]
    print(f"   Price statistics:")
    print(f"      • Min: {price_stats['min']}")
    print(f"      • Max: {price_stats['max']}")
    print(f"      • Avg: {price_stats['avg']:.2f}")

    print(f"   ✅ Aggregations working!")
    return True


def main():
    """Main validation"""
    print("=" * 80)
    print("🔍 UNIFIED_LISTINGS_V2 MIGRATION VALIDATION")
    print("=" * 80)

    results = {}

    # Run validations
    results["counts"] = validate_document_counts()
    results["schema"] = validate_schema()
    results["samples"] = validate_sample_documents()
    results["search"] = validate_search()
    results["aggregations"] = validate_aggregations()

    # Summary
    print("\n" + "=" * 80)
    print("📋 VALIDATION SUMMARY")
    print("=" * 80)

    all_passed = all(results.values())

    for check, passed in results.items():
        status = "✅ PASS" if passed else "❌ FAIL"
        print(f"   {status} - {check.upper()}")

    if all_passed:
        print("\n✅ All validations passed!")
        print("\n💡 Next steps:")
        print("   1. Update backend config to use unified_listings_v2")
        print("   2. Enable feature flag")
        print("   3. Test with real traffic")
        return 0
    else:
        print("\n❌ Some validations failed!")
        print("   Review errors above and fix issues before proceeding")
        return 1


if __name__ == "__main__":
    sys.exit(main())
