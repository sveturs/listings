const { chromium } = require('playwright');
const fs = require('fs').promises;
const path = require('path');

async function saveChatPage() {
  console.log('🔍 Сохранение страницы чата с увеличенной паузой...\n');
  
  const browser = await chromium.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const context = await browser.newContext({
    viewport: { width: 1920, height: 1080 },
  });

  const page = await context.newPage();
  
  try {
    console.log('📄 Переход на страницу чата...');
    await page.goto('http://localhost:3001/ru/chat', { 
      waitUntil: 'networkidle',
      timeout: 30000 
    });
    
    console.log('⏳ Ожидание загрузки страницы (10 секунд)...');
    await page.waitForTimeout(10000);
    
    // Проверяем что загрузилось
    const title = await page.title();
    console.log(`📋 Заголовок страницы: ${title}`);
    
    // Desktop скриншот
    const screenshotPath = './designer-preview-v2/screenshots/chat-desktop-fixed.png';
    await page.screenshot({
      path: screenshotPath,
      fullPage: true,
    });
    console.log(`✅ Скриншот сохранен: ${screenshotPath}`);
    
    // Mobile версия
    await page.setViewportSize({ width: 375, height: 812 });
    await page.waitForTimeout(2000);
    
    const mobileScreenshotPath = './designer-preview-v2/screenshots/chat-mobile-fixed.png';
    await page.screenshot({
      path: mobileScreenshotPath,
      fullPage: true,
    });
    console.log(`✅ Mobile скриншот сохранен: ${mobileScreenshotPath}`);
    
    // PDF
    await page.setViewportSize({ width: 1920, height: 1080 });
    await page.waitForTimeout(1000);
    
    const pdfPath = './designer-preview-v2/pdf/chat-fixed.pdf';
    await page.pdf({
      path: pdfPath,
      format: 'A4',
      printBackground: true,
    });
    console.log(`✅ PDF сохранен: ${pdfPath}`);
    
  } catch (error) {
    console.error(`❌ Ошибка: ${error.message}`);
  } finally {
    await browser.close();
  }
  
  console.log('\n✅ Готово! Проверьте файлы в designer-preview-v2/');
}

saveChatPage().catch(console.error);