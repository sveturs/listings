// Скрипт для тестирования AdminGuard
// Этот скрипт нужно выполнить в консоли браузера на http://localhost:3001

// 1. Сначала проверим текущего пользователя
console.log('=== Тест AdminGuard ===');
const currentUser = JSON.parse(sessionStorage.getItem('user') || '{}');
console.log('Текущий пользователь:', currentUser);

// 2. Создаем тестового пользователя БЕЗ прав администратора
const testUser = {
  id: 99999,
  email: 'test.nonadmin@example.com',
  name: 'Тестовый НЕ администратор',
  is_admin: false, // Критически важно - пользователь НЕ админ
  created_at: new Date().toISOString(),
};

// 3. Сохраняем в sessionStorage
sessionStorage.setItem('user', JSON.stringify(testUser));
console.log('Установлен тестовый пользователь:', testUser);

// 4. Перезагружаем страницу для применения изменений
console.log('Перезагрузка страницы через 2 секунды...');
setTimeout(() => {
  window.location.reload();
}, 2000);

// После перезагрузки выполните:
// window.location.href = '/admin/search';
// AdminGuard должен заблокировать доступ и перенаправить на главную
