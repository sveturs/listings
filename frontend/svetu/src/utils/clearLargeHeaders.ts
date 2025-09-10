// Utility to clear large data from localStorage that might be sent in headers
export function clearLargeHeaders() {
  if (typeof window === 'undefined') return;

  // Clear any large tokens or data that might be causing header size issues
  const keysToCheck = [
    'access_token',
    'refresh_token',
    'id_token',
    'auth_tokens',
    'auth_state',
    'old_tokens',
    'legacy_tokens',
  ];

  keysToCheck.forEach(key => {
    const value = localStorage.getItem(key);
    if (value && value.length > 4096) {
      console.warn(`[Header Size] Removing large item from localStorage: ${key} (${value.length} bytes)`);
      localStorage.removeItem(key);
    }
  });

  // Also check sessionStorage
  keysToCheck.forEach(key => {
    const value = sessionStorage.getItem(key);
    if (value && value.length > 4096) {
      console.warn(`[Header Size] Removing large item from sessionStorage: ${key} (${value.length} bytes)`);
      sessionStorage.removeItem(key);
    }
  });
}