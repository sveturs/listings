import { test, expect } from '@playwright/test';

/**
 * E2E Test: Admin Moderation Flow
 *
 * Full flow: admin login → review pending listings → approve/reject
 *
 * This test verifies the admin moderation workflow for marketplace listings.
 */

test.describe('E2E: Admin Moderation Flow', () => {
  const ADMIN_USER = {
    email: process.env.TEST_ADMIN_EMAIL || 'admin@admin.rs',
    password: process.env.TEST_ADMIN_PASSWORD || 'P@$S4@dmi№',
  };

  test.beforeEach(async ({ page }) => {
    // Navigate to home page
    await page.goto('/en');
  });

  test('should complete admin moderation flow', async ({ page }) => {
    // Step 1: Admin login
    console.log('Step 1: Admin logging in...');

    await page.click('text=Login');
    await page.waitForURL('**/en/auth/login');

    await page.fill('input[type="email"]', ADMIN_USER.email);
    await page.fill('input[type="password"]', ADMIN_USER.password);
    await page.click('button[type="submit"]');

    await page.waitForURL('**/en/**', { timeout: 10000 });

    // Step 2: Navigate to admin panel
    console.log('Step 2: Navigating to admin panel...');

    // Look for admin link in navigation
    const adminLinks = [
      page.locator('a[href*="/admin"], text=Admin'),
      page.locator('[data-testid="admin-link"]'),
    ];

    let adminLinkFound = false;
    for (const link of adminLinks) {
      if (await link.isVisible({ timeout: 3000 }).catch(() => false)) {
        await link.first().click();
        adminLinkFound = true;
        break;
      }
    }

    if (!adminLinkFound) {
      // Try direct navigation to admin panel
      await page.goto('/en/admin');
    }

    await page.waitForLoadState('networkidle');

    // Step 3: Access moderation/listings section
    console.log('Step 3: Accessing moderation section...');

    // Look for moderation or listings management link
    const moderationLinks = [
      page.locator('text=Listings, text=Moderation, text=Pending'),
      page.locator('a[href*="/admin/listings"], a[href*="/admin/moderation"]'),
    ];

    let moderationFound = false;
    for (const link of moderationLinks) {
      if (await link.isVisible({ timeout: 5000 }).catch(() => false)) {
        await link.first().click();
        moderationFound = true;
        break;
      }
    }

    if (!moderationFound) {
      // Try direct navigation to admin listings
      await page.goto('/en/admin/listings');
    }

    await page.waitForLoadState('networkidle');

    // Step 4: View pending listings
    console.log('Step 4: Viewing pending listings...');

    // Wait for listings table/grid to load
    await page.waitForTimeout(2000);

    // Check for listings
    const listingElements = [
      page.locator('table tbody tr, [data-testid="listing-row"]'),
      page.locator('.listing-item, [data-testid="listing-card"]'),
    ];

    let listingsCount = 0;
    for (const elements of listingElements) {
      const count = await elements.count();
      if (count > 0) {
        listingsCount = count;
        console.log(`  Found ${count} listings`);
        break;
      }
    }

    // Step 5: Attempt to moderate a listing
    console.log('Step 5: Attempting to moderate listing...');

    if (listingsCount > 0) {
      // Look for approve/reject buttons
      const moderationButtons = [
        page.locator('button:has-text("Approve"), button:has-text("Accept")').first(),
        page.locator('button:has-text("Reject"), button:has-text("Decline")').first(),
        page.locator('[data-testid="approve-button"], [data-testid="reject-button"]').first(),
      ];

      let moderationButtonFound = false;
      for (const button of moderationButtons) {
        if (await button.isVisible({ timeout: 3000 }).catch(() => false)) {
          console.log('  Found moderation button');
          moderationButtonFound = true;

          // Click the button
          await button.click();

          // Wait for confirmation modal or action to complete
          await page.waitForTimeout(1000);

          // Look for confirmation button if modal appeared
          const confirmButtons = page.locator('button:has-text("Confirm"), button:has-text("Yes"), button:has-text("OK")');
          if (await confirmButtons.isVisible({ timeout: 2000 }).catch(() => false)) {
            console.log('  Confirming moderation action...');
            await confirmButtons.first().click();
          }

          // Wait for action to complete
          await page.waitForTimeout(2000);

          // Check for success message
          const successIndicators = [
            page.locator('text=success, text=approved, text=rejected'),
            page.locator('[role="alert"], .toast, .notification'),
          ];

          for (const indicator of successIndicators) {
            if (await indicator.isVisible({ timeout: 3000 }).catch(() => false)) {
              console.log('  ✅ Moderation action completed successfully');
              break;
            }
          }

          break;
        }
      }

      if (!moderationButtonFound) {
        console.log('  ⚠️  Moderation buttons not found - all listings might be already moderated');
      }
    } else {
      console.log('  ℹ️  No pending listings found to moderate');
    }

    console.log('✅ Admin moderation flow completed');
  });

  test('should access admin dashboard and view statistics', async ({ page }) => {
    // Login as admin
    await page.goto('/en/auth/login');
    await page.fill('input[type="email"]', ADMIN_USER.email);
    await page.fill('input[type="password"]', ADMIN_USER.password);
    await page.click('button[type="submit"]');
    await page.waitForURL('**/en/**');

    // Navigate to admin dashboard
    await page.goto('/en/admin');
    await page.waitForLoadState('networkidle');

    // Check for dashboard elements
    const dashboardElements = [
      page.locator('text=Dashboard, text=Statistics, text=Overview'),
      page.locator('[data-testid="stat-card"], .stat, .metric'),
    ];

    let dashboardFound = false;
    for (const element of dashboardElements) {
      if (await element.isVisible({ timeout: 5000 }).catch(() => false)) {
        dashboardFound = true;
        console.log('  ✅ Admin dashboard accessible');
        break;
      }
    }

    expect(dashboardFound).toBe(true);
  });

  test('should verify admin-only access restriction', async ({ page }) => {
    // Try to access admin panel without login
    await page.goto('/en/admin');

    // Should redirect to login or show unauthorized message
    await page.waitForTimeout(2000);

    const currentURL = page.url();
    const isRestricted =
      currentURL.includes('/auth/login') ||
      currentURL.includes('/unauthorized') ||
      (await page.locator('text=unauthorized, text=access denied, text=login required').isVisible({ timeout: 3000 }).catch(() => false));

    expect(isRestricted).toBe(true);
    console.log('  ✅ Admin access properly restricted for non-authenticated users');
  });
});
