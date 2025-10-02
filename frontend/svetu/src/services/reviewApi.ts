import { apiClient } from '@/services/api-client';
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

    const response = await apiClient.get<any>(`/reviews?${params}`);
    console.log('reviewApi.getReviews: raw response', response.data);

    // Backend возвращает структуру: { data: { data: [...], meta: {...} } }
    const reviewsData = response.data?.data || response.data?.reviews || [];
    const meta = response.data?.meta || {};

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
    const response = await apiClient.get<{ data: Review }>(`/reviews/${id}`);
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data;
  },

  // Get review statistics
  async getReviewStats(
    entityType: string,
    entityId: number
  ): Promise<ReviewStats> {
    const response = await apiClient.get<{ data: ReviewStats }>(
      `/entity/${entityType}/${entityId}/stats`
    );
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data;
  },

  // Get aggregated rating
  async getAggregatedRating(
    entityType: string,
    entityId: number
  ): Promise<AggregatedRating> {
    const endpoint =
      entityType === 'user'
        ? `/users/${entityId}/aggregated-rating`
        : `/storefronts/${entityId}/aggregated-rating`;

    const response = await apiClient.get<{ data: AggregatedRating }>(endpoint);
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data;
  },

  // Check if user can review
  async canReview(
    entityType: string,
    entityId: number
  ): Promise<CanReviewResponse> {
    try {
      const response = await apiClient.get<{ data: CanReviewResponse }>(
        `/review-permission/${entityType}/${entityId}`
      );
      if (!response.data) {
        throw new Error('No data received');
      }
      return response.data.data;
    } catch (error: any) {
      // Return default response instead of throwing error for 401 or 404
      if (error.response?.status === 401 || error.response?.status === 404) {
        return { can_review: false, reason: 'unauthorized' };
      }
      throw error;
    }
  },

  // Create a draft review (step 1)
  async createDraftReview(reviewData: CreateReviewRequest): Promise<Review> {
    const response = await apiClient.post<{ data: Review }>(
      `/reviews/draft`,
      reviewData
    );
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data;
  },

  // Publish a draft review (step 2b)
  async publishReview(reviewId: number): Promise<Review> {
    const response = await apiClient.post<{ data: Review }>(
      `/reviews/${reviewId}/publish`,
      {}
    );
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data;
  },

  // Legacy: Create a new review (single step)
  async createReview(reviewData: CreateReviewRequest): Promise<Review> {
    const response = await apiClient.post<{ data: Review }>(
      `/reviews`,
      reviewData
    );
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data;
  },

  // Upload photos to existing review (step 2a)
  async uploadReviewPhotos(reviewId: number, files: File[]): Promise<string[]> {
    const formData = new FormData();
    files.forEach((file) => {
      formData.append(`photos`, file);
    });

    const response = await apiClient.post<{ data: { photos: string[] } }>(
      `/reviews/${reviewId}/photos`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    );
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data.photos;
  },

  // Update a review
  async updateReview(
    id: number,
    updates: Partial<CreateReviewRequest>
  ): Promise<Review> {
    const response = await apiClient.put<{ data: Review }>(
      `/reviews/${id}`,
      updates
    );
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data;
  },

  // Delete a review
  async deleteReview(id: number): Promise<void> {
    await apiClient.delete(`/reviews/${id}`);
  },

  // Vote on a review
  async voteReview(
    reviewId: number,
    voteType: 'helpful' | 'not_helpful'
  ): Promise<void> {
    await apiClient.post(`/reviews/${reviewId}/vote`, { vote_type: voteType });
  },

  // Confirm review as seller
  async confirmReview(
    reviewId: number,
    notes?: string
  ): Promise<ReviewConfirmation> {
    const response = await apiClient.post<{ data: ReviewConfirmation }>(
      `/reviews/${reviewId}/confirm`,
      { notes }
    );
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data;
  },

  // Dispute a review
  async disputeReview(
    reviewId: number,
    reason: string
  ): Promise<ReviewDispute> {
    const response = await apiClient.post<{ data: ReviewDispute }>(
      `/reviews/${reviewId}/dispute`,
      { reason }
    );
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data;
  },

  // Add response to review
  async addResponse(reviewId: number, response: string): Promise<Review> {
    const res = await apiClient.post<{ data: Review }>(
      `/reviews/${reviewId}/response`,
      { response }
    );
    if (!res.data) {
      throw new Error('No data received');
    }
    return res.data.data;
  },

  // Upload review photos
  async uploadPhotos(files: File[]): Promise<string[]> {
    const formData = new FormData();
    files.forEach((file) => {
      formData.append(`photos`, file);
    });

    const response = await apiClient.post<{ data: { photos: string[] } }>(
      `/reviews/upload-photos`,
      formData,
      {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      }
    );
    if (!response.data) {
      throw new Error('No data received');
    }
    return response.data.data.photos;
  },
};
