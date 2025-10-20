import { defineConfig, devices } from '@playwright/test';

/**
 * See https://playwright.dev/docs/test-configuration.
 */
export default defineConfig({
  testDir: './e2e',
  /* Run tests in files in parallel */
  fullyParallel: false,
  /* Fail the build on CI if you accidentally left test.only in the source code. */
  forbidOnly: !!process.env.CI,
  /* Retry on CI only */
  retries: process.env.CI ? 2 : 0,
  /* Opt out of parallel tests on CI. */
  workers: process.env.CI ? 1 : 1,
  /* Reporter to use. See https://playwright.dev/docs/test-reporters */
  reporter: [
    ['html', { outputFolder: 'playwright-report' }],
    ['json', { outputFile: 'test-results/results.json' }],
    ['list'],
  ],
  /* Global setup - runs once before all tests to establish authentication */
  globalSetup: require.resolve('./e2e/global-setup.ts'),

  /* Maximum time test can run - increased for accessibility tests */
  timeout: 300000, // 5 minutes per test (accessibility tests need more time for axe scans)

  /* Shared settings for all the projects below. See https://playwright.dev/docs/api/class-testoptions. */
  use: {
    /* Base URL to use in actions like `await page.goto('/')`. */
    baseURL: process.env.NEXT_PUBLIC_FRONTEND_URL || 'http://localhost:3001',

    /* Collect trace when retrying the failed test. See https://playwright.dev/docs/trace-viewer */
    trace: 'on-first-retry',

    /* Screenshot on failure */
    screenshot: 'only-on-failure',

    /* Video on failure */
    video: 'retain-on-failure',

    /* Maximum time each action can take */
    actionTimeout: 30000, // Increased for slower pages during accessibility scans

    /* Headless mode - always true for CI, configurable for local */
    headless: process.env.HEADLESS !== 'false',
  },

  /* Configure projects for major browsers */
  projects: [
    {
      name: 'chromium',
      use: {
        ...devices['Desktop Chrome'],
        /* Use saved authentication state from global setup */
        storageState: '.auth/user.json',
        /* Browser launch options for headless mode */
        launchOptions: {
          args: [
            '--no-sandbox',
            '--disable-setuid-sandbox',
            '--disable-dev-shm-usage',
            '--disable-gpu',
          ],
        },
      },
    },
  ],

  /* Run your local dev server before starting the tests */
  // DISABLED: E2E tests are triggered via backend API, frontend must already be running
  // webServer: process.env.CI
  //   ? undefined
  //   : {
  //       command: 'yarn dev',
  //       url: 'http://localhost:3001',
  //       reuseExistingServer: !process.env.CI,
  //       timeout: 120000,
  //     },
});
