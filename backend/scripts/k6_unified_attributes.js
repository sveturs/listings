import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate } from 'k6/metrics';

// Custom metrics
const errorRate = new Rate('errors');
const dualWriteSuccessRate = new Rate('dual_write_success');
const cacheHitRate = new Rate('cache_hits');

export const options = {
  stages: [
    { duration: '30s', target: 10 },   // Warm up
    { duration: '1m', target: 50 },    // Ramp up to 50 users
    { duration: '2m', target: 100 },   // Stay at 100 users
    { duration: '1m', target: 200 },   // Spike to 200 users
    { duration: '2m', target: 100 },   // Back to normal
    { duration: '30s', target: 0 },    // Ramp down
  ],
  thresholds: {
    'http_req_duration': ['p(95)<100', 'p(99)<200'], // 95% of requests under 100ms
    'http_req_failed': ['rate<0.01'],                 // Error rate under 1%
    'errors': ['rate<0.05'],                          // Custom error rate under 5%
    'dual_write_success': ['rate>0.95'],              // Dual write success > 95%
    'cache_hits': ['rate>0.8'],                       // Cache hit rate > 80%
  },
};

const BASE_URL = 'http://localhost:3000';
const CATEGORIES = [1, 2, 3, 4, 5]; // Test category IDs
const LISTINGS = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]; // Test listing IDs

// Test data for attributes
const ATTRIBUTE_VALUES = {
  text: ['Red', 'Blue', 'Green', 'Yellow', 'Black', 'White'],
  number: [10, 20, 30, 40, 50, 100, 200],
  boolean: [true, false],
  select: ['option1', 'option2', 'option3'],
};

export default function() {
  // Test 1: Get category attributes (should hit cache after first request)
  const categoryId = CATEGORIES[Math.floor(Math.random() * CATEGORIES.length)];
  const getAttributesRes = http.get(`${BASE_URL}/api/v2/categories/${categoryId}/attributes`, {
    headers: { 'Content-Type': 'application/json' },
    tags: { name: 'GetCategoryAttributes' },
  });
  
  check(getAttributesRes, {
    'get attributes status is 200': (r) => r.status === 200,
    'get attributes response time < 50ms': (r) => r.timings.duration < 50,
    'attributes response has data': (r) => {
      const body = JSON.parse(r.body);
      return body.data && Array.isArray(body.data);
    },
  });
  
  // Check if response was from cache (header or response time indicates cache hit)
  const fromCache = getAttributesRes.timings.duration < 10;
  cacheHitRate.add(fromCache);
  
  errorRate.add(getAttributesRes.status !== 200);
  
  sleep(0.5);
  
  // Test 2: Save attribute values (tests dual-write mechanism)
  const listingId = LISTINGS[Math.floor(Math.random() * LISTINGS.length)];
  const attributes = JSON.parse(getAttributesRes.body).data || [];
  
  if (attributes.length > 0) {
    const attributeValues = {};
    
    // Generate random values for each attribute
    attributes.slice(0, 3).forEach(attr => {
      switch(attr.type) {
        case 'text':
          attributeValues[attr.id] = ATTRIBUTE_VALUES.text[Math.floor(Math.random() * ATTRIBUTE_VALUES.text.length)];
          break;
        case 'number':
          attributeValues[attr.id] = ATTRIBUTE_VALUES.number[Math.floor(Math.random() * ATTRIBUTE_VALUES.number.length)];
          break;
        case 'boolean':
          attributeValues[attr.id] = ATTRIBUTE_VALUES.boolean[Math.floor(Math.random() * ATTRIBUTE_VALUES.boolean.length)];
          break;
        case 'select':
          attributeValues[attr.id] = ATTRIBUTE_VALUES.select[Math.floor(Math.random() * ATTRIBUTE_VALUES.select.length)];
          break;
        default:
          attributeValues[attr.id] = 'default value';
      }
    });
    
    const saveValuesRes = http.post(
      `${BASE_URL}/api/v2/listings/${listingId}/attributes`,
      JSON.stringify({ attributes: attributeValues }),
      {
        headers: { 
          'Content-Type': 'application/json',
          'Authorization': 'Bearer test-token',
        },
        tags: { name: 'SaveAttributeValues' },
      }
    );
    
    check(saveValuesRes, {
      'save values status is 200': (r) => r.status === 200,
      'save values response time < 100ms': (r) => r.timings.duration < 100,
      'dual write successful': (r) => {
        if (r.status !== 200) return false;
        const body = JSON.parse(r.body);
        return body.dual_write_status === 'success';
      },
    });
    
    dualWriteSuccessRate.add(saveValuesRes.status === 200);
    errorRate.add(saveValuesRes.status !== 200);
  }
  
  sleep(0.5);
  
  // Test 3: Search with attributes (tests filter performance)
  const searchAttribute = attributes.length > 0 ? attributes[0] : null;
  if (searchAttribute) {
    const searchParams = {
      category_id: categoryId,
      attributes: JSON.stringify({
        [searchAttribute.id]: searchAttribute.type === 'boolean' ? true : 'test',
      }),
    };
    
    const searchRes = http.get(`${BASE_URL}/api/v2/marketplace/search?${new URLSearchParams(searchParams)}`, {
      headers: { 'Content-Type': 'application/json' },
      tags: { name: 'SearchWithAttributes' },
    });
    
    check(searchRes, {
      'search status is 200': (r) => r.status === 200,
      'search response time < 200ms': (r) => r.timings.duration < 200,
      'search returns results': (r) => {
        const body = JSON.parse(r.body);
        return body.data && Array.isArray(body.data.listings);
      },
    });
    
    errorRate.add(searchRes.status !== 200);
  }
  
  sleep(1);
  
  // Test 4: Concurrent attribute updates (tests race conditions)
  if (Math.random() < 0.1) { // 10% of VUs do this test
    const batch = [];
    for (let i = 0; i < 5; i++) {
      batch.push([
        'POST',
        `${BASE_URL}/api/v2/listings/${listingId}/attributes`,
        JSON.stringify({ 
          attributes: { 
            [attributes[0]?.id || 1]: `concurrent-${i}` 
          } 
        }),
        { headers: { 'Content-Type': 'application/json' } },
      ]);
    }
    
    const batchRes = http.batch(batch);
    const allSuccessful = batchRes.every(r => r.status === 200);
    
    check(batchRes, {
      'concurrent updates handled': () => allSuccessful,
    });
    
    errorRate.add(!allSuccessful);
  }
  
  sleep(0.5);
}

// Lifecycle hook to check system health after test
export function teardown(data) {
  const healthCheck = http.get(`${BASE_URL}/health/ready`);
  
  if (healthCheck.status !== 200) {
    console.error('System unhealthy after load test!');
  }
  
  // Get final metrics
  const metricsRes = http.get(`${BASE_URL}/metrics`);
  if (metricsRes.status === 200) {
    const metrics = metricsRes.body;
    
    // Parse Prometheus metrics
    const dualWriteErrors = metrics.match(/dual_write_errors_total (\d+)/);
    const cacheHits = metrics.match(/cache_hits_total (\d+)/);
    const cacheMisses = metrics.match(/cache_misses_total (\d+)/);
    
    console.log('Final metrics:');
    if (dualWriteErrors) console.log(`Dual write errors: ${dualWriteErrors[1]}`);
    if (cacheHits && cacheMisses) {
      const hitRate = parseInt(cacheHits[1]) / (parseInt(cacheHits[1]) + parseInt(cacheMisses[1]));
      console.log(`Cache hit rate: ${(hitRate * 100).toFixed(2)}%`);
    }
  }
}