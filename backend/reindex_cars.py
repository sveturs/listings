#!/usr/bin/env python3

import psycopg2
import requests
import json
import sys
import time

# Database connection
DATABASE_URL = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd?sslmode=disable"

def get_car_listings():
    """Get all car listing IDs from the database"""
    conn = psycopg2.connect(DATABASE_URL)
    cur = conn.cursor()

    # Get all car listings (category 1301 = cars)
    cur.execute("""
        SELECT id, title
        FROM marketplace_listings
        WHERE category_id = 1301
        AND status = 'active'
        ORDER BY id
    """)

    listings = cur.fetchall()
    cur.close()
    conn.close()

    return listings

def reindex_via_api():
    """Call the backend API to reindex all listings"""
    url = "http://localhost:3000/api/v1/admin/reindex"

    # Use the JWT token for admin access
    headers = {
        "Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6ImFkbWluQGV4YW1wbGUuY29tIiwiZXhwIjoxNzU5MDczNTUzLCJpYXQiOjE3NTg5ODcxNTMsImlzX2FkbWluIjp0cnVlLCJ1c2VyX2lkIjoxfQ.YXwIj2b9-uoQZE0eGJ_gJRUH5pqqwWPMo73K-Vdf3A4",
        "Content-Type": "application/json"
    }

    try:
        response = requests.post(url, headers=headers, json={"category": "cars"})
        if response.status_code == 200:
            print("✅ Reindex API called successfully")
            return True
        else:
            print(f"❌ API returned status {response.status_code}: {response.text}")
            return False
    except Exception as e:
        print(f"❌ Error calling API: {e}")
        return False

def main():
    print("=== Переиндексация автомобильных объявлений ===")
    print(f"Время начала: {time.strftime('%Y-%m-%d %H:%M:%S')}")

    # Get car listings
    listings = get_car_listings()
    print(f"Найдено {len(listings)} автомобильных объявлений")

    if listings:
        print("Примеры объявлений:")
        for id, title in listings[:5]:
            print(f"  - ID {id}: {title}")

    # Try to reindex via API
    if not reindex_via_api():
        print("\nПытаемся альтернативный метод переиндексации...")
        # Here we would implement direct OpenSearch indexing if needed

    print(f"\nВремя завершения: {time.strftime('%Y-%m-%d %H:%M:%S')}")

if __name__ == "__main__":
    main()