'use client';

import { useTranslations } from 'next-intl';
import { useState } from 'react';
import { apiClient } from '@/services/api-client';

interface Test {
  id: string;
  name: string;
  description: string;
  category: 'quality' | 'unit' | 'integration' | 'build' | 'coverage' | 'functional';
  icon: string;
}

interface TestResult {
  name: string;
  status: 'success' | 'error' | 'warning' | 'running' | 'pending';
  duration?: number;
  output?: string;
  error?: string;
  stats?: {
    passed: number;
    failed: number;
    skipped: number;
    total: number;
  };
}

const TESTS: Test[] = [
  // Code Quality
  {
    id: 'backend-format',
    name: 'Backend Format',
    description: 'Check Go code formatting (gofumpt, goimports)',
    category: 'quality',
    icon: '🎨',
  },
  {
    id: 'backend-lint',
    name: 'Backend Lint',
    description: 'Run golangci-lint for code quality',
    category: 'quality',
    icon: '🔍',
  },
  {
    id: 'frontend-format',
    name: 'Frontend Format',
    description: 'Check Prettier formatting',
    category: 'quality',
    icon: '✨',
  },
  {
    id: 'frontend-lint',
    name: 'Frontend Lint',
    description: 'Run ESLint for code quality',
    category: 'quality',
    icon: '🔎',
  },

  // Unit Tests
  {
    id: 'backend-tests-unit',
    name: 'Backend Unit Tests',
    description: 'Run Go unit tests only',
    category: 'unit',
    icon: '🧪',
  },
  {
    id: 'frontend-tests',
    name: 'Frontend Unit Tests',
    description: 'Run Jest unit tests',
    category: 'unit',
    icon: '🔬',
  },

  // Integration Tests
  {
    id: 'backend-tests-short',
    name: 'Backend Short Tests',
    description: 'Quick tests (skip Integration/OpenSearch)',
    category: 'integration',
    icon: '⚡',
  },
  {
    id: 'backend-tests-full',
    name: 'Backend Full Tests',
    description: 'All tests including integration',
    category: 'integration',
    icon: '🔥',
  },
  {
    id: 'backend-tests-postexpress',
    name: 'Post Express Integration',
    description: 'Test Post Express API integration',
    category: 'integration',
    icon: '📮',
  },
  {
    id: 'backend-tests-cache',
    name: 'Redis Cache Integration',
    description: 'Test Redis caching functionality',
    category: 'integration',
    icon: '💾',
  },
  {
    id: 'backend-tests-opensearch',
    name: 'OpenSearch Integration',
    description: 'Test search indexing and queries',
    category: 'integration',
    icon: '🔍',
  },

  // Build & Type Checking
  {
    id: 'backend-build',
    name: 'Backend Build',
    description: 'Compile Go code',
    category: 'build',
    icon: '🔨',
  },
  {
    id: 'frontend-build',
    name: 'Frontend Build',
    description: 'Next.js production build',
    category: 'build',
    icon: '⚙️',
  },
  {
    id: 'typescript-check',
    name: 'TypeScript Check',
    description: 'Type checking (tsc --noEmit)',
    category: 'build',
    icon: '📘',
  },

  // Coverage
  {
    id: 'backend-tests-coverage',
    name: 'Backend Test Coverage',
    description: 'Run tests with coverage report',
    category: 'coverage',
    icon: '📊',
  },
  {
    id: 'frontend-tests-coverage',
    name: 'Frontend Test Coverage',
    description: 'Jest tests with coverage',
    category: 'coverage',
    icon: '📈',
  },

  // Functional API Tests (Backend API Testing)
  {
    id: 'api-auth-flow',
    name: 'Auth Flow Test',
    description: 'Test authentication endpoints (login, me, logout)',
    category: 'functional',
    icon: '🔐',
  },
  {
    id: 'api-marketplace-crud',
    name: 'Marketplace CRUD',
    description: 'Test marketplace listing operations',
    category: 'functional',
    icon: '🛒',
  },
  {
    id: 'api-categories-fetch',
    name: 'Categories API',
    description: 'Test admin categories endpoints',
    category: 'functional',
    icon: '📁',
  },
  {
    id: 'api-search-functionality',
    name: 'Search API',
    description: 'Test unified search functionality',
    category: 'functional',
    icon: '🔍',
  },
  {
    id: 'api-admin-operations',
    name: 'Admin Operations',
    description: 'Test admin panel endpoints',
    category: 'functional',
    icon: '⚙️',
  },
];

export default function QualityTestsClient({ locale }: { locale: string }) {
  const t = useTranslations('admin.qualityTests');
  const [results, setResults] = useState<Record<string, TestResult>>({});
  const [running, setRunning] = useState<Set<string>>(new Set());
  const [expanded, setExpanded] = useState<Set<string>>(new Set());

  const runTest = async (testId: string) => {
    setRunning((prev) => new Set(prev).add(testId));
    setResults((prev) => ({
      ...prev,
      [testId]: { name: testId, status: 'running' },
    }));

    try {
      // Определяем категорию теста
      const test = TESTS.find((t) => t.id === testId);
      const isFunctional = test?.category === 'functional';

      if (isFunctional) {
        // Functional тесты: вызываем backend API через apiClient
        const response = await apiClient.post('/admin/tests/run', {
          test_suite: 'api-endpoints',
          parallel: false,
        });

        // Backend API возвращает test run ID - нужно подождать результаты
        if (response.data) {
          const runId = response.data.test_run_id;

          // Polling для получения результатов
          let attempts = 0;
          const maxAttempts = 30; // 30 секунд

          while (attempts < maxAttempts) {
            await new Promise((resolve) => setTimeout(resolve, 1000));

            const detailResponse = await apiClient.get(
              `/admin/tests/runs/${runId}`
            );
            const detail = detailResponse.data;

            if (
              detail.status === 'completed' ||
              detail.status === 'failed'
            ) {
              // Тесты из backend выполняются как единый suite
              // Показываем общий результат с статистикой
              setResults((prev) => ({
                ...prev,
                [testId]: {
                  name: testId,
                  status: detail.failed_tests > 0 ? 'error' : 'success',
                  duration: detail.duration_ms,
                  output: `Test suite completed: ${detail.passed_tests} passed, ${detail.failed_tests} failed`,
                  error: detail.failed_tests > 0 ? `${detail.failed_tests} tests failed` : undefined,
                  stats: {
                    passed: detail.passed_tests,
                    failed: detail.failed_tests,
                    skipped: detail.skipped_tests,
                    total: detail.total_tests,
                  },
                },
              }));
              break;
            }

            attempts++;
          }

          if (attempts >= maxAttempts) {
            throw new Error('Test execution timeout');
          }
        }
      } else {
        // Старые тесты: используем Next.js API route
        const response = await fetch('/api/admin/tests', {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
          },
          body: JSON.stringify({ test: testId }),
        });

        const data = await response.json();

        if (data.result) {
          setResults((prev) => ({
            ...prev,
            [testId]: data.result,
          }));
        } else if (data.error) {
          setResults((prev) => ({
            ...prev,
            [testId]: {
              name: testId,
              status: 'error',
              error: data.error,
            },
          }));
        }
      }
    } catch (error) {
      setResults((prev) => ({
        ...prev,
        [testId]: {
          name: testId,
          status: 'error',
          error: error instanceof Error ? error.message : 'Unknown error',
        },
      }));
    } finally {
      setRunning((prev) => {
        const newSet = new Set(prev);
        newSet.delete(testId);
        return newSet;
      });
    }
  };

  const runTestGroup = async (groupId: string) => {
    setRunning((prev) => new Set(prev).add(groupId));

    try {
      const response = await fetch('/api/admin/tests', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ test: groupId }),
      });

      const data = await response.json();

      if (data.results) {
        const newResults: Record<string, TestResult> = {};
        data.results.forEach((result: TestResult, index: number) => {
          const test = TESTS.find((t) => result.name.includes(t.name));
          if (test) {
            newResults[test.id] = result;
          }
        });
        setResults((prev) => ({ ...prev, ...newResults }));
      }
    } catch (error) {
      console.error('Test group execution error:', error);
    } finally {
      setRunning((prev) => {
        const newSet = new Set(prev);
        newSet.delete(groupId);
        return newSet;
      });
    }
  };

  const toggleExpanded = (testId: string) => {
    setExpanded((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(testId)) {
        newSet.delete(testId);
      } else {
        newSet.add(testId);
      }
      return newSet;
    });
  };

  const getStatusBadge = (status: TestResult['status']) => {
    switch (status) {
      case 'success':
        return <span className="badge badge-success">{t('success')}</span>;
      case 'error':
        return <span className="badge badge-error">{t('failed')}</span>;
      case 'warning':
        return <span className="badge badge-warning">{t('warning')}</span>;
      case 'running':
        return <span className="badge badge-info">{t('running')}</span>;
      default:
        return <span className="badge badge-ghost">{t('pending')}</span>;
    }
  };

  const getTestsByCategory = (category: Test['category']) => {
    return TESTS.filter((test) => test.category === category);
  };

  const getCategoryStats = (category: Test['category']) => {
    const tests = getTestsByCategory(category);
    const total = tests.length;
    const completed = tests.filter(
      (t) =>
        results[t.id]?.status !== 'running' &&
        results[t.id]?.status !== 'pending'
    ).length;
    const success = tests.filter(
      (t) => results[t.id]?.status === 'success'
    ).length;
    const failed = tests.filter(
      (t) => results[t.id]?.status === 'error'
    ).length;
    return { total, completed, success, failed };
  };

  const allStats = {
    total: TESTS.length,
    completed: Object.values(results).filter(
      (r) => r.status !== 'running' && r.status !== 'pending'
    ).length,
    success: Object.values(results).filter((r) => r.status === 'success')
      .length,
    failed: Object.values(results).filter((r) => r.status === 'error').length,
    running: Object.values(results).filter((r) => r.status === 'running')
      .length,
  };

  const getCategoryName = (category: Test['category']) => {
    switch (category) {
      case 'quality':
        return t('categoryQuality');
      case 'unit':
        return t('categoryUnit');
      case 'integration':
        return t('categoryIntegration');
      case 'build':
        return t('categoryBuild');
      case 'coverage':
        return t('categoryCoverage');
      case 'functional':
        return t('categoryFunctional') || 'Functional API Tests';
    }
  };

  const getCategoryIcon = (category: Test['category']) => {
    switch (category) {
      case 'quality':
        return '✨';
      case 'unit':
        return '🧪';
      case 'integration':
        return '🔗';
      case 'build':
        return '🔨';
      case 'coverage':
        return '📊';
      case 'functional':
        return '🌐';
    }
  };

  const renderCategory = (category: Test['category']) => {
    const tests = getTestsByCategory(category);
    const stats = getCategoryStats(category);

    return (
      <div key={category} className="mb-8">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-xl font-semibold flex items-center gap-2">
            <span>{getCategoryIcon(category)}</span>
            <span>{getCategoryName(category)}</span>
            <span className="text-sm font-normal text-base-content/60">
              ({stats.completed}/{stats.total})
            </span>
          </h3>
          {stats.total > 0 && (
            <div className="flex gap-2">
              <span className="badge badge-success">{stats.success}</span>
              <span className="badge badge-error">{stats.failed}</span>
            </div>
          )}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
          {tests.map((test) => {
            const result = results[test.id];
            const isRunning = running.has(test.id);
            const isExpanded = expanded.has(test.id);
            const status = result?.status || 'pending';

            let borderColor = '';
            if (status === 'success') borderColor = 'border-success';
            else if (status === 'error') borderColor = 'border-error';
            else if (status === 'warning') borderColor = 'border-warning';
            else if (status === 'running') borderColor = 'border-info';

            return (
              <div
                key={test.id}
                className={`card bg-base-100 shadow-xl ${borderColor ? `border-2 ${borderColor}` : ''}`}
              >
                <div className="card-body p-4">
                  <h4 className="card-title text-sm">
                    <span>{test.icon}</span>
                    <span>{test.name}</span>
                  </h4>
                  <p className="text-xs text-base-content/70">
                    {test.description}
                  </p>

                  <div className="mt-2 space-y-2">
                    <div className="flex items-center justify-between">
                      {getStatusBadge(status)}
                      {result?.duration && (
                        <span className="text-xs text-base-content/60">
                          {(result.duration / 1000).toFixed(2)}s
                        </span>
                      )}
                    </div>

                    {result?.stats && result.stats.total > 0 && (
                      <div className="flex gap-2 text-xs">
                        <span className="badge badge-success badge-sm">
                          ✓ {result.stats.passed}
                        </span>
                        {result.stats.failed > 0 && (
                          <span className="badge badge-error badge-sm">
                            ✗ {result.stats.failed}
                          </span>
                        )}
                        {result.stats.skipped > 0 && (
                          <span className="badge badge-warning badge-sm">
                            ⊝ {result.stats.skipped}
                          </span>
                        )}
                        <span className="badge badge-ghost badge-sm">
                          Σ {result.stats.total}
                        </span>
                      </div>
                    )}
                  </div>

                  <div className="card-actions mt-3">
                    <button
                      className={`btn btn-xs btn-primary ${isRunning ? 'loading' : ''}`}
                      onClick={() => runTest(test.id)}
                      disabled={isRunning}
                    >
                      {isRunning ? t('running') : t('runTest')}
                    </button>
                    {(result?.output || result?.error) && (
                      <button
                        className="btn btn-xs btn-ghost"
                        onClick={() => toggleExpanded(test.id)}
                      >
                        {isExpanded ? '▼' : '▶'} {t('details')}
                      </button>
                    )}
                  </div>

                  {isExpanded && (result?.output || result?.error) && (
                    <div className="mt-3">
                      {result.error && (
                        <div className="mb-2">
                          <p className="text-xs font-semibold text-error mb-1">
                            {t('error')}:
                          </p>
                          <pre className="text-xs bg-base-200 p-2 rounded overflow-auto max-h-48">
                            {result.error}
                          </pre>
                        </div>
                      )}
                      {result.output && (
                        <div>
                          <p className="text-xs font-semibold mb-1">
                            {t('output')}:
                          </p>
                          <pre className="text-xs bg-base-200 p-2 rounded overflow-auto max-h-48">
                            {result.output}
                          </pre>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    );
  };

  return (
    <div className="container mx-auto px-4 py-8">
      <h1 className="text-3xl font-bold mb-2">{t('title')}</h1>
      <p className="text-base-content/70 mb-6">{t('description')}</p>

      {/* Info Alert */}
      <div className="alert alert-info mb-8">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          className="stroke-current shrink-0 w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
        <span>{t('testSuitesInfo')}</span>
      </div>

      {/* Overall Statistics */}
      <div className="stats shadow w-full mb-8 bg-base-200">
        <div className="stat">
          <div className="stat-title">{t('totalTestSuites')}</div>
          <div className="stat-value text-primary">{allStats.total}</div>
          <div className="stat-desc">{t('testSuitesDesc')}</div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('passed')}</div>
          <div className="stat-value text-success">{allStats.success}</div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('failed')}</div>
          <div className="stat-value text-error">{allStats.failed}</div>
        </div>
        <div className="stat">
          <div className="stat-title">{t('progress')}</div>
          <div className="stat-value text-sm">
            {allStats.total > 0
              ? Math.round((allStats.completed / allStats.total) * 100)
              : 0}
            %
          </div>
          <div className="stat-desc">
            {t('suitesCompleted', {
              count: allStats.completed,
              total: allStats.total,
            })}
          </div>
        </div>
      </div>

      {/* Quick Actions */}
      <div className="flex flex-wrap gap-4 mb-8">
        <button
          className={`btn btn-primary ${running.has('all') ? 'loading' : ''}`}
          onClick={() => runTestGroup('all')}
          disabled={running.has('all')}
        >
          {running.has('all') ? t('running') : '🚀 ' + t('runAll')}
        </button>
        <button
          className={`btn btn-outline ${running.has('all-quality') ? 'loading' : ''}`}
          onClick={() => runTestGroup('all-quality')}
          disabled={running.has('all-quality')}
        >
          {running.has('all-quality')
            ? t('running')
            : '✨ ' + t('runAllQuality')}
        </button>
        <button
          className={`btn btn-outline ${running.has('all-unit') ? 'loading' : ''}`}
          onClick={() => runTestGroup('all-unit')}
          disabled={running.has('all-unit')}
        >
          {running.has('all-unit') ? t('running') : '🧪 ' + t('runAllUnit')}
        </button>
        <button
          className={`btn btn-outline ${running.has('all-builds') ? 'loading' : ''}`}
          onClick={() => runTestGroup('all-builds')}
          disabled={running.has('all-builds')}
        >
          {running.has('all-builds') ? t('running') : '🔨 ' + t('runAllBuilds')}
        </button>
      </div>

      {/* Test Categories */}
      <div className="space-y-8">
        {renderCategory('functional')}
        {renderCategory('quality')}
        {renderCategory('unit')}
        {renderCategory('integration')}
        {renderCategory('build')}
        {renderCategory('coverage')}
      </div>

      {/* Results Summary */}
      {allStats.completed > 0 && (
        <div className="mt-12 alert alert-info">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            className="stroke-current shrink-0 w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            ></path>
          </svg>
          <span>
            {allStats.failed === 0 ? (
              <strong>{t('allGood')}</strong>
            ) : (
              <span>
                {t('testsCompleted', {
                  count: allStats.completed,
                  total: allStats.total,
                })}
                {' - '}
                {allStats.success} {t('passed')}, {allStats.failed}{' '}
                {t('failed')}
              </span>
            )}
          </span>
        </div>
      )}
    </div>
  );
}
