'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { tokenManager } from '@/utils/tokenManager';

interface OptimizationParams {
  field_names?: string[];
  item_type: 'marketplace' | 'storefront' | 'global';
  category_id?: number;
  min_sample_size: number;
  confidence_level: number;
  learning_rate: number;
  max_iterations: number;
  analysis_period_days: number;
  auto_apply: boolean;
}

interface OptimizationResult {
  field_name: string;
  current_weight: number;
  optimized_weight: number;
  improvement_score: number;
  confidence_level: number;
  sample_size: number;
  current_ctr: number;
  predicted_ctr: number;
  statistical_significance_level: number;
}

interface OptimizationSession {
  id: number;
  status: 'running' | 'completed' | 'failed' | 'cancelled';
  start_time: string;
  end_time?: string;
  total_fields: number;
  processed_fields: number;
  results?: OptimizationResult[];
  error_message?: string;
  created_by: number;
}

export default function WeightOptimization() {
  const t = useTranslations('admin');
  const [isOptimizing, setIsOptimizing] = useState(false);
  const [currentSession, setCurrentSession] =
    useState<OptimizationSession | null>(null);
  const [optimizationResults, setOptimizationResults] = useState<
    OptimizationResult[]
  >([]);
  const [selectedResults, setSelectedResults] = useState<number[]>([]);
  const [showAdvanced, setShowAdvanced] = useState(false);

  // Параметры оптимизации
  const [params, setParams] = useState<OptimizationParams>({
    item_type: 'global',
    min_sample_size: 100,
    confidence_level: 0.85,
    learning_rate: 0.01,
    max_iterations: 1000,
    analysis_period_days: 30,
    auto_apply: false,
  });

  // Состояния для аналитики
  const [isAnalyzing, setIsAnalyzing] = useState(false);
  const [analysisResults, setAnalysisResults] = useState<OptimizationResult[]>(
    []
  );

  // Polling для обновления статуса оптимизации
  useEffect(() => {
    let interval: NodeJS.Timeout;

    if (currentSession && currentSession.status === 'running') {
      interval = setInterval(async () => {
        try {
          const response = await fetch(
            `/api/v1/admin/search/optimization-status/${currentSession.id}`,
            {
              headers: {
                Authorization: `Bearer ${tokenManager.getAccessToken()}`,
              },
            }
          );

          if (response.ok) {
            const data = await response.json();
            setCurrentSession(data.data);

            if (data.data.status !== 'running') {
              setIsOptimizing(false);
              if (data.data.results) {
                setOptimizationResults(data.data.results);
              }
            }
          }
        } catch (error) {
          console.error('Failed to fetch optimization status:', error);
        }
      }, 2000);
    }

    return () => {
      if (interval) clearInterval(interval);
    };
  }, [currentSession]);

  const startOptimization = async () => {
    try {
      setIsOptimizing(true);
      const response = await fetch('/api/v1/admin/search/optimize-weights', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokenManager.getAccessToken()}`,
        },
        body: JSON.stringify(params),
      });

      if (!response.ok) {
        throw new Error('Failed to start optimization');
      }

      const data = await response.json();

      // Получаем статус сессии
      const sessionResponse = await fetch(
        `/api/v1/admin/search/optimization-status/${data.data.session_id}`,
        {
          headers: {
            Authorization: `Bearer ${tokenManager.getAccessToken()}`,
          },
        }
      );

      if (sessionResponse.ok) {
        const sessionData = await sessionResponse.json();
        setCurrentSession(sessionData.data);
      }
    } catch (error) {
      console.error('Failed to start optimization:', error);
      setIsOptimizing(false);
    }
  };

  const cancelOptimization = async () => {
    if (!currentSession) return;

    try {
      await fetch(
        `/api/v1/admin/search/optimization-cancel/${currentSession.id}`,
        {
          method: 'POST',
          headers: {
            Authorization: `Bearer ${tokenManager.getAccessToken()}`,
          },
        }
      );

      setIsOptimizing(false);
      setCurrentSession(null);
    } catch (error) {
      console.error('Failed to cancel optimization:', error);
    }
  };

  const applySelectedWeights = async () => {
    if (!currentSession || selectedResults.length === 0) return;

    try {
      const response = await fetch('/api/v1/admin/search/apply-weights', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokenManager.getAccessToken()}`,
        },
        body: JSON.stringify({
          session_id: currentSession.id,
          selected_results: selectedResults,
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to apply weights');
      }

      alert(t('messages.weightsApplied'));
      setSelectedResults([]);
    } catch (error) {
      console.error('Failed to apply weights:', error);
      alert(t('messages.weightsApplyError'));
    }
  };

  const analyzeCurrentWeights = async () => {
    try {
      setIsAnalyzing(true);
      const fromDate = new Date();
      fromDate.setDate(fromDate.getDate() - params.analysis_period_days);

      const response = await fetch('/api/v1/admin/search/analyze-weights', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${tokenManager.getAccessToken()}`,
        },
        body: JSON.stringify({
          item_type: params.item_type,
          category_id: params.category_id,
          from_date: fromDate.toISOString().split('T')[0],
          to_date: new Date().toISOString().split('T')[0],
        }),
      });

      if (!response.ok) {
        throw new Error('Failed to analyze weights');
      }

      const data = await response.json();
      setAnalysisResults(data.data.results || []);
    } catch (error) {
      console.error('Failed to analyze weights:', error);
    } finally {
      setIsAnalyzing(false);
    }
  };

  const handleResultSelection = (index: number, checked: boolean) => {
    if (checked) {
      setSelectedResults([...selectedResults, index]);
    } else {
      setSelectedResults(selectedResults.filter((i) => i !== index));
    }
  };

  const getImprovementColor = (score: number) => {
    if (score > 10) return 'text-success';
    if (score > 0) return 'text-warning';
    return 'text-error';
  };

  const getConfidenceColor = (level: number) => {
    if (level > 0.9) return 'text-success';
    if (level > 0.7) return 'text-warning';
    return 'text-error';
  };

  return (
    <div className="space-y-6">
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
                d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z"
              />
            </svg>
            {t('title')}
          </h3>

          <p className="text-base-content/60 mb-4">{t('description')}</p>

          {/* Параметры оптимизации */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('contentType')}</span>
              </label>
              <select
                className="select select-bordered"
                value={params.item_type}
                onChange={(e) =>
                  setParams({ ...params, item_type: e.target.value as any })
                }
                disabled={isOptimizing}
              >
                <option value="global">{t('contentTypes.global')}</option>
                <option value="marketplace">
                  {t('contentTypes.marketplace')}
                </option>
                <option value="storefront">
                  {t('contentTypes.storefront')}
                </option>
              </select>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('analysisPeriod')}</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={params.analysis_period_days}
                onChange={(e) =>
                  setParams({
                    ...params,
                    analysis_period_days: parseInt(e.target.value),
                  })
                }
                min={1}
                max={365}
                disabled={isOptimizing}
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('minSampleSize')}</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={params.min_sample_size}
                onChange={(e) =>
                  setParams({
                    ...params,
                    min_sample_size: parseInt(e.target.value),
                  })
                }
                min={10}
                max={10000}
                disabled={isOptimizing}
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">{t('confidenceLevel')}</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={params.confidence_level}
                onChange={(e) =>
                  setParams({
                    ...params,
                    confidence_level: parseFloat(e.target.value),
                  })
                }
                min={0.5}
                max={0.99}
                step={0.01}
                disabled={isOptimizing}
              />
            </div>
          </div>

          {/* Расширенные параметры */}
          <div className="collapse collapse-arrow bg-base-200">
            <input
              type="checkbox"
              checked={showAdvanced}
              onChange={(e) => setShowAdvanced(e.target.checked)}
            />
            <div className="collapse-title text-sm font-medium">
              {t('advancedParams')}
            </div>
            <div className="collapse-content">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4 pt-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">{t('learningRate')}</span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered input-sm"
                    value={params.learning_rate}
                    onChange={(e) =>
                      setParams({
                        ...params,
                        learning_rate: parseFloat(e.target.value),
                      })
                    }
                    min={0.001}
                    max={1}
                    step={0.001}
                    disabled={isOptimizing}
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">{t('maxIterations')}</span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered input-sm"
                    value={params.max_iterations}
                    onChange={(e) =>
                      setParams({
                        ...params,
                        max_iterations: parseInt(e.target.value),
                      })
                    }
                    min={10}
                    max={10000}
                    disabled={isOptimizing}
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Безопасность */}
          <div className="form-control">
            <label className="cursor-pointer label">
              <span className="label-text">
                <span className="text-warning">⚠️</span> {t('autoApply')}
                <span className="text-xs text-base-content/60 block">
                  {t('autoApplyWarning')}
                </span>
              </span>
              <input
                type="checkbox"
                className="checkbox checkbox-warning"
                checked={params.auto_apply}
                onChange={(e) =>
                  setParams({ ...params, auto_apply: e.target.checked })
                }
                disabled={isOptimizing}
              />
            </label>
          </div>

          {/* Кнопки действий */}
          <div className="flex flex-wrap gap-3 mt-6">
            <button
              className={`btn btn-primary ${isOptimizing ? 'loading' : ''}`}
              onClick={startOptimization}
              disabled={isOptimizing}
            >
              {isOptimizing ? t('optimizing') : t('startOptimization')}
            </button>

            <button
              className={`btn btn-outline ${isAnalyzing ? 'loading' : ''}`}
              onClick={analyzeCurrentWeights}
              disabled={isOptimizing || isAnalyzing}
            >
              {isAnalyzing ? t('analyzing') : t('quickAnalysis')}
            </button>

            {isOptimizing && (
              <button
                className="btn btn-error btn-outline"
                onClick={cancelOptimization}
              >
                {t('cancel')}
              </button>
            )}
          </div>
        </div>
      </div>

      {/* Прогресс оптимизации */}
      {currentSession && isOptimizing && (
        <div className="card bg-base-100 shadow-md">
          <div className="card-body">
            <h4 className="card-title text-lg">{t('progress.title')}</h4>

            <div className="flex items-center gap-4">
              <div className="flex-1">
                <div className="flex justify-between text-sm mb-1">
                  <span>
                    {t('progress.processed')}: {currentSession.processed_fields}{' '}
                    / {currentSession.total_fields}
                  </span>
                  <span>
                    {Math.round(
                      (currentSession.processed_fields /
                        currentSession.total_fields) *
                        100
                    )}
                    %
                  </span>
                </div>
                <progress
                  className="progress progress-primary w-full"
                  value={currentSession.processed_fields}
                  max={currentSession.total_fields}
                ></progress>
              </div>

              <div className="loading loading-spinner loading-md"></div>
            </div>

            <div className="text-sm text-base-content/60">
              <p>
                {t('progress.status')}: {currentSession.status}
              </p>
              <p>
                {t('progress.started')}:{' '}
                {new Date(currentSession.start_time).toLocaleString()}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Результаты оптимизации */}
      {(optimizationResults.length > 0 || analysisResults.length > 0) && (
        <div className="card bg-base-100 shadow-md">
          <div className="card-body">
            <div className="flex justify-between items-center">
              <h4 className="card-title text-lg">
                {optimizationResults.length > 0
                  ? t('results.optimizationTitle')
                  : t('results.analysisTitle')}
              </h4>

              {optimizationResults.length > 0 && selectedResults.length > 0 && (
                <button
                  className="btn btn-success btn-sm"
                  onClick={applySelectedWeights}
                >
                  {t('results.applySelected')} ({selectedResults.length})
                </button>
              )}
            </div>

            <div className="overflow-x-auto">
              <table className="table table-zebra">
                <thead>
                  <tr>
                    {optimizationResults.length > 0 && (
                      <th>
                        <input
                          type="checkbox"
                          className="checkbox checkbox-sm"
                          onChange={(e) => {
                            if (e.target.checked) {
                              setSelectedResults(
                                optimizationResults.map((_, i) => i)
                              );
                            } else {
                              setSelectedResults([]);
                            }
                          }}
                        />
                      </th>
                    )}
                    <th>{t('results.field')}</th>
                    <th>{t('results.currentWeight')}</th>
                    <th>{t('results.optimizedWeight')}</th>
                    <th>{t('results.improvement')}</th>
                    <th>{t('results.confidence')}</th>
                    <th>{t('results.sample')}</th>
                    <th>{t('results.currentCTR')}</th>
                    <th>{t('results.predictedCTR')}</th>
                  </tr>
                </thead>
                <tbody>
                  {(optimizationResults.length > 0
                    ? optimizationResults
                    : analysisResults
                  ).map((result, index) => (
                    <tr
                      key={index}
                      className={
                        selectedResults.includes(index) ? 'bg-primary/10' : ''
                      }
                    >
                      {optimizationResults.length > 0 && (
                        <td>
                          <input
                            type="checkbox"
                            className="checkbox checkbox-sm"
                            checked={selectedResults.includes(index)}
                            onChange={(e) =>
                              handleResultSelection(index, e.target.checked)
                            }
                          />
                        </td>
                      )}
                      <td className="font-medium">{result.field_name}</td>
                      <td>{result.current_weight.toFixed(3)}</td>
                      <td className="font-bold">
                        {result.optimized_weight.toFixed(3)}
                      </td>
                      <td
                        className={getImprovementColor(
                          result.improvement_score
                        )}
                      >
                        {result.improvement_score > 0 ? '+' : ''}
                        {result.improvement_score.toFixed(1)}%
                      </td>
                      <td
                        className={getConfidenceColor(result.confidence_level)}
                      >
                        {(result.confidence_level * 100).toFixed(1)}%
                      </td>
                      <td>{result.sample_size.toLocaleString()}</td>
                      <td>{(result.current_ctr * 100).toFixed(2)}%</td>
                      <td>{(result.predicted_ctr * 100).toFixed(2)}%</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>

            {optimizationResults.length > 0 && (
              <div className="alert alert-warning mt-4">
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
                    d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.98-.833-2.75 0L3.982 16.5c-.77.833.192 2.5 1.732 2.5z"
                  />
                </svg>
                <div>
                  <p className="font-bold">{t('safety.title')}</p>
                  <ul className="text-sm mt-1">
                    <li>• {t('safety.checkResults')}</li>
                    <li>• {t('safety.testSmall')}</li>
                    <li>• {t('safety.backup')}</li>
                    <li>• {t('safety.monitor')}</li>
                  </ul>
                </div>
              </div>
            )}
          </div>
        </div>
      )}
    </div>
  );
}
