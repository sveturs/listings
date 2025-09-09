// Test token storage in browser
// Run this in browser console to test token persistence

console.log('=== Testing Token Storage ===');

// 1. Check current state
console.log('\n1. Current storage state:');
console.log('SessionStorage svetu_* keys:', 
  Object.keys(sessionStorage).filter(k => k.startsWith('svetu_'))
);
console.log('svetu_access_token:', sessionStorage.getItem('svetu_access_token')?.substring(0, 50) + '...');
console.log('svetu_refresh_token:', sessionStorage.getItem('svetu_refresh_token')?.substring(0, 50) + '...');

// 2. Test setting tokens
console.log('\n2. Testing token storage:');
const testAccessToken = 'eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.test_access_token';
const testRefreshToken = 'eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.test_refresh_token';

sessionStorage.setItem('svetu_access_token', testAccessToken);
sessionStorage.setItem('svetu_refresh_token', testRefreshToken);
console.log('Tokens set');

// 3. Check if they persist
console.log('\n3. Checking persistence:');
console.log('svetu_access_token:', sessionStorage.getItem('svetu_access_token')?.substring(0, 50) + '...');
console.log('svetu_refresh_token:', sessionStorage.getItem('svetu_refresh_token')?.substring(0, 50) + '...');

// 4. Check TokenManager
console.log('\n4. TokenManager state:');
if (window.tokenManager) {
  console.log('Access token from TokenManager:', window.tokenManager.getAccessToken()?.substring(0, 50) + '...');
  console.log('Refresh token from TokenManager:', window.tokenManager.getRefreshToken()?.substring(0, 50) + '...');
} else {
  console.log('TokenManager not available in window');
}

// 5. Check for cleanup flag
console.log('\n5. Cleanup status:');
console.log('force_token_cleanup_done:', sessionStorage.getItem('force_token_cleanup_done'));
console.log('token_migration_v2_done:', sessionStorage.getItem('token_migration_v2_done'));

console.log('\n=== Test Complete ===');