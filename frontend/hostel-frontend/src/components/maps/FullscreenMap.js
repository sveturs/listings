// frontend/hostel-frontend/src/components/maps/FullscreenMap.js
import React, { useEffect, useRef } from 'react';
import { Paper } from '@mui/material';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import '../maps/leaflet-icons'; // Импортируем файл с фиксом иконок

const FullscreenMap = ({ latitude, longitude, title }) => {
  // Перемещаем объявление хуков до условной проверки
  const mapContainerRef = useRef(null);
  const mapRef = useRef(null);
  
  // Если координат нет, компонент всё равно отрендерится но карта не будет инициализирована
  const hasCoordinates = latitude && longitude;
  
  useEffect(() => {
    // Проверяем условия внутри хука
    if (!hasCoordinates || !mapContainerRef.current) return;
    
    // Инициализируем карту только если её еще нет
    if (!mapRef.current) {
      mapRef.current = L.map(mapContainerRef.current).setView([latitude, longitude], 15);
      
      // Добавляем слой тайлов OpenStreetMap
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '&copy; <a href="https://www.openstreetmap.org/copyright">OpenStreetMap</a> contributors',
        maxZoom: 19
      }).addTo(mapRef.current);
      
      // Добавляем маркер
      const marker = L.marker([latitude, longitude]).addTo(mapRef.current);
      if (title) {
        marker.bindPopup(title);
      }
    } else {
      // Если карта уже существует, просто обновляем её центр
      mapRef.current.setView([latitude, longitude], 15);
    }
    
    // Очистка при размонтировании компонента
    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
  }, [latitude, longitude, title]);
  
  return (
    <Paper
      sx={{
        position: 'relative',
        width: '100%',
        maxWidth: 1200,
        maxHeight: '90vh',
        overflow: 'hidden'
      }}
    >
      {hasCoordinates ? (
        <div 
          ref={mapContainerRef}
          style={{ width: '100%', height: '80vh' }}
        />
      ) : (
        <div style={{ width: '100%', height: '80vh', display: 'flex', 
                     alignItems: 'center', justifyContent: 'center', 
                     backgroundColor: '#f5f5f5' }}>
          Координаты отсутствуют
        </div>
      )}
    </Paper>
  );
};

export default FullscreenMap;