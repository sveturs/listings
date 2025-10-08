import { Suspense } from 'react';
import { setRequestLocale } from 'next-intl/server';
import { getTranslations } from 'next-intl/server';
import AdminB2CStoresTable from './AdminB2CStoresTable';
import configManager from '@/config';

interface PageProps {
  params: Promise<{ locale: string }>;
}

async function getStorefrontsStats() {
  try {
    const response = await fetch(
      `${configManager.getApiUrl()}/api/v1/b2c`,
      {
        cache: 'no-store',
      }
    );
    const data = await response.json();

    if (data.storefronts) {
      const storefronts = data.storefronts;
      const total = storefronts.length;
      const active = storefronts.filter((s: any) => s.is_active).length;
      const verified = storefronts.filter((s: any) => s.is_verified).length;
      const totalProducts = storefronts.reduce(
        (sum: number, s: any) => sum + (s.products_count || 0),
        0
      );
      const totalSales = storefronts.reduce(
        (sum: number, s: any) => sum + (s.sales_count || 0),
        0
      );

      return { total, active, verified, totalProducts, totalSales };
    }
  } catch (error) {
    console.error('Error fetching storefronts stats:', error);
  }
  return { total: 0, active: 0, verified: 0, totalProducts: 0, totalSales: 0 };
}

// Делаем страницу динамической из-за API вызовов
export const dynamic = 'force-dynamic';

export default async function StorefrontsPage({ params }: PageProps) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations('admin');
  const stats = await getStorefrontsStats();

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold">{t('b2c.title')}</h1>
          <p className="text-base-content/60 mt-1">
            {t('b2c.description')}
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
            {t('b2c.createNew')}
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
            {t('b2c.export')}
          </button>
        </div>
      </div>

      <div className="stats shadow mb-6 w-full">
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
                d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"
              />
            </svg>
          </div>
          <div className="stat-title">{t('b2c.stats.total')}</div>
          <div className="stat-value">{stats.total}</div>
          <div className="stat-desc">{t('b2c.stats.totalDesc')}</div>
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
          <div className="stat-title">{t('b2c.stats.active')}</div>
          <div className="stat-value">{stats.active}</div>
          <div className="stat-desc">{t('b2c.stats.activeDesc')}</div>
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
                d="M9 12l2 2 4-4M7 12l2 2 4-4m0 0a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
          </div>
          <div className="stat-title">{t('b2c.stats.verified')}</div>
          <div className="stat-value">{stats.verified}</div>
          <div className="stat-desc">{t('b2c.stats.verifiedDesc')}</div>
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
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
              />
            </svg>
          </div>
          <div className="stat-title">{t('b2c.stats.products')}</div>
          <div className="stat-value">{stats.totalProducts}</div>
          <div className="stat-desc">{t('b2c.stats.productsDesc')}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-secondary">
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
                d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
          </div>
          <div className="stat-title">{t('b2c.stats.sales')}</div>
          <div className="stat-value">{stats.totalSales.toLocaleString()}</div>
          <div className="stat-desc">{t('b2c.stats.salesDesc')}</div>
        </div>
      </div>

      <Suspense
        fallback={
          <div className="flex justify-center items-center h-64">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        }
      >
        <AdminB2CStoresTable />
      </Suspense>
    </div>
  );
}
