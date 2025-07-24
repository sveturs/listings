'use client';

import React from 'react';
import { Zap, Percent } from 'lucide-react';

interface DiscountStats {
  totalProducts: number;
  discountedProducts: number;
  averageDiscount: number;
  maxDiscount?: number;
}

interface BlackFridayBadgeProps {
  discountStats: DiscountStats;
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}

export const BlackFridayBadge: React.FC<BlackFridayBadgeProps> = ({
  discountStats,
  size = 'md',
  className = '',
}) => {
  const discountedPercent =
    discountStats.totalProducts > 0
      ? (discountStats.discountedProducts / discountStats.totalProducts) * 100
      : 0;

  // Показывать только если >20% товаров со скидками >10%
  if (discountedPercent < 20 || discountStats.averageDiscount < 10) {
    return null;
  }

  const sizeClasses = {
    sm: 'text-xs px-2 py-1',
    md: 'text-sm px-3 py-2',
    lg: 'text-base px-4 py-3',
  };

  const iconClasses = {
    sm: 'w-3 h-3',
    md: 'w-4 h-4',
    lg: 'w-5 h-5',
  };

  // Определяем тип бейджа в зависимости от уровня скидок
  const getBadgeVariant = () => {
    if (discountedPercent >= 50 && discountStats.averageDiscount >= 20) {
      return {
        style: 'bg-gradient-to-r from-red-600 to-black text-white',
        title: 'BLACK FRIDAY',
        subtitle: 'MEGA SALE',
      };
    } else if (discountedPercent >= 30 && discountStats.averageDiscount >= 15) {
      return {
        style: 'bg-gradient-to-r from-orange-600 to-red-600 text-white',
        title: 'HOT DEALS',
        subtitle: 'BIG SALE',
      };
    } else {
      return {
        style: 'bg-gradient-to-r from-yellow-600 to-orange-600 text-white',
        title: 'SALE',
        subtitle: 'DISCOUNTS',
      };
    }
  };

  const variant = getBadgeVariant();

  return (
    <div
      className={`
        ${variant.style} ${sizeClasses[size]} ${className}
        flex items-center gap-2 rounded-full font-bold 
        shadow-lg animate-pulse cursor-pointer
        hover:scale-105 transition-transform
      `}
      title={`${Math.round(discountedPercent)}% товаров со скидками (средняя ${Math.round(discountStats.averageDiscount)}%)`}
    >
      <Zap className={`${iconClasses[size]} text-yellow-300`} />

      <div className="flex flex-col leading-none">
        <span
          className={
            size === 'sm' ? 'text-xs' : size === 'lg' ? 'text-lg' : 'text-sm'
          }
        >
          {variant.title}
        </span>
        {size !== 'sm' && (
          <span className="text-xs opacity-90">{variant.subtitle}</span>
        )}
      </div>

      <div className="flex items-center gap-1">
        <Percent className={iconClasses[size]} />
        <span className={size === 'sm' ? 'text-xs' : 'text-sm'}>
          {Math.round(discountedPercent)}%
        </span>
      </div>
    </div>
  );
};
