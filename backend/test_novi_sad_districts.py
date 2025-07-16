#!/usr/bin/env python3
"""
Тестирование фильтра районов в Новом Саде с помощью Playwright
"""
import asyncio
import time
from playwright.async_api import async_playwright
import os

async def test_novi_sad_districts():
    print("Запуск тестирования фильтра районов в Новом Саде...")
    
    async with async_playwright() as p:
        # Запуск браузера
        print("1. Запуск Chrome браузера...")
        browser = await p.chromium.launch(headless=True)
        page = await browser.new_page()
        
        # Переход на страницу
        url = "http://localhost:3001/ru/map?lat=45.2671&lng=19.8335&zoom=11"
        print(f"2. Переход на страницу: {url}")
        await page.goto(url)
        
        # Ожидание загрузки
        print("3. Ожидание загрузки страницы (3 секунды)...")
        await page.wait_for_timeout(3000)
        
        # Первый скриншот
        print("4. Создание первого скриншота...")
        await page.screenshot(path="/tmp/novi-sad-initial.png", full_page=False)
        print("   Скриншот сохранен: /tmp/novi-sad-initial.png")
        
        # Поиск кнопки с текстом 'По району'
        print("5. Поиск кнопки с текстом 'По району'...")
        district_button = await page.get_by_text("По району").first.wait_for(state="visible", timeout=5000)
        
        # Клик по кнопке
        print("6. Клик по кнопке 'По району'...")
        await page.get_by_text("По району").first.click()
        
        # Ожидание
        print("7. Ожидание открытия селектора районов (2 секунды)...")
        await page.wait_for_timeout(2000)
        
        # Второй скриншот
        print("8. Создание скриншота с открытым селектором районов...")
        await page.screenshot(path="/tmp/novi-sad-district-selector.png", full_page=False)
        print("   Скриншот сохранен: /tmp/novi-sad-district-selector.png")
        
        # Попробуем найти селектор или поле для выбора района
        print("9. Поиск поля выбора района...")
        try:
            # Ищем селект или инпут для выбора района
            district_selector = await page.locator("select, input[placeholder*='район']").first.wait_for(state="visible", timeout=5000)
            print("   Найден элемент для выбора района")
            
            # Делаем скриншот текущего состояния
            await page.screenshot(path="/tmp/novi-sad-district-field.png", full_page=False)
            print("   Скриншот сохранен: /tmp/novi-sad-district-field.png")
            
            # Если это селект, выбираем Лиман
            if await district_selector.evaluate("el => el.tagName === 'SELECT'"):
                print("10. Выбор района 'Лиман' из списка...")
                await district_selector.select_option(label="Лиман")
            else:
                # Если это инпут, вводим текст
                print("10. Ввод названия района 'Лиман'...")
                await district_selector.fill("Лиман")
                await page.wait_for_timeout(1000)
                # Пробуем кликнуть по появившемуся варианту
                await page.get_by_text("Лиман").first.click()
        except Exception as e:
            print(f"   Не удалось найти селектор района: {e}")
            print("   Пропускаем выбор района...")
        
        # Ожидание применения фильтра
        print("11. Ожидание применения фильтра (2 секунды)...")
        await page.wait_for_timeout(2000)
        
        # Финальный скриншот
        print("12. Создание финального скриншота...")
        await page.screenshot(path="/tmp/novi-sad-liman.png", full_page=False)
        print("   Скриншот сохранен: /tmp/novi-sad-liman.png")
        
        # Закрытие браузера
        print("13. Закрытие браузера...")
        await browser.close()
        
        print("\nТестирование завершено!")
        print("\nСозданные скриншоты:")
        
        # Проверяем существование файлов
        import os
        screenshots = [
            ("/tmp/novi-sad-initial.png", "начальное состояние карты"),
            ("/tmp/novi-sad-district-selector.png", "после клика на 'По району'"),
            ("/tmp/novi-sad-district-field.png", "поле выбора района (если найдено)"),
            ("/tmp/novi-sad-liman.png", "финальное состояние карты")
        ]
        
        for path, description in screenshots:
            if os.path.exists(path):
                size = os.path.getsize(path)
                print(f"✓ {path} ({size} байт) - {description}")
            else:
                print(f"✗ {path} - не создан")

if __name__ == "__main__":
    asyncio.run(test_novi_sad_districts())