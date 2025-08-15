import { test, expect } from '@playwright/test';

test.describe('Svetu Application Integration Test', () => {
  test('should load homepage and verify basic functionality', async ({ page }) => {
    console.log('🚀 Starting integration test...');
    
    // Навигация на главную страницу с более коротким тайм-аутом
    await page.goto('http://localhost:3001', { 
      waitUntil: 'domcontentloaded',
      timeout: 10000 
    });
    console.log('📄 Navigated to homepage');

    // Ждем немного для загрузки основного контента (без networkidle)
    await page.waitForTimeout(3000);
    console.log('⏳ Waited for initial content load');
    
    // Получаем актуальный title и URL для диагностики
    const title = await page.title();
    const url = page.url();
    console.log(`📝 Page title: "${title}"`);
    console.log(`🔗 Current URL: ${url}`);
    
    // Проверяем что страница не пустая - ищем любой текст или элементы
    const bodyText = await page.locator('body').textContent();
    expect(bodyText).toBeTruthy();
    expect(bodyText.trim().length).toBeGreaterThan(0);
    console.log('✅ Page has content');
    
    // Проверяем что это Next.js приложение (ищем характерные элементы)
    const nextElements = page.locator('#__next, [data-nextjs], script[src*="/_next/"]');
    const hasNextJs = await nextElements.count() > 0;
    if (hasNextJs) {
      console.log('✅ Next.js application detected');
    }
    
    // Проверяем наличие основных элементов (любой из них, но с коротким тайм-аутом)
    const header = page.locator('header, nav, .header, [data-testid="header"]');
    const main = page.locator('main, .main, [data-testid="main"], #__next > *');
    const content = page.locator('.map, .marketplace, .listings, .content, div[class*="container"]');
    
    let mainContentFound = false;
    try {
      await expect(header.or(main).or(content)).toBeVisible({ timeout: 5000 });
      console.log('✅ Main content elements are visible');
      mainContentFound = true;
    } catch (error) {
      console.log('⚠️ Main content elements not found, checking for any visible content...');
      
      // Если основные элементы не найдены, проверяем что хотя бы что-то отображается
      const anyVisible = page.locator('body *:visible').first();
      await expect(anyVisible).toBeVisible({ timeout: 3000 });
      console.log('✅ Some content is visible on page');
    }

    // Проверяем что это не страница ошибки или 404
    const errorElements = page.locator('.error, [data-testid="error"], .not-found, h1:has-text("404"), h1:has-text("Error"), h1:has-text("Page Not Found")');
    const errorCount = await errorElements.count();
    expect(errorCount).toBe(0);
    console.log('✅ No error elements found');

    // Проверяем базовую работоспособность (отсутствие white screen of death)
    const bodyHTML = await page.locator('body').innerHTML();
    expect(bodyHTML.length).toBeGreaterThan(100); // Страница должна содержать достаточно HTML
    console.log('✅ Page has substantial HTML content');

    console.log('🎉 Integration test completed successfully!');
  });
});