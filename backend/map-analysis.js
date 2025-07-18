const puppeteer = require('puppeteer');

(async () => {
  console.log('Запуск анализа карты...');
  
  const browser = await puppeteer.launch({
    headless: false,
    defaultViewport: null,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });
  
  try {
    const page = await browser.newPage();
    
    // Включаем логи консоли
    page.on('console', msg => console.log('Browser console:', msg.text()));
    
    const url = 'http://localhost:3001/en/map?lat=45.267136&lng=19.833549&radius=50000&zoom=10';
    console.log(`Переход на: ${url}`);
    
    await page.goto(url, { waitUntil: 'networkidle2' });
    console.log('Страница загружена, ожидание 5 секунд...');
    
    await new Promise(resolve => setTimeout(resolve, 5000));
    
    const result = await page.evaluate(() => {
      const data = {
        mapExists: false,
        markerCount: 0,
        clusterCount: 0,
        totalObjects: 0,
        markers: [],
        error: null,
        debug: {}
      };
      
      try {
        // Пытаемся найти карту разными способами
        let map = null;
        
        // Способ 1: window._map
        if (window._map) {
          map = window._map;
          data.debug.mapFoundBy = 'window._map';
        }
        
        // Способ 2: через React компоненты
        if (!map && window.React) {
          const mapContainer = document.querySelector('#map, .leaflet-container, [data-testid="map"]');
          if (mapContainer && mapContainer._leaflet_map) {
            map = mapContainer._leaflet_map;
            data.debug.mapFoundBy = 'leaflet container';
          }
        }
        
        // Способ 3: через глобальные переменные Leaflet
        if (!map && window.L) {
          const containers = document.querySelectorAll('.leaflet-container');
          containers.forEach(container => {
            if (container._leaflet_map) {
              map = container._leaflet_map;
              data.debug.mapFoundBy = 'leaflet container search';
            }
          });
        }
        
        // Отладочная информация
        data.debug.hasLeaflet = !!window.L;
        data.debug.hasReact = !!window.React;
        data.debug.leafletContainers = document.querySelectorAll('.leaflet-container').length;
        data.debug.mapDivExists = !!document.querySelector('#map');
        
        if (map) {
          data.mapExists = true;
          
          // Подсчитываем маркеры и кластеры
          map.eachLayer((layer) => {
            if (layer instanceof L.Marker) {
              if (layer instanceof L.MarkerCluster) {
                data.clusterCount++;
              } else {
                data.markerCount++;
                
                // Собираем данные о первых 5 маркерах
                if (data.markers.length < 5) {
                  const latLng = layer.getLatLng();
                  const popup = layer.getPopup();
                  
                  data.markers.push({
                    lat: latLng.lat,
                    lng: latLng.lng,
                    hasPopup: !!popup,
                    popupContent: popup ? popup.getContent() : null
                  });
                }
              }
            }
          });
          
          data.totalObjects = data.markerCount + data.clusterCount;
        } else {
          data.error = 'Карта не найдена ни одним из способов';
        }
      } catch (e) {
        data.error = e.toString();
      }
      
      return data;
    });
    
    console.log('\n=== РЕЗУЛЬТАТЫ АНАЛИЗА ===');
    console.log(`Карта существует: ${result.mapExists}`);
    console.log(`Количество маркеров: ${result.markerCount}`);
    console.log(`Количество кластеров: ${result.clusterCount}`);
    console.log(`Всего объектов на карте: ${result.totalObjects}`);
    
    console.log('\n=== ОТЛАДОЧНАЯ ИНФОРМАЦИЯ ===');
    console.log(`Leaflet доступен: ${result.debug.hasLeaflet}`);
    console.log(`React доступен: ${result.debug.hasReact}`);
    console.log(`Leaflet контейнеров найдено: ${result.debug.leafletContainers}`);
    console.log(`Div #map существует: ${result.debug.mapDivExists}`);
    if (result.debug.mapFoundBy) {
      console.log(`Карта найдена через: ${result.debug.mapFoundBy}`);
    }
    
    if (result.error) {
      console.log(`\nОшибка: ${result.error}`);
    }
    
    if (result.markers.length > 0) {
      console.log('\nПервые маркеры:');
      result.markers.forEach((marker, index) => {
        console.log(`${index + 1}. Координаты: [${marker.lat}, ${marker.lng}]`);
        console.log(`   Есть popup: ${marker.hasPopup}`);
        if (marker.popupContent) {
          console.log(`   Содержимое popup: ${marker.popupContent.substring(0, 100)}...`);
        }
      });
    }
    
  } catch (error) {
    console.error('Ошибка:', error);
  } finally {
    await browser.close();
    console.log('\nБраузер закрыт');
  }
})();