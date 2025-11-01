#!/usr/bin/env python3
"""
OpenSearch Migration Script - Phase 2 Sprint 2.2
Миграция C2C listings + B2C products → unified_listings_v2

Features:
- Unified mapping (30 fields instead of 85)
- source_type differentiation ('c2c' | 'b2c')
- Attributes unification (nested structure)
- Dry-run mode для тестирования
- Валидация данных после миграции
"""

import json
import sys
import argparse
from typing import Any, Dict, List
import psycopg2
from opensearchpy import OpenSearch, helpers
from datetime import datetime

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


def load_unified_mapping():
    """Загрузить unified mapping из JSON файла"""
    mapping_file = "/p/github.com/sveturs/svetu/backend/scripts/opensearch/unified_mapping_v2.json"
    with open(mapping_file, 'r', encoding='utf-8') as f:
        return json.load(f)


def create_unified_index_v2(os_client, dry_run=False):
    """Создать unified index v2 с новым mapping"""
    if dry_run:
        print(f"   [DRY-RUN] Would create index: {NEW_INDEX}")
        return True

    # Удалить старый индекс если существует
    if os_client.indices.exists(index=NEW_INDEX):
        print(f"   ⚠️  Index {NEW_INDEX} exists. Deleting...")
        os_client.indices.delete(index=NEW_INDEX)

    # Загрузить mapping
    mapping = load_unified_mapping()

    # Создать индекс
    os_client.indices.create(index=NEW_INDEX, body=mapping)
    print(f"   ✅ Created index: {NEW_INDEX}")
    return True


def get_c2c_listings(cursor):
    """Получить C2C listings с атрибутами"""
    query = """
        SELECT
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.currency, l.condition, l.status,
            l.latitude, l.longitude, l.address_city, l.address_country,
            l.created_at, l.updated_at, l.views_count, l.original_language,
            l.negotiable,
            c.name as category_name,
            ARRAY_AGG(DISTINCT ci.id) FILTER (WHERE ci.id IS NOT NULL) as category_path_ids
        FROM c2c_listings l
        LEFT JOIN c2c_categories c ON l.category_id = c.id
        LEFT JOIN c2c_category_inheritance ci ON l.category_id = ci.descendant_id
        WHERE l.status = 'active'
        GROUP BY l.id, l.user_id, l.category_id, l.title, l.description,
                 l.price, l.currency, l.condition, l.status,
                 l.latitude, l.longitude, l.address_city, l.address_country,
                 l.created_at, l.updated_at, l.views_count, l.original_language,
                 l.negotiable, c.name
        ORDER BY l.id
    """
    cursor.execute(query)
    return cursor.fetchall()


def get_c2c_images(cursor, listing_id):
    """Получить изображения C2C listing"""
    query = """
        SELECT id, public_url, thumbnail_url, is_main
        FROM c2c_images
        WHERE listing_id = %s
        ORDER BY is_main DESC, id ASC
    """
    cursor.execute(query, (listing_id,))

    images = []
    display_order = 0
    for row in cursor.fetchall():
        images.append({
            "id": row[0],
            "url": row[1],
            "thumbnail_url": row[2],
            "is_main": row[3],
            "display_order": display_order
        })
        display_order += 1
    return images


def get_c2c_attributes(cursor, listing_id):
    """Получить атрибуты C2C listing в unified формате"""
    query = """
        SELECT
            a.id as attribute_id,
            a.name as attribute_name,
            a.type as attribute_type,
            av.string_value,
            av.numeric_value,
            av.boolean_value,
            av.translations
        FROM c2c_listing_attribute_values av
        JOIN c2c_attributes a ON av.attribute_id = a.id
        WHERE av.listing_id = %s
        ORDER BY a.id
    """
    cursor.execute(query, (listing_id,))

    attributes = []
    for row in cursor.fetchall():
        attr_id, attr_name, attr_type, str_val, num_val, bool_val, translations = row

        # Unified attribute structure
        attribute = {
            "attribute_id": attr_id,
            "attribute_name": attr_name,
            "attribute_type": attr_type
        }

        # Определяем значение на основе типа
        if attr_type == "string" and str_val:
            attribute["value"] = str_val
            attribute["display_value"] = str_val
        elif attr_type == "numeric" and num_val is not None:
            attribute["numeric_value"] = num_val
            attribute["value"] = str(num_val)
            attribute["display_value"] = str(num_val)
        elif attr_type == "boolean" and bool_val is not None:
            attribute["boolean_value"] = bool_val
            attribute["value"] = "true" if bool_val else "false"
            attribute["display_value"] = "Yes" if bool_val else "No"

        # Переводы display_name
        if translations:
            # Предполагаем что translations - это JSON с display_name
            if isinstance(translations, str):
                translations = json.loads(translations)
            attribute["display_name"] = translations.get("en", attr_name)

        attributes.append(attribute)

    return attributes


def build_c2c_document(listing_row, images, attributes):
    """Построить unified документ для C2C listing"""
    (
        listing_id, user_id, category_id, title, description,
        price, currency, condition, status,
        latitude, longitude, city, country,
        created_at, updated_at, views_count, original_language,
        negotiable, category_name, category_path_ids
    ) = listing_row

    doc = {
        "id": listing_id,
        "source_type": "c2c",
        "user_id": user_id,
        "seller_type": "user",
        "category_id": category_id,
        "category_name": category_name,
        "category_path_ids": category_path_ids or [],
        "title": title,
        "description": description,
        "price": price,
        "currency": currency or "RSD",
        "condition": condition,
        "negotiable": negotiable or False,
        "status": status,
        "visibility": "public",
        "views_count": views_count or 0,
        "original_language": original_language or "sr",
        "created_at": created_at.isoformat() if created_at else None,
        "updated_at": updated_at.isoformat() if updated_at else None,

        # Location
        "location": {
            "lat": latitude,
            "lon": longitude
        } if latitude and longitude else None,
        "city": city,
        "country": country or "RS",

        # Images
        "images": images,
        "has_images": len(images) > 0,
        "image_count": len(images),

        # Attributes
        "attributes": attributes
    }

    # Main image URL
    if images:
        main_image = next((img for img in images if img.get("is_main")), images[0])
        doc["primary_image_url"] = main_image["url"]
        doc["thumbnail_url"] = main_image.get("thumbnail_url") or main_image["url"]

    return doc


def get_b2c_products(cursor):
    """Получить B2C products"""
    query = """
        SELECT
            p.id, s.user_id, p.category_id, p.name, p.description,
            p.price, p.currency, p.is_active, p.sku, p.stock_quantity,
            COALESCE(p.individual_latitude, s.latitude) as latitude,
            COALESCE(p.individual_longitude, s.longitude) as longitude,
            COALESCE(p.individual_address, s.city) as city,
            s.country, p.created_at, p.updated_at, p.view_count,
            p.storefront_id,
            c.name as category_name,
            s.name as storefront_name,
            ARRAY_AGG(DISTINCT ci.id) FILTER (WHERE ci.id IS NOT NULL) as category_path_ids
        FROM b2c_products p
        JOIN b2c_stores s ON p.storefront_id = s.id
        LEFT JOIN c2c_categories c ON p.category_id = c.id
        LEFT JOIN c2c_category_inheritance ci ON p.category_id = ci.descendant_id
        WHERE p.is_active = true
        GROUP BY p.id, s.user_id, p.category_id, p.name, p.description,
                 p.price, p.currency, p.is_active, p.sku, p.stock_quantity,
                 latitude, longitude, city, s.country,
                 p.created_at, p.updated_at, p.view_count, p.storefront_id,
                 c.name, s.name
        ORDER BY p.id
    """
    cursor.execute(query)
    return cursor.fetchall()


def get_b2c_images(cursor, product_id):
    """Получить изображения B2C product"""
    query = """
        SELECT id, image_url, thumbnail_url, is_default, display_order
        FROM b2c_product_images
        WHERE storefront_product_id = %s
        ORDER BY is_default DESC, display_order ASC
    """
    cursor.execute(query, (product_id,))

    images = []
    for row in cursor.fetchall():
        images.append({
            "id": row[0],
            "url": row[1],
            "thumbnail_url": row[2],
            "is_main": row[3],
            "display_order": row[4]
        })
    return images


def get_b2c_attributes(cursor, product_id):
    """Получить атрибуты B2C product в unified формате"""
    # TODO: Implement when B2C attributes schema is available
    # For now, return empty list
    return []


def build_b2c_document(product_row, images, attributes):
    """Построить unified документ для B2C product"""
    (
        product_id, user_id, category_id, name, description,
        price, currency, is_active, sku, stock_quantity,
        latitude, longitude, city, country,
        created_at, updated_at, view_count, storefront_id,
        category_name, storefront_name, category_path_ids
    ) = product_row

    # Определить stock_status
    if stock_quantity is None:
        stock_status = "unknown"
    elif stock_quantity > 10:
        stock_status = "in_stock"
    elif stock_quantity > 0:
        stock_status = "low_stock"
    else:
        stock_status = "out_of_stock"

    doc = {
        "id": product_id,
        "source_type": "b2c",
        "user_id": user_id,
        "seller_type": "storefront",
        "storefront_id": storefront_id,
        "storefront_name": storefront_name,
        "category_id": category_id,
        "category_name": category_name,
        "category_path_ids": category_path_ids or [],
        "title": name,
        "description": description,
        "price": price,
        "currency": currency or "RSD",
        "condition": "new",  # B2C always new
        "negotiable": False,  # B2C not negotiable
        "sku": sku,
        "stock_quantity": stock_quantity or 0,
        "stock_status": stock_status,
        "status": "active" if is_active else "inactive",
        "visibility": "public",
        "views_count": view_count or 0,
        "original_language": "sr",
        "created_at": created_at.isoformat() if created_at else None,
        "updated_at": updated_at.isoformat() if updated_at else None,

        # Location
        "location": {
            "lat": latitude,
            "lon": longitude
        } if latitude and longitude else None,
        "city": city,
        "country": country or "RS",

        # Images
        "images": images,
        "has_images": len(images) > 0,
        "image_count": len(images),

        # Attributes
        "attributes": attributes
    }

    # Main image URL
    if images:
        main_image = next((img for img in images if img.get("is_main")), images[0])
        doc["primary_image_url"] = main_image["url"]
        doc["thumbnail_url"] = main_image.get("thumbnail_url") or main_image["url"]

    return doc


def migrate_data(os_client, dry_run=False):
    """Главная функция миграции данных"""
    print("=" * 80)
    print("🔄 Миграция C2C + B2C → unified_listings_v2")
    print("=" * 80)

    if dry_run:
        print("\n⚠️  DRY-RUN MODE - no changes will be made to OpenSearch")

    # Подключение к БД
    print(f"\n📊 Connecting to PostgreSQL ({PG_HOST}:{PG_PORT}/{PG_DATABASE})...")
    conn = get_db_connection()
    cursor = conn.cursor()

    # Подключение к OpenSearch
    print(f"🔍 Connecting to OpenSearch ({OS_HOST}:{OS_PORT})...")
    os_client_instance = get_opensearch_client()

    try:
        # Создать индекс
        print(f"\n🏗️  Creating unified index v2...")
        create_unified_index_v2(os_client_instance, dry_run)

        stats = {
            "c2c_total": 0,
            "c2c_success": 0,
            "c2c_errors": 0,
            "b2c_total": 0,
            "b2c_success": 0,
            "b2c_errors": 0
        }

        # Миграция C2C listings
        print(f"\n📦 Fetching C2C listings from PostgreSQL...")
        c2c_listings = get_c2c_listings(cursor)
        stats["c2c_total"] = len(c2c_listings)
        print(f"✅ Found {stats['c2c_total']} C2C listings")

        c2c_docs = []
        for listing_row in c2c_listings:
            listing_id = listing_row[0]
            title = listing_row[3]

            try:
                # Получаем изображения
                images = get_c2c_images(cursor, listing_id)

                # Получаем атрибуты
                attributes = get_c2c_attributes(cursor, listing_id)

                # Создаем документ
                doc = build_c2c_document(listing_row, images, attributes)
                doc["_id"] = f"c2c_{listing_id}"
                c2c_docs.append(doc)
                stats["c2c_success"] += 1

            except Exception as e:
                print(f"   ❌ C2C #{listing_id} error: {e}")
                stats["c2c_errors"] += 1

        # Миграция B2C products
        print(f"\n📦 Fetching B2C products from PostgreSQL...")
        b2c_products = get_b2c_products(cursor)
        stats["b2c_total"] = len(b2c_products)
        print(f"✅ Found {stats['b2c_total']} B2C products")

        b2c_docs = []
        for product_row in b2c_products:
            product_id = product_row[0]
            name = product_row[3]

            try:
                # Получаем изображения
                images = get_b2c_images(cursor, product_id)

                # Получаем атрибуты
                attributes = get_b2c_attributes(cursor, product_id)

                # Создаем документ
                doc = build_b2c_document(product_row, images, attributes)
                doc["_id"] = f"b2c_{product_id}"
                b2c_docs.append(doc)
                stats["b2c_success"] += 1

            except Exception as e:
                print(f"   ❌ B2C #{product_id} error: {e}")
                stats["b2c_errors"] += 1

        # Bulk indexing
        if not dry_run:
            print(f"\n🚀 Bulk indexing documents to OpenSearch...")

            # Подготовка documents для bulk API
            all_docs = c2c_docs + b2c_docs
            actions = []
            for doc in all_docs:
                doc_id = doc.pop("_id")
                actions.append({
                    "_index": NEW_INDEX,
                    "_id": doc_id,
                    "_source": doc
                })

            # Bulk index
            success_count, errors = helpers.bulk(
                os_client_instance,
                actions,
                stats_only=False,
                raise_on_error=False
            )

            print(f"✅ Indexed {success_count} documents")
            if errors:
                print(f"⚠️  {len(errors)} errors occurred")

            # Refresh индекс
            os_client_instance.indices.refresh(index=NEW_INDEX)

        # Итоги
        print(f"\n" + "=" * 80)
        print(f"✅ Migration completed!")
        print(f"   • C2C: {stats['c2c_success']}/{stats['c2c_total']} success, {stats['c2c_errors']} errors")
        print(f"   • B2C: {stats['b2c_success']}/{stats['b2c_total']} success, {stats['b2c_errors']} errors")
        print(f"   • Total: {stats['c2c_success'] + stats['b2c_success']}/{stats['c2c_total'] + stats['b2c_total']}")
        print(f"=" * 80)

        return stats

    finally:
        cursor.close()
        conn.close()


def validate_migration(os_client):
    """Валидация миграции"""
    print(f"\n📊 Validating migration...")

    # Подсчет документов
    c2c_count = os_client.count(index=NEW_INDEX, body={"query": {"term": {"source_type": "c2c"}}})["count"]
    b2c_count = os_client.count(index=NEW_INDEX, body={"query": {"term": {"source_type": "b2c"}}})["count"]
    total_count = os_client.count(index=NEW_INDEX)["count"]

    print(f"   • C2C documents: {c2c_count}")
    print(f"   • B2C documents: {b2c_count}")
    print(f"   • Total: {total_count}")

    # Проверка атрибутов
    sample_doc = os_client.search(
        index=NEW_INDEX,
        body={"size": 1, "query": {"match_all": {}}}
    )

    if sample_doc["hits"]["total"]["value"] > 0:
        fields = list(sample_doc["hits"]["hits"][0]["_source"].keys())
        print(f"   • Fields in document: {len(fields)}")
        print(f"   • Sample fields: {', '.join(fields[:10])}")

    return True


def main():
    """Main entry point"""
    parser = argparse.ArgumentParser(
        description="Migrate C2C + B2C to unified_listings_v2"
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="Dry run mode (no changes to OpenSearch)"
    )
    parser.add_argument(
        "--validate-only",
        action="store_true",
        help="Only validate existing unified_listings_v2 index"
    )

    args = parser.parse_args()

    os_client = get_opensearch_client()

    if args.validate_only:
        validate_migration(os_client)
        return

    # Run migration
    stats = migrate_data(os_client, dry_run=args.dry_run)

    # Validate if not dry-run
    if not args.dry_run:
        validate_migration(os_client)

    print("\n💡 Next steps:")
    print("   1. Test search in unified_listings_v2")
    print("   2. Update backend to use new index (feature flag)")
    print("   3. Monitor performance")
    print("   4. Delete old indexes after validation")


if __name__ == "__main__":
    main()
