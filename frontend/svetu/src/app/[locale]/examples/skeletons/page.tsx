'use client';

import React, { useState } from 'react';
import { ListingCardSkeleton } from '@/components/ui/skeletons/ListingCardSkeleton';
import { ListingGridSkeleton } from '@/components/ui/skeletons/ListingGridSkeleton';
import { EnhancedListingCardSkeleton } from '@/components/ui/skeletons/EnhancedListingCardSkeleton';
import { EnhancedListingGridSkeleton } from '@/components/ui/skeletons/EnhancedListingGridSkeleton';

export default function SkeletonsExamplePage() {
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');

  return (
    <div className="container mx-auto p-6 max-w-7xl">
      <h1 className="text-3xl font-bold mb-8">Skeleton Loading States</h1>

      {/* Переключатель режима отображения */}
      <div className="mb-8 flex justify-center">
        <div className="join">
          <button
            className={`join-item btn ${viewMode === 'grid' ? 'btn-active' : ''}`}
            onClick={() => setViewMode('grid')}
          >
            Grid View
          </button>
          <button
            className={`join-item btn ${viewMode === 'list' ? 'btn-active' : ''}`}
            onClick={() => setViewMode('list')}
          >
            List View
          </button>
        </div>
      </div>

      {/* Стандартные скелетоны */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-6">Стандартные скелетоны</h2>
        
        <div className="space-y-8">
          <div>
            <h3 className="text-lg font-medium mb-4">Одиночная карточка</h3>
            <div className={viewMode === 'grid' ? 'max-w-sm' : ''}>
              <ListingCardSkeleton viewMode={viewMode} />
            </div>
          </div>

          <div>
            <h3 className="text-lg font-medium mb-4">Сетка карточек</h3>
            <ListingGridSkeleton count={4} viewMode={viewMode} />
          </div>
        </div>
      </section>

      {/* Улучшенные скелетоны */}
      <section className="mb-12">
        <h2 className="text-2xl font-semibold mb-6">Улучшенные скелетоны с анимацией</h2>
        
        <div className="space-y-8">
          <div>
            <h3 className="text-lg font-medium mb-4">Одиночная карточка с shimmer эффектом</h3>
            <div className={viewMode === 'grid' ? 'max-w-sm' : ''}>
              <EnhancedListingCardSkeleton viewMode={viewMode} />
            </div>
          </div>

          <div>
            <h3 className="text-lg font-medium mb-4">Сетка с каскадной анимацией</h3>
            <EnhancedListingGridSkeleton count={8} viewMode={viewMode} />
          </div>
        </div>
      </section>

      {/* Особенности */}
      <section className="card bg-base-200 p-6">
        <h2 className="text-2xl font-semibold mb-4">Особенности улучшенных скелетонов</h2>
        <ul className="list-disc list-inside space-y-2">
          <li><strong>Shimmer эффект:</strong> Плавная анимация градиента, имитирующая загрузку</li>
          <li><strong>Каскадная анимация:</strong> Элементы появляются по очереди с задержкой</li>
          <li><strong>Адаптивность:</strong> Поддержка режимов grid и list</li>
          <li><strong>Производительность:</strong> Используются только CSS анимации</li>
          <li><strong>Темы:</strong> Автоматическая адаптация к светлой/темной теме</li>
        </ul>

        <div className="mt-6">
          <h3 className="text-lg font-medium mb-3">Примеры использования:</h3>
          <div className="mockup-code">
            <pre data-prefix="1"><code>{`import { EnhancedListingGridSkeleton } from '@/components/ui/skeletons';`}</code></pre>
            <pre data-prefix="2"><code>{``}</code></pre>
            <pre data-prefix="3"><code>{`// В компоненте`}</code></pre>
            <pre data-prefix="4"><code>{`{isLoading ? (`}</code></pre>
            <pre data-prefix="5"><code>{`  <EnhancedListingGridSkeleton count={8} viewMode="grid" />`}</code></pre>
            <pre data-prefix="6"><code>{`) : (`}</code></pre>
            <pre data-prefix="7"><code>{`  <ListingGrid items={items} />`}</code></pre>
            <pre data-prefix="8"><code>{`)}`}</code></pre>
          </div>
        </div>
      </section>

      {/* Сравнение анимаций */}
      <section className="mt-8 card bg-base-200 p-6">
        <h2 className="text-2xl font-semibold mb-4">Сравнение типов анимаций</h2>
        
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <h3 className="font-medium mb-3">Pulse (стандарт)</h3>
            <div className="h-24 bg-base-300 rounded animate-pulse"></div>
            <p className="text-sm mt-2 text-base-content/70">
              Простая пульсация прозрачности
            </p>
          </div>

          <div>
            <h3 className="font-medium mb-3">Gentle Pulse</h3>
            <div className="h-24 bg-base-300 rounded animate-gentle-pulse"></div>
            <p className="text-sm mt-2 text-base-content/70">
              Более плавная пульсация
            </p>
          </div>

          <div>
            <h3 className="font-medium mb-3">Shimmer</h3>
            <div className="h-24 bg-base-300 rounded relative overflow-hidden before:absolute before:inset-0 before:-translate-x-full before:animate-[shimmer_1.5s_infinite] before:bg-gradient-to-r before:from-transparent before:via-white/10 before:to-transparent"></div>
            <p className="text-sm mt-2 text-base-content/70">
              Эффект движущегося блика
            </p>
          </div>
        </div>
      </section>
    </div>
  );
}