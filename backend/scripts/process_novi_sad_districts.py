#!/usr/bin/env python3
"""
Обработка найденных границ районов Нови Сада и создание SQL миграций
"""

import json
import logging
import psycopg2
from psycopg2.extras import RealDictCursor
import os
from datetime import datetime
from typing import Dict, List, Optional, Tuple

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

DATABASE_URL = os.getenv('DATABASE_URL', 'postgres://postgres:password@localhost:5432/hostel_db?sslmode=disable')

# Маппинг названий из OSM на названия в БД
DISTRICT_MAPPING = {
    # Прямые соответствия
    'Грбавица': ['Грбавица'],
    'Адице': ['Адице'],
    'Салајка': ['Салајка'],
    'Подбара': ['Подбара'], 
    'Роткварија': ['Роткварија'],
    'Стари град': ['Стари град', 'Стари Град'],
    'Клиса': ['Клиса'],
    'Сремска Каменица': ['Сремска Каменица', 'Sremska Kamenica'],
    'Ветерник': ['Ветерник', 'Veternik', 'Ветерник - Футог'],
    'Футог': ['Футог', 'Futog', 'Нови Сад - Футог'],
    'Петроварадин': ['Петроварадин', 'Petrovaradin', 'Петроварадинска'],
    
    # Сложные соответствия
    'Лиман': ['Лиман', 'Лиманец', 'Liman I', 'Liman II', 'Liman III', 'Liman IV'],
    'Ново насеље': ['Ново насеље', 'Novo Naselje', 'Ново Насеље'],
    'Детелинара': ['Детелинара', 'Detelinara', 'Петроварадин - Центар - Детелинара'],
    'Телеп': ['Телеп', 'Telep', 'Телепска долина'],
    'Каћ': ['Каћ', 'Kać', 'Ченеј - Каћ']
}

def get_db_connection():
    """Получить соединение с БД"""
    return psycopg2.connect(DATABASE_URL)

def load_districts_from_db():
    """Загрузить районы Нови Сада из БД"""
    conn = get_db_connection()
    try:
        with conn.cursor(cursor_factory=RealDictCursor) as cur:
            cur.execute("""
                SELECT d.id, d.name as district_name
                FROM districts d
                JOIN cities c ON d.city_id = c.id
                WHERE c.name = 'Нови Сад'
                ORDER BY d.name
            """)
            return {row['district_name']: row['id'] for row in cur.fetchall()}
    finally:
        conn.close()

def load_found_districts():
    """Загрузить найденные границы районов"""
    with open('/data/hostel-booking-system/geodata/novi_sad_districts_found.geojson', 'r', encoding='utf-8') as f:
        return json.load(f)

def match_district(osm_name: str, db_districts: Dict[str, str]) -> Optional[Tuple[str, str]]:
    """Сопоставить название из OSM с районом в БД"""
    
    # Прямое совпадение
    if osm_name in db_districts:
        return (osm_name, db_districts[osm_name])
    
    # Поиск через маппинг
    for db_name, osm_variants in DISTRICT_MAPPING.items():
        if db_name in db_districts:
            for variant in osm_variants:
                if variant.lower() in osm_name.lower() or osm_name.lower() in variant.lower():
                    return (db_name, db_districts[db_name])
    
    return None

def process_geometry(geometry: Dict) -> Dict:
    """Обработать и валидировать геометрию"""
    if geometry['type'] != 'Polygon':
        logger.warning(f"Неподдерживаемый тип геометрии: {geometry['type']}")
        return None
    
    # Проверяем что полигон замкнут
    coords = geometry['coordinates'][0]
    if coords[0] != coords[-1]:
        coords.append(coords[0])
    
    return {
        'type': 'Polygon',
        'coordinates': [coords]
    }

def select_best_geometry(features: List[Dict], district_name: str) -> Optional[Dict]:
    """Выбрать лучшую геометрию из нескольких вариантов"""
    
    # Приоритеты: 
    # 1. admin_level=10
    # 2. place=suburb
    # 3. Самый большой полигон
    
    candidates = []
    
    for feature in features:
        props = feature['properties']
        geometry = feature.get('geometry')
        
        if not geometry:
            continue
            
        # Проверяем что название подходит
        name = props.get('name', '')
        matched = False
        
        if district_name in DISTRICT_MAPPING:
            for variant in DISTRICT_MAPPING[district_name]:
                if variant.lower() in name.lower() or name.lower() in variant.lower():
                    matched = True
                    break
        
        if matched:
            # Считаем приоритет
            priority = 0
            if props.get('admin_level') == '10':
                priority += 100
            if props.get('place') == 'suburb':
                priority += 50
            
            candidates.append({
                'feature': feature,
                'priority': priority,
                'name': name
            })
    
    if not candidates:
        return None
    
    # Сортируем по приоритету
    candidates.sort(key=lambda x: x['priority'], reverse=True)
    
    # Берем с наивысшим приоритетом
    best = candidates[0]
    logger.info(f"Выбран вариант для {district_name}: {best['name']} (приоритет: {best['priority']})")
    
    return best['feature']['geometry']

def create_migration_sql(matched_districts: List[Tuple[str, str, Dict]]) -> str:
    """Создать SQL миграцию"""
    
    sql_lines = [
        "-- Обновление реальных границ районов Нови Сада",
        f"-- Сгенерировано: {datetime.now().isoformat()}",
        "-- Источник: OpenStreetMap",
        "",
        "BEGIN;",
        "",
        "-- Создаем резервную копию",
        f"CREATE TABLE IF NOT EXISTS districts_novi_sad_backup_{datetime.now().strftime('%Y%m%d')} AS",
        "SELECT d.* FROM districts d",
        "JOIN cities c ON d.city_id = c.id", 
        "WHERE c.name = 'Нови Сад';",
        ""
    ]
    
    update_count = 0
    
    for district_name, district_id, geometry in matched_districts:
        geometry_json = json.dumps(geometry)
        
        sql_lines.append(f"""
-- Район: {district_name}
UPDATE districts 
SET boundary = ST_GeomFromGeoJSON('{geometry_json}'),
    area_km2 = ST_Area(ST_GeomFromGeoJSON('{geometry_json}')::geography) / 1000000,
    updated_at = NOW()
WHERE id = '{district_id}';
""")
        update_count += 1
    
    # Обновляем центральные точки
    sql_lines.append("""
-- Обновляем центральные точки
UPDATE districts 
SET center_point = ST_Centroid(boundary)
WHERE city_id IN (SELECT id FROM cities WHERE name = 'Нови Сад')
  AND boundary IS NOT NULL;
""")
    
    # Добавляем проверки
    sql_lines.extend([
        "",
        "-- Проверка результатов",
        "SELECT d.name as district,",
        "       ST_IsValid(d.boundary) as is_valid,",
        "       d.area_km2,",
        "       ST_NPoints(d.boundary) as num_points,",
        "       ST_Contains(d.boundary, d.center_point) as contains_center",
        "FROM districts d",
        "JOIN cities c ON d.city_id = c.id",
        "WHERE c.name = 'Нови Сад'",
        "ORDER BY d.name;",
        "",
        f"-- Обновлено районов: {update_count}",
        "",
        "COMMIT;"
    ])
    
    return '\n'.join(sql_lines)

def main():
    logger.info("=== Обработка границ районов Нови Сада ===")
    
    # Загружаем данные
    db_districts = load_districts_from_db()
    logger.info(f"Загружено {len(db_districts)} районов из БД")
    
    geojson_data = load_found_districts()
    features = geojson_data['features']
    logger.info(f"Загружено {len(features)} объектов из GeoJSON")
    
    # Группируем features по районам
    features_by_district = {}
    for db_name in db_districts.keys():
        features_by_district[db_name] = []
    
    # Сопоставляем
    for feature in features:
        osm_name = feature['properties'].get('name', '')
        match = match_district(osm_name, db_districts)
        
        if match:
            db_name, district_id = match
            features_by_district[db_name].append(feature)
    
    # Выбираем лучшую геометрию для каждого района
    matched_districts = []
    
    for district_name, district_id in db_districts.items():
        candidates = features_by_district[district_name]
        
        if candidates:
            best_geometry = select_best_geometry(candidates, district_name)
            
            if best_geometry:
                processed_geometry = process_geometry(best_geometry)
                if processed_geometry:
                    matched_districts.append((district_name, district_id, processed_geometry))
                    logger.info(f"✓ Найдена геометрия для {district_name}")
            else:
                logger.warning(f"✗ Не найдена подходящая геометрия для {district_name}")
        else:
            logger.warning(f"✗ Нет кандидатов для {district_name}")
    
    logger.info(f"\nСопоставлено {len(matched_districts)} из {len(db_districts)} районов")
    
    # Создаем SQL миграцию
    if matched_districts:
        sql_content = create_migration_sql(matched_districts)
        
        timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
        sql_filename = f'/data/hostel-booking-system/backend/migrations/110_update_novi_sad_districts_{timestamp}.up.sql'
        
        with open(sql_filename, 'w', encoding='utf-8') as f:
            f.write(sql_content)
        
        logger.info(f"\nСоздан SQL файл: {sql_filename}")
        
        # Создаем также GeoJSON с финальными границами
        final_features = []
        for district_name, district_id, geometry in matched_districts:
            feature = {
                'type': 'Feature',
                'properties': {
                    'name': district_name,
                    'district_id': district_id,
                    'city': 'Нови Сад'
                },
                'geometry': geometry
            }
            final_features.append(feature)
        
        final_geojson = {
            'type': 'FeatureCollection',
            'features': final_features
        }
        
        with open('/data/hostel-booking-system/geodata/novi_sad_districts_final.geojson', 'w', encoding='utf-8') as f:
            json.dump(final_geojson, f, ensure_ascii=False, indent=2)
        
        logger.info("Создан финальный GeoJSON: novi_sad_districts_final.geojson")
    
    logger.info("""
    
    === ГОТОВО ===
    
    Следующие шаги:
    1. Проверить SQL миграцию
    2. Применить миграцию к БД
    3. Проверить результаты на карте
    """)

if __name__ == '__main__':
    main()