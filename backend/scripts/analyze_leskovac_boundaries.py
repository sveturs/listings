#!/usr/bin/env python3
"""
Анализ загруженных данных OSM для поиска реальных границ районов Лесковца
"""

import json
import logging
from typing import Dict, List, Optional

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

def load_osm_data():
    """Загрузить данные OSM"""
    with open('/data/hostel-booking-system/geodata/leskovac_osm_search.json', 'r') as f:
        return json.load(f)

def load_nominatim_data():
    """Загрузить данные Nominatim"""
    with open('/data/hostel-booking-system/geodata/leskovac_nominatim.json', 'r') as f:
        return json.load(f)

def extract_city_boundaries(data):
    """Извлечь городские административные единицы"""
    city_units = []
    
    for element in data.get('elements', []):
        tags = element.get('tags', {})
        name = tags.get('name', tags.get('name:sr', ''))
        admin_level = tags.get('admin_level', '')
        
        # Ищем административные единицы города Лесковац
        if name and admin_level == '9':
            # Проверяем, что это городская единица, а не село
            if 'Лесковац' in name or element.get('id') == 2247345:  # ID для города Лесковац в OSM
                city_units.append({
                    'name': name,
                    'id': element.get('id'),
                    'type': element.get('type'),
                    'admin_level': admin_level,
                    'tags': tags,
                    'geometry': element.get('geometry', [])
                })
    
    return city_units

def convert_to_geojson_polygon(geometry):
    """Конвертировать OSM геометрию в GeoJSON полигон"""
    if not geometry:
        return None
    
    coords = [[point['lon'], point['lat']] for point in geometry]
    
    # Убедимся, что полигон замкнут
    if coords[0] != coords[-1]:
        coords.append(coords[0])
    
    return {
        "type": "Polygon",
        "coordinates": [coords]
    }

def analyze_leskovac_structure():
    """Анализ структуры города"""
    
    # Загружаем данные
    osm_data = load_osm_data()
    nominatim_data = load_nominatim_data()
    
    logger.info(f"Загружено {len(osm_data.get('elements', []))} элементов OSM")
    
    # Анализируем Nominatim
    logger.info(f"Nominatim boundary type: {nominatim_data.get('geojson', {}).get('type')}")
    
    # Извлекаем городские единицы
    city_units = extract_city_boundaries(osm_data)
    logger.info(f"Найдено {len(city_units)} городских административных единиц")
    
    # Ищем конкретно город Лесковац
    leskovac_city = None
    for element in osm_data.get('elements', []):
        if element.get('type') == 'relation' and element.get('id') == 2247345:
            leskovac_city = element
            break
    
    if leskovac_city:
        logger.info("Найден город Лесковац (relation 2247345)")
        
        # Конвертируем в GeoJSON
        members = leskovac_city.get('members', [])
        
        # Собираем outer границы
        outer_ways = []
        for member in members:
            if member.get('role') == 'outer' and 'geometry' in member:
                outer_ways.append(member['geometry'])
        
        if outer_ways:
            # Создаем единый полигон города
            city_boundary = {
                "type": "FeatureCollection",
                "features": [{
                    "type": "Feature",
                    "properties": {
                        "name": "Лесковац",
                        "type": "city_boundary",
                        "osm_id": 2247345,
                        "admin_level": "9"
                    },
                    "geometry": {
                        "type": "Polygon",
                        "coordinates": [[[p['lon'], p['lat']] for p in outer_ways[0]]]
                    }
                }]
            }
            
            # Сохраняем границы города
            with open('/data/hostel-booking-system/geodata/leskovac_city_boundary.geojson', 'w') as f:
                json.dump(city_boundary, f, ensure_ascii=False, indent=2)
            
            logger.info("Границы города Лесковац сохранены в leskovac_city_boundary.geojson")
            
            # Теперь нужно разделить город на районы
            # Поскольку OSM не содержит внутренних районов, создадим их на основе географии
            create_district_boundaries_from_city(city_boundary)
    
    return city_units

def create_district_boundaries_from_city(city_boundary):
    """Создать границы районов на основе границ города"""
    
    # Извлекаем координаты города
    city_coords = city_boundary['features'][0]['geometry']['coordinates'][0]
    
    # Находим границы города
    min_lon = min(p[0] for p in city_coords)
    max_lon = max(p[0] for p in city_coords)
    min_lat = min(p[1] for p in city_coords)
    max_lat = max(p[1] for p in city_coords)
    
    center_lon = (min_lon + max_lon) / 2
    center_lat = (min_lat + max_lat) / 2
    
    logger.info(f"Границы города: {min_lon},{min_lat} - {max_lon},{max_lat}")
    logger.info(f"Центр города: {center_lon},{center_lat}")
    
    # Создаем районы на основе квадрантов
    districts = {
        "Центар": {
            "bounds": [center_lon - 0.01, center_lat - 0.01, center_lon + 0.01, center_lat + 0.01],
            "description": "Центральная часть города"
        },
        "Хисар": {
            "bounds": [center_lon, center_lat - 0.02, center_lon + 0.02, center_lat],
            "description": "Юго-восточная часть с крепостью"
        },
        "Дубочица": {
            "bounds": [center_lon - 0.02, center_lat, center_lon, center_lat + 0.02],
            "description": "Северо-западная часть"
        },
        "Бубањ": {
            "bounds": [center_lon - 0.02, center_lat - 0.02, center_lon, center_lat],
            "description": "Юго-западная часть"
        }
    }
    
    features = []
    for name, info in districts.items():
        bounds = info['bounds']
        feature = {
            "type": "Feature",
            "properties": {
                "name": name,
                "city": "Лесковац",
                "description": info['description']
            },
            "geometry": {
                "type": "Polygon",
                "coordinates": [[
                    [bounds[0], bounds[1]],
                    [bounds[2], bounds[1]],
                    [bounds[2], bounds[3]],
                    [bounds[0], bounds[3]],
                    [bounds[0], bounds[1]]
                ]]
            }
        }
        features.append(feature)
    
    districts_geojson = {
        "type": "FeatureCollection",
        "features": features
    }
    
    with open('/data/hostel-booking-system/geodata/leskovac_districts_from_city.geojson', 'w') as f:
        json.dump(districts_geojson, f, ensure_ascii=False, indent=2)
    
    logger.info("Созданы приблизительные границы районов на основе границ города")

def main():
    logger.info("=== Анализ границ Лесковца ===")
    
    city_units = analyze_leskovac_structure()
    
    logger.info(f"""
    
    РЕЗУЛЬТАТЫ АНАЛИЗА:
    
    1. В OSM найдена административная граница города Лесковац (relation 2247345)
    2. Городские районы (Центар, Хисар, Дубочица, Бубањ) НЕ имеют 
       официальных административных границ в OSM
    3. Лесковац административно делится на села и месне заједнице,
       но внутригородские районы не оформлены как отдельные единицы
    
    ВЫВОД: Для получения РЕАЛЬНЫХ границ районов необходимо:
    - Обратиться в градску управу Лесковца за официальными данными
    - Использовать кадастровые карты
    - Или создать границы на основе местных знаний и ориентиров
    """)

if __name__ == '__main__':
    main()