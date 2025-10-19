/**
 * Accessibility Test: WCAG 2.1 AA Compliance
 * Tests pages for WCAG compliance using axe-core
 * frontend/svetu/e2e/axe/a11y-wcag-compliance.spec.ts
 */

import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

const TEST_ADMIN_EMAIL = process.env.TEST_ADMIN_EMAIL || 'admin@admin.rs';
const TEST_ADMIN_PASSWORD = process.env.TEST_ADMIN_PASSWORD || 'P@$S4@dmi№';
const BASE_URL = process.env.BASE_URL || 'http://localhost:3001';

test.describe('WCAG 2.1 AA Compliance Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Login as admin before each test
    await page.goto(`${BASE_URL}/en/login`);
    await page.fill('input[type="email"]', TEST_ADMIN_EMAIL);
    await page.fill('input[type="password"]', TEST_ADMIN_PASSWORD);
    await page.click('button[type="submit"]');

    // Wait for successful login
    await page.waitForURL(/\/(en|ru|sr)\/(marketplace|admin)/, {
      timeout: 10000,
    });
  });

  test('Homepage should have no accessibility violations', async ({ page }) => {
    await page.goto(`${BASE_URL}/en`);

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Marketplace listing page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/en/marketplace`);
    await page.waitForLoadState('networkidle');

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Admin dashboard should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/en/admin`);
    await page.waitForLoadState('networkidle');

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Admin quality tests page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/en/admin/quality-tests`);
    await page.waitForLoadState('networkidle');

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Search results page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/en/search?query=test`);
    await page.waitForLoadState('networkidle');

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });

  test('Admin categories page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto(`${BASE_URL}/en/admin/categories`);
    await page.waitForLoadState('networkidle');

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    expect(accessibilityScanResults.violations).toEqual([]);
  });
});
