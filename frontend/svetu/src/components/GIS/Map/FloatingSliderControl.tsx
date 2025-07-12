import React, { useState, useRef, useEffect } from 'react';

interface FloatingSliderControlProps {
  mode: 'radius' | 'walking';
  onModeChange: (mode: 'radius' | 'walking') => void;
  walkingTime: number; // в минутах (5-60)
  onWalkingTimeChange: (time: number) => void;
  searchRadius: number; // в метрах (500-50000)
  onRadiusChange: (radius: number) => void;
  isFullscreen?: boolean;
  isMobile?: boolean;
  translations?: {
    walkingAccessibility: string;
    searchRadius: string;
    minutes: string;
    km: string;
    m: string;
    holdForSettings: string;
    singleClickHint: string;
    mobileHint: string;
    desktopHint: string;
  };
}

const FloatingSliderControl: React.FC<FloatingSliderControlProps> = ({
  mode,
  onModeChange,
  walkingTime,
  onWalkingTimeChange,
  searchRadius,
  onRadiusChange,
  isFullscreen = false,
  isMobile = false,
  translations,
}) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [lastTapTime, setLastTapTime] = useState(0);
  const [showTooltip, setShowTooltip] = useState(false);
  const [tempWalkingTime, setTempWalkingTime] = useState(walkingTime);
  const [tempSearchRadius, setTempSearchRadius] = useState(searchRadius);
  const longPressTimerRef = useRef<NodeJS.Timeout | null>(null);
  const firstUseRef = useRef(true);

  // Форматирование значения для компактного отображения
  const getCompactValue = () => {
    const t = translations || {
      minutes: 'min',
      km: 'km',
      m: 'm'
    };

    if (mode === 'walking') {
      return `${walkingTime}'`;
    } else {
      if (searchRadius >= 1000) {
        const km = (searchRadius / 1000).toFixed(0);
        return `${km}${t.km}`;
      }
      return `${searchRadius}${t.m}`;
    }
  };

  // Форматирование значения для полного отображения
  const getDisplayValue = () => {
    const t = translations || {
      minutes: 'min',
      km: 'km',
      m: 'm'
    };

    if (mode === 'walking') {
      return `${displayWalkingTime} ${t.minutes}`;
    } else {
      if (displayRadius >= 1000) {
        return `${(displayRadius / 1000).toFixed(1)} ${t.km}`;
      }
      return `${displayRadius} ${t.m}`;
    }
  };

  // Обработка одиночного клика
  const handleSingleTap = () => {
    const newMode = mode === 'walking' ? 'radius' : 'walking';
    onModeChange(newMode);
  };

  // Обработка двойного клика
  const handleDoubleTap = () => {
    setIsExpanded(!isExpanded);
  };

  // Обработка клика на иконку
  const handleIconClick = (e: React.MouseEvent) => {
    e.preventDefault();
    e.stopPropagation();

    const currentTime = Date.now();
    const timeDiff = currentTime - lastTapTime;

    if (timeDiff < 300 && lastTapTime > 0) {
      // Двойной клик - открываем/закрываем слайдер
      console.log('[FloatingSliderControl] Double click detected');
      handleDoubleTap();
      setLastTapTime(0); // сбрасываем чтобы избежать повторных срабатываний
    } else {
      // Одиночный клик - ждем чтобы убедиться что это не двойной
      setLastTapTime(currentTime);
      setTimeout(() => {
        // Проверяем, что после клика не было второго клика
        if (Date.now() - currentTime >= 300) {
          console.log('[FloatingSliderControl] Single click detected');
          handleSingleTap();
        }
      }, 300);
    }
  };

  // Long press для мобильных
  const handleTouchStart = (e: React.TouchEvent) => {
    e.preventDefault();

    longPressTimerRef.current = setTimeout(() => {
      handleDoubleTap();
      // Вибрация на поддерживающих устройствах
      if ('vibrate' in navigator) {
        navigator.vibrate(50);
      }
    }, 500);
  };

  const handleTouchEnd = (e: React.TouchEvent) => {
    e.preventDefault();

    if (longPressTimerRef.current) {
      clearTimeout(longPressTimerRef.current);
    }
  };

  const handleTouchMove = () => {
    if (longPressTimerRef.current) {
      clearTimeout(longPressTimerRef.current);
    }
  };

  // Обработка изменения слайдера - только обновляем временные значения
  const handleSliderChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const percent = parseFloat(e.target.value);

    if (mode === 'walking') {
      const minutes = Math.round(5 + (percent / 100) * (60 - 5));
      setTempWalkingTime(minutes);
    } else {
      const meters = Math.round(500 + (percent / 100) * (50000 - 500));
      setTempSearchRadius(meters);
    }
  };

  // Обработка отпускания слайдера - применяем изменения
  const handleSliderChangeEnd = () => {
    if (mode === 'walking' && tempWalkingTime !== walkingTime) {
      onWalkingTimeChange(tempWalkingTime);
    } else if (mode === 'radius' && tempSearchRadius !== searchRadius) {
      onRadiusChange(tempSearchRadius);
    }
  };

  // Закрытие по клику вне контрола
  useEffect(() => {
    const handleClickOutside = (e: MouseEvent) => {
      const target = e.target as HTMLElement;
      if (!target.closest('.floating-slider-control')) {
        setIsExpanded(false);
      }
    };

    if (isExpanded) {
      document.addEventListener('click', handleClickOutside);
      return () => document.removeEventListener('click', handleClickOutside);
    }
  }, [isExpanded]);

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Проверяем, что фокус не на input элементах
      if (
        e.target instanceof HTMLInputElement ||
        e.target instanceof HTMLTextAreaElement
      ) {
        return;
      }

      switch (e.key.toLowerCase()) {
        case 'r':
          if (mode !== 'radius') {
            onModeChange('radius');
          }
          break;
        case 'w':
          if (mode !== 'walking') {
            onModeChange('walking');
          }
          break;
        case ' ':
          e.preventDefault();
          setIsExpanded(!isExpanded);
          break;
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [mode, isExpanded, onModeChange]);

  // Показать подсказку для мобильных при первом использовании
  useEffect(() => {
    if (isMobile && firstUseRef.current) {
      const timer = setTimeout(() => {
        setShowTooltip(true);
        setTimeout(() => setShowTooltip(false), 3000);
        firstUseRef.current = false;
      }, 1000);
      return () => clearTimeout(timer);
    }
  }, [isMobile]);

  // Синхронизация временных значений с реальными при изменении режима
  useEffect(() => {
    setTempWalkingTime(walkingTime);
    setTempSearchRadius(searchRadius);
  }, [mode, walkingTime, searchRadius]);

  const icon = mode === 'walking' ? '🚶' : '📏';
  const color = mode === 'walking' ? '#10B981' : '#3B82F6';

  // Позиционирование под нативными контролами Mapbox
  // NavigationControl (88px) + GeolocateControl (32px) + FullscreenControl (32px) + отступы = ~195px
  const positionClasses = 'top-[195px] right-[10px]';

  // Конвертируем значения в проценты для слайдера
  const displayRadius = mode === 'radius' ? tempSearchRadius : searchRadius;
  const displayWalkingTime = mode === 'walking' ? tempWalkingTime : walkingTime;
  const radiusPercent = ((displayRadius - 500) / (50000 - 500)) * 100;
  const walkingPercent = ((displayWalkingTime - 5) / (60 - 5)) * 100;
  const currentPercent = mode === 'walking' ? walkingPercent : radiusPercent;

  return (
    <div
      className={`floating-slider-control fixed ${positionClasses} z-10 transition-all duration-300 ease-in-out`}
      style={{
        width: isExpanded ? (isMobile ? '260px' : '300px') : '29px',
        height: isExpanded ? 'auto' : '29px',
      }}
    >
        {!isExpanded ? (
          // Компактное состояние - стиль как у нативного контрола
          <button
            className="w-full h-full flex items-center justify-center bg-white rounded cursor-pointer transition-all duration-200 hover:bg-gray-50"
            style={{
              boxShadow: '0 1px 4px rgba(0,0,0,0.3)',
              border: 'none',
              padding: 0,
              outline: 'none',
              WebkitTapHighlightColor: 'transparent',
            }}
            onClick={handleIconClick}
            onTouchStart={handleTouchStart}
            onTouchEnd={handleTouchEnd}
            onTouchMove={handleTouchMove}
            title={translations?.singleClickHint || "Click - change mode, double click - expand"}
          >
            <span style={{ fontSize: '18px' }}>{icon}</span>
            
            {/* Индикатор значения */}
            <div
              className="absolute text-white font-semibold rounded shadow-sm"
              style={{
                backgroundColor: color,
                fontSize: '9px',
                lineHeight: '1',
                padding: '1px 4px',
                bottom: '-2px',
                right: '-2px',
                boxShadow: '0 1px 2px rgba(0,0,0,0.3)',
              }}
            >
              {getCompactValue()}
            </div>

            {/* Подсказка для мобильных */}
            {showTooltip && isMobile && (
              <div className="absolute -bottom-8 left-0 bg-gray-800 text-white text-xs px-2 py-1 rounded shadow-lg whitespace-nowrap">
                {translations?.holdForSettings || "Hold for settings"}
                <div className="absolute -top-1 left-4 w-2 h-2 bg-gray-800 transform rotate-45"></div>
              </div>
            )}
          </button>
        ) : (
          // Развернутое состояние
          <div 
            className="bg-white rounded-lg shadow-lg p-3 space-y-3"
            style={{
              boxShadow: '0 2px 10px rgba(0,0,0,0.15)',
            }}
          >
            {/* Заголовок */}
            <div className="flex items-center justify-between">
              <div className="flex items-center gap-2">
                <span className="text-lg transition-all duration-200">
                  {icon}
                </span>
                <span className="text-sm font-medium text-gray-700">
                  {mode === 'walking' ? (translations?.walkingAccessibility || 'Walking accessibility') : (translations?.searchRadius || 'Search radius')}
                </span>
              </div>
              <button
                onClick={() => setIsExpanded(false)}
                className="text-gray-400 hover:text-gray-600 transition-colors p-1"
              >
                <svg
                  className="w-4 h-4"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M6 18L18 6M6 6l12 12"
                  />
                </svg>
              </button>
            </div>

            {/* Слайдер */}
            <div className="space-y-2">
              <input
                type="range"
                min="0"
                max="100"
                value={currentPercent}
                onChange={handleSliderChange}
                onMouseUp={handleSliderChangeEnd}
                onTouchEnd={handleSliderChangeEnd}
                className="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer slider"
                style={{
                  background: `linear-gradient(to right, ${color} 0%, ${color} ${currentPercent}%, #e5e7eb ${currentPercent}%, #e5e7eb 100%)`,
                }}
              />

              {/* Значение и подсказка */}
              <div className="flex justify-between items-center text-xs">
                <span className="text-gray-500">
                  {isMobile
                    ? (translations?.mobileHint || 'Tap icon to change mode')
                    : (translations?.desktopHint || 'R/W - change mode')}
                </span>
                <span className="font-semibold" style={{ color }}>
                  {getDisplayValue()}
                </span>
              </div>
            </div>

            {/* Кликабельная область для переключения режима */}
            <button
              className="absolute inset-x-0 top-0 h-10 cursor-pointer hover:bg-gray-50 hover:bg-opacity-50 transition-all duration-200 rounded-t-lg"
              onClick={handleSingleTap}
              aria-label="Переключить режим измерения"
            />
          </div>
        )}

      <style jsx>{`
        .slider::-webkit-slider-thumb {
          appearance: none;
          width: 16px;
          height: 16px;
          border-radius: 50%;
          background: ${color};
          cursor: pointer;
          border: 2px solid white;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        }

        .slider::-moz-range-thumb {
          width: 16px;
          height: 16px;
          border-radius: 50%;
          background: ${color};
          cursor: pointer;
          border: 2px solid white;
          box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        }
      `}</style>
    </div>
  );
};

export default FloatingSliderControl;
