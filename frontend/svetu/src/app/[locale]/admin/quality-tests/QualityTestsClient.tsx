'use client';

import { useTranslations } from 'next-intl';
import { useState, useEffect } from 'react';
import Link from 'next/link';
import { apiClient } from '@/services/api-client';

interface Test {
  id: string;
  name: string;
  description: string;
  category:
    | 'quality'
    | 'unit'
    | 'integration'
    | 'build'
    | 'coverage'
    | 'functional'
    | 'security'
    | 'performance'
    | 'data-integrity'
    | 'e2e'
    | 'monitoring'
    | 'accessibility';
  icon: string;
  localOnly?: boolean; // Only works on localhost (uses Next.js API route)
}

interface BackendTestResult {
  id: number;
  test_run_id: number;
  test_name: string;
  test_suite: string;
  status: 'passed' | 'failed' | 'skipped';
  duration_ms: number;
  error_msg?: string | null;
  stack_trace?: string | null;
  started_at: string;
  completed_at: string;
}

interface BackendTestLog {
  id: number;
  test_run_id: number;
  level: string;
  message: string;
  timestamp: string;
}

interface BackendTestRunDetail {
  id: number;
  test_suite: string;
  status: string;
  total_tests: number;
  passed_tests: number;
  failed_tests: number;
  skipped_tests: number;
  duration_ms: number;
  started_at: string;
  completed_at?: string | null;
  results?: BackendTestResult[];
  logs?: BackendTestLog[];
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
  failedTests?: BackendTestResult[];
  logs?: BackendTestLog[];
}

const TESTS: Test[] = [
  // Code Quality (LOCAL ONLY - requires Next.js API route)
  {
    id: 'backend-format',
    name: 'Backend Format',
    description: 'Check Go code formatting (gofumpt, goimports)',
    category: 'quality',
    icon: '🎨',
    localOnly: true,
  },
  {
    id: 'backend-lint',
    name: 'Backend Lint',
    description: 'Run golangci-lint for code quality',
    category: 'quality',
    icon: '🔍',
    localOnly: true,
  },
  {
    id: 'frontend-format',
    name: 'Frontend Format',
    description: 'Check Prettier formatting',
    category: 'quality',
    icon: '✨',
    localOnly: true,
  },
  {
    id: 'frontend-lint',
    name: 'Frontend Lint',
    description: 'Run ESLint for code quality',
    category: 'quality',
    icon: '🔎',
    localOnly: true,
  },

  // Unit Tests (LOCAL ONLY - requires Next.js API route)
  {
    id: 'backend-tests-unit',
    name: 'Backend Unit Tests',
    description: 'Run Go unit tests only',
    category: 'unit',
    icon: '🧪',
    localOnly: true,
  },
  {
    id: 'frontend-tests',
    name: 'Frontend Unit Tests',
    description: 'Run Jest unit tests',
    category: 'unit',
    icon: '🔬',
    localOnly: true,
  },

  // New Frontend Unit Tests (LOCAL ONLY - from test coverage improvement plan)
  {
    id: 'frontend-unit-autocomplete-field',
    name: 'AutocompleteAttributeField Tests',
    description:
      'Unit tests for AutocompleteAttributeField component (40 tests)',
    category: 'unit',
    icon: '🎯',
    localOnly: true,
  },
  {
    id: 'frontend-unit-autocomplete-hook',
    name: 'useAttributeAutocomplete Hook Tests',
    description: 'Unit tests for autocomplete hook (35 tests)',
    category: 'unit',
    icon: '🪝',
    localOnly: true,
  },
  {
    id: 'frontend-unit-cars-service',
    name: 'Cars Service Tests',
    description: 'Unit tests for cars API service (45 tests)',
    category: 'unit',
    icon: '🚗',
    localOnly: true,
  },
  {
    id: 'frontend-unit-icon-mapper',
    name: 'Icon Mapper Tests',
    description: 'Unit tests for icon mapping utility (80 tests)',
    category: 'unit',
    icon: '🎨',
    localOnly: true,
  },
  {
    id: 'frontend-unit-env-utils',
    name: 'Environment Utils Tests',
    description: 'Unit tests for environment utilities (60 tests)',
    category: 'unit',
    icon: '⚙️',
    localOnly: true,
  },

  // Integration Tests (LOCAL ONLY - requires Next.js API route)
  {
    id: 'backend-tests-short',
    name: 'Backend Short Tests',
    description: 'Quick tests (skip Integration/OpenSearch)',
    category: 'integration',
    icon: '⚡',
    localOnly: true,
  },
  {
    id: 'backend-tests-full',
    name: 'Backend Full Tests',
    description: 'All tests including integration',
    category: 'integration',
    icon: '🔥',
    localOnly: true,
  },
  {
    id: 'backend-tests-postexpress',
    name: 'Post Express Integration',
    description: 'Test Post Express API integration',
    category: 'integration',
    icon: '📮',
    localOnly: true,
  },

  // Build & Type Checking (LOCAL ONLY - requires Next.js API route)
  {
    id: 'backend-build',
    name: 'Backend Build',
    description: 'Compile Go code',
    category: 'build',
    icon: '🔨',
    localOnly: true,
  },
  {
    id: 'frontend-build',
    name: 'Frontend Build',
    description: 'Next.js production build',
    category: 'build',
    icon: '⚙️',
    localOnly: true,
  },
  {
    id: 'typescript-check',
    name: 'TypeScript Check',
    description: 'Type checking (tsc --noEmit)',
    category: 'build',
    icon: '📘',
    localOnly: true,
  },

  // Coverage (LOCAL ONLY - requires Next.js API route)
  {
    id: 'backend-tests-coverage',
    name: 'Backend Test Coverage',
    description: 'Run tests with coverage report',
    category: 'coverage',
    icon: '📊',
    localOnly: true,
  },
  {
    id: 'frontend-tests-coverage',
    name: 'Frontend Test Coverage',
    description: 'Jest tests with coverage',
    category: 'coverage',
    icon: '📈',
    localOnly: true,
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
  {
    id: 'api-review-creation',
    name: 'Review Creation',
    description: 'Test review creation with rating (draft + publish)',
    category: 'functional',
    icon: '⭐',
  },

  // Negative Test Cases
  {
    id: 'api-auth-invalid-token',
    name: 'Invalid Token Test',
    description: 'Test API rejection with invalid authentication token',
    category: 'functional',
    icon: '🚫',
  },
  {
    id: 'api-auth-missing-token',
    name: 'Missing Token Test',
    description: 'Test API rejection when authentication token is missing',
    category: 'functional',
    icon: '❌',
  },
  {
    id: 'api-admin-unauthorized',
    name: 'Admin Unauthorized',
    description: 'Test admin endpoint rejection for non-admin users',
    category: 'functional',
    icon: '🔒',
  },
  {
    id: 'api-search-invalid-params',
    name: 'Invalid Search Params',
    description: 'Test handling of invalid search parameters',
    category: 'functional',
    icon: '⚠️',
  },

  // Edge Cases
  {
    id: 'api-search-empty-query',
    name: 'Empty Search Query',
    description: 'Test search with empty query string',
    category: 'functional',
    icon: '🔎',
  },
  {
    id: 'api-search-unicode',
    name: 'Unicode Search',
    description: 'Test search with Unicode characters (Cyrillic, Emoji)',
    category: 'functional',
    icon: '🌐',
  },
  {
    id: 'api-listings-extreme-limit',
    name: 'Extreme Limit Values',
    description: 'Test listings with extreme limit values (0, 10000)',
    category: 'functional',
    icon: '🔢',
  },

  // Security Tests
  {
    id: 'security-sql-injection',
    name: 'SQL Injection Protection',
    description: 'Test SQL injection attempts in search and filters',
    category: 'security',
    icon: '💉',
  },
  {
    id: 'security-xss-protection',
    name: 'XSS Protection',
    description: 'Test XSS attempts in user inputs (listings, reviews)',
    category: 'security',
    icon: '🛡️',
  },
  {
    id: 'security-file-upload-validation',
    name: 'File Upload Security',
    description: 'Test file type and size validation, malicious file rejection',
    category: 'security',
    icon: '📎',
  },
  {
    id: 'security-auth-session-expiry',
    name: 'Session Expiry Test',
    description: 'Test JWT token expiration and refresh logic',
    category: 'security',
    icon: '⏰',
  },
  {
    id: 'security-api-rate-limiting',
    name: 'Rate Limiting Enforcement',
    description: 'Test rate limiting enforcement on API endpoints',
    category: 'security',
    icon: '🚦',
  },
  {
    id: 'security-csrf-protection',
    name: 'CSRF Protection',
    description: 'Test CSRF token validation on state-changing requests',
    category: 'security',
    icon: '🔐',
  },

  // Performance Tests
  {
    id: 'performance-api-response-time',
    name: 'API Response Time',
    description: 'Measure API endpoint response times (should be <200ms)',
    category: 'performance',
    icon: '⚡',
  },
  {
    id: 'performance-concurrent-users',
    name: 'Concurrent Users Test',
    description: 'Test system with 10/50/100 concurrent users',
    category: 'performance',
    icon: '👥',
  },
  {
    id: 'performance-database-queries',
    name: 'Database Query Performance',
    description: 'Check for slow database queries (>100ms)',
    category: 'performance',
    icon: '🗄️',
  },
  {
    id: 'performance-memory-usage',
    name: 'Memory Usage Monitoring',
    description: 'Monitor memory usage during test execution',
    category: 'performance',
    icon: '🧠',
  },

  // Data Integrity Tests
  {
    id: 'data-integrity-marketplace-listing',
    name: 'Listing Data Consistency',
    description:
      'Verify listing data matches across DB, cache, and search index',
    category: 'data-integrity',
    icon: '🔄',
  },
  {
    id: 'data-integrity-transaction-rollback',
    name: 'Transaction Rollback Test',
    description: 'Test database transaction rollback on errors',
    category: 'data-integrity',
    icon: '↩️',
  },
  {
    id: 'data-integrity-image-orphan-cleanup',
    name: 'Orphan Image Cleanup',
    description: 'Verify orphaned images are cleaned up from MinIO',
    category: 'data-integrity',
    icon: '🗑️',
  },

  // E2E Tests
  {
    id: 'e2e-user-journey-create-listing',
    name: 'User Journey: Create Listing',
    description: 'Full flow: login → create listing → upload images → publish',
    category: 'e2e',
    icon: '🎬',
  },
  {
    id: 'e2e-user-journey-search-contact',
    name: 'User Journey: Search & Contact',
    description: 'Search → view listing → contact seller',
    category: 'e2e',
    icon: '🛍️',
  },
  {
    id: 'e2e-admin-moderation',
    name: 'Admin Moderation Flow',
    description: 'Admin reviews and approves/rejects listing',
    category: 'e2e',
    icon: '👨‍⚖️',
  },

  // Integration Tests
  {
    id: 'integration-redis-cache',
    name: 'Redis Cache Test',
    description: 'Test Redis cache operations (SET, GET, TTL)',
    category: 'integration',
    icon: '💾',
  },
  {
    id: 'integration-opensearch-index',
    name: 'OpenSearch Test',
    description: 'Test OpenSearch indexing and search functionality',
    category: 'integration',
    icon: '🔍',
  },
  {
    id: 'integration-postgres-connection',
    name: 'PostgreSQL Test',
    description: 'Test PostgreSQL connection and queries',
    category: 'integration',
    icon: '🐘',
  },

  // Monitoring & Observability Tests
  {
    id: 'monitoring-health-endpoints',
    name: 'Health Check Endpoints',
    description: 'Test /health/live and /health/ready endpoints',
    category: 'monitoring',
    icon: '💓',
  },
  {
    id: 'monitoring-metrics-collection',
    name: 'Metrics Collection',
    description: 'Verify Prometheus metrics are being collected',
    category: 'monitoring',
    icon: '📊',
  },
  {
    id: 'monitoring-error-logging',
    name: 'Error Logging Test',
    description: 'Verify errors are properly logged with context',
    category: 'monitoring',
    icon: '📝',
  },

  // Accessibility Tests
  {
    id: 'a11y-wcag-compliance',
    name: 'WCAG 2.1 Compliance',
    description: 'Test WCAG 2.1 AA compliance using axe-core',
    category: 'accessibility',
    icon: '♿',
  },
  {
    id: 'a11y-keyboard-navigation',
    name: 'Keyboard Navigation',
    description: 'Test keyboard navigation on all interactive elements',
    category: 'accessibility',
    icon: '⌨️',
  },
];

const STORAGE_KEY = 'quality-tests-results';

interface QualityTestsClientProps {
  locale: string;
}

export default function QualityTestsClient({
  locale: _locale,
}: QualityTestsClientProps) {
  const t = useTranslations('admin.qualityTests');

  // Initialize with empty state (same for SSR and CSR)
  const [results, setResults] = useState<Record<string, TestResult>>({});
  const [running, setRunning] = useState<Set<string>>(new Set());
  const [expanded, setExpanded] = useState<Set<string>>(new Set());
  const [isHydrated, setIsHydrated] = useState(false);

  // Load from localStorage after hydration (client-side only)
  useEffect(() => {
    setIsHydrated(true);
    try {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) {
        setResults(JSON.parse(stored));
      }
    } catch (error) {
      console.error('Failed to load test results from localStorage:', error);
    }
  }, []);

  // Save results to localStorage whenever they change (client-side only)
  useEffect(() => {
    if (!isHydrated) return; // Don't save during hydration
    try {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(results));
    } catch (error) {
      console.error('Failed to save test results to localStorage:', error);
    }
  }, [results, isHydrated]);

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
      const isSecurity = test?.category === 'security';
      const isPerformance = test?.category === 'performance';
      const isDataIntegrity = test?.category === 'data-integrity';
      const isMonitoring = test?.category === 'monitoring';
      const isE2E = test?.category === 'e2e';
      const isAccessibility = test?.category === 'accessibility';

      if (
        isFunctional ||
        isSecurity ||
        isPerformance ||
        isDataIntegrity ||
        isMonitoring ||
        isE2E ||
        isAccessibility
      ) {
        // Functional, Security, Performance, Data Integrity, Monitoring, E2E и Accessibility тесты: вызываем backend API через apiClient
        // Передаем test_name для запуска конкретного теста
        const testSuite = isSecurity
          ? 'security'
          : isPerformance
            ? 'performance'
            : isDataIntegrity
              ? 'data-integrity'
              : isMonitoring
                ? 'monitoring'
                : isE2E
                  ? 'e2e'
                  : isAccessibility
                    ? 'accessibility'
                    : 'api-endpoints';
        const response = await apiClient.post('/admin/tests/run', {
          test_suite: testSuite,
          test_name: testId,
          parallel: false,
        });

        // Backend API возвращает test run ID - нужно подождать результаты
        if (response.data) {
          const runId = response.data.test_run_id;

          // Polling для получения результатов
          let attempts = 0;
          const maxAttempts = 1500; // 25 минут (для долгих тестов типа accessibility)

          while (attempts < maxAttempts) {
            await new Promise((resolve) => setTimeout(resolve, 1000));

            const detailResponse = await apiClient.get(
              `/admin/tests/runs/${runId}`
            );
            const detail = detailResponse.data as BackendTestRunDetail;

            if (detail.status === 'completed' || detail.status === 'failed') {
              // Получаем упавшие тесты из results
              const failedTests =
                detail.results?.filter((r) => r.status === 'failed') || [];

              // Формируем детальное сообщение об ошибке
              let errorMessage = '';
              if (failedTests.length > 0) {
                errorMessage = `${failedTests.length} test(s) failed:\n\n`;
                failedTests.forEach((test, idx) => {
                  errorMessage += `${idx + 1}. ${test.test_name}`;
                  if (test.error_msg) {
                    errorMessage += `\n   Error: ${test.error_msg}`;
                  }
                  errorMessage += '\n\n';
                });
              }

              // Тесты из backend выполняются как единый suite
              // Показываем общий результат с статистикой
              setResults((prev) => ({
                ...prev,
                [testId]: {
                  name: testId,
                  status: detail.failed_tests > 0 ? 'error' : 'success',
                  duration: detail.duration_ms,
                  output: `Test suite completed: ${detail.passed_tests} passed, ${detail.failed_tests} failed`,
                  error: errorMessage || undefined,
                  stats: {
                    passed: detail.passed_tests,
                    failed: detail.failed_tests,
                    skipped: detail.skipped_tests,
                    total: detail.total_tests,
                  },
                  failedTests: failedTests,
                  logs: detail.logs,
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
        data.results.forEach((result: TestResult) => {
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

  const clearResults = () => {
    setResults({});
    setExpanded(new Set());
    try {
      localStorage.removeItem(STORAGE_KEY);
    } catch (error) {
      // localStorage not available (SSR or private mode)
      console.error('Failed to clear localStorage:', error);
    }
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
      case 'security':
        return t('categorySecurity') || 'Security Tests';
      case 'performance':
        return t('categoryPerformance') || 'Performance Tests';
      case 'data-integrity':
        return t('categoryDataIntegrity') || 'Data Integrity Tests';
      case 'e2e':
        return t('categoryE2E') || 'End-to-End Tests';
      case 'monitoring':
        return t('categoryMonitoring') || 'Monitoring & Observability';
      case 'accessibility':
        return t('categoryAccessibility') || 'Accessibility Tests';
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
      case 'security':
        return '🔒';
      case 'performance':
        return '⚡';
      case 'data-integrity':
        return '🔄';
      case 'e2e':
        return '🎬';
      case 'monitoring':
        return '📊';
      case 'accessibility':
        return '♿';
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
                    <span>{t(`tests.${test.id}.name`) || test.name}</span>
                    {test.localOnly && (
                      <span
                        className="ml-2 cursor-help"
                        title={t('localOnlyTooltip')}
                      >
                        🏠
                      </span>
                    )}
                  </h4>
                  <p className="text-xs text-base-content/70">
                    {t(`tests.${test.id}.description`) || test.description}
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

                  {isExpanded &&
                    (result?.output ||
                      result?.error ||
                      result?.failedTests) && (
                      <div className="mt-3 space-y-3">
                        {result.error && (
                          <div>
                            <p className="text-xs font-semibold text-error mb-1">
                              {t('error')}:
                            </p>
                            <pre className="text-xs bg-base-200 p-2 rounded overflow-auto max-h-48">
                              {result.error}
                            </pre>
                          </div>
                        )}
                        {result.failedTests &&
                          result.failedTests.length > 0 && (
                            <div>
                              <p className="text-xs font-semibold text-error mb-2">
                                Детальная информация об упавших тестах:
                              </p>
                              <div className="space-y-2">
                                {result.failedTests.map((failedTest, idx) => (
                                  <div
                                    key={failedTest.id}
                                    className="bg-error/10 p-2 rounded border border-error/20"
                                  >
                                    <p className="text-xs font-semibold mb-1">
                                      {idx + 1}. {failedTest.test_name}
                                    </p>
                                    {failedTest.error_msg && (
                                      <div className="mb-1">
                                        <p className="text-xs text-error/80">
                                          Error:
                                        </p>
                                        <pre className="text-xs bg-base-200 p-1 rounded overflow-auto max-h-24">
                                          {failedTest.error_msg}
                                        </pre>
                                      </div>
                                    )}
                                    {failedTest.stack_trace && (
                                      <div>
                                        <p className="text-xs text-error/80">
                                          Stack trace:
                                        </p>
                                        <pre className="text-xs bg-base-200 p-1 rounded overflow-auto max-h-32 font-mono">
                                          {failedTest.stack_trace}
                                        </pre>
                                      </div>
                                    )}
                                    <p className="text-xs text-base-content/60 mt-1">
                                      Duration: {failedTest.duration_ms}ms
                                    </p>
                                  </div>
                                ))}
                              </div>
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
                        {result.logs && result.logs.length > 0 && (
                          <div>
                            <p className="text-xs font-semibold mb-1">Logs:</p>
                            <div className="text-xs bg-base-200 p-2 rounded overflow-auto max-h-48 space-y-1">
                              {result.logs.map((log) => (
                                <div
                                  key={log.id}
                                  className={
                                    log.level === 'error'
                                      ? 'text-error'
                                      : log.level === 'warn'
                                        ? 'text-warning'
                                        : 'text-base-content/80'
                                  }
                                >
                                  <span className="font-mono text-base-content/60">
                                    [
                                    {new Date(
                                      log.timestamp
                                    ).toLocaleTimeString()}
                                    ]
                                  </span>{' '}
                                  <span className="font-semibold">
                                    [{log.level.toUpperCase()}]
                                  </span>{' '}
                                  {log.message}
                                </div>
                              ))}
                            </div>
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
      <div className="flex items-center justify-between mb-2">
        <h1 className="text-3xl font-bold">{t('title')}</h1>
        <Link
          href="/admin/test-data-cleanup"
          className="btn btn-outline btn-sm gap-2"
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="h-5 w-5"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M20.25 6.375c0 2.278-3.694 4.125-8.25 4.125S3.75 8.653 3.75 6.375m16.5 0c0-2.278-3.694-4.125-8.25-4.125S3.75 4.097 3.75 6.375m16.5 0v11.25c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125V6.375m16.5 0v3.75m-16.5-3.75v3.75m16.5 0v3.75C20.25 16.153 16.556 18 12 18s-8.25-1.847-8.25-4.125v-3.75m16.5 0c0 2.278-3.694 4.125-8.25 4.125s-8.25-1.847-8.25-4.125"
            />
          </svg>
          {t('cleanup_test_data')}
        </Link>
      </div>
      <p className="text-base-content/70 mb-6">{t('description')}</p>

      {/* Info Alert */}
      <div className="alert alert-info mb-4">
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

      {/* Local Only Tests Warning */}
      <div className="alert alert-warning mb-8">
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
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
          ></path>
        </svg>
        <div>
          <div className="font-bold">{t('localOnlyTestsTitle')}</div>
          <div className="text-sm">{t('localOnlyTestsDescription')}</div>
        </div>
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
        {allStats.completed > 0 && (
          <button className="btn btn-outline btn-error" onClick={clearResults}>
            🗑️ {t('clearResults') || 'Clear Results'}
          </button>
        )}
      </div>

      {/* Test Categories */}
      <div className="space-y-8">
        {renderCategory('functional')}
        {renderCategory('security')}
        {renderCategory('performance')}
        {renderCategory('data-integrity')}
        {renderCategory('e2e')}
        {renderCategory('monitoring')}
        {renderCategory('accessibility')}
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
