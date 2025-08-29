'use client';

import React, { useState, useRef } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';
import { AnimatedSection } from '@/components/ui/AnimatedSection';
import Link from 'next/link';
import configManager from '@/config';

const IdealHomepageV2 = () => {
  const [selectedCategory, setSelectedCategory] = useState<string | null>(null);
  const [searchFocused, setSearchFocused] = useState(false);
  const categoriesRef = useRef<HTMLDivElement>(null);

  // Категории с иконками как в Avito
  const mainCategories = [
    {
      id: 'transport',
      name: 'Транспорт',
      icon: '🚗',
      count: '45K+',
      color: 'from-blue-500 to-blue-600',
    },
    {
      id: 'realestate',
      name: 'Недвижимость',
      icon: '🏠',
      count: '28K+',
      color: 'from-green-500 to-green-600',
    },
    {
      id: 'jobs',
      name: 'Работа',
      icon: '💼',
      count: '12K+',
      color: 'from-purple-500 to-purple-600',
    },
    {
      id: 'services',
      name: 'Услуги',
      icon: '🛠️',
      count: '35K+',
      color: 'from-orange-500 to-orange-600',
    },
    {
      id: 'electronics',
      name: 'Электроника',
      icon: '💻',
      count: '67K+',
      color: 'from-red-500 to-red-600',
    },
    {
      id: 'fashion',
      name: 'Одежда',
      icon: '👕',
      count: '89K+',
      color: 'from-pink-500 to-pink-600',
    },
    {
      id: 'home',
      name: 'Для дома',
      icon: '🏡',
      count: '54K+',
      color: 'from-indigo-500 to-indigo-600',
    },
    {
      id: 'hobby',
      name: 'Хобби',
      icon: '🎨',
      count: '23K+',
      color: 'from-yellow-500 to-yellow-600',
    },
    {
      id: 'pets',
      name: 'Животные',
      icon: '🐕',
      count: '15K+',
      color: 'from-teal-500 to-teal-600',
    },
    {
      id: 'business',
      name: 'Для бизнеса',
      icon: '📊',
      count: '8K+',
      color: 'from-gray-500 to-gray-600',
    },
  ];

  // Популярные поиски
  const popularSearches = [
    'iPhone 15',
    'PlayStation 5',
    'Квартира в центре',
    'BMW X5',
    'MacBook Pro',
    'Работа водителем',
    'Диван',
    'Электросамокат',
    'Кроссовки Nike',
  ];

  // Товары с рейтингами как в Wildberries
  const featuredProducts = [
    {
      id: 1,
      title: 'iPhone 14 Pro Max 256GB',
      price: 899,
      oldPrice: 1199,
      discount: 25,
      image: configManager.buildImageUrl('/listings/7/1753007242863504454.jpg'),
      rating: 4.8,
      reviews: 1234,
      isNew: true,
      isBestseller: true,
    },
    {
      id: 2,
      title: 'Квартира 2-комн, 65м², центр',
      price: 85000,
      priceUnit: '',
      image: configManager.buildImageUrl('/listings/8/1753097303704349399.jpg'),
      rating: 4.9,
      reviews: 67,
      isPremium: true,
    },
    {
      id: 3,
      title: 'MacBook Air M2 13" 512GB',
      price: 1299,
      oldPrice: 1599,
      discount: 19,
      image: configManager.buildImageUrl(
        '/listings/17/1753268215885579893.jpg'
      ),
      rating: 4.9,
      reviews: 892,
    },
    {
      id: 4,
      title: 'AirPods Pro 2 USB-C',
      price: 249,
      oldPrice: 299,
      discount: 17,
      image: configManager.buildImageUrl(
        '/listings/19/1753351396895835946.jpg'
      ),
      rating: 4.7,
      reviews: 3421,
      isNew: true,
    },
    {
      id: 5,
      title: 'BMW X5 2019 xDrive30d',
      price: 45900,
      image: configManager.buildImageUrl(
        '/listings/36/1753721116303907551.jpg'
      ),
      rating: 5.0,
      reviews: 12,
      isPremium: true,
    },
    {
      id: 6,
      title: 'Дом 180м² с участком 15 соток',
      price: 120000,
      image: configManager.buildImageUrl(
        '/listings/27/1753572833638039456.jpg'
      ),
      rating: 4.9,
      reviews: 8,
    },
  ];

  // Коллекции как в Amazon
  const collections = [
    {
      title: 'Топ электроники',
      subtitle: 'Самые популярные гаджеты',
      items: [
        {
          name: 'Смартфоны',
          image: configManager.buildImageUrl(
            '/listings/7/1753007242863504454.jpg'
          ),
          count: '2.5K+',
        },
        {
          name: 'Ноутбуки',
          image: configManager.buildImageUrl(
            '/listings/17/1753268215885579893.jpg'
          ),
          count: '1.8K+',
        },
        {
          name: 'Наушники',
          image: configManager.buildImageUrl(
            '/listings/19/1753351396895835946.jpg'
          ),
          count: '3.2K+',
        },
        {
          name: 'Планшеты',
          image: configManager.buildImageUrl(
            '/listings/28/1753574013161901892.jpg'
          ),
          count: '980+',
        },
      ],
    },
    {
      title: 'Дом и сад',
      subtitle: 'Всё для уюта',
      items: [
        {
          name: 'Мебель',
          image: configManager.buildImageUrl(
            '/listings/29/1753575302423995244.jpg'
          ),
          count: '4.1K+',
        },
        {
          name: 'Декор',
          image: configManager.buildImageUrl(
            '/listings/25/1753550885742188000.jpg'
          ),
          count: '2.7K+',
        },
        {
          name: 'Техника',
          image: configManager.buildImageUrl(
            '/listings/20/1753428897128302370.jpg'
          ),
          count: '1.9K+',
        },
        {
          name: 'Сад',
          image: configManager.buildImageUrl(
            '/listings/26/1753554432788980038.jpg'
          ),
          count: '890+',
        },
      ],
    },
  ];

  return (
    <div className="min-h-screen bg-base-100">
      {/* Header с поиском как в Wildberries */}
      <header className="sticky top-0 z-50 bg-base-100 border-b border-base-300 shadow-sm">
        <div className="container mx-auto px-4">
          <div className="flex items-center gap-4 py-3">
            {/* Logo */}
            <SveTuLogoStatic variant="gradient" width={100} height={32} />

            {/* Кнопка каталога */}
            <button
              className="btn btn-primary btn-sm gap-2"
              onClick={() =>
                setSelectedCategory(selectedCategory ? null : 'all')
              }
            >
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
                  d="M4 6h16M4 12h16M4 18h16"
                />
              </svg>
              <span className="hidden sm:inline">Каталог</span>
            </button>

            {/* Search Bar с подсказками */}
            <div className="flex-1 relative">
              <div
                className={`form-control transition-all ${searchFocused ? 'scale-[1.02]' : ''}`}
              >
                <div className="input-group">
                  <input
                    type="text"
                    placeholder="Поиск среди 2 млн товаров..."
                    className="input input-bordered w-full"
                    onFocus={() => setSearchFocused(true)}
                    onBlur={() =>
                      setTimeout(() => setSearchFocused(false), 200)
                    }
                  />
                  <button className="btn btn-square btn-primary">
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

              {/* Выпадающие подсказки при фокусе */}
              {searchFocused && (
                <div className="absolute top-full left-0 right-0 mt-1 bg-base-100 rounded-lg shadow-xl border border-base-300 p-4 z-50">
                  <p className="text-sm font-semibold mb-2 text-base-content/60">
                    Популярные запросы
                  </p>
                  <div className="flex flex-wrap gap-2">
                    {popularSearches.map((search) => (
                      <button key={search} className="btn btn-sm btn-ghost">
                        {search}
                      </button>
                    ))}
                  </div>
                </div>
              )}
            </div>

            {/* User actions */}
            <div className="flex items-center gap-2">
              <button className="btn btn-ghost btn-circle">
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
                    d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                  />
                </svg>
              </button>
              <button className="btn btn-ghost btn-circle">
                <div className="indicator">
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
                      d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z"
                    />
                  </svg>
                  <span className="badge badge-sm badge-primary indicator-item">
                    3
                  </span>
                </div>
              </button>
              <Link href="/auth/login" className="btn btn-primary btn-sm">
                Войти
              </Link>
            </div>
          </div>
        </div>
      </header>

      {/* Категории сетка как в Avito */}
      <section className="py-4 border-b border-base-200">
        <div className="container mx-auto px-4">
          <div
            ref={categoriesRef}
            className="flex gap-3 overflow-x-auto scrollbar-hide pb-2"
          >
            {mainCategories.map((category) => (
              <button
                key={category.id}
                onClick={() => setSelectedCategory(category.id)}
                className={`flex flex-col items-center gap-1 p-3 rounded-lg transition-all hover:scale-105 min-w-[80px] ${
                  selectedCategory === category.id
                    ? 'bg-gradient-to-r ' +
                      category.color +
                      ' text-white shadow-lg'
                    : 'bg-base-200 hover:bg-base-300'
                }`}
              >
                <span className="text-2xl">{category.icon}</span>
                <span className="text-xs font-semibold whitespace-nowrap">
                  {category.name}
                </span>
                <span className="text-xs opacity-70">{category.count}</span>
              </button>
            ))}
          </div>
        </div>
      </section>

      {/* Hero Banner как в Wildberries */}
      <section className="relative">
        <div className="container mx-auto px-4 py-4">
          <AnimatedSection animation="fadeIn">
            <div className="relative h-48 md:h-64 rounded-2xl overflow-hidden bg-gradient-to-r from-purple-600 to-pink-600">
              <div className="absolute inset-0 flex items-center justify-between p-8">
                <div className="text-white">
                  <h1 className="text-2xl md:text-4xl font-bold mb-2">
                    Черная пятница уже здесь!
                  </h1>
                  <p className="text-lg md:text-xl mb-4 opacity-90">
                    Скидки до 70% на все категории
                  </p>
                  <button className="btn btn-warning btn-lg">
                    Смотреть акции
                  </button>
                </div>
                <div className="hidden md:block text-8xl">🎁</div>
              </div>
              <div className="absolute top-4 right-4 bg-yellow-400 text-black px-4 py-2 rounded-full font-bold">
                -70%
              </div>
            </div>
          </AnimatedSection>
        </div>
      </section>

      {/* Рекомендуемые товары с рейтингами как в Wildberries */}
      <section className="py-6">
        <div className="container mx-auto px-4">
          <AnimatedSection animation="fadeIn">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl md:text-2xl font-bold">🔥 Хиты продаж</h2>
              <button className="btn btn-ghost btn-sm">Все товары →</button>
            </div>
          </AnimatedSection>

          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3">
            {featuredProducts.map((product, idx) => (
              <AnimatedSection
                key={product.id}
                animation="slideUp"
                delay={idx * 0.05}
              >
                <div className="card bg-base-100 border border-base-200 hover:shadow-xl transition-all hover:-translate-y-1 group">
                  <figure className="relative h-48 overflow-hidden bg-base-200">
                    {/* Badges */}
                    <div className="absolute top-2 left-2 z-10 flex flex-col gap-1">
                      {product.isNew && (
                        <div className="badge badge-secondary badge-sm">
                          NEW
                        </div>
                      )}
                      {product.isBestseller && (
                        <div className="badge badge-warning badge-sm">ХИТ</div>
                      )}
                      {product.isPremium && (
                        <div className="badge badge-primary badge-sm">
                          PREMIUM
                        </div>
                      )}
                    </div>

                    {/* Discount */}
                    {product.discount && (
                      <div className="absolute top-2 right-2 z-10">
                        <div className="bg-error text-white rounded-lg px-2 py-1 text-sm font-bold">
                          -{product.discount}%
                        </div>
                      </div>
                    )}

                    {/* Quick actions */}
                    <div className="absolute bottom-2 right-2 z-10 opacity-0 group-hover:opacity-100 transition-opacity">
                      <button className="btn btn-circle btn-sm bg-base-100/80 backdrop-blur">
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
                            d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"
                          />
                        </svg>
                      </button>
                    </div>

                    <img
                      src={product.image}
                      alt={product.title}
                      className="w-full h-full object-cover group-hover:scale-110 transition-transform duration-300"
                    />
                  </figure>

                  <div className="card-body p-3">
                    <h3 className="text-sm font-semibold line-clamp-2">
                      {product.title}
                    </h3>

                    {/* Rating */}
                    <div className="flex items-center gap-1 text-xs">
                      <div className="flex items-center">
                        <span className="text-warning">★</span>
                        <span className="font-semibold">{product.rating}</span>
                      </div>
                      <span className="text-base-content/60">
                        ({product.reviews})
                      </span>
                    </div>

                    {/* Price */}
                    <div className="mt-2">
                      {product.oldPrice && (
                        <div className="text-xs text-base-content/50 line-through">
                          €{product.oldPrice}
                        </div>
                      )}
                      <div className="text-lg font-bold text-primary">
                        €{product.price}
                        {product.priceUnit}
                      </div>
                    </div>

                    {/* Add to cart */}
                    <button className="btn btn-primary btn-sm btn-block mt-2">
                      В корзину
                    </button>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* Коллекции как в Amazon */}
      <section className="py-6 bg-base-200">
        <div className="container mx-auto px-4">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {collections.map((collection, idx) => (
              <AnimatedSection key={idx} animation="fadeIn" delay={idx * 0.1}>
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body">
                    <h3 className="card-title">{collection.title}</h3>
                    <p className="text-sm text-base-content/60">
                      {collection.subtitle}
                    </p>

                    <div className="grid grid-cols-2 gap-3 mt-4">
                      {collection.items.map((item, itemIdx) => (
                        <button
                          key={itemIdx}
                          className="relative group overflow-hidden rounded-lg border border-base-200 hover:border-primary transition-all"
                        >
                          <figure className="h-24 overflow-hidden bg-base-200">
                            <img
                              src={item.image}
                              alt={item.name}
                              className="w-full h-full object-cover group-hover:scale-110 transition-transform"
                            />
                          </figure>
                          <div className="absolute inset-0 bg-gradient-to-t from-black/70 to-transparent flex items-end p-2">
                            <div className="text-white text-left">
                              <p className="text-sm font-semibold">
                                {item.name}
                              </p>
                              <p className="text-xs opacity-80">{item.count}</p>
                            </div>
                          </div>
                        </button>
                      ))}
                    </div>

                    <button className="btn btn-outline btn-sm mt-3">
                      Смотреть всё →
                    </button>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* Картографический сервис */}
      <section className="py-6 bg-gradient-to-br from-blue-50 to-green-50">
        <div className="container mx-auto px-4">
          <AnimatedSection animation="fadeIn">
            <div className="text-center mb-6">
              <h2 className="text-2xl font-bold mb-2">🗺️ Товары на карте</h2>
              <p className="text-base text-base-content/70">
                Найдите нужное рядом с вами или исследуйте другие районы
              </p>
            </div>
          </AnimatedSection>

          <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
            {/* Визуализация карты */}
            <div className="lg:col-span-2">
              <div className="card bg-base-100 shadow-xl h-full">
                <div className="card-body p-0">
                  <div className="relative h-[400px] lg:h-[500px] rounded-xl overflow-hidden">
                    {/* Реальная карта - статическое изображение OpenStreetMap */}
                    <iframe
                      className="absolute inset-0 w-full h-full"
                      frameBorder="0"
                      scrolling="no"
                      marginHeight={0}
                      marginWidth={0}
                      src="https://www.openstreetmap.org/export/embed.html?bbox=20.35,44.75,20.55,44.88&amp;layer=mapnik"
                      style={{ border: 0 }}
                    />

                    {/* Полупрозрачный оверлей для лучшей видимости маркеров */}
                    <div className="absolute inset-0 bg-gradient-to-br from-blue-500/10 to-green-500/10 pointer-events-none" />

                    {/* Маркер пользователя */}
                    <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 z-20">
                      <div className="relative">
                        <div className="w-5 h-5 bg-blue-500 rounded-full shadow-lg border-2 border-white">
                          <div className="absolute inset-0 bg-blue-400 rounded-full animate-ping" />
                        </div>
                      </div>
                      <div className="absolute top-8 left-1/2 -translate-x-1/2 bg-white/95 backdrop-blur px-3 py-1.5 rounded-lg shadow-lg text-sm font-medium whitespace-nowrap">
                        📍 Вы здесь
                      </div>
                    </div>

                    {/* Маркеры товаров с реалистичными локациями для Белграда */}
                    {[
                      {
                        top: '25%',
                        left: '45%',
                        price: '250€',
                        category: 'Электроника',
                        address: 'Knez Mihailova',
                      },
                      {
                        top: '35%',
                        left: '55%',
                        price: '450€',
                        category: 'Мебель',
                        address: 'Vračar',
                      },
                      {
                        top: '55%',
                        left: '40%',
                        price: '120€',
                        category: 'Одежда',
                        address: 'Novi Beograd',
                      },
                      {
                        top: '40%',
                        left: '65%',
                        price: '800€/мес',
                        category: 'Квартира',
                        address: 'Dedinje',
                      },
                      {
                        top: '65%',
                        left: '50%',
                        price: '50€',
                        category: 'Книги',
                        address: 'Zemun',
                      },
                      {
                        top: '30%',
                        left: '35%',
                        price: '350€',
                        category: 'Спорт',
                        address: 'Kalemegdan',
                      },
                      {
                        top: '45%',
                        left: '30%',
                        price: '15,000€',
                        category: 'Авто',
                        address: 'Autokomanda',
                      },
                      {
                        top: '50%',
                        left: '60%',
                        price: '90€',
                        category: 'Хобби',
                        address: 'Zvezdara',
                      },
                    ].map((marker, idx) => (
                      <div
                        key={idx}
                        className="absolute group cursor-pointer z-10"
                        style={{ top: marker.top, left: marker.left }}
                      >
                        <div className="relative">
                          {/* Пин маркера */}
                          <div className="relative">
                            <svg
                              className="w-8 h-8 text-orange-500 drop-shadow-lg transform group-hover:scale-110 transition-transform"
                              fill="currentColor"
                              viewBox="0 0 24 24"
                            >
                              <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7zm0 9.5c-1.38 0-2.5-1.12-2.5-2.5s1.12-2.5 2.5-2.5 2.5 1.12 2.5 2.5-1.12 2.5-2.5 2.5z" />
                            </svg>
                            <div className="absolute top-1 left-1/2 -translate-x-1/2 text-white text-xs font-bold">
                              {idx + 1}
                            </div>
                          </div>

                          {/* Информация при наведении */}
                          <div className="absolute bottom-10 left-1/2 -translate-x-1/2 bg-white/95 backdrop-blur px-4 py-3 rounded-lg shadow-xl opacity-0 group-hover:opacity-100 transition-all duration-300 pointer-events-none z-30 whitespace-nowrap">
                            <p className="font-bold text-lg text-orange-600">
                              {marker.price}
                            </p>
                            <p className="text-sm font-medium text-gray-800">
                              {marker.category}
                            </p>
                            <p className="text-xs text-gray-600 mt-1">
                              📍 {marker.address}
                            </p>
                          </div>
                        </div>
                      </div>
                    ))}

                    {/* Радиус поиска */}
                    <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-64 h-64 border-2 border-blue-400 border-dashed rounded-full opacity-40">
                      <div className="absolute -top-6 left-1/2 -translate-x-1/2 bg-blue-500 text-white text-xs px-2 py-1 rounded-full">
                        Радиус 2 км
                      </div>
                    </div>

                    {/* Легенда карты */}
                    <div className="absolute bottom-4 left-4 bg-white/90 backdrop-blur rounded-lg p-3 shadow-lg">
                      <p className="text-xs font-medium mb-2">Легенда:</p>
                      <div className="space-y-1">
                        <div className="flex items-center gap-2">
                          <div className="w-3 h-3 bg-blue-500 rounded-full"></div>
                          <span className="text-xs">Ваше местоположение</span>
                        </div>
                        <div className="flex items-center gap-2">
                          <svg
                            className="w-4 h-4 text-orange-500"
                            fill="currentColor"
                            viewBox="0 0 24 24"
                          >
                            <path d="M12 2C8.13 2 5 5.13 5 9c0 5.25 7 13 7 13s7-7.75 7-13c0-3.87-3.13-7-7-7z" />
                          </svg>
                          <span className="text-xs">Товары поблизости</span>
                        </div>
                      </div>
                    </div>

                    {/* Контролы карты */}
                    <div className="absolute top-4 right-4 flex flex-col gap-2">
                      <button className="btn btn-circle btn-sm bg-white/90 backdrop-blur shadow-lg hover:bg-white">
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
                            d="M12 6v6m0 0v6m0-6h6m-6 0H6"
                          />
                        </svg>
                      </button>
                      <button className="btn btn-circle btn-sm bg-white/90 backdrop-blur shadow-lg hover:bg-white">
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
                            d="M20 12H4"
                          />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            {/* Статистика и фильтры */}
            <AnimatedSection animation="slideRight">
              <div className="space-y-4">
                {/* Статистика по районам */}
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body">
                    <h3 className="card-title text-lg mb-3">
                      📊 В вашем районе
                    </h3>
                    <div className="space-y-3">
                      <div className="flex justify-between items-center">
                        <span className="text-sm">Всего объявлений</span>
                        <span className="badge badge-primary badge-lg">
                          1,234
                        </span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm">Новых сегодня</span>
                        <span className="badge badge-success">+89</span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-sm">В радиусе 5 км</span>
                        <span className="badge badge-info">567</span>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Быстрые фильтры */}
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body">
                    <h3 className="card-title text-lg mb-3">
                      ⚡ Быстрые фильтры
                    </h3>
                    <div className="flex flex-wrap gap-2">
                      <button className="btn btn-sm btn-outline">
                        До 500€
                      </button>
                      <button className="btn btn-sm btn-outline">
                        Сегодня
                      </button>
                      <button className="btn btn-sm btn-outline">С фото</button>
                      <button className="btn btn-sm btn-outline">Рядом</button>
                      <button className="btn btn-sm btn-outline">Срочно</button>
                    </div>
                  </div>
                </div>

                {/* CTA */}
                <div className="card bg-gradient-to-r from-primary to-secondary text-primary-content">
                  <div className="card-body">
                    <h3 className="card-title text-white">Исследуйте карту</h3>
                    <p className="text-sm opacity-90">
                      Найдите товары в любом районе города
                    </p>
                    <div className="card-actions justify-end mt-3">
                      <Link href="/map" className="btn btn-white btn-sm">
                        Открыть карту
                        <svg
                          className="w-4 h-4 ml-1"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M9 5l7 7-7 7"
                          />
                        </svg>
                      </Link>
                    </div>
                  </div>
                </div>
              </div>
            </AnimatedSection>
          </div>
        </div>
      </section>

      {/* История просмотров как в Avito */}
      <section className="py-6">
        <div className="container mx-auto px-4">
          <AnimatedSection animation="fadeIn">
            <div className="flex items-center justify-between mb-4">
              <h2 className="text-xl font-bold">👁️ Вы смотрели</h2>
              <button className="btn btn-ghost btn-sm">Очистить</button>
            </div>
          </AnimatedSection>

          <div className="flex gap-3 overflow-x-auto scrollbar-hide pb-2">
            {featuredProducts.slice(0, 4).map((product, idx) => (
              <AnimatedSection
                key={idx}
                animation="slideLeft"
                delay={idx * 0.05}
              >
                <div className="card bg-base-100 border border-base-200 w-40 flex-shrink-0">
                  <figure className="h-32 overflow-hidden">
                    <img
                      src={product.image}
                      alt={product.title}
                      className="w-full h-full object-cover"
                    />
                  </figure>
                  <div className="card-body p-2">
                    <p className="text-xs line-clamp-2">{product.title}</p>
                    <p className="text-sm font-bold text-primary">
                      €{product.price}
                    </p>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* Сервисы как в Avito */}
      <section className="py-6 bg-gradient-to-r from-primary/10 to-secondary/10">
        <div className="container mx-auto px-4">
          <AnimatedSection animation="fadeIn">
            <h2 className="text-xl font-bold mb-4 text-center">
              🛡️ Наши сервисы
            </h2>
          </AnimatedSection>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            {[
              {
                icon: '📦',
                title: 'Доставка',
                desc: 'По всей Сербии',
                color: 'from-blue-500 to-blue-600',
              },
              {
                icon: '💳',
                title: 'Оплата онлайн',
                desc: 'Безопасные платежи',
                color: 'from-green-500 to-green-600',
              },
              {
                icon: '🔒',
                title: 'Защита сделок',
                desc: 'Эскроу-сервис',
                color: 'from-purple-500 to-purple-600',
              },
              {
                icon: '🚗',
                title: 'Автотека',
                desc: 'Проверка авто',
                color: 'from-orange-500 to-orange-600',
              },
            ].map((service, idx) => (
              <AnimatedSection key={idx} animation="zoomIn" delay={idx * 0.1}>
                <div className="card bg-base-100 shadow-lg hover:shadow-xl transition-all cursor-pointer group">
                  <div className="card-body items-center text-center p-4">
                    <div
                      className={`w-16 h-16 rounded-full bg-gradient-to-r ${service.color} flex items-center justify-center text-3xl mb-2 group-hover:scale-110 transition-transform`}
                    >
                      {service.icon}
                    </div>
                    <h3 className="font-bold">{service.title}</h3>
                    <p className="text-xs text-base-content/60">
                      {service.desc}
                    </p>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="footer p-6 bg-base-200 text-base-content">
        <div className="container mx-auto">
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <span className="footer-title">Покупателям</span>
              <a className="link link-hover text-sm">Как купить</a>
              <a className="link link-hover text-sm">Доставка</a>
              <a className="link link-hover text-sm">Оплата</a>
              <a className="link link-hover text-sm">Возврат</a>
            </div>
            <div>
              <span className="footer-title">Продавцам</span>
              <a className="link link-hover text-sm">Как продать</a>
              <a className="link link-hover text-sm">Правила</a>
              <a className="link link-hover text-sm">Продвижение</a>
              <a className="link link-hover text-sm">Магазины</a>
            </div>
            <div>
              <span className="footer-title">О компании</span>
              <a className="link link-hover text-sm">О нас</a>
              <a className="link link-hover text-sm">Контакты</a>
              <a className="link link-hover text-sm">Блог</a>
              <a className="link link-hover text-sm">API</a>
            </div>
            <div>
              <span className="footer-title">Приложение</span>
              <div className="grid grid-flow-col gap-2">
                <button className="btn btn-sm">
                  <svg
                    className="w-4 h-4"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                  >
                    <path d="M17.05 20.28c-.98.95-2.05.8-3.08.35-1.09-.46-2.09-.48-3.24 0-1.44.62-2.2.44-3.06-.35C2.79 15.25 3.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09l.01-.01zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z" />
                  </svg>
                </button>
                <button className="btn btn-sm">
                  <svg
                    className="w-4 h-4"
                    viewBox="0 0 24 24"
                    fill="currentColor"
                  >
                    <path d="M3,20.5V3.5C3,2.91 3.34,2.39 3.84,2.15L13.69,12L3.84,21.85C3.34,21.6 3,21.09 3,20.5M16.81,15.12L6.05,21.34L14.54,12.85L16.81,15.12M20.16,10.81C20.5,11.08 20.75,11.5 20.75,12C20.75,12.5 20.53,12.9 20.18,13.18L17.89,14.5L15.39,12L17.89,9.5L20.16,10.81M6.05,2.66L16.81,8.88L14.54,11.15L6.05,2.66Z" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
          <div className="divider"></div>
          <div className="text-center text-sm text-base-content/60">
            © 2024 Sve Tu. Все права защищены.
          </div>
        </div>
      </footer>
    </div>
  );
};

export default IdealHomepageV2;
