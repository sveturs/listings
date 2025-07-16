'use client';

import React from 'react';
import { Marker } from 'react-map-gl';
import { MapMarkerData } from '../types/gis';

interface PriceMarkerProps {
  marker: MapMarkerData;
  onClick?: (marker: MapMarkerData) => void;
  selected?: boolean;
  className?: string;
}

const PriceMarker: React.FC<PriceMarkerProps> = ({
  marker,
  onClick,
  selected = false,
  className = '',
}) => {
  const handleClick = (e: any) => {
    e.originalEvent?.stopPropagation();
    if (onClick) {
      onClick(marker);
    }
  };

  const price = marker.data?.price || marker.metadata?.price;
  const icon = marker.data?.icon || marker.metadata?.icon || '🏠';

  return (
    <Marker
      longitude={marker.position[0]}
      latitude={marker.position[1]}
      anchor="center"
      onClick={handleClick}
    >
      <div className={`relative ${className}`}>
        {/* Основной маркер */}
        <div
          className={`
            flex items-center justify-center
            w-10 h-10 rounded-full
            border-2 border-white
            shadow-lg cursor-pointer
            transition-all duration-200
            ${
              selected
                ? 'bg-green-500 border-green-600 scale-110 shadow-xl'
                : 'bg-blue-500 border-blue-600 hover:scale-105'
            }
          `}
        >
          <span className="text-lg">{icon}</span>
        </div>

        {/* Цена под маркером */}
        {price && (
          <div
            className={`
              absolute top-full left-1/2 transform -translate-x-1/2 mt-1
              px-2 py-1 rounded-lg text-xs font-bold
              border border-gray-200 shadow-md
              bg-white text-gray-800
              whitespace-nowrap z-10
              ${selected ? 'border-green-300 bg-green-50' : ''}
            `}
          >
            {typeof price === 'number'
              ? `${price.toLocaleString()} RSD`
              : `${price} RSD`}
          </div>
        )}

        {/* Пульсирующий эффект для выбранного маркера */}
        {selected && (
          <div
            className="absolute inset-0 rounded-full animate-ping bg-green-400 opacity-30"
            style={{
              width: '40px',
              height: '40px',
            }}
          />
        )}
      </div>
    </Marker>
  );
};

export default PriceMarker;
