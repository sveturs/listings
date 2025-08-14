'use client';

import { useState } from 'react';
import {
  TruckIcon,
  MapPinIcon,
  ClockIcon,
  CheckIcon,
  InformationCircleIcon,
} from '@heroicons/react/24/outline';

interface DeliveryMethod {
  id: string;
  name: string;
  description: string;
  price: string;
  time: string;
  icon: React.ComponentType<any>;
  features: string[];
  popular?: boolean;
}

interface Props {
  onMethodChange: (method: string) => void;
  selectedMethod: string;
}

export default function DeliveryMethodSelector({
  onMethodChange,
  selectedMethod,
}: Props) {
  const [showDetails, setShowDetails] = useState<string | null>(null);

  const deliveryMethods: DeliveryMethod[] = [
    {
      id: 'courier',
      name: 'Курьерская доставка',
      description: 'Доставка курьером BEX Express на ваш адрес',
      price: '250-450 RSD',
      time: '2-3 рабочих дня',
      icon: TruckIcon,
      popular: true,
      features: [
        'Доставка до двери',
        'Отслеживание в реальном времени',
        'СМС уведомления',
        'Возможность оплаты при получении',
      ],
    },
    {
      id: 'parcel-shop',
      name: 'Пункт выдачи',
      description: 'Самовывоз из ближайшего пункта BEX',
      price: '180-280 RSD',
      time: '1-2 рабочих дня',
      icon: MapPinIcon,
      features: [
        '200+ пунктов по Сербии',
        'Удобное расположение',
        'Хранение до 7 дней',
        'Гибкий график работы',
      ],
    },
    {
      id: 'express',
      name: 'Экспресс доставка',
      description: 'Срочная доставка в тот же день',
      price: '800-1200 RSD',
      time: '3-6 часов',
      icon: ClockIcon,
      features: [
        'Доставка в день заказа',
        'Приоритетная обработка',
        'Персональный курьер',
        'Доступно в крупных городах',
      ],
    },
  ];

  const handleMethodSelect = (methodId: string) => {
    onMethodChange(methodId);
  };

  return (
    <div className="space-y-6">
      {/* Method Cards */}
      <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3 sm:gap-4">
        {deliveryMethods.map((method) => (
          <div
            key={method.id}
            className={`
              card cursor-pointer transition-all
              ${
                selectedMethod === method.id
                  ? 'ring-2 ring-primary shadow-xl scale-[1.02]'
                  : 'hover:shadow-lg hover:scale-[1.01]'
              }
              ${method.popular ? 'relative' : ''}
            `}
            onClick={() => handleMethodSelect(method.id)}
          >
            {method.popular && (
              <div className="absolute -top-2 sm:-top-3 left-1/2 -translate-x-1/2 z-10">
                <span className="badge badge-primary badge-xs sm:badge-sm px-2 sm:px-3">
                  Популярно
                </span>
              </div>
            )}

            <div className="card-body p-4 sm:p-6">
              {/* Header */}
              <div className="flex items-start justify-between">
                <div
                  className={`
                  p-2 sm:p-3 rounded-lg
                  ${
                    selectedMethod === method.id
                      ? 'bg-primary text-primary-content'
                      : 'bg-base-200'
                  }
                `}
                >
                  <method.icon className="w-5 h-5 sm:w-6 sm:h-6" />
                </div>
                {selectedMethod === method.id && (
                  <div className="p-1 bg-success text-success-content rounded-full">
                    <CheckIcon className="w-4 h-4 sm:w-5 sm:h-5" />
                  </div>
                )}
              </div>

              {/* Content */}
              <div className="mt-3 sm:mt-4">
                <h3 className="font-semibold text-base sm:text-lg">
                  {method.name}
                </h3>
                <p className="text-xs sm:text-sm text-base-content/60 mt-1">
                  {method.description}
                </p>
              </div>

              {/* Price and Time */}
              <div className="grid grid-cols-2 gap-2 sm:gap-3 mt-3 sm:mt-4 pt-3 sm:pt-4 border-t">
                <div>
                  <div className="text-xs text-base-content/60">Стоимость</div>
                  <div className="font-semibold text-sm sm:text-base text-primary">
                    {method.price}
                  </div>
                </div>
                <div>
                  <div className="text-xs text-base-content/60">Время</div>
                  <div className="font-semibold text-sm sm:text-base">
                    {method.time}
                  </div>
                </div>
              </div>

              {/* Features (shown when selected) */}
              {selectedMethod === method.id && (
                <div className="mt-3 sm:mt-4 pt-3 sm:pt-4 border-t space-y-1 sm:space-y-2">
                  {method.features.map((feature, index) => (
                    <div key={index} className="flex items-start gap-2">
                      <CheckIcon className="w-3 h-3 sm:w-4 sm:h-4 text-success mt-0.5 flex-shrink-0" />
                      <span className="text-xs sm:text-sm">{feature}</span>
                    </div>
                  ))}
                </div>
              )}

              {/* Details Button */}
              <button
                className="btn btn-ghost btn-xs sm:btn-sm mt-3 sm:mt-4"
                onClick={(e) => {
                  e.stopPropagation();
                  setShowDetails(showDetails === method.id ? null : method.id);
                }}
              >
                <InformationCircleIcon className="w-3 h-3 sm:w-4 sm:h-4" />
                <span className="hidden sm:inline">Подробнее</span>
              </button>
            </div>
          </div>
        ))}
      </div>

      {/* Selected Method Details */}
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
                'Курьер свяжется с вами для уточнения времени доставки.'}
              {selectedMethod === 'parcel-shop' &&
                'Выберите удобный пункт выдачи на карте ниже.'}
              {selectedMethod === 'express' &&
                'Доступно только для заказов до 14:00 в рабочие дни.'}
            </p>
          </div>
        </div>
      )}

      {/* Additional Options */}
      <div className="card bg-base-200">
        <div className="card-body p-4 sm:p-6">
          <h4 className="font-semibold text-base sm:text-lg mb-3">
            Дополнительные услуги
          </h4>
          <div className="space-y-3">
            <label className="flex items-start sm:items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                className="checkbox checkbox-primary checkbox-sm sm:checkbox-md mt-0.5 sm:mt-0"
              />
              <div className="flex-1">
                <div className="font-medium text-sm sm:text-base">
                  Страхование посылки
                </div>
                <div className="text-xs sm:text-sm text-base-content/60">
                  До 100,000 RSD (+2% от суммы)
                </div>
              </div>
            </label>
            <label className="flex items-start sm:items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                className="checkbox checkbox-primary checkbox-sm sm:checkbox-md mt-0.5 sm:mt-0"
              />
              <div className="flex-1">
                <div className="font-medium text-sm sm:text-base">
                  Возврат документов
                </div>
                <div className="text-xs sm:text-sm text-base-content/60">
                  Подписанные документы вернутся к вам (+150 RSD)
                </div>
              </div>
            </label>
            <label className="flex items-start sm:items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                className="checkbox checkbox-primary checkbox-sm sm:checkbox-md mt-0.5 sm:mt-0"
                defaultChecked
              />
              <div className="flex-1">
                <div className="font-medium text-sm sm:text-base">
                  СМС уведомления
                </div>
                <div className="text-xs sm:text-sm text-base-content/60">
                  Отслеживание статуса доставки (Бесплатно)
                </div>
              </div>
            </label>
          </div>
        </div>
      </div>
    </div>
  );
}
