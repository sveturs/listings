import React, { useState } from 'react';
import {
  NavigationControl,
  FullscreenControl,
  GeolocateControl,
} from 'react-map-gl';
import { MapControlsConfig } from '../types/gis';
import { useGeolocation } from '../hooks/useGeolocation';

interface MapControlsProps {
  config?: MapControlsConfig;
  onStyleChange?: (style: string) => void;
  className?: string;
  isMobile?: boolean;
  useOpenStreetMap?: boolean;
}

const MapControls: React.FC<MapControlsProps> = ({
  config = {},
  onStyleChange,
  className: _className = '',
  isMobile = false,
  useOpenStreetMap = false,
}) => {
  const [currentStyle, setCurrentStyle] = useState('streets');

  const { getCurrentPosition, loading: geoLoading } = useGeolocation();

  const {
    showZoom = true,
    showCompass = true,
    showFullscreen = true,
    showGeolocate = true,
    showNavigation = true,
    position = 'top-right',
  } = config;

  const mapStyles = [
    {
      id: 'streets',
      name: 'Streets',
      url: 'mapbox://styles/mapbox/streets-v12',
    },
    {
      id: 'satellite',
      name: 'Satellite',
      url: 'mapbox://styles/mapbox/satellite-streets-v12',
    },
    {
      id: 'outdoors',
      name: 'Outdoors',
      url: 'mapbox://styles/mapbox/outdoors-v12',
    },
    { id: 'light', name: 'Light', url: 'mapbox://styles/mapbox/light-v11' },
    { id: 'dark', name: 'Dark', url: 'mapbox://styles/mapbox/dark-v11' },
  ];

  const handleStyleChange = (styleId: string) => {
    const style = mapStyles.find((s) => s.id === styleId);
    if (style && onStyleChange) {
      setCurrentStyle(styleId);
      onStyleChange(style.url);
    }
  };

  const handleGeolocation = async () => {
    try {
      await getCurrentPosition();
    } catch (error) {
      console.error('Geolocation error:', error);
    }
  };

  const controlsPosition = {
    'top-left': 'top-0 left-0',
    'top-right': 'top-0 right-0',
    'bottom-left': 'bottom-0 left-0',
    'bottom-right': 'bottom-0 right-0',
  }[position];

  return (
    <>
      {/* Стандартные контролы Mapbox */}
      <div
        className={`absolute ${controlsPosition} ${
          isMobile ? 'm-2' : 'm-4'
        } z-10 flex flex-col space-y-2`}
      >
        {showNavigation && (
          <NavigationControl
            showCompass={showCompass}
            showZoom={showZoom}
            visualizePitch={true}
          />
        )}

        {showGeolocate && (
          <GeolocateControl
            positionOptions={{ enableHighAccuracy: true }}
            trackUserLocation={true}
            showAccuracyCircle={true}
            showUserLocation={true}
          />
        )}

        {showFullscreen && <FullscreenControl />}
      </div>

      {/* Панель переключения стилей - адаптивная (скрыта для OpenStreetMap) */}
      {!useOpenStreetMap && (
        <div
          className={`absolute ${isMobile ? 'bottom-20' : 'bottom-4'} left-4 z-10`}
        >
          <div className="bg-white rounded-lg shadow-lg p-2">
            <div
              className={`grid ${isMobile ? 'grid-cols-2' : 'grid-cols-3'} gap-1`}
            >
              {mapStyles.slice(0, isMobile ? 4 : 5).map((style) => (
                <button
                  key={style.id}
                  onClick={() => handleStyleChange(style.id)}
                  className={`px-2 py-1 text-xs rounded-md transition-colors ${
                    currentStyle === style.id
                      ? 'bg-primary text-white'
                      : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                  }`}
                >
                  {style.name}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Кнопка геолокации (дополнительная) - скрыта на мобильном */}
      {!isMobile && (
        <div className="absolute bottom-4 right-4 z-10">
          <button
            onClick={handleGeolocation}
            disabled={geoLoading}
            className={`bg-white rounded-lg shadow-lg p-3 text-gray-700 hover:bg-gray-50 transition-colors ${
              geoLoading ? 'opacity-50 cursor-not-allowed' : ''
            }`}
            title="Показать мое местоположение"
          >
            {geoLoading ? (
              <svg
                className="w-5 h-5 animate-spin"
                fill="none"
                viewBox="0 0 24 24"
              >
                <circle
                  className="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  strokeWidth="4"
                ></circle>
                <path
                  className="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
            ) : (
              <svg
                className="w-5 h-5"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                />
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                />
              </svg>
            )}
          </button>
        </div>
      )}
    </>
  );
};

export default MapControls;
