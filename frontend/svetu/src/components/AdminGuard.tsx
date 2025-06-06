'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from '@/i18n/routing';
import { useEffect, useState } from 'react';

interface AdminGuardProps {
  children: React.ReactNode;
  loading?: React.ReactNode;
}

export default function AdminGuard({ children, loading }: AdminGuardProps) {
  const { user, isAuthenticated, isLoading } = useAuth();
  const router = useRouter();
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  useEffect(() => {
    if (isMounted && !isLoading && (!isAuthenticated || !user?.is_admin)) {
      router.push('/');
    }
  }, [isAuthenticated, user?.is_admin, isLoading, router, isMounted]);

  // Показываем loading до тех пор, пока компонент не смонтирован или пока загружается auth
  if (!isMounted || isLoading) {
    return (
      loading || (
        <div className="min-h-screen flex items-center justify-center">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      )
    );
  }

  if (!isAuthenticated || !user?.is_admin) {
    // Показываем loading вместо null для предотвращения hydration mismatch
    return (
      <div className="min-h-screen flex items-center justify-center">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return <>{children}</>;
}
