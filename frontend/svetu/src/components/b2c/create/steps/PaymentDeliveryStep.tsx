'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useCreateStorefrontContext } from '@/contexts/CreateB2CStoreContext';

interface PaymentDeliveryStepProps {
  onNext: () => void;
  onBack: () => void;
}

const paymentMethods = [
  { id: 'cash', icon: '💵' },
  { id: 'card', icon: '💳' },
  { id: 'bank_transfer', icon: '🏦' },
  { id: 'cod', icon: '📦' },
  { id: 'postanska_stedionica', icon: '📮' },
];

const deliveryProviders = [
  'posta_srbije',
  'aks',
  'bex',
  'city_express',
  'd_express',
  'post_express',
  'daily_express',
  'own_delivery',
];

interface DeliveryOption {
  providerName: string;
  deliveryTimeMinutes: number;
  deliveryCostRSD: number;
  freeDeliveryThresholdRSD?: number;
  maxDistanceKm?: number;
}

export default function PaymentDeliveryStep({
  onNext,
  onBack,
}: PaymentDeliveryStepProps) {
  const t = useTranslations('create_storefront');
  const tCommon = useTranslations('common');
  const { formData, updateFormData } = useCreateStorefrontContext();
  const [newDeliveryOption, setNewDeliveryOption] = useState<DeliveryOption>({
    providerName: '',
    deliveryTimeMinutes: 60,
    deliveryCostRSD: 300,
  });

  const togglePaymentMethod = (method: string) => {
    const currentMethods = formData.paymentMethods || [];
    const updated = currentMethods.includes(method)
      ? currentMethods.filter((m) => m !== method)
      : [...currentMethods, method];
    updateFormData({ paymentMethods: updated });
  };

  const addDeliveryOption = () => {
    if (newDeliveryOption.providerName) {
      const currentOptions = formData.deliveryOptions || [];
      updateFormData({
        deliveryOptions: [...currentOptions, { ...newDeliveryOption }],
      });
      setNewDeliveryOption({
        providerName: '',
        deliveryTimeMinutes: 60,
        deliveryCostRSD: 300,
      });
    }
  };

  const removeDeliveryOption = (index: number) => {
    const currentOptions = formData.deliveryOptions || [];
    updateFormData({
      deliveryOptions: currentOptions.filter((_, i) => i !== index),
    });
  };

  const handleNext = () => {
    onNext();
  };

  return (
    <div className="max-w-3xl mx-auto">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">{t('title')}</h2>
          <p className="text-base-content/70 mb-6">{t('subtitle')}</p>

          {/* Payment Methods */}
          <div className="mb-8">
            <h3 className="text-lg font-semibold mb-4">
              {t('payment_methods_title')}
            </h3>
            <div className="grid grid-cols-2 md:grid-cols-3 gap-3">
              {paymentMethods.map((method) => (
                <label key={method.id} className="cursor-pointer">
                  <input
                    type="checkbox"
                    className="checkbox checkbox-primary"
                    checked={
                      formData.paymentMethods?.includes(method.id) || false
                    }
                    onChange={() => togglePaymentMethod(method.id)}
                  />
                  <span className="ml-2">
                    {method.icon} {t(`payment_methods.${method.id}`)}
                  </span>
                </label>
              ))}
            </div>
          </div>

          {/* Delivery Options */}
          <div>
            <h3 className="text-lg font-semibold mb-4">
              {t('delivery_options_title')}
            </h3>

            {/* Existing delivery options */}
            {formData.deliveryOptions &&
              formData.deliveryOptions.length > 0 && (
                <div className="space-y-2 mb-4">
                  {formData.deliveryOptions.map((option, index) => (
                    <div
                      key={index}
                      className="flex items-center justify-between p-3 bg-base-200 rounded-lg"
                    >
                      <div>
                        <span className="font-medium">
                          {t(`delivery_providers.${option.providerName}`)}
                        </span>
                        <span className="text-sm text-base-content/70 ml-2">
                          {option.deliveryTimeMinutes} min •{' '}
                          {option.deliveryCostRSD} RSD
                        </span>
                      </div>
                      <button
                        className="btn btn-ghost btn-sm"
                        onClick={() => removeDeliveryOption(index)}
                      >
                        ✕
                      </button>
                    </div>
                  ))}
                </div>
              )}

            {/* Add new delivery option */}
            <div className="card bg-base-200">
              <div className="card-body p-4">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">{t('provider')}</span>
                    </label>
                    <select
                      className="select select-bordered select-sm"
                      value={newDeliveryOption.providerName}
                      onChange={(e) =>
                        setNewDeliveryOption({
                          ...newDeliveryOption,
                          providerName: e.target.value,
                        })
                      }
                    >
                      <option value="">{tCommon('select')}</option>
                      {deliveryProviders.map((provider) => (
                        <option key={provider} value={provider}>
                          {t(`delivery_providers.${provider}`)}
                        </option>
                      ))}
                    </select>
                  </div>

                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">{t('delivery_time')}</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered input-sm"
                      value={newDeliveryOption.deliveryTimeMinutes}
                      onChange={(e) =>
                        setNewDeliveryOption({
                          ...newDeliveryOption,
                          deliveryTimeMinutes: parseInt(e.target.value) || 0,
                        })
                      }
                    />
                  </div>

                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">{t('delivery_cost')}</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered input-sm"
                      value={newDeliveryOption.deliveryCostRSD}
                      onChange={(e) =>
                        setNewDeliveryOption({
                          ...newDeliveryOption,
                          deliveryCostRSD: parseInt(e.target.value) || 0,
                        })
                      }
                    />
                  </div>

                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">{t('free_threshold')}</span>
                    </label>
                    <input
                      type="number"
                      className="input input-bordered input-sm"
                      placeholder={tCommon('optional')}
                      value={newDeliveryOption.freeDeliveryThresholdRSD || ''}
                      onChange={(e) =>
                        setNewDeliveryOption({
                          ...newDeliveryOption,
                          freeDeliveryThresholdRSD: e.target.value
                            ? parseInt(e.target.value)
                            : undefined,
                        })
                      }
                    />
                  </div>
                </div>

                <button
                  className="btn btn-primary btn-sm mt-3"
                  onClick={addDeliveryOption}
                  disabled={!newDeliveryOption.providerName}
                >
                  {tCommon('add')}
                </button>
              </div>
            </div>
          </div>

          <div className="card-actions justify-between mt-6">
            <button className="btn btn-ghost" onClick={onBack}>
              {tCommon('back')}
            </button>
            <button className="btn btn-primary" onClick={handleNext}>
              {tCommon('next')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
