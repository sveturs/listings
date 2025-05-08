// frontend/hostel-frontend/src/contexts/LocationContext.tsx
import React, { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import axios from '../api/axios';

// Определяем типы для геолокации
export interface LocationData {
  city: string;
  country: string;
  lat: number;
  lon: number;
  [key: string]: any; // Для дополнительных свойств
}

export interface LocationContextType {
  userLocation: LocationData | null;
  setCity: (cityData: LocationData) => void;
  isGeolocating: boolean;
  detectUserLocation: () => Promise<LocationData>;
  locationDismissed: boolean;
  dismissLocationSuggestion: () => void;
}

interface LocationProviderProps {
  children: ReactNode;
}

// Создаем контекст с null в качестве начального значения
const LocationContext = createContext<LocationContextType | null>(null);

export function useLocation(): LocationContextType {
  const context = useContext(LocationContext);
  if (!context) {
    throw new Error('useLocation must be used within a LocationProvider');
  }
  return context;
}

export function LocationProvider({ children }: LocationProviderProps): React.ReactNode {
  const [userLocation, setUserLocation] = useState<LocationData | null>(null);
  const [isGeolocating, setIsGeolocating] = useState<boolean>(false);
  const [locationDismissed, setLocationDismissed] = useState<boolean>(
    localStorage.getItem('locationDismissed') === 'true'
  );

  // Загружаем данные из localStorage при инициализации
  useEffect(() => {
    const savedLocation = localStorage.getItem('userLocation');
    if (savedLocation) {
      try {
        setUserLocation(JSON.parse(savedLocation));
      } catch (e) {
        console.error('Ошибка при загрузке данных местоположения:', e);
        localStorage.removeItem('userLocation');
      }
    }
  }, []);

  // Функция для получения геолокации пользователя
  const detectUserLocation = (): Promise<LocationData> => {
    if (!navigator.geolocation) {
      console.log('Геолокация не поддерживается вашим браузером');
      return Promise.reject('Геолокация не поддерживается');
    }

    setIsGeolocating(true);
    
    return new Promise<LocationData>((resolve, reject) => {
      navigator.geolocation.getCurrentPosition(
        async (position) => {
          const { latitude, longitude } = position.coords;
          
          try {
            // Получаем информацию о городе на основе координат
            const response = await axios.get('/api/v1/geocode/reverse', {
              params: { lat: latitude, lon: longitude }
            });
            
            if (response.data?.data) {
              const locationData: LocationData = {
                city: response.data.data.city,
                country: response.data.data.country,
                lat: latitude,
                lon: longitude
              };
              
              // Сохраняем в localStorage
              localStorage.setItem('userLocation', JSON.stringify(locationData));
              setUserLocation(locationData);
              resolve(locationData);
            } else {
              reject('Не удалось определить город по координатам');
            }
          } catch (error) {
            console.error('Ошибка получения данных о городе:', error);
            reject(error);
          } finally {
            setIsGeolocating(false);
          }
        },
        (error) => {
          console.error('Ошибка геолокации:', error);
          setIsGeolocating(false);
          reject(error);
        },
        { 
          enableHighAccuracy: true,
          timeout: 10000,
          maximumAge: 0
        }
      );
    });
  };

  // Функция для установки города
  const setCity = (cityData: LocationData): void => {
    // Проверяем наличие необходимых данных
    if (!cityData || !cityData.lat || !cityData.lon) {
      console.error('Недостаточно данных для установки города:', cityData);
      return;
    }
    
    const locationData: LocationData = {
      city: cityData.city,
      country: cityData.country,
      lat: cityData.lat,
      lon: cityData.lon
    };
    
    // Сохраняем в localStorage
    localStorage.setItem('userLocation', JSON.stringify(locationData));
    
    // Обновляем состояние
    setUserLocation(locationData);
    
    // Генерируем событие для уведомления других компонентов
    const cityChangeEvent = new CustomEvent('cityChanged', {
      detail: {
        lat: locationData.lat,
        lon: locationData.lon,
        city: locationData.city,
        country: locationData.country
      }
    });
    window.dispatchEvent(cityChangeEvent);
    
    console.log(`Выбран город: ${locationData.city}, координаты: ${locationData.lat}, ${locationData.lon}`);
  };

  const dismissLocationSuggestion = (): void => {
    localStorage.setItem('locationDismissed', 'true');
    setLocationDismissed(true);
  };

  const value: LocationContextType = {
    userLocation,
    setCity,
    isGeolocating,
    detectUserLocation,
    locationDismissed,
    dismissLocationSuggestion
  };

  return (
    <LocationContext.Provider value={value}>
      {children}
    </LocationContext.Provider>
  );
}