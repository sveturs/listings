'use client';

import { useRef } from 'react';
import { useTranslations } from 'next-intl';
import { SearchSuggestion, EnhancedSuggestion } from '@/services/unifiedSearch';
import Image from 'next/image';

interface SearchAutocompleteProps {
  suggestions: (SearchSuggestion | EnhancedSuggestion)[];
  searchHistory: string[];
  trendingSearches: string[];
  showSuggestions: boolean;
  selectedIndex: number;
  query: string;
  isLoading: boolean;
  onSelect: (value: string) => void;
  onCategorySelect?: (categoryId: number) => void;
  onProductSelect?: (productId: number) => void;
}

export default function SearchAutocomplete({
  suggestions,
  searchHistory,
  trendingSearches,
  showSuggestions,
  selectedIndex,
  query,
  isLoading,
  onSelect,
  onCategorySelect,
  onProductSelect,
}: SearchAutocompleteProps) {
  const t = useTranslations('search');
  const ref = useRef<HTMLDivElement>(null);

  if (!showSuggestions) return null;

  const renderSuggestion = (
    suggestion: SearchSuggestion | EnhancedSuggestion,
    index: number
  ) => {
    const isSelected = index === selectedIndex;
    const isEnhanced = 'product' in suggestion || 'metadata' in suggestion;

    return (
      <div
        key={`suggestion-${index}`}
        className={`px-4 py-2 cursor-pointer transition-colors ${
          isSelected ? 'bg-primary text-primary-content' : 'hover:bg-base-200'
        }`}
        onMouseEnter={() => {}}
        onClick={() => {
          if (
            suggestion.type === 'category' &&
            suggestion.category &&
            onCategorySelect
          ) {
            onCategorySelect(suggestion.category.id);
          } else if (suggestion.type === 'product' && onProductSelect) {
            // Обработка обоих форматов: SearchSuggestion и EnhancedSuggestion
            const productId =
              'product' in suggestion && suggestion.product
                ? suggestion.product.id
                : 'product_id' in suggestion
                  ? suggestion.product_id
                  : undefined;
            if (productId) {
              onProductSelect(productId);
            }
          } else {
            onSelect(suggestion.text);
          }
        }}
      >
        <div className="flex items-center gap-3">
          {/* Icon based on type */}
          <div className="flex-shrink-0">
            {suggestion.type === 'category' && (
              <svg
                className="w-5 h-5 text-primary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"
                />
              </svg>
            )}
            {suggestion.type === 'product' && (
              <svg
                className="w-5 h-5 text-secondary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M16 11V7a4 4 0 00-8 0v4M5 9h14l1 12H4L5 9z"
                />
              </svg>
            )}
            {suggestion.type === 'text' && (
              <svg
                className="w-5 h-5 text-base-content/50"
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
          </div>

          {/* Content */}
          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-2">
              <span className="font-medium truncate">
                {highlightMatch(suggestion.text, query)}
              </span>
              {suggestion.type === 'category' && suggestion.category && (
                <span className="badge badge-sm badge-primary">
                  {t('category')}
                </span>
              )}
              {suggestion.type === 'product' && (
                <span className="badge badge-sm badge-secondary">
                  {t('product')}
                </span>
              )}
            </div>

            {/* Additional metadata */}
            {isEnhanced && suggestion.metadata && (
              <div className="text-xs text-base-content/60 mt-1">
                {'category' in suggestion.metadata &&
                  suggestion.metadata.category && (
                    <span>{suggestion.metadata.category}</span>
                  )}
                {'price' in suggestion.metadata &&
                  suggestion.metadata.price && (
                    <span className="ml-2">
                      {new Intl.NumberFormat().format(
                        suggestion.metadata.price
                      )}{' '}
                      RSD
                    </span>
                  )}
                {'count' in suggestion.metadata &&
                  suggestion.metadata.count && (
                    <span className="ml-2">
                      ({suggestion.metadata.count} {t('results')})
                    </span>
                  )}
              </div>
            )}
          </div>

          {/* Product image if available */}
          {suggestion.type === 'product' &&
            suggestion.metadata &&
            'image' in suggestion.metadata &&
            suggestion.metadata.image && (
              <div className="flex-shrink-0">
                <Image
                  src={suggestion.metadata.image}
                  alt=""
                  width={40}
                  height={40}
                  className="w-10 h-10 object-cover rounded"
                />
              </div>
            )}
        </div>
      </div>
    );
  };

  const highlightMatch = (text: string, query: string) => {
    if (!query) return text;

    const regex = new RegExp(`(${query})`, 'gi');
    const parts = text.split(regex);

    return (
      <>
        {parts.map((part, i) =>
          regex.test(part) ? (
            <mark key={i} className="bg-yellow-200 text-inherit">
              {part}
            </mark>
          ) : (
            <span key={i}>{part}</span>
          )
        )}
      </>
    );
  };

  return (
    <div
      ref={ref}
      className="absolute top-full left-0 right-0 mt-1 bg-base-100 rounded-lg shadow-xl border border-base-200 max-h-96 overflow-y-auto z-50"
    >
      {isLoading && (
        <div className="p-4 text-center">
          <span className="loading loading-spinner loading-sm"></span>
        </div>
      )}

      {!isLoading && suggestions.length === 0 && searchHistory.length === 0 && (
        <div className="p-4 text-center text-base-content/60">
          {t('noSuggestions')}
        </div>
      )}

      {/* Suggestions */}
      {suggestions.length > 0 && (
        <div>
          <div className="px-4 py-2 text-xs font-semibold text-base-content/60 uppercase">
            {t('suggestions')}
          </div>
          {suggestions.map((suggestion, index) =>
            renderSuggestion(suggestion, index)
          )}
        </div>
      )}

      {/* Search History */}
      {searchHistory.length > 0 && query.length === 0 && (
        <div>
          <div className="px-4 py-2 text-xs font-semibold text-base-content/60 uppercase flex justify-between items-center">
            <span>{t('searchHistory')}</span>
            <button
              onClick={(e) => {
                e.stopPropagation();
                // Clear history
              }}
              className="text-primary hover:underline"
            >
              {t('clear')}
            </button>
          </div>
          {searchHistory.slice(0, 5).map((item, index) => (
            <div
              key={`history-${index}`}
              className={`px-4 py-2 cursor-pointer hover:bg-base-200 flex items-center gap-3`}
              onClick={() => onSelect(item)}
            >
              <svg
                className="w-5 h-5 text-base-content/50"
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
              <span>{item}</span>
            </div>
          ))}
        </div>
      )}

      {/* Trending Searches */}
      {trendingSearches.length > 0 && query.length === 0 && (
        <div>
          <div className="px-4 py-2 text-xs font-semibold text-base-content/60 uppercase">
            {t('trending')}
          </div>
          {trendingSearches.slice(0, 5).map((item, index) => (
            <div
              key={`trending-${index}`}
              className={`px-4 py-2 cursor-pointer hover:bg-base-200 flex items-center gap-3`}
              onClick={() => onSelect(item)}
            >
              <svg
                className="w-5 h-5 text-warning"
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
              <span>{item}</span>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
