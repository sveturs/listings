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
} from 'chart.js';
import { Line, Bar } from 'react-chartjs-2';
import {
  searchAnalyticsService,
  SearchMetrics,
  ItemPerformance,
  AnalyticsFilters,
} from '@/services/searchAnalytics';

// Регистрируем компоненты Chart.js
ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  BarElement,
  Title,
  Tooltip,
  Legend
);

export default function BehaviorAnalytics() {
  const t = useTranslations('admin.search.behavior');
  const tCommon = useTranslations('common');

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [metrics, setMetrics] = useState<SearchMetrics | null>(null);
  const [itemsPerformance, setItemsPerformance] = useState<ItemPerformance[]>(
    []
  );
  const [filters, setFilters] = useState<AnalyticsFilters>({
    period: 'week',
    limit: 20,
  });

  // Состояния для управления видимостью секций
  const [expandedSections, setExpandedSections] = useState({
    ctr: true,
    topQueries: true,
    itemPerformance: true,
    conversion: true,
  });

  useEffect(() => {
    loadData();
  }, [filters]); // eslint-disable-line react-hooks/exhaustive-deps

  const loadData = async () => {
    try {
      setLoading(true);
      setError(null);

      const [metricsData, itemsData] = await Promise.all([
        searchAnalyticsService.getSearchMetrics(filters),
        searchAnalyticsService.getItemsPerformance(filters),
      ]);

      setMetrics(metricsData);
      setItemsPerformance(itemsData);
    } catch (err) {
      console.error('Failed to load analytics data:', err);
      setError(t('loadError'));
    } finally {
      setLoading(false);
    }
  };

  const toggleSection = (section: keyof typeof expandedSections) => {
    setExpandedSections((prev) => ({
      ...prev,
      [section]: !prev[section],
    }));
  };

  // Данные для графика CTR по времени
  const ctrChartData = {
    labels:
      metrics?.search_trends?.map((trend) =>
        new Date(trend.date).toLocaleDateString('ru-RU', {
          month: 'short',
          day: 'numeric',
        })
      ) || [],
    datasets: [
      {
        label: t('ctr.label'),
        data: metrics?.search_trends?.map((trend) => trend.ctr * 100) || [],
        borderColor: 'rgb(59, 130, 246)',
        backgroundColor: 'rgba(59, 130, 246, 0.1)',
        borderWidth: 3,
        fill: true,
        tension: 0.4,
        pointBackgroundColor: 'rgb(59, 130, 246)',
        pointBorderColor: '#fff',
        pointBorderWidth: 2,
        pointRadius: 6,
        pointHoverRadius: 8,
      },
    ],
  };

  const ctrChartOptions = {
    responsive: true,
    plugins: {
      legend: {
        position: 'top' as const,
      },
      title: {
        display: true,
        text: t('ctr.chartTitle'),
        font: {
          size: 16,
          weight: 'bold' as const,
        },
      },
    },
    scales: {
      y: {
        beginAtZero: true,
        max: 100,
        ticks: {
          callback: function (value: any) {
            return value + '%';
          },
        },
      },
    },
  };

  // Данные для графика кликов и поисков
  const searchClicksData = {
    labels:
      metrics?.search_trends?.map((trend) =>
        new Date(trend.date).toLocaleDateString('ru-RU', {
          month: 'short',
          day: 'numeric',
        })
      ) || [],
    datasets: [
      {
        label: t('searches'),
        data:
          metrics?.search_trends?.map((trend) => trend.searches_count) || [],
        backgroundColor: 'rgba(34, 197, 94, 0.8)',
        borderColor: 'rgb(34, 197, 94)',
        borderWidth: 1,
      },
      {
        label: t('clicks'),
        data: metrics?.search_trends?.map((trend) => trend.clicks_count) || [],
        backgroundColor: 'rgba(239, 68, 68, 0.8)',
        borderColor: 'rgb(239, 68, 68)',
        borderWidth: 1,
      },
    ],
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center py-20">
        <span className="loading loading-spinner loading-lg text-primary"></span>
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
        <button className="btn btn-sm" onClick={loadData}>
          {tCommon('tryAgain')}
        </button>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Фильтры периода */}
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <h3 className="card-title text-lg mb-4">{t('filters.title')}</h3>
          <div className="flex flex-wrap gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('filters.period')}</span>
              </label>
              <select
                className="select select-bordered"
                value={filters.period}
                onChange={(e) =>
                  setFilters((prev) => ({
                    ...prev,
                    period: e.target.value as any,
                  }))
                }
              >
                <option value="day">{t('filters.periods.day')}</option>
                <option value="week">{t('filters.periods.week')}</option>
                <option value="month">{t('filters.periods.month')}</option>
              </select>
            </div>
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('filters.limit')}</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={filters.limit || 20}
                onChange={(e) =>
                  setFilters((prev) => ({
                    ...prev,
                    limit: Number(e.target.value),
                  }))
                }
                min="5"
                max="100"
              />
            </div>
          </div>
        </div>
      </div>

      {/* Общие метрики */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <div className="stat bg-primary text-primary-content rounded-lg">
          <div className="stat-figure">
            <svg className="w-8 h-8" fill="currentColor" viewBox="0 0 20 20">
              <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <div className="stat-title text-primary-content/80">
            {t('metrics.totalClicks')}
          </div>
          <div className="stat-value text-2xl">
            {metrics?.click_metrics?.total_clicks || 0}
          </div>
        </div>

        <div className="stat bg-secondary text-secondary-content rounded-lg">
          <div className="stat-figure">
            <svg className="w-8 h-8" fill="currentColor" viewBox="0 0 20 20">
              <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z" />
            </svg>
          </div>
          <div className="stat-title text-secondary-content/80">
            {t('metrics.avgCTR')}
          </div>
          <div className="stat-value text-2xl">
            {((metrics?.click_metrics?.ctr || 0) * 100).toFixed(1)}%
          </div>
        </div>

        <div className="stat bg-accent text-accent-content rounded-lg">
          <div className="stat-figure">
            <svg className="w-8 h-8" fill="currentColor" viewBox="0 0 20 20">
              <path d="M3 4a1 1 0 011-1h12a1 1 0 011 1v2a1 1 0 01-1 1H4a1 1 0 01-1-1V4zM3 10a1 1 0 011-1h6a1 1 0 011 1v6a1 1 0 01-1 1H4a1 1 0 01-1-1v-6zM14 9a1 1 0 00-1 1v6a1 1 0 001 1h2a1 1 0 001-1v-6a1 1 0 00-1-1h-2z" />
            </svg>
          </div>
          <div className="stat-title text-accent-content/80">
            {t('metrics.avgPosition')}
          </div>
          <div className="stat-value text-2xl">
            {(metrics?.click_metrics?.average_click_position || 0).toFixed(1)}
          </div>
        </div>

        <div className="stat bg-info text-info-content rounded-lg">
          <div className="stat-figure">
            <svg className="w-8 h-8" fill="currentColor" viewBox="0 0 20 20">
              <path
                fillRule="evenodd"
                d="M3 3a1 1 0 000 2v8a2 2 0 002 2h2.586l-1.293 1.293a1 1 0 101.414 1.414L10 15.414l2.293 2.293a1 1 0 001.414-1.414L12.414 15H15a2 2 0 002-2V5a1 1 0 100-2H3zm11.707 4.707a1 1 0 00-1.414-1.414L10 9.586 8.707 8.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                clipRule="evenodd"
              />
            </svg>
          </div>
          <div className="stat-title text-info-content/80">
            {t('metrics.conversionRate')}
          </div>
          <div className="stat-value text-2xl">
            {((metrics?.click_metrics?.conversion_rate || 0) * 100).toFixed(1)}%
          </div>
        </div>
      </div>

      {/* CTR по времени */}
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h3 className="card-title">{t('ctr.title')}</h3>
            <button
              className="btn btn-ghost btn-sm"
              onClick={() => toggleSection('ctr')}
            >
              {expandedSections.ctr ? '−' : '+'}
            </button>
          </div>
          {expandedSections.ctr && (
            <div className="h-80">
              <Line data={ctrChartData} options={ctrChartOptions} />
            </div>
          )}
        </div>
      </div>

      {/* График поисков и кликов */}
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h3 className="card-title">{t('searchClicks.title')}</h3>
            <button
              className="btn btn-ghost btn-sm"
              onClick={() => toggleSection('ctr')}
            >
              {expandedSections.ctr ? '−' : '+'}
            </button>
          </div>
          {expandedSections.ctr && (
            <div className="h-80">
              <Bar
                data={searchClicksData}
                options={{
                  responsive: true,
                  plugins: {
                    legend: {
                      position: 'top' as const,
                    },
                    title: {
                      display: true,
                      text: t('searchClicks.chartTitle'),
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
          )}
        </div>
      </div>

      {/* Топ запросы по CTR */}
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h3 className="card-title">{t('topQueries.title')}</h3>
            <button
              className="btn btn-ghost btn-sm"
              onClick={() => toggleSection('topQueries')}
            >
              {expandedSections.topQueries ? '−' : '+'}
            </button>
          </div>
          {expandedSections.topQueries && (
            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>{t('topQueries.query')}</th>
                    <th>{t('topQueries.searches')}</th>
                    <th>{t('topQueries.ctr')}</th>
                    <th>{t('topQueries.avgPosition')}</th>
                    <th>{t('topQueries.avgResults')}</th>
                  </tr>
                </thead>
                <tbody>
                  {metrics?.top_queries
                    ?.slice(0, filters.limit || 20)
                    .map((query, index) => (
                      <tr key={index}>
                        <td className="font-semibold">{query.query}</td>
                        <td>{query.count}</td>
                        <td>
                          <div className="flex items-center gap-2">
                            <div className="badge badge-primary">
                              {(query.ctr * 100).toFixed(1)}%
                            </div>
                          </div>
                        </td>
                        <td>{query.avg_position?.toFixed(1) || 'N/A'}</td>
                        <td>{query.avg_results}</td>
                      </tr>
                    )) || []}
                </tbody>
              </table>
              {(!metrics?.top_queries || metrics.top_queries.length === 0) && (
                <div className="text-center py-8 text-base-content/60">
                  {tCommon('noData')}
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Производительность товаров */}
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h3 className="card-title">{t('itemPerformance.title')}</h3>
            <button
              className="btn btn-ghost btn-sm"
              onClick={() => toggleSection('itemPerformance')}
            >
              {expandedSections.itemPerformance ? '−' : '+'}
            </button>
          </div>
          {expandedSections.itemPerformance && (
            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>{t('itemPerformance.item')}</th>
                    <th>{t('itemPerformance.impressions')}</th>
                    <th>{t('itemPerformance.clicks')}</th>
                    <th>{t('itemPerformance.ctr')}</th>
                    <th>{t('itemPerformance.avgPosition')}</th>
                    <th>{t('itemPerformance.conversions')}</th>
                  </tr>
                </thead>
                <tbody>
                  {itemsPerformance
                    .slice(0, filters.limit || 20)
                    .map((item, index) => (
                      <tr key={index}>
                        <td
                          className="max-w-xs truncate font-semibold"
                          title={item.item_title}
                        >
                          {item.item_title}
                        </td>
                        <td>{item.impressions}</td>
                        <td>{item.clicks}</td>
                        <td>
                          <div className="badge badge-primary">
                            {(item.ctr * 100).toFixed(1)}%
                          </div>
                        </td>
                        <td>{item.average_position?.toFixed(1) || 'N/A'}</td>
                        <td>{item.conversions}</td>
                      </tr>
                    ))}
                </tbody>
              </table>
              {itemsPerformance.length === 0 && (
                <div className="text-center py-8 text-base-content/60">
                  {tCommon('noData')}
                </div>
              )}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
