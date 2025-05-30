export interface MarketplaceImage {
  id: number;
  public_url: string;
  is_main: boolean;
}

export interface MarketplaceItem {
  id: number;
  title: string;
  description?: string;
  price?: number;
  location?: string;
  city?: string;
  country?: string;
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
  };
  condition?: string;
  status?: string;
  views_count?: number;
  old_price?: number;
  has_discount?: boolean;
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

export interface MarketplaceSearchParams {
  sort_by?: 'date_desc' | 'date_asc' | 'price_asc' | 'price_desc';
  page?: number;
  size?: number;
  category?: string;
  minPrice?: number;
  maxPrice?: number;
  search?: string;
}
