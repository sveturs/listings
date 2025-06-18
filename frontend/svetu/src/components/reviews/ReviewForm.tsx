'use client';

import React, { useState, useRef } from 'react';
import { useLocale, useTranslations } from 'next-intl';
import Image from 'next/image';
import { RatingInput } from './RatingInput';
import {
  useCreateDraftReview,
  useUploadReviewPhotosToReview,
  usePublishReview,
  useUploadReviewPhotos,
} from '@/hooks/useReviews';
import type { CreateReviewRequest } from '@/types/review';

interface ReviewFormProps {
  // Entity data for creating review
  entityType: 'listing' | 'user' | 'storefront';
  entityId: number;
  storefrontId?: number;

  // Callbacks
  onSuccess?: () => void;
  onCancel: () => void;

  // Optional override for old behavior
  legacyOnSubmit?: (data: {
    rating: number;
    comment: string;
    pros?: string;
    cons?: string;
    photos?: string[];
  }) => Promise<void>;
}

export const ReviewForm: React.FC<ReviewFormProps> = ({
  entityType,
  entityId,
  storefrontId,
  onSuccess,
  onCancel,
  legacyOnSubmit,
}) => {
  const locale = useLocale();
  const t = useTranslations('reviews.form');
  const [rating, setRating] = useState(0);
  const [comment, setComment] = useState('');
  const [pros, setPros] = useState('');
  const [cons, setCons] = useState('');
  const [photos, setPhotos] = useState<File[]>([]);
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [dragActive, setDragActive] = useState(false);
  const [currentStep, setCurrentStep] = useState<
    'form' | 'uploading' | 'publishing'
  >('form');
  const [_createdReviewId, setCreatedReviewId] = useState<number | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  // Two-step process hooks
  const createDraftMutation = useCreateDraftReview();
  const uploadPhotosMutation = useUploadReviewPhotosToReview();
  const publishReviewMutation = usePublishReview();

  // Legacy hook for backward compatibility
  const legacyUploadPhotosMutation = useUploadReviewPhotos();

  const handlePhotoSelect = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);

    // Validate files
    const validFiles = files.filter((file) => {
      if (!['image/jpeg', 'image/png', 'image/webp'].includes(file.type)) {
        return false;
      }
      if (file.size > 5 * 1024 * 1024) {
        // 5MB
        return false;
      }
      return true;
    });

    if (validFiles.length + photos.length > 5) {
      setErrors((prev) => ({
        ...prev,
        photos: t('photos.maxFiles'),
      }));
      return;
    }

    setPhotos((prev) => [...prev, ...validFiles]);
    setErrors((prev) => ({ ...prev, photos: '' }));
  };

  const removePhoto = (index: number) => {
    setPhotos((prev) => prev.filter((_, i) => i !== index));
  };

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    const files = Array.from(e.dataTransfer.files).filter(
      (file) =>
        ['image/jpeg', 'image/png', 'image/webp'].includes(file.type) &&
        file.size <= 5 * 1024 * 1024
    );

    if (files.length + photos.length > 5) {
      setErrors((prev) => ({
        ...prev,
        photos: t('photos.maxFiles'),
      }));
      return;
    }

    setPhotos((prev) => [...prev, ...files]);
    setErrors((prev) => ({ ...prev, photos: '' }));
  };

  const validate = () => {
    const newErrors: Record<string, string> = {};

    if (rating === 0) {
      newErrors.rating = t('rating.error');
    }

    if (!comment.trim()) {
      newErrors.comment = t('comment.error');
    } else if (comment.length < 10) {
      newErrors.comment = t('comment.tooShort');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!validate()) {
      return;
    }

    console.log('🔥 ReviewForm handleSubmit:', {
      legacyOnSubmit: !!legacyOnSubmit,
    });

    // Support legacy behavior if legacyOnSubmit is provided
    if (legacyOnSubmit) {
      console.log('🔥 ИСПОЛЬЗУЕТСЯ LEGACY РЕЖИМ!');
      alert('LEGACY РЕЖИМ АКТИВЕН!');

      try {
        let uploadedPhotoUrls: string[] = [];

        // Upload photos if any (legacy way)
        if (photos.length > 0) {
          uploadedPhotoUrls =
            await legacyUploadPhotosMutation.mutateAsync(photos);
        }

        await legacyOnSubmit({
          rating,
          comment,
          pros: pros.trim() || undefined,
          cons: cons.trim() || undefined,
          photos: uploadedPhotoUrls.length > 0 ? uploadedPhotoUrls : undefined,
        });
      } catch (error) {
        console.error('Failed to submit review (legacy):', error);
      }
      return;
    }

    // New two-step process
    console.log('🚀 НОВЫЙ ДВУХЭТАПНЫЙ ПРОЦЕСС ЗАПУСКАЕТСЯ!');
    try {
      setCurrentStep('form');

      // Step 1: Create draft review
      const reviewData: CreateReviewRequest = {
        entity_type: entityType,
        entity_id: entityId,
        rating,
        comment,
        pros: pros.trim() || undefined,
        cons: cons.trim() || undefined,
        original_language: locale,
        ...(storefrontId && { storefront_id: storefrontId }),
      };

      const draftReview = await createDraftMutation.mutateAsync(reviewData);
      setCreatedReviewId(draftReview.id);

      // Step 2a: Upload photos if any
      if (photos.length > 0) {
        setCurrentStep('uploading');
        await uploadPhotosMutation.mutateAsync({
          reviewId: draftReview.id,
          files: photos,
        });
      }

      // Step 2b: Publish the review
      setCurrentStep('publishing');
      await publishReviewMutation.mutateAsync(draftReview.id);

      // Success - call callback
      onSuccess?.();
    } catch (error) {
      console.error('Failed to submit review:', error);
      setCurrentStep('form');
    }
  };

  // Progress indicator for two-step process
  const renderProgressIndicator = () => {
    if (legacyOnSubmit) return null; // Don't show for legacy mode

    const isProcessing =
      createDraftMutation.isPending ||
      uploadPhotosMutation.isPending ||
      publishReviewMutation.isPending;

    if (!isProcessing) return null;

    const steps = [
      {
        key: 'form',
        label: t('progress.creating'),
        active: currentStep === 'form' || createDraftMutation.isPending,
      },
      {
        key: 'uploading',
        label: t('progress.uploading'),
        active: currentStep === 'uploading' && photos.length > 0,
        skip: photos.length === 0,
      },
      {
        key: 'publishing',
        label: t('progress.publishing'),
        active: currentStep === 'publishing' || publishReviewMutation.isPending,
      },
    ].filter((step) => !step.skip);

    return (
      <div className="mb-6 p-4 bg-base-100 rounded-lg border border-base-200">
        <div className="flex items-center justify-between">
          {steps.map((step, index) => (
            <div key={step.key} className="flex items-center flex-1">
              <div
                className={`w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium
                ${
                  step.active
                    ? 'bg-primary text-primary-content'
                    : 'bg-base-200 text-base-content/60'
                }`}
              >
                {step.active ? (
                  <span className="loading loading-spinner loading-xs"></span>
                ) : (
                  index + 1
                )}
              </div>
              <span
                className={`ml-2 text-sm font-medium ${
                  step.active ? 'text-base-content' : 'text-base-content/60'
                }`}
              >
                {step.label}
              </span>
              {index < steps.length - 1 && (
                <div className="flex-1 mx-4 h-px bg-base-200"></div>
              )}
            </div>
          ))}
        </div>
      </div>
    );
  };

  return (
    <div className="space-y-6">
      {renderProgressIndicator()}

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Rating */}
        <div className="space-y-3 animate-in fade-in duration-500">
          <label className="block">
            <span className="text-sm font-medium text-base-content uppercase tracking-wider">
              {t('rating.label')}
              <span className="text-error ml-1">*</span>
            </span>
          </label>
          <div className="bg-base-100 p-4 rounded-lg border border-base-200 transition-all hover:border-primary/30">
            <RatingInput
              value={rating}
              onChange={setRating}
              size="lg"
              required
              error={errors.rating}
            />
          </div>
        </div>

        {/* Comment */}
        <div className="space-y-3 animate-in fade-in duration-500 delay-100">
          <label className="block">
            <span className="text-sm font-medium text-base-content uppercase tracking-wider">
              {t('comment.label')}
              <span className="text-error ml-1">*</span>
            </span>
          </label>
          <div className="relative">
            <textarea
              value={comment}
              onChange={(e) => setComment(e.target.value)}
              className={`w-full min-h-[120px] p-4 rounded-lg border ${
                errors.comment
                  ? 'border-error focus:ring-error'
                  : 'border-base-200 focus:border-primary focus:ring-primary'
              } bg-base-100 resize-none transition-all duration-200 focus:outline-none focus:ring-2`}
              placeholder={t('comment.placeholder')}
              maxLength={500}
            />
            <div className="absolute bottom-3 right-3 text-xs text-base-content/50">
              {comment.length}/500
            </div>
          </div>
          {errors.comment && (
            <p className="text-xs text-error mt-1">{errors.comment}</p>
          )}
        </div>

        {/* Pros & Cons */}
        <div className="grid md:grid-cols-2 gap-4 animate-in fade-in duration-500 delay-200">
          {/* Pros */}
          <div className="space-y-2">
            <label className="flex items-center gap-2">
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
              <span className="text-sm font-medium text-base-content">
                {t('pros.label')}
              </span>
            </label>
            <textarea
              value={pros}
              onChange={(e) => setPros(e.target.value)}
              className="w-full min-h-[80px] p-3 rounded-lg border border-base-200 bg-base-100 
                     resize-none transition-all duration-200 
                     focus:outline-none focus:ring-2 focus:ring-success/20 focus:border-success
                     hover:border-success/30"
              placeholder={t('pros.placeholder')}
              maxLength={300}
            />
          </div>

          {/* Cons */}
          <div className="space-y-2">
            <label className="flex items-center gap-2">
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
              <span className="text-sm font-medium text-base-content">
                {t('cons.label')}
              </span>
            </label>
            <textarea
              value={cons}
              onChange={(e) => setCons(e.target.value)}
              className="w-full min-h-[80px] p-3 rounded-lg border border-base-200 bg-base-100 
                     resize-none transition-all duration-200 
                     focus:outline-none focus:ring-2 focus:ring-warning/20 focus:border-warning
                     hover:border-warning/30"
              placeholder={t('cons.placeholder')}
              maxLength={300}
            />
          </div>
        </div>

        {/* Photos */}
        <div className="space-y-3 animate-in fade-in duration-500 delay-300">
          <div className="flex items-center justify-between">
            <label className="flex items-center gap-2">
              <div className="w-8 h-8 rounded-md bg-primary/10 flex items-center justify-center">
                <svg
                  className="w-4 h-4 text-primary"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
                  />
                </svg>
              </div>
              <span className="text-sm font-medium text-base-content">
                {t('photos.label')}
              </span>
            </label>
            <span className="text-xs text-base-content/60 font-medium">
              {photos.length}/5
            </span>
          </div>

          {/* Photo Grid */}
          {photos.length > 0 && (
            <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-5 gap-2 mb-4">
              {photos.map((photo, index) => (
                <div
                  key={index}
                  className="relative group animate-in fade-in zoom-in-95 duration-300"
                >
                  <Image
                    src={URL.createObjectURL(photo)}
                    alt={`Preview ${index + 1}`}
                    width={80}
                    height={80}
                    className="w-full h-20 object-cover rounded-md transition-all group-hover:brightness-90"
                  />
                  <button
                    type="button"
                    onClick={() => removePhoto(index)}
                    className="absolute -top-1 -right-1 w-5 h-5 bg-error text-error-content 
                           rounded-full opacity-0 group-hover:opacity-100 transition-opacity
                           flex items-center justify-center text-xs hover:scale-110 shadow-sm"
                  >
                    ✕
                  </button>
                </div>
              ))}
            </div>
          )}

          {/* Drag & Drop Area */}
          <div
            className={`relative border-2 border-dashed rounded-lg p-6 text-center transition-all duration-300
            ${
              dragActive
                ? 'border-primary bg-primary/5'
                : photos.length >= 5
                  ? 'border-base-200 bg-base-200/50 cursor-not-allowed'
                  : 'border-base-300 hover:border-primary/50 hover:bg-base-100 cursor-pointer'
            }`}
            onDragEnter={handleDrag}
            onDragLeave={handleDrag}
            onDragOver={handleDrag}
            onDrop={handleDrop}
            onClick={() =>
              !dragActive && photos.length < 5 && fileInputRef.current?.click()
            }
          >
            <input
              ref={fileInputRef}
              type="file"
              multiple
              accept="image/jpeg,image/png,image/webp"
              onChange={handlePhotoSelect}
              className="hidden"
              disabled={photos.length >= 5}
            />

            <svg
              className="w-10 h-10 mx-auto mb-2 text-base-content/30"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1.5}
                d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
              />
            </svg>

            <p className="text-sm font-medium text-base-content/70 mb-1">
              {photos.length >= 5
                ? t('photos.limitReached')
                : t('photos.dragText')}
            </p>
            <p className="text-xs text-base-content/50">
              {t('photos.fileTypes')}
            </p>
          </div>

          {errors.photos && (
            <p className="text-xs text-error mt-1">{errors.photos}</p>
          )}
        </div>

        {/* Actions */}
        <div className="flex gap-3 pt-4 animate-in fade-in duration-500 delay-400">
          <button
            type="submit"
            disabled={
              createDraftMutation.isPending ||
              uploadPhotosMutation.isPending ||
              publishReviewMutation.isPending ||
              legacyUploadPhotosMutation.isPending
            }
            className="btn btn-primary flex-1 h-11 rounded-lg font-medium
                   hover:shadow-md transition-all duration-200
                   disabled:shadow-none"
          >
            {createDraftMutation.isPending ||
            uploadPhotosMutation.isPending ||
            publishReviewMutation.isPending ||
            legacyUploadPhotosMutation.isPending ? (
              <>
                <span className="loading loading-spinner loading-sm"></span>
                {currentStep === 'uploading'
                  ? t('progress.uploading')
                  : currentStep === 'publishing'
                    ? t('progress.publishing')
                    : t('progress.creating')}
              </>
            ) : (
              <>
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
                    d="M5 13l4 4L19 7"
                  />
                </svg>
                {t('submit')}
              </>
            )}
          </button>
          <button
            type="button"
            onClick={onCancel}
            disabled={
              createDraftMutation.isPending ||
              uploadPhotosMutation.isPending ||
              publishReviewMutation.isPending ||
              legacyUploadPhotosMutation.isPending
            }
            className="btn btn-ghost h-11 px-6 rounded-lg font-medium
                   hover:bg-base-200 transition-all duration-200"
          >
            {t('cancel')}
          </button>
        </div>
      </form>
    </div>
  );
};
