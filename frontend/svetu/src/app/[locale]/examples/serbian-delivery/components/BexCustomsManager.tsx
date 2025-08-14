'use client';

import React, { useState } from 'react';
import {
  GlobeAltIcon,
  DocumentTextIcon,
  CurrencyDollarIcon,
  TruckIcon,
  PlusIcon,
  TrashIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';

interface CustomsItem {
  ordinalNumber: number;
  description: string;
  skuCode: string;
  hsCode: string;
  quantity: number;
  valuePerItem: number;
  weightPerItem: number;
  originCountryCode: string;
}

interface CustomsData {
  originSenderName: string;
  originSenderAddress: string;
  originSenderPlace: string;
  originSenderCountryCode: string;
  currencyCode: string;
  DDP: boolean;
  items: CustomsItem[];
}

export default function BexCustomsManager() {
  const [customsData, setCustomsData] = useState<CustomsData>({
    originSenderName: 'Sunyou Cross-Border Logistics Co.',
    originSenderAddress: 'Rm 21-N18 Unit A 11/F Tin Wui Ind Bldg',
    originSenderPlace: 'Hong Kong',
    originSenderCountryCode: 'CN',
    currencyCode: 'USD',
    DDP: true,
    items: [
      {
        ordinalNumber: 1,
        description: 'Electronic components',
        skuCode: 'SKU1234',
        hsCode: '8542.31',
        quantity: 2,
        valuePerItem: 25.99,
        weightPerItem: 0.15,
        originCountryCode: 'CN',
      },
    ],
  });

  const [showSuccess, setShowSuccess] = useState(false);

  const countries = [
    { code: 'CN', name: 'Кина', flag: '🇨🇳' },
    { code: 'DE', name: 'Немачка', flag: '🇩🇪' },
    { code: 'US', name: 'САД', flag: '🇺🇸' },
    { code: 'GB', name: 'УК', flag: '🇬🇧' },
    { code: 'IT', name: 'Италија', flag: '🇮🇹' },
    { code: 'FR', name: 'Француска', flag: '🇫🇷' },
    { code: 'JP', name: 'Јапан', flag: '🇯🇵' },
    { code: 'KR', name: 'Јужна Кореја', flag: '🇰🇷' },
  ];

  const currencies = [
    { code: 'USD', symbol: '$', name: 'US Dollar' },
    { code: 'EUR', symbol: '€', name: 'Euro' },
    { code: 'GBP', symbol: '£', name: 'British Pound' },
    { code: 'CNY', symbol: '¥', name: 'Chinese Yuan' },
    { code: 'JPY', symbol: '¥', name: 'Japanese Yen' },
  ];

  const hsCategories = [
    { code: '8542', name: 'Електронски интегрисани кругови' },
    { code: '6109', name: 'Мајице и сличне одеће' },
    { code: '9503', name: 'Играчке' },
    { code: '8517', name: 'Телефони и телекомуникације' },
    { code: '3304', name: 'Козметика и парфеми' },
    { code: '6403', name: 'Обућа' },
    { code: '4202', name: 'Торбе и кофери' },
    { code: '9102', name: 'Сатови' },
  ];

  const addItem = () => {
    const newItem: CustomsItem = {
      ordinalNumber: customsData.items.length + 1,
      description: '',
      skuCode: '',
      hsCode: '',
      quantity: 1,
      valuePerItem: 0,
      weightPerItem: 0,
      originCountryCode: 'CN',
    };
    setCustomsData({
      ...customsData,
      items: [...customsData.items, newItem],
    });
  };

  const removeItem = (index: number) => {
    const newItems = customsData.items.filter((_, i) => i !== index);
    // Reorder ordinal numbers
    newItems.forEach((item, i) => {
      item.ordinalNumber = i + 1;
    });
    setCustomsData({
      ...customsData,
      items: newItems,
    });
  };

  const updateItem = (index: number, field: keyof CustomsItem, value: any) => {
    const newItems = [...customsData.items];
    newItems[index] = {
      ...newItems[index],
      [field]: value,
    };
    setCustomsData({
      ...customsData,
      items: newItems,
    });
  };

  const calculateTotalValue = () => {
    return customsData.items.reduce((total, item) => {
      return total + item.quantity * item.valuePerItem;
    }, 0);
  };

  const calculateTotalWeight = () => {
    return customsData.items.reduce((total, item) => {
      return total + item.quantity * item.weightPerItem;
    }, 0);
  };

  const handleSubmit = () => {
    setShowSuccess(true);
    setTimeout(() => setShowSuccess(false), 3000);
  };

  return (
    <div className="space-y-6">
      {/* Header Stats */}
      <div className="stats stats-horizontal shadow w-full">
        <div className="stat">
          <div className="stat-figure text-primary">
            <CurrencyDollarIcon className="w-8 h-8" />
          </div>
          <div className="stat-title">Укупна вредност</div>
          <div className="stat-value text-primary">
            {
              currencies.find((c) => c.code === customsData.currencyCode)
                ?.symbol
            }
            {calculateTotalValue().toFixed(2)}
          </div>
          <div className="stat-desc">{customsData.currencyCode}</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-secondary">
            <TruckIcon className="w-8 h-8" />
          </div>
          <div className="stat-title">Укупна тежина</div>
          <div className="stat-value text-secondary">
            {calculateTotalWeight().toFixed(2)} kg
          </div>
          <div className="stat-desc">{customsData.items.length} артикал(а)</div>
        </div>

        <div className="stat">
          <div className="stat-figure text-accent">
            {customsData.DDP ? (
              <CheckCircleIcon className="w-8 h-8" />
            ) : (
              <ExclamationTriangleIcon className="w-8 h-8" />
            )}
          </div>
          <div className="stat-title">DDP статус</div>
          <div className="stat-value text-accent">
            {customsData.DDP ? 'Плаћено' : 'Неплаћено'}
          </div>
          <div className="stat-desc">
            {customsData.DDP ? 'Пошиљалац плаћа' : 'Прималац плаћа'}
          </div>
        </div>
      </div>

      {/* Origin Sender Information */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">
            <GlobeAltIcon className="w-6 h-6" />
            Информације о оригиналном пошиљаоцу
          </h3>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">Назив компаније</span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                value={customsData.originSenderName}
                onChange={(e) =>
                  setCustomsData({
                    ...customsData,
                    originSenderName: e.target.value,
                  })
                }
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Земља порекла</span>
              </label>
              <select
                className="select select-bordered"
                value={customsData.originSenderCountryCode}
                onChange={(e) =>
                  setCustomsData({
                    ...customsData,
                    originSenderCountryCode: e.target.value,
                  })
                }
              >
                {countries.map((country) => (
                  <option key={country.code} value={country.code}>
                    {country.flag} {country.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Адреса</span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                value={customsData.originSenderAddress}
                onChange={(e) =>
                  setCustomsData({
                    ...customsData,
                    originSenderAddress: e.target.value,
                  })
                }
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Место</span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                value={customsData.originSenderPlace}
                onChange={(e) =>
                  setCustomsData({
                    ...customsData,
                    originSenderPlace: e.target.value,
                  })
                }
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Валута</span>
              </label>
              <select
                className="select select-bordered"
                value={customsData.currencyCode}
                onChange={(e) =>
                  setCustomsData({
                    ...customsData,
                    currencyCode: e.target.value,
                  })
                }
              >
                {currencies.map((currency) => (
                  <option key={currency.code} value={currency.code}>
                    {currency.symbol} {currency.code} - {currency.name}
                  </option>
                ))}
              </select>
            </div>

            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">DDP - Delivered Duty Paid</span>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={customsData.DDP}
                  onChange={(e) =>
                    setCustomsData({
                      ...customsData,
                      DDP: e.target.checked,
                    })
                  }
                />
              </label>
              <span className="label-text-alt">
                {customsData.DDP
                  ? 'Пошиљалац плаћа царинске дажбине'
                  : 'Прималац плаћа царинске дажбине'}
              </span>
            </div>
          </div>
        </div>
      </div>

      {/* Items List */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="flex justify-between items-center mb-4">
            <h3 className="card-title">
              <DocumentTextIcon className="w-6 h-6" />
              Садржај пошиљке
            </h3>
            <button onClick={addItem} className="btn btn-primary btn-sm">
              <PlusIcon className="w-4 h-4" />
              Додај артикал
            </button>
          </div>

          <div className="space-y-4">
            {customsData.items.map((item, index) => (
              <div key={index} className="card bg-base-200">
                <div className="card-body">
                  <div className="flex justify-between items-start">
                    <h4 className="font-bold">Артикал #{item.ordinalNumber}</h4>
                    {customsData.items.length > 1 && (
                      <button
                        onClick={() => removeItem(index)}
                        className="btn btn-ghost btn-sm btn-square"
                      >
                        <TrashIcon className="w-4 h-4 text-error" />
                      </button>
                    )}
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Опис производа</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered input-sm"
                        value={item.description}
                        onChange={(e) =>
                          updateItem(index, 'description', e.target.value)
                        }
                        placeholder="нпр. Electronic components"
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">SKU код</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered input-sm"
                        value={item.skuCode}
                        onChange={(e) =>
                          updateItem(index, 'skuCode', e.target.value)
                        }
                        placeholder="SKU1234"
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">HS код</span>
                      </label>
                      <select
                        className="select select-bordered select-sm"
                        value={item.hsCode}
                        onChange={(e) =>
                          updateItem(index, 'hsCode', e.target.value)
                        }
                      >
                        <option value="">Изаберите HS код</option>
                        {hsCategories.map((hs) => (
                          <option key={hs.code} value={hs.code}>
                            {hs.code} - {hs.name}
                          </option>
                        ))}
                      </select>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Количина</span>
                      </label>
                      <input
                        type="number"
                        className="input input-bordered input-sm"
                        min="1"
                        value={item.quantity}
                        onChange={(e) =>
                          updateItem(index, 'quantity', Number(e.target.value))
                        }
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          Цена по комаду ({customsData.currencyCode})
                        </span>
                      </label>
                      <input
                        type="number"
                        className="input input-bordered input-sm"
                        step="0.01"
                        min="0"
                        value={item.valuePerItem}
                        onChange={(e) =>
                          updateItem(
                            index,
                            'valuePerItem',
                            Number(e.target.value)
                          )
                        }
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          Тежина по комаду (kg)
                        </span>
                      </label>
                      <input
                        type="number"
                        className="input input-bordered input-sm"
                        step="0.01"
                        min="0"
                        value={item.weightPerItem}
                        onChange={(e) =>
                          updateItem(
                            index,
                            'weightPerItem',
                            Number(e.target.value)
                          )
                        }
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Земља порекла</span>
                      </label>
                      <select
                        className="select select-bordered select-sm"
                        value={item.originCountryCode}
                        onChange={(e) =>
                          updateItem(index, 'originCountryCode', e.target.value)
                        }
                      >
                        {countries.map((country) => (
                          <option key={country.code} value={country.code}>
                            {country.flag} {country.code}
                          </option>
                        ))}
                      </select>
                    </div>

                    <div className="col-span-2 flex items-end">
                      <div className="stats shadow-sm">
                        <div className="stat py-2 px-4">
                          <div className="stat-title text-xs">Укупно</div>
                          <div className="stat-value text-lg">
                            {
                              currencies.find(
                                (c) => c.code === customsData.currencyCode
                              )?.symbol
                            }
                            {(item.quantity * item.valuePerItem).toFixed(2)}
                          </div>
                          <div className="stat-desc text-xs">
                            {(item.quantity * item.weightPerItem).toFixed(2)} kg
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* Summary and Actions */}
      <div className="card bg-gradient-to-r from-blue-50 to-red-50">
        <div className="card-body">
          <h3 className="card-title mb-4">Резиме царинске декларације</h3>

          <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
            <div>
              <p className="text-sm text-base-content/60">Укупна вредност</p>
              <p className="text-xl font-bold">
                {
                  currencies.find((c) => c.code === customsData.currencyCode)
                    ?.symbol
                }
                {calculateTotalValue().toFixed(2)}
              </p>
            </div>
            <div>
              <p className="text-sm text-base-content/60">Укупна тежина</p>
              <p className="text-xl font-bold">
                {calculateTotalWeight().toFixed(2)} kg
              </p>
            </div>
            <div>
              <p className="text-sm text-base-content/60">Број артикала</p>
              <p className="text-xl font-bold">{customsData.items.length}</p>
            </div>
            <div>
              <p className="text-sm text-base-content/60">DDP статус</p>
              <p className="text-xl font-bold">
                {customsData.DDP ? '✅ Плаћено' : '❌ Неплаћено'}
              </p>
            </div>
          </div>

          {showSuccess && (
            <div className="alert alert-success mb-4">
              <CheckCircleIcon className="w-6 h-6" />
              <div>
                <h3 className="font-bold">
                  Царинска декларација успешно послата!
                </h3>
                <div className="text-xs">Број пошиљке: 245815264</div>
              </div>
            </div>
          )}

          <button
            onClick={handleSubmit}
            className="btn btn-primary btn-lg btn-block"
          >
            <GlobeAltIcon className="w-5 h-5" />
            Пошаљи царинску декларацију
          </button>
        </div>
      </div>

      {/* API Request Preview */}
      <div className="collapse collapse-arrow bg-base-200">
        <input type="checkbox" />
        <div className="collapse-title text-xl font-medium">
          📡 API Request Preview
        </div>
        <div className="collapse-content">
          <div className="mockup-code">
            <pre data-prefix="$">
              <code className="text-xs">
                POST https://api.bex.rs:62502/postShipmentsCustoms
              </code>
            </pre>
            <pre data-prefix=">">
              <code className="text-xs text-warning">
                Content-Type: application/json
              </code>
            </pre>
            <pre data-prefix=">">
              <code className="text-xs text-info">
                X-AUTH-TOKEN: your-api-token
              </code>
            </pre>
            <pre>
              <code className="text-xs">
                {JSON.stringify(
                  {
                    parcels: [
                      {
                        No: 1,
                        OperatorNumber: 'LP000164786PP',
                        items: customsData.items,
                      },
                    ],
                    customs: [
                      {
                        originSenderName: customsData.originSenderName,
                        originSenderAddress: customsData.originSenderAddress,
                        originSenderPlace: customsData.originSenderPlace,
                        originSenderCountryCode:
                          customsData.originSenderCountryCode,
                        currencyCode: customsData.currencyCode,
                        DDP: customsData.DDP,
                      },
                    ],
                  },
                  null,
                  2
                )}
              </code>
            </pre>
          </div>
        </div>
      </div>
    </div>
  );
}
