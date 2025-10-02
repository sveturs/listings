'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { toast } from '@/utils/toast';
import { apiClient } from '@/services/api-client';

interface TranslationConflict {
  id: number;
  key: string;
  module: string;
  language: string;
  frontend_value: string;
  database_value: string;
  last_modified_frontend: string;
  last_modified_database: string;
  conflict_type:
    | 'value_mismatch'
    | 'missing_in_frontend'
    | 'missing_in_database';
  resolved: boolean;
  resolution?: 'use_frontend' | 'use_database' | 'use_custom';
  custom_value?: string;
  resolved_at?: string;
  resolved_by?: number;
}

interface ConflictResolverProps {
  onConflictResolved?: () => void;
}

export default function ConflictResolver({
  onConflictResolved,
}: ConflictResolverProps) {
  const _t = useTranslations('admin');
  const [conflicts, setConflicts] = useState<TranslationConflict[]>([]);
  const [loading, setLoading] = useState(true);
  const [selectedConflicts, setSelectedConflicts] = useState<Set<number>>(
    new Set()
  );
  const [resolutions, setResolutions] = useState<
    Record<
      number,
      {
        resolution: 'use_frontend' | 'use_database' | 'use_custom';
        custom_value?: string;
      }
    >
  >({});
  const [filter, setFilter] = useState<'all' | 'unresolved' | 'resolved'>(
    'unresolved'
  );
  const [searchTerm, setSearchTerm] = useState('');
  const [isResolving, setIsResolving] = useState(false);

  useEffect(() => {
    fetchConflicts();
  }, []);

  const fetchConflicts = async () => {
    setLoading(true);
    try {
      const response = await apiClient.get('/admin/translations/sync/conflicts');

      if (response.data) {
        setConflicts(response.data.data || []);
      } else {
        toast.error('Ошибка загрузки конфликтов');
      }
    } catch (error) {
      console.error('Error fetching conflicts:', error);
      toast.error('Ошибка при загрузке конфликтов');
    } finally {
      setLoading(false);
    }
  };

  const handleSelectAll = () => {
    const filteredConflicts = getFilteredConflicts();
    if (selectedConflicts.size === filteredConflicts.length) {
      setSelectedConflicts(new Set());
    } else {
      setSelectedConflicts(new Set(filteredConflicts.map((c) => c.id)));
    }
  };

  const handleSelectConflict = (id: number) => {
    const newSelected = new Set(selectedConflicts);
    if (newSelected.has(id)) {
      newSelected.delete(id);
    } else {
      newSelected.add(id);
    }
    setSelectedConflicts(newSelected);
  };

  const handleResolutionChange = (
    conflictId: number,
    resolution: 'use_frontend' | 'use_database' | 'use_custom',
    customValue?: string
  ) => {
    setResolutions((prev) => ({
      ...prev,
      [conflictId]: {
        resolution,
        custom_value: customValue,
      },
    }));
  };

  const handleBatchResolve = async (
    resolution: 'use_frontend' | 'use_database'
  ) => {
    if (selectedConflicts.size === 0) {
      toast.error('Выберите конфликты для разрешения');
      return;
    }

    setIsResolving(true);
    try {
      const conflictResolutions = Array.from(selectedConflicts).map((id) => ({
        conflict_id: id,
        resolution,
      }));

      const response = await apiClient.post(
        '/admin/translations/sync/conflicts/resolve',
        { resolutions: conflictResolutions }
      );

      if (response.data) {
        toast.success(`Разрешено ${selectedConflicts.size} конфликтов`);
        setSelectedConflicts(new Set());
        await fetchConflicts();
        onConflictResolved?.();
      } else {
        toast.error('Ошибка при разрешении конфликтов');
      }
    } catch (error) {
      console.error('Error resolving conflicts:', error);
      toast.error('Ошибка при разрешении конфликтов');
    } finally {
      setIsResolving(false);
    }
  };

  const handleIndividualResolve = async (conflictId: number) => {
    const resolution = resolutions[conflictId];
    if (!resolution) {
      toast.error('Выберите способ разрешения');
      return;
    }

    if (resolution.resolution === 'use_custom' && !resolution.custom_value) {
      toast.error('Введите кастомное значение');
      return;
    }

    setIsResolving(true);
    try {
      const response = await apiClient.post(
        '/admin/translations/sync/conflicts/resolve',
        {
          resolutions: [
            {
              conflict_id: conflictId,
              resolution: resolution.resolution,
              custom_value: resolution.custom_value,
            },
          ],
        }
      );

      if (response.data) {
        toast.success('Конфликт разрешен');
        await fetchConflicts();
        onConflictResolved?.();
      } else {
        toast.error('Ошибка при разрешении конфликта');
      }
    } catch (error) {
      console.error('Error resolving conflict:', error);
      toast.error('Ошибка при разрешении конфликта');
    } finally {
      setIsResolving(false);
    }
  };

  const getFilteredConflicts = () => {
    return conflicts.filter((conflict) => {
      if (filter === 'resolved' && !conflict.resolved) return false;
      if (filter === 'unresolved' && conflict.resolved) return false;

      if (searchTerm) {
        const search = searchTerm.toLowerCase();
        return (
          conflict.key.toLowerCase().includes(search) ||
          conflict.module.toLowerCase().includes(search) ||
          conflict.frontend_value?.toLowerCase().includes(search) ||
          conflict.database_value?.toLowerCase().includes(search)
        );
      }

      return true;
    });
  };

  const getConflictTypeLabel = (type: string) => {
    switch (type) {
      case 'value_mismatch':
        return { label: 'Несоответствие', color: 'badge-warning' };
      case 'missing_in_frontend':
        return { label: 'Нет в Frontend', color: 'badge-error' };
      case 'missing_in_database':
        return { label: 'Нет в БД', color: 'badge-info' };
      default:
        return { label: type, color: 'badge-ghost' };
    }
  };

  const formatDate = (dateString: string) => {
    if (!dateString) return '-';
    return new Date(dateString).toLocaleString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  const filteredConflicts = getFilteredConflicts();
  const unresolvedCount = conflicts.filter((c) => !c.resolved).length;

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex justify-between items-center flex-wrap gap-4">
        <div>
          <h3 className="text-xl font-semibold">Конфликты синхронизации</h3>
          <p className="text-sm text-base-content/70 mt-1">
            Всего: {conflicts.length} | Не разрешено: {unresolvedCount}
          </p>
        </div>

        <div className="flex gap-2">
          <button
            onClick={() => fetchConflicts()}
            className="btn btn-ghost btn-sm"
          >
            🔄 Обновить
          </button>
        </div>
      </div>

      {/* Filters */}
      <div className="flex flex-wrap gap-4">
        <div className="join">
          <button
            onClick={() => setFilter('all')}
            className={`join-item btn btn-sm ${filter === 'all' ? 'btn-active' : ''}`}
          >
            Все ({conflicts.length})
          </button>
          <button
            onClick={() => setFilter('unresolved')}
            className={`join-item btn btn-sm ${filter === 'unresolved' ? 'btn-active' : ''}`}
          >
            Не разрешенные ({unresolvedCount})
          </button>
          <button
            onClick={() => setFilter('resolved')}
            className={`join-item btn btn-sm ${filter === 'resolved' ? 'btn-active' : ''}`}
          >
            Разрешенные ({conflicts.length - unresolvedCount})
          </button>
        </div>

        <input
          type="text"
          placeholder="Поиск по ключу или значению..."
          value={searchTerm}
          onChange={(e) => setSearchTerm(e.target.value)}
          className="input input-bordered input-sm w-64"
        />
      </div>

      {/* Batch Actions */}
      {selectedConflicts.size > 0 && (
        <div className="alert alert-info">
          <div className="flex justify-between items-center w-full">
            <span>Выбрано конфликтов: {selectedConflicts.size}</span>
            <div className="flex gap-2">
              <button
                onClick={() => handleBatchResolve('use_frontend')}
                className="btn btn-sm btn-primary"
                disabled={isResolving}
              >
                Использовать Frontend
              </button>
              <button
                onClick={() => handleBatchResolve('use_database')}
                className="btn btn-sm btn-secondary"
                disabled={isResolving}
              >
                Использовать БД
              </button>
              <button
                onClick={() => setSelectedConflicts(new Set())}
                className="btn btn-sm btn-ghost"
              >
                Отменить выбор
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Conflicts List */}
      <div className="space-y-4">
        {filteredConflicts.length === 0 ? (
          <div className="text-center py-8 text-base-content/50">
            {searchTerm
              ? 'Конфликты не найдены'
              : 'Нет конфликтов для отображения'}
          </div>
        ) : (
          <>
            {/* Select All */}
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-3">
                <input
                  type="checkbox"
                  checked={
                    selectedConflicts.size === filteredConflicts.length &&
                    filteredConflicts.length > 0
                  }
                  onChange={handleSelectAll}
                  className="checkbox checkbox-primary"
                />
                <span className="label-text">Выбрать все</span>
              </label>
            </div>

            {/* Conflict Cards */}
            {filteredConflicts.map((conflict) => {
              const typeInfo = getConflictTypeLabel(conflict.conflict_type);
              const resolution = resolutions[conflict.id];

              return (
                <div
                  key={conflict.id}
                  className={`card bg-base-100 shadow-sm border ${
                    conflict.resolved
                      ? 'border-success/30 bg-success/5'
                      : 'border-base-300'
                  }`}
                >
                  <div className="card-body">
                    {/* Header */}
                    <div className="flex items-start justify-between">
                      <div className="flex items-start gap-3">
                        {!conflict.resolved && (
                          <input
                            type="checkbox"
                            checked={selectedConflicts.has(conflict.id)}
                            onChange={() => handleSelectConflict(conflict.id)}
                            className="checkbox checkbox-primary mt-1"
                          />
                        )}
                        <div>
                          <div className="flex items-center gap-2 flex-wrap">
                            <h4 className="font-mono text-sm font-semibold">
                              {conflict.key}
                            </h4>
                            <div className="badge badge-ghost badge-sm">
                              {conflict.module}
                            </div>
                            <div className="badge badge-primary badge-sm">
                              {conflict.language.toUpperCase()}
                            </div>
                            <div className={`badge ${typeInfo.color} badge-sm`}>
                              {typeInfo.label}
                            </div>
                            {conflict.resolved && (
                              <div className="badge badge-success badge-sm">
                                ✓ Разрешен
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Values Comparison */}
                    <div className="grid md:grid-cols-2 gap-4 mt-4">
                      <div className="space-y-2">
                        <div className="text-xs font-semibold text-primary">
                          Frontend значение:
                        </div>
                        <div className="p-3 bg-base-200 rounded-lg">
                          <div className="text-sm break-all">
                            {conflict.frontend_value || (
                              <span className="text-base-content/50">
                                Отсутствует
                              </span>
                            )}
                          </div>
                          <div className="text-xs text-base-content/50 mt-1">
                            {formatDate(conflict.last_modified_frontend)}
                          </div>
                        </div>
                      </div>

                      <div className="space-y-2">
                        <div className="text-xs font-semibold text-secondary">
                          База данных значение:
                        </div>
                        <div className="p-3 bg-base-200 rounded-lg">
                          <div className="text-sm break-all">
                            {conflict.database_value || (
                              <span className="text-base-content/50">
                                Отсутствует
                              </span>
                            )}
                          </div>
                          <div className="text-xs text-base-content/50 mt-1">
                            {formatDate(conflict.last_modified_database)}
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Resolution Controls */}
                    {!conflict.resolved && (
                      <div className="mt-4 space-y-3">
                        <div className="divider text-xs">Способ разрешения</div>

                        <div className="flex flex-wrap gap-2">
                          <label className="label cursor-pointer gap-2">
                            <input
                              type="radio"
                              name={`resolution-${conflict.id}`}
                              className="radio radio-primary radio-sm"
                              checked={
                                resolution?.resolution === 'use_frontend'
                              }
                              onChange={() =>
                                handleResolutionChange(
                                  conflict.id,
                                  'use_frontend'
                                )
                              }
                            />
                            <span className="label-text text-sm">Frontend</span>
                          </label>

                          <label className="label cursor-pointer gap-2">
                            <input
                              type="radio"
                              name={`resolution-${conflict.id}`}
                              className="radio radio-secondary radio-sm"
                              checked={
                                resolution?.resolution === 'use_database'
                              }
                              onChange={() =>
                                handleResolutionChange(
                                  conflict.id,
                                  'use_database'
                                )
                              }
                            />
                            <span className="label-text text-sm">
                              База данных
                            </span>
                          </label>

                          <label className="label cursor-pointer gap-2">
                            <input
                              type="radio"
                              name={`resolution-${conflict.id}`}
                              className="radio radio-accent radio-sm"
                              checked={resolution?.resolution === 'use_custom'}
                              onChange={() =>
                                handleResolutionChange(
                                  conflict.id,
                                  'use_custom'
                                )
                              }
                            />
                            <span className="label-text text-sm">
                              Кастомное
                            </span>
                          </label>
                        </div>

                        {resolution?.resolution === 'use_custom' && (
                          <textarea
                            value={resolution.custom_value || ''}
                            onChange={(e) =>
                              handleResolutionChange(
                                conflict.id,
                                'use_custom',
                                e.target.value
                              )
                            }
                            className="textarea textarea-bordered w-full"
                            rows={2}
                            placeholder="Введите кастомное значение..."
                          />
                        )}

                        <button
                          onClick={() => handleIndividualResolve(conflict.id)}
                          className="btn btn-primary btn-sm"
                          disabled={!resolution || isResolving}
                        >
                          {isResolving ? (
                            <>
                              <span className="loading loading-spinner loading-xs"></span>
                              Разрешение...
                            </>
                          ) : (
                            'Разрешить конфликт'
                          )}
                        </button>
                      </div>
                    )}

                    {/* Resolution Info */}
                    {conflict.resolved && conflict.resolution && (
                      <div className="mt-4 p-3 bg-success/10 rounded-lg">
                        <div className="text-xs text-success font-semibold mb-1">
                          Разрешение:
                        </div>
                        <div className="text-sm">
                          {conflict.resolution === 'use_frontend' &&
                            'Использовано значение Frontend'}
                          {conflict.resolution === 'use_database' &&
                            'Использовано значение из БД'}
                          {conflict.resolution === 'use_custom' && (
                            <div>
                              <div>Кастомное значение:</div>
                              <div className="mt-1 p-2 bg-base-200 rounded">
                                {conflict.custom_value}
                              </div>
                            </div>
                          )}
                        </div>
                        {conflict.resolved_at && (
                          <div className="text-xs text-base-content/50 mt-2">
                            Разрешено: {formatDate(conflict.resolved_at)}
                          </div>
                        )}
                      </div>
                    )}
                  </div>
                </div>
              );
            })}
          </>
        )}
      </div>
    </div>
  );
}
