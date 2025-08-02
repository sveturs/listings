export interface MarketplaceImage {
  id: number;
  public_url: string;
  is_main: boolean;
}

export interface MarketplaceItem {
  id: number;
  user_id: number;
  title: string;
  description?: string;
  price?: number;
  location?: string;
  city?: string;
  country?: string;
  location_privacy?: 'exact' | 'street' | 'district' | 'city';
  // Переводы адресов (могут приходить из backend)
  translations?: {
    location?: {
      sr?: string;
      en?: string;
      ru?: string;
    };
    city?: {
      sr?: string;
      en?: string;
      ru?: string;
    };
    country?: {
      sr?: string;
      en?: string;
      ru?: string;
    };
  };
  images?: MarketplaceImage[];
  created_at: string;
  updated_at?: string;
  user?: {
    id: number;
    name: string;
    email: string;
    picture_url?: string;
  };
  category?: {
    id: number;
    name: string;
    slug: string;
    translations?: Record<string, string>;
  };
  condition?: string;
  status?: string;
  metadata?: {
    seo?: {
      keywords?: string;
      slug?: string;
    };
    [key: string]: any;
  };
  views_count?: number;
  old_price?: number;
  has_discount?: boolean;
  // Добавляем для поддержки unified search
  product_type?: 'marketplace' | 'storefront';
  storefront_id?: number;
  storefront_slug?: string;
  storefront?: {
    id: number;
    name: string;
    slug: string;
  };
  // Атрибуты объявления
  attributes?: Array<{
    attribute_id?: number;
    attribute_name?: string;
    name?: string;
    value?: string | number | boolean;
    text_value?: string;
    numeric_value?: number;
    boolean_value?: boolean;
  }>;
}

export interface MarketplaceSearchResponse {
  data: MarketplaceItem[];
  meta: {
    total: number;
    page: number;
    size: number;
    has_more: boolean;
    total_pages: number;
  };
}

export interface ListingFormData {
  title: string;
  description: string;
  price: number;
  currency: 'RSD' | 'EUR' | 'HRK' | 'MKD';
  condition: 'new' | 'used' | 'refurbished';
  category?: {
    id: number;
    name: string;
    slug: string;
    translations?: Record<string, string>;
  };
  trust?: {
    phoneVerified: boolean;
    preferredMeetingType: 'personal' | 'pickup' | 'delivery';
    meetingLocations: string[];
    availableHours: string;
    localReputation: number;
  };
  payment?: {
    methods: string[];
    codEnabled: boolean;
    codPrice: number;
    personalMeeting: boolean;
    deliveryOptions: string[];
    negotiablePrice: boolean;
    bundleDeals: boolean;
  };
  localization?: {
    script: 'cyrillic' | 'latin' | 'mixed';
    language: 'sr' | 'hr' | 'sl' | 'mk';
    traditionalUnits: boolean;
    regionalPhrases: string[];
  };
  pijaca?: {
    vendorStallStyle: string;
    neighborhoodGroup?: string;
    regularCustomers: boolean;
    traditionalStyle: boolean;
  };
  location?: {
    latitude: number;
    longitude: number;
    address: string;
    city: string;
    region: string;
    country: string;
  };
}

export interface MarketplaceSearchParams {
  sort_by?: 'date_desc' | 'date_asc' | 'price_asc' | 'price_desc';
  page?: number;
  size?: number;
  category?: string;
  minPrice?: number;
  maxPrice?: number;
  search?: string;
}

export interface ListingCardData {
  id: number;
  title: string;
  price: number;
  old_price?: number;
  location?: string;
  main_image?: string;
  has_discount?: boolean;
  discount_percentage?: number;
  created_at: string;
  views_count?: number;
  favorites_count?: number;
  is_favorite?: boolean;
}
