'use client';

import { useState, useEffect } from 'react';
import {
  TruckIcon,
  HomeModernIcon,
  BuildingStorefrontIcon,
  CheckCircleIcon,
  ClockIcon,
  CurrencyDollarIcon,
  MapPinIcon,
  ScaleIcon,
  ShieldCheckIcon,
} from '@heroicons/react/24/outline';
import { motion, AnimatePresence } from 'framer-motion';
import configManager from '@/config';

interface DeliveryOption {
  id: string;
  name: string;
  description: string;
  icon: React.ComponentType<{ className?: string }>;
  price: number;
  estimatedDays: string;
  features: string[];
  color: string;
  popular?: boolean;
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

export default function BEXDeliverySelector({
  onMethodChange,
  selectedMethod = 'courier',
  weight = 1,
  insuranceAmount = 0,
  hasCOD = false,
  recipientCity,
  className = '',
}: Props) {
  const [calculatedPrices, setCalculatedPrices] = useState<
    Record<string, number>
  >({});
  const [loading, setLoading] = useState(false);

  // Расчет категории отправления по весу
  const getShipmentCategory = (weightKg: number) => {
    if (weightKg <= 0.5) return 1; // Документы
    if (weightKg <= 1) return 2; // Посылка до 1кг
    if (weightKg <= 2) return 3; // Посылка до 2кг
    return 31; // Посылка за кг
  };

  // Локальный расчет базовой цены
  const calculateBasePrice = (method: string, weightKg: number) => {
    if (method === 'warehouse_pickup') return 0;

    const category = getShipmentCategory(weightKg);
    let basePrice = 0;

    switch (category) {
      case 1:
        basePrice = 250; // Документы
        break;
      case 2:
        basePrice = 350; // до 1кг
        break;
      case 3:
        basePrice = 450; // до 2кг
        break;
      case 31:
        basePrice = 450 + Math.ceil(weightKg - 2) * 100; // 450 + 100/кг после 2кг
        break;
    }

    // Добавляем комиссию за наложенный платеж
    if (hasCOD) {
      basePrice += 150 + Math.ceil(insuranceAmount * 0.01);
    }

    // Добавляем стоимость страхования
    if (insuranceAmount > 0) {
      basePrice += Math.ceil(insuranceAmount * 0.005);
    }

    return basePrice;
  };

  const deliveryOptions: DeliveryOption[] = [
    {
      id: 'courier',
      name: 'Курьерская доставка',
      description: 'Доставка курьером до двери',
      icon: TruckIcon,
      price: calculateBasePrice('courier', weight),
      estimatedDays: weight > 5 ? '2-3' : '1-2',
      features: [
        'Доставка до двери',
        'СМС-уведомления',
        'Отслеживание в реальном времени',
        'Возможность оплаты при получении',
      ],
      color: 'primary',
      popular: true,
    },
    {
      id: 'pickup_point',
      name: 'Пункт выдачи BEX',
      description: 'Самовывоз из пункта выдачи BEX',
      icon: BuildingStorefrontIcon,
      price: calculateBasePrice('pickup_point', weight) * 0.8, // Скидка 20%
      estimatedDays: '1-2',
      features: [
        'Более 200 пунктов по Сербии',
        'Удобное время получения',
        'СМС-уведомления о прибытии',
        'Хранение до 7 дней',
      ],
      color: 'secondary',
    },
    {
      id: 'warehouse_pickup',
      name: 'Самовывоз со склада',
      description: 'Забрать со склада отправителя',
      icon: HomeModernIcon,
      price: 0,
      estimatedDays: '0',
      features: [
        'Бесплатно',
        'Мгновенное получение',
        'Адрес: Мике Манојловића 53, Нови Сад',
        'Пн-Пт: 9:00-18:00',
      ],
      color: 'accent',
    },
  ];

  // Загрузка точных цен с сервера
  useEffect(() => {
    const fetchRates = async () => {
      if (!recipientCity) return;

      setLoading(true);
      try {
        const apiUrl = configManager.get('api.url');
        const response = await fetch(`${apiUrl}/api/v1/bex/calculate-rate`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            weight_kg: weight,
            shipment_category: getShipmentCategory(weight),
            cod_amount: hasCOD ? insuranceAmount : 0,
            insurance_amount: insuranceAmount,
            recipient_city: recipientCity,
          }),
        });

        if (response.ok) {
          const data = await response.json();
          if (data.success && data.data) {
            setCalculatedPrices({
              courier: data.data.total_price,
              pickup_point: data.data.total_price * 0.8,
              warehouse_pickup: 0,
            });
          }
        }
      } catch (error) {
        console.error('Failed to fetch BEX rates:', error);
      } finally {
        setLoading(false);
      }
    };

    fetchRates();
  }, [weight, insuranceAmount, hasCOD, recipientCity]);

  const selectedOption = deliveryOptions.find(
    (opt) => opt.id === selectedMethod
  );

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Опции доставки */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {deliveryOptions.map((option) => {
          const Icon = option.icon;
          const isSelected = selectedMethod === option.id;
          const price = calculatedPrices[option.id] ?? option.price;

          return (
            <motion.div
              key={option.id}
              whileHover={{ scale: 1.02 }}
              whileTap={{ scale: 0.98 }}
              onClick={() => onMethodChange(option.id)}
              className={`
                relative cursor-pointer rounded-xl border-2 p-6 transition-all
                ${
                  isSelected
                    ? `border-${option.color} bg-${option.color}/5 shadow-lg`
                    : 'border-base-300 hover:border-base-content/20'
                }
              `}
            >
              {option.popular && (
                <div className="absolute -top-3 left-1/2 -translate-x-1/2">
                  <span className="badge badge-primary badge-sm">
                    Популярный выбор
                  </span>
                </div>
              )}

              <div className="flex items-start gap-4">
                <div
                  className={`
                  p-3 rounded-lg bg-${option.color}/10
                  ${isSelected ? `text-${option.color}` : 'text-base-content/60'}
                `}
                >
                  <Icon className="w-6 h-6" />
                </div>

                <div className="flex-1">
                  <h3 className="font-semibold text-lg mb-1">{option.name}</h3>
                  <p className="text-sm text-base-content/60 mb-3">
                    {option.description}
                  </p>

                  {/* Цена и время */}
                  <div className="flex items-center gap-4 mb-3">
                    <div className="flex items-center gap-1">
                      <CurrencyDollarIcon className="w-4 h-4 text-primary" />
                      <span className="font-bold text-lg">
                        {loading ? (
                          <span className="loading loading-spinner loading-xs"></span>
                        ) : price === 0 ? (
                          'Бесплатно'
                        ) : (
                          `${price} RSD`
                        )}
                      </span>
                    </div>

                    {option.estimatedDays !== '0' && (
                      <div className="flex items-center gap-1 text-sm">
                        <ClockIcon className="w-4 h-4 text-info" />
                        <span>{option.estimatedDays} дня</span>
                      </div>
                    )}
                  </div>

                  {/* Особенности */}
                  <ul className="space-y-1">
                    {option.features.map((feature, idx) => (
                      <li key={idx} className="flex items-start gap-2 text-sm">
                        <CheckCircleIcon className="w-4 h-4 text-success mt-0.5 flex-shrink-0" />
                        <span className="text-base-content/70">{feature}</span>
                      </li>
                    ))}
                  </ul>
                </div>

                {/* Индикатор выбора */}
                {isSelected && (
                  <CheckCircleIcon
                    className={`w-6 h-6 text-${option.color} flex-shrink-0`}
                  />
                )}
              </div>
            </motion.div>
          );
        })}
      </div>

      {/* Дополнительная информация */}
      <AnimatePresence mode="wait">
        {selectedOption && (
          <motion.div
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            className="card bg-base-200 shadow-lg"
          >
            <div className="card-body p-4">
              <div className="flex items-center gap-3">
                <ShieldCheckIcon className="w-5 h-5 text-success" />
                <div className="flex-1">
                  <h4 className="font-semibold">Гарантии BEX Express</h4>
                  <p className="text-sm text-base-content/60">
                    Страхование груза • Отслеживание 24/7 • Компенсация при
                    утере
                  </p>
                </div>
              </div>

              {weight && (
                <div className="flex items-center gap-3 mt-3 pt-3 border-t border-base-300">
                  <ScaleIcon className="w-5 h-5 text-info" />
                  <div className="text-sm">
                    <span className="text-base-content/60">Вес посылки: </span>
                    <span className="font-medium">{weight} кг</span>
                    <span className="text-base-content/60"> • Категория: </span>
                    <span className="font-medium">
                      {getShipmentCategory(weight) === 1 && 'Документы'}
                      {getShipmentCategory(weight) === 2 && 'Легкая посылка'}
                      {getShipmentCategory(weight) === 3 &&
                        'Стандартная посылка'}
                      {getShipmentCategory(weight) === 31 && 'Тяжелая посылка'}
                    </span>
                  </div>
                </div>
              )}

              {recipientCity && (
                <div className="flex items-center gap-3 mt-2">
                  <MapPinIcon className="w-5 h-5 text-warning" />
                  <div className="text-sm">
                    <span className="text-base-content/60">
                      Город доставки:{' '}
                    </span>
                    <span className="font-medium">{recipientCity}</span>
                  </div>
                </div>
              )}
            </div>
          </motion.div>
        )}
      </AnimatePresence>
    </div>
  );
}
