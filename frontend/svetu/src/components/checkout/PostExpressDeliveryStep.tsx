'use client';

import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useTranslations } from 'next-intl';
import {
  PostExpressDeliverySelector,
  PostExpressAddressForm,
  PostExpressOfficeSelector,
} from '@/components/delivery/postexpress';

// Validation schema for Post Express delivery
const postExpressDeliverySchema = z
  .object({
    method: z.enum(['courier', 'pickup_point', 'warehouse_pickup'], {
      required_error: 'Please select a delivery method',
    }),
    address: z
      .object({
        recipient_name: z.string().min(2, 'Recipient name is required'),
        recipient_phone: z.string().min(10, 'Phone number is required'),
        recipient_email: z.string().email().optional().or(z.literal('')),
        city: z.string().min(2, 'City is required'),
        postal_code: z.string().optional(),
        street_address: z.string().optional(),
        apartment: z.string().optional(),
        floor: z.string().optional(),
        entrance: z.string().optional(),
        note: z.string().optional(),
      })
      .refine(
        (data) => {
          // Требуем адрес только для курьерской доставки
          if (data.street_address && data.street_address.length < 5) {
            return false;
          }
          return true;
        },
        {
          message: 'Street address is required for courier delivery',
          path: ['street_address'],
        }
      ),
    office: z
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
        total_price: z.number(),
        estimated_days: z.number(),
      })
      .optional(),
    cod_amount: z.number().optional(),
  })
  .refine(
    (data) => {
      // Валидация в зависимости от метода доставки
      if (data.method === 'courier') {
        return (
          data.address.street_address && data.address.street_address.length >= 5
        );
      } else if (data.method === 'pickup_point') {
        return data.office && data.office.id > 0;
      }
      return true; // для warehouse_pickup дополнительных требований нет
    },
    {
      message: 'Required fields for selected delivery method are missing',
      path: ['method'],
    }
  );

export type PostExpressDeliveryData = z.infer<typeof postExpressDeliverySchema>;

interface Props {
  onDataChange: (
    data: PostExpressDeliveryData | null,
    isValid: boolean
  ) => void;
  initialData?: Partial<PostExpressDeliveryData>;
  orderWeight?: number;
  orderValue?: number;
  allowCOD?: boolean;
  className?: string;
}

export default function PostExpressDeliveryStep({
  onDataChange,
  initialData,
  orderWeight = 1,
  orderValue = 1000,
  allowCOD = true,
  className = '',
}: Props) {
  const t = useTranslations('delivery');
  const tCheckout = useTranslations('checkout');

  const [deliveryData, setDeliveryData] = useState<
    Partial<PostExpressDeliveryData>
  >({
    method: 'courier',
    cod_amount: 0,
    ...initialData,
  });

  const form = useForm<PostExpressDeliveryData>({
    resolver: zodResolver(postExpressDeliverySchema),
    defaultValues: {
      method: deliveryData.method || 'courier',
      address: deliveryData.address || {
        recipient_name: '',
        recipient_phone: '',
        recipient_email: '',
        city: '',
        postal_code: '',
        street_address: '',
        apartment: '',
        floor: '',
        entrance: '',
        note: '',
      },
      office: deliveryData.office,
      rate: deliveryData.rate,
      cod_amount: deliveryData.cod_amount || 0,
    } as PostExpressDeliveryData,
    mode: 'onChange',
  });

  // Update parent when data changes
  useEffect(() => {
    const subscription = form.watch((data) => {
      const isValid = form.formState.isValid;
      onDataChange(data as PostExpressDeliveryData, isValid);
    });
    return () => subscription.unsubscribe();
  }, [form, onDataChange]);

  // Initialize form with initial data
  useEffect(() => {
    if (initialData) {
      form.reset(initialData as PostExpressDeliveryData);
      setDeliveryData(initialData);
    }
  }, [initialData, form]);

  const handleMethodChange = (method: string) => {
    const updatedData = { ...deliveryData, method: method as any };
    setDeliveryData(updatedData);
    form.setValue('method', method as any);

    // Reset office when changing away from pickup_point
    if (method !== 'pickup_point') {
      form.setValue('office', undefined);
    }
  };

  const handleAddressChange = (address: any) => {
    const updatedData = { ...deliveryData, address };
    setDeliveryData(updatedData);
    form.setValue('address', address);
  };

  const handleOfficeSelect = (office: any) => {
    const updatedData = { ...deliveryData, office };
    setDeliveryData(updatedData);
    form.setValue('office', office);
  };

  // Unused but kept for future implementation
  // const handleRateCalculated = (rate: any) => {
  //   const updatedData = { ...deliveryData, rate };
  //   setDeliveryData(updatedData);
  //   form.setValue('rate', rate);
  // };

  const currentMethod = form.watch('method');
  const currentAddress = form.watch('address');

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Заголовок */}
      <div className="text-center">
        <h2 className="text-2xl font-bold mb-2">
          {tCheckout('shipping.title')}
        </h2>
        <p className="text-base-content/70">
          Выберите удобный способ получения заказа
        </p>
      </div>

      {/* Выбор способа доставки */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body p-6">
          <h3 className="font-semibold text-lg mb-4">{t('selectMethod')}</h3>

          <PostExpressDeliverySelector
            onMethodChange={handleMethodChange}
            selectedMethod={currentMethod}
            weight={orderWeight}
            insuranceAmount={orderValue}
            hasCOD={allowCOD && (deliveryData.cod_amount || 0) > 0}
            recipientCity={currentAddress?.city}
          />

          {/* Ошибки валидации метода */}
          {form.formState.errors.method && (
            <div className="alert alert-error mt-4">
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
              {currentMethod === 'courier'
                ? t('selectAddress')
                : t('recipientDetails')}
            </h3>

            <PostExpressAddressForm
              onAddressChange={handleAddressChange}
              initialAddress={currentAddress}
              deliveryMethod={currentMethod}
            />

            {/* Ошибки валидации адреса */}
            {form.formState.errors.address && (
              <div className="alert alert-error mt-4">
                <span>
                  {form.formState.errors.address.street_address?.message ||
                    form.formState.errors.address.recipient_name?.message ||
                    form.formState.errors.address.recipient_phone?.message ||
                    form.formState.errors.address.city?.message ||
                    'Please fill in all required fields'}
                </span>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Выбор отделения для pickup_point */}
      {currentMethod === 'pickup_point' && currentAddress?.city && (
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body p-6">
            <h3 className="font-semibold text-lg mb-4">{t('selectOffice')}</h3>

            <PostExpressOfficeSelector
              selectedCity={currentAddress.city}
              onOfficeSelect={handleOfficeSelect}
              selectedOffice={deliveryData.office}
            />

            {/* Ошибки валидации отделения */}
            {form.formState.errors.office && (
              <div className="alert alert-error mt-4">
                <span>{t('errors.officeRequired')}</span>
              </div>
            )}
          </div>
        </div>
      )}

      {/* Наложенный платеж (если разрешен) */}
      {allowCOD && currentMethod !== 'warehouse_pickup' && (
        <div className="card bg-base-100 shadow-lg">
          <div className="card-body p-6">
            <h3 className="font-semibold text-lg mb-4">Дополнительные опции</h3>

            <div className="form-control">
              <label className="label cursor-pointer">
                <span className="label-text">
                  Наложенный платеж (оплата при получении)
                </span>
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
                  }}
                />
              </label>
              {(deliveryData.cod_amount || 0) > 0 && (
                <label className="label">
                  <span className="label-text-alt">
                    Сумма к оплате: {deliveryData.cod_amount} RSD + комиссия 45
                    RSD
                  </span>
                </label>
              )}
            </div>
          </div>
        </div>
      )}

      {/* Информация о стоимости и времени доставки */}
      {deliveryData.rate && (
        <div className="card bg-gradient-to-r from-primary/5 to-secondary/5 shadow-lg">
          <div className="card-body p-6">
            <h3 className="font-semibold text-lg mb-4">Стоимость доставки</h3>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div className="space-y-2">
                <div className="flex justify-between">
                  <span>Базовая стоимость:</span>
                  <span>{deliveryData.rate.base_price?.toFixed(0)} RSD</span>
                </div>
                {(deliveryData.cod_amount || 0) > 0 && (
                  <div className="flex justify-between">
                    <span>Комиссия за наложенный платеж:</span>
                    <span>45 RSD</span>
                  </div>
                )}
                <div className="divider my-1"></div>
                <div className="flex justify-between text-lg font-bold">
                  <span>Итого:</span>
                  <span className="text-primary">
                    {currentMethod === 'warehouse_pickup'
                      ? 'Бесплатно'
                      : `${(deliveryData.rate.total_price + ((deliveryData.cod_amount || 0) > 0 ? 45 : 0)).toFixed(0)} RSD`}
                  </span>
                </div>
              </div>

              <div className="space-y-2">
                <div className="text-center p-4 bg-base-200 rounded-lg">
                  <div className="text-sm text-base-content/70">
                    Ожидаемое время доставки
                  </div>
                  <div className="text-xl font-bold text-primary">
                    {deliveryData.rate.estimated_days || '1-2'}{' '}
                    {t('businessDays')}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Validation Summary */}
      <div className="text-center">
        {form.formState.isValid ? (
          <div className="alert alert-success">
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M5 13l4 4L19 7"
              />
            </svg>
            <span>Все данные заполнены правильно</span>
          </div>
        ) : (
          <div className="alert alert-warning">
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16c-.77.833.192 2.5 1.732 2.5z"
              />
            </svg>
            <span>Заполните все обязательные поля для продолжения</span>
          </div>
        )}
      </div>
    </div>
  );
}
