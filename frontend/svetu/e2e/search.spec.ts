import { test, expect } from '@playwright/test';

test.describe('Search Page E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/en/search');
  });

  test('should load search page with all components', async ({ page }) => {
    // Check main components are visible
    await expect(page.locator('input[placeholder*="Search"]')).toBeVisible();
    await expect(page.getByText('Categories')).toBeVisible();
    await expect(page.getByText('Filters')).toBeVisible();
    await expect(page.getByText('Location')).toBeVisible();
    await expect(page.getByText('Condition')).toBeVisible();
  });

  test('should perform basic search', async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="Search"]');

    // Type search query
    await searchInput.fill('laptop');
    await searchInput.press('Enter');

    // Wait for results
    await page.waitForResponse(
      (response) =>
        response.url().includes('/api/v1/marketplace/search') &&
        response.status() === 200
    );

    // Check results are displayed
    await expect(page.getByText(/Found \d+ results/)).toBeVisible();
  });

  test('should show autocomplete suggestions', async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="Search"]');

    // Start typing
    await searchInput.fill('iph');

    // Wait for suggestions to appear
    await page.waitForSelector('[role="listbox"]', { timeout: 5000 });

    // Check suggestions are visible
    const suggestions = page.locator('[role="option"]');
    await expect(suggestions).toHaveCount(3, { timeout: 5000 });
  });

  test('should filter by category', async ({ page }) => {
    // Wait for categories to load
    await page.waitForSelector('text=Electronics', { timeout: 10000 });

    // Expand Electronics category
    const electronicsExpander = page
      .locator('button')
      .filter({ hasText: /^$/ })
      .near(page.getByText('Electronics'))
      .first();
    await electronicsExpander.click();

    // Wait for subcategories
    await page.waitForSelector('text=Smartphones', { timeout: 5000 });

    // Select Smartphones category
    const smartphonesCheckbox = page
      .locator('input[type="checkbox"]')
      .near(page.getByText('Smartphones'));
    await smartphonesCheckbox.check();

    // Verify filter is applied
    await expect(page.getByText('Clear (1)')).toBeVisible();
  });

  test('should filter by location', async ({ page }) => {
    const locationInput = page.locator('input[placeholder*="Enter city"]');

    // Type city name
    await locationInput.fill('Novi Sad');

    // Wait for suggestions (if any)
    await page.waitForTimeout(500);

    // Select first suggestion if available
    const citySuggestion = page.locator('text=Novi Sad').first();
    if (await citySuggestion.isVisible()) {
      await citySuggestion.click();
    }

    // Check radius slider
    const radiusSlider = page.locator('input[type="range"]');
    await expect(radiusSlider).toBeVisible();

    // Change radius
    await radiusSlider.fill('25');
    await expect(page.getByText('25 km')).toBeVisible();
  });

  test('should filter by condition', async ({ page }) => {
    // Click on condition dropdown/select
    const conditionSelect = page
      .locator('select')
      .filter({ has: page.locator('option[value="new"]') });

    if (await conditionSelect.isVisible()) {
      await conditionSelect.selectOption('new');
    } else {
      // Alternative: click on radio button
      const newCondition = page.locator('input[type="radio"][value="new"]');
      await newCondition.check();
    }

    // Verify filter is applied
    await page.waitForTimeout(500);
  });

  test('should clear all filters', async ({ page }) => {
    // Apply some filters first
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('test');
    await searchInput.press('Enter');

    // Wait for results
    await page.waitForTimeout(1000);

    // Clear filters button should be visible
    const clearButton = page.getByText('Clear filters');
    if (await clearButton.isVisible()) {
      await clearButton.click();

      // Verify filters are cleared
      await expect(searchInput).toHaveValue('');
    }
  });

  test('should expand and collapse all categories', async ({ page }) => {
    // Wait for categories to load
    await page.waitForSelector('text=Electronics', { timeout: 10000 });

    // Click expand all button
    const expandAllButton = page.locator('button[title="Expand all"]');
    await expandAllButton.click();

    // Check subcategories are visible
    await expect(page.getByText('Smartphones')).toBeVisible();
    await expect(page.getByText('Laptops')).toBeVisible();

    // Click collapse all button
    const collapseAllButton = page.locator('button[title="Collapse all"]');
    await collapseAllButton.click();

    // Check subcategories are hidden
    await expect(page.getByText('Smartphones')).not.toBeVisible();
    await expect(page.getByText('Laptops')).not.toBeVisible();
  });

  test('should navigate through search suggestions with keyboard', async ({
    page,
  }) => {
    const searchInput = page.locator('input[placeholder*="Search"]');

    // Type to trigger suggestions
    await searchInput.fill('phone');

    // Wait for suggestions
    await page.waitForSelector('[role="listbox"]', { timeout: 5000 });

    // Navigate with arrow keys
    await searchInput.press('ArrowDown');
    await searchInput.press('ArrowDown');

    // Select with Enter
    await searchInput.press('Enter');

    // Verify search was performed
    await page.waitForResponse(
      (response) =>
        response.url().includes('/api/v1/marketplace/search') &&
        response.status() === 200
    );
  });

  test('should show no results message', async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="Search"]');

    // Search for something unlikely to have results
    await searchInput.fill('xyzabc123456789');
    await searchInput.press('Enter');

    // Wait for response
    await page.waitForResponse(
      (response) =>
        response.url().includes('/api/v1/marketplace/search') &&
        response.status() === 200
    );

    // Check no results message
    await expect(page.getByText(/No results|No Results/)).toBeVisible();
  });

  test('should maintain filters after search', async ({ page }) => {
    // Select a category first
    await page.waitForSelector('text=Electronics', { timeout: 10000 });
    const electronicsCheckbox = page
      .locator('input[type="checkbox"]')
      .near(page.getByText('Electronics'));
    await electronicsCheckbox.check();

    // Perform search
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('laptop');
    await searchInput.press('Enter');

    // Wait for results
    await page.waitForResponse(
      (response) =>
        response.url().includes('/api/v1/marketplace/search') &&
        response.status() === 200
    );

    // Verify category is still selected
    await expect(electronicsCheckbox).toBeChecked();
  });

  test('should update URL with search parameters', async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="Search"]');

    // Perform search
    await searchInput.fill('laptop');
    await searchInput.press('Enter');

    // Check URL contains search query
    await expect(page).toHaveURL(/.*[?&]q=laptop/);
  });

  test('should show loading state during search', async ({ page }) => {
    const searchInput = page.locator('input[placeholder*="Search"]');

    // Start search
    await searchInput.fill('test');
    const searchPromise = searchInput.press('Enter');

    // Check for loading indicator (skeleton or spinner)
    const loadingIndicator = page
      .locator('.skeleton, .loading, [role="status"]')
      .first();

    // Wait for either loading to appear or search to complete
    await Promise.race([
      loadingIndicator
        .waitFor({ state: 'visible', timeout: 1000 })
        .catch(() => {}),
      searchPromise,
    ]);

    // Wait for results
    await page.waitForResponse(
      (response) =>
        response.url().includes('/api/v1/marketplace/search') &&
        response.status() === 200
    );
  });

  test('should handle API errors gracefully', async ({ page }) => {
    // Intercept API call and return error
    await page.route('**/api/v1/marketplace/search*', (route) => {
      route.fulfill({
        status: 500,
        body: JSON.stringify({ error: 'Internal Server Error' }),
      });
    });

    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('test');
    await searchInput.press('Enter');

    // Check error message is displayed
    await expect(page.getByText(/error|Error|failed|Failed/i)).toBeVisible({
      timeout: 5000,
    });
  });
});

test.describe('Search Page Mobile Tests', () => {
  test.use({ viewport: { width: 375, height: 667 } });

  test('should be responsive on mobile', async ({ page }) => {
    await page.goto('/en/search');

    // Check search input is visible
    await expect(page.locator('input[placeholder*="Search"]')).toBeVisible();

    // Filters might be in a collapsible menu on mobile
    const filterButton = page.getByText('Filters');
    if (await filterButton.isVisible()) {
      await filterButton.click();
      await expect(page.getByText('Categories')).toBeVisible();
    }
  });

  test('should handle touch interactions', async ({ page }) => {
    await page.goto('/en/search');

    // Tap on search input
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.tap();

    // Type using virtual keyboard
    await searchInput.fill('mobile test');

    // Submit search
    await searchInput.press('Enter');

    // Wait for results
    await page
      .waitForResponse(
        (response) =>
          response.url().includes('/api/v1/marketplace/search') &&
          response.status() === 200,
        { timeout: 10000 }
      )
      .catch(() => {});
  });
});
