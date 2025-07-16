const { chromium } = require('@playwright/test');

(async () => {
  let browser;
  try {
    console.log('Запуск браузера...');
    browser = await chromium.launch({ headless: true });
    const page = await browser.newPage();

    console.log('Переход на страницу карты...');
    await page.goto('http://localhost:3001/ru/map', {
      waitUntil: 'domcontentloaded',
      timeout: 60000,
    });

    console.log('Ожидание загрузки страницы...');
    await page.waitForTimeout(10000);

    // Проверка элементов
    const results = {};

    try {
      await page.waitForSelector('.leaflet-container', { timeout: 30000 });
      const leafletContainers = await page.$$('.leaflet-container');
      results.leafletContainer = leafletContainers.length;
      console.log(`✓ Найдено .leaflet-container: ${leafletContainers.length}`);
    } catch (e) {
      results.leafletContainer = 0;
      console.log('✗ .leaflet-container не найден');
    }

    try {
      await page.waitForSelector('.leaflet-marker-icon', { timeout: 30000 });
      const markers = await page.$$('.leaflet-marker-icon');
      results.markers = markers.length;
      console.log(`✓ Найдено .leaflet-marker-icon: ${markers.length}`);
    } catch (e) {
      results.markers = 0;
      console.log('✗ .leaflet-marker-icon не найден');
    }

    try {
      await page.waitForSelector('.rc-slider', { timeout: 30000 });
      const sliders = await page.$$('.rc-slider');
      results.sliders = sliders.length;
      console.log(`✓ Найдено .rc-slider: ${sliders.length}`);
    } catch (e) {
      results.sliders = 0;
      console.log('✗ .rc-slider не найден');
    }

    // Дополнительное ожидание
    console.log('Ожидание 2 секунды...');
    await page.waitForTimeout(2000);

    // Скриншот
    console.log('Создание скриншота...');
    await page.screenshot({ path: '/tmp/map-loaded.png', fullPage: true });

    console.log('\nИтоговые результаты:');
    console.log(JSON.stringify(results, null, 2));
  } catch (error) {
    console.error('Ошибка:', error.message);
  } finally {
    if (browser) {
      await browser.close();
    }
  }
})();
