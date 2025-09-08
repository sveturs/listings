'use client';

import { useEffect, useState, Suspense } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

function OAuthCallbackContent() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { login } = useAuth();
  const [error, setError] = useState<string | null>(null);
  const [isProcessing, setIsProcessing] = useState(true);

  useEffect(() => {
    const handleOAuthCallback = async () => {
      try {
        const code = searchParams.get('code');
        const state = searchParams.get('state');
        const error = searchParams.get('error');
        const errorDescription = searchParams.get('error_description');

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

        console.log('[OAuth Callback] Processing callback with code and state');

        // Exchange code for token with Auth Service
        const response = await fetch(`${process.env.NEXT_PUBLIC_API_URL}/api/v1/auth/oauth/callback`, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({
            code,
            state,
            provider: 'google', // Can be dynamic based on state or URL
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

        // This OAuth callback page is not used in our flow
        // Backend handles OAuth and redirects with token to homepage
        // Commenting out incorrect login call
        /* await login({
          access_token: data.access_token,
          refresh_token: data.refresh_token,
          user: data.user,
        }); */
        
        console.log('[OAuth Callback] This page should not be reached in normal flow');
        console.log('[OAuth Callback] Backend handles OAuth and redirects to homepage with token');

        // Redirect to home or intended destination
        const redirectTo = localStorage.getItem('oauth_redirect_to') || '/';
        localStorage.removeItem('oauth_redirect_to');
        
        console.log('[OAuth Callback] Login successful, redirecting to:', redirectTo);
        router.push(redirectTo);
      } catch (err) {
        console.error('[OAuth Callback] Error processing callback:', err);
        setError(err instanceof Error ? err.message : 'An unexpected error occurred');
        setIsProcessing(false);
      }
    };

    handleOAuthCallback();
  }, [searchParams, login, router]);

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
              onClick={() => router.push('/auth/login')}
              className="btn btn-primary"
            >
              Try Again
            </button>
            <button
              onClick={() => router.push('/')}
              className="btn btn-ghost"
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