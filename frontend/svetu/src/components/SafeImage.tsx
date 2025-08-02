'use client';

import Image, { ImageProps } from 'next/image';
import { useState } from 'react';
import { getSafeImageUrl } from '@/utils/imageUtils';

interface SafeImageProps extends Omit<ImageProps, 'src'> {
  src: string | null | undefined;
  fallback?: React.ReactNode;
}

export default function SafeImage({
  src,
  alt,
  fallback,
  onError,
  ...props
}: SafeImageProps) {
  const [hasError, setHasError] = useState(false);
  const [isLoading, setIsLoading] = useState(true);

  // Для относительных путей не преобразуем в полный URL - позволяем Next.js rewrites работать
  const urlToCheck = src || '';
  const safeUrl = getSafeImageUrl(urlToCheck);

  // Отладка для всех изображений в SimilarListings
  if (
    src &&
    (src.includes('177') || (window as any).SimilarListingsImageDebug)
  ) {
    console.log('SafeImage debug:', {
      src,
      urlToCheck,
      safeUrl,
      hasError,
      isLoading,
    });
  }

  // Если URL небезопасный или произошла ошибка, показываем fallback
  if (!safeUrl || hasError) {
    if (fallback) {
      return <>{fallback}</>;
    }

    // Стандартный placeholder
    return (
      <div
        className="bg-base-200 flex items-center justify-center"
        style={{
          width: props.width || '100%',
          height: props.height || '100%',
          minHeight: '100px',
        }}
      >
        <svg
          className="w-12 h-12 text-base-content/20"
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
    );
  }

  return (
    <>
      {isLoading && (
        <div
          className="bg-base-200 animate-pulse"
          style={{
            width: props.width || '100%',
            height: props.height || '100%',
            minHeight: '100px',
          }}
        />
      )}
      <Image
        {...props}
        src={safeUrl}
        alt={alt}
        sizes={props.sizes || (props.fill ? '100vw' : undefined)}
        onLoad={() => {
          if (
            src &&
            (src.includes('177') || (window as any).SimilarListingsImageDebug)
          ) {
            console.log('SafeImage onLoad fired:', src);
          }
          setIsLoading(false);
        }}
        onError={(e) => {
          if (
            src &&
            (src.includes('177') || (window as any).SimilarListingsImageDebug)
          ) {
            console.log('SafeImage onError fired:', src, e);
          }
          setHasError(true);
          setIsLoading(false);
          if (onError) {
            onError(e);
          }
        }}
        style={{
          ...props.style,
          opacity: isLoading ? 0 : 1,
          transition: 'opacity 0.3s ease-in-out',
        }}
      />
    </>
  );
}
