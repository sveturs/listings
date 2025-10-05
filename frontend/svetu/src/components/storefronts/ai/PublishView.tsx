'use client';

import React, { useState } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import { useCreateAIProduct } from '@/contexts/CreateAIProductContext';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import { toast } from '@/utils/toast';
import { apiClient } from '@/services/api-client';

interface PublishViewProps {
  storefrontId: number | null;
  storefrontSlug: string;
}

export default function PublishView({
  storefrontId,
  storefrontSlug,
}: PublishViewProps) {
  const t = useTranslations('storefronts');
  const locale = useLocale();
  const router = useRouter();
  const { state, setError, setProcessing, setView, reset } =
    useCreateAIProduct();
  const [isPublishing, setIsPublishing] = useState(false);

  // Получаем контент на текущей локали для отображения
  const getLocalizedContent = () => {
    const translations = state.aiData.translations || {};
    if (translations[locale]) {
      return {
        title: translations[locale].title,
        description: translations[locale].description,
      };
    }
    return {
      title: state.aiData.title,
      description: state.aiData.description,
    };
  };

  const localizedContent = getLocalizedContent();

  // Debug logging
  React.useEffect(() => {
    console.log('[PublishView] storefrontId:', storefrontId);
    console.log('[PublishView] translations:', state.aiData.translations);
    console.log(
      '[PublishView] translations count:',
      Object.keys(state.aiData.translations).length
    );
  }, [storefrontId, state.aiData.translations]);

  const handlePublish = async () => {
    if (!storefrontId) {
      setError('Storefront ID not found');
      return;
    }

    setIsPublishing(true);
    setProcessing(true);
    setError(null);

    try {
      // 1. Create product first (without images)
      const productData = {
        name: state.aiData.title,
        description: state.aiData.description,
        price: state.aiData.price,
        currency: state.aiData.currency || 'RSD',
        category_id: state.aiData.categoryId,
        stock_quantity: state.aiData.stockQuantity,
        is_active: true,
        attributes: state.aiData.attributes,
        has_individual_location: !!state.aiData.location,
        individual_address: state.aiData.location?.address,
        individual_latitude: state.aiData.location?.latitude,
        individual_longitude: state.aiData.location?.longitude,

        // Variants
        has_variants: state.aiData.hasVariants,
        variants: state.aiData.hasVariants ? state.aiData.variants : undefined,

        // Translations - convert from {"en": {"title": "...", "description": "..."}} to same format
        translations: state.aiData.translations || {},
      };

      const createResponse = await apiClient.post(
        `/storefronts/slug/${storefrontSlug}/products`,
        productData
      );

      if (!createResponse.data) {
        throw new Error('Failed to create product');
      }

      // Backend возвращает product напрямую, а не обернутый в data
      const productId = createResponse.data.id;

      if (!productId) {
        console.error('[PublishView] Invalid response:', createResponse.data);
        throw new Error('Product ID not returned from server');
      }

      // 2. Upload images to the created product
      const uploadPromises = state.imageFiles.map(async (image, index) => {
        const formData = new FormData();
        formData.append('image', image);

        // Первое изображение делаем главным
        if (index === 0) {
          formData.append('is_main', 'true');
        }
        formData.append('display_order', String(index));

        // НЕ устанавливаем Content-Type вручную для FormData!
        // Браузер сам добавит правильный заголовок с boundary
        return apiClient.post(
          `/storefronts/slug/${storefrontSlug}/products/${productId}/images`,
          formData
        );
      });

      const uploadResults = await Promise.all(uploadPromises);

      // Check if all uploads succeeded
      const failedUploads = uploadResults.filter((res) => !res.data);
      if (failedUploads.length > 0) {
        console.warn(
          `${failedUploads.length} image(s) failed to upload, but product was created`
        );
      }

      toast.success(
        t('productCreatedSuccess') || 'Product created successfully!'
      );

      // Reset context and redirect
      reset();
      router.push(
        `/${locale}/storefronts/${storefrontSlug}/products/${productId}`
      );
    } catch (error: any) {
      console.error('Publish error:', error);
      setError(error.message || 'Failed to publish product');
      toast.error(error.message || 'Failed to publish product');
    } finally {
      setIsPublishing(false);
      setProcessing(false);
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold mb-2">
          {t('publishProduct') || 'Publish Your Product'}
        </h2>
        <p className="text-base-content/70">
          {t('publishDescription') || 'Review final product before publishing'}
        </p>
      </div>

      {/* Product Preview Card */}
      <div className="card bg-base-200 shadow-xl">
        <figure className="relative h-64">
          {state.images[0] && (
            <Image
              src={state.images[0]}
              alt="Product"
              fill
              className="object-cover"
            />
          )}
        </figure>
        <div className="card-body">
          <h3 className="card-title text-xl">{localizedContent.title}</h3>

          <div className="badge badge-primary badge-lg my-2">
            {state.aiData.price} {state.aiData.currency}
          </div>

          <p className="text-base-content/80 line-clamp-3">
            {localizedContent.description}
          </p>

          <div className="divider"></div>

          <div className="grid grid-cols-2 gap-3 text-sm">
            <div>
              <span className="text-base-content/60">
                {t('category') || 'Category'}:
              </span>
              <p className="font-semibold">{state.aiData.category}</p>
            </div>
            <div>
              <span className="text-base-content/60">
                {t('stockQuantity') || 'Stock'}:
              </span>
              <p className="font-semibold">{state.aiData.stockQuantity}</p>
            </div>
            <div>
              <span className="text-base-content/60">
                {t('condition') || 'Condition'}:
              </span>
              <p className="font-semibold capitalize">
                {state.aiData.condition}
              </p>
            </div>
            <div>
              <span className="text-base-content/60">
                {t('images') || 'Images'}:
              </span>
              <p className="font-semibold">{state.images.length}</p>
            </div>
            {state.aiData.hasVariants && (
              <div>
                <span className="text-base-content/60">
                  {t('variants') || 'Variants'}:
                </span>
                <p className="font-semibold">{state.aiData.variants.length}</p>
              </div>
            )}
          </div>

          {state.aiData.location && (
            <div className="mt-3">
              <span className="text-base-content/60">
                {t('location') || 'Location'}:
              </span>
              <p className="font-semibold">{state.aiData.location.address}</p>
            </div>
          )}
        </div>
      </div>

      {/* AI Metadata Info */}
      <div className="alert">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          className="stroke-info shrink-0 w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          />
        </svg>
        <div className="text-sm">
          <div className="font-semibold">
            {t('aiGeneratedContent') || 'AI-Generated Content'}
          </div>
          <div>
            {t('aiTranslationsAvailable', {
              count: Object.keys(state.aiData.translations).length,
            }) ||
              `Translations available in ${Object.keys(state.aiData.translations).length} languages`}
          </div>
        </div>
      </div>

      {/* Translations Preview */}
      {Object.keys(state.aiData.translations).length > 0 && (
        <div className="collapse collapse-arrow bg-base-200">
          <input type="checkbox" />
          <div className="collapse-title font-medium">
            {t('viewTranslations') || 'View Translations'} (
            {Object.keys(state.aiData.translations).length})
          </div>
          <div className="collapse-content">
            <div className="space-y-4">
              {Object.entries(state.aiData.translations).map(
                ([lang, translation]: [string, any]) => (
                  <div key={lang} className="card bg-base-100 shadow">
                    <div className="card-body">
                      <h4 className="card-title text-sm">
                        {lang.toUpperCase()}
                      </h4>
                      <div className="space-y-2">
                        <div>
                          <span className="text-xs text-base-content/60">
                            {t('productName') || 'Title'}:
                          </span>
                          <p className="font-semibold">{translation.title}</p>
                        </div>
                        <div>
                          <span className="text-xs text-base-content/60">
                            {t('description') || 'Description'}:
                          </span>
                          <p className="text-sm">{translation.description}</p>
                        </div>
                        {translation.address && (
                          <div>
                            <span className="text-xs text-base-content/60">
                              {t('address') || 'Address'}:
                            </span>
                            <p className="text-sm">
                              {translation.address}
                              {translation.city && ` (${translation.city})`}
                            </p>
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                )
              )}
            </div>
          </div>
        </div>
      )}

      {/* Variants Preview */}
      {state.aiData.hasVariants && state.aiData.variants.length > 0 && (
        <div className="collapse collapse-arrow bg-base-200">
          <input type="checkbox" />
          <div className="collapse-title font-medium">
            {t('viewVariants') || 'View Variants'} (
            {state.aiData.variants.length})
          </div>
          <div className="collapse-content">
            <div className="overflow-x-auto">
              <table className="table table-sm">
                <thead>
                  <tr>
                    <th>{t('sku') || 'SKU'}</th>
                    <th>{t('attributes') || 'Attributes'}</th>
                    <th>{t('price') || 'Price'}</th>
                    <th>{t('stock') || 'Stock'}</th>
                  </tr>
                </thead>
                <tbody>
                  {state.aiData.variants.map((variant, idx) => (
                    <tr key={idx}>
                      <td>{variant.sku || '-'}</td>
                      <td>
                        {Object.entries(variant.variant_attributes || {})
                          .map(([k, v]) => `${k}: ${v}`)
                          .join(', ') || '-'}
                      </td>
                      <td>{variant.price || state.aiData.price}</td>
                      <td>{variant.stock_quantity}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          </div>
        </div>
      )}

      {/* Actions */}
      <div className="flex justify-between gap-3">
        <button
          onClick={() => setView('variants')}
          disabled={isPublishing}
          className="btn btn-outline"
        >
          {t('back') || 'Back'}
        </button>
        <button
          onClick={handlePublish}
          disabled={isPublishing || !storefrontId}
          className="btn btn-primary btn-lg px-12"
        >
          {isPublishing ? (
            <>
              <span className="loading loading-spinner"></span>
              {t('publishing') || 'Publishing...'}
            </>
          ) : (
            <>
              {t('publishNow') || 'Publish Now'}
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
                  d="M5 13l4 4L19 7"
                />
              </svg>
            </>
          )}
        </button>
      </div>
    </div>
  );
}
