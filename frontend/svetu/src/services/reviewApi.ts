import config from '@/config';
import { AuthService } from '@/services/auth';
import type {
  Review,
  ReviewsFilter,
  CreateReviewRequest,
  ReviewStats,
  AggregatedRating,
  CanReviewResponse,
  ReviewConfirmation,
  ReviewDispute,
} from '@/types/review';

const API_BASE = config.getApiUrl() + '/api/v1';

export const reviewApi = {
  // Get reviews with filters
  async getReviews(filters: ReviewsFilter): Promise<{
    reviews: Review[];
    totalPages: number;
    currentPage: number;
    total: number;
  }> {
    const params = new URLSearchParams();
    Object.entries(filters).forEach(([key, value]) => {
      if (value !== undefined) {
        params.append(key, String(value));
      }
    });

    const response = await fetch(`${API_BASE}/reviews?${params}`, {
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to fetch reviews');
    }

    const data = await response.json();
    console.log('reviewApi.getReviews: raw response', data);

    // Backend возвращает структуру: { data: { data: [...], meta: {...} } }
    const reviewsData = data.data?.data || data.data?.reviews || [];
    const meta = data.data?.meta || {};

    return {
      reviews: reviewsData,
      totalPages:
        meta.total_pages || Math.ceil((meta.total || 0) / (meta.limit || 10)),
      currentPage: meta.page || meta.current_page || 1,
      total: meta.total || 0,
    };
  },

  // Get review by ID
  async getReviewById(id: number): Promise<Review> {
    const response = await fetch(`${API_BASE}/reviews/${id}`, {
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to fetch review');
    }

    const data = await response.json();
    return data.data;
  },

  // Get review statistics
  async getReviewStats(
    entityType: string,
    entityId: number
  ): Promise<ReviewStats> {
    const response = await fetch(
      `${API_BASE}/entity/${entityType}/${entityId}/stats`,
      {
        credentials: 'include',
      }
    );

    if (!response.ok) {
      throw new Error('Failed to fetch review stats');
    }

    const data = await response.json();
    return data.data;
  },

  // Get aggregated rating
  async getAggregatedRating(
    entityType: string,
    entityId: number
  ): Promise<AggregatedRating> {
    const endpoint =
      entityType === 'user'
        ? `${API_BASE}/users/${entityId}/aggregated-rating`
        : `${API_BASE}/storefronts/${entityId}/aggregated-rating`;

    const response = await fetch(endpoint, {
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to fetch aggregated rating');
    }

    const data = await response.json();
    return data.data;
  },

  // Check if user can review
  async canReview(
    entityType: string,
    entityId: number
  ): Promise<CanReviewResponse> {
    const response = await fetch(
      `${API_BASE}/reviews/can-review/${entityType}/${entityId}`,
      {
        credentials: 'include',
      }
    );

    if (!response.ok) {
      // Return default response instead of throwing error for 401
      if (response.status === 401) {
        return { can_review: false, reason: 'unauthorized' };
      }
      throw new Error('Failed to check review permission');
    }

    const data = await response.json();
    return data.data;
  },

  // Create a draft review (step 1)
  async createDraftReview(reviewData: CreateReviewRequest): Promise<Review> {
    const csrfToken = await AuthService.getCsrfToken();

    const response = await fetch(`${API_BASE}/reviews/draft`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
      body: JSON.stringify(reviewData),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Failed to create draft review');
    }

    const data = await response.json();
    return data.data;
  },

  // Publish a draft review (step 2b)
  async publishReview(reviewId: number): Promise<Review> {
    const csrfToken = await AuthService.getCsrfToken();

    const response = await fetch(`${API_BASE}/reviews/${reviewId}/publish`, {
      method: 'POST',
      headers: {
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Failed to publish review');
    }

    const data = await response.json();
    return data.data;
  },

  // Legacy: Create a new review (single step)
  async createReview(reviewData: CreateReviewRequest): Promise<Review> {
    // Get CSRF token for POST request
    const csrfToken = await AuthService.getCsrfToken();

    const response = await fetch(`${API_BASE}/reviews`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
      body: JSON.stringify(reviewData),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.message || 'Failed to create review');
    }

    const data = await response.json();
    return data.data;
  },

  // Upload photos to existing review (step 2a)
  async uploadReviewPhotos(reviewId: number, files: File[]): Promise<string[]> {
    const csrfToken = await AuthService.getCsrfToken();

    const formData = new FormData();
    files.forEach((file) => {
      formData.append(`photos`, file);
    });

    const response = await fetch(`${API_BASE}/reviews/${reviewId}/photos`, {
      method: 'POST',
      headers: {
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
      body: formData,
    });

    if (!response.ok) {
      throw new Error('Failed to upload photos');
    }

    const data = await response.json();
    return data.data.photos;
  },

  // Update a review
  async updateReview(
    id: number,
    updates: Partial<CreateReviewRequest>
  ): Promise<Review> {
    const csrfToken = await AuthService.getCsrfToken();

    const response = await fetch(`${API_BASE}/reviews/${id}`, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
      body: JSON.stringify(updates),
    });

    if (!response.ok) {
      throw new Error('Failed to update review');
    }

    const data = await response.json();
    return data.data;
  },

  // Delete a review
  async deleteReview(id: number): Promise<void> {
    const csrfToken = await AuthService.getCsrfToken();

    const response = await fetch(`${API_BASE}/reviews/${id}`, {
      method: 'DELETE',
      headers: {
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
    });

    if (!response.ok) {
      throw new Error('Failed to delete review');
    }
  },

  // Vote on a review
  async voteReview(
    reviewId: number,
    voteType: 'helpful' | 'not_helpful'
  ): Promise<void> {
    const csrfToken = await AuthService.getCsrfToken();

    const response = await fetch(`${API_BASE}/reviews/${reviewId}/vote`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
      body: JSON.stringify({ vote_type: voteType }),
    });

    if (!response.ok) {
      throw new Error('Failed to vote on review');
    }
  },

  // Confirm review as seller
  async confirmReview(
    reviewId: number,
    notes?: string
  ): Promise<ReviewConfirmation> {
    const csrfToken = await AuthService.getCsrfToken();

    const response = await fetch(`${API_BASE}/reviews/${reviewId}/confirm`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
      body: JSON.stringify({ notes }),
    });

    if (!response.ok) {
      throw new Error('Failed to confirm review');
    }

    const data = await response.json();
    return data.data;
  },

  // Dispute a review
  async disputeReview(
    reviewId: number,
    reason: string
  ): Promise<ReviewDispute> {
    const csrfToken = await AuthService.getCsrfToken();

    const response = await fetch(`${API_BASE}/reviews/${reviewId}/dispute`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
      body: JSON.stringify({ reason }),
    });

    if (!response.ok) {
      throw new Error('Failed to dispute review');
    }

    const data = await response.json();
    return data.data;
  },

  // Add response to review
  async addResponse(reviewId: number, response: string): Promise<Review> {
    const csrfToken = await AuthService.getCsrfToken();

    const res = await fetch(`${API_BASE}/reviews/${reviewId}/response`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
      body: JSON.stringify({ response }),
    });

    if (!res.ok) {
      throw new Error('Failed to add response');
    }

    const data = await res.json();
    return data.data;
  },

  // Upload review photos
  async uploadPhotos(files: File[]): Promise<string[]> {
    const csrfToken = await AuthService.getCsrfToken();

    const formData = new FormData();
    files.forEach((file) => {
      formData.append(`photos`, file);
    });

    const response = await fetch(`${API_BASE}/reviews/upload-photos`, {
      method: 'POST',
      headers: {
        ...(csrfToken ? { 'X-CSRF-Token': csrfToken } : {}),
      },
      credentials: 'include',
      body: formData,
    });

    if (!response.ok) {
      throw new Error('Failed to upload photos');
    }

    const data = await response.json();
    return data.data.photos;
  },
};
