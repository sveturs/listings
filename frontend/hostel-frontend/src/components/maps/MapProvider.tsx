// frontend/hostel-frontend/src/components/maps/MapProvider.ts
import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import 'leaflet/dist/leaflet.css';

// Определяем интерфейс для контекста карты
interface MapContextType {
  isLoaded: boolean;
}

// Создаем контекст для карты
const MapContext = createContext<MapContextType | null>(null);

interface MapProviderProps {
  children: ReactNode;
}

export const MapProvider: React.FC<MapProviderProps> = ({ children }) => {
  const [isLoaded, setIsLoaded] = useState<boolean>(false);

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
      const preloadImage = (url: string): void => {
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

export const useMap = (): MapContextType => {
  const context = useContext(MapContext);
  if (!context) {
    throw new Error('useMap must be used within a MapProvider');
  }
  return context;
};

export default MapProvider;