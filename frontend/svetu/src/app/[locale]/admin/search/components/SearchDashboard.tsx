'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from 'react-hot-toast';
import { tokenManager } from '@/utils/tokenManager';
import Pagination from '@/components/admin/Pagination';

interface DashboardStats {
  totalSearches: number;
  avgResponseTime: number;
  topQueries: Array<{
    query: string;
    count: number;
    avgResultsCount: number;
    avgClickPosition: number;
    lastSearched: string;
  }>;
  categoryDistribution?: Array<{
    category: string;
    count: number;
  }>;
  metrics?: {
    totalSearches: number;
    uniqueQueries: number;
    avgSearchTime: number;
    zeroResultsRate: number;
    clickThroughRate: number;
  };
  popularSearches?: Array<{
    query: string;
    count: number;
    avgResults: number;
    lastSearched: string;
  }>;
  deviceStats?: {
    [key: string]: number;
  };
  totalTopQueries?: number;
}

export default function SearchDashboard() {
  const t = useTranslations();
  const [stats, setStats] = useState<DashboardStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [currentPage, setCurrentPage] = useState(1);
  const [itemsPerPage, setItemsPerPage] = useState(25);

  useEffect(() => {
    fetchDashboardStats();
  }, [currentPage, itemsPerPage]);

  const fetchDashboardStats = async () => {
    try {
      setLoading(true);
      const accessToken = await tokenManager.getAccessToken();
      const offset = (currentPage - 1) * itemsPerPage;
      const response = await fetch(`/api/admin/search/analytics?offset=${offset}&limit=${itemsPerPage}`, {
        headers: {
          Authorization: `Bearer ${accessToken}`,
        },
      });
      if (!response.ok) {
        throw new Error('Failed to fetch dashboard stats');
      }
      const result = await response.json();
      // Извлекаем данные из обертки
      const data = result.data || result;
      setStats(data);
    } catch (error) {
      console.error('Error fetching dashboard stats:', error);
      toast.error(t('admin.search.dashboard.error'));
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  if (!stats) {
    return (
      <div className="alert alert-warning">
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
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.98-.833-2.75 0L3.098 16.5c-.77.833.192 2.5 1.732 2.5z"
          />
        </svg>
        <span>{t('admin.search.dashboard.noData')}</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        <div className="stat bg-base-100 rounded-box">
          <div className="stat-title">
            {t('admin.search.dashboard.totalSearches')}
          </div>
          <div className="stat-value">
            {(stats.totalSearches || 0).toLocaleString()}
          </div>
          <div className="stat-desc">
            {t('admin.search.dashboard.last30Days')}
          </div>
        </div>

        <div className="stat bg-base-100 rounded-box">
          <div className="stat-title">
            {t('admin.search.dashboard.avgResponseTime')}
          </div>
          <div className="stat-value">
            {(stats.avgResponseTime || 0).toFixed(0)}ms
          </div>
          <div className="stat-desc">
            {t('admin.search.dashboard.searchPerformance')}
          </div>
        </div>

        <div className="stat bg-base-100 rounded-box">
          <div className="stat-title">
            {t('admin.search.dashboard.topQueries')}
          </div>
          <div className="stat-value">{stats.topQueries?.length || 0}</div>
          <div className="stat-desc">
            {t('admin.search.dashboard.uniqueQueries')}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">
              {t('admin.search.dashboard.popularQueries')}
            </h2>
            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    <th>{t('admin.search.dashboard.query')}</th>
                    <th>{t('admin.search.dashboard.count')}</th>
                    <th>{t('admin.search.dashboard.relevance')}</th>
                  </tr>
                </thead>
                <tbody>
                  {(stats.topQueries || []).map((query, index) => (
                    <tr key={index}>
                      <td className="font-mono">{query.query}</td>
                      <td>{query.count}</td>
                      <td>
                        <div className="badge badge-primary">
                          {(query.avgResultsCount || 0).toFixed(1)}
                        </div>
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            {stats.totalTopQueries && stats.totalTopQueries > itemsPerPage && (
              <Pagination
                currentPage={currentPage}
                totalPages={Math.ceil(stats.totalTopQueries / itemsPerPage)}
                totalItems={stats.totalTopQueries}
                itemsPerPage={itemsPerPage}
                onPageChange={(page) => {
                  setCurrentPage(page);
                }}
                onItemsPerPageChange={(items) => {
                  setItemsPerPage(items);
                  setCurrentPage(1);
                }}
              />
            )}
          </div>
        </div>

        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h2 className="card-title">
              {t('admin.search.dashboard.deviceDistribution')}
            </h2>
            {stats.deviceStats && Object.keys(stats.deviceStats).length > 0 ? (
              <div className="space-y-2">
                {Object.entries(stats.deviceStats).map(([device, count], index) => {
                  const total = Object.values(stats.deviceStats || {}).reduce((sum, val) => sum + val, 0);
                  const percentage = (count / total) * 100;
                  return (
                    <div key={index} className="flex items-center justify-between">
                      <span className="text-sm capitalize">{device || 'Unknown'}</span>
                      <div className="flex items-center gap-2">
                        <div className="w-32 bg-base-300 rounded-full h-2">
                          <div
                            className="bg-primary h-2 rounded-full"
                            style={{ width: `${percentage}%` }}
                          ></div>
                        </div>
                        <span className="text-xs text-base-content/70">
                          {count}
                        </span>
                      </div>
                    </div>
                  );
                })}
              </div>
            ) : (
              <p className="text-sm text-base-content/70">
                {t('admin.search.dashboard.noDeviceData')}
              </p>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
