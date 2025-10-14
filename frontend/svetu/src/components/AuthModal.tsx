'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useParams } from 'next/navigation';
import { XMarkIcon } from '@heroicons/react/24/outline';
import { useAuth } from '@/contexts/AuthContext';
import configManager from '@/config';

interface AuthModalProps {
  isOpen: boolean;
  onClose: () => void;
  onSuccess?: () => void;
  requireAdmin?: boolean;
}

export default function AuthModal({
  isOpen,
  onClose,
  onSuccess: _onSuccess,
  requireAdmin = false,
}: AuthModalProps) {
  const t = useTranslations('auth');
  const { login, isLoading } = useAuth();
  const params = useParams();
  const locale = params?.locale || 'en';

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isOpen) {
      // Reset form when modal closes
      setEmail('');
      setPassword('');
      setError(null);
    }
  }, [isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setIsSubmitting(true);

    try {
      // For now, just redirect to login page
      // This component needs to be refactored to work with new auth flow
      await login(email, password);
      onClose();
    } catch {
      setError(t('errors.unexpectedError'));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleGoogleLogin = async () => {
    try {
      setError(null);

      // Construct the OAuth redirect URL with locale
      const currentPath = window.location.pathname;
      const returnUrl = encodeURIComponent(currentPath);

      // Pass locale and returnUrl as query parameters to backend
      const oauthUrl = `${configManager.getApiUrl()}/api/v1/auth/google?locale=${locale}&return_url=${returnUrl}`;

      console.log('[AuthModal] === GOOGLE LOGIN INITIATED ===');
      console.log('[AuthModal] Current locale:', locale);
      console.log('[AuthModal] Return URL:', currentPath);
      console.log('[AuthModal] OAuth URL:', oauthUrl);
      console.log('[AuthModal] === END GOOGLE LOGIN INITIATED ===');

      // Redirect to Google OAuth
      window.location.href = oauthUrl;
    } catch (error) {
      console.error('[AuthModal] Google login error:', error);
      setError(t('errors.googleLoginFailed'));
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
        onClick={onClose}
      />

      {/* Modal */}
      <div className="flex min-h-full items-center justify-center p-4">
        <div className="relative bg-base-100 rounded-lg shadow-xl max-w-md w-full p-6">
          {/* Close button */}
          <button
            onClick={onClose}
            className="absolute top-4 right-4 text-base-content/60 hover:text-base-content"
          >
            <XMarkIcon className="h-6 w-6" />
          </button>

          {/* Title */}
          <h2 className="text-2xl font-bold mb-6">
            {requireAdmin ? t('adminLogin') : t('login')}
          </h2>

          {/* Error message */}
          {error && (
            <div className="alert alert-error mb-4">
              <span>{error}</span>
            </div>
          )}

          {/* Login form */}
          <form onSubmit={handleSubmit} className="space-y-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('email')}</span>
              </label>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="input input-bordered"
                placeholder={t('emailPlaceholder')}
                required
                disabled={isSubmitting}
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('password')}</span>
              </label>
              <input
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="input input-bordered"
                placeholder={t('passwordPlaceholder')}
                required
                disabled={isSubmitting}
              />
            </div>

            <button
              type="submit"
              className="btn btn-primary w-full"
              disabled={isSubmitting || isLoading}
            >
              {isSubmitting ? (
                <span className="loading loading-spinner loading-sm"></span>
              ) : (
                t('loginButton')
              )}
            </button>
          </form>

          {/* Divider */}
          <div className="divider my-6">{t('or')}</div>

          {/* Google login */}
          <button
            onClick={handleGoogleLogin}
            className="btn btn-outline w-full"
            disabled={isSubmitting || isLoading}
          >
            <svg className="w-5 h-5 mr-2" viewBox="0 0 24 24">
              <path
                fill="currentColor"
                d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
              />
              <path
                fill="currentColor"
                d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
              />
              <path
                fill="currentColor"
                d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
              />
              <path
                fill="currentColor"
                d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
              />
            </svg>
            {t('loginWithGoogle')}
          </button>

          {/* Additional info for admin login */}
          {requireAdmin && (
            <div className="alert alert-info mt-4">
              <span className="text-sm">{t('adminLoginInfo')}</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
