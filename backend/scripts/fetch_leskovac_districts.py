#!/usr/bin/env python3
"""
Скрипт для загрузки границ районов Лесковца и создания GeoJSON файлов
"""

import json
import requests
import os
import logging
from datetime import datetime

# Настройка логирования
logging.basicConfig(level=logging.INFO, format='%(asctime)s - %(levelname)s - %(message)s')
logger = logging.getLogger(__name__)

# Создаем директорию для геоданных
GEODATA_DIR = '/data/hostel-booking-system/geodata'
os.makedirs(GEODATA_DIR, exist_ok=True)

def create_leskovac_districts():
    """
    Создаем приблизительные границы районов Лесковца на основе их центральных точек
    Поскольку OSM может не иметь детальных границ городских районов малых городов,
    создадим реалистичные полигоны вокруг известных центров
    """
    
    # Районы Лесковца с их центрами из БД
    districts = [
        {
            'id': 'b4ba8cfc-d42e-45cb-8a40-d88f3d47b474',
            'name': 'Бубањ',
            'center': [21.9347, 42.9787]
        },
        {
            'id': '7cd05e68-8a97-4585-9753-49c555b7cef6',
            'name': 'Дубочица', 
            'center': [21.9347, 43.0087]
        },
        {
            'id': '72c6c046-daee-4867-afac-a917b5a0f60c',
            'name': 'Хисар',
            'center': [21.9547, 42.9887]
        },
        {
            'id': 'b5a0afd3-c687-4853-a7dc-dd79e67dd78a',
            'name': 'Центар',
            'center': [21.9447, 42.9987]
        }
    ]
    
    # Создаем GeoJSON коллекцию
    feature_collection = {
        "type": "FeatureCollection",
        "features": []
    }
    
    for district in districts:
        # Создаем реалистичный полигон для каждого района
        # Размер зависит от типа района
        if district['name'] == 'Центар':
            # Центральный район обычно меньше
            radius_lon = 0.008
            radius_lat = 0.006
        else:
            # Периферийные районы больше
            radius_lon = 0.012
            radius_lat = 0.010
        
        # Создаем полигон неправильной формы для реалистичности
        center_lon, center_lat = district['center']
        
        # Генерируем точки полигона (неправильная форма)
        points = []
        
        # Северная часть
        points.append([center_lon - radius_lon * 0.8, center_lat + radius_lat])
        points.append([center_lon - radius_lon * 0.4, center_lat + radius_lat * 1.1])
        points.append([center_lon, center_lat + radius_lat * 0.9])
        points.append([center_lon + radius_lon * 0.4, center_lat + radius_lat * 1.05])
        points.append([center_lon + radius_lon * 0.9, center_lat + radius_lat * 0.8])
        
        # Восточная часть
        points.append([center_lon + radius_lon * 1.1, center_lat + radius_lat * 0.4])
        points.append([center_lon + radius_lon * 0.95, center_lat])
        points.append([center_lon + radius_lon * 1.05, center_lat - radius_lat * 0.3])
        
        # Южная часть
        points.append([center_lon + radius_lon * 0.7, center_lat - radius_lat * 0.9])
        points.append([center_lon + radius_lon * 0.3, center_lat - radius_lat * 1.05])
        points.append([center_lon - radius_lon * 0.2, center_lat - radius_lat * 0.95])
        points.append([center_lon - radius_lon * 0.6, center_lat - radius_lat * 1.0])
        
        # Западная часть
        points.append([center_lon - radius_lon * 1.05, center_lat - radius_lat * 0.5])
        points.append([center_lon - radius_lon * 0.9, center_lat - radius_lat * 0.1])
        points.append([center_lon - radius_lon * 1.1, center_lat + radius_lat * 0.3])
        points.append([center_lon - radius_lon * 0.95, center_lat + radius_lat * 0.6])
        
        # Замыкаем полигон
        points.append(points[0])
        
        # Создаем feature
        feature = {
            "type": "Feature",
            "properties": {
                "name": district['name'],
                "id": district['id'],
                "city": "Лесковац",
                "created_at": datetime.now().isoformat(),
                "source": "approximated_from_center"
            },
            "geometry": {
                "type": "Polygon",
                "coordinates": [points]
            }
        }
        
        feature_collection['features'].append(feature)
        logger.info(f"Создан полигон для района {district['name']}")
    
    # Сохраняем в файл
    output_file = os.path.join(GEODATA_DIR, 'leskovac_districts.geojson')
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(feature_collection, f, ensure_ascii=False, indent=2)
    
    logger.info(f"Сохранен файл: {output_file}")
    return output_file


def try_fetch_from_osm():
    """
    Пытаемся получить границы из OSM (если они есть)
    """
    overpass_url = "http://overpass-api.de/api/interpreter"
    
    # Запрос для поиска административных границ в районе Лесковца
    query = """
    [out:json][timeout:25];
    // Ищем город Лесковац
    (
      node["name"="Лесковац"]["place"="city"];
      way["name"="Лесковац"]["place"="city"];
      relation["name"="Лесковац"]["place"="city"];
    );
    // Получаем область вокруг города
    (._;)->.city;
    // Ищем административные границы
    (
      relation["boundary"="administrative"]["admin_level"~"[8-10]"](around.city:5000);
      way["boundary"="administrative"](around.city:5000);
    );
    out geom;
    """
    
    try:
        logger.info("Попытка загрузки данных из OSM...")
        response = requests.post(overpass_url, data={'data': query}, timeout=30)
        
        if response.status_code == 200:
            data = response.json()
            elements = data.get('elements', [])
            logger.info(f"Получено {len(elements)} элементов из OSM")
            
            if elements:
                # Сохраняем сырые данные для анализа
                with open(os.path.join(GEODATA_DIR, 'leskovac_osm_raw.json'), 'w') as f:
                    json.dump(data, f, indent=2)
                logger.info("Сырые данные OSM сохранены для анализа")
            
            return elements
        else:
            logger.error(f"Ошибка OSM API: {response.status_code}")
            return []
            
    except Exception as e:
        logger.error(f"Ошибка при запросе к OSM: {e}")
        return []


def main():
    logger.info("Начало загрузки границ районов Лесковца")
    
    # Сначала пробуем OSM
    osm_data = try_fetch_from_osm()
    
    # В любом случае создаем приблизительные границы
    geojson_file = create_leskovac_districts()
    
    logger.info(f"""
    Завершено! Создан файл: {geojson_file}
    
    Следующие шаги:
    1. Запустить: python3 scripts/process_geojson_boundaries.py
    2. Проверить сгенерированный SQL файл
    3. Применить миграцию к базе данных
    """)


if __name__ == '__main__':
    main()