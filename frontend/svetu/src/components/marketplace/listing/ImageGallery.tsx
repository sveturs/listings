'use client';

import { useState, useRef, useEffect, useCallback } from 'react';
import Image from 'next/image';
import { useLocale } from 'next-intl';
import config from '@/config';

interface ImageGalleryProps {
  images: Array<{
    id: number;
    public_url: string;
    is_video?: boolean;
  }>;
  title: string;
}

export default function ImageGallery({ images, title }: ImageGalleryProps) {
  const locale = useLocale();
  const [selectedIndex, setSelectedIndex] = useState(0);
  const [showLightbox, setShowLightbox] = useState(false);
  const [isZoomed, setIsZoomed] = useState(false);
  const [zoomPosition, setZoomPosition] = useState({ x: 50, y: 50 });
  const imageRef = useRef<HTMLDivElement>(null);
  const [imageNaturalDimensions, setImageNaturalDimensions] = useState({
    width: 0,
    height: 0,
  });

  const navigateImage = useCallback(
    (direction: number) => {
      setSelectedIndex((prev) => {
        const newIndex = prev + direction;
        if (newIndex < 0) return images.length - 1;
        if (newIndex >= images.length) return 0;
        return newIndex;
      });
    },
    [images.length]
  );

  // Load natural dimensions of current image
  useEffect(() => {
    const currentImage = images[selectedIndex];
    if (!currentImage || currentImage.id === 0 || currentImage.is_video) {
      return;
    }

    const img = new window.Image();
    img.onload = () => {
      setImageNaturalDimensions({
        width: img.naturalWidth,
        height: img.naturalHeight,
      });
    };
    img.src = config.buildImageUrl(currentImage.public_url);
  }, [selectedIndex, images]);

  // Keyboard navigation
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if (!showLightbox) return;

      if (e.key === 'Escape') {
        setShowLightbox(false);
      } else if (e.key === 'ArrowLeft') {
        navigateImage(-1);
      } else if (e.key === 'ArrowRight') {
        navigateImage(1);
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, [showLightbox, navigateImage]);

  const handleMouseMove = (e: React.MouseEvent<HTMLDivElement>) => {
    if (
      !imageRef.current ||
      !imageNaturalDimensions.width ||
      !imageNaturalDimensions.height
    )
      return;

    const container = imageRef.current;
    const rect = container.getBoundingClientRect();

    // Размеры контейнера
    const containerWidth = rect.width;
    const containerHeight = rect.height;

    // Соотношение сторон изображения
    const imageAspectRatio =
      imageNaturalDimensions.width / imageNaturalDimensions.height;
    // Соотношение сторон контейнера (aspect-[4/3] = 4/3 = 1.333...)
    const containerAspectRatio = 4 / 3;

    let renderWidth,
      renderHeight,
      offsetX = 0,
      offsetY = 0;

    if (imageAspectRatio > containerAspectRatio) {
      // Изображение шире контейнера - ограничено по ширине
      renderWidth = containerWidth;
      renderHeight = containerWidth / imageAspectRatio;
      offsetY = (containerHeight - renderHeight) / 2;
    } else {
      // Изображение выше контейнера - ограничено по высоте
      renderHeight = containerHeight;
      renderWidth = containerHeight * imageAspectRatio;
      offsetX = (containerWidth - renderWidth) / 2;
    }

    // Позиция курсора относительно контейнера
    const mouseX = e.clientX - rect.left;
    const mouseY = e.clientY - rect.top;

    // Проверяем, находится ли курсор над изображением
    const isOverImage =
      mouseX >= offsetX &&
      mouseX <= offsetX + renderWidth &&
      mouseY >= offsetY &&
      mouseY <= offsetY + renderHeight;

    if (!isOverImage) {
      setIsZoomed(false);
      return;
    }

    // Рассчитываем позицию относительно изображения
    const x = ((mouseX - offsetX) / renderWidth) * 100;
    const y = ((mouseY - offsetY) / renderHeight) * 100;

    setIsZoomed(true);
    setZoomPosition({ x, y });
  };

  const renderImage = (image: (typeof images)[0], isMain = false) => {
    if (image.id === 0) {
      return (
        // eslint-disable-next-line @next/next/no-img-element
        <img
          src="/placeholder-listing.jpg"
          alt={title}
          className={
            isMain
              ? 'w-full h-full object-contain'
              : 'w-full h-full object-cover'
          }
        />
      );
    }

    if (image.is_video) {
      return (
        <div className="w-full h-full flex items-center justify-center bg-base-300">
          <svg
            className="w-16 h-16 text-base-content/40"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
            />
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        </div>
      );
    }

    return (
      <Image
        src={config.buildImageUrl(image.public_url)}
        alt={`${title} ${isMain ? '' : selectedIndex + 1}`}
        fill
        className={isMain ? 'object-contain' : 'object-cover'}
        priority={isMain && selectedIndex === 0}
        sizes={
          isMain
            ? '(max-width: 768px) 100vw, 66vw'
            : '(max-width: 640px) 80px, 80px'
        }
      />
    );
  };

  return (
    <>
      <div className="bg-base-200 rounded-2xl overflow-hidden">
        {/* Main Image */}
        <div
          ref={imageRef}
          className="relative aspect-[4/3] w-full cursor-zoom-in overflow-hidden"
          onClick={() => setShowLightbox(true)}
          onMouseLeave={() => setIsZoomed(false)}
          onMouseMove={handleMouseMove}
        >
          {renderImage(images[selectedIndex], true)}

          {/* Image counter */}
          {images.length > 1 && (
            <div className="absolute top-4 right-4 bg-base-100/80 backdrop-blur-sm px-3 py-1 rounded-full text-sm font-medium">
              {selectedIndex + 1} / {images.length}
            </div>
          )}

          {/* Navigation arrows */}
          {images.length > 1 && (
            <>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  navigateImage(-1);
                }}
                className="absolute left-4 top-1/2 -translate-y-1/2 btn btn-circle btn-sm bg-base-100/80 hover:bg-base-100"
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
                    d="M15 19l-7-7 7-7"
                  />
                </svg>
              </button>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  navigateImage(1);
                }}
                className="absolute right-4 top-1/2 -translate-y-1/2 btn btn-circle btn-sm bg-base-100/80 hover:bg-base-100"
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
                    d="M9 5l7 7-7 7"
                  />
                </svg>
              </button>
            </>
          )}

          {/* Zoom indicator */}
          <div className="absolute bottom-4 left-4 bg-base-100/80 backdrop-blur-sm px-3 py-1 rounded-full text-sm">
            <svg
              className="w-4 h-4 inline mr-1"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v3m0 0v3m0-3h3m-3 0H7"
              />
            </svg>
            {locale === 'ru' ? 'Нажмите для увеличения' : 'Click to zoom'}
          </div>

          {/* Zoom preview (desktop only) */}
          {isZoomed &&
            images[selectedIndex].id !== 0 &&
            !images[selectedIndex].is_video && (
              <div
                className="absolute pointer-events-none w-64 h-64 border-2 border-primary rounded-lg overflow-hidden shadow-xl hidden lg:block z-10"
                style={{
                  right: '1rem',
                  bottom: '1rem',
                  backgroundImage: `url(${config.buildImageUrl(images[selectedIndex].public_url)})`,
                  backgroundPosition: `${zoomPosition.x}% ${zoomPosition.y}%`,
                  backgroundSize: '700%',
                  backgroundRepeat: 'no-repeat',
                }}
              />
            )}
        </div>

        {/* Thumbnail Gallery */}
        {images.length > 1 && (
          <div className="flex gap-2 p-4 overflow-x-auto scrollbar-thin scrollbar-thumb-base-300">
            {images.map((image, index) => (
              <button
                key={image.id}
                onClick={() => setSelectedIndex(index)}
                className={`relative flex-shrink-0 w-20 h-20 rounded-lg overflow-hidden transition-all ${
                  selectedIndex === index
                    ? 'ring-2 ring-primary scale-105'
                    : 'opacity-70 hover:opacity-100'
                }`}
              >
                {renderImage(image)}
                {image.is_video && (
                  <div className="absolute inset-0 flex items-center justify-center bg-black/30">
                    <svg
                      className="w-6 h-6 text-white"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"
                      />
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                  </div>
                )}
              </button>
            ))}
          </div>
        )}
      </div>

      {/* Lightbox */}
      {showLightbox && (
        <div
          className="fixed inset-0 z-50 bg-black/95 flex items-center justify-center p-4"
          onClick={() => setShowLightbox(false)}
        >
          <button
            className="absolute top-4 right-4 btn btn-circle btn-ghost text-white"
            onClick={() => setShowLightbox(false)}
          >
            <svg
              className="w-6 h-6"
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

          <div
            className="relative max-w-7xl max-h-[90vh] w-full h-full flex items-center justify-center"
            onClick={(e) => e.stopPropagation()}
          >
            {renderImage(images[selectedIndex], true)}
          </div>

          {/* Lightbox navigation */}
          {images.length > 1 && (
            <>
              <button
                onClick={(e) => {
                  e.stopPropagation();
                  navigateImage(-1);
                }}
                className="absolute left-4 top-1/2 -translate-y-1/2 btn btn-circle btn-lg btn-ghost text-white"
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
                onClick={(e) => {
                  e.stopPropagation();
                  navigateImage(1);
                }}
                className="absolute right-4 top-1/2 -translate-y-1/2 btn btn-circle btn-lg btn-ghost text-white"
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

          {/* Lightbox thumbnails */}
          {images.length > 1 && (
            <div className="absolute bottom-4 left-1/2 -translate-x-1/2 flex gap-2 max-w-full overflow-x-auto p-2 bg-black/50 backdrop-blur-sm rounded-lg">
              {images.map((image, index) => (
                <button
                  key={image.id}
                  onClick={(e) => {
                    e.stopPropagation();
                    setSelectedIndex(index);
                  }}
                  className={`relative flex-shrink-0 w-16 h-16 rounded overflow-hidden transition-all ${
                    selectedIndex === index
                      ? 'ring-2 ring-white opacity-100'
                      : 'opacity-50 hover:opacity-80'
                  }`}
                >
                  {renderImage(image)}
                </button>
              ))}
            </div>
          )}
        </div>
      )}
    </>
  );
}
