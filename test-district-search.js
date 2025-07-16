const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: false });
  const page = await browser.newPage();
  
  // Слушаем консольные сообщения
  page.on('console', msg => {
    const text = msg.text();
    // Показываем только наши логи с эмодзи
    if (text.includes('🔍') || text.includes('🌍') || text.includes('🏙️') || 
        text.includes('📦') || text.includes('🗺️') || text.includes('📡') ||
        text.includes('📍')) {
      console.log('[CONSOLE]', text);
    }
  });

  console.log('📂 Открываем карту...');
  await page.goto('http://localhost:3001/ru/map');
  await page.waitForTimeout(3000);

  console.log('🔄 Переключаемся на поиск по районам...');
  // Ищем кнопку "По району"
  const districtButton = await page.locator('button:has-text("По району")');
  if (await districtButton.isVisible()) {
    await districtButton.click();
    await page.waitForTimeout(1000);
    
    console.log('📋 Выбираем район...');
    // Ищем select с районами
    const districtSelect = await page.locator('select').first();
    if (await districtSelect.isVisible()) {
      // Выбираем район "Врачар"
      await districtSelect.selectOption({ label: 'Врачар' });
      await page.waitForTimeout(2000);
      
      // Делаем скриншот
      await page.screenshot({ path: 'district-selected.png', fullPage: true });
      console.log('📸 Скриншот сохранен: district-selected.png');
    }
  }

  // Ждем еще немного для обработки всех логов
  await page.waitForTimeout(3000);
  
  await browser.close();
  console.log('✅ Тест завершен');
})();