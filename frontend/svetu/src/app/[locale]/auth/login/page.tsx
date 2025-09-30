'use client';

import { useEffect, Suspense } from 'react';
import { useAuth } from '@/contexts/AuthContext';

function LoginPageContent() {
  const { loginWithGoogle } = useAuth();

  useEffect(() => {
    // Инициируем логин через Google OAuth
    loginWithGoogle();
  }, [loginWithGoogle]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-base-100">
      <div className="text-center">
        <div className="loading loading-spinner loading-lg text-primary"></div>
        <p className="mt-4 text-base-content/70">Redirecting to login...</p>
      </div>
    </div>
  );
}

export default function LoginPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center bg-base-100">
          <div className="loading loading-spinner loading-lg text-primary"></div>
        </div>
      }
    >
      <LoginPageContent />
    </Suspense>
  );
}
