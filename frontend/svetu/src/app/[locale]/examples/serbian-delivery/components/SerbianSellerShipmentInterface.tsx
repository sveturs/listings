'use client';

import { useState } from 'react';
import { PrinterIcon } from '@heroicons/react/24/outline';

export default function SerbianSellerShipmentInterface() {
  const [formData, setFormData] = useState({
    courier: 'aks',
    recipient: 'Марко Петровић',
    phone: '066/123-456',
    address: 'Булевар Ослобођења 15, Нови Сад',
    weight: '0.5',
    cod: '2500',
    insurance: true,
  });

  const couriers = [
    { id: 'aks', name: 'AKS Express', price: 200, color: 'orange' },
    { id: 'post-express', name: 'Post Express', price: 150, color: 'blue' },
    { id: 'city-express', name: 'City Express', price: 180, color: 'green' },
    { id: 'yettel-post', name: 'Yettel Post', price: 120, color: 'purple' },
  ];

  const selectedCourier = couriers.find((c) => c.id === formData.courier);

  const calculatePrice = () => {
    const basePrice = selectedCourier?.price || 0;
    const weightExtra = parseFloat(formData.weight) > 1 ? 50 : 0;
    const insuranceExtra = formData.insurance ? 30 : 0;
    return basePrice + weightExtra + insuranceExtra;
  };

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h3 className="card-title mb-4">Нова пошиљка - српски курири</h3>

        <div className="space-y-4">
          {/* Courier Selection */}
          <div className="form-control">
            <label className="label">
              <span className="label-text font-semibold">Курирска служба</span>
            </label>
            <select
              className="select select-bordered"
              value={formData.courier}
              onChange={(e) =>
                setFormData({ ...formData, courier: e.target.value })
              }
            >
              {couriers.map((courier) => (
                <option key={courier.id} value={courier.id}>
                  {courier.name} - {courier.price} РСД
                </option>
              ))}
            </select>
          </div>

          {/* Recipient Info */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">Име примаоца</span>
              </label>
              <input
                type="text"
                className="input input-bordered"
                value={formData.recipient}
                onChange={(e) =>
                  setFormData({ ...formData, recipient: e.target.value })
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
                value={formData.phone}
                onChange={(e) =>
                  setFormData({ ...formData, phone: e.target.value })
                }
                placeholder="06x/xxx-xxx"
              />
            </div>
          </div>

          <div className="form-control">
            <label className="label">
              <span className="label-text">Адреса за доставу</span>
            </label>
            <input
              type="text"
              className="input input-bordered"
              value={formData.address}
              onChange={(e) =>
                setFormData({ ...formData, address: e.target.value })
              }
            />
          </div>

          {/* Package Info */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="form-control">
              <label className="label">
                <span className="label-text">Тежина (кг)</span>
              </label>
              <input
                type="number"
                step="0.1"
                className="input input-bordered"
                value={formData.weight}
                onChange={(e) =>
                  setFormData({ ...formData, weight: e.target.value })
                }
              />
            </div>

            <div className="form-control">
              <label className="label">
                <span className="label-text">Поштарина (РСД)</span>
              </label>
              <input
                type="number"
                className="input input-bordered"
                value={formData.cod}
                onChange={(e) =>
                  setFormData({ ...formData, cod: e.target.value })
                }
              />
            </div>
          </div>

          {/* Insurance */}
          <div className="form-control">
            <label className="cursor-pointer label justify-start gap-3">
              <input
                type="checkbox"
                className="checkbox"
                checked={formData.insurance}
                onChange={(e) =>
                  setFormData({ ...formData, insurance: e.target.checked })
                }
              />
              <span className="label-text">Додатно осигурање (+30 РСД)</span>
            </label>
          </div>

          {/* Price Summary */}
          <div className="bg-base-200 p-4 rounded-lg">
            <div className="flex justify-between items-center mb-2">
              <span>Основна цена:</span>
              <span>{selectedCourier?.price} РСД</span>
            </div>
            {parseFloat(formData.weight) > 1 && (
              <div className="flex justify-between items-center mb-2">
                <span>Додатна тежина:</span>
                <span>+50 РСД</span>
              </div>
            )}
            {formData.insurance && (
              <div className="flex justify-between items-center mb-2">
                <span>Осигурање:</span>
                <span>+30 РСД</span>
              </div>
            )}
            <hr className="my-2" />
            <div className="flex justify-between items-center font-bold">
              <span>Укупно:</span>
              <span className="text-primary">{calculatePrice()} РСД</span>
            </div>
          </div>

          {/* Actions */}
          <div className="flex gap-2">
            <button className="btn btn-primary flex-1 gap-2">
              <PrinterIcon className="w-5 h-5" />
              Штампај налепницу
            </button>
            <button className="btn btn-outline">Сачувај</button>
          </div>

          {/* Serbian specific tips */}
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
              <h4 className="font-semibold">
                Важне напомене за српско тржиште:
              </h4>
              <ul className="text-sm mt-1 space-y-1">
                <li>• AKS не ради недељом</li>
                <li>• Post Express пунктови раде до 22:00</li>
                <li>• Максимална поштарина: 300.000 РСД</li>
                <li>• Обавезна је лична карта при преузимању</li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
