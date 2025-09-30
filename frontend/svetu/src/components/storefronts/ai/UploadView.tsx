'use client';

import React, { useRef, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateAIProduct } from '@/contexts/CreateAIProductContext';
import Image from 'next/image';
import {
  PhotoIcon,
  XMarkIcon,
  ArrowUpTrayIcon,
} from '@heroicons/react/24/outline';

export default function UploadView() {
  const t = useTranslations('storefronts');
  const { state, setImages, setView, setError } = useCreateAIProduct();
  const fileInputRef = useRef<HTMLInputElement>(null);
  const [isDragging, setIsDragging] = useState(false);

  const handleFileSelect = (files: FileList | null) => {
    if (!files || files.length === 0) return;

    const fileArray = Array.from(files);
    const validFiles = fileArray.filter((file) => {
      if (!file.type.startsWith('image/')) {
        setError(`${file.name} is not an image file`);
        return false;
      }
      if (file.size > 10 * 1024 * 1024) {
        // 10MB limit
        setError(`${file.name} is too large (max 10MB)`);
        return false;
      }
      return true;
    });

    if (validFiles.length === 0) return;

    // Limit to 10 images total
    const totalImages = state.imageFiles.length + validFiles.length;
    if (totalImages > 10) {
      setError('Maximum 10 images allowed');
      return;
    }

    // Create blob URLs for preview
    const newImageUrls = validFiles.map((file) => URL.createObjectURL(file));
    const allImageUrls = [...state.images, ...newImageUrls];
    const allImageFiles = [...state.imageFiles, ...validFiles];

    setImages(allImageUrls, allImageFiles);
    setError(null);
  };

  const handleRemoveImage = (index: number) => {
    const newImageUrls = state.images.filter((_, i) => i !== index);
    const newImageFiles = state.imageFiles.filter((_, i) => i !== index);

    // Revoke the blob URL to free memory
    URL.revokeObjectURL(state.images[index]);

    setImages(newImageUrls, newImageFiles);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setIsDragging(false);
    handleFileSelect(e.dataTransfer.files);
  };

  const handleContinue = () => {
    if (state.imageFiles.length === 0) {
      setError('Please upload at least one image');
      return;
    }
    setView('process');
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold text-base-content mb-2">
          {t('uploadProductImages') || 'Upload Product Images'}
        </h2>
        <p className="text-base-content/70">
          {t('uploadImagesDescription') ||
            'Upload 1-10 high-quality images of your product. AI will analyze them to create a perfect listing.'}
        </p>
      </div>

      {/* Upload Area */}
      <div
        className={`border-2 border-dashed rounded-lg p-8 text-center cursor-pointer transition-colors ${
          isDragging
            ? 'border-primary bg-primary/5'
            : 'border-base-300 hover:border-primary/50 hover:bg-base-200'
        }`}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => fileInputRef.current?.click()}
      >
        <input
          ref={fileInputRef}
          type="file"
          accept="image/*"
          multiple
          className="hidden"
          onChange={(e) => handleFileSelect(e.target.files)}
        />

        <div className="flex flex-col items-center gap-4">
          <div className="w-16 h-16 rounded-full bg-primary/10 flex items-center justify-center">
            <ArrowUpTrayIcon className="w-8 h-8 text-primary" />
          </div>

          <div>
            <p className="text-lg font-semibold text-base-content mb-1">
              {t('dragDropImages') || 'Drag & drop images here'}
            </p>
            <p className="text-base-content/60 text-sm">
              {t('orClickToSelect') || 'or click to select files'}
            </p>
          </div>

          <div className="text-sm text-base-content/50">
            {t('imageRequirements') ||
              'PNG, JPG, WebP up to 10MB • Maximum 10 images'}
          </div>
        </div>
      </div>

      {/* Image Preview Grid */}
      {state.images.length > 0 && (
        <div>
          <h3 className="text-lg font-semibold text-base-content mb-3">
            {t('selectedImages') || 'Selected Images'} ({state.images.length}
            /10)
          </h3>
          <div className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4">
            {state.images.map((imageUrl, index) => (
              <div key={index} className="relative group aspect-square">
                <Image
                  src={imageUrl}
                  alt={`Product image ${index + 1}`}
                  fill
                  className="rounded-lg object-cover"
                />
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    handleRemoveImage(index);
                  }}
                  className="absolute top-2 right-2 p-1.5 bg-error text-error-content rounded-full opacity-0 group-hover:opacity-100 transition-opacity shadow-lg hover:scale-110"
                >
                  <XMarkIcon className="w-4 h-4" />
                </button>
                {index === 0 && (
                  <div className="absolute bottom-2 left-2 px-2 py-1 bg-primary text-primary-content text-xs rounded-md font-semibold">
                    Main
                  </div>
                )}
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Tips */}
      <div className="bg-info/10 border border-info/20 rounded-lg p-4">
        <h4 className="font-semibold text-info mb-2 flex items-center gap-2">
          <PhotoIcon className="w-5 h-5" />
          {t('photoTips') || 'Photo Tips'}
        </h4>
        <ul className="text-sm text-base-content/70 space-y-1 list-disc list-inside">
          <li>{t('photoTip1') || 'Use good lighting and clear focus'}</li>
          <li>{t('photoTip2') || 'Show product from different angles'}</li>
          <li>{t('photoTip3') || 'Include close-ups of important details'}</li>
          <li>
            {t('photoTip4') || 'First image will be used as main product photo'}
          </li>
        </ul>
      </div>

      {/* Actions */}
      <div className="flex justify-end gap-3">
        <button
          onClick={handleContinue}
          disabled={state.imageFiles.length === 0}
          className="btn btn-primary px-8"
        >
          {t('continueToAI') || 'Continue to AI Processing'}
          <svg
            className="w-5 h-5 ml-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M13 7l5 5m0 0l-5 5m5-5H6"
            />
          </svg>
        </button>
      </div>
    </div>
  );
}
