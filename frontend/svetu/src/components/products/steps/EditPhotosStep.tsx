'use client';

import { useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useEditProduct } from '@/contexts/EditProductContext';
import Image from 'next/image';
import {
  PhotoIcon,
  PlusIcon,
  TrashIcon,
  StarIcon,
  ArrowUpTrayIcon,
} from '@heroicons/react/24/outline';
import { StarIcon as StarSolidIcon } from '@heroicons/react/24/solid';

interface EditPhotosStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function EditPhotosStep({
  onNext,
  onBack,
}: EditPhotosStepProps) {
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const {
    state,
    addNewImage,
    removeNewImage,
    markImageForDeletion,
    restoreImage,
    setMainImage,
    completeStep,
    setError,
    clearError,
  } = useEditProduct();

  const [isDragOver, setIsDragOver] = useState(false);

  // Получаем все изображения (существующие + новые)
  const visibleExistingImages = state.existingImages.filter(
    (img) => !state.imagesToDelete.includes(img.id)
  );
  const totalImages = visibleExistingImages.length + state.newImages.length;

  const handleFileSelect = useCallback(
    (files: FileList | null) => {
      if (!files) return;

      clearError('images');

      for (let i = 0; i < files.length; i++) {
        const file = files[i];

        // Проверка типа файла
        if (!file.type.startsWith('image/')) {
          setError('images', t('products.errors.invalidImageType'));
          continue;
        }

        // Проверка размера файла (максимум 10MB)
        if (file.size > 10 * 1024 * 1024) {
          setError('images', t('products.errors.imageTooLarge'));
          continue;
        }

        addNewImage(file);
      }
    },
    [addNewImage, setError, clearError, t]
  );

  const handleDragOver = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(true);
  }, []);

  const handleDragLeave = useCallback((e: React.DragEvent) => {
    e.preventDefault();
    setIsDragOver(false);
  }, []);

  const handleDrop = useCallback(
    (e: React.DragEvent) => {
      e.preventDefault();
      setIsDragOver(false);

      const files = e.dataTransfer.files;
      handleFileSelect(files);
    },
    [handleFileSelect]
  );

  const handleNext = () => {
    if (totalImages === 0) {
      setError('images', t('products.errors.atLeastOneImage'));
      return;
    }

    completeStep(4);
    onNext();
  };

  const getImageUrl = (file: File): string => {
    return URL.createObjectURL(file);
  };

  const isMainImage = (type: 'existing' | 'new', index: number): boolean => {
    if (type === 'existing') {
      return state.existingImages[index]?.is_main || false;
    }
    // Для новых изображений главное определяется позицией (первое = главное если нет главного среди существующих)
    const hasMainExisting = visibleExistingImages.some((img) => img.is_main);
    return !hasMainExisting && index === 0;
  };

  return (
    <div className="space-y-6">
      {/* Заголовок */}
      <div className="text-center">
        <PhotoIcon className="w-16 h-16 text-primary mx-auto mb-4" />
        <h3 className="text-2xl font-bold text-base-content mb-2">
          {t('products.steps.photos')}
        </h3>
        <p className="text-base-content/70">{t('photosStepDescription')}</p>
      </div>

      {/* Статистика изображений */}
      <div className="bg-base-200 rounded-2xl p-4">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-primary/20 rounded-full flex items-center justify-center">
              <PhotoIcon className="w-5 h-5 text-primary" />
            </div>
            <div>
              <p className="font-semibold text-base-content">
                {t('totalImages')}
              </p>
              <p className="text-sm text-base-content/60">
                {totalImages} {t('images')}
              </p>
            </div>
          </div>
          {state.imagesToDelete.length > 0 && (
            <div className="text-right">
              <p className="text-sm text-error">
                {state.imagesToDelete.length} {t('toDelete')}
              </p>
            </div>
          )}
        </div>
      </div>

      {/* Зона загрузки */}
      <div
        className={`
          relative border-2 border-dashed rounded-2xl p-8 text-center transition-all
          ${isDragOver ? 'border-primary bg-primary/5' : 'border-base-300'}
          ${totalImages >= 10 ? 'opacity-50 cursor-not-allowed' : 'hover:border-primary hover:bg-primary/5 cursor-pointer'}
        `}
        onDragOver={handleDragOver}
        onDragLeave={handleDragLeave}
        onDrop={handleDrop}
        onClick={() => {
          if (totalImages < 10) {
            const input = document.createElement('input');
            input.type = 'file';
            input.accept = 'image/*';
            input.multiple = true;
            input.onchange = (e) => {
              const target = e.target as HTMLInputElement;
              handleFileSelect(target.files);
            };
            input.click();
          }
        }}
      >
        <div className="flex flex-col items-center">
          <ArrowUpTrayIcon className="w-12 h-12 text-base-content/40 mb-4" />
          <p className="text-lg font-semibold text-base-content mb-2">
            {totalImages >= 10 ? t('maxImagesReached') : t('dragDropImages')}
          </p>
          <p className="text-base-content/60 text-sm">
            {t('imageRequirements')}
          </p>
        </div>
      </div>

      {/* Список изображений */}
      {(visibleExistingImages.length > 0 || state.newImages.length > 0) && (
        <div className="space-y-4">
          <h4 className="text-lg font-semibold text-base-content">
            {t('selectedImages')}
          </h4>

          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
            {/* Существующие изображения */}
            {state.existingImages.map((image, index) => {
              const isMarkedForDeletion = state.imagesToDelete.includes(
                image.id
              );
              const isMain = image.is_main;

              return (
                <div
                  key={`existing-${image.id}`}
                  className={`
                    relative group rounded-2xl overflow-hidden border-2 transition-all
                    ${isMarkedForDeletion ? 'border-error bg-error/10 opacity-50' : 'border-base-300'}
                    ${isMain ? 'ring-2 ring-primary ring-offset-2' : ''}
                  `}
                >
                  <div className="aspect-square relative">
                    {image.url ? (
                      <Image
                        src={image.url}
                        alt={`Product image ${index + 1}`}
                        fill
                        className="object-cover"
                      />
                    ) : (
                      <div className="w-full h-full bg-base-200 flex items-center justify-center">
                        <span className="text-base-content/50">No image</span>
                      </div>
                    )}

                    {/* Оверлей действий */}
                    <div className="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-all flex items-center justify-center opacity-0 group-hover:opacity-100">
                      <div className="flex gap-2">
                        {!isMarkedForDeletion && (
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              setMainImage('existing', index);
                            }}
                            className="btn btn-sm btn-circle bg-primary/80 hover:bg-primary border-none text-primary-content"
                            title={t('setAsMain')}
                          >
                            {isMain ? (
                              <StarSolidIcon className="w-4 h-4" />
                            ) : (
                              <StarIcon className="w-4 h-4" />
                            )}
                          </button>
                        )}

                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            if (isMarkedForDeletion) {
                              restoreImage(image.id);
                            } else {
                              markImageForDeletion(image.id);
                            }
                          }}
                          className={`btn btn-sm btn-circle border-none text-white ${
                            isMarkedForDeletion
                              ? 'bg-success/80 hover:bg-success'
                              : 'bg-error/80 hover:bg-error'
                          }`}
                          title={
                            isMarkedForDeletion
                              ? t('restoreImage')
                              : t('deleteImage')
                          }
                        >
                          {isMarkedForDeletion ? (
                            <PlusIcon className="w-4 h-4" />
                          ) : (
                            <TrashIcon className="w-4 h-4" />
                          )}
                        </button>
                      </div>
                    </div>
                  </div>

                  {/* Индикаторы */}
                  <div className="absolute top-2 left-2 flex gap-1">
                    {isMain && !isMarkedForDeletion && (
                      <span className="px-2 py-1 bg-primary text-primary-content text-xs rounded-full font-medium">
                        {t('main')}
                      </span>
                    )}
                    {isMarkedForDeletion && (
                      <span className="px-2 py-1 bg-error text-error-content text-xs rounded-full font-medium">
                        {t('toDelete')}
                      </span>
                    )}
                  </div>
                </div>
              );
            })}

            {/* Новые изображения */}
            {state.newImages.map((image, index) => {
              const isMain = isMainImage('new', index);

              return (
                <div
                  key={`new-${index}`}
                  className={`
                    relative group rounded-2xl overflow-hidden border-2 border-success
                    ${isMain ? 'ring-2 ring-primary ring-offset-2' : ''}
                  `}
                >
                  <div className="aspect-square relative">
                    <Image
                      src={getImageUrl(image)}
                      alt={`New product image ${index + 1}`}
                      fill
                      className="object-cover"
                    />

                    {/* Оверлей действий */}
                    <div className="absolute inset-0 bg-black/0 group-hover:bg-black/40 transition-all flex items-center justify-center opacity-0 group-hover:opacity-100">
                      <div className="flex gap-2">
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setMainImage('new', index);
                          }}
                          className="btn btn-sm btn-circle bg-primary/80 hover:bg-primary border-none text-primary-content"
                          title={t('setAsMain')}
                        >
                          {isMain ? (
                            <StarSolidIcon className="w-4 h-4" />
                          ) : (
                            <StarIcon className="w-4 h-4" />
                          )}
                        </button>

                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            removeNewImage(index);
                          }}
                          className="btn btn-sm btn-circle bg-error/80 hover:bg-error border-none text-white"
                          title={t('removeImage')}
                        >
                          <TrashIcon className="w-4 h-4" />
                        </button>
                      </div>
                    </div>
                  </div>

                  {/* Индикаторы */}
                  <div className="absolute top-2 left-2 flex gap-1">
                    <span className="px-2 py-1 bg-success text-success-content text-xs rounded-full font-medium">
                      {t('new')}
                    </span>
                    {isMain && (
                      <span className="px-2 py-1 bg-primary text-primary-content text-xs rounded-full font-medium">
                        {t('main')}
                      </span>
                    )}
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      )}

      {/* Ошибки */}
      {state.errors.images && (
        <div className="alert alert-error">
          <svg
            className="w-6 h-6"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.96-.833-2.73 0L3.084 16.5c-.77.833.192 2.5 1.732 2.5z"
            />
          </svg>
          <span>{state.errors.images}</span>
        </div>
      )}

      {/* Кнопки навигации */}
      <div className="flex justify-between">
        <button
          onClick={onBack}
          className="btn btn-outline btn-lg"
          disabled={state.isSaving}
        >
          {tCommon('back')}
        </button>
        <button
          onClick={handleNext}
          className="btn btn-primary btn-lg"
          disabled={totalImages === 0 || state.isSaving}
        >
          {tCommon('continue')}
        </button>
      </div>
    </div>
  );
}
