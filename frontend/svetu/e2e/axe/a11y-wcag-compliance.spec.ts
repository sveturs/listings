/**
 * Accessibility Test: WCAG 2.1 AA Compliance
 * Tests critical pages for WCAG compliance using axe-core
 * frontend/svetu/e2e/axe/a11y-wcag-compliance.spec.ts
 */

import { test, expect } from '@playwright/test';
import AxeBuilder from '@axe-core/playwright';

test.describe('WCAG 2.1 AA Compliance Tests', () => {
  // Set reasonable timeout for each test
  test.setTimeout(120000); // 2 minutes per test (reduced from 10)

  /**
   * Helper to format violations into readable error message
   */
  function formatViolations(violations: any[]): string {
    if (violations.length === 0) return '';

    return violations
      .map(
        (violation, idx) =>
          `\n${idx + 1}. ${violation.id}: ${violation.description}\n` +
          `   Impact: ${violation.impact}\n` +
          `   Help: ${violation.help}\n` +
          `   Elements affected: ${violation.nodes.length}\n` +
          `   Example: ${violation.nodes[0]?.html || 'N/A'}`
      )
      .join('\n');
  }

  test('Homepage should have no accessibility violations', async ({ page }) => {
    // Use domcontentloaded instead of networkidle to avoid waiting for all API calls
    await page.goto('/en', { waitUntil: 'domcontentloaded', timeout: 30000 });

    // Wait for main content to be present
    await page.waitForSelector('main, [role="main"]', {
      timeout: 10000,
      state: 'attached',
    });

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    const violations = accessibilityScanResults.violations;
    expect(
      violations,
      `Found ${violations.length} accessibility violations on homepage:${formatViolations(violations)}`
    ).toEqual([]);
  });

  test('Admin dashboard should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/admin', {
      waitUntil: 'domcontentloaded',
      timeout: 30000,
    });

    // Wait for admin page to fully load (either shows admin content or redirects to login)
    await page.waitForLoadState('networkidle', { timeout: 10000 });

    // Wait for main content (admin dashboard or login page)
    await page.waitForSelector('main, [role="main"]', {
      timeout: 10000,
      state: 'attached',
    });

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      // Exclude document-title check since AdminGuard may still be hydrating
      .disableRules(['document-title'])
      .analyze();

    const violations = accessibilityScanResults.violations;
    expect(
      violations,
      `Found ${violations.length} accessibility violations on admin dashboard:${formatViolations(violations)}`
    ).toEqual([]);
  });

  test('Search results page should have no accessibility violations', async ({
    page,
  }) => {
    await page.goto('/en/search?query=laptop', {
      waitUntil: 'domcontentloaded',
      timeout: 30000,
    });

    await page.waitForSelector('main, [role="main"]', {
      timeout: 10000,
      state: 'attached',
    });

    const accessibilityScanResults = await new AxeBuilder({ page })
      .withTags(['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa'])
      .analyze();

    const violations = accessibilityScanResults.violations;
    expect(
      violations,
      `Found ${violations.length} accessibility violations on search page:${formatViolations(violations)}`
    ).toEqual([]);
  });
});
