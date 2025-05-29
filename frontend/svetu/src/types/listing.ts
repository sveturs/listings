export interface Listing {
  id: string | number;
  title: string;
  description: string;
  price: number;
  old_price?: number;
  currency?: string;
  location?: string;
  city?: string;
  country?: string;
  address?: string;
  latitude?: number;
  longitude?: number;
  coordinates?: {
    latitude: number;
    longitude: number;
  };
  category?: Category;
  category_id?: string | number;
  subcategory?: Subcategory;
  attributes?: AttributeValue[];
  images?: ListingImage[];
  status?: ListingStatus;
  condition?: '' | 'new' | 'used' | string;
  userId?: string;
  user_id?: string | number;
  user?: User;
  storeId?: string;
  store?: Store;
  storefront_id?: number | string;
  storefront_name?: string;
  storefrontName?: string;
  isPartOfStorefront?: boolean;
  storefrontItemCount?: number;
  viewCount?: number;
  views_count?: number;
  favoriteCount?: number;
  isFavorite?: boolean;
  rating?: number;
  reviews_count?: number;
  average_rating?: number;
  review_count?: number;
  show_on_map?: boolean;
  has_discount?: boolean;
  metadata?: ListingMetadata;
  original_language?: string;
  translations?: {
    [language: string]: {
      [field: string]: string;
    };
  };
  created_at?: string;
  updated_at?: string;
  createdAt?: string;
  updatedAt?: string;
  [key: string]: unknown;
}

export interface ListingImage {
  id?: string | number;
  url?: string;
  thumbnailUrl?: string;
  order?: number;
  file_path?: string;
  is_main?: boolean;
  storage_type?: string;
  public_url?: string;
  [key: string]: unknown;
}

export interface Category {
  id: string;
  name: string;
  slug: string;
  icon?: string;
  parentId?: string;
}

export interface Subcategory {
  id: string;
  name: string;
  slug: string;
  categoryId: string;
}

export interface Attribute {
  id?: string;
  attribute_id?: number | string;
  name?: string;
  attribute_name?: string;
  display_name?: string;
  display_value?: string;
  numeric_value?: number;
  type?: AttributeType;
  required?: boolean;
  options?: string[];
  translations?: {
    [language: string]: string;
  };
  [key: string]: unknown;
}

export interface AttributeValue {
  attributeId?: string;
  attribute_id?: number | string;
  attribute?: Attribute;
  value?: string | number | boolean;
  display_value?: string;
  numeric_value?: number;
}

export enum AttributeType {
  TEXT = 'text',
  NUMBER = 'number',
  SELECT = 'select',
  MULTISELECT = 'multiselect',
  BOOLEAN = 'boolean',
  DATE = 'date',
}

export enum ListingStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  SOLD = 'sold',
  DRAFT = 'draft',
}

export interface User {
  id: string;
  name: string;
  avatar?: string;
  rating?: number;
}

export interface Store {
  id: string;
  name: string;
  logo?: string;
  rating?: number;
}

export interface ListingFilters {
  categoryId?: string;
  category_id?: string | number;
  subcategoryId?: string;
  minPrice?: number;
  min_price?: string | number;
  maxPrice?: number;
  max_price?: string | number;
  location?: string;
  city?: string;
  country?: string;
  search?: string;
  query?: string;
  attributes?: Record<string, unknown>;
  attributeFilters?: Record<string, unknown>;
  status?: ListingStatus;
  condition?: '' | 'new' | 'used' | string;
  userId?: string;
  user_id?: string;
  storeId?: string;
  storefront_id?: string | number;
  distance?: string;
  latitude?: number | null;
  longitude?: number | null;
  sort_by?: string;
  view_mode?: string;
  size?: number;
  [key: string]: unknown;
}

export interface ListingSort {
  field: 'price' | 'createdAt' | 'viewCount' | 'favoriteCount';
  order: 'asc' | 'desc';
}

export interface DiscountInfo {
  discount_percent?: number;
  previous_price?: number;
  percent?: number;
  oldPrice?: number;
  has_price_history?: boolean;
  [key: string]: unknown;
}

export interface ListingMetadata {
  discount?: DiscountInfo;
  promotions?: Record<string, unknown>;
  [key: string]: unknown;
}

export interface PaginatedResponse<T> {
  items: T[];
  total: number;
  page: number;
  pageSize: number;
  totalPages: number;
}