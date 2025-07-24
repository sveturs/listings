'use client';

import React from 'react';
import { ListingCardSkeleton } from './ListingCardSkeleton';

interface ListingGridSkeletonProps {
  count?: number;
  viewMode?: 'grid' | 'list';
}

export const ListingGridSkeleton: React.FC<ListingGridSkeletonProps> = ({
  count = 8,
  viewMode = 'grid',
}) => {
  const gridClass =
    viewMode === 'grid'
      ? 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4'
      : 'flex flex-col gap-4';

  return (
    <div className={gridClass}>
      {Array.from({ length: count }).map((_, index) => (
        <ListingCardSkeleton key={index} viewMode={viewMode} />
      ))}
    </div>
  );
};
