'use client';

import { useEffect, useRef, useState, useCallback } from 'react';
import { useTranslations } from 'next-intl';
import { FiTruck, FiCheckCircle, FiAlertTriangle } from 'react-icons/fi';

interface MapMarker {
  id: string;
  lat: number;
  lng: number;
  type: 'sender' | 'receiver' | 'in_transit' | 'delivered' | 'problem';
  shipment_id: number;
  tracking_number: string;
  status: string;
  city: string;
  address: string;
  delivery_time?: string;
  problem_description?: string;
}

interface DeliveryMapProps {
  markers: MapMarker[];
  onMarkerClick?: (marker: MapMarker) => void;
  height?: string;
  showControls?: boolean;
}

export default function DeliveryMap({
  markers,
  onMarkerClick,
  height = '400px',
  showControls = true,
}: DeliveryMapProps) {
  const t = useTranslations('admin');
  const mapRef = useRef<HTMLDivElement>(null);
  const mapInstanceRef = useRef<any>(null);
  const markersLayerRef = useRef<any>(null);
  const [leafletLoaded, setLeafletLoaded] = useState(false);
  const [selectedMarkerType, setSelectedMarkerType] = useState<string>('all');
  const [filteredMarkers, setFilteredMarkers] = useState<MapMarker[]>(markers);

  // Динамическая загрузка Leaflet
  useEffect(() => {
    const loadLeaflet = async () => {
      if (typeof window === 'undefined') return;

      try {
        // Загружаем CSS
        if (!document.querySelector('link[href*="leaflet.css"]')) {
          const link = document.createElement('link');
          link.rel = 'stylesheet';
          link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css';
          document.head.appendChild(link);
        }

        // Загружаем JS
        const L = await import('leaflet');

        // Исправляем проблему с иконками по умолчанию
        delete (L.Icon.Default.prototype as any)._getIconUrl;
        L.Icon.Default.mergeOptions({
          iconRetinaUrl:
            'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon-2x.png',
          iconUrl:
            'https://unpkg.com/leaflet@1.9.4/dist/images/marker-icon.png',
          shadowUrl:
            'https://unpkg.com/leaflet@1.9.4/dist/images/marker-shadow.png',
        });

        setLeafletLoaded(true);

        if (mapRef.current && !mapInstanceRef.current) {
          initializeMap(L);
          // После инициализации карты обновляем маркеры
          setTimeout(() => updateMarkers(L), 100);
        }
      } catch (error) {
        console.error('Error loading Leaflet:', error);
      }
    };

    loadLeaflet();
  }, [initializeMap, updateMarkers]);

  // Инициализация карты
  const initializeMap = useCallback((L: typeof import('leaflet')) => {
    if (!mapRef.current || mapInstanceRef.current) return;

    const map = L.map(mapRef.current, {
      center: [44.0165, 21.0059], // Сербия
      zoom: 7,
      zoomControl: true,
    });

    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '© OpenStreetMap contributors',
    }).addTo(map);

    mapInstanceRef.current = map;
    markersLayerRef.current = L.layerGroup().addTo(map);
  }, []);

  // Обновление маркеров на карте
  const updateMarkers = useCallback(
    (L: typeof import('leaflet')) => {
      if (!mapInstanceRef.current || !markersLayerRef.current || !leafletLoaded)
        return;

      markersLayerRef.current.clearLayers();

      filteredMarkers.forEach((marker) => {
        const icon = createCustomIcon(L, marker.type);

        const leafletMarker = L.marker([marker.lat, marker.lng], { icon })
          .bindPopup(createPopupContent(marker))
          .on('click', () => {
            if (onMarkerClick) {
              onMarkerClick(marker);
            }
          });

        markersLayerRef.current.addLayer(leafletMarker);
      });

      // Подгоняем границы карты под маркеры
      if (filteredMarkers.length > 0) {
        const group = new (L as any).featureGroup(
          markersLayerRef.current.getLayers()
        );
        mapInstanceRef.current.fitBounds(group.getBounds().pad(0.1));
      }
    },
    [
      filteredMarkers,
      leafletLoaded,
      onMarkerClick,
      createCustomIcon,
      createPopupContent,
    ]
  );

  // Создание кастомной иконки
  const createCustomIcon = (L: typeof import('leaflet'), type: string) => {
    const colors = {
      sender: '#3b82f6', // синий
      receiver: '#10b981', // зеленый
      in_transit: '#f59e0b', // желтый
      delivered: '#10b981', // зеленый
      problem: '#ef4444', // красный
    };

    const color = colors[type as keyof typeof colors] || '#6b7280';

    return L.divIcon({
      className: 'custom-marker',
      html: `
        <div style="
          background-color: ${color};
          width: 24px;
          height: 24px;
          border-radius: 50%;
          border: 2px solid white;
          box-shadow: 0 2px 4px rgba(0,0,0,0.3);
          display: flex;
          align-items: center;
          justify-content: center;
          color: white;
          font-size: 12px;
        ">
          ${getMarkerIcon(type)}
        </div>
      `,
      iconSize: [24, 24],
      iconAnchor: [12, 12],
    });
  };

  // Получение иконки для маркера
  const getMarkerIcon = (type: string) => {
    switch (type) {
      case 'sender':
        return '📤';
      case 'receiver':
        return '📍';
      case 'in_transit':
        return '🚚';
      case 'delivered':
        return '✅';
      case 'problem':
        return '⚠️';
      default:
        return '📍';
    }
  };

  // Создание содержимого попапа
  const createPopupContent = (marker: MapMarker) => {
    return `
      <div style="min-width: 200px;">
        <h3 style="margin: 0 0 8px 0; font-size: 14px; font-weight: bold;">
          ${marker.tracking_number}
        </h3>
        <p style="margin: 4px 0; font-size: 12px;">
          <strong>Город:</strong> ${marker.city}
        </p>
        <p style="margin: 4px 0; font-size: 12px;">
          <strong>Адрес:</strong> ${marker.address}
        </p>
        <p style="margin: 4px 0; font-size: 12px;">
          <strong>Статус:</strong> ${getStatusLabel(marker.status)}
        </p>
        ${
          marker.delivery_time
            ? `
          <p style="margin: 4px 0; font-size: 12px;">
            <strong>Время доставки:</strong> ${marker.delivery_time}
          </p>
        `
            : ''
        }
        ${
          marker.problem_description
            ? `
          <p style="margin: 4px 0; font-size: 12px; color: red;">
            <strong>Проблема:</strong> ${marker.problem_description}
          </p>
        `
            : ''
        }
      </div>
    `;
  };

  // Получение локализованной метки статуса
  const getStatusLabel = (status: string) => {
    const labels: { [key: string]: string } = {
      pending: 'Ожидает',
      in_transit: 'В пути',
      delivered: 'Доставлено',
      problem: 'Проблема',
      returned: 'Возвращено',
    };
    return labels[status] || status;
  };

  // Фильтрация маркеров
  useEffect(() => {
    if (selectedMarkerType === 'all') {
      setFilteredMarkers(markers);
    } else {
      setFilteredMarkers(
        markers.filter((marker) => marker.type === selectedMarkerType)
      );
    }
  }, [selectedMarkerType, markers]);

  // Обновление маркеров при изменении фильтрованного списка
  useEffect(() => {
    if (leafletLoaded) {
      // Небольшая задержка для корректного обновления
      setTimeout(() => {
        import('leaflet').then((L) => {
          updateMarkers(L);
        });
      }, 100);
    }
  }, [filteredMarkers, leafletLoaded, updateMarkers]);

  return (
    <div className="space-y-4">
      {/* Фильтры */}
      {showControls && (
        <div className="flex flex-wrap gap-2">
          <button
            className={`btn btn-sm ${selectedMarkerType === 'all' ? 'btn-primary' : 'btn-outline'}`}
            onClick={() => setSelectedMarkerType('all')}
          >
            {t('map.all')} ({markers.length})
          </button>
          <button
            className={`btn btn-sm ${selectedMarkerType === 'in_transit' ? 'btn-warning' : 'btn-outline'}`}
            onClick={() => setSelectedMarkerType('in_transit')}
          >
            <FiTruck className="w-4 h-4" />
            {t('map.in_transit')} (
            {markers.filter((m) => m.type === 'in_transit').length})
          </button>
          <button
            className={`btn btn-sm ${selectedMarkerType === 'delivered' ? 'btn-success' : 'btn-outline'}`}
            onClick={() => setSelectedMarkerType('delivered')}
          >
            <FiCheckCircle className="w-4 h-4" />
            {t('map.delivered')} (
            {markers.filter((m) => m.type === 'delivered').length})
          </button>
          <button
            className={`btn btn-sm ${selectedMarkerType === 'problem' ? 'btn-error' : 'btn-outline'}`}
            onClick={() => setSelectedMarkerType('problem')}
          >
            <FiAlertTriangle className="w-4 h-4" />
            {t('map.problems')} (
            {markers.filter((m) => m.type === 'problem').length})
          </button>
        </div>
      )}

      {/* Карта */}
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body p-0">
          <div
            ref={mapRef}
            style={{ height, width: '100%' }}
            className="rounded-lg"
          />
          {!leafletLoaded && (
            <div className="absolute inset-0 flex items-center justify-center bg-base-200 rounded-lg">
              <div className="flex flex-col items-center gap-2">
                <span className="loading loading-spinner loading-lg"></span>
                <span>{t('map.loading')}</span>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Легенда */}
      {showControls && (
        <div className="card bg-base-100 shadow-sm">
          <div className="card-body">
            <h3 className="card-title text-sm">{t('map.legend')}</h3>
            <div className="grid grid-cols-2 md:grid-cols-5 gap-2 text-xs">
              <div className="flex items-center gap-2">
                <div className="w-4 h-4 rounded-full bg-blue-500"></div>
                <span>{t('map.sender')}</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-4 h-4 rounded-full bg-green-500"></div>
                <span>{t('map.receiver')}</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-4 h-4 rounded-full bg-yellow-500"></div>
                <span>{t('map.in_transit')}</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-4 h-4 rounded-full bg-green-500"></div>
                <span>{t('map.delivered')}</span>
              </div>
              <div className="flex items-center gap-2">
                <div className="w-4 h-4 rounded-full bg-red-500"></div>
                <span>{t('map.problem')}</span>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
