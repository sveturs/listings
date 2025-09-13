import React, {
  useState,
  useEffect,
  useRef,
  useMemo,
  useCallback,
} from 'react';
import useMobileOptimization from '@/hooks/useMobileOptimization';

interface VirtualizedListProps<T> {
  items: T[];
  itemHeight: number; // Высота одного элемента в пикселях
  containerHeight: number; // Высота контейнера
  renderItem: (item: T, index: number) => React.ReactNode;
  keyExtractor?: (item: T, index: number) => string;
  overscan?: number; // Количество элементов для предзагрузки за пределами viewport
  onEndReached?: () => void; // Callback для подгрузки данных
  onEndReachedThreshold?: number; // Когда вызывать onEndReached (в пикселях до конца)
  className?: string;
  loadingComponent?: React.ReactNode;
  emptyComponent?: React.ReactNode;
}

function VirtualizedList<T>({
  items,
  itemHeight,
  containerHeight,
  renderItem,
  keyExtractor = (item, index) => index.toString(),
  overscan = 5,
  onEndReached,
  onEndReachedThreshold = 100,
  className = '',
  loadingComponent,
  emptyComponent,
}: VirtualizedListProps<T>) {
  const [scrollTop, setScrollTop] = useState(0);
  const [isScrolling, setIsScrolling] = useState(false);
  const scrollElementRef = useRef<HTMLDivElement>(null);
  const scrollTimeout = useRef<NodeJS.Timeout>();

  const { settings, throttledScrollHandler } = useMobileOptimization();

  // Вычисляем видимый диапазон
  const visibleRange = useMemo(() => {
    const start = Math.max(0, Math.floor(scrollTop / itemHeight) - overscan);
    const visibleCount = Math.ceil(containerHeight / itemHeight);
    const end = Math.min(items.length - 1, start + visibleCount + overscan);

    return { start, end };
  }, [scrollTop, itemHeight, containerHeight, overscan, items.length]);

  // Видимые элементы
  const visibleItems = useMemo(() => {
    return items
      .slice(visibleRange.start, visibleRange.end + 1)
      .map((item, index) => ({
        item,
        index: visibleRange.start + index,
      }));
  }, [items, visibleRange.start, visibleRange.end]);

  // Общая высота всех элементов
  const totalHeight = items.length * itemHeight;

  // Offset для первого видимого элемента
  const offsetY = visibleRange.start * itemHeight;

  // Обработчик скролла с throttling
  const handleScroll = throttledScrollHandler(
    useCallback(() => {
      const scrollElement = scrollElementRef.current;
      if (!scrollElement) return;

      const newScrollTop = scrollElement.scrollTop;
      setScrollTop(newScrollTop);

      // Проверяем достижение конца списка
      if (onEndReached) {
        const distanceFromEnd = totalHeight - (newScrollTop + containerHeight);
        if (distanceFromEnd <= onEndReachedThreshold) {
          onEndReached();
        }
      }

      // Устанавливаем состояние скроллинга
      setIsScrolling(true);
      if (scrollTimeout.current) {
        clearTimeout(scrollTimeout.current);
      }
      scrollTimeout.current = setTimeout(() => {
        setIsScrolling(false);
      }, 150);
    }, [onEndReached, totalHeight, containerHeight, onEndReachedThreshold])
  );

  // Добавляем обработчик скролла
  useEffect(() => {
    const scrollElement = scrollElementRef.current;
    if (!scrollElement) return;

    scrollElement.addEventListener('scroll', handleScroll, { passive: true });

    return () => {
      scrollElement.removeEventListener('scroll', handleScroll);
      if (scrollTimeout.current) {
        clearTimeout(scrollTimeout.current);
      }
    };
  }, [handleScroll]);

  // Если список пуст
  if (items.length === 0) {
    return (
      <div
        className={`flex items-center justify-center ${className}`}
        style={{ height: containerHeight }}
      >
        {emptyComponent || (
          <div className="text-gray-500 text-center">
            <div className="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4 mx-auto">
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
                  d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2M4 13h2m13-4v4H7V9z"
                />
              </svg>
            </div>
            <p className="text-lg font-medium mb-2">Список пуст</p>
            <p className="text-sm">Элементы не найдены</p>
          </div>
        )}
      </div>
    );
  }

  return (
    <div
      ref={scrollElementRef}
      className={`overflow-auto ${className} ${isScrolling ? 'pointer-events-none' : ''}`}
      style={{
        height: containerHeight,
        // Оптимизация для мобильных
        WebkitOverflowScrolling: 'touch',
        scrollBehavior: isScrolling ? 'auto' : 'smooth',
      }}
    >
      {/* Контейнер для всех элементов */}
      <div style={{ height: totalHeight, position: 'relative' }}>
        {/* Видимые элементы */}
        <div style={{ transform: `translateY(${offsetY}px)` }}>
          {visibleItems.map(({ item, index }) => (
            <div
              key={keyExtractor(item, index)}
              style={{
                height: itemHeight,
                // Добавляем класс для оптимизации
                willChange: isScrolling ? 'transform' : 'auto',
              }}
              className="w-full"
            >
              {renderItem(item, index)}
            </div>
          ))}
        </div>

        {/* Индикатор загрузки внизу списка */}
        {loadingComponent && onEndReached && (
          <div
            style={{
              position: 'absolute',
              bottom: 0,
              left: 0,
              right: 0,
              height: 60,
              display: 'flex',
              alignItems: 'center',
              justifyContent: 'center',
            }}
          >
            {loadingComponent}
          </div>
        )}
      </div>

      {/* Индикатор скроллинга для лучшего UX на мобильных */}
      {isScrolling && settings.enableLazyLoading && (
        <div className="fixed top-1/2 right-4 transform -translate-y-1/2 bg-black/70 text-white px-2 py-1 rounded text-xs font-medium z-50">
          {Math.round(
            (scrollTop / Math.max(1, totalHeight - containerHeight)) * 100
          )}
          %
        </div>
      )}
    </div>
  );
}

export default VirtualizedList;
