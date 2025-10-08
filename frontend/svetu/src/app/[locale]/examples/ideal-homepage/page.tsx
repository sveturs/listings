'use client';

import React, { useState } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';
import { AnimatedSection } from '@/components/ui/AnimatedSection';
import Link from 'next/link';
import configManager from '@/config';

const IdealHomepage = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');

  const categories = [
    { id: 'all', name: 'Все категории', icon: '🔍' },
    { id: 'realestate', name: 'Недвижимость', icon: '🏠' },
    { id: 'auto', name: 'Автомобили', icon: '🚗' },
    { id: 'electronics', name: 'Электроника', icon: '💻' },
    { id: 'fashion', name: 'Одежда', icon: '👕' },
    { id: 'services', name: 'Услуги', icon: '🛠️' },
    { id: 'hobby', name: 'Хобби', icon: '🎨' },
    { id: 'jobs', name: 'Работа', icon: '💼' },
  ];

  const popularProducts = [
    {
      id: 1,
      title: 'iPhone 14 Pro Max 256GB',
      price: 899,
      location: 'Белград',
      image: configManager.buildImageUrl('/listings/7/1753007242863504454.jpg'),
      isNew: true,
      views: 234,
    },
    {
      id: 2,
      title: 'Квартира 2-комнатная, центр',
      price: 650,
      priceUnit: '/месяц',
      location: 'Нови Сад',
      image: configManager.buildImageUrl('/listings/8/1753097303704349399.jpg'),
      isPromoted: true,
      views: 567,
    },
    {
      id: 3,
      title: 'MacBook Pro M2 13"',
      price: 1299,
      location: 'Ниш',
      image: configManager.buildImageUrl(
        '/listings/17/1753268215885579893.jpg'
      ),
      discount: 10,
      oldPrice: 1449,
      views: 189,
    },
    {
      id: 4,
      title: 'AirPods Pro 2',
      price: 249,
      location: 'Белград',
      image: configManager.buildImageUrl(
        '/listings/19/1753351396895835946.jpg'
      ),
      isNew: true,
      views: 412,
    },
    {
      id: 5,
      title: 'BMW X5 2019',
      price: 45900,
      location: 'Белград',
      image: configManager.buildImageUrl(
        '/listings/36/1753721116303907551.jpg'
      ),
      isPromoted: true,
      views: 892,
    },
    {
      id: 6,
      title: 'Дом с участком',
      price: 120000,
      location: 'Земун',
      image: configManager.buildImageUrl(
        '/listings/27/1753572833638039456.jpg'
      ),
      isNew: true,
      views: 445,
    },
    {
      id: 7,
      title: 'Samsung Galaxy S23',
      price: 1099,
      location: 'Нови Сад',
      image: configManager.buildImageUrl(
        '/listings/28/1753574013161901892.jpg'
      ),
      discount: 15,
      oldPrice: 1299,
      views: 678,
    },
    {
      id: 8,
      title: 'Офисное кресло',
      price: 799,
      location: 'Белград',
      image: configManager.buildImageUrl(
        '/listings/29/1753575302423995244.jpg'
      ),
      views: 234,
    },
  ];

  const features = [
    {
      icon: '🤖',
      title: 'AI-помощник',
      description: 'Умный анализ фото для создания объявлений',
      color: 'from-violet-500 to-purple-500',
    },
    {
      icon: '🔒',
      title: 'Эскроу-защита',
      description: 'Безопасные сделки с гарантией возврата',
      color: 'from-blue-500 to-cyan-500',
    },
    {
      icon: '🗺️',
      title: 'Приватность',
      description: 'Контроль видимости вашего адреса',
      color: 'from-green-500 to-emerald-500',
    },
    {
      icon: '💬',
      title: 'Живой чат',
      description: 'Общение с анимированными эмодзи',
      color: 'from-pink-500 to-rose-500',
    },
  ];

  const stats = [
    { value: '2.5M+', label: 'Активных пользователей' },
    { value: '500K+', label: 'Объявлений' },
    { value: '50K+', label: 'Ежедневных сделок' },
    { value: '99.9%', label: 'Безопасность' },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-100 to-base-200">
      {/* Header */}
      <header className="navbar bg-base-100/80 backdrop-blur-md sticky top-0 z-50 shadow-sm">
        <div className="navbar-start">
          <SveTuLogoStatic variant="gradient" width={120} height={40} />
        </div>
        <nav className="navbar-center hidden lg:flex">
          <ul className="menu menu-horizontal px-1">
            <li>
              <a href="#categories">Категории</a>
            </li>
            <li>
              <a href="#how-it-works">Как это работает</a>
            </li>
            <li>
              <a href="/b2c">Магазины</a>
            </li>
            <li>
              <a href="#about">О нас</a>
            </li>
          </ul>
        </nav>
        <div className="navbar-end">
          <button className="btn btn-ghost btn-circle lg:hidden">
            <svg
              className="w-6 h-6"
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
          </button>
          <Link
            href="/auth/login"
            className="btn btn-ghost btn-sm hidden lg:flex"
          >
            Войти
          </Link>
          <Link href="/auth/register" className="btn btn-primary btn-sm ml-2">
            Регистрация
          </Link>
        </div>
      </header>

      {/* Hero Section */}
      <section className="hero min-h-[40vh] relative overflow-hidden">
        <div className="hero-content text-center py-8">
          <AnimatedSection animation="fadeIn">
            <div className="max-w-4xl">
              <h1 className="text-4xl lg:text-5xl font-bold mb-4 bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
                Sve Tu - Всё здесь!
              </h1>
              <p className="text-lg lg:text-xl mb-6 text-base-content/80">
                Покупайте и продавайте с комфортом на главной платформе Сербии
              </p>

              {/* Search Bar */}
              <AnimatedSection animation="slideUp" delay={0.2}>
                <div className="card bg-base-100 shadow-2xl p-4 mb-6">
                  <div className="flex flex-col lg:flex-row gap-4">
                    <div className="form-control flex-1">
                      <div className="input-group">
                        <span className="bg-base-200">
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
                        </span>
                        <input
                          type="text"
                          placeholder="Что вы ищете?"
                          className="input input-bordered w-full"
                          value={searchQuery}
                          onChange={(e) => setSearchQuery(e.target.value)}
                        />
                      </div>
                    </div>
                    <select
                      className="select select-bordered"
                      value={selectedCategory}
                      onChange={(e) => setSelectedCategory(e.target.value)}
                    >
                      {categories.map((cat) => (
                        <option key={cat.id} value={cat.id}>
                          {cat.icon} {cat.name}
                        </option>
                      ))}
                    </select>
                    <button className="btn btn-primary">
                      <svg
                        className="w-5 h-5 mr-2"
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
                      Найти
                    </button>
                  </div>
                </div>
              </AnimatedSection>

              {/* Quick Categories */}
              <AnimatedSection animation="fadeIn" delay={0.3}>
                <div className="flex flex-wrap justify-center gap-3">
                  {categories.slice(1, 7).map((cat) => (
                    <button
                      key={cat.id}
                      className="btn btn-outline btn-sm gap-2 hover:scale-105 transition-transform"
                    >
                      <span className="text-lg">{cat.icon}</span>
                      {cat.name}
                    </button>
                  ))}
                </div>
              </AnimatedSection>
            </div>
          </AnimatedSection>
        </div>

        {/* Animated Background Elements */}
        <div className="absolute inset-0 -z-10 overflow-hidden">
          <div className="absolute -top-40 -right-40 w-80 h-80 bg-primary/20 rounded-full blur-3xl animate-pulse"></div>
          <div className="absolute -bottom-40 -left-40 w-80 h-80 bg-secondary/20 rounded-full blur-3xl animate-pulse delay-1000"></div>
        </div>
      </section>

      {/* Popular Products */}
      <section className="py-6 px-4">
        <div className="container mx-auto max-w-7xl">
          <AnimatedSection animation="fadeIn">
            <div className="text-center mb-4">
              <h2 className="text-2xl lg:text-3xl font-bold mb-2">
                🔥 Горячие предложения
              </h2>
              <p className="text-base text-base-content/70">
                Самые популярные товары прямо сейчас
              </p>
            </div>
          </AnimatedSection>

          <div className="grid grid-cols-2 md:grid-cols-4 lg:grid-cols-8 gap-3">
            {popularProducts.slice(0, 8).map((product, idx) => (
              <AnimatedSection
                key={product.id}
                animation="slideUp"
                delay={idx * 0.05}
              >
                <div className="card bg-base-100 shadow hover:shadow-lg transition-all hover:-translate-y-1 card-compact">
                  {product.isNew && (
                    <div className="badge badge-secondary badge-xs absolute top-1 left-1 z-10">
                      New
                    </div>
                  )}
                  {product.discount && (
                    <div className="badge badge-error badge-xs absolute top-1 right-1 z-10">
                      -{product.discount}%
                    </div>
                  )}
                  <figure className="relative h-32 overflow-hidden">
                    <img
                      src={product.image}
                      alt={product.title}
                      className="w-full h-full object-cover hover:scale-110 transition-transform duration-300"
                    />
                  </figure>
                  <div className="card-body p-2">
                    <h3 className="text-xs font-semibold line-clamp-2">
                      {product.title}
                    </h3>
                    <div className="text-xs text-base-content/60">
                      {product.location}
                    </div>
                    <div className="text-sm font-bold text-primary">
                      €{product.price}
                    </div>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* Categories Grid */}
      <section className="py-6 bg-base-200">
        <div className="container mx-auto max-w-7xl px-4">
          <AnimatedSection animation="fadeIn">
            <h2 className="text-2xl lg:text-3xl font-bold mb-4 text-center">
              📋 Популярные категории
            </h2>
          </AnimatedSection>
          <div className="grid grid-cols-4 md:grid-cols-8 gap-3">
            {[
              { icon: '🏠', name: 'Недвижимость', count: '12K+' },
              { icon: '🚗', name: 'Авто', count: '8.5K+' },
              { icon: '💻', name: 'Электроника', count: '15K+' },
              { icon: '👕', name: 'Одежда', count: '22K+' },
              { icon: '🎪', name: 'Работа', count: '5K+' },
              { icon: '🎨', name: 'Хобби', count: '7K+' },
              { icon: '🐶', name: 'Животные', count: '3K+' },
              { icon: '🏁', name: 'Спорт', count: '6K+' },
            ].map((cat, idx) => (
              <AnimatedSection key={idx} animation="zoomIn" delay={idx * 0.05}>
                <div className="card bg-base-100 shadow hover:shadow-lg transition-all cursor-pointer hover:-translate-y-1">
                  <div className="card-body p-3 text-center">
                    <div className="text-2xl mb-1">{cat.icon}</div>
                    <h3 className="text-xs font-semibold">{cat.name}</h3>
                    <p className="text-xs text-base-content/60">{cat.count}</p>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* Hot Deals Section */}
      <section className="py-6 bg-gradient-to-r from-red-50 to-orange-50">
        <div className="container mx-auto max-w-7xl px-4">
          <AnimatedSection animation="fadeIn">
            <div className="text-center mb-4">
              <h2 className="text-2xl lg:text-3xl font-bold mb-2">
                🏷️ Товары со скидкой
              </h2>
              <p className="text-base text-base-content/70">
                Успейте купить по выгодной цене!
              </p>
            </div>
          </AnimatedSection>

          <div className="grid grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-3">
            {[
              {
                title: 'Samsung TV 55"',
                price: 599,
                oldPrice: 899,
                discount: 33,
                image: configManager.buildImageUrl(
                  '/listings/20/1753428897128302370.jpg'
                ),
              },
              {
                title: 'Nike Air Max',
                price: 89,
                oldPrice: 149,
                discount: 40,
                image: configManager.buildImageUrl(
                  '/listings/21/1753445822265644326.jpg'
                ),
              },
              {
                title: 'PlayStation 5',
                price: 449,
                oldPrice: 549,
                discount: 18,
                image: configManager.buildImageUrl(
                  '/listings/23/1753548160663849380.jpg'
                ),
              },
              {
                title: 'iPad Air',
                price: 549,
                oldPrice: 699,
                discount: 21,
                image: configManager.buildImageUrl(
                  '/listings/24/1753549239639443133.jpg'
                ),
              },
              {
                title: 'Canon EOS R5',
                price: 2999,
                oldPrice: 3999,
                discount: 25,
                image: configManager.buildImageUrl(
                  '/listings/25/1753550885742188000.jpg'
                ),
              },
              {
                title: 'Dyson V15',
                price: 399,
                oldPrice: 599,
                discount: 33,
                image: configManager.buildImageUrl(
                  '/listings/26/1753554432788980038.jpg'
                ),
              },
            ].map((product, idx) => (
              <AnimatedSection key={idx} animation="zoomIn" delay={idx * 0.05}>
                <div className="card bg-base-100 shadow-lg hover:shadow-2xl transition-all hover:-translate-y-1">
                  <figure className="relative h-24">
                    <img
                      src={product.image}
                      alt={product.title}
                      className="w-full h-full object-cover"
                    />
                    <div className="absolute top-1 right-1 badge badge-error badge-sm">
                      -{product.discount}%
                    </div>
                  </figure>
                  <div className="card-body p-2">
                    <h3 className="text-xs font-semibold line-clamp-2">
                      {product.title}
                    </h3>
                    <div>
                      <div className="text-sm font-bold text-error">
                        €{product.price}
                      </div>
                      <div className="text-xs text-base-content/50 line-through">
                        €{product.oldPrice}
                      </div>
                    </div>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>

          <AnimatedSection animation="fadeIn" delay={0.3}>
            <div className="text-center mt-4">
              <div className="countdown font-mono text-xl inline-flex items-center gap-1">
                <span style={{ '--value': 23 } as any}></span>:
                <span style={{ '--value': 45 } as any}></span>:
                <span style={{ '--value': 12 } as any}></span>
                <span className="text-sm text-base-content/60 ml-2">
                  до конца акции
                </span>
              </div>
            </div>
          </AnimatedSection>
        </div>
      </section>

      {/* Black Friday Storefronts */}
      <section className="py-8 bg-black text-white">
        <div className="container mx-auto max-w-7xl px-4">
          <AnimatedSection animation="fadeIn">
            <div className="text-center mb-6">
              <h2 className="text-3xl lg:text-4xl font-bold mb-4">
                ⚡ Черная пятница в витринах
              </h2>
              <p className="text-lg opacity-80">
                Официальные магазины с мега-скидками
              </p>
            </div>
          </AnimatedSection>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {[
              {
                name: 'TechnoWorld',
                category: 'Электроника',
                discount: 'до 70%',
                followers: '125K',
                logo: '🖥️',
                color: 'from-blue-600 to-purple-600',
              },
              {
                name: 'FashionHub',
                category: 'Одежда и обувь',
                discount: 'до 80%',
                followers: '89K',
                logo: '👗',
                color: 'from-pink-600 to-red-600',
              },
              {
                name: 'HomeDecor',
                category: 'Дом и сад',
                discount: 'до 60%',
                followers: '67K',
                logo: '🏡',
                color: 'from-green-600 to-teal-600',
              },
              {
                name: 'SportMaster',
                category: 'Спорт и отдых',
                discount: 'до 50%',
                followers: '234K',
                logo: '⚽',
                color: 'from-orange-600 to-red-600',
              },
            ].map((store, idx) => (
              <AnimatedSection key={idx} animation="slideUp" delay={idx * 0.1}>
                <div className="card bg-gray-900 border border-gray-800 hover:border-yellow-400 transition-all">
                  <div className="card-body">
                    <div className="flex items-start justify-between mb-4">
                      <div
                        className={`w-16 h-16 rounded-xl bg-gradient-to-r ${store.color} flex items-center justify-center text-3xl`}
                      >
                        {store.logo}
                      </div>
                      <div className="badge badge-warning badge-lg">
                        BLACK FRIDAY
                      </div>
                    </div>
                    <h3 className="text-xl font-bold">{store.name}</h3>
                    <p className="text-sm opacity-70">{store.category}</p>
                    <div className="text-3xl font-bold text-yellow-400 my-2">
                      {store.discount}
                    </div>
                    <div className="flex items-center justify-between mt-auto">
                      <span className="text-sm opacity-60">
                        {store.followers} подписчиков
                      </span>
                      <button className="btn btn-sm btn-warning">
                        В магазин
                      </button>
                    </div>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* Trending Searches */}
      <section className="py-6 bg-gradient-to-r from-purple-50 to-pink-50">
        <div className="container mx-auto max-w-7xl px-4">
          <AnimatedSection animation="fadeIn">
            <h2 className="text-xl lg:text-2xl font-bold mb-3 text-center">
              🎆 Что сейчас ищут
            </h2>
          </AnimatedSection>
          <div className="flex flex-wrap justify-center gap-2">
            {[
              'iPhone 15 Pro',
              'PS5',
              'Квартира в центре',
              'MacBook Air',
              'Електросамокат',
              'Диван',
              'AirPods',
              'Кроссовки Nike',
              'Холодильник',
              'Велосипед',
            ].map((search, idx) => (
              <AnimatedSection key={idx} animation="fadeIn" delay={idx * 0.03}>
                <button className="btn btn-sm btn-outline hover:btn-primary">
                  {search}
                </button>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* Map Section */}
      <section className="py-6 px-4">
        <div className="container mx-auto max-w-7xl">
          <AnimatedSection animation="fadeIn">
            <div className="text-center mb-6">
              <h2 className="text-3xl lg:text-4xl font-bold mb-4">
                🗺️ Товары рядом с вами
              </h2>
              <p className="text-lg text-base-content/70">
                Найдите то, что нужно, в вашем районе
              </p>
            </div>
          </AnimatedSection>

          <AnimatedSection animation="zoomIn">
            <div className="card bg-base-100 shadow-2xl overflow-hidden">
              <div className="relative h-64 bg-gradient-to-br from-blue-50 to-green-50">
                {/* Map Background */}
                <div className="absolute inset-0 opacity-20">
                  {[...Array(10)].map((_, i) => (
                    <div
                      key={i}
                      className="absolute w-full border-t border-gray-300"
                      style={{ top: `${i * 10}%` }}
                    ></div>
                  ))}
                  {[...Array(10)].map((_, i) => (
                    <div
                      key={i}
                      className="absolute h-full border-l border-gray-300"
                      style={{ left: `${i * 10}%` }}
                    ></div>
                  ))}
                </div>

                {/* Map Markers */}
                <div className="absolute top-1/4 left-1/4 transform -translate-x-1/2 -translate-y-1/2">
                  <div className="relative">
                    <div className="absolute -inset-4 bg-blue-500 rounded-full opacity-30 animate-ping"></div>
                    <div className="relative bg-blue-500 text-white rounded-full w-12 h-12 flex items-center justify-center shadow-lg">
                      <svg
                        className="w-6 h-6"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm0-2a6 6 0 100-12 6 6 0 000 12z"
                          clipRule="evenodd"
                        />
                      </svg>
                    </div>
                  </div>
                  <div className="absolute top-14 left-1/2 transform -translate-x-1/2 whitespace-nowrap">
                    <div className="bg-blue-500 text-white rounded-lg px-3 py-1 text-sm shadow-lg">
                      Вы здесь
                    </div>
                  </div>
                </div>

                {/* Product Markers */}
                {[
                  { left: '40%', top: '30%', price: '€899', category: '💻' },
                  { left: '60%', top: '40%', price: '€650', category: '🏠' },
                  { left: '30%', top: '60%', price: '€249', category: '👕' },
                  { left: '70%', top: '20%', price: '€1299', category: '🚗' },
                  { left: '50%', top: '70%', price: '€89', category: '🎮' },
                ].map((marker, idx) => (
                  <div
                    key={idx}
                    className="absolute cursor-pointer transform hover:scale-110 transition-transform"
                    style={{ left: marker.left, top: marker.top }}
                  >
                    <div className="relative">
                      <div className="bg-white rounded-full w-10 h-10 flex items-center justify-center shadow-lg border-2 border-primary">
                        <span className="text-xl">{marker.category}</span>
                      </div>
                      <div className="absolute -top-8 left-1/2 transform -translate-x-1/2">
                        <div className="bg-primary text-white rounded px-2 py-1 text-xs font-bold whitespace-nowrap">
                          {marker.price}
                        </div>
                      </div>
                    </div>
                  </div>
                ))}

                {/* Map Controls */}
                <div className="absolute bottom-4 right-4 space-y-2">
                  <button className="btn btn-circle btn-sm bg-white shadow-lg">
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
                  <button className="btn btn-circle btn-sm bg-white shadow-lg">
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

                {/* Info Overlay */}
                <div className="absolute top-4 left-4 card bg-white/90 backdrop-blur-sm shadow-lg">
                  <div className="card-body p-4">
                    <h3 className="font-bold">В радиусе 5 км:</h3>
                    <div className="stats stats-vertical shadow-none bg-transparent">
                      <div className="stat p-2">
                        <div className="stat-value text-2xl">234</div>
                        <div className="stat-desc">объявления</div>
                      </div>
                      <div className="stat p-2">
                        <div className="stat-value text-2xl">12</div>
                        <div className="stat-desc">магазинов</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
              <div className="card-body">
                <div className="flex flex-col sm:flex-row gap-4 items-center justify-between">
                  <div>
                    <h3 className="text-xl font-bold">
                      Интерактивная карта товаров
                    </h3>
                    <p className="text-sm text-base-content/60">
                      Настройте фильтры и найдите нужное рядом
                    </p>
                  </div>
                  <button className="btn btn-primary">
                    Открыть карту
                    <svg
                      className="w-5 h-5 ml-2"
                      fill="none"
                      stroke="currentColor"
                      viewBox="0 0 24 24"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"
                      />
                    </svg>
                  </button>
                </div>
              </div>
            </div>
          </AnimatedSection>
        </div>
      </section>

      {/* Features */}
      <section className="py-8 bg-base-200">
        <div className="container mx-auto max-w-7xl px-4">
          <AnimatedSection animation="fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-3xl lg:text-4xl font-bold mb-4">
                Почему выбирают нас?
              </h2>
              <p className="text-lg text-base-content/70">
                Инновационные функции для вашего удобства
              </p>
            </div>
          </AnimatedSection>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {features.map((feature, idx) => (
              <AnimatedSection key={idx} animation="zoomIn" delay={idx * 0.1}>
                <div className="card bg-base-100 shadow-xl hover:shadow-2xl transition-all group">
                  <div className="card-body items-center text-center">
                    <div
                      className={`w-16 h-16 rounded-2xl bg-gradient-to-r ${feature.color} flex items-center justify-center text-3xl mb-3 group-hover:scale-110 transition-transform`}
                    >
                      {feature.icon}
                    </div>
                    <h3 className="card-title">{feature.title}</h3>
                    <p className="text-sm text-base-content/70">
                      {feature.description}
                    </p>
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* How it Works */}
      <section id="how-it-works" className="py-8 px-4 bg-base-200">
        <div className="container mx-auto max-w-7xl">
          <AnimatedSection animation="fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-3xl lg:text-4xl font-bold mb-4">
                Как это работает?
              </h2>
              <p className="text-lg text-base-content/70">
                4 сценария использования платформы
              </p>
            </div>
          </AnimatedSection>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {/* C2C Seller */}
            <AnimatedSection animation="slideUp">
              <div className="card bg-gradient-to-br from-primary/10 to-primary/5 border border-primary/20 h-full">
                <div className="card-body">
                  <div className="text-center mb-4">
                    <div className="text-4xl mb-2">👤➡️👤</div>
                    <h3 className="text-xl font-bold">Продавец C2C</h3>
                    <p className="text-sm text-base-content/60">
                      Частные объявления
                    </p>
                  </div>
                  <ul className="space-y-3">
                    <li className="flex items-start gap-2">
                      <span className="text-primary">📸</span>
                      <div>
                        <p className="font-semibold text-sm">Сделайте фото</p>
                        <p className="text-xs text-base-content/60">
                          AI создаст описание
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary">💰</span>
                      <div>
                        <p className="font-semibold text-sm">Назначьте цену</p>
                        <p className="text-xs text-base-content/60">
                          С учетом рынка
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary">🚀</span>
                      <div>
                        <p className="font-semibold text-sm">Публикация</p>
                        <p className="text-xs text-base-content/60">
                          Мгновенно онлайн
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-primary">✅</span>
                      <div>
                        <p className="font-semibold text-sm">Получите деньги</p>
                        <p className="text-xs text-base-content/60">
                          Через эскроу
                        </p>
                      </div>
                    </li>
                  </ul>
                </div>
              </div>
            </AnimatedSection>

            {/* B2C Seller */}
            <AnimatedSection animation="slideUp" delay={0.1}>
              <div className="card bg-gradient-to-br from-secondary/10 to-secondary/5 border border-secondary/20 h-full">
                <div className="card-body">
                  <div className="text-center mb-4">
                    <div className="text-4xl mb-2">🏪➡️👤</div>
                    <h3 className="text-xl font-bold">Продавец B2C</h3>
                    <p className="text-sm text-base-content/60">
                      Витрина магазина
                    </p>
                  </div>
                  <ul className="space-y-3">
                    <li className="flex items-start gap-2">
                      <span className="text-secondary">🏪</span>
                      <div>
                        <p className="font-semibold text-sm">
                          Создайте витрину
                        </p>
                        <p className="text-xs text-base-content/60">
                          Брендированный магазин
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-secondary">📦</span>
                      <div>
                        <p className="font-semibold text-sm">
                          Загрузите каталог
                        </p>
                        <p className="text-xs text-base-content/60">
                          Массовая загрузка
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-secondary">📊</span>
                      <div>
                        <p className="font-semibold text-sm">Аналитика</p>
                        <p className="text-xs text-base-content/60">
                          Dashboard с метриками
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-secondary">💳</span>
                      <div>
                        <p className="font-semibold text-sm">Автоплатежи</p>
                        <p className="text-xs text-base-content/60">
                          Интеграция с банками
                        </p>
                      </div>
                    </li>
                  </ul>
                </div>
              </div>
            </AnimatedSection>

            {/* C2C Buyer */}
            <AnimatedSection animation="slideUp" delay={0.2}>
              <div className="card bg-gradient-to-br from-accent/10 to-accent/5 border border-accent/20 h-full">
                <div className="card-body">
                  <div className="text-center mb-4">
                    <div className="text-4xl mb-2">👤⬅️👤</div>
                    <h3 className="text-xl font-bold">Покупатель C2C</h3>
                    <p className="text-sm text-base-content/60">
                      Покупки у частных лиц
                    </p>
                  </div>
                  <ul className="space-y-3">
                    <li className="flex items-start gap-2">
                      <span className="text-accent">🔍</span>
                      <div>
                        <p className="font-semibold text-sm">Поиск товара</p>
                        <p className="text-xs text-base-content/60">
                          AI-рекомендации
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-accent">💬</span>
                      <div>
                        <p className="font-semibold text-sm">Чат с продавцом</p>
                        <p className="text-xs text-base-content/60">
                          Живые эмодзи
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-accent">🤝</span>
                      <div>
                        <p className="font-semibold text-sm">Договор встречи</p>
                        <p className="text-xs text-base-content/60">
                          Или доставка
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-accent">🔒</span>
                      <div>
                        <p className="font-semibold text-sm">
                          Безопасная сделка
                        </p>
                        <p className="text-xs text-base-content/60">
                          Эскроу-защита
                        </p>
                      </div>
                    </li>
                  </ul>
                </div>
              </div>
            </AnimatedSection>

            {/* B2C Buyer */}
            <AnimatedSection animation="slideUp" delay={0.3}>
              <div className="card bg-gradient-to-br from-info/10 to-info/5 border border-info/20 h-full">
                <div className="card-body">
                  <div className="text-center mb-4">
                    <div className="text-4xl mb-2">👤⬅️🏪</div>
                    <h3 className="text-xl font-bold">Покупатель B2C</h3>
                    <p className="text-sm text-base-content/60">
                      Покупки в магазинах
                    </p>
                  </div>
                  <ul className="space-y-3">
                    <li className="flex items-start gap-2">
                      <span className="text-info">🛍️</span>
                      <div>
                        <p className="font-semibold text-sm">Витрины брендов</p>
                        <p className="text-xs text-base-content/60">
                          Официальные магазины
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-info">🎯</span>
                      <div>
                        <p className="font-semibold text-sm">Акции и скидки</p>
                        <p className="text-xs text-base-content/60">
                          Черная пятница
                        </p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-info">🚚</span>
                      <div>
                        <p className="font-semibold text-sm">
                          Быстрая доставка
                        </p>
                        <p className="text-xs text-base-content/60">От 1 дня</p>
                      </div>
                    </li>
                    <li className="flex items-start gap-2">
                      <span className="text-info">🛡️</span>
                      <div>
                        <p className="font-semibold text-sm">Гарантия</p>
                        <p className="text-xs text-base-content/60">
                          Возврат и обмен
                        </p>
                      </div>
                    </li>
                  </ul>
                </div>
              </div>
            </AnimatedSection>
          </div>
        </div>
      </section>

      {/* Stats */}
      <section className="py-8 bg-gradient-to-r from-primary to-secondary text-primary-content">
        <div className="container mx-auto max-w-6xl px-4">
          <AnimatedSection animation="fadeIn">
            <div className="text-center mb-8">
              <h2 className="text-3xl lg:text-4xl font-bold mb-4">
                Наши достижения
              </h2>
              <p className="text-lg opacity-90">Цифры говорят сами за себя</p>
            </div>
          </AnimatedSection>

          <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
            {stats.map((stat, idx) => (
              <AnimatedSection key={idx} animation="slideUp" delay={idx * 0.1}>
                <div className="text-center">
                  <div className="text-4xl lg:text-5xl font-bold mb-2">
                    {stat.value}
                  </div>
                  <div className="text-sm lg:text-base opacity-90">
                    {stat.label}
                  </div>
                </div>
              </AnimatedSection>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-10 px-4">
        <div className="container mx-auto max-w-4xl">
          <AnimatedSection animation="zoomIn">
            <div className="card bg-gradient-to-r from-primary to-secondary text-primary-content shadow-2xl">
              <div className="card-body text-center p-8">
                <h2 className="text-3xl lg:text-4xl font-bold mb-6">
                  Готовы начать?
                </h2>
                <p className="text-xl mb-8 opacity-90">
                  Присоединяйтесь к миллионам пользователей уже сегодня
                </p>
                <div className="flex flex-col sm:flex-row gap-4 justify-center">
                  <button className="btn btn-lg bg-white text-primary hover:bg-gray-100">
                    <svg
                      className="w-6 h-6 mr-2"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path d="M10 2a5 5 0 00-5 5v2a2 2 0 00-2 2v5a2 2 0 002 2h10a2 2 0 002-2v-5a2 2 0 00-2-2H7V7a3 3 0 015.905-.75 1 1 0 001.937-.5A5.002 5.002 0 0010 2z" />
                    </svg>
                    Создать аккаунт
                  </button>
                  <button className="btn btn-lg btn-outline border-white text-white hover:bg-white hover:text-primary">
                    <svg
                      className="w-6 h-6 mr-2"
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
                    Разместить объявление
                  </button>
                </div>
              </div>
            </div>
          </AnimatedSection>
        </div>
      </section>

      {/* Download App */}
      <section className="py-8 bg-base-200">
        <div className="container mx-auto max-w-6xl px-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
            <AnimatedSection animation="slideLeft">
              <div>
                <h2 className="text-3xl lg:text-4xl font-bold mb-6">
                  Sve Tu всегда с вами
                </h2>
                <p className="text-lg text-base-content/70 mb-8">
                  Скачайте мобильное приложение и получите доступ к миллионам
                  товаров в вашем кармане
                </p>
                <div className="flex flex-wrap gap-4">
                  <button className="btn btn-neutral gap-2">
                    <svg
                      className="w-6 h-6"
                      viewBox="0 0 24 24"
                      fill="currentColor"
                    >
                      <path d="M17.05 20.28c-.98.95-2.05.8-3.08.35-1.09-.46-2.09-.48-3.24 0-1.44.62-2.2.44-3.06-.35C2.79 15.25 3.51 7.59 9.05 7.31c1.35.07 2.29.74 3.08.8 1.18-.24 2.31-.93 3.57-.84 1.51.12 2.65.72 3.4 1.8-3.12 1.87-2.38 5.98.48 7.13-.57 1.5-1.31 2.99-2.54 4.09l.01-.01zM12.03 7.25c-.15-2.23 1.66-4.07 3.74-4.25.29 2.58-2.34 4.5-3.74 4.25z" />
                    </svg>
                    App Store
                  </button>
                  <button className="btn btn-neutral gap-2">
                    <svg
                      className="w-6 h-6"
                      viewBox="0 0 24 24"
                      fill="currentColor"
                    >
                      <path d="M3,20.5V3.5C3,2.91 3.34,2.39 3.84,2.15L13.69,12L3.84,21.85C3.34,21.6 3,21.09 3,20.5M16.81,15.12L6.05,21.34L14.54,12.85L16.81,15.12M20.16,10.81C20.5,11.08 20.75,11.5 20.75,12C20.75,12.5 20.53,12.9 20.18,13.18L17.89,14.5L15.39,12L17.89,9.5L20.16,10.81M6.05,2.66L16.81,8.88L14.54,11.15L6.05,2.66Z" />
                    </svg>
                    Google Play
                  </button>
                </div>
                <div className="mt-8 flex items-center gap-6">
                  <div className="flex -space-x-2">
                    {[1, 2, 3, 4, 5].map((i) => (
                      <div key={i} className="avatar">
                        <div className="w-10 rounded-full ring ring-base-100">
                          <img
                            src={`https://ui-avatars.com/api/?name=User${i}&background=random`}
                            alt="User"
                          />
                        </div>
                      </div>
                    ))}
                  </div>
                  <div>
                    <div className="flex items-center gap-1">
                      <div className="rating rating-sm">
                        {[1, 2, 3, 4, 5].map((star) => (
                          <input
                            key={star}
                            type="radio"
                            className="mask mask-star-2 bg-orange-400"
                            checked={star <= 5}
                            readOnly
                          />
                        ))}
                      </div>
                      <span className="font-bold">4.9</span>
                    </div>
                    <p className="text-sm text-base-content/60">
                      из 50K+ отзывов
                    </p>
                  </div>
                </div>
              </div>
            </AnimatedSection>

            <AnimatedSection animation="slideRight">
              <div className="mockup-phone">
                <div className="camera"></div>
                <div className="display">
                  <div className="artboard artboard-demo phone-1 bg-base-100">
                    <div className="w-full h-full flex items-center justify-center">
                      <img
                        src={configManager.buildImageUrl(
                          '/listings/7/1753007242863504454.jpg'
                        )}
                        alt="App Screenshot"
                        className="w-full h-full object-cover"
                      />
                    </div>
                  </div>
                </div>
              </div>
            </AnimatedSection>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="footer p-10 bg-neutral text-neutral-content">
        <div className="container mx-auto max-w-7xl">
          <div className="footer">
            <div>
              <SveTuLogoStatic variant="minimal" width={120} height={40} />
              <p className="mt-4">
                Sve Tu Platform
                <br />
                Всё что нужно в одном месте
              </p>
              <div className="flex gap-4 mt-4">
                <a className="btn btn-circle btn-sm">
                  <svg
                    className="w-5 h-5"
                    fill="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z" />
                  </svg>
                </a>
                <a className="btn btn-circle btn-sm">
                  <svg
                    className="w-5 h-5"
                    fill="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path d="M23.953 4.57a10 10 0 01-2.825.775 4.958 4.958 0 002.163-2.723c-.951.555-2.005.959-3.127 1.184a4.92 4.92 0 00-8.384 4.482C7.69 8.095 4.067 6.13 1.64 3.162a4.822 4.822 0 00-.666 2.475c0 1.71.87 3.213 2.188 4.096a4.904 4.904 0 01-2.228-.616v.06a4.923 4.923 0 003.946 4.827 4.996 4.996 0 01-2.212.085 4.936 4.936 0 004.604 3.417 9.867 9.867 0 01-6.102 2.105c-.39 0-.779-.023-1.17-.067a13.995 13.995 0 007.557 2.209c9.053 0 13.998-7.496 13.998-13.985 0-.21 0-.42-.015-.63A9.935 9.935 0 0024 4.59z" />
                  </svg>
                </a>
                <a className="btn btn-circle btn-sm">
                  <svg
                    className="w-5 h-5"
                    fill="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zM5.838 12a6.162 6.162 0 1112.324 0 6.162 6.162 0 01-12.324 0zM12 16a4 4 0 110-8 4 4 0 010 8zm4.965-10.405a1.44 1.44 0 112.881.001 1.44 1.44 0 01-2.881-.001z" />
                  </svg>
                </a>
              </div>
            </div>
            <div>
              <span className="footer-title">Покупателям</span>
              <a className="link link-hover">Как купить</a>
              <a className="link link-hover">Безопасность</a>
              <a className="link link-hover">Доставка</a>
              <a className="link link-hover">Гарантии</a>
            </div>
            <div>
              <span className="footer-title">Продавцам</span>
              <a className="link link-hover">Как продать</a>
              <a className="link link-hover">Тарифы</a>
              <a className="link link-hover">Продвижение</a>
              <a className="link link-hover">Аналитика</a>
            </div>
            <div>
              <span className="footer-title">Помощь</span>
              <a className="link link-hover">Частые вопросы</a>
              <a className="link link-hover">Служба поддержки</a>
              <a className="link link-hover">Правила</a>
              <a className="link link-hover">Блог</a>
            </div>
          </div>
          <div className="divider"></div>
          <div className="footer items-center">
            <div className="grid-flow-col gap-4">
              <p>© 2024 Sve Tu. Все права защищены.</p>
            </div>
            <div className="grid-flow-col gap-4 md:place-self-center md:justify-self-end">
              <a className="link link-hover">Условия использования</a>
              <a className="link link-hover">Политика конфиденциальности</a>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
};

export default IdealHomepage;
