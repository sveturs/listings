#!/usr/bin/env python3
"""
Скрипт для загрузки административных границ районов Сербии из OpenStreetMap
и подготовки SQL для обновления базы данных.
"""

import json
import requests
import psycopg2
from psycopg2.extras import RealDictCursor
import os
from typing import Dict, List, Tuple, Optional
import logging
from datetime import datetime
import time

# Настройка логирования
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

# Конфигурация
DATABASE_URL = os.getenv('DATABASE_URL', 'postgres://postgres:password@localhost:5432/hostel_db?sslmode=disable')
OVERPASS_URL = "http://overpass-api.de/api/interpreter"

# Маппинг городов на их OSM relation ID или координаты
CITY_MAPPING = {
    'Београд': {
        'osm_id': 1677007,  # OSM relation ID для Белграда
        'name_variants': ['Belgrade', 'Beograd', 'Београд'],
        'admin_level': 10  # Уровень для городских районов
    },
    'Нови Сад': {
        'osm_id': 2614767,
        'name_variants': ['Novi Sad', 'Нови Сад'],
        'admin_level': 10
    },
    'Ниш': {
        'osm_id': 2614144,
        'name_variants': ['Niš', 'Nis', 'Ниш'],
        'admin_level': 10
    },
    'Крагујевац': {
        'osm_id': 2716291,
        'name_variants': ['Kragujevac', 'Крагујевац'],
        'admin_level': 10
    },
    'Суботица': {
        'osm_id': 2614405,
        'name_variants': ['Subotica', 'Суботица'],
        'admin_level': 10
    }
}

# Маппинг названий районов (OSM -> БД)
DISTRICT_NAME_MAPPING = {
    # Белград
    'Земун': ['Zemun', 'Земун'],
    'Стари Град': ['Stari Grad', 'Стари град', 'Old Town'],
    'Врачар': ['Vračar', 'Врачар'],
    'Савски венац': ['Savski Venac', 'Савски венац'],
    'Палилула': ['Palilula', 'Палилула'],
    'Нови Београд': ['Novi Beograd', 'New Belgrade', 'Нови Београд'],
    'Чукарица': ['Čukarica', 'Чукарица'],
    'Звездара': ['Zvezdara', 'Звездара'],
    'Раковица': ['Rakovica', 'Раковица'],
    'Вождовац': ['Voždovac', 'Вождовац'],
    # Нови Сад
    'Лиман': ['Liman', 'Лиман'],
    'Грбавица': ['Grbavica', 'Грбавица'],
    'Ново насеље': ['Novo Naselje', 'Ново насеље'],
    'Детелинара': ['Detelinara', 'Детелинара'],
    'Телеп': ['Telep', 'Телеп'],
    'Адице': ['Adice', 'Адице'],
    'Салајка': ['Salajka', 'Салајка'],
}


def get_db_connection():
    """Получить соединение с базой данных"""
    return psycopg2.connect(DATABASE_URL)


def get_districts_from_db() -> Dict[str, List[Dict]]:
    """Получить список районов из базы данных сгруппированных по городам"""
    conn = get_db_connection()
    try:
        with conn.cursor(cursor_factory=RealDictCursor) as cur:
            cur.execute("""
                SELECT d.id, d.name as district_name, d.center_point,
                       c.id as city_id, c.name as city_name,
                       ST_X(d.center_point) as lon, ST_Y(d.center_point) as lat
                FROM districts d
                JOIN cities c ON d.city_id = c.id
                ORDER BY c.name, d.name
            """)
            
            districts_by_city = {}
            for row in cur.fetchall():
                city_name = row['city_name']
                if city_name not in districts_by_city:
                    districts_by_city[city_name] = []
                districts_by_city[city_name].append(dict(row))
            
            return districts_by_city
    finally:
        conn.close()


def fetch_osm_boundaries(city_name: str, city_config: Dict) -> List[Dict]:
    """Загрузить границы районов города из OpenStreetMap"""
    logger.info(f"Загрузка границ для города {city_name}")
    
    # Overpass QL запрос для получения административных границ
    query = f"""
    [out:json][timeout:60];
    (
        // Получаем все административные единицы внутри города
        area({city_config['osm_id'] + 3600000000})->.city;
        (
            relation["boundary"="administrative"]["admin_level"="{city_config['admin_level']}"](area.city);
            relation["boundary"="administrative"]["admin_level"="9"](area.city);
            relation["boundary"="administrative"]["admin_level"="8"](area.city);
        );
    );
    out geom;
    """
    
    try:
        response = requests.post(OVERPASS_URL, data={'data': query}, timeout=120)
        response.raise_for_status()
        data = response.json()
        
        districts = []
        for element in data.get('elements', []):
            if element['type'] == 'relation':
                tags = element.get('tags', {})
                name = tags.get('name', '') or tags.get('name:sr', '') or tags.get('name:sr-Cyrl', '')
                
                if name:
                    # Конвертируем геометрию в GeoJSON
                    geometry = convert_osm_to_geojson(element)
                    if geometry:
                        districts.append({
                            'name': name,
                            'osm_id': element['id'],
                            'geometry': geometry,
                            'tags': tags
                        })
        
        logger.info(f"Найдено {len(districts)} районов для города {city_name}")
        return districts
        
    except Exception as e:
        logger.error(f"Ошибка при загрузке данных для {city_name}: {e}")
        return []


def convert_osm_to_geojson(osm_element: Dict) -> Optional[Dict]:
    """Конвертировать OSM элемент в GeoJSON геометрию"""
    if 'members' not in osm_element:
        return None
    
    # Собираем координаты из way членов
    coordinates = []
    
    for member in osm_element['members']:
        if member['type'] == 'way' and member['role'] == 'outer':
            if 'geometry' in member:
                way_coords = [[node['lon'], node['lat']] for node in member['geometry']]
                if way_coords:
                    coordinates.append(way_coords)
    
    if coordinates:
        # Для простоты берем только первый внешний контур
        return {
            'type': 'Polygon',
            'coordinates': coordinates[0:1]  # Берем только первый полигон
        }
    
    return None


def match_district_names(osm_name: str, db_districts: List[Dict]) -> Optional[Dict]:
    """Сопоставить название района из OSM с районом в БД"""
    osm_name_lower = osm_name.lower().strip()
    
    # Проверяем прямое совпадение
    for district in db_districts:
        db_name_lower = district['district_name'].lower().strip()
        if osm_name_lower == db_name_lower:
            return district
    
    # Проверяем через маппинг
    for db_name, variants in DISTRICT_NAME_MAPPING.items():
        if osm_name in variants or osm_name_lower in [v.lower() for v in variants]:
            for district in db_districts:
                if district['district_name'] == db_name:
                    return district
    
    # Частичное совпадение
    for district in db_districts:
        db_name_lower = district['district_name'].lower().strip()
        if osm_name_lower in db_name_lower or db_name_lower in osm_name_lower:
            return district
    
    return None


def generate_sql_updates(districts_data: Dict[str, List[Tuple[Dict, Dict]]]) -> str:
    """Генерировать SQL для обновления границ районов"""
    sql_lines = [
        "-- SQL скрипт для обновления границ районов Сербии",
        f"-- Сгенерировано: {datetime.now().isoformat()}",
        "-- Источник: OpenStreetMap",
        "",
        "BEGIN;",
        "",
        "-- Создаем резервную копию текущих границ",
        "CREATE TABLE IF NOT EXISTS districts_boundaries_backup AS ",
        "SELECT id, name, boundary, updated_at FROM districts;",
        ""
    ]
    
    update_count = 0
    for city_name, matches in districts_data.items():
        if matches:
            sql_lines.append(f"-- Обновления для города {city_name}")
            
            for osm_district, db_district in matches:
                district_id = db_district['id']
                district_name = db_district['district_name']
                geometry_json = json.dumps(osm_district['geometry'])
                
                sql_lines.append(f"""
-- Район: {district_name} (OSM ID: {osm_district['osm_id']})
UPDATE districts 
SET boundary = ST_GeomFromGeoJSON('{geometry_json}'),
    updated_at = NOW()
WHERE id = '{district_id}';
""")
                update_count += 1
            
            sql_lines.append("")
    
    sql_lines.extend([
        f"-- Всего обновлений: {update_count}",
        "",
        "-- Проверка результатов",
        "SELECT c.name as city, d.name as district, ",
        "       ST_IsValid(d.boundary) as is_valid,",
        "       ST_Area(d.boundary::geography) / 1000000 as area_km2",
        "FROM districts d",
        "JOIN cities c ON d.city_id = c.id",
        "WHERE d.boundary IS NOT NULL",
        "ORDER BY c.name, d.name;",
        "",
        "COMMIT;"
    ])
    
    return '\n'.join(sql_lines)


def main():
    """Основная функция"""
    logger.info("Начало работы скрипта загрузки границ районов")
    
    # Получаем данные из БД
    districts_by_city = get_districts_from_db()
    logger.info(f"Загружено {sum(len(d) for d in districts_by_city.values())} районов из БД")
    
    # Результаты сопоставления
    matched_districts = {}
    
    # Обрабатываем каждый город
    for city_name, city_config in CITY_MAPPING.items():
        if city_name in districts_by_city:
            # Загружаем границы из OSM
            osm_districts = fetch_osm_boundaries(city_name, city_config)
            
            if osm_districts:
                matched = []
                db_districts = districts_by_city[city_name]
                
                # Сопоставляем районы
                for osm_district in osm_districts:
                    db_match = match_district_names(osm_district['name'], db_districts)
                    if db_match:
                        matched.append((osm_district, db_match))
                        logger.info(f"Сопоставлен: {osm_district['name']} -> {db_match['district_name']}")
                    else:
                        logger.warning(f"Не найдено соответствие для: {osm_district['name']}")
                
                matched_districts[city_name] = matched
                logger.info(f"Сопоставлено {len(matched)} из {len(db_districts)} районов для {city_name}")
            
            # Задержка между запросами к OSM API
            time.sleep(2)
    
    # Генерируем SQL
    sql_content = generate_sql_updates(matched_districts)
    
    # Сохраняем результаты
    timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
    sql_filename = f'/data/hostel-booking-system/backend/migrations/109_update_district_boundaries_{timestamp}.up.sql'
    
    with open(sql_filename, 'w', encoding='utf-8') as f:
        f.write(sql_content)
    
    logger.info(f"SQL скрипт сохранен: {sql_filename}")
    
    # Сохраняем детальный отчет
    report_filename = f'/data/hostel-booking-system/memory-bank/district_boundaries_report_{timestamp}.json'
    report_data = {
        'timestamp': datetime.now().isoformat(),
        'matched_districts': {
            city: [(osm['name'], db['district_name']) for osm, db in matches]
            for city, matches in matched_districts.items()
        },
        'statistics': {
            'total_db_districts': sum(len(d) for d in districts_by_city.values()),
            'total_matched': sum(len(m) for m in matched_districts.values()),
            'cities_processed': len(matched_districts)
        }
    }
    
    with open(report_filename, 'w', encoding='utf-8') as f:
        json.dump(report_data, f, ensure_ascii=False, indent=2)
    
    logger.info(f"Отчет сохранен: {report_filename}")
    logger.info("Работа скрипта завершена")


if __name__ == '__main__':
    main()