#!/usr/bin/env python3

import os
import sys
import requests
import psycopg2
import json
import time
from datetime import datetime

# Конфигурация
DB_URL = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"
OPENSEARCH_URL = "http://localhost:9200"
BACKEND_URL = "http://localhost:3000"

def check_opensearch():
    """Проверяет состояние OpenSearch"""
    try:
        response = requests.get(f"{OPENSEARCH_URL}/_cat/indices/marketplace_listings?format=json")
        if response.status_code == 200:
            data = response.json()[0]
            print(f"✅ OpenSearch индекс 'marketplace_listings' существует")
            print(f"   Документов: {data.get('docs.count', 0)}")
            print(f"   Размер: {data.get('store.size', 'N/A')}")
            print(f"   Статус: {data.get('status', 'N/A')}")
        else:
            print("❌ Индекс не найден")
            return False
    except Exception as e:
        print(f"❌ Ошибка подключения к OpenSearch: {e}")
        return False
    return True

def get_listings_count():
    """Получает количество объявлений в БД"""
    conn = psycopg2.connect(DB_URL)
    cur = conn.cursor()

    # Общее количество
    cur.execute("SELECT COUNT(*) FROM marketplace_listings WHERE status = 'active'")
    total = cur.fetchone()[0]

    # Количество автомобилей
    cur.execute("SELECT COUNT(*) FROM marketplace_listings WHERE status = 'active' AND category_id IN (1301, 1303)")
    cars = cur.fetchone()[0]

    # Количество с атрибутами
    cur.execute("""
        SELECT COUNT(DISTINCT ml.id)
        FROM marketplace_listings ml
        JOIN listing_attribute_values lav ON ml.id = lav.listing_id
        WHERE ml.status = 'active'
    """)
    with_attrs = cur.fetchone()[0]

    cur.close()
    conn.close()

    return total, cars, with_attrs

def reindex_direct():
    """Прямая переиндексация через OpenSearch API"""
    print("\n🔄 Начинаем прямую переиндексацию...")

    conn = psycopg2.connect(DB_URL)
    cur = conn.cursor()

    # Получаем все активные объявления с атрибутами и изображениями
    query = """
        WITH listing_images AS (
            SELECT
                listing_id,
                json_agg(
                    json_build_object(
                        'id', id,
                        'url', public_url,
                        'is_main', is_main,
                        'file_name', file_name
                    ) ORDER BY is_main DESC, id
                ) as images
            FROM marketplace_images
            GROUP BY listing_id
        ),
        listing_attributes AS (
            SELECT
                lav.listing_id,
                json_agg(
                    json_build_object(
                        'attribute_id', lav.attribute_id,
                        'attribute_name', ua.name,
                        'display_name', ua.display_name,
                        'attribute_type', ua.attribute_type,
                        'text_value', lav.text_value,
                        'numeric_value', lav.numeric_value,
                        'boolean_value', lav.boolean_value,
                        'json_value', lav.json_value,
                        'unit', lav.unit
                    )
                ) as attributes
            FROM listing_attribute_values lav
            JOIN unified_attributes ua ON lav.attribute_id = ua.id
            GROUP BY lav.listing_id
        )
        SELECT
            ml.id,
            ml.title,
            ml.description,
            ml.price,
            ml.category_id,
            ml.user_id,
            ml.status,
            ml.address_city,
            ml.address_country,
            ml.created_at,
            ml.updated_at,
            ml.condition,
            ml.views_count,
            ml.location,
            ml.show_on_map,
            ml.original_language,
            COALESCE(la.attributes, '[]'::json) as attributes,
            COALESCE(li.images, '[]'::json) as images
        FROM marketplace_listings ml
        LEFT JOIN listing_attributes la ON ml.id = la.listing_id
        LEFT JOIN listing_images li ON ml.id = li.listing_id
        WHERE ml.status = 'active'
        ORDER BY ml.id
        LIMIT 100
    """

    cur.execute(query)
    listings = cur.fetchall()

    print(f"📊 Найдено {len(listings)} объявлений для индексации")

    success_count = 0
    error_count = 0

    for listing in listings:
        listing_id = listing[0]

        # Формируем документ для индексации
        doc = {
            "id": listing_id,
            "title": listing[1],
            "description": listing[2],
            "price": float(listing[3]) if listing[3] else 0,
            "category_id": listing[4],
            "user_id": listing[5],
            "status": listing[6],
            "city": listing[7],
            "country": listing[8],
            "created_at": listing[9].isoformat() if listing[9] else None,
            "updated_at": listing[10].isoformat() if listing[10] else None,
            "condition": listing[11],
            "views_count": listing[12],
            "location": listing[13],
            "show_on_map": listing[14],
            "original_language": listing[15],
            "average_rating": 0,  # Позже можно добавить из отзывов
            "review_count": 0,    # Позже можно добавить из отзывов
            "attributes": [],
            "images": []
        }

        # Обрабатываем атрибуты
        attributes_json = listing[16]
        if attributes_json and isinstance(attributes_json, list):
            for attr in attributes_json:
                attr_doc = {
                    "attribute_id": attr.get("attribute_id"),
                    "attribute_name": attr.get("attribute_name"),
                    "display_name": attr.get("display_name"),
                    "attribute_type": attr.get("attribute_type"),
                }

                # Добавляем значение в зависимости от типа
                if attr.get("text_value"):
                    attr_doc["text_value"] = attr["text_value"]
                    attr_doc["text_value_lowercase"] = attr["text_value"].lower()
                if attr.get("numeric_value") is not None:
                    attr_doc["numeric_value"] = float(attr["numeric_value"])
                if attr.get("boolean_value") is not None:
                    attr_doc["boolean_value"] = attr["boolean_value"]
                if attr.get("json_value"):
                    attr_doc["json_value"] = json.dumps(attr["json_value"])
                if attr.get("unit"):
                    attr_doc["unit"] = attr["unit"]

                doc["attributes"].append(attr_doc)

        # Обрабатываем изображения
        images_json = listing[17]
        if images_json and isinstance(images_json, list):
            for img in images_json:
                img_doc = {
                    "id": img.get("id"),
                    "url": img.get("url"),
                    "is_main": img.get("is_main", False),
                    "file_name": img.get("file_name")
                }
                doc["images"].append(img_doc)

        # Индексируем в OpenSearch
        try:
            response = requests.put(
                f"{OPENSEARCH_URL}/marketplace_listings/_doc/{listing_id}",
                json=doc,
                headers={"Content-Type": "application/json"}
            )

            if response.status_code in [200, 201]:
                success_count += 1
                if doc["attributes"]:
                    print(f"✅ [{success_count}/{len(listings)}] Объявление {listing_id} проиндексировано с {len(doc['attributes'])} атрибутами")
                else:
                    print(f"✅ [{success_count}/{len(listings)}] Объявление {listing_id} проиндексировано без атрибутов")
            else:
                error_count += 1
                print(f"❌ [{success_count + error_count}/{len(listings)}] Ошибка индексации {listing_id}: {response.text}")

        except Exception as e:
            error_count += 1
            print(f"❌ [{success_count + error_count}/{len(listings)}] Ошибка индексации {listing_id}: {e}")

    cur.close()
    conn.close()

    print(f"\n📊 Результаты переиндексации:")
    print(f"   ✅ Успешно: {success_count}")
    print(f"   ❌ Ошибки: {error_count}")

    # Проверяем результат
    time.sleep(1)
    response = requests.get(f"{OPENSEARCH_URL}/marketplace_listings/_count")
    if response.status_code == 200:
        count = response.json()["count"]
        print(f"\n📈 Документов в индексе после переиндексации: {count}")

def verify_attributes():
    """Проверяет наличие атрибутов в индексе"""
    print("\n🔍 Проверка атрибутов в индексе...")

    # Ищем автомобили с атрибутами
    query = {
        "query": {
            "bool": {
                "must": [
                    {"term": {"category_id": 1301}},
                    {"exists": {"field": "attributes"}}
                ]
            }
        },
        "size": 1
    }

    response = requests.post(
        f"{OPENSEARCH_URL}/marketplace_listings/_search",
        json=query,
        headers={"Content-Type": "application/json"}
    )

    if response.status_code == 200:
        data = response.json()
        if data["hits"]["total"]["value"] > 0:
            doc = data["hits"]["hits"][0]["_source"]
            print(f"✅ Найдено объявление с атрибутами:")
            print(f"   ID: {doc['id']}")
            print(f"   Название: {doc['title']}")
            print(f"   Атрибутов: {len(doc.get('attributes', []))}")

            if doc.get("attributes"):
                print("   Примеры атрибутов:")
                for attr in doc["attributes"][:3]:
                    print(f"     - {attr.get('attribute_name', 'N/A')}: {attr.get('text_value') or attr.get('numeric_value') or attr.get('boolean_value', 'N/A')}")
        else:
            print("❌ Объявления с атрибутами не найдены в индексе")
    else:
        print(f"❌ Ошибка поиска: {response.text}")

def main():
    print("=" * 60)
    print("🚗 ПЕРЕИНДЕКСАЦИЯ МАРКЕТПЛЕЙСА С АТРИБУТАМИ")
    print("=" * 60)

    # Проверяем состояние
    print("\n📊 Текущее состояние:")
    total, cars, with_attrs = get_listings_count()
    print(f"   Всего объявлений в БД: {total}")
    print(f"   Автомобилей: {cars}")
    print(f"   С атрибутами: {with_attrs}")

    if not check_opensearch():
        print("❌ OpenSearch недоступен")
        return

    # Запускаем переиндексацию
    reindex_direct()

    # Проверяем результаты
    verify_attributes()

    print("\n✅ Переиндексация завершена!")

if __name__ == "__main__":
    main()