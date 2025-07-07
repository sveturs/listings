'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from '@/i18n/routing';
import { useEffect, useState } from 'react';

interface AdminGuardProps {
  children: React.ReactNode;
  loading?: React.ReactNode;
}

export default function AdminGuard({ children, loading }: AdminGuardProps) {
  const { user, isLoading } = useAuth();
  const router = useRouter();
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    if (!isLoading) {
      if (!user) {
        router.push('/auth/login');
        return;
      }

      // Проверяем является ли пользователь администратором
      if (!user.is_admin) {
        router.push('/');
        return;
      }

      setIsChecking(false);
    }
  }, [user, isLoading, router]);

  if (isLoading || isChecking) {
    return (
      loading || <div className="loading loading-spinner loading-lg"></div>
    );
  }

  if (!user || !user.is_admin) {
    return null;
  }

  return <>{children}</>;
}
