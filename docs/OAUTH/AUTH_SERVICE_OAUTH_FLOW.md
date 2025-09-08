# OAuth Flow with Auth Service Microservice

## Current Architecture

The OAuth flow works through a proxy pattern where the main backend acts as a gateway to the Auth Service:

### 1. OAuth Initiation Flow
1. User clicks "Login with Google" on frontend (port 3001)
2. Frontend calls `AuthService.loginWithGoogle()` which redirects to:
   - `http://localhost:3000/api/v1/auth/oauth/google`
3. Backend proxy (`auth_proxy.go`) forwards the request to Auth Service:
   - `http://localhost:28080/api/v1/auth/oauth/google`
4. Auth Service generates OAuth URL and returns 302 redirect to Google
5. Backend proxy passes the 302 redirect back to browser
6. Browser redirects to Google OAuth consent page

### 2. OAuth Callback Flow
The callback flow is split between backend handling and frontend display:

#### Registered Redirect URI
- **Google Cloud Console**: `http://localhost:3000/auth/google/callback`
- This is the URI registered with Google and cannot be changed without Google Cloud Console access

#### Callback Handling
1. Google redirects to: `http://localhost:3000/auth/google/callback?code=XXX&state=YYY`
2. Backend has legacy handler that processes the callback
3. Backend exchanges code for tokens with Google
4. Backend creates session and redirects to frontend

### 3. Current Issues

#### Issue 1: Redirect URI Mismatch
- **Problem**: Tried to use `http://localhost:3001/ru/auth/oauth/callback` 
- **Error**: Google returns `redirect_uri_mismatch` error
- **Solution**: Must use the registered URI `http://localhost:3000/auth/google/callback`

#### Issue 2: Mixed Legacy and Microservice Code
- **Problem**: Backend has both legacy OAuth handlers and new proxy to Auth Service
- **Confusion**: Two competing OAuth flows exist simultaneously
- **Solution**: Need to fully migrate to Auth Service or keep using legacy

### 4. Recommended Solution

Since the redirect URI `http://localhost:3000/auth/google/callback` is already registered in Google Cloud Console, the best approach is:

1. **Keep the registered redirect URI** as is
2. **Use the backend legacy OAuth handler** for now
3. **Or update the Auth Service** to handle callbacks at the registered URL

### 5. Testing Credentials

```
Email: boxmail386@gmail.com
Password: Dimasik1
```

### 6. Environment Configuration

```env
# Auth Service (.env)
GOOGLE_CLIENT_ID=917315728307-au9ga5fl7o3bbid9nv7e4l92gut194pq.apps.googleusercontent.com
GOOGLE_CLIENT_SECRET=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
GOOGLE_REDIRECT_URL=http://localhost:3000/auth/google/callback
```

### 7. Next Steps

To complete the OAuth flow:

1. **Option A**: Use legacy backend OAuth (already working)
   - Keep using the existing backend OAuth handlers
   - Don't proxy OAuth requests to Auth Service
   
2. **Option B**: Update Auth Service redirect handling
   - Configure Auth Service to accept callbacks at `/auth/google/callback`
   - Update backend proxy to forward these callbacks to Auth Service
   
3. **Option C**: Register new redirect URI in Google Cloud Console
   - Need access to Google Cloud Console
   - Register `http://localhost:3001/ru/auth/oauth/callback`
   - Update all configurations to use the new URI

## Current Status

- ✅ OAuth initiation works (redirects to Google)
- ❌ OAuth callback fails (redirect_uri_mismatch)
- 🔧 Need to align redirect URI across all components