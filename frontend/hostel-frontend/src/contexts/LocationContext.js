// frontend/hostel-frontend/src/contexts/LocationContext.js
import React, { createContext, useState, useContext, useEffect } from 'react';
import axios from '../api/axios';

// Значения по умолчанию для Белграда
const defaultLocation = {
  city: 'Белград',
  country: 'Сербия',
  lat: 44.8178,
  lon: 20.4570,
};

const LocationContext = createContext();

export const useLocation = () => useContext(LocationContext);

export const LocationProvider = ({ children }) => {
  const [userLocation, setUserLocation] = useState(() => {
    // Пытаемся загрузить из localStorage
    const saved = localStorage.getItem('userLocation');
    if (saved) {
      try {
        return JSON.parse(saved);
      } catch (e) {
        console.error('Failed to parse saved location', e);
      }
    }
    return null;
  });
  
  const [locationAsked, setLocationAsked] = useState(() => {
    return localStorage.getItem('locationAsked') === 'true';
  });

  const [locationDismissed, setLocationDismissed] = useState(() => {
    return localStorage.getItem('locationDismissed') === 'true';
  });

  // Сохраняем в localStorage при изменении
  useEffect(() => {
    if (userLocation) {
      localStorage.setItem('userLocation', JSON.stringify(userLocation));
    }
  }, [userLocation]);

  // Отмечаем, что уже запрашивали местоположение
  useEffect(() => {
    localStorage.setItem('locationAsked', String(locationAsked));
  }, [locationAsked]);

  useEffect(() => {
    localStorage.setItem('locationDismissed', String(locationDismissed));
  }, [locationDismissed]);

  // Функция для определения местоположения
  const detectLocation = async () => {
    setLocationAsked(true);
    
    if (!navigator.geolocation) {
      console.log('Geolocation is not supported by this browser');
      setUserLocation(defaultLocation);
      return;
    }
    
    try {
      const position = await new Promise((resolve, reject) => {
        navigator.geolocation.getCurrentPosition(resolve, reject, {
          enableHighAccuracy: true,
          timeout: 5000,
          maximumAge: 0
        });
      });
      
      const { latitude, longitude } = position.coords;
      
      // Получаем город по координатам
      const response = await axios.get('/api/v1/geocode/reverse', {
        params: { lat: latitude, lon: longitude }
      });
      
      if (response.data?.data) {
        setUserLocation(response.data.data);
      } else {
        setUserLocation(defaultLocation);
      }
    } catch (error) {
      console.error('Error getting location:', error);
      setUserLocation(defaultLocation);
    }
  };

  // Функция для изменения города вручную
  const setCity = (location) => {
    setUserLocation(location);
    setLocationDismissed(true);
    
    // Генерируем событие смены города, чтобы все компоненты могли отреагировать
    const cityChangeEvent = new CustomEvent('cityChanged', { 
      detail: { 
        lat: location.lat, 
        lon: location.lon,
        city: location.city,
        country: location.country 
      } 
    });
    window.dispatchEvent(cityChangeEvent);
  };

  // Функция для отклонения предложения изменить город
  const dismissLocationSuggestion = () => {
    setLocationDismissed(true);
  };

  // Проверяем, нужно ли запросить местоположение при первой загрузке
  useEffect(() => {
    if (!userLocation && !locationAsked) {
      detectLocation();
    } else if (!userLocation && locationAsked) {
      // Если пользователь отказался ранее, используем местоположение по умолчанию
      setUserLocation(defaultLocation);
    }
  }, [userLocation, locationAsked]);

  const value = {
    userLocation: userLocation || defaultLocation,
    setCity,
    detectLocation,
    locationAsked,
    locationDismissed,
    dismissLocationSuggestion
  };

  return (
    <LocationContext.Provider value={value}>
      {children}
    </LocationContext.Provider>
  );
};

export default LocationContext;