'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';
import { tokenManager } from '@/utils/tokenManager';
import Pagination from '@/components/admin/Pagination';

interface SearchQuery {
  query: string;
  count: number;
  avgResultsCount: number;
  avgClickPosition: number;
  lastSearched: string;
}

interface SearchMetrics {
  totalSearches: number;
  uniqueQueries: number;
  avgSearchTime: number;
  zeroResultsRate: number;
  clickThroughRate: number;
}

interface TimeRange {
  label: string;
  value: '24h' | '7d' | '30d' | '90d';
}

export default function SearchAnalytics() {
  const t = useTranslations();
  const [metrics, setMetrics] = useState<SearchMetrics | null>(null);
  const [topQueries, setTopQueries] = useState<SearchQuery[]>([]);
  const [zeroResultQueries, setZeroResultQueries] = useState<SearchQuery[]>([]);
  const [timeRange, setTimeRange] = useState<TimeRange['value']>('7d');
  const [loading, setLoading] = useState(true);
  const [currentPageTop, setCurrentPageTop] = useState(1);
  const [currentPageZero, setCurrentPageZero] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(25);
  const [totalTopQueries, setTotalTopQueries] = useState(0);
  const [totalZeroQueries, setTotalZeroQueries] = useState(0);

  const timeRanges: TimeRange[] = [
    { label: t('admin.search.analytics.last24h'), value: '24h' },
    { label: t('admin.search.analytics.last7d'), value: '7d' },
    { label: t('admin.search.analytics.last30d'), value: '30d' },
    { label: t('admin.search.analytics.last90d'), value: '90d' },
  ];

  useEffect(() => {
    fetchAnalytics();
  }, [timeRange, currentPageTop, currentPageZero, itemsPerPage]);

  const fetchAnalytics = async () => {
    setLoading(true);
    try {
      const accessToken = await tokenManager.getAccessToken();
      const offsetTop = (currentPageTop - 1) * itemsPerPage;
      const offsetZero = (currentPageZero - 1) * itemsPerPage;
      const response = await fetch(
        `/api/v1/admin/search/analytics?range=${timeRange}&offsetTop=${offsetTop}&offsetZero=${offsetZero}&limit=${itemsPerPage}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        }
      );
      if (!response.ok) throw new Error('Failed to fetch analytics');
      const result = await response.json();
      const data = result.data || {};

      setMetrics(
        data.metrics || {
          totalSearches: 0,
          uniqueQueries: 0,
          avgSearchTime: 0,
          zeroResultsRate: 0,
          clickThroughRate: 0,
        }
      );
      setTopQueries(data.topQueries || []);
      setZeroResultQueries(data.zeroResultQueries || []);
      setTotalTopQueries(data.totalTopQueries || data.topQueries?.length || 0);
      setTotalZeroQueries(data.totalZeroQueries || data.zeroResultQueries?.length || 0);
    } catch (error) {
      console.error('Error fetching analytics:', error);
      toast.error(t('admin.search.analytics.fetchError'));
    } finally {
      setLoading(false);
    }
  };

  const exportData = async () => {
    try {
      const accessToken = await tokenManager.getAccessToken();
      const response = await fetch(
        `/api/v1/admin/search/analytics/export?range=${timeRange}`,
        {
          headers: {
            Authorization: `Bearer ${accessToken}`,
          },
        }
      );
      if (!response.ok) throw new Error('Failed to export data');

      const blob = await response.blob();
      const url = window.URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `search-analytics-${timeRange}.csv`;
      document.body.appendChild(a);
      a.click();
      window.URL.revokeObjectURL(url);
      document.body.removeChild(a);

      toast.success(t('admin.search.analytics.exportSuccess'));
    } catch (error) {
      console.error('Error exporting data:', error);
      toast.error(t('admin.search.analytics.exportError'));
    }
  };

  if (loading) {
    return <div className="loading loading-spinner loading-lg"></div>;
  }

  return (
    <div className="space-y-6">
      <div className="flex justify-between items-center">
        <div className="tabs tabs-boxed">
          {timeRanges.map((range) => (
            <button
              key={range.value}
              className={`tab ${timeRange === range.value ? 'tab-active' : ''}`}
              onClick={() => setTimeRange(range.value)}
            >
              {range.label}
            </button>
          ))}
        </div>
        <button className="btn btn-outline btn-sm" onClick={exportData}>
          {t('admin.search.analytics.export')}
        </button>
      </div>

      {metrics && (
        <div className="stats shadow">
          <div className="stat">
            <div className="stat-title">
              {t('admin.search.analytics.totalSearches')}
            </div>
            <div className="stat-value">
              {metrics.totalSearches.toLocaleString()}
            </div>
          </div>
          <div className="stat">
            <div className="stat-title">
              {t('admin.search.analytics.uniqueQueries')}
            </div>
            <div className="stat-value">
              {metrics.uniqueQueries.toLocaleString()}
            </div>
          </div>
          <div className="stat">
            <div className="stat-title">
              {t('admin.search.analytics.avgSearchTime')}
            </div>
            <div className="stat-value">
              {metrics.avgSearchTime.toFixed(2)}ms
            </div>
          </div>
          <div className="stat">
            <div className="stat-title">
              {t('admin.search.analytics.zeroResultsRate')}
            </div>
            <div className="stat-value text-error">
              {(metrics.zeroResultsRate * 100).toFixed(1)}%
            </div>
          </div>
          <div className="stat">
            <div className="stat-title">
              {t('admin.search.analytics.clickThroughRate')}
            </div>
            <div className="stat-value text-success">
              {(metrics.clickThroughRate * 100).toFixed(1)}%
            </div>
          </div>
        </div>
      )}

      <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title">
              {t('admin.search.analytics.topQueries')}
            </h3>
            <div className="overflow-x-auto">
              <table className="table table-compact">
                <thead>
                  <tr>
                    <th>{t('admin.search.analytics.query')}</th>
                    <th>{t('admin.search.analytics.count')}</th>
                    <th>{t('admin.search.analytics.avgResults')}</th>
                    <th>{t('admin.search.analytics.avgClickPos')}</th>
                  </tr>
                </thead>
                <tbody>
                  {topQueries.map((query, index) => (
                    <tr key={index}>
                      <td className="font-mono">{query.query}</td>
                      <td>{query.count}</td>
                      <td>{query.avgResultsCount.toFixed(0)}</td>
                      <td>{query.avgClickPosition.toFixed(1)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {totalTopQueries > itemsPerPage && (
              <Pagination
                currentPage={currentPageTop}
                totalPages={Math.ceil(totalTopQueries / itemsPerPage)}
                totalItems={totalTopQueries}
                itemsPerPage={itemsPerPage}
                onPageChange={(page) => {
                  setCurrentPageTop(page);
                }}
                onItemsPerPageChange={(items) => {
                  setItemsPerPage(items);
                  setCurrentPageTop(1);
                  setCurrentPageZero(1);
                }}
              />
            )}
          </div>
        </div>

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title text-error">
              {t('admin.search.analytics.zeroResultQueries')}
            </h3>
            <div className="overflow-x-auto">
              <table className="table table-compact">
                <thead>
                  <tr>
                    <th>{t('admin.search.analytics.query')}</th>
                    <th>{t('admin.search.analytics.count')}</th>
                    <th>{t('admin.search.analytics.lastSearched')}</th>
                  </tr>
                </thead>
                <tbody>
                  {zeroResultQueries.map((query, index) => (
                    <tr key={index}>
                      <td className="font-mono">{query.query}</td>
                      <td>{query.count}</td>
                      <td>
                        {new Date(query.lastSearched).toLocaleDateString()}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {totalZeroQueries > itemsPerPage && (
              <Pagination
                currentPage={currentPageZero}
                totalPages={Math.ceil(totalZeroQueries / itemsPerPage)}
                totalItems={totalZeroQueries}
                itemsPerPage={itemsPerPage}
                onPageChange={(page) => {
                  setCurrentPageZero(page);
                }}
                onItemsPerPageChange={(items) => {
                  setItemsPerPage(items);
                  setCurrentPageTop(1);
                  setCurrentPageZero(1);
                }}
              />
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
