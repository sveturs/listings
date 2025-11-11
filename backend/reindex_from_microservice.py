#!/usr/bin/env python3
"""
Microservice Listings Reindexing Script
Реиндексация listings из микросервиса PostgreSQL в OpenSearch
"""

import json
import psycopg2
from opensearchpy import OpenSearch
from datetime import datetime

# PostgreSQL configuration - МИКРОСЕРВИС LISTINGS
PG_HOST = "localhost"
PG_PORT = 35434  # Порт микросервиса
PG_USER = "listings_user"
PG_PASSWORD = "listings_secret"
PG_DATABASE = "listings_dev_db"

# OpenSearch configuration
OS_HOST = "localhost"
OS_PORT = 9200
OS_INDEX = "marketplace_listings"  # Правильный индекс

def get_db_connection():
    """Создать подключение к PostgreSQL микросервиса"""
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

def create_marketplace_index(os_client):
    """Создать marketplace индекс"""
    index_body = {
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
                "source_type": {"type": "keyword"},
                "document_type": {"type": "keyword"},
                "title": {"type": "text", "analyzer": "serbian_analyzer"},
                "description": {"type": "text", "analyzer": "serbian_analyzer"},
                "price": {"type": "float"},
                "condition": {"type": "keyword"},
                "status": {"type": "keyword"},
                "category_id": {"type": "integer"},
                "user_id": {"type": "integer"},
                "storefront_id": {"type": "integer"},
                "created_at": {"type": "date"},
                "updated_at": {"type": "date"},
                "published_at": {"type": "date"},
                "location": {"type": "geo_point"},
                "images": {
                    "type": "nested",
                    "properties": {
                        "id": {"type": "integer"},
                        "url": {"type": "keyword"},
                        "thumbnail_url": {"type": "keyword"},
                        "is_primary": {"type": "boolean"},
                        "display_order": {"type": "integer"}
                    }
                },
                "category": {
                    "properties": {
                        "id": {"type": "integer"},
                        "name": {"type": "text"},
                        "slug": {"type": "keyword"}
                    }
                },
                "storefront": {
                    "properties": {
                        "id": {"type": "integer"},
                        "user_id": {"type": "integer"},
                        "name": {"type": "text"},
                        "slug": {"type": "keyword"}
                    }
                },
                "translations": {"type": "object", "enabled": False}
            }
        }
    }

    # Удалить старый индекс если существует
    if os_client.indices.exists(index=OS_INDEX):
        os_client.indices.delete(index=OS_INDEX)
        print(f"   🗑️  Удален старый индекс: {OS_INDEX}")

    # Создать новый индекс
    os_client.indices.create(index=OS_INDEX, body=index_body)
    print(f"   ✅ Создан новый индекс: {OS_INDEX}")

def get_listings(cursor):
    """Получить активные listings из микросервиса"""
    query = """
        SELECT
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.status, l.source_type,
            l.created_at, l.updated_at, l.published_at,
            ll.latitude, ll.longitude, ll.city, ll.country,
            c.name as category_name, c.slug as category_slug,
            l.storefront_id
        FROM listings l
        LEFT JOIN listing_locations ll ON l.id = ll.listing_id
        LEFT JOIN c2c_categories c ON l.category_id = c.id
        WHERE l.status = 'active' AND l.published_at IS NOT NULL
        ORDER BY l.id
    """
    cursor.execute(query)
    return cursor.fetchall()

def get_listing_images(cursor, listing_id):
    """Получить изображения listing из микросервиса"""
    query = """
        SELECT id, url, thumbnail_url, is_primary, display_order
        FROM listing_images
        WHERE listing_id = %s
        ORDER BY is_primary DESC, display_order ASC
    """
    cursor.execute(query, (listing_id,))

    images = []
    for row in cursor.fetchall():
        images.append({
            "id": row[0],
            "url": row[1],
            "thumbnail_url": row[2],
            "is_primary": row[3],
            "display_order": row[4]
        })
    return images

def get_storefront(cursor, storefront_id):
    """Получить данные storefront"""
    if not storefront_id:
        return None

    query = """
        SELECT id, user_id, name, slug
        FROM storefronts
        WHERE id = %s
    """
    cursor.execute(query, (storefront_id,))
    row = cursor.fetchone()

    if not row:
        return None

    return {
        "id": row[0],
        "user_id": row[1],
        "name": row[2],
        "slug": row[3]
    }

def build_listing_document(listing_row, images, storefront):
    """Построить документ listing для OpenSearch"""
    (
        listing_id, user_id, category_id, title, description,
        price, status, source_type,
        created_at, updated_at, published_at,
        latitude, longitude, city, country,
        category_name, category_slug,
        storefront_id
    ) = listing_row

    # condition зависит от source_type
    condition = "new" if source_type == "b2c" else "used"

    doc = {
        "id": listing_id,
        "source_type": source_type,
        "document_type": "listing",  # Важно для фильтрации в поиске
        "user_id": user_id,
        "category_id": category_id,
        "title": title,
        "description": description,
        "price": float(price) if price else None,
        "condition": condition,
        "status": status,
        "created_at": created_at.isoformat() if created_at else None,
        "updated_at": updated_at.isoformat() if updated_at else None,
        "published_at": published_at.isoformat() if published_at else None,

        # Location
        "location": {
            "lat": float(latitude),
            "lon": float(longitude)
        } if latitude and longitude else None,
        "city": city,
        "country": country,

        # Category
        "category": {
            "id": category_id,
            "name": category_name,
            "slug": category_slug
        } if category_name else None,

        # Storefront (для B2C)
        "storefront_id": storefront_id,
        "storefront": storefront,

        # Images
        "images": images,
        "has_images": len(images) > 0,
        "image_count": len(images),
    }

    # Main image URL
    if images and len(images) > 0:
        primary_image = next((img for img in images if img.get("is_primary")), None)
        if primary_image:
            doc["image_url"] = primary_image["url"]
            doc["thumbnail_url"] = primary_image.get("thumbnail_url")
        else:
            doc["image_url"] = images[0]["url"]
            doc["thumbnail_url"] = images[0].get("thumbnail_url")

    return doc

def reindex_from_microservice():
    """Главная функция реиндексации"""
    print("=" * 80)
    print("🔄 Реиндексация Listings из Микросервиса")
    print("=" * 80)

    # Подключение к БД микросервиса
    print(f"\n📊 Подключение к PostgreSQL микросервиса ({PG_HOST}:{PG_PORT}/{PG_DATABASE})...")
    conn = get_db_connection()
    cursor = conn.cursor()

    # Подключение к OpenSearch
    print(f"🔍 Подключение к OpenSearch ({OS_HOST}:{OS_PORT})...")
    os_client = get_opensearch_client()

    try:
        # Создать индекс
        print(f"\n🏗️  Создание индекса {OS_INDEX}...")
        create_marketplace_index(os_client)

        # Реиндексация listings
        print(f"\n📦 Получение активных listings из микросервиса...")
        listings = get_listings(cursor)
        print(f"✅ Найдено {len(listings)} активных listings")

        success = 0
        errors = 0

        for listing_row in listings:
            listing_id = listing_row[0]
            title = listing_row[3]
            status = listing_row[6]
            storefront_id = listing_row[17]

            try:
                print(f"\n🔄 Listing #{listing_id}: {title} (status: {status}, storefront_id: {storefront_id})")

                # Получаем изображения
                images = get_listing_images(cursor, listing_id)
                print(f"   📸 Изображений: {len(images)}")

                # Получаем storefront если есть
                storefront = None
                if storefront_id:
                    storefront = get_storefront(cursor, storefront_id)
                    if storefront:
                        print(f"   🏪 Storefront: {storefront['name']}")

                # Создаем документ
                doc = build_listing_document(listing_row, images, storefront)

                # Определяем ID документа
                source_type = doc["source_type"]
                doc_id = f"{source_type}_{listing_id}"

                # Индексируем
                response = os_client.index(
                    index=OS_INDEX,
                    id=doc_id,
                    body=doc,
                    refresh=False
                )

                if response.get("result") in ["created", "updated"]:
                    print(f"   ✅ Проиндексирован как {doc_id}")
                    success += 1
                else:
                    print(f"   ⚠️  Неожиданный результат: {response}")
                    errors += 1

            except Exception as e:
                print(f"   ❌ Ошибка: {e}")
                errors += 1

        # Refresh индекс
        print(f"\n🔄 Обновление индекса...")
        os_client.indices.refresh(index=OS_INDEX)

        # Итоги
        print(f"\n" + "=" * 80)
        print(f"✅ Реиндексация завершена!")
        print(f"   • Успешно: {success}/{len(listings)}")
        print(f"   • Ошибок: {errors}")
        print(f"=" * 80)

    finally:
        cursor.close()
        conn.close()
        print(f"\n🔌 Подключения закрыты")

if __name__ == "__main__":
    reindex_from_microservice()
