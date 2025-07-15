#!/usr/bin/env python3
"""
Скрипт для обработки GeoJSON файлов с границами районов
и генерации SQL для обновления базы данных.
"""

import json
import psycopg2
from psycopg2.extras import RealDictCursor
import os
import glob
from typing import Dict, List, Optional
import logging
from datetime import datetime
from pathlib import Path

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Конфигурация
DATABASE_URL = os.getenv('DATABASE_URL', 'postgres://postgres:password@localhost:5432/hostel_db?sslmode=disable')
GEOJSON_DIR = '/data/hostel-booking-system/geodata'

# Маппинг названий районов для сопоставления
DISTRICT_NAME_VARIANTS = {
    # Белград
    'Земун': ['Zemun', 'Земун', 'ZEMUN'],
    'Стари Град': ['Stari Grad', 'Стари град', 'STARI GRAD', 'Old Town'],
    'Врачар': ['Vračar', 'Врачар', 'VRAČAR', 'Vracar'],
    'Савски венац': ['Savski Venac', 'Савски венац', 'SAVSKI VENAC'],
    'Палилула': ['Palilula', 'Палилула', 'PALILULA'],
    'Нови Београд': ['Novi Beograd', 'New Belgrade', 'Нови Београд', 'NOVI BEOGRAD'],
    'Чукарица': ['Čukarica', 'Чукарица', 'ČUKARICA', 'Cukarica'],
    'Звездара': ['Zvezdara', 'Звездара', 'ZVEZDARA'],
    'Раковица': ['Rakovica', 'Раковица', 'RAKOVICA'],
    'Вождовац': ['Voždovac', 'Вождовац', 'VOŽDOVAC', 'Vozdovac'],
    'Барајево': ['Barajevo', 'Барајево', 'BARAJEVO'],
    'Гроцка': ['Grocka', 'Гроцка', 'GROCKA'],
    'Лазаревац': ['Lazarevac', 'Лазаревац', 'LAZAREVAC'],
    'Младеновац': ['Mladenovac', 'Младеновац', 'MLADENOVAC'],
    'Обреновац': ['Obrenovac', 'Обреновац', 'OBRENOVAC'],
    'Сопот': ['Sopot', 'Сопот', 'SOPOT'],
    'Сурчин': ['Surčin', 'Сурчин', 'SURČIN', 'Surcin'],
    # Нови Сад
    'Лиман': ['Liman', 'Лиман', 'LIMAN'],
    'Грбавица': ['Grbavica', 'Грбавица', 'GRBAVICA'],
    'Ново насеље': ['Novo Naselje', 'Ново насеље', 'NOVO NASELJE'],
    'Детелинара': ['Detelinara', 'Детелинара', 'DETELINARA'],
    'Телеп': ['Telep', 'Телеп', 'TELEP'],
}


def get_db_connection():
    """Получить соединение с базой данных"""
    return psycopg2.connect(DATABASE_URL)


def get_districts_from_db() -> Dict[str, Dict]:
    """Получить список районов из базы данных"""
    conn = get_db_connection()
    try:
        with conn.cursor(cursor_factory=RealDictCursor) as cur:
            cur.execute("""
                SELECT d.id, d.name as district_name,
                       c.id as city_id, c.name as city_name
                FROM districts d
                JOIN cities c ON d.city_id = c.id
                ORDER BY c.name, d.name
            """)
            
            districts = {}
            for row in cur.fetchall():
                # Создаем ключ для быстрого поиска
                for variant in DISTRICT_NAME_VARIANTS.get(row['district_name'], [row['district_name']]):
                    key = f"{row['city_name'].lower()}_{variant.lower()}"
                    districts[key] = dict(row)
            
            return districts
    finally:
        conn.close()


def process_geojson_file(filepath: str, districts_db: Dict) -> List[Dict]:
    """Обработать GeoJSON файл и сопоставить с районами в БД"""
    logger.info(f"Обработка файла: {filepath}")
    
    with open(filepath, 'r', encoding='utf-8') as f:
        data = json.load(f)
    
    matched = []
    
    if data.get('type') == 'FeatureCollection':
        features = data.get('features', [])
    elif data.get('type') == 'Feature':
        features = [data]
    else:
        logger.warning(f"Неизвестный тип GeoJSON: {data.get('type')}")
        return matched
    
    for feature in features:
        properties = feature.get('properties', {})
        geometry = feature.get('geometry', {})
        
        # Извлекаем название района
        district_name = (
            properties.get('name') or 
            properties.get('NAME') or 
            properties.get('name:sr') or
            properties.get('name:sr-Cyrl') or
            properties.get('naziv') or
            properties.get('NAZIV')
        )
        
        # Извлекаем название города
        city_name = (
            properties.get('city') or
            properties.get('CITY') or
            properties.get('grad') or
            properties.get('GRAD') or
            'Београд'  # По умолчанию Белград
        )
        
        if district_name and geometry:
            # Пытаемся найти соответствие в БД
            db_match = find_district_match(district_name, city_name, districts_db)
            
            if db_match:
                matched.append({
                    'db_district': db_match,
                    'geometry': geometry,
                    'source_name': district_name,
                    'properties': properties
                })
                logger.info(f"Сопоставлен: {district_name} -> {db_match['district_name']}")
            else:
                logger.warning(f"Не найдено соответствие для: {district_name} ({city_name})")
    
    return matched


def find_district_match(name: str, city: str, districts_db: Dict) -> Optional[Dict]:
    """Найти соответствие района в БД"""
    # Прямой поиск
    key = f"{city.lower()}_{name.lower()}"
    if key in districts_db:
        return districts_db[key]
    
    # Поиск через варианты названий
    for db_name, variants in DISTRICT_NAME_VARIANTS.items():
        if name in variants or name.lower() in [v.lower() for v in variants]:
            # Пробуем разные варианты города
            for city_variant in [city, 'Београд', 'Belgrade', 'Нови Сад', 'Novi Sad']:
                key = f"{city_variant.lower()}_{db_name.lower()}"
                if key in districts_db:
                    return districts_db[key]
    
    return None


def simplify_geometry(geometry: Dict, tolerance: float = 0.0001) -> Dict:
    """Упростить геометрию для оптимизации"""
    # Здесь можно добавить алгоритм упрощения полигонов
    # Пока возвращаем как есть
    return geometry


def generate_sql_updates(matched_data: List[Dict]) -> str:
    """Генерировать SQL для обновления границ районов"""
    sql_lines = [
        "-- SQL скрипт для обновления границ районов из GeoJSON файлов",
        f"-- Сгенерировано: {datetime.now().isoformat()}",
        "",
        "BEGIN;",
        "",
        "-- Создаем резервную копию текущих границ",
        f"CREATE TABLE IF NOT EXISTS districts_boundaries_backup_{datetime.now().strftime('%Y%m%d')} AS ",
        "SELECT id, name, boundary, updated_at FROM districts;",
        ""
    ]
    
    # Группируем по городам
    by_city = {}
    for match in matched_data:
        city = match['db_district']['city_name']
        if city not in by_city:
            by_city[city] = []
        by_city[city].append(match)
    
    update_count = 0
    for city_name, matches in by_city.items():
        sql_lines.append(f"-- Обновления для города {city_name}")
        
        for match in matches:
            district = match['db_district']
            geometry = simplify_geometry(match['geometry'])
            geometry_json = json.dumps(geometry)
            
            sql_lines.append(f"""
-- Район: {district['district_name']} (источник: {match['source_name']})
UPDATE districts 
SET boundary = ST_GeomFromGeoJSON('{geometry_json}'),
    area_km2 = ST_Area(ST_GeomFromGeoJSON('{geometry_json}')::geography) / 1000000,
    updated_at = NOW()
WHERE id = '{district['id']}';
""")
            update_count += 1
        
        sql_lines.append("")
    
    # Добавляем проверки
    sql_lines.extend([
        f"-- Всего обновлений: {update_count}",
        "",
        "-- Проверка валидности геометрии",
        "SELECT c.name as city, d.name as district, ",
        "       ST_IsValid(d.boundary) as is_valid,",
        "       ST_Area(d.boundary::geography) / 1000000 as area_km2,",
        "       ST_NPoints(d.boundary) as num_points",
        "FROM districts d",
        "JOIN cities c ON d.city_id = c.id",
        "WHERE d.updated_at > NOW() - INTERVAL '1 minute'",
        "ORDER BY c.name, d.name;",
        "",
        "-- Исправление невалидных геометрий",
        "UPDATE districts",
        "SET boundary = ST_MakeValid(boundary)",
        "WHERE NOT ST_IsValid(boundary) AND boundary IS NOT NULL;",
        "",
        "COMMIT;"
    ])
    
    return '\n'.join(sql_lines)


def create_sample_belgrade_geojson():
    """Создать пример GeoJSON файла для районов Белграда"""
    sample_data = {
        "type": "FeatureCollection",
        "features": [
            {
                "type": "Feature",
                "properties": {
                    "name": "Стари Град",
                    "city": "Београд",
                    "description": "Исторический центр Белграда"
                },
                "geometry": {
                    "type": "Polygon",
                    "coordinates": [[
                        [20.4500, 44.8200],
                        [20.4550, 44.8250],
                        [20.4650, 44.8200],
                        [20.4700, 44.8150],
                        [20.4650, 44.8100],
                        [20.4550, 44.8050],
                        [20.4500, 44.8100],
                        [20.4450, 44.8150],
                        [20.4500, 44.8200]
                    ]]
                }
            },
            {
                "type": "Feature",
                "properties": {
                    "name": "Врачар",
                    "city": "Београд"
                },
                "geometry": {
                    "type": "Polygon",
                    "coordinates": [[
                        [20.4650, 44.8100],
                        [20.4700, 44.8150],
                        [20.4800, 44.8100],
                        [20.4900, 44.8050],
                        [20.4900, 44.7950],
                        [20.4800, 44.7900],
                        [20.4700, 44.7950],
                        [20.4650, 44.8000],
                        [20.4650, 44.8100]
                    ]]
                }
            }
        ]
    }
    
    # Создаем директорию если её нет
    os.makedirs(GEOJSON_DIR, exist_ok=True)
    
    filepath = os.path.join(GEOJSON_DIR, 'belgrade_districts_sample.geojson')
    with open(filepath, 'w', encoding='utf-8') as f:
        json.dump(sample_data, f, ensure_ascii=False, indent=2)
    
    logger.info(f"Создан пример файла: {filepath}")


def main():
    """Основная функция"""
    logger.info("Начало обработки GeoJSON файлов с границами районов")
    
    # Получаем данные из БД
    districts_db = get_districts_from_db()
    logger.info(f"Загружено {len(districts_db)} вариантов названий районов из БД")
    
    # Создаем пример файла если директория пуста
    if not os.path.exists(GEOJSON_DIR) or not os.listdir(GEOJSON_DIR):
        create_sample_belgrade_geojson()
    
    # Ищем все GeoJSON файлы
    geojson_files = []
    if os.path.exists(GEOJSON_DIR):
        geojson_files.extend(glob.glob(os.path.join(GEOJSON_DIR, '*.geojson')))
        geojson_files.extend(glob.glob(os.path.join(GEOJSON_DIR, '*.json')))
    
    # Добавляем тестовый файл из корня проекта
    test_file = '/data/hostel-booking-system/zemun_test.geojson'
    if os.path.exists(test_file):
        geojson_files.append(test_file)
    
    logger.info(f"Найдено {len(geojson_files)} GeoJSON файлов")
    
    # Обрабатываем все файлы
    all_matched = []
    for filepath in geojson_files:
        matched = process_geojson_file(filepath, districts_db)
        all_matched.extend(matched)
    
    logger.info(f"Всего сопоставлено {len(all_matched)} районов")
    
    if all_matched:
        # Генерируем SQL
        sql_content = generate_sql_updates(all_matched)
        
        # Сохраняем результаты
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        sql_filename = f'/data/hostel-booking-system/backend/migrations/109_update_district_boundaries_{timestamp}.up.sql'
        
        with open(sql_filename, 'w', encoding='utf-8') as f:
            f.write(sql_content)
        
        logger.info(f"SQL скрипт сохранен: {sql_filename}")
        
        # Сохраняем отчет
        report_filename = f'/data/hostel-booking-system/memory-bank/geojson_processing_report_{timestamp}.json'
        report_data = {
            'timestamp': datetime.now().isoformat(),
            'files_processed': geojson_files,
            'matched_districts': [
                {
                    'db_name': m['db_district']['district_name'],
                    'db_city': m['db_district']['city_name'],
                    'source_name': m['source_name']
                }
                for m in all_matched
            ],
            'total_matched': len(all_matched)
        }
        
        with open(report_filename, 'w', encoding='utf-8') as f:
            json.dump(report_data, f, ensure_ascii=False, indent=2)
        
        logger.info(f"Отчет сохранен: {report_filename}")
    else:
        logger.warning("Не найдено соответствий для обновления")
    
    logger.info("Работа завершена")


if __name__ == '__main__':
    main()