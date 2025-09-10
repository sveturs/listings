'use client';

import { Suspense } from 'react';
import { useEffect, useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { tokenManager } from '@/utils/tokenManager';
import { clearLargeHeaders } from '@/utils/clearLargeHeaders';

function GoogleCallbackContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [status, setStatus] = useState('Processing login...');
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const handleOAuthCallback = async () => {
      try {
        // Clear any large headers that might cause issues
        clearLargeHeaders();
        
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
          const errorDescription = urlParams.get('error_description');
          console.error('[OAuth Callback] OAuth error:', {
            error: oauthError,
            description: errorDescription,
          });
          setError(`Authentication failed: ${errorDescription || oauthError}`);
          setStatus('Redirecting to home...');
          setTimeout(() => {
            router.push('/?error=oauth_failed');
          }, 2000);
          return;
        }

        // Check if we have a code
        if (!code) {
          console.error('[OAuth Callback] No code parameter found');
          setError('No authorization code received');
          setStatus('Redirecting to home...');
          setTimeout(() => {
            router.push('/?error=no_code');
          }, 2000);
          return;
        }

        // Exchange code for token via backend
        setStatus('Exchanging authorization code...');
        console.log('[OAuth Callback] Exchanging code for token...');

        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/google/callback?code=${code}&state=${state || ''}`,
          {
            method: 'GET',
            credentials: 'include',
            headers: {
              Accept: 'application/json',
            },
          }
        );

        console.log('[OAuth Callback] Exchange response status:', response.status);

        if (!response.ok) {
          const errorText = await response.text();
          console.error('[OAuth Callback] Exchange failed:', {
            status: response.status,
            statusText: response.statusText,
            body: errorText,
          });
          setError(`Authentication failed: ${response.statusText}`);
          setStatus('Redirecting to login...');
          setTimeout(() => {
            router.push('/?error=exchange_failed');
          }, 2000);
          return;
        }

        const data = await response.json();
        console.log('[OAuth Callback] Exchange successful:', data);

        if (data.data?.access_token) {
          // Store the token
          tokenManager.setAccessToken(data.data.access_token);
          console.log('[OAuth Callback] Token stored successfully');

          setStatus('Login successful! Redirecting...');

          // Redirect to the originally requested page or home
          const redirectTo = state || '/';
          console.log('[OAuth Callback] Redirecting to:', redirectTo);

          // Small delay for user feedback
          setTimeout(() => {
            router.push(redirectTo);
          }, 1000);
        } else {
          console.error('[OAuth Callback] No token in response:', data);
          setError('No authentication token received');
          setStatus('Redirecting to login...');
          setTimeout(() => {
            router.push('/?error=no_token');
          }, 2000);
        }
      } catch (err) {
        console.error('[OAuth Callback] Error during callback handling:', err);
        setError(err instanceof Error ? err.message : 'Authentication failed');
        setStatus('Redirecting to login...');
        setTimeout(() => {
          router.push('/?error=callback_error');
        }, 2000);
      }
    };

    handleOAuthCallback();
  }, [router, searchParams]);

  return (
    <div className="flex items-center justify-center min-h-screen bg-base-200">
      <div className="card bg-base-100 shadow-xl p-8 max-w-md w-full">
        <div className="text-center">
          {!error ? (
            <>
              <div className="loading loading-spinner loading-lg text-primary mb-4"></div>
              <h2 className="text-2xl font-bold mb-2">{status}</h2>
              <p className="text-base-content/70">Please wait...</p>
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
              <p className="text-base-content/70 mb-2">{error}</p>
              <p className="text-sm text-base-content/50">{status}</p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}

export default function GoogleCallbackPage() {
  return (
    <Suspense fallback={
      <div className="flex items-center justify-center min-h-screen bg-base-200">
        <div className="loading loading-spinner loading-lg text-primary"></div>
      </div>
    }>
      <GoogleCallbackContent />
    </Suspense>
  );
}