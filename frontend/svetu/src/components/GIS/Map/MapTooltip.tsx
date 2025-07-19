import React from 'react';
import { useTranslations } from 'next-intl';
import Image from 'next/image';
import type { MapMarkerData } from '../types/gis';

interface MapTooltipProps {
  marker: MapMarkerData & { distance?: number };
  visible: boolean;
  position: { x: number; y: number };
}

const MapTooltip: React.FC<MapTooltipProps> = ({
  marker,
  visible,
  position,
}) => {
  const t = useTranslations();

  if (!visible) return null;

  // Форматирование цены
  const formatPrice = (price: number, currency: string = 'RSD') => {
    return new Intl.NumberFormat('sr-RS', {
      style: 'currency',
      currency: currency,
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(price);
  };

  // Форматирование расстояния
  const formatDistance = (distance?: number) => {
    if (!distance) return null;

    if (distance < 1) {
      return `${Math.round(distance * 1000)} ${t('gis.meters')}`;
    }

    return `${distance.toFixed(1)} ${t('gis.kilometers')}`;
  };

  return (
    <div
      className="absolute z-50 pointer-events-none"
      style={{
        left: `${position.x}px`,
        top: `${position.y}px`,
        transform: 'translate(-50%, -100%)',
        marginTop: '-10px',
      }}
    >
      <div className="bg-base-100 rounded-lg shadow-lg p-3 max-w-xs">
        <div className="flex items-start gap-3">
          {/* Миниатюра */}
          {marker.imageUrl && (
            <div className="avatar">
              <div className="w-16 h-16 rounded">
                <Image
                  src={marker.imageUrl}
                  alt={marker.title}
                  fill
                  className="object-cover"
                />
              </div>
            </div>
          )}

          {/* Информация */}
          <div className="flex-1">
            <h3 className="font-semibold text-sm line-clamp-1">
              {marker.title}
            </h3>

            {marker.description && (
              <p className="text-xs text-base-content/70 line-clamp-2 mt-1">
                {marker.description}
              </p>
            )}

            <div className="flex items-center gap-2 mt-2">
              {/* Цена */}
              {marker.metadata?.price && (
                <span className="badge badge-primary badge-sm">
                  {formatPrice(marker.metadata.price, marker.metadata.currency)}
                </span>
              )}

              {/* Расстояние */}
              {marker.distance !== undefined && (
                <span className="badge badge-ghost badge-sm">
                  📍 {formatDistance(marker.distance)}
                </span>
              )}
            </div>

            {/* Категория */}
            {marker.metadata?.category && (
              <div className="text-xs text-base-content/60 mt-1">
                {marker.metadata.category}
              </div>
            )}
          </div>
        </div>

        {/* Стрелка вниз */}
        <div className="absolute -bottom-2 left-1/2 transform -translate-x-1/2">
          <div className="w-0 h-0 border-l-[8px] border-l-transparent border-r-[8px] border-r-transparent border-t-[8px] border-t-base-100"></div>
        </div>
      </div>
    </div>
  );
};

export default MapTooltip;
