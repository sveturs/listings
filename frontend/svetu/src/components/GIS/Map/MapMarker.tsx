import React, { useMemo } from 'react';
import { Marker } from 'react-map-gl';
import { MapMarkerData } from '../types/gis';

interface MapMarkerProps {
  marker: MapMarkerData;
  onClick?: (marker: MapMarkerData) => void;
  selected?: boolean;
  className?: string;
}

const MapMarker: React.FC<MapMarkerProps> = ({
  marker,
  onClick,
  selected = false,
  className = '',
}) => {
  const markerStyle = useMemo(() => {
    const baseStyle = {
      cursor: 'pointer',
      transition: 'all 0.2s ease',
    };

    switch (marker.type) {
      case 'listing':
        return {
          ...baseStyle,
          backgroundColor: selected ? '#10b981' : '#3b82f6',
          color: 'white',
          border: selected ? '2px solid #059669' : '2px solid #1d4ed8',
          borderRadius: '50%',
          width: selected ? '32px' : '24px',
          height: selected ? '32px' : '24px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: selected ? '16px' : '12px',
          fontWeight: '600',
          boxShadow: selected
            ? '0 4px 12px rgba(0, 0, 0, 0.3)'
            : '0 2px 6px rgba(0, 0, 0, 0.2)',
        };
      case 'user':
        return {
          ...baseStyle,
          backgroundColor: selected ? '#f59e0b' : '#f97316',
          color: 'white',
          border: selected ? '2px solid #d97706' : '2px solid #ea580c',
          borderRadius: '50%',
          width: selected ? '28px' : '20px',
          height: selected ? '28px' : '20px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: selected ? '14px' : '10px',
          fontWeight: '600',
          boxShadow: selected
            ? '0 4px 12px rgba(0, 0, 0, 0.3)'
            : '0 2px 6px rgba(0, 0, 0, 0.2)',
        };
      case 'poi':
        return {
          ...baseStyle,
          backgroundColor: selected ? '#dc2626' : '#ef4444',
          color: 'white',
          border: selected ? '2px solid #b91c1c' : '2px solid #dc2626',
          borderRadius: '50%',
          width: selected ? '24px' : '18px',
          height: selected ? '24px' : '18px',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          fontSize: selected ? '12px' : '8px',
          fontWeight: '600',
          boxShadow: selected
            ? '0 4px 12px rgba(0, 0, 0, 0.3)'
            : '0 2px 6px rgba(0, 0, 0, 0.2)',
        };
      default:
        return baseStyle;
    }
  }, [marker.type, selected]);

  const markerIcon = useMemo(() => {
    switch (marker.type) {
      case 'listing':
        return '🏠';
      case 'user':
        return '👤';
      case 'poi':
        return '📍';
      default:
        return '📍';
    }
  }, [marker.type]);

  const handleClick = (e: any) => {
    e.originalEvent?.stopPropagation();
    if (onClick) {
      onClick(marker);
    }
  };

  return (
    <Marker
      longitude={marker.position[0]}
      latitude={marker.position[1]}
      anchor="center"
      onClick={handleClick}
    >
      <div
        className={`map-marker ${className}`}
        style={markerStyle}
        title={marker.title}
      >
        {marker.type === 'listing' && marker.data?.price ? (
          <span className="text-xs font-bold">{marker.data.price}€</span>
        ) : (
          <span>{markerIcon}</span>
        )}
      </div>

      {/* Пульсирующий эффект для выбранного маркера */}
      {selected && (
        <div
          className="absolute inset-0 rounded-full animate-ping"
          style={{
            backgroundColor: (markerStyle as any).backgroundColor || '#3b82f6',
            opacity: 0.5,
            width: '32px',
            height: '32px',
            top: '50%',
            left: '50%',
            transform: 'translate(-50%, -50%)',
          }}
        />
      )}
    </Marker>
  );
};

export default MapMarker;
