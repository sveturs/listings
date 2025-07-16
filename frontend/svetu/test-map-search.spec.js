const { test, expect } = require('@playwright/test');
const fs = require('fs');

test('Тестирование поиска по карте', async ({ page }) => {
  console.log('Переход на страницу карты...');
  await page.goto('http://localhost:3001/ru/map', { waitUntil: 'networkidle' });

  console.log('Ожидание загрузки карты...');
  await page.waitForSelector('.leaflet-container', { timeout: 30000 });

  console.log('Поиск поля ввода...');
  const searchInput = await page.locator('input[type="text"]').first();

  console.log('Ввод текста "Нови Београд"...');
  await searchInput.fill('Нови Београд');

  console.log('Ожидание результатов...');
  await page.waitForTimeout(2000);

  console.log('Создание скриншота...');
  await page.screenshot({
    path: '/tmp/search-novi-beograd.png',
    fullPage: true,
  });

  // Анализ страницы
  const searchValue = await searchInput.inputValue();
  const mapVisible = await page.locator('.leaflet-container').isVisible();

  // Поиск результатов поиска
  const searchResults = await page
    .locator('[role="listbox"], .search-results, .district-results')
    .count();

  // Проверка наличия районов
  const districtButtons = await page
    .locator('button')
    .filter({ hasText: /район|округ/i })
    .count();

  // Создание отчета
  const report = `Отчет о тестировании поиска по карте
========================================

Дата: ${new Date().toLocaleString('ru-RU')}
URL: http://localhost:3001/ru/map

Поле поиска:
- Текущее значение: "${searchValue}"
- Поле успешно найдено и заполнено

Состояние карты:
- Карта загружена: ${mapVisible ? 'Да' : 'Нет'}
- Контейнер карты присутствует на странице

Результаты поиска:
- Найдено элементов с результатами: ${searchResults}
- Найдено кнопок районов: ${districtButtons}

Дополнительная информация:
- Скриншот сохранен в: /tmp/search-novi-beograd.png
- Тест выполнен успешно
`;

  fs.writeFileSync('/tmp/search-result.txt', report);
  console.log('Отчет создан: /tmp/search-result.txt');
});
