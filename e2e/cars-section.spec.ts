import { test, expect } from '@playwright/test';

test.describe('Automotive Section', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the cars page
    await page.goto('http://localhost:3001/en/cars');
  });

  test('should display cars page with proper heading', async ({ page }) => {
    // Check page title
    await expect(page).toHaveTitle(/Cars/);

    // Check main heading
    const heading = page.getByRole('heading', { name: /Cars Marketplace/i });
    await expect(heading).toBeVisible();
  });

  test('should load and display car makes', async ({ page }) => {
    // Wait for car makes to load
    await page.waitForSelector('[data-testid="car-makes-list"]', { timeout: 10000 });

    // Check that at least one make is displayed
    const carMakes = page.locator('[data-testid="car-make-item"]');
    await expect(carMakes).toHaveCount({ min: 1 });
  });

  test('should filter cars by make', async ({ page }) => {
    // Wait for filters to load
    await page.waitForSelector('[data-testid="car-filters"]', { timeout: 10000 });

    // Select BMW from make dropdown
    const makeSelect = page.locator('select[name="car-make"]');
    await makeSelect.selectOption('bmw');

    // Wait for results to update
    await page.waitForTimeout(1000);

    // Check that filtered results are displayed
    const results = page.locator('[data-testid="car-listing-card"]');
    const count = await results.count();

    // If there are results, they should all be BMW
    if (count > 0) {
      const firstResult = results.first();
      await expect(firstResult).toContainText(/BMW/i);
    }
  });

  test('should filter cars by year range', async ({ page }) => {
    // Set year from
    const yearFrom = page.locator('input[name="year-from"]');
    await yearFrom.fill('2015');

    // Set year to
    const yearTo = page.locator('input[name="year-to"]');
    await yearTo.fill('2020');

    // Apply filters
    const applyButton = page.getByRole('button', { name: /Apply Filters/i });
    await applyButton.click();

    // Wait for results to update
    await page.waitForTimeout(1000);

    // Check that results are within the year range
    const yearTexts = await page.locator('[data-testid="car-year"]').allTextContents();
    yearTexts.forEach(yearText => {
      const year = parseInt(yearText);
      expect(year).toBeGreaterThanOrEqual(2015);
      expect(year).toBeLessThanOrEqual(2020);
    });
  });

  test('should toggle between grid and list view', async ({ page }) => {
    // Default should be grid view
    const gridView = page.locator('[data-testid="grid-view"]');
    await expect(gridView).toHaveClass(/active/);

    // Switch to list view
    const listViewButton = page.getByRole('button', { name: /List View/i });
    await listViewButton.click();

    // Check that view changed
    const listView = page.locator('[data-testid="list-view"]');
    await expect(listView).toHaveClass(/active/);

    // Check that cards are displayed in list format
    const cards = page.locator('[data-testid="car-listing-card"]');
    const firstCard = cards.first();
    await expect(firstCard).toHaveClass(/list-format/);
  });

  test('should navigate to car detail page', async ({ page }) => {
    // Wait for listings to load
    await page.waitForSelector('[data-testid="car-listing-card"]', { timeout: 10000 });

    // Click on the first car listing
    const firstListing = page.locator('[data-testid="car-listing-card"]').first();
    await firstListing.click();

    // Should navigate to detail page
    await expect(page).toHaveURL(/\/listing\/\d+/);

    // Detail page should have car information
    const carTitle = page.getByRole('heading', { level: 1 });
    await expect(carTitle).toBeVisible();
  });

  test('should show car quick filters', async ({ page }) => {
    // Check for quick filter buttons
    const newCarsFilter = page.getByRole('button', { name: /New Cars/i });
    await expect(newCarsFilter).toBeVisible();

    const under10kFilter = page.getByRole('button', { name: /Under.*10,000/i });
    await expect(under10kFilter).toBeVisible();

    const lowMileageFilter = page.getByRole('button', { name: /Low Mileage/i });
    await expect(lowMileageFilter).toBeVisible();
  });

  test('should reset filters', async ({ page }) => {
    // Apply some filters
    const makeSelect = page.locator('select[name="car-make"]');
    await makeSelect.selectOption('bmw');

    const yearFrom = page.locator('input[name="year-from"]');
    await yearFrom.fill('2015');

    // Click reset button
    const resetButton = page.getByRole('button', { name: /Reset Filters/i });
    await resetButton.click();

    // Check that filters are cleared
    await expect(makeSelect).toHaveValue('');
    await expect(yearFrom).toHaveValue('');
  });

  test('should display car statistics', async ({ page }) => {
    // Check for statistics section
    const statsSection = page.locator('[data-testid="car-statistics"]');
    await expect(statsSection).toBeVisible();

    // Check for total listings count
    const totalListings = page.locator('[data-testid="total-car-listings"]');
    await expect(totalListings).toContainText(/\d+/);

    // Check for popular makes
    const popularMakes = page.locator('[data-testid="popular-makes"]');
    await expect(popularMakes).toBeVisible();
  });

  test('should handle empty search results gracefully', async ({ page }) => {
    // Search for something that won't return results
    const makeSelect = page.locator('select[name="car-make"]');
    await makeSelect.selectOption('non-existent-make');

    // Wait for results to update
    await page.waitForTimeout(1000);

    // Should show no results message
    const noResults = page.getByText(/No cars found/i);
    await expect(noResults).toBeVisible();
  });
});

test.describe('Car VIN Decoder', () => {
  test('should navigate to VIN decoder page', async ({ page }) => {
    await page.goto('http://localhost:3001/en/cars');

    // Find and click VIN decoder link
    const vinDecoderLink = page.getByRole('link', { name: /VIN Decoder/i });
    await vinDecoderLink.click();

    // Should navigate to VIN decoder page
    await expect(page).toHaveURL(/\/cars\/vin-decoder/);

    // Page should have VIN input field
    const vinInput = page.locator('input[placeholder*="VIN"]');
    await expect(vinInput).toBeVisible();
  });

  test('should validate VIN format', async ({ page }) => {
    await page.goto('http://localhost:3001/en/cars/vin-decoder');

    // Enter invalid VIN (too short)
    const vinInput = page.locator('input[placeholder*="VIN"]');
    await vinInput.fill('ABC123');

    // Click decode button
    const decodeButton = page.getByRole('button', { name: /Decode/i });
    await decodeButton.click();

    // Should show error message
    const errorMessage = page.getByText(/Invalid VIN/i);
    await expect(errorMessage).toBeVisible();
  });
});

test.describe('Integration with Map', () => {
  test('should show car filters on map for automotive categories', async ({ page }) => {
    // Go to map page
    await page.goto('http://localhost:3001/en/map');

    // Select automotive category (10100)
    const categorySelect = page.locator('[data-testid="category-selector"]');
    await categorySelect.selectOption('10100');

    // Wait for filters to update
    await page.waitForTimeout(1000);

    // Car-specific filters should appear
    const carMakeFilter = page.locator('select[name="car-make"]');
    await expect(carMakeFilter).toBeVisible();

    const yearFilter = page.locator('input[name="year-from"]');
    await expect(yearFilter).toBeVisible();
  });
});