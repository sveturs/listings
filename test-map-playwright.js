const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();
  
  // Массив для сбора логов с эмодзи
  const emojiLogs = [];
  
  // Настройка перехвата консольных логов
  page.on('console', message => {
    const text = message.text();
    // Проверяем наличие эмодзи в логе
    const emojiRegex = /[\u{1F300}-\u{1F9FF}]|[\u{2600}-\u{26FF}]|[\u{2700}-\u{27BF}]/u;
    if (emojiRegex.test(text)) {
      emojiLogs.push({
        type: message.type(),
        text: text,
        time: new Date().toISOString()
      });
    }
  });
  
  console.log('Открываю страницу карты...');
  await page.goto('http://localhost:3001/ru/map');
  
  // Ждем загрузки карты
  console.log('Жду загрузки карты...');
  await page.waitForSelector('.leaflet-container', { timeout: 30000 });
  await page.waitForTimeout(3000); // Дополнительное ожидание для полной загрузки
  
  // Нажимаем кнопку "По району"
  console.log('Переключаюсь на поиск по районам...');
  await page.click('button:has-text("По району")');
  await page.waitForTimeout(1000);
  
  // Выбираем район Врачар
  console.log('Выбираю район Врачар...');
  await page.click('text=Врачар');
  
  // Ждем 3 секунды после выбора района
  console.log('Жду 3 секунды после выбора района...');
  await page.waitForTimeout(3000);
  
  // Делаем скриншот
  console.log('Делаю скриншот...');
  await page.screenshot({ 
    path: '/data/hostel-booking-system/district-test.png',
    fullPage: true 
  });
  
  // Выводим все собранные логи с эмодзи
  console.log('\n=== КОНСОЛЬНЫЕ ЛОГИ С ЭМОДЗИ ===');
  if (emojiLogs.length === 0) {
    console.log('Логи с эмодзи не найдены');
  } else {
    emojiLogs.forEach((log, index) => {
      console.log(`\n[${index + 1}] ${log.time}`);
      console.log(`Тип: ${log.type}`);
      console.log(`Текст: ${log.text}`);
    });
  }
  
  await browser.close();
  console.log('\nТест завершен!');
})();