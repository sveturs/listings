'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';

interface RatingInputProps {
  value: number;
  onChange: (rating: number) => void;
  maxRating?: number;
  size?: 'sm' | 'md' | 'lg';
  disabled?: boolean;
  className?: string;
  required?: boolean;
  error?: string;
}

export const RatingInput: React.FC<RatingInputProps> = ({
  value,
  onChange,
  maxRating = 5,
  size = 'md',
  disabled = false,
  className = '',
  required = false,
  error,
}) => {
  const [hoverRating, setHoverRating] = useState(0);
  const t = useTranslations('reviews.rating.labels');

  const sizeClasses = {
    sm: 'w-6 h-6',
    md: 'w-8 h-8',
    lg: 'w-12 h-12',
  };

  const handleClick = (rating: number) => {
    if (!disabled) {
      onChange(rating);
    }
  };

  const handleMouseEnter = (rating: number) => {
    if (!disabled) {
      setHoverRating(rating);
    }
  };

  const handleMouseLeave = () => {
    setHoverRating(0);
  };

  const renderStars = () => {
    const stars = [];
    const activeRating = hoverRating || value;

    for (let i = 1; i <= maxRating; i++) {
      const isActive = i <= activeRating;
      const isHovered = i <= hoverRating;

      stars.push(
        <button
          key={i}
          type="button"
          onClick={() => handleClick(i)}
          onMouseEnter={() => handleMouseEnter(i)}
          onMouseLeave={handleMouseLeave}
          disabled={disabled}
          className={`${disabled ? 'cursor-not-allowed' : 'cursor-pointer'} 
            transform transition-all duration-200 hover:scale-110
            focus:outline-none focus:scale-110`}
          aria-label={`Rate ${i} out of ${maxRating}`}
        >
          <svg
            className={`${sizeClasses[size]} ${
              isActive
                ? 'text-warning fill-warning drop-shadow-lg'
                : 'text-base-300 fill-base-300'
            } transition-all duration-300 ${
              isHovered && !disabled ? 'scale-125 rotate-12' : ''
            }`}
            viewBox="0 0 20 20"
          >
            <path d="M10 15l-5.878 3.09 1.123-6.545L.489 6.91l6.572-.955L10 0l2.939 5.955 6.572.955-4.756 4.635 1.123 6.545z" />
          </svg>
        </button>
      );
    }

    return stars;
  };

  const getRatingText = () => {
    const rating = hoverRating || value;
    const texts = {
      1: { text: t('terrible'), color: 'text-error' },
      2: { text: t('bad'), color: 'text-warning' },
      3: { text: t('normal'), color: 'text-info' },
      4: { text: t('good'), color: 'text-success' },
      5: { text: t('excellent'), color: 'text-success' },
    };
    return texts[rating as keyof typeof texts] || null;
  };

  const ratingInfo = getRatingText();

  return (
    <div className={className}>
      <div className="flex items-center gap-3">
        <div className="flex gap-1">{renderStars()}</div>
        {ratingInfo && (
          <span
            className={`text-lg font-medium ${ratingInfo.color} 
                          animate-in fade-in slide-in-from-left-2 duration-300`}
          >
            {ratingInfo.text}
          </span>
        )}
      </div>
      {required && value === 0 && error && (
        <p className="mt-2 text-sm text-error animate-in fade-in">{error}</p>
      )}
    </div>
  );
};
