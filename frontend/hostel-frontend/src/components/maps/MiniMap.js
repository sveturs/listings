// frontend/hostel-frontend/src/components/maps/SimpleMiniMap.js
import React, { useEffect, useRef } from 'react';
import { Box, IconButton } from '@mui/material';
import { Maximize2 } from 'lucide-react';
import L from 'leaflet';
import 'leaflet/dist/leaflet.css';

// Исправляем проблему с маркерами Leaflet в React
delete L.Icon.Default.prototype._getIconUrl;
L.Icon.Default.mergeOptions({
  iconRetinaUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon-2x.png',
  iconUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png',
  shadowUrl: 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png',
});

const SimpleMiniMap = ({ latitude, longitude, title = 'Местоположение', onExpand }) => {
  const mapDivRef = useRef(null);

  useEffect(() => {
    // Проверяем, что координаты есть и DOM-элемент готов
    if (!latitude || !longitude || !mapDivRef.current) return;

    // Очищаем содержимое div, чтобы избежать дублирования карты
    mapDivRef.current.innerHTML = '';

    console.log('Initializing simple mini map...');
    
    // Создаем новую карту
    const map = L.map(mapDivRef.current).setView([latitude, longitude], 14);
    
    // Добавляем слой тайлов
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
      attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors'
    }).addTo(map);
    
    // Добавляем маркер
    L.marker([latitude, longitude])
      .addTo(map)
      .bindPopup(title);
    
    // Сохраняем ссылку на карту для будущей очистки
    const mapInstance = map;
    
    // Функция очистки при размонтировании компонента
    return () => {
      if (mapInstance) {
        console.log('Cleaning up simple mini map...');
        mapInstance.remove();
      }
    };
  }, [latitude, longitude, title]);

  return (
    <Box sx={{ width: '100%', height: 200, borderRadius: '4px', overflow: 'hidden', position: 'relative' }}>
      <div 
        ref={mapDivRef}
        style={{ width: '100%', height: '100%' }}
      />
      
      {onExpand && (
        <IconButton
          onClick={onExpand}
          sx={{
            position: 'absolute',
            top: 8,
            right: 8,
            bgcolor: 'background.paper',
            '&:hover': {
              bgcolor: 'background.paper',
            },
            zIndex: 1000 // Убедимся, что кнопка отображается поверх карты
          }}
        >
          <Maximize2 size={20} />
        </IconButton>
      )}
    </Box>
  );
};

export default SimpleMiniMap;