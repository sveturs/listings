'use client';

import React, { useState } from 'react';
import Image from 'next/image';
import { useLocale, useTranslations } from 'next-intl';
import { formatDistanceToNow } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import { RatingDisplay } from './RatingDisplay';
import { ImageGallery } from './ImageGallery';
import {
  useVoteReview,
  useConfirmReview,
  useDisputeReview,
} from '@/hooks/useReviews';
import type { Review, ReviewsFilter } from '@/types/review';

interface ReviewListProps {
  reviews: Review[];
  totalPages: number;
  currentPage: number;
  onPageChange: (page: number) => void;
  onSortChange: (sortBy: string) => void;
  onFilterChange: (filters: Partial<ReviewsFilter>) => void;
  currentUserId?: number;
  sellerId?: number;
}

export const ReviewList: React.FC<ReviewListProps> = ({
  reviews,
  totalPages,
  currentPage,
  onPageChange,
  onSortChange,
  onFilterChange,
  currentUserId,
  sellerId,
}) => {
  const locale = useLocale();
  const t = useTranslations('reviews');
  const [selectedRating, setSelectedRating] = useState<number | null>(null);
  const [showOnlyWithPhotos, setShowOnlyWithPhotos] = useState(false);
  const [showOnlyVerified, setShowOnlyVerified] = useState(false);
  const [disputeReviewId, setDisputeReviewId] = useState<number | null>(null);
  const [disputeReason, setDisputeReason] = useState('');

  // Image gallery state
  const [galleryImages, setGalleryImages] = useState<string[]>([]);
  const [galleryInitialIndex, setGalleryInitialIndex] = useState(0);
  const [isGalleryOpen, setIsGalleryOpen] = useState(false);

  const voteReviewMutation = useVoteReview();
  const confirmReviewMutation = useConfirmReview();
  const disputeReviewMutation = useDisputeReview();

  const dateLocale = locale === 'ru' ? ru : enUS;

  const formatDate = (date: string) => {
    return formatDistanceToNow(new Date(date), {
      addSuffix: true,
      locale: dateLocale,
    });
  };

  const handleVote = async (
    reviewId: number,
    voteType: 'helpful' | 'not_helpful'
  ) => {
    try {
      await voteReviewMutation.mutateAsync({ reviewId, voteType });
    } catch (error) {
      console.error('Failed to vote:', error);
    }
  };

  const handleConfirm = async (reviewId: number) => {
    try {
      await confirmReviewMutation.mutateAsync({ reviewId });
    } catch (error) {
      console.error('Failed to confirm review:', error);
    }
  };

  const handleDispute = async () => {
    if (!disputeReviewId || !disputeReason.trim()) return;

    try {
      await disputeReviewMutation.mutateAsync({
        reviewId: disputeReviewId,
        reason: disputeReason,
      });
      setDisputeReviewId(null);
      setDisputeReason('');
    } catch (error) {
      console.error('Failed to dispute review:', error);
    }
  };

  const handleRatingFilter = (rating: number | null) => {
    setSelectedRating(rating);
    if (rating === null) {
      onFilterChange({ min_rating: undefined, max_rating: undefined });
    } else {
      onFilterChange({ min_rating: rating, max_rating: rating });
    }
  };

  const openImageGallery = (images: string[], initialIndex: number = 0) => {
    setGalleryImages(images);
    setGalleryInitialIndex(initialIndex);
    setIsGalleryOpen(true);
  };

  const closeImageGallery = () => {
    setIsGalleryOpen(false);
  };

  // Debug: log reviews data
  console.log('ReviewList: reviews data', reviews);

  return (
    <div className="space-y-6">
      {/* Filters */}
      <div className="bg-base-100 rounded-lg shadow-sm border border-base-200 overflow-hidden">
        <div className="p-4 lg:p-6">
          <div className="flex flex-wrap gap-4">
            {/* Sort */}
            <div className="form-control">
              <label className="label">
                <span className="label-text text-xs font-medium uppercase tracking-wider text-base-content/70">
                  {t('list.sortBy')}
                </span>
              </label>
              <select
                className="select select-bordered select-sm min-w-[140px]"
                onChange={(e) => onSortChange(e.target.value)}
              >
                <option value="date">{t('list.byDate')}</option>
                <option value="rating">{t('list.byRating')}</option>
                <option value="likes">{t('list.byHelpfulness')}</option>
              </select>
            </div>

            {/* Rating filter */}
            <div className="form-control">
              <label className="label">
                <span className="label-text text-xs font-medium uppercase tracking-wider text-base-content/70">
                  {t('list.rating')}
                </span>
              </label>
              <div className="flex gap-1">
                <button
                  onClick={() => handleRatingFilter(null)}
                  className={`px-3 py-1.5 text-sm font-medium rounded-md transition-all duration-200 ${
                    selectedRating === null
                      ? 'bg-primary text-primary-content shadow-sm'
                      : 'bg-base-200 text-base-content hover:bg-base-300'
                  }`}
                >
                  {t('list.all')}
                </button>
                {[5, 4, 3, 2, 1].map((rating) => (
                  <button
                    key={rating}
                    onClick={() => handleRatingFilter(rating)}
                    className={`px-3 py-1.5 text-sm font-medium rounded-md transition-all duration-200 flex items-center gap-1 ${
                      selectedRating === rating
                        ? 'bg-primary text-primary-content shadow-sm'
                        : 'bg-base-200 text-base-content hover:bg-base-300'
                    }`}
                  >
                    {rating}
                    <svg
                      className="w-3.5 h-3.5"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  </button>
                ))}
              </div>
            </div>

            {/* Toggles */}
            <div className="flex items-end gap-4">
              <label className="label cursor-pointer">
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm mr-2"
                  checked={showOnlyWithPhotos}
                  onChange={(e) => {
                    setShowOnlyWithPhotos(e.target.checked);
                    // Note: This would require backend support for photo filter
                  }}
                />
                <span className="label-text text-sm">{t('list.withPhotos')}</span>
              </label>

              <label className="label cursor-pointer">
                <input
                  type="checkbox"
                  className="checkbox checkbox-sm mr-2"
                  checked={showOnlyVerified}
                  onChange={(e) => {
                    setShowOnlyVerified(e.target.checked);
                    // Note: This would require backend support for verified filter
                  }}
                />
                <span className="label-text text-sm">{t('list.verified')}</span>
              </label>
            </div>
          </div>
        </div>
      </div>

      {/* Reviews */}
      <div className="space-y-4">
        {reviews.map((review, index) => (
          <div
            key={review.id}
            className="bg-base-100 rounded-lg shadow-sm border border-base-200 overflow-hidden
                     hover:shadow-md transition-all duration-300
                     animate-in fade-in-50 slide-in-from-bottom-2"
            style={{ animationDelay: `${index * 50}ms` }}
          >
            <div className="p-4 lg:p-6">
              {/* Header */}
              <div className="flex items-start justify-between">
                <div className="flex items-center gap-3">
                  {review.user?.avatar ? (
                    <Image
                      src={review.user.avatar}
                      alt={review.user.name}
                      width={40}
                      height={40}
                      className="w-10 h-10 rounded-full object-cover"
                    />
                  ) : (
                    <div className="w-10 h-10 rounded-full bg-gradient-to-br from-primary/20 to-secondary/20 flex items-center justify-center">
                      <span className="text-sm font-semibold text-base-content">
                        {review.user?.name?.[0]?.toUpperCase() || '?'}
                      </span>
                    </div>
                  )}
                  <div>
                    <div className="font-medium text-base-content">
                      {review.user?.name || t('list.anonymous')}
                    </div>
                    <div className="text-xs text-base-content/60">
                      {formatDate(review.created_at)}
                    </div>
                  </div>
                </div>
                <div className="flex items-center gap-2">
                  <RatingDisplay
                    rating={review.rating}
                    size="sm"
                    showValue={false}
                  />
                  {review.is_verified_purchase && (
                    <span className="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium bg-success/10 text-success rounded-md">
                      <svg
                        className="w-3 h-3"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      {t('list.verifiedPurchase')}
                    </span>
                  )}
                  {review.seller_confirmed && (
                    <span className="inline-flex items-center gap-1 px-2 py-1 text-xs font-medium bg-info/10 text-info rounded-md">
                      <svg
                        className="w-3 h-3"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M6.267 3.455a3.066 3.066 0 001.745-.723 3.066 3.066 0 013.976 0 3.066 3.066 0 001.745.723 3.066 3.066 0 012.812 2.812c.051.643.304 1.254.723 1.745a3.066 3.066 0 010 3.976 3.066 3.066 0 00-.723 1.745 3.066 3.066 0 01-2.812 2.812 3.066 3.066 0 00-1.745.723 3.066 3.066 0 01-3.976 0 3.066 3.066 0 00-1.745-.723 3.066 3.066 0 01-2.812-2.812 3.066 3.066 0 00-.723-1.745 3.066 3.066 0 010-3.976 3.066 3.066 0 00.723-1.745 3.066 3.066 0 012.812-2.812zm7.44 5.252a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      {t('list.sellerConfirmed')}
                    </span>
                  )}
                </div>
              </div>

              {/* Content */}
              <div className="mt-4 space-y-4">
                {review.comment && (
                  <p className="text-sm text-base-content/80 leading-relaxed">
                    {review.translations?.[locale]?.comment || review.comment}
                  </p>
                )}

                {review.pros && (
                  <div className="flex gap-3">
                    <div className="flex-shrink-0 mt-0.5">
                      <div className="w-8 h-8 rounded-md bg-success/10 flex items-center justify-center">
                        <svg
                          className="w-4 h-4 text-success"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                      </div>
                    </div>
                    <div className="flex-1">
                      <span className="text-xs font-medium text-success uppercase tracking-wider">
                        {t('list.pros')}
                      </span>
                      <p className="text-sm text-base-content/80 mt-1 leading-relaxed">
                        {review.translations?.[locale]?.pros || review.pros}
                      </p>
                    </div>
                  </div>
                )}

                {review.cons && (
                  <div className="flex gap-3">
                    <div className="flex-shrink-0 mt-0.5">
                      <div className="w-8 h-8 rounded-md bg-warning/10 flex items-center justify-center">
                        <svg
                          className="w-4 h-4 text-warning"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                      </div>
                    </div>
                    <div className="flex-1">
                      <span className="text-xs font-medium text-warning uppercase tracking-wider">
                        {t('list.cons')}
                      </span>
                      <p className="text-sm text-base-content/80 mt-1 leading-relaxed">
                        {review.translations?.[locale]?.cons || review.cons}
                      </p>
                    </div>
                  </div>
                )}

                {/* Photos */}
                {review.photos && review.photos.length > 0 && (
                  <div className="flex flex-wrap gap-2 mt-4">
                    {review.photos.map((photo, index) => (
                      <div
                        key={index}
                        className="relative group cursor-pointer"
                        onClick={() => openImageGallery(review.photos!, index)}
                      >
                        <Image
                          src={photo}
                          alt={`Review photo ${index + 1}`}
                          width={80}
                          height={80}
                          className="w-20 h-20 object-cover rounded-md transition-all duration-200
                                   group-hover:brightness-75 group-hover:scale-105"
                        />
                        <div className="absolute inset-0 flex items-center justify-center opacity-0 group-hover:opacity-100 transition-opacity duration-200 rounded-md bg-black/20">
                          <svg
                            className="w-6 h-6 text-white drop-shadow-lg"
                            fill="none"
                            viewBox="0 0 24 24"
                            stroke="currentColor"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7"
                            />
                          </svg>
                        </div>
                        {review.photos!.length > 1 && (
                          <div className="absolute top-1 right-1 bg-black/60 text-white text-xs px-1.5 py-0.5 rounded-md font-medium">
                            {index + 1}/{review.photos!.length}
                          </div>
                        )}
                      </div>
                    ))}
                  </div>
                )}
              </div>

              {/* Actions */}
              <div className="flex items-center justify-between mt-4 pt-4 border-t border-base-200/50">
                <div className="flex items-center gap-4">
                  {/* Helpful votes */}
                  <div className="flex items-center gap-2">
                    <span className="text-sm text-base-content/60">
                      {t('list.helpful')}
                    </span>
                    <button
                      onClick={() => handleVote(review.id, 'helpful')}
                      disabled={!currentUserId || voteReviewMutation.isPending}
                      className={`inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-md
                                transition-all duration-200 ${
                                  review.current_user_vote === 'helpful'
                                    ? 'bg-success/10 text-success'
                                    : 'bg-base-200 text-base-content/70 hover:bg-base-300'
                                } disabled:opacity-50 disabled:cursor-not-allowed`}
                    >
                      <svg
                        className="w-3.5 h-3.5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"
                        />
                      </svg>
                      {review.helpful_votes || 0}
                    </button>
                    <button
                      onClick={() => handleVote(review.id, 'not_helpful')}
                      disabled={!currentUserId || voteReviewMutation.isPending}
                      className={`inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium rounded-md
                                transition-all duration-200 ${
                                  review.current_user_vote === 'not_helpful'
                                    ? 'bg-error/10 text-error'
                                    : 'bg-base-200 text-base-content/70 hover:bg-base-300'
                                } disabled:opacity-50 disabled:cursor-not-allowed`}
                    >
                      <svg
                        className="w-3.5 h-3.5 rotate-180"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M14 10h4.764a2 2 0 011.789 2.894l-3.5 7A2 2 0 0115.263 21h-4.017c-.163 0-.326-.02-.485-.06L7 20m7-10V5a2 2 0 00-2-2h-.095c-.5 0-.905.405-.905.905 0 .714-.211 1.412-.608 2.006L7 11v9m7-10h-2M7 20H5a2 2 0 01-2-2v-6a2 2 0 012-2h2.5"
                        />
                      </svg>
                      {review.not_helpful_votes || 0}
                    </button>
                  </div>
                </div>

                <div className="flex items-center gap-2">
                  {/* Seller actions */}
                  {currentUserId === sellerId && !review.seller_confirmed && (
                    <button
                      onClick={() => handleConfirm(review.id)}
                      disabled={confirmReviewMutation.isPending}
                      className="inline-flex items-center gap-1.5 px-4 py-1.5 text-xs font-medium
                               bg-success text-success-content rounded-md
                               hover:bg-success/90 transition-colors duration-200
                               disabled:opacity-50 disabled:cursor-not-allowed"
                    >
                      <svg
                        className="w-3.5 h-3.5"
                        fill="none"
                        viewBox="0 0 24 24"
                        stroke="currentColor"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth={2}
                          d="M5 13l4 4L19 7"
                        />
                      </svg>
                      {t('list.confirm')}
                    </button>
                  )}

                  {/* Dispute */}
                  {currentUserId &&
                    (currentUserId === sellerId ||
                      currentUserId === review.entity_id) &&
                    !review.has_active_dispute && (
                      <button
                        onClick={() => setDisputeReviewId(review.id)}
                        className="inline-flex items-center gap-1.5 px-3 py-1.5 text-xs font-medium
                                 text-base-content/70 rounded-md
                                 hover:bg-base-200 transition-colors duration-200"
                      >
                        <svg
                          className="w-3.5 h-3.5"
                          fill="none"
                          viewBox="0 0 24 24"
                          stroke="currentColor"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                          />
                        </svg>
                        {t('list.dispute')}
                      </button>
                    )}
                </div>
              </div>

              {/* Seller response */}
              {review.responses && review.responses.length > 0 && (
                <div className="mt-4 ml-12 p-3 bg-base-200/30 rounded-md border-l-2 border-primary/20">
                  {review.responses.map((response) => (
                    <div key={response.id}>
                      <div className="flex items-center gap-2 mb-2">
                        <div className="w-6 h-6 rounded-full bg-primary/10 flex items-center justify-center">
                          <svg
                            className="w-3 h-3 text-primary"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                          >
                            <path
                              fillRule="evenodd"
                              d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z"
                              clipRule="evenodd"
                            />
                          </svg>
                        </div>
                        <span className="text-xs font-medium text-base-content">
                          {response.user?.name || t('list.seller')}
                        </span>
                        <span className="text-xs text-base-content/50">
                          {formatDate(response.created_at)}
                        </span>
                      </div>
                      <p className="text-sm text-base-content/70 leading-relaxed ml-8">
                        {response.response}
                      </p>
                    </div>
                  ))}
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex justify-center mt-8">
          <div className="inline-flex items-center gap-1">
            <button
              onClick={() => onPageChange(currentPage - 1)}
              disabled={currentPage === 1}
              className="p-2 rounded-md bg-base-100 border border-base-200 hover:bg-base-200
                       disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
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
                  d="M15 19l-7-7 7-7"
                />
              </svg>
            </button>
            {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
              const page = i + 1;
              return (
                <button
                  key={page}
                  onClick={() => onPageChange(page)}
                  className={`px-3 py-2 text-sm font-medium rounded-md transition-colors duration-200 ${
                    currentPage === page
                      ? 'bg-primary text-primary-content'
                      : 'bg-base-100 border border-base-200 hover:bg-base-200'
                  }`}
                >
                  {page}
                </button>
              );
            })}
            {totalPages > 5 && (
              <>
                <span className="px-2 text-base-content/50">...</span>
                <button
                  onClick={() => onPageChange(totalPages)}
                  className={`px-3 py-2 text-sm font-medium rounded-md transition-colors duration-200 ${
                    currentPage === totalPages
                      ? 'bg-primary text-primary-content'
                      : 'bg-base-100 border border-base-200 hover:bg-base-200'
                  }`}
                >
                  {totalPages}
                </button>
              </>
            )}
            <button
              onClick={() => onPageChange(currentPage + 1)}
              disabled={currentPage === totalPages}
              className="p-2 rounded-md bg-base-100 border border-base-200 hover:bg-base-200
                       disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
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
                  d="M9 5l7 7-7 7"
                />
              </svg>
            </button>
          </div>
        </div>
      )}

      {/* Dispute Modal */}
      {disputeReviewId && (
        <div className="modal modal-open">
          <div className="modal-box max-w-md">
            <button
              onClick={() => {
                setDisputeReviewId(null);
                setDisputeReason('');
              }}
              className="btn btn-sm btn-circle btn-ghost absolute right-2 top-2"
            >
              ✕
            </button>
            <h3 className="font-bold text-lg mb-4">{t('list.disputeReview')}</h3>
            <div className="space-y-4">
              <div>
                <label className="text-sm text-base-content/70 mb-2 block">
                  {t('list.disputeReason')}
                </label>
                <textarea
                  value={disputeReason}
                  onChange={(e) => setDisputeReason(e.target.value)}
                  className="w-full min-h-[120px] p-3 rounded-lg border border-base-200
                           bg-base-100 resize-none transition-all duration-200
                           focus:outline-none focus:ring-2 focus:ring-primary/20 focus:border-primary"
                  placeholder={t('list.disputePlaceholder')}
                />
              </div>
              <div className="flex gap-3 justify-end">
                <button
                  onClick={() => {
                    setDisputeReviewId(null);
                    setDisputeReason('');
                  }}
                  className="btn btn-ghost"
                >
                  {t('list.cancel')}
                </button>
                <button
                  onClick={handleDispute}
                  disabled={
                    !disputeReason.trim() || disputeReviewMutation.isPending
                  }
                  className="btn btn-primary"
                >
                  {disputeReviewMutation.isPending && (
                    <span className="loading loading-spinner loading-sm mr-2"></span>
                  )}
                  {t('list.submit')}
                </button>
              </div>
            </div>
          </div>
          <div
            className="modal-backdrop bg-black/50"
            onClick={() => {
              setDisputeReviewId(null);
              setDisputeReason('');
            }}
          ></div>
        </div>
      )}

      {/* Image Gallery */}
      <ImageGallery
        images={galleryImages}
        initialIndex={galleryInitialIndex}
        isOpen={isGalleryOpen}
        onClose={closeImageGallery}
      />
    </div>
  );
};
