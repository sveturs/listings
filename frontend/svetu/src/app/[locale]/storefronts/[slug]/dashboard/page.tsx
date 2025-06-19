'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchStorefrontBySlug,
  fetchStorefrontAnalytics,
} from '@/store/slices/storefrontSlice';
import { usePageViewTracking } from '@/hooks/useAnalytics';
import Link from 'next/link';
import {
  ChartBarIcon,
  ShoppingBagIcon,
  UserGroupIcon,
  StarIcon,
  CurrencyDollarIcon,
  CalendarIcon,
  ArrowTrendingUpIcon,
  ArrowTrendingDownIcon,
  ArrowLeftIcon,
} from '@heroicons/react/24/outline';
import { Line, Bar, Doughnut } from 'react-chartjs-2';
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

export default function StorefrontDashboardPage() {
  const t = useTranslations();
  const params = useParams();
  const dispatch = useAppDispatch();
  const slug = params.slug as string;

  const { currentStorefront, isLoading, analytics, isLoadingAnalytics } =
    useAppSelector((state) => state.storefronts);

  const [timeRange, setTimeRange] = useState<'week' | 'month' | 'year'>(
    'month'
  );

  // Отслеживание просмотра страницы дашборда
  usePageViewTracking(currentStorefront?.id, 'dashboard');

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

  // Prepare real data for charts
  const prepareChartData = () => {
    if (!analytics || analytics.length === 0) {
      // Return empty data structure
      return {
        viewsData: {
          labels: [],
          datasets: [
            {
              label: t('storefronts.views'),
              data: [],
              fill: true,
              borderColor: 'rgb(75, 192, 192)',
              backgroundColor: 'rgba(75, 192, 192, 0.1)',
              tension: 0.4,
            },
          ],
        },
        salesData: {
          labels: [],
          datasets: [
            {
              label: t('storefronts.sales'),
              data: [],
              backgroundColor: 'rgba(99, 102, 241, 0.8)',
              borderColor: 'rgb(99, 102, 241)',
              borderWidth: 1,
            },
          ],
        },
        stats: {
          totalRevenue: 0,
          revenueChange: 0,
          totalOrders: 0,
          ordersChange: 0,
          totalCustomers: 0,
          customersChange: 0,
          avgOrderValue: 0,
          avgOrderChange: 0,
        },
      };
    }

    // Process analytics data
    const sortedAnalytics = [...analytics].sort(
      (a, b) => new Date(a.date || 0).getTime() - new Date(b.date || 0).getTime()
    );

    const labels = sortedAnalytics.map((a) =>
      new Date(a.date || 0).toLocaleDateString('en-US', { weekday: 'short' })
    );

    const pageViews = sortedAnalytics.map((a) => a.page_views || 0);
    const orders = sortedAnalytics.map((a) => a.orders_count || 0);
    const revenue = sortedAnalytics.map((a) => a.revenue || 0);

    // Calculate totals and changes
    const totalRevenue = revenue.reduce((sum, r) => sum + r, 0);
    const totalOrders = orders.reduce((sum, o) => sum + o, 0);
    const totalCustomers = sortedAnalytics.reduce(
      (sum, a) => sum + (a.unique_visitors || 0),
      0
    );
    const avgOrderValue = totalOrders > 0 ? totalRevenue / totalOrders : 0;

    // Calculate changes (comparing last period with previous)
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

    const firstHalfOrders = firstHalf.reduce(
      (sum, a) => sum + (a.orders_count || 0),
      0
    );
    const secondHalfOrders = secondHalf.reduce(
      (sum, a) => sum + (a.orders_count || 0),
      0
    );
    const ordersChange =
      firstHalfOrders > 0
        ? ((secondHalfOrders - firstHalfOrders) / firstHalfOrders) * 100
        : 0;

    return {
      viewsData: {
        labels,
        datasets: [
          {
            label: t('storefronts.views'),
            data: pageViews,
            fill: true,
            borderColor: 'rgb(75, 192, 192)',
            backgroundColor: 'rgba(75, 192, 192, 0.1)',
            tension: 0.4,
          },
        ],
      },
      salesData: {
        labels,
        datasets: [
          {
            label: t('storefronts.sales'),
            data: orders,
            backgroundColor: 'rgba(99, 102, 241, 0.8)',
            borderColor: 'rgb(99, 102, 241)',
            borderWidth: 1,
          },
        ],
      },
      stats: {
        totalRevenue,
        revenueChange,
        totalOrders,
        ordersChange,
        totalCustomers,
        customersChange: 0, // TODO: Calculate based on previous period
        avgOrderValue,
        avgOrderChange: 0, // TODO: Calculate based on previous period
      },
    };
  };

  const { viewsData, salesData, stats } = prepareChartData();

  // Process category data from analytics
  const categoryData = (() => {
    if (!analytics || analytics.length === 0) {
      return {
        labels: [],
        datasets: [{ data: [], backgroundColor: [], borderWidth: 0 }],
      };
    }

    // Aggregate top categories from all analytics
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

  const chartOptions = {
    responsive: true,
    maintainAspectRatio: false,
    plugins: {
      legend: {
        display: false,
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
                href="/profile/storefronts"
                className="btn btn-ghost btn-sm btn-square"
              >
                <ArrowLeftIcon className="w-5 h-5" />
              </Link>
              <div>
                <h1 className="text-2xl font-bold">{currentStorefront.name}</h1>
                <p className="text-sm text-base-content/60">
                  {t('storefronts.dashboardTitle')}
                </p>
              </div>
            </div>

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
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8">
        {/* Stats Grid */}
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
                      {Math.abs(stats.revenueChange)}%
                    </span>
                  </div>
                </div>
                <div className="p-3 bg-primary/10 rounded-xl">
                  <CurrencyDollarIcon className="w-8 h-8 text-primary" />
                </div>
              </div>
            </div>
          </div>

          {/* Total Orders */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-base-content/60">
                    {t('storefronts.totalOrders')}
                  </p>
                  <p className="text-2xl font-bold mt-1">{stats.totalOrders}</p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.ordersChange > 0 ? (
                      <ArrowTrendingUpIcon className="w-4 h-4 text-success" />
                    ) : (
                      <ArrowTrendingDownIcon className="w-4 h-4 text-error" />
                    )}
                    <span
                      className={`text-sm ${stats.ordersChange > 0 ? 'text-success' : 'text-error'}`}
                    >
                      {Math.abs(stats.ordersChange)}%
                    </span>
                  </div>
                </div>
                <div className="p-3 bg-secondary/10 rounded-xl">
                  <ShoppingBagIcon className="w-8 h-8 text-secondary" />
                </div>
              </div>
            </div>
          </div>

          {/* Total Customers */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <div className="flex items-center justify-between">
                <div>
                  <p className="text-sm text-base-content/60">
                    {t('storefronts.totalCustomers')}
                  </p>
                  <p className="text-2xl font-bold mt-1">
                    {stats.totalCustomers}
                  </p>
                  <div className="flex items-center gap-1 mt-2">
                    {stats.customersChange > 0 ? (
                      <ArrowTrendingUpIcon className="w-4 h-4 text-success" />
                    ) : (
                      <ArrowTrendingDownIcon className="w-4 h-4 text-error" />
                    )}
                    <span
                      className={`text-sm ${stats.customersChange > 0 ? 'text-success' : 'text-error'}`}
                    >
                      {Math.abs(stats.customersChange)}%
                    </span>
                  </div>
                </div>
                <div className="p-3 bg-accent/10 rounded-xl">
                  <UserGroupIcon className="w-8 h-8 text-accent" />
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
                      {Math.abs(stats.avgOrderChange)}%
                    </span>
                  </div>
                </div>
                <div className="p-3 bg-warning/10 rounded-xl">
                  <ChartBarIcon className="w-8 h-8 text-warning" />
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Charts */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {/* Views Chart */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.viewsOverTime')}
              </h2>
              <div className="h-64">
                <Line data={viewsData} options={chartOptions} />
              </div>
            </div>
          </div>

          {/* Sales Chart */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.salesOverTime')}
              </h2>
              <div className="h-64">
                <Bar data={salesData} options={chartOptions} />
              </div>
            </div>
          </div>
        </div>

        {/* Additional Charts */}
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

          {/* Recent Activity */}
          <div className="card bg-base-100 shadow-md lg:col-span-2">
            <div className="card-body">
              <h2 className="card-title text-lg mb-4">
                {t('storefronts.recentActivity')}
              </h2>
              <div className="space-y-3">
                {[1, 2, 3, 4, 5].map((i) => (
                  <div
                    key={i}
                    className="flex items-center gap-4 p-3 bg-base-200 rounded-lg"
                  >
                    <div className="p-2 bg-primary/10 rounded-lg">
                      <ShoppingBagIcon className="w-5 h-5 text-primary" />
                    </div>
                    <div className="flex-1">
                      <p className="text-sm font-medium">New order #123{i}</p>
                      <p className="text-xs text-base-content/60">
                        2 hours ago
                      </p>
                    </div>
                    <p className="text-sm font-semibold">$99.99</p>
                  </div>
                ))}
              </div>
              <button className="btn btn-ghost btn-sm w-full mt-4">
                {t('common.viewAll')}
              </button>
            </div>
          </div>
        </div>

        {/* Quick Actions */}
        <div className="mt-8">
          <h2 className="text-xl font-bold mb-4">
            {t('storefronts.quickActions')}
          </h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <Link
              href={`/storefronts/${currentStorefront.slug}/products/new`}
              className="btn btn-outline"
            >
              <ShoppingBagIcon className="w-5 h-5" />
              {t('storefronts.addProduct')}
            </Link>
            <Link
              href={`/storefronts/${currentStorefront.slug}/orders`}
              className="btn btn-outline"
            >
              <CalendarIcon className="w-5 h-5" />
              {t('storefronts.viewOrders')}
            </Link>
            <Link
              href={`/storefronts/${currentStorefront.slug}/customers`}
              className="btn btn-outline"
            >
              <UserGroupIcon className="w-5 h-5" />
              {t('storefronts.viewCustomers')}
            </Link>
            <Link
              href={`/storefronts/${currentStorefront.slug}/settings`}
              className="btn btn-outline"
            >
              <StarIcon className="w-5 h-5" />
              {t('storefronts.settings')}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
