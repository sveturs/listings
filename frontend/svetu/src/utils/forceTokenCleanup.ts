/**
 * Force Token Cleanup
 * One-time aggressive cleanup of all old authentication tokens
 * This ensures complete removal of HS256 tokens after Auth Service migration
 */

export function forceTokenCleanup(): boolean {
  if (typeof window === 'undefined') return false;

  console.log('[ForceTokenCleanup] Starting aggressive token cleanup...');

  let cleanedSomething = false;

  // Clear all cookies that might contain tokens
  const cookiesToClear = [
    'jwt_token',
    'session_token',
    'refresh_token',
    'access_token',
    'user_id',
    'user_email',
    'session_id',
    'auth_token',
  ];

  document.cookie.split(';').forEach((cookie) => {
    const eqPos = cookie.indexOf('=');
    const name = eqPos > -1 ? cookie.substring(0, eqPos).trim() : cookie.trim();

    if (
      cookiesToClear.includes(name) ||
      name.includes('token') ||
      name.includes('auth')
    ) {
      // Clear for all possible domains
      document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/`;
      document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=${window.location.hostname}`;
      document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=.${window.location.hostname}`;

      if (window.location.hostname === 'localhost') {
        document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=localhost`;
        document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=.localhost`;
        document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=127.0.0.1`;
      }

      console.log(`[ForceTokenCleanup] Cleared cookie: ${name}`);
      cleanedSomething = true;
    }
  });

  // Clear localStorage items
  const localStorageKeys = Object.keys(localStorage);
  localStorageKeys.forEach((key) => {
    if (
      key.includes('token') ||
      key.includes('auth') ||
      key.includes('user') ||
      key === 'access_token' ||
      key === 'refresh_token' ||
      key === 'jwt_token'
    ) {
      localStorage.removeItem(key);
      console.log(`[ForceTokenCleanup] Cleared localStorage: ${key}`);
      cleanedSomething = true;
    }
  });

  // Clear sessionStorage items
  const sessionStorageKeys = Object.keys(sessionStorage);
  sessionStorageKeys.forEach((key) => {
    // Don't clear the migration flag or cleanup flag
    if (key === 'token_migration_v2_done' || key === 'force_token_cleanup_done')
      return;

    // Always keep svetu_* keys - they are for Auth Service
    if (key.startsWith('svetu_')) {
      console.log(`[ForceTokenCleanup] Keeping Auth Service key: ${key}`);
      return;
    }

    // Keep client_id for rate limiting
    if (key === 'client_id') {
      console.log(`[ForceTokenCleanup] Keeping client_id for rate limiting`);
      return;
    }

    // Only remove keys that contain auth-related strings
    if (
      key.includes('token') ||
      key.includes('auth') ||
      key.includes('user') ||
      key.includes('jwt') ||
      key.includes('session')
    ) {
      // Check if it's a JWT token
      const value = sessionStorage.getItem(key);
      if (value && value.includes('.')) {
        try {
          const headerB64 = value.split('.')[0];
          const headerJson = atob(
            headerB64.replace(/-/g, '+').replace(/_/g, '/')
          );
          const header = JSON.parse(headerJson);

          if (header.alg === 'HS256') {
            sessionStorage.removeItem(key);
            console.log(
              `[ForceTokenCleanup] Cleared HS256 token from sessionStorage: ${key}`
            );
            cleanedSomething = true;
          } else if (header.alg === 'RS256') {
            // Even RS256 tokens should be removed if they're not svetu_* keys
            // because we only trust svetu_* namespace for Auth Service
            sessionStorage.removeItem(key);
            console.log(
              `[ForceTokenCleanup] Cleared non-svetu RS256 token: ${key}`
            );
            cleanedSomething = true;
          }
        } catch {
          // Not a valid JWT, remove it
          sessionStorage.removeItem(key);
          console.log(
            `[ForceTokenCleanup] Cleared invalid token from sessionStorage: ${key}`
          );
          cleanedSomething = true;
        }
      } else {
        // Remove non-token auth data
        sessionStorage.removeItem(key);
        console.log(`[ForceTokenCleanup] Cleared sessionStorage: ${key}`);
        cleanedSomething = true;
      }
    }
  });

  // Clear any indexed DB auth data
  if (window.indexedDB) {
    const dbNames = ['auth', 'tokens', 'session'];
    dbNames.forEach((dbName) => {
      try {
        indexedDB.deleteDatabase(dbName);
        console.log(`[ForceTokenCleanup] Deleted IndexedDB: ${dbName}`);
        // eslint-disable-next-line @typescript-eslint/no-unused-vars
      } catch (_e) {
        // Ignore errors
      }
    });
  }

  if (cleanedSomething) {
    console.log('[ForceTokenCleanup] Cleanup completed - old tokens removed');

    // Set a flag to indicate cleanup was performed
    sessionStorage.setItem(
      'force_token_cleanup_done',
      new Date().toISOString()
    );
  } else {
    console.log('[ForceTokenCleanup] No old tokens found');
  }

  return cleanedSomething;
}

// ОТКЛЮЧЕНО: автоматическая очистка при загрузке модуля
// Это удаляло валидные токены после OAuth авторизации
// forceTokenCleanup() должна вызываться только при явной необходимости
if (typeof window !== 'undefined' && process.env.NODE_ENV === 'development') {
  console.log(
    '[ForceTokenCleanup] Module loaded, auto-cleanup DISABLED to preserve OAuth tokens'
  );
  console.log(
    '[ForceTokenCleanup] Call forceTokenCleanup() manually when needed'
  );
}
