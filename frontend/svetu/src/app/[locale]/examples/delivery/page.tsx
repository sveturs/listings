'use client';

import { useState } from 'react';
import {
  TruckIcon,
  MapPinIcon,
  CreditCardIcon,
  ShieldCheckIcon,
  ClockIcon,
  CubeIcon,
  ChartBarIcon,
  UserGroupIcon,
  BuildingStorefrontIcon,
  ArrowRightIcon,
  BanknotesIcon,
  DocumentCheckIcon,
  ArrowPathIcon,
} from '@heroicons/react/24/outline';
import { StarIcon } from '@heroicons/react/24/solid';

// Import delivery components
import DeliveryMethodSelector from './components/DeliveryMethodSelector';
import TrackingWidget from './components/TrackingWidget';
import SellerShipmentInterface from './components/SellerShipmentInterface';
import ParcelShopMap from './components/ParcelShopMap';
import DeliveryCalculator from './components/DeliveryCalculator';
import BulkShipmentManager from './components/BulkShipmentManager';

export default function DeliveryExamplesPage() {
  const [activeTab, setActiveTab] = useState('overview');
  const [selectedDeliveryMethod, setSelectedDeliveryMethod] =
    useState('courier');

  const features = [
    {
      icon: TruckIcon,
      title: 'Курьерская доставка',
      description: 'Доставка на адрес получателя курьером BEX Express',
      badge: 'Популярно',
    },
    {
      icon: MapPinIcon,
      title: 'Пункты выдачи',
      description: 'Сеть из 200+ пунктов самовывоза по всей Сербии',
      badge: 'Удобно',
    },
    {
      icon: BanknotesIcon,
      title: 'Оплата при получении',
      description: 'Безопасная оплата наличными или картой при получении',
      badge: 'COD',
    },
    {
      icon: ShieldCheckIcon,
      title: 'Страхование',
      description: 'Защита посылки на сумму до 100,000 RSD',
      badge: 'Защита',
    },
  ];

  const stats = [
    { label: 'Городов покрытия', value: '150+', icon: MapPinIcon },
    { label: 'Пунктов выдачи', value: '200+', icon: BuildingStorefrontIcon },
    { label: 'Среднее время доставки', value: '2-3 дня', icon: ClockIcon },
    { label: 'Доставлено посылок', value: '10M+', icon: CubeIcon },
  ];

  const testimonials = [
    {
      name: 'Марко Петрович',
      role: 'Продавец электроники',
      rating: 5,
      text: 'С BEX Express я могу легко отправлять товары по всей Сербии. Отслеживание работает отлично!',
      avatar: '👨‍💼',
    },
    {
      name: 'Ана Йованович',
      role: 'Покупатель',
      rating: 5,
      text: 'Очень удобно выбрать пункт выдачи рядом с работой. Забираю посылки по пути домой.',
      avatar: '👩',
    },
    {
      name: 'Стефан Николич',
      role: 'Владелец магазина',
      rating: 5,
      text: 'Массовая обработка заказов экономит кучу времени. API работает стабильно.',
      avatar: '👨‍💻',
    },
  ];

  const tabs = [
    { id: 'overview', label: 'Обзор', icon: ChartBarIcon },
    { id: 'customer', label: 'Для покупателя', icon: UserGroupIcon },
    { id: 'seller', label: 'Для продавца', icon: BuildingStorefrontIcon },
    { id: 'tracking', label: 'Отслеживание', icon: MapPinIcon },
    { id: 'calculator', label: 'Калькулятор', icon: CreditCardIcon },
  ];

  return (
    <div className="min-h-screen bg-gradient-to-b from-base-100 to-base-200">
      {/* Hero Section */}
      <div className="bg-gradient-to-r from-primary to-secondary text-primary-content">
        <div className="container mx-auto px-4 py-6 md:py-12">
          <div className="flex flex-col sm:flex-row items-center gap-3 mb-4">
            <div className="p-3 bg-white/20 rounded-xl backdrop-blur-sm">
              <TruckIcon className="w-6 h-6 sm:w-8 sm:h-8" />
            </div>
            <div className="text-center sm:text-left">
              <h1 className="text-2xl sm:text-3xl md:text-4xl font-bold">
                Интеграция доставки BEX Express
              </h1>
              <p className="text-sm sm:text-base text-primary-content/80 mt-2">
                Современное решение для доставки товаров на маркетплейсе
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Tabs Navigation */}
      <div className="bg-base-100 border-b sticky top-0 z-40 backdrop-blur-lg bg-opacity-90">
        <div className="container mx-auto px-2 sm:px-4">
          <div className="flex gap-1 sm:gap-2 overflow-x-auto py-2 sm:py-4 scrollbar-hide">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`
                  flex items-center gap-1 sm:gap-2 px-2 sm:px-4 py-2 rounded-lg transition-all
                  whitespace-nowrap text-xs sm:text-sm font-medium min-w-fit
                  ${
                    activeTab === tab.id
                      ? 'bg-primary text-primary-content shadow-lg'
                      : 'bg-base-200 hover:bg-base-300'
                  }
                `}
              >
                <tab.icon className="w-4 h-4 sm:w-5 sm:h-5 flex-shrink-0" />
                <span className="hidden xs:inline">{tab.label}</span>
              </button>
            ))}
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="container mx-auto px-4 py-4 sm:py-8">
        {activeTab === 'overview' && (
          <div className="space-y-8">
            {/* Features Grid */}
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 sm:gap-6">
              {features.map((feature, index) => (
                <div
                  key={index}
                  className="card bg-base-100 shadow-xl hover:shadow-2xl transition-all"
                >
                  <div className="card-body">
                    <div className="flex items-start justify-between">
                      <div className="p-3 bg-primary/10 rounded-lg">
                        <feature.icon className="w-6 h-6 text-primary" />
                      </div>
                      {feature.badge && (
                        <div className="badge badge-primary badge-sm">
                          {feature.badge}
                        </div>
                      )}
                    </div>
                    <h3 className="card-title text-lg mt-4">{feature.title}</h3>
                    <p className="text-base-content/70">
                      {feature.description}
                    </p>
                  </div>
                </div>
              ))}
            </div>

            {/* Stats */}
            <div className="bg-base-100 rounded-2xl shadow-xl p-4 sm:p-8">
              <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6 text-center">
                Ключевые показатели
              </h2>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
                {stats.map((stat, index) => (
                  <div key={index} className="text-center">
                    <div className="flex justify-center mb-3">
                      <div className="p-3 bg-secondary/10 rounded-full">
                        <stat.icon className="w-8 h-8 text-secondary" />
                      </div>
                    </div>
                    <div className="text-xl sm:text-3xl font-bold text-primary">
                      {stat.value}
                    </div>
                    <div className="text-xs sm:text-sm text-base-content/60 mt-1">
                      {stat.label}
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Process Steps */}
            <div className="bg-gradient-to-r from-primary/5 to-secondary/5 rounded-2xl p-4 sm:p-8">
              <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6">
                Как это работает
              </h2>
              <div className="grid grid-cols-2 md:grid-cols-4 gap-4 sm:gap-6">
                {[
                  {
                    step: '1',
                    title: 'Оформление',
                    desc: 'Покупатель выбирает способ доставки',
                  },
                  {
                    step: '2',
                    title: 'Подготовка',
                    desc: 'Продавец готовит посылку к отправке',
                  },
                  {
                    step: '3',
                    title: 'Доставка',
                    desc: 'BEX Express доставляет товар',
                  },
                  {
                    step: '4',
                    title: 'Получение',
                    desc: 'Покупатель получает и оплачивает',
                  },
                ].map((item, index) => (
                  <div key={index} className="relative">
                    {index < 3 && (
                      <ArrowRightIcon className="absolute top-8 -right-3 w-6 h-6 text-primary/30 hidden md:block" />
                    )}
                    <div className="text-center">
                      <div className="w-12 h-12 sm:w-16 sm:h-16 bg-primary text-primary-content rounded-full flex items-center justify-center text-lg sm:text-2xl font-bold mx-auto mb-2 sm:mb-4">
                        {item.step}
                      </div>
                      <h3 className="text-sm sm:text-base font-semibold mb-1 sm:mb-2">
                        {item.title}
                      </h3>
                      <p className="text-xs sm:text-sm text-base-content/60">
                        {item.desc}
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>

            {/* Testimonials */}
            <div>
              <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6">
                Отзывы пользователей
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-3 gap-4 sm:gap-6">
                {testimonials.map((testimonial, index) => (
                  <div key={index} className="card bg-base-100 shadow-xl">
                    <div className="card-body">
                      <div className="flex items-center gap-3 mb-4">
                        <div className="text-4xl">{testimonial.avatar}</div>
                        <div>
                          <div className="font-semibold">
                            {testimonial.name}
                          </div>
                          <div className="text-sm text-base-content/60">
                            {testimonial.role}
                          </div>
                        </div>
                      </div>
                      <div className="flex gap-1 mb-3">
                        {[...Array(testimonial.rating)].map((_, i) => (
                          <StarIcon key={i} className="w-5 h-5 text-warning" />
                        ))}
                      </div>
                      <p className="text-base-content/80 italic">
                        &ldquo;{testimonial.text}&rdquo;
                      </p>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        )}

        {activeTab === 'customer' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Выбор способа доставки</h2>
              <p>
                Покупатели могут выбрать удобный способ получения товара при
                оформлении заказа.
              </p>
            </div>
            <DeliveryMethodSelector
              onMethodChange={setSelectedDeliveryMethod}
              selectedMethod={selectedDeliveryMethod}
            />

            {selectedDeliveryMethod === 'parcel-shop' && (
              <div>
                <h3 className="text-xl font-semibold mb-4">
                  Карта пунктов выдачи
                </h3>
                <ParcelShopMap />
              </div>
            )}
          </div>
        )}

        {activeTab === 'seller' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Интерфейс продавца</h2>
              <p>
                Удобные инструменты для управления отправками и массовой
                обработки заказов.
              </p>
            </div>

            <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 sm:gap-8">
              <div>
                <h3 className="text-xl font-semibold mb-4">
                  Создание отправления
                </h3>
                <SellerShipmentInterface />
              </div>

              <div>
                <h3 className="text-xl font-semibold mb-4">
                  Массовая обработка
                </h3>
                <BulkShipmentManager />
              </div>
            </div>
          </div>
        )}

        {activeTab === 'tracking' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Отслеживание посылки</h2>
              <p>
                Покупатели и продавцы могут отслеживать статус доставки в
                реальном времени.
              </p>
            </div>
            <TrackingWidget />
          </div>
        )}

        {activeTab === 'calculator' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Калькулятор стоимости доставки</h2>
              <p>
                Рассчитайте стоимость доставки в зависимости от параметров
                посылки и маршрута.
              </p>
            </div>
            <DeliveryCalculator />
          </div>
        )}
      </div>

      {/* CTA Section */}
      <div className="bg-gradient-to-r from-primary to-secondary text-primary-content mt-8 sm:mt-16">
        <div className="container mx-auto px-4 py-6 sm:py-12">
          <div className="text-center">
            <h2 className="text-xl sm:text-3xl font-bold mb-2 sm:mb-4">
              Готовы подключить доставку?
            </h2>
            <p className="text-sm sm:text-xl mb-4 sm:mb-8 opacity-90">
              Интегрируйте BEX Express и расширьте географию ваших продаж
            </p>
            <div className="flex flex-col sm:flex-row gap-2 sm:gap-4 justify-center">
              <button className="btn btn-sm sm:btn-lg bg-white text-primary hover:bg-white/90">
                <DocumentCheckIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                Документация API
              </button>
              <button className="btn btn-sm sm:btn-lg btn-outline border-white text-white hover:bg-white/20">
                <ArrowPathIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                Начать интеграцию
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
