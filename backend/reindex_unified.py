#!/usr/bin/env python3
"""
Unified Listings Reindexing Script
Реиндексация unified listings (C2C + B2C) из PostgreSQL в OpenSearch
"""

import json
import psycopg2
from opensearchpy import OpenSearch
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
OS_INDEX = "unified_listings"

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

def create_unified_index(os_client):
    """Создать unified индекс"""
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
                "location": {"type": "geo_point"},
                "images": {
                    "type": "nested",
                    "properties": {
                        "id": {"type": "integer"},
                        "url": {"type": "keyword"},
                        "thumbnail_url": {"type": "keyword"},
                        "is_main": {"type": "boolean"},
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

def get_c2c_listings(cursor):
    """Получить C2C listings"""
    query = """
        SELECT
            l.id, l.user_id, l.category_id, l.title, l.description,
            l.price, l.condition, l.status, l.latitude, l.longitude,
            l.created_at, l.updated_at, l.address_city, l.address_country,
            l.original_language, l.views_count,
            c.name as category_name, c.slug as category_slug
        FROM c2c_listings l
        LEFT JOIN c2c_categories c ON l.category_id = c.id
        WHERE l.status = 'active'
        ORDER BY l.id
    """
    cursor.execute(query)
    return cursor.fetchall()

def get_c2c_images(cursor, listing_id):
    """Получить изображения C2C listing"""
    query = """
        SELECT id, public_url, is_main
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
            "is_main": row[2],
            "display_order": display_order
        })
        display_order += 1
    return images

def get_c2c_translations(cursor, listing_id):
    """Получить переводы C2C listing"""
    query = """
        SELECT language, field_name, translated_text
        FROM translations
        WHERE entity_type = 'marketplace_listing' AND entity_id = %s
    """
    cursor.execute(query, (listing_id,))

    translations = {}
    for row in cursor.fetchall():
        lang = row[0]
        field = row[1]
        text = row[2]

        if lang not in translations:
            translations[lang] = {}
        translations[lang][field] = text

    return translations

def get_b2c_products(cursor):
    """Получить B2C products"""
    query = """
        SELECT
            p.id, s.user_id, p.category_id, p.name, p.description,
            p.price, p.is_active, p.created_at, p.updated_at,
            p.storefront_id, p.view_count,
            COALESCE(p.individual_latitude, s.latitude) as latitude,
            COALESCE(p.individual_longitude, s.longitude) as longitude,
            COALESCE(p.individual_address, s.city) as city,
            s.country,
            c.name as category_name, c.slug as category_slug,
            s.name as storefront_name, s.slug as storefront_slug
        FROM b2c_products p
        JOIN b2c_stores s ON p.storefront_id = s.id
        LEFT JOIN c2c_categories c ON p.category_id = c.id
        WHERE p.is_active = true
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

def get_b2c_translations(cursor, product_id):
    """Получить переводы B2C product"""
    query = """
        SELECT language, field_name, translated_text
        FROM translations
        WHERE entity_type = 'storefront_product' AND entity_id = %s
    """
    cursor.execute(query, (product_id,))

    translations = {}
    for row in cursor.fetchall():
        lang = row[0]
        field = row[1]
        text = row[2]

        if lang not in translations:
            translations[lang] = {}
        translations[lang][field] = text

    return translations

def build_c2c_document(listing_row, images, translations):
    """Построить документ C2C для OpenSearch"""
    (
        listing_id, user_id, category_id, title, description,
        price, condition, status, latitude, longitude,
        created_at, updated_at, city, country,
        original_language, views_count,
        category_name, category_slug
    ) = listing_row

    doc = {
        "id": listing_id,
        "source_type": "c2c",
        "user_id": user_id,
        "category_id": category_id,
        "title": title,
        "description": description,
        "price": price,
        "condition": condition,
        "status": status,
        "views_count": views_count,
        "original_language": original_language,
        "created_at": created_at.isoformat() if created_at else None,
        "updated_at": updated_at.isoformat() if updated_at else None,

        # Location
        "location": {
            "lat": latitude,
            "lon": longitude
        } if latitude and longitude else None,
        "city": city,
        "country": country,

        # Category
        "category": {
            "id": category_id,
            "name": category_name,
            "slug": category_slug
        } if category_name else None,

        # Images
        "images": images,
        "has_images": len(images) > 0,
        "image_count": len(images),

        # Translations
        "translations": translations if translations else {}
    }

    # Main image URL
    if images and len(images) > 0:
        main_image = next((img for img in images if img.get("is_main")), None)
        if main_image:
            doc["image_url"] = main_image["url"]
        else:
            doc["image_url"] = images[0]["url"]

    return doc

def build_b2c_document(product_row, images, translations):
    """Построить документ B2C для OpenSearch"""
    (
        product_id, user_id, category_id, name, description,
        price, is_active, created_at, updated_at,
        storefront_id, view_count,
        latitude, longitude, city, country,
        category_name, category_slug,
        storefront_name, storefront_slug
    ) = product_row

    doc = {
        "id": product_id,
        "source_type": "b2c",
        "user_id": user_id,
        "category_id": category_id,
        "title": name,
        "description": description,
        "price": price,
        "condition": "new",
        "status": "active" if is_active else "inactive",
        "views_count": view_count,
        "original_language": "sr",
        "storefront_id": storefront_id,
        "created_at": created_at.isoformat() if created_at else None,
        "updated_at": updated_at.isoformat() if updated_at else None,

        # Location
        "location": {
            "lat": latitude,
            "lon": longitude
        } if latitude and longitude else None,
        "city": city,
        "country": country,

        # Category
        "category": {
            "id": category_id,
            "name": category_name,
            "slug": category_slug
        } if category_name else None,

        # Storefront
        "storefront": {
            "id": storefront_id,
            "user_id": user_id,
            "name": storefront_name,
            "slug": storefront_slug
        } if storefront_name else None,

        # Images
        "images": images,
        "has_images": len(images) > 0,
        "image_count": len(images),

        # Translations
        "translations": translations if translations else {}
    }

    # Main image URL
    if images and len(images) > 0:
        main_image = next((img for img in images if img.get("is_main")), None)
        if main_image:
            doc["image_url"] = main_image["url"]
            doc["thumbnail_url"] = main_image.get("thumbnail_url")
        else:
            doc["image_url"] = images[0]["url"]
            doc["thumbnail_url"] = images[0].get("thumbnail_url")

    return doc

def reindex_unified():
    """Главная функция реиндексации"""
    print("=" * 80)
    print("🔄 Реиндексация Unified Listings (C2C + B2C)")
    print("=" * 80)

    # Подключение к БД
    print(f"\n📊 Подключение к PostgreSQL ({PG_HOST}:{PG_PORT}/{PG_DATABASE})...")
    conn = get_db_connection()
    cursor = conn.cursor()

    # Подключение к OpenSearch
    print(f"🔍 Подключение к OpenSearch ({OS_HOST}:{OS_PORT})...")
    os_client = get_opensearch_client()

    try:
        # Создать индекс
        print(f"\n🏗️  Создание unified индекса...")
        create_unified_index(os_client)

        # Реиндексация C2C listings
        print(f"\n📦 Получение C2C listings из PostgreSQL...")
        c2c_listings = get_c2c_listings(cursor)
        print(f"✅ Найдено {len(c2c_listings)} C2C listings")

        c2c_success = 0
        c2c_errors = 0

        for listing_row in c2c_listings:
            listing_id = listing_row[0]
            title = listing_row[3]

            try:
                print(f"\n🔄 C2C #{listing_id}: {title}")

                # Получаем изображения
                images = get_c2c_images(cursor, listing_id)
                print(f"   📸 Изображений: {len(images)}")

                # Получаем переводы
                translations = get_c2c_translations(cursor, listing_id)
                if translations:
                    print(f"   🌐 Переводы: {', '.join(translations.keys())}")

                # Создаем документ
                doc = build_c2c_document(listing_row, images, translations)

                # Индексируем
                response = os_client.index(
                    index=OS_INDEX,
                    id=f"c2c_{listing_id}",
                    body=doc,
                    refresh=False
                )

                if response.get("result") in ["created", "updated"]:
                    print(f"   ✅ Проиндексирован")
                    c2c_success += 1
                else:
                    print(f"   ⚠️  Неожиданный результат: {response}")
                    c2c_errors += 1

            except Exception as e:
                print(f"   ❌ Ошибка: {e}")
                c2c_errors += 1

        # Реиндексация B2C products
        print(f"\n📦 Получение B2C products из PostgreSQL...")
        b2c_products = get_b2c_products(cursor)
        print(f"✅ Найдено {len(b2c_products)} B2C products")

        b2c_success = 0
        b2c_errors = 0

        for product_row in b2c_products:
            product_id = product_row[0]
            name = product_row[3]

            try:
                print(f"\n🔄 B2C #{product_id}: {name}")

                # Получаем изображения
                images = get_b2c_images(cursor, product_id)
                print(f"   📸 Изображений: {len(images)}")

                # Получаем переводы
                translations = get_b2c_translations(cursor, product_id)
                if translations:
                    print(f"   🌐 Переводы: {', '.join(translations.keys())}")

                # Создаем документ
                doc = build_b2c_document(product_row, images, translations)

                # Индексируем
                response = os_client.index(
                    index=OS_INDEX,
                    id=f"b2c_{product_id}",
                    body=doc,
                    refresh=False
                )

                if response.get("result") in ["created", "updated"]:
                    print(f"   ✅ Проиндексирован")
                    b2c_success += 1
                else:
                    print(f"   ⚠️  Неожиданный результат: {response}")
                    b2c_errors += 1

            except Exception as e:
                print(f"   ❌ Ошибка: {e}")
                b2c_errors += 1

        # Refresh индекс
        print(f"\n🔄 Обновление индекса...")
        os_client.indices.refresh(index=OS_INDEX)

        # Итоги
        total_success = c2c_success + b2c_success
        total_errors = c2c_errors + b2c_errors
        total_items = len(c2c_listings) + len(b2c_products)

        print(f"\n" + "=" * 80)
        print(f"✅ Реиндексация завершена!")
        print(f"   • C2C успешно: {c2c_success}/{len(c2c_listings)}")
        print(f"   • B2C успешно: {b2c_success}/{len(b2c_products)}")
        print(f"   • Всего успешно: {total_success}/{total_items}")
        print(f"   • Ошибок: {total_errors}")
        print(f"=" * 80)

    finally:
        cursor.close()
        conn.close()
        print(f"\n🔌 Подключения закрыты")

if __name__ == "__main__":
    reindex_unified()
