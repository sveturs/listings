'use client';

import React from 'react';
import { EnhancedListingCardSkeleton } from './EnhancedListingCardSkeleton';

interface EnhancedListingGridSkeletonProps {
  count?: number;
  viewMode?: 'grid' | 'list';
}

export const EnhancedListingGridSkeleton: React.FC<EnhancedListingGridSkeletonProps> = ({
  count = 8,
  viewMode = 'grid',
}) => {
  // Генерируем разные задержки для каждого элемента
  const items = Array.from({ length: count }, (_, i) => ({
    id: i,
    delay: i * 100, // Каскадная анимация
  }));

  return (
    <div
      className={
        viewMode === 'grid'
          ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4'
          : 'flex flex-col gap-4'
      }
    >
      {items.map((item) => (
        <div
          key={item.id}
          style={{
            animationDelay: `${item.delay}ms`,
          }}
          className="animate-gentle-pulse"
        >
          <EnhancedListingCardSkeleton viewMode={viewMode} />
        </div>
      ))}
    </div>
  );
};