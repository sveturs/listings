const { chromium } = require('playwright');

async function testFuzzySearch() {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();

  console.log('Открываю главную страницу...');
  await page.goto('http://localhost:3001');
  await page.waitForLoadState('networkidle');

  // Тест 1: Поиск с опечаткой "квортира" вместо "квартира"
  console.log('\nТест 1: Поиск "квортира" (с опечаткой)');
  
  // Сначала без нечеткого поиска
  await page.fill('input[type="search"]', 'квортира');
  await page.keyboard.press('Enter');
  await page.waitForTimeout(2000);
  
  // Делаем скриншот результатов без нечеткого поиска
  await page.screenshot({ 
    path: 'test-results/search-kvortira-without-fuzzy.png',
    fullPage: true 
  });
  
  // Включаем нечеткий поиск
  const fuzzyCheckbox = await page.locator('label:has-text("Нечеткий поиск")');
  await fuzzyCheckbox.click();
  await page.waitForTimeout(2000);
  
  // Делаем скриншот с нечетким поиском
  await page.screenshot({ 
    path: 'test-results/search-kvortira-with-fuzzy.png',
    fullPage: true 
  });

  // Тест 2: Поиск "телифон" (правильно написано)
  console.log('\nТест 2: Поиск "телифон"');
  await page.fill('input[type="search"]', 'телифон');
  await page.keyboard.press('Enter');
  await page.waitForTimeout(2000);
  
  await page.screenshot({ 
    path: 'test-results/search-telefon-with-fuzzy.png',
    fullPage: true 
  });

  // Тест 3: Поиск "машына" вместо "машина"
  console.log('\nТест 3: Поиск "машына" (с опечаткой)');
  await page.fill('input[type="search"]', 'машына');
  await page.keyboard.press('Enter');
  await page.waitForTimeout(2000);
  
  await page.screenshot({ 
    path: 'test-results/search-mashyna-with-fuzzy.png',
    fullPage: true 
  });

  // Выключаем нечеткий поиск для сравнения
  await fuzzyCheckbox.click();
  await page.waitForTimeout(2000);
  
  await page.screenshot({ 
    path: 'test-results/search-mashyna-without-fuzzy.png',
    fullPage: true 
  });

  // Тест 4: Поиск "аппартаменты" вместо "апартаменты"
  console.log('\nТест 4: Поиск "аппартаменты" (лишняя буква)');
  await fuzzyCheckbox.click(); // Включаем обратно
  await page.fill('input[type="search"]', 'аппартаменты');
  await page.keyboard.press('Enter');
  await page.waitForTimeout(2000);
  
  await page.screenshot({ 
    path: 'test-results/search-appartamenty-with-fuzzy.png',
    fullPage: true 
  });

  console.log('\nТестирование завершено!');
  await browser.close();
}

// Создаем директорию для результатов
const fs = require('fs');
if (!fs.existsSync('test-results')) {
  fs.mkdirSync('test-results');
}

testFuzzySearch().catch(console.error);