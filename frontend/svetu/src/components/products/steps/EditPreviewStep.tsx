'use client';

import { useState } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useRouter } from 'next/navigation';
import { useEditProduct } from '@/contexts/EditProductContext';
import { storefrontProductsService } from '@/services/storefrontProducts';
import Image from 'next/image';
import {
  EyeIcon,
  PencilIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
  MapPinIcon,
  TagIcon,
  CurrencyDollarIcon,
  PhotoIcon,
} from '@heroicons/react/24/outline';

interface EditPreviewStepProps {
  onBack: () => void;
  storefrontSlug: string;
  productId: number;
}

export default function EditPreviewStep({
  onBack,
  storefrontSlug,
  productId,
}: EditPreviewStepProps) {
  const t = useTranslations('storefronts');
  const tStorefronts.products.errors = useTranslations('storefronts');
  const tStorefronts.products.steps = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const router = useRouter();
  const { state, dispatch, goToStep } = useEditProduct();
  const [isSubmitting, setIsSubmitting] = useState(false);

  // Подготовка данных для отправки
  const getVisibleImages = () => {
    return state.existingImages.filter(
      (img) => !state.imagesToDelete.includes(img.id)
    );
  };

  const getMainImageUrl = () => {
    const visibleExisting = getVisibleImages();
    const mainExisting = visibleExisting.find((img) => img.is_main);
    if (mainExisting) return mainExisting.url;

    if (state.newImages.length > 0) {
      return URL.createObjectURL(state.newImages[0]);
    }

    return visibleExisting[0]?.url || '/images/placeholder-product.jpg';
  };

  const handleSubmit = async () => {
    setIsSubmitting(true);
    dispatch({ type: 'SET_SAVING', payload: true });

    try {
      // 1. Обновляем основную информацию о продукте
      await storefrontProductsService.updateProduct(
        storefrontSlug,
        productId,
        state.productData
      );

      // 2. Загружаем новые изображения если есть
      if (state.newImages.length > 0) {
        console.log('New images to upload:', state.newImages);
        const hasMainExisting = getVisibleImages().some((img) => img.is_main);
        const mainImageIndex = hasMainExisting ? undefined : 0; // Первое новое изображение как главное, если нет главного среди существующих

        await storefrontProductsService.uploadProductImages(
          storefrontSlug,
          productId,
          state.newImages,
          mainImageIndex
        );
      } else {
        console.log('No new images to upload');
      }

      // 3. Удаляем помеченные изображения
      for (const imageId of state.imagesToDelete) {
        await storefrontProductsService.deleteProductImage(
          storefrontSlug,
          productId,
          imageId
        );
      }

      // 4. Обновляем главное изображение среди существующих если изменилось
      const visibleExisting = getVisibleImages();
      const currentMain = state.originalProduct?.images?.find(
        (img) => img.is_default
      );
      const newMain = visibleExisting.find((img) => img.is_main);

      if (newMain && (!currentMain || currentMain.id !== newMain.id)) {
        await storefrontProductsService.setMainImage(
          storefrontSlug,
          productId,
          newMain.id
        );
      }

      dispatch({ type: 'SET_UNSAVED_CHANGES', payload: false });

      // Перенаправляем к списку продуктов
      router.push(`/${locale}/storefronts/${storefrontSlug}/products`);
    } catch (error: any) {
      console.error('Error updating product:', error);
      dispatch({
        type: 'SET_ERROR',
        payload: {
          field: 'submit',
          message:
            error.message || tStorefronts.products.errors('updateFailed'),
        },
      });
    } finally {
      setIsSubmitting(false);
      dispatch({ type: 'SET_SAVING', payload: false });
    }
  };

  const visibleImages = getVisibleImages();
  const totalImages = visibleImages.length + state.newImages.length;

  return (
    <div className="space-y-8">
      {/* Заголовок */}
      <div className="text-center">
        <EyeIcon className="w-16 h-16 text-primary mx-auto mb-4" />
        <h3 className="text-2xl font-bold text-base-content mb-2">
          {tStorefronts.products.steps('preview')}
        </h3>
        <p className="text-base-content/70">
          {t('previewStepDescription')}
        </p>
      </div>

      {/* Предупреждение о несохраненных изменениях */}
      {state.hasUnsavedChanges && (
        <div className="bg-warning/10 border border-warning rounded-2xl p-4">
          <div className="flex items-start gap-3">
            <ExclamationTriangleIcon className="w-6 h-6 text-warning flex-shrink-0 mt-0.5" />
            <div>
              <h4 className="font-semibold text-warning">
                {t('unsavedChanges')}
              </h4>
              <p className="text-warning/80 text-sm mt-1">
                {t('unsavedChangesDescription')}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Превью продукта */}
      <div className="bg-base-100 rounded-3xl shadow-xl overflow-hidden">
        {/* Изображение */}
        <div className="relative">
          <div className="aspect-[16/10] relative">
            <Image
              src={getMainImageUrl()}
              alt={state.productData.name || ''}
              fill
              className="object-cover"
            />
            {totalImages > 1 && (
              <div className="absolute bottom-4 right-4 bg-black/60 text-white px-3 py-1 rounded-full text-sm backdrop-blur-sm">
                +{totalImages - 1} {t('morePhotos')}
              </div>
            )}
          </div>
        </div>

        {/* Информация о продукте */}
        <div className="p-6 space-y-4">
          {/* Название и цена */}
          <div className="flex items-start justify-between gap-4">
            <div className="flex-1">
              <h3 className="text-2xl font-bold text-base-content mb-2">
                {state.productData.name}
              </h3>
              <div className="flex items-center gap-2 text-base-content/70">
                <TagIcon className="w-4 h-4" />
                <span>{state.category?.name}</span>
              </div>
            </div>
            <div className="text-right">
              <div className="flex items-center gap-2 text-2xl font-bold text-primary">
                <CurrencyDollarIcon className="w-6 h-6" />
                <span>{state.productData.price?.toLocaleString()}</span>
              </div>
              <p className="text-sm text-base-content/60">RSD</p>
            </div>
          </div>

          {/* Описание */}
          <div>
            <h4 className="font-semibold text-base-content mb-2">
              {t('description')}
            </h4>
            <p className="text-base-content/80 whitespace-pre-wrap">
              {state.productData.description}
            </p>
          </div>

          {/* Местоположение */}
          {state.location && (
            <div className="flex items-center gap-3 p-3 bg-base-200 rounded-xl">
              <MapPinIcon className="w-5 h-5 text-base-content/60" />
              <div>
                <p className="font-medium text-base-content">
                  {state.location.useStorefrontLocation
                    ? t('storefrontLocation')
                    : state.location.individualAddress}
                </p>
                <p className="text-sm text-base-content/60">
                  {state.location.privacyLevel &&
                    t(
                      `storefronts.products.privacy.${state.location.privacyLevel}`
                    )}
                </p>
              </div>
            </div>
          )}

          {/* Атрибуты */}
          {Object.keys(state.attributes).length > 0 && (
            <div>
              <h4 className="font-semibold text-base-content mb-3">
                {t('attributes')}
              </h4>
              <div className="grid grid-cols-2 gap-3">
                {Object.entries(state.attributes).map(([key, value]) => (
                  <div key={key} className="bg-base-200 rounded-xl p-3">
                    <p className="text-sm text-base-content/60 mb-1">
                      {/* Здесь должно быть название атрибута из категории */}
                      {key}
                    </p>
                    <p className="font-medium text-base-content">{value}</p>
                  </div>
                ))}
              </div>
            </div>
          )}

          {/* Дополнительная информация */}
          <div className="flex flex-wrap gap-4 pt-4 border-t border-base-300">
            <div className="flex items-center gap-2">
              <span className="text-sm text-base-content/60">
                {t('stockQuantity')}:
              </span>
              <span className="font-medium text-base-content">
                {state.productData.stock_quantity}
              </span>
            </div>
            {state.productData.sku && (
              <div className="flex items-center gap-2">
                <span className="text-sm text-base-content/60">SKU:</span>
                <span className="font-medium text-base-content">
                  {state.productData.sku}
                </span>
              </div>
            )}
            {state.productData.barcode && (
              <div className="flex items-center gap-2">
                <span className="text-sm text-base-content/60">
                  {t('barcode')}:
                </span>
                <span className="font-medium text-base-content">
                  {state.productData.barcode}
                </span>
              </div>
            )}
          </div>
        </div>
      </div>

      {/* Статус активности */}
      <div className="bg-base-200 rounded-2xl p-6">
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div
              className={`w-12 h-12 rounded-full flex items-center justify-center ${
                state.productData.is_active ? 'bg-success/20' : 'bg-error/20'
              }`}
            >
              <CheckCircleIcon
                className={`w-6 h-6 ${
                  state.productData.is_active ? 'text-success' : 'text-error'
                }`}
              />
            </div>
            <div>
              <p className="font-semibold text-base-content">
                {state.productData.is_active
                  ? t('active')
                  : t('inactive')}
              </p>
              <p className="text-sm text-base-content/60">
                {state.productData.is_active
                  ? t('visibleToCustomers')
                  : t('hiddenFromCustomers')}
              </p>
            </div>
          </div>
          <button
            onClick={() => goToStep(1)}
            className="btn btn-outline btn-sm"
            disabled={isSubmitting}
          >
            <PencilIcon className="w-4 h-4 mr-2" />
            {tCommon('edit')}
          </button>
        </div>
      </div>

      {/* Изменения в изображениях */}
      {(state.newImages.length > 0 || state.imagesToDelete.length > 0) && (
        <div className="bg-info/10 border border-info rounded-2xl p-6">
          <div className="flex items-start gap-3">
            <PhotoIcon className="w-6 h-6 text-info flex-shrink-0 mt-0.5" />
            <div>
              <h4 className="font-semibold text-info mb-2">
                {t('imageChanges')}
              </h4>
              <div className="space-y-1">
                {state.newImages.length > 0 && (
                  <p className="text-info/80 text-sm">
                    • {state.newImages.length}{' '}
                    {t('newImagesToAdd')}
                  </p>
                )}
                {state.imagesToDelete.length > 0 && (
                  <p className="text-info/80 text-sm">
                    • {state.imagesToDelete.length}{' '}
                    {t('imagesToDelete')}
                  </p>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Ошибки */}
      {state.errors.submit && (
        <div className="alert alert-error">
          <ExclamationTriangleIcon className="w-6 h-6" />
          <span>{state.errors.submit}</span>
        </div>
      )}

      {/* Кнопки навигации */}
      <div className="flex justify-between">
        <button
          onClick={onBack}
          className="btn btn-outline btn-lg"
          disabled={isSubmitting}
        >
          {tCommon('back')}
        </button>
        <button
          onClick={handleSubmit}
          className="btn btn-primary btn-lg"
          disabled={isSubmitting}
        >
          {isSubmitting ? (
            <>
              <span className="loading loading-spinner loading-sm mr-2"></span>
              {t('updating')}
            </>
          ) : (
            <>
              <CheckCircleIcon className="w-5 h-5 mr-2" />
              {t('updateProduct')}
            </>
          )}
        </button>
      </div>
    </div>
  );
}
