'use client';

import { useState } from 'react';
import {
  TruckIcon,
  MapPinIcon,
  ClockIcon,
  CurrencyDollarIcon,
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
  company: string;
  color: string;
}

interface Props {
  onMethodChange: (method: string) => void;
  selectedMethod: string;
}

export default function SerbianDeliveryMethodSelector({
  onMethodChange,
  selectedMethod,
}: Props) {
  const [showDetails, setShowDetails] = useState<string | null>(null);

  const deliveryMethods: DeliveryMethod[] = [
    {
      id: 'aks',
      name: 'AKS курирска служба',
      description: 'Најбржа достава у Београду и Новом Саду',
      price: '200-350 РСД',
      time: '1-2 радна дана',
      icon: TruckIcon,
      popular: true,
      company: 'AKS',
      color: 'bg-orange-500',
      features: [
        'Достава на кућну адресу',
        'СМС обавештење',
        'Плаћање готовином при преузимању',
        'Осигурање до 50.000 РСД',
        'Бесплатно преупаковање',
      ],
    },
    {
      id: 'post-express',
      name: 'Post Express пункт',
      description: 'Преузми у једном од 500+ Post Express пунктова',
      price: '150-200 РСД',
      time: '1-3 радна дана',
      icon: MapPinIcon,
      company: 'Post Express',
      color: 'bg-blue-500',
      features: [
        'Мрежа од 500+ пунктова',
        'Радно време до 22:00',
        'СМС и имејл обавештење',
        'Чување 7 дана бесплатно',
        'Могућност промене пункта',
      ],
    },
    {
      id: 'city-express',
      name: 'City Express',
      description: 'Брза достава по целој Србији',
      price: '180-300 РСД',
      time: '2-4 радна дана',
      icon: TruckIcon,
      company: 'City Express',
      color: 'bg-green-500',
      features: [
        'Покривеност 150+ градова',
        'Праћење пошиљке онлајн',
        'Плаћање картицом при преузимању',
        'Осигурање до 100.000 РСД',
      ],
    },
    {
      id: 'yettel-post',
      name: 'Yettel Post пакомат',
      description: 'Аутоматски пунктови за преузимање 24/7',
      price: '120-180 РСД',
      time: '1-2 радна дана',
      icon: MapPinIcon,
      company: 'Yettel Post',
      color: 'bg-purple-500',
      features: [
        'Доступан 24/7',
        'Без чекања у реду',
        'QR код за преузимање',
        'Најјефтинија опција',
        'Чување до 10 дана',
      ],
    },
  ];

  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {deliveryMethods.map((method) => (
          <div
            key={method.id}
            className={`
              card cursor-pointer transition-all duration-300 hover:shadow-lg
              ${
                selectedMethod === method.id
                  ? 'ring-2 ring-primary shadow-lg bg-base-100'
                  : 'bg-base-100 hover:bg-base-50'
              }
            `}
            onClick={() => onMethodChange(method.id)}
          >
            <div className="card-body p-4">
              <div className="flex items-start justify-between mb-3">
                <div className="flex items-center gap-3">
                  <div className={`p-2 rounded-lg ${method.color} text-white`}>
                    <method.icon className="w-6 h-6" />
                  </div>
                  <div>
                    <div className="flex items-center gap-2">
                      <h3 className="font-semibold text-lg">{method.name}</h3>
                      {method.popular && (
                        <div className="badge badge-primary badge-sm">
                          Популарно
                        </div>
                      )}
                    </div>
                    <p className="text-sm text-base-content/70 font-medium">
                      {method.company}
                    </p>
                  </div>
                </div>
                <div
                  className={`
                  w-6 h-6 rounded-full border-2 flex items-center justify-center
                  ${
                    selectedMethod === method.id
                      ? 'border-primary bg-primary'
                      : 'border-base-300'
                  }
                `}
                >
                  {selectedMethod === method.id && (
                    <CheckIcon className="w-4 h-4 text-primary-content" />
                  )}
                </div>
              </div>

              <p className="text-base-content/80 mb-3">{method.description}</p>

              <div className="flex items-center justify-between mb-3">
                <div className="flex items-center gap-1 text-primary">
                  <CurrencyDollarIcon className="w-5 h-5" />
                  <span className="font-semibold">{method.price}</span>
                </div>
                <div className="flex items-center gap-1 text-base-content/60">
                  <ClockIcon className="w-5 h-5" />
                  <span className="text-sm">{method.time}</span>
                </div>
              </div>

              <div className="flex items-center justify-between">
                <button
                  onClick={(e) => {
                    e.stopPropagation();
                    setShowDetails(
                      showDetails === method.id ? null : method.id
                    );
                  }}
                  className="btn btn-ghost btn-sm gap-1"
                >
                  <InformationCircleIcon className="w-4 h-4" />
                  Детаљи
                </button>
              </div>

              {/* Details Panel */}
              {showDetails === method.id && (
                <div className="mt-3 pt-3 border-t">
                  <h4 className="font-medium mb-2">Укључено у цену:</h4>
                  <ul className="space-y-1">
                    {method.features.map((feature, index) => (
                      <li
                        key={index}
                        className="flex items-center gap-2 text-sm"
                      >
                        <CheckIcon className="w-4 h-4 text-success flex-shrink-0" />
                        <span>{feature}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          </div>
        ))}
      </div>

      {/* Selected Method Info */}
      {selectedMethod && (
        <div className="alert alert-info">
          <InformationCircleIcon className="w-6 h-6" />
          <div>
            <h3 className="font-semibold">
              Изабрали сте:{' '}
              {deliveryMethods.find((m) => m.id === selectedMethod)?.name}
            </h3>
            <p className="text-sm">
              Трошак доставе:{' '}
              <span className="font-medium">
                {deliveryMethods.find((m) => m.id === selectedMethod)?.price}
              </span>
              {' • '}
              Време доставе:{' '}
              <span className="font-medium">
                {deliveryMethods.find((m) => m.id === selectedMethod)?.time}
              </span>
            </p>
          </div>
        </div>
      )}

      {/* Serbian Specific Info */}
      <div className="bg-gradient-to-r from-blue-50 to-red-50 p-4 rounded-lg">
        <h4 className="font-semibold text-blue-700 mb-2">
          🇷🇸 Специфичности српског тржишта
        </h4>
        <div className="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
          <div>
            <strong>Плаћање пошљом:</strong> 95% купаца бира ову опцију
          </div>
          <div>
            <strong>Најпопуларнији градови:</strong> Београд, Нови Сад, Ниш
          </div>
          <div>
            <strong>Радно време:</strong> Курирске службе раде од 8-18h
          </div>
          <div>
            <strong>Викенд достава:</strong> Само AKS доставља суботом
          </div>
        </div>
      </div>
    </div>
  );
}
