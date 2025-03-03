// frontend/hostel-frontend/src/components/maps/MapProvider.js
import React, { createContext, useContext, useState, useEffect } from 'react';
import 'leaflet/dist/leaflet.css';

// Создаем контекст для карты
const MapContext = createContext(null);

export const MapProvider = ({ children }) => {
  const [isLoaded, setIsLoaded] = useState(false);

  // Эффект для подтверждения, что Leaflet доступен в браузере
  useEffect(() => {
    // В случае с Leaflet нам не нужно ожидать загрузки внешних скриптов,
    // поскольку они импортируются как npm-пакеты
    
    // Проверяем, что мы в браузере (не на сервере)
    if (typeof window !== 'undefined') {
      // Отмечаем карты как загруженные
      setIsLoaded(true);
      
      // Необходимо загрузить маркеры Leaflet
      const iconRetinaUrl = 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon-2x.png';
      const iconUrl = 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-icon.png';
      const shadowUrl = 'https://unpkg.com/leaflet@1.7.1/dist/images/marker-shadow.png';
      
      // Предзагрузка изображений
      const preloadImage = (url) => {
        const img = new Image();
        img.src = url;
      };
      
      preloadImage(iconRetinaUrl);
      preloadImage(iconUrl);
      preloadImage(shadowUrl);
      
      // Добавляем стили для Leaflet в head документа, если они еще не добавлены
      if (!document.querySelector('link[href*="leaflet.css"]')) {
        const linkElement = document.createElement('link');
        linkElement.rel = 'stylesheet';
        linkElement.href = 'https://unpkg.com/leaflet@1.7.1/dist/leaflet.css';
        
        document.head.appendChild(linkElement);
      }
    }
  }, []);

  return (
    <MapContext.Provider value={{ isLoaded }}>
      {children}
    </MapContext.Provider>
  );
};

export const useMap = () => {
  const context = useContext(MapContext);
  if (!context) {
    throw new Error('useMap must be used within a MapProvider');
  }
  return context;
};

export default MapProvider;