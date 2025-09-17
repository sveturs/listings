'use client';

import { useEffect, useRef } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { tokenManager } from '@/utils/tokenManager';
import { decodeUserFromToken } from '@/utils/jwtDecode';

export default function OAuthProcessor() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const processedRef = useRef(false);

  useEffect(() => {
    if (processedRef.current) {
      return;
    }
    processedRef.current = true;

    // Читаем токен напрямую из URL в клиентском компоненте
    const token = searchParams.get('token');
    const error = searchParams.get('error');

    if (error) {
      router.push(`/?error=${error}`);
      return;
    }

    if (token) {
      tokenManager.setAccessToken(token);

      // Decode and cache user data
      const userData = decodeUserFromToken(token);
      if (userData) {
        try {
          localStorage.setItem('svetu_user', JSON.stringify(userData));
        } catch {
          // Ignore cache errors
        }
      }

      // Пропускаем вызов /auth/session - его нет в backend
      // Cookies должны устанавливаться самим Auth Service при OAuth callback
      console.log(
        '[OAuthProcessor] Token saved, cookies should be set by Auth Service'
      );

      // Wait a bit to ensure token is saved before redirect
      setTimeout(() => {
        // Проверяем и дублируем сохранение для надежности
        const storedToken = localStorage.getItem('svetu_access_token');
        if (!storedToken) {
          localStorage.setItem('svetu_access_token', token);
        }
        router.push('/');
      }, 500); // Небольшая задержка для надежности
    } else {
      router.push('/?error=no_auth_data');
    }
  }, [searchParams, router]);

  const token = searchParams.get('token');
  const error = searchParams.get('error');

  return (
    <div className="flex items-center justify-center min-h-screen bg-base-200">
      <div className="card bg-base-100 shadow-xl p-8 max-w-md w-full">
        <div className="text-center">
          {!error && token ? (
            <>
              <div className="loading loading-spinner loading-lg text-primary mb-4"></div>
              <h2 className="text-2xl font-bold mb-2">Login successful!</h2>
              <p className="text-base-content/70">Redirecting...</p>
            </>
          ) : (
            <>
              <div className="text-error mb-4">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  className="h-16 w-16 mx-auto"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                  />
                </svg>
              </div>
              <h2 className="text-2xl font-bold mb-2 text-error">
                Authentication Failed
              </h2>
              <p className="text-base-content/70 mb-2">
                {error || 'No authorization data received'}
              </p>
              <p className="text-sm text-base-content/50">
                Redirecting to home...
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
