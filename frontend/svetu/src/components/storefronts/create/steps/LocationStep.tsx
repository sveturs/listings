'use client';

import { useState, useCallback, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import dynamic from 'next/dynamic';
import { useCreateStorefrontContext } from '@/contexts/CreateStorefrontContext';
import { LatLngExpression } from 'leaflet';

// Dynamic import for SSR compatibility
const LocationMap = dynamic(
  () => import('@/components/storefronts/create/LocationMap'),
  { ssr: false }
);

interface LocationStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function LocationStep({ onNext, onBack }: LocationStepProps) {
  const t = useTranslations();
  const { formData, updateFormData } = useCreateStorefrontContext();
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [searchingLocation, setSearchingLocation] = useState(false);

  const validate = () => {
    const newErrors: Record<string, string> = {};

    if (!formData.address || formData.address.length < 5) {
      newErrors.address = t('create_storefront.errors.address_required');
    }

    if (!formData.city || formData.city.length < 2) {
      newErrors.city = t('create_storefront.errors.city_required');
    }

    if (!formData.postalCode || formData.postalCode.length < 4) {
      newErrors.postalCode = t('create_storefront.errors.postal_code_required');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleNext = () => {
    if (validate()) {
      onNext();
    }
  };

  const handleLocationSelect = useCallback(
    (lat: number, lng: number) => {
      updateFormData({ latitude: lat, longitude: lng });
    },
    [updateFormData]
  );

  const searchCoordinates = async () => {
    if (!formData.address || !formData.city) {
      return;
    }

    setSearchingLocation(true);
    try {
      const query = `${formData.address}, ${formData.city}, ${formData.country}`;
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(query)}&limit=1`
      );
      const data = await response.json();

      if (data && data.length > 0) {
        const lat = parseFloat(data[0].lat);
        const lon = parseFloat(data[0].lon);
        updateFormData({ latitude: lat, longitude: lon });
      }
    } catch (error) {
      console.error('Error searching coordinates:', error);
    } finally {
      setSearchingLocation(false);
    }
  };

  useEffect(() => {
    if (
      formData.address &&
      formData.city &&
      !formData.latitude &&
      !formData.longitude
    ) {
      searchCoordinates();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [formData.address, formData.city]);

  const position: LatLngExpression =
    formData.latitude && formData.longitude
      ? [formData.latitude, formData.longitude]
      : [44.8125, 20.4612]; // Belgrade center as default

  return (
    <div className="max-w-4xl mx-auto">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">
            {t('create_storefront.location.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('create_storefront.location.subtitle')}
          </p>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    {t('create_storefront.location.address')}
                  </span>
                  <span className="label-text-alt text-error">*</span>
                </label>
                <input
                  type="text"
                  placeholder={t(
                    'create_storefront.location.address_placeholder'
                  )}
                  className={`input input-bordered w-full ${errors.address ? 'input-error' : ''}`}
                  value={formData.address}
                  onChange={(e) => updateFormData({ address: e.target.value })}
                />
                {errors.address && (
                  <label className="label">
                    <span className="label-text-alt text-error">
                      {errors.address}
                    </span>
                  </label>
                )}
              </div>

              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    {t('create_storefront.location.city')}
                  </span>
                  <span className="label-text-alt text-error">*</span>
                </label>
                <input
                  type="text"
                  placeholder={t('create_storefront.location.city_placeholder')}
                  className={`input input-bordered w-full ${errors.city ? 'input-error' : ''}`}
                  value={formData.city}
                  onChange={(e) => updateFormData({ city: e.target.value })}
                />
                {errors.city && (
                  <label className="label">
                    <span className="label-text-alt text-error">
                      {errors.city}
                    </span>
                  </label>
                )}
              </div>

              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    {t('create_storefront.location.postal_code')}
                  </span>
                  <span className="label-text-alt text-error">*</span>
                </label>
                <input
                  type="text"
                  placeholder={t(
                    'create_storefront.location.postal_code_placeholder'
                  )}
                  className={`input input-bordered w-full ${errors.postalCode ? 'input-error' : ''}`}
                  value={formData.postalCode}
                  onChange={(e) =>
                    updateFormData({ postalCode: e.target.value })
                  }
                />
                {errors.postalCode && (
                  <label className="label">
                    <span className="label-text-alt text-error">
                      {errors.postalCode}
                    </span>
                  </label>
                )}
              </div>

              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    {t('create_storefront.location.country')}
                  </span>
                </label>
                <select
                  className="select select-bordered w-full"
                  value={formData.country}
                  onChange={(e) => updateFormData({ country: e.target.value })}
                >
                  <option value="RS">{t('countries.RS')}</option>
                  <option value="BA">{t('countries.BA')}</option>
                  <option value="HR">{t('countries.HR')}</option>
                  <option value="ME">{t('countries.ME')}</option>
                  <option value="MK">{t('countries.MK')}</option>
                </select>
              </div>

              <button
                className={`btn btn-outline btn-sm w-full ${searchingLocation ? 'loading' : ''}`}
                onClick={searchCoordinates}
                disabled={
                  searchingLocation || !formData.address || !formData.city
                }
              >
                {searchingLocation
                  ? t('common.searching')
                  : t('create_storefront.location.find_on_map')}
              </button>
            </div>

            <div className="h-96 rounded-lg overflow-hidden">
              <LocationMap
                position={position}
                onLocationSelect={handleLocationSelect}
              />
            </div>
          </div>

          <div className="card-actions justify-between mt-6">
            <button className="btn btn-ghost" onClick={onBack}>
              {t('common.back')}
            </button>
            <button className="btn btn-primary" onClick={handleNext}>
              {t('common.next')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
