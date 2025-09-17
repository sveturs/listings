"use client";

import { useEffect, useRef, useState } from 'react';
import mapboxgl from 'mapbox-gl';
import 'mapbox-gl/dist/mapbox-gl.css';
import { useTranslations } from 'next-intl';

// Mapbox token из .env
mapboxgl.accessToken = process.env.NEXT_PUBLIC_MAPBOX_ACCESS_TOKEN || '';

interface DeliveryData {
  id: number;
  pickup_latitude: number;
  pickup_longitude: number;
  delivery_latitude: number;
  delivery_longitude: number;
  pickup_address: string;
  delivery_address: string;
  courier_location?: {
    latitude: number;
    longitude: number;
    speed: number;
    heading: number;
  };
  nearby_items?: Array<{
    id: number;
    title: string;
    price: number;
    currency: string;
    location: {
      latitude: number;
      longitude: number;
    };
    distance_meters: number;
    store_name?: string;
    store_logo?: string;
  }>;
}

interface TrackingMapProps {
  delivery: DeliveryData;
  onRequestETA?: () => void;
}

export function TrackingMap({ delivery, onRequestETA }: TrackingMapProps) {
  const mapContainer = useRef<HTMLDivElement>(null);
  const map = useRef<mapboxgl.Map | null>(null);
  const courierMarker = useRef<mapboxgl.Marker | null>(null);
  const [mapLoaded, setMapLoaded] = useState(false);
  const t = useTranslations('tracking.map');

  // Инициализация карты
  useEffect(() => {
    if (!mapContainer.current || map.current) return;

    try {
      // Центр между точками получения и доставки
      const centerLat = (delivery.pickup_latitude + delivery.delivery_latitude) / 2;
      const centerLng = (delivery.pickup_longitude + delivery.delivery_longitude) / 2;

      map.current = new mapboxgl.Map({
        container: mapContainer.current,
        style: 'mapbox://styles/mapbox/streets-v12',
        center: [centerLng, centerLat],
        zoom: 12,
        language: 'ru'
      });

      map.current.on('load', () => {
        setMapLoaded(true);
        addInitialMarkers();
        fitBounds();
      });

      // Контролы навигации
      map.current.addControl(new mapboxgl.NavigationControl(), 'top-right');
      
      // Контроль масштаба
      map.current.addControl(new mapboxgl.ScaleControl({
        maxWidth: 100,
        unit: 'metric'
      }), 'bottom-left');

    } catch (error) {
      console.error('Error initializing map:', error);
    }

    return () => {
      if (map.current) {
        map.current.remove();
        map.current = null;
      }
    };
  }, []);

  // Добавление начальных маркеров
  const addInitialMarkers = () => {
    if (!map.current) return;

    // Маркер точки получения
    const pickupEl = document.createElement('div');
    pickupEl.className = 'pickup-marker';
    pickupEl.innerHTML = '🏢';
    pickupEl.style.fontSize = '24px';
    
    new mapboxgl.Marker(pickupEl)
      .setLngLat([delivery.pickup_longitude, delivery.pickup_latitude])
      .setPopup(
        new mapboxgl.Popup({ offset: 25 })
          .setHTML(`
            <div class="p-2">
              <div class="font-semibold text-sm">${t('pickup')}</div>
              <div class="text-xs text-gray-600">${delivery.pickup_address}</div>
            </div>
          `)
      )
      .addTo(map.current);

    // Маркер точки доставки
    const deliveryEl = document.createElement('div');
    deliveryEl.className = 'delivery-marker';
    deliveryEl.innerHTML = '🏠';
    deliveryEl.style.fontSize = '24px';
    
    new mapboxgl.Marker(deliveryEl)
      .setLngLat([delivery.delivery_longitude, delivery.delivery_latitude])
      .setPopup(
        new mapboxgl.Popup({ offset: 25 })
          .setHTML(`
            <div class="p-2">
              <div class="font-semibold text-sm">${t('delivery')}</div>
              <div class="text-xs text-gray-600">${delivery.delivery_address}</div>
            </div>
          `)
      )
      .addTo(map.current);
  };

  // Обновление маркера курьера
  useEffect(() => {
    if (!map.current || !mapLoaded || !delivery.courier_location) return;

    const { latitude, longitude, speed, heading } = delivery.courier_location;

    // Создаём или обновляем маркер курьера
    if (!courierMarker.current) {
      const courierEl = document.createElement('div');
      courierEl.className = 'courier-marker';
      courierEl.innerHTML = '🚴‍♂️';
      courierEl.style.fontSize = '32px';
      courierEl.style.transform = `rotate(${heading}deg)`;
      courierEl.style.transition = 'transform 0.5s ease';
      
      courierMarker.current = new mapboxgl.Marker(courierEl)
        .setLngLat([longitude, latitude])
        .setPopup(
          new mapboxgl.Popup({ offset: 25 })
            .setHTML(`
              <div class="p-2">
                <div class="font-semibold text-sm">${t('courier')}</div>
                <div class="text-xs text-gray-600">
                  📈 ${speed.toFixed(1)} км/ч<br/>
                  🧭 ${heading}°
                </div>
              </div>
            `)
        )
        .addTo(map.current);
    } else {
      // Плавное обновление позиции
      courierMarker.current.setLngLat([longitude, latitude]);
      
      // Обновляем направление
      const markerElement = courierMarker.current.getElement();
      if (markerElement) {
        markerElement.style.transform = `rotate(${heading}deg)`;
      }
      
      // Обновляем popup
      courierMarker.current.getPopup().setHTML(`
        <div class="p-2">
          <div class="font-semibold text-sm">${t('courier')}</div>
          <div class="text-xs text-gray-600">
            📈 ${speed.toFixed(1)} км/ч<br/>
            🧭 ${heading}°
          </div>
        </div>
      `);
    }
  }, [delivery.courier_location, mapLoaded]);

  // Отображение товаров поблизости
  useEffect(() => {
    if (!map.current || !mapLoaded || !delivery.nearby_items) return;

    // Удаляем предыдущие маркеры товаров
    const existingMarkers = document.querySelectorAll('.nearby-item-marker');
    existingMarkers.forEach(marker => marker.remove());

    // Добавляем новые маркеры
    delivery.nearby_items.forEach(item => {
      const itemEl = document.createElement('div');
      itemEl.className = 'nearby-item-marker';
      itemEl.innerHTML = '🛒';
      itemEl.style.fontSize = '18px';
      itemEl.style.cursor = 'pointer';
      
      new mapboxgl.Marker(itemEl)
        .setLngLat([item.location.longitude, item.location.latitude])
        .setPopup(
          new mapboxgl.Popup({ offset: 25 })
            .setHTML(`
              <div class="p-3 max-w-xs">
                <div class="font-semibold text-sm mb-1">${item.title}</div>
                <div class="text-lg font-bold text-green-600 mb-1">
                  ${item.price} ${item.currency}
                </div>
                ${item.store_name ? `<div class="text-xs text-gray-600 mb-1">🏢 ${item.store_name}</div>` : ''}
                <div class="text-xs text-gray-500">
                  📍 ${(item.distance_meters / 1000).toFixed(1)} км
                </div>
                <button class="btn btn-xs btn-primary mt-2 w-full">
                  👁️ ${t('viewItem')}
                </button>
              </div>
            `)
        )
        .addTo(map.current!);
    });
  }, [delivery.nearby_items, mapLoaded]);

  // Подгонка границ карты
  const fitBounds = () => {
    if (!map.current) return;

    const bounds = new mapboxgl.LngLatBounds();
    bounds.extend([delivery.pickup_longitude, delivery.pickup_latitude]);
    bounds.extend([delivery.delivery_longitude, delivery.delivery_latitude]);
    
    if (delivery.courier_location) {
      bounds.extend([delivery.courier_location.longitude, delivery.courier_location.latitude]);
    }

    map.current.fitBounds(bounds, {
      padding: 50,
      maxZoom: 15
    });
  };

  // Центрирование на курьере
  const centerOnCourier = () => {
    if (!map.current || !delivery.courier_location) return;
    
    map.current.flyTo({
      center: [delivery.courier_location.longitude, delivery.courier_location.latitude],
      zoom: 16,
      duration: 1000
    });
  };

  return (
    <div className="relative w-full h-96 lg:h-[500px]">
      <div ref={mapContainer} className="w-full h-full rounded-lg" />
      
      {/* Кнопки управления */}
      <div className="absolute top-4 left-4 flex flex-col gap-2 z-10">
        <button
          onClick={fitBounds}
          className="btn btn-sm bg-white border border-gray-300 text-gray-700 hover:bg-gray-50"
          title={t('fitBounds')}
        >
          🗺️
        </button>
        
        {delivery.courier_location && (
          <button
            onClick={centerOnCourier}
            className="btn btn-sm bg-white border border-gray-300 text-gray-700 hover:bg-gray-50"
            title={t('centerOnCourier')}
          >
            🚴‍♂️
          </button>
        )}
        
        {onRequestETA && (
          <button
            onClick={onRequestETA}
            className="btn btn-sm btn-primary"
            title={t('requestETA')}
          >
            ⏱️
          </button>
        )}
      </div>
      
      {/* Легенда */}
      <div className="absolute bottom-4 right-4 bg-white p-3 rounded-lg shadow-lg text-xs z-10">
        <div className="font-semibold mb-2">{t('legend.title')}</div>
        <div className="flex flex-col gap-1">
          <div className="flex items-center gap-2">
            <span>🏢</span>
            <span>{t('legend.pickup')}</span>
          </div>
          <div className="flex items-center gap-2">
            <span>🏠</span>
            <span>{t('legend.delivery')}</span>
          </div>
          {delivery.courier_location && (
            <div className="flex items-center gap-2">
              <span>🚴‍♂️</span>
              <span>{t('legend.courier')}</span>
            </div>
          )}
          {delivery.nearby_items && delivery.nearby_items.length > 0 && (
            <div className="flex items-center gap-2">
              <span>🛒</span>
              <span>{t('legend.items')}</span>
            </div>
          )}
        </div>
      </div>
      
      {!mapLoaded && (
        <div className="absolute inset-0 bg-gray-100 flex items-center justify-center rounded-lg">
          <div className="text-center">
            <div className="loading loading-spinner loading-lg text-primary"></div>
            <p className="mt-2 text-sm text-gray-600">{t('loading')}</p>
          </div>
        </div>
      )}
    </div>
  );
}