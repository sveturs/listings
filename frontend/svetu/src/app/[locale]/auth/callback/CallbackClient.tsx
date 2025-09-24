'use client';

import { useEffect } from 'react';
import { useRouter, useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

export default function CallbackClient() {
  const router = useRouter();
  const searchParams = useSearchParams();
  const { refreshSession } = useAuth();

  useEffect(() => {
    const handleCallback = async () => {
      // Get success status and error from URL
      const success = searchParams?.get('success');
      const error = searchParams?.get('error');
      const returnUrl = searchParams?.get('return_url') || '/';

      if (error) {
        console.error('[AuthCallback] Authentication error:', error);
        router.push('/');
        return;
      }

      if (success === 'true') {
        console.log('[AuthCallback] OAuth successful');
        console.log('[AuthCallback] Return URL:', returnUrl);

        try {
          // Cookies are already set by backend, just refresh session
          await refreshSession();

          // Decode return URL and redirect
          const decodedReturnUrl = decodeURIComponent(returnUrl);
          console.log('[AuthCallback] Redirecting to:', decodedReturnUrl);
          router.push(decodedReturnUrl);
        } catch (err) {
          console.error('[AuthCallback] Error refreshing session:', err);
          router.push('/');
        }
      } else {
        console.error('[AuthCallback] Missing success flag');
        router.push('/');
      }
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