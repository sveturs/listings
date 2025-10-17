'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';

interface TestResult {
  name: string;
  status: 'success' | 'error' | 'warning' | 'running' | 'pending';
  duration?: number;
  output?: string;
  error?: string;
}

interface Test {
  id: string;
  name: string;
  description: string;
  category: 'backend' | 'frontend';
  icon: string;
}

const TESTS: Test[] = [
  {
    id: 'backend-format',
    name: 'Backend Format',
    description: 'Check Go code formatting (gofumpt, goimports)',
    category: 'backend',
    icon: '🎨',
  },
  {
    id: 'backend-lint',
    name: 'Backend Lint',
    description: 'Run golangci-lint on Go code',
    category: 'backend',
    icon: '🔍',
  },
  {
    id: 'backend-build',
    name: 'Backend Build',
    description: 'Compile all Go packages',
    category: 'backend',
    icon: '🔨',
  },
  {
    id: 'backend-tests',
    name: 'Backend Tests',
    description: 'Run Go unit tests',
    category: 'backend',
    icon: '🧪',
  },
  {
    id: 'frontend-format',
    name: 'Frontend Format',
    description: 'Check code formatting (Prettier)',
    category: 'frontend',
    icon: '🎨',
  },
  {
    id: 'frontend-lint',
    name: 'Frontend Lint',
    description: 'Run ESLint on TypeScript code',
    category: 'frontend',
    icon: '🔍',
  },
  {
    id: 'frontend-tests',
    name: 'Frontend Tests',
    description: 'Run Jest unit tests',
    category: 'frontend',
    icon: '🧪',
  },
  {
    id: 'frontend-build',
    name: 'Frontend Build',
    description: 'Build production Next.js app',
    category: 'frontend',
    icon: '🔨',
  },
];

export default function QualityTestsClient({ locale }: { locale: string }) {
  const t = useTranslations('admin.qualityTests');
  const [results, setResults] = useState<Record<string, TestResult>>({});
  const [running, setRunning] = useState<Set<string>>(new Set());
  const [expandedTests, setExpandedTests] = useState<Set<string>>(new Set());
  const [runningAll, setRunningAll] = useState(false);

  const runTest = async (testId: string) => {
    setRunning((prev) => new Set(prev).add(testId));
    setResults((prev) => ({
      ...prev,
      [testId]: { name: testId, status: 'running' },
    }));

    try {
      const response = await fetch('/api/admin/tests', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ test: testId }),
      });

      const data = await response.json();

      if (data.result) {
        setResults((prev) => ({
          ...prev,
          [testId]: data.result,
        }));
      } else {
        setResults((prev) => ({
          ...prev,
          [testId]: {
            name: testId,
            status: 'error',
            error: data.error || 'Unknown error',
          },
        }));
      }
    } catch (error) {
      setResults((prev) => ({
        ...prev,
        [testId]: {
          name: testId,
          status: 'error',
          error: error instanceof Error ? error.message : 'Network error',
        },
      }));
    } finally {
      setRunning((prev) => {
        const next = new Set(prev);
        next.delete(testId);
        return next;
      });
    }
  };

  const runAllTests = async () => {
    setRunningAll(true);
    const allTestIds = TESTS.map((t) => t.id);

    // Mark all as running
    allTestIds.forEach((id) => {
      setRunning((prev) => new Set(prev).add(id));
      setResults((prev) => ({
        ...prev,
        [id]: { name: id, status: 'running' },
      }));
    });

    try {
      const response = await fetch('/api/admin/tests', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ test: 'all' }),
      });

      const data = await response.json();

      if (data.results && Array.isArray(data.results)) {
        const resultsMap: Record<string, TestResult> = {};
        data.results.forEach((result: TestResult, index: number) => {
          const testId = allTestIds[index];
          resultsMap[testId] = result;
        });
        setResults(resultsMap);
      }
    } catch (error) {
      console.error('Error running all tests:', error);
    } finally {
      setRunningAll(false);
      setRunning(new Set());
    }
  };

  const toggleExpand = (testId: string) => {
    setExpandedTests((prev) => {
      const next = new Set(prev);
      if (next.has(testId)) {
        next.delete(testId);
      } else {
        next.add(testId);
      }
      return next;
    });
  };

  const getStatusBadge = (status: TestResult['status']) => {
    switch (status) {
      case 'success':
        return <span className="badge badge-success gap-2">✓ Success</span>;
      case 'error':
        return <span className="badge badge-error gap-2">✗ Failed</span>;
      case 'warning':
        return <span className="badge badge-warning gap-2">⚠ Warning</span>;
      case 'running':
        return (
          <span className="badge badge-info gap-2">
            <span className="loading loading-spinner loading-xs"></span> Running
          </span>
        );
      case 'pending':
      default:
        return <span className="badge badge-ghost gap-2">○ Pending</span>;
    }
  };

  const getStatusColor = (status: TestResult['status']) => {
    switch (status) {
      case 'success':
        return 'border-success';
      case 'error':
        return 'border-error';
      case 'warning':
        return 'border-warning';
      case 'running':
        return 'border-info';
      default:
        return 'border-base-300';
    }
  };

  const backendTests = TESTS.filter((t) => t.category === 'backend');
  const frontendTests = TESTS.filter((t) => t.category === 'frontend');

  const totalTests = TESTS.length;
  const completedTests = Object.values(results).filter(
    (r) => r.status === 'success' || r.status === 'error' || r.status === 'warning'
  ).length;
  const successTests = Object.values(results).filter((r) => r.status === 'success').length;
  const failedTests = Object.values(results).filter((r) => r.status === 'error').length;

  return (
    <div className="container mx-auto px-4 py-8 max-w-7xl">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-4xl font-bold mb-2">🧪 Code Quality Tests</h1>
        <p className="text-base-content/70">
          Comprehensive testing suite for backend and frontend code quality
        </p>
      </div>

      {/* Dashboard */}
      <div className="stats shadow w-full mb-8 bg-base-200">
        <div className="stat">
          <div className="stat-figure text-primary">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="inline-block w-8 h-8 stroke-current"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              ></path>
            </svg>
          </div>
          <div className="stat-title">Total Tests</div>
          <div className="stat-value text-primary">{totalTests}</div>
          <div className="stat-desc">Backend + Frontend</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-success">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="inline-block w-8 h-8 stroke-current"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
          </div>
          <div className="stat-title">Passed</div>
          <div className="stat-value text-success">{successTests}</div>
          <div className="stat-desc">{completedTests > 0 ? `${Math.round((successTests / completedTests) * 100)}% success rate` : 'No tests run yet'}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-error">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="inline-block w-8 h-8 stroke-current"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M6 18L18 6M6 6l12 12"
              ></path>
            </svg>
          </div>
          <div className="stat-title">Failed</div>
          <div className="stat-value text-error">{failedTests}</div>
          <div className="stat-desc">{failedTests > 0 ? 'Need attention' : 'All good!'}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-info">
            <div className="radial-progress text-info" style={{ '--value': completedTests > 0 ? (completedTests / totalTests) * 100 : 0, '--size': '4rem' } as React.CSSProperties}>
              {completedTests > 0 ? Math.round((completedTests / totalTests) * 100) : 0}%
            </div>
          </div>
          <div className="stat-title">Progress</div>
          <div className="stat-value text-info">
            {completedTests}/{totalTests}
          </div>
          <div className="stat-desc">Tests completed</div>
        </div>
      </div>

      {/* Run All Button */}
      <div className="flex justify-center mb-8">
        <button
          onClick={runAllTests}
          disabled={runningAll || running.size > 0}
          className="btn btn-primary btn-lg gap-2"
        >
          {runningAll ? (
            <>
              <span className="loading loading-spinner"></span>
              Running All Tests...
            </>
          ) : (
            <>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              Run All Tests
            </>
          )}
        </button>
      </div>

      {/* Backend Tests */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-4">
          <div className="badge badge-primary badge-lg">Backend</div>
          <h2 className="text-2xl font-bold">Go / Backend Tests</h2>
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {backendTests.map((test) => {
            const result = results[test.id];
            const isRunning = running.has(test.id);
            const isExpanded = expandedTests.has(test.id);
            const hasOutput = result && (result.output || result.error);

            return (
              <div
                key={test.id}
                className={`card bg-base-100 shadow-xl border-2 ${result ? getStatusColor(result.status) : 'border-base-300'}`}
              >
                <div className="card-body">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="card-title text-lg">
                        <span className="text-2xl mr-2">{test.icon}</span>
                        {test.name}
                      </h3>
                      <p className="text-sm text-base-content/70 mt-1">{test.description}</p>
                    </div>
                    <div>{result && getStatusBadge(result.status)}</div>
                  </div>

                  {result && result.duration && (
                    <div className="text-xs text-base-content/60 mt-2">
                      Duration: {(result.duration / 1000).toFixed(2)}s
                    </div>
                  )}

                  <div className="card-actions justify-between items-center mt-4">
                    <div>
                      {hasOutput && (
                        <button onClick={() => toggleExpand(test.id)} className="btn btn-ghost btn-sm gap-2">
                          {isExpanded ? '▼' : '▶'} Details
                        </button>
                      )}
                    </div>
                    <button
                      onClick={() => runTest(test.id)}
                      disabled={isRunning || runningAll}
                      className="btn btn-primary btn-sm"
                    >
                      {isRunning ? (
                        <>
                          <span className="loading loading-spinner loading-xs"></span>
                          Running...
                        </>
                      ) : (
                        'Run Test'
                      )}
                    </button>
                  </div>

                  {isExpanded && hasOutput && (
                    <div className="mt-4 p-4 bg-base-300 rounded-lg overflow-auto max-h-96">
                      <pre className="text-xs whitespace-pre-wrap break-words">
                        {result.output}
                        {result.error && (
                          <div className="text-error mt-2 font-bold">
                            Error: {result.error}
                          </div>
                        )}
                      </pre>
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>

      {/* Frontend Tests */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-4">
          <div className="badge badge-secondary badge-lg">Frontend</div>
          <h2 className="text-2xl font-bold">Next.js / React Tests</h2>
        </div>
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
          {frontendTests.map((test) => {
            const result = results[test.id];
            const isRunning = running.has(test.id);
            const isExpanded = expandedTests.has(test.id);
            const hasOutput = result && (result.output || result.error);

            return (
              <div
                key={test.id}
                className={`card bg-base-100 shadow-xl border-2 ${result ? getStatusColor(result.status) : 'border-base-300'}`}
              >
                <div className="card-body">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <h3 className="card-title text-lg">
                        <span className="text-2xl mr-2">{test.icon}</span>
                        {test.name}
                      </h3>
                      <p className="text-sm text-base-content/70 mt-1">{test.description}</p>
                    </div>
                    <div>{result && getStatusBadge(result.status)}</div>
                  </div>

                  {result && result.duration && (
                    <div className="text-xs text-base-content/60 mt-2">
                      Duration: {(result.duration / 1000).toFixed(2)}s
                    </div>
                  )}

                  <div className="card-actions justify-between items-center mt-4">
                    <div>
                      {hasOutput && (
                        <button onClick={() => toggleExpand(test.id)} className="btn btn-ghost btn-sm gap-2">
                          {isExpanded ? '▼' : '▶'} Details
                        </button>
                      )}
                    </div>
                    <button
                      onClick={() => runTest(test.id)}
                      disabled={isRunning || runningAll}
                      className="btn btn-secondary btn-sm"
                    >
                      {isRunning ? (
                        <>
                          <span className="loading loading-spinner loading-xs"></span>
                          Running...
                        </>
                      ) : (
                        'Run Test'
                      )}
                    </button>
                  </div>

                  {isExpanded && hasOutput && (
                    <div className="mt-4 p-4 bg-base-300 rounded-lg overflow-auto max-h-96">
                      <pre className="text-xs whitespace-pre-wrap break-words">
                        {result.output}
                        {result.error && (
                          <div className="text-error mt-2 font-bold">
                            Error: {result.error}
                          </div>
                        )}
                      </pre>
                    </div>
                  )}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </div>
  );
}
