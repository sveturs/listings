'use client';

import React, { useState } from 'react';
import { useAuth } from '@/contexts/AuthContext';
import { useLocale } from 'next-intl';
import {
  useReviews,
  useReviewStats,
  useCanReview,
  useCreateReview,
} from '@/hooks/useReviews';
import { RatingStats } from './RatingStats';
import { ReviewList } from './ReviewList';
import { ReviewForm } from './ReviewForm';
import type { ReviewsFilter } from '@/types/review';

interface ReviewsSectionProps {
  entityType: 'listing' | 'user' | 'storefront';
  entityId: number;
  sellerId?: number;
  storefrontId?: number;
}

export const ReviewsSection: React.FC<ReviewsSectionProps> = ({
  entityType,
  entityId,
  sellerId,
  storefrontId,
}) => {
  const { user } = useAuth();
  const locale = useLocale();
  const [showReviewForm, setShowReviewForm] = useState(false);
  const [filters, setFilters] = useState<ReviewsFilter>({
    entity_type: entityType,
    entity_id: entityId,
    page: 1,
    limit: 10,
    sort_by: 'date',
    sort_order: 'desc',
  });

  // Fetch data
  const { data: reviewsData, isLoading: reviewsLoading } = useReviews(filters);
  const { data: statsData, isLoading: statsLoading } = useReviewStats(
    entityType,
    entityId
  );
  const { data: canReviewData } = useCanReview(entityType, entityId, user?.id);
  const createReviewMutation = useCreateReview();

  // Check if user can write a review
  const canWriteReview =
    user && canReviewData?.can_review && user.id !== sellerId;

  // Debug logging
  console.log('ReviewsSection debug:', {
    user: user?.id,
    canReviewData,
    sellerId,
    canWriteReview,
    reviewsData,
    statsData,
    reviewsLoading,
    statsLoading,
  });

  const handleCreateReview = async (reviewData: any) => {
    try {
      await createReviewMutation.mutateAsync({
        ...reviewData,
        entity_type: entityType,
        entity_id: entityId,
        storefront_id: storefrontId,
        original_language: locale,
      });
      setShowReviewForm(false);
    } catch (error) {
      console.error('Failed to create review:', error);
    }
  };

  const handlePageChange = (page: number) => {
    setFilters((prev) => ({ ...prev, page }));
  };

  const handleSortChange = (sortBy: string) => {
    setFilters((prev) => ({
      ...prev,
      sort_by: sortBy as 'date' | 'rating' | 'likes',
      page: 1,
    }));
  };

  const handleFilterChange = (newFilters: Partial<ReviewsFilter>) => {
    setFilters((prev) => ({ ...prev, ...newFilters, page: 1 }));
  };

  if (reviewsLoading || statsLoading) {
    return (
      <div className="space-y-6">
        <div className="skeleton h-48"></div>
        <div className="skeleton h-96"></div>
      </div>
    );
  }

  const reviews = reviewsData?.reviews || [];
  const totalPages = reviewsData?.totalPages || 0;
  const stats = statsData || {
    total_reviews: 0,
    average_rating: 0,
    verified_reviews: 0,
    rating_distribution: {},
    photo_reviews: 0,
  };

  return (
    <div className="space-y-6">
      {/* Reviews Header */}
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold flex items-center gap-2">
          <svg
            className="w-6 h-6 text-primary"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"
            />
          </svg>
          {locale === 'ru' ? 'Отзывы' : 'Reviews'}
          {stats.total_reviews > 0 && (
            <span className="text-base-content/60 text-base font-normal">
              ({stats.total_reviews})
            </span>
          )}
        </h2>

        {canWriteReview && !showReviewForm && (
          <button
            onClick={() => setShowReviewForm(true)}
            className="btn btn-primary rounded-lg px-4 font-medium"
          >
            <svg
              className="w-4 h-4 mr-2"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 4v16m8-8H4"
              />
            </svg>
            {locale === 'ru' ? 'Написать отзыв' : 'Write a review'}
          </button>
        )}
      </div>

      {/* Can't review message */}
      {user && canReviewData && !canReviewData.can_review && (
        <div className="bg-warning/5 border border-warning/20 rounded-lg p-4">
          <div className="flex items-start gap-3">
            <svg
              className="w-5 h-5 text-warning flex-shrink-0 mt-0.5"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              />
            </svg>
            <span className="text-sm text-base-content/80">
              {canReviewData.reason === 'already_reviewed'
                ? locale === 'ru'
                  ? 'Вы уже оставили отзыв'
                  : 'You have already reviewed this item'
                : canReviewData.reason === 'insufficient_chat_activity'
                  ? locale === 'ru'
                    ? 'Для отзыва необходимо обменяться минимум 5 сообщениями с продавцом'
                    : 'You need to exchange at least 5 messages with the seller to leave a review'
                  : locale === 'ru'
                    ? 'Вы не можете оставить отзыв'
                    : 'You cannot leave a review'}
            </span>
          </div>
        </div>
      )}

      {/* Review Form */}
      {showReviewForm && (
        <div className="bg-base-100 rounded-lg shadow-sm border border-base-200 overflow-hidden">
          <div className="p-4 lg:p-6">
            <div className="flex items-center justify-between mb-4">
              <h3 className="text-lg font-semibold text-base-content">
                {locale === 'ru' ? 'Ваш отзыв' : 'Your review'}
              </h3>
              <button
                onClick={() => setShowReviewForm(false)}
                className="btn btn-ghost btn-sm btn-circle"
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>
            <ReviewForm
              onSubmit={handleCreateReview}
              onCancel={() => setShowReviewForm(false)}
              isSubmitting={createReviewMutation.isPending}
            />
          </div>
        </div>
      )}

      {/* Rating Statistics */}
      {stats.total_reviews > 0 && <RatingStats stats={stats} />}

      {/* Reviews List */}
      {reviews.length > 0 ? (
        <ReviewList
          reviews={reviews}
          totalPages={totalPages}
          currentPage={filters.page || 1}
          onPageChange={handlePageChange}
          onSortChange={handleSortChange}
          onFilterChange={handleFilterChange}
          currentUserId={user?.id}
          sellerId={sellerId}
        />
      ) : (
        <div className="text-center py-12">
          <svg
            className="w-16 h-16 mx-auto mb-4 text-base-content/20"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z"
            />
          </svg>
          <p className="text-base-content/60">
            {locale === 'ru'
              ? 'Пока нет отзывов. Будьте первым!'
              : 'No reviews yet. Be the first!'}
          </p>
        </div>
      )}
    </div>
  );
};
