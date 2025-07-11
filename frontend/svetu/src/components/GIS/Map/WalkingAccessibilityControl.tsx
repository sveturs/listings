import React, { useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';

interface WalkingAccessibilityControlProps {
  mode: 'radius' | 'walking';
  onModeChange: (mode: 'radius' | 'walking') => void;
  walkingTime: number; // в минутах
  onWalkingTimeChange: (time: number) => void;
  searchRadius: number; // в метрах
  onRadiusChange: (radius: number) => void;
}

const WalkingAccessibilityControl: React.FC<
  WalkingAccessibilityControlProps
> = ({
  mode,
  onModeChange,
  walkingTime,
  onWalkingTimeChange,
  searchRadius,
  onRadiusChange,
}) => {
  const t = useTranslations();
  const [isOpen, setIsOpen] = useState(false);

  const walkingTimes = [5, 10, 15, 20, 30]; // минуты
  const radiusOptions = [1000, 2000, 5000, 10000, 20000]; // метры

  const handleModeChange = useCallback(
    (newMode: 'radius' | 'walking') => {
      onModeChange(newMode);

      // При переключении в режим пешей доступности, конвертируем радиус в время
      if (newMode === 'walking' && mode === 'radius') {
        // Предполагаем среднюю скорость ходьбы 5 км/ч (83.33 м/мин)
        const estimatedTime = Math.round(searchRadius / 83.33);
        const closestTime = walkingTimes.reduce((prev, curr) =>
          Math.abs(curr - estimatedTime) < Math.abs(prev - estimatedTime)
            ? curr
            : prev
        );
        onWalkingTimeChange(closestTime);
      } else if (newMode === 'radius' && mode === 'walking') {
        // При переключении обратно, конвертируем время в радиус
        const estimatedRadius = walkingTime * 83.33;
        const closestRadius = radiusOptions.reduce((prev, curr) =>
          Math.abs(curr - estimatedRadius) < Math.abs(prev - estimatedRadius)
            ? curr
            : prev
        );
        onRadiusChange(closestRadius);
      }
    },
    [
      mode,
      searchRadius,
      walkingTime,
      walkingTimes,
      radiusOptions,
      onModeChange,
      onWalkingTimeChange,
      onRadiusChange,
    ]
  );

  const formatDistance = (meters: number) => {
    if (meters < 1000) {
      return `${meters} ${t('gis.meters')}`;
    }
    return `${(meters / 1000).toFixed(1)} ${t('gis.kilometers')}`;
  };

  return (
    <div className="relative">
      {/* Кнопка открытия */}
      <button
        className="btn btn-sm btn-ghost"
        onClick={() => setIsOpen(!isOpen)}
        title={t('gis.accessibilitySettings')}
      >
        {mode === 'walking' ? '🚶' : '📍'}
        <span className="ml-1">
          {mode === 'walking'
            ? `${walkingTime} ${t('gis.minutes')}`
            : formatDistance(searchRadius)}
        </span>
      </button>

      {/* Выпадающее меню */}
      {isOpen && (
        <div className="absolute right-0 mt-2 w-72 bg-base-100 rounded-lg shadow-lg p-4 z-50">
          <h3 className="font-semibold mb-3">{t('gis.searchAreaSettings')}</h3>

          {/* Переключатель режима */}
          <div className="tabs tabs-boxed mb-3">
            <button
              className={`tab ${mode === 'radius' ? 'tab-active' : ''}`}
              onClick={() => handleModeChange('radius')}
            >
              📍 {t('gis.byRadius')}
            </button>
            <button
              className={`tab ${mode === 'walking' ? 'tab-active' : ''}`}
              onClick={() => handleModeChange('walking')}
            >
              🚶 {t('gis.byWalkingTime')}
            </button>
          </div>

          {/* Контролы для режима радиуса */}
          {mode === 'radius' && (
            <div>
              <label className="label">
                <span className="label-text">{t('gis.searchRadius')}</span>
                <span className="label-text-alt">
                  {formatDistance(searchRadius)}
                </span>
              </label>
              <div className="flex flex-wrap gap-2">
                {radiusOptions.map((radius) => (
                  <button
                    key={radius}
                    className={`btn btn-xs ${searchRadius === radius ? 'btn-primary' : 'btn-ghost'}`}
                    onClick={() => onRadiusChange(radius)}
                  >
                    {formatDistance(radius)}
                  </button>
                ))}
              </div>
              <input
                type="range"
                min="500"
                max="20000"
                step="500"
                value={searchRadius}
                onChange={(e) => onRadiusChange(Number(e.target.value))}
                className="range range-primary range-xs mt-3"
              />
            </div>
          )}

          {/* Контролы для режима пешей доступности */}
          {mode === 'walking' && (
            <div>
              <label className="label">
                <span className="label-text">{t('gis.walkingTime')}</span>
                <span className="label-text-alt">
                  {walkingTime} {t('gis.minutes')}
                </span>
              </label>
              <div className="flex flex-wrap gap-2">
                {walkingTimes.map((time) => (
                  <button
                    key={time}
                    className={`btn btn-xs ${walkingTime === time ? 'btn-primary' : 'btn-ghost'}`}
                    onClick={() => onWalkingTimeChange(time)}
                  >
                    {time} {t('gis.min')}
                  </button>
                ))}
              </div>
              <div className="text-xs text-base-content/60 mt-3">
                <p>{t('gis.walkingSpeedNote')}</p>
                <p className="mt-1">
                  ≈ {formatDistance(walkingTime * 83.33)}{' '}
                  {t('gis.atAverageSpeed')}
                </p>
              </div>
            </div>
          )}

          {/* Кнопка закрытия */}
          <button
            className="btn btn-sm btn-block mt-4"
            onClick={() => setIsOpen(false)}
          >
            {t('common.close')}
          </button>
        </div>
      )}
    </div>
  );
};

export default WalkingAccessibilityControl;
