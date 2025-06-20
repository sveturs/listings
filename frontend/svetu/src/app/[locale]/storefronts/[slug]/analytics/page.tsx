'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTranslations, useLocale } from 'next-intl';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchStorefrontBySlug,
  fetchStorefrontAnalytics,
} from '@/store/slices/storefrontSlice';
import { useAuth } from '@/contexts/AuthContext';
import Link from 'next/link';
import {
  ArrowLeftIcon,
  ChartBarIcon,
  ArrowTrendingUpIcon,
  ArrowTrendingDownIcon,
  CurrencyDollarIcon,
  ShoppingBagIcon,
  EyeIcon,
  DocumentChartBarIcon,
  GlobeAltIcon,
} from '@heroicons/react/24/outline';
import { Line, Bar, Doughnut, Pie } from 'react-chartjs-2';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler,
} from 'chart.js';

// Register ChartJS components
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend,
  Filler
);

export default function StorefrontAnalyticsPage() {
  const t = useTranslations();
  const locale = useLocale();
  const router = useRouter();
  const params = useParams();
  const dispatch = useAppDispatch();
  const slug = params.slug as string;
  const { user } = useAuth();

  const { currentStorefront, isLoading, analytics, isLoadingAnalytics } =
    useAppSelector((state) => state.storefronts);

  const [accessDenied, setAccessDenied] = useState(false);
  const [timeRange, setTimeRange] = useState<'week' | 'month' | 'year'>(
    'month'
  );

  useEffect(() => {
    if (slug) {
      dispatch(fetchStorefrontBySlug(slug));
    }
  }, [dispatch, slug]);

  useEffect(() => {
    if (currentStorefront?.id) {
      const endDate = new Date();
      const startDate = new Date();

      switch (timeRange) {
        case 'week':
          startDate.setDate(startDate.getDate() - 7);
          break;
        case 'month':
          startDate.setMonth(startDate.getMonth() - 1);
          break;
        case 'year':
          startDate.setFullYear(startDate.getFullYear() - 1);
          break;
      }

      dispatch(
        fetchStorefrontAnalytics({
          id: currentStorefront.id,
          from: startDate.toISOString(),
          to: endDate.toISOString(),
        })
      );
    }
  }, [dispatch, currentStorefront?.id, timeRange]);

  // Проверка доступа
  useEffect(() => {
    if (!isLoading && currentStorefront) {
      if (!user) {
        setAccessDenied(true);
        router.push(`/${locale}/storefronts/${slug}`);
        return;
      }

      if (currentStorefront.user_id !== user.id) {
        setAccessDenied(true);
        router.push(`/${locale}/storefronts/${slug}`);
      }
    }
  }, [currentStorefront, user, isLoading, router, slug, locale]);

  if (accessDenied) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">🔒</div>
          <h1 className="text-2xl font-bold mb-2">
            {t('common.accessDenied')}
          </h1>
          <p className="text-base-content/60">
            {t('storefronts.analyticsAccessDenied')}
          </p>
        </div>
      </div>
    );
  }

  if (isLoading || !currentStorefront || isLoadingAnalytics) {
    return (
      <div className="min-h-screen bg-base-200">
        <div className="container mx-auto px-4 py-8">
          <div className="flex justify-center items-center min-h-[600px]">
            <div className="text-center">
              <span className="loading loading-spinner loading-lg text-primary"></span>
              <p className="mt-4 text-base-content/60">
                {t('admin.common.loading')}
              </p>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Prepare analytics data
  const prepareChartData = () => {
    if (!analytics || analytics.length === 0) {
      return {
        viewsData: { labels: [], datasets: [] },
        salesData: { labels: [], datasets: [] },
        conversionData: { labels: [], datasets: [] },
        revenueData: { labels: [], datasets: [] },
        stats: {
          totalRevenue: 0,
          revenueChange: 0,
          totalOrders: 0,
          ordersChange: 0,
          totalCustomers: 0,
          customersChange: 0,
          avgOrderValue: 0,
          avgOrderChange: 0,
          conversionRate: 0,
          conversionChange: 0,
          totalViews: 0,
          viewsChange: 0,
        },
      };
    }

    const sortedAnalytics = [...analytics].sort(
      (a, b) =>
        new Date(a.date || 0).getTime() - new Date(b.date || 0).getTime()
    );

    const labels = sortedAnalytics.map((a) =>
      new Date(a.date || 0).toLocaleDateString('en-US', {
        month: 'short',
        day: 'numeric',
      })
    );

    const pageViews = sortedAnalytics.map((a) => a.page_views || 0);
    const uniqueVisitors = sortedAnalytics.map((a) => a.unique_visitors || 0);
    const orders = sortedAnalytics.map((a) => a.orders_count || 0);
    const revenue = sortedAnalytics.map((a) => a.revenue || 0);
    const conversionRates = sortedAnalytics.map((a) => {
      const views = a.page_views || 1;
      const orders = a.orders_count || 0;
      return ((orders / views) * 100).toFixed(2);
    });

    // Calculate totals
    const totalRevenue = revenue.reduce((sum, r) => sum + r, 0);
    const totalOrders = orders.reduce((sum, o) => sum + o, 0);
    const totalCustomers = sortedAnalytics.reduce(
      (sum, a) => sum + (a.unique_visitors || 0),
      0
    );
    const totalViews = pageViews.reduce((sum, v) => sum + v, 0);
    const avgOrderValue = totalOrders > 0 ? totalRevenue / totalOrders : 0;
    const conversionRate =
      totalViews > 0 ? (totalOrders / totalViews) * 100 : 0;

    // Calculate changes
    const midPoint = Math.floor(sortedAnalytics.length / 2);
    const firstHalf = sortedAnalytics.slice(0, midPoint);
    const secondHalf = sortedAnalytics.slice(midPoint);

    const firstHalfRevenue = firstHalf.reduce(
      (sum, a) => sum + (a.revenue || 0),
      0
    );
    const secondHalfRevenue = secondHalf.reduce(
      (sum, a) => sum + (a.revenue || 0),
      0
    );
    const revenueChange =
      firstHalfRevenue > 0
        ? ((secondHalfRevenue - firstHalfRevenue) / firstHalfRevenue) * 100
        : 0;

    const firstHalfViews = firstHalf.reduce(
      (sum, a) => sum + (a.page_views || 0),
      0
    );
    const secondHalfViews = secondHalf.reduce(
      (sum, a) => sum + (a.page_views || 0),
      0
    );
    const viewsChange =
      firstHalfViews > 0
        ? ((secondHalfViews - firstHalfViews) / firstHalfViews) * 100
        : 0;

    return {
      viewsData: {
        labels,
        datasets: [
          {
            label: t('storefronts.pageViews'),
            data: pageViews,
            fill: true,
            borderColor: 'rgb(75, 192, 192)',
            backgroundColor: 'rgba(75, 192, 192, 0.1)',
            tension: 0.4,
          },
          {
            label: t('storefronts.uniqueVisitors'),
            data: uniqueVisitors,
            fill: true,
            borderColor: 'rgb(153, 102, 255)',
            backgroundColor: 'rgba(153, 102, 255, 0.1)',
            tension: 0.4,
          },
        ],
      },
      salesData: {
        labels,
        datasets: [
          {
            label: t('storefronts.orders'),
            data: orders,
            backgroundColor: 'rgba(99, 102, 241, 0.8)',
            borderColor: 'rgb(99, 102, 241)',
            borderWidth: 1,
          },
        ],
      },
      conversionData: {
        labels,
        datasets: [
          {
            label: t('storefronts.conversionRate') + ' (%)',
            data: conversionRates,
            fill: false,
            borderColor: 'rgb(255, 99, 132)',
            backgroundColor: 'rgba(255, 99, 132, 0.1)',
            tension: 0.4,
          },
        ],
      },
      revenueData: {
        labels,
        datasets: [
          {
            label: t('storefronts.revenue'),
            data: revenue,
            fill: true,
            borderColor: 'rgb(34, 197, 94)',
            backgroundColor: 'rgba(34, 197, 94, 0.1)',
            tension: 0.4,
          },
        ],
      },
      stats: {
        totalRevenue,
        revenueChange,
        totalOrders,
        ordersChange: 0,
        totalCustomers,
        customersChange: 0,
        avgOrderValue,
        avgOrderChange: 0,
        conversionRate,
        conversionChange: 0,
        totalViews,
        viewsChange,
      },
    };
  };

  const { viewsData, salesData, conversionData, revenueData, stats } =
    prepareChartData();

  // Category data
  const categoryData = (() => {
    if (!analytics || analytics.length === 0) {
      return { labels: [], datasets: [] };
    }

    const categoryMap = new Map();
    analytics.forEach((a) => {
      if (a.top_categories && typeof a.top_categories === 'object') {
        const categories = Array.isArray(a.top_categories)
          ? a.top_categories
          : [];
        categories.forEach((cat: any) => {
          if (cat.name && cat.count) {
            const current = categoryMap.get(cat.name) || 0;
            categoryMap.set(cat.name, current + cat.count);
          }
        });
      }
    });

    const sortedCategories = Array.from(categoryMap.entries())
      .sort((a, b) => b[1] - a[1])
      .slice(0, 5);

    return {
      labels: sortedCategories.map(([name]) => name),
      datasets: [
        {
          data: sortedCategories.map(([, count]) => count),
          backgroundColor: [
            'rgba(255, 99, 132, 0.8)',
            'rgba(54, 162, 235, 0.8)',
            'rgba(255, 206, 86, 0.8)',
            'rgba(75, 192, 192, 0.8)',
            'rgba(153, 102, 255, 0.8)',
          ],
          borderWidth: 0,
        },
      ],
    };
  })();

  // Device data (mock for now)
  const deviceData = {
    labels: [
      t('storefronts.mobile'),
      t('storefronts.desktop'),
      t('storefronts.tablet'),
    ],
    datasets: [
      {
        data: [45, 40, 15],
        backgroundColor: [
          'rgba(99, 102, 241, 0.8)',
          'rgba(34, 197, 94, 0.8)',
          'rgba(251, 146, 60, 0.8)',
        ],
        borderWidth: 0,
      },
    ],
  };

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: true,
        position: 'top' as const,
      },
    },
    scales: {
      y: {
        beginAtZero: true,
      },
    },
  };

  const doughnutOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        position: 'right' as const,
      },
    },
  };

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-base-100 shadow-sm border-b border-base-300">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Link
                href={`/${locale}/storefronts/${slug}/dashboard`}
                className="btn btn-ghost btn-sm btn-square"
              >
                <ArrowLeftIcon className="w-5 h-5" />
              </Link>
              <div>
                <h1 className="text-2xl font-bold">{currentStorefront.name}</h1>
                <p className="text-sm text-base-content/60">
                  {t('storefronts.detailedAnalytics')}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-4">
              {/* Time Range Selector */}
              <div className="btn-group">
                <button
                  className={`btn btn-sm ${timeRange === 'week' ? 'btn-primary' : 'btn-ghost'}`}
                  onClick={() => setTimeRange('week')}
                >
                  {t('common.week')}
                </button>
                <button
                  className={`btn btn-sm ${timeRange === 'month' ? 'btn-primary' : 'btn-ghost'}`}
                  onClick={() => setTimeRange('month')}
                >
                  {t('common.month')}
                </button>
                <button
                  className={`btn btn-sm ${timeRange === 'year' ? 'btn-primary' : 'btn-ghost'}`}
                  onClick={() => setTimeRange('year')}
                >
                  {t('common.year')}
                </button>
              </div>

              {/* Export Button */}
              <button className="btn btn-outline btn-sm">
                <DocumentChartBarIcon className="w-4 h-4" />
                {t('storefronts.exportReport')}
              </button>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Key Metrics Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
          {/* Total Revenue */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-base-content/60">
                    {t('storefronts.totalRevenue')}
                  </p>
                  <p className="text-2xl font-bold mt-1">
                    ${stats.totalRevenue.toFixed(2)}
                  </p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.revenueChange > 0 ? (
                      <ArrowTrendingUpIcon className="w-4 h-4 text-success" />
                    ) : (
                      <ArrowTrendingDownIcon className="w-4 h-4 text-error" />
                    )}
                    <span
                      className={`text-sm ${stats.revenueChange > 0 ? 'text-success' : 'text-error'}`}
                    >
                      {Math.abs(stats.revenueChange).toFixed(1)}%
                    </span>
                  </div>
                </div>
                <div className="p-3 bg-primary/10 rounded-xl">
                  <CurrencyDollarIcon className="w-8 h-8 text-primary" />
                </div>
              </div>
            </div>
          </div>

          {/* Total Views */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-base-content/60">
                    {t('storefronts.totalViews')}
                  </p>
                  <p className="text-2xl font-bold mt-1">
                    {stats.totalViews.toLocaleString()}
                  </p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.viewsChange > 0 ? (
                      <ArrowTrendingUpIcon className="w-4 h-4 text-success" />
                    ) : (
                      <ArrowTrendingDownIcon className="w-4 h-4 text-error" />
                    )}
                    <span
                      className={`text-sm ${stats.viewsChange > 0 ? 'text-success' : 'text-error'}`}
                    >
                      {Math.abs(stats.viewsChange).toFixed(1)}%
                    </span>
                  </div>
                </div>
                <div className="p-3 bg-secondary/10 rounded-xl">
                  <EyeIcon className="w-8 h-8 text-secondary" />
                </div>
              </div>
            </div>
          </div>

          {/* Conversion Rate */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-base-content/60">
                    {t('storefronts.conversionRate')}
                  </p>
                  <p className="text-2xl font-bold mt-1">
                    {stats.conversionRate.toFixed(2)}%
                  </p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.conversionChange > 0 ? (
                      <ArrowTrendingUpIcon className="w-4 h-4 text-success" />
                    ) : (
                      <ArrowTrendingDownIcon className="w-4 h-4 text-error" />
                    )}
                    <span
                      className={`text-sm ${stats.conversionChange > 0 ? 'text-success' : 'text-error'}`}
                    >
                      {Math.abs(stats.conversionChange).toFixed(1)}%
                    </span>
                  </div>
                </div>
                <div className="p-3 bg-accent/10 rounded-xl">
                  <ChartBarIcon className="w-8 h-8 text-accent" />
                </div>
              </div>
            </div>
          </div>

          {/* Average Order Value */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-base-content/60">
                    {t('storefronts.avgOrderValue')}
                  </p>
                  <p className="text-2xl font-bold mt-1">
                    ${stats.avgOrderValue.toFixed(2)}
                  </p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.avgOrderChange > 0 ? (
                      <ArrowTrendingUpIcon className="w-4 h-4 text-success" />
                    ) : (
                      <ArrowTrendingDownIcon className="w-4 h-4 text-error" />
                    )}
                    <span
                      className={`text-sm ${stats.avgOrderChange > 0 ? 'text-success' : 'text-error'}`}
                    >
                      {Math.abs(stats.avgOrderChange).toFixed(1)}%
                    </span>
                  </div>
                </div>
                <div className="p-3 bg-warning/10 rounded-xl">
                  <ShoppingBagIcon className="w-8 h-8 text-warning" />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Charts Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Traffic Over Time */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.trafficOverTime')}
              </h2>
              <div className="h-64">
                <Line data={viewsData} options={chartOptions} />
              </div>
            </div>
          </div>

          {/* Revenue Over Time */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.revenueOverTime')}
              </h2>
              <div className="h-64">
                <Line data={revenueData} options={chartOptions} />
              </div>
            </div>
          </div>

          {/* Orders Over Time */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.ordersOverTime')}
              </h2>
              <div className="h-64">
                <Bar data={salesData} options={chartOptions} />
              </div>
            </div>
          </div>

          {/* Conversion Rate Trend */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.conversionTrend')}
              </h2>
              <div className="h-64">
                <Line data={conversionData} options={chartOptions} />
              </div>
            </div>
          </div>
        </div>

        {/* Additional Analytics */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Sales by Category */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.salesByCategory')}
              </h2>
              <div className="h-64">
                <Doughnut data={categoryData} options={doughnutOptions} />
              </div>
            </div>
          </div>

          {/* Traffic by Device */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.trafficByDevice')}
              </h2>
              <div className="h-64">
                <Pie data={deviceData} options={doughnutOptions} />
              </div>
            </div>
          </div>

          {/* Top Products */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.topProducts')}
              </h2>
              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <span className="text-sm">iPhone 15 Pro Max</span>
                  <div className="text-right">
                    <p className="text-sm font-semibold">
                      23 {t('storefronts.sold')}
                    </p>
                    <p className="text-xs text-base-content/60">$29,897</p>
                  </div>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">MacBook Pro M3</span>
                  <div className="text-right">
                    <p className="text-sm font-semibold">
                      18 {t('storefronts.sold')}
                    </p>
                    <p className="text-xs text-base-content/60">$43,182</p>
                  </div>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Sony WH-1000XM5</span>
                  <div className="text-right">
                    <p className="text-sm font-semibold">
                      31 {t('storefronts.sold')}
                    </p>
                    <p className="text-xs text-base-content/60">$12,399</p>
                  </div>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">Samsung Galaxy S24</span>
                  <div className="text-right">
                    <p className="text-sm font-semibold">
                      15 {t('storefronts.sold')}
                    </p>
                    <p className="text-xs text-base-content/60">$14,985</p>
                  </div>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">iPad Pro 12.9&quot;</span>
                  <div className="text-right">
                    <p className="text-sm font-semibold">
                      12 {t('storefronts.sold')}
                    </p>
                    <p className="text-xs text-base-content/60">$11,988</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Traffic Sources */}
        <div className="card bg-base-100 shadow-md mt-6">
          <div className="card-body">
            <h2 className="card-title text-lg mb-4">
              {t('storefronts.trafficSources')}
            </h2>
            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>{t('storefronts.source')}</th>
                    <th>{t('storefronts.visitors')}</th>
                    <th>{t('storefronts.pageViews')}</th>
                    <th>{t('storefronts.conversionRate')}</th>
                    <th>{t('storefronts.revenue')}</th>
                  </tr>
                </thead>
                <tbody>
                  <tr>
                    <td>
                      <div className="flex items-center gap-2">
                        <GlobeAltIcon className="w-4 h-4" />
                        {t('storefronts.direct')}
                      </div>
                    </td>
                    <td>1,234</td>
                    <td>3,456</td>
                    <td>2.8%</td>
                    <td>$2,345</td>
                  </tr>
                  <tr>
                    <td>
                      <div className="flex items-center gap-2">
                        <GlobeAltIcon className="w-4 h-4" />
                        {t('storefronts.search')}
                      </div>
                    </td>
                    <td>987</td>
                    <td>2,345</td>
                    <td>3.2%</td>
                    <td>$3,123</td>
                  </tr>
                  <tr>
                    <td>
                      <div className="flex items-center gap-2">
                        <GlobeAltIcon className="w-4 h-4" />
                        {t('storefronts.social')}
                      </div>
                    </td>
                    <td>567</td>
                    <td>1,234</td>
                    <td>1.5%</td>
                    <td>$789</td>
                  </tr>
                  <tr>
                    <td>
                      <div className="flex items-center gap-2">
                        <GlobeAltIcon className="w-4 h-4" />
                        {t('storefronts.referral')}
                      </div>
                    </td>
                    <td>345</td>
                    <td>890</td>
                    <td>4.1%</td>
                    <td>$1,234</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
