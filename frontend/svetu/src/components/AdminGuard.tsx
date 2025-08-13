'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useState, useEffect } from 'react';
import AuthModal from './AuthModal';
import { useTranslations } from 'next-intl';

interface AdminGuardProps {
  children: React.ReactNode;
  loading?: React.ReactNode;
}

export default function AdminGuard({ children, loading }: AdminGuardProps) {
  const { user, isLoading } = useAuth();
  const t = useTranslations('admin');
  const [showAuthModal, setShowAuthModal] = useState(false);
  const [isChecking, setIsChecking] = useState(true);

  useEffect(() => {
    if (!isLoading) {
      if (!user) {
        console.log('[AdminGuard] No user found, showing auth modal');
        setShowAuthModal(true);
        setIsChecking(false);
      } else if (!user.is_admin) {
        console.log('[AdminGuard] User is not admin', {
          user_id: user.id,
          is_admin: user.is_admin,
        });
        setIsChecking(false);
      } else {
        console.log('[AdminGuard] Admin access granted', {
          user_id: user.id,
          is_admin: user.is_admin,
        });
        setIsChecking(false);
      }
    }
  }, [user, isLoading]);

  // Show loading state while checking auth
  if (isLoading || isChecking) {
    return (
      loading || (
        <div className="flex items-center justify-center min-h-screen">
          <div className="loading loading-spinner loading-lg"></div>
        </div>
      )
    );
  }

  // Show auth modal if no user
  if (!user) {
    return (
      <>
        <AuthModal
          isOpen={showAuthModal}
          onClose={() => setShowAuthModal(false)}
          onSuccess={() => {
            setShowAuthModal(false);
            // Re-check after successful login
            setIsChecking(true);
          }}
          requireAdmin={true}
        />
        <div className="flex items-center justify-center min-h-screen">
          <div className="text-center">
            <h2 className="text-2xl font-bold mb-4">{t('authRequired')}</h2>
            <p className="mb-4">{t('authRequiredDescription')}</p>
            <button
              onClick={() => setShowAuthModal(true)}
              className="btn btn-primary"
            >
              {t('login')}
            </button>
          </div>
        </div>
      </>
    );
  }

  // Show unauthorized message if user is not admin
  if (!user.is_admin) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <div className="alert alert-error max-w-md">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <div>
            <h3 className="font-bold">{t('accessDenied')}</h3>
            <div className="text-xs">{t('adminOnly')}</div>
          </div>
        </div>
      </div>
    );
  }

  // User is authenticated and is admin
  return <>{children}</>;
}
