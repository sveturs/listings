import React, { useState, useRef, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { SearchBar } from '@/components/SearchBar';

interface MobileSearchProps {
  isOpen: boolean;
  onClose: () => void;
  onSearch: (query: string) => void;
  searchQuery: string;
  _onSearchChange?: (query: string) => void;
  isSearching?: boolean;
  recentSearches?: string[];
  onRecentSearchClick?: (query: string) => void;
  onClearRecentSearches?: () => void;
}

const MobileSearch: React.FC<MobileSearchProps> = ({
  isOpen,
  onClose,
  onSearch,
  searchQuery,
  _onSearchChange,
  isSearching = false,
  recentSearches = [],
  onRecentSearchClick,
  onClearRecentSearches,
}) => {
  const t = useTranslations('map');
  const [_inputFocused, _setInputFocused] = useState(false);
  const searchRef = useRef<HTMLInputElement>(null);

  // Автофокус при открытии
  useEffect(() => {
    if (isOpen && searchRef.current) {
      setTimeout(() => {
        searchRef.current?.focus();
      }, 300); // Задержка для завершения анимации
    }
  }, [isOpen]);

  // Популярные поисковые запросы
  const popularSearches = [
    { query: 'Нови Београд', icon: '🏢' },
    { query: 'Земун', icon: '🏘️' },
    { query: 'Врачар', icon: '🏛️' },
    { query: 'Савски венац', icon: '🌊' },
    { query: 'Звездара', icon: '⭐' },
    { query: 'Нови Сад', icon: '🌉' },
  ];

  const quickActions = [
    {
      title: t('search.quickActions.nearMe'),
      icon: '📍',
      action: 'geolocation',
    },
    {
      title: t('search.quickActions.center'),
      icon: '🏛️',
      action: 'center',
    },
    {
      title: t('search.quickActions.newDistricts'),
      icon: '🏢',
      action: 'new-districts',
    },
  ];

  const handleQuickAction = (action: string) => {
    switch (action) {
      case 'geolocation':
        // Триггерим геолокацию через родительский компонент
        onSearch('geolocation');
        break;
      case 'center':
        onSearch('Београд центар');
        break;
      case 'new-districts':
        onSearch('Нови Београд');
        break;
    }
    onClose();
  };

  if (!isOpen) return null;

  return (
    <>
      {/* Backdrop */}
      <div
        className="fixed inset-0 bg-black/50 z-[60] md:hidden"
        onClick={onClose}
      />

      {/* Search Modal */}
      <div className="fixed inset-x-0 top-0 bg-white z-[61] md:hidden animate-slideInFromTop">
        {/* Header */}
        <div className="flex items-center gap-3 p-4 border-b border-gray-200 bg-white">
          <button
            onClick={onClose}
            className="p-2 hover:bg-gray-100 rounded-full transition-colors"
          >
            <svg
              className="w-5 h-5 text-gray-600"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M15 19l-7-7 7-7"
              />
            </svg>
          </button>

          <div className="flex-1">
            <SearchBar
              initialQuery={searchQuery}
              onSearch={(query) => {
                onSearch(query);
                onClose();
              }}
              placeholder={t('search.placeholder')}
              className="w-full"
            />
          </div>

          {isSearching && (
            <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-600" />
          )}
        </div>

        {/* Content */}
        <div className="max-h-[calc(100vh-120px)] overflow-y-auto">
          {/* Быстрые действия */}
          <div className="p-4 border-b border-gray-100">
            <h3 className="text-sm font-semibold text-gray-700 mb-3">
              {t('search.quickActions.title')}
            </h3>
            <div className="grid grid-cols-1 gap-2">
              {quickActions.map((action, index) => (
                <button
                  key={index}
                  onClick={() => handleQuickAction(action.action)}
                  className="flex items-center gap-3 p-3 hover:bg-gray-50 rounded-lg transition-colors text-left"
                >
                  <span className="text-xl">{action.icon}</span>
                  <span className="text-gray-700">{action.title}</span>
                </button>
              ))}
            </div>
          </div>

          {/* Последние поиски */}
          {recentSearches.length > 0 && (
            <div className="p-4 border-b border-gray-100">
              <div className="flex items-center justify-between mb-3">
                <h3 className="text-sm font-semibold text-gray-700">
                  {t('search.recent.title')}
                </h3>
                <button
                  onClick={onClearRecentSearches}
                  className="text-xs text-blue-600 hover:text-blue-700"
                >
                  {t('search.recent.clear')}
                </button>
              </div>
              <div className="space-y-1">
                {recentSearches.slice(0, 5).map((query, index) => (
                  <button
                    key={index}
                    onClick={() => {
                      onRecentSearchClick?.(query);
                      onClose();
                    }}
                    className="flex items-center gap-3 w-full p-2 hover:bg-gray-50 rounded-lg transition-colors text-left"
                  >
                    <svg
                      className="w-4 h-4 text-gray-400 flex-shrink-0"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                      />
                    </svg>
                    <span className="text-gray-700 truncate">{query}</span>
                  </button>
                ))}
              </div>
            </div>
          )}

          {/* Популярные районы */}
          <div className="p-4">
            <h3 className="text-sm font-semibold text-gray-700 mb-3">
              {t('search.popular.title')}
            </h3>
            <div className="grid grid-cols-2 gap-2">
              {popularSearches.map((item, index) => (
                <button
                  key={index}
                  onClick={() => {
                    onSearch(item.query);
                    onClose();
                  }}
                  className="flex items-center gap-2 p-3 bg-gray-50 hover:bg-gray-100 rounded-lg transition-colors text-left"
                >
                  <span className="text-lg">{item.icon}</span>
                  <span className="text-sm text-gray-700 truncate">
                    {item.query}
                  </span>
                </button>
              ))}
            </div>
          </div>

          {/* Советы */}
          <div className="p-4 bg-blue-50 border-t border-blue-100">
            <div className="flex items-start gap-3">
              <div className="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center flex-shrink-0 mt-0.5">
                <svg
                  className="w-4 h-4 text-blue-600"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  />
                </svg>
              </div>
              <div>
                <h4 className="text-sm font-semibold text-blue-900 mb-1">
                  {t('search.tips.title')}
                </h4>
                <p className="text-xs text-blue-700 leading-relaxed">
                  {t('search.tips.description')}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <style jsx>{`
        @keyframes slideInFromTop {
          from {
            transform: translateY(-100%);
          }
          to {
            transform: translateY(0);
          }
        }

        .animate-slideInFromTop {
          animation: slideInFromTop 0.3s ease-out;
        }
      `}</style>
    </>
  );
};

export default MobileSearch;
