// Путь: src/components/maps/MiniMap.tsx
import React, { useEffect, useRef, useState } from 'react';
import { Box, Paper, IconButton } from '@mui/material';
import { Maximize2, X } from 'lucide-react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import FullscreenMap from './FullscreenMap';

interface MiniMapProps {
  latitude?: string | number;
  longitude?: string | number;
  address?: string;
  onClick?: () => void;
  onExpand?: () => void;
}

interface Marker {
  latitude: number;
  longitude: number;
  title?: string;
}

const MiniMap: React.FC<MiniMapProps> = ({ latitude, longitude, address, onClick, onExpand }) => {
  const mapRef = useRef<HTMLDivElement | null>(null);
  const mapInstanceRef = useRef<L.Map | null>(null);
  // Добавляем состояние для открытия полноэкранной карты
  const [expandedMapOpen, setExpandedMapOpen] = useState<boolean>(false);

  useEffect(() => {
    // Проверяем валидность координат перед инициализацией карты
    if (!latitude || !longitude || isNaN(Number(latitude)) || isNaN(Number(longitude))) {
      console.log('Invalid map coordinates:', { latitude, longitude });
      return; // Пропускаем инициализацию карты, если координаты невалидны
    }

    // Добавляем небольшую задержку перед инициализацией карты, чтобы DOM успел полностью смонтироваться
    const initMapTimer = setTimeout(() => {
      // Конвертируем строковые координаты в числа для точности
      const lat = Number(latitude);
      const lng = Number(longitude);

      // Если карта уже инициализирована, удаляем её перед созданием новой
      if (mapInstanceRef.current) {
        mapInstanceRef.current.remove();
        mapInstanceRef.current = null;
      }

      // Проверяем, существует ли DOM-элемент
      if (!mapRef.current) {
        console.log('Map container element not found');
        return;
      }

      try {
        // Создаем карту с защитой от ошибок и улучшенной производительностью
        const map = L.map(mapRef.current, {
          center: [lat, lng],
          zoom: 13,
          scrollWheelZoom: false,
          dragging: false,
          touchZoom: false,
          doubleClickZoom: false,
          boxZoom: false,
          keyboard: false,
          zoomControl: false,
          attributionControl: false,
          fadeAnimation: true,           // Плавные переходы прозрачности
          zoomAnimation: true,           // Плавное масштабирование
          markerZoomAnimation: true,     // Плавная анимация маркеров
          preferCanvas: true,            // Использовать Canvas для лучшей производительности
          renderer: L.canvas()            // Явно указываем Canvas renderer
        });
        mapInstanceRef.current = map;

        // Добавляем OpenStreetMap тайлы с оптимизированными параметрами
        L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
          maxZoom: 19,
          attribution: '© OpenStreetMap contributors',
          subdomains: 'abc',            // Распределение запросов
          updateWhenIdle: true,         // Обновление тайлов только в состоянии покоя
          updateWhenZooming: false,     // Отключаем обновление при масштабировании
          keepBuffer: 2                 // Буфер тайлов
        }).addTo(map);

        // Добавляем маркер
        L.marker([lat, lng]).addTo(map);
      } catch (error) {
        console.error('Error initializing map:', error);
      }
    }, 100); // задержка в 100мс

    // Правильно обрабатываем очистку при размонтировании
    return () => {
      clearTimeout(initMapTimer);
      if (mapInstanceRef.current) {
        try {
          mapInstanceRef.current.remove();
        } catch (err) {
          console.error('Error removing map:', err);
        }
        mapInstanceRef.current = null;
      }
    };
  }, [latitude, longitude]); // Пересоздаем карту только при изменении координат

  // Функция для открытия полноэкранной карты
  const handleExpandMap = (): void => {
    setExpandedMapOpen(true);
    if (onExpand) onExpand();
  };

  // Функция для закрытия полноэкранной карты
  const handleCloseExpandedMap = (): void => {
    setExpandedMapOpen(false);
  };

  // Если нет координат, не рендерим карту
  if (!latitude || !longitude || isNaN(Number(latitude)) || isNaN(Number(longitude))) {
    return null;
  }

  return (
    <Paper elevation={1} sx={{ overflow: 'hidden', borderRadius: 1 }}>
      <Box
        sx={{
          position: 'relative',
          height: 200,
          width: '100%'
        }}
      >
        {/* Кнопка "Развернуть" */}
        <IconButton
          onClick={handleExpandMap}
          sx={{
            position: 'absolute',
            top: 8,
            right: 8,
            bgcolor: 'background.paper',
            '&:hover': {
              bgcolor: 'background.paper',
            },
            zIndex: 1000,
            boxShadow: '0 2px 6px rgba(0,0,0,0.1)'
          }}
          size="small"
        >
          <Maximize2 size={20} />
        </IconButton>
        
        <Box
          ref={mapRef}
          onClick={() => onClick && onClick()}
          sx={{
            height: '100%',
            width: '100%',
            cursor: onClick ? 'pointer' : 'default'
          }}
        />
      </Box>
      
      {address && (
        <Box sx={{ p: 1, fontSize: '0.875rem', color: 'text.secondary' }}>
          {address}
        </Box>
      )}

      {/* Полноэкранная карта */}
      {expandedMapOpen && (
        <Box
          sx={{
            position: 'fixed',
            top: 0,
            left: 0,
            right: 0,
            bottom: 0,
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            bgcolor: 'rgba(0, 0, 0, 0.5)',
            zIndex: 1300,
            p: 3
          }}
          onClick={handleCloseExpandedMap}
        >
          <Paper
            sx={{
              width: '90%',
              maxWidth: 1200,
              maxHeight: '90vh',
              bgcolor: 'background.paper',
              borderRadius: 2,
              boxShadow: 24,
              position: 'relative',
              overflow: 'hidden'
            }}
            onClick={(e) => e.stopPropagation()}
          >
            <Box sx={{ position: 'absolute', top: 8, right: 8, zIndex: 1050 }}>
              <IconButton
                onClick={handleCloseExpandedMap}
                sx={{
                  bgcolor: 'background.paper',
                  '&:hover': {
                    bgcolor: 'background.paper',
                  },
                  boxShadow: '0 2px 6px rgba(0,0,0,0.2)'
                }}
              >
                <X size={20} />
              </IconButton>
            </Box>

            <FullscreenMap
              latitude={Number(latitude)}
              longitude={Number(longitude)}
              title={address || "Местоположение"}
              markers={[
                {
                  latitude: Number(latitude),
                  longitude: Number(longitude),
                  title: address || "Местоположение"
                }
              ]}
            />
          </Paper>
        </Box>
      )}
    </Paper>
  );
};

export default MiniMap;