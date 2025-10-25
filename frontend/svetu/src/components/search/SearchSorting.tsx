'use client';

import { useTranslations } from 'next-intl';
import { useBehaviorTracking } from '@/hooks/useBehaviorTracking';

interface SearchSortingProps {
  currentSort: string;
  searchQuery: string;
  resultsCount: number;
  onSortChange: (sort: string) => void;
  variant?: 'select' | 'buttons';
  className?: string;
}

/**
 * Компонент сортировки поиска с интегрированным трекингом
 * Автоматически отслеживает изменения сортировки
 */
export default function SearchSorting({
  currentSort,
  searchQuery,
  resultsCount,
  onSortChange,
  variant = 'select',
  className = '',
}: SearchSortingProps) {
  const t = useTranslations('search');
  const { trackSearchSortChanged } = useBehaviorTracking();

  const sortOptions = [
    { value: 'relevance', label: t('relevance') },
    { value: 'price', label: t('price') },
    { value: 'date', label: t('date') },
    { value: 'popularity', label: t('popularity') },
  ];

  const handleSortChange = async (newSort: string) => {
    const previousSort = currentSort;

    // Трекинг изменения сортировки
    if (searchQuery && newSort !== previousSort) {
      try {
        await trackSearchSortChanged({
          search_query: searchQuery,
          sort_type: newSort,
          previous_sort: previousSort,
          results_count: resultsCount,
        });
      } catch (error) {
        console.error('Failed to track sort change:', error);
      }
    }

    // Применяем изменение сортировки
    onSortChange(newSort);
  };

  if (variant === 'buttons') {
    return (
      <div className={`flex flex-wrap gap-2 ${className}`}>
        {sortOptions.map((option) => (
          <button
            key={option.value}
            className={`btn btn-xs sm:btn-sm ${
              currentSort === option.value
                ? 'btn-primary shadow-lg'
                : 'btn-ghost hover:btn-primary hover:btn-outline'
            } transition-all duration-200`}
            onClick={() => handleSortChange(option.value)}
          >
            {option.value === 'price' && (
              <svg
                className="w-4 h-4 mr-1"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            )}
            {option.value === 'date' && (
              <svg
                className="w-4 h-4 mr-1"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"
                />
              </svg>
            )}
            {option.value === 'popularity' && (
              <svg
                className="w-4 h-4 mr-1"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"
                />
              </svg>
            )}
            {option.value === 'relevance' && (
              <svg
                className="w-4 h-4 mr-1"
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
            )}
            {option.label}
          </button>
        ))}
      </div>
    );
  }

  // Select variant
  return (
    <div className={className}>
      <label className="label">
        <span className="label-text font-medium flex items-center gap-2">
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
              d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12"
            />
          </svg>
          {t('sortBy')}
        </span>
      </label>
      <select
        className="select select-bordered w-full bg-base-100"
        value={currentSort}
        onChange={(e) => handleSortChange(e.target.value)}
        aria-label={t('sortBy')}
      >
        {sortOptions.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
    </div>
  );
}
