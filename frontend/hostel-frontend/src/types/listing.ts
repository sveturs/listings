// frontend/hostel-frontend/src/types/listing.ts
// Общий интерфейс для всех объявлений в приложении

export interface ListingImage {
  id?: string | number;
  file_path?: string;
  is_main?: boolean;
  storage_type?: string;
  public_url?: string;
  [key: string]: any;
}

export interface Attribute {
  attribute_id: number | string;
  attribute_name: string;
  display_name: string;
  display_value: string;
  numeric_value?: number;
  translations?: {
    [language: string]: string;
  };
  [key: string]: any;
}

export interface DiscountInfo {
  discount_percent?: number;
  previous_price?: number;
  percent?: number;
  oldPrice?: number;
  [key: string]: any;
}

export interface ListingMetadata {
  discount?: DiscountInfo;
  promotions?: Record<string, any>;
  [key: string]: any;
}

// Универсальный интерфейс Listing, который включает все возможные поля
export interface Listing {
  id: string | number;
  title: string;
  description?: string;
  price: number;
  old_price?: number;
  category_id?: string | number;
  user_id?: string | number;
  city?: string;
  country?: string;
  location?: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  coordinates?: {
    latitude: number;
    longitude: number;
  };
  status?: 'active' | 'inactive' | string;
  condition?: '' | 'new' | 'used' | string;
  created_at?: string;
  updated_at?: string;
  images?: Array<ListingImage | string>;
  attributes?: Attribute[];
  metadata?: ListingMetadata;
  has_discount?: boolean;
  original_language?: string;
  translations?: {
    [language: string]: {
      [field: string]: string;
    };
  };
  storefront_id?: number | string;
  storefront_name?: string;
  isPartOfStorefront?: boolean;
  isUniqueLocation?: boolean;
  storefrontName?: string;
  storefrontItemCount?: number;
  show_on_map?: boolean;
  rating?: number;
  reviews_count?: number;
  average_rating?: number;
  review_count?: number;
  views_count?: number;
  [key: string]: any; // Для динамического доступа к свойствам
}

// Универсальные интерфейсы для фильтров
export interface FilterState {
  query: string;
  category_id: string;
  min_price: string;
  max_price: string;
  city: string;
  country: string;
  condition: string;
  sort_by: string;
  distance: string;
  latitude?: number | null;
  longitude?: number | null;
  attributeFilters?: Record<string, any>;
  view_mode?: string;
  size?: number;
  [key: string]: any;
}

export interface FilterOptions {
  query?: string;
  category_id?: number | string;
  min_price?: string | number;
  max_price?: string | number;
  distance?: string;
  condition?: '' | 'new' | 'used' | string;
  latitude?: number;
  longitude?: number;
  [key: string]: any;
}