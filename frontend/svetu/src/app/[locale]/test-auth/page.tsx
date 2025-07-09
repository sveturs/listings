'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';

export default function TestAuthPage() {
  const router = useRouter();
  const [status, setStatus] = useState<string>('');
  const [error, setError] = useState<string>('');

  const setNonAdminUser = () => {
    try {
      const testUser = {
        id: 99999,
        email: 'test.nonadmin@example.com',
        name: 'Тестовый НЕ администратор',
        is_admin: false,
        created_at: new Date().toISOString(),
      };

      // Сохраняем в sessionStorage
      sessionStorage.setItem('user', JSON.stringify(testUser));

      // Также сохраняем токены для имитации авторизации
      localStorage.setItem('access_token', 'fake-access-token-for-testing');
      localStorage.setItem('refresh_token', 'fake-refresh-token-for-testing');

      setStatus(`✓ Установлен пользователь БЕЗ прав админа: ${testUser.email}`);
      setError('');

      // Показываем данные в консоли
      console.log('[TestAuth] User set:', testUser);
    } catch (err) {
      setError(`Ошибка: ${err}`);
    }
  };

  const setAdminUser = () => {
    try {
      const adminUser = {
        id: 1,
        email: 'admin@example.com',
        name: 'Администратор',
        is_admin: true,
        created_at: new Date().toISOString(),
      };

      sessionStorage.setItem('user', JSON.stringify(adminUser));
      localStorage.setItem('access_token', 'fake-access-token-for-admin');
      localStorage.setItem('refresh_token', 'fake-refresh-token-for-admin');

      setStatus(
        `✓ Установлен пользователь С правами админа: ${adminUser.email}`
      );
      setError('');

      console.log('[TestAuth] Admin set:', adminUser);
    } catch (err) {
      setError(`Ошибка: ${err}`);
    }
  };

  const clearAuth = () => {
    sessionStorage.removeItem('user');
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    setStatus('✓ Все данные авторизации очищены');
    setError('');
    console.log('[TestAuth] Auth cleared');
  };

  const getCurrentUser = () => {
    const userStr = sessionStorage.getItem('user');
    if (userStr) {
      try {
        const user = JSON.parse(userStr);
        return user;
      } catch {
        return null;
      }
    }
    return null;
  };

  const testAdminAccess = () => {
    // Переход на админ-панель
    router.push('/admin/search');
  };

  const currentUser = getCurrentUser();

  return (
    <div className="container mx-auto p-8">
      <h1 className="text-3xl font-bold mb-6">Тестирование AdminGuard</h1>

      <div className="bg-base-200 p-6 rounded-lg mb-6">
        <h2 className="text-xl font-semibold mb-4">Текущий пользователь:</h2>
        {currentUser ? (
          <pre className="bg-base-300 p-4 rounded overflow-x-auto">
            {JSON.stringify(currentUser, null, 2)}
          </pre>
        ) : (
          <p className="text-gray-500">Нет авторизованного пользователя</p>
        )}
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
        <button onClick={setNonAdminUser} className="btn btn-warning">
          Установить обычного пользователя (НЕ админ)
        </button>

        <button onClick={setAdminUser} className="btn btn-success">
          Установить администратора
        </button>

        <button onClick={clearAuth} className="btn btn-error">
          Очистить авторизацию
        </button>

        <button onClick={testAdminAccess} className="btn btn-primary">
          Перейти в админ-панель
        </button>
      </div>

      {status && (
        <div className="alert alert-success mb-4">
          <span>{status}</span>
        </div>
      )}

      {error && (
        <div className="alert alert-error mb-4">
          <span>{error}</span>
        </div>
      )}

      <div className="bg-base-200 p-6 rounded-lg">
        <h3 className="text-lg font-semibold mb-2">Инструкция:</h3>
        <ol className="list-decimal list-inside space-y-2">
          <li>
            Нажмите &quot;Установить обычного пользователя&quot; для создания
            пользователя без прав админа
          </li>
          <li>Нажмите &quot;Перейти в админ-панель&quot;</li>
          <li>
            AdminGuard должен заблокировать доступ и перенаправить на главную
          </li>
          <li>Проверьте консоль браузера (F12) для просмотра логов</li>
        </ol>
      </div>
    </div>
  );
}
