import { ABTest } from '@/components/abtest/ABTestProvider';

// Пример конфигурации A/B тестов
export const abTests: ABTest[] = [
  {
    id: 'listing-card-layout',
    name: 'Listing Card Layout Test',
    description: 'Testing different layouts for listing cards',
    status: 'running',
    allocation: 100, // 100% трафика участвует в тесте
    variants: [
      {
        id: 'control',
        name: 'Current Layout',
        weight: 50,
        isControl: true,
        config: {
          showStats: false,
          showBadges: true,
        },
      },
      {
        id: 'with-stats',
        name: 'With Statistics',
        weight: 50,
        config: {
          showStats: true,
          showBadges: true,
        },
      },
    ],
    targeting: {
      urls: ['/cars', '/c2c', '/real-estate'],
    },
    metrics: [
      {
        name: 'Click-through Rate',
        type: 'conversion',
        goal: 5,
      },
      {
        name: 'Add to Favorites',
        type: 'engagement',
      },
    ],
  },
  {
    id: 'image-loading-strategy',
    name: 'Image Loading Strategy',
    description:
      'Testing lazy loading vs eager loading for above-the-fold images',
    status: 'running',
    allocation: 50, // Только 50% трафика
    variants: [
      {
        id: 'lazy',
        name: 'Lazy Loading',
        weight: 50,
        config: {
          priority: false,
          loading: 'lazy',
        },
      },
      {
        id: 'eager',
        name: 'Eager Loading',
        weight: 50,
        config: {
          priority: true,
          loading: 'eager',
        },
      },
    ],
    metrics: [
      {
        name: 'Page Load Time',
        type: 'custom',
      },
      {
        name: 'First Contentful Paint',
        type: 'custom',
      },
    ],
  },
  {
    id: 'cta-button-color',
    name: 'CTA Button Color Test',
    description: 'Testing different colors for call-to-action buttons',
    status: 'paused',
    variants: [
      {
        id: 'primary',
        name: 'Primary Color',
        weight: 33,
        isControl: true,
        config: {
          buttonClass: 'btn-primary',
        },
      },
      {
        id: 'accent',
        name: 'Accent Color',
        weight: 33,
        config: {
          buttonClass: 'btn-accent',
        },
      },
      {
        id: 'success',
        name: 'Success Color',
        weight: 34,
        config: {
          buttonClass: 'btn-success',
        },
      },
    ],
    targeting: {
      devices: ['mobile', 'tablet'],
    },
    metrics: [
      {
        name: 'Button Click Rate',
        type: 'conversion',
        goal: 10,
      },
    ],
  },
  {
    id: 'search-algorithm',
    name: 'Search Algorithm Test',
    description: 'Testing different search ranking algorithms',
    status: 'draft',
    variants: [
      {
        id: 'relevance',
        name: 'Relevance-based',
        weight: 50,
        isControl: true,
        config: {
          algorithm: 'relevance',
          boostFactor: 1.0,
        },
      },
      {
        id: 'popularity',
        name: 'Popularity-based',
        weight: 50,
        config: {
          algorithm: 'popularity',
          boostFactor: 1.5,
        },
      },
    ],
    targeting: {
      urls: ['/search', '/c2c/search'],
    },
    metrics: [
      {
        name: 'Search Result Clicks',
        type: 'engagement',
      },
      {
        name: 'Zero Results Rate',
        type: 'custom',
      },
      {
        name: 'Search-to-Contact',
        type: 'conversion',
        goal: 3,
      },
    ],
  },
];

// Функция для получения активных тестов
export const getActiveTests = (): ABTest[] => {
  return abTests.filter((test) => test.status === 'running');
};

// Функция для получения теста по ID
export const getTestById = (id: string): ABTest | undefined => {
  return abTests.find((test) => test.id === id);
};
