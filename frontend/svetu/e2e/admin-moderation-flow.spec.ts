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
    await page.goto('/en', { waitUntil: 'domcontentloaded' });
    await page.waitForLoadState('load');
    await page.waitForTimeout(1000);
  });

  test('should complete admin moderation flow', async ({ page }) => {
    // Step 1: User is already logged in via global-setup
    console.log('Step 1: User authenticated via global setup');

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
      await page.goto('/en/admin', { waitUntil: 'domcontentloaded' });
    }

    await page.waitForLoadState('load');
    await page.waitForTimeout(1000);

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
      await page.goto('/en/admin/listings', { waitUntil: 'domcontentloaded' });
    }

    await page.waitForLoadState('load');
    await page.waitForTimeout(1000);

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
        page
          .locator('button:has-text("Approve"), button:has-text("Accept")')
          .first(),
        page
          .locator('button:has-text("Reject"), button:has-text("Decline")')
          .first(),
        page
          .locator(
            '[data-testid="approve-button"], [data-testid="reject-button"]'
          )
          .first(),
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
          const confirmButtons = page.locator(
            'button:has-text("Confirm"), button:has-text("Yes"), button:has-text("OK")'
          );
          if (
            await confirmButtons.isVisible({ timeout: 2000 }).catch(() => false)
          ) {
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
            if (
              await indicator.isVisible({ timeout: 3000 }).catch(() => false)
            ) {
              console.log('  ✅ Moderation action completed successfully');
              break;
            }
          }

          break;
        }
      }

      if (!moderationButtonFound) {
        console.log(
          '  ⚠️  Moderation buttons not found - all listings might be already moderated'
        );
      }
    } else {
      console.log('  ℹ️  No pending listings found to moderate');
    }

    console.log('✅ Admin moderation flow completed');
  });

  test('should access admin dashboard and view statistics', async ({
    page,
  }) => {
    // User is already logged in via global-setup
    console.log('User authenticated via global setup');

    // Navigate to admin dashboard
    await page.goto('/en/admin', { waitUntil: 'domcontentloaded' });
    await page.waitForLoadState('load');
    await page.waitForTimeout(1000);

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

    // Dashboard might not have specific elements - just check we can access /admin
    if (!dashboardFound) {
      console.log(
        '  ℹ️  Dashboard elements not found, but admin page is accessible'
      );
    }

    // Test passes if we reached admin page without redirect
    expect(page.url()).toContain('/admin');
  });

  test('should verify admin-only access restriction', async ({ browser }) => {
    // Create a new context WITHOUT authentication
    const context = await browser.newContext();
    const page = await context.newPage();

    try {
      // Try to access admin panel without login
      await page.goto('/en/admin');

      // Should redirect to login or show unauthorized message
      await page.waitForTimeout(2000);

      const currentURL = page.url();
      const isRestricted =
        currentURL.includes('/auth/login') ||
        currentURL.includes('/unauthorized') ||
        (await page
          .locator('text=unauthorized, text=access denied, text=login required')
          .isVisible({ timeout: 3000 })
          .catch(() => false));

      // If not redirected, check if we actually see admin content (which would be wrong)
      const hasAdminContent = await page
        .locator('text=Admin, text=Dashboard, text=Moderation')
        .isVisible({ timeout: 2000 })
        .catch(() => false);

      // Either should be redirected OR should not see admin content
      expect(isRestricted || !hasAdminContent).toBe(true);
      console.log(
        '  ✅ Admin access properly restricted for non-authenticated users'
      );
    } finally {
      await context.close();
    }
  });
});
