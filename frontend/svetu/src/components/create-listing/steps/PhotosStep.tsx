'use client';

import { useState, useRef } from 'react';
import { useTranslations } from 'next-intl';
import Image from 'next/image';
import { useCreateListing } from '@/contexts/CreateListingContext';
import { toast } from '@/utils/toast';
import { ImageGallery } from '@/components/reviews/ImageGallery';

interface PhotosStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function PhotosStep({ onNext, onBack }: PhotosStepProps) {
  const t = useTranslations('create_listing');
  const tCommon = useTranslations('common');
  const { state, dispatch } = useCreateListing();
  const [photos, setPhotos] = useState<string[]>(state.images || []);
  const [mainPhotoIndex, setMainPhotoIndex] = useState(
    state.mainImageIndex || 0
  );
  const [uploading, setUploading] = useState(false);
  const [dragOver, setDragOver] = useState(false);
  const [galleryOpen, setGalleryOpen] = useState(false);
  const [galleryIndex, setGalleryIndex] = useState(0);
  const fileInputRef = useRef<HTMLInputElement>(null);

  const handleFileSelect = async (files: FileList | null) => {
    if (!files) return;

    setUploading(true);
    const newPhotos: string[] = [];
    let skippedFiles = 0;

    for (let i = 0; i < Math.min(files.length, 8 - photos.length); i++) {
      const file = files[i];

      // TODO: Интегрировать загрузку изображений через MinIO
      // ВРЕМЕННОЕ РЕШЕНИЕ: Локальная обработка изображений
      // См. TODO #19: Интегрировать загрузку изображений через MinIO
      // Проверка размера файла (макс 10MB)
      if (file.size > 10 * 1024 * 1024) {
        skippedFiles++;
        toast.error(`${file.name} превышает 10MB и был пропущен`);
        continue;
      }

      // Создаем превью
      const reader = new FileReader();
      reader.onload = (e) => {
        if (e.target?.result) {
          newPhotos.push(e.target.result as string);
          if (newPhotos.length === Math.min(files.length, 8 - photos.length)) {
            const updatedPhotos = [...photos, ...newPhotos];
            setPhotos(updatedPhotos);
            dispatch({ type: 'SET_IMAGES', payload: updatedPhotos });
            setUploading(false);
          }
        }
      };
      reader.readAsDataURL(file);
    }

    // Если все файлы были пропущены или файлов нет
    if (
      files.length === 0 ||
      (skippedFiles > 0 &&
        newPhotos.length === 0 &&
        Math.min(files.length, 8 - photos.length) === skippedFiles)
    ) {
      setUploading(false);
      if (skippedFiles > 0) {
        toast.error('Все файлы превышают максимальный размер 10MB');
      }
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
    handleFileSelect(e.dataTransfer.files);
  };

  const handleDragOver = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(true);
  };

  const handleDragLeave = (e: React.DragEvent) => {
    e.preventDefault();
    setDragOver(false);
  };

  const removePhoto = (index: number) => {
    const newPhotos = photos.filter((_, i) => i !== index);
    setPhotos(newPhotos);
    dispatch({ type: 'SET_IMAGES', payload: newPhotos });

    // Adjusting main photo index if needed
    if (mainPhotoIndex >= newPhotos.length) {
      const newMainIndex = Math.max(0, newPhotos.length - 1);
      setMainPhotoIndex(newMainIndex);
      dispatch({ type: 'SET_MAIN_IMAGE', payload: newMainIndex });
    }
  };

  const setAsMainPhoto = (index: number) => {
    setMainPhotoIndex(index);
    dispatch({ type: 'SET_MAIN_IMAGE', payload: index });
  };

  const canProceed = photos.length > 0;

  const openGallery = (index: number) => {
    setGalleryIndex(index);
    setGalleryOpen(true);
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            📸 {t('title')}
          </h2>
          <p className="text-base-content/70 mb-6">{t('description')}</p>

          {/* Зона загрузки */}
          <div
            className={`
              border-2 border-dashed rounded-lg p-8 text-center transition-all duration-200
              ${
                dragOver
                  ? 'border-primary bg-primary/5'
                  : 'border-base-300 hover:border-primary/50'
              }
              ${uploading ? 'opacity-50 pointer-events-none' : 'cursor-pointer'}
            `}
            onDrop={handleDrop}
            onDragOver={handleDragOver}
            onDragLeave={handleDragLeave}
            onClick={() => !uploading && fileInputRef.current?.click()}
          >
            <input
              ref={fileInputRef}
              type="file"
              multiple
              accept="image/*"
              className="hidden"
              onChange={(e) => handleFileSelect(e.target.files)}
              disabled={uploading || photos.length >= 8}
            />

            {uploading ? (
              <div className="space-y-2">
                <div className="loading loading-spinner loading-lg mx-auto"></div>
                <p className="text-sm">{t('uploading')}</p>
              </div>
            ) : (
              <div className="space-y-2">
                <div className="text-4xl">📷</div>
                <p className="font-medium">{t('upload_instruction')}</p>
                <p className="text-sm text-base-content/60">
                  {t('supported_formats')}
                </p>
                <p className="text-xs text-base-content/50">
                  {t('max_size')} • {photos.length}/8 {t('photos')}
                </p>
              </div>
            )}
          </div>

          {/* Галерея загруженных фото */}
          {photos.length > 0 && (
            <div className="mt-6">
              <h3 className="font-medium mb-3 flex items-center gap-2">
                🖼️ {t('gallery')}
                <span className="badge badge-primary badge-sm">
                  {photos.length}
                </span>
              </h3>

              <div className="grid grid-cols-2 sm:grid-cols-3 gap-3">
                {photos.map((photo, index) => (
                  <div
                    key={index}
                    className={`
                      relative group rounded-lg overflow-hidden border-2 transition-all
                      ${
                        index === mainPhotoIndex
                          ? 'border-primary shadow-lg'
                          : 'border-base-300'
                      }
                    `}
                  >
                    <div
                      className="cursor-pointer relative"
                      onClick={() => openGallery(index)}
                    >
                      <Image
                        src={photo}
                        alt={`Photo ${index + 1}`}
                        width={200}
                        height={128}
                        className="w-full h-24 sm:h-32 object-cover"
                      />
                      {/* Zoom indicator on hover */}
                      <div className="absolute inset-0 bg-black/0 group-hover:bg-black/30 transition-all duration-200 flex items-center justify-center opacity-0 group-hover:opacity-100">
                        <svg
                          className="w-8 h-8 text-white"
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
                      </div>
                    </div>

                    {/* Overlay с действиями */}
                    <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity duration-200 flex items-center justify-center gap-2 z-10 pointer-events-none group-hover:pointer-events-auto">
                      {index !== mainPhotoIndex && (
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setAsMainPhoto(index);
                          }}
                          className="btn btn-xs btn-primary"
                          title={t('set_main')}
                        >
                          ⭐
                        </button>
                      )}
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          removePhoto(index);
                        }}
                        className="btn btn-xs btn-error"
                        title={t('remove')}
                      >
                        🗑️
                      </button>
                    </div>

                    {/* Индикатор главного фото */}
                    {index === mainPhotoIndex && (
                      <div className="absolute top-1 left-1 bg-primary text-primary-content text-xs px-2 py-1 rounded">
                        {t('main')}
                      </div>
                    )}
                  </div>
                ))}
              </div>

              {/* Подсказка об оптимизации для слабого интернета */}
              <div className="alert alert-info mt-4">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  className="stroke-current shrink-0 w-6 h-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  ></path>
                </svg>
                <div className="text-sm">
                  <p className="font-medium">💡 {t('photos.tips.title')}</p>
                  <ul className="text-xs mt-2 space-y-1">
                    <li>• {t('photos.tips.quality')}</li>
                    <li>• {t('photos.tips.lighting')}</li>
                    <li>• {t('photos.tips.angles')}</li>
                    <li>• {t('photos.tips.defects')}</li>
                  </ul>
                </div>
              </div>
            </div>
          )}

          {/* Совет для региональных пользователей */}
          <div className="alert alert-warning mt-4">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="stroke-current shrink-0 w-6 h-6"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.982 16.5c-.77.833.192 2.5 1.732 2.5z"
              ></path>
            </svg>
            <div className="text-sm">
              <p className="font-medium">🤝 {t('photos.trust_tip.title')}</p>
              <p className="text-xs mt-1">
                {t('photos.trust_tip.description')}
              </p>
            </div>
          </div>

          {/* Кнопки навигации */}
          <div className="card-actions justify-between mt-6">
            <button className="btn btn-outline" onClick={onBack}>
              ← {tCommon('back')}
            </button>
            <button
              className={`btn btn-primary ${!canProceed ? 'btn-disabled' : ''}`}
              onClick={onNext}
              disabled={!canProceed}
            >
              {tCommon('continue')} →
            </button>
          </div>
        </div>
      </div>

      {/* Image Gallery Modal */}
      <ImageGallery
        images={photos}
        initialIndex={galleryIndex}
        isOpen={galleryOpen}
        onClose={() => setGalleryOpen(false)}
      />
    </div>
  );
}
