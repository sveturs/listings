'use client';

import { useEffect, useRef } from 'react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import { BEXTrackingEvent } from './BEXTracker';

// Исправляем проблему с иконками Leaflet
delete (L.Icon.Default.prototype as any)._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: '/leaflet/marker-icon-2x.png',
  iconUrl: '/leaflet/marker-icon.png',
  shadowUrl: '/leaflet/marker-shadow.png',
});

interface Location {
  city: string;
  address?: string;
  latitude?: number;
  longitude?: number;
}

interface Props {
  currentLocation?: Location;
  deliveryLocation?: Location;
  events?: BEXTrackingEvent[];
  className?: string;
}

export default function BEXMap({
  currentLocation,
  deliveryLocation,
  events = [],
  className = '',
}: Props) {
  const mapContainer = useRef<HTMLDivElement>(null);
  const map = useRef<L.Map | null>(null);
  const markersRef = useRef<L.Marker[]>([]);
  const polylineRef = useRef<L.Polyline | null>(null);

  useEffect(() => {
    if (!mapContainer.current || map.current) return;

    // Инициализация карты
    map.current = L.map(mapContainer.current).setView([44.7866, 20.4489], 7); // Центр Сербии

    // Добавляем тайлы OpenStreetMap
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors',
      maxZoom: 19,
    }).addTo(map.current);

    return () => {
      if (map.current) {
        map.current.remove();
        map.current = null;
      }
    };
  }, []);

  useEffect(() => {
    if (!map.current) return;

    // Очищаем старые маркеры
    markersRef.current.forEach((marker) => marker.remove());
    markersRef.current = [];

    // Удаляем старый маршрут
    if (polylineRef.current) {
      polylineRef.current.remove();
      polylineRef.current = null;
    }

    const bounds = L.latLngBounds([]);
    const routePoints: L.LatLngExpression[] = [];

    // Добавляем маркер текущего местоположения
    if (currentLocation?.latitude && currentLocation?.longitude) {
      const currentIcon = L.divIcon({
        className: 'custom-div-icon',
        html: `
          <div class="relative">
            <div class="absolute -top-8 -left-8 w-16 h-16">
              <div class="w-full h-full rounded-full bg-primary/20 animate-ping"></div>
            </div>
            <svg class="w-8 h-8 text-primary" fill="currentColor" viewBox="0 0 20 20">
              <path d="M8 16.5a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zM15 16.5a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z"/>
              <path d="M3 4a1 1 0 00-1 1v10a1 1 0 001 1h1.05a2.5 2.5 0 014.9 0H10a1 1 0 001-1V5a1 1 0 00-1-1H3zM14 7a1 1 0 00-1 1v6.05A2.5 2.5 0 0115.95 16H17a1 1 0 001-1v-5a1 1 0 00-.293-.707l-2-2A1 1 0 0015 7h-1z"/>
            </svg>
          </div>
        `,
        iconSize: [32, 32],
        iconAnchor: [16, 16],
      });

      const marker = L.marker(
        [currentLocation.latitude, currentLocation.longitude],
        { icon: currentIcon }
      ).addTo(map.current!).bindPopup(`
          <div class="p-2">
            <h4 class="font-semibold">Текущее местоположение</h4>
            <p class="text-sm">${currentLocation.address || currentLocation.city}</p>
          </div>
        `);

      markersRef.current.push(marker);
      bounds.extend([currentLocation.latitude, currentLocation.longitude]);
      routePoints.push([currentLocation.latitude, currentLocation.longitude]);
    }

    // Добавляем промежуточные точки из событий
    const locationEvents = events.filter((e) => e.location);
    const uniqueLocations = new Map<
      string,
      { lat: number; lng: number; events: BEXTrackingEvent[] }
    >();

    // Генерируем координаты для событий (в реальном приложении они должны приходить с сервера)
    locationEvents.forEach((event) => {
      if (event.location) {
        // Простая логика для демо - размещаем точки между текущим и конечным местоположением
        let lat = 44.7866; // По умолчанию центр Сербии
        let lng = 20.4489;

        if (event.location.includes('Белград')) {
          lat = 44.8176;
          lng = 20.4633;
        } else if (event.location.includes('Нови Сад')) {
          lat = 45.2671;
          lng = 19.8335;
        } else if (event.location.includes('Ниш')) {
          lat = 43.3209;
          lng = 21.8954;
        }

        const key = `${lat},${lng}`;
        if (!uniqueLocations.has(key)) {
          uniqueLocations.set(key, { lat, lng, events: [] });
        }
        uniqueLocations.get(key)!.events.push(event);
      }
    });

    // Добавляем маркеры для промежуточных точек
    uniqueLocations.forEach(({ lat, lng, events }) => {
      const eventIcon = L.divIcon({
        className: 'custom-div-icon',
        html: `
          <div class="w-3 h-3 bg-base-300 rounded-full border-2 border-white shadow-lg"></div>
        `,
        iconSize: [12, 12],
        iconAnchor: [6, 6],
      });

      const popupContent = `
        <div class="p-2 max-w-xs">
          <h4 class="font-semibold mb-2">${events[0].location}</h4>
          ${events
            .map(
              (e) => `
            <div class="text-sm mb-1">
              <div class="font-medium">${e.status_text}</div>
              <div class="text-gray-600">${new Date(e.timestamp).toLocaleString('ru-RU')}</div>
            </div>
          `
            )
            .join('')}
        </div>
      `;

      const marker = L.marker([lat, lng], { icon: eventIcon })
        .addTo(map.current!)
        .bindPopup(popupContent);

      markersRef.current.push(marker);
      bounds.extend([lat, lng]);
      routePoints.push([lat, lng]);
    });

    // Добавляем маркер пункта назначения
    if (deliveryLocation?.latitude && deliveryLocation?.longitude) {
      const deliveryIcon = L.divIcon({
        className: 'custom-div-icon',
        html: `
          <div class="relative">
            <svg class="w-8 h-8 text-success" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M5.05 4.05a7 7 0 119.9 9.9L10 18.9l-4.95-4.95a7 7 0 010-9.9zM10 11a2 2 0 100-4 2 2 0 000 4z" clip-rule="evenodd"/>
            </svg>
          </div>
        `,
        iconSize: [32, 32],
        iconAnchor: [16, 32],
      });

      const marker = L.marker(
        [deliveryLocation.latitude, deliveryLocation.longitude],
        { icon: deliveryIcon }
      ).addTo(map.current!).bindPopup(`
          <div class="p-2">
            <h4 class="font-semibold">Пункт назначения</h4>
            <p class="text-sm">${deliveryLocation.address}</p>
            <p class="text-sm text-gray-600">${deliveryLocation.city}</p>
          </div>
        `);

      markersRef.current.push(marker);
      bounds.extend([deliveryLocation.latitude, deliveryLocation.longitude]);
      routePoints.push([deliveryLocation.latitude, deliveryLocation.longitude]);
    }

    // Рисуем маршрут
    if (routePoints.length > 1) {
      // Анимированная пунктирная линия для планируемого маршрута
      polylineRef.current = L.polyline(routePoints, {
        color: '#3B82F6',
        weight: 3,
        opacity: 0.6,
        dashArray: '10, 10',
        className: 'animate-dash',
      }).addTo(map.current!);

      // Добавляем стрелки направления
      const decorator = (L as any).polylineDecorator(polylineRef.current, {
        patterns: [
          {
            offset: '50%',
            repeat: 100,
            symbol: (L as any).Symbol.arrowHead({
              pixelSize: 12,
              polygon: false,
              pathOptions: {
                stroke: true,
                color: '#3B82F6',
                weight: 2,
              },
            }),
          },
        ],
      });
      decorator.addTo(map.current!);
    }

    // Подстраиваем масштаб карты под все маркеры
    if (bounds.isValid()) {
      map.current!.fitBounds(bounds, { padding: [50, 50] });
    }
  }, [currentLocation, deliveryLocation, events]);

  return (
    <>
      <style jsx global>{`
        .custom-div-icon {
          background: transparent;
          border: none;
        }

        @keyframes dash {
          to {
            stroke-dashoffset: -20;
          }
        }

        .animate-dash {
          animation: dash 1s linear infinite;
        }

        .leaflet-popup-content-wrapper {
          border-radius: 0.5rem;
          box-shadow: 0 10px 25px rgba(0, 0, 0, 0.1);
        }

        .leaflet-popup-content {
          margin: 0;
          min-width: 200px;
        }
      `}</style>

      <div
        ref={mapContainer}
        className={`h-96 rounded-lg shadow-inner ${className}`}
        style={{ minHeight: '400px' }}
      />
    </>
  );
}
