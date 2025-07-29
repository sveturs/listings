'use client';

import React, { useState, useCallback, useMemo } from 'react';
import { formatAddressWithPrivacy } from '@/utils/addressPrivacy';

export interface LocationPrivacyLevel {
  id: 'exact' | 'street' | 'district' | 'city';
  label: string;
  description: string;
  radiusMeters: number;
  icon: string;
  example?: string;
}

export interface LocationPrivacySettingsWithAddressProps {
  selectedLevel: LocationPrivacyLevel['id'];
  onLevelChange: (level: LocationPrivacyLevel['id']) => void;
  location?: { lat: number; lng: number };
  fullAddress?: string;
  showPreview?: boolean;
  className?: string;
}

const PRIVACY_LEVELS: LocationPrivacyLevel[] = [
  {
    id: 'exact',
    label: 'Точный адрес',
    description: 'Показывается точное местоположение вашего объявления',
    radiusMeters: 0,
    icon: '🎯',
  },
  {
    id: 'street',
    label: 'Улица',
    description: 'Местоположение размыто в пределах ±150 метров',
    radiusMeters: 150,
    icon: '🏠',
  },
  {
    id: 'district',
    label: 'Район',
    description: 'Местоположение размыто в пределах ±750 метров',
    radiusMeters: 750,
    icon: '🏘️',
  },
  {
    id: 'city',
    label: 'Только город',
    description: 'Показывается только город, размытие ±5 км',
    radiusMeters: 5000,
    icon: '🏙️',
  },
];

export default function LocationPrivacySettingsWithAddress({
  selectedLevel,
  onLevelChange,
  location,
  fullAddress,
  showPreview = true,
  className = '',
}: LocationPrivacySettingsWithAddressProps) {
  const [hoveredLevel, setHoveredLevel] = useState<
    LocationPrivacyLevel['id'] | null
  >(null);

  const selectedLevelData = PRIVACY_LEVELS.find(
    (level) => level.id === selectedLevel
  );
  const previewLevel = hoveredLevel
    ? PRIVACY_LEVELS.find((level) => level.id === hoveredLevel)
    : selectedLevelData;

  const handleLevelSelect = useCallback(
    (levelId: LocationPrivacyLevel['id']) => {
      onLevelChange(levelId);
    },
    [onLevelChange]
  );

  // Форматируем адрес в зависимости от уровня приватности
  const getFormattedAddress = useCallback((levelId: LocationPrivacyLevel['id']) => {
    if (!fullAddress) return null;
    
    const options = {
      showHouseNumber: levelId === 'exact',
      showStreet: levelId === 'exact' || levelId === 'street',
      showCity: true,
      showRegion: levelId !== 'city',
      showCountry: levelId === 'city'
    };
    
    return formatAddressWithPrivacy(fullAddress, options);
  }, [fullAddress]);

  // Примеры адресов для каждого уровня
  const addressExamples = useMemo(() => {
    const examples: Record<LocationPrivacyLevel['id'], string> = {
      exact: fullAddress || 'Улица Княза Милоша 15, Белград',
      street: getFormattedAddress('street') || 'Улица Княза Милоша, Белград',
      district: getFormattedAddress('district') || 'Савски венац, Белград',
      city: getFormattedAddress('city') || 'Белград, Сербия',
    };
    return examples;
  }, [fullAddress, getFormattedAddress]);

  // Получение цвета для уровня приватности
  const getLevelColor = (
    levelId: LocationPrivacyLevel['id'],
    isSelected: boolean
  ) => {
    const baseColors = {
      exact: isSelected
        ? 'border-error bg-error/10'
        : 'border-error/30 hover:border-error hover:bg-error/5',
      street: isSelected
        ? 'border-warning bg-warning/10'
        : 'border-warning/30 hover:border-warning hover:bg-warning/5',
      district: isSelected
        ? 'border-info bg-info/10'
        : 'border-info/30 hover:border-info hover:bg-info/5',
      city: isSelected
        ? 'border-success bg-success/10'
        : 'border-success/30 hover:border-success hover:bg-success/5',
    };
    return baseColors[levelId];
  };

  return (
    <div className={`space-y-4 ${className}`}>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
        {PRIVACY_LEVELS.map((level) => {
          const isSelected = selectedLevel === level.id;
          return (
            <button
              key={level.id}
              type="button"
              onClick={() => handleLevelSelect(level.id)}
              onMouseEnter={() => setHoveredLevel(level.id)}
              onMouseLeave={() => setHoveredLevel(null)}
              className={`
                relative p-4 rounded-lg border-2 transition-all duration-200
                ${getLevelColor(level.id, isSelected)}
                ${isSelected ? 'shadow-lg scale-[1.02]' : 'cursor-pointer'}
              `}
            >
              <div className="flex items-start gap-3">
                <span className="text-2xl">{level.icon}</span>
                <div className="flex-1 text-left">
                  <h4 className="font-semibold text-base-content">
                    {level.label}
                  </h4>
                  <p className="text-sm text-base-content/70 mt-1">
                    {level.description}
                  </p>
                  {/* Показываем примеры */}
                  <p className="text-xs text-base-content/50 mt-2 italic">
                    Пример: {addressExamples[level.id]}
                  </p>
                </div>
                {isSelected && (
                  <svg
                    className="w-5 h-5 text-success absolute top-2 right-2"
                    fill="currentColor"
                    viewBox="0 0 20 20"
                  >
                    <path
                      fillRule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                      clipRule="evenodd"
                    />
                  </svg>
                )}
              </div>
            </button>
          );
        })}
      </div>

      {showPreview && previewLevel && (
        <div className="mt-6 p-4 bg-base-200 rounded-lg">
          <h4 className="font-semibold mb-2 flex items-center gap-2">
            <span className="text-xl">{previewLevel.icon}</span>
            Как будет отображаться ваш адрес:
          </h4>
          <p className="text-lg font-medium text-base-content">
            {addressExamples[previewLevel.id]}
          </p>
          {previewLevel.radiusMeters > 0 && (
            <p className="text-sm text-base-content/60 mt-2">
              ± {previewLevel.radiusMeters < 1000
                ? `${previewLevel.radiusMeters} м`
                : `${previewLevel.radiusMeters / 1000} км`}
            </p>
          )}
        </div>
      )}
    </div>
  );
}