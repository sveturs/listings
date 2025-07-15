'use client';

import React, { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useVisibleCitiesContext } from '../GIS/contexts/VisibleCitiesContext';
import type { components as _components } from '@/types/generated/api';

// Временные интерфейсы до исправления API типов
interface District {
  id: string;
  name: string;
  geometry?: any;
  boundary?: {
    coordinates: number[][][];
  };
  bounds?: [number, number, number, number];
  population?: number;
  area?: number;
  area_km2?: number;
}

interface Municipality {
  id: string;
  name: string;
  districts?: District[];
}

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

  // Используем контекст для получения районов текущего города
  const { availableDistricts, closestCity, loading: citiesLoading } = useVisibleCitiesContext();

  // Логирование для отладки
  console.log('🏗️ DistrictSelector render:', {
    availableDistricts: availableDistricts.length,
    closestCity: closestCity?.city.name,
    citiesLoading
  });

  const [municipalities, setMunicipalities] = useState<Municipality[]>([]);
  const [loadingMunicipalities, setLoadingMunicipalities] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Сбрасываем выбранный район, если изменился город
  useEffect(() => {
    if (selectedDistrictId) {
      // Проверяем, есть ли выбранный район в списке доступных районов
      const stillAvailable = availableDistricts.some(
        (district) => district.id === selectedDistrictId
      );

      if (!stillAvailable) {
        onDistrictChange('');
        onMunicipalityChange('');
      }
    }
  }, [availableDistricts, selectedDistrictId, onDistrictChange, onMunicipalityChange]);

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
          disabled={citiesLoading}
        >
          <option value="">{t('allDistricts')}</option>
          {closestCity && (
            <optgroup label={`${t('districtsIn')} ${closestCity.city.name}`}>
              {availableDistricts.map((district) => (
                <option key={district.id} value={district.id}>
                  {district.name}
                </option>
              ))}
            </optgroup>
          )}
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
      {(citiesLoading || loadingMunicipalities) && (
        <div className="flex justify-center">
          <span className="loading loading-spinner loading-sm"></span>
        </div>
      )}
    </div>
  );
}
