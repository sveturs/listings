/**
 * Token Migration Utility
 * Detects and cleans up old HS256 tokens after Auth Service migration
 */

export class TokenMigration {
  private static readonly OLD_TOKEN_KEYS = [
    'jwt_token',
    'session_token',
    'refresh_token',
    'access_token',
    'svetu_access_token'
  ];

  /**
   * Check if a JWT token uses HS256 algorithm
   */
  private static isHS256Token(token: string): boolean {
    try {
      const headerB64 = token.split('.')[0];
      if (!headerB64) return false;
      
      const headerJson = atob(headerB64.replace(/-/g, '+').replace(/_/g, '/'));
      const header = JSON.parse(headerJson);
      
      return header.alg === 'HS256';
    } catch {
      // If we can't parse it, assume it's an old token
      return true;
    }
  }

  /**
   * Check if there are any old HS256 tokens in browser storage
   */
  static hasOldTokens(): boolean {
    if (typeof window === 'undefined') return false;

    // Check cookies
    const cookies = document.cookie.split(';');
    for (const cookie of cookies) {
      const [name, value] = cookie.trim().split('=');
      if (this.OLD_TOKEN_KEYS.includes(name) && value) {
        if (this.isHS256Token(value)) {
          console.log(`[TokenMigration] Found old HS256 token in cookie: ${name}`);
          return true;
        }
      }
    }

    // Check localStorage
    for (const key of this.OLD_TOKEN_KEYS) {
      const token = localStorage.getItem(key);
      if (token && this.isHS256Token(token)) {
        console.log(`[TokenMigration] Found old HS256 token in localStorage: ${key}`);
        return true;
      }
    }

    // Check sessionStorage
    for (const key of this.OLD_TOKEN_KEYS) {
      const token = sessionStorage.getItem(key);
      if (token && this.isHS256Token(token)) {
        console.log(`[TokenMigration] Found old HS256 token in sessionStorage: ${key}`);
        return true;
      }
    }

    return false;
  }

  /**
   * Clear all authentication data from browser storage
   */
  static clearAllAuthData(): void {
    if (typeof window === 'undefined') return;

    console.log('[TokenMigration] Clearing all authentication data...');

    // Clear all cookies
    document.cookie.split(';').forEach((c) => {
      const eqPos = c.indexOf('=');
      const name = eqPos > -1 ? c.substring(0, eqPos).trim() : c.trim();
      
      // Clear auth-related cookies
      if (this.OLD_TOKEN_KEYS.includes(name) || 
          name === 'user_id' || 
          name === 'user_email' ||
          name === 'session_id') {
        // Clear for all possible domains
        document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/`;
        document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=${window.location.hostname}`;
        document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=.${window.location.hostname}`;
        
        // Also try localhost variations
        if (window.location.hostname === 'localhost') {
          document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=localhost`;
          document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=.localhost`;
        }
      }
    });

    // Clear localStorage auth data
    const localStorageKeysToRemove = [
      ...this.OLD_TOKEN_KEYS,
      'user',
      'user_id',
      'user_email',
      'auth_state'
    ];
    
    localStorageKeysToRemove.forEach(key => {
      localStorage.removeItem(key);
    });

    // Clear sessionStorage auth data
    const sessionStorageKeysToRemove = [
      ...this.OLD_TOKEN_KEYS,
      'svetu_user',
      'client_id',
      'oauth_return_to',
      'auth_state'
    ];
    
    sessionStorageKeysToRemove.forEach(key => {
      sessionStorage.removeItem(key);
    });

    console.log('[TokenMigration] Authentication data cleared');
  }

  /**
   * Run the migration check and cleanup if needed
   * Returns true if migration was performed
   */
  static runMigration(): boolean {
    if (typeof window === 'undefined') return false;

    // Check if we've already run migration in this session
    const migrationFlag = sessionStorage.getItem('token_migration_v2_done');
    if (migrationFlag === 'true') {
      return false;
    }

    if (this.hasOldTokens()) {
      console.log('[TokenMigration] Old HS256 tokens detected, performing migration...');
      
      // Clear all old auth data
      this.clearAllAuthData();
      
      // Set flag to prevent multiple migrations in same session
      sessionStorage.setItem('token_migration_v2_done', 'true');
      
      console.log('[TokenMigration] Migration completed. User needs to re-authenticate.');
      
      return true;
    }

    // Set flag even if no migration needed to avoid repeated checks
    sessionStorage.setItem('token_migration_v2_done', 'true');
    
    return false;
  }

  /**
   * Check if user needs to re-authenticate after migration
   */
  static needsReauthentication(): boolean {
    if (typeof window === 'undefined') return false;

    // Check if we have any valid RS256 tokens
    const checkToken = (token: string | null): boolean => {
      if (!token) return false;
      
      try {
        const headerB64 = token.split('.')[0];
        if (!headerB64) return false;
        
        const headerJson = atob(headerB64.replace(/-/g, '+').replace(/_/g, '/'));
        const header = JSON.parse(headerJson);
        
        return header.alg === 'RS256';
      } catch {
        return false;
      }
    };

    // Check various storage locations for valid RS256 tokens
    const hasValidToken = 
      checkToken(sessionStorage.getItem('svetu_access_token')) ||
      checkToken(localStorage.getItem('access_token'));

    return !hasValidToken;
  }
}

// Auto-run migration on module load (only in browser)
if (typeof window !== 'undefined') {
  // Run migration check immediately and after a short delay
  const migrated = TokenMigration.runMigration();
  if (migrated) {
    console.log('[TokenMigration] Automatic migration completed (immediate)');
  }
  
  // Also run after a delay to catch any late-loaded tokens
  setTimeout(() => {
    const migrated2 = TokenMigration.runMigration();
    if (migrated2) {
      console.log('[TokenMigration] Automatic migration completed (delayed)');
    }
  }, 100);
}