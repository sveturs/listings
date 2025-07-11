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
  onSearch?: (query: string) => void;
  className?: string;
  isMobile?: boolean;
  useOpenStreetMap?: boolean;
}

const MapControls: React.FC<MapControlsProps> = ({
  config = {},
  onStyleChange,
  onSearch,
  className = '',
  isMobile = false,
  useOpenStreetMap = false,
}) => {
  const [searchQuery, setSearchQuery] = useState('');
  const [isSearchFocused, setIsSearchFocused] = useState(false);
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

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    if (onSearch && searchQuery.trim()) {
      onSearch(searchQuery.trim());
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
      {/* Поисковая панель - скрыта на мобильном */}
      {!isMobile && (
        <div className={`absolute top-4 left-4 right-4 z-10 ${className}`}>
          <form onSubmit={handleSearch} className="relative">
            <div
              className={`relative transition-all duration-200 ${
                isSearchFocused ? 'ring-2 ring-primary' : ''
              }`}
            >
              <input
                type="text"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onFocus={() => setIsSearchFocused(true)}
                onBlur={() => setIsSearchFocused(false)}
                placeholder="Поиск на карте..."
                className="w-full px-4 py-3 pl-12 pr-20 text-sm bg-white rounded-lg shadow-lg border border-gray-200 focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent"
              />
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg
                  className="h-5 w-5 text-gray-400"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  />
                </svg>
              </div>
              <button
                type="submit"
                className="absolute inset-y-0 right-0 pr-3 flex items-center text-primary hover:text-primary-dark"
              >
                <svg
                  className="h-5 w-5"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13 7l5 5m0 0l-5 5m5-5H6"
                  />
                </svg>
              </button>
            </div>
          </form>
        </div>
      )}

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

      {/* Легенда карты - скрыта на мобильном */}
      {!isMobile && (
        <div className="absolute top-20 right-4 z-10">
          <div className="bg-white rounded-lg shadow-lg p-3 max-w-48">
            <h4 className="font-semibold text-gray-900 text-sm mb-2">
              Легенда
            </h4>
            <div className="space-y-1 text-xs">
              <div className="flex items-center">
                <div className="w-4 h-4 bg-blue-500 rounded-full mr-2"></div>
                <span className="text-gray-700">Жилье</span>
              </div>
              <div className="flex items-center">
                <div className="w-4 h-4 bg-orange-500 rounded-full mr-2"></div>
                <span className="text-gray-700">Пользователи</span>
              </div>
              <div className="flex items-center">
                <div className="w-4 h-4 bg-red-500 rounded-full mr-2"></div>
                <span className="text-gray-700">Достопримечательности</span>
              </div>
            </div>
          </div>
        </div>
      )}
    </>
  );
};

export default MapControls;
