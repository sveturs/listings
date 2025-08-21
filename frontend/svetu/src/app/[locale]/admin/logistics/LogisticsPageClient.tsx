'use client';

import { useEffect, useState } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useRouter } from 'next/navigation';
import {
  FaBox,
  FaTruck,
  FaCheckCircle,
  FaExclamationTriangle,
  FaClock,
  FaChartLine,
} from 'react-icons/fa';
import { tokenManager } from '@/utils/tokenManager';
import config from '@/config';

interface DashboardStats {
  today_shipments: number;
  today_delivered: number;
  active_shipments: number;
  problem_shipments: number;
  avg_delivery_time: number;
  delivery_success_rate: number;
  weekly_deliveries: Array<{
    date: string;
    shipments: number;
    delivered: number;
    in_transit: number;
    problems: number;
  }>;
  status_distribution: {
    pending: number;
    processing: number;
    in_transit: number;
    delivered: number;
    cancelled: number;
    returned: number;
  };
  courier_performance: Array<{
    name: string;
    shipments: number;
    delivered: number;
    success_rate: number;
    avg_time: number;
  }>;
}

export default function LogisticsPageClient() {
  const t = useTranslations('admin');
  const router = useRouter();
  const locale = useLocale();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchDashboardStats();
  }, []);

  const fetchDashboardStats = async () => {
    try {
      const token = tokenManager.getAccessToken();
      const headers: HeadersInit = {
        'Content-Type': 'application/json',
      };

      if (token) {
        headers['Authorization'] = `Bearer ${token}`;
      }

      const response = await fetch(
        `${config.getApiUrl()}/api/v1/admin/logistics/dashboard`,
        {
          method: 'GET',
          headers,
          credentials: 'include',
        }
      );

      if (!response.ok) {
        throw new Error('Failed to fetch dashboard stats');
      }

      const result = await response.json();
      setStats(result.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex items-center justify-center min-h-screen">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error) {
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
          <span>{error}</span>
        </div>
      </div>
    );
  }

  if (!stats) {
    return null;
  }

  return (
    <div className="container mx-auto">
      <h1 className="text-3xl font-bold mb-8">{t('dashboard')}</h1>

      {/* Основные метрики */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4 mb-8">
        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-primary">
            <FaBox className="text-3xl" />
          </div>
          <div className="stat-title">{t('stats.today_shipments')}</div>
          <div className="stat-value text-primary">{stats.today_shipments}</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-success">
            <FaCheckCircle className="text-3xl" />
          </div>
          <div className="stat-title">{t('stats.today_delivered')}</div>
          <div className="stat-value text-success">{stats.today_delivered}</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-info">
            <FaTruck className="text-3xl" />
          </div>
          <div className="stat-title">{t('stats.active_shipments')}</div>
          <div className="stat-value text-info">{stats.active_shipments}</div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-warning">
            <FaExclamationTriangle className="text-3xl" />
          </div>
          <div className="stat-title">{t('stats.problem_shipments')}</div>
          <div className="stat-value text-warning">
            {stats.problem_shipments}
          </div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-secondary">
            <FaClock className="text-3xl" />
          </div>
          <div className="stat-title">{t('stats.avg_delivery_time')}</div>
          <div className="stat-value text-secondary">
            {stats.avg_delivery_time.toFixed(1)}ч
          </div>
        </div>

        <div className="stat bg-base-100 rounded-lg shadow">
          <div className="stat-figure text-accent">
            <FaChartLine className="text-3xl" />
          </div>
          <div className="stat-title">{t('stats.delivery_success_rate')}</div>
          <div className="stat-value text-accent">
            {stats.delivery_success_rate.toFixed(1)}%
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
        {/* График доставок за неделю */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">{t('stats.weekly_deliveries')}</h2>
            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>{t('table.date')}</th>
                    <th>{t('table.sent')}</th>
                    <th>{t('table.delivered')}</th>
                    <th>{t('table.in_transit')}</th>
                    <th>{t('table.problems')}</th>
                  </tr>
                </thead>
                <tbody>
                  {stats.weekly_deliveries.map((day) => (
                    <tr key={day.date}>
                      <td>{new Date(day.date).toLocaleDateString('ru-RU')}</td>
                      <td>{day.shipments}</td>
                      <td className="text-success">{day.delivered}</td>
                      <td className="text-info">{day.in_transit}</td>
                      <td className={day.problems > 0 ? 'text-warning' : ''}>
                        {day.problems}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>

        {/* Распределение по статусам */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title mb-4">
              {t('stats.status_distribution')}
            </h2>
            <div className="space-y-3">
              <div className="flex justify-between items-center">
                <span className="badge badge-ghost">{t('status.pending')}</span>
                <span className="font-bold">
                  {stats.status_distribution.pending}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="badge badge-info">
                  {t('status.processing')}
                </span>
                <span className="font-bold">
                  {stats.status_distribution.processing}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="badge badge-primary">
                  {t('status.in_transit')}
                </span>
                <span className="font-bold">
                  {stats.status_distribution.in_transit}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="badge badge-success">
                  {t('status.delivered')}
                </span>
                <span className="font-bold">
                  {stats.status_distribution.delivered}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="badge badge-error">
                  {t('status.cancelled')}
                </span>
                <span className="font-bold">
                  {stats.status_distribution.cancelled}
                </span>
              </div>
              <div className="flex justify-between items-center">
                <span className="badge badge-warning">
                  {t('status.returned')}
                </span>
                <span className="font-bold">
                  {stats.status_distribution.returned}
                </span>
              </div>
            </div>
          </div>
        </div>

        {/* Производительность курьеров */}
        <div className="card bg-base-100 shadow-xl lg:col-span-2">
          <div className="card-body">
            <h2 className="card-title mb-4">
              {t('stats.courier_performance')}
            </h2>
            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>{t('table.courier_service')}</th>
                    <th>{t('table.shipments')}</th>
                    <th>{t('table.delivered')}</th>
                    <th>{t('table.success_rate')}</th>
                    <th>{t('table.avg_time')}</th>
                  </tr>
                </thead>
                <tbody>
                  {stats.courier_performance.map((courier) => (
                    <tr key={courier.name}>
                      <td className="font-medium">{courier.name}</td>
                      <td>{courier.shipments}</td>
                      <td>{courier.delivered}</td>
                      <td>
                        <div className="flex items-center gap-2">
                          <progress
                            className="progress progress-success w-20"
                            value={courier.success_rate}
                            max="100"
                          ></progress>
                          <span className="text-sm">
                            {courier.success_rate.toFixed(1)}%
                          </span>
                        </div>
                      </td>
                      <td>{courier.avg_time.toFixed(1)}ч</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>

      {/* Быстрые действия */}
      <div className="mt-8">
        <h2 className="text-xl font-bold mb-4">{t('quickActions.title')}</h2>
        <div className="flex flex-wrap gap-4">
          <button
            className="btn btn-primary"
            onClick={() => router.push(`/${locale}/admin/logistics/shipments`)}
          >
            {t('view_all_shipments')}
          </button>
          <button
            className="btn btn-warning"
            onClick={() => router.push(`/${locale}/admin/logistics/problems`)}
          >
            {t('problem_shipments')} ({stats.problem_shipments})
          </button>
          <button className="btn btn-secondary">{t('export_report')}</button>
          <button className="btn btn-ghost">{t('monitoring_settings')}</button>
        </div>
      </div>
    </div>
  );
}
