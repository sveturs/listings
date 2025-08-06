'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';

interface PaymentDeliveryStepProps {
  onNext: () => void;
  onBack: () => void;
}

export default function PaymentDeliveryStep({
  onNext,
  onBack,
}: PaymentDeliveryStepProps) {
  const t = useTranslations('payment');
  const tCommon = useTranslations('common');
  const tDelivery = useTranslations('delivery');
  const [formData, setFormData] = useState({
    paymentMethods: ['cod'], // cod, bank_transfer, cash
    codPrice: 250, // Стоимость наложенного платежа
    personalMeeting: true,
    deliveryOptions: [] as string[],
    negotiablePrice: false,
    bundleDeals: false,
  });

  const paymentOptions = [
    {
      id: 'cod',
      label: 'payment.cod',
      icon: '📦',
      description: 'payment.cod_desc',
      popular: true,
    },
    {
      id: 'cash',
      label: 'payment.cash',
      icon: '💵',
      description: 'payment.cash_desc',
    },
    {
      id: 'bank_transfer',
      label: 'payment.bank_transfer',
      icon: '🏦',
      description: 'payment.bank_transfer_desc',
    },
  ];

  // TODO: Загружать список курьерских служб и тарифов из API
  // ВРЕМЕННОЕ РЕШЕНИЕ: Хардкодные службы доставки для демонстрации
  // См. TODO #16: Создать API для списка курьерских служб и их тарифов
  const deliveryServices = [
    { id: 'post_srbije', name: 'Пошта Србије', fee: 250 },
    { id: 'aks', name: 'AKS Express', fee: 300 },
    { id: 'bex', name: 'BEX Express', fee: 280 },
    { id: 'city_express', name: 'City Express', fee: 320 },
  ];

  const togglePaymentMethod = (methodId: string) => {
    setFormData((prev) => ({
      ...prev,
      paymentMethods: prev.paymentMethods.includes(methodId)
        ? prev.paymentMethods.filter((m) => m !== methodId)
        : [...prev.paymentMethods, methodId],
    }));
  };

  const calculateCODTotal = (price: number) => {
    return price + formData.codPrice;
  };

  return (
    <div className="max-w-2xl mx-auto">
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4 flex items-center">
            💳 {t('setup_title')}
          </h2>
          <p className="text-base-content/70 mb-6">{t('setup_description')}</p>

          {/* Способы оплаты */}
          <div className="space-y-6">
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  💳 {t('methods_title')}
                </span>
                <span className="label-text-alt text-error">*</span>
              </label>

              <div className="grid gap-3">
                {paymentOptions.map((option) => (
                  <label key={option.id} className="cursor-pointer">
                    <input
                      type="checkbox"
                      checked={formData.paymentMethods.includes(option.id)}
                      onChange={() => togglePaymentMethod(option.id)}
                      className="sr-only"
                    />
                    <div
                      className={`
                      card border-2 transition-all duration-200
                      ${
                        formData.paymentMethods.includes(option.id)
                          ? 'border-primary bg-primary/5'
                          : 'border-base-300 hover:border-primary/50'
                      }
                    `}
                    >
                      <div className="card-body p-4">
                        <div className="flex items-start gap-3">
                          <span className="text-2xl">{option.icon}</span>
                          <div className="flex-1">
                            <div className="flex items-center gap-2">
                              <h3 className="font-medium">{t(option.label)}</h3>
                              {option.popular && (
                                <span className="badge badge-primary badge-sm">
                                  {tCommon('popular')}
                                </span>
                              )}
                            </div>
                            <p className="text-sm text-base-content/60 mt-1">
                              {t(option.description)}
                            </p>
                          </div>
                          {formData.paymentMethods.includes(option.id) && (
                            <svg
                              className="w-6 h-6 text-primary"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                                clipRule="evenodd"
                              />
                            </svg>
                          )}
                        </div>
                      </div>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            {/* Калькулятор наложенного платежа */}
            {formData.paymentMethods.includes('cod') && (
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
                  <h3 className="font-medium">{t('cod_calculator')}</h3>
                  <div className="text-sm mt-2 space-y-1">
                    <p>
                      • {t('cod_fee')}: {formData.codPrice} РСД
                    </p>
                    <p>
                      • {t('example')}: 1.000 РСД + {formData.codPrice} РСД ={' '}
                      {calculateCODTotal(1000)} РСД
                    </p>
                  </div>
                </div>
              </div>
            )}

            {/* Лична предаја */}
            <div className="form-control">
              <label className="label cursor-pointer">
                <div className="flex items-center gap-3">
                  <span className="text-2xl">🤝</span>
                  <div>
                    <span className="label-text font-medium">
                      {tDelivery('personal_handover')}
                    </span>
                    <p className="text-sm text-base-content/60">
                      {tDelivery('personal_handover_desc')}
                    </p>
                  </div>
                </div>
                <input
                  type="checkbox"
                  className="toggle toggle-primary"
                  checked={formData.personalMeeting}
                  onChange={(e) =>
                    setFormData((prev) => ({
                      ...prev,
                      personalMeeting: e.target.checked,
                    }))
                  }
                />
              </label>
            </div>

            {/* Опции доставки */}
            <div className="form-control">
              <label className="label">
                <span className="label-text font-medium">
                  🚚 {tDelivery('services')}
                </span>
              </label>

              <div className="space-y-2">
                {deliveryServices.map((service) => (
                  <label key={service.id} className="cursor-pointer">
                    <input
                      type="checkbox"
                      checked={formData.deliveryOptions.includes(service.id)}
                      onChange={(e) => {
                        const newOptions = e.target.checked
                          ? [...formData.deliveryOptions, service.id]
                          : formData.deliveryOptions.filter(
                              (o) => o !== service.id
                            );
                        setFormData((prev) => ({
                          ...prev,
                          deliveryOptions: newOptions,
                        }));
                      }}
                      className="sr-only"
                    />
                    <div
                      className={`
                      flex items-center justify-between p-3 rounded-lg border-2 transition-all
                      ${
                        formData.deliveryOptions.includes(service.id)
                          ? 'border-primary bg-primary/5'
                          : 'border-base-300 hover:border-primary/50'
                      }
                    `}
                    >
                      <div className="flex items-center gap-3">
                        <span className="text-lg">📦</span>
                        <span className="font-medium">{service.name}</span>
                      </div>
                      <div className="flex items-center gap-2">
                        <span className="text-sm text-base-content/70">
                          {service.fee} РСД
                        </span>
                        {formData.deliveryOptions.includes(service.id) && (
                          <svg
                            className="w-5 h-5 text-primary"
                            fill="currentColor"
                            viewBox="0 0 20 20"
                          >
                            <path
                              fillRule="evenodd"
                              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                              clipRule="evenodd"
                            />
                          </svg>
                        )}
                      </div>
                    </div>
                  </label>
                ))}
              </div>
            </div>

            {/* Дополнительные опции */}
            <div className="space-y-3">
              <div className="form-control">
                <label className="label cursor-pointer">
                  <div className="flex items-center gap-3">
                    <span className="text-xl">💬</span>
                    <span className="label-text">{t('negotiable_price')}</span>
                  </div>
                  <input
                    type="checkbox"
                    className="toggle toggle-sm"
                    checked={formData.negotiablePrice}
                    onChange={(e) =>
                      setFormData((prev) => ({
                        ...prev,
                        negotiablePrice: e.target.checked,
                      }))
                    }
                  />
                </label>
              </div>

              <div className="form-control">
                <label className="label cursor-pointer">
                  <div className="flex items-center gap-3">
                    <span className="text-xl">📦</span>
                    <span className="label-text">{t('bundle_deals')}</span>
                  </div>
                  <input
                    type="checkbox"
                    className="toggle toggle-sm"
                    checked={formData.bundleDeals}
                    onChange={(e) =>
                      setFormData((prev) => ({
                        ...prev,
                        bundleDeals: e.target.checked,
                      }))
                    }
                  />
                </label>
              </div>
            </div>
          </div>

          {/* Кнопки навигации */}
          <div className="card-actions justify-between mt-6">
            <button className="btn btn-outline" onClick={onBack}>
              ← {tCommon('back')}
            </button>
            <button
              className="btn btn-primary"
              onClick={onNext}
              disabled={formData.paymentMethods.length === 0}
            >
              {tCommon('continue')} →
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
