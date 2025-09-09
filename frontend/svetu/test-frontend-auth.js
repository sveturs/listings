// Test frontend authentication flow
// This simulates what the React app does

const API_BASE = 'http://localhost:3000';

async function testLogin() {
  console.log('=== Testing Frontend Auth Flow ===\n');
  
  // 1. Login
  console.log('1. Attempting login...');
  const loginResponse = await fetch(`${API_BASE}/api/v1/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      email: 'test@example.com',
      password: 'password123'
    }),
    credentials: 'include'
  });
  
  const loginData = await loginResponse.json();
  
  if (!loginResponse.ok) {
    console.error('Login failed:', loginData);
    return;
  }
  
  console.log('✓ Login successful');
  console.log('  Access token:', loginData.access_token?.substring(0, 50) + '...');
  console.log('  Refresh token:', loginData.refresh_token?.substring(0, 50) + '...');
  
  // 2. Save tokens like TokenManager does
  console.log('\n2. Saving tokens to sessionStorage...');
  if (typeof sessionStorage !== 'undefined') {
    sessionStorage.setItem('svetu_access_token', loginData.access_token);
    sessionStorage.setItem('svetu_refresh_token', loginData.refresh_token);
    console.log('✓ Tokens saved');
  } else {
    console.log('Running in Node.js - sessionStorage not available');
  }
  
  // 3. Test session endpoint
  console.log('\n3. Testing session endpoint...');
  const sessionResponse = await fetch(`${API_BASE}/api/v1/auth/session`, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${loginData.access_token}`
    },
    credentials: 'include'
  });
  
  const sessionData = await sessionResponse.json();
  
  if (sessionResponse.ok && sessionData.data?.authenticated) {
    console.log('✓ Session valid');
    console.log('  User:', sessionData.data.user?.email);
  } else {
    console.log('✗ Session invalid');
    console.log('  Response:', sessionData);
  }
  
  // 4. Test refresh
  console.log('\n4. Testing refresh token...');
  const refreshResponse = await fetch(`${API_BASE}/api/v1/auth/refresh`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': `Bearer ${loginData.refresh_token}`
    },
    body: JSON.stringify({
      refresh_token: loginData.refresh_token
    }),
    credentials: 'include'
  });
  
  const refreshData = await refreshResponse.json();
  
  if (refreshResponse.ok) {
    console.log('✓ Refresh successful');
    console.log('  New access token:', refreshData.access_token?.substring(0, 50) + '...');
    console.log('  New refresh token:', refreshData.refresh_token?.substring(0, 50) + '...');
  } else {
    console.log('✗ Refresh failed');
    console.log('  Response:', refreshData);
  }
  
  console.log('\n=== Test Complete ===');
}

// Run the test
testLogin().catch(console.error);