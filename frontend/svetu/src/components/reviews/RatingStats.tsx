'use client';

import React from 'react';
import { useTranslations } from 'next-intl';
import { RatingDisplay } from './RatingDisplay';
import type { ReviewStats } from '@/types/review';

interface RatingStatsProps {
  stats: ReviewStats;
  className?: string;
}

export const RatingStats: React.FC<RatingStatsProps> = ({
  stats,
  className = '',
}) => {
  const t = useTranslations('reviews');
  const maxCount = Math.max(...Object.values(stats.rating_distribution || {}));

  const getRatingPercentage = (count: number) => {
    if (stats.total_reviews === 0) return 0;
    return (count / stats.total_reviews) * 100;
  };

  const getBarWidth = (count: number) => {
    if (maxCount === 0) return 0;
    return (count / maxCount) * 100;
  };

  return (
    <div className={`bg-base-100 rounded-lg p-6 ${className}`}>
      <div className="flex flex-col md:flex-row gap-6">
        {/* Overall Rating */}
        <div className="text-center md:text-left">
          <div className="text-4xl font-bold text-primary">
            {stats.average_rating.toFixed(1)}
          </div>
          <RatingDisplay
            rating={stats.average_rating}
            size="lg"
            showValue={false}
            className="mt-2"
          />
          <div className="text-sm text-base-content/70 mt-2">
            {stats.total_reviews} {t('stats.reviews')}
          </div>
          {stats.verified_reviews > 0 && (
            <div className="text-xs text-success mt-1">
              {Math.round((stats.verified_reviews / stats.total_reviews) * 100)}
              % {t('stats.verified')}
            </div>
          )}
        </div>

        {/* Rating Distribution */}
        <div className="flex-1">
          <h4 className="text-sm font-medium text-base-content/70 mb-3">
            {t('stats.ratingDistribution')}
          </h4>
          <div className="space-y-2">
            {[5, 4, 3, 2, 1].map((rating) => {
              const count = stats.rating_distribution?.[rating] || 0;
              const percentage = getRatingPercentage(count);
              const barWidth = getBarWidth(count);

              return (
                <div key={rating} className="flex items-center gap-3">
                  <div className="flex items-center gap-1 w-12">
                    <span className="text-sm">{rating}</span>
                    <svg
                      className="w-3 h-3 text-yellow-400 fill-current"
                      viewBox="0 0 20 20"
                    >
                      <path d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z" />
                    </svg>
                  </div>
                  <div className="flex-1 bg-base-200 rounded-full h-2 overflow-hidden">
                    <div
                      className="bg-yellow-400 h-full transition-all duration-300"
                      style={{ width: `${barWidth}%` }}
                    />
                  </div>
                  <div className="text-sm text-base-content/70 w-12 text-right">
                    {count}
                  </div>
                  <div className="text-xs text-base-content/50 w-10 text-right">
                    {percentage.toFixed(0)}%
                  </div>
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {/* Additional Stats */}
      {stats.photo_reviews > 0 && (
        <div className="mt-4 pt-4 border-t border-base-200">
          <div className="flex items-center gap-2 text-sm text-base-content/70">
            <svg
              className="w-4 h-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z"
              />
            </svg>
            <span>
              {stats.photo_reviews} {t('stats.withPhotos')}
            </span>
          </div>
        </div>
      )}
    </div>
  );
};
