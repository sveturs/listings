'use client';

import { useState } from 'react';
import {
  CalculatorIcon,
  MapPinIcon,
  CubeIcon,
  BanknotesIcon,
  TruckIcon,
  ClockIcon,
  InformationCircleIcon,
} from '@heroicons/react/24/outline';

export default function DeliveryCalculator() {
  const [isCalculating, setIsCalculating] = useState(false);
  const [showResults, setShowResults] = useState(false);
  const [formData, setFormData] = useState({
    fromCity: '',
    toCity: '',
    weight: '',
    length: '',
    width: '',
    height: '',
    value: '',
    deliveryType: 'standard',
    insurance: false,
    cod: false,
    codAmount: '',
  });

  const calculateDelivery = () => {
    setIsCalculating(true);
    setTimeout(() => {
      setIsCalculating(false);
      setShowResults(true);
    }, 1500);
  };

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>
  ) => {
    const { name, value, type } = e.target;
    if (type === 'checkbox') {
      setFormData((prev) => ({
        ...prev,
        [name]: (e.target as HTMLInputElement).checked,
      }));
    } else {
      setFormData((prev) => ({
        ...prev,
        [name]: value,
      }));
    }
  };

  const calculateVolume = () => {
    const l = parseFloat(formData.length) || 0;
    const w = parseFloat(formData.width) || 0;
    const h = parseFloat(formData.height) || 0;
    return ((l * w * h) / 1000000).toFixed(3); // Convert to m³
  };

  const calculateVolumetricWeight = () => {
    const l = parseFloat(formData.length) || 0;
    const w = parseFloat(formData.width) || 0;
    const h = parseFloat(formData.height) || 0;
    return ((l * w * h) / 5000).toFixed(2); // Standard volumetric divisor
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      {/* Calculator Form */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="flex items-center gap-3 mb-4">
            <div className="p-2 bg-primary/10 rounded-lg">
              <CalculatorIcon className="w-6 h-6 text-primary" />
            </div>
            <h3 className="card-title">Калькулятор доставки</h3>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            {/* Route */}
            <div className="space-y-4">
              <div className="text-sm font-semibold text-base-content/60">
                Маршрут
              </div>

              <div>
                <label className="label">
                  <span className="label-text">Откуда</span>
                </label>
                <select
                  name="fromCity"
                  className="select select-bordered w-full"
                  value={formData.fromCity}
                  onChange={handleInputChange}
                >
                  <option value="">Выберите город</option>
                  <option value="belgrade">Белград</option>
                  <option value="novi-sad">Нови Сад</option>
                  <option value="nis">Ниш</option>
                  <option value="kragujevac">Крагуевац</option>
                  <option value="subotica">Суботица</option>
                </select>
              </div>

              <div>
                <label className="label">
                  <span className="label-text">Куда</span>
                </label>
                <select
                  name="toCity"
                  className="select select-bordered w-full"
                  value={formData.toCity}
                  onChange={handleInputChange}
                >
                  <option value="">Выберите город</option>
                  <option value="belgrade">Белград</option>
                  <option value="novi-sad">Нови Сад</option>
                  <option value="nis">Ниш</option>
                  <option value="kragujevac">Крагуевац</option>
                  <option value="subotica">Суботица</option>
                </select>
              </div>

              <div>
                <label className="label">
                  <span className="label-text">Тип доставки</span>
                </label>
                <select
                  name="deliveryType"
                  className="select select-bordered w-full"
                  value={formData.deliveryType}
                  onChange={handleInputChange}
                >
                  <option value="standard">Стандартная (2-3 дня)</option>
                  <option value="express">Экспресс (1 день)</option>
                  <option value="same-day">В тот же день</option>
                </select>
              </div>
            </div>

            {/* Package Parameters */}
            <div className="space-y-4">
              <div className="text-sm font-semibold text-base-content/60">
                Параметры посылки
              </div>

              <div>
                <label className="label">
                  <span className="label-text">Вес (кг)</span>
                </label>
                <input
                  type="number"
                  name="weight"
                  placeholder="0.0"
                  className="input input-bordered w-full"
                  value={formData.weight}
                  onChange={handleInputChange}
                  step="0.1"
                />
              </div>

              <div>
                <label className="label">
                  <span className="label-text">Размеры (см)</span>
                </label>
                <div className="grid grid-cols-3 gap-2">
                  <input
                    type="number"
                    name="length"
                    placeholder="Длина"
                    className="input input-bordered"
                    value={formData.length}
                    onChange={handleInputChange}
                  />
                  <input
                    type="number"
                    name="width"
                    placeholder="Ширина"
                    className="input input-bordered"
                    value={formData.width}
                    onChange={handleInputChange}
                  />
                  <input
                    type="number"
                    name="height"
                    placeholder="Высота"
                    className="input input-bordered"
                    value={formData.height}
                    onChange={handleInputChange}
                  />
                </div>
                {formData.length && formData.width && formData.height && (
                  <div className="text-xs text-base-content/60 mt-2">
                    Объем: {calculateVolume()} м³ | Объемный вес:{' '}
                    {calculateVolumetricWeight()} кг
                  </div>
                )}
              </div>

              <div>
                <label className="label">
                  <span className="label-text">
                    Объявленная стоимость (RSD)
                  </span>
                </label>
                <input
                  type="number"
                  name="value"
                  placeholder="0"
                  className="input input-bordered w-full"
                  value={formData.value}
                  onChange={handleInputChange}
                />
              </div>
            </div>
          </div>

          {/* Additional Services */}
          <div className="mt-6">
            <div className="text-sm font-semibold text-base-content/60 mb-3">
              Дополнительные услуги
            </div>

            <div className="space-y-3">
              <label className="flex items-start sm:items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  name="insurance"
                  className="checkbox checkbox-primary checkbox-sm sm:checkbox-md mt-0.5 sm:mt-0"
                  checked={formData.insurance}
                  onChange={handleInputChange}
                />
                <div className="flex-1">
                  <div className="font-medium text-sm sm:text-base">
                    Страхование посылки
                  </div>
                  <div className="text-xs text-base-content/60">
                    2% от объявленной стоимости
                  </div>
                </div>
              </label>

              <label className="flex items-start sm:items-center gap-3 cursor-pointer">
                <input
                  type="checkbox"
                  name="cod"
                  className="checkbox checkbox-primary checkbox-sm sm:checkbox-md mt-0.5 sm:mt-0"
                  checked={formData.cod}
                  onChange={handleInputChange}
                />
                <div className="flex-1">
                  <div className="font-medium text-sm sm:text-base">
                    Наложенный платеж (COD)
                  </div>
                  <div className="text-xs text-base-content/60">
                    Получатель оплачивает при получении
                  </div>
                </div>
              </label>

              {formData.cod && (
                <div className="ml-7">
                  <input
                    type="number"
                    name="codAmount"
                    placeholder="Сумма наложенного платежа (RSD)"
                    className="input input-bordered w-full"
                    value={formData.codAmount}
                    onChange={handleInputChange}
                  />
                </div>
              )}
            </div>
          </div>

          {/* Calculate Button */}
          <button
            className={`btn btn-primary btn-lg w-full mt-6 ${isCalculating ? 'loading' : ''}`}
            onClick={calculateDelivery}
            disabled={isCalculating}
          >
            {!isCalculating && <CalculatorIcon className="w-5 h-5" />}
            Рассчитать стоимость
          </button>
        </div>
      </div>

      {/* Results */}
      {showResults && (
        <>
          {/* Main Result Card */}
          <div className="card bg-gradient-to-r from-primary/10 to-secondary/10">
            <div className="card-body">
              <h3 className="text-xl font-bold mb-4">Результаты расчета</h3>

              <div className="grid md:grid-cols-3 gap-6">
                <div className="text-center">
                  <div className="p-3 bg-primary/20 rounded-full inline-flex mb-2">
                    <BanknotesIcon className="w-8 h-8 text-primary" />
                  </div>
                  <div className="text-3xl font-bold text-primary">450 RSD</div>
                  <div className="text-sm text-base-content/60">
                    Стоимость доставки
                  </div>
                </div>

                <div className="text-center">
                  <div className="p-3 bg-secondary/20 rounded-full inline-flex mb-2">
                    <ClockIcon className="w-8 h-8 text-secondary" />
                  </div>
                  <div className="text-3xl font-bold">2-3 дня</div>
                  <div className="text-sm text-base-content/60">
                    Время доставки
                  </div>
                </div>

                <div className="text-center">
                  <div className="p-3 bg-accent/20 rounded-full inline-flex mb-2">
                    <MapPinIcon className="w-8 h-8 text-accent" />
                  </div>
                  <div className="text-3xl font-bold">~180 км</div>
                  <div className="text-sm text-base-content/60">Расстояние</div>
                </div>
              </div>
            </div>
          </div>

          {/* Detailed Breakdown */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h4 className="font-semibold mb-4">Детальная расшифровка</h4>

              <div className="space-y-3">
                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-2">
                    <TruckIcon className="w-5 h-5 text-base-content/40" />
                    <span>Базовая стоимость доставки</span>
                  </div>
                  <span className="font-medium">350 RSD</span>
                </div>

                <div className="flex justify-between items-center">
                  <div className="flex items-center gap-2">
                    <CubeIcon className="w-5 h-5 text-base-content/40" />
                    <span>Надбавка за вес (5 кг)</span>
                  </div>
                  <span className="font-medium">50 RSD</span>
                </div>

                {formData.insurance && (
                  <div className="flex justify-between items-center">
                    <div className="flex items-center gap-2">
                      <span className="ml-7">
                        Страхование (2% от 10,000 RSD)
                      </span>
                    </div>
                    <span className="font-medium">200 RSD</span>
                  </div>
                )}

                {formData.cod && (
                  <div className="flex justify-between items-center">
                    <div className="flex items-center gap-2">
                      <span className="ml-7">
                        Обработка наложенного платежа
                      </span>
                    </div>
                    <span className="font-medium">50 RSD</span>
                  </div>
                )}

                <div className="divider"></div>

                <div className="flex justify-between items-center">
                  <span className="font-semibold text-lg">Итого</span>
                  <span className="font-bold text-xl text-primary">
                    450 RSD
                  </span>
                </div>
              </div>

              {/* Info Alert */}
              <div className="alert alert-info mt-6">
                <InformationCircleIcon className="w-5 h-5" />
                <div>
                  <div className="font-semibold">
                    Это предварительный расчет
                  </div>
                  <div className="text-sm">
                    Окончательная стоимость может отличаться в зависимости от
                    фактического веса и объема посылки.
                  </div>
                </div>
              </div>

              {/* Alternative Options */}
              <div className="mt-6">
                <h5 className="font-semibold mb-3">Альтернативные варианты</h5>

                <div className="grid md:grid-cols-2 gap-3">
                  <div className="p-4 bg-base-200 rounded-lg">
                    <div className="flex justify-between items-start">
                      <div>
                        <div className="font-medium">Экспресс доставка</div>
                        <div className="text-sm text-base-content/60">
                          1 рабочий день
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="font-bold text-primary">750 RSD</div>
                        <div className="text-xs text-base-content/60">+67%</div>
                      </div>
                    </div>
                  </div>

                  <div className="p-4 bg-base-200 rounded-lg">
                    <div className="flex justify-between items-start">
                      <div>
                        <div className="font-medium">Самовывоз из пункта</div>
                        <div className="text-sm text-base-content/60">
                          2-3 рабочих дня
                        </div>
                      </div>
                      <div className="text-right">
                        <div className="font-bold text-success">280 RSD</div>
                        <div className="text-xs text-base-content/60">-38%</div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              {/* Action Buttons */}
              <div className="flex gap-3 mt-6">
                <button className="btn btn-primary flex-1">
                  Оформить доставку
                </button>
                <button
                  className="btn btn-ghost"
                  onClick={() => setShowResults(false)}
                >
                  Новый расчет
                </button>
              </div>
            </div>
          </div>
        </>
      )}
    </div>
  );
}
