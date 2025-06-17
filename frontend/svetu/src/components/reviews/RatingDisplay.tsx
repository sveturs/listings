'use client';

import React from 'react';

interface RatingDisplayProps {
  rating: number;
  maxRating?: number;
  size?: 'sm' | 'md' | 'lg';
  showValue?: boolean;
  totalReviews?: number;
  className?: string;
}

export const RatingDisplay: React.FC<RatingDisplayProps> = ({
  rating,
  maxRating = 5,
  size = 'md',
  showValue = true,
  totalReviews,
  className = '',
}) => {
  const sizeClasses = {
    sm: 'w-4 h-4',
    md: 'w-5 h-5',
    lg: 'w-6 h-6',
  };

  const textSizeClasses = {
    sm: 'text-sm',
    md: 'text-base',
    lg: 'text-lg',
  };

  const renderStars = () => {
    const stars = [];
    const fullStars = Math.floor(rating);
    const hasHalfStar = rating % 1 !== 0;

    for (let i = 1; i <= maxRating; i++) {
      if (i <= fullStars) {
        // Full star
        stars.push(
          <svg
            key={i}
            className={`${sizeClasses[size]} text-yellow-400 fill-current`}
            viewBox="0 0 20 20"
          >
            <path d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z" />
          </svg>
        );
      } else if (i === fullStars + 1 && hasHalfStar) {
        // Half star
        stars.push(
          <div key={i} className="relative">
            <svg
              className={`${sizeClasses[size]} text-gray-300 fill-current`}
              viewBox="0 0 20 20"
            >
              <path d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z" />
            </svg>
            <svg
              className={`${sizeClasses[size]} text-yellow-400 fill-current absolute top-0 left-0`}
              viewBox="0 0 20 20"
              style={{ clipPath: 'inset(0 50% 0 0)' }}
            >
              <path d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z" />
            </svg>
          </div>
        );
      } else {
        // Empty star
        stars.push(
          <svg
            key={i}
            className={`${sizeClasses[size]} text-gray-300 fill-current`}
            viewBox="0 0 20 20"
          >
            <path d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z" />
          </svg>
        );
      }
    }

    return stars;
  };

  return (
    <div className={`flex items-center gap-1 ${className}`}>
      <div className="flex">{renderStars()}</div>
      {showValue && (
        <span className={`${textSizeClasses[size]} text-gray-600 ml-1`}>
          {rating.toFixed(1)}
        </span>
      )}
      {totalReviews !== undefined && (
        <span className={`${textSizeClasses[size]} text-gray-500 ml-1`}>
          ({totalReviews})
        </span>
      )}
    </div>
  );
};
