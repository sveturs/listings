'use client';

import { useState, useEffect } from 'react';
import {
  ChevronLeftIcon,
  ChevronRightIcon,
  CheckIcon,
  TruckIcon,
  MapPinIcon,
  BuildingStorefrontIcon,
} from '@heroicons/react/24/outline';
import { useTranslations } from 'next-intl';

// Import our components
import PostExpressDeliverySelector from './PostExpressDeliverySelector';
import PostExpressAddressForm from './PostExpressAddressForm';
import PostExpressOfficeSelector from './PostExpressOfficeSelector';
import PostExpressRateCalculator from './PostExpressRateCalculator';

interface DeliveryData {
  method: string;
  address?: any;
  office?: any;
  rate?: any;
  weight: number;
  declaredValue: number;
  codAmount: number;
}

interface Props {
  onDeliveryComplete: (deliveryData: DeliveryData) => void;
  initialData?: Partial<DeliveryData>;
  orderWeight?: number;
  orderValue?: number;
  allowCOD?: boolean;
  className?: string;
}

const STEPS = [
  { id: 'method', title: 'Способ доставки', icon: TruckIcon },
  { id: 'details', title: 'Детали доставки', icon: MapPinIcon },
  { id: 'confirmation', title: 'Подтверждение', icon: CheckIcon },
];

export default function PostExpressDeliveryFlow({
  onDeliveryComplete,
  initialData,
  orderWeight = 1,
  orderValue = 1000,
  allowCOD = true,
  className = '',
}: Props) {
  // const t = useTranslations('delivery');
  // allowCOD parameter will be used in future implementations
  const [currentStep, setCurrentStep] = useState(0);
  const [deliveryData, setDeliveryData] = useState<DeliveryData>({
    method: '',
    weight: orderWeight,
    declaredValue: orderValue,
    codAmount: 0,
    ...initialData,
  });

  useEffect(() => {
    if (initialData) {
      setDeliveryData((prev) => ({ ...prev, ...initialData }));
    }
  }, [initialData]);

  const handleMethodChange = (method: string) => {
    setDeliveryData((prev) => ({ ...prev, method }));
  };

  const handleAddressChange = (address: any) => {
    setDeliveryData((prev) => ({ ...prev, address }));
  };

  const handleOfficeSelect = (office: any) => {
    setDeliveryData((prev) => ({ ...prev, office }));
  };

  const handleRateCalculated = (rate: any) => {
    setDeliveryData((prev) => ({ ...prev, rate }));
  };

  const canProceedToNext = () => {
    switch (currentStep) {
      case 0: // Method selection
        return deliveryData.method !== '';
      case 1: // Details
        if (deliveryData.method === 'courier') {
          return (
            deliveryData.address &&
            deliveryData.address.recipient_name &&
            deliveryData.address.recipient_phone &&
            deliveryData.address.city &&
            deliveryData.address.street_address
          );
        } else if (deliveryData.method === 'pickup_point') {
          return (
            deliveryData.address &&
            deliveryData.address.recipient_name &&
            deliveryData.address.recipient_phone &&
            deliveryData.address.city &&
            deliveryData.office
          );
        } else if (deliveryData.method === 'warehouse_pickup') {
          return (
            deliveryData.address &&
            deliveryData.address.recipient_name &&
            deliveryData.address.recipient_phone
          );
        }
        return false;
      case 2: // Confirmation
        return true;
      default:
        return false;
    }
  };

  const handleNext = () => {
    if (currentStep < STEPS.length - 1) {
      setCurrentStep(currentStep + 1);
    } else {
      // Final step - complete delivery setup
      onDeliveryComplete(deliveryData);
    }
  };

  const handlePrevious = () => {
    if (currentStep > 0) {
      setCurrentStep(currentStep - 1);
    }
  };

  const getStepIcon = (stepIndex: number) => {
    const step = STEPS[stepIndex];
    if (stepIndex < currentStep) {
      return CheckIcon;
    }
    return step.icon;
  };

  const getStepColor = (stepIndex: number) => {
    if (stepIndex < currentStep) {
      return 'bg-success text-success-content';
    } else if (stepIndex === currentStep) {
      return 'bg-primary text-primary-content';
    } else {
      return 'bg-base-300 text-base-content';
    }
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Progress Steps */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body p-6">
          <div className="flex items-center justify-between">
            {STEPS.map((step, index) => {
              const StepIcon = getStepIcon(index);

              return (
                <div key={step.id} className="flex items-center">
                  {/* Step circle */}
                  <div className="flex flex-col items-center">
                    <div
                      className={`
                      w-12 h-12 rounded-full flex items-center justify-center transition-all
                      ${getStepColor(index)}
                    `}
                    >
                      <StepIcon className="w-6 h-6" />
                    </div>
                    <div className="mt-2 text-sm font-medium text-center">
                      {step.title}
                    </div>
                  </div>

                  {/* Connector line */}
                  {index < STEPS.length - 1 && (
                    <div
                      className={`
                      flex-1 h-1 mx-4 transition-all
                      ${index < currentStep ? 'bg-success' : 'bg-base-300'}
                    `}
                    />
                  )}
                </div>
              );
            })}
          </div>
        </div>
      </div>

      {/* Step Content */}
      <div className="min-h-[600px]">
        {currentStep === 0 && (
          <PostExpressDeliverySelector
            onMethodChange={handleMethodChange}
            selectedMethod={deliveryData.method}
            weight={deliveryData.weight}
            insuranceAmount={deliveryData.declaredValue}
            hasCOD={deliveryData.codAmount > 0}
            recipientCity={deliveryData.address?.city}
          />
        )}

        {currentStep === 1 && (
          <div className="space-y-6">
            {/* Address Form */}
            <PostExpressAddressForm
              onAddressChange={handleAddressChange}
              initialAddress={deliveryData.address}
              deliveryMethod={deliveryData.method}
            />

            {/* Office Selector for pickup points */}
            {deliveryData.method === 'pickup_point' && (
              <PostExpressOfficeSelector
                selectedCity={deliveryData.address?.city}
                onOfficeSelect={handleOfficeSelect}
                selectedOffice={deliveryData.office}
              />
            )}

            {/* Rate Calculator */}
            <PostExpressRateCalculator
              onRateCalculated={handleRateCalculated}
              initialParams={{
                weight: deliveryData.weight,
                declaredValue: deliveryData.declaredValue,
                codAmount: deliveryData.codAmount,
                recipientCity: deliveryData.address?.city,
              }}
            />
          </div>
        )}

        {currentStep === 2 && (
          <div className="space-y-6">
            {/* Confirmation Summary */}
            <div className="card bg-base-100 shadow-lg">
              <div className="card-body p-6">
                <h3 className="text-xl font-bold mb-6">
                  Подтверждение доставки
                </h3>

                <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
                  {/* Method and Address */}
                  <div className="space-y-4">
                    <div className="card bg-base-200">
                      <div className="card-body p-4">
                        <h4 className="font-semibold mb-3 flex items-center gap-2">
                          {deliveryData.method === 'courier' && (
                            <TruckIcon className="w-5 h-5" />
                          )}
                          {deliveryData.method === 'pickup_point' && (
                            <MapPinIcon className="w-5 h-5" />
                          )}
                          {deliveryData.method === 'warehouse_pickup' && (
                            <BuildingStorefrontIcon className="w-5 h-5" />
                          )}
                          Способ доставки
                        </h4>
                        <div className="text-lg font-medium">
                          {deliveryData.method === 'courier' &&
                            'Курьерская доставка'}
                          {deliveryData.method === 'pickup_point' &&
                            'Почтовое отделение'}
                          {deliveryData.method === 'warehouse_pickup' &&
                            'Склад Sve Tu'}
                        </div>
                      </div>
                    </div>

                    <div className="card bg-base-200">
                      <div className="card-body p-4">
                        <h4 className="font-semibold mb-3">Получатель</h4>
                        <div className="space-y-1">
                          <div>{deliveryData.address?.recipient_name}</div>
                          <div className="text-sm text-base-content/70">
                            {deliveryData.address?.recipient_phone}
                          </div>
                          {deliveryData.address?.recipient_email && (
                            <div className="text-sm text-base-content/70">
                              {deliveryData.address.recipient_email}
                            </div>
                          )}
                        </div>
                      </div>
                    </div>

                    <div className="card bg-base-200">
                      <div className="card-body p-4">
                        <h4 className="font-semibold mb-3">Адрес доставки</h4>
                        <div className="space-y-1">
                          {deliveryData.method === 'courier' && (
                            <>
                              <div>{deliveryData.address?.street_address}</div>
                              {deliveryData.address?.apartment && (
                                <div className="text-sm">
                                  Кв. {deliveryData.address.apartment}
                                  {deliveryData.address.floor &&
                                    `, ${deliveryData.address.floor} этаж`}
                                  {deliveryData.address.entrance &&
                                    `, подъезд ${deliveryData.address.entrance}`}
                                </div>
                              )}
                              <div>
                                {deliveryData.address?.city},{' '}
                                {deliveryData.address?.postal_code}
                              </div>
                            </>
                          )}

                          {deliveryData.method === 'pickup_point' &&
                            deliveryData.office && (
                              <>
                                <div className="font-medium">
                                  {deliveryData.office.name}
                                </div>
                                <div>{deliveryData.office.address}</div>
                                <div>
                                  {deliveryData.office.city},{' '}
                                  {deliveryData.office.postal_code}
                                </div>
                                <div className="text-sm text-base-content/70">
                                  Код отделения: {deliveryData.office.code}
                                </div>
                              </>
                            )}

                          {deliveryData.method === 'warehouse_pickup' && (
                            <>
                              <div className="font-medium">Склад Sve Tu</div>
                              <div>Микија Манојловића 53</div>
                              <div>Нови Сад, 21000</div>
                              <div className="text-sm text-base-content/70">
                                Пн-Пт 09:00-19:00, Сб 10:00-16:00
                              </div>
                            </>
                          )}
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Cost Breakdown */}
                  <div>
                    {deliveryData.rate && (
                      <div className="card bg-gradient-to-r from-primary/5 to-secondary/5">
                        <div className="card-body p-4">
                          <h4 className="font-semibold mb-4">
                            Стоимость доставки
                          </h4>

                          <div className="space-y-3">
                            <div className="flex justify-between">
                              <span>Базовая стоимость:</span>
                              <span>
                                {deliveryData.rate.base_price?.toFixed(0)} RSD
                              </span>
                            </div>

                            {deliveryData.rate.fuel_surcharge > 0 && (
                              <div className="flex justify-between">
                                <span>Топливная надбавка:</span>
                                <span>
                                  {deliveryData.rate.fuel_surcharge.toFixed(0)}{' '}
                                  RSD
                                </span>
                              </div>
                            )}

                            {deliveryData.rate.cod_fee > 0 && (
                              <div className="flex justify-between">
                                <span>Наложенный платеж:</span>
                                <span>
                                  {deliveryData.rate.cod_fee.toFixed(0)} RSD
                                </span>
                              </div>
                            )}

                            {deliveryData.rate.insurance_fee > 0 && (
                              <div className="flex justify-between">
                                <span>Доп. страхование:</span>
                                <span>
                                  {deliveryData.rate.insurance_fee.toFixed(0)}{' '}
                                  RSD
                                </span>
                              </div>
                            )}

                            <div className="divider my-2"></div>

                            <div className="flex justify-between text-lg font-bold">
                              <span>Итого:</span>
                              <span className="text-primary">
                                {deliveryData.method === 'warehouse_pickup'
                                  ? 'Бесплатно'
                                  : `${deliveryData.rate.total_price?.toFixed(0)} RSD`}
                              </span>
                            </div>

                            <div className="text-sm text-base-content/70 text-center">
                              Ожидаемое время:{' '}
                              {deliveryData.rate.estimated_days || '1-2'}{' '}
                              рабочих дня
                            </div>
                          </div>
                        </div>
                      </div>
                    )}

                    {/* Additional notes */}
                    {deliveryData.address?.note && (
                      <div className="card bg-base-200 mt-4">
                        <div className="card-body p-4">
                          <h4 className="font-semibold mb-2">Комментарий</h4>
                          <div className="text-sm">
                            {deliveryData.address.note}
                          </div>
                        </div>
                      </div>
                    )}
                  </div>
                </div>
              </div>
            </div>
          </div>
        )}
      </div>

      {/* Navigation Buttons */}
      <div className="flex justify-between items-center pt-6 border-t">
        <button
          className={`btn btn-outline ${currentStep === 0 ? 'btn-disabled' : ''}`}
          onClick={handlePrevious}
          disabled={currentStep === 0}
        >
          <ChevronLeftIcon className="w-5 h-5" />
          Назад
        </button>

        <div className="text-sm text-base-content/60">
          Шаг {currentStep + 1} из {STEPS.length}
        </div>

        <button
          className={`btn btn-primary ${!canProceedToNext() ? 'btn-disabled' : ''}`}
          onClick={handleNext}
          disabled={!canProceedToNext()}
        >
          {currentStep === STEPS.length - 1 ? 'Подтвердить доставку' : 'Далее'}
          {currentStep < STEPS.length - 1 && (
            <ChevronRightIcon className="w-5 h-5" />
          )}
        </button>
      </div>
    </div>
  );
}
