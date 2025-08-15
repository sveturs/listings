'use client';

import { useState } from 'react';
import {
  TruckIcon,
  MapPinIcon,
  ClockIcon,
  QrCodeIcon,
  ChartBarIcon,
  UserGroupIcon,
  BuildingStorefrontIcon,
  MagnifyingGlassIcon,
  CalculatorIcon,
} from '@heroicons/react/24/outline';

// Import Post Express components
import {
  PostExpressDeliveryFlow,
  PostExpressDeliverySelector,
  PostExpressRateCalculator,
  PostExpressTracker,
  PostExpressPickupCode,
} from '@/components/delivery/postexpress';

export default function PostExpressExamplesPage() {
  const [activeTab, setActiveTab] = useState('overview');
  const [deliveryData, setDeliveryData] = useState<any>(null);

  // Mock data for demonstrations
  const mockPickupOrder = {
    id: 1,
    pickup_code: 'PE-NS-240815-001',
    status: 'ready',
    created_at: '2024-08-15T10:00:00Z',
    expires_at: '2024-08-22T18:00:00Z',
    customer_name: 'Петар Петрович',
    customer_phone: '+381 60 123 4567',
    items_count: 3,
    total_amount: 4500,
    warehouse: {
      code: 'NS-MAIN-01',
      name: 'Склад Sve Tu - Нови Сад',
      address: 'Микија Манојловића 53',
      phone: '+381 21 123 456',
      working_hours: {
        monday: '09:00-19:00',
        tuesday: '09:00-19:00',
        wednesday: '09:00-19:00',
        thursday: '09:00-19:00',
        friday: '09:00-19:00',
        saturday: '10:00-16:00',
        sunday: 'Закрыто',
      },
    },
    notes: 'Заказ готов к выдаче. Товары проверены и упакованы.',
  };

  const features = [
    {
      icon: TruckIcon,
      title: 'Курьерская доставка',
      description: 'Доставка на адрес получателя курьером Post Express',
      badge: 'Популярно',
    },
    {
      icon: MapPinIcon,
      title: 'Отделения Post Express',
      description: 'Сеть из 180+ отделений по всей Сербии',
      badge: 'Удобно',
    },
    {
      icon: BuildingStorefrontIcon,
      title: 'Склад Sve Tu',
      description: 'Бесплатный самовывоз со склада в Нови Саде',
      badge: 'Бесплатно',
    },
    {
      icon: QrCodeIcon,
      title: 'QR коды самовывоза',
      description: 'Удобные коды для получения товаров на складе',
      badge: 'Технологично',
    },
  ];

  const stats = [
    { label: 'Городов покрытия', value: '180+', icon: MapPinIcon },
    { label: 'Отделений', value: '200+', icon: BuildingStorefrontIcon },
    { label: 'Среднее время доставки', value: '1-2 дня', icon: ClockIcon },
    { label: 'Точность доставки', value: '99.5%', icon: ChartBarIcon },
  ];

  const tabs = [
    { id: 'overview', label: 'Обзор', icon: ChartBarIcon },
    { id: 'flow', label: 'Процесс доставки', icon: TruckIcon },
    { id: 'selector', label: 'Выбор доставки', icon: UserGroupIcon },
    { id: 'calculator', label: 'Калькулятор', icon: CalculatorIcon },
    { id: 'tracking', label: 'Отслеживание', icon: MagnifyingGlassIcon },
    { id: 'pickup', label: 'Код самовывоза', icon: QrCodeIcon },
  ];

  const handleDeliveryComplete = (data: any) => {
    setDeliveryData(data);
    console.log('Delivery data:', data);
  };

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
                Интеграция Post Express
              </h1>
              <p className="text-sm sm:text-base text-primary-content/80 mt-2">
                Национальный почтовый оператор Сербии для маркетплейса Sve Tu
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
                Ключевые показатели Post Express
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

            {/* Advantages */}
            <div className="bg-gradient-to-r from-primary/5 to-secondary/5 rounded-2xl p-4 sm:p-8">
              <h2 className="text-xl sm:text-2xl font-bold mb-4 sm:mb-6">
                Преимущества Post Express
              </h2>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 sm:gap-6">
                <div className="space-y-2">
                  <h3 className="font-semibold">🚚 Курьерская доставка</h3>
                  <p className="text-sm text-base-content/70">
                    Услуга &quot;Данас за сутра&quot; - доставка до 19:00
                    следующего дня
                  </p>
                </div>
                <div className="space-y-2">
                  <h3 className="font-semibold">📍 Широкая сеть</h3>
                  <p className="text-sm text-base-content/70">
                    180+ городов и населенных пунктов по всей Сербии
                  </p>
                </div>
                <div className="space-y-2">
                  <h3 className="font-semibold">💰 Наложенный платеж</h3>
                  <p className="text-sm text-base-content/70">
                    Безопасная оплата при получении с комиссией всего 45 RSD
                  </p>
                </div>
                <div className="space-y-2">
                  <h3 className="font-semibold">🛡️ Страхование</h3>
                  <p className="text-sm text-base-content/70">
                    Базовое страхование до 15,000 RSD включено
                  </p>
                </div>
                <div className="space-y-2">
                  <h3 className="font-semibold">📱 SMS уведомления</h3>
                  <p className="text-sm text-base-content/70">
                    Информирование о всех этапах доставки
                  </p>
                </div>
                <div className="space-y-2">
                  <h3 className="font-semibold">🏪 Склад Sve Tu</h3>
                  <p className="text-sm text-base-content/70">
                    Бесплатный самовывоз с возможностью примерки
                  </p>
                </div>
              </div>
            </div>
          </div>
        )}

        {activeTab === 'flow' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Полный процесс оформления доставки</h2>
              <p>
                Пошаговый процесс выбора и настройки доставки Post Express с
                автоматическим расчетом стоимости и валидацией данных.
              </p>
            </div>
            <PostExpressDeliveryFlow
              onDeliveryComplete={handleDeliveryComplete}
              orderWeight={2.5}
              orderValue={3500}
              allowCOD={true}
            />
            {deliveryData && (
              <div className="alert alert-success">
                <TruckIcon className="w-5 h-5" />
                <div>
                  <h4 className="font-semibold">Доставка настроена!</h4>
                  <p className="text-sm">
                    Способ: {deliveryData.method},
                    {deliveryData.rate &&
                      ` Стоимость: ${deliveryData.rate.total_price || 0} RSD`}
                  </p>
                </div>
              </div>
            )}
          </div>
        )}

        {activeTab === 'selector' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Выбор способа доставки</h2>
              <p>
                Красивые карточки для выбора между курьерской доставкой,
                отделением Post Express или складом Sve Tu.
              </p>
            </div>
            <PostExpressDeliverySelector
              onMethodChange={(method) =>
                console.log('Selected method:', method)
              }
              weight={1.5}
              insuranceAmount={2000}
              hasCOD={false}
              recipientCity="Белград"
            />
          </div>
        )}

        {activeTab === 'calculator' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Калькулятор стоимости доставки</h2>
              <p>
                Интерактивный калькулятор для расчета точной стоимости доставки
                с учетом всех параметров и дополнительных услуг.
              </p>
            </div>
            <PostExpressRateCalculator
              onRateCalculated={(rate) => console.log('Calculated rate:', rate)}
              initialParams={{
                weight: 1.2,
                declaredValue: 2500,
                recipientCity: 'Суботица',
              }}
            />
          </div>
        )}

        {activeTab === 'tracking' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Отслеживание посылки</h2>
              <p>
                Современный интерфейс для отслеживания статуса доставки с
                подробной историей событий и автообновлением.
              </p>
            </div>
            <PostExpressTracker
              initialTrackingNumber=""
              onTrackingUpdate={(shipment) =>
                console.log('Tracking update:', shipment)
              }
            />
          </div>
        )}

        {activeTab === 'pickup' && (
          <div className="space-y-8">
            <div className="prose max-w-none">
              <h2>Код самовывоза</h2>
              <p>
                Удобный интерфейс для отображения кода самовывоза с QR кодом,
                инструкциями и информацией о складе.
              </p>
            </div>
            <PostExpressPickupCode
              pickupOrder={mockPickupOrder}
              onStatusUpdate={(status) =>
                console.log('Status updated:', status)
              }
            />
          </div>
        )}
      </div>

      {/* CTA Section */}
      <div className="bg-gradient-to-r from-primary to-secondary text-primary-content mt-8 sm:mt-16">
        <div className="container mx-auto px-4 py-6 sm:py-12">
          <div className="text-center">
            <h2 className="text-xl sm:text-3xl font-bold mb-2 sm:mb-4">
              Готовы подключить Post Express?
            </h2>
            <p className="text-sm sm:text-xl mb-4 sm:mb-8 opacity-90">
              Интегрируйте национального почтового оператора и расширьте
              географию продаж
            </p>
            <div className="flex flex-col sm:flex-row gap-2 sm:gap-4 justify-center">
              <button className="btn btn-sm sm:btn-lg bg-white text-primary hover:bg-white/90">
                <TruckIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                Начать интеграцию
              </button>
              <button className="btn btn-sm sm:btn-lg btn-outline border-white text-white hover:bg-white/20">
                <MapPinIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                Карта покрытия
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
