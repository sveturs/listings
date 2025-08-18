'use client';

import { useState, useEffect } from 'react';
import {
  CalculatorIcon,
  ScaleIcon,
  CurrencyDollarIcon,
  MapPinIcon,
  TruckIcon,
  // BuildingStorefrontIcon,
  InformationCircleIcon,
  CheckIcon,
  // BanknotesIcon,
  ShieldCheckIcon,
  ClockIcon,
} from '@heroicons/react/24/outline';
// import { useTranslations } from 'next-intl';

interface RateCalculation {
  base_price: number;
  cod_fee: number;
  insurance_fee: number;
  fuel_surcharge: number;
  total_price: number;
  estimated_days: number;
  services: {
    name: string;
    price: number;
    included: boolean;
  }[];
}

interface Props {
  onRateCalculated?: (rate: RateCalculation) => void;
  initialParams?: {
    weight?: number;
    declaredValue?: number;
    codAmount?: number;
    senderCity?: string;
    recipientCity?: string;
  };
  className?: string;
}

export default function PostExpressRateCalculator({
  onRateCalculated,
  initialParams,
  className = '',
}: Props) {
  // const t = useTranslations('delivery');
  const [loading, setLoading] = useState(false);
  const [rate, setRate] = useState<RateCalculation | null>(null);
  const [error, setError] = useState<string | null>(null);

  const [params, setParams] = useState({
    weight: initialParams?.weight || 1,
    declared_value: initialParams?.declaredValue || 1000,
    cod_amount: initialParams?.codAmount || 0,
    sender_city: initialParams?.senderCity || 'Нови Сад',
    sender_postal_code: '21000',
    recipient_city: initialParams?.recipientCity || '',
    recipient_postal_code: '',
    express_delivery: false,
    insurance_extra: false,
    signature_required: false,
  });

  const [deliveryMethods] = useState([
    {
      id: 'standard',
      name: 'Стандартная доставка',
      description: '1-2 рабочих дня',
      icon: TruckIcon,
      multiplier: 1,
    },
    {
      id: 'express',
      name: 'Экспресс доставка',
      description: 'До 19:00 следующего дня',
      icon: ClockIcon,
      multiplier: 1.5,
    },
  ]);

  const [selectedMethod, setSelectedMethod] = useState('standard');

  useEffect(() => {
    if (initialParams) {
      setParams((prev) => ({ ...prev, ...initialParams }));
    }
  }, [initialParams]);

  // Автоматический пересчет при изменении параметров
  useEffect(() => {
    if (params.recipient_city && params.weight > 0) {
      const timer = setTimeout(() => {
        calculateRate();
      }, 500);

      return () => clearTimeout(timer);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [params, selectedMethod]);

  const calculateRate = async () => {
    if (!params.recipient_city || params.weight <= 0) return;

    setLoading(true);
    setError(null);

    try {
      const requestData = {
        ...params,
        express_delivery: selectedMethod === 'express',
        recipient_postal_code: params.recipient_postal_code || '11000', // default Belgrade
      };

      const response = await fetch('/api/v1/postexpress/calculate-rate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(requestData),
      });

      const data = await response.json();

      if (data.success) {
        const calculatedRate = data.data;
        setRate(calculatedRate);
        onRateCalculated?.(calculatedRate);
      } else {
        setError(data.message || 'Ошибка расчета стоимости');
      }
    } catch (err) {
      setError('Ошибка при расчете стоимости доставки');
      console.error('Rate calculation error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleParamChange = (field: string, value: any) => {
    setParams((prev) => ({ ...prev, [field]: value }));
  };

  const getSelectedMethodInfo = () => {
    return deliveryMethods.find((m) => m.id === selectedMethod);
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Заголовок */}
      <div className="text-center">
        <h3 className="text-xl font-bold mb-2 flex items-center justify-center gap-2">
          <CalculatorIcon className="w-6 h-6 text-primary" />
          Калькулятор стоимости доставки
        </h3>
        <p className="text-base-content/70">
          Рассчитайте точную стоимость доставки Post Express
        </p>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Параметры отправления */}
        <div className="space-y-6">
          {/* Основные параметры */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body p-6">
              <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
                <ScaleIcon className="w-5 h-5 text-primary" />
                Параметры посылки
              </h4>

              <div className="space-y-4">
                {/* Вес */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">Вес (кг) *</span>
                  </label>
                  <input
                    type="number"
                    min="0.1"
                    max="30"
                    step="0.1"
                    className="input input-bordered focus:input-primary"
                    value={params.weight}
                    onChange={(e) =>
                      handleParamChange(
                        'weight',
                        parseFloat(e.target.value) || 0
                      )
                    }
                  />
                  <label className="label">
                    <span className="label-text-alt">Максимум 30 кг</span>
                  </label>
                </div>

                {/* Объявленная ценность */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      Объявленная ценность (RSD)
                    </span>
                  </label>
                  <input
                    type="number"
                    min="0"
                    max="100000"
                    className="input input-bordered focus:input-primary"
                    value={params.declared_value}
                    onChange={(e) =>
                      handleParamChange(
                        'declared_value',
                        parseInt(e.target.value) || 0
                      )
                    }
                  />
                  <label className="label">
                    <span className="label-text-alt">
                      Для страхования посылки
                    </span>
                  </label>
                </div>

                {/* Наложенный платеж */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      Наложенный платеж (RSD)
                    </span>
                  </label>
                  <input
                    type="number"
                    min="0"
                    className="input input-bordered focus:input-primary"
                    value={params.cod_amount}
                    onChange={(e) =>
                      handleParamChange(
                        'cod_amount',
                        parseInt(e.target.value) || 0
                      )
                    }
                  />
                  <label className="label">
                    <span className="label-text-alt">
                      0 если оплата не при получении
                    </span>
                  </label>
                </div>
              </div>
            </div>
          </div>

          {/* Маршрут */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body p-6">
              <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
                <MapPinIcon className="w-5 h-5 text-primary" />
                Маршрут доставки
              </h4>

              <div className="space-y-4">
                {/* Город отправления */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      Город отправления
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered bg-base-200 cursor-not-allowed"
                    value={params.sender_city}
                    disabled
                  />
                  <label className="label">
                    <span className="label-text-alt">
                      Склад Sve Tu в Нови Саде
                    </span>
                  </label>
                </div>

                {/* Город получения */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      Город получения *
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered focus:input-primary"
                    placeholder="Белград, Суботица, Ниш..."
                    value={params.recipient_city}
                    onChange={(e) =>
                      handleParamChange('recipient_city', e.target.value)
                    }
                  />
                </div>

                {/* Почтовый индекс */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-medium">
                      Почтовый индекс получателя
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered focus:input-primary"
                    placeholder="11000"
                    value={params.recipient_postal_code}
                    onChange={(e) =>
                      handleParamChange('recipient_postal_code', e.target.value)
                    }
                  />
                  <label className="label">
                    <span className="label-text-alt">
                      Оставьте пустым для автоопределения
                    </span>
                  </label>
                </div>
              </div>
            </div>
          </div>

          {/* Способ доставки */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body p-6">
              <h4 className="font-semibold text-lg mb-4">Способ доставки</h4>

              <div className="space-y-3">
                {deliveryMethods.map((method) => (
                  <div
                    key={method.id}
                    className={`
                      card cursor-pointer transition-all border-2
                      ${
                        selectedMethod === method.id
                          ? 'border-primary bg-primary/5'
                          : 'border-transparent hover:border-primary/30'
                      }
                    `}
                    onClick={() => setSelectedMethod(method.id)}
                  >
                    <div className="card-body p-4">
                      <div className="flex items-center gap-3">
                        <method.icon className="w-6 h-6 text-primary" />
                        <div className="flex-1">
                          <div className="font-medium">{method.name}</div>
                          <div className="text-sm text-base-content/60">
                            {method.description}
                          </div>
                        </div>
                        {selectedMethod === method.id && (
                          <CheckIcon className="w-5 h-5 text-success" />
                        )}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </div>
          </div>
        </div>

        {/* Результат расчета */}
        <div className="space-y-6">
          {loading && (
            <div className="card bg-base-100 shadow-lg">
              <div className="card-body p-6 text-center">
                <span className="loading loading-spinner loading-lg"></span>
                <p className="mt-4">Расчет стоимости...</p>
              </div>
            </div>
          )}

          {error && (
            <div className="alert alert-error">
              <InformationCircleIcon className="w-5 h-5" />
              <span>{error}</span>
            </div>
          )}

          {rate && !loading && (
            <>
              {/* Основная стоимость */}
              <div className="card bg-gradient-to-r from-primary/5 to-secondary/5 shadow-lg">
                <div className="card-body p-6">
                  <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
                    <CurrencyDollarIcon className="w-5 h-5 text-primary" />
                    Стоимость доставки
                  </h4>

                  <div className="text-center mb-6">
                    <div className="text-4xl font-bold text-primary mb-2">
                      {rate.total_price.toFixed(0)} RSD
                    </div>
                    <div className="text-base-content/70">
                      {getSelectedMethodInfo()?.name} • {rate.estimated_days}{' '}
                      рабочих дня
                    </div>
                  </div>

                  {/* Детализация */}
                  <div className="space-y-3">
                    <div className="flex justify-between items-center">
                      <span>Базовая стоимость:</span>
                      <span className="font-medium">
                        {rate.base_price.toFixed(0)} RSD
                      </span>
                    </div>

                    {rate.fuel_surcharge > 0 && (
                      <div className="flex justify-between items-center">
                        <span>Топливная надбавка:</span>
                        <span className="font-medium">
                          {rate.fuel_surcharge.toFixed(0)} RSD
                        </span>
                      </div>
                    )}

                    {rate.cod_fee > 0 && (
                      <div className="flex justify-between items-center">
                        <span>Комиссия за наложенный платеж:</span>
                        <span className="font-medium">
                          {rate.cod_fee.toFixed(0)} RSD
                        </span>
                      </div>
                    )}

                    {rate.insurance_fee > 0 && (
                      <div className="flex justify-between items-center">
                        <span>Дополнительное страхование:</span>
                        <span className="font-medium">
                          {rate.insurance_fee.toFixed(0)} RSD
                        </span>
                      </div>
                    )}

                    <div className="divider my-2"></div>

                    <div className="flex justify-between items-center text-lg font-bold">
                      <span>Итого:</span>
                      <span className="text-primary">
                        {rate.total_price.toFixed(0)} RSD
                      </span>
                    </div>
                  </div>
                </div>
              </div>

              {/* Включенные услуги */}
              {rate.services && rate.services.length > 0 && (
                <div className="card bg-base-100 shadow-lg">
                  <div className="card-body p-6">
                    <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
                      <ShieldCheckIcon className="w-5 h-5 text-primary" />
                      Включенные услуги
                    </h4>

                    <div className="space-y-3">
                      {rate.services.map((service, index) => (
                        <div
                          key={index}
                          className="flex items-center justify-between"
                        >
                          <div className="flex items-center gap-2">
                            <CheckIcon className="w-5 h-5 text-success" />
                            <span>{service.name}</span>
                          </div>
                          <span
                            className={
                              service.included
                                ? 'text-success'
                                : 'text-base-content'
                            }
                          >
                            {service.included
                              ? 'Включено'
                              : `+${service.price} RSD`}
                          </span>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )}

              {/* Дополнительная информация */}
              <div className="card bg-gradient-to-r from-info/5 to-info/10">
                <div className="card-body p-6">
                  <h4 className="font-semibold text-lg mb-4 flex items-center gap-2">
                    <InformationCircleIcon className="w-5 h-5 text-info" />
                    Условия доставки
                  </h4>

                  <div className="space-y-2 text-sm">
                    <div>
                      • Страхование до 15,000 RSD включено в базовую стоимость
                    </div>
                    <div>• SMS уведомления о статусе доставки</div>
                    <div>• Хранение в отделении до 5 рабочих дней</div>
                    <div>• Доставка в рабочие дни с 08:00 до 17:00</div>
                    {params.cod_amount > 0 && (
                      <div className="text-warning">
                        • Оплата наличными или картой при получении
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </>
          )}

          {/* Пустое состояние */}
          {!rate && !loading && !error && (
            <div className="card bg-base-100 shadow-lg">
              <div className="card-body p-6 text-center">
                <CalculatorIcon className="w-16 h-16 mx-auto text-base-content/30 mb-4" />
                <h3 className="text-lg font-semibold mb-2">
                  Калькулятор доставки
                </h3>
                <p className="text-base-content/60">
                  Заполните параметры посылки для расчета стоимости доставки
                </p>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
