'use client';

import React, { useEffect, useRef, useState, useCallback } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import { useAddressGeocoding } from '@/hooks/useAddressGeocoding';
import configManager from '@/config';

export interface AddressConfirmationMapProps {
  address: string;
  initialLocation?: { lat: number; lng: number };
  onLocationConfirm: (location: {
    lat: number;
    lng: number;
    address: string;
    confidence: number;
  }) => void;
  onLocationChange?: (location: { lat: number; lng: number }) => void;
  editable?: boolean;
  zoom?: number;
  height?: string;
  className?: string;
}

export default function AddressConfirmationMap({
  address,
  initialLocation,
  onLocationConfirm,
  onLocationChange,
  editable = true,
  zoom = 16,
  height = '400px',
  className = '',
}: AddressConfirmationMapProps) {
  const mapContainer = useRef<HTMLDivElement>(null);
  const map = useRef<mapboxgl.Map | null>(null);
  const marker = useRef<mapboxgl.Marker | null>(null);

  const [currentLocation, setCurrentLocation] = useState(initialLocation);
  const [_isReversing, _setIsReversing] = useState(false);
  const [_confidence, _setConfidence] = useState(0);
  const [_reverseAddress, _setReverseAddress] = useState('');

  const { reverseGeocode: _reverseGeocode, validateAddress } =
    useAddressGeocoding();

  // Инициализация карты
  useEffect(() => {
    if (!mapContainer.current || map.current) return;

    // Получаем токен из ConfigManager
    const mapboxToken = configManager.getMapboxToken();
    if (!mapboxToken) {
      console.error('MAPBOX_ACCESS_TOKEN is not set');
      return;
    }

    mapboxgl.accessToken = mapboxToken;

    map.current = new mapboxgl.Map({
      container: mapContainer.current,
      style: 'mapbox://styles/mapbox/streets-v12',
      center: currentLocation
        ? [currentLocation.lng, currentLocation.lat]
        : [20.4651, 44.8176], // Belgrade default
      zoom: currentLocation ? zoom : 10,
      attributionControl: false,
    });

    // Добавляем контролы
    map.current.addControl(new mapboxgl.NavigationControl(), 'top-right');
    map.current.addControl(
      new mapboxgl.AttributionControl({
        compact: true,
      }),
      'bottom-right'
    );

    return () => {
      if (map.current) {
        map.current.remove();
        map.current = null;
      }
    };
  }, [currentLocation, zoom]);

  // Обновление маркера при изменении локации
  useEffect(() => {
    if (!map.current || !currentLocation) return;

    // Удаляем старый маркер
    if (marker.current) {
      marker.current.remove();
    }

    // Создаем новый маркер
    const markerElement = document.createElement('div');
    markerElement.className = 'custom-marker';
    markerElement.innerHTML = `
      <div class="relative">
        <div class="w-8 h-8 bg-primary rounded-full border-4 border-white shadow-lg flex items-center justify-center">
          <svg class="w-4 h-4 text-primary-content" fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z"/>
          </svg>
        </div>
        <div class="absolute top-full left-1/2 transform -translate-x-1/2 w-1 h-3 bg-primary"></div>
      </div>
    `;

    marker.current = new mapboxgl.Marker({
      element: markerElement,
      draggable: editable,
    })
      .setLngLat([currentLocation.lng, currentLocation.lat])
      .addTo(map.current);

    // Обработчик перетаскивания
    if (editable) {
      marker.current.on('dragend', () => {
        if (!marker.current) return;

        const lngLat = marker.current.getLngLat();
        const newLocation = { lat: lngLat.lat, lng: lngLat.lng };

        setCurrentLocation(newLocation);
        onLocationChange?.(newLocation);

        // TODO: Обратное геокодирование при перетаскивании
        // handleReverseGeocode(newLocation);
      });
    }

    // Центрируем карту на маркере
    map.current.flyTo({
      center: [currentLocation.lng, currentLocation.lat],
      zoom: zoom,
      duration: 1000,
    });
  }, [currentLocation, editable, zoom, onLocationChange]);

  // Обновление локации при изменении адреса
  useEffect(() => {
    if (!address.trim()) return;

    const geocodeAddress = async () => {
      try {
        const result = await validateAddress(address);

        if (result.success && result.location) {
          const newLocation = {
            lat: result.location.lat,
            lng: result.location.lng,
          };

          setCurrentLocation(newLocation);
          _setReverseAddress(result.formatted_address || address);
          _setConfidence(result.confidence || 0.5);
          onLocationChange?.(newLocation);
        }
      } catch (error) {
        console.error('Address validation error:', error);
      }
    };

    geocodeAddress();
  }, [address, validateAddress, onLocationChange]);

  // Клик по карте для установки маркера
  useEffect(() => {
    if (!map.current || !editable) return;

    const handleMapClick = (e: mapboxgl.MapMouseEvent) => {
      const newLocation = {
        lat: e.lngLat.lat,
        lng: e.lngLat.lng,
      };

      setCurrentLocation(newLocation);
      onLocationChange?.(newLocation);
      // TODO: handleReverseGeocode(newLocation);
    };

    map.current.on('click', handleMapClick);

    return () => {
      map.current?.off('click', handleMapClick);
    };
  }, [editable, onLocationChange]);

  // Подтверждение локации
  const handleConfirmLocation = useCallback(() => {
    if (!currentLocation) return;

    onLocationConfirm({
      lat: currentLocation.lat,
      lng: currentLocation.lng,
      address: _reverseAddress || address,
      confidence: _confidence,
    });
  }, [
    currentLocation,
    _reverseAddress,
    address,
    _confidence,
    onLocationConfirm,
  ]);

  // Определение цвета по уровню доверия
  const _getConfidenceColor = (conf: number) => {
    if (conf >= 0.8) return 'text-success';
    if (conf >= 0.6) return 'text-warning';
    return 'text-error';
  };

  const getConfidenceBadgeClass = (conf: number) => {
    if (conf >= 0.8) return 'badge-success';
    if (conf >= 0.6) return 'badge-warning';
    return 'badge-error';
  };

  return (
    <div className={`w-full ${className}`}>
      {/* Карта */}
      <div
        ref={mapContainer}
        className="w-full rounded-lg border border-base-300 shadow-sm"
        style={{ height }}
      />

      {/* Информационная панель */}
      {currentLocation && (
        <div className="mt-4 p-4 bg-base-100 border border-base-300 rounded-lg">
          {/* Заголовок */}
          <div className="flex items-center justify-between mb-3">
            <h3 className="text-lg font-semibold">
              Подтверждение местоположения
            </h3>

            {/* Показатель доверия */}
            <div className={`badge ${getConfidenceBadgeClass(_confidence)}`}>
              Точность: {Math.round(_confidence * 100)}%
            </div>
          </div>

          {/* Адрес */}
          <div className="space-y-2">
            <div>
              <span className="text-sm font-medium text-base-content/70">
                Введенный адрес:
              </span>
              <p className="text-sm mt-1">{address}</p>
            </div>

            {_reverseAddress && _reverseAddress !== address && (
              <div>
                <span className="text-sm font-medium text-base-content/70">
                  Адрес по координатам:
                  {_isReversing && (
                    <span className="loading loading-spinner loading-xs ml-2"></span>
                  )}
                </span>
                <p className="text-sm mt-1">{_reverseAddress}</p>
              </div>
            )}
          </div>

          {/* Координаты */}
          <div className="mt-3 grid grid-cols-2 gap-4 text-xs text-base-content/50">
            <div>
              <span className="font-medium">Широта:</span>{' '}
              {currentLocation.lat.toFixed(6)}
            </div>
            <div>
              <span className="font-medium">Долгота:</span>{' '}
              {currentLocation.lng.toFixed(6)}
            </div>
          </div>

          {/* Инструкции */}
          {editable && (
            <div className="mt-3 p-3 bg-info/10 border border-info/20 rounded-lg">
              <p className="text-sm text-info-content">
                💡 <strong>Подсказка:</strong> Вы можете перетащить маркер или
                кликнуть по карте для точной настройки местоположения.
              </p>
            </div>
          )}

          {/* Предупреждения */}
          {_confidence < 0.7 && (
            <div className="mt-3 p-3 bg-warning/10 border border-warning/20 rounded-lg">
              <p className="text-sm text-warning-content">
                ⚠️ <strong>Низкая точность:</strong> Рекомендуется проверить и
                при необходимости скорректировать местоположение на карте.
              </p>
            </div>
          )}

          {/* Кнопки */}
          <div className="mt-4 flex gap-2">
            <button
              onClick={handleConfirmLocation}
              disabled={!currentLocation}
              className="btn btn-primary flex-1"
            >
              <svg
                className="w-4 h-4 mr-2"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M5 13l4 4L19 7"
                />
              </svg>
              Подтвердить местоположение
            </button>

            <button
              onClick={() => {
                // Сброс к исходному состоянию
                if (initialLocation) {
                  setCurrentLocation(initialLocation);
                  _setReverseAddress('');
                  _setConfidence(0);
                }
              }}
              className="btn btn-ghost"
              disabled={!initialLocation}
            >
              Сбросить
            </button>
          </div>
        </div>
      )}

      {/* Состояние без локации */}
      {!currentLocation && (
        <div className="mt-4 p-6 text-center bg-base-100 border border-base-300 rounded-lg">
          <div className="text-base-content/50 mb-2">
            <svg
              className="w-12 h-12 mx-auto mb-2"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1}
                d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1}
                d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>
          </div>

          <h3 className="text-lg font-medium mb-2">Укажите местоположение</h3>
          <p className="text-sm text-base-content/70">
            Введите адрес в поле выше или кликните по карте для установки
            местоположения
          </p>
        </div>
      )}
    </div>
  );
}
