'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import {
  TrashIcon,
  ArrowPathIcon,
  ExclamationTriangleIcon,
  CheckCircleIcon,
  CircleStackIcon,
  ChartBarIcon,
} from '@heroicons/react/24/outline';

interface DataTypeStats {
  count: number;
  size_mb: string;
  size_bytes: number;
  by_type?: Record<string, number>;
}

interface TestDataStats {
  test_logs: DataTypeStats;
  test_results: DataTypeStats;
  behavior_events: DataTypeStats;
  category_feedback: DataTypeStats;
  price_history: DataTypeStats;
  ai_decisions: DataTypeStats;
  total_size_mb: string;
  total_size_bytes: number;
  database_size_mb: string;
  collected_at: string;
}

interface CleanupResponse {
  success: boolean;
  deleted_count: Record<string, number>;
  freed_space_mb: string;
  message: string;
  cleaned_at: string;
}

interface DataType {
  id: string;
  label: string;
  description: string;
  icon: string;
  color: string;
}

export default function TestDataCleanupClient() {
  const t = useTranslations('admin.test_data_cleanup');
  const [stats, setStats] = useState<TestDataStats | null>(null);
  const [loading, setLoading] = useState(true);
  const [cleaning, setCleaning] = useState(false);
  const [selectedTypes, setSelectedTypes] = useState<Set<string>>(new Set());
  const [showConfirm, setShowConfirm] = useState(false);
  const [cleanupResult, setCleanupResult] = useState<CleanupResponse | null>(
    null
  );

  const dataTypes: DataType[] = [
    {
      id: 'test_runs',
      label: t('data_types.test_runs.label'),
      description: t('data_types.test_runs.description'),
      icon: '🏃',
      color: 'indigo',
    },
    {
      id: 'logs',
      label: t('data_types.logs.label'),
      description: t('data_types.logs.description'),
      icon: '📋',
      color: 'blue',
    },
    {
      id: 'results',
      label: t('data_types.results.label'),
      description: t('data_types.results.description'),
      icon: '✅',
      color: 'green',
    },
    {
      id: 'behavior',
      label: t('data_types.behavior.label'),
      description: t('data_types.behavior.description'),
      icon: '👤',
      color: 'purple',
    },
    {
      id: 'feedback',
      label: t('data_types.feedback.label'),
      description: t('data_types.feedback.description'),
      icon: '💬',
      color: 'yellow',
    },
    {
      id: 'price_history',
      label: t('data_types.price_history.label'),
      description: t('data_types.price_history.description'),
      icon: '💰',
      color: 'orange',
    },
    {
      id: 'ai_decisions',
      label: t('data_types.ai_decisions.label'),
      description: t('data_types.ai_decisions.description'),
      icon: '🤖',
      color: 'pink',
    },
  ];

  const fetchStats = async () => {
    try {
      setLoading(true);
      const response = await apiClient.get<TestDataStats>(
        '/admin/tests/data/stats'
      );
      setStats(response.data);
    } catch (error) {
      console.error('Failed to fetch stats:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchStats();
  }, []);

  const toggleType = (typeId: string) => {
    const newSelected = new Set(selectedTypes);
    if (newSelected.has(typeId)) {
      newSelected.delete(typeId);
    } else {
      newSelected.add(typeId);
    }
    setSelectedTypes(newSelected);
  };

  const selectAll = () => {
    setSelectedTypes(new Set(dataTypes.map((t) => t.id)));
  };

  const deselectAll = () => {
    setSelectedTypes(new Set());
  };

  const handleCleanup = async () => {
    if (selectedTypes.size === 0) return;

    try {
      setCleaning(true);
      const typesParam = Array.from(selectedTypes).join(',');
      const response = await apiClient.delete<CleanupResponse>(
        `/admin/tests/data/cleanup?types=${typesParam}`
      );
      setCleanupResult(response.data);
      setShowConfirm(false);
      setSelectedTypes(new Set());
      // Refresh stats
      await fetchStats();
    } catch (error) {
      console.error('Failed to cleanup:', error);
    } finally {
      setCleaning(false);
    }
  };

  const getStatsForType = (typeId: string): DataTypeStats | null => {
    if (!stats) return null;
    switch (typeId) {
      case 'test_runs':
        return (stats as any).test_runs;
      case 'logs':
        return stats.test_logs;
      case 'results':
        return stats.test_results;
      case 'behavior':
        return stats.behavior_events;
      case 'feedback':
        return stats.category_feedback;
      case 'price_history':
        return stats.price_history;
      case 'ai_decisions':
        return stats.ai_decisions;
      default:
        return null;
    }
  };

  const formatNumber = (num: number) => {
    return new Intl.NumberFormat().format(num);
  };

  if (loading) {
    return (
      <div className="flex h-96 items-center justify-center">
        <ArrowPathIcon className="h-8 w-8 animate-spin text-blue-500" />
      </div>
    );
  }

  return (
    <div className="space-y-6 p-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            {t('title')}
          </h1>
          <p className="mt-2 text-gray-600 dark:text-gray-400">
            {t('description')}
          </p>
        </div>
        <button
          onClick={fetchStats}
          className="flex items-center gap-2 rounded-lg bg-blue-500 px-4 py-2 text-white hover:bg-blue-600"
        >
          <ArrowPathIcon className="h-5 w-5" />
          {t('refresh')}
        </button>
      </div>

      {/* Overall Stats */}
      {stats && (
        <div className="grid grid-cols-1 gap-4 md:grid-cols-3">
          <div className="rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 p-6 text-white shadow-lg">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-blue-100">
                  {t('stats.total_size')}
                </p>
                <p className="mt-1 text-3xl font-bold">{stats.total_size_mb}</p>
              </div>
              <CircleStackIcon className="h-12 w-12 text-blue-200" />
            </div>
          </div>

          <div className="rounded-lg bg-gradient-to-br from-green-500 to-green-600 p-6 text-white shadow-lg">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-green-100">
                  {t('stats.database_size')}
                </p>
                <p className="mt-1 text-3xl font-bold">
                  {stats.database_size_mb}
                </p>
              </div>
              <ChartBarIcon className="h-12 w-12 text-green-200" />
            </div>
          </div>

          <div className="rounded-lg bg-gradient-to-br from-purple-500 to-purple-600 p-6 text-white shadow-lg">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-purple-100">
                  {t('stats.total_records')}
                </p>
                <p className="mt-1 text-3xl font-bold">
                  {formatNumber(
                    stats.test_logs.count +
                      stats.test_results.count +
                      stats.behavior_events.count +
                      stats.category_feedback.count +
                      stats.price_history.count +
                      stats.ai_decisions.count
                  )}
                </p>
              </div>
              <TrashIcon className="h-12 w-12 text-purple-200" />
            </div>
          </div>
        </div>
      )}

      {/* Cleanup Result */}
      {cleanupResult && (
        <div className="rounded-lg bg-green-50 p-4 dark:bg-green-900/20">
          <div className="flex items-start">
            <CheckCircleIcon className="h-6 w-6 text-green-500" />
            <div className="ml-3">
              <h3 className="text-sm font-medium text-green-800 dark:text-green-200">
                {t('cleanup_success')}
              </h3>
              <div className="mt-2 text-sm text-green-700 dark:text-green-300">
                <p>{cleanupResult.message}</p>
                <p className="mt-1">
                  {t('freed_space')}:{' '}
                  <strong>{cleanupResult.freed_space_mb}</strong>
                </p>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Data Types Grid */}
      <div>
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-xl font-semibold text-gray-900 dark:text-white">
            {t('select_data_types')}
          </h2>
          <div className="flex gap-2">
            <button
              onClick={selectAll}
              className="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400"
            >
              {t('select_all')}
            </button>
            <span className="text-gray-400">|</span>
            <button
              onClick={deselectAll}
              className="text-sm text-blue-600 hover:text-blue-700 dark:text-blue-400"
            >
              {t('deselect_all')}
            </button>
          </div>
        </div>

        <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-3">
          {dataTypes.map((type) => {
            const typeStats = getStatsForType(type.id);
            const isSelected = selectedTypes.has(type.id);

            return (
              <div
                key={type.id}
                onClick={() => toggleType(type.id)}
                className={`cursor-pointer rounded-lg border-2 p-4 transition-all ${
                  isSelected
                    ? 'border-blue-500 bg-blue-50 dark:bg-blue-900/20'
                    : 'border-gray-200 bg-white hover:border-gray-300 dark:border-gray-700 dark:bg-gray-800'
                }`}
              >
                <div className="flex items-start justify-between">
                  <div className="flex-1">
                    <div className="flex items-center gap-2">
                      <span className="text-2xl">{type.icon}</span>
                      <h3 className="font-semibold text-gray-900 dark:text-white">
                        {type.label}
                      </h3>
                    </div>
                    <p className="mt-1 text-sm text-gray-600 dark:text-gray-400">
                      {type.description}
                    </p>

                    {typeStats && (
                      <div className="mt-3 space-y-1">
                        <div className="flex justify-between text-sm">
                          <span className="text-gray-500 dark:text-gray-400">
                            {t('records')}:
                          </span>
                          <span className="font-medium text-gray-900 dark:text-white">
                            {formatNumber(typeStats.count)}
                          </span>
                        </div>
                        <div className="flex justify-between text-sm">
                          <span className="text-gray-500 dark:text-gray-400">
                            {t('size')}:
                          </span>
                          <span className="font-medium text-gray-900 dark:text-white">
                            {typeStats.size_mb}
                          </span>
                        </div>

                        {/* Behavior events breakdown */}
                        {type.id === 'behavior' && typeStats.by_type && (
                          <div className="mt-2 border-t pt-2 dark:border-gray-700">
                            <p className="mb-1 text-xs font-medium text-gray-500 dark:text-gray-400">
                              {t('event_types')}:
                            </p>
                            {Object.entries(typeStats.by_type).map(
                              ([eventType, count]) => (
                                <div
                                  key={eventType}
                                  className="flex justify-between text-xs"
                                >
                                  <span className="text-gray-500 dark:text-gray-400">
                                    {eventType}:
                                  </span>
                                  <span className="text-gray-700 dark:text-gray-300">
                                    {formatNumber(count)}
                                  </span>
                                </div>
                              )
                            )}
                          </div>
                        )}
                      </div>
                    )}
                  </div>

                  <input
                    type="checkbox"
                    checked={isSelected}
                    onChange={() => {}}
                    className="h-5 w-5 rounded border-gray-300 text-blue-600 focus:ring-blue-500"
                  />
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Actions */}
      <div className="flex items-center justify-between rounded-lg bg-gray-50 p-4 dark:bg-gray-800">
        <div>
          <p className="text-sm text-gray-600 dark:text-gray-400">
            {selectedTypes.size > 0
              ? t('selected_count', { count: selectedTypes.size })
              : t('no_selection')}
          </p>
        </div>
        <button
          onClick={() => setShowConfirm(true)}
          disabled={selectedTypes.size === 0 || cleaning}
          className="flex items-center gap-2 rounded-lg bg-red-600 px-6 py-3 font-medium text-white hover:bg-red-700 disabled:cursor-not-allowed disabled:bg-gray-400"
        >
          {cleaning ? (
            <>
              <ArrowPathIcon className="h-5 w-5 animate-spin" />
              {t('cleaning')}
            </>
          ) : (
            <>
              <TrashIcon className="h-5 w-5" />
              {t('cleanup_button')}
            </>
          )}
        </button>
      </div>

      {/* Confirmation Modal */}
      {showConfirm && (
        <div className="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50">
          <div className="w-full max-w-md rounded-lg bg-white p-6 shadow-xl dark:bg-gray-800">
            <div className="flex items-start">
              <ExclamationTriangleIcon className="h-6 w-6 text-red-500" />
              <div className="ml-3 flex-1">
                <h3 className="text-lg font-medium text-gray-900 dark:text-white">
                  {t('confirm_title')}
                </h3>
                <p className="mt-2 text-sm text-gray-600 dark:text-gray-400">
                  {t('confirm_message')}
                </p>
                <ul className="mt-3 list-inside list-disc text-sm text-gray-600 dark:text-gray-400">
                  {Array.from(selectedTypes).map((typeId) => {
                    const type = dataTypes.find((t) => t.id === typeId);
                    return type ? <li key={typeId}>{type.label}</li> : null;
                  })}
                </ul>
              </div>
            </div>

            <div className="mt-6 flex gap-3">
              <button
                onClick={() => setShowConfirm(false)}
                className="flex-1 rounded-lg border border-gray-300 px-4 py-2 text-gray-700 hover:bg-gray-50 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-700"
              >
                {t('cancel')}
              </button>
              <button
                onClick={handleCleanup}
                className="flex-1 rounded-lg bg-red-600 px-4 py-2 text-white hover:bg-red-700"
              >
                {t('confirm_cleanup')}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Info */}
      <div className="rounded-lg border border-blue-200 bg-blue-50 p-4 dark:border-blue-800 dark:bg-blue-900/20">
        <p className="text-sm text-blue-800 dark:text-blue-200">
          {t('info_message')}
        </p>
      </div>
    </div>
  );
}
