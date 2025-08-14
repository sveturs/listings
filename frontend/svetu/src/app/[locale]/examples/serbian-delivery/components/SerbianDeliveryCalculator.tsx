'use client';

import { useState } from 'react';
import { CalculatorIcon } from '@heroicons/react/24/outline';

export default function SerbianDeliveryCalculator() {
  const [from, setFrom] = useState('Београд');
  const [to, setTo] = useState('Нови Сад');
  const [weight, setWeight] = useState('0.5');
  const [cod, setCod] = useState('2000');

  const cities = [
    'Београд',
    'Нови Сад',
    'Ниш',
    'Крагујевац',
    'Суботица',
    'Зрењанин',
    'Панчево',
  ];

  const couriers = [
    {
      name: 'BexExpress',
      basePrice: 170,
      perKg: 45,
      codRate: 0.018,
      color: 'bg-purple-100',
      textColor: 'text-purple-600',
    },
    {
      name: 'Post Express',
      basePrice: 150,
      perKg: 40,
      codRate: 0.025,
      color: 'bg-blue-100',
      textColor: 'text-blue-600',
    },
    {
      name: 'City Express',
      basePrice: 180,
      perKg: 45,
      codRate: 0.02,
      color: 'bg-green-100',
      textColor: 'text-green-600',
    },
    {
      name: 'Yettel Post',
      basePrice: 120,
      perKg: 35,
      codRate: 0.03,
      color: 'bg-orange-100',
      textColor: 'text-orange-600',
    },
  ];

  const calculatePrice = (courier: any) => {
    const basePrice = courier.basePrice;
    const weightPrice =
      parseFloat(weight) > 1 ? (parseFloat(weight) - 1) * courier.perKg : 0;
    const codPrice = parseFloat(cod) * courier.codRate;
    return Math.round(basePrice + weightPrice + codPrice);
  };

  return (
    <div className="max-w-4xl mx-auto space-y-6">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title text-2xl mb-4">
            <CalculatorIcon className="w-8 h-8 text-primary" />
            Калкулатор трошкова доставе
          </h3>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">Од</span>
              </label>
              <select
                className="select select-bordered"
                value={from}
                onChange={(e) => setFrom(e.target.value)}
              >
                {cities.map((city) => (
                  <option key={city} value={city}>
                    {city}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">До</span>
              </label>
              <select
                className="select select-bordered"
                value={to}
                onChange={(e) => setTo(e.target.value)}
              >
                {cities.map((city) => (
                  <option key={city} value={city}>
                    {city}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Тежина (кг)</span>
              </label>
              <input
                type="number"
                step="0.1"
                className="input input-bordered"
                value={weight}
                onChange={(e) => setWeight(e.target.value)}
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Поштарина (РСД)</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={cod}
                onChange={(e) => setCod(e.target.value)}
              />
            </div>
          </div>
        </div>
      </div>

      {/* Results */}
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        {couriers.map((courier) => (
          <div
            key={courier.name}
            className={`card ${courier.color} shadow-xl border-2 border-base-200`}
          >
            <div className="card-body">
              <h4 className={`card-title text-lg ${courier.textColor}`}>
                {courier.name}
              </h4>

              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Основна цена:</span>
                  <span>{courier.basePrice} РСД</span>
                </div>
                {parseFloat(weight) > 1 && (
                  <div className="flex justify-between">
                    <span>Додатна тежина:</span>
                    <span>
                      +{Math.round((parseFloat(weight) - 1) * courier.perKg)}{' '}
                      РСД
                    </span>
                  </div>
                )}
                <div className="flex justify-between">
                  <span>Провизија поштарине:</span>
                  <span>
                    +{Math.round(parseFloat(cod) * courier.codRate)} РСД
                  </span>
                </div>
                <hr />
                <div className="flex justify-between font-bold text-lg">
                  <span>Укупно:</span>
                  <span className="text-primary">
                    {calculatePrice(courier)} РСД
                  </span>
                </div>
              </div>

              <div className="mt-4 flex items-center justify-between">
                <div className="badge badge-outline badge-sm">
                  {courier.name === 'BexExpress'
                    ? '1-2 дана'
                    : courier.name === 'Post Express'
                      ? '1-3 дана'
                      : courier.name === 'City Express'
                        ? '2-4 дана'
                        : '1-2 дана'}
                </div>
                {courier.name === 'BexExpress' && (
                  <div className="badge badge-primary badge-sm">
                    API интеграција
                  </div>
                )}
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="alert alert-info">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
          className="stroke-current shrink-0 w-6 h-6"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth="2"
            d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
          ></path>
        </svg>
        <div>
          <h4 className="font-semibold">Напомене:</h4>
          <ul className="text-sm mt-1">
            <li>• Цене су оријентационе и могу да се разликују</li>
            <li>• Максимална тежина пакета: 30кг</li>
            <li>• Осигурање се наплаћује додатно</li>
          </ul>
        </div>
      </div>
    </div>
  );
}
