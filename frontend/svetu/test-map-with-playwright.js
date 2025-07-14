const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();

  console.log('Открываю страницу карты...');
  await page.goto('http://localhost:3001/ru/map', {
    waitUntil: 'domcontentloaded',
    timeout: 60000,
  });

  console.log('Страница загружена, проверяю элементы...');

  const results = {
    leafletContainer: null,
    leafletMarkers: null,
    rcSliders: null,
  };

  // Проверка .leaflet-container
  try {
    console.log('Ожидание .leaflet-container...');
    await page.waitForSelector('.leaflet-container', { timeout: 30000 });
    const containers = await page.$$('.leaflet-container');
    results.leafletContainer = containers.length;
    console.log(`✓ Найдено .leaflet-container: ${containers.length}`);
  } catch (error) {
    console.log('✗ .leaflet-container не найден за 30 секунд');
    results.leafletContainer = 0;
  }

  // Проверка .leaflet-marker-icon
  try {
    console.log('Ожидание .leaflet-marker-icon...');
    await page.waitForSelector('.leaflet-marker-icon', { timeout: 30000 });
    const markers = await page.$$('.leaflet-marker-icon');
    results.leafletMarkers = markers.length;
    console.log(`✓ Найдено .leaflet-marker-icon: ${markers.length}`);
  } catch (error) {
    console.log('✗ .leaflet-marker-icon не найден за 30 секунд');
    results.leafletMarkers = 0;
  }

  // Проверка .rc-slider
  try {
    console.log('Ожидание .rc-slider...');
    await page.waitForSelector('.rc-slider', { timeout: 30000 });
    const sliders = await page.$$('.rc-slider');
    results.rcSliders = sliders.length;
    console.log(`✓ Найдено .rc-slider: ${sliders.length}`);
  } catch (error) {
    console.log('✗ .rc-slider не найден за 30 секунд');
    results.rcSliders = 0;
  }

  // Ждем 2 секунды после появления всех элементов
  console.log('Ожидание 2 секунды...');
  await page.waitForTimeout(2000);

  // Делаем скриншот
  console.log('Создание скриншота...');
  await page.screenshot({ path: '/tmp/map-loaded.png', fullPage: true });
  console.log('✓ Скриншот сохранен в /tmp/map-loaded.png');

  console.log('\nИтоговые результаты:');
  console.log(
    `- Контейнеры карты (.leaflet-container): ${results.leafletContainer}`
  );
  console.log(`- Маркеры (.leaflet-marker-icon): ${results.leafletMarkers}`);
  console.log(`- Слайдеры (.rc-slider): ${results.rcSliders}`);

  await browser.close();
})();
