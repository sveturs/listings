'use client';

import React, { useState, useEffect, useRef, ImgHTMLAttributes } from 'react';
import Image from 'next/image';

interface OptimizedImageProps
  extends Omit<ImgHTMLAttributes<HTMLImageElement>, 'src'> {
  src: string;
  alt: string;
  width?: number;
  height?: number;
  priority?: boolean;
  lazy?: boolean;
  placeholder?: 'blur' | 'empty' | 'skeleton';
  blurDataURL?: string;
  quality?: number;
  sizes?: string;
  formats?: ('webp' | 'avif' | 'jpg' | 'png')[];
  onLoad?: () => void;
  onError?: () => void;
  fallbackSrc?: string;
  aspectRatio?: string;
  objectFit?: 'contain' | 'cover' | 'fill' | 'none' | 'scale-down';
  enableProgressive?: boolean;
}

interface ImageState {
  isLoading: boolean;
  isInView: boolean;
  hasError: boolean;
  currentSrc: string;
  isLowQuality: boolean;
}

/**
 * Mobile-optimized image component with progressive loading
 */
export default function OptimizedImage({
  src,
  alt,
  width,
  height,
  priority = false,
  lazy = true,
  placeholder = 'blur',
  blurDataURL,
  quality = 75,
  sizes,
  formats = ['webp', 'avif', 'jpg'],
  onLoad,
  onError,
  fallbackSrc = '/images/placeholder.jpg',
  aspectRatio = '16/9',
  objectFit = 'cover',
  enableProgressive = true,
  className = '',
  ...props
}: OptimizedImageProps) {
  const [state, setState] = useState<ImageState>({
    isLoading: true,
    isInView: !lazy || priority,
    hasError: false,
    currentSrc: src,
    isLowQuality: enableProgressive,
  });

  const imageRef = useRef<HTMLDivElement>(null);
  const observerRef = useRef<IntersectionObserver | null>(null);

  // Device and connection detection
  const getOptimalImageSettings = () => {
    const connection =
      (navigator as any).connection ||
      (navigator as any).mozConnection ||
      (navigator as any).webkitConnection;

    const isSlowConnection =
      connection?.effectiveType === '2g' ||
      connection?.effectiveType === 'slow-2g';
    const isSaveData = connection?.saveData || false;
    const devicePixelRatio = window.devicePixelRatio || 1;

    // Determine optimal quality
    let optimalQuality = quality;
    if (isSaveData) {
      optimalQuality = 40;
    } else if (isSlowConnection) {
      optimalQuality = 60;
    } else if (devicePixelRatio > 2) {
      optimalQuality = Math.min(quality, 85); // High DPI but cap quality
    }

    // Determine optimal format
    const supportsWebP = checkWebPSupport();
    const supportsAVIF = checkAVIFSupport();

    let optimalFormat = 'jpg';
    if (supportsAVIF && formats.includes('avif')) {
      optimalFormat = 'avif';
    } else if (supportsWebP && formats.includes('webp')) {
      optimalFormat = 'webp';
    }

    return {
      quality: optimalQuality,
      format: optimalFormat,
      shouldUseLowRes: isSlowConnection || isSaveData,
      devicePixelRatio,
    };
  };

  // Check WebP support
  const checkWebPSupport = (): boolean => {
    if (typeof window === 'undefined') return false;

    const canvas = document.createElement('canvas');
    canvas.width = canvas.height = 1;
    return canvas.toDataURL('image/webp').indexOf('image/webp') === 0;
  };

  // Check AVIF support
  const checkAVIFSupport = (): boolean => {
    if (typeof window === 'undefined') return false;

    // Simple check - more comprehensive check would use actual image loading
    return CSS.supports('background-image', 'url("data:image/avif;base64,")');
  };

  // Generate responsive image URLs
  const generateImageUrls = (baseSrc: string) => {
    const settings = getOptimalImageSettings();
    const baseUrl = baseSrc.split('?')[0];
    const extension = baseUrl.split('.').pop();
    const nameWithoutExt = baseUrl.substring(0, baseUrl.lastIndexOf('.'));

    // Generate URLs for different sizes
    const widths = [320, 640, 768, 1024, 1280, 1920];
    const urls: Record<number, string> = {};

    widths.forEach((w) => {
      const params = new URLSearchParams({
        w: w.toString(),
        q: settings.quality.toString(),
        fmt: settings.format,
      });

      urls[w] = `${nameWithoutExt}.${extension}?${params.toString()}`;
    });

    return urls;
  };

  // Generate srcSet string
  const generateSrcSet = (): string => {
    const urls = generateImageUrls(src);
    return Object.entries(urls)
      .map(([width, url]) => `${url} ${width}w`)
      .join(', ');
  };

  // Generate sizes string
  const generateSizes = (): string => {
    if (sizes) return sizes;

    // Default responsive sizes
    return '(max-width: 640px) 100vw, (max-width: 768px) 90vw, (max-width: 1024px) 80vw, 70vw';
  };

  // Progressive image loading
  const loadProgressiveImage = async () => {
    if (!enableProgressive) return;

    const settings = getOptimalImageSettings();

    try {
      // Load low quality first
      if (settings.shouldUseLowRes) {
        const lowQualitySrc = src.replace(/(\.[^.]+)$/, '-low$1');
        await loadImage(lowQualitySrc);
        setState((prev) => ({
          ...prev,
          currentSrc: lowQualitySrc,
          isLowQuality: true,
        }));
      }

      // Load full quality
      await loadImage(src);
      setState((prev) => ({
        ...prev,
        currentSrc: src,
        isLoading: false,
        isLowQuality: false,
      }));

      onLoad?.();
    } catch (error) {
      handleImageError();
    }
  };

  // Load image promise
  const loadImage = (imageSrc: string): Promise<void> => {
    return new Promise((resolve, reject) => {
      const img = new Image();
      img.onload = () => resolve();
      img.onerror = reject;
      img.src = imageSrc;
    });
  };

  // Handle image error
  const handleImageError = () => {
    setState((prev) => ({
      ...prev,
      hasError: true,
      currentSrc: fallbackSrc,
      isLoading: false,
    }));
    onError?.();
  };

  // Intersection Observer for lazy loading
  useEffect(() => {
    if (!lazy || priority || state.isInView) return;

    const options: IntersectionObserverInit = {
      rootMargin: '50px', // Start loading 50px before visible
      threshold: 0.01,
    };

    observerRef.current = new IntersectionObserver((entries) => {
      entries.forEach((entry) => {
        if (entry.isIntersecting) {
          setState((prev) => ({ ...prev, isInView: true }));
          observerRef.current?.disconnect();
        }
      });
    }, options);

    if (imageRef.current) {
      observerRef.current.observe(imageRef.current);
    }

    return () => {
      observerRef.current?.disconnect();
    };
  }, [lazy, priority, state.isInView]);

  // Load image when in view
  useEffect(() => {
    if (state.isInView && !state.hasError) {
      loadProgressiveImage();
    }
  }, [state.isInView]);

  // Generate blur placeholder
  const getPlaceholder = () => {
    if (placeholder === 'blur' && blurDataURL) {
      return blurDataURL;
    }

    if (placeholder === 'blur' && !blurDataURL) {
      // Generate simple blur placeholder
      return `data:image/svg+xml;base64,${Buffer.from(
        `<svg width="${width || 100}" height="${height || 100}" xmlns="http://www.w3.org/2000/svg">
          <rect width="100%" height="100%" fill="#e5e7eb"/>
        </svg>`
      ).toString('base64')}`;
    }

    return undefined;
  };

  // Render skeleton placeholder
  const renderSkeleton = () => (
    <div
      className={`animate-pulse bg-gray-200 ${className}`}
      style={{
        aspectRatio,
        width: width || '100%',
        height: height || 'auto',
      }}
    />
  );

  // Render loading state
  if (!state.isInView && lazy) {
    if (placeholder === 'skeleton') {
      return renderSkeleton();
    }

    return (
      <div
        ref={imageRef}
        className={className}
        style={{
          aspectRatio,
          width: width || '100%',
          height: height || 'auto',
          backgroundColor: '#f3f4f6',
        }}
      />
    );
  }

  // Use Next.js Image for optimization
  if (!state.hasError && typeof window === 'undefined') {
    return (
      <Image
        src={state.currentSrc}
        alt={alt}
        width={width || 500}
        height={height || 500}
        priority={priority}
        quality={getOptimalImageSettings().quality}
        placeholder={placeholder === 'blur' ? 'blur' : 'empty'}
        blurDataURL={getPlaceholder()}
        sizes={generateSizes()}
        className={`${className} ${state.isLowQuality ? 'blur-sm' : ''}`}
        style={{ objectFit }}
        onLoad={() => {
          setState((prev) => ({ ...prev, isLoading: false }));
          onLoad?.();
        }}
        onError={handleImageError}
        {...props}
      />
    );
  }

  // Native img with optimizations
  return (
    <div ref={imageRef} className="relative" style={{ aspectRatio }}>
      {state.isLoading && placeholder === 'skeleton' && renderSkeleton()}

      <picture>
        {/* AVIF source */}
        {formats.includes('avif') && checkAVIFSupport() && (
          <source
            type="image/avif"
            srcSet={generateSrcSet().replace(/\.(jpg|png|webp)/g, '.avif')}
            sizes={generateSizes()}
          />
        )}

        {/* WebP source */}
        {formats.includes('webp') && checkWebPSupport() && (
          <source
            type="image/webp"
            srcSet={generateSrcSet().replace(/\.(jpg|png)/g, '.webp')}
            sizes={generateSizes()}
          />
        )}

        {/* Fallback image */}
        <img
          src={state.currentSrc}
          alt={alt}
          width={width}
          height={height}
          loading={lazy && !priority ? 'lazy' : 'eager'}
          decoding="async"
          className={`${className} ${state.isLoading ? 'opacity-0' : 'opacity-100'} ${state.isLowQuality ? 'blur-sm' : ''} transition-opacity duration-300`}
          style={{
            objectFit,
            width: width || '100%',
            height: height || 'auto',
          }}
          onLoad={() => {
            setState((prev) => ({ ...prev, isLoading: false }));
            onLoad?.();
          }}
          onError={handleImageError}
          {...props}
        />
      </picture>
    </div>
  );
}

/**
 * Image optimization utilities
 */
export const ImageOptimizer = {
  /**
   * Convert image to WebP format
   */
  async convertToWebP(file: File, quality: number = 0.8): Promise<Blob> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = (e) => {
        const img = new Image();

        img.onload = () => {
          const canvas = document.createElement('canvas');
          canvas.width = img.width;
          canvas.height = img.height;

          const ctx = canvas.getContext('2d');
          if (!ctx) {
            reject(new Error('Failed to get canvas context'));
            return;
          }

          ctx.drawImage(img, 0, 0);

          canvas.toBlob(
            (blob) => {
              if (blob) {
                resolve(blob);
              } else {
                reject(new Error('Failed to convert to WebP'));
              }
            },
            'image/webp',
            quality
          );
        };

        img.onerror = reject;
        img.src = e.target?.result as string;
      };

      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  },

  /**
   * Resize image to maximum dimensions
   */
  async resizeImage(
    file: File,
    maxWidth: number,
    maxHeight: number
  ): Promise<Blob> {
    return new Promise((resolve, reject) => {
      const reader = new FileReader();

      reader.onload = (e) => {
        const img = new Image();

        img.onload = () => {
          let { width, height } = img;

          // Calculate new dimensions
          if (width > maxWidth || height > maxHeight) {
            const ratio = Math.min(maxWidth / width, maxHeight / height);
            width *= ratio;
            height *= ratio;
          }

          const canvas = document.createElement('canvas');
          canvas.width = width;
          canvas.height = height;

          const ctx = canvas.getContext('2d');
          if (!ctx) {
            reject(new Error('Failed to get canvas context'));
            return;
          }

          // Enable image smoothing
          ctx.imageSmoothingEnabled = true;
          ctx.imageSmoothingQuality = 'high';

          ctx.drawImage(img, 0, 0, width, height);

          canvas.toBlob(
            (blob) => {
              if (blob) {
                resolve(blob);
              } else {
                reject(new Error('Failed to resize image'));
              }
            },
            file.type,
            0.9
          );
        };

        img.onerror = reject;
        img.src = e.target?.result as string;
      };

      reader.onerror = reject;
      reader.readAsDataURL(file);
    });
  },

  /**
   * Generate blur placeholder
   */
  async generateBlurPlaceholder(
    src: string,
    size: number = 40
  ): Promise<string> {
    return new Promise((resolve) => {
      const img = new Image();

      img.onload = () => {
        const canvas = document.createElement('canvas');
        const ratio = img.width / img.height;

        canvas.width = size;
        canvas.height = size / ratio;

        const ctx = canvas.getContext('2d');
        if (!ctx) {
          resolve('');
          return;
        }

        ctx.filter = 'blur(5px)';
        ctx.drawImage(img, 0, 0, canvas.width, canvas.height);

        resolve(canvas.toDataURL('image/jpeg', 0.4));
      };

      img.onerror = () => resolve('');
      img.src = src;
    });
  },
};
