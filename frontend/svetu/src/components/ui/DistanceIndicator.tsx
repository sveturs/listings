'use client';

import React from 'react';
import { MapPin } from 'lucide-react';
import { Tooltip } from './Tooltip';

interface DistanceIndicatorProps {
  distance: number; // в километрах
  className?: string;
  showTooltip?: boolean;
  size?: 'sm' | 'md' | 'lg';
}

export const DistanceIndicator: React.FC<DistanceIndicatorProps> = ({
  distance,
  className = '',
  showTooltip = true,
  size = 'md',
}) => {
  // Размеры в зависимости от size
  const sizeClasses = {
    sm: 'text-xs gap-1',
    md: 'text-sm gap-1.5',
    lg: 'text-base gap-2',
  };

  const iconSizes = {
    sm: 'w-3 h-3',
    md: 'w-4 h-4',
    lg: 'w-5 h-5',
  };

  // Форматирование расстояния
  const formatDistance = (km: number): string => {
    if (km < 1) {
      return `${Math.round(km * 1000)}м`;
    }
    return `${km.toFixed(1)}км`;
  };

  // Цвет в зависимости от расстояния
  const getColor = (km: number): string => {
    if (km <= 1) return 'text-success';
    if (km <= 5) return 'text-info';
    if (km <= 15) return 'text-warning';
    return 'text-base-content/60';
  };

  const colorClass = getColor(distance);

  const content = (
    <div
      className={`inline-flex items-center ${sizeClasses[size]} ${colorClass} ${className}`}
    >
      <MapPin className={iconSizes[size]} />
      <span className="font-medium">{formatDistance(distance)}</span>
    </div>
  );

  if (!showTooltip) {
    return content;
  }

  // Подробная информация для tooltip
  const getTooltipContent = () => {
    const walkingTime = Math.round((distance / 5) * 60); // 5 км/ч скорость ходьбы

    return (
      <div className="space-y-2 p-2">
        <div className="font-medium">Расстояние от вас</div>
        <div className="text-sm space-y-1">
          <div>
            {distance < 1
              ? `${Math.round(distance * 1000)} метров`
              : `${distance.toFixed(1)} километров`}
          </div>
          <div className="text-base-content/70">
            Пешком: ~
            {walkingTime < 60
              ? `${walkingTime} мин`
              : `${Math.round(walkingTime / 60)} ч`}
          </div>
          {distance <= 1 && (
            <div className="text-success">В шаговой доступности!</div>
          )}
        </div>
      </div>
    );
  };

  return <Tooltip content={getTooltipContent()}>{content}</Tooltip>;
};
