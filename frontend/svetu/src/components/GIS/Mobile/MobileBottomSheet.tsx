import React, { useState, useRef, useEffect, useCallback } from 'react';
import { MapMarkerData } from '@/components/GIS/types/gis';
import { useTranslations } from 'next-intl';
import useMobileOptimization from '@/hooks/useMobileOptimization';

interface MobileBottomSheetProps {
  isOpen: boolean;
  onClose: () => void;
  markers: MapMarkerData[];
  isLoading?: boolean;
  onMarkerClick?: (marker: MapMarkerData) => void;
}

type SheetState = 'collapsed' | 'peek' | 'expanded';

const MobileBottomSheet: React.FC<MobileBottomSheetProps> = ({
  isOpen,
  onClose,
  markers,
  isLoading = false,
  onMarkerClick,
}) => {
  const t = useTranslations('map');
  const { optimizeImageUrl, settings } = useMobileOptimization();
  const [sheetState, setSheetState] = useState<SheetState>('peek');
  const [startY, setStartY] = useState(0);
  const [currentY, setCurrentY] = useState(0);
  const [isDragging, setIsDragging] = useState(false);
  const [isMounted, setIsMounted] = useState(false);
  const sheetRef = useRef<HTMLDivElement>(null);

  // Устанавливаем флаг монтирования
  useEffect(() => {
    setIsMounted(true);
  }, []);

  // Получаем высоту в пикселях
  const getSheetHeight = useCallback(
    (state: SheetState) => {
      if (!isMounted) return 0;
      // Высоты для разных состояний (в процентах от высоты экрана)
      const SHEET_HEIGHTS = {
        collapsed: 0,
        peek: 20, // 20% от высоты экрана
        expanded: 85, // 85% от высоты экрана
      };
      return (window.innerHeight * SHEET_HEIGHTS[state]) / 100;
    },
    [isMounted]
  );

  // Обработчики touch событий
  const handleTouchStart = useCallback(
    (e: TouchEvent) => {
      if (!isOpen) return;
      setStartY(e.touches[0].clientY);
      setCurrentY(e.touches[0].clientY);
      setIsDragging(true);
    },
    [isOpen]
  );

  const handleTouchMove = useCallback(
    (e: TouchEvent) => {
      if (!isDragging) return;
      e.preventDefault();
      const touchY = e.touches[0].clientY;
      setCurrentY(touchY);
    },
    [isDragging]
  );

  const handleTouchEnd = useCallback(() => {
    if (!isDragging) return;
    setIsDragging(false);

    const deltaY = currentY - startY;
    const threshold = 50; // Минимальное расстояние для изменения состояния

    if (Math.abs(deltaY) < threshold) return;

    if (deltaY > 0) {
      // Свайп вниз
      if (sheetState === 'expanded') {
        setSheetState('peek');
      } else if (sheetState === 'peek') {
        setSheetState('collapsed');
        setTimeout(() => onClose(), 300);
      }
    } else {
      // Свайп вверх
      if (sheetState === 'peek') {
        setSheetState('expanded');
      } else if (sheetState === 'collapsed') {
        setSheetState('peek');
      }
    }
  }, [isDragging, currentY, startY, sheetState, onClose]);

  // Добавляем обработчики событий
  useEffect(() => {
    const sheet = sheetRef.current;
    if (!sheet) return;

    sheet.addEventListener('touchstart', handleTouchStart, { passive: false });
    sheet.addEventListener('touchmove', handleTouchMove, { passive: false });
    sheet.addEventListener('touchend', handleTouchEnd);

    return () => {
      sheet.removeEventListener('touchstart', handleTouchStart);
      sheet.removeEventListener('touchmove', handleTouchMove);
      sheet.removeEventListener('touchend', handleTouchEnd);
    };
  }, [handleTouchStart, handleTouchMove, handleTouchEnd]);

  // Устанавливаем состояние peek при открытии
  useEffect(() => {
    if (isOpen && sheetState === 'collapsed') {
      setSheetState('peek');
    } else if (!isOpen) {
      setSheetState('collapsed');
    }
  }, [isOpen, sheetState]);

  // Вычисляем текущую высоту с учетом перетаскивания
  const getCurrentHeight = () => {
    const baseHeight = getSheetHeight(sheetState);
    if (!isDragging) return baseHeight;

    const dragOffset = startY - currentY;
    const newHeight = baseHeight + dragOffset;

    // Ограничиваем высоту в пределах допустимых значений
    const minHeight = getSheetHeight('peek');
    const maxHeight = getSheetHeight('expanded');

    return Math.max(minHeight, Math.min(maxHeight, newHeight));
  };

  if (!isOpen && sheetState === 'collapsed') return null;

  return (
    <>
      {/* Backdrop - только для expanded состояния */}
      {sheetState === 'expanded' && (
        <div
          className="fixed inset-0 bg-black/30 z-40 md:hidden"
          onClick={() => setSheetState('peek')}
        />
      )}

      {/* Bottom Sheet */}
      <div
        ref={sheetRef}
        className={`fixed bottom-0 left-0 right-0 bg-white rounded-t-2xl shadow-2xl z-50 md:hidden transform transition-transform duration-300 ease-out ${
          sheetState === 'collapsed' ? 'translate-y-full' : 'translate-y-0'
        }`}
        style={{
          height: `${getCurrentHeight()}px`,
          transform: `translateY(${sheetState === 'collapsed' ? '100%' : '0'})`,
        }}
      >
        {/* Ручка для перетаскивания */}
        <div className="w-full flex justify-center pt-3 pb-2">
          <div className="w-10 h-1 bg-gray-300 rounded-full" />
        </div>

        {/* Заголовок */}
        <div className="px-4 pb-3 border-b border-gray-200">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <h3 className="text-lg font-semibold text-gray-900">
                {t('results.title')}
              </h3>
              <div className="bg-blue-100 text-blue-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
                {isLoading ? '...' : markers.length}
              </div>
            </div>
            {sheetState === 'expanded' && (
              <button
                onClick={() => setSheetState('peek')}
                className="p-1 hover:bg-gray-100 rounded-full transition-colors"
              >
                <svg
                  className="w-5 h-5 text-gray-500"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M19 9l-7 7-7-7"
                  />
                </svg>
              </button>
            )}
          </div>

          {sheetState === 'peek' && (
            <p className="text-sm text-gray-500 mt-1">{t('results.swipeUp')}</p>
          )}
        </div>

        {/* Контент */}
        <div className="flex-1 overflow-y-auto">
          {isLoading ? (
            <div className="flex items-center justify-center py-8">
              <div className="flex items-center gap-3">
                <div className="animate-spin rounded-full h-5 w-5 border-b-2 border-blue-600"></div>
                <span className="text-gray-600">{t('common.loading')}</span>
              </div>
            </div>
          ) : markers.length === 0 ? (
            <div className="flex flex-col items-center justify-center py-12">
              <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
                <svg
                  className="w-8 h-8 text-gray-400"
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
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">
                {t('results.empty.title')}
              </h3>
              <p className="text-gray-500 text-center">
                {t('results.empty.description')}
              </p>
            </div>
          ) : (
            <div className="px-4 max-h-96 overflow-y-auto">
              {markers
                .slice(0, settings.maxMarkersCount)
                .map((marker, index) => (
                  <div
                    key={`${marker.id}-${index}`}
                    className="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow cursor-pointer mb-3"
                    onClick={() => onMarkerClick?.(marker)}
                  >
                    <div className="flex gap-3">
                      {/* Изображение или иконка */}
                      <div className="flex-shrink-0">
                        {marker.imageUrl ? (
                          <img
                            src={optimizeImageUrl(marker.imageUrl, 64, 64)}
                            alt={marker.title}
                            className="w-16 h-16 object-cover rounded-lg"
                            loading="lazy"
                            decoding="async"
                          />
                        ) : (
                          <div className="w-16 h-16 bg-gray-100 rounded-lg flex items-center justify-center">
                            <span className="text-2xl">
                              {marker.metadata?.icon || '📦'}
                            </span>
                          </div>
                        )}
                      </div>

                      {/* Информация */}
                      <div className="flex-1 min-w-0">
                        <h4 className="text-base font-medium text-gray-900 truncate mb-1">
                          {marker.title}
                        </h4>

                        {marker.metadata?.price && (
                          <p className="text-lg font-semibold text-blue-600 mb-1">
                            {new Intl.NumberFormat('sr-RS').format(
                              marker.metadata.price
                            )}{' '}
                            {marker.metadata.currency || 'RSD'}
                          </p>
                        )}

                        {marker.metadata?.category && (
                          <p className="text-sm text-gray-500 mb-1">
                            {marker.metadata.category}
                          </p>
                        )}

                        {marker.data?.address && (
                          <p className="text-sm text-gray-400 truncate">
                            📍 {marker.data.address}
                          </p>
                        )}

                        {/* Дополнительная информация */}
                        <div className="flex items-center gap-3 mt-2 text-xs text-gray-400">
                          {marker.data?.views_count && (
                            <span>👁 {marker.data.views_count}</span>
                          )}
                          {marker.data?.rating && (
                            <span>⭐ {marker.data.rating}</span>
                          )}
                        </div>
                      </div>

                      {/* Стрелка */}
                      <div className="flex-shrink-0 flex items-center">
                        <svg
                          className="w-5 h-5 text-gray-400"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M9 5l7 7-7 7"
                          />
                        </svg>
                      </div>
                    </div>
                  </div>
                ))}
            </div>
          )}

          {/* Показываем сообщение если маркеров больше чем лимит */}
          {markers.length > settings.maxMarkersCount &&
            sheetState === 'expanded' && (
              <div className="px-4 pb-2">
                <div className="bg-amber-50 border border-amber-200 rounded-lg p-3">
                  <div className="flex items-center gap-2">
                    <svg
                      className="w-4 h-4 text-amber-600 flex-shrink-0"
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
                    <p className="text-sm text-amber-800">
                      Показано {settings.maxMarkersCount} из {markers.length}{' '}
                      результатов.
                      <span className="font-medium">
                        {' '}
                        Уточните фильтры для лучшей производительности.
                      </span>
                    </p>
                  </div>
                </div>
              </div>
            )}
        </div>

        {/* Кнопки действий в peek режиме */}
        {sheetState === 'peek' && markers.length > 0 && (
          <div className="px-4 py-3 border-t border-gray-200 bg-gray-50">
            <button
              onClick={() => setSheetState('expanded')}
              className="w-full bg-blue-600 text-white py-3 rounded-lg font-medium hover:bg-blue-700 transition-colors"
            >
              {t('results.viewAll')} ({markers.length})
            </button>
          </div>
        )}
      </div>
    </>
  );
};

export default MobileBottomSheet;
