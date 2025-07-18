'use client';

import { useState, useCallback, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import dynamic from 'next/dynamic';
import SmartAddressInput from '@/components/GIS/SmartAddressInput';
import { AddressGeocodingResult } from '@/hooks/useAddressGeocoding';
import { MagnifyingGlassIcon, MapPinIcon } from '@heroicons/react/24/outline';

// Динамически импортируем карту для избежания SSR проблем
const Map = dynamic(() => import('react-map-gl'), { ssr: false });
const Marker = dynamic(() => import('react-map-gl').then((mod) => mod.Marker), {
  ssr: false,
});

interface LocationData {
  latitude: number;
  longitude: number;
  address: string;
  city: string;
  region: string;
  country: string;
  confidence: number;
}

interface LocationPickerProps {
  value?: LocationData;
  onChange: (location: LocationData) => void;
  placeholder?: string;
  height?: string;
  showCurrentLocation?: boolean;
  defaultCountry?: string;
  mapboxToken?: string;
}

export default function LocationPicker({
  value,
  onChange,
  placeholder = 'Введите адрес или выберите точку на карте',
  height = '400px',
  showCurrentLocation = true,
  defaultCountry = 'Србија',
  mapboxToken = process.env.NEXT_PUBLIC_MAPBOX_TOKEN,
}: LocationPickerProps) {
  const t = useTranslations();
  const [mode, setMode] = useState<'search' | 'map'>('search');
  const [address, setAddress] = useState(value?.address || '');
  const [mapLocation, setMapLocation] = useState({
    latitude: value?.latitude || 44.8167,
    longitude: value?.longitude || 20.4667,
    zoom: value ? 14 : 10,
  });
  const [selectedLocation, setSelectedLocation] = useState<LocationData | null>(
    value || null
  );
  const [isReverseGeocoding, setIsReverseGeocoding] = useState(false);

  // Обновляем состояние при изменении value
  useEffect(() => {
    if (value) {
      setAddress(value.address);
      setSelectedLocation(value);
      setMapLocation({
        latitude: value.latitude,
        longitude: value.longitude,
        zoom: 14,
      });
    }
  }, [value]);

  // Обратное геокодирование для получения адреса по координатам
  const reverseGeocode = useCallback(
    async (lat: number, lng: number) => {
      setIsReverseGeocoding(true);
      try {
        const response = await fetch(
          `https://api.mapbox.com/geocoding/v5/mapbox.places/${lng},${lat}.json?` +
            `access_token=${mapboxToken}&language=ru&types=address,place`
        );

        if (!response.ok) throw new Error('Geocoding failed');

        const data = await response.json();
        if (data.features && data.features.length > 0) {
          const feature = data.features[0];

          // Извлекаем компоненты адреса
          const context = feature.context || [];
          const city =
            context.find((c: any) => c.id.includes('place'))?.text || '';
          const region =
            context.find((c: any) => c.id.includes('region'))?.text || '';
          const country =
            context.find((c: any) => c.id.includes('country'))?.text ||
            defaultCountry;

          const locationData: LocationData = {
            latitude: lat,
            longitude: lng,
            address:
              feature.place_name || `${lat.toFixed(6)}, ${lng.toFixed(6)}`,
            city: city,
            region: region,
            country: country,
            confidence: 0.9, // Высокая уверенность для точки, выбранной пользователем
          };

          setAddress(locationData.address);
          setSelectedLocation(locationData);
          onChange(locationData);
        }
      } catch (error) {
        console.error('Reverse geocoding error:', error);
        // Если геокодирование не удалось, используем координаты
        const locationData: LocationData = {
          latitude: lat,
          longitude: lng,
          address: `${lat.toFixed(6)}, ${lng.toFixed(6)}`,
          city: '',
          region: '',
          country: defaultCountry,
          confidence: 0.5,
        };

        setAddress(locationData.address);
        setSelectedLocation(locationData);
        onChange(locationData);
      } finally {
        setIsReverseGeocoding(false);
      }
    },
    [mapboxToken, defaultCountry, onChange]
  );

  // Обработка выбора адреса из поиска
  const handleAddressSelect = useCallback(
    (value: string, result?: AddressGeocodingResult) => {
      setAddress(value);

      if (result) {
        // Извлекаем city из результата геокодирования
        const city = result.city || result.properties?.city || '';
        const region = result.region || result.properties?.region || '';
        const country =
          result.country || result.properties?.country || defaultCountry;

        const locationData: LocationData = {
          latitude: result.location.lat,
          longitude: result.location.lng,
          address: result.address || value,
          city: city,
          region: region,
          country: country,
          confidence: result.confidence,
        };

        setSelectedLocation(locationData);
        setMapLocation({
          latitude: result.location.lat,
          longitude: result.location.lng,
          zoom: 16,
        });
        onChange(locationData);
      }
    },
    [defaultCountry, onChange]
  );

  // Обработка клика по карте
  const handleMapClick = useCallback(
    (event: any) => {
      const { lng, lat } = event.lngLat;
      reverseGeocode(lat, lng);
    },
    [reverseGeocode]
  );

  return (
    <div className="space-y-4">
      {/* Переключатель режимов */}
      <div className="flex gap-2">
        <button
          type="button"
          onClick={() => setMode('search')}
          className={`btn btn-sm flex-1 ${mode === 'search' ? 'btn-primary' : 'btn-outline'}`}
        >
          <MagnifyingGlassIcon className="w-4 h-4 mr-1" />
          Поиск по адресу
        </button>
        <button
          type="button"
          onClick={() => setMode('map')}
          className={`btn btn-sm flex-1 ${mode === 'map' ? 'btn-primary' : 'btn-outline'}`}
        >
          <MapPinIcon className="w-4 h-4 mr-1" />
          Выбрать на карте
        </button>
      </div>

      {/* Режим поиска */}
      {mode === 'search' && (
        <div className="space-y-4">
          <SmartAddressInput
            value={address}
            onChange={handleAddressSelect}
            onLocationSelect={(locationData) => {
              handleAddressSelect(locationData.address, {
                id: `location-${Date.now()}`,
                text: locationData.address,
                place_name: locationData.address,
                location: { lat: locationData.lat, lng: locationData.lng },
                address_components: {
                  formatted: locationData.address,
                  country: defaultCountry,
                },
                confidence: locationData.confidence,
                place_types: ['address'],
                address: locationData.address,
                city: '', // SmartAddressInput не предоставляет city
                region: '',
                country: defaultCountry,
              });
            }}
            placeholder={placeholder}
            showCurrentLocation={showCurrentLocation}
            country={['rs', 'hr', 'ba', 'me']}
            language="ru"
          />

          {selectedLocation && (
            <div className="p-4 bg-success/10 border border-success/20 rounded-lg">
              <h4 className="font-medium text-success-content mb-2">
                ✅ Местоположение выбрано
              </h4>
              <div className="text-sm text-success-content/80 space-y-2">
                <div>
                  <label className="font-medium">Адрес:</label>
                  <input
                    type="text"
                    value={selectedLocation.address}
                    onChange={(e) => {
                      const updatedLocation = {
                        ...selectedLocation,
                        address: e.target.value,
                      };
                      setSelectedLocation(updatedLocation);
                      onChange(updatedLocation);
                    }}
                    className="input input-sm input-bordered w-full mt-1"
                    placeholder="Введите адрес"
                  />
                </div>
                <div>
                  <label className="font-medium">Город:</label>
                  <input
                    type="text"
                    value={selectedLocation.city}
                    onChange={(e) => {
                      const updatedLocation = {
                        ...selectedLocation,
                        city: e.target.value,
                      };
                      setSelectedLocation(updatedLocation);
                      onChange(updatedLocation);
                    }}
                    className="input input-sm input-bordered w-full mt-1"
                    placeholder="Введите город"
                  />
                </div>
                <p>
                  <strong>Координаты:</strong>{' '}
                  {selectedLocation.latitude.toFixed(6)},{' '}
                  {selectedLocation.longitude.toFixed(6)}
                </p>
              </div>
            </div>
          )}
        </div>
      )}

      {/* Режим карты */}
      {mode === 'map' && (
        <div className="space-y-4">
          <div className="alert alert-info">
            <MapPinIcon className="w-5 h-5" />
            <span>Кликните на карте, чтобы выбрать местоположение</span>
          </div>

          <div
            style={{ height: height }}
            className="relative rounded-lg overflow-hidden"
          >
            <Map
              {...mapLocation}
              onMove={(evt) => setMapLocation(evt.viewState)}
              onClick={handleMapClick}
              mapboxAccessToken={mapboxToken}
              mapStyle="mapbox://styles/mapbox/streets-v12"
              attributionControl={false}
            >
              {selectedLocation && (
                <Marker
                  longitude={selectedLocation.longitude}
                  latitude={selectedLocation.latitude}
                  anchor="bottom"
                >
                  <div className="relative">
                    <MapPinIcon className="w-8 h-8 text-primary drop-shadow-lg" />
                    {isReverseGeocoding && (
                      <div className="absolute -top-8 left-1/2 transform -translate-x-1/2">
                        <div className="loading loading-spinner loading-sm"></div>
                      </div>
                    )}
                  </div>
                </Marker>
              )}
            </Map>
          </div>

          {selectedLocation && (
            <div className="p-4 bg-success/10 border border-success/20 rounded-lg">
              <h4 className="font-medium text-success-content mb-2">
                ✅ Местоположение выбрано
              </h4>
              <div className="text-sm text-success-content/80 space-y-2">
                <div>
                  <label className="font-medium">Адрес:</label>
                  <input
                    type="text"
                    value={selectedLocation.address}
                    onChange={(e) => {
                      const updatedLocation = {
                        ...selectedLocation,
                        address: e.target.value,
                      };
                      setSelectedLocation(updatedLocation);
                      onChange(updatedLocation);
                    }}
                    className="input input-sm input-bordered w-full mt-1"
                    placeholder="Введите адрес"
                  />
                </div>
                <div>
                  <label className="font-medium">Город:</label>
                  <input
                    type="text"
                    value={selectedLocation.city}
                    onChange={(e) => {
                      const updatedLocation = {
                        ...selectedLocation,
                        city: e.target.value,
                      };
                      setSelectedLocation(updatedLocation);
                      onChange(updatedLocation);
                    }}
                    className="input input-sm input-bordered w-full mt-1"
                    placeholder="Введите город"
                  />
                </div>
                <p>
                  <strong>Координаты:</strong>{' '}
                  {selectedLocation.latitude.toFixed(6)},{' '}
                  {selectedLocation.longitude.toFixed(6)}
                </p>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
}
