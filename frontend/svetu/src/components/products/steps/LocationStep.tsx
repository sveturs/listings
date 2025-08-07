'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateProduct } from '@/contexts/CreateProductContext';
import LocationPicker from '@/components/GIS/LocationPicker';
import LocationPrivacySettings from '@/components/GIS/LocationPrivacySettings';

interface LocationStepProps {
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

export default function LocationStep({ onNext, onBack }: LocationStepProps) {
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const { state, setLocation, setError, clearError } = useCreateProduct();

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
  >(state.location?.privacyLevel || 'exact');
  const [showOnMap, setShowOnMap] = useState(state.location?.showOnMap ?? true);
  const [showPrivacySettings, setShowPrivacySettings] = useState(false);

  // Сохраняем изменения в контекст
  useEffect(() => {
    const locationData = {
      useStorefrontLocation,
      ...(individualLocation && !useStorefrontLocation
        ? {
            individualAddress: individualLocation.address,
            latitude: individualLocation.latitude,
            longitude: individualLocation.longitude,
            city: individualLocation.city,
            region: individualLocation.region,
            country: individualLocation.country,
          }
        : {}),
      privacyLevel,
      showOnMap,
    };

    setLocation(locationData);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [useStorefrontLocation, individualLocation, privacyLevel, showOnMap]);

  const handleLocationTypeChange = (useStorefront: boolean) => {
    setUseStorefrontLocation(useStorefront);
    if (useStorefront) {
      clearError('location');
    }
  };

  const handleLocationChange = (locationData: LocationData) => {
    setIndividualLocation(locationData);
    clearError('location');
  };

  const validateAndProceed = () => {
    if (!useStorefrontLocation && !individualLocation) {
      setError('location', t('locationRequired'));
      return;
    }
    onNext();
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="text-center mb-8">
        <h2 className="text-3xl font-bold text-base-content mb-4">
          📍 {t('productLocation')}
        </h2>
        <p className="text-lg text-base-content/70">
          {t('locationDescription')}
        </p>
      </div>

      {/* Выбор типа адреса */}
      <div className="card bg-base-100 shadow-xl mb-6">
        <div className="card-body">
          <h3 className="text-xl font-semibold mb-4 flex items-center gap-2">
            <span className="text-2xl">🏪</span>
            {t('locationType')}
          </h3>

          <div className="space-y-4">
            {/* Использовать адрес витрины */}
            <label className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer">
              <div className="card-body p-4">
                <div className="flex items-start gap-4">
                  <input
                    type="radio"
                    name="locationType"
                    className="radio radio-primary mt-1"
                    checked={useStorefrontLocation}
                    onChange={() => handleLocationTypeChange(true)}
                  />
                  <div className="flex-1">
                    <h4 className="font-medium text-base-content">
                      {t('useStorefrontLocation')}
                    </h4>
                    <p className="text-sm text-base-content/70 mt-1">
                      {t('products.useStorefrontLocationDescription')}
                    </p>
                    <div className="mt-3 p-3 bg-info/10 rounded-lg">
                      <p className="text-sm text-info-content flex items-start gap-2">
                        <span>💡</span>
                        <span>
                          <strong>{t('example')}:</strong>{' '}
                          {t('storefrontLocationExample')}
                        </span>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </label>

            {/* Указать индивидуальный адрес */}
            <label className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer">
              <div className="card-body p-4">
                <div className="flex items-start gap-4">
                  <input
                    type="radio"
                    name="locationType"
                    className="radio radio-primary mt-1"
                    checked={!useStorefrontLocation}
                    onChange={() => handleLocationTypeChange(false)}
                  />
                  <div className="flex-1">
                    <h4 className="font-medium text-base-content">
                      {t('useIndividualLocation')}
                    </h4>
                    <p className="text-sm text-base-content/70 mt-1">
                      {t('products.useIndividualLocationDescription')}
                    </p>
                    <div className="mt-3 p-3 bg-info/10 rounded-lg">
                      <p className="text-sm text-info-content flex items-start gap-2">
                        <span>💡</span>
                        <span>
                          <strong>{t('example')}:</strong>{' '}
                          {t('individualLocationExample')}
                        </span>
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </label>
          </div>
        </div>
      </div>

      {/* Выбор индивидуального адреса */}
      {!useStorefrontLocation && (
        <div className="card bg-base-100 shadow-xl mb-6">
          <div className="card-body">
            <h3 className="text-xl font-semibold mb-4 flex items-center gap-2">
              <span className="text-2xl">📍</span>
              {t('selectProductLocation')}
            </h3>

            <LocationPicker
              value={individualLocation}
              onChange={handleLocationChange}
              placeholder={t('locationPlaceholder')}
              height="400px"
              showCurrentLocation={true}
              defaultCountry={state.location?.country || 'Србија'}
            />

            {state.errors.location && (
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
                <span>{state.errors.location}</span>
              </div>
            )}

            {individualLocation && (
              <div className="mt-6">
                <button
                  type="button"
                  onClick={() => setShowPrivacySettings(!showPrivacySettings)}
                  className="btn btn-outline btn-sm"
                >
                  🛡️ {t('privacySettings')}
                  <svg
                    className={`w-4 h-4 ml-2 transition-transform ${showPrivacySettings ? 'rotate-180' : ''}`}
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M19 9l-7 7-7-7"
                    />
                  </svg>
                </button>

                {showPrivacySettings && (
                  <div className="mt-4 p-4 bg-base-200 rounded-lg">
                    <LocationPrivacySettings
                      selectedLevel={privacyLevel}
                      onLevelChange={setPrivacyLevel}
                      location={{
                        lat: individualLocation.latitude,
                        lng: individualLocation.longitude,
                      }}
                      showPreview={true}
                    />

                    <div className="form-control mt-4">
                      <label className="label cursor-pointer">
                        <span className="label-text font-medium">
                          {t('showOnMap')}
                        </span>
                        <input
                          type="checkbox"
                          checked={showOnMap}
                          onChange={(e) => setShowOnMap(e.target.checked)}
                          className="checkbox checkbox-primary"
                        />
                      </label>
                      <p className="text-sm text-base-content/60 mt-1">
                        {t('showOnMapDescription')}
                      </p>
                    </div>
                  </div>
                )}
              </div>
            )}
          </div>
        </div>
      )}

      {/* Информационный блок */}
      <div className="alert alert-info mb-6">
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
          />
        </svg>
        <div>
          <h4 className="font-semibold">{t('privacyNote')}</h4>
          <p className="text-sm mt-1">{t('privacyNoteDescription')}</p>
        </div>
      </div>

      {/* Навигация */}
      <div className="flex justify-between items-center mt-8">
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
          onClick={validateAndProceed}
          className="btn btn-primary btn-lg px-8"
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
