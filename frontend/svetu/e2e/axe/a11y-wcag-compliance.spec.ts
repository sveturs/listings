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
  // Set timeout for each test in this suite
  test.setTimeout(300000); // 5 minutes per test

  test('Homepage should have no accessibility violations', async ({ page }) => {
    await page.goto('/en', { waitUntil: 'domcontentloaded', timeout: 60000 });
    await page.waitForLoadState('load');
    // Wait for main content to be visible
    await page
      .waitForSelector('main, body', { timeout: 10000 })
      .catch(() => {});
    await page.waitForTimeout(2000); // Allow dynamic content to settle

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Marketplace listing page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/marketplace', {
      waitUntil: 'domcontentloaded',
      timeout: 30000,
    });
    await page.waitForLoadState('load');
    // Wait for dynamic content to load
    await page.waitForSelector('main', { timeout: 10000 }).catch(() => {});

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Admin dashboard should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/admin', {
      waitUntil: 'domcontentloaded',
      timeout: 30000,
    });
    await page.waitForLoadState('load');
    await page.waitForSelector('main', { timeout: 10000 }).catch(() => {});

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Admin quality tests page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/admin/quality-tests', {
      waitUntil: 'domcontentloaded',
      timeout: 30000,
    });
    await page.waitForLoadState('load');
    await page.waitForSelector('main', { timeout: 10000 }).catch(() => {});

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Search results page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/search?query=test', {
      waitUntil: 'domcontentloaded',
      timeout: 30000,
    });
    await page.waitForLoadState('load');
    await page.waitForSelector('main', { timeout: 10000 }).catch(() => {});

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Admin categories page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/admin/categories', {
      waitUntil: 'domcontentloaded',
      timeout: 30000,
    });
    await page.waitForLoadState('load');
    await page.waitForSelector('main', { timeout: 10000 }).catch(() => {});

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });
});
