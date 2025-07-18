const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: true });
  const page = await browser.newPage();
  
  try {
    console.log('Открываю страницу карты...');
    await page.goto('http://localhost:3001/en/map?lat=45.267136&lng=19.833549&radius=50000&zoom=10', {
      waitUntil: 'domcontentloaded',
      timeout: 60000
    });
    
    console.log('Жду загрузки карты...');
    await page.waitForTimeout(5000);
    
    console.log('Делаю скриншот...');
    await page.screenshot({ path: 'map-check.png', fullPage: true });
    
    console.log('Проверяю консоль на ошибки...');
    page.on('console', msg => {
      if (msg.type() === 'error') {
        console.log('Console error:', msg.text());
      }
    });
    
    console.log('Проверяю маркеры...');
    const mapData = await page.evaluate(() => {
      const result = {
        mapExists: false,
        markersCount: 0,
        clustersCount: 0,
        totalLayers: 0,
        mapDetails: null
      };
      
      try {
        if (window._map) {
          result.mapExists = true;
          let markerCount = 0;
          let clusterCount = 0;
          let totalLayers = 0;
          
          window._map.eachLayer((layer) => {
            totalLayers++;
            if (layer.options && layer.options.icon) {
              if (layer._childCount) {
                clusterCount++;
              } else {
                markerCount++;
              }
            }
          });
          
          result.markersCount = markerCount;
          result.clustersCount = clusterCount;
          result.totalLayers = totalLayers;
          result.mapDetails = {
            center: window._map.getCenter(),
            zoom: window._map.getZoom(),
            bounds: window._map.getBounds()
          };
        }
      } catch (e) {
        result.error = e.toString();
      }
      
      return result;
    });
    
    console.log('\n=== ОТЧЕТ О ПРОВЕРКЕ КАРТЫ ===');
    console.log('Карта существует:', mapData.mapExists);
    console.log('Количество маркеров:', mapData.markersCount);
    console.log('Количество кластеров:', mapData.clustersCount);
    console.log('Всего слоев на карте:', mapData.totalLayers);
    if (mapData.mapDetails) {
      console.log('Центр карты:', mapData.mapDetails.center);
      console.log('Уровень зума:', mapData.mapDetails.zoom);
    }
    if (mapData.error) {
      console.log('Ошибка:', mapData.error);
    }
    console.log('Скриншот сохранен в: map-check.png');
    
  } catch (error) {
    console.error('Ошибка выполнения:', error);
  } finally {
    await browser.close();
  }
})();