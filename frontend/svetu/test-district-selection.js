const { chromium } = require('playwright');

(async () => {
  const browser = await chromium.launch({ headless: false });
  const context = await browser.newContext();
  const page = await context.newPage();

  // Собираем логи консоли
  const consoleLogs = [];
  page.on('console', (msg) => {
    const text = msg.text();
    // Фильтруем логи с эмодзи
    if (text.match(/[\u{1F300}-\u{1F9FF}]/gu)) {
      consoleLogs.push({
        type: msg.type(),
        text: text,
      });
    }
  });

  console.log('📍 Открываю страницу карты...');
  await page.goto('http://localhost:3001/ru/map');

  // Ждем загрузки карты
  await page.waitForTimeout(3000);

  console.log('🔍 Кликаю на кнопку "По району"...');
  await page.click('button:has-text("По району")');

  // Ждем появления интерфейса выбора района
  await page.waitForTimeout(2000);

  // Делаем скриншот после клика
  await page.screenshot({ path: 'district-click-test.png' });
  console.log('📸 Сохранен скриншот после клика: district-click-test.png');

  // Проверяем, появился ли селектор районов
  const districtSelector = await page.$('select.select-bordered');
  if (districtSelector) {
    console.log('✅ Селектор районов найден!');

    // Получаем все опции
    const options = await page.$$eval('select.select-bordered option', (opts) =>
      opts.map((opt) => ({ value: opt.value, text: opt.textContent }))
    );

    console.log('📋 Доступные районы:', options);

    // Ищем Врачар
    const vracarOption = options.find(
      (opt) => opt.text && opt.text.includes('Врачар')
    );
    if (vracarOption) {
      console.log('🎯 Найден район Врачар:', vracarOption);

      // Выбираем Врачар
      await page.selectOption('select.select-bordered', vracarOption.value);
      console.log('✅ Выбран район Врачар');

      // Ждем 3 секунды
      await page.waitForTimeout(3000);

      // Делаем финальный скриншот
      await page.screenshot({ path: 'district-test.png' });
      console.log('📸 Сохранен финальный скриншот: district-test.png');
    } else {
      console.log('❌ Район Врачар не найден в списке');
    }
  } else {
    console.log('❌ Селектор районов не появился');

    // Проверяем, есть ли сообщение об ошибке
    const errorMessage = await page.$('.alert-error');
    if (errorMessage) {
      const errorText = await errorMessage.textContent();
      console.log('⚠️ Найдено сообщение об ошибке:', errorText);
    }

    // Проверяем, есть ли информационное сообщение
    const infoMessage = await page.$('.alert-info');
    if (infoMessage) {
      const infoText = await infoMessage.textContent();
      console.log('ℹ️ Найдено информационное сообщение:', infoText);
    }
  }

  console.log('\n📝 Логи консоли с эмодзи:');
  consoleLogs.forEach((log, index) => {
    console.log(`${index + 1}. [${log.type.toUpperCase()}] ${log.text}`);
  });

  await browser.close();
})();
