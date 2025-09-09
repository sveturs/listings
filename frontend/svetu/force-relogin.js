// Принудительный сброс аутентификации и перелогин
console.log('🔄 Принудительный сброс аутентификации...');

// 1. Очистим все токены из хранилища
console.log('🧹 Очищаем токены...');
Object.keys(sessionStorage).forEach((key) => {
  if (
    key.includes('token') ||
    key.includes('auth') ||
    key.includes('user') ||
    key.includes('svetu')
  ) {
    sessionStorage.removeItem(key);
    console.log(`  Удален из sessionStorage: ${key}`);
  }
});

Object.keys(localStorage).forEach((key) => {
  if (key.includes('token') || key.includes('auth') || key.includes('user')) {
    localStorage.removeItem(key);
    console.log(`  Удален из localStorage: ${key}`);
  }
});

// 2. Очистим cookies
console.log('🍪 Очищаем cookies...');
document.cookie.split(';').forEach((cookie) => {
  const [name] = cookie.trim().split('=');
  if (
    name.includes('token') ||
    name.includes('auth') ||
    name.includes('user') ||
    name.includes('jwt') ||
    name.includes('refresh')
  ) {
    // Очищаем для всех возможных доменов и путей
    document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/`;
    document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=${window.location.hostname}`;
    document.cookie = `${name}=;expires=${new Date(0).toUTCString()};path=/;domain=.${window.location.hostname}`;
    console.log(`  Удален cookie: ${name}`);
  }
});

// 3. Очистим tokenManager если доступен
if (window.tokenManager) {
  console.log('🔧 Очищаем tokenManager...');
  window.tokenManager.clearTokens();
}

// 4. Перенаправим на Google OAuth
console.log('🔑 Перенаправляем на Google OAuth...');
window.location.href = '/api/v1/auth/google';
