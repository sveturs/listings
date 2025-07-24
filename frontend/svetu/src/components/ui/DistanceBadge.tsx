'use client';

import React from 'react';
import { MapPin } from 'lucide-react';

interface DistanceBadgeProps {
  distance: number; // в километрах
  className?: string;
  variant?: 'default' | 'compact' | 'detailed';
  showIcon?: boolean;
}

export const DistanceBadge: React.FC<DistanceBadgeProps> = ({
  distance,
  className = '',
  variant = 'default',
  showIcon = true,
}) => {
  // Форматирование расстояния
  const formatDistance = (km: number): string => {
    if (km < 1) {
      return `${Math.round(km * 1000)} м`;
    }
    if (km < 10) {
      return `${km.toFixed(1)} км`;
    }
    return `${Math.round(km)} км`;
  };

  // Определение цвета в зависимости от расстояния
  const getColorClass = (km: number): string => {
    if (km <= 1) return 'text-success bg-success/10 border-success/20';
    if (km <= 5) return 'text-info bg-info/10 border-info/20';
    if (km <= 15) return 'text-warning bg-warning/10 border-warning/20';
    return 'text-base-content/60 bg-base-200 border-base-300';
  };

  // Получение описания расстояния
  const getDistanceDescription = (km: number): string => {
    if (km <= 0.5) return 'Очень близко';
    if (km <= 1) return 'В шаговой доступности';
    if (km <= 3) return 'Рядом';
    if (km <= 5) return 'Недалеко';
    if (km <= 10) return 'В районе';
    return 'Далеко';
  };

  if (variant === 'compact') {
    return (
      <span
        className={`text-xs font-medium ${getColorClass(distance)} ${className}`}
      >
        {formatDistance(distance)}
      </span>
    );
  }

  if (variant === 'detailed') {
    return (
      <div className={`flex flex-col gap-1 ${className}`}>
        <div className="flex items-center gap-1.5">
          {showIcon && <MapPin className="w-4 h-4" />}
          <span className="font-semibold">{formatDistance(distance)}</span>
        </div>
        <span className="text-xs text-base-content/60">
          {getDistanceDescription(distance)}
        </span>
      </div>
    );
  }

  // Default variant
  return (
    <div
      className={`
        inline-flex items-center gap-1.5 px-2.5 py-1 
        rounded-full border text-sm font-medium
        ${getColorClass(distance)}
        ${className}
      `}
    >
      {showIcon && <MapPin className="w-3.5 h-3.5" />}
      <span>{formatDistance(distance)}</span>
    </div>
  );
};
