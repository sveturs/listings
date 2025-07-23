'use client';

import React from 'react';

interface ListingCardSkeletonProps {
  viewMode?: 'grid' | 'list';
}

export const ListingCardSkeleton: React.FC<ListingCardSkeletonProps> = ({
  viewMode = 'grid',
}) => {
  if (viewMode === 'list') {
    return (
      <div className="card bg-base-100 shadow-sm animate-pulse">
        <div className="card-body p-4">
          <div className="flex gap-4">
            {/* Изображение слева */}
            <div className="w-32 h-32 flex-shrink-0 bg-base-300 rounded-lg" />

            {/* Информация */}
            <div className="flex-grow">
              <div className="flex justify-between items-start gap-4">
                <div className="flex-grow space-y-3">
                  {/* Категория */}
                  <div className="h-3 bg-base-300 rounded w-20" />

                  {/* Заголовок */}
                  <div className="h-5 bg-base-300 rounded w-3/4" />

                  {/* Описание */}
                  <div className="space-y-2">
                    <div className="h-3 bg-base-300 rounded w-full" />
                    <div className="h-3 bg-base-300 rounded w-2/3" />
                  </div>

                  {/* Продавец */}
                  <div className="flex items-center gap-2">
                    <div className="w-5 h-5 bg-base-300 rounded-full" />
                    <div className="h-3 bg-base-300 rounded w-24" />
                  </div>

                  {/* Статистика */}
                  <div className="flex items-center gap-3">
                    <div className="h-3 bg-base-300 rounded w-20" />
                    <div className="h-3 bg-base-300 rounded w-16" />
                    <div className="h-3 bg-base-300 rounded w-12" />
                  </div>
                </div>

                {/* Цена и действия */}
                <div className="flex flex-col items-end gap-2">
                  {/* Цена */}
                  <div className="text-right space-y-1">
                    <div className="h-6 bg-base-300 rounded w-24" />
                  </div>

                  {/* Кнопки */}
                  <div className="flex gap-2">
                    <div className="w-8 h-8 bg-base-300 rounded-full" />
                    <div className="w-16 h-8 bg-base-300 rounded" />
                    <div className="w-16 h-8 bg-base-300 rounded" />
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    );
  }

  // Grid view skeleton
  return (
    <div className="card card-compact bg-base-100 shadow-sm animate-pulse">
      {/* Изображение */}
      <figure className="relative aspect-square bg-base-300">
        {/* Badges placeholder */}
        <div className="absolute top-2 left-2 flex flex-col gap-1">
          <div className="h-5 w-16 bg-base-200 rounded" />
          <div className="h-5 w-12 bg-base-200 rounded" />
        </div>

        {/* Favorite button placeholder */}
        <div className="absolute top-2 right-2">
          <div className="w-8 h-8 bg-base-200 rounded-full" />
        </div>

        {/* Photo count placeholder */}
        <div className="absolute bottom-2 right-2">
          <div className="h-5 w-14 bg-base-200 rounded" />
        </div>
      </figure>

      <div className="card-body p-3 space-y-2">
        {/* Категория */}
        <div className="h-3 bg-base-300 rounded w-1/3" />

        {/* Заголовок */}
        <div className="space-y-1">
          <div className="h-4 bg-base-300 rounded w-full" />
          <div className="h-4 bg-base-300 rounded w-3/4" />
        </div>

        {/* Продавец */}
        <div className="flex items-center gap-2">
          <div className="w-5 h-5 bg-base-300 rounded-full" />
          <div className="h-3 bg-base-300 rounded w-20" />
          <div className="h-3 bg-base-300 rounded w-8" />
        </div>

        {/* Локация и статистика */}
        <div className="space-y-1">
          <div className="h-3 bg-base-300 rounded w-2/3" />
          <div className="flex items-center gap-3">
            <div className="h-3 bg-base-300 rounded w-16" />
            <div className="h-3 bg-base-300 rounded w-12" />
          </div>
        </div>

        {/* Цена и действия */}
        <div className="flex justify-between items-end mt-2">
          <div className="space-y-1">
            <div className="h-5 bg-base-300 rounded w-20" />
            <div className="h-3 bg-base-300 rounded w-24" />
          </div>

          <div className="flex gap-1">
            <div className="w-6 h-6 bg-base-300 rounded" />
            <div className="w-6 h-6 bg-base-300 rounded" />
          </div>
        </div>
      </div>
    </div>
  );
};
