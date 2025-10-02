'use client';

import { Suspense, useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import AdminListingsTable from './AdminListingsTable';
import { apiClient } from '@/services/api-client';

interface ListingsStats {
  total: number;
  active: number;
  pending: number;
  views: number;
}

// Делаем страницу динамической из-за API вызовов
export const dynamic = 'force-dynamic';

export default function ListingsPage() {
  const t = useTranslations('admin');
  const [stats, setStats] = useState<ListingsStats>({
    total: 0,
    active: 0,
    pending: 0,
    views: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    async function fetchStats() {
      try {
        console.log('[Admin Listings] Fetching statistics');

        const response = await apiClient.get('/admin/listings/statistics');

        console.log('[Admin Listings] Statistics data:', response.data);

        if (response.data?.success && response.data.data) {
          setStats({
            total: response.data.data.total || 0,
            active: response.data.data.active || 0,
            pending: response.data.data.pending || 0,
            views: response.data.data.views || 0,
          });
        }
      } catch (error) {
        console.error('Error fetching admin statistics:', error);
      } finally {
        setLoading(false);
      }
    }

    fetchStats();
  }, []);

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
          <div className="stat-value">
            {loading ? (
              <span className="loading loading-dots loading-sm"></span>
            ) : (
              stats.total
            )}
          </div>
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
          <div className="stat-value">
            {loading ? (
              <span className="loading loading-dots loading-sm"></span>
            ) : (
              stats.active
            )}
          </div>
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
          <div className="stat-value">
            {loading ? (
              <span className="loading loading-dots loading-sm"></span>
            ) : (
              stats.pending
            )}
          </div>
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
          <div className="stat-value">
            {loading ? (
              <span className="loading loading-dots loading-sm"></span>
            ) : (
              stats.views.toLocaleString()
            )}
          </div>
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
