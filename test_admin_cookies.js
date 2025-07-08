const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();
  
  console.log('1. Открываю http://localhost:3001');
  await page.goto('http://localhost:3001');
  await page.waitForTimeout(2000);
  
  console.log('2. Очищаю все cookies для домена localhost:3001');
  await context.clearCookies();
  console.log('Cookies очищены');
  
  console.log('3. Обновляю страницу');
  await page.reload();
  await page.waitForTimeout(2000);
  
  console.log('4. Перехожу на http://localhost:3001/admin/search');
  await page.goto('http://localhost:3001/admin/search');
  await page.waitForTimeout(3000);
  
  console.log('5. Делаю скриншот');
  await page.screenshot({ path: 'admin-search-after-cookies-clear.png', fullPage: true });
  console.log('Скриншот сохранен как admin-search-after-cookies-clear.png');
  
  await browser.close();
})();