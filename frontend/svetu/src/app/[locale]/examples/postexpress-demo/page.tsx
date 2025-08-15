'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import {
  TruckIcon,
  MapPinIcon,
  BuildingOfficeIcon,
  CalculatorIcon,
  DocumentTextIcon,
  QrCodeIcon,
  ChartBarIcon,
  CurrencyDollarIcon,
  ClockIcon,
  CheckCircleIcon,
  XCircleIcon,
  ArrowPathIcon,
  MagnifyingGlassIcon,
} from '@heroicons/react/24/outline';

export default function PostExpressDemoPage() {
  const t = useTranslations('delivery');
  const [activeTab, setActiveTab] = useState('overview');
  const [selectedMethod, setSelectedMethod] = useState('courier');
  const [weight, setWeight] = useState(2);
  const [codAmount, setCodAmount] = useState(5000);
  const [trackingNumber, setTrackingNumber] = useState('');
  const [showCalculation, setShowCalculation] = useState(false);

  // Тарифы из коммерческого предложения
  const rates = {
    2: 340,
    5: 450,
    10: 580,
    20: 790,
  };

  const getRate = (weight: number) => {
    if (weight <= 2) return rates[2];
    if (weight <= 5) return rates[5];
    if (weight <= 10) return rates[10];
    return rates[20];
  };

  const calculateTotal = () => {
    const baseRate = getRate(weight);
    const codFee = codAmount > 0 ? 45 : 0;
    const insuranceFee = codAmount > 15000 ? (codAmount - 15000) * 0.01 : 0;
    return {
      base: baseRate,
      cod: codFee,
      insurance: insuranceFee,
      total: baseRate + codFee + insuranceFee,
    };
  };

  const tabs = [
    { id: 'overview', label: 'Преглед', icon: ChartBarIcon },
    { id: 'delivery', label: 'Достава', icon: TruckIcon },
    { id: 'calculator', label: 'Калкулатор', icon: CalculatorIcon },
    { id: 'tracking', label: 'Праћење', icon: MapPinIcon },
    { id: 'api', label: 'API статус', icon: DocumentTextIcon },
  ];

  return (
    <div className="min-h-screen bg-base-200">
      {/* Hero Section */}
      <div className="hero bg-gradient-to-r from-blue-600 to-blue-800 text-white">
        <div className="hero-content text-center py-12">
          <div className="max-w-md">
            <h1 className="text-5xl font-bold mb-4">
              Post Express интеграција
            </h1>
            <p className="text-xl">
              Комплетна интеграција са WSP API за вашу платформу
            </p>
            <div className="mt-6 flex justify-center gap-4">
              <div className="badge badge-lg badge-warning gap-2">
                <CheckCircleIcon className="w-4 h-4" />
                Production Ready
              </div>
              <div className="badge badge-lg badge-success gap-2">
                <CheckCircleIcon className="w-4 h-4" />
                API Integrated
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="tabs tabs-boxed justify-center p-4 bg-base-100">
        {tabs.map((tab) => (
          <a
            key={tab.id}
            className={`tab tab-lg gap-2 ${activeTab === tab.id ? 'tab-active' : ''}`}
            onClick={() => setActiveTab(tab.id)}
          >
            <tab.icon className="w-5 h-5" />
            {tab.label}
          </a>
        ))}
      </div>

      {/* Content */}
      <div className="container mx-auto p-6">
        {/* Overview Tab */}
        {activeTab === 'overview' && (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title text-primary">
                  <CurrencyDollarIcon className="w-6 h-6" />
                  Тарифе
                </h2>
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span>До 2кг:</span>
                    <span className="font-bold">340 РСД</span>
                  </div>
                  <div className="flex justify-between">
                    <span>2-5кг:</span>
                    <span className="font-bold">450 РСД</span>
                  </div>
                  <div className="flex justify-between">
                    <span>5-10кг:</span>
                    <span className="font-bold">580 РСД</span>
                  </div>
                  <div className="flex justify-between">
                    <span>10-20кг:</span>
                    <span className="font-bold">790 РСД</span>
                  </div>
                </div>
                <div className="divider"></div>
                <div className="text-sm text-success">✓ Без ПДВ-а</div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title text-primary">
                  <DocumentTextIcon className="w-6 h-6" />
                  Откупнина
                </h2>
                <div className="space-y-3">
                  <div className="alert alert-info">
                    <span>
                      Провизија: <strong>45 РСД</strong> фиксно
                    </span>
                  </div>
                  <div className="text-sm">
                    • Наплата од примаоца без провизије
                    <br />
                    • Уплата истог дана
                    <br />
                    • Осигурање до 15.000 РСД укључено
                    <br />• Преко 15.000: +1% на вишак
                  </div>
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title text-primary">
                  <ClockIcon className="w-6 h-6" />
                  Рокови доставе
                </h2>
                <div className="space-y-3">
                  <div className="badge badge-lg badge-success w-full justify-center">
                    Данас за сутра
                  </div>
                  <div className="text-sm space-y-2">
                    <div className="flex items-center gap-2">
                      <CheckCircleIcon className="w-4 h-4 text-success" />
                      Стандардне: 1 дан
                    </div>
                    <div className="flex items-center gap-2">
                      <ClockIcon className="w-4 h-4 text-warning" />
                      Нестандардне: 2-3 дана
                    </div>
                    <div className="flex items-center gap-2">
                      <BuildingOfficeIcon className="w-4 h-4 text-info" />
                      Чување 5 радних дана
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}

        {/* Delivery Tab */}
        {activeTab === 'delivery' && (
          <div className="space-y-6">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">Изаберите начин доставе</h2>
                <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mt-4">
                  <label
                    className={`card bordered cursor-pointer ${selectedMethod === 'courier' ? 'ring-2 ring-primary' : ''}`}
                  >
                    <input
                      type="radio"
                      name="delivery"
                      className="hidden"
                      checked={selectedMethod === 'courier'}
                      onChange={() => setSelectedMethod('courier')}
                    />
                    <div className="card-body text-center">
                      <TruckIcon className="w-12 h-12 mx-auto text-primary" />
                      <h3 className="font-bold">Курирска достава</h3>
                      <p className="text-sm">Достава на кућну адресу</p>
                    </div>
                  </label>

                  <label
                    className={`card bordered cursor-pointer ${selectedMethod === 'office' ? 'ring-2 ring-primary' : ''}`}
                  >
                    <input
                      type="radio"
                      name="delivery"
                      className="hidden"
                      checked={selectedMethod === 'office'}
                      onChange={() => setSelectedMethod('office')}
                    />
                    <div className="card-body text-center">
                      <BuildingOfficeIcon className="w-12 h-12 mx-auto text-primary" />
                      <h3 className="font-bold">Пошта</h3>
                      <p className="text-sm">Преузимање у пошти</p>
                    </div>
                  </label>

                  <label
                    className={`card bordered cursor-pointer ${selectedMethod === 'pickup' ? 'ring-2 ring-primary' : ''}`}
                  >
                    <input
                      type="radio"
                      name="delivery"
                      className="hidden"
                      checked={selectedMethod === 'pickup'}
                      onChange={() => setSelectedMethod('pickup')}
                    />
                    <div className="card-body text-center">
                      <MapPinIcon className="w-12 h-12 mx-auto text-primary" />
                      <h3 className="font-bold">Pickup пункт</h3>
                      <p className="text-sm">500+ локација</p>
                    </div>
                  </label>
                </div>

                {selectedMethod === 'courier' && (
                  <div className="mt-6 p-4 bg-base-200 rounded-lg">
                    <h3 className="font-bold mb-3">Адреса доставе</h3>
                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <input
                        type="text"
                        placeholder="Име и презиме"
                        className="input input-bordered"
                      />
                      <input
                        type="text"
                        placeholder="Телефон"
                        className="input input-bordered"
                      />
                      <input
                        type="text"
                        placeholder="Улица и број"
                        className="input input-bordered"
                      />
                      <input
                        type="text"
                        placeholder="Град"
                        className="input input-bordered"
                      />
                      <input
                        type="text"
                        placeholder="Поштански број"
                        className="input input-bordered"
                      />
                      <input
                        type="text"
                        placeholder="Напомена"
                        className="input input-bordered"
                      />
                    </div>
                  </div>
                )}

                {selectedMethod === 'office' && (
                  <div className="mt-6 p-4 bg-base-200 rounded-lg">
                    <h3 className="font-bold mb-3">Изаберите пошту</h3>
                    <select className="select select-bordered w-full">
                      <option>21101 Нови Сад 1 - Народних хероја 2</option>
                      <option>21102 Нови Сад 2 - Булевар ослобођења 100</option>
                      <option>21103 Нови Сад 3 - Футошка 14</option>
                      <option>11000 Београд 1 - Таковска 2</option>
                      <option>11010 Београд 6 - Савска 2</option>
                    </select>
                    <div className="alert alert-info mt-4">
                      <span>Пошиљка ће бити доступна 5 радних дана</span>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Calculator Tab */}
        {activeTab === 'calculator' && (
          <div className="max-w-2xl mx-auto">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">Калкулатор цене доставе</h2>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Тежина пакета (кг)</span>
                  </label>
                  <input
                    type="range"
                    min="0.5"
                    max="20"
                    step="0.5"
                    value={weight}
                    onChange={(e) => setWeight(Number(e.target.value))}
                    className="range range-primary"
                  />
                  <div className="w-full flex justify-between text-xs px-2">
                    <span>0.5кг</span>
                    <span className="font-bold text-lg">{weight}кг</span>
                    <span>20кг</span>
                  </div>
                </div>

                <div className="form-control mt-4">
                  <label className="label">
                    <span className="label-text">Откупнина (РСД)</span>
                  </label>
                  <input
                    type="number"
                    value={codAmount}
                    onChange={(e) => setCodAmount(Number(e.target.value))}
                    className="input input-bordered"
                    placeholder="Износ откупнине"
                  />
                </div>

                <button
                  className="btn btn-primary mt-6"
                  onClick={() => setShowCalculation(true)}
                >
                  <CalculatorIcon className="w-5 h-5" />
                  Израчунај цену
                </button>

                {showCalculation && (
                  <div className="mt-6 p-4 bg-base-200 rounded-lg">
                    <h3 className="font-bold mb-3">Детаљи калкулације</h3>
                    <div className="space-y-2">
                      <div className="flex justify-between">
                        <span>Основна цена ({weight}кг):</span>
                        <span className="font-bold">
                          {calculateTotal().base} РСД
                        </span>
                      </div>
                      {codAmount > 0 && (
                        <div className="flex justify-between">
                          <span>Провизија откупнине:</span>
                          <span className="font-bold">
                            {calculateTotal().cod} РСД
                          </span>
                        </div>
                      )}
                      {calculateTotal().insurance > 0 && (
                        <div className="flex justify-between">
                          <span>Додатно осигурање:</span>
                          <span className="font-bold">
                            {calculateTotal().insurance.toFixed(2)} РСД
                          </span>
                        </div>
                      )}
                      <div className="divider"></div>
                      <div className="flex justify-between text-lg font-bold text-primary">
                        <span>УКУПНО:</span>
                        <span>{calculateTotal().total.toFixed(2)} РСД</span>
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {/* Tracking Tab */}
        {activeTab === 'tracking' && (
          <div className="max-w-2xl mx-auto">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">Праћење пошиљке</h2>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Унесите број за праћење</span>
                  </label>
                  <div className="input-group">
                    <input
                      type="text"
                      placeholder="RE123456789RS"
                      className="input input-bordered w-full"
                      value={trackingNumber}
                      onChange={(e) => setTrackingNumber(e.target.value)}
                    />
                    <button className="btn btn-primary">
                      <MagnifyingGlassIcon className="w-5 h-5" />
                      Пронађи
                    </button>
                  </div>
                </div>

                {trackingNumber && (
                  <div className="mt-6">
                    <ul className="steps steps-vertical">
                      <li className="step step-primary">
                        <div className="text-left ml-4">
                          <div className="font-bold">Пошиљка регистрована</div>
                          <div className="text-sm text-base-content/70">
                            15.08.2025 09:00 - Нови Сад
                          </div>
                        </div>
                      </li>
                      <li className="step step-primary">
                        <div className="text-left ml-4">
                          <div className="font-bold">Преузето од пошиљаоца</div>
                          <div className="text-sm text-base-content/70">
                            15.08.2025 14:30 - Нови Сад
                          </div>
                        </div>
                      </li>
                      <li className="step">
                        <div className="text-left ml-4">
                          <div className="font-bold">У транзиту</div>
                          <div className="text-sm text-base-content/70">
                            Очекивано време
                          </div>
                        </div>
                      </li>
                      <li className="step">
                        <div className="text-left ml-4">
                          <div className="font-bold">Достављено</div>
                          <div className="text-sm text-base-content/70">
                            Очекивано сутра
                          </div>
                        </div>
                      </li>
                    </ul>
                  </div>
                )}
              </div>
            </div>
          </div>
        )}

        {/* API Status Tab */}
        {activeTab === 'api' && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">
                  Имплементиране WSP API транзакције
                </h2>
                <div className="space-y-2">
                  {[
                    {
                      id: 3,
                      name: 'GetNaselje',
                      desc: 'Претрага насеља',
                      status: true,
                    },
                    {
                      id: 10,
                      name: 'GetPostanskeJedinice',
                      desc: 'Листа пошта',
                      status: true,
                    },
                    {
                      id: 15,
                      name: 'PracenjePosiljke',
                      desc: 'Праћење',
                      status: true,
                    },
                    {
                      id: 20,
                      name: 'StampaNalepnice',
                      desc: 'Штампа налепнице',
                      status: true,
                    },
                    {
                      id: 25,
                      name: 'StorniranjePosiljke',
                      desc: 'Сторнирање',
                      status: true,
                    },
                    {
                      id: 63,
                      name: 'CreatePosiljka',
                      desc: 'Креирање пошиљке',
                      status: true,
                    },
                    {
                      id: 73,
                      name: 'Manifest',
                      desc: 'Манифест',
                      status: false,
                    },
                  ].map((api) => (
                    <div
                      key={api.id}
                      className="flex items-center justify-between p-2 rounded-lg bg-base-200"
                    >
                      <div>
                        <span className="font-mono text-sm">ID {api.id}:</span>
                        <span className="ml-2 font-bold">{api.name}</span>
                        <span className="ml-2 text-sm text-base-content/70">
                          ({api.desc})
                        </span>
                      </div>
                      {api.status ? (
                        <CheckCircleIcon className="w-5 h-5 text-success" />
                      ) : (
                        <XCircleIcon className="w-5 h-5 text-error" />
                      )}
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h2 className="card-title">Статус интеграције</h2>
                <div className="space-y-4">
                  <div className="alert alert-success">
                    <CheckCircleIcon className="w-6 h-6" />
                    <span>Backend модул - 100% готов</span>
                  </div>
                  <div className="alert alert-success">
                    <CheckCircleIcon className="w-6 h-6" />
                    <span>База података - Миграције готове</span>
                  </div>
                  <div className="alert alert-success">
                    <CheckCircleIcon className="w-6 h-6" />
                    <span>Frontend компоненте - Готове</span>
                  </div>
                  <div className="alert alert-warning">
                    <ArrowPathIcon className="w-6 h-6" />
                    <span>Production credentials - Чека се</span>
                  </div>
                </div>

                <div className="divider"></div>

                <div className="card bg-base-200">
                  <div className="card-body">
                    <h3 className="font-bold">Потребно за production:</h3>
                    <ul className="text-sm space-y-1">
                      <li>• Username за WSP API</li>
                      <li>• Password за WSP API</li>
                      <li>• Потписан уговор</li>
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="bg-base-300 mt-12 p-6">
        <div className="container mx-auto text-center">
          <div className="flex justify-center gap-4 mb-4">
            <div className="badge badge-lg gap-2">
              <DocumentTextIcon className="w-4 h-4" />
              Понуда: 2025-sl од 31.07.2025
            </div>
            <div className="badge badge-lg gap-2">
              <MapPinIcon className="w-4 h-4" />
              Контакт: prodaja@posta.rs
            </div>
          </div>
          <p className="text-sm text-base-content/70">
            Post Express WSP API интеграција - Sve Tu Platform
          </p>
        </div>
      </div>
    </div>
  );
}
