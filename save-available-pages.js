const { chromium } = require('playwright');
const fs = require('fs').promises;
const path = require('path');

// Сначала проверим, какие страницы доступны
const PAGES_TO_CHECK = [
  { url: 'http://localhost:3001/ru', name: 'homepage' },
  { url: 'http://localhost:3001/ru/search', name: 'search' },
  { url: 'http://localhost:3001/ru/create-listing-choice', name: 'create-listing-choice' },
  { url: 'http://localhost:3001/ru/examples', name: 'examples' },
  { url: 'http://localhost:3001/ru/examples/ideal-homepage', name: 'ideal-homepage' },
  { url: 'http://localhost:3001/ru/examples/ideal-homepage-v2', name: 'ideal-homepage-v2' },
  { url: 'http://localhost:3001/ru/examples/storefront-dashboard', name: 'storefront-dashboard' },
  { url: 'http://localhost:3001/ru/examples/listing-creation-ux-v2', name: 'listing-creation-ux-v2' },
  { url: 'http://localhost:3001/ru/examples/animated-chat', name: 'animated-chat' },
  { url: 'http://localhost:3001/ru/map', name: 'map' },
  { url: 'http://localhost:3001/ru/profile', name: 'profile' },
  { url: 'http://localhost:3001/ru/chat', name: 'chat' },
  { url: 'http://localhost:3001/ru/examples/quick-view', name: 'quick-view' },
  { url: 'http://localhost:3001/ru/examples/smart-search', name: 'smart-search' },
  { url: 'http://localhost:3001/ru/examples/map-privacy', name: 'map-privacy' },
];

async function checkAndSavePages() {
  console.log('🔍 Проверка доступных страниц...\n');
  
  const browser = await chromium.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
  });

  const availablePages = [];
  
  // Проверяем доступность страниц
  for (const pageInfo of PAGES_TO_CHECK) {
    const page = await context.newPage();
    try {
      const response = await page.goto(pageInfo.url, { 
        waitUntil: 'domcontentloaded',
        timeout: 10000 
      });
      
      if (response.status() === 200) {
        console.log(`✅ ${pageInfo.name}: Доступна`);
        availablePages.push(pageInfo);
      } else {
        console.log(`❌ ${pageInfo.name}: Статус ${response.status()}`);
      }
    } catch (error) {
      console.log(`❌ ${pageInfo.name}: Ошибка - ${error.message}`);
    } finally {
      await page.close();
    }
  }

  console.log(`\n📊 Найдено доступных страниц: ${availablePages.length}\n`);

  if (availablePages.length === 0) {
    console.log('❌ Нет доступных страниц для сохранения');
    await browser.close();
    return;
  }

  // Создаем директории
  const outputDir = './designer-preview-v2';
  const screenshotsDir = path.join(outputDir, 'screenshots');
  const pdfDir = path.join(outputDir, 'pdf');
  
  await fs.mkdir(outputDir, { recursive: true });
  await fs.mkdir(screenshotsDir, { recursive: true });
  await fs.mkdir(pdfDir, { recursive: true });

  console.log('📸 Начинаем сохранение страниц...\n');

  // Сохраняем только доступные страницы
  for (const pageInfo of availablePages) {
    console.log(`📄 Сохранение: ${pageInfo.name}`);
    
    const page = await context.newPage();
    
    try {
      await page.goto(pageInfo.url, { 
        waitUntil: 'domcontentloaded',
        timeout: 10000 
      });
      
      // Даем странице время на загрузку
      // Для чата и других динамических страниц нужно больше времени
      const waitTime = pageInfo.name === 'chat' || pageInfo.name === 'animated-chat' ? 5000 : 2000;
      await page.waitForTimeout(waitTime);
      
      // Для страницы чата ждем загрузку компонентов
      if (pageInfo.name === 'chat') {
        try {
          // Ждем появления чат-интерфейса
          await page.waitForSelector('[data-testid="chat-container"], .chat-container, #chat-root', { 
            timeout: 5000,
            state: 'visible' 
          });
          await page.waitForTimeout(1000); // Дополнительная пауза после загрузки
        } catch (e) {
          console.log('  ⚠️  Чат интерфейс не найден, продолжаем...');
        }
      }
      
      // Desktop скриншот
      await page.screenshot({
        path: path.join(screenshotsDir, `${pageInfo.name}-desktop.png`),
        fullPage: true,
      });
      
      // Mobile версия
      await page.setViewportSize({ width: 375, height: 812 });
      await page.waitForTimeout(500);
      
      await page.screenshot({
        path: path.join(screenshotsDir, `${pageInfo.name}-mobile.png`),
        fullPage: true,
      });
      
      // PDF
      await page.setViewportSize({ width: 1920, height: 1080 });
      await page.pdf({
        path: path.join(pdfDir, `${pageInfo.name}.pdf`),
        format: 'A4',
        printBackground: true,
      });
      
      console.log(`  ✅ Сохранено успешно\n`);
      
    } catch (error) {
      console.error(`  ❌ Ошибка: ${error.message}\n`);
    } finally {
      await page.close();
    }
  }

  await browser.close();

  // Создаем index.html
  const indexHtml = `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Sve Tu - Доступные страницы</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      max-width: 1200px;
      margin: 0 auto;
      padding: 20px;
      background: #f5f5f5;
    }
    h1 { color: #333; }
    .info {
      background: #e3f2fd;
      padding: 15px;
      border-radius: 8px;
      margin-bottom: 30px;
    }
    .pages-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 20px;
    }
    .page-card {
      background: white;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .page-card h3 { margin-top: 0; }
    .links a {
      display: inline-block;
      margin: 5px 0;
      color: #0066cc;
      text-decoration: none;
    }
    .links a:hover { text-decoration: underline; }
  </style>
</head>
<body>
  <h1>Sve Tu Platform - Доступные страницы</h1>
  
  <div class="info">
    <p>Сохранено страниц: ${availablePages.length}</p>
    <p>Формат: PNG скриншоты (Desktop/Mobile) + PDF</p>
  </div>
  
  <div class="pages-grid">
${availablePages.map(page => `
    <div class="page-card">
      <h3>${page.name}</h3>
      <div class="links">
        <a href="screenshots/${page.name}-desktop.png" target="_blank">🖥️ Desktop</a><br>
        <a href="screenshots/${page.name}-mobile.png" target="_blank">📱 Mobile</a><br>
        <a href="pdf/${page.name}.pdf" target="_blank">📄 PDF</a>
      </div>
    </div>
`).join('')}
  </div>
</body>
</html>`;

  await fs.writeFile(path.join(outputDir, 'index.html'), indexHtml, 'utf-8');
  
  console.log(`
✅ Готово!
📁 Файлы сохранены в: ${outputDir}/
🌐 Откройте index.html для навигации
`);
}

checkAndSavePages().catch(console.error);