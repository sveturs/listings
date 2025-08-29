#!/usr/bin/env python3
"""
Скрипт для аудита и синхронизации изображений между MinIO и PostgreSQL
"""

import psycopg2
import subprocess
import json
import re
from datetime import datetime

# Конфигурация
DB_URL = "postgres://postgres:mX3g1XGhMRUZEX3l@localhost:5432/svetubd"
MINIO_ALIAS = "myminio"
MINIO_BUCKET = "listings"
STOREFRONT_BUCKET = "storefront-products"

def run_minio_cmd(cmd):
    """Выполняет команду MinIO через docker"""
    try:
        result = subprocess.run(
            f"docker exec minio mc {cmd}",
            shell=True, capture_output=True, text=True
        )
        if result.returncode == 0:
            return result.stdout.strip()
        return None
    except Exception as e:
        print(f"Ошибка MinIO: {e}")
        return None

def get_minio_files(listing_id):
    """Получает список файлов из MinIO для объявления"""
    output = run_minio_cmd(f"ls {MINIO_ALIAS}/{MINIO_BUCKET}/{listing_id}/")
    if not output:
        return []
    
    files = []
    for line in output.split('\n'):
        if line:
            # Парсим вывод mc ls
            parts = line.split()
            if len(parts) >= 5:
                filename = parts[-1]
                if filename and not filename.endswith('/'):
                    files.append(filename)
    return files

def audit_marketplace_listings():
    """Аудит изображений marketplace_listings"""
    print("=" * 60)
    print("АУДИТ MARKETPLACE LISTINGS")
    print("=" * 60)
    
    conn = psycopg2.connect(DB_URL)
    cur = conn.cursor()
    
    # Получаем все объявления с изображениями
    cur.execute("""
        SELECT DISTINCT listing_id 
        FROM marketplace_images 
        ORDER BY listing_id
    """)
    
    listings = cur.fetchall()
    total_listings = len(listings)
    problems = []
    
    print(f"Найдено объявлений с изображениями: {total_listings}\n")
    
    for (listing_id,) in listings:
        # Изображения из БД
        cur.execute("""
            SELECT id, file_path, file_name, is_main, public_url
            FROM marketplace_images 
            WHERE listing_id = %s 
            ORDER BY is_main DESC, id
        """, (listing_id,))
        
        db_images = cur.fetchall()
        db_count = len(db_images)
        
        # Изображения из MinIO
        minio_files = get_minio_files(listing_id)
        minio_count = len(minio_files)
        
        # Проверка синхронизации
        if db_count != minio_count:
            problems.append({
                'listing_id': listing_id,
                'issue': 'count_mismatch',
                'db_count': db_count,
                'minio_count': minio_count,
                'db_files': [img[2] for img in db_images],
                'minio_files': minio_files
            })
            print(f"❌ Объявление {listing_id}: БД={db_count}, MinIO={minio_count}")
        else:
            # Проверка путей
            has_wrong_path = False
            wrong_paths = []
            
            for img_id, file_path, file_name, is_main, public_url in db_images:
                # Проверяем наличие IP или неправильных путей
                if public_url and ('100.88.44.15' in public_url or 
                                  'localhost:9000' in public_url or
                                  not public_url.startswith('http')):
                    has_wrong_path = True
                    wrong_paths.append((img_id, public_url))
            
            if has_wrong_path:
                problems.append({
                    'listing_id': listing_id,
                    'issue': 'wrong_paths',
                    'wrong_paths': wrong_paths
                })
                print(f"⚠️  Объявление {listing_id}: неправильные пути")
            else:
                print(f"✅ Объявление {listing_id}: OK")
    
    cur.close()
    conn.close()
    
    return problems

def audit_storefront_products():
    """Аудит изображений storefront_products"""
    print("\n" + "=" * 60)
    print("АУДИТ STOREFRONT PRODUCTS")
    print("=" * 60)
    
    conn = psycopg2.connect(DB_URL)
    cur = conn.cursor()
    
    # Получаем все товары с изображениями
    cur.execute("""
        SELECT DISTINCT storefront_product_id 
        FROM storefront_product_images 
        ORDER BY storefront_product_id
    """)
    
    products = cur.fetchall()
    total_products = len(products)
    problems = []
    
    print(f"Найдено товаров с изображениями: {total_products}\n")
    
    for (product_id,) in products:
        # Изображения из БД
        cur.execute("""
            SELECT id, image_url, thumbnail_url, is_default
            FROM storefront_product_images 
            WHERE storefront_product_id = %s 
            ORDER BY display_order
        """, (product_id,))
        
        product_images = cur.fetchall()
        
        # Проверка URL
        has_wrong_url = False
        wrong_urls = []
        
        for img_id, image_url, thumbnail_url, is_default in product_images:
            if image_url and ('100.88.44.15' in image_url or 
                             not image_url.startswith(('http://', 'https://', '/'))):
                has_wrong_url = True
                wrong_urls.append((img_id, image_url))
        
        if has_wrong_url:
            problems.append({
                'product_id': product_id,
                'issue': 'wrong_urls',
                'wrong_urls': wrong_urls
            })
            print(f"⚠️  Товар {product_id}: неправильные URL")
        else:
            print(f"✅ Товар {product_id}: OK")
    
    cur.close()
    conn.close()
    
    return problems

def save_report(marketplace_problems, storefront_problems):
    """Сохраняет отчет в файл"""
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    report_file = f"/tmp/images_audit_{timestamp}.json"
    
    report = {
        'timestamp': timestamp,
        'marketplace_problems': marketplace_problems,
        'storefront_problems': storefront_problems,
        'total_problems': len(marketplace_problems) + len(storefront_problems)
    }
    
    with open(report_file, 'w') as f:
        json.dump(report, f, indent=2, ensure_ascii=False)
    
    print(f"\n📄 Отчет сохранен в: {report_file}")
    return report_file

def main():
    print("🔍 Запуск аудита изображений...")
    print("=" * 60)
    
    # Аудит marketplace
    marketplace_problems = audit_marketplace_listings()
    
    # Аудит storefronts
    storefront_problems = audit_storefront_products()
    
    # Сохраняем отчет
    report_file = save_report(marketplace_problems, storefront_problems)
    
    # Итоги
    print("\n" + "=" * 60)
    print("ИТОГИ АУДИТА")
    print("=" * 60)
    print(f"Найдено проблем в marketplace: {len(marketplace_problems)}")
    print(f"Найдено проблем в storefronts: {len(storefront_problems)}")
    print(f"Всего проблем: {len(marketplace_problems) + len(storefront_problems)}")
    
    # Показываем примеры проблем
    if marketplace_problems:
        print("\n📋 Примеры проблем в marketplace:")
        for problem in marketplace_problems[:5]:
            if problem['issue'] == 'count_mismatch':
                print(f"  • Объявление {problem['listing_id']}: несоответствие количества файлов")
                print(f"    БД: {problem['db_files']}")
                print(f"    MinIO: {problem['minio_files']}")
            elif problem['issue'] == 'wrong_paths':
                print(f"  • Объявление {problem['listing_id']}: неправильные пути")
                for img_id, path in problem['wrong_paths'][:2]:
                    print(f"    ID {img_id}: {path}")

if __name__ == "__main__":
    main()