'use client';

import { useState, useEffect, useRef } from 'react';
import Image from 'next/image';

interface LazyImageProps {
  src: string;
  alt: string;
  fill?: boolean;
  width?: number;
  height?: number;
  className?: string;
  objectFit?: 'contain' | 'cover' | 'fill' | 'none' | 'scale-down';
  sizes?: string;
  priority?: boolean;
  placeholder?: 'blur' | 'empty';
  blurDataURL?: string;
  onLoad?: () => void;
  onError?: () => void;
}

const LazyImage: React.FC<LazyImageProps> = ({
  src,
  alt,
  fill = false,
  width,
  height,
  className = '',
  objectFit = 'cover',
  sizes,
  priority = false,
  placeholder,
  blurDataURL,
  onLoad,
  onError,
}) => {
  const [isLoaded, setIsLoaded] = useState(false);
  const [isInView, setIsInView] = useState(false);
  const [hasError, setHasError] = useState(false);
  const imgRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (!imgRef.current || priority) {
      setIsInView(true);
      return;
    }

    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            setIsInView(true);
            observer.disconnect();
          }
        });
      },
      {
        threshold: 0.01,
        rootMargin: '50px',
      }
    );

    observer.observe(imgRef.current);

    return () => {
      observer.disconnect();
    };
  }, [priority]);

  const handleLoad = () => {
    setIsLoaded(true);
    onLoad?.();
  };

  const handleError = () => {
    setHasError(true);
    onError?.();
  };

  // Placeholder для скелетона
  const skeletonClass = 'animate-pulse bg-base-300';

  if (hasError) {
    return (
      <div
        ref={imgRef}
        className={`${className} flex items-center justify-center bg-base-200 text-base-content/30`}
        style={
          !fill
            ? {
                width: width || 'auto',
                height: height || 'auto',
              }
            : undefined
        }
      >
        <svg
          className="w-12 h-12"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
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
    <div ref={imgRef} className="relative">
      {/* Skeleton loader */}
      {!isLoaded && isInView && (
        <div
          className={`absolute inset-0 ${skeletonClass} ${className}`}
          style={
            !fill
              ? {
                  width: width || 'auto',
                  height: height || 'auto',
                }
              : undefined
          }
        />
      )}

      {/* Actual image */}
      {isInView && (
        <>
          {fill ? (
            <Image
              src={src}
              alt={alt}
              fill
              className={`${className} ${
                !isLoaded ? 'opacity-0' : 'opacity-100'
              } transition-opacity duration-300`}
              style={{ objectFit }}
              sizes={sizes}
              priority={priority}
              placeholder={placeholder}
              blurDataURL={blurDataURL}
              onLoad={handleLoad}
              onError={handleError}
            />
          ) : (
            <Image
              src={src}
              alt={alt}
              width={width || 500}
              height={height || 500}
              className={`${className} ${
                !isLoaded ? 'opacity-0' : 'opacity-100'
              } transition-opacity duration-300`}
              style={{ objectFit }}
              sizes={sizes}
              priority={priority}
              placeholder={placeholder}
              blurDataURL={blurDataURL}
              onLoad={handleLoad}
              onError={handleError}
            />
          )}
        </>
      )}

      {/* Low-quality placeholder for non-priority images */}
      {!isInView && !priority && (
        <div
          className={`${skeletonClass} ${className}`}
          style={
            !fill
              ? {
                  width: width || 'auto',
                  height: height || 'auto',
                }
              : undefined
          }
        />
      )}
    </div>
  );
};

export default LazyImage;
