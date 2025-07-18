const puppeteer = require('puppeteer');
const fs = require('fs');

async function checkMap() {
    const browser = await puppeteer.launch({
        headless: 'new',
        args: ['--no-sandbox', '--disable-setuid-sandbox']
    });
    
    try {
        const page = await browser.newPage();
        await page.setViewport({ width: 1920, height: 1080 });
        
        const url = 'http://localhost:3001/en/map?lat=45.2551&lng=19.8452&radius=50000&zoom=12';
        console.log(`Opening ${url}...`);
        
        await page.goto(url, { waitUntil: 'networkidle2', timeout: 30000 });
        
        // Wait for map to load
        console.log('Waiting for map to load...');
        await new Promise(resolve => setTimeout(resolve, 5000));
        
        // Take screenshot
        await page.screenshot({ path: '/tmp/map-screenshot.png' });
        console.log('Screenshot saved to /tmp/map-screenshot.png');
        
        // Count map elements
        const results = await page.evaluate(() => {
            const clusterSmall = document.querySelectorAll('.marker-cluster-small').length;
            const clusterMedium = document.querySelectorAll('.marker-cluster-medium').length;
            const clusterLarge = document.querySelectorAll('.marker-cluster-large').length;
            const singleMarkers = document.querySelectorAll('.leaflet-marker-icon:not(.marker-cluster)').length;
            
            return {
                clusterSmall,
                clusterMedium,
                clusterLarge,
                singleMarkers,
                totalClusters: clusterSmall + clusterMedium + clusterLarge,
                totalMarkers: clusterSmall + clusterMedium + clusterLarge + singleMarkers
            };
        });
        
        console.log('Map elements count:', results);
        
        // Try to click on a marker
        let popupOpened = false;
        try {
            const marker = await page.$('.leaflet-marker-icon');
            if (marker) {
                await marker.click();
                await new Promise(resolve => setTimeout(resolve, 2000));
                
                // Check if popup opened
                const popup = await page.$('.leaflet-popup');
                popupOpened = popup !== null;
            }
        } catch (e) {
            console.log('Error clicking marker:', e.message);
        }
        
        // Save results
        const report = {
            timestamp: new Date().toISOString(),
            url: url,
            clusters: {
                small: results.clusterSmall,
                medium: results.clusterMedium,
                large: results.clusterLarge,
                total: results.totalClusters
            },
            singleMarkers: results.singleMarkers,
            totalMarkers: results.totalMarkers,
            popupOpened: popupOpened,
            status: 'success'
        };
        
        fs.writeFileSync('/tmp/map-check-results.txt', JSON.stringify(report, null, 2));
        console.log('Results saved to /tmp/map-check-results.txt');
        
        return report;
        
    } catch (error) {
        console.error('Error:', error);
        const errorReport = {
            timestamp: new Date().toISOString(),
            url: 'http://localhost:3001/en/map?lat=45.2551&lng=19.8452&radius=50000&zoom=12',
            error: error.message,
            status: 'error'
        };
        fs.writeFileSync('/tmp/map-check-results.txt', JSON.stringify(errorReport, null, 2));
    } finally {
        await browser.close();
    }
}

checkMap().then(report => {
    console.log('\n=== Map Check Report ===');
    console.log(JSON.stringify(report, null, 2));
}).catch(err => {
    console.error('Fatal error:', err);
});