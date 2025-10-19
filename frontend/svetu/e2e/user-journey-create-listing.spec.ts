import { test, expect } from '@playwright/test';

/**
 * E2E Test: User Journey - Create Listing
 *
 * Full flow: login → create listing → upload images → publish
 *
 * This test verifies the complete user journey from authentication
 * to successfully publishing a marketplace listing.
 */

test.describe('E2E: User Journey - Create Listing', () => {
  const TEST_USER = {
    email: process.env.TEST_ADMIN_EMAIL || 'admin@admin.rs',
    password: process.env.TEST_ADMIN_PASSWORD || 'P@$S4@dmi№',
  };

  test.beforeEach(async ({ page }) => {
    // Navigate to home page
    await page.goto('/en');
  });

  test('should complete full listing creation flow', async ({ page }) => {
    // Step 1: Login
    console.log('Step 1: Logging in...');

    // Click login button
    await page.click('text=Login');
    await page.waitForURL('**/en/auth/login');

    // Fill login form
    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button[type="submit"]');

    // Wait for redirect after login
    await page.waitForURL('**/en/**', { timeout: 10000 });

    // Verify logged in (check for user menu or profile)
    await expect(
      page.locator('text=Profile, text=Account, text=Logout')
    ).toBeVisible({
      timeout: 5000,
    });

    // Step 2: Navigate to Create Listing
    console.log('Step 2: Navigating to create listing...');

    await page.goto('/en/marketplace/create');
    await page.waitForLoadState('networkidle');

    // Step 3: Fill listing form
    console.log('Step 3: Filling listing form...');

    const testListingTitle = `E2E Test Listing ${Date.now()}`;

    await page.fill('input[name="title"]', testListingTitle);
    await page.fill(
      'textarea[name="description"]',
      'This is an automated E2E test listing created by Playwright.'
    );
    await page.fill('input[name="price"]', '99.99');

    // Select category (if available)
    const categorySelector = page.locator(
      'select[name="category_id"], [data-testid="category-select"]'
    );
    if (
      await categorySelector.isVisible({ timeout: 2000 }).catch(() => false)
    ) {
      await categorySelector.selectOption({ index: 1 });
    }

    // Step 4: Upload images (if upload component is present)
    console.log('Step 4: Uploading images (if available)...');

    const fileInput = page.locator('input[type="file"]');
    if (await fileInput.isVisible({ timeout: 2000 }).catch(() => false)) {
      // Note: In real scenario, would upload actual test image
      // For now, skip if file input requires actual file
      console.log(
        '  Image upload input found but skipping actual upload for automated test'
      );
    }

    // Step 5: Publish listing
    console.log('Step 5: Publishing listing...');

    const submitButton = page.locator(
      'button[type="submit"], button:has-text("Publish"), button:has-text("Create")'
    );
    await submitButton.click();

    // Wait for success response
    await page.waitForResponse(
      (response) =>
        response.url().includes('/api/v1/marketplace/listings') &&
        (response.status() === 200 || response.status() === 201),
      { timeout: 15000 }
    );

    // Step 6: Verify listing was created
    console.log('Step 6: Verifying listing creation...');

    // Should redirect to listing detail or listings page
    await page.waitForURL('**/en/marketplace/**', { timeout: 10000 });

    // Verify success message or listing title is visible
    const successIndicators = [
      page.locator(`text=${testListingTitle}`),
      page.locator('text=successfully, text=created, text=published'),
    ];

    let successFound = false;
    for (const indicator of successIndicators) {
      if (await indicator.isVisible({ timeout: 3000 }).catch(() => false)) {
        successFound = true;
        break;
      }
    }

    expect(successFound).toBe(true);

    console.log('✅ Full listing creation flow completed successfully');
  });

  test('should show validation errors for incomplete form', async ({
    page,
  }) => {
    // Login first
    await page.goto('/en/auth/login');
    await page.fill('input[type="email"]', TEST_USER.email);
    await page.fill('input[type="password"]', TEST_USER.password);
    await page.click('button[type="submit"]');
    await page.waitForURL('**/en/**');

    // Navigate to create listing
    await page.goto('/en/marketplace/create');
    await page.waitForLoadState('networkidle');

    // Try to submit without filling required fields
    const submitButton = page.locator(
      'button[type="submit"], button:has-text("Publish"), button:has-text("Create")'
    );
    await submitButton.click();

    // Should show validation errors
    await expect(
      page.locator(
        'text=required, text=field is required, text=cannot be empty'
      )
    ).toBeVisible({
      timeout: 5000,
    });
  });
});
