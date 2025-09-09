'use client';

import { useEffect, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { tokenManager } from '@/utils/tokenManager';

export default function GoogleCallbackPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [status, setStatus] = useState('Processing login...');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleOAuthCallback = async () => {
      try {
        // Get parameters from the current URL
        const urlParams = new URLSearchParams(window.location.search);
        const code = urlParams.get('code');
        const state = urlParams.get('state');
        const oauthError = urlParams.get('error');

        console.log('[OAuth Callback] Params from URL:', {
          code,
          state,
          oauthError,
        });
        console.log('[OAuth Callback] Full URL:', window.location.href);

        // Handle OAuth error
        if (oauthError) {
          console.error('[OAuth Callback] OAuth error:', oauthError);
          setError('Authorization failed. Please try again.');
          return;
        }

        // If no code, try searchParams as fallback
        let finalCode = code;
        let finalState = state;

        if (!code) {
          console.log('[OAuth Callback] No code in URL, trying searchParams');
          finalCode = searchParams.get('code');
          finalState = searchParams.get('state');
          const spError = searchParams.get('error');

          if (spError) {
            console.error(
              '[OAuth Callback] OAuth error from searchParams:',
              spError
            );
            setError('Authorization failed. Please try again.');
            return;
          }
        }

        if (!finalCode) {
          console.error('[OAuth Callback] No authorization code found');
          setError('No authorization code received. Please try again.');
          return;
        }

        console.log('[OAuth Callback] Processing authorization code...');
        setStatus('Exchanging authorization code for tokens...');

        // Build the backend callback URL - используем относительный путь для proxy
        const backendCallbackUrl = `/api/v1/auth/google/callback?code=${encodeURIComponent(finalCode)}${
          finalState ? `&state=${encodeURIComponent(finalState)}` : ''
        }`;

        console.log('[OAuth Callback] Calling backend:', backendCallbackUrl);

        // Call backend to complete OAuth flow and get tokens
        const response = await fetch(backendCallbackUrl, {
          method: 'GET',
          credentials: 'include', // Important for cookies
          redirect: 'manual', // Не следовать за редиректами автоматически
          headers: {
            Accept: 'application/json',
          },
        });

        console.log(
          '[OAuth Callback] Backend response status:',
          response.status
        );

        // Обработка редиректа от Auth Service (302/307)
        if (
          response.type === 'opaqueredirect' ||
          response.status === 302 ||
          response.status === 307
        ) {
          console.log(
            '[OAuth Callback] Got redirect response from Auth Service'
          );

          // При редиректе Auth Service должен установить cookies
          // Даём небольшую задержку и пробуем refresh
          await new Promise((resolve) => setTimeout(resolve, 500));

          // Пробуем получить токен через refresh
          const refreshResponse = await fetch('/api/v1/auth/refresh', {
            method: 'POST',
            credentials: 'include',
            headers: {
              'Content-Type': 'application/json',
            },
          });

          if (refreshResponse.ok) {
            const refreshData = await refreshResponse.json();
            console.log(
              '[OAuth Callback] Got tokens via refresh after redirect:',
              refreshData
            );

            if (refreshData.access_token) {
              tokenManager.setAccessToken(refreshData.access_token);
              if (refreshData.refresh_token) {
                tokenManager.setRefreshToken(refreshData.refresh_token);
              }

              // Сохраняем пользователя в sessionStorage
              if (refreshData.user) {
                sessionStorage.setItem(
                  'svetu_user',
                  JSON.stringify(refreshData.user)
                );
              }

              // Триггерим событие для немедленного обновления AuthContext
              if (typeof window !== 'undefined') {
                window.dispatchEvent(
                  new CustomEvent('tokenChanged', {
                    detail: {
                      token: refreshData.access_token,
                      action: 'set',
                    },
                  })
                );
              }

              setStatus('Login successful! Redirecting...');

              // Немедленный редирект с параметром для показа уведомления
              router.replace('/?login=success');
              return;
            }
          }

          // Если refresh не сработал, пробуем извлечь токен из Location заголовка
          const location = response.headers.get('Location');
          if (location) {
            const urlParams = new URLSearchParams(new URL(location).search);
            const token = urlParams.get('token');
            if (token) {
              console.log('[OAuth Callback] Got token from redirect location');
              tokenManager.setAccessToken(token);
              setStatus('Login successful! Redirecting...');
              setTimeout(() => {
                router.replace('/');
              }, 1000);
              return;
            }
          }
        }

        if (
          !response.ok &&
          response.status !== 302 &&
          response.status !== 307
        ) {
          const errorText = await response.text();
          console.error('[OAuth Callback] Backend error:', errorText);
          setError('Authentication failed. Please try again.');
          return;
        }

        // Handle Auth Service response
        console.log(
          '[OAuth Callback] Backend response status:',
          response.status
        );
        console.log(
          '[OAuth Callback] Response headers:',
          Object.fromEntries(response.headers.entries())
        );

        // Auth Service может возвращать либо JSON с токенами, либо HTML redirect
        const contentType = response.headers.get('content-type') || '';
        console.log('[OAuth Callback] Response content-type:', contentType);

        if (contentType.includes('application/json')) {
          // JSON response - попытка получить токены напрямую
          try {
            const data = await response.json();
            console.log('[OAuth Callback] Received JSON response:', data);

            // Проверяем разные форматы ответа
            let accessToken = data.access_token;
            if (!accessToken && data.data) {
              accessToken = data.data.access_token;
            }
            if (!accessToken && data.success && data.data) {
              accessToken = data.data.access_token;
            }

            if (accessToken) {
              console.log(
                '[OAuth Callback] Got access token from JSON response, saving to TokenManager'
              );
              tokenManager.setAccessToken(accessToken);

              // Также сохраняем refresh токен если есть
              const refreshToken =
                data.refresh_token || (data.data && data.data.refresh_token);
              if (refreshToken) {
                tokenManager.setRefreshToken(refreshToken);
              }

              setStatus('Login successful! Redirecting...');
              setTimeout(() => {
                router.replace('/');
              }, 1000);
              return;
            } else {
              console.warn(
                '[OAuth Callback] JSON response without access token:',
                data
              );
              // Fallback to refresh flow
            }
          } catch (jsonError) {
            console.error(
              '[OAuth Callback] Failed to parse JSON response:',
              jsonError
            );
            // Fallback to refresh flow
          }
        }

        // Fallback: попробовать refresh flow
        // Auth Service должен был установить httpOnly refresh_token cookie
        console.log('[OAuth Callback] Trying refresh flow for access token...');
        setStatus('Completing authentication...');

        // Небольшая задержка для обеспечения что cookie установлены
        await new Promise((resolve) => setTimeout(resolve, 500));

        try {
          const refreshResponse = await fetch('/api/v1/auth/refresh', {
            method: 'POST',
            credentials: 'include', // Критично для httpOnly cookies
            headers: {
              'Content-Type': 'application/json',
            },
          });

          console.log(
            '[OAuth Callback] Refresh response status:',
            refreshResponse.status
          );

          if (refreshResponse.ok) {
            const refreshData = await refreshResponse.json();
            console.log('[OAuth Callback] Refresh response:', refreshData);

            // Обрабатываем разные форматы ответа refresh
            let accessToken = refreshData.access_token;
            if (!accessToken && refreshData.data) {
              accessToken = refreshData.data.access_token;
            }

            if (accessToken) {
              console.log(
                '[OAuth Callback] Got access token from refresh, saving to TokenManager'
              );
              tokenManager.setAccessToken(accessToken);

              // Также сохраняем refresh токен если есть
              const refreshToken =
                refreshData.refresh_token ||
                (refreshData.data && refreshData.data.refresh_token);
              if (refreshToken) {
                tokenManager.setRefreshToken(refreshToken);
              }

              setStatus('Login successful! Redirecting...');
              setTimeout(() => {
                router.replace('/');
              }, 1000);
              return;
            } else {
              console.error(
                '[OAuth Callback] No access token in refresh response:',
                refreshData
              );
            }
          } else {
            const errorText = await refreshResponse.text();
            console.error(
              '[OAuth Callback] Refresh failed:',
              refreshResponse.status,
              errorText
            );

            // Специальная обработка 401 - возможно cookie не установлены
            if (refreshResponse.status === 401) {
              console.error(
                '[OAuth Callback] 401 on refresh - likely no refresh token cookie set by Auth Service'
              );
              setError(
                'Authentication setup incomplete. Please contact support.'
              );
              return;
            }
          }
        } catch (refreshError) {
          console.error('[OAuth Callback] Error during refresh:', refreshError);
        }

        // Если дошли до сюда - что-то пошло не так
        console.error('[OAuth Callback] All token acquisition methods failed');
        setError('Authentication failed. Please try again.');
      } catch (error) {
        console.error('[OAuth Callback] Error during OAuth callback:', error);
        setError('An error occurred during login. Please try again.');
      }
    };

    // Small delay to ensure everything is loaded
    const timer = setTimeout(handleOAuthCallback, 100);
    return () => clearTimeout(timer);
  }, [searchParams, router]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-base-200">
      <div className="card w-96 bg-base-100 shadow-xl">
        <div className="card-body items-center text-center">
          {error ? (
            <>
              <div className="text-error text-4xl mb-4">⚠️</div>
              <h2 className="card-title mt-4 text-error">Login Failed</h2>
              <p className="text-base-content/70 mb-4">{error}</p>
              <button
                className="btn btn-primary"
                onClick={() => (window.location.href = '/')}
              >
                Return to Home
              </button>
            </>
          ) : (
            <>
              <span className="loading loading-spinner loading-lg text-primary"></span>
              <h2 className="card-title mt-4">{status}</h2>
              <p className="text-base-content/70">
                Please wait while we complete your Google login
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}
