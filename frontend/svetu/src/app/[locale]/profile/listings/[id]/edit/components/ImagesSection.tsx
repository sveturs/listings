'use client';

import { useTranslations } from 'next-intl';
import { useState, useRef, useCallback } from 'react';
import { apiClient } from '@/services/api-client';
import Image from 'next/image';
import configManager from '@/config';

interface Image {
  id: number;
  file_path: string;
  file_name: string;
  is_main: boolean;
  public_url: string;
}

interface ImagesSectionProps {
  listingId: number;
  images: Image[];
  onImagesChange: (images: Image[]) => void;
}

export function ImagesSection({
  listingId,
  images,
  onImagesChange,
}: ImagesSectionProps) {
  const t = useTranslations('profile');
  const [uploading, setUploading] = useState(false);
  const [dragActive, setDragActive] = useState(false);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleDrag = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  }, []);

  const handleFiles = useCallback(
    async (files: FileList) => {
      setUploading(true);

      try {
        const formData = new FormData();
        Array.from(files).forEach((file) => {
          formData.append('file', file);
        });

        // Если нет изображений, первое будет главным
        if (images.length === 0) {
          formData.append('main_image_index', '0');
        }

        const response = await apiClient.post(
          `/api/v1/marketplace/listings/${listingId}/images`,
          formData
          // Не устанавливаем Content-Type для FormData - браузер сделает это автоматически с boundary
        );

        console.log('Upload response:', response.data);

        if (response.data?.data?.images) {
          onImagesChange([...images, ...response.data.data.images]);
        } else if (response.data?.images) {
          onImagesChange([...images, ...response.data.images]);
        } else {
          console.warn('No images data in response:', response.data);
        }
      } catch (error) {
        console.error('Error uploading images:', error);
        // TODO: Show error toast
      } finally {
        setUploading(false);
      }
    },
    [images, listingId, onImagesChange]
  );

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      e.stopPropagation();
      setDragActive(false);

      if (e.dataTransfer.files && e.dataTransfer.files[0]) {
        handleFiles(e.dataTransfer.files);
      }
    },
    [handleFiles]
  );

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      handleFiles(e.target.files);
    }
  };

  const handleDelete = async (imageId: number) => {
    if (!confirm(t('images.deleteConfirm'))) return;

    try {
      await apiClient.delete(
        `/api/v1/marketplace/listings/${listingId}/images/${imageId}`
      );

      onImagesChange(images.filter((img) => img.id !== imageId));
    } catch (error) {
      console.error('Error deleting image:', error);
      // TODO: Show error toast
    }
  };

  const getImageUrl = (image: Image) => {
    return configManager.buildImageUrl(image.public_url);
  };

  const handleSetMain = async (imageId: number) => {
    try {
      // Обновляем локально
      const updatedImages = images.map((img) => ({
        ...img,
        is_main: img.id === imageId,
      }));

      onImagesChange(updatedImages);

      // Показываем уведомление что нужно сохранить изменения
      console.log(
        'Main image changed. Please save the listing to apply changes.'
      );
      // TODO: Show info toast about saving
    } catch (error) {
      console.error('Error setting main image:', error);
      // TODO: Show error toast
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-medium mb-4">{t('images.title')}</h3>

        {/* Upload area */}
        <div
          className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
            dragActive
              ? 'border-primary bg-primary/10'
              : 'border-base-300 hover:border-base-content/30'
          } ${uploading ? 'opacity-50 pointer-events-none' : ''}`}
          onDragEnter={handleDrag}
          onDragLeave={handleDrag}
          onDragOver={handleDrag}
          onDrop={handleDrop}
        >
          <input
            ref={fileInputRef}
            type="file"
            multiple
            accept="image/*"
            onChange={handleFileInput}
            className="hidden"
          />

          <svg
            className="mx-auto h-12 w-12 text-base-content/50"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
            />
          </svg>

          <p className="mt-2 text-sm text-base-content/70">
            {t('images.dragDrop')}
          </p>

          <button
            type="button"
            onClick={() => fileInputRef.current?.click()}
            className="btn btn-primary btn-sm mt-4"
            disabled={uploading}
          >
            {uploading ? (
              <>
                <span className="loading loading-spinner loading-xs"></span>
                {t('images.uploading')}
              </>
            ) : (
              t('images.selectFiles')
            )}
          </button>
        </div>

        {/* Images grid */}
        {images.length > 0 && (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 mt-6">
            {images.map((image, index) => (
              <div key={image.id} className="relative group">
                <div className="aspect-square rounded-lg overflow-hidden bg-base-200">
                  <Image
                    src={getImageUrl(image)}
                    alt={image.file_name}
                    fill
                    className="object-cover"
                  />
                </div>

                {/* Overlay controls */}
                <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity rounded-lg flex items-center justify-center gap-2">
                  {!image.is_main && (
                    <button
                      onClick={() => handleSetMain(image.id)}
                      className="btn btn-sm btn-primary"
                      title={t('images.setAsMain')}
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
                          d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"
                        />
                      </svg>
                    </button>
                  )}
                  <button
                    onClick={() => handleDelete(image.id)}
                    className="btn btn-sm btn-error"
                    title={t('images.delete')}
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
                        d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"
                      />
                    </svg>
                  </button>
                </div>

                {/* Main badge */}
                {image.is_main && (
                  <div className="absolute top-2 left-2 badge badge-primary badge-sm">
                    {t('images.main')}
                  </div>
                )}

                {/* Index */}
                <div className="absolute top-2 right-2 badge badge-neutral badge-sm">
                  {index + 1}
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Help text */}
        <div className="text-sm text-base-content/70 mt-4">
          <p>{t('images.helpText')}</p>
          <p className="mt-1">{t('images.maxSize')}</p>
        </div>
      </div>
    </div>
  );
}
