'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useTranslations, useLocale } from 'next-intl';
import { useCreateProduct } from '@/contexts/CreateProductContext';
import { apiClient } from '@/services/api-client';
import { toast } from '@/utils/toast';
import { getTranslatedAttribute } from '@/utils/translatedAttribute';
import Image from 'next/image';
import type { components } from '@/types/generated/api';
import {
  formatAddressWithPrivacy,
  type LocationPrivacyLevel,
} from '@/utils/addressUtils';

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
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const router = useRouter();
  const { state } = useCreateProduct();
  const [submitting, setSubmitting] = useState(false);
  const [categoryAttributes, setCategoryAttributes] = useState<
    CategoryAttribute[]
  >([]);

  const uploadImages = async (productId: number, images: File[]) => {
    const uploadPromises = images.map(async (image, index) => {
      const formData = new FormData();
      formData.append('image', image);

      // Первое изображение делаем главным
      if (index === 0) {
        formData.append('is_main', 'true');
      }
      formData.append('display_order', String(index));

      return apiClient.post(
        `/api/v1/storefronts/slug/${storefrontSlug}/products/${productId}/images`,
        formData
      );
    });

    await Promise.all(uploadPromises);
  };

  // Маппинг значений privacy level из frontend в backend
  const mapPrivacyLevel = (frontendLevel: string): string => {
    const mapping: Record<string, string> = {
      exact: 'exact',
      street: 'street',
      district: 'district',
      city: 'city',
    };
    return mapping[frontendLevel] || 'exact';
  };

  const handleSubmit = async () => {
    try {
      setSubmitting(true);

      // Подготавливаем данные для отправки, включая местоположение и варианты
      const productDataWithLocation = {
        ...state.productData,
        ...(state.location && !state.location.useStorefrontLocation
          ? {
              has_individual_location: true,
              individual_address: state.location.individualAddress,
              individual_latitude: state.location.latitude,
              individual_longitude: state.location.longitude,
              location_privacy: mapPrivacyLevel(
                state.location.privacyLevel || 'exact'
              ),
              show_on_map: state.location.showOnMap,
            }
          : {
              has_individual_location: false,
            }),
        // Добавляем варианты если они есть
        ...(state.hasVariants && state.variants.length > 0
          ? {
              has_variants: true,
              variants: state.variants,
              variant_settings: state.variantSettings,
            }
          : {
              has_variants: false,
            }),
      };

      // Создаем товар
      const productResponse = await apiClient.post(
        `/api/v1/storefronts/slug/${storefrontSlug}/products`,
        productDataWithLocation
      );

      if (productResponse.data) {
        const productData = productResponse.data.data || productResponse.data;

        // Загрузка изображений
        if (state.images.length > 0 && productData.id) {
          try {
            await uploadImages(productData.id, state.images);
          } catch (imageError) {
            console.error('Failed to upload images:', imageError);
            // Продолжаем даже если не удалось загрузить изображения
            toast.warning(t('productCreatedButImagesError'));
          }
        }

        toast.success(t('productCreated'));
        router.push(`/${locale}/storefronts/${storefrontSlug}/products`);
      }
    } catch (error: any) {
      console.error('Failed to create product:', error);
      toast.error(error.response?.data?.error || t('errorCreatingProduct'));
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
      return value ? tCommon('yes') : tCommon('no');
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
          {t('previewProduct')}
        </h2>
        <p className="text-lg text-base-content/70">
          {t('previewProductDescription')}
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
                  <Image
                    src={URL.createObjectURL(state.images[0])}
                    alt={state.productData.name}
                    fill
                    className="object-cover"
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
                          <Image
                            src={URL.createObjectURL(image)}
                            alt={`${state.productData.name} ${index + 2}`}
                            fill
                            className="object-cover"
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
              <div className="stat-title">{t('totalPhotos')}</div>
              <div className="stat-value text-primary">
                {state.images.length}
              </div>
            </div>
            <div className="stat">
              <div className="stat-title">{t('category')}</div>
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
                      {state.productData.stock_quantity} {t('inStock')}
                    </span>
                  )}
              </div>

              {/* Детали */}
              <div className="grid grid-cols-2 gap-4 text-sm">
                {state.productData.sku && (
                  <div>
                    <span className="font-semibold">{t('sku')}:</span>
                    <span className="ml-2">{state.productData.sku}</span>
                  </div>
                )}
                {state.productData.barcode && (
                  <div>
                    <span className="font-semibold">{t('barcode')}:</span>
                    <span className="ml-2">{state.productData.barcode}</span>
                  </div>
                )}
                <div>
                  <span className="font-semibold">{t('statusLabel')}:</span>
                  <span
                    className={`ml-2 ${state.productData.is_active ? 'text-success' : 'text-error'}`}
                  >
                    {state.productData.is_active ? t('active') : t('inactive')}
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
                  {t('specifications')}
                </h4>

                <div className="space-y-3">
                  {Object.entries(state.attributes).map(([id, value]) => {
                    if (!value || (Array.isArray(value) && value.length === 0))
                      return null;

                    // Находим атрибут для получения переведенного названия
                    const attribute = categoryAttributes.find(
                      (attr) => attr.id === parseInt(id)
                    );
                    const { displayName, getOptionLabel } =
                      attribute && attribute.id
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

          {/* Варианты */}
          {state.hasVariants && state.variants.length > 0 && (
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h4 className="card-title text-lg mb-4 flex items-center gap-2">
                  <span className="text-xl">🎯</span>
                  {t('variants')}
                </h4>

                <div className="overflow-x-auto">
                  <table className="table table-sm">
                    <thead>
                      <tr>
                        <th>{t('variant')}</th>
                        <th>{t('sku')}</th>
                        <th>{t('price')}</th>
                        <th>{t('stock')}</th>
                        <th>{t('default')}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {state.variants.map((variant, index) => (
                        <tr key={index}>
                          <td>
                            {Object.entries(variant.variant_attributes)
                              .map(([k, v]) => `${k}: ${v}`)
                              .join(', ')}
                          </td>
                          <td>{variant.sku || '-'}</td>
                          <td>
                            {variant.price
                              ? formatPrice(
                                  variant.price,
                                  state.productData.currency || 'RSD'
                                )
                              : t('basePrice')}
                          </td>
                          <td>{variant.stock_quantity}</td>
                          <td>
                            {variant.is_default && (
                              <span className="badge badge-primary badge-sm">
                                {t('default')}
                              </span>
                            )}
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
                <div className="text-sm text-base-content/60 mt-2">
                  {t('totalVariants')}: {state.variants.length}
                </div>
              </div>
            </div>
          )}

          {/* Местоположение */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h4 className="card-title text-lg mb-4 flex items-center gap-2">
                <span className="text-xl">📍</span>
                {t('location')}
              </h4>

              <div className="space-y-3">
                <div className="flex justify-between items-center py-2 border-b border-base-200">
                  <span className="font-medium text-base-content/70">
                    {t('locationType')}:
                  </span>
                  <span className="text-base-content">
                    {state.location?.useStorefrontLocation
                      ? t('storefrontLocation')
                      : t('individualLocation')}
                  </span>
                </div>

                {!state.location?.useStorefrontLocation &&
                  state.location?.individualAddress && (
                    <>
                      <div className="flex justify-between items-center py-2 border-b border-base-200">
                        <span className="font-medium text-base-content/70">
                          {t('address')}:
                        </span>
                        <span className="text-base-content text-right max-w-[60%]">
                          {formatAddressWithPrivacy(
                            state.location.individualAddress || '',
                            state.location.privacyLevel as LocationPrivacyLevel
                          )}
                        </span>
                      </div>

                      <div className="flex justify-between items-center py-2 border-b border-base-200">
                        <span className="font-medium text-base-content/70">
                          {t('privacyLevel')}:
                        </span>
                        <span className="text-base-content">
                          {state.location.privacyLevel === 'exact' &&
                            t('products.privacy.exact')}
                          {state.location.privacyLevel === 'street' &&
                            t('products.privacy.street')}
                          {state.location.privacyLevel === 'district' &&
                            t('products.privacy.district')}
                          {state.location.privacyLevel === 'city' &&
                            t('products.privacy.city')}
                        </span>
                      </div>

                      <div className="flex justify-between items-center py-2">
                        <span className="font-medium text-base-content/70">
                          {t('showOnMap')}:
                        </span>
                        <span
                          className={`text-base-content ${state.location.showOnMap ? 'text-success' : 'text-error'}`}
                        >
                          {state.location.showOnMap
                            ? tCommon('yes')
                            : tCommon('no')}
                        </span>
                      </div>
                    </>
                  )}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Подтверждение */}
      <div className="card bg-gradient-to-r from-primary/10 to-secondary/10 border border-primary/20 shadow-xl mb-8">
        <div className="card-body">
          <div className="flex items-center gap-4">
            <div className="text-5xl">🎉</div>
            <div>
              <h3 className="text-xl font-bold text-primary mb-2">
                {t('readyToPublish')}
              </h3>
              <p className="text-base-content/70">{t('publishConfirmation')}</p>
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
          {tCommon('back')}
        </button>

        <div className="flex gap-4">
          <button className="btn btn-outline btn-lg px-8" disabled={submitting}>
            {t('saveDraft')}
          </button>

          <button
            onClick={handleSubmit}
            disabled={submitting}
            className="btn btn-primary btn-lg px-8"
          >
            {submitting ? (
              <>
                <span className="loading loading-spinner loading-sm mr-2"></span>
                {t('creating')}
              </>
            ) : (
              <>
                {t('publishProduct')}
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
