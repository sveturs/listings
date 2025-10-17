import { NextRequest, NextResponse } from 'next/server';
import { exec } from 'child_process';
import { promisify } from 'util';

const execPromise = promisify(exec);

export const maxDuration = 300; // 5 minutes timeout

interface TestResult {
  name: string;
  status: 'success' | 'error' | 'warning';
  duration: number;
  output: string;
  error?: string;
}

// Backend tests
async function runBackendFormat(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise('cd /data/hostel-booking-system/backend && make format', {
      timeout: 60000,
    });
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
    const { stdout, stderr } = await execPromise('cd /data/hostel-booking-system/backend && make lint', {
      timeout: 120000,
    });
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
    const { stdout, stderr } = await execPromise('cd /data/hostel-booking-system/backend && go build ./...', {
      timeout: 120000,
    });
    return {
      name: 'Backend Build (go build)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : '') || 'Build successful (no output)',
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

async function runBackendTests(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise('cd /data/hostel-booking-system/backend && make test-short', {
      timeout: 180000,
    });
    return {
      name: 'Backend Tests (make test-short)',
      status: 'success',
      duration: Date.now() - start,
      output: stdout + (stderr ? `\n${stderr}` : ''),
    };
  } catch (error: unknown) {
    const err = error as { message: string; stdout?: string; stderr?: string };
    return {
      name: 'Backend Tests (make test-short)',
      status: 'error',
      duration: Date.now() - start,
      output: err.stdout || '',
      error: err.message + (err.stderr ? `\n${err.stderr}` : ''),
    };
  }
}

// Frontend tests
async function runFrontendFormat(): Promise<TestResult> {
  const start = Date.now();
  try {
    const { stdout, stderr } = await execPromise('cd /data/hostel-booking-system/frontend/svetu && yarn format', {
      timeout: 60000,
    });
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
    const { stdout, stderr } = await execPromise('cd /data/hostel-booking-system/frontend/svetu && yarn lint', {
      timeout: 120000,
    });
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
    const { stdout, stderr } = await execPromise('cd /data/hostel-booking-system/frontend/svetu && yarn build', {
      timeout: 180000,
    });
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

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();
    const { test } = body;

    if (!test) {
      return NextResponse.json({ error: 'Test name required' }, { status: 400 });
    }

    let result: TestResult;

    switch (test) {
      case 'backend-format':
        result = await runBackendFormat();
        break;
      case 'backend-lint':
        result = await runBackendLint();
        break;
      case 'backend-build':
        result = await runBackendBuild();
        break;
      case 'backend-tests':
        result = await runBackendTests();
        break;
      case 'frontend-format':
        result = await runFrontendFormat();
        break;
      case 'frontend-lint':
        result = await runFrontendLint();
        break;
      case 'frontend-tests':
        result = await runFrontendTests();
        break;
      case 'frontend-build':
        result = await runFrontendBuild();
        break;
      case 'all':
        // Run all tests
        const results = await Promise.all([
          runBackendFormat(),
          runBackendLint(),
          runBackendBuild(),
          runBackendTests(),
          runFrontendFormat(),
          runFrontendLint(),
          runFrontendTests(),
          runFrontendBuild(),
        ]);
        return NextResponse.json({ results });
      default:
        return NextResponse.json({ error: 'Invalid test name' }, { status: 400 });
    }

    return NextResponse.json({ result });
  } catch (error: unknown) {
    const err = error as Error;
    console.error('Test execution error:', err);
    return NextResponse.json({ error: err.message || 'Internal server error' }, { status: 500 });
  }
}
