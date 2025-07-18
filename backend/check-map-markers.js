const puppeteer = require('puppeteer');

(async () => {
  const browser = await puppeteer.launch({ 
    headless: false,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });
  const page = await browser.newPage();
  
  await page.goto('http://localhost:3001/en/map?lat=45.267136&lng=19.833549&radius=50000&zoom=10');
  await page.waitForFunction(() => window._map, { timeout: 10000 });
  
  const mapData = await page.evaluate(() => {
    if (!window._map) {
      return { error: 'Map not found' };
    }
    
    let markers = 0;
    let clusters = 0;
    const markerDetails = [];
    
    window._map.eachLayer((layer) => {
      if (layer instanceof L.Marker && !(layer instanceof L.MarkerCluster)) {
        markers++;
        if (markerDetails.length < 5) {
          const pos = layer.getLatLng();
          markerDetails.push({
            lat: pos.lat,
            lng: pos.lng,
            options: layer.options
          });
        }
      } else if (layer instanceof L.MarkerCluster) {
        clusters++;
      }
    });
    
    return {
      markers,
      clusters,
      markerDetails,
      mapCenter: window._map.getCenter(),
      mapZoom: window._map.getZoom()
    };
  });
  
  console.log('Map Analysis Results:');
  console.log(JSON.stringify(mapData, null, 2));
  
  const consoleErrors = await page.evaluate(() => {
    return window.__consoleErrors || [];
  });
  
  if (consoleErrors.length > 0) {
    console.log('\nConsole Errors:');
    consoleErrors.forEach(err => console.log(err));
  }
  
  await browser.close();
})();