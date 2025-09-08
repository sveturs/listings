'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { tokenManager } from '@/utils/tokenManager';
import configManager from '@/config';

export default function AdminAuthPage() {
  const router = useRouter();
  const [loading, setLoading] = useState(false);

  const handleDemoLogin = async () => {
    setLoading(true);

    try {
      // Получаем демо токен от сервера
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(`${apiUrl}/api/v1/admin/demo-token`, {
        method: 'POST',
      });

      if (!response.ok) {
        // Используем статический токен как запасной вариант
        const demoToken = 'demo-admin-token';
        tokenManager.setAccessToken(demoToken);
      } else {
        const data = await response.json();
        if (data.token) {
          tokenManager.setAccessToken(data.token);
        }
      }

      // Редирект на страницу переводов
      setTimeout(() => {
        router.push('/admin/translations');
      }, 500);
    } catch (error) {
      console.error('Auth error:', error);
      // Используем демо токен для тестирования
      const demoToken = 'demo-admin-token';
      tokenManager.setAccessToken(demoToken);

      setTimeout(() => {
        router.push('/admin/translations');
      }, 500);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-base-200">
      <div className="card w-96 bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">Админ панель</h2>

          <div className="alert alert-info mb-4">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="stroke-current shrink-0 w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
            <span>Демо-режим для тестирования системы переводов</span>
          </div>

          <button
            onClick={handleDemoLogin}
            disabled={loading}
            className="btn btn-primary"
          >
            {loading ? (
              <>
                <span className="loading loading-spinner loading-sm"></span>
                Вход...
              </>
            ) : (
              '🔐 Войти как администратор'
            )}
          </button>

          <div className="divider">ИЛИ</div>

          <button
            onClick={() => router.push('/demo-translations')}
            className="btn btn-ghost"
          >
            Демо без авторизации
          </button>
        </div>
      </div>
    </div>
  );
}
