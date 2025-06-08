'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';

interface LocationStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function LocationStep({ onNext, onBack }: LocationStepProps) {
  const t = useTranslations();
  const { state, dispatch } = useCreateListing();
  const [formData, setFormData] = useState({
    address: state.location?.address || '',
    city: state.location?.city || '',
    region: state.location?.region || '',
    country: state.location?.country || 'Србија',
    exactLocation: false,
    approximateArea: '',
    safeMeetingPlaces: [] as string[],
  });

  const [currentLocation, setCurrentLocation] = useState<{
    lat: number;
    lng: number;
  } | null>(null);
  const [loadingLocation, setLoadingLocation] = useState(false);

  // TODO: Загружать список городов из API
  // ВРЕМЕННОЕ РЕШЕНИЕ: Хардкодные города Сербии
  const serbianCities = [
    'Београд',
    'Нови Сад',
    'Ниш',
    'Крагујевац',
    'Суботица',
    'Панчево',
    'Чачак',
    'Нови Пазар',
    'Зрењанин',
    'Лесковац',
    'Врање',
    'Кикинда',
    'Сомбор',
    'Ужице',
    'Пријепоље',
    'Смедерево',
    'Валјево',
    'Пирот',
    'Бор',
    'Шабац',
    'Остало',
  ];

  const serbianRegions = [
    'Београд (главни град)',
    'Војводина',
    'Централна Србија',
    'Јужна и источна Србија',
  ];

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
    const locationData = {
      latitude: currentLocation?.lat || 0,
      longitude: currentLocation?.lng || 0,
      address: formData.address,
      city: formData.city,
      region: formData.region,
      country: formData.country,
    };

    dispatch({ type: 'SET_LOCATION', payload: locationData });
  }, [formData, currentLocation, dispatch]);

  const getCurrentLocation = () => {
    setLoadingLocation(true);

    if ('geolocation' in navigator) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          setCurrentLocation({
            lat: position.coords.latitude,
            lng: position.coords.longitude,
          });
          setFormData((prev) => ({ ...prev, exactLocation: true }));
          setLoadingLocation(false);
        },
        (error) => {
          console.error('Error getting location:', error);
          setLoadingLocation(false);
        },
        { enableHighAccuracy: true, timeout: 10000 }
      );
    } else {
      setLoadingLocation(false);
    }
  };

  const toggleSafeMeetingPlace = (place: string) => {
    setFormData((prev) => ({
      ...prev,
      safeMeetingPlaces: prev.safeMeetingPlaces.includes(place)
        ? prev.safeMeetingPlaces.filter((p) => p !== place)
        : [...prev.safeMeetingPlaces, place],
    }));
  };

  const canProceed = formData.city && formData.region;

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            📍 {t('create_listing.location.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('create_listing.location.description')}
          </p>

          <div className="space-y-6">
            {/* Страна */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  🌍 {t('create_listing.location.country')}
                </span>
                <span className="label-text-alt text-error">*</span>
              </label>
              <select
                className="select select-bordered"
                value={formData.country}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, country: e.target.value }))
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

            {/* Регион */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  🗺️ {t('create_listing.location.region')}
                </span>
                <span className="label-text-alt text-error">*</span>
              </label>
              <select
                className="select select-bordered"
                value={formData.region}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, region: e.target.value }))
                }
              >
                <option value="">{t('common.select')}</option>
                {serbianRegions.map((region) => (
                  <option key={region} value={region}>
                    {region}
                  </option>
                ))}
              </select>
            </div>

            {/* Град */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  🏙️ {t('create_listing.location.city')}
                </span>
                <span className="label-text-alt text-error">*</span>
              </label>
              <select
                className="select select-bordered"
                value={formData.city}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, city: e.target.value }))
                }
              >
                <option value="">{t('common.select')}</option>
                {serbianCities.map((city) => (
                  <option key={city} value={city}>
                    {city}
                  </option>
                ))}
              </select>
              {formData.city === 'Остало' && (
                <input
                  type="text"
                  placeholder={t('create_listing.location.custom_city')}
                  className="input input-bordered mt-2"
                  onChange={(e) =>
                    setFormData((prev) => ({ ...prev, city: e.target.value }))
                  }
                />
              )}
            </div>

            {/* Адреса */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  🏠 {t('create_listing.location.address')}
                </span>
              </label>
              <input
                type="text"
                placeholder={t('create_listing.location.address_placeholder')}
                className="input input-bordered"
                value={formData.address}
                onChange={(e) =>
                  setFormData((prev) => ({ ...prev, address: e.target.value }))
                }
              />
              <label className="label">
                <span className="label-text-alt text-base-content/60">
                  {t('create_listing.location.address_hint')}
                </span>
              </label>
            </div>

            {/* Тачна локација */}
            <div className="form-control">
              <label className="label cursor-pointer">
                <div className="flex items-center gap-3">
                  <span className="text-xl">🎯</span>
                  <div>
                    <span className="label-text font-medium">
                      {t('create_listing.location.exact_location')}
                    </span>
                    <p className="text-sm text-base-content/60">
                      {t('create_listing.location.exact_location_desc')}
                    </p>
                  </div>
                </div>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={formData.exactLocation}
                  onChange={(e) =>
                    setFormData((prev) => ({
                      ...prev,
                      exactLocation: e.target.checked,
                    }))
                  }
                />
              </label>
            </div>

            {/* Кнопка получения текущей локации */}
            {formData.exactLocation && (
              <div className="card border border-primary/20 bg-primary/5">
                <div className="card-body p-4">
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="font-medium text-sm">
                        {t('create_listing.location.get_current')}
                      </h4>
                      <p className="text-xs text-base-content/60 mt-1">
                        {t('create_listing.location.get_current_desc')}
                      </p>
                    </div>
                    <button
                      onClick={getCurrentLocation}
                      disabled={loadingLocation}
                      className={`btn btn-sm btn-primary ${loadingLocation ? 'loading' : ''}`}
                    >
                      {loadingLocation ? '' : '📍'}{' '}
                      {t('create_listing.location.get_location')}
                    </button>
                  </div>
                  {currentLocation && (
                    <div className="mt-2 text-xs text-success">
                      ✅ {t('create_listing.location.location_obtained')}
                    </div>
                  )}
                </div>
              </div>
            )}

            {/* Безбедна места за састанак */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  🛡️ {t('create_listing.location.safe_meeting_places')}
                </span>
              </label>
              <p className="text-sm text-base-content/60 mb-3">
                {t('create_listing.location.safe_meeting_desc')}
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

            {/* Региональная подсказка о доверии */}
            <div className="alert alert-info">
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
                ></path>
              </svg>
              <div className="text-sm">
                <p className="font-medium">
                  🤝 {t('create_listing.location.trust_tip.title')}
                </p>
                <ul className="text-xs mt-2 space-y-1">
                  <li>
                    • {t('create_listing.location.trust_tip.public_places')}
                  </li>
                  <li>• {t('create_listing.location.trust_tip.daylight')}</li>
                  <li>
                    • {t('create_listing.location.trust_tip.known_areas')}
                  </li>
                </ul>
              </div>
            </div>
          </div>

          {/* Кнопки навигации */}
          <div className="card-actions justify-between mt-6">
            <button className="btn btn-outline" onClick={onBack}>
              ← {t('common.back')}
            </button>
            <button
              className={`btn btn-primary ${!canProceed ? 'btn-disabled' : ''}`}
              onClick={onNext}
              disabled={!canProceed}
            >
              {t('common.continue')} →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
