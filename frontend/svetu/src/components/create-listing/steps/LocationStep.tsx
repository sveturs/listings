'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';
import SmartAddressInput from '@/components/GIS/SmartAddressInput';
import AddressConfirmationMap from '@/components/GIS/AddressConfirmationMap';
import LocationPrivacySettings from '@/components/GIS/LocationPrivacySettings';
import { AddressGeocodingResult } from '@/hooks/useAddressGeocoding';

interface LocationStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function LocationStep({ onNext, onBack }: LocationStepProps) {
  const t = useTranslations();
  const { state, dispatch } = useCreateListing();
  const [step, setStep] = useState<'input' | 'confirm' | 'privacy'>('input');
  const [address, setAddress] = useState(state.location?.address || '');
  const [location, setLocation] = useState<
    { lat: number; lng: number } | undefined
  >(
    state.location?.latitude && state.location?.longitude
      ? { lat: state.location.latitude, lng: state.location.longitude }
      : undefined
  );
  const [confidence, setConfidence] = useState(0);
  const [privacyLevel, setPrivacyLevel] = useState<
    'exact' | 'street' | 'district' | 'city'
  >('street');
  const [formData, setFormData] = useState({
    country: state.location?.country || 'Србија',
    safeMeetingPlaces: [] as string[],
  });

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

  useEffect(() => {
    if (location && address) {
      const locationData = {
        latitude: location.lat,
        longitude: location.lng,
        address: address,
        city: '', // Будет извлечено из адреса через геокодирование
        region: '',
        country: formData.country,
        privacyLevel: privacyLevel,
        confidence: confidence,
      };

      dispatch({ type: 'SET_LOCATION', payload: locationData });
    }
  }, [location, address, formData.country, privacyLevel, confidence, dispatch]);

  const toggleSafeMeetingPlace = (place: string) => {
    setFormData((prev) => ({
      ...prev,
      safeMeetingPlaces: prev.safeMeetingPlaces.includes(place)
        ? prev.safeMeetingPlaces.filter((p) => p !== place)
        : [...prev.safeMeetingPlaces, place],
    }));
  };

  const canProceed = address && location && step === 'privacy';

  const handleAddressChange = (
    value: string,
    result?: AddressGeocodingResult
  ) => {
    setAddress(value);

    if (result) {
      setLocation({
        lat: result.location.lat,
        lng: result.location.lng,
      });
      setConfidence(result.confidence);
    }
  };

  const handleLocationSelect = (locationData: {
    lat: number;
    lng: number;
    address: string;
    confidence: number;
  }) => {
    setLocation({ lat: locationData.lat, lng: locationData.lng });
    setAddress(locationData.address);
    setConfidence(locationData.confidence);
    setStep('confirm');
  };

  const handleLocationConfirm = (locationData: {
    lat: number;
    lng: number;
    address: string;
    confidence: number;
  }) => {
    setLocation({ lat: locationData.lat, lng: locationData.lng });
    setAddress(locationData.address);
    setConfidence(locationData.confidence);
    setStep('privacy');
  };

  const handleLocationChange = (newLocation: { lat: number; lng: number }) => {
    setLocation(newLocation);
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-8">
        <h2 className="text-2xl font-bold mb-2 flex items-center">
          📍 {t('create_listing.location.title')}
        </h2>
        <p className="text-base-content/70">
          Используйте умный ввод адресов для точного указания местоположения
        </p>
      </div>

      {/* Навигация по шагам */}
      <div className="mb-8">
        <div className="flex justify-center">
          <div className="steps">
            <div
              className={`step ${step === 'input' ? 'step-primary' : ''} ${location ? 'step-success' : ''}`}
              onClick={() => setStep('input')}
            >
              Ввод адреса
            </div>
            <div
              className={`step ${step === 'confirm' ? 'step-primary' : ''} ${step === 'privacy' ? 'step-success' : ''}`}
              onClick={() => location && setStep('confirm')}
            >
              Подтверждение
            </div>
            <div
              className={`step ${step === 'privacy' ? 'step-primary' : ''}`}
              onClick={() => location && setStep('privacy')}
            >
              Приватность
            </div>
          </div>
        </div>
      </div>

      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          {/* Шаг 1: Ввод адреса */}
          {step === 'input' && (
            <div className="space-y-6">
              <h3 className="text-lg font-semibold mb-4 flex items-center">
                <span className="text-2xl mr-2">📍</span>
                Шаг 1: Введите адрес
              </h3>

              <SmartAddressInput
                value={address}
                onChange={handleAddressChange}
                onLocationSelect={handleLocationSelect}
                placeholder="Начните вводить адрес (например: Београд, Кнез Михаилова)"
                showCurrentLocation={true}
                country={['rs', 'hr', 'ba', 'me']}
                language="ru"
              />

              {location && (
                <div className="mt-4 p-4 bg-success/10 border border-success/20 rounded-lg">
                  <h4 className="font-medium text-success-content mb-2">
                    ✅ Адрес найден!
                  </h4>
                  <div className="text-sm text-success-content/80 space-y-1">
                    <p>
                      <strong>Адрес:</strong> {address}
                    </p>
                    <p>
                      <strong>Координаты:</strong> {location.lat.toFixed(6)},{' '}
                      {location.lng.toFixed(6)}
                    </p>
                    <p>
                      <strong>Точность:</strong> {Math.round(confidence * 100)}%
                    </p>
                  </div>

                  <div className="mt-3">
                    <button
                      className="btn btn-primary"
                      onClick={() => setStep('confirm')}
                    >
                      Перейти к подтверждению
                    </button>
                  </div>
                </div>
              )}

              {/* Страна */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">🌍 Страна</span>
                </label>
                <select
                  className="select select-bordered"
                  value={formData.country}
                  onChange={(e) =>
                    setFormData((prev) => ({
                      ...prev,
                      country: e.target.value,
                    }))
                  }
                >
                  <option value="Србија">🇷🇸 Србија</option>
                  <option value="Хрватска">🇭🇷 Хрватска</option>
                  <option value="Босна и Херцеговина">
                    🇧🇦 Босна и Херцеговина
                  </option>
                  <option value="Црна Гора">🇲🇪 Црна Гора</option>
                  <option value="Словенија">🇸🇮 Словенија</option>
                  <option value="Македонија">🇲🇰 Македонија</option>
                </select>
              </div>
            </div>
          )}

          {/* Шаг 2: Подтверждение на карте */}
          {step === 'confirm' && location && (
            <div className="space-y-6">
              <h3 className="text-lg font-semibold mb-4 flex items-center">
                <span className="text-2xl mr-2">🗺️</span>
                Шаг 2: Подтвердите местоположение на карте
              </h3>

              <AddressConfirmationMap
                address={address}
                initialLocation={location}
                onLocationConfirm={handleLocationConfirm}
                onLocationChange={handleLocationChange}
                editable={true}
                zoom={16}
                height="500px"
              />
            </div>
          )}

          {/* Шаг 3: Настройки приватности */}
          {step === 'privacy' && location && (
            <div className="space-y-6">
              <h3 className="text-lg font-semibold mb-4 flex items-center">
                <span className="text-2xl mr-2">🛡️</span>
                Шаг 3: Настройки приватности
              </h3>

              <LocationPrivacySettings
                selectedLevel={privacyLevel}
                onLevelChange={setPrivacyLevel}
                location={location}
                showPreview={true}
              />

              {/* Безбедна места за састанак */}
              <div className="form-control">
                <label className="label">
                  <span className="label-text font-medium">
                    🛡️ Безбедна места за састанак
                  </span>
                </label>
                <p className="text-sm text-base-content/60 mb-3">
                  Препоручите безбедна места за састанак у вашој близини
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
                          formData.safeMeetingPlaces.includes(place)
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
            </div>
          )}

          {/* Кнопки навигации */}
          <div className="card-actions justify-between mt-6">
            <button className="btn btn-outline" onClick={onBack}>
              ← {t('common.back')}
            </button>

            {step === 'confirm' && (
              <button
                className="btn btn-outline"
                onClick={() => setStep('input')}
              >
                ← Назад к вводу
              </button>
            )}

            {step === 'privacy' && (
              <button
                className="btn btn-outline"
                onClick={() => setStep('confirm')}
              >
                ← Назад к карте
              </button>
            )}

            <button
              className={`btn btn-primary ${!canProceed ? 'btn-disabled' : ''}`}
              onClick={onNext}
              disabled={!canProceed}
            >
              {step === 'privacy'
                ? 'Сохранить местоположение'
                : t('common.continue')}{' '}
              →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
