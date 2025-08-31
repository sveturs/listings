import { Suspense } from 'react';
import { setRequestLocale } from 'next-intl/server';
import { getTranslations } from 'next-intl/server';
import AdminListingsTable from './AdminListingsTable';

interface PageProps {
  params: Promise<{ locale: string }>;
}

async function getListingsStats() {
  try {
    const response = await fetch(
      'http://localhost:3000/api/v1/marketplace/listings?limit=1',
      {
        cache: 'no-store',
      }
    );
    const data = await response.json();
    if (data.success && data.data?.meta) {
      return {
        total: data.data.meta.total || 0,
        // Эти значения можно будет расширить, когда API будет возвращать больше статистики
        active: Math.floor((data.data.meta.total || 0) * 0.7), // Примерное значение
        pending: Math.floor((data.data.meta.total || 0) * 0.1), // Примерное значение
        views: Math.floor((data.data.meta.total || 0) * 150), // Примерное значение
      };
    }
  } catch (error) {
    console.error('Error fetching stats:', error);
  }
  return { total: 0, active: 0, pending: 0, views: 0 };
}

// Делаем страницу динамической из-за API вызовов
export const dynamic = 'force-dynamic';

export default async function ListingsPage({ params }: PageProps) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations('admin');
  const stats = await getListingsStats();

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold">{t('listings.title')}</h1>
          <p className="text-base-content/60 mt-1">
            {t('listings.description')}
          </p>
        </div>

        <div className="flex gap-2">
          <button className="btn btn-primary">
            <svg
              className="w-5 h-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 4v16m8-8H4"
              />
            </svg>
            {t('listings.createNew')}
          </button>
          <button className="btn btn-outline">
            <svg
              className="w-5 h-5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"
              />
            </svg>
            {t('listings.export')}
          </button>
        </div>
      </div>

      <div className="stats shadow mb-6">
        <div className="stat">
          <div className="stat-figure text-primary">
            <svg
              className="w-8 h-8"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
              />
            </svg>
          </div>
          <div className="stat-title">{t('listings.stats.total')}</div>
          <div className="stat-value">{stats.total}</div>
          <div className="stat-desc">{t('listings.stats.totalDesc')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-success">
            <svg
              className="w-8 h-8"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <div className="stat-title">{t('listings.stats.active')}</div>
          <div className="stat-value">{stats.active}</div>
          <div className="stat-desc">{t('listings.stats.activeDesc')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-warning">
            <svg
              className="w-8 h-8"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <div className="stat-title">{t('listings.stats.pending')}</div>
          <div className="stat-value">{stats.pending}</div>
          <div className="stat-desc">{t('listings.stats.pendingDesc')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-info">
            <svg
              className="w-8 h-8"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
              />
            </svg>
          </div>
          <div className="stat-title">{t('listings.stats.views')}</div>
          <div className="stat-value">{stats.views.toLocaleString()}</div>
          <div className="stat-desc">{t('listings.stats.viewsDesc')}</div>
        </div>
      </div>

      <Suspense
        fallback={
          <div className="flex justify-center items-center h-64">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        }
      >
        <AdminListingsTable />
      </Suspense>
    </div>
  );
}
