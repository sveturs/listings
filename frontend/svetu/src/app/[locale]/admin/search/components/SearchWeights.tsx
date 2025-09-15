'use client';

import { useState, useEffect, useCallback } from 'react';
import WeightOptimization from './WeightOptimization';
import { tokenManager } from '@/utils/tokenManager';
import configManager from '@/config';

interface SearchWeight {
  id: number;
  field_name: string;
  weight: number;
  search_type: string;
  item_type: string;
  category_id?: number;
  description?: string;
  is_active: boolean;
  version: number;
  created_at: string;
  updated_at: string;
}

export default function SearchWeights() {
  const [weights, setWeights] = useState<SearchWeight[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'manual' | 'optimization'>(
    'manual'
  );
  const [selectedItemType, setSelectedItemType] = useState<
    'global' | 'marketplace' | 'storefront'
  >('global');
  const [editingWeight, setEditingWeight] = useState<SearchWeight | null>(null);
  const [newWeight, setNewWeight] = useState<number>(0);

  const loadWeights = useCallback(async () => {
    try {
      setLoading(true);
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/admin/search/weights?item_type=${selectedItemType}`,
        {
          headers: {
            Authorization: `Bearer ${tokenManager.getAccessToken()}`,
          },
        }
      );

      if (!response.ok) {
        throw new Error('Failed to load weights');
      }

      const data = await response.json();
      setWeights(data.data || []);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to load weights');
    } finally {
      setLoading(false);
    }
  }, [selectedItemType]);

  useEffect(() => {
    loadWeights();
  }, [loadWeights]);

  const updateWeight = async (weightId: number, newWeight: number) => {
    try {
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/admin/search/weights/${weightId}`,
        {
          method: 'PUT',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${tokenManager.getAccessToken()}`,
          },
          body: JSON.stringify({ weight: newWeight }),
        }
      );

      if (!response.ok) {
        throw new Error('Failed to update weight');
      }

      await loadWeights();
      setEditingWeight(null);
    } catch (err) {
      alert(
        'Ошибка при обновлении веса: ' +
          (err instanceof Error ? err.message : 'Unknown error')
      );
    }
  };

  const createBackup = async () => {
    try {
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/admin/search/backup-weights`,
        {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            Authorization: `Bearer ${tokenManager.getAccessToken()}`,
          },
          body: JSON.stringify({
            item_type: selectedItemType,
          }),
        }
      );

      if (!response.ok) {
        throw new Error('Failed to create backup');
      }

      alert('Резервная копия весов создана успешно!');
    } catch (err) {
      alert(
        'Ошибка при создании резервной копии: ' +
          (err instanceof Error ? err.message : 'Unknown error')
      );
    }
  };

  const getWeightColor = (weight: number) => {
    if (weight >= 0.8) return 'text-success';
    if (weight >= 0.5) return 'text-warning';
    return 'text-error';
  };

  const getFieldTypeIcon = (searchType: string) => {
    switch (searchType) {
      case 'fulltext':
        return '📝';
      case 'fuzzy':
        return '🔍';
      case 'exact':
        return '🎯';
      default:
        return '❓';
    }
  };

  return (
    <div className="space-y-6">
      {/* Заголовок и навигация */}
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <h3 className="card-title flex items-center gap-2">
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4"
              />
            </svg>
            Управление весами поиска
          </h3>
          <p className="text-base-content/60">
            Настройка весов полей поиска для управления релевантностью
            результатов. Используйте ручную настройку или автоматическую
            оптимизацию на основе ML.
          </p>

          {/* Табы */}
          <div className="tabs tabs-boxed bg-base-200 w-fit mt-4">
            <button
              className={`tab ${activeTab === 'manual' ? 'tab-active' : ''}`}
              onClick={() => setActiveTab('manual')}
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
                  d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4"
                />
              </svg>
              Ручная настройка
            </button>
            <button
              className={`tab ${activeTab === 'optimization' ? 'tab-active' : ''}`}
              onClick={() => setActiveTab('optimization')}
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
                  d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
                />
              </svg>
              Автоматическая оптимизация
            </button>
          </div>
        </div>
      </div>

      {/* Ручная настройка */}
      {activeTab === 'manual' && (
        <>
          {/* Фильтры и действия */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <div className="flex flex-wrap justify-between items-center gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Тип контента</span>
                  </label>
                  <select
                    className="select select-bordered"
                    value={selectedItemType}
                    onChange={(e) => setSelectedItemType(e.target.value as any)}
                  >
                    <option value="global">Глобальные веса</option>
                    <option value="marketplace">Маркетплейс</option>
                    <option value="storefront">Магазины</option>
                  </select>
                </div>

                <div className="flex gap-2">
                  <button
                    className="btn btn-outline btn-sm"
                    onClick={createBackup}
                  >
                    <svg
                      className="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M8 7H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-3m-1 4l-3-3m0 0l-3 3m3-3v12"
                      />
                    </svg>
                    Создать бэкап
                  </button>

                  <button
                    className="btn btn-primary btn-sm"
                    onClick={loadWeights}
                  >
                    <svg
                      className="w-4 h-4"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
                      />
                    </svg>
                    Обновить
                  </button>
                </div>
              </div>
            </div>
          </div>

          {/* Таблица весов */}
          <div className="card bg-base-100 shadow-md">
            <div className="card-body">
              <h4 className="card-title text-lg">Веса полей поиска</h4>

              {loading ? (
                <div className="flex justify-center py-8">
                  <div className="loading loading-spinner loading-lg"></div>
                </div>
              ) : error ? (
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
              ) : (
                <div className="overflow-x-auto">
                  <table className="table table-zebra">
                    <thead>
                      <tr>
                        <th>Поле</th>
                        <th>Тип поиска</th>
                        <th>Текущий вес</th>
                        <th>Описание</th>
                        <th>Статус</th>
                        <th>Действия</th>
                      </tr>
                    </thead>
                    <tbody>
                      {weights.map((weight) => (
                        <tr key={weight.id}>
                          <td className="font-medium">{weight.field_name}</td>
                          <td>
                            <div className="flex items-center gap-2">
                              <span>
                                {getFieldTypeIcon(weight.search_type)}
                              </span>
                              <span className="badge badge-sm badge-outline">
                                {weight.search_type}
                              </span>
                            </div>
                          </td>
                          <td>
                            {editingWeight?.id === weight.id ? (
                              <div className="flex items-center gap-2">
                                <input
                                  type="number"
                                  className="input input-xs input-bordered w-20"
                                  value={newWeight}
                                  onChange={(e) =>
                                    setNewWeight(parseFloat(e.target.value))
                                  }
                                  min={0}
                                  max={1}
                                  step={0.01}
                                />
                                <button
                                  className="btn btn-xs btn-success"
                                  onClick={() =>
                                    updateWeight(weight.id, newWeight)
                                  }
                                >
                                  ✓
                                </button>
                                <button
                                  className="btn btn-xs btn-error"
                                  onClick={() => setEditingWeight(null)}
                                >
                                  ✕
                                </button>
                              </div>
                            ) : (
                              <span
                                className={`font-bold ${getWeightColor(weight.weight)}`}
                              >
                                {weight.weight.toFixed(3)}
                              </span>
                            )}
                          </td>
                          <td className="text-sm text-base-content/60 max-w-xs truncate">
                            {weight.description || 'Нет описания'}
                          </td>
                          <td>
                            <div
                              className={`badge badge-sm ${weight.is_active ? 'badge-success' : 'badge-error'}`}
                            >
                              {weight.is_active ? 'Активен' : 'Неактивен'}
                            </div>
                          </td>
                          <td>
                            <button
                              className="btn btn-xs btn-ghost"
                              onClick={() => {
                                setEditingWeight(weight);
                                setNewWeight(weight.weight);
                              }}
                              disabled={editingWeight !== null}
                            >
                              <svg
                                className="w-3 h-3"
                                fill="none"
                                stroke="currentColor"
                                viewBox="0 0 24 24"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  strokeWidth={2}
                                  d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"
                                />
                              </svg>
                            </button>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>

                  {weights.length === 0 && (
                    <div className="text-center py-8 text-base-content/60">
                      Веса для выбранного типа контента не найдены
                    </div>
                  )}
                </div>
              )}

              <div className="alert alert-info mt-4">
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
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
                <div>
                  <p className="font-bold">Как работают веса:</p>
                  <ul className="text-sm mt-1">
                    <li>
                      • <strong>0.8-1.0</strong> - Критически важные поля
                      (название, основные характеристики)
                    </li>
                    <li>
                      • <strong>0.5-0.8</strong> - Важные поля (описание,
                      категория)
                    </li>
                    <li>
                      • <strong>0.0-0.5</strong> - Дополнительные поля (теги,
                      метаданные)
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </>
      )}

      {/* Автоматическая оптимизация */}
      {activeTab === 'optimization' && <WeightOptimization />}
    </div>
  );
}
