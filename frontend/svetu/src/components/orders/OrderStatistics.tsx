'use client';

import React, { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
// import { apiClient } from '@/services/api-client';
import { balanceService } from '@/services/balance';
import { useAuthContext } from '@/contexts/AuthContext';

interface OrderStats {
  totalOrders: number;
  activeOrders: number;
  completedOrders: number;
  totalSpent: number;
  averageOrderValue: number;
  thisMonthOrders: number;
  thisMonthSpent: number;
  lastMonthOrders: number;
  lastMonthSpent: number;
}

export default function OrderStatistics() {
  const t = useTranslations('orders.statistics');
  const { isAuthenticated } = useAuthContext();
  const [stats, setStats] = useState<OrderStats | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!isAuthenticated) {
      setLoading(false);
      return;
    }

    const fetchStatistics = async () => {
      try {
        // Здесь будет реальный запрос к API
        // const response = await apiClient.get('/api/v1/orders/statistics');
        // setStats(response.data.data);

        // Временные данные для демонстрации
        setStats({
          totalOrders: 24,
          activeOrders: 3,
          completedOrders: 18,
          totalSpent: 125000,
          averageOrderValue: 5208,
          thisMonthOrders: 5,
          thisMonthSpent: 32500,
          lastMonthOrders: 8,
          lastMonthSpent: 45000,
        });
      } catch (error) {
        console.error('Failed to fetch order statistics:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchStatistics();
  }, [isAuthenticated]);

  if (!isAuthenticated || loading) {
    return (
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-8">
        {[...Array(4)].map((_, i) => (
          <div key={i} className="card bg-base-100 shadow">
            <div className="card-body">
              <div className="skeleton h-4 w-24 mb-2"></div>
              <div className="skeleton h-8 w-32"></div>
              <div className="skeleton h-3 w-20 mt-2"></div>
            </div>
          </div>
        ))}
      </div>
    );
  }

  if (!stats) return null;

  const statCards = [
    {
      title: t('totalSpent'),
      value: balanceService.formatAmount(stats.totalSpent, 'RSD'),
      subtitle: `${stats.totalOrders} ${t('orders')}`,
      icon: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={1.5}
          stroke="currentColor"
          className="w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M2.25 8.25h19.5M2.25 9h19.5m-16.5 5.25h6m-6 2.25h3m-3.75 3h15a2.25 2.25 0 002.25-2.25V6.75A2.25 2.25 0 0019.5 4.5h-15a2.25 2.25 0 00-2.25 2.25v10.5A2.25 2.25 0 004.5 19.5z"
          />
        </svg>
      ),
      color: 'text-primary',
    },
    {
      title: t('activeOrders'),
      value: stats.activeOrders.toString(),
      subtitle: t('inProgress'),
      icon: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={1.5}
          stroke="currentColor"
          className="w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M16.023 9.348h4.992v-.001M2.985 19.644v-4.992m0 0h4.992m-4.993 0l3.181 3.183a8.25 8.25 0 0013.803-3.7M4.031 9.865a8.25 8.25 0 0113.803-3.7l3.181 3.182m0-4.991v4.99"
          />
        </svg>
      ),
      color: 'text-warning',
    },
    {
      title: t('completedOrders'),
      value: stats.completedOrders.toString(),
      subtitle: t('allTime'),
      icon: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={1.5}
          stroke="currentColor"
          className="w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M9 12.75L11.25 15 15 9.75M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
      ),
      color: 'text-success',
    },
    {
      title: t('averageOrderValue'),
      value: balanceService.formatAmount(stats.averageOrderValue, 'RSD'),
      subtitle: t('perOrder'),
      icon: (
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          strokeWidth={1.5}
          stroke="currentColor"
          className="w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            d="M3 13.125C3 12.504 3.504 12 4.125 12h2.25c.621 0 1.125.504 1.125 1.125v6.75C7.5 20.496 6.996 21 6.375 21h-2.25A1.125 1.125 0 013 19.875v-6.75zM9.75 8.625c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125v11.25c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V8.625zM16.5 4.125c0-.621.504-1.125 1.125-1.125h2.25C20.496 3 21 3.504 21 4.125v15.75c0 .621-.504 1.125-1.125 1.125h-2.25a1.125 1.125 0 01-1.125-1.125V4.125z"
          />
        </svg>
      ),
      color: 'text-info',
    },
  ];

  return (
    <div className="space-y-6">
      {/* Main Stats */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        {statCards.map((stat, index) => (
          <div
            key={index}
            className="card bg-base-100 shadow hover:shadow-lg transition-shadow"
          >
            <div className="card-body">
              <div className="flex items-center justify-between">
                <div className="flex-1">
                  <p className="text-sm text-base-content/70">{stat.title}</p>
                  <p className="text-2xl font-bold">{stat.value}</p>
                  <p className="text-xs text-base-content/60 mt-1">
                    {stat.subtitle}
                  </p>
                </div>
                <div className={`${stat.color} opacity-20`}>{stat.icon}</div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Monthly Comparison */}
      <div className="card bg-base-100 shadow">
        <div className="card-body">
          <h3 className="card-title text-lg">{t('monthlyComparison')}</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* This Month */}
            <div className="space-y-2">
              <div className="flex justify-between items-center">
                <span className="text-sm text-base-content/70">
                  {t('thisMonth')}
                </span>
                <span className="badge badge-primary badge-sm">
                  {t('current')}
                </span>
              </div>
              <div className="bg-base-200 rounded-lg p-4">
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm">{t('orders')}</span>
                  <span className="font-semibold">{stats.thisMonthOrders}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">{t('spent')}</span>
                  <span className="font-semibold">
                    {balanceService.formatAmount(stats.thisMonthSpent, 'RSD')}
                  </span>
                </div>
              </div>
            </div>

            {/* Last Month */}
            <div className="space-y-2">
              <div className="flex justify-between items-center">
                <span className="text-sm text-base-content/70">
                  {t('lastMonth')}
                </span>
              </div>
              <div className="bg-base-200 rounded-lg p-4">
                <div className="flex justify-between items-center mb-2">
                  <span className="text-sm">{t('orders')}</span>
                  <span className="font-semibold">{stats.lastMonthOrders}</span>
                </div>
                <div className="flex justify-between items-center">
                  <span className="text-sm">{t('spent')}</span>
                  <span className="font-semibold">
                    {balanceService.formatAmount(stats.lastMonthSpent, 'RSD')}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Trend Indicator */}
          {stats.thisMonthSpent > stats.lastMonthSpent ? (
            <div className="alert alert-success mt-4">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-6 h-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M2.25 18L9 11.25l4.306 4.307a11.95 11.95 0 015.814-5.519l2.74-1.22m0 0l-5.94-2.28m5.94 2.28l-2.28 5.941"
                />
              </svg>
              <span>
                {t('spendingIncrease', {
                  percent: Math.round(
                    ((stats.thisMonthSpent - stats.lastMonthSpent) /
                      stats.lastMonthSpent) *
                      100
                  ),
                })}
              </span>
            </div>
          ) : stats.thisMonthSpent < stats.lastMonthSpent ? (
            <div className="alert alert-warning mt-4">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-6 h-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M2.25 6L9 12.75l4.286-4.286a11.948 11.948 0 014.306 6.43l.776 2.898m0 0l3.182-5.511m-3.182 5.51l-5.511-3.181"
                />
              </svg>
              <span>
                {t('spendingDecrease', {
                  percent: Math.round(
                    ((stats.lastMonthSpent - stats.thisMonthSpent) /
                      stats.lastMonthSpent) *
                      100
                  ),
                })}
              </span>
            </div>
          ) : null}
        </div>
      </div>
    </div>
  );
}
