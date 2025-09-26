'use client';

import React, { useState } from 'react';
import Image from 'next/image';
import { Car } from 'lucide-react';

interface OptimizedCarImageProps {
  src?: string | null;
  alt: string;
  className?: string;
  sizes?: string;
  priority?: boolean;
  onLoad?: () => void;
  onError?: () => void;
}

/**
 * Оптимизированный компонент изображения автомобиля с lazy loading
 * и анимацией загрузки
 */
export const OptimizedCarImage: React.FC<OptimizedCarImageProps> = ({
  src,
  alt,
  className = '',
  sizes = '(max-width: 640px) 100vw, (max-width: 1024px) 50vw, 33vw',
  priority = false,
  onLoad,
  onError,
}) => {
  const [imageLoading, setImageLoading] = useState(true);
  const [imageError, setImageError] = useState(false);

  const handleImageLoad = () => {
    setImageLoading(false);
    onLoad?.();
  };

  const handleImageError = () => {
    setImageLoading(false);
    setImageError(true);
    onError?.();
  };

  // Если нет изображения или произошла ошибка
  if (!src || imageError) {
    return (
      <div
        className={`w-full h-full bg-base-200 flex items-center justify-center ${className}`}
      >
        <Car className="w-12 h-12 text-base-content/30" />
      </div>
    );
  }

  return (
    <div className={`relative w-full h-full ${className}`}>
      {/* Skeleton loader */}
      {imageLoading && (
        <div className="absolute inset-0 bg-base-200 animate-pulse flex items-center justify-center">
          <div className="loading loading-spinner loading-md"></div>
        </div>
      )}

      {/* Blur placeholder для быстрой загрузки */}
      <div
        className={`absolute inset-0 bg-gradient-to-br from-base-200 to-base-300 ${
          imageLoading ? 'opacity-100' : 'opacity-0'
        } transition-opacity duration-300`}
      />

      {/* Основное изображение */}
      <Image
        src={src}
        alt={alt}
        fill
        sizes={sizes}
        priority={priority}
        loading={priority ? 'eager' : 'lazy'}
        className={`object-cover ${
          imageLoading ? 'opacity-0' : 'opacity-100'
        } transition-opacity duration-300`}
        onLoad={handleImageLoad}
        onError={handleImageError}
        placeholder="empty" // Используем пустой placeholder, так как у нас свой skeleton
        quality={85} // Оптимальное качество для автомобильных фото
      />
    </div>
  );
};

/**
 * Компонент галереи изображений с lazy loading для множественных фото
 */
interface CarImageGalleryProps {
  images?: Array<{ thumbnail_url?: string; url?: string; id?: number }>;
  title: string;
  className?: string;
  showCount?: number;
}

export const CarImageGallery: React.FC<CarImageGalleryProps> = ({
  images,
  title,
  className = '',
  showCount = 1,
}) => {
  // const [currentIndex, setCurrentIndex] = useState(0);

  if (!images || images.length === 0) {
    return (
      <div
        className={`bg-base-200 flex items-center justify-center ${className}`}
      >
        <Car className="w-12 h-12 text-base-content/30" />
      </div>
    );
  }

  const visibleImages = images.slice(0, showCount);
  const hasMore = images.length > showCount;

  // Для одного изображения
  if (showCount === 1) {
    const image = images[0];
    return (
      <div className={`relative ${className}`}>
        <OptimizedCarImage
          src={image.thumbnail_url || image.url}
          alt={title}
          priority={false}
        />
        {hasMore && (
          <div className="absolute bottom-2 right-2 badge badge-neutral">
            +{images.length - 1} фото
          </div>
        )}
      </div>
    );
  }

  // Для нескольких изображений (галерея)
  return (
    <div
      className={`grid grid-cols-${Math.min(showCount, 2)} gap-1 ${className}`}
    >
      {visibleImages.map((image, index) => (
        <div
          key={image.id || index}
          className={`relative ${index === 0 && showCount > 2 ? 'col-span-2' : ''}`}
        >
          <OptimizedCarImage
            src={image.thumbnail_url || image.url}
            alt={`${title} - фото ${index + 1}`}
            priority={index === 0} // Приоритет только для первого изображения
          />
        </div>
      ))}
      {hasMore && (
        <div className="absolute bottom-2 right-2 badge badge-neutral">
          +{images.length - showCount} фото
        </div>
      )}
    </div>
  );
};

export default OptimizedCarImage;
