'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useTranslations } from 'next-intl';
import { useRouter } from '@/i18n/routing';
import { useEffect } from 'react';

export default function AdminPage() {
  const { user, isAuthenticated, isLoading } = useAuth();
  const t = useTranslations('admin');
  const router = useRouter();

  // Handle authentication redirect
  useEffect(() => {
    if (!isLoading && (!isAuthenticated || !user?.is_admin)) {
      router.push('/');
    }
  }, [isAuthenticated, user?.is_admin, isLoading, router]);

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  if (!isAuthenticated || !user?.is_admin) {
    return null; // Redirect handled by useEffect
  }

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-8">{t('title')}</h1>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t('sections.users')}</h2>
            <p>{t('sections.usersDescription')}</p>
            <div className="card-actions justify-end">
              <button className="btn btn-primary" disabled title="Coming soon">
                {t('manage')}
              </button>
            </div>
          </div>
        </div>

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t('sections.listings')}</h2>
            <p>{t('sections.listingsDescription')}</p>
            <div className="card-actions justify-end">
              <button className="btn btn-primary" disabled title="Coming soon">
                {t('manage')}
              </button>
            </div>
          </div>
        </div>

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">{t('sections.categories')}</h2>
            <p>{t('sections.categoriesDescription')}</p>
            <div className="card-actions justify-end">
              <button className="btn btn-primary" disabled title="Coming soon">
                {t('manage')}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
