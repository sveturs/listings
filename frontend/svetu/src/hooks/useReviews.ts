import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { useDispatch } from 'react-redux';
import { useEffect } from 'react';
import { reviewApi } from '@/services/reviewApi';
import type { ReviewsFilter, CreateReviewRequest } from '@/types/review';
import {
  setCurrentReview,
  updateReviewInList,
} from '@/store/slices/reviewsSlice';

// Query keys
const QUERY_KEYS = {
  reviews: (filters: ReviewsFilter) => ['reviews', filters],
  review: (id: number) => ['review', id],
  stats: (entityType: string, entityId: number) => [
    'review-stats',
    entityType,
    entityId,
  ],
  aggregatedRating: (entityType: string, entityId: number) => [
    'aggregated-rating',
    entityType,
    entityId,
  ],
  canReview: (entityType: string, entityId: number) => [
    'can-review',
    entityType,
    entityId,
  ],
};

// Fetch reviews hook
export const useReviews = (filters: ReviewsFilter) => {
  return useQuery({
    queryKey: QUERY_KEYS.reviews(filters),
    queryFn: async () => {
      console.log('useReviews: fetching with filters', filters);
      const result = await reviewApi.getReviews(filters);
      console.log('useReviews: result', result);
      return result;
    },
    staleTime: 5 * 60 * 1000, // 5 minutes
  });
};

// Fetch single review hook
export const useReview = (id: number) => {
  const dispatch = useDispatch();

  const query = useQuery({
    queryKey: QUERY_KEYS.review(id),
    queryFn: () => reviewApi.getReviewById(id),
    enabled: !!id,
  });

  // Handle success in effect
  useEffect(() => {
    if (query.data) {
      dispatch(setCurrentReview(query.data));
    }
  }, [query.data, dispatch]);

  return query;
};

// Fetch review stats hook
export const useReviewStats = (entityType: string, entityId: number) => {
  return useQuery({
    queryKey: QUERY_KEYS.stats(entityType, entityId),
    queryFn: () => reviewApi.getReviewStats(entityType, entityId),
    enabled: !!entityType && !!entityId,
    staleTime: 5 * 60 * 1000,
  });
};

// Fetch aggregated rating hook
export const useAggregatedRating = (entityType: string, entityId: number) => {
  return useQuery({
    queryKey: QUERY_KEYS.aggregatedRating(entityType, entityId),
    queryFn: () => reviewApi.getAggregatedRating(entityType, entityId),
    enabled:
      !!entityType &&
      !!entityId &&
      (entityType === 'user' || entityType === 'storefront'),
    staleTime: 5 * 60 * 1000,
  });
};

// Check if user can review hook
export const useCanReview = (
  entityType: string,
  entityId: number,
  userId?: number
) => {
  return useQuery({
    queryKey: QUERY_KEYS.canReview(entityType, entityId),
    queryFn: () => reviewApi.canReview(entityType, entityId),
    enabled: !!entityType && !!entityId && !!userId, // Only run if user is authenticated
    staleTime: 60 * 1000, // 1 minute
  });
};

// Create review mutation
export const useCreateReview = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (reviewData: CreateReviewRequest) =>
      reviewApi.createReview(reviewData),
    onSuccess: (data, variables) => {
      // Invalidate related queries
      queryClient.invalidateQueries({ queryKey: ['reviews'] });
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.stats(variables.entity_type, variables.entity_id),
      });
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.canReview(
          variables.entity_type,
          variables.entity_id
        ),
      });

      // Invalidate aggregated ratings if applicable
      if (
        variables.entity_type === 'listing' &&
        data.entity_origin_type &&
        data.entity_origin_id
      ) {
        queryClient.invalidateQueries({
          queryKey: QUERY_KEYS.aggregatedRating(
            data.entity_origin_type,
            data.entity_origin_id
          ),
        });
      }
    },
  });
};

// Update review mutation
export const useUpdateReview = () => {
  const queryClient = useQueryClient();
  const dispatch = useDispatch();

  return useMutation({
    mutationFn: ({
      id,
      updates,
    }: {
      id: number;
      updates: Partial<CreateReviewRequest>;
    }) => reviewApi.updateReview(id, updates),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['reviews'] });
      queryClient.invalidateQueries({ queryKey: QUERY_KEYS.review(data.id) });
      dispatch(updateReviewInList(data));
    },
  });
};

// Delete review mutation
export const useDeleteReview = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (id: number) => reviewApi.deleteReview(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['reviews'] });
    },
  });
};

// Vote review mutation
export const useVoteReview = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      reviewId,
      voteType,
    }: {
      reviewId: number;
      voteType: 'helpful' | 'not_helpful';
    }) => reviewApi.voteReview(reviewId, voteType),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['reviews'] });
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.review(variables.reviewId),
      });
    },
  });
};

// Confirm review mutation
export const useConfirmReview = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ reviewId, notes }: { reviewId: number; notes?: string }) =>
      reviewApi.confirmReview(reviewId, notes),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['reviews'] });
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.review(variables.reviewId),
      });
    },
  });
};

// Dispute review mutation
export const useDisputeReview = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ reviewId, reason }: { reviewId: number; reason: string }) =>
      reviewApi.disputeReview(reviewId, reason),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ['reviews'] });
      queryClient.invalidateQueries({
        queryKey: QUERY_KEYS.review(variables.reviewId),
      });
    },
  });
};

// Add response mutation
export const useAddReviewResponse = () => {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({
      reviewId,
      response,
    }: {
      reviewId: number;
      response: string;
    }) => reviewApi.addResponse(reviewId, response),
    onSuccess: (data) => {
      queryClient.invalidateQueries({ queryKey: ['reviews'] });
      queryClient.invalidateQueries({ queryKey: QUERY_KEYS.review(data.id) });
    },
  });
};

// Upload photos mutation
export const useUploadReviewPhotos = () => {
  return useMutation({
    mutationFn: (files: File[]) => reviewApi.uploadPhotos(files),
  });
};
