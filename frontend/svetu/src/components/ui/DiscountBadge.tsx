'use client';

import React from 'react';
import { TrendingDown } from 'lucide-react';

interface DiscountBadgeProps {
  oldPrice: number;
  currentPrice: number;
  onClick?: (e: React.MouseEvent) => void;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const DiscountBadge: React.FC<DiscountBadgeProps> = ({
  oldPrice,
  currentPrice,
  onClick,
  size = 'md',
  className = '',
}) => {
  const discountPercent = Math.round(
    ((oldPrice - currentPrice) / oldPrice) * 100
  );

  // Не показывать скидки менее 5%
  if (discountPercent < 5) return null;

  const sizeClasses = {
    sm: 'text-xs px-2 py-1',
    md: 'text-sm px-3 py-1.5',
    lg: 'text-base px-4 py-2',
  };

  return (
    <button
      onClick={onClick}
      className={`
        badge badge-error gap-1 cursor-pointer 
        hover:scale-105 transition-transform
        ${sizeClasses[size]} ${className}
      `}
      title="Нажмите, чтобы увидеть историю цены"
      type="button"
    >
      <TrendingDown className="w-3 h-3" />-{discountPercent}%
    </button>
  );
};
