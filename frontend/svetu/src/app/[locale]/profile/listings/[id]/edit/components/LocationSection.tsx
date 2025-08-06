'use client';

import { useTranslations } from 'next-intl';
import { useState } from 'react';
import dynamic from 'next/dynamic';

// Динамически импортируем LocationPicker чтобы избежать SSR проблем
const LocationPicker = dynamic(
  () => import('@/components/GIS/LocationPicker'),
  {
    ssr: false,
    loading: () => (
      <div className="h-64 bg-base-200 animate-pulse rounded-lg" />
    ),
  }
);

interface LocationSectionProps {
  data: {
    city: string;
    country: string;
    location: string;
    latitude?: number;
    longitude?: number;
    show_on_map: boolean;
  };
  errors?: Record<string, string>;
  onChange: (field: string, value: string | number | boolean) => void;
}

export function LocationSection({
  data,
  errors = {},
  onChange,
}: LocationSectionProps) {
  const t = useTranslations('profile');
  const [privacyLevel, setPrivacyLevel] = useState<'exact' | 'area' | 'city'>(
    data.show_on_map ? 'exact' : 'area'
  );

  const handleLocationSelect = (location: {
    latitude: number;
    longitude: number;
    address: string;
    city: string;
    region: string;
    country: string;
    confidence: number;
    addressMultilingual?: {
      sr: string;
      en: string;
      ru: string;
    };
  }) => {
    onChange('location', location.address);
    onChange('city', location.city);
    onChange('country', location.country);
    onChange('latitude', location.latitude);
    onChange('longitude', location.longitude);
  };

  const handlePrivacyChange = (level: 'exact' | 'area' | 'city') => {
    setPrivacyLevel(level);
    onChange('show_on_map', level === 'exact');
  };

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold mb-4 flex items-center gap-2">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            strokeWidth={1.5}
            stroke="currentColor"
            className="w-5 h-5 text-primary"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"
            />
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"
            />
          </svg>
          {t('location.title')}
        </h3>

        {/* Map - Always visible */}
        <div className="mb-6">
          <div className="mb-2 flex justify-between items-center">
            <span className="text-sm font-medium">
              {t('location.clickToSelect')}
            </span>
            {data.latitude && data.longitude && (
              <span className="text-xs text-base-content/60">
                {t('location.coordinates')}: {data.latitude.toFixed(6)},{' '}
                {data.longitude.toFixed(6)}
              </span>
            )}
          </div>
          <div className="h-96 rounded-lg overflow-hidden border-2 border-base-300 shadow-lg">
            <LocationPicker
              value={
                data.latitude && data.longitude
                  ? {
                      latitude: data.latitude,
                      longitude: data.longitude,
                      address: data.location,
                      city: data.city,
                      region: '',
                      country: data.country,
                      confidence: 1,
                    }
                  : undefined
              }
              onChange={handleLocationSelect}
              height="384px"
            />
          </div>
        </div>

        {/* Privacy Settings */}
        <div className="card bg-base-100 border border-base-300 mb-6">
          <div className="card-body">
            <h4 className="card-title text-base flex items-center gap-2">
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="w-5 h-5"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z"
                />
              </svg>
              {t('location.privacySettings')}
            </h4>
            <div className="space-y-3">
              <label className="cursor-pointer">
                <div className="flex items-start gap-3">
                  <input
                    type="radio"
                    name="privacy"
                    className="radio radio-primary mt-1"
                    checked={privacyLevel === 'exact'}
                    onChange={() => handlePrivacyChange('exact')}
                  />
                  <div className="flex-1">
                    <div className="font-medium">
                      {t('location.privacy.exact')}
                    </div>
                    <div className="text-sm text-base-content/60">
                      {t('location.privacy.exactDesc')}
                    </div>
                  </div>
                </div>
              </label>

              <label className="cursor-pointer">
                <div className="flex items-start gap-3">
                  <input
                    type="radio"
                    name="privacy"
                    className="radio radio-primary mt-1"
                    checked={privacyLevel === 'area'}
                    onChange={() => handlePrivacyChange('area')}
                  />
                  <div className="flex-1">
                    <div className="font-medium">
                      {t('location.privacy.area')}
                    </div>
                    <div className="text-sm text-base-content/60">
                      {t('location.privacy.areaDesc')}
                    </div>
                  </div>
                </div>
              </label>

              <label className="cursor-pointer">
                <div className="flex items-start gap-3">
                  <input
                    type="radio"
                    name="privacy"
                    className="radio radio-primary mt-1"
                    checked={privacyLevel === 'city'}
                    onChange={() => handlePrivacyChange('city')}
                  />
                  <div className="flex-1">
                    <div className="font-medium">
                      {t('location.privacy.city')}
                    </div>
                    <div className="text-sm text-base-content/60">
                      {t('location.privacy.cityDesc')}
                    </div>
                  </div>
                </div>
              </label>
            </div>
          </div>
        </div>

        {/* Location fields */}
        <div className="space-y-4">
          <h4 className="font-medium text-base">
            {t('location.addressDetails')}
          </h4>

          <div className="form-control">
            <label className="label">
              <span className="label-text font-semibold flex items-center gap-2">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  strokeWidth={1.5}
                  stroke="currentColor"
                  className="w-4 h-4"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    d="M2.25 12l8.954-8.955c.44-.439 1.152-.439 1.591 0L21.75 12M4.5 9.75v10.125c0 .621.504 1.125 1.125 1.125H9.75v-4.875c0-.621.504-1.125 1.125-1.125h2.25c.621 0 1.125.504 1.125 1.125V21h4.125c.621 0 1.125-.504 1.125-1.125V9.75M8.25 21h8.25"
                  />
                </svg>
                {t('fields.location')}
              </span>
            </label>
            <input
              type="text"
              className={`input input-bordered w-full ${errors.location ? 'input-error' : ''} focus:input-primary transition-all`}
              value={data.location}
              onChange={(e) => onChange('location', e.target.value)}
              placeholder={t('fields.locationPlaceholder')}
            />
            {errors.location && (
              <label className="label">
                <span className="label-text-alt text-error">
                  {errors.location}
                </span>
              </label>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text font-semibold flex items-center gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z"
                    />
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z"
                    />
                  </svg>
                  {t('fields.city')}
                </span>
              </label>
              <input
                type="text"
                className={`input input-bordered w-full ${errors.city ? 'input-error' : ''} focus:input-primary transition-all`}
                value={data.city}
                onChange={(e) => onChange('city', e.target.value)}
                placeholder={t('fields.cityPlaceholder')}
              />
              {errors.city && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {errors.city}
                  </span>
                </label>
              )}
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text font-semibold flex items-center gap-2">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    strokeWidth={1.5}
                    stroke="currentColor"
                    className="w-4 h-4"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      d="M12 21a9.004 9.004 0 008.716-6.747M12 21a9.004 9.004 0 01-8.716-6.747M12 21c2.485 0 4.5-4.03 4.5-9S14.485 3 12 3m0 18c-2.485 0-4.5-4.03-4.5-9S9.515 3 12 3m0 0a8.997 8.997 0 017.843 4.582M12 3a8.997 8.997 0 00-7.843 4.582m15.686 0A11.953 11.953 0 0112 10.5c-2.998 0-5.74-1.1-7.843-2.918m15.686 0A8.959 8.959 0 0121 12c0 .778-.099 1.533-.284 2.253m0 0A17.919 17.919 0 0112 16.5c-3.162 0-6.133-.815-8.716-2.247m0 0A9.015 9.015 0 013 12c0-1.605.42-3.113 1.157-4.418"
                    />
                  </svg>
                  {t('fields.country')}
                </span>
              </label>
              <input
                type="text"
                className={`input input-bordered w-full ${errors.country ? 'input-error' : ''} focus:input-primary transition-all`}
                value={data.country}
                onChange={(e) => onChange('country', e.target.value)}
                placeholder={t('fields.countryPlaceholder')}
              />
              {errors.country && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {errors.country}
                  </span>
                </label>
              )}
            </div>
          </div>
        </div>

        {/* Privacy info */}
        <div className="alert shadow-lg mt-6">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            className="stroke-info shrink-0 w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            ></path>
          </svg>
          <div>
            <h3 className="font-bold">{t('location.privacyTitle')}</h3>
            <div className="text-xs">{t('location.privacyNote')}</div>
          </div>
        </div>
      </div>
    </div>
  );
}
