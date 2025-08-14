'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import {
  TranslationAuditLog,
  AuditStatistics,
} from '@/services/translationAdminApi';
import {
  ClockIcon,
  UserIcon,
  DocumentTextIcon,
  ChartBarIcon,
  FunnelIcon,
  ArrowPathIcon,
} from '@heroicons/react/24/outline';

export default function AuditLogViewer() {
  const _t = useTranslations('admin.translations');

  const [logs, setLogs] = useState<TranslationAuditLog[]>([]);
  const [statistics, setStatistics] = useState<AuditStatistics | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [activeTab, setActiveTab] = useState<'logs' | 'stats'>('logs');

  // Filters
  const [filters, setFilters] = useState({
    user_id: '',
    action: '',
    entity_type: '',
    start_date: '',
    end_date: '',
  });

  const ACTIONS = [
    'create',
    'update',
    'delete',
    'translate',
    'approve',
    'rollback',
    'import',
    'export',
  ];

  const ENTITY_TYPES = ['translation', 'category', 'attribute', 'listing'];

  useEffect(() => {
    if (activeTab === 'logs') {
      loadLogs();
    } else {
      loadStatistics();
    }
  }, [activeTab, currentPage, filters]);

  const loadLogs = async () => {
    try {
      setIsLoading(true);

      // TODO: Fix audit API integration
      // const cleanFilters = Object.fromEntries(
      //   Object.entries(filters).filter(([_, value]) => value !== '')
      // );

      // const response = await translationAdminApi.audit.getLogs(
      //   currentPage,
      //   20,
      //   cleanFilters
      // );

      // setLogs(response.data);
      // setTotalPages(response.total_pages);

      // Temporary mock data
      setLogs([]);
      setTotalPages(1);
    } catch (error) {
      console.error('Failed to load audit logs:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const loadStatistics = async () => {
    try {
      setIsLoading(true);
      // TODO: Fix audit API integration
      // const stats = await translationAdminApi.audit.getStatistics();
      // setStatistics(stats);

      // Temporary mock data
      setStatistics({
        total_actions: 0,
        actions_by_type: {},
        actions_by_user: [],
        actions_by_date: [],
        recent_actions: [],
      });
    } catch (error) {
      console.error('Failed to load audit statistics:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleFilterChange = (key: string, value: string) => {
    setFilters((prev) => ({ ...prev, [key]: value }));
    setCurrentPage(1); // Reset to first page when filters change
  };

  const clearFilters = () => {
    setFilters({
      user_id: '',
      action: '',
      entity_type: '',
      start_date: '',
      end_date: '',
    });
    setCurrentPage(1);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString('ru-RU');
  };

  const getActionColor = (action: string) => {
    const colors: Record<string, string> = {
      create: 'badge-success',
      update: 'badge-warning',
      delete: 'badge-error',
      translate: 'badge-info',
      approve: 'badge-primary',
      rollback: 'badge-secondary',
      import: 'badge-accent',
      export: 'badge-neutral',
    };
    return colors[action] || 'badge-ghost';
  };

  const getActionIcon = (action: string) => {
    switch (action) {
      case 'create':
        return '➕';
      case 'update':
        return '✏️';
      case 'delete':
        return '🗑️';
      case 'translate':
        return '🌐';
      case 'approve':
        return '✅';
      case 'rollback':
        return '↩️';
      case 'import':
        return '📥';
      case 'export':
        return '📤';
      default:
        return '📝';
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Журнал аудита</h2>
          <p className="text-base-content/60 mt-1">
            История всех операций с переводами
          </p>
        </div>
        <button
          className="btn btn-ghost btn-sm gap-2"
          onClick={activeTab === 'logs' ? loadLogs : loadStatistics}
          disabled={isLoading}
        >
          <ArrowPathIcon className="h-4 w-4" />
          Обновить
        </button>
      </div>

      {/* Tabs */}
      <div className="tabs tabs-boxed">
        <button
          className={`tab ${activeTab === 'logs' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('logs')}
        >
          <DocumentTextIcon className="h-4 w-4 mr-2" />
          Логи операций
        </button>
        <button
          className={`tab ${activeTab === 'stats' ? 'tab-active' : ''}`}
          onClick={() => setActiveTab('stats')}
        >
          <ChartBarIcon className="h-4 w-4 mr-2" />
          Статистика
        </button>
      </div>

      {/* Logs Tab */}
      {activeTab === 'logs' && (
        <>
          {/* Filters */}
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              <div className="flex items-center gap-2 mb-4">
                <FunnelIcon className="h-5 w-5" />
                <h3 className="font-semibold">Фильтры</h3>
                <button
                  className="btn btn-ghost btn-xs ml-auto"
                  onClick={clearFilters}
                >
                  Сбросить
                </button>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">ID пользователя</span>
                  </label>
                  <input
                    type="number"
                    placeholder="123"
                    className="input input-bordered input-sm"
                    value={filters.user_id}
                    onChange={(e) =>
                      handleFilterChange('user_id', e.target.value)
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Действие</span>
                  </label>
                  <select
                    className="select select-bordered select-sm"
                    value={filters.action}
                    onChange={(e) =>
                      handleFilterChange('action', e.target.value)
                    }
                  >
                    <option value="">Все действия</option>
                    {ACTIONS.map((action) => (
                      <option key={action} value={action}>
                        {getActionIcon(action)} {action}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Тип сущности</span>
                  </label>
                  <select
                    className="select select-bordered select-sm"
                    value={filters.entity_type}
                    onChange={(e) =>
                      handleFilterChange('entity_type', e.target.value)
                    }
                  >
                    <option value="">Все типы</option>
                    {ENTITY_TYPES.map((type) => (
                      <option key={type} value={type}>
                        {type}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Дата от</span>
                  </label>
                  <input
                    type="date"
                    className="input input-bordered input-sm"
                    value={filters.start_date}
                    onChange={(e) =>
                      handleFilterChange('start_date', e.target.value)
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Дата до</span>
                  </label>
                  <input
                    type="date"
                    className="input input-bordered input-sm"
                    value={filters.end_date}
                    onChange={(e) =>
                      handleFilterChange('end_date', e.target.value)
                    }
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Logs List */}
          <div className="card bg-base-100 shadow-sm">
            <div className="card-body">
              {isLoading ? (
                <div className="flex justify-center py-8">
                  <span className="loading loading-spinner loading-lg"></span>
                </div>
              ) : logs.length === 0 ? (
                <div className="text-center py-8 text-base-content/60">
                  Записи не найдены
                </div>
              ) : (
                <div className="space-y-3">
                  {logs.map((log) => (
                    <div
                      key={log.id}
                      className="border border-base-300 rounded-lg p-4"
                    >
                      <div className="flex items-start justify-between">
                        <div className="flex-1">
                          <div className="flex items-center gap-2 mb-2">
                            <div
                              className={`badge ${getActionColor(log.action)}`}
                            >
                              {getActionIcon(log.action)} {log.action}
                            </div>

                            {log.entity_type && (
                              <div className="badge badge-outline">
                                {log.entity_type}
                                {log.entity_id && ` #${log.entity_id}`}
                              </div>
                            )}
                          </div>

                          <div className="flex items-center gap-4 text-sm text-base-content/60 mb-3">
                            <div className="flex items-center gap-1">
                              <ClockIcon className="h-3 w-3" />
                              {formatDate(log.created_at)}
                            </div>

                            {log.user_id && (
                              <div className="flex items-center gap-1">
                                <UserIcon className="h-3 w-3" />
                                Пользователь {log.user_id}
                              </div>
                            )}

                            {log.ip_address && (
                              <div className="text-xs">
                                IP: {log.ip_address}
                              </div>
                            )}
                          </div>

                          {(log.old_value || log.new_value) && (
                            <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                              {log.old_value && (
                                <div className="bg-error/10 p-2 rounded">
                                  <div className="text-xs text-error font-medium mb-1">
                                    Было:
                                  </div>
                                  <div className="text-sm">{log.old_value}</div>
                                </div>
                              )}

                              {log.new_value && (
                                <div className="bg-success/10 p-2 rounded">
                                  <div className="text-xs text-success font-medium mb-1">
                                    Стало:
                                  </div>
                                  <div className="text-sm">{log.new_value}</div>
                                </div>
                              )}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              )}

              {/* Pagination */}
              {totalPages > 1 && (
                <div className="flex justify-center mt-6">
                  <div className="join">
                    <button
                      className="join-item btn"
                      onClick={() =>
                        setCurrentPage((prev) => Math.max(prev - 1, 1))
                      }
                      disabled={currentPage === 1}
                    >
                      «
                    </button>

                    <button className="join-item btn">
                      Страница {currentPage} из {totalPages}
                    </button>

                    <button
                      className="join-item btn"
                      onClick={() =>
                        setCurrentPage((prev) => Math.min(prev + 1, totalPages))
                      }
                      disabled={currentPage === totalPages}
                    >
                      »
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </>
      )}

      {/* Statistics Tab */}
      {activeTab === 'stats' && (
        <div className="space-y-6">
          {isLoading ? (
            <div className="flex justify-center py-8">
              <span className="loading loading-spinner loading-lg"></span>
            </div>
          ) : !statistics ? (
            <div className="text-center py-8 text-base-content/60">
              Статистика не найдена
            </div>
          ) : (
            <>
              {/* Overview */}
              <div className="card bg-base-100 shadow-sm">
                <div className="card-body">
                  <h3 className="card-title">Общая статистика</h3>
                  <div className="stat">
                    <div className="stat-title">Всего операций</div>
                    <div className="stat-value text-primary">
                      {statistics.total_actions.toLocaleString()}
                    </div>
                  </div>
                </div>
              </div>

              {/* Actions by Type */}
              <div className="card bg-base-100 shadow-sm">
                <div className="card-body">
                  <h3 className="card-title">По типу операций</h3>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    {Object.entries(statistics.actions_by_type).map(
                      ([action, count]) => (
                        <div key={action} className="stat">
                          <div className="stat-figure text-2xl">
                            {getActionIcon(action)}
                          </div>
                          <div className="stat-title">{action}</div>
                          <div className="stat-value text-sm">{count}</div>
                        </div>
                      )
                    )}
                  </div>
                </div>
              </div>

              {/* Top Users */}
              <div className="card bg-base-100 shadow-sm">
                <div className="card-body">
                  <h3 className="card-title">Активность пользователей</h3>
                  <div className="space-y-2">
                    {statistics.actions_by_user
                      .sort((a, b) => b.action_count - a.action_count)
                      .slice(0, 10)
                      .map((user) => (
                        <div
                          key={user.user_id}
                          className="flex justify-between items-center p-2 bg-base-200 rounded"
                        >
                          <span className="flex items-center gap-2">
                            <UserIcon className="h-4 w-4" />
                            {user.user_name}
                          </span>
                          <span className="font-medium">
                            {user.action_count} операций
                          </span>
                        </div>
                      ))}
                  </div>
                </div>
              </div>

              {/* Recent Actions */}
              {statistics.recent_actions.length > 0 && (
                <div className="card bg-base-100 shadow-sm">
                  <div className="card-body">
                    <h3 className="card-title">Последние операции</h3>
                    <div className="space-y-2">
                      {statistics.recent_actions.slice(0, 10).map((log) => (
                        <div
                          key={log.id}
                          className="flex justify-between items-center p-2 bg-base-200 rounded text-sm"
                        >
                          <div className="flex items-center gap-2">
                            <div
                              className={`badge badge-sm ${getActionColor(log.action)}`}
                            >
                              {getActionIcon(log.action)} {log.action}
                            </div>
                            {log.entity_type && (
                              <span>
                                {log.entity_type} #{log.entity_id}
                              </span>
                            )}
                          </div>
                          <span className="text-base-content/60">
                            {formatDate(log.created_at)}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )}
            </>
          )}
        </div>
      )}
    </div>
  );
}
