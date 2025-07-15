#!/usr/bin/env python3
import asyncio
from playwright.async_api import async_playwright
import os

async def test_map():
    """Test map functionality with screenshots"""
    
    async with async_playwright() as p:
        # Launch browser
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()
        
        try:
            print("1. Opening map page...")
            await page.goto('http://localhost:3001/en/map', wait_until='networkidle')
            await page.wait_for_timeout(3000)
            
            print("2. Taking initial screenshot...")
            await page.screenshot(path='/data/hostel-booking-system/backend/initial-map.png', full_page=True)
            print("   ✓ Initial screenshot saved")
            
            print("3. Finding map container...")
            map_element = await page.query_selector('.leaflet-container, [class*="map"]')
            if not map_element:
                print("   ✗ Map container not found!")
                return
            
            print("4. Clicking map center...")
            box = await map_element.bounding_box()
            if box:
                center_x = box['x'] + box['width'] / 2
                center_y = box['y'] + box['height'] / 2
                await page.mouse.click(center_x, center_y)
                print("   ✓ Clicked map center")
            
            print("5. Zooming in...")
            for i in range(4):
                await page.mouse.wheel(0, -100)
                await page.wait_for_timeout(500)
            print("   ✓ Zoomed in 4 times")
            
            print("6. Waiting for tiles to load...")
            await page.wait_for_timeout(3000)
            
            print("7. Taking zoomed screenshot...")
            await page.screenshot(path='/data/hostel-booking-system/backend/zoomed-map.png', full_page=True)
            print("   ✓ Zoomed screenshot saved")
            
            print("8. Searching for District/Район elements...")
            # Search for elements containing District or Район
            district_elements = await page.query_selector_all('text=/District|Район/i')
            print(f"   Found {len(district_elements)} elements with District/Район text")
            
            # Get all text content to search more broadly
            all_text = await page.evaluate('() => document.body.innerText')
            if 'District' in all_text or 'Район' in all_text:
                print("   ✓ Found District/Район text in page content")
                
                # Try to find specific elements
                elements_with_text = await page.evaluate('''() => {
                    const elements = [];
                    const allElements = document.querySelectorAll('*');
                    for (const el of allElements) {
                        const text = el.textContent || '';
                        if (text.includes('District') || text.includes('Район')) {
                            elements.push({
                                tag: el.tagName,
                                class: el.className,
                                text: text.substring(0, 100)
                            });
                        }
                    }
                    return elements.slice(0, 10); // Return first 10 matches
                }''')
                
                if elements_with_text:
                    print("\n   Elements containing District/Район:")
                    for el in elements_with_text:
                        print(f"     - {el['tag']} (class: {el['class']}): {el['text'][:50]}...")
            else:
                print("   ✗ No District/Район text found")
            
            print("\n9. Taking final screenshot...")
            await page.screenshot(path='/data/hostel-booking-system/backend/final-state.png', full_page=True)
            print("   ✓ Final screenshot saved")
            
            print("\n✅ All tests completed successfully!")
            
        except Exception as e:
            print(f"\n❌ Error occurred: {str(e)}")
            
        finally:
            await browser.close()

if __name__ == "__main__":
    asyncio.run(test_map())