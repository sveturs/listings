'use client';

import React, { useState, useEffect, useCallback, useRef } from 'react';
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
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [touchStart, setTouchStart] = useState<number | null>(null);
  const [touchEnd, setTouchEnd] = useState<number | null>(null);
  const [isMobile, setIsMobile] = useState(false);
  const galleryRef = useRef<HTMLDivElement>(null);

  // Minimum swipe distance (in pixels)
  const minSwipeDistance = 50;

  // Navigation functions
  const goToNext = useCallback(() => {
    setCurrentIndex((prev) => (prev + 1) % images.length);
  }, [images.length]);

  const goToPrevious = useCallback(() => {
    setCurrentIndex((prev) => (prev - 1 + images.length) % images.length);
  }, [images.length]);

  // Toggle fullscreen
  const toggleFullscreen = useCallback(async () => {
    if (!document.fullscreenElement) {
      try {
        await galleryRef.current?.requestFullscreen();
        setIsFullscreen(true);
      } catch (err) {
        console.error('Error attempting to enable fullscreen:', err);
      }
    } else {
      try {
        await document.exitFullscreen();
        setIsFullscreen(false);
      } catch (err) {
        console.error('Error attempting to exit fullscreen:', err);
      }
    }
  }, []);

  // Handle touch start
  const onTouchStart = (e: React.TouchEvent) => {
    setTouchEnd(null);
    setTouchStart(e.targetTouches[0].clientX);
  };

  // Handle touch move
  const onTouchMove = (e: React.TouchEvent) => {
    setTouchEnd(e.targetTouches[0].clientX);
  };

  // Handle touch end
  const onTouchEnd = () => {
    if (!touchStart || !touchEnd) return;

    const distance = touchStart - touchEnd;
    const isLeftSwipe = distance > minSwipeDistance;
    const isRightSwipe = distance < -minSwipeDistance;

    if (isLeftSwipe && images.length > 1) {
      goToNext();
    }
    if (isRightSwipe && images.length > 1) {
      goToPrevious();
    }
  };

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
          if (document.fullscreenElement) {
            document.exitFullscreen();
          } else {
            onClose();
          }
          break;
        case 'ArrowLeft':
          goToPrevious();
          break;
        case 'ArrowRight':
          goToNext();
          break;
        case 'f':
        case 'F':
          toggleFullscreen();
          break;
      }
    };

    document.addEventListener('keydown', handleKeyPress);
    document.body.style.overflow = 'hidden';

    return () => {
      document.removeEventListener('keydown', handleKeyPress);
      document.body.style.overflow = '';
    };
  }, [isOpen, onClose, goToNext, goToPrevious, toggleFullscreen]);

  // Handle wheel navigation
  useEffect(() => {
    if (!isOpen) return;

    const handleWheel = (e: WheelEvent) => {
      e.preventDefault();
      if (images.length <= 1) return;

      if (e.deltaY > 0) {
        goToNext();
      } else {
        goToPrevious();
      }
    };

    const container = galleryRef.current;
    if (container) {
      container.addEventListener('wheel', handleWheel, { passive: false });
    }

    return () => {
      if (container) {
        container.removeEventListener('wheel', handleWheel);
      }
    };
  }, [isOpen, goToNext, goToPrevious, images.length]);

  // Handle fullscreen change
  useEffect(() => {
    const handleFullscreenChange = () => {
      setIsFullscreen(!!document.fullscreenElement);
    };

    document.addEventListener('fullscreenchange', handleFullscreenChange);
    return () => {
      document.removeEventListener('fullscreenchange', handleFullscreenChange);
    };
  }, []);

  // Check if mobile
  useEffect(() => {
    const checkMobile = () => {
      setIsMobile(window.innerWidth <= 768);
    };

    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, []);

  if (!isOpen) return null;

  return (
    <div
      ref={galleryRef}
      className="fixed inset-0 z-[101] flex items-center justify-center bg-base-300/5 backdrop-blur-sm"
      onTouchStart={onTouchStart}
      onTouchMove={onTouchMove}
      onTouchEnd={onTouchEnd}
    >
      {/* Background overlay */}
      <div
        className="absolute inset-0 cursor-pointer"
        onClick={onClose}
        aria-label={locale === 'ru' ? 'Закрыть галерею' : 'Close gallery'}
      />

      {/* Gallery container */}
      <div
        className={`relative ${isFullscreen ? 'w-screen h-screen' : 'max-w-5xl max-h-[90vh] w-full mx-4'} bg-base-100/10 backdrop-blur-md ${!isFullscreen && 'rounded-2xl'} shadow-2xl overflow-hidden`}
      >
        {/* Header with controls */}
        <div className="flex items-center justify-between p-4 bg-base-100/20 backdrop-blur-sm border-b border-base-200/20">
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

          <div className="flex items-center gap-2">
            <button
              onClick={toggleFullscreen}
              className="btn btn-sm btn-ghost btn-circle"
              aria-label={
                locale === 'ru' ? 'Полноэкранный режим' : 'Fullscreen mode'
              }
            >
              <svg
                className="w-5 h-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                {isFullscreen ? (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M9 9V4.5M9 9H4.5M9 9L3.75 3.75M9 15v4.5M9 15H4.5M9 15l-5.25 5.25M15 9h4.5M15 9V4.5M15 9l5.25-5.25M15 15h4.5M15 15v4.5m0-4.5l5.25 5.25"
                  />
                ) : (
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4"
                  />
                )}
              </svg>
            </button>
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
        </div>

        {/* Main image container */}
        <div
          className="relative bg-transparent p-6"
          style={{
            height: isFullscreen ? 'calc(100vh - 80px)' : 'calc(90vh - 160px)',
          }}
        >
          <div
            className="relative w-full h-full flex items-center justify-center cursor-pointer"
            onClick={toggleFullscreen}
          >
            <Image
              src={images[currentIndex]}
              alt={`Review photo ${currentIndex + 1}`}
              fill
              className="object-contain rounded-lg drop-shadow-2xl"
              sizes="(max-width: 1200px) 100vw, 1200px"
              priority
            />
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

        {/* Footer with thumbnails and close button - hide in fullscreen on mobile */}
        {(!isFullscreen || !isMobile) && (
          <div className="px-4 pt-4 pb-12 bg-transparent">
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
        )}
      </div>

      {/* Swipe indicator for mobile */}
      {images.length > 1 && isMobile && (
        <div className="absolute bottom-4 left-1/2 -translate-x-1/2 text-base-content/50 text-sm">
          {locale === 'ru' ? 'Свайпните для навигации' : 'Swipe to navigate'}
        </div>
      )}
    </div>
  );
};
