// Отладка токенов в браузере
console.log('🔍 Отладка токенов:');

// Проверим sessionStorage
console.log('📦 sessionStorage:');
Object.keys(sessionStorage).forEach((key) => {
  if (
    key.includes('token') ||
    key.includes('auth') ||
    key.includes('user') ||
    key.includes('svetu')
  ) {
    const value = sessionStorage.getItem(key);
    console.log(`  ${key}: ${value ? value.substring(0, 50) + '...' : 'null'}`);

    // Если это токен, попробуем его декодировать
    if (value && value.includes('.')) {
      try {
        const [headerB64] = value.split('.');
        const header = JSON.parse(
          atob(headerB64.replace(/-/g, '+').replace(/_/g, '/'))
        );
        console.log(`    Алгоритм: ${header.alg}, Тип: ${header.typ}`);
      } catch (e) {
        console.log(`    Не удалось декодировать: ${e.message}`);
      }
    }
  }
});

// Проверим localStorage
console.log('📦 localStorage:');
Object.keys(localStorage).forEach((key) => {
  if (key.includes('token') || key.includes('auth') || key.includes('user')) {
    const value = localStorage.getItem(key);
    console.log(`  ${key}: ${value ? value.substring(0, 50) + '...' : 'null'}`);
  }
});

// Проверим cookies
console.log('🍪 cookies:');
document.cookie.split(';').forEach((cookie) => {
  const [name, value] = cookie.trim().split('=');
  if (
    name.includes('token') ||
    name.includes('auth') ||
    name.includes('user') ||
    name.includes('jwt')
  ) {
    console.log(
      `  ${name}: ${value ? value.substring(0, 50) + '...' : 'null'}`
    );
  }
});

// Если есть tokenManager, проверим его
if (window.tokenManager) {
  console.log('🔧 TokenManager:');
  const token = window.tokenManager.getAccessToken();
  console.log(
    `  Текущий токен: ${token ? token.substring(0, 50) + '...' : 'null'}`
  );

  if (token) {
    console.log(`  Токен истек: ${window.tokenManager.isTokenExpired()}`);
  }
}
