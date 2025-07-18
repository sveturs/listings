const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();
  
  await page.goto('http://localhost:3001/en/map?lat=45.267136&lng=19.833549&radius=50000&zoom=10');
  
  // Ждем загрузки карты
  await page.waitForFunction(() => {
    return window.L && document.querySelector('.leaflet-container');
  }, { timeout: 10000 });
  
  await page.waitForTimeout(5000);
  
  const result = await page.evaluate(() => {
    // Проверяем разные способы доступа к карте
    let map = window._map || window.map;
    
    // Если не нашли, ищем через Leaflet
    if (!map && window.L) {
      const containers = document.querySelectorAll('.leaflet-container');
      if (containers.length > 0) {
        // Получаем карту через _leaflet_id
        const mapId = containers[0]._leaflet_id;
        if (mapId && window.L._maps) {
          map = window.L._maps[mapId];
        }
      }
    }
    
    if (!map) {
      return { 
        error: 'Map not found',
        hasLeaflet: !!window.L,
        hasContainer: !!document.querySelector('.leaflet-container'),
        globalKeys: Object.keys(window).filter(k => k.includes('map') || k.includes('Map'))
      };
    }
    
    let markerCount = 0;
    let clusterCount = 0;
    let markers = [];
    
    // Подсчет маркеров
    map.eachLayer((layer) => {
      if (layer instanceof L.Marker && !(layer instanceof L.MarkerCluster)) {
        markerCount++;
        const pos = layer.getLatLng();
        markers.push({ lat: pos.lat, lng: pos.lng });
      }
      if (layer instanceof L.MarkerCluster) {
        clusterCount++;
      }
    });
    
    return {
      mapExists: true,
      center: map.getCenter(),
      zoom: map.getZoom(),
      bounds: map.getBounds(),
      markerCount,
      clusterCount,
      totalLayers: Object.keys(map._layers).length,
      markers: markers.slice(0, 5) // Первые 5 маркеров
    };
  });
  
  console.log(JSON.stringify(result, null, 2));
  
  await browser.close();
})();