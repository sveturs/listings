'use client';

import React from 'react';
import { Marker } from 'react-map-gl';

interface PriceClusterProps {
  /** Координаты центра кластера */
  coordinates: [number, number];
  /** Количество объявлений в кластере */
  count: number;
  /** Минимальная цена в кластере */
  minPrice?: number;
  /** Максимальная цена в кластере */
  maxPrice?: number;
  /** Количество объявлений с ценами в кластере */
  listingsCount?: number;
  /** Обработчик клика */
  onClick?: () => void;
  className?: string;
}

const PriceCluster: React.FC<PriceClusterProps> = ({
  coordinates,
  count,
  minPrice,
  maxPrice,
  listingsCount = 0,
  onClick,
  className = '',
}) => {
  // Размер кластера зависит от количества элементов
  const getClusterSize = (count: number): number => {
    if (count < 10) return 48;
    if (count < 50) return 58;
    return 68;
  };

  // Цвет кластера зависит от количества элементов
  const getClusterColor = (count: number): string => {
    if (count < 10) return 'bg-blue-500';
    if (count < 50) return 'bg-green-500';
    return 'bg-red-500';
  };

  const size = getClusterSize(count);
  const colorClass = getClusterColor(count);

  // Форматируем диапазон цен
  const formatPriceRange = (): string => {
    if (!minPrice || !maxPrice || listingsCount === 0) return '';

    if (minPrice === maxPrice) {
      return `${minPrice.toLocaleString()} RSD`;
    }

    return `${minPrice.toLocaleString()}-${maxPrice.toLocaleString()} RSD`;
  };

  const priceRange = formatPriceRange();

  return (
    <Marker
      longitude={coordinates[0]}
      latitude={coordinates[1]}
      anchor="center"
      onClick={onClick}
    >
      <div className={`relative ${className}`}>
        {/* Пульсирующий фон */}
        <div
          className={`absolute inset-0 ${colorClass} rounded-full opacity-20 animate-pulse`}
          style={{ width: `${size}px`, height: `${size}px` }}
        />

        {/* Основной кластер */}
        <div
          className={`
            relative flex items-center justify-center
            ${colorClass} rounded-full shadow-lg border-2 border-white
            cursor-pointer transition-all duration-200 hover:scale-110
          `}
          style={{ width: `${size}px`, height: `${size}px` }}
        >
          <span className="text-white font-bold text-sm md:text-base">
            {count > 99 ? '99+' : count}
          </span>
        </div>

        {/* Диапазон цен под кластером */}
        {priceRange && (
          <div
            className="
              absolute top-full left-1/2 transform -translate-x-1/2 mt-2
              px-3 py-1 rounded-lg text-xs font-bold
              border border-gray-200 shadow-md
              bg-white text-gray-800
              whitespace-nowrap z-10
            "
          >
            {priceRange}
          </div>
        )}

        {/* Количество объявлений с ценами (если отличается от общего количества) */}
        {listingsCount > 0 && listingsCount !== count && (
          <div
            className="
              absolute -top-1 -right-1
              px-1 py-0.5 rounded-full text-xs font-bold
              bg-orange-500 text-white
              border border-white shadow-sm
            "
          >
            {listingsCount}
          </div>
        )}
      </div>
    </Marker>
  );
};

export default PriceCluster;
