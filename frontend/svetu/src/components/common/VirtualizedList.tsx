import React, { useState, useRef, useCallback } from 'react';

export interface VirtualizedListProps<T> {
  items: T[];
  itemHeight: number;
  containerHeight: number;
  renderItem: (item: T, index: number) => React.ReactNode;
  onLoadMore?: () => void;
  hasMore?: boolean;
  loading?: boolean;
  threshold?: number;
  className?: string;
}

export function VirtualizedList<T>({
  items,
  itemHeight,
  containerHeight,
  renderItem,
  onLoadMore,
  hasMore = false,
  loading = false,
  threshold = 5,
  className = '',
}: VirtualizedListProps<T>) {
  const [scrollTop, setScrollTop] = useState(0);
  const containerRef = useRef<HTMLDivElement>(null);

  // Calculate visible range
  const visibleCount = Math.ceil(containerHeight / itemHeight);
  const bufferSize = Math.max(5, Math.ceil(visibleCount / 2));

  const startIndex = Math.max(
    0,
    Math.floor(scrollTop / itemHeight) - bufferSize
  );
  const endIndex = Math.min(
    items.length,
    Math.ceil((scrollTop + containerHeight) / itemHeight) + bufferSize
  );

  // Virtual items to render
  const visibleItems = items.slice(startIndex, endIndex);

  // Total height for scrollbar
  const totalHeight = items.length * itemHeight;

  // Offset for visible items
  const offsetY = startIndex * itemHeight;

  const handleScroll = useCallback(
    (event: React.UIEvent<HTMLDivElement>) => {
      const scrollTop = event.currentTarget.scrollTop;
      setScrollTop(scrollTop);

      // Check if we need to load more data
      if (onLoadMore && hasMore && !loading) {
        const scrollBottom = scrollTop + containerHeight;
        const shouldLoadMore =
          totalHeight - scrollBottom < itemHeight * threshold;

        if (shouldLoadMore) {
          onLoadMore();
        }
      }
    },
    [
      onLoadMore,
      hasMore,
      loading,
      containerHeight,
      totalHeight,
      itemHeight,
      threshold,
    ]
  );

  return (
    <div
      ref={containerRef}
      className={`overflow-auto ${className}`}
      style={{ height: containerHeight }}
      onScroll={handleScroll}
    >
      <div style={{ height: totalHeight, position: 'relative' }}>
        <div
          style={{
            transform: `translateY(${offsetY}px)`,
            position: 'absolute',
            top: 0,
            left: 0,
            right: 0,
          }}
        >
          {visibleItems.map((item, index) => (
            <div
              key={startIndex + index}
              style={{ height: itemHeight }}
              className="flex items-center"
            >
              {renderItem(item, startIndex + index)}
            </div>
          ))}

          {loading && (
            <div className="flex justify-center items-center py-4">
              <span className="loading loading-spinner loading-md"></span>
              <span className="ml-2">Загрузка...</span>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
