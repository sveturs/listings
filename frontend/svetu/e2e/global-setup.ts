/**
 * Global setup for Playwright tests
 * Performs authentication once and saves the browser state
 * frontend/svetu/e2e/global-setup.ts
 */

import { chromium, FullConfig } from '@playwright/test';
import path from 'path';

const TEST_USER = {
  email: process.env.TEST_ADMIN_EMAIL || 'admin@admin.rs',
  password: process.env.TEST_ADMIN_PASSWORD || 'P@$S4@dmi№',
};

const FRONTEND_URL =
  process.env.NEXT_PUBLIC_FRONTEND_URL || 'http://localhost:3001';
const BACKEND_URL = process.env.BACKEND_INTERNAL_URL || 'http://localhost:3000';

async function globalSetup(config: FullConfig) {
  console.log('🔐 Global Setup: Starting authentication...');

  const browser = await chromium.launch();
  const context = await browser.newContext();
  const page = await context.newPage();

  try {
    // Login via backend API to get JWT token
    console.log('📡 Authenticating via backend API...');

    const loginResponse = await page.request.post(
      `${BACKEND_URL}/api/v1/auth/login`,
      {
        data: {
          email: TEST_USER.email,
          password: TEST_USER.password,
        },
      }
    );

    if (!loginResponse.ok()) {
      const errorText = await loginResponse.text();
      throw new Error(`Login failed: ${loginResponse.status()} ${errorText}`);
    }

    const loginData = await loginResponse.json();
    const accessToken = loginData.access_token || loginData.token;

    if (!accessToken) {
      throw new Error('No access token received from login');
    }

    console.log('✓ JWT token received');

    // Set auth cookie in browser context
    // IMPORTANT: Must use 'access_token' to match what /api/auth/session expects
    await context.addCookies([
      {
        name: 'access_token',
        value: accessToken,
        domain: 'localhost',
        path: '/',
        httpOnly: false,
        secure: false,
        sameSite: 'Lax',
      },
    ]);

    // Save storage state immediately without navigation
    // Navigation causes timeout and closes context, so skip it
    console.log('💾 Saving storage state with auth cookie...');
    const storageStatePath = path.join(__dirname, '..', '.auth', 'user.json');
    await context.storageState({ path: storageStatePath });
    console.log(`✓ Saved auth state with access_token cookie`);

    console.log(`✓ Storage state saved to: ${storageStatePath}`);
    console.log('✅ Global setup complete!');
  } catch (error) {
    console.error('❌ Global setup failed:', error);
    throw error;
  } finally {
    await context.close();
    await browser.close();
  }
}

export default globalSetup;
