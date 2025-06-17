'use client';

import React, { useState, useEffect, useCallback } from 'react';
import Image from 'next/image';
import { useLocale } from 'next-intl';

interface ImageGalleryProps {
  images: string[];
  initialIndex?: number;
  isOpen: boolean;
  onClose: () => void;
}

export const ImageGallery: React.FC<ImageGalleryProps> = ({
  images,
  initialIndex = 0,
  isOpen,
  onClose,
}) => {
  const locale = useLocale();
  const [currentIndex, setCurrentIndex] = useState(initialIndex);

  // Navigation functions
  const goToNext = useCallback(() => {
    setCurrentIndex((prev) => (prev + 1) % images.length);
  }, [images.length]);

  const goToPrevious = useCallback(() => {
    setCurrentIndex((prev) => (prev - 1 + images.length) % images.length);
  }, [images.length]);

  // Update current index when initialIndex changes
  useEffect(() => {
    setCurrentIndex(initialIndex);
  }, [initialIndex]);

  // Handle keyboard navigation
  useEffect(() => {
    if (!isOpen) return;

    const handleKeyPress = (e: KeyboardEvent) => {
      switch (e.key) {
        case 'Escape':
          onClose();
          break;
        case 'ArrowLeft':
          goToPrevious();
          break;
        case 'ArrowRight':
          goToNext();
          break;
      }
    };

    document.addEventListener('keydown', handleKeyPress);
    document.body.style.overflow = 'hidden';

    return () => {
      document.removeEventListener('keydown', handleKeyPress);
      document.body.style.overflow = '';
    };
  }, [isOpen, onClose, goToNext, goToPrevious]);

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-[101] flex items-center justify-center bg-base-300/20 backdrop-blur-sm">
      {/* Background overlay */}
      <div
        className="absolute inset-0 cursor-pointer"
        onClick={onClose}
        aria-label={locale === 'ru' ? 'Закрыть галерею' : 'Close gallery'}
      />

      {/* Gallery container */}
      <div className="relative max-w-5xl max-h-[90vh] w-full mx-4 bg-base-100 rounded-2xl shadow-2xl overflow-hidden">
        {/* Header with controls */}
        <div className="flex items-center justify-between p-4 bg-base-100 border-b border-base-200">
          <div className="flex items-center gap-4">
            <div className="text-base-content text-sm font-medium">
              {locale === 'ru' ? 'Фотография' : 'Photo'} {currentIndex + 1}{' '}
              {locale === 'ru' ? 'из' : 'of'} {images.length}
            </div>
            {images.length > 1 && (
              <div className="flex gap-1">
                <button
                  onClick={goToPrevious}
                  className="btn btn-sm btn-ghost btn-circle"
                  aria-label={
                    locale === 'ru'
                      ? 'Предыдущее изображение'
                      : 'Previous image'
                  }
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
                <button
                  onClick={goToNext}
                  className="btn btn-sm btn-ghost btn-circle"
                  aria-label={
                    locale === 'ru' ? 'Следующее изображение' : 'Next image'
                  }
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
            )}
          </div>

          <button
            onClick={onClose}
            className="btn btn-sm btn-ghost btn-circle"
            aria-label={locale === 'ru' ? 'Закрыть галерею' : 'Close gallery'}
          >
            <svg
              className="w-5 h-5"
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

        {/* Main image container */}
        <div
          className="relative bg-base-200/10 backdrop-blur-md p-6"
          style={{ height: 'calc(90vh - 160px)' }}
        >
          <div className="relative w-full h-full bg-base-100 rounded-xl shadow-inner">
            <div className="absolute inset-2 flex items-center justify-center">
              <div className="relative w-full h-full">
                <Image
                  src={images[currentIndex]}
                  alt={`Review photo ${currentIndex + 1}`}
                  fill
                  className="object-contain rounded-lg"
                  sizes="(max-width: 1200px) 100vw, 1200px"
                  priority
                />
              </div>
            </div>
          </div>

          {/* Large navigation arrows for easier clicking */}
          {images.length > 1 && (
            <>
              <button
                onClick={goToPrevious}
                className="absolute left-4 top-1/2 -translate-y-1/2 p-4 rounded-full bg-base-100/90 hover:bg-base-200/90 text-base-content shadow-lg transition-all duration-200 hover:scale-110 border border-base-300 backdrop-blur-sm"
                aria-label={
                  locale === 'ru' ? 'Предыдущее изображение' : 'Previous image'
                }
              >
                <svg
                  className="w-8 h-8"
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

              <button
                onClick={goToNext}
                className="absolute right-4 top-1/2 -translate-y-1/2 p-4 rounded-full bg-base-100/90 hover:bg-base-200/90 text-base-content shadow-lg transition-all duration-200 hover:scale-110 border border-base-300 backdrop-blur-sm"
                aria-label={
                  locale === 'ru' ? 'Следующее изображение' : 'Next image'
                }
              >
                <svg
                  className="w-8 h-8"
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
            </>
          )}
        </div>

        {/* Footer with thumbnails and close button */}
        <div className="px-4 pt-4 pb-12 bg-base-100 border-t border-base-200">
          {images.length > 1 && (
            <div className="flex justify-center gap-2 mb-16 max-w-full overflow-x-auto pb-2">
              {images.map((image, index) => (
                <button
                  key={index}
                  onClick={() => setCurrentIndex(index)}
                  className={`relative flex-shrink-0 w-16 h-16 rounded-lg overflow-hidden transition-all duration-200 border-2
                    ${
                      currentIndex === index
                        ? 'border-primary scale-105'
                        : 'border-base-300 hover:border-primary/50 hover:scale-105 opacity-70 hover:opacity-100'
                    }`}
                >
                  <Image
                    src={image}
                    alt={`Thumbnail ${index + 1}`}
                    width={64}
                    height={64}
                    className="w-full h-full object-cover"
                  />
                </button>
              ))}
            </div>
          )}

          <div className="flex justify-center">
            <button onClick={onClose} className="btn btn-primary gap-2">
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
              {locale === 'ru' ? 'Закрыть' : 'Close'}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
};
