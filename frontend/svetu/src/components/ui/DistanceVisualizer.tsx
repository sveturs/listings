'use client';

import React from 'react';
import { MapPin, Navigation, Clock, Car, PersonStanding } from 'lucide-react';

interface DistanceVisualizerProps {
  distance: number; // в километрах
  address?: string;
  className?: string;
  showMap?: boolean;
  showTravelTime?: boolean;
}

export const DistanceVisualizer: React.FC<DistanceVisualizerProps> = ({
  distance,
  address,
  className = '',
  showMap = true,
  showTravelTime = true,
}) => {
  // Расчет времени в пути
  const getTravelTime = (km: number) => {
    const walkingSpeed = 5; // км/ч
    const drivingSpeed = 30; // км/ч в городе
    
    const walkingMinutes = Math.round((km / walkingSpeed) * 60);
    const drivingMinutes = Math.round((km / drivingSpeed) * 60);
    
    return { walking: walkingMinutes, driving: drivingMinutes };
  };

  const { walking, driving } = getTravelTime(distance);

  // Визуальное представление расстояния
  const getVisualRadius = (km: number): number => {
    // Масштабирование для визуализации (макс 100px радиус)
    if (km <= 1) return 30;
    if (km <= 5) return 50;
    if (km <= 10) return 70;
    if (km <= 20) return 85;
    return 100;
  };

  const radius = getVisualRadius(distance);

  return (
    <div className={`space-y-3 ${className}`}>
      {/* Визуальная карта */}
      {showMap && (
        <div className="relative h-32 bg-base-200 rounded-lg overflow-hidden">
          {/* Круги расстояния */}
          <svg className="absolute inset-0 w-full h-full">
            {/* Внешние круги для масштаба */}
            <circle
              cx="50%"
              cy="50%"
              r="90"
              fill="none"
              stroke="currentColor"
              strokeWidth="1"
              className="text-base-300"
              strokeDasharray="2 2"
            />
            <circle
              cx="50%"
              cy="50%"
              r="60"
              fill="none"
              stroke="currentColor"
              strokeWidth="1"
              className="text-base-300"
              strokeDasharray="2 2"
            />
            <circle
              cx="50%"
              cy="50%"
              r="30"
              fill="none"
              stroke="currentColor"
              strokeWidth="1"
              className="text-base-300"
              strokeDasharray="2 2"
            />
            
            {/* Радиус расстояния */}
            <circle
              cx="50%"
              cy="50%"
              r={radius}
              fill="currentColor"
              className="text-primary/20"
            />
            <circle
              cx="50%"
              cy="50%"
              r={radius}
              fill="none"
              stroke="currentColor"
              strokeWidth="2"
              className="text-primary"
            />
          </svg>
          
          {/* Центральная точка (вы) */}
          <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2">
            <div className="relative">
              <div className="w-3 h-3 bg-primary rounded-full animate-ping absolute"></div>
              <div className="w-3 h-3 bg-primary rounded-full relative"></div>
            </div>
            <span className="absolute top-4 left-1/2 -translate-x-1/2 text-xs font-medium whitespace-nowrap">
              Вы здесь
            </span>
          </div>
          
          {/* Точка назначения */}
          <div 
            className="absolute"
            style={{
              top: `${50 - (radius * 0.7)}%`,
              left: `${50 + (radius * 0.5)}%`,
              transform: 'translate(-50%, -50%)'
            }}
          >
            <MapPin className="w-5 h-5 text-error" />
          </div>
          
          {/* Легенда расстояния */}
          <div className="absolute top-2 right-2 bg-base-100/90 backdrop-blur-sm rounded px-2 py-1">
            <span className="text-xs font-medium">
              {distance < 1 ? `${Math.round(distance * 1000)} м` : `${distance.toFixed(1)} км`}
            </span>
          </div>
        </div>
      )}
      
      {/* Информация о расстоянии */}
      <div className="space-y-2">
        {/* Адрес */}
        {address && (
          <div className="flex items-start gap-2">
            <Navigation className="w-4 h-4 text-base-content/60 mt-0.5" />
            <p className="text-sm text-base-content/80 line-clamp-2">{address}</p>
          </div>
        )}
        
        {/* Время в пути */}
        {showTravelTime && (
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1.5">
              <PersonStanding className="w-4 h-4 text-base-content/60" />
              <span className="text-sm">
                {walking < 60 ? `${walking} мин` : `${Math.round(walking / 60)} ч`}
              </span>
            </div>
            
            {distance > 2 && (
              <div className="flex items-center gap-1.5">
                <Car className="w-4 h-4 text-base-content/60" />
                <span className="text-sm">
                  {driving < 60 ? `${driving} мин` : `${Math.round(driving / 60)} ч`}
                </span>
              </div>
            )}
          </div>
        )}
        
        {/* Статус близости */}
        <div className="flex items-center gap-2">
          <div className={`w-2 h-2 rounded-full ${
            distance <= 1 ? 'bg-success' :
            distance <= 5 ? 'bg-info' :
            distance <= 15 ? 'bg-warning' :
            'bg-base-300'
          }`}></div>
          <span className="text-sm font-medium">
            {distance <= 0.5 ? 'Очень близко' :
             distance <= 1 ? 'В шаговой доступности' :
             distance <= 3 ? 'Рядом с вами' :
             distance <= 5 ? 'Недалеко' :
             distance <= 10 ? 'В вашем районе' :
             'Далеко'}
          </span>
        </div>
      </div>
    </div>
  );
};