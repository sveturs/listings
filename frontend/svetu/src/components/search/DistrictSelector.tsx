'use client';

import React, { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import type { components } from '@/types/generated/api';

type District = components['schemas']['types.District'];
type Municipality = components['schemas']['types.Municipality'];

interface DistrictSelectorProps {
  selectedDistrictId?: string;
  selectedMunicipalityId?: string;
  onDistrictChange?: (districtId: string | null) => void;
  onMunicipalityChange?: (municipalityId: string | null) => void;
  className?: string;
}

export function DistrictSelector({
  selectedDistrictId,
  selectedMunicipalityId,
  onDistrictChange,
  onMunicipalityChange,
  className = '',
}: DistrictSelectorProps) {
  const t = useTranslations('search');
  const [districts, setDistricts] = useState<District[]>([]);
  const [municipalities, setMunicipalities] = useState<Municipality[]>([]);
  const [loadingDistricts, setLoadingDistricts] = useState(true);
  const [loadingMunicipalities, setLoadingMunicipalities] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Загрузка районов
  useEffect(() => {
    const fetchDistricts = async () => {
      try {
        setLoadingDistricts(true);
        const response = await fetch('/api/v1/gis/districts?country_code=RS');
        if (!response.ok) {
          throw new Error('Failed to fetch districts');
        }
        const data = await response.json();
        setDistricts(data.data || []);
      } catch (err) {
        console.error('Error fetching districts:', err);
        setError(t('errors.loadingDistricts'));
      } finally {
        setLoadingDistricts(false);
      }
    };

    fetchDistricts();
  }, [t]);

  // Загрузка муниципалитетов при выборе района
  useEffect(() => {
    if (!selectedDistrictId) {
      setMunicipalities([]);
      return;
    }

    const fetchMunicipalities = async () => {
      try {
        setLoadingMunicipalities(true);
        const response = await fetch(
          `/api/v1/gis/municipalities?district_id=${selectedDistrictId}`
        );
        if (!response.ok) {
          throw new Error('Failed to fetch municipalities');
        }
        const data = await response.json();
        setMunicipalities(data.data || []);
      } catch (err) {
        console.error('Error fetching municipalities:', err);
        setError(t('errors.loadingMunicipalities'));
      } finally {
        setLoadingMunicipalities(false);
      }
    };

    fetchMunicipalities();
  }, [selectedDistrictId, t]);

  const handleDistrictChange = (e: React.ChangeEvent<HTMLSelectElement>) => {
    const value = e.target.value || null;
    onDistrictChange?.(value);
    // Сброс выбранного муниципалитета при смене района
    onMunicipalityChange?.(null);
  };

  const handleMunicipalityChange = (
    e: React.ChangeEvent<HTMLSelectElement>
  ) => {
    const value = e.target.value || null;
    onMunicipalityChange?.(value);
  };

  if (error) {
    return (
      <div className={`alert alert-error ${className}`}>
        <span>{error}</span>
      </div>
    );
  }

  return (
    <div className={`flex flex-col gap-4 ${className}`}>
      {/* Выбор района */}
      <div className="form-control">
        <label className="label">
          <span className="label-text">{t('district')}</span>
        </label>
        <select
          className="select select-bordered w-full"
          value={selectedDistrictId || ''}
          onChange={handleDistrictChange}
          disabled={loadingDistricts}
        >
          <option value="">{t('allDistricts')}</option>
          {districts.map((district) => (
            <option key={district.id} value={district.id}>
              {district.name}
            </option>
          ))}
        </select>
      </div>

      {/* Выбор муниципалитета (только если выбран район) */}
      {selectedDistrictId && municipalities.length > 0 && (
        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('municipality')}</span>
          </label>
          <select
            className="select select-bordered w-full"
            value={selectedMunicipalityId || ''}
            onChange={handleMunicipalityChange}
            disabled={loadingMunicipalities}
          >
            <option value="">{t('allMunicipalities')}</option>
            {municipalities.map((municipality) => (
              <option key={municipality.id} value={municipality.id}>
                {municipality.name}
              </option>
            ))}
          </select>
        </div>
      )}

      {/* Индикатор загрузки */}
      {(loadingDistricts || loadingMunicipalities) && (
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-sm"></span>
        </div>
      )}
    </div>
  );
}
