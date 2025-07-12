'use client';

import React, { useState, useCallback } from 'react';

export interface LocationPrivacyLevel {
  id: 'exact' | 'street' | 'district' | 'city';
  label: string;
  description: string;
  radiusMeters: number;
  icon: string;
  example: string;
}

export interface LocationPrivacySettingsProps {
  selectedLevel: LocationPrivacyLevel['id'];
  onLevelChange: (level: LocationPrivacyLevel['id']) => void;
  location?: { lat: number; lng: number };
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
    example: 'Улица Княза Милоша 15, Белград',
  },
  {
    id: 'street',
    label: 'Улица',
    description: 'Местоположение размыто в пределах ±150 метров',
    radiusMeters: 150,
    icon: '🏠',
    example: 'Улица Княза Милоша, Белград',
  },
  {
    id: 'district',
    label: 'Район',
    description: 'Местоположение размыто в пределах ±750 метров',
    radiusMeters: 750,
    icon: '🏘️',
    example: 'Савски венец, Белград',
  },
  {
    id: 'city',
    label: 'Только город',
    description: 'Показывается только город, размытие ±5 км',
    radiusMeters: 5000,
    icon: '🏙️',
    example: 'Белград, Сербия',
  },
];

export default function LocationPrivacySettings({
  selectedLevel,
  onLevelChange,
  location,
  showPreview = true,
  className = '',
}: LocationPrivacySettingsProps) {
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

  // Иконка для уровня безопасности
  const getSecurityIcon = (levelId: LocationPrivacyLevel['id']) => {
    switch (levelId) {
      case 'exact':
        return '🔓'; // Открытый замок
      case 'street':
        return '🔐'; // Частично закрытый
      case 'district':
        return '🔒'; // Закрытый замок
      case 'city':
        return '🔐'; // Самый безопасный
      default:
        return '🔒';
    }
  };

  return (
    <div className={`w-full ${className}`}>
      {/* Заголовок */}
      <div className="mb-6">
        <h3 className="text-lg font-semibold mb-2">
          Настройки приватности местоположения
        </h3>
        <p className="text-sm text-base-content/70">
          Выберите, как точно показывать местоположение вашего объявления другим
          пользователям
        </p>
      </div>

      {/* Карточки уровней приватности */}
      <div className="grid gap-4 mb-6">
        {PRIVACY_LEVELS.map((level) => {
          const isSelected = selectedLevel === level.id;
          const isHovered = hoveredLevel === level.id;

          return (
            <div
              key={level.id}
              className={`
                relative cursor-pointer p-4 border-2 rounded-lg transition-all duration-200
                ${getLevelColor(level.id, isSelected)}
                ${isSelected ? 'ring-2 ring-offset-2 ring-base-300' : ''}
              `}
              onClick={() => handleLevelSelect(level.id)}
              onMouseEnter={() => setHoveredLevel(level.id)}
              onMouseLeave={() => setHoveredLevel(null)}
            >
              {/* Радио кнопка и заголовок */}
              <div className="flex items-start justify-between mb-2">
                <div className="flex items-center">
                  <input
                    type="radio"
                    name="privacy-level"
                    value={level.id}
                    checked={isSelected}
                    onChange={() => handleLevelSelect(level.id)}
                    className="radio radio-primary mr-3"
                  />

                  <div className="flex items-center">
                    <span className="text-2xl mr-2">{level.icon}</span>
                    <div>
                      <h4 className="font-medium text-base">{level.label}</h4>
                      {level.radiusMeters > 0 && (
                        <span className="text-xs text-base-content/50">
                          ±{level.radiusMeters}м
                        </span>
                      )}
                    </div>
                  </div>
                </div>

                {/* Иконка безопасности */}
                <div className="flex items-center">
                  <span className="text-xl mr-1">
                    {getSecurityIcon(level.id)}
                  </span>
                  <div className="text-xs text-base-content/50 text-right">
                    {level.id === 'exact' && 'Низкая\nприватность'}
                    {level.id === 'street' && 'Средняя\nприватность'}
                    {level.id === 'district' && 'Высокая\nприватность'}
                    {level.id === 'city' && 'Максимальная\nприватность'}
                  </div>
                </div>
              </div>

              {/* Описание */}
              <p className="text-sm text-base-content/70 mb-2 ml-8">
                {level.description}
              </p>

              {/* Пример */}
              <div className="ml-8">
                <span className="text-xs font-medium text-base-content/50">
                  Пример отображения:
                </span>
                <p className="text-xs text-base-content/60 italic mt-1">
                  "{level.example}"
                </p>
              </div>

              {/* Индикатор выбора */}
              {isSelected && (
                <div className="absolute top-2 right-2">
                  <div className="w-6 h-6 bg-primary rounded-full flex items-center justify-center">
                    <svg
                      className="w-3 h-3 text-primary-content"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                        clipRule="evenodd"
                      />
                    </svg>
                  </div>
                </div>
              )}
            </div>
          );
        })}
      </div>

      {/* Превью выбранного уровня */}
      {showPreview && previewLevel && location && (
        <div className="p-4 bg-base-200 rounded-lg">
          <h4 className="font-medium mb-3 flex items-center">
            <span className="text-lg mr-2">👁️</span>
            Превью: как видят другие пользователи
          </h4>

          <div className="space-y-3">
            {/* Визуализация зоны */}
            <div className="flex items-center justify-between p-3 bg-base-100 rounded border">
              <div>
                <div className="font-medium text-sm">{previewLevel.label}</div>
                <div className="text-xs text-base-content/70">
                  {previewLevel.radiusMeters === 0
                    ? 'Точное местоположение'
                    : `Размытие в радиусе ${previewLevel.radiusMeters}м`}
                </div>
              </div>

              <div className="text-right">
                <div className="text-xs text-base-content/50">Координаты:</div>
                <div className="text-xs font-mono">
                  {previewLevel.radiusMeters === 0
                    ? `${location.lat.toFixed(6)}, ${location.lng.toFixed(6)}`
                    : '●●●.●●●●●●, ●●●.●●●●●●'}
                </div>
              </div>
            </div>

            {/* Предупреждения */}
            {previewLevel.id === 'exact' && (
              <div className="p-3 bg-warning/10 border border-warning/20 rounded-lg">
                <div className="flex items-start">
                  <span className="text-warning text-lg mr-2">⚠️</span>
                  <div className="text-sm">
                    <div className="font-medium text-warning-content mb-1">
                      Внимание!
                    </div>
                    <p className="text-warning-content/80">
                      Ваш точный адрес будет виден всем пользователям.
                      Рекомендуется для бизнеса, но не для личных объявлений.
                    </p>
                  </div>
                </div>
              </div>
            )}

            {previewLevel.id === 'city' && (
              <div className="p-3 bg-info/10 border border-info/20 rounded-lg">
                <div className="flex items-start">
                  <span className="text-info text-lg mr-2">💡</span>
                  <div className="text-sm">
                    <div className="font-medium text-info-content mb-1">
                      Подсказка
                    </div>
                    <p className="text-info-content/80">
                      Максимальная приватность. Покупатели смогут связаться с
                      вами для уточнения точного местоположения.
                    </p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Информационная панель */}
      <div className="mt-6 p-4 bg-base-100 border border-base-300 rounded-lg">
        <h5 className="font-medium mb-2 flex items-center">
          <span className="mr-2">🛡️</span>
          Рекомендации по безопасности
        </h5>

        <ul className="text-sm text-base-content/70 space-y-1">
          <li>
            • <strong>Для бизнеса:</strong> используйте "Точный адрес" для
            магазинов и офисов
          </li>
          <li>
            • <strong>Для дома:</strong> рекомендуется "Улица" или "Район" для
            защиты приватности
          </li>
          <li>
            • <strong>Для встреч:</strong> "Район" позволяет договориться о
            точном месте отдельно
          </li>
          <li>
            • Вы всегда можете изменить настройки приватности в любое время
          </li>
        </ul>
      </div>
    </div>
  );
}
