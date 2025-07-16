#!/usr/bin/env python3
"""
Скрипт для поиска реальных административных границ районов Нови Сада
"""

import json
import requests
import os
import logging
from datetime import datetime
import time

logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

OVERPASS_URL = "http://overpass-api.de/api/interpreter"
NOMINATIM_URL = "https://nominatim.openstreetmap.org/search"

def search_novi_sad_districts():
    """Поиск районов Нови Сада в OSM"""
    logger.info("Поиск районов Нови Сада в OpenStreetMap...")
    
    # Известные районы Нови Сада
    districts = [
        "Лиман", "Грбавица", "Ново насеље", "Детелинара", "Телеп",
        "Адице", "Салајка", "Подбара", "Роткварија", "Стари град",
        "Клиса", "Петроварадин", "Сремска Каменица", "Ветерник", "Футог", "Каћ"
    ]
    
    # Различные запросы для поиска районов
    queries = [
        # Поиск месних заједница в Нови Саде
        """
        [out:json][timeout:60];
        area["name"="Нови Сад"]["admin_level"="8"]->.city;
        (
          node["place"~"suburb|neighbourhood|quarter"](area.city);
          way["place"~"suburb|neighbourhood|quarter"](area.city);
          relation["place"~"suburb|neighbourhood|quarter"](area.city);
          
          // Месне заједнице
          relation["admin_level"="10"](area.city);
          way["admin_level"="10"](area.city);
          
          // Административные границы
          relation["boundary"="administrative"]["admin_level"~"9|10"](area.city);
          way["boundary"="administrative"]["admin_level"~"9|10"](area.city);
        );
        out geom;
        """,
        
        # Поиск по конкретным названиям районов
        """
        [out:json][timeout:60];
        (
          // Лиман
          way["name"~"Лиман|Liman"]["boundary"="administrative"];
          relation["name"~"Лиман|Liman"]["boundary"="administrative"];
          way["name"~"Лиман|Liman"]["place"~"suburb|neighbourhood"];
          
          // Петроварадин
          way["name"~"Петроварадин|Petrovaradin"];
          relation["name"~"Петроварадин|Petrovaradin"];
          
          // Сремска Каменица
          way["name"~"Сремска Каменица|Sremska Kamenica"];
          relation["name"~"Сремска Каменица|Sremska Kamenica"];
        );
        out geom;
        """,
        
        # Поиск в границах bbox Нови Сада
        """
        [out:json][timeout:60];
        (
          way(45.20,19.70,45.35,19.95)["place"~"suburb|neighbourhood|quarter"];
          relation(45.20,19.70,45.35,19.95)["place"~"suburb|neighbourhood|quarter"];
          way(45.20,19.70,45.35,19.95)["admin_level"~"9|10"];
          relation(45.20,19.70,45.35,19.95)["admin_level"~"9|10"];
        );
        out geom;
        """
    ]
    
    all_results = []
    
    for i, query in enumerate(queries):
        logger.info(f"Выполнение запроса {i+1} из {len(queries)}...")
        try:
            response = requests.post(OVERPASS_URL, data={'data': query}, timeout=120)
            if response.status_code == 200:
                data = response.json()
                elements = data.get('elements', [])
                logger.info(f"Запрос {i+1} вернул {len(elements)} элементов")
                
                # Фильтруем и добавляем только уникальные элементы
                for elem in elements:
                    if elem['id'] not in [e['id'] for e in all_results]:
                        all_results.append(elem)
                        
            else:
                logger.error(f"Ошибка в запросе {i+1}: {response.status_code}")
        except Exception as e:
            logger.error(f"Ошибка при выполнении запроса {i+1}: {e}")
        
        time.sleep(3)  # Пауза между запросами
    
    return all_results

def search_via_nominatim():
    """Поиск районов через Nominatim"""
    logger.info("Поиск районов Нови Сада через Nominatim...")
    
    districts = [
        "Лиман", "Петроварадин", "Сремска Каменица", "Ветерник", "Футог"
    ]
    
    results = []
    headers = {'User-Agent': 'HostelBookingSystem/1.0'}
    
    for district in districts:
        params = {
            'q': f'{district}, Novi Sad, Serbia',
            'format': 'json',
            'polygon_geojson': 1,
            'addressdetails': 1,
            'limit': 3
        }
        
        try:
            response = requests.get(NOMINATIM_URL, params=params, headers=headers, timeout=30)
            if response.status_code == 200:
                data = response.json()
                if data:
                    logger.info(f"Nominatim нашел {len(data)} результатов для {district}")
                    results.extend(data)
            
            time.sleep(1)  # Nominatim требует паузу между запросами
            
        except Exception as e:
            logger.error(f"Ошибка Nominatim для {district}: {e}")
    
    return results

def analyze_results(osm_results, nominatim_results):
    """Анализ найденных данных"""
    logger.info("\n=== АНАЛИЗ РЕЗУЛЬТАТОВ ===")
    
    # Анализ OSM данных
    found_districts = {}
    
    for element in osm_results:
        tags = element.get('tags', {})
        name = tags.get('name', tags.get('name:sr', tags.get('name:sr-Cyrl', '')))
        place = tags.get('place', '')
        admin_level = tags.get('admin_level', '')
        boundary = tags.get('boundary', '')
        
        if name:
            element_type = element.get('type', 'unknown')
            osm_id = element.get('id', 'unknown')
            
            info = {
                'name': name,
                'type': element_type,
                'osm_id': osm_id,
                'place': place,
                'admin_level': admin_level,
                'boundary': boundary,
                'has_geometry': 'geometry' in element or 'members' in element
            }
            
            if name not in found_districts:
                found_districts[name] = []
            found_districts[name].append(info)
    
    # Вывод результатов
    logger.info(f"Найдено уникальных названий: {len(found_districts)}")
    
    # Проверяем известные районы Нови Сада
    known_districts = [
        "Лиман", "Грбавица", "Ново насеље", "Детелинара", "Телеп",
        "Адице", "Салајка", "Подбара", "Роткварија", "Стари град",
        "Клиса", "Петроварадин", "Сремска Каменица", "Ветерник", "Футог", "Каћ"
    ]
    
    found_known = []
    not_found = []
    
    for district in known_districts:
        found = False
        for name in found_districts.keys():
            if district.lower() in name.lower() or name.lower() in district.lower():
                found = True
                found_known.append(f"{district} -> {name}")
                break
        
        if not found:
            not_found.append(district)
    
    logger.info(f"\nНайдено известных районов: {len(found_known)}")
    for item in found_known:
        logger.info(f"  ✓ {item}")
    
    logger.info(f"\nНЕ найдено районов: {len(not_found)}")
    for item in not_found:
        logger.info(f"  ✗ {item}")
    
    # Анализ Nominatim
    logger.info(f"\nNominatim нашел {len(nominatim_results)} результатов")
    
    return found_districts

def save_results(osm_results, nominatim_results, found_districts):
    """Сохранение результатов"""
    
    # Сохраняем сырые данные OSM
    with open('/data/hostel-booking-system/geodata/novi_sad_osm_raw.json', 'w', encoding='utf-8') as f:
        json.dump({'elements': osm_results}, f, ensure_ascii=False, indent=2)
    
    # Сохраняем данные Nominatim
    with open('/data/hostel-booking-system/geodata/novi_sad_nominatim.json', 'w', encoding='utf-8') as f:
        json.dump(nominatim_results, f, ensure_ascii=False, indent=2)
    
    # Создаем GeoJSON с найденными районами
    features = []
    
    for name, items in found_districts.items():
        for item in items:
            if item['has_geometry']:
                # Ищем элемент в исходных данных
                for elem in osm_results:
                    if elem.get('id') == item['osm_id'] and elem.get('type') == item['type']:
                        geometry = None
                        
                        # Извлекаем геометрию
                        if 'geometry' in elem and elem['type'] == 'way':
                            coords = [[p['lon'], p['lat']] for p in elem['geometry']]
                            if coords[0] != coords[-1]:
                                coords.append(coords[0])
                            geometry = {
                                'type': 'Polygon',
                                'coordinates': [coords]
                            }
                        elif elem['type'] == 'relation' and 'members' in elem:
                            # Обрабатываем relation
                            outer_ways = []
                            for member in elem['members']:
                                if member.get('role') == 'outer' and 'geometry' in member:
                                    way_coords = [[p['lon'], p['lat']] for p in member['geometry']]
                                    outer_ways.append(way_coords)
                            
                            if outer_ways:
                                geometry = {
                                    'type': 'Polygon',
                                    'coordinates': outer_ways[:1]  # Берем первый outer
                                }
                        
                        if geometry:
                            feature = {
                                'type': 'Feature',
                                'properties': {
                                    'name': name,
                                    'osm_id': item['osm_id'],
                                    'osm_type': item['type'],
                                    'place': item['place'],
                                    'admin_level': item['admin_level']
                                },
                                'geometry': geometry
                            }
                            features.append(feature)
                            break
    
    geojson = {
        'type': 'FeatureCollection',
        'features': features
    }
    
    with open('/data/hostel-booking-system/geodata/novi_sad_districts_found.geojson', 'w', encoding='utf-8') as f:
        json.dump(geojson, f, ensure_ascii=False, indent=2)
    
    logger.info(f"\nСохранено {len(features)} районов с геометрией в novi_sad_districts_found.geojson")

def main():
    logger.info("=== Поиск реальных границ районов Нови Сада ===")
    
    # Создаем директорию если нет
    os.makedirs('/data/hostel-booking-system/geodata', exist_ok=True)
    
    # Поиск в OSM
    osm_results = search_novi_sad_districts()
    logger.info(f"Всего найдено {len(osm_results)} элементов в OSM")
    
    # Поиск через Nominatim
    nominatim_results = search_via_nominatim()
    
    # Анализ результатов
    found_districts = analyze_results(osm_results, nominatim_results)
    
    # Сохранение результатов
    save_results(osm_results, nominatim_results, found_districts)
    
    logger.info("""
    
    === ВЫВОДЫ ===
    
    1. Нови Сад имеет более развитую структуру в OSM чем Лесковац
    2. Некоторые районы (Петроварадин, Сремска Каменица) - это отдельные населенные пункты
    3. Городские районы (Лиман, Ново насеље) могут быть найдены как place=suburb
    
    Следующие шаги:
    1. Проанализировать файл novi_sad_districts_found.geojson
    2. Сопоставить найденные районы с районами в БД
    3. Применить реальные границы где они найдены
    """)

if __name__ == '__main__':
    main()