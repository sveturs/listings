'use client';

import { useEffect } from 'react';
import { useSearchParams } from 'next/navigation';

export default function GoogleCallbackPage() {
  const searchParams = useSearchParams();

  useEffect(() => {
    const handleOAuthCallback = () => {
      // Get parameters from the current URL
      const urlParams = new URLSearchParams(window.location.search);
      const code = urlParams.get('code');
      const state = urlParams.get('state');
      const error = urlParams.get('error');

      console.log('OAuth callback params from URL:', { code, state, error });
      console.log('Full URL:', window.location.href);
      console.log('Search string:', window.location.search);
      
      // If no code and no error, parameters might not be ready yet
      if (!code && !error) {
        console.log('No parameters found in URL');
        // Try searchParams as fallback
        const spCode = searchParams.get('code');
        const spState = searchParams.get('state');
        const spError = searchParams.get('error');
        
        console.log('searchParams values:', { spCode, spState, spError });
        
        if (!spCode && !spError) {
          console.log('No parameters in searchParams either, waiting...');
          return;
        }
        
        // Use searchParams values if available
        if (spCode || spError) {
          const backendUrl = new URL('/api/v1/auth/google/callback', window.location.origin);
          if (spCode) backendUrl.searchParams.set('code', spCode);
          if (spState) backendUrl.searchParams.set('state', spState);
          if (spError) backendUrl.searchParams.set('error', spError);
          
          console.log('Redirecting with searchParams:', backendUrl.toString());
          window.location.href = backendUrl.toString();
          return;
        }
      }

      // Build the backend callback URL
      const backendCallbackUrl = new URL('/api/v1/auth/google/callback', window.location.origin);
      
      // Add all query parameters to the backend URL
      if (code) backendCallbackUrl.searchParams.set('code', code);
      if (state) backendCallbackUrl.searchParams.set('state', state);
      if (error) backendCallbackUrl.searchParams.set('error', error);

      console.log('Redirecting to backend:', backendCallbackUrl.toString());
      
      // Redirect to backend to complete OAuth flow
      window.location.href = backendCallbackUrl.toString();
    };

    // Small delay to ensure everything is loaded
    const timer = setTimeout(handleOAuthCallback, 100);
    return () => clearTimeout(timer);
  }, [searchParams]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-base-200">
      <div className="card w-96 bg-base-100 shadow-xl">
        <div className="card-body items-center text-center">
          <span className="loading loading-spinner loading-lg text-primary"></span>
          <h2 className="card-title mt-4">Processing login...</h2>
          <p className="text-base-content/70">Please wait while we complete your Google login</p>
        </div>
      </div>
    </div>
  );
}