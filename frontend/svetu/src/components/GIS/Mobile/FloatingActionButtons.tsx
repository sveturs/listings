import React, { useState } from 'react';
import { useTranslations } from 'next-intl';

interface FloatingActionButtonsProps {
  onSearchClick: () => void;
  onFiltersClick: () => void;
  onGeolocationClick: () => void;
  onShowAllClick: () => void;
  markersCount: number;
  isLoading?: boolean;
  hasFilters?: boolean;
}

const FloatingActionButtons: React.FC<FloatingActionButtonsProps> = ({
  onSearchClick,
  onFiltersClick,
  onGeolocationClick,
  onShowAllClick,
  markersCount,
  isLoading = false,
  hasFilters = false,
}) => {
  const t = useTranslations('map');
  const [isExpanded, setIsExpanded] = useState(false);

  const toggleExpanded = () => {
    setIsExpanded(!isExpanded);
  };

  return (
    <>
      {/* Компактный поиск сверху */}
      <div className="fixed top-16 left-4 right-4 z-20 md:hidden">
        <div className="bg-white/95 backdrop-blur-sm rounded-full shadow-lg px-4 py-3 flex items-center gap-3">
          <svg
            className="w-5 h-5 text-gray-400 flex-shrink-0"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
          <button
            onClick={onSearchClick}
            className="flex-1 text-left text-gray-500 text-sm truncate"
          >
            {t('search.placeholder', 'Введите адрес или район...')}
          </button>
          {isLoading && (
            <div className="animate-spin rounded-full h-4 w-4 border-b-2 border-blue-600 flex-shrink-0"></div>
          )}
        </div>
      </div>

      {/* Кнопки действий справа */}
      <div className="fixed top-32 right-4 flex flex-col items-end gap-3 z-30 md:hidden">
        {/* Основная кнопка меню */}
        <button
          onClick={toggleExpanded}
          className={`bg-white shadow-xl rounded-full p-4 hover:shadow-2xl transition-all duration-300 ${
            isExpanded ? 'rotate-45' : ''
          } ${isLoading ? 'animate-pulse' : ''}`}
          style={{
            boxShadow: '0 8px 32px rgba(0, 0, 0, 0.12)',
          }}
          aria-label={t('filters.title', 'Меню')}
        >
          <svg
            className="w-6 h-6 text-gray-700 transition-transform duration-300"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M12 6v6m0 0v6m0-6h6m-6 0H6"
            />
          </svg>
        </button>

        {/* Вторичные кнопки */}
        <div
          className={`flex flex-col gap-3 transition-all duration-300 ${
            isExpanded
              ? 'opacity-100 translate-y-0'
              : 'opacity-0 translate-y-4 pointer-events-none'
          }`}
        >
          {/* Кнопка фильтров с индикатором */}
          <button
            onClick={() => {
              onFiltersClick();
              setIsExpanded(false);
            }}
            className={`relative rounded-full p-3 shadow-lg hover:shadow-xl transition-all duration-200 ${
              hasFilters
                ? 'bg-orange-600 hover:bg-orange-700 text-white'
                : 'bg-white hover:bg-gray-50 text-gray-700'
            }`}
            title={t('filters.title', 'Фильтры')}
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.414A1 1 0 013 6.707V4z"
              />
            </svg>
            {hasFilters && (
              <div className="absolute -top-1 -right-1 w-3 h-3 bg-red-500 rounded-full animate-pulse" />
            )}
          </button>

          {/* Кнопка геолокации */}
          <button
            onClick={() => {
              onGeolocationClick();
              setIsExpanded(false);
            }}
            className="bg-green-600 hover:bg-green-700 text-white rounded-full p-3 shadow-lg hover:shadow-xl transition-all duration-200"
            title={t('geolocation.findMe', 'Моё местоположение')}
          >
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
              />
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
              />
            </svg>
          </button>

          {/* Кнопка "Показать все" с счетчиком */}
          {markersCount > 0 && (
            <button
              onClick={() => {
                onShowAllClick();
                setIsExpanded(false);
              }}
              className="bg-purple-600 hover:bg-purple-700 text-white rounded-full p-3 shadow-lg hover:shadow-xl transition-all duration-200 relative"
              title={t('showAll.title', 'Показать все результаты')}
            >
              <svg
                className="w-5 h-5"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"
                />
              </svg>
              {markersCount > 0 && (
                <div className="absolute -top-2 -right-2 bg-red-500 text-white text-xs font-bold rounded-full w-6 h-6 flex items-center justify-center">
                  {markersCount > 99 ? '99+' : markersCount}
                </div>
              )}
            </button>
          )}
        </div>
      </div>

      {/* Кнопка результатов внизу экрана */}
      {markersCount > 0 && (
        <button
          onClick={onShowAllClick}
          className="fixed bottom-20 left-4 right-4 bg-blue-600 hover:bg-blue-700 text-white py-4 px-6 rounded-full shadow-xl hover:shadow-2xl transition-all duration-300 flex items-center justify-center gap-3 font-medium z-20 md:hidden"
        >
          <svg
            className="w-5 h-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"
            />
          </svg>
          <span>
            {t('results.showCount', { count: markersCount })}
          </span>
        </button>
      )}
    </>
  );
};

export default FloatingActionButtons;