 // frontend/hostel-frontend/src/components/maps/FullscreenMap.js
import React, { useEffect, useRef } from 'react';
import { Paper } from '@mui/material';
import 'leaflet/dist/leaflet.css';
import L from 'leaflet';
import '../maps/leaflet-icons'; // Импортируем файл с фиксом иконок

const FullscreenMap = ({ latitude, longitude, title, markers = [] }) => {
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
      
      // Если есть список маркеров, добавляем их
      if (markers && markers.length > 0) {
        const markerGroup = L.featureGroup();
        
        markers.forEach(marker => {
          const leafletMarker = L.marker([marker.latitude, marker.longitude]);
          
          // Если у маркера есть всплывающая подсказка или заголовок
          if (marker.tooltip) {
            leafletMarker.bindTooltip(marker.tooltip);
          }
          
          if (marker.title) {
            leafletMarker.bindPopup(marker.title);
          }
          
          leafletMarker.addTo(mapRef.current);
          markerGroup.addLayer(leafletMarker);
        });
        
        // Масштабируем карту, чтобы видеть все маркеры
        if (markers.length > 1) {
          mapRef.current.fitBounds(markerGroup.getBounds(), {
            padding: [50, 50],
            maxZoom: 15
          });
        }
      } else if (title) {
        // Если маркеров нет, но есть центральная точка с заголовком
        const marker = L.marker([latitude, longitude]).addTo(mapRef.current);
        marker.bindPopup(title);
      }
    } else {
      // Если карта уже существует, обновляем маркеры
      // Сначала очищаем все существующие маркеры
      mapRef.current.eachLayer(layer => {
        if (layer instanceof L.Marker) {
          mapRef.current.removeLayer(layer);
        }
      });
      
      // Затем добавляем новые маркеры
      if (markers && markers.length > 0) {
        const markerGroup = L.featureGroup();
        
        markers.forEach(marker => {
          const leafletMarker = L.marker([marker.latitude, marker.longitude]);
          
          if (marker.tooltip) {
            leafletMarker.bindTooltip(marker.tooltip);
          }
          
          if (marker.title) {
            leafletMarker.bindPopup(marker.title);
          }
          
          leafletMarker.addTo(mapRef.current);
          markerGroup.addLayer(leafletMarker);
        });
        
        // Масштабируем карту, чтобы видеть все маркеры
        if (markers.length > 1) {
          mapRef.current.fitBounds(markerGroup.getBounds(), {
            padding: [50, 50],
            maxZoom: 15
          });
        } else {
          // Просто центрируем на указанной точке
          mapRef.current.setView([latitude, longitude], 15);
        }
      } else {
        // Просто центрируем карту на указанной точке
        mapRef.current.setView([latitude, longitude], 15);
        
        // Если есть заголовок, добавляем маркер с ним
        if (title) {
          const marker = L.marker([latitude, longitude]).addTo(mapRef.current);
          marker.bindPopup(title);
        }
      }
    }
    
    // Очистка при размонтировании компонента
    return () => {
      if (mapRef.current) {
        mapRef.current.remove();
        mapRef.current = null;
      }
    };
  }, [latitude, longitude, title, markers]);
  
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