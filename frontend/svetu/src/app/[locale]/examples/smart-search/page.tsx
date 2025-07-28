'use client';

import React, { useState } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

const SmartSearch = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [activeFilters, setActiveFilters] = useState<string[]>([]);
  const [priceRange, setPriceRange] = useState([0, 1000]);
  const [sortBy, setSortBy] = useState('relevance');
  const [viewMode, setViewMode] = useState<'grid' | 'list'>('grid');

  const categories = [
    { id: 'electronics', name: 'Электроника', icon: '💻', count: 234 },
    { id: 'realestate', name: 'Недвижимость', icon: '🏠', count: 156 },
    { id: 'auto', name: 'Автомобили', icon: '🚗', count: 89 },
    { id: 'fashion', name: 'Одежда', icon: '👕', count: 412 },
    { id: 'services', name: 'Услуги', icon: '🛠️', count: 178 },
    { id: 'hobby', name: 'Хобби', icon: '🎨', count: 267 },
  ];

  const filters = {
    condition: ['Новое', 'Б/у', 'На запчасти'],
    delivery: ['Самовывоз', 'Доставка', 'Почта'],
    payment: ['Наличные', 'Карта', 'Перевод'],
  };

  const searchResults = [
    {
      id: 1,
      title: 'iPhone 14 Pro Max 256GB',
      price: 899,
      location: 'Белград',
      image:
        '/api/minio/download?fileName=listings/0a47e66f-d8da-459f-a2ba-8e2b85ae0163/38ad29e6-7b07-4bfc-9db2-d965cb6b966f.jpg',
      rating: 4.8,
      reviews: 23,
      isPromoted: true,
    },
    {
      id: 2,
      title: 'Квартира 2-комнатная, центр',
      price: 650,
      location: 'Нови Сад',
      image:
        '/api/minio/download?fileName=listings/0c1fc30d-5d84-485f-a86a-5c5dc37f8b97/4b8b8e48-ddd8-4c04-ad8e-00c4b4d10d26.jpg',
      rating: 4.9,
      reviews: 12,
    },
    {
      id: 3,
      title: 'MacBook Pro M2 13"',
      price: 1299,
      location: 'Ниш',
      image:
        '/api/minio/download?fileName=listings/0c91d2f7-53f7-4bff-87fe-d7e82dc3e2f0/3b26f07f-c5d6-4ff7-ba56-06ec69bb7f4d.jpg',
      rating: 5.0,
      reviews: 8,
    },
  ];

  const recentSearches = ['iPhone', 'Квартира центр', 'MacBook', 'Велосипед'];
  const popularTags = ['Срочно', 'Новое', 'С гарантией', 'Торг', 'Обмен'];

  const handleFilterToggle = (filter: string) => {
    setActiveFilters((prev) =>
      prev.includes(filter)
        ? prev.filter((f) => f !== filter)
        : [...prev, filter]
    );
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      {/* Header */}
      <div className="navbar bg-base-100 shadow-lg sticky top-0 z-50">
        <div className="navbar-start">
          <SveTuLogoStatic variant="gradient" width={120} height={40} />
        </div>
        <div className="navbar-center flex-1 px-4">
          <div className="form-control w-full max-w-2xl">
            <div className="input-group">
              <input
                type="text"
                placeholder="Поиск товаров и услуг..."
                className="input input-bordered w-full"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
              <button className="btn btn-primary">
                <svg
                  className="w-5 h-5"
                  fill="none"
                  stroke="currentColor"
                  viewBox="0 0 24 24"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                  />
                </svg>
              </button>
            </div>
          </div>
        </div>
        <div className="navbar-end">
          <div className="btn-group">
            <button
              className={`btn btn-sm ${viewMode === 'grid' ? 'btn-active' : ''}`}
              onClick={() => setViewMode('grid')}
            >
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path d="M5 3a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2V5a2 2 0 00-2-2H5zM5 11a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2v-2a2 2 0 00-2-2H5zM11 5a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V5zM13 11a2 2 0 00-2 2v2a2 2 0 002 2h2a2 2 0 002-2v-2a2 2 0 00-2-2h-2z" />
              </svg>
            </button>
            <button
              className={`btn btn-sm ${viewMode === 'list' ? 'btn-active' : ''}`}
              onClick={() => setViewMode('list')}
            >
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z"
                  clipRule="evenodd"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6">
        {/* Quick Search Suggestions */}
        <AnimatedSection animation="fadeIn">
          <div className="mb-6">
            <div className="flex items-center gap-4 flex-wrap">
              <span className="text-sm text-base-content/60">Недавние:</span>
              {recentSearches.map((search, idx) => (
                <button
                  key={idx}
                  className="btn btn-sm btn-ghost"
                  onClick={() => setSearchQuery(search)}
                >
                  {search}
                </button>
              ))}
              <div className="divider divider-horizontal"></div>
              <span className="text-sm text-base-content/60">Популярное:</span>
              {popularTags.map((tag, idx) => (
                <span
                  key={idx}
                  className="badge badge-outline cursor-pointer hover:badge-primary"
                >
                  {tag}
                </span>
              ))}
            </div>
          </div>
        </AnimatedSection>

        <div className="flex gap-6">
          {/* Filters Sidebar */}
          <AnimatedSection animation="slideLeft" className="w-80">
            <div className="card bg-base-100 shadow-xl sticky top-24">
              <div className="card-body">
                <h3 className="card-title mb-4">🔍 Фильтры</h3>

                {/* Categories */}
                <div className="mb-6">
                  <h4 className="font-semibold mb-3">Категории</h4>
                  <div className="space-y-2">
                    {categories.map((cat) => (
                      <label
                        key={cat.id}
                        className="flex items-center gap-3 cursor-pointer hover:bg-base-200 p-2 rounded-lg"
                      >
                        <input
                          type="checkbox"
                          className="checkbox checkbox-primary checkbox-sm"
                        />
                        <span className="text-xl">{cat.icon}</span>
                        <span className="flex-1">{cat.name}</span>
                        <span className="badge badge-sm">{cat.count}</span>
                      </label>
                    ))}
                  </div>
                </div>

                <div className="divider"></div>

                {/* Price Range */}
                <div className="mb-6">
                  <h4 className="font-semibold mb-3">Цена (€)</h4>
                  <div className="flex gap-2 mb-3">
                    <input
                      type="number"
                      placeholder="От"
                      className="input input-bordered input-sm w-full"
                      value={priceRange[0]}
                      onChange={(e) =>
                        setPriceRange([Number(e.target.value), priceRange[1]])
                      }
                    />
                    <input
                      type="number"
                      placeholder="До"
                      className="input input-bordered input-sm w-full"
                      value={priceRange[1]}
                      onChange={(e) =>
                        setPriceRange([priceRange[0], Number(e.target.value)])
                      }
                    />
                  </div>
                  <input
                    type="range"
                    min="0"
                    max="2000"
                    value={priceRange[1]}
                    className="range range-primary range-sm"
                    onChange={(e) =>
                      setPriceRange([priceRange[0], Number(e.target.value)])
                    }
                  />
                </div>

                <div className="divider"></div>

                {/* Other Filters */}
                {Object.entries(filters).map(([key, values]) => (
                  <div key={key} className="mb-4">
                    <h4 className="font-semibold mb-2 capitalize">
                      {key === 'condition'
                        ? 'Состояние'
                        : key === 'delivery'
                          ? 'Доставка'
                          : 'Оплата'}
                    </h4>
                    <div className="flex flex-wrap gap-2">
                      {values.map((value) => (
                        <button
                          key={value}
                          className={`badge ${activeFilters.includes(value) ? 'badge-primary' : 'badge-outline'} cursor-pointer`}
                          onClick={() => handleFilterToggle(value)}
                        >
                          {value}
                        </button>
                      ))}
                    </div>
                  </div>
                ))}

                <button className="btn btn-primary btn-block mt-4">
                  Применить фильтры
                </button>
                <button className="btn btn-ghost btn-block">
                  Сбросить все
                </button>
              </div>
            </div>
          </AnimatedSection>

          {/* Search Results */}
          <div className="flex-1">
            <AnimatedSection animation="fadeIn">
              <div className="flex justify-between items-center mb-6">
                <div>
                  <h2 className="text-2xl font-bold">
                    Найдено: 1,234 объявления
                  </h2>
                  <p className="text-sm text-base-content/60">
                    по запросу "{searchQuery || 'все товары'}"
                  </p>
                </div>
                <select
                  className="select select-bordered select-sm"
                  value={sortBy}
                  onChange={(e) => setSortBy(e.target.value)}
                >
                  <option value="relevance">По релевантности</option>
                  <option value="price-asc">Цена: по возрастанию</option>
                  <option value="price-desc">Цена: по убыванию</option>
                  <option value="date">Сначала новые</option>
                  <option value="rating">По рейтингу</option>
                </select>
              </div>
            </AnimatedSection>

            {/* Results Grid */}
            <div
              className={
                viewMode === 'grid'
                  ? 'grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6'
                  : 'space-y-4'
              }
            >
              {searchResults.map((item, idx) => (
                <AnimatedSection
                  key={item.id}
                  animation="slideUp"
                  delay={idx * 0.1}
                >
                  {viewMode === 'grid' ? (
                    <div className="card bg-base-100 shadow-xl hover:shadow-2xl transition-all hover:-translate-y-1">
                      {item.isPromoted && (
                        <div className="badge badge-warning absolute top-2 right-2 z-10">
                          Промо
                        </div>
                      )}
                      <figure className="relative h-48">
                        <img
                          src={item.image}
                          alt={item.title}
                          className="w-full h-full object-cover"
                        />
                      </figure>
                      <div className="card-body">
                        <h3 className="card-title text-lg">{item.title}</h3>
                        <div className="flex items-center gap-2 text-sm">
                          <svg
                            className="w-4 h-4"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                            />
                            <path
                              strokeLinecap="round"
                              strokeLinejoin="round"
                              strokeWidth={2}
                              d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"
                            />
                          </svg>
                          <span>{item.location}</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <div className="rating rating-sm">
                            {[1, 2, 3, 4, 5].map((star) => (
                              <input
                                key={star}
                                type="radio"
                                className="mask mask-star-2 bg-orange-400"
                                checked={star <= Math.floor(item.rating)}
                                readOnly
                              />
                            ))}
                          </div>
                          <span className="text-sm">({item.reviews})</span>
                        </div>
                        <div className="card-actions justify-between items-center mt-4">
                          <span className="text-2xl font-bold">
                            €{item.price}
                          </span>
                          <button className="btn btn-primary btn-sm">
                            Подробнее
                          </button>
                        </div>
                      </div>
                    </div>
                  ) : (
                    <div className="card bg-base-100 shadow-xl">
                      <div className="card-body">
                        <div className="flex gap-4">
                          <figure className="w-32 h-32 flex-shrink-0">
                            <img
                              src={item.image}
                              alt={item.title}
                              className="w-full h-full object-cover rounded-lg"
                            />
                          </figure>
                          <div className="flex-1">
                            <div className="flex justify-between items-start">
                              <div>
                                <h3 className="text-xl font-bold">
                                  {item.title}
                                </h3>
                                <div className="flex items-center gap-4 mt-2 text-sm text-base-content/60">
                                  <span className="flex items-center gap-1">
                                    <svg
                                      className="w-4 h-4"
                                      fill="none"
                                      stroke="currentColor"
                                      viewBox="0 0 24 24"
                                    >
                                      <path
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth={2}
                                        d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"
                                      />
                                    </svg>
                                    {item.location}
                                  </span>
                                  <div className="flex items-center gap-1">
                                    <div className="rating rating-sm">
                                      {[1, 2, 3, 4, 5].map((star) => (
                                        <input
                                          key={star}
                                          type="radio"
                                          className="mask mask-star-2 bg-orange-400"
                                          checked={
                                            star <= Math.floor(item.rating)
                                          }
                                          readOnly
                                        />
                                      ))}
                                    </div>
                                    <span>({item.reviews})</span>
                                  </div>
                                </div>
                              </div>
                              <div className="text-right">
                                <div className="text-2xl font-bold">
                                  €{item.price}
                                </div>
                                {item.isPromoted && (
                                  <span className="badge badge-warning badge-sm">
                                    Промо
                                  </span>
                                )}
                              </div>
                            </div>
                            <div className="mt-4 flex gap-2">
                              <button className="btn btn-primary btn-sm">
                                Подробнее
                              </button>
                              <button className="btn btn-ghost btn-sm">
                                В избранное
                              </button>
                            </div>
                          </div>
                        </div>
                      </div>
                    </div>
                  )}
                </AnimatedSection>
              ))}
            </div>

            {/* Pagination */}
            <AnimatedSection animation="fadeIn" delay={0.5}>
              <div className="flex justify-center mt-8">
                <div className="join">
                  <button className="join-item btn">«</button>
                  <button className="join-item btn btn-active">1</button>
                  <button className="join-item btn">2</button>
                  <button className="join-item btn">3</button>
                  <button className="join-item btn btn-disabled">...</button>
                  <button className="join-item btn">41</button>
                  <button className="join-item btn">»</button>
                </div>
              </div>
            </AnimatedSection>
          </div>
        </div>
      </div>
    </div>
  );
};

export default SmartSearch;
