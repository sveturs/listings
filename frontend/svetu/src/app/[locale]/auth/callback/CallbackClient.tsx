'use client';

import { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { tokenManager } from '@/utils/tokenManager';
import { useAuth } from '@/contexts/AuthContext';

export default function CallbackClient() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { refreshSession } = useAuth();

  useEffect(() => {
    const handleCallback = async () => {
      // Получаем токен из URL (поддерживаем оба варианта: auth_token и token)
      const authToken = searchParams?.get('auth_token') || searchParams?.get('token');
      const returnUrl = searchParams?.get('returnUrl') || searchParams?.get('state') || '/';

      if (authToken) {
        console.log('[AuthCallback] Received auth token from OAuth callback');
        // Сохраняем токен
        tokenManager.setAccessToken(authToken);

        // Обновляем сессию для загрузки данных пользователя
        await refreshSession();
      }

      // Редиректим на нужную страницу
      router.push(returnUrl);
    };

    handleCallback();
  }, [searchParams, router, refreshSession]);

  return (
    <div className="flex min-h-screen items-center justify-center">
      <div className="text-center">
        <div className="loading loading-spinner loading-lg mb-4"></div>
        <p>Completing authentication...</p>
      </div>
    </div>
  );
}
