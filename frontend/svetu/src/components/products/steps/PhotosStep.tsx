'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateProduct } from '@/contexts/CreateProductContext';
import Image from 'next/image';
import { toast } from '@/utils/toast';

interface PhotosStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function PhotosStep({ onNext, onBack }: PhotosStepProps) {
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const { state, setImages, setError, clearError } = useCreateProduct();
  const [images, setImagesState] = useState<File[]>(state.images || []);
  const [previews, setPreviews] = useState<string[]>([]);

  useEffect(() => {
    // Создаем превью для существующих изображений
    const newPreviews: string[] = [];
    images.forEach((file) => {
      const reader = new FileReader();
      reader.onloadend = () => {
        newPreviews.push(reader.result as string);
        if (newPreviews.length === images.length) {
          setPreviews([...newPreviews]);
        }
      };
      reader.readAsDataURL(file);
    });

    if (images.length === 0) {
      setPreviews([]);
    }
  }, [images]);

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);

    // Проверяем общее количество
    if (files.length + images.length > 10) {
      toast.error(t('maxImagesError'));
      return;
    }

    // Валидируем файлы
    const invalidFiles = files.filter((file) => {
      const isValidType = file.type.startsWith('image/');
      const isValidSize = file.size <= 10 * 1024 * 1024; // 10MB
      return !isValidType || !isValidSize;
    });

    if (invalidFiles.length > 0) {
      toast.error(t('invalidImagesError'));
      return;
    }

    const newImages = [...images, ...files];
    setImagesState(newImages);
    setImages(newImages);
    clearError('images');
  };

  const removeImage = (index: number) => {
    const newImages = images.filter((_, i) => i !== index);
    const newPreviews = previews.filter((_, i) => i !== index);
    setImagesState(newImages);
    setImages(newImages);
    setPreviews(newPreviews);
  };

  const moveImage = (fromIndex: number, toIndex: number) => {
    const newImages = [...images];
    const newPreviews = [...previews];

    [newImages[fromIndex], newImages[toIndex]] = [
      newImages[toIndex],
      newImages[fromIndex],
    ];
    [newPreviews[fromIndex], newPreviews[toIndex]] = [
      newPreviews[toIndex],
      newPreviews[fromIndex],
    ];

    setImagesState(newImages);
    setImages(newImages);
    setPreviews(newPreviews);
  };

  const validateImages = (): boolean => {
    if (images.length === 0) {
      setError('images', t('imagesRequired'));
      return false;
    }
    return true;
  };

  const handleNext = () => {
    console.log('PhotosStep handleNext called');
    console.log('Images count:', images.length);
    if (validateImages()) {
      console.log('Validation passed, calling onNext');
      onNext();
    } else {
      console.log('Validation failed');
    }
  };

  return (
    <div className="max-w-6xl mx-auto">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold text-base-content mb-4">
          {t('productPhotos')}
        </h2>
        <p className="text-lg text-base-content/70">
          {t('productPhotosDescription')}
        </p>
      </div>

      {/* Статистика */}
      <div className="stats shadow mb-8 w-full">
        <div className="stat">
          <div className="stat-figure text-primary">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="inline-block w-8 h-8 stroke-current"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
              ></path>
            </svg>
          </div>
          <div className="stat-title">
            {t('photosUploaded')}
          </div>
          <div className="stat-value text-primary">{images.length}</div>
          <div className="stat-desc">
            {t('maxPhotos')}: 10
          </div>
        </div>

        <div className="stat">
          <div className="stat-figure text-secondary">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              className="inline-block w-8 h-8 stroke-current"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 10V3L4 14h7v7l9-11h-7z"
              ></path>
            </svg>
          </div>
          <div className="stat-title">
            {t('recommendedSize')}
          </div>
          <div className="stat-value text-secondary">1200px</div>
          <div className="stat-desc">
            {t('aspectRatio')}: 1:1
          </div>
        </div>
      </div>

      {/* Загрузка изображений */}
      <div className="card bg-base-100 shadow-xl mb-8">
        <div className="card-body">
          <h3 className="card-title text-xl mb-4 flex items-center gap-2">
            <span className="text-2xl">📸</span>
            {t('uploadPhotos')}
          </h3>

          {/* Дроп-зона */}
          <div className="border-2 border-dashed border-base-300 rounded-2xl p-8 text-center hover:border-primary transition-colors cursor-pointer">
            <input
              type="file"
              accept="image/jpeg,image/png,image/gif,image/webp"
              multiple
              onChange={handleImageChange}
              className="hidden"
              id="image-upload"
            />
            <label htmlFor="image-upload" className="cursor-pointer">
              <div className="text-6xl mb-4">📷</div>
              <h4 className="text-xl font-semibold mb-2">
                {t('dragDropPhotos')}
              </h4>
              <p className="text-base-content/60 mb-4">
                {t('supportedFormats')}: JPEG, PNG, GIF,
                WebP
              </p>
              <div className="btn btn-primary">
                {t('chooseFiles')}
              </div>
            </label>
          </div>

          {/* Ошибка */}
          {state.errors.images && (
            <div className="alert alert-error mt-4">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                className="stroke-current shrink-0 h-6 w-6"
                fill="none"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth="2"
                  d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
              <span>{state.errors.images}</span>
            </div>
          )}
        </div>
      </div>

      {/* Превью изображений */}
      {images.length > 0 && (
        <div className="card bg-base-100 shadow-xl mb-8">
          <div className="card-body">
            <h3 className="card-title text-xl mb-4 flex items-center gap-2">
              <span className="text-2xl">🖼️</span>
              {t('photoPreview')}
            </h3>

            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4">
              {previews.map((preview, index) => (
                <div key={index} className="relative group">
                  <div className="aspect-square bg-base-200 rounded-xl overflow-hidden">
                    <Image
                      src={preview}
                      alt={`Preview ${index + 1}`}
                      fill
                      className="object-cover"
                    />
                  </div>

                  {/* Overlay с кнопками */}
                  <div className="absolute inset-0 bg-black/50 opacity-0 group-hover:opacity-100 transition-opacity rounded-xl flex items-center justify-center gap-2">
                    {/* Главное фото */}
                    {index === 0 && (
                      <div className="badge badge-primary">
                        {t('mainPhoto')}
                      </div>
                    )}

                    {/* Кнопки управления */}
                    <div className="absolute top-2 right-2 flex gap-1">
                      {index > 0 && (
                        <button
                          onClick={() => moveImage(index, index - 1)}
                          className="btn btn-circle btn-sm btn-ghost text-white hover:bg-white/20"
                        >
                          ←
                        </button>
                      )}
                      {index < images.length - 1 && (
                        <button
                          onClick={() => moveImage(index, index + 1)}
                          className="btn btn-circle btn-sm btn-ghost text-white hover:bg-white/20"
                        >
                          →
                        </button>
                      )}
                      <button
                        onClick={() => removeImage(index)}
                        className="btn btn-circle btn-sm btn-error text-white"
                      >
                        ×
                      </button>
                    </div>
                  </div>

                  {/* Номер фото */}
                  <div className="absolute bottom-2 left-2 bg-black/70 text-white text-xs px-2 py-1 rounded">
                    {index + 1}
                  </div>
                </div>
              ))}
            </div>

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
              <span className="text-sm">
                💡 {t('photoTips')}
              </span>
            </div>
          </div>
        </div>
      )}

      {/* Навигация */}
      <div className="flex justify-between items-center">
        <button onClick={onBack} className="btn btn-outline btn-lg px-8">
          <svg
            className="w-5 h-5 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M15 19l-7-7 7-7"
            />
          </svg>
          {tCommon('back')}
        </button>

        <button
          onClick={handleNext}
          disabled={images.length === 0}
          className={`btn btn-lg px-8 ${images.length > 0 ? 'btn-primary' : 'btn-disabled'}`}
        >
          {tCommon('next')}
          <svg
            className="w-5 h-5 ml-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
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
    </div>
  );
}
