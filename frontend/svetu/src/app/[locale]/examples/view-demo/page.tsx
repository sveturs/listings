'use client';

import React, { useState } from 'react';
import Link from 'next/link';
import { useLocale } from 'next-intl';
import {
  Eye,
  ShoppingBag,
  Package,
  CreditCard,
  Map,
  Image as ImageIcon,
  Layout,
  Maximize2,
  MessageSquare,
  TrendingUp,
  Settings,
  Grid,
  List,
  Layers,
} from 'lucide-react';
import { QuickView } from '@/components/ui/QuickView';
import { PriceHistoryModal } from '@/components/marketplace/PriceHistoryModal';
import VariantSelectionModal from '@/components/cart/VariantSelectionModal';
import { EnhancedListingCard } from '@/components/marketplace/EnhancedListingCard';
import { UnifiedProductCard } from '@/components/common/UnifiedProductCard';
import { ImageGallery } from '@/components/reviews/ImageGallery';
import type { MarketplaceItem } from '@/types/marketplace';
import type { UnifiedProduct } from '@/types/unified-product';

export default function ViewDemoPage() {
  const locale = useLocale();
  const [showQuickView, setShowQuickView] = useState(false);
  const [showPriceHistory, setShowPriceHistory] = useState(false);
  const [showVariantModal, setShowVariantModal] = useState(false);
  const [showImageGallery, setShowImageGallery] = useState(false);
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');
  const [gridColumns, setGridColumns] = useState<1 | 2 | 3 | 4>(3);

  // Примеры реальных ID из базы данных
  const marketplaceListingId = 106; // Volkswagen Atlas Cross Sport
  const storefrontSlug = 'agenstvo'; // Используем существующую витрину
  const storefrontProductId = 111; // телефон - реальный товар из витрины
  const orderId = 57; // Реальный заказ из storefront_orders

  // Тестовые данные для карточек
  const sampleMarketplaceListing: MarketplaceItem = {
    id: 106,
    user_id: 8,
    category_id: 1003,
    title: 'Volkswagen Atlas Cross Sport',
    description: 'Отличный автомобиль в хорошем состоянии',
    price: 500000,
    condition: 'used',
    status: 'active',
    location: 'Белград',
    latitude: 44.8178131,
    longitude: 20.4568974,
    views_count: 125,
    created_at: '2025-08-02T22:11:25Z',
    updated_at: '2025-08-02T22:11:25Z',
    images: [
      {
        id: 49,
        is_main: true,
        public_url:
          'http://localhost:9000/listings/106/1754172685620502083.jpg',
      },
    ],
    user: {
      id: 8,
      name: 'EmailEmail',
      email: 'EmailEmail@EmailEmail.ru',
      picture_url: '',
    },
    category: {
      id: 1003,
      name: 'Automobili',
      slug: 'automotive',
    },
    attributes: [
      {
        id: 2204,
        attribute_id: 2204,
        name: 'fuel_type',
        attribute_name: 'fuel_type',
        display_name: 'Gorivo',
        attribute_type: 'select',
        text_value: 'petrol',
        display_value: 'Benzin',
        value: 'petrol',
        is_required: true,
        show_in_card: true,
        show_in_list: false,
      },
      {
        id: 2205,
        attribute_id: 2205,
        name: 'transmission',
        attribute_name: 'transmission',
        display_name: 'Menjač',
        attribute_type: 'select',
        text_value: 'automatic',
        display_value: 'Automatik',
        value: 'automatic',
        is_required: true,
        show_in_card: true,
        show_in_list: false,
      },
    ],
    is_favorite: false,
    show_on_map: true,
    has_discount: false,
  };

  const sampleUnifiedProduct: UnifiedProduct = {
    id: 1,
    type: 'storefront',
    name: 'iPhone 15 Pro',
    description: 'Новый iPhone 15 Pro с гарантией',
    price: 145000,
    currency: 'RSD',
    images: [
      {
        id: '1',
        url: 'http://localhost:9000/listings/109/1754253315899650123.jpg',
        isMain: true,
      },
    ],
    condition: 'new',
    stockStatus: 'in_stock',
    stockQuantity: 10,
    category: {
      id: 1001,
      name: 'Elektronika',
      slug: 'electronics',
    },
    seller: {
      id: 1,
      name: 'Demo Store',
      rating: 4.8,
      reviewsCount: 125,
    },
    location: {
      city: 'Белград',
      country: 'Сербия',
      latitude: 44.8178131,
      longitude: 20.4568974,
    },
    createdAt: '2025-08-16T13:17:39Z',
    updatedAt: '2025-08-16T13:17:39Z',
    storefront: {
      id: 1,
      name: 'Агентство недвижимости',
      slug: 'agenstvo',
    },
    isFavorite: false,
    viewsCount: 250,
    variants: [],
  };

  const mockQuickViewProduct = {
    id: '106',
    title: 'Volkswagen Atlas Cross Sport',
    price: '500,000 RSD',
    description:
      'Отличный автомобиль в хорошем состоянии. Экономичный двигатель, автоматическая коробка передач.',
    images: ['http://localhost:9000/listings/106/1754172685620502083.jpg'],
    category: 'Automobili',
    seller: {
      name: 'EmailEmail',
      rating: 4.5,
      totalReviews: 23,
      avatar: undefined,
    },
    location: {
      address: 'Белград, Сербия',
      distance: 5.2,
    },
    stats: {
      views: 125,
      favorites: 8,
    },
    condition: 'used' as const,
    storefrontId: undefined,
    storefrontName: undefined,
    storefrontSlug: undefined,
    stockQuantity: 1,
  };

  return (
    <div className="min-h-screen bg-base-100">
      <div className="container mx-auto px-4 py-8">
        {/* Заголовок */}
        <div className="text-center mb-12">
          <h1 className="text-4xl font-bold mb-4">
            🔍 Демонстрация всех вариантов просмотра
          </h1>
          <p className="text-lg text-base-content/70">
            Все способы отображения деталей объявлений и товаров в системе
          </p>
        </div>

        {/* 1. Основные страницы деталей */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
            <Layout className="w-6 h-6" />
            Основные страницы деталей
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <Link
              href={`/${locale}/marketplace/${marketplaceListingId}`}
              className="card bg-base-200 hover:bg-base-300 transition-colors"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">📄 Страница объявления</h3>
                <p className="text-sm text-base-content/70">
                  /marketplace/[id] - Полная информация об объявлении
                </p>
                <div className="badge badge-primary">
                  ID: {marketplaceListingId}
                </div>
              </div>
            </Link>

            <Link
              href={`/${locale}/storefronts/${storefrontSlug}/products/${storefrontProductId}`}
              className="card bg-base-200 hover:bg-base-300 transition-colors"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">🏪 Товар витрины</h3>
                <p className="text-sm text-base-content/70">
                  /storefronts/[slug]/products/[id]
                </p>
                <div className="badge badge-secondary">Витрина</div>
              </div>
            </Link>

            <Link
              href={`/${locale}/profile/orders/${orderId}`}
              className="card bg-base-200 hover:bg-base-300 transition-colors"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">📦 Детали заказа</h3>
                <p className="text-sm text-base-content/70">
                  /profile/orders/[id] - Информация о заказе
                </p>
                <div className="badge badge-info">Заказ #{orderId}</div>
              </div>
            </Link>

            <Link
              href={`/${locale}/profile/orders`}
              className="card bg-base-200 hover:bg-base-300 transition-colors"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">👤 Все заказы</h3>
                <p className="text-sm text-base-content/70">
                  /profile/orders - Список заказов
                </p>
                <div className="badge badge-success">Профиль</div>
              </div>
            </Link>

            <Link
              href={`/${locale}/marketplace/${marketplaceListingId}/buy`}
              className="card bg-base-200 hover:bg-base-300 transition-colors"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">💳 Страница покупки</h3>
                <p className="text-sm text-base-content/70">
                  /marketplace/[id]/buy
                </p>
                <div className="badge badge-warning">Оформление</div>
              </div>
            </Link>

            <Link
              href={`/${locale}/map`}
              className="card bg-base-200 hover:bg-base-300 transition-colors"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">🗺️ Карта с объявлениями</h3>
                <p className="text-sm text-base-content/70">
                  /map - Маркеры и кластеры
                </p>
                <div className="badge badge-accent">Карта</div>
              </div>
            </Link>
          </div>
        </section>

        {/* 2. Модальные окна и попапы */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
            <Maximize2 className="w-6 h-6" />
            Модальные окна и попапы
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <button
              onClick={() => setShowQuickView(true)}
              className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">👁️ QuickView</h3>
                <p className="text-sm text-base-content/70">
                  Быстрый просмотр товара без перехода
                </p>
                <div className="badge badge-primary">Нажмите для демо</div>
              </div>
            </button>

            <button
              onClick={() => setShowPriceHistory(true)}
              className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">📈 История цены</h3>
                <p className="text-sm text-base-content/70">
                  График изменения цены товара
                </p>
                <div className="badge badge-secondary">Нажмите для демо</div>
              </div>
            </button>

            <button
              onClick={() => setShowVariantModal(true)}
              className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">🎨 Выбор вариантов</h3>
                <p className="text-sm text-base-content/70">
                  Выбор размера, цвета и других опций
                </p>
                <div className="badge badge-info">Нажмите для демо</div>
              </div>
            </button>

            <button
              onClick={() => setShowImageGallery(true)}
              className="card bg-base-200 hover:bg-base-300 transition-colors cursor-pointer"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">🖼️ Галерея изображений</h3>
                <p className="text-sm text-base-content/70">
                  Полноэкранный просмотр фото
                </p>
                <div className="badge badge-success">Нажмите для демо</div>
              </div>
            </button>
          </div>
        </section>

        {/* 3. Переключатель режима отображения */}
        <section className="mb-8">
          <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
            <Settings className="w-6 h-6" />
            Настройки отображения карточек
          </h2>
          <div className="flex flex-wrap gap-4 mb-6">
            <div className="btn-group">
              <button
                className={`btn ${viewMode === 'grid' ? 'btn-active' : ''}`}
                onClick={() => setViewMode('grid')}
              >
                <Grid className="w-4 h-4 mr-2" />
                Сетка
              </button>
              <button
                className={`btn ${viewMode === 'list' ? 'btn-active' : ''}`}
                onClick={() => setViewMode('list')}
              >
                <List className="w-4 h-4 mr-2" />
                Список
              </button>
            </div>

            {viewMode === 'grid' && (
              <div className="btn-group">
                {[1, 2, 3, 4].map((cols) => (
                  <button
                    key={cols}
                    className={`btn ${gridColumns === cols ? 'btn-active' : ''}`}
                    onClick={() => setGridColumns(cols as 1 | 2 | 3 | 4)}
                  >
                    {cols} кол.
                  </button>
                ))}
              </div>
            )}
          </div>
        </section>

        {/* 4. Карточки товаров */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
            <Layers className="w-6 h-6" />
            Карточки товаров
          </h2>

          <div className="space-y-8">
            {/* EnhancedListingCard */}
            <div>
              <h3 className="text-xl font-semibold mb-4">
                EnhancedListingCard - Расширенная карточка маркетплейса
              </h3>
              <div
                className={
                  viewMode === 'grid'
                    ? `grid grid-cols-${gridColumns} gap-4`
                    : 'space-y-4'
                }
              >
                <EnhancedListingCard
                  item={sampleMarketplaceListing}
                  locale={locale}
                  viewMode={viewMode}
                  gridColumns={gridColumns}
                  onToggleFavorite={async (id) => {
                    console.log('Toggle favorite:', id);
                  }}
                />
              </div>
            </div>

            {/* UnifiedProductCard */}
            <div>
              <h3 className="text-xl font-semibold mb-4">
                UnifiedProductCard - Универсальная карточка
              </h3>
              <div
                className={
                  viewMode === 'grid'
                    ? `grid grid-cols-${gridColumns} gap-4`
                    : 'space-y-4'
                }
              >
                <UnifiedProductCard
                  product={sampleUnifiedProduct}
                  locale={locale}
                  viewMode={viewMode}
                  gridColumns={gridColumns}
                  onToggleFavorite={async (id) => {
                    console.log('Toggle favorite:', id);
                  }}
                />
              </div>
            </div>
          </div>
        </section>

        {/* 5. Дополнительные компоненты */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6 flex items-center gap-2">
            <Package className="w-6 h-6" />
            Дополнительные компоненты просмотра
          </h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <Link
              href={`/${locale}/create-listing-smart`}
              className="card bg-base-200 hover:bg-base-300 transition-colors"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">
                  ✏️ Предпросмотр при создании
                </h3>
                <p className="text-sm text-base-content/70">
                  PreviewStep - как будет выглядеть объявление
                </p>
                <div className="badge badge-warning">Создание</div>
              </div>
            </Link>

            <Link
              href={`/${locale}/storefronts/create`}
              className="card bg-base-200 hover:bg-base-300 transition-colors"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">🏪 Создание витрины</h3>
                <p className="text-sm text-base-content/70">
                  Предпросмотр витрины перед публикацией
                </p>
                <div className="badge badge-accent">Витрина</div>
              </div>
            </Link>

            <Link
              href={`/${locale}/chat?listing_id=${marketplaceListingId}`}
              className="card bg-base-200 hover:bg-base-300 transition-colors"
            >
              <div className="card-body">
                <h3 className="card-title text-lg">💬 Чат с продавцом</h3>
                <p className="text-sm text-base-content/70">
                  Обсуждение товара в чате
                </p>
                <div className="badge badge-info">Сообщения</div>
              </div>
            </Link>
          </div>
        </section>

        {/* 6. Полезные ссылки */}
        <section className="mb-12">
          <h2 className="text-2xl font-bold mb-6">
            🔗 Полезные ссылки для тестирования
          </h2>
          <div className="overflow-x-auto">
            <table className="table table-zebra">
              <thead>
                <tr>
                  <th>Тип</th>
                  <th>Описание</th>
                  <th>Ссылка</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Маркетплейс</td>
                  <td>Главная страница с товарами</td>
                  <td>
                    <Link
                      href={`/${locale}/marketplace`}
                      className="link link-primary"
                    >
                      /marketplace
                    </Link>
                  </td>
                </tr>
                <tr>
                  <td>Витрины</td>
                  <td>Список всех витрин</td>
                  <td>
                    <Link
                      href={`/${locale}/storefronts`}
                      className="link link-primary"
                    >
                      /storefronts
                    </Link>
                  </td>
                </tr>
                <tr>
                  <td>Карта</td>
                  <td>Интерактивная карта с объявлениями</td>
                  <td>
                    <Link href={`/${locale}/map`} className="link link-primary">
                      /map
                    </Link>
                  </td>
                </tr>
                <tr>
                  <td>Профиль</td>
                  <td>Личный кабинет с заказами</td>
                  <td>
                    <Link
                      href={`/${locale}/profile`}
                      className="link link-primary"
                    >
                      /profile
                    </Link>
                  </td>
                </tr>
                <tr>
                  <td>Корзина</td>
                  <td>Корзина покупок</td>
                  <td>
                    <Link
                      href={`/${locale}/cart`}
                      className="link link-primary"
                    >
                      /cart
                    </Link>
                  </td>
                </tr>
                <tr>
                  <td>Избранное</td>
                  <td>Сохраненные товары</td>
                  <td>
                    <Link
                      href={`/${locale}/favorites`}
                      className="link link-primary"
                    >
                      /favorites
                    </Link>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </section>
      </div>

      {/* Модальные окна */}
      {showQuickView && (
        <QuickView
          isOpen={showQuickView}
          onClose={() => setShowQuickView(false)}
          product={mockQuickViewProduct}
          onAddToCart={() => {
            console.log('Add to cart from QuickView');
            setShowQuickView(false);
          }}
          onContact={() => {
            console.log('Contact seller from QuickView');
            setShowQuickView(false);
          }}
        />
      )}

      {showPriceHistory && (
        <PriceHistoryModal
          isOpen={showPriceHistory}
          onClose={() => setShowPriceHistory(false)}
          listingId={marketplaceListingId}
          currentPrice={500000}
        />
      )}

      {showVariantModal && (
        <VariantSelectionModal
          isOpen={showVariantModal}
          onClose={() => setShowVariantModal(false)}
          productId={1}
          productName="iPhone 15 Pro"
          productImage="http://localhost:9000/listings/109/1754253315899650123.jpg"
          storefrontSlug="agenstvo"
          basePrice={145000}
          baseCurrency="RSD"
          onAddToCart={(variant, quantity) => {
            console.log(
              'Add to cart with variant:',
              variant,
              'quantity:',
              quantity
            );
            setShowVariantModal(false);
          }}
        />
      )}

      {showImageGallery && (
        <ImageGallery
          images={[
            'http://localhost:9000/listings/106/1754172685620502083.jpg',
            'http://localhost:9000/listings/109/1754253315899650123.jpg',
            'http://localhost:9000/listings/110/1754410715141164922.jpg',
          ]}
          currentIndex={0}
          onClose={() => setShowImageGallery(false)}
        />
      )}
    </div>
  );
}
