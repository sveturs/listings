'use client';

import { forwardRef } from 'react';
import { useTranslations } from 'next-intl';

interface InfiniteScrollTriggerProps {
  loading: boolean;
  hasMore: boolean;
  onLoadMore: () => void;
  showButton?: boolean;
  className?: string;
}

export const InfiniteScrollTrigger = forwardRef<
  HTMLDivElement,
  InfiniteScrollTriggerProps
>(
  (
    { loading, hasMore, onLoadMore, showButton = true, className = '' },
    ref
  ) => {
    const t = useTranslations('common');

    return (
      <>
        {/* Invisible trigger element for IntersectionObserver */}
        <div
          ref={ref}
          className={`w-full h-20 mt-8 ${className}`}
          aria-hidden="true"
        >
          {loading && (
            <div className="flex justify-center items-center h-full">
              <span className="loading loading-spinner loading-lg"></span>
            </div>
          )}
        </div>

        {/* Fallback button for manual loading */}
        {showButton && hasMore && !loading && (
          <div className="text-center mb-8">
            <button
              className="btn btn-outline btn-sm"
              onClick={onLoadMore}
              aria-label={t('loadMore')}
            >
              {t('loadMore')}
              <svg
                className="w-4 h-4 ml-2"
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
          </div>
        )}
      </>
    );
  }
);

InfiniteScrollTrigger.displayName = 'InfiniteScrollTrigger';

export default InfiniteScrollTrigger;
