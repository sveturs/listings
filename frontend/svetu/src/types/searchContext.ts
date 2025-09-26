export type SearchContextType =
  | 'automotive'
  | 'real-estate'
  | 'electronics'
  | 'services'
  | 'fashion'
  | 'jobs'
  | 'default';

export interface SearchContextConfig {
  id: SearchContextType;
  heroTitle: string;
  heroDescription: string;
  heroIcon: string;
  bgGradient: string;
  accentColor: string;
  statsToShow?: string[];
  quickFilters?: Array<{
    id: string;
    label: string;
    filters: Record<string, any>;
  }>;
  customBanner?: {
    title: string;
    subtitle: string;
    cta?: string;
    icon?: string;
  };
  showAdvancedFilters?: boolean;
  filterComponents?: string[];
}

export const SEARCH_CONTEXTS: Record<SearchContextType, SearchContextConfig> = {
  automotive: {
    id: 'automotive',
    heroTitle: 'contexts.automotive.heroTitle',
    heroDescription: 'contexts.automotive.heroDescription',
    heroIcon: '🚗',
    bgGradient: 'from-red-600 to-orange-600',
    accentColor: 'red',
    statsToShow: ['totalCars', 'brands', 'newToday'],
    quickFilters: [
      {
        id: 'new-cars',
        label: 'search.contexts.automotive.filters.new',
        filters: { condition: 'new' },
      },
      {
        id: 'under-5000',
        label: 'search.contexts.automotive.filters.under5000',
        filters: { price_max: 5000 },
      },
      {
        id: 'diesel',
        label: 'search.contexts.automotive.filters.diesel',
        filters: { car_fuel_type: 'diesel' },
      },
      {
        id: 'automatic',
        label: 'search.contexts.automotive.filters.automatic',
        filters: { car_transmission: 'automatic' },
      },
    ],
    customBanner: {
      title: 'contexts.automotive.bannerTitle',
      subtitle: 'contexts.automotive.bannerSubtitle',
      cta: 'contexts.automotive.bannerCta',
      icon: '🔧',
    },
    showAdvancedFilters: true,
    filterComponents: ['CarFilters', 'BaseFilters'],
  },
  'real-estate': {
    id: 'real-estate',
    heroTitle: 'contexts.realEstate.heroTitle',
    heroDescription: 'contexts.realEstate.heroDescription',
    heroIcon: '🏠',
    bgGradient: 'from-blue-600 to-indigo-600',
    accentColor: 'blue',
    statsToShow: ['properties', 'newListings', 'priceRange'],
    quickFilters: [
      {
        id: 'for-rent',
        label: 'search.contexts.realEstate.filters.rent',
        filters: { listing_type: 'rent' },
      },
      {
        id: 'for-sale',
        label: 'search.contexts.realEstate.filters.sale',
        filters: { listing_type: 'sale' },
      },
      {
        id: 'apartments',
        label: 'search.contexts.realEstate.filters.apartments',
        filters: { property_type: 'apartment' },
      },
      {
        id: 'houses',
        label: 'search.contexts.realEstate.filters.houses',
        filters: { property_type: 'house' },
      },
    ],
    showAdvancedFilters: true,
    filterComponents: ['RealEstateFilters', 'LocationFilter', 'BaseFilters'],
  },
  electronics: {
    id: 'electronics',
    heroTitle: 'contexts.electronics.heroTitle',
    heroDescription: 'contexts.electronics.heroDescription',
    heroIcon: '💻',
    bgGradient: 'from-purple-600 to-pink-600',
    accentColor: 'purple',
    statsToShow: ['products', 'brands', 'deals'],
    quickFilters: [
      {
        id: 'smartphones',
        label: 'search.contexts.electronics.filters.smartphones',
        filters: { subcategory: 'smartphones' },
      },
      {
        id: 'laptops',
        label: 'search.contexts.electronics.filters.laptops',
        filters: { subcategory: 'laptops' },
      },
      {
        id: 'gaming',
        label: 'search.contexts.electronics.filters.gaming',
        filters: { subcategory: 'gaming' },
      },
      {
        id: 'warranty',
        label: 'search.contexts.electronics.filters.warranty',
        filters: { has_warranty: true },
      },
    ],
    customBanner: {
      title: 'contexts.electronics.bannerTitle',
      subtitle: 'contexts.electronics.bannerSubtitle',
      icon: '⚡',
    },
    filterComponents: ['ElectronicsFilters', 'BaseFilters'],
  },
  services: {
    id: 'services',
    heroTitle: 'contexts.services.heroTitle',
    heroDescription: 'contexts.services.heroDescription',
    heroIcon: '🛠️',
    bgGradient: 'from-green-600 to-teal-600',
    accentColor: 'green',
    statsToShow: ['providers', 'categories', 'reviews'],
    quickFilters: [
      {
        id: 'home-repair',
        label: 'search.contexts.services.filters.homeRepair',
        filters: { service_category: 'home-repair' },
      },
      {
        id: 'beauty',
        label: 'search.contexts.services.filters.beauty',
        filters: { service_category: 'beauty' },
      },
      {
        id: 'education',
        label: 'search.contexts.services.filters.education',
        filters: { service_category: 'education' },
      },
      {
        id: 'verified',
        label: 'search.contexts.services.filters.verified',
        filters: { verified_provider: true },
      },
    ],
    filterComponents: ['ServiceFilters', 'LocationFilter', 'BaseFilters'],
  },
  fashion: {
    id: 'fashion',
    heroTitle: 'contexts.fashion.heroTitle',
    heroDescription: 'contexts.fashion.heroDescription',
    heroIcon: '👗',
    bgGradient: 'from-pink-600 to-rose-600',
    accentColor: 'pink',
    statsToShow: ['items', 'brands', 'newArrivals'],
    quickFilters: [
      {
        id: 'women',
        label: 'search.contexts.fashion.filters.women',
        filters: { gender: 'women' },
      },
      {
        id: 'men',
        label: 'search.contexts.fashion.filters.men',
        filters: { gender: 'men' },
      },
      {
        id: 'shoes',
        label: 'search.contexts.fashion.filters.shoes',
        filters: { category: 'shoes' },
      },
      {
        id: 'sale',
        label: 'search.contexts.fashion.filters.sale',
        filters: { on_sale: true },
      },
    ],
    filterComponents: ['FashionFilters', 'BaseFilters'],
  },
  jobs: {
    id: 'jobs',
    heroTitle: 'contexts.jobs.heroTitle',
    heroDescription: 'contexts.jobs.heroDescription',
    heroIcon: '💼',
    bgGradient: 'from-indigo-600 to-blue-600',
    accentColor: 'indigo',
    statsToShow: ['openPositions', 'companies', 'categories'],
    quickFilters: [
      {
        id: 'full-time',
        label: 'search.contexts.jobs.filters.fullTime',
        filters: { employment_type: 'full-time' },
      },
      {
        id: 'remote',
        label: 'search.contexts.jobs.filters.remote',
        filters: { remote: true },
      },
      {
        id: 'it',
        label: 'search.contexts.jobs.filters.it',
        filters: { industry: 'it' },
      },
      {
        id: 'entry-level',
        label: 'search.contexts.jobs.filters.entryLevel',
        filters: { experience_level: 'entry' },
      },
    ],
    filterComponents: ['JobFilters', 'LocationFilter', 'BaseFilters'],
  },
  default: {
    id: 'default',
    heroTitle: 'search.title',
    heroDescription: 'search.description',
    heroIcon: '🔍',
    bgGradient: 'from-primary to-primary-focus',
    accentColor: 'primary',
    filterComponents: ['BaseFilters'],
  },
};

export function getSearchContext(contextId?: string): SearchContextConfig {
  if (!contextId) return SEARCH_CONTEXTS.default;

  const context = SEARCH_CONTEXTS[contextId as SearchContextType];
  return context || SEARCH_CONTEXTS.default;
}
