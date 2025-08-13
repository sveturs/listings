'use client';

import { useState } from 'react';
import {
  PrinterIcon,
  DocumentDuplicateIcon,
  CheckCircleIcon,
  InformationCircleIcon,
  CurrencyDollarIcon,
  ShieldCheckIcon,
  DocumentCheckIcon,
} from '@heroicons/react/24/outline';

export default function SellerShipmentInterface() {
  const [isCreating, setIsCreating] = useState(false);
  const [shipmentCreated, setShipmentCreated] = useState(false);
  const [formData, setFormData] = useState({
    recipientName: '',
    recipientPhone: '',
    recipientAddress: '',
    recipientCity: '',
    weight: '',
    dimensions: '',
    value: '',
    cod: '',
    insurance: false,
    returnDocs: false,
    comment: '',
  });

  const handleCreateShipment = () => {
    setIsCreating(true);
    setTimeout(() => {
      setIsCreating(false);
      setShipmentCreated(true);
    }, 2000);
  };

  const handleInputChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
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

  if (shipmentCreated) {
    return (
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="text-center py-8">
            <CheckCircleIcon className="w-16 h-16 text-success mx-auto mb-4" />
            <h3 className="text-xl font-bold mb-2">Отправление создано!</h3>
            <p className="text-base-content/60 mb-4">
              Номер отслеживания: 170123457
            </p>

            <div className="flex gap-3 justify-center mb-6">
              <button className="btn btn-primary">
                <PrinterIcon className="w-5 h-5" />
                Печать этикетки
              </button>
              <button className="btn btn-outline">
                <DocumentDuplicateIcon className="w-5 h-5" />
                Копировать номер
              </button>
            </div>

            <button
              className="btn btn-ghost"
              onClick={() => {
                setShipmentCreated(false);
                setFormData({
                  recipientName: '',
                  recipientPhone: '',
                  recipientAddress: '',
                  recipientCity: '',
                  weight: '',
                  dimensions: '',
                  value: '',
                  cod: '',
                  insurance: false,
                  returnDocs: false,
                  comment: '',
                });
              }}
            >
              Создать новое отправление
            </button>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body p-4 sm:p-6">
        <h3 className="card-title text-base sm:text-lg mb-3 sm:mb-4">
          Новое отправление
        </h3>

        <div className="space-y-4">
          {/* Order Info (Auto-filled) */}
          <div className="alert alert-info">
            <InformationCircleIcon className="w-4 h-4 sm:w-5 sm:h-5 flex-shrink-0" />
            <div>
              <div className="font-semibold text-sm sm:text-base">
                Заказ #12345
              </div>
              <div className="text-xs sm:text-sm">
                Данные заполнены автоматически из заказа
              </div>
            </div>
          </div>

          {/* Recipient Information */}
          <div className="space-y-3">
            <div className="text-sm font-semibold text-base-content/60">
              Получатель
            </div>

            <input
              type="text"
              name="recipientName"
              placeholder="ФИО получателя"
              className="input input-bordered input-sm sm:input-md w-full"
              value={formData.recipientName}
              onChange={handleInputChange}
            />

            <input
              type="tel"
              name="recipientPhone"
              placeholder="Телефон получателя"
              className="input input-bordered input-sm sm:input-md w-full"
              value={formData.recipientPhone}
              onChange={handleInputChange}
            />

            <select
              name="recipientCity"
              className="select select-bordered select-sm sm:select-md w-full"
              value={formData.recipientCity}
              onChange={handleInputChange}
            >
              <option value="">Выберите город</option>
              <option value="belgrade">Белград</option>
              <option value="novi-sad">Нови Сад</option>
              <option value="nis">Ниш</option>
              <option value="kragujevac">Крагуевац</option>
            </select>

            <input
              type="text"
              name="recipientAddress"
              placeholder="Адрес доставки"
              className="input input-bordered input-sm sm:input-md w-full"
              value={formData.recipientAddress}
              onChange={handleInputChange}
            />
          </div>

          {/* Package Information */}
          <div className="space-y-3">
            <div className="text-sm font-semibold text-base-content/60">
              Параметры посылки
            </div>

            <div className="grid grid-cols-2 gap-2 sm:gap-3">
              <input
                type="text"
                name="weight"
                placeholder="Вес (кг)"
                className="input input-bordered input-sm sm:input-md"
                value={formData.weight}
                onChange={handleInputChange}
              />
              <input
                type="text"
                name="dimensions"
                placeholder="Размеры (см)"
                className="input input-bordered input-sm sm:input-md"
                value={formData.dimensions}
                onChange={handleInputChange}
              />
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 sm:gap-3">
              <div>
                <label className="label py-1">
                  <span className="label-text text-xs">
                    Объявленная стоимость
                  </span>
                </label>
                <input
                  type="text"
                  name="value"
                  placeholder="0 RSD"
                  className="input input-bordered input-sm sm:input-md w-full"
                  value={formData.value}
                  onChange={handleInputChange}
                />
              </div>
              <div>
                <label className="label py-1">
                  <span className="label-text text-xs">Наложенный платеж</span>
                </label>
                <input
                  type="text"
                  name="cod"
                  placeholder="0 RSD"
                  className="input input-bordered input-sm sm:input-md w-full"
                  value={formData.cod}
                  onChange={handleInputChange}
                />
              </div>
            </div>
          </div>

          {/* Additional Services */}
          <div className="space-y-3">
            <div className="text-sm font-semibold text-base-content/60">
              Дополнительные услуги
            </div>

            <label className="flex items-start sm:items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                name="insurance"
                className="checkbox checkbox-primary checkbox-sm sm:checkbox-md mt-0.5 sm:mt-0"
                checked={formData.insurance}
                onChange={handleInputChange}
              />
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <ShieldCheckIcon className="w-3 h-3 sm:w-4 sm:h-4" />
                  <span className="font-medium text-sm sm:text-base">
                    Страхование
                  </span>
                </div>
                <div className="text-xs text-base-content/60">
                  +2% от объявленной стоимости
                </div>
              </div>
            </label>

            <label className="flex items-start sm:items-center gap-3 cursor-pointer">
              <input
                type="checkbox"
                name="returnDocs"
                className="checkbox checkbox-primary checkbox-sm sm:checkbox-md mt-0.5 sm:mt-0"
                checked={formData.returnDocs}
                onChange={handleInputChange}
              />
              <div className="flex-1">
                <div className="flex items-center gap-2">
                  <DocumentCheckIcon className="w-3 h-3 sm:w-4 sm:h-4" />
                  <span className="font-medium text-sm sm:text-base">
                    Возврат документов
                  </span>
                </div>
                <div className="text-xs text-base-content/60">+150 RSD</div>
              </div>
            </label>
          </div>

          {/* Comment */}
          <div>
            <label className="label py-1">
              <span className="label-text text-xs">
                Комментарий для курьера
              </span>
            </label>
            <textarea
              name="comment"
              className="textarea textarea-bordered textarea-sm sm:textarea-md w-full"
              placeholder="Например: Позвонить за 30 минут до доставки"
              rows={2}
              value={formData.comment}
              onChange={handleInputChange}
            />
          </div>

          {/* Cost Summary */}
          <div className="bg-base-200 rounded-lg p-3 sm:p-4">
            <div className="space-y-2 text-xs sm:text-sm">
              <div className="flex justify-between">
                <span>Базовая стоимость доставки</span>
                <span>350 RSD</span>
              </div>
              {formData.insurance && (
                <div className="flex justify-between">
                  <span>Страхование</span>
                  <span>70 RSD</span>
                </div>
              )}
              {formData.returnDocs && (
                <div className="flex justify-between">
                  <span>Возврат документов</span>
                  <span>150 RSD</span>
                </div>
              )}
              <div className="divider my-1"></div>
              <div className="flex justify-between font-bold">
                <span>Итого</span>
                <span className="text-primary">
                  {350 +
                    (formData.insurance ? 70 : 0) +
                    (formData.returnDocs ? 150 : 0)}{' '}
                  RSD
                </span>
              </div>
            </div>
          </div>

          {/* Action Buttons */}
          <div className="flex gap-2 sm:gap-3">
            <button
              className={`btn btn-primary btn-sm sm:btn-md flex-1 ${isCreating ? 'loading' : ''}`}
              onClick={handleCreateShipment}
              disabled={isCreating}
            >
              {!isCreating && (
                <CheckCircleIcon className="w-4 h-4 sm:w-5 sm:h-5" />
              )}
              <span className="text-xs sm:text-base">Создать отправление</span>
            </button>
            <button className="btn btn-ghost btn-sm sm:btn-md">
              <span className="text-xs sm:text-base">Отмена</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
