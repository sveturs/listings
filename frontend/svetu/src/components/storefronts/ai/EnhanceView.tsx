'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations, useLocale, NextIntlClientProvider } from 'next-intl';
import { useCreateAIProduct } from '@/contexts/CreateAIProductContext';
import { CategoryTreeSelector } from '@/components/common/CategoryTreeSelector';
import LocationPicker from '@/components/GIS/LocationPicker';
import LocationPrivacySettings from '@/components/GIS/LocationPrivacySettings';
import { useAddressGeocoding } from '@/hooks/useAddressGeocoding';
import { MapPin, ChevronDown } from 'lucide-react';

interface EnhanceViewProps {
  storefrontId: number | null;
  storefrontSlug: string;
}

export default function EnhanceView({
  storefrontId: _storefrontId,
  storefrontSlug: _storefrontSlug,
}: EnhanceViewProps) {
  const t = useTranslations('storefronts');
  const locale = useLocale();
  const {
    state,
    setAIData,
    setView,
    setUseStorefrontLocation,
    setLocationPrivacyLevel,
    setShowOnMap,
  } = useCreateAIProduct();

  // Геокодирование для получения локализованного адреса
  const geocoding = useAddressGeocoding({
    country: ['rs'],
    language: locale,
  });

  // Получаем контент на текущей локали
  const getLocalizedContent = () => {
    const translations = state.aiData.translations || {};
    // Если есть перевод на текущую локаль, используем его
    if (translations[locale]) {
      return {
        title: translations[locale].title,
        description: translations[locale].description,
      };
    }
    // Иначе используем оригинал (английский)
    return {
      title: state.aiData.title,
      description: state.aiData.description,
    };
  };

  const localizedContent = getLocalizedContent();

  const [editedData, setEditedData] = useState({
    title: localizedContent.title,
    description: localizedContent.description,
    price: state.aiData.price,
    stockQuantity: state.aiData.stockQuantity,
    categoryId: state.aiData.categoryId,
    category: state.aiData.category,
  });

  // Location states
  const hasExifLocation = state.aiData.location?.source === 'exif';
  // ВАЖНО: Если есть EXIF, ВСЕГДА используем его по умолчанию (игнорируем state.useStorefrontLocation)
  const [localUseStorefrontLocation, setLocalUseStorefrontLocation] = useState(
    hasExifLocation ? false : (state.useStorefrontLocation ?? true)
  );
  const [useExifLocation, setUseExifLocation] = useState(hasExifLocation);
  const [individualLocation, setIndividualLocation] = useState<
    | {
        latitude: number;
        longitude: number;
        address: string;
        city: string;
        region: string;
        country: string;
        confidence: number;
      }
    | undefined
  >(
    state.aiData.location && state.aiData.location.source !== 'storefront'
      ? {
          latitude: state.aiData.location.latitude,
          longitude: state.aiData.location.longitude,
          address: state.aiData.location.address,
          city: state.aiData.location.city,
          region: state.aiData.location.region,
          country: 'Србија',
          confidence: 0.9,
        }
      : undefined
  );
  const [localPrivacyLevel, setLocalPrivacyLevel] = useState<
    'exact' | 'street' | 'district' | 'city'
  >(state.locationPrivacyLevel || 'exact');
  const [localShowOnMap, setLocalShowOnMap] = useState(state.showOnMap ?? true);
  const [showPrivacySettings, setShowPrivacySettings] = useState(false);

  // Состояние для локализованного EXIF адреса
  const [localizedExifAddress, setLocalizedExifAddress] = useState<{
    address: string;
    city: string;
    region: string;
  } | null>(null);

  // Эффект для загрузки локализованного адреса при изменении локали или EXIF данных
  useEffect(() => {
    const loadLocalizedAddress = async () => {
      if (
        state.aiData.location?.source === 'exif' &&
        state.aiData.location?.latitude &&
        state.aiData.location?.longitude
      ) {
        try {
          const geocodedAddress = await geocoding.reverseGeocode(
            state.aiData.location.latitude,
            state.aiData.location.longitude
          );

          if (geocodedAddress) {
            setLocalizedExifAddress({
              address:
                geocodedAddress.address_components?.formatted ||
                geocodedAddress.place_name ||
                state.aiData.location.address,
              city:
                geocodedAddress.address_components?.city ||
                state.aiData.location.city,
              region:
                geocodedAddress.address_components?.district ||
                geocodedAddress.region ||
                state.aiData.location.region,
            });
          } else {
            // Fallback к оригинальным данным
            setLocalizedExifAddress({
              address: state.aiData.location.address,
              city: state.aiData.location.city,
              region: state.aiData.location.region,
            });
          }
        } catch (error) {
          console.error('Failed to geocode EXIF location:', error);
          // Fallback к оригинальным данным
          setLocalizedExifAddress({
            address: state.aiData.location.address,
            city: state.aiData.location.city,
            region: state.aiData.location.region,
          });
        }
      }
    };

    loadLocalizedAddress();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [state.aiData.location, locale]);

  const handleLocationTypeChange = (type: 'storefront' | 'exif' | 'manual') => {
    if (type === 'storefront') {
      setLocalUseStorefrontLocation(true);
      setUseStorefrontLocation(true);
      setUseExifLocation(false);
    } else if (type === 'exif') {
      setLocalUseStorefrontLocation(false);
      setUseStorefrontLocation(false);
      setUseExifLocation(true);
    } else {
      setLocalUseStorefrontLocation(false);
      setUseStorefrontLocation(false);
      setUseExifLocation(false);
    }
  };

  const handleLocationChange = (locationData: any) => {
    setIndividualLocation({
      latitude: locationData.latitude,
      longitude: locationData.longitude,
      address: locationData.address,
      city: locationData.city,
      region: locationData.region,
      country: locationData.country || 'Србија',
      confidence: locationData.confidence || 0.9,
    });
  };

  const handleSave = () => {
    console.log(
      '[EnhanceView] handleSave - state.aiData.location:',
      state.aiData.location
    );
    console.log(
      '[EnhanceView] handleSave - localUseStorefrontLocation:',
      localUseStorefrontLocation
    );
    console.log('[EnhanceView] handleSave - useExifLocation:', useExifLocation);

    // Определяем location данные
    let locationData = null;
    if (localUseStorefrontLocation) {
      // Используем адрес витрины - location = null (будет использоваться storefront.address)
      locationData = null;
      console.log(
        '[EnhanceView] Using storefront location, setting location to null'
      );
    } else if (useExifLocation && state.aiData.location) {
      // Используем EXIF location
      locationData = {
        ...state.aiData.location,
        source: 'exif' as const,
      };
      console.log('[EnhanceView] Using EXIF location:', locationData);
    } else if (individualLocation) {
      // Используем manually entered location
      locationData = {
        ...individualLocation,
        source: 'manual' as const,
      };
      console.log('[EnhanceView] Using manual location:', locationData);
    } else if (state.aiData.location) {
      // Сохраняем существующий location
      locationData = state.aiData.location;
      console.log('[EnhanceView] Preserving existing location:', locationData);
    }

    console.log('[EnhanceView] Final locationData to save:', locationData);

    // Сохраняем основные данные (сохраняем translations и другие существующие поля)
    setAIData({
      ...state.aiData, // Сохраняем все существующие поля
      title: editedData.title,
      description: editedData.description,
      price: editedData.price,
      stockQuantity: editedData.stockQuantity,
      categoryId: editedData.categoryId,
      category: editedData.category,
      location: locationData,
    });

    // Сохраняем настройки локации
    setLocationPrivacyLevel(localPrivacyLevel);
    setShowOnMap(localShowOnMap);
    setUseStorefrontLocation(localUseStorefrontLocation);

    setView('variants');
  };

  const handleCategoryChange = async (categoryId: number | number[]) => {
    const id = Array.isArray(categoryId) ? categoryId[0] : categoryId;
    // Загружаем информацию о категории для получения имени
    try {
      const response = await fetch(
        `/api/v1/marketplace/categories?page=1&limit=1000`
      );
      if (response.ok) {
        const data = await response.json();
        const category = data.data?.find((cat: any) => cat.id === id);
        if (category) {
          setEditedData((prev) => ({
            ...prev,
            categoryId: id,
            category: category.name,
          }));
        }
      }
    } catch (error) {
      console.error('Failed to load category name:', error);
      setEditedData((prev) => ({ ...prev, categoryId: id }));
    }
  };

  return (
    <div className="space-y-6">
      <div>
        <h2 className="text-2xl font-bold mb-2">
          {t('enhanceProduct') || 'Enhance Your Product'}
        </h2>
        <p className="text-base-content/70">
          {t('enhanceDescription') || 'Review and edit AI-generated content'}
        </p>
      </div>

      {/* Title */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-semibold">
            {t('productName') || 'Title'}
          </span>
        </label>
        <input
          type="text"
          value={editedData.title}
          onChange={(e) =>
            setEditedData((prev) => ({ ...prev, title: e.target.value }))
          }
          className="input input-bordered w-full"
          placeholder={t('productName') || 'Product title'}
        />
      </div>

      {/* Description */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-semibold">
            {t('description') || 'Description'}
          </span>
        </label>
        <textarea
          value={editedData.description}
          onChange={(e) =>
            setEditedData((prev) => ({ ...prev, description: e.target.value }))
          }
          className="textarea textarea-bordered h-32"
          placeholder={t('productDescription') || 'Product description'}
        />
      </div>

      {/* Price & Stock */}
      <div className="grid grid-cols-2 gap-4">
        <div className="form-control">
          <label className="label">
            <span className="label-text font-semibold">
              {t('price') || 'Price'}
            </span>
          </label>
          <input
            type="number"
            value={editedData.price}
            onChange={(e) =>
              setEditedData((prev) => ({
                ...prev,
                price: Number(e.target.value),
              }))
            }
            className="input input-bordered"
            min="0"
          />
        </div>
        <div className="form-control">
          <label className="label">
            <span className="label-text font-semibold">
              {t('stockQuantity') || 'Stock'}
            </span>
          </label>
          <input
            type="number"
            value={editedData.stockQuantity}
            onChange={(e) =>
              setEditedData((prev) => ({
                ...prev,
                stockQuantity: Number(e.target.value),
              }))
            }
            className="input input-bordered"
            min="0"
          />
        </div>
      </div>

      {/* Category Selection */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-semibold">
            {t('category') || 'Category'}
          </span>
          <span className="label-text-alt text-info">
            {t('aiDetected') || 'AI Detected'}: {state.aiData.category}
          </span>
        </label>
        <NextIntlClientProvider
          locale={locale}
          messages={{
            marketplace: {
              selectCategory: t('selectCategory'),
              searchCategories: t('searchCategories'),
              categoriesSelected: t.raw('categoriesSelected'),
              apply: t('apply'),
              cancel: t('cancel'),
              categoriesLoadError: t('categoriesLoadError'),
              noCategoriesFound: t('noCategoriesFound'),
            },
          }}
        >
          <CategoryTreeSelector
            value={editedData.categoryId}
            onChange={handleCategoryChange}
            placeholder={t('selectCategory') || 'Select category'}
            showPath={true}
            allowParentSelection={false}
          />
        </NextIntlClientProvider>
      </div>

      {/* Location Selection */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-semibold flex items-center gap-2">
            <MapPin className="w-4 h-4" />
            {t('productLocation') || 'Location'}
          </span>
        </label>

        <div className="space-y-3">
          {/* Использовать адрес витрины */}
          <label className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer">
            <div className="card-body p-4">
              <div className="flex items-start gap-3">
                <input
                  type="radio"
                  name="locationType"
                  className="radio radio-primary mt-1"
                  checked={localUseStorefrontLocation}
                  onChange={() => handleLocationTypeChange('storefront')}
                />
                <div className="flex-1">
                  <h4 className="font-medium text-base-content">
                    {t('useStorefrontLocation') || 'Использовать адрес витрины'}
                  </h4>
                  <p className="text-sm text-base-content/70 mt-1">
                    {t('products.useStorefrontLocationDescription') ||
                      'Адрес будет взят из настроек вашей витрины'}
                  </p>
                </div>
              </div>
            </div>
          </label>

          {/* Использовать адрес из EXIF */}
          {state.aiData.location?.source === 'exif' && (
            <label className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer">
              <div className="card-body p-4">
                <div className="flex items-start gap-3">
                  <input
                    type="radio"
                    name="locationType"
                    className="radio radio-primary mt-1"
                    checked={useExifLocation}
                    onChange={() => handleLocationTypeChange('exif')}
                  />
                  <div className="flex-1">
                    <h4 className="font-medium text-base-content">
                      {t('useExifLocation') ||
                        'Использовать адрес из метаданных фото'}
                    </h4>
                    <p className="text-sm text-base-content/70 mt-1">
                      {localizedExifAddress?.address ||
                        state.aiData.location?.address ||
                        'Адрес из фото'}
                    </p>
                    <div className="mt-2 text-xs text-info flex items-center gap-1">
                      <MapPin className="w-3 h-3" />
                      <span>
                        {localizedExifAddress?.city ||
                          state.aiData.location?.city}
                        {localizedExifAddress?.region ||
                        state.aiData.location?.region
                          ? `, ${localizedExifAddress?.region || state.aiData.location?.region}`
                          : ''}
                      </span>
                      <span className="text-base-content/50 ml-2">
                        ({state.aiData.location?.latitude.toFixed(6)},{' '}
                        {state.aiData.location?.longitude.toFixed(6)})
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </label>
          )}

          {/* Указать индивидуальный адрес */}
          <label className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer">
            <div className="card-body p-4">
              <div className="flex items-start gap-3">
                <input
                  type="radio"
                  name="locationType"
                  className="radio radio-primary mt-1"
                  checked={!localUseStorefrontLocation && !useExifLocation}
                  onChange={() => handleLocationTypeChange('manual')}
                />
                <div className="flex-1">
                  <h4 className="font-medium text-base-content">
                    {t('useIndividualLocation') ||
                      'Указать индивидуальный адрес'}
                  </h4>
                  <p className="text-sm text-base-content/70 mt-1">
                    {t('products.useIndividualLocationDescription') ||
                      'Выбрать адрес на карте вручную'}
                  </p>
                </div>
              </div>
            </div>
          </label>
        </div>

        {/* Location Picker для индивидуального адреса */}
        {!localUseStorefrontLocation && !useExifLocation && (
          <div className="mt-4 p-4 bg-base-200 rounded-lg">
            <LocationPicker
              value={individualLocation}
              onChange={handleLocationChange}
              placeholder={
                t('locationPlaceholder') || 'Выберите адрес на карте'
              }
              height="300px"
              showCurrentLocation={true}
              defaultCountry="Србија"
            />

            {/* Privacy Settings */}
            {individualLocation && (
              <div className="mt-4">
                <button
                  type="button"
                  onClick={() => setShowPrivacySettings(!showPrivacySettings)}
                  className="btn btn-outline btn-sm w-full"
                >
                  🛡️ {t('privacySettings') || 'Настройки приватности'}
                  <ChevronDown
                    className={`w-4 h-4 ml-2 transition-transform ${showPrivacySettings ? 'rotate-180' : ''}`}
                  />
                </button>

                {showPrivacySettings && (
                  <div className="mt-4 p-4 bg-base-100 rounded-lg">
                    <LocationPrivacySettings
                      selectedLevel={localPrivacyLevel}
                      onLevelChange={setLocalPrivacyLevel}
                      location={{
                        lat: individualLocation.latitude,
                        lng: individualLocation.longitude,
                      }}
                      showPreview={true}
                    />

                    <div className="form-control mt-4">
                      <label className="label cursor-pointer">
                        <span className="label-text font-medium">
                          {t('showOnMap') || 'Показывать на карте'}
                        </span>
                        <input
                          type="checkbox"
                          checked={localShowOnMap}
                          onChange={(e) => setLocalShowOnMap(e.target.checked)}
                          className="checkbox checkbox-primary"
                        />
                      </label>
                      <p className="text-sm text-base-content/60 mt-1">
                        {t('showOnMapDescription') ||
                          'Разрешить показывать товар на карте в результатах поиска'}
                      </p>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        )}

        {/* Privacy Settings для EXIF адреса */}
        {useExifLocation && state.aiData.location && (
          <div className="mt-4">
            <button
              type="button"
              onClick={() => setShowPrivacySettings(!showPrivacySettings)}
              className="btn btn-outline btn-sm w-full"
            >
              🛡️ {t('privacySettings') || 'Настройки приватности'}
              <ChevronDown
                className={`w-4 h-4 ml-2 transition-transform ${showPrivacySettings ? 'rotate-180' : ''}`}
              />
            </button>

            {showPrivacySettings && (
              <div className="mt-4 p-4 bg-base-200 rounded-lg">
                <LocationPrivacySettings
                  selectedLevel={localPrivacyLevel}
                  onLevelChange={setLocalPrivacyLevel}
                  location={{
                    lat: state.aiData.location.latitude,
                    lng: state.aiData.location.longitude,
                  }}
                  showPreview={true}
                />

                <div className="form-control mt-4">
                  <label className="label cursor-pointer">
                    <span className="label-text font-medium">
                      {t('showOnMap') || 'Показывать на карте'}
                    </span>
                    <input
                      type="checkbox"
                      checked={localShowOnMap}
                      onChange={(e) => setLocalShowOnMap(e.target.checked)}
                      className="checkbox checkbox-primary"
                    />
                  </label>
                  <p className="text-sm text-base-content/60 mt-1">
                    {t('showOnMapDescription') ||
                      'Разрешить показывать товар на карте в результатах поиска'}
                  </p>
                </div>
              </div>
            )}
          </div>
        )}
      </div>

      {/* Actions */}
      <div className="flex justify-between gap-3">
        <button onClick={() => setView('process')} className="btn btn-outline">
          {t('back') || 'Back'}
        </button>
        <button onClick={handleSave} className="btn btn-primary px-8">
          {t('continueToVariants') || 'Continue to Variants'}
        </button>
      </div>
    </div>
  );
}
