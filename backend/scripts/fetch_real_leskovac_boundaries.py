#!/usr/bin/env python3
"""
Скрипт для поиска реальных административных границ Лесковца из различных источников
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

def search_leskovac_in_nominatim():
    """Поиск Лесковца через Nominatim API"""
    logger.info("Поиск Лесковца через Nominatim...")
    
    params = {
        'q': 'Leskovac, Serbia',
        'format': 'json',
        'polygon_geojson': 1,
        'addressdetails': 1,
        'limit': 5
    }
    
    headers = {
        'User-Agent': 'HostelBookingSystem/1.0'
    }
    
    try:
        response = requests.get(NOMINATIM_URL, params=params, headers=headers, timeout=30)
        if response.status_code == 200:
            data = response.json()
            logger.info(f"Nominatim вернул {len(data)} результатов")
            
            for item in data:
                if item.get('type') in ['city', 'town', 'administrative']:
                    logger.info(f"Найден: {item.get('display_name')} (OSM ID: {item.get('osm_id')})")
                    if 'geojson' in item:
                        return item
            
        else:
            logger.error(f"Ошибка Nominatim: {response.status_code}")
    except Exception as e:
        logger.error(f"Ошибка при запросе к Nominatim: {e}")
    
    return None


def search_leskovac_districts_osm():
    """Поиск районов/месних заједница Лесковца в OSM"""
    logger.info("Поиск административных единиц Лесковца в OSM...")
    
    # Различные запросы для поиска
    queries = [
        # Поиск по месним заједницама (местным сообществам)
        """
        [out:json][timeout:30];
        area["name"~"Leskovac|Лесковац"]["admin_level"="8"]->.city;
        (
          node["place"~"suburb|neighbourhood|quarter"](area.city);
          way["place"~"suburb|neighbourhood|quarter"](area.city);
          relation["place"~"suburb|neighbourhood|quarter"](area.city);
        );
        out geom;
        """,
        
        # Поиск административных границ внутри bbox Лесковца
        """
        [out:json][timeout:30];
        (
          node["name"~"Leskovac|Лесковац"]["place"="city"];
        )->.city;
        (
          way(around.city:5000)["boundary"="administrative"]["admin_level"~"9|10"];
          relation(around.city:5000)["boundary"="administrative"]["admin_level"~"9|10"];
        );
        out geom;
        """,
        
        # Поиск месних заједница через другие теги
        """
        [out:json][timeout:30];
        (
          way(42.96,21.91,43.02,21.97)["admin_level"~"9|10"];
          relation(42.96,21.91,43.02,21.97)["admin_level"~"9|10"];
          way(42.96,21.91,43.02,21.97)["boundary"~"administrative|political"];
          relation(42.96,21.91,43.02,21.97)["boundary"~"administrative|political"];
        );
        out geom;
        """
    ]
    
    all_results = []
    
    for i, query in enumerate(queries):
        logger.info(f"Выполнение запроса {i+1}...")
        try:
            response = requests.post(OVERPASS_URL, data={'data': query}, timeout=60)
            if response.status_code == 200:
                data = response.json()
                elements = data.get('elements', [])
                logger.info(f"Запрос {i+1} вернул {len(elements)} элементов")
                all_results.extend(elements)
            else:
                logger.error(f"Ошибка в запросе {i+1}: {response.status_code}")
        except Exception as e:
            logger.error(f"Ошибка при выполнении запроса {i+1}: {e}")
        
        time.sleep(2)  # Пауза между запросами
    
    return all_results


def search_serbian_geodata():
    """Поиск в сербских источниках геоданных"""
    logger.info("Поиск в альтернативных источниках...")
    
    # Проверяем наличие загруженных файлов GADM или других источников
    potential_files = [
        '/data/hostel-booking-system/geodata/serbia_admin.geojson',
        '/data/hostel-booking-system/geodata/gadm36_SRB_2.json',
        '/data/hostel-booking-system/geodata/serbia_municipalities.geojson'
    ]
    
    for filepath in potential_files:
        if os.path.exists(filepath):
            logger.info(f"Найден файл: {filepath}")
            with open(filepath, 'r', encoding='utf-8') as f:
                data = json.load(f)
                return data
    
    return None


def analyze_leskovac_structure():
    """
    Анализ структуры Лесковца на основе известной информации
    
    Лесковац делится на месне заједнице (местные сообщества), 
    а не на классические городские районы как крупные города.
    """
    logger.info("""
    Известная информация о Лесковце:
    - Город в южной Сербии
    - Население около 60,000 человек
    - Административно делится на месне заједнице (местные сообщества)
    - Основные части города:
      * Центар (центр города)
      * Хисар (историческая крепость)
      * Дубочица (северная часть)
      * Бубањ (восточная часть)
    
    Эти районы могут не иметь официальных административных границ в OSM.
    """)
    
    # Создаем приблизительные, но более реалистичные границы
    # основанные на реальной географии Лесковца
    districts_info = {
        'Центар': {
            'description': 'Центральная часть города с главной площадью',
            'approximate_bounds': [[21.935, 42.993], [21.955, 43.003]]
        },
        'Хисар': {
            'description': 'Район вокруг средневековой крепости Хисар',
            'approximate_bounds': [[21.945, 42.983], [21.965, 42.993]]
        },
        'Дубочица': {
            'description': 'Северный район вдоль реки Ветерница',
            'approximate_bounds': [[21.925, 43.003], [21.945, 43.013]]
        },
        'Бубањ': {
            'description': 'Восточный район, место мемориала Бубањ',
            'approximate_bounds': [[21.925, 42.973], [21.945, 42.983]]
        }
    }
    
    return districts_info


def main():
    logger.info("=== Поиск реальных административных границ Лесковца ===")
    
    # 1. Поиск через Nominatim
    nominatim_result = search_leskovac_in_nominatim()
    if nominatim_result:
        with open('/data/hostel-booking-system/geodata/leskovac_nominatim.json', 'w', encoding='utf-8') as f:
            json.dump(nominatim_result, f, ensure_ascii=False, indent=2)
        logger.info("Данные Nominatim сохранены")
    
    # 2. Поиск районов в OSM
    osm_results = search_leskovac_districts_osm()
    if osm_results:
        with open('/data/hostel-booking-system/geodata/leskovac_osm_search.json', 'w', encoding='utf-8') as f:
            json.dump({'elements': osm_results}, f, ensure_ascii=False, indent=2)
        logger.info(f"Найдено {len(osm_results)} элементов в OSM")
        
        # Анализируем результаты
        for element in osm_results:
            tags = element.get('tags', {})
            name = tags.get('name', tags.get('name:sr', 'Без названия'))
            admin_level = tags.get('admin_level', 'N/A')
            boundary = tags.get('boundary', tags.get('place', 'N/A'))
            logger.info(f"  - {name} (admin_level: {admin_level}, type: {boundary})")
    
    # 3. Поиск в локальных файлах
    local_data = search_serbian_geodata()
    
    # 4. Анализ структуры города
    districts_info = analyze_leskovac_structure()
    
    logger.info("""
    
    === ВЫВОДЫ ===
    
    Лесковац, как малый город Сербии, вероятно не имеет официальных 
    административных границ районов в публичных геоданных.
    
    Варианты действий:
    1. Использовать границы месних заједница если они найдены в OSM
    2. Запросить официальные данные у города Лесковац
    3. Создать границы на основе местной географии и известных ориентиров
    
    Для получения РЕАЛЬНЫХ границ необходимо:
    - Обратиться в градску управу Лесковца
    - Или использовать официальные кадастровые данные РГЗ
    """)


if __name__ == '__main__':
    main()