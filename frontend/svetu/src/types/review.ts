export interface Review {
  id: number;
  user_id: number;
  entity_type: 'listing' | 'user' | 'storefront';
  entity_id: number;
  entity_origin_type?: string;
  entity_origin_id?: number;
  rating: number;
  comment?: string;
  pros?: string;
  cons?: string;
  photos?: string[];
  likes_count: number;
  is_verified_purchase: boolean;
  status: 'pending' | 'published' | 'deleted' | 'spam';
  created_at: string;
  updated_at: string;
  helpful_votes: number;
  not_helpful_votes: number;
  seller_confirmed: boolean;
  has_active_dispute: boolean;
  original_language: string;
  translations?: Record<string, Record<string, string>>;

  // Дополнительные поля для отображения
  user?: {
    id: number;
    name: string;
    avatar?: string;
  };
  responses?: ReviewResponse[];
  votes_count?: {
    helpful: number;
    not_helpful: number;
  };
  current_user_vote?: 'helpful' | 'not_helpful' | null;
}

export interface ReviewResponse {
  id: number;
  review_id: number;
  user_id: number;
  response: string;
  created_at: string;
  updated_at: string;
  user?: {
    id: number;
    name: string;
    avatar?: string;
  };
}

export interface CreateReviewRequest {
  entity_type: 'listing' | 'user' | 'storefront';
  entity_id: number;
  rating: number;
  storefront_id?: number;
  comment: string;
  pros?: string;
  cons?: string;
  photos?: string[];
  original_language: string;
}

export interface ReviewsFilter {
  entity_type?: string;
  entity_id?: number;
  user_id?: number;
  min_rating?: number;
  max_rating?: number;
  status?: string;
  sort_by?: 'rating' | 'date' | 'likes';
  sort_order?: 'asc' | 'desc';
  page?: number;
  limit?: number;
}

export interface ReviewStats {
  total_reviews: number;
  average_rating: number;
  verified_reviews: number;
  rating_distribution: Record<number, number>;
  photo_reviews: number;
}

export interface AggregatedRating {
  entity_id: number;
  entity_type: string;
  total_reviews: number;
  average: number;
  recent_rating: number;
  recent_reviews: number;
  recent_trend: 'up' | 'down' | 'stable';
  verified_percentage: number;
  last_review_at: string;
  distribution: Record<string, number>;
  breakdown: {
    listing_reviews?: number;
    storefront_reviews?: number;
    direct_reviews?: number;
  };
}

export interface CanReviewResponse {
  can_review: boolean;
  reason?: string;
  chat_activity?: {
    total_messages: number;
    buyer_messages: number;
    seller_messages: number;
    last_message_at: string;
    first_message_at: string;
  };
}

export interface ReviewConfirmation {
  review_id: number;
  seller_id: number;
  confirmed: boolean;
  confirmed_at?: string;
  notes?: string;
}

export interface ReviewDispute {
  id: number;
  review_id: number;
  user_id: number;
  reason: string;
  status: 'pending' | 'resolved' | 'rejected';
  admin_notes?: string;
  created_at: string;
  resolved_at?: string;
  resolved_by?: number;
}
