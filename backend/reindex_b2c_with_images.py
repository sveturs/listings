#!/usr/bin/env python3
"""
B2C Products Reindexing Script
Реиндексация B2C товаров с изображениями из PostgreSQL в OpenSearch
"""

import json
import psycopg2
from opensearchpy import OpenSearch
from datetime import datetime

# PostgreSQL configuration
PG_HOST = "localhost"
PG_PORT = 5432
PG_USER = "postgres"
PG_PASSWORD = "mX3g1XGhMRUZEX3l"
PG_DATABASE = "svetubd"

# OpenSearch configuration
OS_HOST = "localhost"
OS_PORT = 9200
OS_INDEX = "b2c_products"

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

def get_product_images(cursor, product_id):
    """Получить изображения товара"""
    query = """
        SELECT
            id, storefront_product_id, image_url, thumbnail_url,
            display_order, is_default, created_at
        FROM storefront_product_images
        WHERE storefront_product_id = %s
        ORDER BY display_order ASC, id ASC
    """
    cursor.execute(query, (product_id,))

    images = []
    for row in cursor.fetchall():
        image = {
            "id": row[0],
            "url": row[2],  # image_url
            "thumbnail_url": row[3],
            "is_main": row[5],
            "is_default": row[5],
            "display_order": row[4],
        }
        images.append(image)

    return images

def get_product_translations(cursor, product_id):
    """Получить переводы товара"""
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

def get_all_products(cursor):
    """Получить все B2C товары"""
    query = """
        SELECT
            p.id, p.storefront_id, p.name, p.description, p.price, p.currency,
            p.category_id, p.sku, p.barcode, p.stock_quantity, p.stock_status,
            p.is_active, p.attributes, p.view_count, p.sold_count,
            p.created_at, p.updated_at,
            p.has_individual_location, p.individual_address, p.individual_latitude,
            p.individual_longitude, p.location_privacy, p.show_on_map, p.has_variants,
            c.name as category_name, c.slug as category_slug,
            s.user_id, s.name as storefront_name, s.slug as storefront_slug
        FROM storefront_products p
        LEFT JOIN c2c_categories c ON p.category_id = c.id
        LEFT JOIN storefronts s ON p.storefront_id = s.id
        ORDER BY p.id
    """
    cursor.execute(query)
    return cursor.fetchall()

def build_opensearch_document(product_row, images, translations):
    """Построить документ для OpenSearch"""
    (
        product_id, storefront_id, name, description, price, currency,
        category_id, sku, barcode, stock_quantity, stock_status,
        is_active, attributes, view_count, sold_count,
        created_at, updated_at,
        has_individual_location, individual_address, individual_latitude,
        individual_longitude, location_privacy, show_on_map, has_variants,
        category_name, category_slug,
        user_id, storefront_name, storefront_slug
    ) = product_row

    doc = {
        "id": product_id,
        "storefront_id": storefront_id,
        "name": name,
        "description": description,
        "price": price,
        "currency": currency,
        "category_id": category_id,
        "stock_quantity": stock_quantity,
        "stock_status": stock_status,
        "is_active": is_active,
        "view_count": view_count,
        "sold_count": sold_count,
        "has_variants": has_variants,
        "created_at": created_at.isoformat() if created_at else None,
        "updated_at": updated_at.isoformat() if updated_at else None,

        # Category info
        "category": {
            "id": category_id,
            "name": category_name,
            "slug": category_slug,
        } if category_name else None,

        # Storefront info
        "storefront": {
            "id": storefront_id,
            "user_id": user_id,
            "name": storefront_name,
            "slug": storefront_slug,
        } if storefront_name else None,

        # Images
        "images": images,
        "has_images": len(images) > 0,
        "image_count": len(images),

        # Translations
        "translations": translations if translations else {},

        # SKU/Barcode
        "sku": sku,
        "barcode": barcode,

        # Location
        "has_individual_location": has_individual_location,
        "individual_address": individual_address,
        "location": {
            "lat": individual_latitude,
            "lon": individual_longitude,
        } if individual_latitude and individual_longitude else None,
        "show_on_map": show_on_map,

        # Attributes
        "attributes": attributes if attributes else {},

        # Type marker
        "entity_type": "b2c_product",
        "listing_type": "b2c",
    }

    # Image URL для обратной совместимости
    if images and len(images) > 0:
        # Ищем главное изображение
        main_image = next((img for img in images if img.get("is_main")), None)
        if main_image:
            doc["image_url"] = main_image["url"]
            doc["thumbnail_url"] = main_image.get("thumbnail_url")
        else:
            # Берем первое
            doc["image_url"] = images[0]["url"]
            doc["thumbnail_url"] = images[0].get("thumbnail_url")

    return doc

def reindex_products():
    """Главная функция реиндексации"""
    print("=" * 80)
    print("🔄 Реиндексация B2C товаров с изображениями")
    print("=" * 80)

    # Подключение к БД
    print(f"\n📊 Подключение к PostgreSQL ({PG_HOST}:{PG_PORT}/{PG_DATABASE})...")
    conn = get_db_connection()
    cursor = conn.cursor()

    # Подключение к OpenSearch
    print(f"🔍 Подключение к OpenSearch ({OS_HOST}:{OS_PORT})...")
    os_client = get_opensearch_client()

    try:
        # Получаем все товары
        print(f"\n📦 Получение товаров из PostgreSQL...")
        products = get_all_products(cursor)
        print(f"✅ Найдено {len(products)} товаров")

        # Реиндексируем каждый товар
        success_count = 0
        error_count = 0

        for product_row in products:
            product_id = product_row[0]
            product_name = product_row[2]

            try:
                print(f"\n🔄 Обработка товара #{product_id}: {product_name}")

                # Получаем изображения
                images = get_product_images(cursor, product_id)
                print(f"   📸 Найдено изображений: {len(images)}")
                if images:
                    for img in images:
                        is_main = "✓" if img.get("is_main") else " "
                        print(f"      [{is_main}] {img['url']}")

                # Получаем переводы
                translations = get_product_translations(cursor, product_id)
                if translations:
                    print(f"   🌐 Переводы: {', '.join(translations.keys())}")

                # Создаем документ
                doc = build_opensearch_document(product_row, images, translations)

                # Индексируем в OpenSearch
                response = os_client.index(
                    index=OS_INDEX,
                    id=product_id,
                    body=doc,
                    refresh=True
                )

                if response.get("result") in ["created", "updated"]:
                    print(f"   ✅ Проиндексирован: {response.get('result')}")
                    success_count += 1
                else:
                    print(f"   ⚠️  Неожиданный результат: {response}")
                    error_count += 1

            except Exception as e:
                print(f"   ❌ Ошибка индексации товара #{product_id}: {e}")
                error_count += 1

        # Итоги
        print(f"\n" + "=" * 80)
        print(f"✅ Реиндексация завершена!")
        print(f"   • Успешно: {success_count}")
        print(f"   • Ошибок: {error_count}")
        print(f"=" * 80)

    finally:
        cursor.close()
        conn.close()
        print(f"\n🔌 Подключения закрыты")

if __name__ == "__main__":
    reindex_products()
