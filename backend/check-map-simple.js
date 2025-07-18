const puppeteer = require('puppeteer');

(async () => {
  console.log('Starting map check...');
  
  try {
    const browser = await puppeteer.launch({
      headless: false,
      args: ['--no-sandbox', '--disable-setuid-sandbox']
    });
    
    const page = await browser.newPage();
    console.log('Opening page...');
    
    await page.goto('http://localhost:3001/en/map?lat=45.267136&lng=19.833549&radius=50000&zoom=10');
    console.log('Page loaded, waiting for map...');
    
    await new Promise(resolve => setTimeout(resolve, 5000));
    
    // Take screenshot
    await page.screenshot({ path: 'map-check-2.png' });
    console.log('Screenshot saved');
    
    // Check map existence
    const hasMap = await page.evaluate(() => {
      return typeof window._map !== 'undefined';
    });
    
    console.log('Map found:', hasMap);
    
    if (hasMap) {
      // Count markers
      const markerCount = await page.evaluate(() => {
        let count = 0;
        window._map.eachLayer((layer) => {
          if (layer instanceof L.Marker) {
            count++;
          }
        });
        return count;
      });
      
      console.log('Marker count:', markerCount);
    }
    
    await browser.close();
    console.log('Done');
  } catch (error) {
    console.error('Error:', error);
  }
})();