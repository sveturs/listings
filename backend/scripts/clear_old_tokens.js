// Script to clear old HS256 tokens from browser storage
// This should be embedded in the frontend to automatically clean up old tokens

(function() {
    console.log('[Token Migration] Starting cleanup of old HS256 tokens...');
    
    // Check if we have old tokens
    const checkForOldTokens = () => {
        // Check cookies
        const cookies = document.cookie.split(';');
        let hasOldTokens = false;
        
        for (const cookie of cookies) {
            const [name, value] = cookie.trim().split('=');
            if (name === 'jwt_token' || name === 'session_token' || name === 'refresh_token') {
                if (value) {
                    try {
                        // Decode JWT header to check algorithm
                        const headerB64 = value.split('.')[0];
                        const header = JSON.parse(atob(headerB64));
                        if (header.alg === 'HS256') {
                            console.log(`[Token Migration] Found old HS256 token in cookie: ${name}`);
                            hasOldTokens = true;
                        }
                    } catch (e) {
                        // Invalid token format, should be cleaned up anyway
                        hasOldTokens = true;
                    }
                }
            }
        }
        
        // Check localStorage
        const localToken = localStorage.getItem('access_token');
        if (localToken) {
            try {
                const headerB64 = localToken.split('.')[0];
                const header = JSON.parse(atob(headerB64));
                if (header.alg === 'HS256') {
                    console.log('[Token Migration] Found old HS256 token in localStorage');
                    hasOldTokens = true;
                }
            } catch (e) {
                hasOldTokens = true;
            }
        }
        
        // Check sessionStorage
        const sessionToken = sessionStorage.getItem('svetu_access_token');
        if (sessionToken) {
            try {
                const headerB64 = sessionToken.split('.')[0];
                const header = JSON.parse(atob(headerB64));
                if (header.alg === 'HS256') {
                    console.log('[Token Migration] Found old HS256 token in sessionStorage');
                    hasOldTokens = true;
                }
            } catch (e) {
                hasOldTokens = true;
            }
        }
        
        return hasOldTokens;
    };
    
    // Clear all authentication data
    const clearAuthData = () => {
        console.log('[Token Migration] Clearing all authentication data...');
        
        // Clear cookies
        document.cookie.split(';').forEach(function(c) { 
            const eqPos = c.indexOf('=');
            const name = eqPos > -1 ? c.substring(0, eqPos).trim() : c.trim();
            // Clear auth-related cookies
            if (name === 'jwt_token' || name === 'session_token' || name === 'refresh_token' || 
                name === 'user_id' || name === 'user_email') {
                document.cookie = name + '=;expires=' + new Date(0).toUTCString() + ';path=/';
                document.cookie = name + '=;expires=' + new Date(0).toUTCString() + ';path=/;domain=' + window.location.hostname;
                document.cookie = name + '=;expires=' + new Date(0).toUTCString() + ';path=/;domain=.' + window.location.hostname;
            }
        });
        
        // Clear localStorage auth data
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('user');
        
        // Clear sessionStorage auth data
        sessionStorage.removeItem('svetu_access_token');
        sessionStorage.removeItem('svetu_user');
        sessionStorage.removeItem('client_id');
        
        console.log('[Token Migration] Authentication data cleared');
    };
    
    // Main migration logic
    if (checkForOldTokens()) {
        console.log('[Token Migration] Old HS256 tokens detected, cleaning up...');
        clearAuthData();
        
        // Show notification to user
        if (typeof window !== 'undefined' && window.alert) {
            alert('Your session has expired. Please log in again to continue.');
        }
        
        // Redirect to login page
        window.location.href = '/auth/login?message=session_expired';
    } else {
        console.log('[Token Migration] No old tokens found, system is up to date');
    }
})();