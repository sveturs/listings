'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations, useLocale } from 'next-intl';
import { useCreateProduct } from '@/contexts/CreateProductContext';
import { apiClient } from '@/services/api-client';
import { toast } from '@/utils/toast';
import { getTranslatedAttribute } from '@/utils/translatedAttribute';
import type { components } from '@/types/generated/api';

type CategoryAttribute =
  components['schemas']['backend_internal_domain_models.CategoryAttribute'];

interface PreviewStepProps {
  onBack: () => void;
  storefrontSlug: string;
}

export default function PreviewStep({
  onBack,
  storefrontSlug,
}: PreviewStepProps) {
  const t = useTranslations();
  const locale = useLocale();
  const router = useRouter();
  const { state } = useCreateProduct();
  const [submitting, setSubmitting] = useState(false);
  const [categoryAttributes, setCategoryAttributes] = useState<
    CategoryAttribute[]
  >([]);

  const handleSubmit = async () => {
    try {
      setSubmitting(true);

      // Создаем товар
      const productResponse = await apiClient.post(
        `/api/v1/storefronts/${storefrontSlug}/products`,
        state.productData
      );

      if (productResponse.data) {
        // TODO: Загрузка изображений
        // if (state.images.length > 0) {
        //   await uploadImages(productResponse.data.id, state.images);
        // }

        toast.success(t('storefronts.products.productCreated'));
        router.push(`/${locale}/storefronts/${storefrontSlug}/products`);
      }
    } catch (error: any) {
      console.error('Failed to create product:', error);
      toast.error(
        error.response?.data?.error ||
          t('storefronts.products.errorCreatingProduct')
      );
    } finally {
      setSubmitting(false);
    }
  };

  // Загружаем атрибуты категории для получения метаданных
  useEffect(() => {
    const loadAttributes = async () => {
      if (!state.category) return;

      try {
        const response = await apiClient.get(
          `/api/v1/marketplace/categories/${state.category.id}/attributes`
        );

        if (response.data) {
          const responseData = response.data.data || response.data;
          if (Array.isArray(responseData)) {
            setCategoryAttributes(responseData);
          }
        }
      } catch (error) {
        console.error('Failed to load attributes:', error);
      }
    };

    loadAttributes();
  }, [state.category]);

  const formatPrice = (price: number, currency: string) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency,
    }).format(price);
  };

  const renderAttributeValue = (
    value: any,
    attribute?: CategoryAttribute,
    getOptionLabel?: (value: string) => string
  ): string => {
    if (Array.isArray(value)) {
      return value
        .map((v) => (getOptionLabel ? getOptionLabel(v) : v))
        .join(', ');
    }
    if (typeof value === 'boolean') {
      return value ? t('common.yes') : t('common.no');
    }

    // Для числовых значений добавляем единицы измерения если есть
    if (typeof value === 'number' && attribute?.attribute_type === 'number') {
      const options = attribute.options as any;
      if (options?.unit) {
        return `${value} ${options.unit}`;
      }
    }

    // Для select опций используем переведенное значение
    if (
      typeof value === 'string' &&
      attribute?.attribute_type === 'select' &&
      getOptionLabel
    ) {
      return getOptionLabel(value);
    }

    return String(value || '');
  };

  return (
    <div className="max-w-6xl mx-auto">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold text-base-content mb-4">
          {t('storefronts.products.previewProduct')}
        </h2>
        <p className="text-lg text-base-content/70">
          {t('storefronts.products.previewProductDescription')}
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        {/* Левая колонка - изображения */}
        <div className="space-y-6">
          {/* Главное изображение */}
          {state.images.length > 0 && (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body p-0">
                <div className="aspect-square bg-base-200 rounded-t-2xl overflow-hidden">
                  <img
                    src={URL.createObjectURL(state.images[0])}
                    alt={state.productData.name}
                    className="w-full h-full object-cover"
                  />
                </div>

                {/* Миниатюры */}
                {state.images.length > 1 && (
                  <div className="p-4">
                    <div className="grid grid-cols-4 gap-2">
                      {state.images.slice(1, 5).map((image, index) => (
                        <div
                          key={index}
                          className="aspect-square bg-base-200 rounded-lg overflow-hidden"
                        >
                          <img
                            src={URL.createObjectURL(image)}
                            alt={`${state.productData.name} ${index + 2}`}
                            className="w-full h-full object-cover"
                          />
                        </div>
                      ))}
                      {state.images.length > 5 && (
                        <div className="aspect-square bg-base-200 rounded-lg flex items-center justify-center">
                          <span className="text-sm font-semibold">
                            +{state.images.length - 4}
                          </span>
                        </div>
                      )}
                    </div>
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Статистика */}
          <div className="stats shadow">
            <div className="stat">
              <div className="stat-title">
                {t('storefronts.products.totalPhotos')}
              </div>
              <div className="stat-value text-primary">
                {state.images.length}
              </div>
            </div>
            <div className="stat">
              <div className="stat-title">
                {t('storefronts.products.category')}
              </div>
              <div className="stat-value text-sm">{state.category?.name}</div>
            </div>
          </div>
        </div>

        {/* Правая колонка - информация о товаре */}
        <div className="space-y-6">
          {/* Основная информация */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h3 className="card-title text-2xl mb-4">
                {state.productData.name}
              </h3>

              <div className="badge badge-secondary mb-4">
                {state.category?.name}
              </div>

              <p className="text-base-content/80 mb-6">
                {state.productData.description}
              </p>

              {/* Цена */}
              <div className="flex items-baseline gap-2 mb-6">
                <span className="text-4xl font-bold text-primary">
                  {formatPrice(
                    state.productData.price || 0,
                    state.productData.currency || 'RSD'
                  )}
                </span>
                {state.productData.stock_quantity !== undefined &&
                  state.productData.stock_quantity > 0 && (
                    <span className="text-sm text-success">
                      {state.productData.stock_quantity}{' '}
                      {t('storefronts.products.inStock')}
                    </span>
                  )}
              </div>

              {/* Детали */}
              <div className="grid grid-cols-2 gap-4 text-sm">
                {state.productData.sku && (
                  <div>
                    <span className="font-semibold">
                      {t('storefronts.products.sku')}:
                    </span>
                    <span className="ml-2">{state.productData.sku}</span>
                  </div>
                )}
                {state.productData.barcode && (
                  <div>
                    <span className="font-semibold">
                      {t('storefronts.products.barcode')}:
                    </span>
                    <span className="ml-2">{state.productData.barcode}</span>
                  </div>
                )}
                <div>
                  <span className="font-semibold">
                    {t('storefronts.products.status')}:
                  </span>
                  <span
                    className={`ml-2 ${state.productData.is_active ? 'text-success' : 'text-error'}`}
                  >
                    {state.productData.is_active
                      ? t('storefronts.products.active')
                      : t('storefronts.products.inactive')}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Атрибуты */}
          {Object.keys(state.attributes).length > 0 && (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h4 className="card-title text-lg mb-4 flex items-center gap-2">
                  <span className="text-xl">⚡</span>
                  {t('storefronts.products.specifications')}
                </h4>

                <div className="space-y-3">
                  {Object.entries(state.attributes).map(([id, value]) => {
                    if (!value || (Array.isArray(value) && value.length === 0))
                      return null;

                    // Находим атрибут для получения переведенного названия
                    const attribute = categoryAttributes.find(
                      (attr) => attr.id === parseInt(id)
                    );
                    const { displayName, getOptionLabel } = attribute && attribute.id
                      ? getTranslatedAttribute(attribute as any, locale)
                      : {
                          displayName: `Attribute ${id}`,
                          getOptionLabel: (v: string) => v,
                        };

                    return (
                      <div
                        key={id}
                        className="flex justify-between items-center py-2 border-b border-base-200 last:border-b-0"
                      >
                        <span className="font-medium text-base-content/70">
                          {displayName}:
                        </span>
                        <span className="text-base-content">
                          {renderAttributeValue(
                            value,
                            attribute,
                            getOptionLabel
                          )}
                        </span>
                      </div>
                    );
                  })}
                </div>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Подтверждение */}
      <div className="card bg-gradient-to-r from-primary/10 to-secondary/10 border border-primary/20 shadow-xl mb-8">
        <div className="card-body">
          <div className="flex items-center gap-4">
            <div className="text-5xl">🎉</div>
            <div>
              <h3 className="text-xl font-bold text-primary mb-2">
                {t('storefronts.products.readyToPublish')}
              </h3>
              <p className="text-base-content/70">
                {t('storefronts.products.publishConfirmation')}
              </p>
            </div>
          </div>
        </div>
      </div>

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
          {t('common.back')}
        </button>

        <div className="flex gap-4">
          <button className="btn btn-outline btn-lg px-8" disabled={submitting}>
            {t('storefronts.products.saveDraft')}
          </button>

          <button
            onClick={handleSubmit}
            disabled={submitting}
            className="btn btn-primary btn-lg px-8"
          >
            {submitting ? (
              <>
                <span className="loading loading-spinner loading-sm mr-2"></span>
                {t('storefronts.products.creating')}
              </>
            ) : (
              <>
                {t('storefronts.products.publishProduct')}
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
    </div>
  );
}
