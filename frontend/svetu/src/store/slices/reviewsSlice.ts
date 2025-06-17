import { createSlice, createAsyncThunk } from '@reduxjs/toolkit';
import type { PayloadAction } from '@reduxjs/toolkit';
import type {
  Review,
  ReviewsFilter,
  ReviewStats,
  AggregatedRating,
  CanReviewResponse,
  CreateReviewRequest,
} from '@/types/review';
import { reviewApi } from '@/services/reviewApi';

interface ReviewsState {
  reviews: Review[];
  currentReview: Review | null;
  stats: ReviewStats | null;
  aggregatedRating: AggregatedRating | null;
  canReview: CanReviewResponse | null;
  filters: ReviewsFilter;
  loading: boolean;
  error: string | null;
  totalPages: number;
  currentPage: number;
}

const initialState: ReviewsState = {
  reviews: [],
  currentReview: null,
  stats: null,
  aggregatedRating: null,
  canReview: null,
  filters: {
    page: 1,
    limit: 10,
    sort_by: 'date',
    sort_order: 'desc',
  },
  loading: false,
  error: null,
  totalPages: 0,
  currentPage: 1,
};

// Async thunks
export const fetchReviews = createAsyncThunk(
  'reviews/fetchReviews',
  async (filters: ReviewsFilter) => {
    const response = await reviewApi.getReviews(filters);
    return response;
  }
);

export const fetchReviewStats = createAsyncThunk(
  'reviews/fetchStats',
  async ({
    entityType,
    entityId,
  }: {
    entityType: string;
    entityId: number;
  }) => {
    const response = await reviewApi.getReviewStats(entityType, entityId);
    return response;
  }
);

export const fetchAggregatedRating = createAsyncThunk(
  'reviews/fetchAggregatedRating',
  async ({
    entityType,
    entityId,
  }: {
    entityType: string;
    entityId: number;
  }) => {
    const response = await reviewApi.getAggregatedRating(entityType, entityId);
    return response;
  }
);

export const checkCanReview = createAsyncThunk(
  'reviews/checkCanReview',
  async ({
    entityType,
    entityId,
  }: {
    entityType: string;
    entityId: number;
  }) => {
    const response = await reviewApi.canReview(entityType, entityId);
    return response;
  }
);

export const createReview = createAsyncThunk(
  'reviews/createReview',
  async (reviewData: CreateReviewRequest) => {
    const response = await reviewApi.createReview(reviewData);
    return response;
  }
);

export const voteReview = createAsyncThunk(
  'reviews/voteReview',
  async ({
    reviewId,
    voteType,
  }: {
    reviewId: number;
    voteType: 'helpful' | 'not_helpful';
  }) => {
    const response = await reviewApi.voteReview(reviewId, voteType);
    return response;
  }
);

export const confirmReview = createAsyncThunk(
  'reviews/confirmReview',
  async ({ reviewId, notes }: { reviewId: number; notes?: string }) => {
    const response = await reviewApi.confirmReview(reviewId, notes);
    return response;
  }
);

export const disputeReview = createAsyncThunk(
  'reviews/disputeReview',
  async ({ reviewId, reason }: { reviewId: number; reason: string }) => {
    const response = await reviewApi.disputeReview(reviewId, reason);
    return response;
  }
);

const reviewsSlice = createSlice({
  name: 'reviews',
  initialState,
  reducers: {
    setFilters: (state, action: PayloadAction<Partial<ReviewsFilter>>) => {
      state.filters = { ...state.filters, ...action.payload };
    },
    clearError: (state) => {
      state.error = null;
    },
    setCurrentReview: (state, action: PayloadAction<Review | null>) => {
      state.currentReview = action.payload;
    },
    updateReviewInList: (state, action: PayloadAction<Review>) => {
      const index = state.reviews.findIndex((r) => r.id === action.payload.id);
      if (index !== -1) {
        state.reviews[index] = action.payload;
      }
    },
  },
  extraReducers: (builder) => {
    // Fetch reviews
    builder
      .addCase(fetchReviews.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(fetchReviews.fulfilled, (state, action) => {
        state.loading = false;
        state.reviews = action.payload.reviews;
        state.totalPages = action.payload.totalPages;
        state.currentPage = action.payload.currentPage;
      })
      .addCase(fetchReviews.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to fetch reviews';
      });

    // Fetch stats
    builder.addCase(fetchReviewStats.fulfilled, (state, action) => {
      state.stats = action.payload;
    });

    // Fetch aggregated rating
    builder.addCase(fetchAggregatedRating.fulfilled, (state, action) => {
      state.aggregatedRating = action.payload;
    });

    // Check can review
    builder.addCase(checkCanReview.fulfilled, (state, action) => {
      state.canReview = action.payload;
    });

    // Create review
    builder
      .addCase(createReview.pending, (state) => {
        state.loading = true;
        state.error = null;
      })
      .addCase(createReview.fulfilled, (state, action) => {
        state.loading = false;
        state.reviews.unshift(action.payload);
        // Reset canReview since user just submitted a review
        state.canReview = { can_review: false, reason: 'already_reviewed' };
      })
      .addCase(createReview.rejected, (state, action) => {
        state.loading = false;
        state.error = action.error.message || 'Failed to create review';
      });

    // Vote review
    builder.addCase(voteReview.fulfilled, (state, action) => {
      const review = state.reviews.find(
        (r) => r.id === action.meta.arg.reviewId
      );
      if (review) {
        // Update vote counts based on response
        const { voteType } = action.meta.arg;
        if (voteType === 'helpful') {
          review.helpful_votes += 1;
          if (review.current_user_vote === 'not_helpful') {
            review.not_helpful_votes -= 1;
          }
        } else {
          review.not_helpful_votes += 1;
          if (review.current_user_vote === 'helpful') {
            review.helpful_votes -= 1;
          }
        }
        review.current_user_vote = voteType;
      }
    });
  },
});

export const { setFilters, clearError, setCurrentReview, updateReviewInList } =
  reviewsSlice.actions;
export default reviewsSlice.reducer;
