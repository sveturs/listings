/**
 * Accessibility Test: WCAG 2.1 AA Compliance
 * Tests pages for WCAG compliance using axe-core
 * frontend/svetu/e2e/axe/a11y-wcag-compliance.spec.ts
 */

import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

const TEST_ADMIN_EMAIL = process.env.TEST_ADMIN_EMAIL || 'admin@admin.rs';
const TEST_ADMIN_PASSWORD = process.env.TEST_ADMIN_PASSWORD || 'P@$S4@dmi№';

test.describe('WCAG 2.1 AA Compliance Tests', () => {
  test.beforeEach(async ({ page }) => {
    // User is already logged in via global-setup
    // Just wait a bit for the page to be ready
    await page.waitForTimeout(500);
  });

  test('Homepage should have no accessibility violations', async ({ page }) => {
    await page.goto('/en');

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Marketplace listing page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/marketplace', { waitUntil: 'domcontentloaded' });
    await page.waitForLoadState('load');
    await page.waitForTimeout(1000);

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Admin dashboard should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/admin', { waitUntil: 'domcontentloaded' });
    await page.waitForLoadState('load');
    await page.waitForTimeout(1000);

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Admin quality tests page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/admin/quality-tests', { waitUntil: 'domcontentloaded' });
    await page.waitForLoadState('load');
    await page.waitForTimeout(1000);

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Search results page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/search?query=test', { waitUntil: 'domcontentloaded' });
    await page.waitForLoadState('load');
    await page.waitForTimeout(1000);

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Admin categories page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/admin/categories', { waitUntil: 'domcontentloaded' });
    await page.waitForLoadState('load');
    await page.waitForTimeout(1000);

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });
});
