#!/usr/bin/env python3
import asyncio
from playwright.async_api import async_playwright

async def check_pages():
    """Check available pages"""
    
    urls = [
        'http://localhost:3001/',
        'http://localhost:3001/en',
        'http://localhost:3001/en/map',
        'http://localhost:3001/map',
        'http://localhost:3001/en/test-map-simple',
        'http://localhost:3001/test-map-simple'
    ]
    
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()
        
        for url in urls:
            try:
                print(f"\nChecking {url}...")
                response = await page.goto(url, wait_until='domcontentloaded', timeout=10000)
                status = response.status if response else 'No response'
                title = await page.title()
                print(f"  Status: {status}")
                print(f"  Title: {title}")
                
                # Check for map elements
                map_elements = await page.query_selector_all('.leaflet-container, [class*="map"], #map')
                if map_elements:
                    print(f"  ✓ Found {len(map_elements)} map elements")
                
            except Exception as e:
                print(f"  ✗ Error: {str(e)}")
        
        await browser.close()

if __name__ == "__main__":
    asyncio.run(check_pages())