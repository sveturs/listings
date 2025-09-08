'use client';

import { useEffect, useState, Suspense } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { tokenManager } from '@/utils/tokenManager';

function OAuthCallbackContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [error, setError] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(true);
  
  useEffect(() => {
    const handleOAuthCallback = async () => {
      try {
        const code = searchParams.get('code');
        const state = searchParams.get('state');
        const error = searchParams.get('error');
        const errorDescription = searchParams.get('error_description');

        console.log('[OAuth Callback] Processing:', { code: !!code, state: !!state, error });

        // Handle OAuth errors
        if (error) {
          console.error('[OAuth Callback] Error from provider:', error, errorDescription);
          setError(errorDescription || error);
          setIsProcessing(false);
          return;
        }

        // Validate required parameters
        if (!code || !state) {
          console.error('[OAuth Callback] Missing required parameters');
          setError('Invalid OAuth callback parameters');
          setIsProcessing(false);
          return;
        }

        // Exchange code for token with Auth Service through backend proxy
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/oauth/callback`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            code,
            state,
            provider: 'google',
          }),
        });

        if (!response.ok) {
          const errorData = await response.json().catch(() => ({}));
          throw new Error(errorData.message || 'Failed to exchange code for token');
        }

        const data = await response.json();
        
        if (!data.access_token) {
          throw new Error('No access token received from server');
        }

        console.log('[OAuth Callback] Successfully received tokens');

        // Save access token
        tokenManager.setAccessToken(data.access_token);
        // Note: refresh_token is handled internally by tokenManager if needed

        // Redirect to home with token in URL for AuthContext to handle
        console.log('[OAuth Callback] Login successful, redirecting to home');
        router.push(`/ru?auth_token=${encodeURIComponent(data.access_token)}`);
      } catch (err) {
        console.error('[OAuth Callback] Error processing callback:', err);
        setError(err instanceof Error ? err.message : 'An unexpected error occurred');
        setIsProcessing(false);
      }
    };

    handleOAuthCallback();
  }, [searchParams, router]);

  if (isProcessing) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-center">
          <div className="loading loading-spinner loading-lg"></div>
          <h2 className="mt-4 text-xl font-semibold">Processing login...</h2>
          <p className="mt-2 text-gray-600">Please wait while we complete your authentication</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="max-w-md w-full bg-red-50 border border-red-200 rounded-lg p-6">
          <h2 className="text-xl font-semibold text-red-800 mb-2">Authentication Failed</h2>
          <p className="text-red-600 mb-4">{error}</p>
          <div className="flex gap-4">
            <button
              onClick={() => router.push('/ru')}
              className="btn btn-primary"
            >
              Go Home
            </button>
          </div>
        </div>
      </div>
    );
  }

  return null;
}

export default function OAuthCallbackPage() {
  return (
    <Suspense fallback={
      <div className="min-h-screen flex items-center justify-center">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    }>
      <OAuthCallbackContent />
    </Suspense>
  );
}