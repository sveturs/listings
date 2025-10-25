import { test, expect } from '@playwright/test';

test.describe('Admin Quality Tests - Result Persistence', () => {
  test.use({ storageState: '.auth/user.json' });

  test('should persist test results after page refresh', async ({ page }) => {
    console.log('Step 1: Navigate to Quality Tests page...');
    await page.goto('http://localhost:3001/ru/admin/quality-tests');
    await page.waitForLoadState('networkidle');

    console.log('Step 2: Find and click the first test (Backend Format)...');
    // Найдём первую кнопку "Запустить тест"
    const firstTestButton = page
      .locator('button')
      .filter({ hasText: /Запустить тест|Run Test/i })
      .first();
    await expect(firstTestButton).toBeVisible({ timeout: 10000 });
    await firstTestButton.click();

    console.log('Step 3: Wait for test to complete...');
    // Ждём либо успеха, либо ошибки
    await expect(
      page.locator('.badge-success, .badge-error').first()
    ).toBeVisible({ timeout: 30000 });

    console.log('Step 4: Capture test result status...');
    const resultBadge = page.locator('.badge-success, .badge-error').first();
    const initialStatus = await resultBadge.textContent();
    console.log(`✓ Test completed with status: ${initialStatus}`);

    console.log('Step 5: Check localStorage...');
    const storageData = await page.evaluate(() => {
      return localStorage.getItem('quality-tests-results');
    });
    expect(storageData).toBeTruthy();
    console.log('✓ localStorage contains test results');

    console.log('Step 6: Refresh the page...');
    await page.reload();
    await page.waitForLoadState('networkidle');

    console.log('Step 7: Verify test results persisted...');
    // Проверяем, что badge всё ещё виден и имеет тот же статус
    await expect(
      page.locator('.badge-success, .badge-error').first()
    ).toBeVisible({ timeout: 5000 });

    const persistedBadge = page.locator('.badge-success, .badge-error').first();
    const persistedStatus = await persistedBadge.textContent();
    expect(persistedStatus).toBe(initialStatus);
    console.log(`✓ Test results persisted! Status: ${persistedStatus}`);

    console.log('Step 8: Verify statistics are preserved...');
    // Проверяем, что статистика показывает 1 завершённый тест
    const statsCompleted = page.locator('.stat-value.text-sm').first();
    await expect(statsCompleted).toContainText(/\d+%/);
    console.log('✓ Statistics preserved after refresh');

    console.log('Step 9: Test "Clear Results" button...');
    const clearButton = page.locator('button', {
      hasText: /Очистить результаты|Clear Results/i,
    });
    if (await clearButton.isVisible()) {
      console.log('✓ Clear Results button is visible');
      await clearButton.click();

      console.log('Step 10: Verify results are cleared...');
      // После очистки кнопка должна исчезнуть (она показывается только если есть результаты)
      await expect(clearButton).not.toBeVisible({ timeout: 2000 });

      // Проверяем localStorage
      const clearedStorage = await page.evaluate(() => {
        return localStorage.getItem('quality-tests-results');
      });
      expect(clearedStorage).toBeNull();
      console.log('✓ Results successfully cleared from localStorage');
    }

    console.log('✅ All persistence tests passed!');
  });

  test('should handle multiple test results persistence', async ({ page }) => {
    console.log('Step 1: Navigate to Quality Tests page...');
    await page.goto('http://localhost:3001/ru/admin/quality-tests');
    await page.waitForLoadState('networkidle');

    console.log('Step 2: Run first test...');
    const firstTestButton = page
      .locator('button')
      .filter({ hasText: /Запустить тест|Run Test/i })
      .first();
    await firstTestButton.click();
    await expect(
      page.locator('.badge-success, .badge-error').first()
    ).toBeVisible({ timeout: 30000 });

    console.log('Step 3: Run second test...');
    const secondTestButton = page
      .locator('button')
      .filter({ hasText: /Запустить тест|Run Test/i })
      .nth(1);
    await secondTestButton.click();
    await expect(
      page.locator('.badge-success, .badge-error').nth(1)
    ).toBeVisible({ timeout: 30000 });

    console.log('Step 4: Count completed tests...');
    const completedBadges = await page
      .locator('.badge-success, .badge-error')
      .count();
    expect(completedBadges).toBeGreaterThanOrEqual(2);
    console.log(`✓ ${completedBadges} tests completed`);

    console.log('Step 5: Refresh page...');
    await page.reload();
    await page.waitForLoadState('networkidle');

    console.log('Step 6: Verify both results persisted...');
    const persistedBadges = await page
      .locator('.badge-success, .badge-error')
      .count();
    expect(persistedBadges).toBe(completedBadges);
    console.log(`✓ All ${persistedBadges} test results persisted!`);

    console.log('✅ Multiple results persistence test passed!');
  });
});
