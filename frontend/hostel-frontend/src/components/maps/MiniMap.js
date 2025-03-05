// Путь: src/components/maps/MiniMap.js
import React, { useEffect, useRef, useState } from 'react';
import { Box, Paper, IconButton, Modal } from '@mui/material';
import { Maximize2, X } from 'lucide-react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';
import FullscreenMap from './FullscreenMap';

const MiniMap = ({ latitude, longitude, address }) => {
  const mapRef = useRef(null);
  const mapInstanceRef = useRef(null);
  // Добавляем состояние для открытия полноэкранной карты
  const [expandedMapOpen, setExpandedMapOpen] = useState(false);

  useEffect(() => {
    // Проверяем валидность координат перед инициализацией карты
    if (!latitude || !longitude || isNaN(Number(latitude)) || isNaN(Number(longitude))) {
      console.log('Invalid map coordinates:', { latitude, longitude });
      return; // Пропускаем инициализацию карты, если координаты невалидны
    }
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
      // Создаем карту с защитой от ошибок
      const map = L.map(mapRef.current, {
        center: [lat, lng],
        zoom: 13,
        scrollWheelZoom: false,
        dragging: false,
        touchZoom: false,
        doubleClickZoom: false,
        boxZoom: false,
        tap: false,
        keyboard: false,
        zoomControl: false,
        attributionControl: false
      });
      mapInstanceRef.current = map;
      // Добавляем OpenStreetMap тайлы
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        maxZoom: 19,
        attribution: '© OpenStreetMap contributors'
      }).addTo(map);
      // Добавляем маркер
      L.marker([lat, lng]).addTo(map);
      // Правильно обрабатываем очистку при размонтировании
      return () => {
        if (mapInstanceRef.current) {
          mapInstanceRef.current.remove();
          mapInstanceRef.current = null;
        }
      };
    } catch (error) {
      console.error('Error initializing map:', error);
    }
  }, [latitude, longitude]); // Пересоздаем карту только при изменении координат

  // Функция для открытия полноэкранной карты
  const handleExpandMap = () => {
    setExpandedMapOpen(true);
  };

  // Функция для закрытия полноэкранной карты
  const handleCloseExpandedMap = () => {
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
          sx={{
            height: '100%',
            width: '100%'
          }}
        />
      </Box>
      
      {address && (
        <Box sx={{ p: 1, fontSize: '0.875rem', color: 'text.secondary' }}>
          {address}
        </Box>
      )}

      {/* Модальное окно с полноэкранной картой */}
      <Modal
        open={expandedMapOpen}
        onClose={handleCloseExpandedMap}
        aria-labelledby="expanded-map-modal"
        aria-describedby="expanded-map-view"
        sx={{
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          p: 3
        }}
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
      </Modal>
    </Paper>
  );
};

export default MiniMap;