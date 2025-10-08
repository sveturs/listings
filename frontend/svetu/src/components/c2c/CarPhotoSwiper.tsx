'use client';

import React, { useState, useRef, useEffect } from 'react';
import { ChevronLeft, ChevronRight, Car, X, Maximize2 } from 'lucide-react';
import { OptimizedCarImage } from '@/components/common/OptimizedCarImage';

interface CarPhotoSwiperProps {
  images?: Array<{
    id?: number;
    url?: string;
    thumbnail_url?: string;
  }>;
  title: string;
  showThumbnails?: boolean;
  showFullscreenButton?: boolean;
  autoPlay?: boolean;
  autoPlayInterval?: number;
}

/**
 * Компонент свайпера фотографий автомобиля с поддержкой touch-жестов
 */
export const CarPhotoSwiper: React.FC<CarPhotoSwiperProps> = ({
  images = [],
  title,
  showThumbnails = true,
  showFullscreenButton = true,
  autoPlay = false,
  autoPlayInterval = 5000,
}) => {
  const [currentIndex, setCurrentIndex] = useState(0);
  const [isFullscreen, setIsFullscreen] = useState(false);
  const [isDragging, setIsDragging] = useState(false);
  const [startX, setStartX] = useState(0);
  const [translateX, setTranslateX] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);
  const autoPlayRef = useRef<NodeJS.Timeout | undefined>(undefined);

  // Автопроигрывание
  useEffect(() => {
    if (autoPlay && images.length > 1) {
      autoPlayRef.current = setInterval(() => {
        setCurrentIndex((prev) => (prev + 1) % images.length);
      }, autoPlayInterval);
    }

    return () => {
      if (autoPlayRef.current) {
        clearInterval(autoPlayRef.current);
      }
    };
  }, [autoPlay, autoPlayInterval, images.length]);

  // Обработка touch-событий
  const handleTouchStart = (e: React.TouchEvent) => {
    setIsDragging(true);
    setStartX(e.touches[0].clientX);

    // Останавливаем автопроигрывание при взаимодействии
    if (autoPlayRef.current) {
      clearInterval(autoPlayRef.current);
    }
  };

  const handleTouchMove = (e: React.TouchEvent) => {
    if (!isDragging) return;

    const currentX = e.touches[0].clientX;
    const diff = currentX - startX;
    setTranslateX(diff);
  };

  const handleTouchEnd = () => {
    if (!isDragging) return;

    const threshold = 50; // Минимальное расстояние свайпа

    if (Math.abs(translateX) > threshold) {
      if (translateX > 0) {
        // Свайп вправо - предыдущее фото
        handlePrevious();
      } else {
        // Свайп влево - следующее фото
        handleNext();
      }
    }

    setIsDragging(false);
    setTranslateX(0);

    // Возобновляем автопроигрывание
    if (autoPlay && images.length > 1) {
      autoPlayRef.current = setInterval(() => {
        setCurrentIndex((prev) => (prev + 1) % images.length);
      }, autoPlayInterval);
    }
  };

  // Обработка мышки для десктопа
  const handleMouseDown = (e: React.MouseEvent) => {
    setIsDragging(true);
    setStartX(e.clientX);
    e.preventDefault();
  };

  const handleMouseMove = (e: React.MouseEvent) => {
    if (!isDragging) return;

    const currentX = e.clientX;
    const diff = currentX - startX;
    setTranslateX(diff);
  };

  const handleMouseUp = () => {
    handleTouchEnd();
  };

  const handleMouseLeave = () => {
    if (isDragging) {
      handleTouchEnd();
    }
  };

  const handlePrevious = () => {
    setCurrentIndex((prev) => (prev - 1 + images.length) % images.length);
  };

  const handleNext = () => {
    setCurrentIndex((prev) => (prev + 1) % images.length);
  };

  const handleThumbnailClick = (index: number) => {
    setCurrentIndex(index);
  };

  const toggleFullscreen = () => {
    setIsFullscreen(!isFullscreen);
  };

  if (!images || images.length === 0) {
    return (
      <div className="w-full h-full bg-base-200 flex items-center justify-center rounded-lg">
        <Car className="w-16 h-16 text-base-content/30" />
      </div>
    );
  }

  const currentImage = images[currentIndex];

  return (
    <>
      <div className="relative w-full h-full">
        {/* Main image container */}
        <div
          ref={containerRef}
          className="relative w-full h-full overflow-hidden rounded-lg cursor-grab active:cursor-grabbing"
          onTouchStart={handleTouchStart}
          onTouchMove={handleTouchMove}
          onTouchEnd={handleTouchEnd}
          onMouseDown={handleMouseDown}
          onMouseMove={handleMouseMove}
          onMouseUp={handleMouseUp}
          onMouseLeave={handleMouseLeave}
        >
          <div
            className="relative w-full h-full transition-transform duration-300 ease-out"
            style={{
              transform: `translateX(${isDragging ? translateX : 0}px)`,
            }}
          >
            <OptimizedCarImage
              src={currentImage.url || currentImage.thumbnail_url}
              alt={`${title} - фото ${currentIndex + 1}`}
              priority={currentIndex === 0}
            />
          </div>

          {/* Navigation arrows for desktop */}
          {images.length > 1 && (
            <>
              <button
                onClick={handlePrevious}
                className="absolute left-2 top-1/2 -translate-y-1/2 btn btn-circle btn-sm bg-base-100/80 hover:bg-base-100 backdrop-blur hidden lg:flex"
                aria-label="Previous photo"
              >
                <ChevronLeft className="w-5 h-5" />
              </button>
              <button
                onClick={handleNext}
                className="absolute right-2 top-1/2 -translate-y-1/2 btn btn-circle btn-sm bg-base-100/80 hover:bg-base-100 backdrop-blur hidden lg:flex"
                aria-label="Next photo"
              >
                <ChevronRight className="w-5 h-5" />
              </button>
            </>
          )}

          {/* Image counter */}
          {images.length > 1 && (
            <div className="absolute bottom-2 left-2 badge badge-neutral badge-sm">
              {currentIndex + 1} / {images.length}
            </div>
          )}

          {/* Fullscreen button */}
          {showFullscreenButton && (
            <button
              onClick={toggleFullscreen}
              className="absolute top-2 right-2 btn btn-circle btn-sm bg-base-100/80 hover:bg-base-100 backdrop-blur"
              aria-label="Fullscreen"
            >
              <Maximize2 className="w-4 h-4" />
            </button>
          )}

          {/* Dots indicator for mobile */}
          {images.length > 1 && (
            <div className="absolute bottom-2 left-1/2 -translate-x-1/2 flex gap-1 lg:hidden">
              {images.map((_, index) => (
                <button
                  key={index}
                  onClick={() => handleThumbnailClick(index)}
                  className={`w-2 h-2 rounded-full transition-all ${
                    index === currentIndex ? 'bg-primary w-4' : 'bg-base-100/60'
                  }`}
                  aria-label={`Go to photo ${index + 1}`}
                />
              ))}
            </div>
          )}
        </div>

        {/* Thumbnails for desktop */}
        {showThumbnails && images.length > 1 && (
          <div className="hidden lg:flex gap-2 mt-2 overflow-x-auto">
            {images.map((image, index) => (
              <button
                key={image.id || index}
                onClick={() => handleThumbnailClick(index)}
                className={`relative w-20 h-20 rounded-lg overflow-hidden flex-shrink-0 transition-all ${
                  index === currentIndex
                    ? 'ring-2 ring-primary ring-offset-2'
                    : 'opacity-70 hover:opacity-100'
                }`}
              >
                <OptimizedCarImage
                  src={image.thumbnail_url || image.url}
                  alt={`${title} - миниатюра ${index + 1}`}
                />
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Fullscreen modal */}
      {isFullscreen && (
        <div className="fixed inset-0 z-50 bg-black flex items-center justify-center">
          <button
            onClick={toggleFullscreen}
            className="absolute top-4 right-4 btn btn-circle btn-ghost text-white"
            aria-label="Close fullscreen"
          >
            <X className="w-6 h-6" />
          </button>

          <div className="relative w-full h-full flex items-center justify-center p-4">
            <CarPhotoSwiper
              images={images}
              title={title}
              showThumbnails={false}
              showFullscreenButton={false}
              autoPlay={false}
            />
          </div>

          {/* Navigation in fullscreen */}
          <button
            onClick={handlePrevious}
            className="absolute left-4 top-1/2 -translate-y-1/2 btn btn-circle btn-ghost text-white"
            aria-label="Previous photo"
          >
            <ChevronLeft className="w-8 h-8" />
          </button>
          <button
            onClick={handleNext}
            className="absolute right-4 top-1/2 -translate-y-1/2 btn btn-circle btn-ghost text-white"
            aria-label="Next photo"
          >
            <ChevronRight className="w-8 h-8" />
          </button>

          {/* Counter in fullscreen */}
          <div className="absolute bottom-4 left-1/2 -translate-x-1/2 text-white">
            {currentIndex + 1} / {images.length}
          </div>
        </div>
      )}
    </>
  );
};

export default CarPhotoSwiper;
