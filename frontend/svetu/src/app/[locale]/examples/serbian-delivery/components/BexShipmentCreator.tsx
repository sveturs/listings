'use client';

import React, { useState } from 'react';
import {
  CubeIcon,
  UserIcon,
  MapPinIcon,
  DocumentCheckIcon,
  TruckIcon,
} from '@heroicons/react/24/outline';

interface ShipmentData {
  sender: {
    type: number;
    firstName: string;
    lastName: string;
    phone: string;
    place: string;
    street: string;
    houseNumber: string;
    apartment: string;
    timeFrom: string;
    timeTo: string;
  };
  receiver: {
    firstName: string;
    lastName: string;
    phone: string;
    place: string;
    street: string;
    houseNumber: string;
    apartment: string;
    preNotification: number;
    comment: string;
  };
  shipment: {
    category: number;
    weight: number;
    packages: number;
    payType: number;
    insurance: number;
    cashOnDelivery: number;
    publicComment: string;
    privateComment: string;
    personalDelivery: boolean;
    returnInvoices: boolean;
    returnConfirmation: boolean;
    returnPackage: boolean;
  };
}

export default function BexShipmentCreator() {
  const [activeStep, setActiveStep] = useState(1);
  const [showSuccess, setShowSuccess] = useState(false);
  const [shipmentData, setShipmentData] = useState<ShipmentData>({
    sender: {
      type: 1,
      firstName: 'Марко',
      lastName: 'Петровић',
      phone: '0648516928',
      place: 'БЕОГРАД',
      street: 'Кнез Михаилова',
      houseNumber: '25',
      apartment: '5',
      timeFrom: '09:00',
      timeTo: '17:00',
    },
    receiver: {
      firstName: 'Ана',
      lastName: 'Јовановић',
      phone: '0651234567',
      place: 'НОВИ САД',
      street: 'Дунавска',
      houseNumber: '10',
      apartment: '',
      preNotification: 30,
      comment: 'Позвати пре доставе',
    },
    shipment: {
      category: 1,
      weight: 0,
      packages: 1,
      payType: 6,
      insurance: 5000,
      cashOnDelivery: 2500,
      publicComment: 'Пажљиво - ломљиво',
      privateComment: 'Поруџбина #12345',
      personalDelivery: false,
      returnInvoices: false,
      returnConfirmation: false,
      returnPackage: false,
    },
  });

  const categories = [
    { id: 1, name: 'Документ', icon: '📄' },
    { id: 2, name: 'Документ у коверти', icon: '✉️' },
    { id: 3, name: 'Документ у ПВЦ', icon: '📋' },
    { id: 4, name: 'Коверат А3', icon: '📨' },
    { id: 5, name: 'Коверат А4', icon: '📧' },
    { id: 31, name: 'Пакет до 50кг', icon: '📦', needsWeight: true },
    { id: 32, name: 'Палета до 200кг', icon: '🚛', needsWeight: true },
  ];

  const payTypes = [
    { id: 1, name: 'Пошиљалац готовина', icon: '💵' },
    { id: 2, name: 'Прималац готовина', icon: '💴' },
    { id: 6, name: 'Купац преко банке', icon: '🏦' },
  ];

  const preNotifications = [
    { value: 0, label: 'Без обавештења' },
    { value: 1, label: '1 минут' },
    { value: 5, label: '5 минута' },
    { value: 15, label: '15 минута' },
    { value: 30, label: '30 минута' },
    { value: 45, label: '45 минута' },
    { value: 60, label: '1 сат' },
  ];

  const handleCreateShipment = () => {
    setShowSuccess(true);
    setTimeout(() => {
      setShowSuccess(false);
    }, 3000);
  };

  const steps = [
    { id: 1, name: 'Пошиљалац', icon: UserIcon },
    { id: 2, name: 'Прималац', icon: MapPinIcon },
    { id: 3, name: 'Детаљи', icon: CubeIcon },
    { id: 4, name: 'Потврда', icon: DocumentCheckIcon },
  ];

  return (
    <div className="space-y-6">
      {/* Progress Steps */}
      <div className="flex justify-center">
        <ul className="steps steps-horizontal">
          {steps.map((step) => (
            <li
              key={step.id}
              className={`step ${activeStep >= step.id ? 'step-primary' : ''}`}
              onClick={() => setActiveStep(step.id)}
            >
              <span className="hidden sm:inline">{step.name}</span>
            </li>
          ))}
        </ul>
      </div>

      {/* Step Content */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          {/* Step 1: Sender */}
          {activeStep === 1 && (
            <div className="space-y-6">
              <h3 className="text-xl font-bold flex items-center gap-2">
                <UserIcon className="w-6 h-6" />
                Информације о пошиљаоцу
              </h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Тип пошиљаоца</span>
                  </label>
                  <select
                    className="select select-bordered"
                    value={shipmentData.sender.type}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          type: Number(e.target.value),
                        },
                      })
                    }
                  >
                    <option value="1">Физичко лице</option>
                    <option value="2">Правно лице</option>
                    <option value="3">Матични број</option>
                  </select>
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Телефон</span>
                  </label>
                  <input
                    type="tel"
                    className="input input-bordered"
                    value={shipmentData.sender.phone}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          phone: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Име</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.sender.firstName}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          firstName: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Презиме</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.sender.lastName}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          lastName: e.target.value,
                        },
                      })
                    }
                  />
                </div>
              </div>

              <div className="divider">Адреса преузимања</div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Место</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.sender.place}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          place: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Улица</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.sender.street}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          street: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Број</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.sender.houseNumber}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          houseNumber: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Стан</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.sender.apartment}
                    placeholder="Опционално"
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          apartment: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Време од</span>
                  </label>
                  <input
                    type="time"
                    className="input input-bordered"
                    value={shipmentData.sender.timeFrom}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          timeFrom: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Време до</span>
                  </label>
                  <input
                    type="time"
                    className="input input-bordered"
                    value={shipmentData.sender.timeTo}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        sender: {
                          ...shipmentData.sender,
                          timeTo: e.target.value,
                        },
                      })
                    }
                  />
                </div>
              </div>
            </div>
          )}

          {/* Step 2: Receiver */}
          {activeStep === 2 && (
            <div className="space-y-6">
              <h3 className="text-xl font-bold flex items-center gap-2">
                <MapPinIcon className="w-6 h-6" />
                Информације о примаоцу
              </h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Име</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.receiver.firstName}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        receiver: {
                          ...shipmentData.receiver,
                          firstName: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Презиме</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.receiver.lastName}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        receiver: {
                          ...shipmentData.receiver,
                          lastName: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Телефон</span>
                  </label>
                  <input
                    type="tel"
                    className="input input-bordered"
                    value={shipmentData.receiver.phone}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        receiver: {
                          ...shipmentData.receiver,
                          phone: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Претходно обавештење</span>
                  </label>
                  <select
                    className="select select-bordered"
                    value={shipmentData.receiver.preNotification}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        receiver: {
                          ...shipmentData.receiver,
                          preNotification: Number(e.target.value),
                        },
                      })
                    }
                  >
                    {preNotifications.map((notif) => (
                      <option key={notif.value} value={notif.value}>
                        {notif.label}
                      </option>
                    ))}
                  </select>
                </div>
              </div>

              <div className="divider">Адреса доставе</div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Место</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.receiver.place}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        receiver: {
                          ...shipmentData.receiver,
                          place: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Улица</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.receiver.street}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        receiver: {
                          ...shipmentData.receiver,
                          street: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Број</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.receiver.houseNumber}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        receiver: {
                          ...shipmentData.receiver,
                          houseNumber: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Стан</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.receiver.apartment}
                    placeholder="Опционално"
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        receiver: {
                          ...shipmentData.receiver,
                          apartment: e.target.value,
                        },
                      })
                    }
                  />
                </div>
              </div>

              <div className="form-control">
                <label className="label">
                  <span className="label-text">Коментар за курира</span>
                </label>
                <textarea
                  className="textarea textarea-bordered"
                  value={shipmentData.receiver.comment}
                  onChange={(e) =>
                    setShipmentData({
                      ...shipmentData,
                      receiver: {
                        ...shipmentData.receiver,
                        comment: e.target.value,
                      },
                    })
                  }
                  placeholder="нпр. позвати 3. дугме"
                />
              </div>
            </div>
          )}

          {/* Step 3: Shipment Details */}
          {activeStep === 3 && (
            <div className="space-y-6">
              <h3 className="text-xl font-bold flex items-center gap-2">
                <CubeIcon className="w-6 h-6" />
                Детаљи пошиљке
              </h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Категорија</span>
                  </label>
                  <select
                    className="select select-bordered"
                    value={shipmentData.shipment.category}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        shipment: {
                          ...shipmentData.shipment,
                          category: Number(e.target.value),
                        },
                      })
                    }
                  >
                    {categories.map((cat) => (
                      <option key={cat.id} value={cat.id}>
                        {cat.icon} {cat.name}
                      </option>
                    ))}
                  </select>
                </div>

                {(shipmentData.shipment.category === 31 ||
                  shipmentData.shipment.category === 32) && (
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Тежина (кг)</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered"
                      value={shipmentData.shipment.weight}
                      onChange={(e) =>
                        setShipmentData({
                          ...shipmentData,
                          shipment: {
                            ...shipmentData.shipment,
                            weight: Number(e.target.value),
                          },
                        })
                      }
                    />
                  </div>
                )}

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Број пакета</span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered"
                    min="1"
                    max="99"
                    value={shipmentData.shipment.packages}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        shipment: {
                          ...shipmentData.shipment,
                          packages: Number(e.target.value),
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Начин плаћања</span>
                  </label>
                  <select
                    className="select select-bordered"
                    value={shipmentData.shipment.payType}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        shipment: {
                          ...shipmentData.shipment,
                          payType: Number(e.target.value),
                        },
                      })
                    }
                  >
                    {payTypes.map((type) => (
                      <option key={type.id} value={type.id}>
                        {type.icon} {type.name}
                      </option>
                    ))}
                  </select>
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Осигурање (РСД)</span>
                    <span className="label-text-alt">макс. 100,000</span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered"
                    max="100000"
                    value={shipmentData.shipment.insurance}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        shipment: {
                          ...shipmentData.shipment,
                          insurance: Number(e.target.value),
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Откупнина (РСД)</span>
                    <span className="label-text-alt">макс. 1,000,000</span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered"
                    max="1000000"
                    value={shipmentData.shipment.cashOnDelivery}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        shipment: {
                          ...shipmentData.shipment,
                          cashOnDelivery: Number(e.target.value),
                        },
                      })
                    }
                  />
                </div>
              </div>

              <div className="divider">Додатне опције</div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label cursor-pointer">
                    <span className="label-text">Лична достава</span>
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={shipmentData.shipment.personalDelivery}
                      onChange={(e) =>
                        setShipmentData({
                          ...shipmentData,
                          shipment: {
                            ...shipmentData.shipment,
                            personalDelivery: e.target.checked,
                          },
                        })
                      }
                    />
                  </label>
                </div>

                <div className="form-control">
                  <label className="label cursor-pointer">
                    <span className="label-text">Враћање фактура</span>
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={shipmentData.shipment.returnInvoices}
                      onChange={(e) =>
                        setShipmentData({
                          ...shipmentData,
                          shipment: {
                            ...shipmentData.shipment,
                            returnInvoices: e.target.checked,
                          },
                        })
                      }
                    />
                  </label>
                </div>

                <div className="form-control">
                  <label className="label cursor-pointer">
                    <span className="label-text">Враћање потврде</span>
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={shipmentData.shipment.returnConfirmation}
                      onChange={(e) =>
                        setShipmentData({
                          ...shipmentData,
                          shipment: {
                            ...shipmentData.shipment,
                            returnConfirmation: e.target.checked,
                          },
                        })
                      }
                    />
                  </label>
                </div>

                <div className="form-control">
                  <label className="label cursor-pointer">
                    <span className="label-text">Повратни пакет</span>
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={shipmentData.shipment.returnPackage}
                      onChange={(e) =>
                        setShipmentData({
                          ...shipmentData,
                          shipment: {
                            ...shipmentData.shipment,
                            returnPackage: e.target.checked,
                          },
                        })
                      }
                    />
                  </label>
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Јавна напомена</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.shipment.publicComment}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        shipment: {
                          ...shipmentData.shipment,
                          publicComment: e.target.value,
                        },
                      })
                    }
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">Приватна напомена</span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={shipmentData.shipment.privateComment}
                    onChange={(e) =>
                      setShipmentData({
                        ...shipmentData,
                        shipment: {
                          ...shipmentData.shipment,
                          privateComment: e.target.value,
                        },
                      })
                    }
                  />
                </div>
              </div>
            </div>
          )}

          {/* Step 4: Confirmation */}
          {activeStep === 4 && (
            <div className="space-y-6">
              <h3 className="text-xl font-bold flex items-center gap-2">
                <DocumentCheckIcon className="w-6 h-6" />
                Потврда пошиљке
              </h3>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <div className="card bg-base-200">
                  <div className="card-body">
                    <h4 className="font-bold">📤 Пошиљалац</h4>
                    <p>
                      {shipmentData.sender.firstName}{' '}
                      {shipmentData.sender.lastName}
                    </p>
                    <p className="text-sm">
                      {shipmentData.sender.place}, {shipmentData.sender.street}{' '}
                      {shipmentData.sender.houseNumber}
                    </p>
                    <p className="text-sm">📞 {shipmentData.sender.phone}</p>
                    <p className="text-sm">
                      ⏰ {shipmentData.sender.timeFrom} -{' '}
                      {shipmentData.sender.timeTo}
                    </p>
                  </div>
                </div>

                <div className="card bg-base-200">
                  <div className="card-body">
                    <h4 className="font-bold">📥 Прималац</h4>
                    <p>
                      {shipmentData.receiver.firstName}{' '}
                      {shipmentData.receiver.lastName}
                    </p>
                    <p className="text-sm">
                      {shipmentData.receiver.place},{' '}
                      {shipmentData.receiver.street}{' '}
                      {shipmentData.receiver.houseNumber}
                    </p>
                    <p className="text-sm">📞 {shipmentData.receiver.phone}</p>
                    {shipmentData.receiver.comment && (
                      <p className="text-sm">
                        💬 {shipmentData.receiver.comment}
                      </p>
                    )}
                  </div>
                </div>
              </div>

              <div className="card bg-base-200">
                <div className="card-body">
                  <h4 className="font-bold mb-4">📦 Детаљи пошиљке</h4>
                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
                    <div>
                      <p className="text-sm text-base-content/60">Категорија</p>
                      <p className="font-semibold">
                        {
                          categories.find(
                            (c) => c.id === shipmentData.shipment.category
                          )?.name
                        }
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-base-content/60">
                        Број пакета
                      </p>
                      <p className="font-semibold">
                        {shipmentData.shipment.packages}
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-base-content/60">Осигурање</p>
                      <p className="font-semibold">
                        {shipmentData.shipment.insurance.toLocaleString()} РСД
                      </p>
                    </div>
                    <div>
                      <p className="text-sm text-base-content/60">Откупнина</p>
                      <p className="font-semibold">
                        {shipmentData.shipment.cashOnDelivery.toLocaleString()}{' '}
                        РСД
                      </p>
                    </div>
                  </div>
                </div>
              </div>

              {showSuccess && (
                <div className="alert alert-success">
                  <TruckIcon className="w-6 h-6" />
                  <div>
                    <h3 className="font-bold">Пошиљка успешно креирана!</h3>
                    <div className="text-xs">Број пошиљке: 170123456</div>
                  </div>
                </div>
              )}

              <button
                onClick={handleCreateShipment}
                className="btn btn-primary btn-lg btn-block"
              >
                <TruckIcon className="w-5 h-5" />
                Креирај пошиљку
              </button>
            </div>
          )}

          {/* Navigation */}
          <div className="card-actions justify-between mt-6">
            <button
              className="btn btn-ghost"
              onClick={() => setActiveStep(Math.max(1, activeStep - 1))}
              disabled={activeStep === 1}
            >
              ← Назад
            </button>
            <button
              className="btn btn-primary"
              onClick={() => setActiveStep(Math.min(4, activeStep + 1))}
              disabled={activeStep === 4}
            >
              Даље →
            </button>
          </div>
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
                POST https://api.bex.rs:62502/postShipments
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
                    shipmentslist: [
                      {
                        shipmentId: 0,
                        serviceSpeed: 1,
                        shipmentType: 1,
                        shipmentCategory: shipmentData.shipment.category,
                        totalPackages: shipmentData.shipment.packages,
                        payType: shipmentData.shipment.payType,
                        insuranceAmount: shipmentData.shipment.insurance,
                        payToSender: shipmentData.shipment.cashOnDelivery,
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
