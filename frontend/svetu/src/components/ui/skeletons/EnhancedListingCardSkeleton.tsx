'use client';

import React from 'react';

interface EnhancedListingCardSkeletonProps {
  viewMode?: 'grid' | 'list';
}

export const EnhancedListingCardSkeleton: React.FC<EnhancedListingCardSkeletonProps> = ({
  viewMode = 'grid',
}) => {
  // Общий класс для анимированных блоков
  const shimmerClass = 'relative overflow-hidden bg-base-300 before:absolute before:inset-0 before:-translate-x-full before:animate-[shimmer_1.5s_infinite] before:bg-gradient-to-r before:from-transparent before:via-white/10 before:to-transparent';

  if (viewMode === 'list') {
    return (
      <div className="card bg-base-100 shadow-sm">
        <div className="card-body p-4">
          <div className="flex gap-4">
            {/* Изображение слева */}
            <div className={`w-32 h-32 flex-shrink-0 rounded-lg ${shimmerClass}`} />

            {/* Информация */}
            <div className="flex-grow">
              <div className="flex justify-between items-start gap-4">
                <div className="flex-grow space-y-3">
                  {/* Категория */}
                  <div className={`h-3 rounded w-20 ${shimmerClass}`} />

                  {/* Заголовок */}
                  <div className={`h-5 rounded w-3/4 ${shimmerClass}`} />

                  {/* Описание */}
                  <div className="space-y-2">
                    <div className={`h-3 rounded w-full ${shimmerClass}`} />
                    <div className={`h-3 rounded w-2/3 ${shimmerClass}`} />
                  </div>

                  {/* Продавец */}
                  <div className="flex items-center gap-2">
                    <div className={`w-5 h-5 rounded-full ${shimmerClass}`} />
                    <div className={`h-3 rounded w-24 ${shimmerClass}`} />
                  </div>

                  {/* Статистика */}
                  <div className="flex items-center gap-3">
                    <div className={`h-3 rounded w-20 ${shimmerClass}`} />
                    <div className={`h-3 rounded w-16 ${shimmerClass}`} />
                    <div className={`h-3 rounded w-12 ${shimmerClass}`} />
                  </div>
                </div>

                {/* Цена и действия */}
                <div className="flex flex-col items-end gap-2">
                  {/* Цена */}
                  <div className="text-right space-y-1">
                    <div className={`h-6 rounded w-24 ${shimmerClass}`} />
                  </div>

                  {/* Кнопки */}
                  <div className="flex gap-2">
                    <div className={`w-8 h-8 rounded-full ${shimmerClass}`} />
                    <div className={`w-16 h-8 rounded ${shimmerClass}`} />
                    <div className={`w-16 h-8 rounded ${shimmerClass}`} />
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
    <div className="card card-compact bg-base-100 shadow-sm">
      {/* Изображение */}
      <figure className={`relative aspect-square ${shimmerClass}`}>
        {/* Badges placeholder */}
        <div className="absolute top-2 left-2 flex flex-col gap-1">
          <div className="h-5 w-16 bg-base-200/50 rounded animate-pulse" />
          <div className="h-5 w-12 bg-base-200/50 rounded animate-pulse delay-75" />
        </div>

        {/* Favorite button placeholder */}
        <div className="absolute top-2 right-2">
          <div className="w-8 h-8 bg-base-200/50 rounded-full animate-pulse delay-150" />
        </div>

        {/* Photo count placeholder */}
        <div className="absolute bottom-2 right-2">
          <div className="h-5 w-14 bg-base-200/50 rounded animate-pulse delay-100" />
        </div>
      </figure>

      <div className="card-body p-3 space-y-2">
        {/* Категория */}
        <div className={`h-3 rounded w-1/3 ${shimmerClass}`} />

        {/* Заголовок */}
        <div className="space-y-1">
          <div className={`h-4 rounded w-full ${shimmerClass}`} />
          <div className={`h-4 rounded w-3/4 ${shimmerClass}`} />
        </div>

        {/* Продавец */}
        <div className="flex items-center gap-2">
          <div className={`w-5 h-5 rounded-full ${shimmerClass}`} />
          <div className={`h-3 rounded w-20 ${shimmerClass}`} />
          <div className={`h-3 rounded w-8 ${shimmerClass}`} />
        </div>

        {/* Локация и статистика */}
        <div className="space-y-1">
          <div className={`h-3 rounded w-2/3 ${shimmerClass}`} />
          <div className="flex items-center gap-3">
            <div className={`h-3 rounded w-16 ${shimmerClass}`} />
            <div className={`h-3 rounded w-12 ${shimmerClass}`} />
          </div>
        </div>

        {/* Цена и действия */}
        <div className="flex justify-between items-end mt-2">
          <div className="space-y-1">
            <div className={`h-5 rounded w-20 ${shimmerClass}`} />
            <div className={`h-3 rounded w-24 ${shimmerClass}`} />
          </div>

          <div className="flex gap-1">
            <div className={`w-6 h-6 rounded ${shimmerClass}`} />
            <div className={`w-6 h-6 rounded ${shimmerClass}`} />
          </div>
        </div>
      </div>
    </div>
  );
};