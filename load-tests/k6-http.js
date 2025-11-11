/**
 * K6 HTTP Load Test for Listings Microservice
 *
 * Tests HTTP endpoints under various load scenarios:
 * - Warmup: 10 RPS for 30s
 * - Ramp-up: 10 to 100 RPS over 1m
 * - Sustained: 100 RPS for 2m
 * - Peak: 200 RPS for 1m
 * - Cool-down: 200 to 0 RPS over 30s
 *
 * Success Criteria:
 * - p95 latency < 100ms
 * - Error rate < 1%
 * - 100 RPS without degradation
 */

import http from 'k6/http';
import { check, group, sleep } from 'k6';
import { Counter, Rate, Trend } from 'k6/metrics';

// ============================================================================
// Configuration
// ============================================================================

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8086';
const TEST_STOREFRONT_ID = __ENV.STOREFRONT_ID || '1';
const TEST_PRODUCT_ID = __ENV.PRODUCT_ID || '1';
const TEST_CATEGORY_ID = __ENV.CATEGORY_ID || '1';

// Custom metrics
const errorRate = new Rate('errors');
const healthCheckDuration = new Trend('health_check_duration');
const storefrontsDuration = new Trend('storefronts_duration');
const listingsDuration = new Trend('listings_duration');
const successCounter = new Counter('success_count');
const errorCounter = new Counter('error_count');

// ============================================================================
// Load Test Stages
// ============================================================================

export const options = {
  stages: [
    // Warmup Phase
    { duration: '30s', target: 10 },  // Warmup: 10 RPS for 30s

    // Ramp-up Phase
    { duration: '1m', target: 100 },  // Ramp-up: 10 to 100 RPS over 1m

    // Sustained Load Phase
    { duration: '2m', target: 100 },  // Sustained: 100 RPS for 2m

    // Peak Load Phase
    { duration: '1m', target: 200 },  // Peak: 200 RPS for 1m

    // Cool-down Phase
    { duration: '30s', target: 0 },   // Cool-down: 200 to 0 RPS over 30s
  ],

  thresholds: {
    // Success criteria
    'http_req_duration': ['p(95)<100'],     // p95 latency < 100ms
    'errors': ['rate<0.01'],                 // Error rate < 1%
    'http_req_failed': ['rate<0.01'],        // Failed requests < 1%

    // Additional quality gates
    'http_req_duration{group:::Health Check}': ['p(95)<50'],
    'http_req_duration{group:::Storefronts API}': ['p(95)<150'],
    'http_req_duration{group:::Listings API}': ['p(95)<150'],
  },

  // HTTP configuration
  noConnectionReuse: false,
  userAgent: 'K6-LoadTest/1.0',
};

// ============================================================================
// Test Scenarios
// ============================================================================

export default function() {
  // Distribute load across different endpoints
  const scenario = Math.random();

  if (scenario < 0.3) {
    testHealthEndpoint();
  } else if (scenario < 0.65) {
    testStorefrontsAPI();
  } else {
    testListingsAPI();
  }

  // Think time: simulate user reading/processing response
  sleep(Math.random() * 2 + 0.5); // 0.5-2.5s
}

// ============================================================================
// Health Check Tests
// ============================================================================

function testHealthEndpoint() {
  group('Health Check', () => {
    const startTime = new Date();
    const response = http.get(`${BASE_URL}/health`, {
      tags: { name: 'HealthCheck' },
    });

    healthCheckDuration.add(new Date() - startTime);

    const success = check(response, {
      'health check status is 200': (r) => r.status === 200,
      'health check has status field': (r) => {
        try {
          const body = JSON.parse(r.body);
          return body.status !== undefined;
        } catch (e) {
          return false;
        }
      },
      'response time < 50ms': (r) => r.timings.duration < 50,
    });

    errorRate.add(!success);
    if (success) {
      successCounter.add(1);
    } else {
      errorCounter.add(1);
      console.error(`Health check failed: ${response.status} - ${response.body}`);
    }
  });
}

// ============================================================================
// Storefronts API Tests
// ============================================================================

function testStorefrontsAPI() {
  group('Storefronts API', () => {
    const startTime = new Date();

    // Test 1: List storefronts with pagination
    const listResponse = http.get(`${BASE_URL}/api/v1/storefronts?limit=10&offset=0`, {
      tags: { name: 'ListStorefronts' },
    });

    storefrontsDuration.add(new Date() - startTime);

    const listSuccess = check(listResponse, {
      'storefronts list status is 200': (r) => r.status === 200,
      'storefronts list has data array': (r) => {
        try {
          const body = JSON.parse(r.body);
          return Array.isArray(body.storefronts);
        } catch (e) {
          return false;
        }
      },
      'response time < 150ms': (r) => r.timings.duration < 150,
    });

    errorRate.add(!listSuccess);
    if (listSuccess) {
      successCounter.add(1);

      // Test 2: Get specific storefront (if we got results)
      try {
        const storefronts = JSON.parse(listResponse.body).storefronts;
        if (storefronts && storefronts.length > 0) {
          const storefrontId = storefronts[0].id || TEST_STOREFRONT_ID;

          const getResponse = http.get(`${BASE_URL}/api/v1/storefronts/${storefrontId}`, {
            tags: { name: 'GetStorefront' },
          });

          const getSuccess = check(getResponse, {
            'get storefront status is 200': (r) => r.status === 200,
            'get storefront has id field': (r) => {
              try {
                const body = JSON.parse(r.body);
                return body.storefront && body.storefront.id !== undefined;
              } catch (e) {
                return false;
              }
            },
          });

          errorRate.add(!getSuccess);
          if (getSuccess) {
            successCounter.add(1);
          } else {
            errorCounter.add(1);
          }
        }
      } catch (e) {
        console.error(`Failed to parse storefronts response: ${e.message}`);
        errorCounter.add(1);
      }
    } else {
      errorCounter.add(1);
      console.error(`Storefronts list failed: ${listResponse.status} - ${listResponse.body}`);
    }
  });
}

// ============================================================================
// Listings API Tests
// ============================================================================

function testListingsAPI() {
  group('Listings API', () => {
    const startTime = new Date();

    // Test 1: List listings with filters
    const listResponse = http.get(
      `${BASE_URL}/api/v1/listings?category_id=${TEST_CATEGORY_ID}&limit=20&offset=0`,
      {
        tags: { name: 'ListListings' },
      }
    );

    listingsDuration.add(new Date() - startTime);

    const listSuccess = check(listResponse, {
      'listings list status is 200': (r) => r.status === 200,
      'listings list has data array': (r) => {
        try {
          const body = JSON.parse(r.body);
          return Array.isArray(body.listings);
        } catch (e) {
          return false;
        }
      },
      'response time < 150ms': (r) => r.timings.duration < 150,
    });

    errorRate.add(!listSuccess);
    if (listSuccess) {
      successCounter.add(1);

      // Test 2: Get specific listing (if we got results)
      try {
        const listings = JSON.parse(listResponse.body).listings;
        if (listings && listings.length > 0) {
          const listingId = listings[0].id || TEST_PRODUCT_ID;

          const getResponse = http.get(`${BASE_URL}/api/v1/listings/${listingId}`, {
            tags: { name: 'GetListing' },
          });

          const getSuccess = check(getResponse, {
            'get listing status is 200': (r) => r.status === 200,
            'get listing has id field': (r) => {
              try {
                const body = JSON.parse(r.body);
                return body.listing && body.listing.id !== undefined;
              } catch (e) {
                return false;
              }
            },
          });

          errorRate.add(!getSuccess);
          if (getSuccess) {
            successCounter.add(1);
          } else {
            errorCounter.add(1);
          }
        }
      } catch (e) {
        console.error(`Failed to parse listings response: ${e.message}`);
        errorCounter.add(1);
      }
    } else {
      errorCounter.add(1);
      console.error(`Listings list failed: ${listResponse.status} - ${listResponse.body}`);
    }
  });
}

// ============================================================================
// Test Setup and Teardown
// ============================================================================

export function setup() {
  console.log('========================================');
  console.log('Starting HTTP Load Test');
  console.log('========================================');
  console.log(`Base URL: ${BASE_URL}`);
  console.log(`Test Duration: ~5 minutes`);
  console.log(`Max RPS: 200`);
  console.log('');

  // Verify service is available
  const healthCheck = http.get(`${BASE_URL}/health`);
  if (healthCheck.status !== 200) {
    throw new Error(`Service not available: ${healthCheck.status} - ${healthCheck.body}`);
  }

  console.log('Service health check: OK');
  console.log('Starting load test...');
  console.log('');
}

export function teardown(data) {
  console.log('');
  console.log('========================================');
  console.log('Load Test Completed');
  console.log('========================================');
}
