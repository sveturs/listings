const { chromium } = require('playwright');

(async () => {
  console.log('Запуск браузера...');
  const browser = await chromium.launch({
    headless: false,
    args: ['--no-sandbox', '--disable-setuid-sandbox'],
  });

  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
  });

  const page = await context.newPage();

  try {
    console.log('Переход на страницу карты...');
    await page.goto('http://localhost:3001/ru/map', {
      waitUntil: 'networkidle',
      timeout: 30000,
    });

    console.log('Ожидание загрузки карты...');
    await page.waitForSelector('.leaflet-container', { timeout: 30000 });

    // Даем карте время полностью загрузиться
    await page.waitForTimeout(3000);

    console.log('Поиск поля ввода...');
    // Ищем поле поиска - пробуем разные селекторы
    const searchSelectors = [
      'input[placeholder*="Поиск"]',
      'input[placeholder*="поиск"]',
      'input[type="search"]',
      'input[type="text"]',
      '.search-input',
      '[data-testid="search-input"]',
    ];

    let searchInput = null;
    for (const selector of searchSelectors) {
      try {
        searchInput = await page.waitForSelector(selector, { timeout: 5000 });
        if (searchInput) {
          console.log(`Найдено поле поиска по селектору: ${selector}`);
          break;
        }
      } catch (e) {
        continue;
      }
    }

    if (!searchInput) {
      throw new Error('Не удалось найти поле поиска');
    }

    console.log('Ввод текста "Нови Београд"...');
    await searchInput.click();
    await searchInput.fill('Нови Београд');

    // Ждем появления результатов поиска
    await page.waitForTimeout(2000);

    console.log('Создание скриншота...');
    await page.screenshot({
      path: '/tmp/search-novi-beograd.png',
      fullPage: true,
    });

    // Сохраняем информацию о результатах
    const pageContent = await page.content();
    const results = {
      timestamp: new Date().toISOString(),
      url: page.url(),
      searchTerm: 'Нови Београд',
      screenshotPath: '/tmp/search-novi-beograd.png',
      hasSearchInput: !!searchInput,
      pageTitle: await page.title(),
    };

    const fs = require('fs');
    fs.writeFileSync(
      '/tmp/search-result.txt',
      JSON.stringify(results, null, 2)
    );

    console.log('Тест завершен успешно!');
    console.log('Скриншот сохранен в: /tmp/search-novi-beograd.png');
    console.log('Результаты сохранены в: /tmp/search-result.txt');
  } catch (error) {
    console.error('Ошибка при выполнении теста:', error);

    // Сохраняем скриншот ошибки
    await page.screenshot({
      path: '/tmp/search-error.png',
      fullPage: true,
    });

    const fs = require('fs');
    fs.writeFileSync(
      '/tmp/search-result.txt',
      JSON.stringify(
        {
          error: error.message,
          timestamp: new Date().toISOString(),
          errorScreenshot: '/tmp/search-error.png',
        },
        null,
        2
      )
    );
  } finally {
    await browser.close();
  }
})();
