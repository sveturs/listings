#!/usr/bin/env python3
import asyncio
from playwright.async_api import async_playwright

async def debug_map():
    """Debug what's on the map page"""
    
    async with async_playwright() as p:
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()
        
        try:
            print("Opening map page...")
            await page.goto('http://localhost:3001/en/storefronts-map', wait_until='networkidle')
            await page.wait_for_timeout(5000)
            
            # Take screenshot
            await page.screenshot(path='/data/hostel-booking-system/backend/debug-page.png', full_page=True)
            print("Screenshot saved to debug-page.png")
            
            # Get page title
            title = await page.title()
            print(f"Page title: {title}")
            
            # Check for common map selectors
            selectors = [
                '.leaflet-container',
                '[class*="map"]',
                '#map',
                'div[id*="map"]',
                'canvas',
                '[class*="Map"]',
                '.maplibregl-map',
                '.mapboxgl-map'
            ]
            
            print("\nChecking for map elements:")
            for selector in selectors:
                elements = await page.query_selector_all(selector)
                if elements:
                    print(f"  ✓ Found {len(elements)} elements matching '{selector}'")
                    # Get first element details
                    first = elements[0]
                    box = await first.bounding_box()
                    if box:
                        print(f"    Size: {box['width']}x{box['height']}")
            
            # Get all divs with their classes
            print("\nAll div elements with classes:")
            divs = await page.evaluate('''() => {
                const divs = document.querySelectorAll('div[class]');
                return Array.from(divs).slice(0, 20).map(div => ({
                    class: div.className,
                    id: div.id || 'no-id',
                    width: div.offsetWidth,
                    height: div.offsetHeight
                })).filter(d => d.width > 100 && d.height > 100);
            }''')
            
            for div in divs:
                print(f"  - {div['class']} (id: {div['id']}, size: {div['width']}x{div['height']})")
            
            # Check page content
            text_content = await page.evaluate('() => document.body.innerText')
            print(f"\nPage text length: {len(text_content)} chars")
            if 'District' in text_content or 'Район' in text_content:
                print("✓ Found District/Район in page text")
            
        except Exception as e:
            print(f"Error: {str(e)}")
            
        finally:
            await browser.close()

if __name__ == "__main__":
    asyncio.run(debug_map())