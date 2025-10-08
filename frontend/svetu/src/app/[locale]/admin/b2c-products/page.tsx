import { Suspense } from 'react';
import { setRequestLocale } from 'next-intl/server';
import { getTranslations } from 'next-intl/server';
import AdminB2CProductsTable from './AdminB2CProductsTable';

interface PageProps {
  params: Promise<{ locale: string }>;
}

async function getProductsStats() {
  try {
    // Получаем список всех витрин для админ панели
    const storefrontsResponse = await fetch(
      'http://localhost:3000/api/v1/b2c',
      {
        cache: 'no-store',
      }
    );
    const storefrontsData = await storefrontsResponse.json();

    // API возвращает витрины в поле storefronts, а не data
    const storefronts =
      storefrontsData.storefronts || storefrontsData.data || [];

    if (!storefronts || storefronts.length === 0) {
      return { total: 0, active: 0, outOfStock: 0, sales: 0 };
    }

    let total = 0;
    let active = 0;
    let outOfStock = 0;
    let sales = 0;

    // Получаем статистику для каждой витрины
    for (const storefront of storefronts) {
      if (!storefront.slug) continue;

      try {
        const response = await fetch(
          `http://localhost:3000/api/v1/b2c/slug/${storefront.slug}/products?limit=100`,
          {
            cache: 'no-store',
          }
        );

        if (response.ok) {
          const data = await response.json();
          // API возвращает массив товаров напрямую, а не в поле data
          // Обрабатываем случай когда API возвращает null для пустых результатов
          let products = [];

          if (data === null || data === undefined) {
            products = [];
          } else if (Array.isArray(data)) {
            products = data;
          } else if (data && data.data) {
            products = Array.isArray(data.data) ? data.data : [];
          }

          total += products.length;
          active += products.filter((p: any) => p.is_active !== false).length;
          outOfStock += products.filter((p: any) => {
            const stock = p.stock || p.stock_quantity || 0;
            return !p.unlimited_stock && stock === 0;
          }).length;
          sales += products.reduce(
            (sum: number, p: any) => sum + (p.sales_count || p.sold_count || 0),
            0
          );
        }
      } catch (err) {
        console.error(
          `Error fetching products for storefront ${storefront.slug}:`,
          err
        );
      }
    }

    return { total, active, outOfStock, sales };
  } catch (error) {
    console.error('Error fetching stats:', error);
  }
  return { total: 0, active: 0, outOfStock: 0, sales: 0 };
}

// Делаем страницу динамической из-за API вызовов
export const dynamic = 'force-dynamic';

export default async function B2CProductsPage({ params }: PageProps) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations('admin');
  const stats = await getProductsStats();

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <div>
          <h1 className="text-3xl font-bold">
            {t('storefrontProducts.title')}
          </h1>
          <p className="text-base-content/60 mt-1">
            {t('storefrontProducts.description')}
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
            {t('storefrontProducts.createNew')}
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
            {t('storefrontProducts.export')}
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
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              />
            </svg>
            {t('storefrontProducts.import')}
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
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
              />
            </svg>
          </div>
          <div className="stat-title">
            {t('storefrontProducts.stats.total')}
          </div>
          <div className="stat-value">{stats.total}</div>
          <div className="stat-desc">
            {t('storefrontProducts.stats.totalDesc')}
          </div>
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
          <div className="stat-title">
            {t('storefrontProducts.stats.active')}
          </div>
          <div className="stat-value">{stats.active}</div>
          <div className="stat-desc">
            {t('storefrontProducts.stats.activeDesc')}
          </div>
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
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
          </div>
          <div className="stat-title">
            {t('storefrontProducts.stats.outOfStock')}
          </div>
          <div className="stat-value">{stats.outOfStock}</div>
          <div className="stat-desc">
            {t('storefrontProducts.stats.outOfStockDesc')}
          </div>
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
                d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
              />
            </svg>
          </div>
          <div className="stat-title">
            {t('storefrontProducts.stats.sales')}
          </div>
          <div className="stat-value">{stats.sales.toLocaleString()}</div>
          <div className="stat-desc">
            {t('storefrontProducts.stats.salesDesc')}
          </div>
        </div>
      </div>

      <Suspense
        fallback={
          <div className="flex justify-center items-center h-64">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        }
      >
        <AdminB2CProductsTable />
      </Suspense>
    </div>
  );
}
