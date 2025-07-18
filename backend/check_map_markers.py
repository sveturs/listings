#!/usr/bin/env python3
import time
import json
from selenium import webdriver
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC

def check_map():
    # Настройка Chrome
    chrome_options = Options()
    chrome_options.add_argument('--headless')
    chrome_options.add_argument('--no-sandbox')
    chrome_options.add_argument('--disable-dev-shm-usage')
    chrome_options.add_argument('--window-size=1920,1080')
    
    driver = webdriver.Chrome(options=chrome_options)
    
    try:
        # Открываем страницу карты
        url = "http://localhost:3001/en/map?lat=45.267136&lng=19.833549&radius=50000&zoom=10"
        print(f"Открываю страницу: {url}")
        driver.get(url)
        
        # Ждем загрузки карты
        print("Жду загрузки карты...")
        time.sleep(5)
        
        # Делаем скриншот
        driver.save_screenshot("/data/hostel-booking-system/backend/map-check.png")
        print("Скриншот сохранен: map-check.png")
        
        # Проверяем наличие объекта window._map
        map_exists = driver.execute_script("return typeof window._map !== 'undefined';")
        print(f"Объект window._map существует: {map_exists}")
        
        if map_exists:
            # Получаем информацию о маркерах
            marker_info = driver.execute_script("""
                let markers = [];
                let clusters = [];
                
                if (window._map) {
                    window._map.eachLayer(function(layer) {
                        if (layer instanceof L.Marker && !(layer instanceof L.MarkerCluster)) {
                            markers.push({
                                lat: layer.getLatLng().lat,
                                lng: layer.getLatLng().lng
                            });
                        } else if (layer instanceof L.MarkerCluster) {
                            clusters.push({
                                childCount: layer._childCount,
                                bounds: layer.getBounds()
                            });
                        }
                    });
                }
                
                return {
                    markers: markers,
                    clusters: clusters,
                    totalMarkers: markers.length,
                    totalClusters: clusters.length
                };
            """)
            
            print(f"\nИнформация о маркерах:")
            print(f"Количество маркеров: {marker_info['totalMarkers']}")
            print(f"Количество кластеров: {marker_info['totalClusters']}")
            
            # Проверяем консоль на ошибки
            logs = driver.get_log('browser')
            errors = [log for log in logs if log['level'] == 'SEVERE']
            
            if errors:
                print(f"\nНайдено {len(errors)} ошибок в консоли:")
                for error in errors[:5]:  # Показываем первые 5 ошибок
                    print(f"- {error['message']}")
            else:
                print("\nОшибок в консоли не найдено")
                
            # Проверяем загрузку Leaflet
            leaflet_loaded = driver.execute_script("return typeof L !== 'undefined';")
            print(f"\nLeaflet загружен: {leaflet_loaded}")
            
            # Проверяем наличие MarkerCluster
            markercluster_loaded = driver.execute_script("return typeof L.MarkerClusterGroup !== 'undefined';")
            print(f"MarkerCluster загружен: {markercluster_loaded}")
            
        else:
            print("Объект window._map не найден!")
            
    except Exception as e:
        print(f"Произошла ошибка: {e}")
        
    finally:
        driver.quit()
        print("\nПроверка завершена")

if __name__ == "__main__":
    check_map()