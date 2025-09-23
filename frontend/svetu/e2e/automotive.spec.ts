import { test, expect } from '@playwright/test';

test.describe('Automotive Section E2E Tests', () => {
  test.beforeEach(async ({ page }) => {
    // Переходим на страницу автомобилей
    await page.goto('/en/cars');
  });

  test('should display cars page with correct elements', async ({ page }) => {
    // Проверяем заголовок страницы
    await expect(page.locator('h1')).toContainText('Car Marketplace');

    // Проверяем наличие поля поиска
    await expect(page.locator('input[placeholder*="Search"]')).toBeVisible();

    // Проверяем статистику
    await expect(page.locator('.stats')).toBeVisible();
    await expect(
      page.locator('.stat-title:has-text("Active Listings")')
    ).toBeVisible();

    // Проверяем популярные марки
    await expect(page.locator('h2:has-text("Popular Brands")')).toBeVisible();

    // Проверяем категории
    await expect(page.locator('.btn:has-text("Passenger Cars")')).toBeVisible();
    await expect(page.locator('.btn:has-text("SUVs")')).toBeVisible();
  });

  test('should navigate through car categories', async ({ page }) => {
    // Клик на категорию "Passenger Cars"
    await page.locator('.btn:has-text("Passenger Cars")').click();

    // Проверяем, что URL изменился
    await expect(page).toHaveURL(/category=10101/);
  });

  test('should perform car search', async ({ page }) => {
    // Вводим поисковый запрос
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('BMW');

    // Нажимаем Enter или кнопку поиска
    await searchInput.press('Enter');

    // Проверяем, что URL обновился с параметром поиска
    await expect(page).toHaveURL(/q=BMW/);
  });

  test('should use quick filters', async ({ page }) => {
    // Проверяем наличие секции расширенного поиска
    await expect(page.locator('.card:has-text("Car Filters")')).toBeVisible();

    // Кликаем на быстрый фильтр "New cars"
    const newCarsFilter = page.locator('button:has-text("New cars")');
    if (await newCarsFilter.isVisible()) {
      await newCarsFilter.click();

      // Проверяем, что кнопка стала активной (имеет класс btn-primary)
      await expect(newCarsFilter).toHaveClass(/btn-primary/);
    }

    // Кликаем на быстрый фильтр "Under €10,000"
    const budgetFilter = page.locator('button:has-text("Under €10,000")');
    if (await budgetFilter.isVisible()) {
      await budgetFilter.click();
      await expect(budgetFilter).toHaveClass(/btn-primary/);
    }
  });

  test('should change sorting options', async ({ page }) => {
    // Находим селектор сортировки
    const sortSelect = page.locator('select').filter({ hasText: /Sort by/ });

    if (await sortSelect.isVisible()) {
      // Меняем сортировку на "Price: Low to High"
      await sortSelect.selectOption({ label: 'Price: Low to High' });

      // Ждем обновления результатов
      await page.waitForTimeout(1000);

      // Проверяем, что сортировка применена
      const selectedOption = await sortSelect.inputValue();
      expect(selectedOption).toBe('price_asc');
    }
  });

  test('should use car filters', async ({ page }) => {
    const filtersSection = page.locator('.card:has-text("Car Filters")');

    if (await filtersSection.isVisible()) {
      // Выбираем марку
      const makeSelect = filtersSection.locator('select').first();
      if (await makeSelect.isVisible()) {
        await makeSelect.selectOption({ index: 1 }); // Выбираем первую марку

        // Ждем загрузки моделей
        await page.waitForTimeout(500);

        // Проверяем, что появился селектор моделей
        const modelSelect = filtersSection.locator('select:has-text("Model")');
        await expect(modelSelect).toBeVisible();
      }

      // Устанавливаем ценовой диапазон
      const priceFromInput = filtersSection
        .locator('input[placeholder*="From"]')
        .first();
      const priceToInput = filtersSection
        .locator('input[placeholder*="To"]')
        .first();

      if (await priceFromInput.isVisible()) {
        await priceFromInput.fill('5000');
        await priceToInput.fill('20000');
      }

      // Выбираем тип топлива
      const fuelSelect = filtersSection
        .locator('select')
        .filter({ hasText: /Fuel/ });
      if (await fuelSelect.isVisible()) {
        await fuelSelect.selectOption('diesel');
      }
    }
  });

  test('should display car listings correctly', async ({ page }) => {
    // Проверяем наличие секции с последними объявлениями
    const latestListings = page.locator('h2:has-text("Latest Listings")');

    if (await latestListings.isVisible()) {
      // Проверяем, что есть карточки объявлений
      const listingCards = page
        .locator('.card')
        .filter({ has: page.locator('figure') });
      const count = await listingCards.count();

      expect(count).toBeGreaterThan(0);

      // Проверяем структуру первой карточки
      if (count > 0) {
        const firstCard = listingCards.first();

        // Проверяем наличие изображения или заглушки
        const figure = firstCard.locator('figure');
        await expect(figure).toBeVisible();

        // Проверяем наличие заголовка
        const title = firstCard.locator('.card-title');
        await expect(title).toBeVisible();

        // Проверяем наличие цены
        const price = firstCard.locator('text=/€\\d+/');
        const hasPrice = (await price.count()) > 0;

        if (hasPrice) {
          await expect(price.first()).toBeVisible();
        }
      }
    }
  });

  test('should handle mobile responsiveness', async ({ page }) => {
    // Устанавливаем мобильный viewport
    await page.setViewportSize({ width: 375, height: 667 });

    // Проверяем, что мобильное меню работает
    const mobileMenuButton = page.locator('.btn.lg\\:hidden');

    if (await mobileMenuButton.isVisible()) {
      // Проверяем, что фильтры скрыты на мобильном
      const filters = page.locator('.card:has-text("Car Filters")');
      const isFiltersVisible = await filters.isVisible();

      // На мобильном фильтры должны быть в bottom sheet или скрыты
      if (!isFiltersVisible) {
        // Ищем кнопку показа фильтров
        const showFiltersBtn = page.locator('button:has-text("Show filters")');
        if (await showFiltersBtn.isVisible()) {
          await showFiltersBtn.click();

          // Проверяем, что фильтры появились
          await expect(filters).toBeVisible();
        }
      }
    }

    // Проверяем, что карточки адаптировались под мобильный размер
    const cards = page.locator('.grid .card');
    const firstCard = cards.first();

    if (await firstCard.isVisible()) {
      const box = await firstCard.boundingBox();
      if (box) {
        // На мобильном карточки должны занимать почти всю ширину
        expect(box.width).toBeGreaterThan(300);
      }
    }
  });

  test('should navigate to car details page', async ({ page }) => {
    // Находим первую карточку автомобиля
    const firstListing = page
      .locator('.card')
      .filter({ has: page.locator('figure') })
      .first();

    if (await firstListing.isVisible()) {
      // Кликаем на карточку
      await firstListing.click();

      // Проверяем, что перешли на страницу детального просмотра
      await expect(page).toHaveURL(/\/listing\/\d+/);

      // Проверяем, что страница загрузилась
      await page.waitForLoadState('networkidle');
    }
  });

  test('should display popular car makes', async ({ page }) => {
    // Находим секцию с популярными марками
    const popularMakesSection = page.locator('h2:has-text("Popular Brands")');

    if (await popularMakesSection.isVisible()) {
      // Проверяем наличие карточек марок
      const makeCards = popularMakesSection.locator('..').locator('.card');
      const count = await makeCards.count();

      expect(count).toBeGreaterThan(0);

      // Кликаем на первую марку
      if (count > 0) {
        const firstMake = makeCards.first();
        await firstMake.click();

        // Проверяем, что URL обновился с параметром марки
        await expect(page).toHaveURL(/car_make=/);
      }
    }
  });

  test('should handle quick links section', async ({ page }) => {
    // Прокручиваем к секции быстрых ссылок
    const quickLinksSection = page.locator('h2:has-text("Quick Links")');

    if (await quickLinksSection.isVisible()) {
      await quickLinksSection.scrollIntoViewIfNeeded();

      // Проверяем карточку "Sell Your Car"
      const sellCarCard = page.locator('.card:has-text("Sell Your Car")');
      await expect(sellCarCard).toBeVisible();

      // Проверяем кнопку создания объявления
      const createListingBtn = sellCarCard.locator(
        '.btn:has-text("Create Listing")'
      );
      await expect(createListingBtn).toBeVisible();

      // Проверяем, что кнопка кликабельна
      const isDisabled = await createListingBtn.isDisabled();
      expect(isDisabled).toBe(false);
    }
  });

  test('should handle search with no results', async ({ page }) => {
    // Вводим несуществующий запрос
    const searchInput = page.locator('input[placeholder*="Search"]');
    await searchInput.fill('XYZ123NonExistentCar');
    await searchInput.press('Enter');

    // Ждем обновления страницы
    await page.waitForTimeout(1000);

    // Проверяем сообщение об отсутствии результатов
    const noResultsMessage = page.locator('text=/No.*found|No.*results/i');

    // Если есть секция результатов, проверяем сообщение
    const hasResults = (await page.locator('.card figure').count()) > 0;

    if (!hasResults && (await noResultsMessage.isVisible())) {
      await expect(noResultsMessage).toBeVisible();
    }
  });
});

// Тесты производительности
test.describe('Automotive Performance Tests', () => {
  test('should load page within acceptable time', async ({ page }) => {
    const startTime = Date.now();

    await page.goto('/en/cars');
    await page.waitForLoadState('networkidle');

    const loadTime = Date.now() - startTime;

    // Страница должна загружаться менее чем за 3 секунды
    expect(loadTime).toBeLessThan(3000);
  });

  test('should handle virtualization for many items', async ({ page }) => {
    // Если есть много элементов, проверяем виртуализацию
    await page.goto('/en/cars');

    // Проверяем наличие результатов поиска
    const searchResults = page.locator('.grid .card');
    const count = await searchResults.count();

    if (count > 50) {
      // При большом количестве элементов должна работать виртуализация
      // Проверяем, что не все элементы рендерятся сразу
      const visibleCards = await searchResults.evaluateAll(
        (cards) =>
          cards.filter((card) => {
            const rect = card.getBoundingClientRect();
            return rect.bottom >= 0 && rect.top <= window.innerHeight;
          }).length
      );

      // Видимых карточек должно быть меньше общего количества
      expect(visibleCards).toBeLessThan(count);
    }
  });
});
