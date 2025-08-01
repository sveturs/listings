const { chromium } = require('playwright');
const fs = require('fs').promises;
const path = require('path');

// Конфигурация страниц для сохранения
const PAGES_TO_SAVE = [
  { url: 'http://localhost:3001/ru', name: 'homepage' },
  { url: 'http://localhost:3001/ru/search', name: 'search' },
  { url: 'http://localhost:3001/ru/create-listing-choice', name: 'create-listing-choice' },
  { url: 'http://localhost:3001/ru/create-listing', name: 'create-listing' },
  { url: 'http://localhost:3001/ru/create-listing-ai', name: 'create-listing-ai' },
  { url: 'http://localhost:3001/ru/create-listing-smart', name: 'create-listing-smart' },
  { url: 'http://localhost:3001/ru/examples/ideal-homepage', name: 'ideal-homepage' },
  { url: 'http://localhost:3001/ru/examples/ideal-homepage-v2', name: 'ideal-homepage-v2' },
  { url: 'http://localhost:3001/ru/examples/storefront-dashboard', name: 'storefront-dashboard' },
  { url: 'http://localhost:3001/ru/examples/listing-creation-ux-v2', name: 'listing-creation-ux-v2' },
  { url: 'http://localhost:3001/ru/examples/animated-chat', name: 'animated-chat' },
  { url: 'http://localhost:3001/ru/map', name: 'map' },
  { url: 'http://localhost:3001/ru/profile', name: 'profile' },
  { url: 'http://localhost:3001/ru/chat', name: 'chat' },
  // Добавьте другие страницы по необходимости
];

async function saveSiteForDesigner() {
  console.log('🚀 Запуск сохранения сайта для дизайнера...\n');
  
  // Создаем директории
  const outputDir = './designer-preview';
  const screenshotsDir = path.join(outputDir, 'screenshots');
  const pdfDir = path.join(outputDir, 'pdf');
  const fullPageDir = path.join(outputDir, 'full-pages');
  
  await fs.mkdir(outputDir, { recursive: true });
  await fs.mkdir(screenshotsDir, { recursive: true });
  await fs.mkdir(pdfDir, { recursive: true });
  await fs.mkdir(fullPageDir, { recursive: true });

  const browser = await chromium.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
    deviceScaleFactor: 2, // Retina качество
  });

  for (const pageInfo of PAGES_TO_SAVE) {
    console.log(`📄 Обработка страницы: ${pageInfo.name} (${pageInfo.url})`);
    
    const page = await context.newPage();
    
    try {
      // Переходим на страницу
      await page.goto(pageInfo.url, { 
        waitUntil: 'domcontentloaded',
        timeout: 15000 
      });
      
      // Ждем загрузки всех изображений
      await page.waitForTimeout(3000);
      
      // 1. Скриншоты (desktop)
      await page.screenshot({
        path: path.join(screenshotsDir, `${pageInfo.name}-desktop.png`),
        fullPage: true,
        animations: 'disabled'
      });
      console.log(`  ✅ Desktop скриншот сохранен`);
      
      // 2. Mobile версия
      await page.setViewportSize({ width: 375, height: 812 }); // iPhone 11 Pro
      await page.waitForTimeout(1000);
      
      await page.screenshot({
        path: path.join(screenshotsDir, `${pageInfo.name}-mobile.png`),
        fullPage: true,
        animations: 'disabled'
      });
      console.log(`  ✅ Mobile скриншот сохранен`);
      
      // 3. Tablet версия
      await page.setViewportSize({ width: 768, height: 1024 }); // iPad
      await page.waitForTimeout(1000);
      
      await page.screenshot({
        path: path.join(screenshotsDir, `${pageInfo.name}-tablet.png`),
        fullPage: true,
        animations: 'disabled'
      });
      console.log(`  ✅ Tablet скриншот сохранен`);
      
      // Возвращаем desktop viewport для PDF
      await page.setViewportSize({ width: 1920, height: 1080 });
      await page.waitForTimeout(1000);
      
      // 4. PDF версия
      await page.pdf({
        path: path.join(pdfDir, `${pageInfo.name}.pdf`),
        format: 'A4',
        printBackground: true,
        margin: { top: '20px', bottom: '20px', left: '20px', right: '20px' }
      });
      console.log(`  ✅ PDF сохранен`);
      
      // 5. Полный HTML со встроенными стилями и изображениями
      const content = await page.content();
      
      // Получаем все стили
      const styles = await page.evaluate(() => {
        const styleSheets = Array.from(document.styleSheets);
        let css = '';
        
        styleSheets.forEach(sheet => {
          try {
            const rules = Array.from(sheet.cssRules || sheet.rules || []);
            rules.forEach(rule => {
              css += rule.cssText + '\n';
            });
          } catch (e) {
            // Игнорируем CORS ошибки
          }
        });
        
        return css;
      });
      
      // Встраиваем стили в HTML
      const fullHtml = content.replace('</head>', `<style>${styles}</style></head>`);
      
      await fs.writeFile(
        path.join(fullPageDir, `${pageInfo.name}.html`),
        fullHtml,
        'utf-8'
      );
      console.log(`  ✅ HTML страница сохранена\n`);
      
    } catch (error) {
      console.error(`  ❌ Ошибка при обработке ${pageInfo.name}:`, error.message);
    } finally {
      await page.close();
    }
  }

  await browser.close();
  
  // Создаем index.html для навигации
  const indexHtml = `<!DOCTYPE html>
<html lang="ru">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Sve Tu - Превью для дизайнера</title>
  <style>
    body {
      font-family: Arial, sans-serif;
      max-width: 1200px;
      margin: 0 auto;
      padding: 20px;
      background: #f5f5f5;
    }
    h1 { color: #333; margin-bottom: 30px; }
    .pages-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
      gap: 20px;
      margin-bottom: 40px;
    }
    .page-card {
      background: white;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 0 2px 4px rgba(0,0,0,0.1);
    }
    .page-card h3 { margin-top: 0; color: #0066cc; }
    .links { display: flex; flex-direction: column; gap: 10px; }
    .links a {
      color: #0066cc;
      text-decoration: none;
      padding: 5px 10px;
      border: 1px solid #0066cc;
      border-radius: 4px;
      display: inline-block;
      transition: all 0.3s;
    }
    .links a:hover {
      background: #0066cc;
      color: white;
    }
    .info {
      background: #e3f2fd;
      padding: 15px;
      border-radius: 8px;
      margin-bottom: 30px;
    }
  </style>
</head>
<body>
  <h1>Sve Tu Platform - Превью для дизайнера</h1>
  
  <div class="info">
    <h3>📋 Информация о материалах:</h3>
    <ul>
      <li><strong>Screenshots</strong> - полные скриншоты страниц в 3х разрешениях (Desktop, Tablet, Mobile)</li>
      <li><strong>PDF</strong> - версии страниц для печати и аннотаций</li>
      <li><strong>HTML</strong> - интерактивные версии страниц (могут работать некорректно без сервера)</li>
    </ul>
  </div>
  
  <div class="pages-grid">
${PAGES_TO_SAVE.map(page => `
    <div class="page-card">
      <h3>${page.name}</h3>
      <div class="links">
        <a href="screenshots/${page.name}-desktop.png" target="_blank">🖥️ Desktop Screenshot</a>
        <a href="screenshots/${page.name}-tablet.png" target="_blank">📱 Tablet Screenshot</a>
        <a href="screenshots/${page.name}-mobile.png" target="_blank">📱 Mobile Screenshot</a>
        <a href="pdf/${page.name}.pdf" target="_blank">📄 PDF версия</a>
        <a href="full-pages/${page.name}.html" target="_blank">🌐 HTML версия</a>
      </div>
    </div>
`).join('')}
  </div>
</body>
</html>`;

  await fs.writeFile(path.join(outputDir, 'index.html'), indexHtml, 'utf-8');
  
  console.log(`
✅ Сохранение завершено!

📁 Все файлы сохранены в директорию: ${outputDir}/

Структура:
- index.html - главная страница с навигацией
- screenshots/ - скриншоты всех страниц (desktop, tablet, mobile)
- pdf/ - PDF версии страниц
- full-pages/ - HTML версии страниц

🚀 Для просмотра откройте ${outputDir}/index.html в браузере

💡 Совет: Заархивируйте папку ${outputDir} и отправьте дизайнеру
`);
}

// Запускаем
saveSiteForDesigner().catch(console.error);