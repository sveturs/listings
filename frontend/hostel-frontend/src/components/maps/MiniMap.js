
// Путь: src/components/maps/MiniMap.js

import React, { useEffect, useRef } from 'react';
import { Box, Paper } from '@mui/material';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

const MiniMap = ({ latitude, longitude, address }) => {
  const mapRef = useRef(null);
  const mapInstanceRef = useRef(null);

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

  // Если нет координат, не рендерим карту
  if (!latitude || !longitude || isNaN(Number(latitude)) || isNaN(Number(longitude))) {
    return null;
  }

  return (
    <Paper elevation={1} sx={{ overflow: 'hidden', borderRadius: 1 }}>
      <Box
        ref={mapRef}
        sx={{
          height: 200,
          width: '100%'
        }}
      />
      {address && (
        <Box sx={{ p: 1, fontSize: '0.875rem', color: 'text.secondary' }}>
          {address}
        </Box>
      )}
    </Paper>
  );
};

export default MiniMap;