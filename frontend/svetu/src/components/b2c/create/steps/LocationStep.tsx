'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateStorefrontContext } from '@/contexts/CreateB2CStoreContext';
import LocationPicker from '@/components/GIS/LocationPicker';

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
  const t = useTranslations('create_storefront');
  const tCommon = useTranslations('common');
  const { formData, updateFormData } = useCreateStorefrontContext();
  const [errors, setErrors] = useState<Record<string, string>>({});
  const [location, setLocation] = useState<LocationData | undefined>(
    formData.latitude && formData.longitude
      ? {
          latitude: formData.latitude,
          longitude: formData.longitude,
          address: formData.address || '',
          city: formData.city || '',
          region: '',
          country: formData.country || 'RS',
          confidence: 0.9,
        }
      : undefined
  );
  const [postalCode, setPostalCode] = useState(formData.postalCode || '');
  const [additionalInfo, setAdditionalInfo] = useState({
    floor: '',
    suite: '',
    hasParking: false,
    hasElevator: false,
    accessibilityNotes: '',
  });

  // Обновляем данные формы при изменении location
  useEffect(() => {
    if (location) {
      updateFormData({
        latitude: location.latitude,
        longitude: location.longitude,
        address: location.address,
        city: location.city,
        country: location.country,
      });
    }
  }, [location, updateFormData]);

  const handleLocationChange = (locationData: LocationData) => {
    setLocation(locationData);
    setErrors({}); // Очищаем ошибки при выборе нового местоположения
  };

  const validate = () => {
    const newErrors: Record<string, string> = {};

    if (!location) {
      newErrors.location = 'Необходимо выбрать местоположение';
    }

    if (!location?.address || location.address.length < 5) {
      newErrors.address = t('address_required');
    }

    if (!location?.city || location.city.length < 2) {
      newErrors.city = t('city_required');
    }

    if (!postalCode || postalCode.length < 4) {
      newErrors.postalCode = t('postal_code_required');
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleNext = () => {
    if (validate()) {
      updateFormData({ postalCode });
      onNext();
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">{t('location.title')}</h2>
          <p className="text-base-content/70 mb-6">{t('location.subtitle')}</p>

          {/* Location selection */}
          <div className="mb-6">
            <h3 className="text-lg font-semibold mb-4">
              📍 {t('location.storefront_location')}
            </h3>
            <LocationPicker
              value={location}
              onChange={handleLocationChange}
              placeholder={t('location.address_placeholder')}
              height="400px"
              showCurrentLocation={false}
              defaultCountry="Србија"
            />
            {errors.location && (
              <p className="text-error text-sm mt-2">{errors.location}</p>
            )}
          </div>

          {/* Additional information */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              {/* Postal code */}
              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    {t('location.postal_code')}
                  </span>
                  <span className="label-text-alt text-error">*</span>
                </label>
                <input
                  type="text"
                  placeholder={t('location.postal_code_placeholder')}
                  className={`input input-bordered w-full ${errors.postalCode ? 'input-error' : ''}`}
                  value={postalCode}
                  onChange={(e) => setPostalCode(e.target.value)}
                />
                {errors.postalCode && (
                  <label className="label">
                    <span className="label-text-alt text-error">
                      {errors.postalCode}
                    </span>
                  </label>
                )}
              </div>

              {/* Floor */}
              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">{t('location.floor')}</span>
                </label>
                <input
                  type="text"
                  placeholder={t('location.floor_placeholder')}
                  className="input input-bordered w-full"
                  value={additionalInfo.floor}
                  onChange={(e) =>
                    setAdditionalInfo({
                      ...additionalInfo,
                      floor: e.target.value,
                    })
                  }
                />
              </div>

              {/* Unit number */}
              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    {t('location.unit_number')}
                  </span>
                </label>
                <input
                  type="text"
                  placeholder={t('location.unit_number_placeholder')}
                  className="input input-bordered w-full"
                  value={additionalInfo.suite}
                  onChange={(e) =>
                    setAdditionalInfo({
                      ...additionalInfo,
                      suite: e.target.value,
                    })
                  }
                />
              </div>
            </div>

            <div className="space-y-4">
              {/* Amenities */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">{t('location.amenities')}</span>
                </label>
                <div className="space-y-2">
                  <label className="label cursor-pointer justify-start gap-3">
                    <input
                      type="checkbox"
                      className="checkbox"
                      checked={additionalInfo.hasParking}
                      onChange={(e) =>
                        setAdditionalInfo({
                          ...additionalInfo,
                          hasParking: e.target.checked,
                        })
                      }
                    />
                    <span className="label-text">
                      🚗 {t('location.has_parking')}
                    </span>
                  </label>
                  <label className="label cursor-pointer justify-start gap-3">
                    <input
                      type="checkbox"
                      className="checkbox"
                      checked={additionalInfo.hasElevator}
                      onChange={(e) =>
                        setAdditionalInfo({
                          ...additionalInfo,
                          hasElevator: e.target.checked,
                        })
                      }
                    />
                    <span className="label-text">
                      🛗 {t('location.has_elevator')}
                    </span>
                  </label>
                </div>
              </div>

              {/* Accessibility notes */}
              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    {t('location.accessibility_notes')}
                  </span>
                </label>
                <textarea
                  className="textarea textarea-bordered h-24"
                  placeholder={t('location.accessibility_notes_placeholder')}
                  value={additionalInfo.accessibilityNotes}
                  onChange={(e) =>
                    setAdditionalInfo({
                      ...additionalInfo,
                      accessibilityNotes: e.target.value,
                    })
                  }
                />
              </div>
            </div>
          </div>

          {/* Selected location information */}
          {location && (
            <div className="mt-6 p-4 bg-info/10 border border-info/20 rounded-lg">
              <h4 className="font-medium text-info-content mb-2">
                📍 {t('location.selected_location')}
              </h4>
              <div className="text-sm text-info-content/80 grid grid-cols-1 md:grid-cols-2 gap-2">
                <div>
                  <p>
                    <strong>{t('location.address')}:</strong> {location.address}
                  </p>
                  {location.city && (
                    <p>
                      <strong>{t('location.city')}:</strong> {location.city}
                    </p>
                  )}
                </div>
                <div>
                  <p>
                    <strong>{t('location.coordinates')}:</strong>{' '}
                    {location.latitude.toFixed(6)},{' '}
                    {location.longitude.toFixed(6)}
                  </p>
                  <p>
                    <strong>{t('location.accuracy')}:</strong>{' '}
                    {Math.round(location.confidence * 100)}%
                  </p>
                </div>
              </div>
            </div>
          )}

          <div className="card-actions justify-between mt-6">
            <button className="btn btn-ghost" onClick={onBack}>
              {tCommon('back')}
            </button>
            <button className="btn btn-primary" onClick={handleNext}>
              {tCommon('next')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
