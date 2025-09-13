'use client';

import React, { useState } from 'react';
import {
  TruckIcon,
  MapPinIcon,
  BuildingOfficeIcon,
  ClockIcon,
  CurrencyDollarIcon,
  CheckCircleIcon,
  InformationCircleIcon,
} from '@heroicons/react/24/outline';

interface DeliveryProvider {
  id: string;
  name: string;
  logo?: string;
  methods: DeliveryMethod[];
}

interface DeliveryMethod {
  id: string;
  name: string;
  description: string;
  icon: React.ElementType;
  basePrice: number;
  estimatedDays: string;
  codAvailable: boolean;
  insuranceIncluded: number;
}

interface DeliverySelection {
  providerId: string;
  methodId: string;
  price: number;
}

interface DeliverySelectorProps {
  storefrontId: number;
  storefrontName: string;
  subtotal: number;
  weight?: number;
  onDeliveryChange: (selection: DeliverySelection | null) => void;
}

const providers: DeliveryProvider[] = [
  {
    id: 'post_express',
    name: 'Post Express',
    logo: '/images/post-express-logo.png',
    methods: [
      {
        id: 'courier',
        name: 'Курьерская доставка',
        description: 'Доставка на адрес получателя',
        icon: TruckIcon,
        basePrice: 340,
        estimatedDays: '1-2',
        codAvailable: true,
        insuranceIncluded: 15000,
      },
      {
        id: 'office',
        name: 'Почтовое отделение',
        description: 'Самовывоз из отделения',
        icon: BuildingOfficeIcon,
        basePrice: 290,
        estimatedDays: '1-2',
        codAvailable: true,
        insuranceIncluded: 15000,
      },
    ],
  },
  {
    id: 'bex',
    name: 'BEX Express',
    logo: '/images/bex-logo.png',
    methods: [
      {
        id: 'standard',
        name: 'Стандартная доставка',
        description: 'Доставка курьером BEX',
        icon: TruckIcon,
        basePrice: 380,
        estimatedDays: '1-3',
        codAvailable: true,
        insuranceIncluded: 20000,
      },
      {
        id: 'parcel_shop',
        name: 'Parcel Shop',
        description: '500+ пунктов выдачи',
        icon: MapPinIcon,
        basePrice: 320,
        estimatedDays: '1-3',
        codAvailable: false,
        insuranceIncluded: 20000,
      },
    ],
  },
  {
    id: 'sve_tu',
    name: 'Самовывоз Sve Tu',
    methods: [
      {
        id: 'warehouse',
        name: 'Склад в Нови Саде',
        description: 'Бесплатно при заказе от 2000 RSD',
        icon: MapPinIcon,
        basePrice: 0,
        estimatedDays: 'Сразу',
        codAvailable: false,
        insuranceIncluded: 0,
      },
    ],
  },
];

export default function DeliverySelector({
  storefrontId,
  storefrontName,
  subtotal,
  weight = 1,
  onDeliveryChange,
}: DeliverySelectorProps) {
  // const t = useTranslations('cart');
  const [selectedProvider, setSelectedProvider] = useState<string | null>(null);
  const [selectedMethod, setSelectedMethod] = useState<string | null>(null);
  const [isExpanded, setIsExpanded] = useState(false);

  // Рассчитываем цену доставки Post Express на основе веса
  const calculatePostExpressPrice = (basePrice: number): number => {
    if (weight <= 2) return basePrice;
    if (weight <= 5) return basePrice + 110; // 450 RSD для курьера
    if (weight <= 10) return basePrice + 240; // 580 RSD для курьера
    return basePrice + 450; // 790 RSD для курьера
  };

  // Рассчитываем финальную цену доставки
  const calculateDeliveryPrice = (
    provider: DeliveryProvider,
    method: DeliveryMethod
  ): number => {
    let price = method.basePrice;

    // Специальная логика для Post Express
    if (provider.id === 'post_express') {
      price = calculatePostExpressPrice(method.basePrice);
    }

    // Бесплатная доставка для самовывоза при заказе от 2000 RSD
    if (provider.id === 'sve_tu' && subtotal >= 2000) {
      price = 0;
    } else if (provider.id === 'sve_tu' && subtotal < 2000) {
      price = 100; // Платный самовывоз при заказе менее 2000 RSD
    }

    // Бесплатная доставка при заказе от 5000 RSD (кроме самовывоза)
    if (provider.id !== 'sve_tu' && subtotal >= 5000) {
      price = 0;
    }

    return price;
  };

  // Обработка выбора доставки
  const handleSelection = (providerId: string, methodId: string) => {
    const provider = providers.find((p) => p.id === providerId);
    const method = provider?.methods.find((m) => m.id === methodId);

    if (provider && method) {
      const price = calculateDeliveryPrice(provider, method);
      setSelectedProvider(providerId);
      setSelectedMethod(methodId);
      onDeliveryChange({
        providerId,
        methodId,
        price,
      });
    }
  };

  // Получаем выбранную доставку
  const getSelectedDelivery = () => {
    if (!selectedProvider || !selectedMethod) return null;
    const provider = providers.find((p) => p.id === selectedProvider);
    const method = provider?.methods.find((m) => m.id === selectedMethod);
    if (!provider || !method) return null;
    return {
      provider,
      method,
      price: calculateDeliveryPrice(provider, method),
    };
  };

  const selected = getSelectedDelivery();

  return (
    <div className="card bg-base-200">
      <div className="card-body">
        <div className="flex items-center justify-between mb-4">
          <h3 className="text-lg font-semibold flex items-center gap-2">
            <TruckIcon className="w-5 h-5" />
            Доставка для {storefrontName}
          </h3>
          {selected && !isExpanded && (
            <button
              onClick={() => setIsExpanded(true)}
              className="btn btn-ghost btn-sm"
            >
              Изменить
            </button>
          )}
        </div>

        {/* Компактный вид выбранной доставки */}
        {selected && !isExpanded && (
          <div className="p-4 bg-base-100 rounded-lg">
            <div className="flex justify-between items-start">
              <div className="flex gap-3">
                <selected.method.icon className="w-5 h-5 text-primary mt-1" />
                <div>
                  <p className="font-medium">{selected.provider.name}</p>
                  <p className="text-sm text-base-content/60">
                    {selected.method.name}
                  </p>
                  <p className="text-xs text-base-content/60 mt-1">
                    {selected.method.estimatedDays}{' '}
                    {selected.method.estimatedDays === 'Сразу' ? '' : 'дня'}
                  </p>
                </div>
              </div>
              <div className="text-right">
                {selected.price === 0 ? (
                  <span className="text-success font-semibold">Бесплатно</span>
                ) : (
                  <span className="font-semibold">{selected.price} RSD</span>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Развернутый вид выбора доставки */}
        {(!selected || isExpanded) && (
          <div className="space-y-4">
            {providers.map((provider) => (
              <div key={provider.id} className="space-y-2">
                <div className="flex items-center gap-2">
                  {provider.logo && (
                    <img
                      src={provider.logo}
                      alt={provider.name}
                      className="h-6 object-contain"
                    />
                  )}
                  <h4 className="font-medium text-sm">{provider.name}</h4>
                </div>

                <div className="grid gap-2">
                  {provider.methods.map((method) => {
                    const price = calculateDeliveryPrice(provider, method);
                    const isSelected =
                      selectedProvider === provider.id &&
                      selectedMethod === method.id;

                    return (
                      <label
                        key={method.id}
                        className={`card bg-base-100 cursor-pointer transition-all ${
                          isSelected ? 'ring-2 ring-primary' : 'hover:shadow-md'
                        }`}
                      >
                        <input
                          type="radio"
                          name={`delivery-${storefrontId}`}
                          className="hidden"
                          checked={isSelected}
                          onChange={() =>
                            handleSelection(provider.id, method.id)
                          }
                        />
                        <div className="card-body p-3">
                          <div className="flex items-start gap-3">
                            <method.icon
                              className={`w-5 h-5 mt-0.5 ${
                                isSelected
                                  ? 'text-primary'
                                  : 'text-base-content/60'
                              }`}
                            />
                            <div className="flex-1">
                              <div className="flex justify-between items-start">
                                <div>
                                  <p className="font-medium text-sm">
                                    {method.name}
                                  </p>
                                  <p className="text-xs text-base-content/60 mt-0.5">
                                    {method.description}
                                  </p>
                                </div>
                                <div className="text-right ml-4">
                                  {price === 0 ? (
                                    <span className="text-success font-semibold text-sm">
                                      Бесплатно
                                    </span>
                                  ) : (
                                    <span className="font-semibold text-sm">
                                      {price} RSD
                                    </span>
                                  )}
                                </div>
                              </div>

                              <div className="flex flex-wrap gap-2 mt-2">
                                <span className="badge badge-sm badge-ghost">
                                  <ClockIcon className="w-3 h-3 mr-1" />
                                  {method.estimatedDays}{' '}
                                  {method.estimatedDays === 'Сразу'
                                    ? ''
                                    : 'дня'}
                                </span>
                                {method.codAvailable && (
                                  <span className="badge badge-sm badge-ghost">
                                    <CurrencyDollarIcon className="w-3 h-3 mr-1" />
                                    Наложенный платеж
                                  </span>
                                )}
                                {method.insuranceIncluded > 0 && (
                                  <span className="badge badge-sm badge-ghost">
                                    <CheckCircleIcon className="w-3 h-3 mr-1" />
                                    Страховка до{' '}
                                    {method.insuranceIncluded / 1000}k
                                  </span>
                                )}
                              </div>
                            </div>
                          </div>
                        </div>
                      </label>
                    );
                  })}
                </div>
              </div>
            ))}

            {/* Информация о бесплатной доставке */}
            {subtotal < 5000 && (
              <div className="alert alert-info py-2">
                <InformationCircleIcon className="w-5 h-5" />
                <span className="text-sm">
                  Бесплатная доставка при заказе от 5000 RSD (осталось{' '}
                  {(5000 - subtotal).toFixed(0)} RSD)
                </span>
              </div>
            )}

            {selected && isExpanded && (
              <button
                onClick={() => setIsExpanded(false)}
                className="btn btn-primary btn-sm btn-block"
              >
                Применить выбор
              </button>
            )}
          </div>
        )}

        {/* Информация о весе */}
        {weight > 1 && (
          <div className="text-xs text-base-content/60 mt-2">
            Общий вес: {weight.toFixed(1)} кг
          </div>
        )}
      </div>
    </div>
  );
}
