'use client';

import React, { useState, useEffect, useCallback, useRef } from 'react';
import Image from 'next/image';
import { useTranslations } from 'next-intl';

export interface ImageItem {
  src: string;
  alt?: string;
  title?: string;
  description?: string;
}

export interface ImageGalleryProps {
  images: string[] | ImageItem[];
  initialIndex?: number;
  isOpen: boolean;
  onClose: () => void;
  onImageChange?: (index: number) => void;
  showThumbnails?: boolean;
  showFullscreenButton?: boolean;
  showNavigationButtons?: boolean;
  showCounter?: boolean;
  showCloseButton?: boolean;
  enableKeyboardNavigation?: boolean;
  enableWheelNavigation?: boolean;
  enableSwipeNavigation?: boolean;
  minSwipeDistance?: number;
  overlayClassName?: string;
  containerClassName?: string;
  imageClassName?: string;
  thumbnailClassName?: string;
  translationPrefix?: string;
  customHeader?: React.ReactNode;
  customFooter?: React.ReactNode;
  renderImage?: (image: ImageItem, index: number) => React.ReactNode;
  renderThumbnail?: (image: ImageItem, index: number) => React.ReactNode;
  imageAltPrefix?: string;
  imageSizes?: string;
  thumbnailWidth?: number;
  thumbnailHeight?: number;
  zIndex?: number;
  closeOnOverlayClick?: boolean;
  showSwipeIndicator?: boolean;
  mobileBreakpoint?: number;
}

const normalizeImages = (images: string[] | ImageItem[]): ImageItem[] => {
  return images.map((image, index) => {
    if (typeof image === 'string') {
      return { src: image, alt: `Image ${index + 1}` };
    }
    return image;
  });
};

export const ImageGallery: React.FC<ImageGalleryProps> = ({
  images,
  initialIndex = 0,
  isOpen,
  onClose,
  onImageChange,
  showThumbnails = true,
  showFullscreenButton = true,
  showNavigationButtons = true,
  showCounter = true,
  showCloseButton = true,
  enableKeyboardNavigation = true,
  enableWheelNavigation = true,
  enableSwipeNavigation = true,
  minSwipeDistance = 50,
  overlayClassName = 'bg-base-300/5 backdrop-blur-sm',
  containerClassName = '',
  imageClassName = 'object-contain rounded-lg drop-shadow-2xl',
  thumbnailClassName = '',
  translationPrefix = 'reviews.gallery',
  customHeader,
  customFooter,
  renderImage,
  renderThumbnail,
  imageAltPrefix,
  imageSizes = '(max-width: 1200px) 100vw, 1200px',
  thumbnailWidth = 64,
  thumbnailHeight = 64,
  zIndex = 101,
  closeOnOverlayClick = true,
  showSwipeIndicator = true,
  mobileBreakpoint = 768,
}) => {
  const t = useTranslations(translationPrefix);
  const normalizedImages = normalizeImages(images);
  const [currentIndex, setCurrentIndex] = useState(initialIndex);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [touchStart, setTouchStart] = useState<number | null>(null);
  const [touchEnd, setTouchEnd] = useState<number | null>(null);
  const [isMobile, setIsMobile] = useState(false);
  const galleryRef = useRef<HTMLDivElement>(null);

  // Navigation functions
  const goToNext = useCallback(() => {
    const newIndex = (currentIndex + 1) % normalizedImages.length;
    setCurrentIndex(newIndex);
    onImageChange?.(newIndex);
  }, [currentIndex, normalizedImages.length, onImageChange]);

  const goToPrevious = useCallback(() => {
    const newIndex =
      (currentIndex - 1 + normalizedImages.length) % normalizedImages.length;
    setCurrentIndex(newIndex);
    onImageChange?.(newIndex);
  }, [currentIndex, normalizedImages.length, onImageChange]);

  const goToIndex = useCallback(
    (index: number) => {
      setCurrentIndex(index);
      onImageChange?.(index);
    },
    [onImageChange]
  );

  // Toggle fullscreen
  const toggleFullscreen = useCallback(async () => {
    if (!showFullscreenButton) return;

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
  }, [showFullscreenButton]);

  // Handle touch events
  const onTouchStart = (e: React.TouchEvent) => {
    if (!enableSwipeNavigation) return;
    setTouchEnd(null);
    setTouchStart(e.targetTouches[0].clientX);
  };

  const onTouchMove = (e: React.TouchEvent) => {
    if (!enableSwipeNavigation) return;
    setTouchEnd(e.targetTouches[0].clientX);
  };

  const onTouchEnd = () => {
    if (!enableSwipeNavigation || !touchStart || !touchEnd) return;

    const distance = touchStart - touchEnd;
    const isLeftSwipe = distance > minSwipeDistance;
    const isRightSwipe = distance < -minSwipeDistance;

    if (isLeftSwipe && normalizedImages.length > 1) {
      goToNext();
    }
    if (isRightSwipe && normalizedImages.length > 1) {
      goToPrevious();
    }
  };

  // Update current index when initialIndex changes
  useEffect(() => {
    setCurrentIndex(initialIndex);
  }, [initialIndex]);

  // Handle keyboard navigation
  useEffect(() => {
    if (!isOpen || !enableKeyboardNavigation) return;

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
  }, [
    isOpen,
    enableKeyboardNavigation,
    onClose,
    goToNext,
    goToPrevious,
    toggleFullscreen,
  ]);

  // Handle wheel navigation
  useEffect(() => {
    if (!isOpen || !enableWheelNavigation) return;

    const handleWheel = (e: WheelEvent) => {
      e.preventDefault();
      if (normalizedImages.length <= 1) return;

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
  }, [
    isOpen,
    enableWheelNavigation,
    goToNext,
    goToPrevious,
    normalizedImages.length,
  ]);

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
      setIsMobile(window.innerWidth <= mobileBreakpoint);
    };

    checkMobile();
    window.addEventListener('resize', checkMobile);
    return () => window.removeEventListener('resize', checkMobile);
  }, [mobileBreakpoint]);

  if (!isOpen) return null;

  const currentImage = normalizedImages[currentIndex];
  const imageAlt = imageAltPrefix
    ? `${imageAltPrefix} ${currentIndex + 1}`
    : currentImage.alt || `Image ${currentIndex + 1}`;

  return (
    <div
      ref={galleryRef}
      className={`fixed inset-0 z-[${zIndex}] flex items-center justify-center ${overlayClassName}`}
      onTouchStart={onTouchStart}
      onTouchMove={onTouchMove}
      onTouchEnd={onTouchEnd}
    >
      {/* Background overlay */}
      <div
        className="absolute inset-0 cursor-pointer"
        onClick={closeOnOverlayClick ? onClose : undefined}
        aria-label={t('closeGallery')}
      />

      {/* Gallery container */}
      <div
        className={`relative ${isFullscreen ? 'w-screen h-screen' : 'max-w-5xl max-h-[90vh] w-full mx-4'} bg-base-100/10 backdrop-blur-md ${!isFullscreen && 'rounded-2xl'} shadow-2xl overflow-hidden ${containerClassName}`}
      >
        {/* Header */}
        {customHeader || (
          <div className="flex items-center justify-between p-4 bg-base-100/20 backdrop-blur-sm border-b border-base-200/20">
            <div className="flex items-center gap-4">
              {showCounter && (
                <div className="text-base-content text-sm font-medium">
                  {t('photo')} {currentIndex + 1} {t('of')}{' '}
                  {normalizedImages.length}
                </div>
              )}
              {showNavigationButtons && normalizedImages.length > 1 && (
                <div className="flex gap-1">
                  <button
                    onClick={goToPrevious}
                    className="btn btn-sm btn-ghost btn-circle"
                    aria-label={t('previousImage')}
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
                    aria-label={t('nextImage')}
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
              {showFullscreenButton && (
                <button
                  onClick={toggleFullscreen}
                  className="btn btn-sm btn-ghost btn-circle"
                  aria-label={t('fullscreenMode')}
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
              )}
              {showCloseButton && (
                <button
                  onClick={onClose}
                  className="btn btn-sm btn-ghost btn-circle"
                  aria-label={t('closeGallery')}
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
              )}
            </div>
          </div>
        )}

        {/* Main image container */}
        <div
          className="relative bg-transparent p-6"
          style={{
            height: isFullscreen ? 'calc(100vh - 80px)' : 'calc(90vh - 160px)',
          }}
        >
          {/* Image info overlay */}
          {(currentImage.title || currentImage.description) && (
            <div className="absolute top-8 left-8 right-8 z-10 bg-base-100/80 backdrop-blur-sm rounded-lg p-4">
              {currentImage.title && (
                <h3 className="text-lg font-semibold text-base-content mb-1">
                  {currentImage.title}
                </h3>
              )}
              {currentImage.description && (
                <p className="text-sm text-base-content/80">
                  {currentImage.description}
                </p>
              )}
            </div>
          )}

          <div
            className="relative w-full h-full flex items-center justify-center cursor-pointer"
            onClick={showFullscreenButton ? toggleFullscreen : undefined}
          >
            {renderImage ? (
              renderImage(currentImage, currentIndex)
            ) : (
              <Image
                src={currentImage.src}
                alt={imageAlt}
                fill
                className={imageClassName}
                sizes={imageSizes}
                priority
              />
            )}
          </div>

          {/* Large navigation arrows */}
          {showNavigationButtons && normalizedImages.length > 1 && (
            <>
              <button
                onClick={goToPrevious}
                className="absolute left-4 top-1/2 -translate-y-1/2 p-4 rounded-full bg-base-100/90 hover:bg-base-200/90 text-base-content shadow-lg transition-all duration-200 hover:scale-110 border border-base-300 backdrop-blur-sm"
                aria-label={t('previousImage')}
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
                aria-label={t('nextImage')}
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

        {/* Footer */}
        {customFooter ||
          ((!isFullscreen || !isMobile) && (
            <div className="px-4 pt-4 pb-12 bg-transparent">
              {showThumbnails && normalizedImages.length > 1 && (
                <div className="flex justify-center gap-2 mb-16 max-w-full overflow-x-auto pb-2">
                  {normalizedImages.map((image, index) => (
                    <button
                      key={index}
                      onClick={() => goToIndex(index)}
                      className={`relative flex-shrink-0 w-16 h-16 rounded-lg overflow-hidden transition-all duration-200 border-2 ${thumbnailClassName}
                        ${
                          currentIndex === index
                            ? 'border-primary scale-105'
                            : 'border-base-300 hover:border-primary/50 hover:scale-105 opacity-70 hover:opacity-100'
                        }`}
                    >
                      {renderThumbnail ? (
                        renderThumbnail(image, index)
                      ) : (
                        <Image
                          src={image.src}
                          alt={`Thumbnail ${index + 1}`}
                          width={thumbnailWidth}
                          height={thumbnailHeight}
                          className="w-full h-full object-cover"
                        />
                      )}
                    </button>
                  ))}
                </div>
              )}

              {showCloseButton && (
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
                    {t('close')}
                  </button>
                </div>
              )}
            </div>
          ))}
      </div>

      {/* Swipe indicator for mobile */}
      {showSwipeIndicator &&
        normalizedImages.length > 1 &&
        isMobile &&
        enableSwipeNavigation && (
          <div className="absolute bottom-4 left-1/2 -translate-x-1/2 text-base-content/50 text-sm">
            {t('swipeToNavigate')}
          </div>
        )}
    </div>
  );
};
