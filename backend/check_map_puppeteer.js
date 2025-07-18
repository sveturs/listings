const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({ 
    headless: false,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });
  const page = await browser.newPage();
  
  // Collect console messages
  const consoleLogs = [];
  const networkErrors = [];
  
  page.on('console', msg => {
    consoleLogs.push({
      type: msg.type(),
      text: msg.text(),
      location: msg.location()
    });
  });
  
  page.on('pageerror', error => {
    consoleLogs.push({
      type: 'pageerror',
      text: error.message
    });
  });
  
  page.on('requestfailed', request => {
    networkErrors.push({
      url: request.url(),
      failure: request.failure()
    });
  });
  
  console.log('Navigating to map...');
  await page.goto('http://localhost:3001/en/map?lat=45.267136&lng=19.833549&radius=50000&zoom=10', {
    waitUntil: 'networkidle2',
    timeout: 60000
  });
  
  console.log('Waiting for map to load...');
  await new Promise(r => setTimeout(r, 10000));
  
  console.log('Taking screenshot...');
  await page.screenshot({ path: '/data/hostel-booking-system/backend/map-check.png', fullPage: true });
  
  console.log('Checking map state...');
  const mapState = await page.evaluate(() => {
    const result = {
      hasMap: typeof window._map !== 'undefined',
      hasLeaflet: typeof L !== 'undefined',
      layersCount: 0,
      hasTileLayer: false,
      markers: 0,
      clusters: 0,
      markerInfo: [],
      errors: [],
      mapContainer: null
    };
    
    // Check map container
    const mapDiv = document.getElementById('map') || document.querySelector('[id*="map"]') || document.querySelector('.map');
    if (mapDiv) {
      result.mapContainer = {
        id: mapDiv.id,
        className: mapDiv.className,
        hasChildren: mapDiv.children.length > 0,
        childrenCount: mapDiv.children.length
      };
    }
    
    // Check all possible map references
    const possibleMaps = ['window._map', 'window.map', 'window.leafletMap', 'window.mapInstance'];
    for (const mapRef of possibleMaps) {
      try {
        const map = eval(mapRef);
        if (map && typeof map === 'object') {
          result.foundMapIn = mapRef;
          result.hasMap = true;
          if (!window._map) window._map = map;
          break;
        }
      } catch (e) {}
    }
    
    if (window._map) {
      result.layersCount = Object.keys(window._map._layers).length;
      
      // Count markers and clusters properly
      window._map.eachLayer((layer) => {
        if (layer.options && layer.options.icon) {
          if (layer._childCount) {
            result.clusters++;
          } else {
            result.markers++;
            if (result.markerInfo.length < 5 && layer.getLatLng) {
              const latlng = layer.getLatLng();
              result.markerInfo.push({
                lat: latlng.lat,
                lng: latlng.lng
              });
            }
          }
        }
      });
    }
    
    return result;
  });
  
  console.log('\n=== РЕЗУЛЬТАТЫ АНАЛИЗА МАРКЕРОВ ===');
  console.log('Map exists:', mapState.hasMap);
  console.log('Leaflet loaded:', mapState.hasLeaflet);
  console.log('Markers:', mapState.markers);
  console.log('Clusters:', mapState.clusters);
  
  if (mapState.markerInfo.length > 0) {
    console.log('\nFirst 5 markers:');
    mapState.markerInfo.forEach((info, i) => {
      console.log(`  ${i + 1}. [${info.lat.toFixed(6)}, ${info.lng.toFixed(6)}]`);
    });
  }
  
  console.log('\nMap container:', mapState.mapContainer);
  
  // Show console errors
  const errors = consoleLogs.filter(log => log.type === 'error' || log.type === 'pageerror');
  console.log('\n=== ОШИБКИ КОНСОЛИ ===');
  if (errors.length > 0) {
    errors.forEach(err => {
      console.log(`[${err.type}] ${err.text}`);
      if (err.location) {
        console.log(`  Location: ${err.location.url}:${err.location.lineNumber}`);
      }
    });
  } else {
    console.log('Ошибок не найдено');
  }
  
  if (networkErrors.length > 0) {
    console.log('\n=== NETWORK ERRORS ===');
    networkErrors.forEach(err => {
      console.log(`Failed: ${err.url}`);
      console.log(`Reason: ${err.failure ? err.failure.errorText : 'Unknown'}`);
    });
  }
  
  await browser.close();
})();