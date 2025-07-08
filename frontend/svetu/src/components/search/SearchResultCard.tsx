'use client';

import { ReactNode } from 'react';
import { useBehaviorTracking } from '@/hooks/useBehaviorTracking';

interface SearchResultCardProps {
  children: ReactNode;
  searchQuery: string;
  itemId: string;
  position: number;
  totalResults: number;
  searchStartTime: number;
  className?: string;
  onClick?: () => void;
}

/**
 * Обертка для карточек результатов поиска с интегрированным трекингом
 * Автоматически отслеживает клики по результатам поиска
 */
export default function SearchResultCard({
  children,
  searchQuery,
  itemId,
  position,
  totalResults,
  searchStartTime,
  className = '',
  onClick,
}: SearchResultCardProps) {
  const { trackResultClicked } = useBehaviorTracking();

  const handleClick = async () => {
    // Трекинг клика по результату поиска
    if (searchQuery) {
      try {
        await trackResultClicked({
          search_query: searchQuery,
          clicked_item_id: itemId,
          click_position: position,
          total_results: totalResults,
          click_time_from_search_ms: Date.now() - searchStartTime,
        });
      } catch (error) {
        console.error('Failed to track result click:', error);
      }
    }

    // Вызываем переданный onClick, если есть
    if (onClick) {
      onClick();
    }
  };

  return (
    <div className={className} onClick={handleClick}>
      {children}
    </div>
  );
}
