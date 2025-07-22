'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';
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
  const t = useTranslations();
  const { state, dispatch } = useCreateListing();
  const [step, setStep] = useState<'select' | 'privacy'>('select');
  const [location, setLocation] = useState<LocationData | undefined>(
    state.location?.latitude && state.location?.longitude
      ? {
          latitude: state.location.latitude,
          longitude: state.location.longitude,
          address: state.location.address || '',
          city: state.location.city || '',
          region: state.location.region || '',
          country: state.location.country || 'Србија',
          confidence: 0.9,
        }
      : undefined
  );
  const [privacyLevel, setPrivacyLevel] = useState<
    'exact' | 'street' | 'district' | 'city'
  >('street');
  const [safeMeetingPlaces, setSafeMeetingPlaces] = useState<string[]>([]);

  const safeMeetingOptions = [
    'Тржни центар',
    'Главни трг',
    'Аутобуска станица',
    'Железничка станица',
    'Паркинг тржног центра',
    'Кафе/ресторан',
    'Банка',
    'Пошта',
    'Бензинска пумпа',
  ];

  // Обновляем данные в контексте при изменении location
  useEffect(() => {
    if (location) {
      dispatch({
        type: 'SET_LOCATION',
        payload: {
          latitude: location.latitude,
          longitude: location.longitude,
          address: location.address,
          city: location.city,
          region: location.region,
          country: location.country,
        },
      });
    }
  }, [location, dispatch]);

  const handleLocationChange = (locationData: LocationData) => {
    setLocation(locationData);
    // Убираем автоматический переход - пользователь сам решит когда переходить
  };

  const toggleSafeMeetingPlace = (place: string) => {
    setSafeMeetingPlaces((prev) =>
      prev.includes(place) ? prev.filter((p) => p !== place) : [...prev, place]
    );
  };

  const canProceed = location !== undefined;

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-8">
        <h2 className="text-2xl font-bold mb-2 flex items-center">
          📍 {t('create_listing.location.title')}
        </h2>
        <p className="text-base-content/70">
          Выберите местоположение объявления - введите адрес или укажите точку
          на карте
        </p>
      </div>

      {/* Навигация по шагам */}
      <div className="mb-8">
        <div className="flex justify-center">
          <div className="steps">
            <div
              className={`step ${step === 'select' ? 'step-primary' : ''} ${location ? 'step-success' : ''}`}
              onClick={() => setStep('select')}
            >
              Выбор местоположения
            </div>
            <div
              className={`step ${step === 'privacy' ? 'step-primary' : ''}`}
              onClick={() => location && setStep('privacy')}
            >
              Настройки приватности
            </div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          {/* Шаг 1: Выбор местоположения */}
          {step === 'select' && (
            <div className="space-y-6">
              <h3 className="text-lg font-semibold mb-4 flex items-center">
                <span className="text-2xl mr-2">📍</span>
                Шаг 1: Выберите местоположение
              </h3>

              <LocationPicker
                value={location}
                onChange={handleLocationChange}
                placeholder="Начните вводить адрес (например: Београд, Кнез Михаилова)"
                height="500px"
                showCurrentLocation={true}
                defaultCountry={state.location?.country || 'Србија'}
              />

              {location && (
                <div className="mt-6 flex justify-end">
                  <button
                    type="button"
                    onClick={() => setStep('privacy')}
                    className="btn btn-primary"
                  >
                    Далее к настройкам приватности →
                  </button>
                </div>
              )}
            </div>
          )}

          {/* Шаг 2: Настройки приватности */}
          {step === 'privacy' && location && (
            <div className="space-y-6">
              <h3 className="text-lg font-semibold mb-4 flex items-center">
                <span className="text-2xl mr-2">🛡️</span>
                Шаг 2: Настройки приватности
              </h3>

              <LocationPrivacySettings
                selectedLevel={privacyLevel}
                onLevelChange={setPrivacyLevel}
                location={{ lat: location.latitude, lng: location.longitude }}
                showPreview={true}
              />

              {/* Безопасные места для встречи */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    🛡️ Безопасные места для встречи
                  </span>
                </label>
                <p className="text-sm text-base-content/60 mb-3">
                  Рекомендуйте безопасные места для встречи в вашей близости
                </p>

                <div className="grid grid-cols-2 gap-2">
                  {safeMeetingOptions.map((place) => (
                    <button
                      key={place}
                      type="button"
                      onClick={() => toggleSafeMeetingPlace(place)}
                      className={`
                        btn btn-sm text-xs
                        ${
                          safeMeetingPlaces.includes(place)
                            ? 'btn-primary'
                            : 'btn-outline'
                        }
                      `}
                    >
                      {place}
                    </button>
                  ))}
                </div>
              </div>

              {/* Информация о выбранном местоположении */}
              <div className="p-4 bg-info/10 border border-info/20 rounded-lg">
                <h4 className="font-medium text-info-content mb-2">
                  📍 Выбранное местоположение
                </h4>
                <div className="text-sm text-info-content/80 space-y-1">
                  <p>
                    <strong>Адрес:</strong> {location.address}
                  </p>
                  {location.city && (
                    <p>
                      <strong>Город:</strong> {location.city}
                    </p>
                  )}
                  <p>
                    <strong>Координаты:</strong> {location.latitude.toFixed(6)},{' '}
                    {location.longitude.toFixed(6)}
                  </p>
                  <p>
                    <strong>Уровень приватности:</strong>{' '}
                    {privacyLevel === 'exact'
                      ? 'Точный адрес'
                      : privacyLevel === 'street'
                        ? 'Улица'
                        : privacyLevel === 'district'
                          ? 'Район'
                          : 'Город'}
                  </p>
                </div>
              </div>
            </div>
          )}

          {/* Кнопки навигации */}
          <div className="card-actions justify-between mt-6">
            <button className="btn btn-outline" onClick={onBack}>
              ← {t('common.back')}
            </button>

            {step === 'privacy' && (
              <button
                className="btn btn-outline"
                onClick={() => setStep('select')}
              >
                ← Изменить местоположение
              </button>
            )}

            <button
              className={`btn btn-primary ${!canProceed ? 'btn-disabled' : ''}`}
              onClick={() => {
                if (step === 'select' && location) {
                  setStep('privacy');
                } else {
                  onNext();
                }
              }}
              disabled={!canProceed}
            >
              {step === 'privacy'
                ? 'Сохранить и продолжить'
                : 'Далее к настройкам приватности'}{' '}
              →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
