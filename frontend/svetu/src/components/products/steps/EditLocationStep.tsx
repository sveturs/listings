'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useEditProduct } from '@/contexts/EditProductContext';
import LocationPicker from '@/components/GIS/LocationPicker';
import LocationPrivacySettings from '@/components/GIS/LocationPrivacySettings';

interface EditLocationStepProps {
  onNext: () => void;
  onBack: () => void;
}

interface LocationData {
  latitude: number;
  longitude: number;
  address: string;
  city: string;
  region: string;
  country: string;
  confidence: number;
}

export default function EditLocationStep({
  onNext,
  onBack,
}: EditLocationStepProps) {
  const t = useTranslations('storefronts.products');
  const tCommon = useTranslations('common');
  const { state, setLocation, setError, clearError } = useEditProduct();

  // Инициализация состояния из контекста
  const [useStorefrontLocation, setUseStorefrontLocation] = useState(
    state.location?.useStorefrontLocation ?? true
  );
  const [individualLocation, setIndividualLocation] = useState<
    LocationData | undefined
  >(
    state.location?.latitude &&
      state.location?.longitude &&
      !state.location?.useStorefrontLocation
      ? {
          latitude: state.location.latitude,
          longitude: state.location.longitude,
          address: state.location.individualAddress || '',
          city: state.location.city || '',
          region: state.location.region || '',
          country: state.location.country || 'Србија',
          confidence: 0.9,
        }
      : undefined
  );
  const [privacyLevel, setPrivacyLevel] = useState<
    'exact' | 'street' | 'district' | 'city'
  >(
    state.location?.privacyLevel === 'street' ||
      state.location?.privacyLevel === 'district' ||
      state.location?.privacyLevel === 'city' ||
      state.location?.privacyLevel === 'exact'
      ? state.location.privacyLevel
      : 'exact'
  );

  // Сохранение в контексте при изменениях
  useEffect(() => {
    setLocation({
      useStorefrontLocation,
      individualAddress: individualLocation?.address || '',
      latitude: individualLocation?.latitude,
      longitude: individualLocation?.longitude,
      city: individualLocation?.city || '',
      region: individualLocation?.region || '',
      country: individualLocation?.country || '',
      privacyLevel,
      showOnMap: true,
    });
  }, [useStorefrontLocation, individualLocation, privacyLevel, setLocation]);

  const handleLocationTypeChange = (useStorefront: boolean) => {
    setUseStorefrontLocation(useStorefront);
    clearError('location');

    if (useStorefront) {
      setIndividualLocation(undefined);
    }
  };

  const handleLocationSelect = (location: LocationData) => {
    setIndividualLocation(location);
    clearError('location');
  };

  const handleNext = () => {
    // Валидация
    if (!useStorefrontLocation && !individualLocation) {
      setError('location', t('locationRequired'));
      return;
    }

    clearError('location');
    onNext();
  };

  return (
    <div className="w-full">
      <div className="text-center mb-6 lg:mb-8">
        <h2 className="text-xl lg:text-3xl font-bold text-base-content mb-2 lg:mb-4">
          {t('productLocation')}
        </h2>
        <p className="text-sm lg:text-lg text-base-content/70">
          {t('locationDescription')}
        </p>
      </div>

      {/* Тип местоположения */}
      <div className="mb-6 lg:mb-8">
        <h3 className="text-lg font-bold text-base-content mb-4">
          {t('locationType')}
        </h3>

        <div className="space-y-4">
          <div
            className={`
              p-4 lg:p-6 rounded-lg lg:rounded-xl border-2 cursor-pointer 
              transition-all duration-200 hover:shadow-lg
              ${
                useStorefrontLocation
                  ? 'border-primary bg-primary/10 shadow-lg'
                  : 'border-base-300 bg-base-100 hover:border-primary/50'
              }
            `}
            onClick={() => handleLocationTypeChange(true)}
          >
            <div className="flex items-start gap-4">
              <input
                type="radio"
                checked={useStorefrontLocation}
                onChange={() => handleLocationTypeChange(true)}
                className="radio radio-primary flex-shrink-0 mt-1"
              />
              <div className="flex-1">
                <h4 className="text-lg font-semibold text-base-content mb-2">
                  {t('useStorefrontLocation')}
                </h4>
                <p className="text-base-content/70 mb-3">
                  {t('useStorefrontLocationDescription')}
                </p>
                <div className="bg-info/10 border border-info/30 rounded-lg p-3">
                  <div className="flex items-center gap-2 text-sm text-info">
                    <span>💡</span>
                    <span className="font-medium">
                      {t('example')}
                    </span>
                  </div>
                  <p className="text-sm text-info/80 mt-1">
                    {t('storefrontLocationExample')}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div
            className={`
              p-4 lg:p-6 rounded-lg lg:rounded-xl border-2 cursor-pointer 
              transition-all duration-200 hover:shadow-lg
              ${
                !useStorefrontLocation
                  ? 'border-primary bg-primary/10 shadow-lg'
                  : 'border-base-300 bg-base-100 hover:border-primary/50'
              }
            `}
            onClick={() => handleLocationTypeChange(false)}
          >
            <div className="flex items-start gap-4">
              <input
                type="radio"
                checked={!useStorefrontLocation}
                onChange={() => handleLocationTypeChange(false)}
                className="radio radio-primary flex-shrink-0 mt-1"
              />
              <div className="flex-1">
                <h4 className="text-lg font-semibold text-base-content mb-2">
                  {t('useIndividualLocation')}
                </h4>
                <p className="text-base-content/70 mb-3">
                  {t('useIndividualLocationDescription')}
                </p>
                <div className="bg-warning/10 border border-warning/30 rounded-lg p-3">
                  <div className="flex items-center gap-2 text-sm text-warning">
                    <span>💡</span>
                    <span className="font-medium">
                      {t('example')}
                    </span>
                  </div>
                  <p className="text-sm text-warning/80 mt-1">
                    {t('individualLocationExample')}
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Выбор индивидуального адреса */}
      {!useStorefrontLocation && (
        <div className="mb-6 lg:mb-8">
          <h3 className="text-lg font-bold text-base-content mb-4">
            {t('selectProductLocation')}
          </h3>

          <LocationPicker
            onChange={handleLocationSelect}
            value={individualLocation}
            placeholder={t('locationPlaceholder')}
          />

          {individualLocation && (
            <div className="mt-4 p-4 bg-success/10 border border-success/30 rounded-lg">
              <div className="flex items-start gap-3">
                <span className="text-success text-lg">📍</span>
                <div>
                  <h4 className="font-semibold text-success">
                    {t('address')}
                  </h4>
                  <p className="text-success/80">
                    {individualLocation.address}
                  </p>
                  <p className="text-xs text-success/60 mt-1">
                    {individualLocation.latitude.toFixed(6)},{' '}
                    {individualLocation.longitude.toFixed(6)}
                  </p>
                </div>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Настройки приватности */}
      {!useStorefrontLocation && individualLocation && (
        <div className="mb-6 lg:mb-8">
          <h3 className="text-lg font-bold text-base-content mb-4">
            {t('privacySettings')}
          </h3>

          <div className="bg-base-200/50 rounded-xl p-4 lg:p-6">
            <div className="flex items-center gap-3 mb-4">
              <span className="text-2xl">🔒</span>
              <div>
                <h4 className="font-semibold text-base-content">
                  {t('privacyNote')}
                </h4>
                <p className="text-sm text-base-content/70">
                  {t('privacyNoteDescription')}
                </p>
              </div>
            </div>

            <LocationPrivacySettings
              selectedLevel={privacyLevel}
              onLevelChange={setPrivacyLevel}
              location={
                individualLocation
                  ? {
                      lat: individualLocation.latitude,
                      lng: individualLocation.longitude,
                    }
                  : undefined
              }
              showPreview={true}
            />
          </div>
        </div>
      )}

      {/* Ошибки */}
      {state.errors.location && (
        <div className="mb-6 alert alert-error">
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
          <span>{state.errors.location}</span>
        </div>
      )}

      {/* Кнопки навигации */}
      <div className="flex justify-between">
        <button onClick={onBack} className="btn btn-outline btn-lg">
          {tCommon('back')}
        </button>
        <button onClick={handleNext} className="btn btn-primary btn-lg">
          {tCommon('continue')}
        </button>
      </div>
    </div>
  );
}
