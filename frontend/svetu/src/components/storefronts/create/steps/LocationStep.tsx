'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateStorefrontContext } from '@/contexts/CreateStorefrontContext';
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
  const tCreate_storefront.location = useTranslations('create_storefront');
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
          <h2 className="card-title text-2xl mb-4">
            {tCreate_storefront.location('title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {tCreate_storefront.location('subtitle')}
          </p>

          {/* Выбор местоположения */}
          <div className="mb-6">
            <h3 className="text-lg font-semibold mb-4">
              📍 Местоположение витрины
            </h3>
            <LocationPicker
              value={location}
              onChange={handleLocationChange}
              placeholder="Введите адрес вашей витрины или выберите точку на карте"
              height="400px"
              showCurrentLocation={false}
              defaultCountry="Србија"
            />
            {errors.location && (
              <p className="text-error text-sm mt-2">{errors.location}</p>
            )}
          </div>

          {/* Дополнительная информация */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              {/* Почтовый индекс */}
              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    {tCreate_storefront.location('postal_code')}
                  </span>
                  <span className="label-text-alt text-error">*</span>
                </label>
                <input
                  type="text"
                  placeholder={t(
                    'create_storefront.location.postal_code_placeholder'
                  )}
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

              {/* Этаж */}
              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">Этаж (необязательно)</span>
                </label>
                <input
                  type="text"
                  placeholder="Например: 2, приземље, подрум"
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

              {/* Номер помещения */}
              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    Номер помещения/офиса (необязательно)
                  </span>
                </label>
                <input
                  type="text"
                  placeholder="Например: 12, A3, Локал 5"
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
              {/* Удобства */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text">Удобства</span>
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
                    <span className="label-text">🚗 Есть парковка</span>
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
                    <span className="label-text">🛗 Есть лифт</span>
                  </label>
                </div>
              </div>

              {/* Заметки о доступности */}
              <div className="form-control w-full">
                <label className="label">
                  <span className="label-text">
                    Заметки о доступности (необязательно)
                  </span>
                </label>
                <textarea
                  className="textarea textarea-bordered h-24"
                  placeholder="Например: вход со двора, рампа для инвалидов, широкие двери"
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

          {/* Информация о выбранном местоположении */}
          {location && (
            <div className="mt-6 p-4 bg-info/10 border border-info/20 rounded-lg">
              <h4 className="font-medium text-info-content mb-2">
                📍 Выбранное местоположение
              </h4>
              <div className="text-sm text-info-content/80 grid grid-cols-1 md:grid-cols-2 gap-2">
                <div>
                  <p>
                    <strong>Адрес:</strong> {location.address}
                  </p>
                  {location.city && (
                    <p>
                      <strong>Город:</strong> {location.city}
                    </p>
                  )}
                </div>
                <div>
                  <p>
                    <strong>Координаты:</strong> {location.latitude.toFixed(6)},{' '}
                    {location.longitude.toFixed(6)}
                  </p>
                  <p>
                    <strong>Точность:</strong>{' '}
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
