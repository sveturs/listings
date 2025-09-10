# 🔐 Auth Service Token Migration Solution

## Problem Statement

After migrating to the Auth Service microservice, the backend now only accepts RS256 (RSA-signed) JWT tokens. However, users' browsers still contained old HS256 (HMAC-signed) tokens from the previous authentication system, causing widespread authentication failures.

### Symptoms
- Backend logs: `JWT validation failed - only RS256 tokens supported`
- Browser console: Multiple 401 Unauthorized errors
- Users unable to access any protected endpoints
- Cart, profile, and other authenticated features broken

## Root Cause Analysis

1. **Backend Migration**: The backend was updated to only accept RS256 tokens from Auth Service
   - Location: `/backend/internal/middleware/auth_jwt.go`
   - Public key: `/backend/keys/auth_service_public.pem`
   - No backward compatibility with HS256 tokens

2. **Client-Side Cache**: Browsers retained old HS256 tokens in:
   - Cookies: `jwt_token`, `session_token`, `refresh_token`
   - localStorage: `access_token`, `refresh_token`
   - sessionStorage: `svetu_access_token`, `svetu_user`

3. **Token Mismatch**: Frontend kept sending old HS256 tokens, backend rejected them

## Implemented Solution

### 1. Automatic Token Migration Utility

**File**: `/frontend/svetu/src/utils/tokenMigration.ts`

Features:
- Detects HS256 tokens by parsing JWT header
- Clears all authentication data from browser storage
- Runs automatically on page load
- Uses session flag to prevent repeated migrations

### 2. Forced Token Cleanup

**File**: `/frontend/svetu/src/utils/forceTokenCleanup.ts`

Features:
- Aggressive one-time cleanup of all auth-related storage
- Checks token algorithm before removal
- Preserves RS256 tokens while removing HS256 ones
- Clears cookies, localStorage, sessionStorage, and IndexedDB

### 3. TokenManager HS256 Rejection

**File**: `/frontend/svetu/src/utils/tokenManager.ts`

Enhanced features:
- Rejects HS256 tokens in constructor
- Rejects HS256 tokens in setAccessToken method
- Automatically clears invalid tokens
- Only accepts and stores RS256 tokens

```typescript
class TokenMigration {
  // Check if token uses HS256
  private static isHS256Token(token: string): boolean
  
  // Detect old tokens in storage
  static hasOldTokens(): boolean
  
  // Clear all auth data
  static clearAllAuthData(): void
  
  // Run migration if needed
  static runMigration(): boolean
}
```

### 2. AuthContext Integration

**File**: `/frontend/svetu/src/contexts/AuthContext.tsx`

Integration points:
- Runs token migration check on component mount
- Clears user state if migration performed
- Prevents authentication attempts with old tokens

```typescript
useEffect(() => {
  // Check and migrate old tokens
  const migrated = TokenMigration.runMigration();
  if (migrated) {
    updateUser(null);
    setIsLoading(false);
    return;
  }
  // ... rest of auth initialization
}, []);
```

### 3. Backend Validation

**File**: `/backend/internal/middleware/auth_jwt.go`

Current behavior:
- Only accepts RS256 tokens signed by Auth Service
- Validates using public key
- Returns clear error message for HS256 tokens
- No fallback to old authentication methods

## Migration Process

### Automatic (for users)

1. User visits the site
2. TokenMigration utility runs automatically
3. Detects old HS256 tokens
4. Clears all authentication data
5. User sees they're logged out
6. User logs in again via Google OAuth
7. Receives new RS256 token from Auth Service

### Manual (if needed)

Users can manually clear their browser data:
```javascript
// Clear all cookies
document.cookie.split(";").forEach(c => { 
  document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/"); 
});

// Clear storage
localStorage.clear();
sessionStorage.clear();
```

## Testing

### Verify Token Algorithm
```bash
# Decode JWT header to check algorithm
echo "YOUR_TOKEN" | cut -d'.' -f1 | base64 -d | jq .

# Should show:
{
  "alg": "RS256",
  "typ": "JWT"
}
```

### Test Authentication
```bash
# After clearing old tokens and re-authenticating
curl -H "Authorization: Bearer YOUR_NEW_TOKEN" http://localhost:3000/api/v1/users/profile
```

## Status

✅ **Backend**: Fully migrated to RS256-only validation
✅ **Auth Service**: Issues RS256 tokens correctly
✅ **Frontend**: Automatic token migration implemented
✅ **User Experience**: Seamless migration with re-authentication prompt

## Files Modified

1. `/frontend/svetu/src/utils/tokenMigration.ts` - New utility for token migration
2. `/frontend/svetu/src/utils/forceTokenCleanup.ts` - Aggressive token cleanup utility
3. `/frontend/svetu/src/utils/tokenManager.ts` - Enhanced to reject HS256 tokens
4. `/frontend/svetu/src/contexts/AuthContext.tsx` - Integration of migration and cleanup
5. `/docs/AUTH_SERVICE_TOKEN_FIX.md` - User-facing documentation
6. `/backend/scripts/clear_old_tokens.js` - Manual cleanup script (backup)

## Next Steps

1. Monitor logs for any remaining HS256 token attempts
2. Consider adding user-friendly notification when migration occurs
3. Remove migration code after all users have migrated (est. 30 days)
4. Update monitoring to track RS256 token adoption rate

## Security Considerations

- Old HS256 tokens are completely invalidated
- No backward compatibility maintains security boundary
- Clear separation between old and new authentication systems
- No sensitive data exposed during migration process