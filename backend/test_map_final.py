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
            await page.goto('http://localhost:3001/en/map', wait_until='domcontentloaded', timeout=15000)
            await page.wait_for_timeout(5000)  # Wait for map to load
            
            print("2. Taking initial screenshot...")
            await page.screenshot(path='/data/hostel-booking-system/backend/initial-map.png', full_page=True)
            print("   ✓ Initial screenshot saved")
            
            print("3. Finding map container...")
            # Try multiple selectors
            map_element = None
            selectors = ['.leaflet-container', '#map', '[class*="map"]', 'div[id*="map"]']
            
            for selector in selectors:
                map_element = await page.query_selector(selector)
                if map_element:
                    print(f"   ✓ Found map using selector: {selector}")
                    break
            
            if not map_element:
                print("   ✗ Map container not found!")
                # Take debug screenshot
                await page.screenshot(path='/data/hostel-booking-system/backend/debug-no-map.png', full_page=True)
                
                # Log page structure
                divs = await page.evaluate('''() => {
                    return Array.from(document.querySelectorAll('div')).map(d => ({
                        id: d.id,
                        class: d.className,
                        width: d.offsetWidth,
                        height: d.offsetHeight
                    })).filter(d => d.width > 200 && d.height > 200);
                }''')
                print("\n   Large divs on page:")
                for div in divs[:10]:
                    print(f"     - id: {div['id']}, class: {div['class']}, size: {div['width']}x{div['height']}")
            else:
                print("4. Clicking map center...")
                box = await map_element.bounding_box()
                if box:
                    center_x = box['x'] + box['width'] / 2
                    center_y = box['y'] + box['height'] / 2
                    await page.mouse.click(center_x, center_y)
                    print("   ✓ Clicked map center")
                
                print("5. Zooming in...")
                # Move to map center first
                await page.mouse.move(center_x, center_y)
                
                # Zoom in using wheel
                for i in range(4):
                    await page.mouse.wheel(0, -100)
                    await page.wait_for_timeout(1000)
                print("   ✓ Zoomed in 4 times")
                
                print("6. Waiting for tiles to load...")
                await page.wait_for_timeout(3000)
                
                print("7. Taking zoomed screenshot...")
                await page.screenshot(path='/data/hostel-booking-system/backend/zoomed-map.png', full_page=True)
                print("   ✓ Zoomed screenshot saved")
            
            print("8. Searching for District/Район elements...")
            # Get all text content
            all_text = await page.evaluate('() => document.body.innerText')
            
            # Search for District or Район
            has_district = 'District' in all_text or 'Район' in all_text
            
            if has_district:
                print("   ✓ Found District/Район text in page content")
                
                # Find specific elements containing the text
                elements_info = await page.evaluate('''() => {
                    const results = [];
                    const allElements = document.querySelectorAll('*');
                    
                    for (const el of allElements) {
                        if (el.children.length === 0) { // Only leaf nodes
                            const text = el.textContent || '';
                            if (text.includes('District') || text.includes('Район')) {
                                results.push({
                                    tag: el.tagName,
                                    class: el.className || 'no-class',
                                    id: el.id || 'no-id',
                                    text: text.trim().substring(0, 100)
                                });
                            }
                        }
                    }
                    return results.slice(0, 20);
                }''')
                
                if elements_info:
                    print(f"\n   Found {len(elements_info)} elements with District/Район:")
                    for idx, el in enumerate(elements_info[:10], 1):
                        print(f"     {idx}. {el['tag']} (class: {el['class']}, id: {el['id']})")
                        print(f"        Text: {el['text']}")
            else:
                print("   ✗ No District/Район text found")
            
            print("\n9. Taking final screenshot...")
            await page.screenshot(path='/data/hostel-booking-system/backend/final-state.png', full_page=True)
            print("   ✓ Final screenshot saved")
            
            print("\n✅ All tests completed!")
            
            # Summary
            print("\n📊 Summary:")
            print(f"   - Map found: {'Yes' if map_element else 'No'}")
            print(f"   - District/Район text found: {'Yes' if has_district else 'No'}")
            print("   - Screenshots saved:")
            print("     • initial-map.png")
            if map_element:
                print("     • zoomed-map.png")
            else:
                print("     • debug-no-map.png")
            print("     • final-state.png")
            
        except Exception as e:
            print(f"\n❌ Error occurred: {str(e)}")
            # Take error screenshot
            await page.screenshot(path='/data/hostel-booking-system/backend/error-state.png', full_page=True)
            print("   Error screenshot saved as error-state.png")
            
        finally:
            await browser.close()

if __name__ == "__main__":
    asyncio.run(test_map())