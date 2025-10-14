'use client';

import { useState, useEffect } from 'react';
import UniversalListingCard from '@/components/universal/cards/UniversalListingCard';
import UniversalFilters from '@/components/universal/filters/UniversalFilters';
import UniversalCreditCalculator from '@/components/universal/calculators/UniversalCreditCalculator';
import RecommendationsEngine from '@/components/universal/recommendations/RecommendationsEngine';
import { useAppDispatch } from '@/store/hooks';
import { initializeCompare } from '@/store/slices/universalCompareSlice';

export default function TestUniversalPage() {
  const dispatch = useAppDispatch();
  const [filters, setFilters] = useState<Record<string, any>>({});
  const [selectedCategory, setSelectedCategory] = useState('cars');
  const [selectedLayout, setSelectedLayout] = useState<'grid' | 'list'>('grid');

  useEffect(() => {
    // Инициализация сравнения при загрузке
    dispatch(initializeCompare());
  }, [dispatch]);

  // Тестовые данные для разных категорий
  const testListings = {
    cars: {
      id: 1,
      title: 'Toyota Camry 2020 Executive',
      price: 25000,
      currency: '€',
      images: ['/images/car-placeholder.jpg'],
      location: {
        city: 'Belgrade',
        district: 'New Belgrade',
      },
      category: 'cars',
      categorySlug: 'cars',
      createdAt: '2025-09-27T08:00:00.000Z', // 2 часа назад (фиксированное время)
      customFields: [
        { label: 'Year', value: '2020' },
        { label: 'Mileage', value: '45,000 km' },
        { label: 'Fuel', value: 'Gasoline' },
        { label: 'Transmission', value: 'Automatic' },
      ],
      badges: [
        { type: 'new' as const, label: 'New' },
        { type: 'discount' as const, label: 'Sale', value: '-15%' },
      ],
      stats: {
        views: 234,
        favorites: 12,
        calls: 5,
      },
      attributes: {
        year: 2020,
        make: 'Toyota',
        model: 'Camry',
        mileage: 45000,
        fuelType: 'gasoline',
        transmission: 'automatic',
      },
    },
    real_estate: {
      id: 2,
      title: 'Modern Apartment in City Center',
      price: 150000,
      currency: '€',
      images: ['/images/property-placeholder.jpg'],
      location: {
        city: 'Novi Sad',
        district: 'Center',
      },
      category: 'real_estate',
      categorySlug: 'real-estate',
      createdAt: '2025-09-26T10:00:00.000Z', // 1 день назад (фиксированное время)
      customFields: [
        { label: 'Area', value: '85 m²' },
        { label: 'Rooms', value: '3' },
        { label: 'Floor', value: '5/10' },
        { label: 'Year', value: '2019' },
      ],
      badges: [
        { type: 'verified' as const, label: 'Verified' },
        { type: 'recommended' as const, label: 'Top' },
      ],
      stats: {
        views: 567,
        favorites: 34,
      },
    },
    electronics: {
      id: 3,
      title: 'iPhone 14 Pro Max 256GB',
      price: 1200,
      currency: '€',
      images: ['/images/product-placeholder.jpg'],
      location: {
        city: 'Kragujevac',
      },
      category: 'electronics',
      categorySlug: 'electronics',
      createdAt: '2025-09-27T10:00:00.000Z', // сегодня (фиксированное время)
      customFields: [
        { label: 'Brand', value: 'Apple' },
        { label: 'Storage', value: '256 GB' },
        { label: 'Condition', value: 'New' },
        { label: 'Warranty', value: '24 months' },
      ],
      badges: [{ type: 'urgent' as const, label: 'Urgent' }],
      stats: {
        views: 89,
        favorites: 7,
      },
    },
  };

  const currentListing =
    testListings[selectedCategory as keyof typeof testListings] ||
    testListings.cars;

  return (
    <div className="min-h-screen bg-base-100">
      <div className="container mx-auto p-4 space-y-8">
        <div className="bg-base-200 p-6 rounded-lg">
          <h1 className="text-3xl font-bold mb-4">
            🧪 Тестирование универсальных компонентов
          </h1>
          <p className="text-base-content/70">
            Эта страница создана для проверки всех универсальных компонентов
            маркетплейса
          </p>
        </div>

        {/* Контролы */}
        <div className="bg-base-200 p-4 rounded-lg flex flex-wrap gap-4">
          <select
            className="select select-bordered"
            value={selectedCategory}
            onChange={(e) => setSelectedCategory(e.target.value)}
          >
            <option value="cars">🚗 Автомобили</option>
            <option value="real_estate">🏠 Недвижимость</option>
            <option value="electronics">📱 Электроника</option>
          </select>

          <select
            className="select select-bordered"
            value={selectedLayout}
            onChange={(e) =>
              setSelectedLayout(e.target.value as 'grid' | 'list')
            }
          >
            <option value="grid">Grid Layout</option>
            <option value="list">List Layout</option>
          </select>
        </div>

        {/* 1. UniversalListingCard */}
        <section className="space-y-4">
          <div className="divider">
            <h2 className="text-2xl font-semibold">1️⃣ UniversalListingCard</h2>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div>
              <h3 className="text-lg mb-3 font-medium">
                📊 {selectedLayout === 'grid' ? 'Grid' : 'List'} Layout
              </h3>
              <UniversalListingCard
                data={currentListing}
                type={selectedCategory as any}
                layout={selectedLayout}
                showBadges={true}
                showFavorite={true}
                showCompare={true}
                showStats={true}
                onQuickView={() => alert('Quick View clicked!')}
              />
            </div>

            <div>
              <h3 className="text-lg mb-3 font-medium">📋 Данные карточки</h3>
              <div className="bg-base-300 p-4 rounded-lg overflow-auto max-h-96">
                <pre className="text-xs">
                  {JSON.stringify(currentListing, null, 2)}
                </pre>
              </div>
            </div>
          </div>
        </section>

        {/* 2. UniversalFilters */}
        <section className="space-y-4">
          <div className="divider">
            <h2 className="text-2xl font-semibold">2️⃣ UniversalFilters</h2>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            <div className="lg:col-span-1">
              <h3 className="text-lg mb-3 font-medium">
                🎛️ Фильтры ({selectedCategory})
              </h3>
              <div className="bg-base-200 p-4 rounded-lg">
                <UniversalFilters
                  category={selectedCategory}
                  filters={filters}
                  onFiltersChange={setFilters}
                  layout="vertical"
                  config={{
                    showPriceRange: true,
                    showCondition: true,
                    showLocation: true,
                  }}
                />
              </div>
            </div>

            <div className="lg:col-span-2">
              <h3 className="text-lg mb-3 font-medium">📊 Активные фильтры</h3>
              <div className="bg-base-300 p-4 rounded-lg">
                <pre className="text-sm overflow-auto">
                  {Object.keys(filters).length > 0
                    ? JSON.stringify(filters, null, 2)
                    : 'Фильтры не выбраны'}
                </pre>
              </div>

              {Object.keys(filters).length > 0 && (
                <button
                  className="btn btn-error btn-sm mt-3"
                  onClick={() => setFilters({})}
                >
                  Очистить все фильтры
                </button>
              )}
            </div>
          </div>
        </section>

        {/* 3. UniversalCreditCalculator */}
        <section className="space-y-4">
          <div className="divider">
            <h2 className="text-2xl font-semibold">
              3️⃣ UniversalCreditCalculator
            </h2>
          </div>

          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <div>
              <h3 className="text-lg mb-3 font-medium">
                💳 Калькулятор для {selectedCategory}
              </h3>
              <UniversalCreditCalculator
                price={currentListing.price}
                category={selectedCategory}
                onApply={(calculation) => {
                  console.log('Credit calculation:', calculation);
                  alert(`
                    Расчет кредита:
                    - Ежемесячный платеж: €${calculation.monthlyPayment.toFixed(2)}
                    - Общая сумма: €${calculation.totalPayment.toFixed(2)}
                    - Переплата: €${calculation.totalInterest.toFixed(2)}
                  `);
                }}
              />
            </div>

            <div>
              <h3 className="text-lg mb-3 font-medium">ℹ️ Информация</h3>
              <div className="space-y-3">
                <div className="alert alert-info">
                  <span>
                    Калькулятор автоматически адаптируется под тип товара
                  </span>
                </div>
                <div className="alert">
                  <span>
                    <strong>Типы калькулятора:</strong>
                    <br />
                    • cars - Автокредит (до 84 месяцев)
                    <br />
                    • real_estate - Ипотека (до 360 месяцев)
                    <br />• electronics - Рассрочка (0% до 24 месяцев)
                  </span>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* 4. RecommendationsEngine */}
        <section className="space-y-4">
          <div className="divider">
            <h2 className="text-2xl font-semibold">4️⃣ RecommendationsEngine</h2>
          </div>

          <div className="space-y-8">
            {/* Similar Items */}
            <div>
              <RecommendationsEngine
                type="similar"
                category={selectedCategory}
                currentItemId={currentListing.id}
                limit={4}
                layout="grid"
                showTitle={true}
                showDescription={true}
              />
            </div>

            {/* Trending */}
            <div>
              <RecommendationsEngine
                type="trending"
                category={selectedCategory}
                limit={4}
                layout="carousel"
                showTitle={true}
                showDescription={false}
              />
            </div>

            {/* Personal */}
            <div>
              <RecommendationsEngine
                type="personal"
                category={selectedCategory}
                userId={1}
                limit={3}
                layout="list"
                showTitle={true}
                showDescription={true}
              />
            </div>
          </div>
        </section>

        {/* 5. Статус компонентов */}
        <section className="space-y-4">
          <div className="divider">
            <h2 className="text-2xl font-semibold">5️⃣ Статус компонентов</h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <div className="card bg-success text-success-content">
              <div className="card-body">
                <h3 className="card-title">✅ Готовые компоненты</h3>
                <ul className="text-sm space-y-1">
                  <li>• UniversalListingCard</li>
                  <li>• UniversalFilters</li>
                  <li>• UniversalCreditCalculator</li>
                  <li>• RecommendationsEngine</li>
                  <li>• universalCompareSlice</li>
                </ul>
              </div>
            </div>

            <div className="card bg-warning text-warning-content">
              <div className="card-body">
                <h3 className="card-title">⏳ В процессе</h3>
                <ul className="text-sm space-y-1">
                  <li>• Интеграция в страницы</li>
                  <li>• API endpoints</li>
                  <li>• Реальные данные</li>
                </ul>
              </div>
            </div>

            <div className="card bg-info text-info-content">
              <div className="card-body">
                <h3 className="card-title">📋 Планируется</h3>
                <ul className="text-sm space-y-1">
                  <li>• VIN декодер</li>
                  <li>• Сервисная книжка</li>
                  <li>• 360° просмотр</li>
                  <li>• ML рекомендации</li>
                </ul>
              </div>
            </div>
          </div>
        </section>

        {/* Footer */}
        <div className="mt-12 p-6 bg-base-200 rounded-lg text-center">
          <p className="text-base-content/70">
            Тестовая страница универсальных компонентов маркетплейса
            <br />
            Создана: 27.09.2025 | Версия: 1.0
          </p>
        </div>
      </div>
    </div>
  );
}
