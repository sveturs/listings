'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Legend,
  BarElement,
  ArcElement,
} from 'chart.js';
import { Line, Doughnut } from 'react-chartjs-2';
import {
  searchAnalyticsService,
  SearchMetrics,
  AnalyticsFilters,
} from '@/services/searchAnalytics';

// Регистрируем компоненты Chart.js
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  ArcElement,
  Title,
  Tooltip,
  Legend
);

type TimeRange = '24h' | '7d' | '30d' | '90d';

interface OverviewMetrics {
  totalSearches: number;
  uniqueQueries: number;
  avgResponseTime: number;
  zeroresultRate: number;
  avgCTR: number;
  conversionRate: number;
}

export default function SearchAnalytics() {
  const t = useTranslations('admin');

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [timeRange, setTimeRange] = useState<TimeRange>('7d');
  const [metrics, setMetrics] = useState<SearchMetrics | null>(null);
  const [overviewMetrics, setOverviewMetrics] = useState<OverviewMetrics>({
    totalSearches: 0,
    uniqueQueries: 0,
    avgResponseTime: 0,
    zeroresultRate: 0,
    avgCTR: 0,
    conversionRate: 0,
  });

  useEffect(() => {
    loadData();
  }, [timeRange]); // eslint-disable-line react-hooks/exhaustive-deps

  const getFiltersForTimeRange = (range: TimeRange): AnalyticsFilters => {
    const now = new Date();
    const filters: AnalyticsFilters = {};

    switch (range) {
      case '24h':
        filters.period_start = new Date(
          now.getTime() - 24 * 60 * 60 * 1000
        ).toISOString();
        break;
      case '7d':
        filters.period_start = new Date(
          now.getTime() - 7 * 24 * 60 * 60 * 1000
        ).toISOString();
        break;
      case '30d':
        filters.period_start = new Date(
          now.getTime() - 30 * 24 * 60 * 60 * 1000
        ).toISOString();
        break;
      case '90d':
        filters.period_start = new Date(
          now.getTime() - 90 * 24 * 60 * 60 * 1000
        ).toISOString();
        break;
    }

    filters.period_end = now.toISOString();
    return filters;
  };

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);

      const filters = getFiltersForTimeRange(timeRange);
      const data = await searchAnalyticsService.getSearchMetrics(filters);

      setMetrics(data);

      // Вычисляем обзорные метрики
      const zeroResultQueries =
        data.top_queries?.filter((q) => q.avg_results === 0).length || 0;
      const totalQueriesCount = data.top_queries?.length || 0;
      const zeroResultRate =
        totalQueriesCount > 0
          ? (zeroResultQueries / totalQueriesCount) * 100
          : 0;

      // Вычисляем средний CTR из топ запросов
      const avgCTR =
        data.top_queries && data.top_queries.length > 0
          ? data.top_queries.reduce((sum, q) => sum + (q.ctr || 0), 0) /
            data.top_queries.length
          : 0;

      setOverviewMetrics({
        totalSearches: data.total_searches || 0,
        uniqueQueries: data.unique_searches || 0,
        avgResponseTime: 100, // Временное значение, так как backend пока не возвращает это поле
        zeroresultRate: zeroResultRate,
        avgCTR: avgCTR,
        conversionRate: 0, // Временное значение, так как backend пока не возвращает это поле
      });
    } catch (err) {
      console.error('Failed to load search analytics:', err);
      setError(t('loadError'));
    } finally {
      setLoading(false);
    }
  };

  const handleExportCSV = async () => {
    if (!metrics) return;

    const csvData = [
      ['Метрика', 'Значение'],
      ['Общие поиски', overviewMetrics.totalSearches.toString()],
      ['Уникальные запросы', overviewMetrics.uniqueQueries.toString()],
      ['Среднее время ответа (мс)', overviewMetrics.avgResponseTime.toString()],
      [
        'Процент нулевых результатов',
        `${overviewMetrics.zeroresultRate.toFixed(2)}%`,
      ],
      ['Средний CTR', `${overviewMetrics.avgCTR.toFixed(2)}%`],
      [
        'Коэффициент конверсии',
        `${overviewMetrics.conversionRate.toFixed(2)}%`,
      ],
      [''],
      ['Топ запросы', ''],
      ['Запрос', 'Количество', 'CTR', 'Средняя позиция', 'Средние результаты'],
      ...(metrics.top_queries?.map((q) => [
        q.query,
        q.count.toString(),
        `${(q.ctr || 0).toFixed(2)}%`,
        '1.0',
        q.avg_results.toString(),
      ]) || []),
    ];

    const csvContent = csvData.map((row) => row.join(',')).join('\n');
    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');

    if (link.download !== undefined) {
      const url = URL.createObjectURL(blob);
      link.setAttribute('href', url);
      link.setAttribute('download', `search-analytics-${timeRange}.csv`);
      link.style.visibility = 'hidden';
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    }
  };

  // Данные для графика трендов поиска
  const searchTrendsData = {
    labels:
      metrics?.search_trends?.map((trend) =>
        new Date(trend.date).toLocaleDateString('ru-RU', {
          month: 'short',
          day: 'numeric',
        })
      ) || [],
    datasets: [
      {
        label: 'Поиски',
        data:
          metrics?.search_trends?.map((trend) => trend.searches_count) || [],
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        tension: 0.1,
      },
      {
        label: 'Клики',
        data: metrics?.search_trends?.map((trend) => trend.clicks_count) || [],
        borderColor: 'rgb(34, 197, 94)',
        backgroundColor: 'rgba(34, 197, 94, 0.1)',
        tension: 0.1,
      },
    ],
  };

  // Данные для графика CTR
  const ctrTrendsData = {
    labels:
      metrics?.search_trends?.map((trend) =>
        new Date(trend.date).toLocaleDateString('ru-RU', {
          month: 'short',
          day: 'numeric',
        })
      ) || [],
    datasets: [
      {
        label: 'CTR (%)',
        data:
          metrics?.search_trends?.map((trend) => (trend.ctr || 0).toFixed(2)) ||
          [],
        borderColor: 'rgb(168, 85, 247)',
        backgroundColor: 'rgba(168, 85, 247, 0.1)',
        tension: 0.1,
      },
    ],
  };

  // Данные для топ запросов
  const topQueriesData = {
    labels: metrics?.top_queries?.slice(0, 10).map((q) => q.query) || [],
    datasets: [
      {
        label: 'Количество поисков',
        data: metrics?.top_queries?.slice(0, 10).map((q) => q.count) || [],
        backgroundColor: [
          '#3B82F6',
          '#10B981',
          '#8B5CF6',
          '#F59E0B',
          '#EF4444',
          '#06B6D4',
          '#84CC16',
          '#F97316',
          '#EC4899',
          '#6366F1',
        ],
      },
    ],
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-96">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
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
    );
  }

  return (
    <div className="space-y-6">
      {/* Заголовок и фильтры */}
      <div className="flex flex-col lg:flex-row lg:items-center lg:justify-between">
        <div>
          <h3 className="text-2xl font-bold">{t('title')}</h3>
          <p className="text-base-content/60">{t('description')}</p>
        </div>

        <div className="flex items-center gap-4 mt-4 lg:mt-0">
          {/* Выбор временного диапазона */}
          <div className="form-control">
            <label className="label">
              <span className="label-text">{t('timeRange')}</span>
            </label>
            <select
              className="select select-bordered select-sm"
              value={timeRange}
              onChange={(e) => setTimeRange(e.target.value as TimeRange)}
            >
              <option value="24h">{t('ranges.24h')}</option>
              <option value="7d">{t('ranges.7d')}</option>
              <option value="30d">{t('ranges.30d')}</option>
              <option value="90d">{t('ranges.90d')}</option>
            </select>
          </div>

          {/* Кнопка экспорта */}
          <button
            className="btn btn-outline btn-sm"
            onClick={handleExportCSV}
            disabled={!metrics}
          >
            <svg
              className="w-4 h-4 mr-2"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
              />
            </svg>
            {t('exportCSV')}
          </button>
        </div>
      </div>

      {/* Обзорные метрики */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-4">
        <div className="stat bg-base-100 shadow-md rounded-lg">
          <div className="stat-title">{t('metrics.totalSearches')}</div>
          <div className="stat-value text-primary">
            {overviewMetrics.totalSearches.toLocaleString()}
          </div>
        </div>

        <div className="stat bg-base-100 shadow-md rounded-lg">
          <div className="stat-title">{t('metrics.uniqueQueries')}</div>
          <div className="stat-value text-secondary">
            {overviewMetrics.uniqueQueries.toLocaleString()}
          </div>
        </div>

        <div className="stat bg-base-100 shadow-md rounded-lg">
          <div className="stat-title">{t('metrics.avgResponseTime')}</div>
          <div className="stat-value text-accent">
            {overviewMetrics.avgResponseTime.toFixed(0)} ms
          </div>
        </div>

        <div className="stat bg-base-100 shadow-md rounded-lg">
          <div className="stat-title">{t('metrics.zeroResults')}</div>
          <div className="stat-value text-warning">
            {overviewMetrics.zeroresultRate.toFixed(1)}%
          </div>
        </div>

        <div className="stat bg-base-100 shadow-md rounded-lg">
          <div className="stat-title">{t('metrics.avgCTR')}</div>
          <div className="stat-value text-success">
            {overviewMetrics.avgCTR.toFixed(2)}%
          </div>
        </div>

        <div className="stat bg-base-100 shadow-md rounded-lg">
          <div className="stat-title">{t('metrics.conversionRate')}</div>
          <div className="stat-value text-info">
            {overviewMetrics.conversionRate.toFixed(2)}%
          </div>
        </div>
      </div>

      {/* Графики */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* График трендов поиска */}
        <div className="card bg-base-100 shadow-md">
          <div className="card-body">
            <h4 className="card-title">{t('charts.searchTrends')}</h4>
            <div className="h-64">
              <Line
                data={searchTrendsData}
                options={{
                  responsive: true,
                  maintainAspectRatio: false,
                  plugins: {
                    legend: {
                      position: 'top' as const,
                    },
                  },
                  scales: {
                    y: {
                      beginAtZero: true,
                    },
                  },
                }}
              />
            </div>
          </div>
        </div>

        {/* График CTR */}
        <div className="card bg-base-100 shadow-md">
          <div className="card-body">
            <h4 className="card-title">{t('charts.ctrTrends')}</h4>
            <div className="h-64">
              <Line
                data={ctrTrendsData}
                options={{
                  responsive: true,
                  maintainAspectRatio: false,
                  plugins: {
                    legend: {
                      position: 'top' as const,
                    },
                  },
                  scales: {
                    y: {
                      beginAtZero: true,
                      max: 100,
                    },
                  },
                }}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Топ запросы и анализ релевантности */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Топ запросы - таблица */}
        <div className="card bg-base-100 shadow-md lg:col-span-2">
          <div className="card-body">
            <h4 className="card-title">{t('topQueries.title')}</h4>
            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>{t('topQueries.query')}</th>
                    <th>{t('topQueries.count')}</th>
                    <th>{t('topQueries.ctr')}</th>
                    <th>{t('topQueries.avgPosition')}</th>
                    <th>{t('topQueries.avgResults')}</th>
                    <th>{t('topQueries.relevance')}</th>
                  </tr>
                </thead>
                <tbody>
                  {metrics?.top_queries?.slice(0, 15).map((query, index) => {
                    const relevanceScore =
                      query.avg_results > 0
                        ? Math.min(
                            100,
                            query.ctr * 100 + query.avg_results / 10
                          )
                        : 0;
                    const relevanceColor =
                      relevanceScore > 70
                        ? 'text-success'
                        : relevanceScore > 40
                          ? 'text-warning'
                          : 'text-error';

                    return (
                      <tr key={index}>
                        <td className="font-medium">{query.query}</td>
                        <td>{query.count}</td>
                        <td>{(query.ctr || 0).toFixed(2)}%</td>
                        <td>1.0</td>
                        <td>{query.avg_results}</td>
                        <td className={`font-semibold ${relevanceColor}`}>
                          {relevanceScore.toFixed(0)}%
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          </div>
        </div>

        {/* Топ запросы - график */}
        <div className="card bg-base-100 shadow-md">
          <div className="card-body">
            <h4 className="card-title">{t('charts.topQueries')}</h4>
            <div className="h-64">
              <Doughnut
                data={topQueriesData}
                options={{
                  responsive: true,
                  maintainAspectRatio: false,
                  plugins: {
                    legend: {
                      position: 'bottom' as const,
                      labels: {
                        boxWidth: 12,
                        font: {
                          size: 10,
                        },
                      },
                    },
                  },
                }}
              />
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
