import { NextRequest, NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';

const execPromise = promisify(exec);

export const maxDuration = 600; // 10 minutes timeout for comprehensive tests

interface TestResult {
  name: string;
  status: 'success' | 'error' | 'warning';
  duration: number;
  output: string;
  error?: string;
  stats?: {
    passed: number;
    failed: number;
    skipped: number;
    total: number;
  };
}

// Backend tests
async function runBackendFormat(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && make format',
      {
        timeout: 60000,
      }
    );
    return {
      name: 'Backend Format (make format)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Backend Format (make format)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runBackendLint(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && make lint',
      {
        timeout: 120000,
      }
    );
    const hasIssues = stdout.includes('issues') && !stdout.includes('0 issues');
    return {
      name: 'Backend Lint (make lint)',
      status: hasIssues ? 'warning' : 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Backend Lint (make lint)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runBackendBuild(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && go build ./...',
      {
        timeout: 120000,
      }
    );
    return {
      name: 'Backend Build (go build)',
      status: 'success',
      duration: Date.now() - start,
      output:
        stdout + (stderr ? `\n${stderr}` : '') ||
        'Build successful (no output)',
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Backend Build (go build)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runBackendTestsShort(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && make test-short',
      {
        timeout: 180000,
      }
    );
    return {
      name: 'Backend Short Tests (make test-short)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Backend Short Tests (make test-short)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runBackendTestsUnit(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && make test-unit',
      {
        timeout: 180000,
      }
    );
    return {
      name: 'Backend Unit Tests (make test-unit)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Backend Unit Tests (make test-unit)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runBackendTestsFull(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && make test',
      {
        timeout: 300000,
      }
    );
    return {
      name: 'Backend Full Tests (make test)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Backend Full Tests (make test)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runBackendTestsCoverage(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && make test-coverage',
      {
        timeout: 300000,
      }
    );
    return {
      name: 'Backend Test Coverage (make test-coverage)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Backend Test Coverage (make test-coverage)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runBackendTestsPostExpress(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && make test-postexpress',
      {
        timeout: 60000,
      }
    );
    return {
      name: 'Backend Post Express Test (make test-postexpress)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Backend Post Express Test (make test-postexpress)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runBackendTestsCache(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && go test ./internal/cache -v',
      {
        timeout: 120000,
      }
    );
    const stats = parseGoTestOutput(stdout);
    return {
      name: 'Redis Cache Integration (go test ./internal/cache)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
      stats,
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    const stats = parseGoTestOutput(err.stdout || '');
    return {
      name: 'Redis Cache Integration (go test ./internal/cache)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
      stats,
    };
  }
}

async function runBackendTestsOpenSearch(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/backend && go test ./pkg/transliteration -v -run Integration',
      {
        timeout: 120000,
      }
    );
    const stats = parseGoTestOutput(stdout);
    return {
      name: 'OpenSearch Integration (transliteration tests)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
      stats,
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    const stats = parseGoTestOutput(err.stdout || '');
    return {
      name: 'OpenSearch Integration (transliteration tests)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
      stats,
    };
  }
}

// Helper function to parse Go test output
function parseGoTestOutput(output: string): {
  passed: number;
  failed: number;
  skipped: number;
  total: number;
} {
  const passMatches = output.match(/--- PASS:/g);
  const failMatches = output.match(/--- FAIL:/g);
  const skipMatches = output.match(/--- SKIP:/g);

  const passed = passMatches ? passMatches.length : 0;
  const failed = failMatches ? failMatches.length : 0;
  const skipped = skipMatches ? skipMatches.length : 0;

  return {
    passed,
    failed,
    skipped,
    total: passed + failed + skipped,
  };
}

// Frontend tests
async function runFrontendFormat(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/frontend/svetu && yarn format',
      {
        timeout: 60000,
      }
    );
    return {
      name: 'Frontend Format (yarn format)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Frontend Format (yarn format)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runFrontendLint(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/frontend/svetu && yarn lint',
      {
        timeout: 120000,
      }
    );
    const hasErrors = stdout.includes('error') || stderr.includes('error');
    return {
      name: 'Frontend Lint (yarn lint)',
      status: hasErrors ? 'error' : 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Frontend Lint (yarn lint)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runFrontendTests(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/frontend/svetu && yarn test --watchAll=false',
      {
        timeout: 180000,
      }
    );
    const hasFailed = stdout.includes('failed') || stderr.includes('failed');
    return {
      name: 'Frontend Tests (yarn test)',
      status: hasFailed ? 'error' : 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Frontend Tests (yarn test)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runFrontendBuild(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/frontend/svetu && yarn build',
      {
        timeout: 180000,
      }
    );
    return {
      name: 'Frontend Build (yarn build)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Frontend Build (yarn build)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runFrontendTestsCoverage(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/frontend/svetu && yarn test:coverage --watchAll=false',
      {
        timeout: 240000,
      }
    );
    const hasFailed = stdout.includes('failed') || stderr.includes('failed');
    return {
      name: 'Frontend Test Coverage (yarn test:coverage)',
      status: hasFailed ? 'error' : 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Frontend Test Coverage (yarn test:coverage)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

async function runTypeScriptCheck(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise(
      'cd /data/hostel-booking-system/frontend/svetu && yarn tsc --noEmit',
      {
        timeout: 120000,
      }
    );
    return {
      name: 'TypeScript Check (tsc --noEmit)',
      status: 'success',
      duration: Date.now() - start,
      output:
        stdout + (stderr ? `\n${stderr}` : '') || 'No TypeScript errors found',
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'TypeScript Check (tsc --noEmit)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { test } = body;

    if (!test) {
      return NextResponse.json(
        { error: 'Test name required' },
        { status: 400 }
      );
    }

    let result: TestResult;

    switch (test) {
      // Code Quality
      case 'backend-format':
        result = await runBackendFormat();
        break;
      case 'backend-lint':
        result = await runBackendLint();
        break;
      case 'frontend-format':
        result = await runFrontendFormat();
        break;
      case 'frontend-lint':
        result = await runFrontendLint();
        break;

      // Unit Tests
      case 'backend-tests-unit':
        result = await runBackendTestsUnit();
        break;
      case 'frontend-tests':
        result = await runFrontendTests();
        break;

      // Integration Tests
      case 'backend-tests-short':
        result = await runBackendTestsShort();
        break;
      case 'backend-tests-full':
        result = await runBackendTestsFull();
        break;
      case 'backend-tests-postexpress':
        result = await runBackendTestsPostExpress();
        break;
      case 'backend-tests-cache':
        result = await runBackendTestsCache();
        break;
      case 'backend-tests-opensearch':
        result = await runBackendTestsOpenSearch();
        break;

      // Build & Type Checking
      case 'backend-build':
        result = await runBackendBuild();
        break;
      case 'frontend-build':
        result = await runFrontendBuild();
        break;
      case 'typescript-check':
        result = await runTypeScriptCheck();
        break;

      // Coverage
      case 'backend-tests-coverage':
        result = await runBackendTestsCoverage();
        break;
      case 'frontend-tests-coverage':
        result = await runFrontendTestsCoverage();
        break;

      // Run all quality tests (fast)
      case 'all-quality':
        const qualityResults = await Promise.all([
          runBackendFormat(),
          runBackendLint(),
          runFrontendFormat(),
          runFrontendLint(),
        ]);
        return NextResponse.json({ results: qualityResults });

      // Run all unit tests
      case 'all-unit':
        const unitResults = await Promise.all([
          runBackendTestsUnit(),
          runFrontendTests(),
        ]);
        return NextResponse.json({ results: unitResults });

      // Run all builds
      case 'all-builds':
        const buildResults = await Promise.all([
          runBackendBuild(),
          runFrontendBuild(),
          runTypeScriptCheck(),
        ]);
        return NextResponse.json({ results: buildResults });

      // Run everything (comprehensive)
      case 'all':
        const results = await Promise.all([
          runBackendFormat(),
          runBackendLint(),
          runBackendBuild(),
          runBackendTestsUnit(),
          runFrontendFormat(),
          runFrontendLint(),
          runFrontendTests(),
          runFrontendBuild(),
          runTypeScriptCheck(),
        ]);
        return NextResponse.json({ results });

      default:
        return NextResponse.json(
          { error: 'Invalid test name' },
          { status: 400 }
        );
    }

    return NextResponse.json({ result });
  } catch (error: unknown) {
    const err = error as Error;
    console.error('Test execution error:', err);
    return NextResponse.json(
      { error: err.message || 'Internal server error' },
      { status: 500 }
    );
  }
}
