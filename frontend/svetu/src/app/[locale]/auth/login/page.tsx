'use client';

import { useEffect } from 'react';
import { useSearchParams } from 'next/navigation';
import { useAuth } from '@/contexts/AuthContext';

export default function LoginPage() {
  const searchParams = useSearchParams();
  const { login } = useAuth();

  useEffect(() => {
    // Получаем URL для редиректа после логина
    const redirect = searchParams.get('redirect') || '/';

    // Инициируем логин через Google OAuth
    login(redirect);
  }, [searchParams, login]);

  return (
    <div className="min-h-screen flex items-center justify-center bg-base-100">
      <div className="text-center">
        <div className="loading loading-spinner loading-lg text-primary"></div>
        <p className="mt-4 text-base-content/70">Redirecting to login...</p>
      </div>
    </div>
  );
}
