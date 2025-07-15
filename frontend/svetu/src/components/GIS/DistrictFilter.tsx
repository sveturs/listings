'use client';

import React, { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { useVisibleCitiesContext } from './contexts/VisibleCitiesContext';

interface District {
  id: string;
  name: string;
  country_code: string;
  population?: number;
  area_km2?: number;
}

interface Municipality {
  id: string;
  name: string;
  district_id: string;
  country_code: string;
  population?: number;
  area_km2?: number;
}

interface DistrictFilterProps {
  onDistrictSelect?: (district: District | null) => void;
  onMunicipalitySelect?: (municipality: Municipality | null) => void;
  className?: string;
  disabled?: boolean;
}

export default function DistrictFilter({
  onDistrictSelect,
  onMunicipalitySelect,
  className = '',
  disabled = false,
}: DistrictFilterProps) {
  const t = useTranslations('gis');

  // Используем контекст для получения районов текущего города
  const { availableDistricts, closestCity, loading: citiesLoading } = useVisibleCitiesContext();

  // Логирование для отладки
  console.log('🏗️ DistrictFilter render:', {
    availableDistricts: availableDistricts.length,
    closestCity: closestCity?.city.name,
    citiesLoading
  });

  const [municipalities, setMunicipalities] = useState<Municipality[]>([]);
  const [selectedDistrict, setSelectedDistrict] = useState<District | null>(
    null
  );
  const [selectedMunicipality, setSelectedMunicipality] =
    useState<Municipality | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Загрузка муниципалитетов для выбранного района
  const loadMunicipalities = useCallback(
    async (districtId: string) => {
      setLoading(true);
      setError(null);

      try {
        const response = await fetch(
          `/api/v1/gis/municipalities?district_id=${districtId}&country_code=RS`
        );
        if (!response.ok) {
          throw new Error(`HTTP ${response.status}: ${response.statusText}`);
        }

        const data = await response.json();
        setMunicipalities(data.data || []);
      } catch (err) {
        console.error('Failed to load municipalities:', err);
        setError(t('district_filter.load_error'));
      } finally {
        setLoading(false);
      }
    },
    [t]
  );

  // Сбрасываем выбранный район, если изменился город
  useEffect(() => {
    if (selectedDistrict) {
      // Проверяем, есть ли выбранный район в списке доступных районов
      const stillAvailable = availableDistricts.some(
        (district) => district.id === selectedDistrict.id
      );

      if (!stillAvailable) {
        setSelectedDistrict(null);
        setSelectedMunicipality(null);
        if (onDistrictSelect) onDistrictSelect(null);
        if (onMunicipalitySelect) onMunicipalitySelect(null);
      }
    }
  }, [availableDistricts, selectedDistrict, onDistrictSelect, onMunicipalitySelect]);

  // Обработка выбора района
  const handleDistrictChange = useCallback(
    (districtId: string) => {
      if (!districtId) {
        setSelectedDistrict(null);
        setSelectedMunicipality(null);
        setMunicipalities([]);
        onDistrictSelect?.(null);
        onMunicipalitySelect?.(null);
        return;
      }

      const district = availableDistricts.find((d) => d.id === districtId) || null;
      setSelectedDistrict(district);
      setSelectedMunicipality(null);
      onDistrictSelect?.(district);
      onMunicipalitySelect?.(null);

      if (district) {
        loadMunicipalities(district.id);
      }
    },
    [availableDistricts, onDistrictSelect, onMunicipalitySelect, loadMunicipalities]
  );

  // Обработка выбора муниципалитета
  const handleMunicipalityChange = useCallback(
    (municipalityId: string) => {
      if (!municipalityId) {
        setSelectedMunicipality(null);
        onMunicipalitySelect?.(null);
        return;
      }

      const municipality =
        municipalities.find((m) => m.id === municipalityId) || null;
      setSelectedMunicipality(municipality);
      onMunicipalitySelect?.(municipality);
    },
    [municipalities, onMunicipalitySelect]
  );

  return (
    <div className={`space-y-3 ${className}`}>
      {/* Выбор района */}
      <div>
        <label className="block text-sm font-medium text-base-content mb-1">
          {t('district_filter.district')}
        </label>
        <select
          value={selectedDistrict?.id || ''}
          onChange={(e) => handleDistrictChange(e.target.value)}
          className="select select-bordered select-sm w-full"
          disabled={disabled || loading || citiesLoading}
        >
          <option value="">{t('district_filter.select_district')}</option>
          {closestCity && (
            <optgroup label={`${t('district_filter.districts_in')} ${closestCity.city.name}`}>
              {availableDistricts.map((district) => (
                <option key={district.id} value={district.id}>
                  {district.name}
                  {district.population &&
                    ` (${district.population.toLocaleString()})`}
                </option>
              ))}
            </optgroup>
          )}
        </select>
      </div>

      {/* Выбор муниципалитета */}
      {selectedDistrict && (
        <div>
          <label className="block text-sm font-medium text-base-content mb-1">
            {t('district_filter.municipality')}
          </label>
          <select
            value={selectedMunicipality?.id || ''}
            onChange={(e) => handleMunicipalityChange(e.target.value)}
            className="select select-bordered select-sm w-full"
            disabled={disabled || loading}
          >
            <option value="">{t('district_filter.select_municipality')}</option>
            {municipalities.map((municipality) => (
              <option key={municipality.id} value={municipality.id}>
                {municipality.name}
                {municipality.population &&
                  ` (${municipality.population.toLocaleString()})`}
              </option>
            ))}
          </select>
        </div>
      )}

      {/* Индикатор загрузки */}
      {loading && (
        <div className="flex items-center space-x-2 text-sm text-base-content/70">
          <div className="loading loading-spinner loading-xs"></div>
          <span>{t('district_filter.loading')}</span>
        </div>
      )}

      {/* Ошибка */}
      {error && (
        <div className="alert alert-error alert-sm">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-4 w-4"
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
          <span className="text-xs">{error}</span>
        </div>
      )}
    </div>
  );
}
