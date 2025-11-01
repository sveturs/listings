/**
 * k6_canary_test.js
 *
 * k6 load test script for 1% canary traffic validation
 * Tests marketplace search API with realistic traffic patterns
 *
 * Usage:
 *   # Install k6 (if not already installed)
 *   sudo apt-get install k6
 *
 *   # Run 24-hour canary test
 *   k6 run k6_canary_test.js
 *
 *   # Run 1-hour test (for quick validation)
 *   k6 run -e DURATION=1h k6_canary_test.js
 *
 *   # Run with custom VUs
 *   k6 run -e VUS=20 k6_canary_test.js
 */

import http from 'k6/http';
import { sleep, check } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

// Custom metrics
const canaryRequests = new Counter('canary_requests');
const monolithRequests = new Counter('monolith_requests');
const errorRate = new Rate('error_rate');
const latencyP95 = new Trend('latency_p95');

// Configuration
const DURATION = __ENV.DURATION || '24h';  // Default: 24 hours
const VUS = parseInt(__ENV.VUS || '10');   // Default: 10 concurrent users
const BASE_URL = __ENV.BASE_URL || 'https://devapi.svetu.rs';

// Test configuration
export const options = {
  // Duration
  duration: DURATION,

  // Virtual users (constant rate)
  vus: VUS,

  // Thresholds (validation criteria)
  thresholds: {
    // HTTP request duration: 95% of requests should complete within 200ms
    http_req_duration: ['p(95)<200'],

    // Error rate: < 0.1% (1 error per 1000 requests)
    http_req_failed: ['rate<0.001'],

    // Request rate: at least 1 req/s
    http_reqs: ['rate>1'],

    // Custom metrics
    'error_rate': ['rate<0.001'],
    'latency_p95': ['p(95)<200'],
  },

  // Graceful shutdown
  gracefulStop: '30s',

  // Summary
  summaryTrendStats: ['avg', 'min', 'med', 'max', 'p(90)', 'p(95)', 'p(99)'],
};

/**
 * Test scenario: Random marketplace searches
 */
export default function () {
  // Random search queries (realistic user behavior)
  const queries = [
    'test',
    'laptop',
    'phone',
    'apartment',
    'car',
    'bike',
    'clothes',
    'furniture',
    'electronics',
    'books',
  ];

  const randomQuery = queries[Math.floor(Math.random() * queries.length)];

  // Random category IDs (from svetu marketplace)
  const categories = [
    1301,  // Example category
    1302,
    1303,
    null,  // No category filter
  ];

  const randomCategory = categories[Math.floor(Math.random() * categories.length)];

  // Build URL
  let url = `${BASE_URL}/api/v1/marketplace/search?q=${randomQuery}&limit=20`;
  if (randomCategory) {
    url += `&category_id=${randomCategory}`;
  }

  // Execute request
  const response = http.get(url, {
    tags: {
      name: 'marketplace_search',
      query: randomQuery,
    },
  });

  // Record metrics
  const success = check(response, {
    'status is 200': (r) => r.status === 200,
    'response has data': (r) => r.json('data') !== undefined,
    'latency < 200ms': (r) => r.timings.duration < 200,
    'latency < 500ms': (r) => r.timings.duration < 500,
  });

  // Track error rate
  errorRate.add(!success);

  // Track latency
  latencyP95.add(response.timings.duration);

  // Detect canary vs monolith (by response headers or other markers)
  // NOTE: This requires backend to add a response header like "X-Served-By: microservice|monolith"
  const servedBy = response.headers['X-Served-By'] || 'unknown';
  if (servedBy === 'microservice') {
    canaryRequests.add(1);
  } else if (servedBy === 'monolith') {
    monolithRequests.add(1);
  }

  // Log errors
  if (!success) {
    console.error(`Request failed: ${url} - Status: ${response.status} - Duration: ${response.timings.duration}ms`);
  }

  // Sleep between requests (1-3 seconds - realistic user behavior)
  sleep(Math.random() * 2 + 1);
}

/**
 * Setup function (runs once before test)
 */
export function setup() {
  console.log('====================================================================');
  console.log('k6 Canary Load Test');
  console.log('====================================================================');
  console.log(`Base URL: ${BASE_URL}`);
  console.log(`Duration: ${DURATION}`);
  console.log(`Virtual Users: ${VUS}`);
  console.log(`Start time: ${new Date().toISOString()}`);
  console.log('====================================================================');
  console.log('');

  // Verify backend is accessible
  const healthCheck = http.get(`${BASE_URL}/`);
  if (healthCheck.status !== 200) {
    console.error(`ERROR: Backend health check failed (status: ${healthCheck.status})`);
    console.error('Aborting test...');
    throw new Error('Backend is not accessible');
  }

  console.log('✅ Backend health check passed');
  console.log('');

  return { startTime: new Date().toISOString() };
}

/**
 * Teardown function (runs once after test)
 */
export function teardown(data) {
  console.log('');
  console.log('====================================================================');
  console.log('k6 Canary Load Test Completed');
  console.log('====================================================================');
  console.log(`Start time: ${data.startTime}`);
  console.log(`End time: ${new Date().toISOString()}`);
  console.log('====================================================================');
  console.log('');
  console.log('Next steps:');
  console.log('1. Review k6 summary output above');
  console.log('2. Check Prometheus metrics for detailed analysis');
  console.log('3. Run validation script: ./validate_canary_traffic.sh');
  console.log('');
}

/**
 * Custom summary handler
 */
export function handleSummary(data) {
  const timestamp = new Date().toISOString().replace(/[:.]/g, '-');
  const reportFile = `/tmp/k6_canary_report_${timestamp}.json`;

  console.log(`Saving detailed report to: ${reportFile}`);

  return {
    [reportFile]: JSON.stringify(data, null, 2),
    stdout: textSummary(data, { indent: '  ', enableColors: true }),
  };
}

/**
 * Text summary helper
 */
function textSummary(data, opts = {}) {
  const indent = opts.indent || '';
  const enableColors = opts.enableColors !== false;

  const lines = [];

  // Summary
  lines.push('');
  lines.push('====================================================================');
  lines.push('k6 Load Test Summary');
  lines.push('====================================================================');
  lines.push('');

  // Metrics
  const metrics = data.metrics;

  // HTTP requests
  if (metrics.http_reqs) {
    lines.push(`${indent}HTTP Requests:`);
    lines.push(`${indent}  Total: ${metrics.http_reqs.values.count}`);
    lines.push(`${indent}  Rate: ${metrics.http_reqs.values.rate.toFixed(2)} req/s`);
    lines.push('');
  }

  // Request duration
  if (metrics.http_req_duration) {
    const duration = metrics.http_req_duration.values;
    lines.push(`${indent}Request Duration:`);
    lines.push(`${indent}  Avg: ${duration.avg.toFixed(2)}ms`);
    lines.push(`${indent}  Min: ${duration.min.toFixed(2)}ms`);
    lines.push(`${indent}  Med: ${duration.med.toFixed(2)}ms`);
    lines.push(`${indent}  Max: ${duration.max.toFixed(2)}ms`);
    lines.push(`${indent}  P90: ${duration['p(90)'].toFixed(2)}ms`);
    lines.push(`${indent}  P95: ${duration['p(95)'].toFixed(2)}ms (threshold: <200ms)`);
    lines.push(`${indent}  P99: ${duration['p(99)'].toFixed(2)}ms`);
    lines.push('');
  }

  // Failed requests
  if (metrics.http_req_failed) {
    const failRate = metrics.http_req_failed.values.rate * 100;
    const status = failRate < 0.1 ? '✅' : '❌';
    lines.push(`${indent}Error Rate:`);
    lines.push(`${indent}  ${status} ${failRate.toFixed(3)}% (threshold: <0.1%)`);
    lines.push('');
  }

  // Canary distribution (if available)
  if (metrics.canary_requests && metrics.monolith_requests) {
    const canary = metrics.canary_requests.values.count;
    const monolith = metrics.monolith_requests.values.count;
    const total = canary + monolith;
    const canaryPercent = total > 0 ? (canary / total * 100) : 0;

    lines.push(`${indent}Traffic Distribution:`);
    lines.push(`${indent}  Monolith: ${monolith} (${(100 - canaryPercent).toFixed(2)}%)`);
    lines.push(`${indent}  Canary: ${canary} (${canaryPercent.toFixed(2)}%)`);
    lines.push(`${indent}  Expected: ~1% to canary`);
    lines.push('');
  }

  lines.push('====================================================================');
  lines.push('');

  return lines.join('\n');
}
