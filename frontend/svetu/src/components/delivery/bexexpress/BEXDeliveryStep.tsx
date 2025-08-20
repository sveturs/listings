'use client';

import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import BEXDeliverySelector from './BEXDeliverySelector';
import BEXAddressForm, { type BEXAddressData } from './BEXAddressForm';
import BEXParcelShopSelector from './BEXParcelShopSelector';
import {
  TruckIcon,
  ClockIcon,
  ShieldCheckIcon,
  InformationCircleIcon,
  CheckCircleIcon,
  ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';

// Схема валидации для BEX доставки
const bexDeliverySchema = z
  .object({
    method: z.enum(['courier', 'pickup_point', 'warehouse_pickup'], {
      required_error: 'Выберите способ доставки',
    }),
    address: z
      .object({
        recipient_name: z.string().min(2, 'Имя получателя обязательно'),
        recipient_phone: z.string().min(10, 'Телефон обязателен'),
        recipient_email: z.string().email().optional().or(z.literal('')),
        city: z.string().min(2, 'Город обязателен'),
        postal_code: z.string().optional(),
        street_address: z.string().optional(),
        street_number: z.string().optional(),
        apartment: z.string().optional(),
        floor: z.string().optional(),
        entrance: z.string().optional(),
        note: z.string().optional(),
        municipality_id: z.number().optional(),
        place_id: z.number().optional(),
        street_id: z.number().optional(),
      })
      .refine(
        (data) => {
          // Для курьерской доставки требуем адрес
          if (data.street_address && data.street_address.length < 3) {
            return false;
          }
          return true;
        },
        {
          message: 'Укажите адрес доставки',
          path: ['street_address'],
        }
      ),
    parcel_shop: z
      .object({
        id: z.number(),
        code: z.string(),
        name: z.string(),
        address: z.string(),
        city: z.string(),
        postal_code: z.string(),
        phone: z.string().optional(),
        latitude: z.number().optional(),
        longitude: z.number().optional(),
        working_hours: z.record(z.string()),
        services: z.array(z.string()),
        distance: z.number().optional(),
      })
      .optional(),
    rate: z
      .object({
        base_price: z.number(),
        cod_fee: z.number().optional(),
        insurance_fee: z.number().optional(),
        total_price: z.number(),
        estimated_days: z.number(),
      })
      .optional(),
    cod_amount: z.number().optional(),
    insurance_amount: z.number().optional(),
  })
  .refine(
    (data) => {
      // Валидация в зависимости от метода доставки
      if (data.method === 'courier') {
        return (
          data.address.street_address &&
          data.address.street_address.length >= 3 &&
          data.address.city &&
          data.address.recipient_name &&
          data.address.recipient_phone
        );
      } else if (data.method === 'pickup_point') {
        return (
          data.parcel_shop &&
          data.parcel_shop.id > 0 &&
          data.address.recipient_name &&
          data.address.recipient_phone
        );
      }
      // Для warehouse_pickup нужны только контактные данные
      return data.address.recipient_name && data.address.recipient_phone;
    },
    {
      message:
        'Заполните все обязательные поля для выбранного способа доставки',
      path: ['method'],
    }
  );

export type BEXDeliveryData = z.infer<typeof bexDeliverySchema>;

interface Props {
  onDataChange: (data: BEXDeliveryData | null, isValid: boolean) => void;
  initialData?: Partial<BEXDeliveryData>;
  orderWeight?: number;
  orderValue?: number;
  allowCOD?: boolean;
  className?: string;
}

export default function BEXDeliveryStep({
  onDataChange,
  initialData,
  orderWeight = 1,
  orderValue = 1000,
  allowCOD = true,
  className = '',
}: Props) {
  const [deliveryData, setDeliveryData] = useState<Partial<BEXDeliveryData>>({
    method: 'courier',
    cod_amount: 0,
    insurance_amount: orderValue,
    ...initialData,
  });

  const [calculatingRate, setCalculatingRate] = useState(false);

  const form = useForm<BEXDeliveryData>({
    resolver: zodResolver(bexDeliverySchema),
    defaultValues: {
      method: deliveryData.method || 'courier',
      address: deliveryData.address || {
        recipient_name: '',
        recipient_phone: '',
        recipient_email: '',
        city: '',
        postal_code: '',
        street_address: '',
        street_number: '',
        apartment: '',
        floor: '',
        entrance: '',
        note: '',
      },
      parcel_shop: deliveryData.parcel_shop,
      rate: deliveryData.rate,
      cod_amount: deliveryData.cod_amount || 0,
      insurance_amount: deliveryData.insurance_amount || orderValue,
    } as BEXDeliveryData,
    mode: 'onChange',
  });

  // Обновление родительского компонента при изменениях
  useEffect(() => {
    const subscription = form.watch((data) => {
      const isValid = form.formState.isValid;
      onDataChange(data as BEXDeliveryData, isValid);
    });
    return () => subscription.unsubscribe();
  }, [form, onDataChange]);

  // Инициализация формы с начальными данными
  useEffect(() => {
    if (initialData) {
      form.reset(initialData as BEXDeliveryData);
      setDeliveryData(initialData);
    }
  }, [initialData, form]);

  // Расчет стоимости доставки
  const calculateRate = async () => {
    if (
      !deliveryData.address?.city ||
      deliveryData.method === 'warehouse_pickup'
    ) {
      return;
    }

    setCalculatingRate(true);
    try {
      const response = await fetch('/api/v1/bex/calculate-rate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          weight_kg: orderWeight,
          cod_amount: deliveryData.cod_amount || 0,
          insurance_amount: deliveryData.insurance_amount || 0,
          recipient_city: deliveryData.address.city,
          delivery_method: deliveryData.method,
        }),
      });

      if (response.ok) {
        const data = await response.json();
        if (data.success && data.data) {
          const rate = {
            base_price: data.data.base_price,
            cod_fee: data.data.cod_fee || 0,
            insurance_fee: data.data.insurance_fee || 0,
            total_price: data.data.total_price,
            estimated_days: data.data.estimated_days || 2,
          };

          setDeliveryData({ ...deliveryData, rate });
          form.setValue('rate', rate);
        }
      }
    } catch (error) {
      console.error('Failed to calculate rate:', error);
    } finally {
      setCalculatingRate(false);
    }
  };

  // Обработчики изменений
  const handleMethodChange = (method: string) => {
    const updatedData = { ...deliveryData, method: method as any };
    setDeliveryData(updatedData);
    form.setValue('method', method as any);

    // Сбрасываем пункт выдачи при смене метода
    if (method !== 'pickup_point') {
      form.setValue('parcel_shop', undefined);
    }

    // Пересчитываем стоимость
    if (updatedData.address?.city) {
      calculateRate();
    }
  };

  const handleAddressChange = (address: BEXAddressData) => {
    const updatedData = { ...deliveryData, address };
    setDeliveryData(updatedData);
    form.setValue('address', address);

    // Пересчитываем стоимость при изменении города
    if (address.city !== deliveryData.address?.city) {
      calculateRate();
    }
  };

  const handleParcelShopSelect = (shop: any) => {
    const updatedData = { ...deliveryData, parcel_shop: shop };
    setDeliveryData(updatedData);
    form.setValue('parcel_shop', shop);
  };

  const currentMethod = form.watch('method');
  const currentAddress = form.watch('address');
  const currentRate = form.watch('rate');

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Заголовок */}
      <div className="text-center">
        <div className="flex items-center justify-center gap-3 mb-3">
          <TruckIcon className="w-8 h-8 text-primary" />
          <h2 className="text-2xl font-bold">Доставка BEX Express</h2>
        </div>
        <p className="text-base-content/70">
          Надежная доставка по всей Сербии с отслеживанием в реальном времени
        </p>
      </div>

      {/* Информация о заказе */}
      <div className="card bg-info/10 border border-info/20">
        <div className="card-body p-4">
          <div className="flex items-start gap-3">
            <InformationCircleIcon className="w-5 h-5 text-info mt-0.5 flex-shrink-0" />
            <div className="text-sm">
              <p className="font-semibold mb-1">Информация о посылке</p>
              <div className="flex flex-wrap gap-4 text-base-content/70">
                <span>
                  Вес: <strong>{orderWeight} кг</strong>
                </span>
                <span>
                  Стоимость товаров: <strong>{orderValue} RSD</strong>
                </span>
                {allowCOD && (
                  <span>
                    Наложенный платеж: <strong>Доступен</strong>
                  </span>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Выбор способа доставки */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body p-6">
          <h3 className="font-semibold text-lg mb-4">
            1. Выберите способ доставки
          </h3>

          <BEXDeliverySelector
            onMethodChange={handleMethodChange}
            selectedMethod={currentMethod}
            weight={orderWeight}
            insuranceAmount={deliveryData.insurance_amount}
            hasCOD={(deliveryData.cod_amount || 0) > 0}
            recipientCity={currentAddress?.city}
          />

          {/* Ошибки валидации метода */}
          {form.formState.errors.method && (
            <div className="alert alert-error mt-4">
              <ExclamationTriangleIcon className="w-5 h-5" />
              <span>{form.formState.errors.method.message}</span>
            </div>
          )}
        </div>
      </div>

      {/* Данные получателя и адрес */}
      {currentMethod && (
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body p-6">
            <h3 className="font-semibold text-lg mb-4">
              2.{' '}
              {currentMethod === 'courier'
                ? 'Укажите адрес и данные получателя'
                : 'Укажите данные получателя'}
            </h3>

            <BEXAddressForm
              onAddressChange={handleAddressChange}
              initialAddress={currentAddress}
              deliveryMethod={currentMethod}
            />

            {/* Ошибки валидации адреса */}
            {form.formState.errors.address && (
              <div className="alert alert-error mt-4">
                <ExclamationTriangleIcon className="w-5 h-5" />
                <span>
                  {form.formState.errors.address.street_address?.message ||
                    form.formState.errors.address.recipient_name?.message ||
                    form.formState.errors.address.recipient_phone?.message ||
                    form.formState.errors.address.city?.message ||
                    'Заполните все обязательные поля'}
                </span>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Выбор пункта выдачи */}
      {currentMethod === 'pickup_point' && currentAddress?.city && (
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body p-6">
            <h3 className="font-semibold text-lg mb-4">
              3. Выберите пункт выдачи
            </h3>

            <BEXParcelShopSelector
              selectedCity={currentAddress.city}
              onShopSelect={handleParcelShopSelect}
              selectedShop={deliveryData.parcel_shop}
            />

            {/* Ошибки валидации пункта */}
            {form.formState.errors.parcel_shop &&
              currentMethod === 'pickup_point' && (
                <div className="alert alert-error mt-4">
                  <ExclamationTriangleIcon className="w-5 h-5" />
                  <span>Выберите пункт выдачи</span>
                </div>
              )}
          </div>
        </div>
      )}

      {/* Дополнительные опции */}
      {currentMethod && currentMethod !== 'warehouse_pickup' && (
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body p-6">
            <h3 className="font-semibold text-lg mb-4">
              Дополнительные услуги
            </h3>

            <div className="space-y-4">
              {/* Наложенный платеж */}
              {allowCOD && (
                <div className="form-control">
                  <label className="label cursor-pointer justify-start gap-4">
                    <input
                      type="checkbox"
                      className="checkbox checkbox-primary"
                      checked={(deliveryData.cod_amount || 0) > 0}
                      onChange={(e) => {
                        const codAmount = e.target.checked ? orderValue : 0;
                        const updatedData = {
                          ...deliveryData,
                          cod_amount: codAmount,
                        };
                        setDeliveryData(updatedData);
                        form.setValue('cod_amount', codAmount);
                        calculateRate();
                      }}
                    />
                    <div className="flex-1">
                      <span className="label-text font-medium">
                        Наложенный платеж (оплата при получении)
                      </span>
                      <p className="text-sm text-base-content/60">
                        Оплатите товар при получении наличными или картой
                      </p>
                    </div>
                    {(deliveryData.cod_amount || 0) > 0 && (
                      <div className="text-right">
                        <div className="text-sm text-base-content/60">
                          К оплате:
                        </div>
                        <div className="font-bold">
                          {deliveryData.cod_amount} RSD
                        </div>
                        <div className="text-xs text-warning">
                          + комиссия 150 RSD
                        </div>
                      </div>
                    )}
                  </label>
                </div>
              )}

              {/* Страхование */}
              <div className="form-control">
                <label className="label cursor-pointer justify-start gap-4">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-primary"
                    checked={(deliveryData.insurance_amount || 0) > 0}
                    onChange={(e) => {
                      const insuranceAmount = e.target.checked ? orderValue : 0;
                      const updatedData = {
                        ...deliveryData,
                        insurance_amount: insuranceAmount,
                      };
                      setDeliveryData(updatedData);
                      form.setValue('insurance_amount', insuranceAmount);
                      calculateRate();
                    }}
                  />
                  <div className="flex-1">
                    <span className="label-text font-medium">
                      Страхование груза
                    </span>
                    <p className="text-sm text-base-content/60">
                      Полная компенсация стоимости при утере или повреждении
                    </p>
                  </div>
                  {(deliveryData.insurance_amount || 0) > 0 && (
                    <div className="text-right">
                      <div className="text-sm text-base-content/60">
                        Страховая сумма:
                      </div>
                      <div className="font-bold">
                        {deliveryData.insurance_amount} RSD
                      </div>
                      <div className="text-xs text-success">Комиссия 0.5%</div>
                    </div>
                  )}
                </label>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Стоимость и время доставки */}
      {currentRate && currentMethod !== 'warehouse_pickup' && (
        <div className="card bg-gradient-to-r from-primary/5 to-secondary/5 shadow-lg">
          <div className="card-body p-6">
            <h3 className="font-semibold text-lg mb-4">Стоимость доставки</h3>

            {calculatingRate ? (
              <div className="flex justify-center items-center h-24">
                <span className="loading loading-spinner loading-lg"></span>
              </div>
            ) : (
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                {/* Расчет стоимости */}
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span>Базовая стоимость:</span>
                    <span>{currentRate.base_price} RSD</span>
                  </div>
                  {(deliveryData.cod_amount || 0) > 0 &&
                    currentRate.cod_fee && (
                      <div className="flex justify-between">
                        <span>Комиссия за наложенный платеж:</span>
                        <span>{currentRate.cod_fee} RSD</span>
                      </div>
                    )}
                  {(deliveryData.insurance_amount || 0) > 0 &&
                    currentRate.insurance_fee && (
                      <div className="flex justify-between">
                        <span>Страхование:</span>
                        <span>{currentRate.insurance_fee} RSD</span>
                      </div>
                    )}
                  <div className="divider my-2"></div>
                  <div className="flex justify-between text-lg font-bold">
                    <span>Итого:</span>
                    <span className="text-primary">
                      {currentRate.total_price} RSD
                    </span>
                  </div>
                </div>

                {/* Время доставки */}
                <div className="flex items-center justify-center">
                  <div className="text-center p-4 bg-base-200 rounded-lg w-full">
                    <ClockIcon className="w-8 h-8 mx-auto mb-2 text-info" />
                    <div className="text-sm text-base-content/70">
                      Срок доставки
                    </div>
                    <div className="text-2xl font-bold text-primary">
                      {currentRate.estimated_days}{' '}
                      {currentRate.estimated_days === 1 ? 'день' : 'дня'}
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Гарантии */}
            <div className="flex items-center gap-3 mt-4 pt-4 border-t border-base-300">
              <ShieldCheckIcon className="w-5 h-5 text-success" />
              <div className="text-sm">
                <span className="font-medium">Гарантии BEX Express: </span>
                <span className="text-base-content/70">
                  Отслеживание 24/7 • СМС-уведомления • Страхование груза
                </span>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Статус валидации */}
      <div className="text-center">
        {form.formState.isValid ? (
          <div className="alert alert-success">
            <CheckCircleIcon className="w-6 h-6" />
            <span>Все данные для доставки заполнены корректно</span>
          </div>
        ) : (
          <div className="alert alert-warning">
            <ExclamationTriangleIcon className="w-6 h-6" />
            <span>Заполните все обязательные поля для продолжения</span>
          </div>
        )}
      </div>
    </div>
  );
}
