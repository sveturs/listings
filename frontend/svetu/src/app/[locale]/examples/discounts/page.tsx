'use client';

import React, { useState } from 'react';
import { PageTransition } from '@/components/ui/PageTransition';
import { DiscountBadge } from '@/components/ui/DiscountBadge';
import { PriceHistoryModal } from '@/components/c2c/PriceHistoryModal';
import { BlackFridayBadge } from '@/components/b2c/BlackFridayBadge';
import { ArrowLeft, TrendingDown, Calendar, ShoppingBag } from 'lucide-react';
import Link from 'next/link';

export default function DiscountsDemo() {
  const [showPriceHistory, setShowPriceHistory] = useState(false);
  const [selectedListingId, setSelectedListingId] = useState<number>(1);

  // Примеры данных для демо
  const discountExamples = [
    {
      oldPrice: 15000,
      currentPrice: 10500,
      size: 'sm' as const,
      label: 'Малый бейдж',
    },
    {
      oldPrice: 25000,
      currentPrice: 19999,
      size: 'md' as const,
      label: 'Средний бейдж',
    },
    {
      oldPrice: 50000,
      currentPrice: 35000,
      size: 'lg' as const,
      label: 'Большой бейдж',
    },
  ];

  const storefrontExamples = [
    {
      name: 'Tech Store',
      stats: {
        totalProducts: 100,
        discountedProducts: 60,
        averageDiscount: 25,
      },
      size: 'sm' as const,
    },
    {
      name: 'Fashion Hub',
      stats: {
        totalProducts: 200,
        discountedProducts: 80,
        averageDiscount: 18,
      },
      size: 'md' as const,
    },
    {
      name: 'Electronics World',
      stats: {
        totalProducts: 150,
        discountedProducts: 120,
        averageDiscount: 30,
      },
      size: 'lg' as const,
    },
  ];

  const handlePriceHistoryClick = (listingId: number) => {
    setSelectedListingId(listingId);
    setShowPriceHistory(true);
  };

  return (
    <PageTransition>
      <div className="container mx-auto px-4 py-8 max-w-6xl">
        {/* Заголовок */}
        <div className="flex items-center gap-4 mb-8">
          <Link href="/en/examples" className="btn btn-ghost btn-circle">
            <ArrowLeft className="w-5 h-5" />
          </Link>
          <div>
            <h1 className="text-3xl font-bold flex items-center gap-3">
              <TrendingDown className="w-8 h-8 text-error" />
              Система скидок
            </h1>
            <p className="text-base-content/70 mt-2">
              Демонстрация компонентов для отображения скидок и истории цен
            </p>
          </div>
        </div>

        <div className="grid gap-8">
          {/* DiscountBadge Demo */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title text-2xl">
                <TrendingDown className="w-6 h-6" />
                DiscountBadge
              </h2>
              <p className="text-base-content/70 mb-6">
                Интерактивные бейджи скидок с возможностью просмотра истории цен
              </p>

              <div className="grid md:grid-cols-3 gap-6">
                {discountExamples.map((example, index) => (
                  <div key={index} className="card bg-base-200">
                    <div className="card-body items-center text-center">
                      <h3 className="card-title text-lg">{example.label}</h3>

                      <div className="mockup-window bg-base-300 w-full max-w-xs">
                        <div className="flex justify-center px-4 py-16 bg-base-200">
                          <div className="flex flex-col items-center gap-3">
                            <div className="text-sm line-through text-base-content/50">
                              {example.oldPrice.toLocaleString()} RSD
                            </div>
                            <div className="flex items-center gap-2">
                              <div className="text-xl font-bold">
                                {example.currentPrice.toLocaleString()} RSD
                              </div>
                              <DiscountBadge
                                oldPrice={example.oldPrice}
                                currentPrice={example.currentPrice}
                                size={example.size}
                                onClick={() =>
                                  handlePriceHistoryClick(index + 1)
                                }
                              />
                            </div>
                          </div>
                        </div>
                      </div>

                      <div className="text-xs text-base-content/60 mt-2">
                        Скидка:{' '}
                        {Math.round(
                          ((example.oldPrice - example.currentPrice) /
                            example.oldPrice) *
                            100
                        )}
                        %
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              <div className="alert alert-info mt-6">
                <TrendingDown className="w-5 h-5" />
                <span>
                  Нажмите на любой бейдж скидки, чтобы открыть модалку с
                  историей цен
                </span>
              </div>
            </div>
          </div>

          {/* BlackFridayBadge Demo */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title text-2xl">
                <ShoppingBag className="w-6 h-6" />
                BlackFridayBadge
              </h2>
              <p className="text-base-content/70 mb-6">
                Динамические бейджи для витрин с множественными скидками
              </p>

              <div className="grid md:grid-cols-3 gap-6">
                {storefrontExamples.map((example, index) => (
                  <div key={index} className="card bg-base-200">
                    <div className="card-body">
                      <h3 className="card-title text-lg">{example.name}</h3>

                      <div className="stats stats-vertical shadow mb-4">
                        <div className="stat">
                          <div className="stat-title">Всего товаров</div>
                          <div className="stat-value text-sm">
                            {example.stats.totalProducts}
                          </div>
                        </div>
                        <div className="stat">
                          <div className="stat-title">Со скидками</div>
                          <div className="stat-value text-sm">
                            {example.stats.discountedProducts}
                          </div>
                        </div>
                        <div className="stat">
                          <div className="stat-title">Средняя скидка</div>
                          <div className="stat-value text-sm">
                            {example.stats.averageDiscount}%
                          </div>
                        </div>
                      </div>

                      <div className="flex justify-center">
                        <BlackFridayBadge
                          discountStats={example.stats}
                          size={example.size}
                        />
                      </div>
                    </div>
                  </div>
                ))}
              </div>

              <div className="alert alert-warning mt-6">
                <Calendar className="w-5 h-5" />
                <span>
                  Бейджи автоматически появляются только если &gt;20% товаров
                  имеют скидки &gt;10%
                </span>
              </div>
            </div>
          </div>

          {/* Интеграция в карточки товаров */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title text-2xl">
                Интеграция в карточки товаров
              </h2>
              <p className="text-base-content/70 mb-6">
                Система скидок интегрирована в EnhancedListingCard
              </p>

              <div className="mockup-browser border bg-base-300">
                <div className="mockup-browser-toolbar">
                  <div className="input">http://localhost:3001/c2c</div>
                </div>
                <div className="flex justify-center px-4 py-16 bg-base-200">
                  <div className="grid md:grid-cols-2 gap-4 max-w-md">
                    {/* Пример карточки со скидкой */}
                    <div className="card card-compact bg-base-100 shadow">
                      <figure className="relative">
                        <div className="w-full h-32 bg-gradient-to-br from-blue-400 to-purple-500 flex items-center justify-center text-white font-bold">
                          ТОВАР 1
                        </div>
                        <div className="absolute top-2 left-2">
                          <div className="badge badge-success badge-sm">
                            б/у
                          </div>
                        </div>
                      </figure>
                      <div className="card-body">
                        <h3 className="card-title text-sm">Товар со скидкой</h3>
                        <div className="flex items-center gap-2 mb-1">
                          <p className="text-xs line-through text-base-content/50">
                            15,000 RSD
                          </p>
                          <DiscountBadge
                            oldPrice={15000}
                            currentPrice={10500}
                            size="sm"
                            onClick={() => handlePriceHistoryClick(1)}
                          />
                        </div>
                        <div className="text-lg font-bold">10,500 RSD</div>
                      </div>
                    </div>

                    {/* Пример обычной карточки */}
                    <div className="card card-compact bg-base-100 shadow">
                      <figure className="relative">
                        <div className="w-full h-32 bg-gradient-to-br from-green-400 to-blue-500 flex items-center justify-center text-white font-bold">
                          ТОВАР 2
                        </div>
                        <div className="absolute top-2 left-2">
                          <div className="badge badge-primary badge-sm">
                            новый
                          </div>
                        </div>
                      </figure>
                      <div className="card-body">
                        <h3 className="card-title text-sm">Товар без скидки</h3>
                        <div className="text-lg font-bold">8,500 RSD</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Технические детали */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title text-2xl">Технические детали</h2>

              <div className="grid md:grid-cols-2 gap-6">
                <div>
                  <h3 className="text-lg font-semibold mb-3">DiscountBadge</h3>
                  <ul className="space-y-2 text-sm">
                    <li>• Автоматически вычисляет процент скидки</li>
                    <li>• Скрывается при скидке &lt;5%</li>
                    <li>• Интерактивный клик для истории цен</li>
                    <li>• Размеры: sm, md, lg</li>
                    <li>• Hover эффекты</li>
                  </ul>
                </div>

                <div>
                  <h3 className="text-lg font-semibold mb-3">
                    PriceHistoryModal
                  </h3>
                  <ul className="space-y-2 text-sm">
                    <li>• График изменения цены (Chart.js)</li>
                    <li>• Детекция манипуляций с ценой</li>
                    <li>• Статистика: мин/макс/текущая цена</li>
                    <li>• API: /api/v1/c2c/listings/:id/price-history</li>
                    <li>• Автоматические предупреждения</li>
                  </ul>
                </div>

                <div>
                  <h3 className="text-lg font-semibold mb-3">
                    BlackFridayBadge
                  </h3>
                  <ul className="space-y-2 text-sm">
                    <li>• 3 варианта: Sale, Hot Deals, Black Friday</li>
                    <li>• Автоматическая категоризация</li>
                    <li>• Анимированные градиенты</li>
                    <li>• Процент товаров со скидками</li>
                    <li>• Условная логика отображения</li>
                  </ul>
                </div>

                <div>
                  <h3 className="text-lg font-semibold mb-3">Backend API</h3>
                  <ul className="space-y-2 text-sm">
                    <li>• Таблица price_history</li>
                    <li>• Автоматические триггеры</li>
                    <li>• Защита от манипуляций</li>
                    <li>• Материализованные представления</li>
                    <li>• Оптимизированные индексы</li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>

        {/* Навигация */}
        <div className="mt-8 flex justify-between">
          <Link href="/en/examples" className="btn btn-outline">
            ← Вернуться к примерам
          </Link>
          <Link href="/en/examples/transitions" className="btn btn-primary">
            Следующий пример →
          </Link>
        </div>
      </div>

      {/* Price History Modal */}
      {showPriceHistory && (
        <PriceHistoryModal
          listingId={selectedListingId}
          isOpen={showPriceHistory}
          onClose={() => setShowPriceHistory(false)}
        />
      )}
    </PageTransition>
  );
}
