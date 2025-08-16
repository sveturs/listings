'use client';

import { useState, useEffect } from 'react';
import {
  TruckIcon,
  MapPinIcon,
  BuildingStorefrontIcon,
  CheckIcon,
  InformationCircleIcon,
  BanknotesIcon,
  ShieldCheckIcon,
  ClockIcon,
} from '@heroicons/react/24/outline';
import { useTranslations } from 'next-intl';

interface DeliveryMethod {
  id: string;
  name: string;
  description: string;
  priceRange: string;
  timeRange: string;
  icon: React.ComponentType<any>;
  features: string[];
  popular?: boolean;
  available?: boolean;
}

interface Props {
  onMethodChange: (method: string) => void;
  selectedMethod?: string;
  weight?: number;
  insuranceAmount?: number;
  hasCOD?: boolean;
  recipientCity?: string;
  className?: string;
}

export default function PostExpressDeliverySelector({
  onMethodChange,
  selectedMethod = '',
  weight = 0,
  insuranceAmount = 0,
  hasCOD = false,
  recipientCity,
  className = '',
}: Props) {
  // const t = useTranslations('delivery');
  const [calculatedRates, setCalculatedRates] = useState<Record<string, any>>(
    {}
  );
  const [loading, setLoading] = useState(false);

  const deliveryMethods: DeliveryMethod[] = [
    {
      id: 'courier',
      name: 'Курьерская доставка',
      description: 'Доставка на адрес получателя курьером Post Express',
      priceRange: '340-790 RSD',
      timeRange: '1-2 рабочих дня',
      icon: TruckIcon,
      popular: true,
      available: true,
      features: [
        'Доставка "Данас за сутра"',
        'Доставка до 19:00',
        'SMS уведомления',
        'Страхование до 15,000 RSD включено',
        'Поддержка наложенного платежа',
      ],
    },
    {
      id: 'pickup_point',
      name: 'Почтовое отделение',
      description: 'Самовывоз из ближайшего отделения Post Express',
      priceRange: '340-790 RSD',
      timeRange: '1-2 рабочих дня',
      icon: MapPinIcon,
      available: true,
      features: [
        '180+ городов покрытия',
        'Сеть почтовых отделений',
        'Хранение 5 рабочих дней',
        'Удобные часы работы',
        'Поддержка наложенного платежа',
      ],
    },
    {
      id: 'warehouse_pickup',
      name: 'Склад Sve Tu',
      description: 'Самовывоз со склада в Нови Саде',
      priceRange: 'Бесплатно',
      timeRange: 'Готов к выдаче',
      icon: BuildingStorefrontIcon,
      available: true,
      features: [
        'Бесплатный самовывоз',
        'Возможность примерки',
        'Консолидация заказов',
        'Пн-Пт 09:00-19:00, Сб 10:00-16:00',
        'Микија Манојловића 53, Нови Сад',
      ],
    },
  ];

  // Расчет стоимости доставки
  useEffect(() => {
    if (weight > 0 && recipientCity) {
      calculateDeliveryRates();
    }
  }, [weight, insuranceAmount, hasCOD, recipientCity]);

  const calculateDeliveryRates = async () => {
    setLoading(true);
    try {
      const response = await fetch('/api/v1/postexpress/calculate-rate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          weight_kg: weight,
          declared_value: insuranceAmount,
          cod_amount: hasCOD ? 100 : 0, // Примерная сумма
          recipient_city: recipientCity,
          sender_postal_code: '21000', // Склад в Нови Саде
          recipient_postal_code: '11000', // Примерный код для расчета
        }),
      });

      const data = await response.json();
      if (data.success) {
        setCalculatedRates({
          courier: data.data,
          pickup_point: data.data,
          warehouse_pickup: { total_price: 0, base_price: 0 },
        });
      }
    } catch (error) {
      console.error('Failed to calculate delivery rates:', error);
    } finally {
      setLoading(false);
    }
  };

  const getMethodPrice = (methodId: string) => {
    const rate = calculatedRates[methodId];
    if (!rate) return null;

    if (methodId === 'warehouse_pickup') {
      return 'Бесплатно';
    }

    return `${rate.total_price?.toFixed(0) || '0'} RSD`;
  };

  const handleMethodSelect = (methodId: string) => {
    onMethodChange(methodId);
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Заголовок */}
      <div className="text-center">
        <h3 className="text-xl font-bold mb-2">Выберите способ доставки</h3>
        <p className="text-base-content/70">
          Post Express - национальный почтовый оператор Сербии
        </p>
      </div>

      {/* Карточки способов доставки */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-4">
        {deliveryMethods.map((method) => {
          const isSelected = selectedMethod === method.id;
          const calculatedPrice = getMethodPrice(method.id);

          return (
            <div
              key={method.id}
              className={`
                card cursor-pointer transition-all duration-200 relative overflow-hidden
                ${
                  isSelected
                    ? 'ring-2 ring-primary shadow-xl scale-[1.02] bg-primary/5'
                    : 'hover:shadow-lg hover:scale-[1.01] bg-base-100'
                }
                ${!method.available ? 'opacity-50 cursor-not-allowed' : ''}
                ${method.popular ? 'border-primary border-2' : ''}
              `}
              onClick={() => method.available && handleMethodSelect(method.id)}
            >
              {/* Популярный бейдж */}
              {method.popular && (
                <div className="absolute -top-2 left-1/2 -translate-x-1/2 z-10">
                  <span className="badge badge-primary badge-sm px-3 shadow-lg">
                    Популярно
                  </span>
                </div>
              )}

              <div className="card-body p-6">
                {/* Хедер с иконкой */}
                <div className="flex items-start justify-between mb-4">
                  <div
                    className={`
                    p-3 rounded-xl transition-colors
                    ${
                      isSelected
                        ? 'bg-primary text-primary-content shadow-lg'
                        : 'bg-primary/10 text-primary'
                    }
                  `}
                  >
                    <method.icon className="w-6 h-6" />
                  </div>
                  {isSelected && (
                    <div className="p-1.5 bg-success text-success-content rounded-full shadow-lg">
                      <CheckIcon className="w-5 h-5" />
                    </div>
                  )}
                </div>

                {/* Название и описание */}
                <div className="mb-4">
                  <h4 className="font-semibold text-lg mb-2">{method.name}</h4>
                  <p className="text-sm text-base-content/70">
                    {method.description}
                  </p>
                </div>

                {/* Цена и время */}
                <div className="grid grid-cols-2 gap-4 mb-4 p-3 bg-base-200/50 rounded-lg">
                  <div>
                    <div className="text-xs text-base-content/60 mb-1">
                      Стоимость
                    </div>
                    <div className="font-bold text-lg text-primary">
                      {loading ? (
                        <span className="loading loading-dots loading-sm"></span>
                      ) : (
                        calculatedPrice || method.priceRange
                      )}
                    </div>
                  </div>
                  <div>
                    <div className="text-xs text-base-content/60 mb-1">
                      Время доставки
                    </div>
                    <div className="font-semibold text-sm">
                      {method.timeRange}
                    </div>
                  </div>
                </div>

                {/* Особенности (показываем всегда) */}
                <div className="space-y-2">
                  {method.features
                    .slice(0, isSelected ? method.features.length : 3)
                    .map((feature, index) => (
                      <div key={index} className="flex items-start gap-2">
                        <CheckIcon className="w-4 h-4 text-success mt-0.5 flex-shrink-0" />
                        <span className="text-sm text-base-content/80">
                          {feature}
                        </span>
                      </div>
                    ))}
                  {!isSelected && method.features.length > 3 && (
                    <div className="text-xs text-base-content/60 mt-2">
                      +{method.features.length - 3} преимуществ
                    </div>
                  )}
                </div>

                {/* Недоступен */}
                {!method.available && (
                  <div className="absolute inset-0 bg-base-300/80 flex items-center justify-center">
                    <div className="text-center">
                      <InformationCircleIcon className="w-8 h-8 mx-auto mb-2 text-base-content/60" />
                      <div className="text-sm font-medium">
                        Временно недоступно
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>

      {/* Информация о выбранном методе */}
      {selectedMethod && (
        <div className="alert alert-info">
          <InformationCircleIcon className="w-5 h-5" />
          <div>
            <h4 className="font-semibold">
              Выбран способ:{' '}
              {deliveryMethods.find((m) => m.id === selectedMethod)?.name}
            </h4>
            <p className="text-sm mt-1">
              {selectedMethod === 'courier' &&
                'Курьер доставит посылку на указанный адрес. Вы получите SMS с уведомлением о доставке.'}
              {selectedMethod === 'pickup_point' &&
                'Выберите удобное почтовое отделение для получения посылки. Хранение до 5 рабочих дней.'}
              {selectedMethod === 'warehouse_pickup' &&
                'Заберите заказ бесплатно со склада Sve Tu в Нови Саде. Возможна примерка товаров.'}
            </p>
          </div>
        </div>
      )}

      {/* Дополнительные услуги */}
      <div className="card bg-gradient-to-r from-base-200/50 to-base-300/30">
        <div className="card-body p-6">
          <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
            <ShieldCheckIcon className="w-5 h-5 text-primary" />
            Дополнительные услуги Post Express
          </h4>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
            {/* Страхование */}
            <div className="flex items-center gap-3 p-3 bg-base-100 rounded-lg">
              <div className="p-2 bg-green-100 rounded-lg">
                <ShieldCheckIcon className="w-5 h-5 text-green-600" />
              </div>
              <div className="flex-1">
                <div className="font-medium text-sm">Страхование</div>
                <div className="text-xs text-base-content/60">
                  До 15,000 RSD включено
                </div>
              </div>
            </div>

            {/* Наложенный платеж */}
            <div className="flex items-center gap-3 p-3 bg-base-100 rounded-lg">
              <div className="p-2 bg-blue-100 rounded-lg">
                <BanknotesIcon className="w-5 h-5 text-blue-600" />
              </div>
              <div className="flex-1">
                <div className="font-medium text-sm">Наложенный платеж</div>
                <div className="text-xs text-base-content/60">
                  Комиссия 45 RSD
                </div>
              </div>
            </div>

            {/* Быстрая доставка */}
            <div className="flex items-center gap-3 p-3 bg-base-100 rounded-lg">
              <div className="p-2 bg-orange-100 rounded-lg">
                <ClockIcon className="w-5 h-5 text-orange-600" />
              </div>
              <div className="flex-1">
                <div className="font-medium text-sm">Данас за сутра</div>
                <div className="text-xs text-base-content/60">
                  До 19:00 следующего дня
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
