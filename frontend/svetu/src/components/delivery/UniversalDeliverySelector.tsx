'use client';

import { useState, useEffect } from 'react';
import {
  TruckIcon,
  ClockIcon,
  CheckIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  ArrowPathIcon,
  ShieldCheckIcon,
  MapPinIcon,
} from '@heroicons/react/24/outline';
import { useTranslations } from 'next-intl';
import {
  DeliveryQuote,
  CalculationRequest,
  CalculationResponse,
} from '@/types/delivery';
import configManager from '@/config';

interface Props {
  calculationRequest: CalculationRequest;
  onQuoteSelected?: (quote: DeliveryQuote) => void;
  selectedQuoteId?: number;
  autoCalculate?: boolean;
  showComparison?: boolean;
  className?: string;
}

const PROVIDER_LOGOS: Record<string, string> = {
  'post-express': '/images/delivery/post-express-logo.png',
  'bex-express': '/images/delivery/bex-logo.png',
  'aks-express': '/images/delivery/aks-logo.png',
  'd-express': '/images/delivery/d-express-logo.png',
  'city-express': '/images/delivery/city-express-logo.png',
  dhl: '/images/delivery/dhl-logo.png',
};

export default function UniversalDeliverySelector({
  calculationRequest,
  onQuoteSelected,
  selectedQuoteId,
  autoCalculate = true,
  showComparison = true,
  className = '',
}: Props) {
  const _t = useTranslations('delivery');
  const [quotes, setQuotes] = useState<DeliveryQuote[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [calculationResponse, setCalculationResponse] =
    useState<CalculationResponse | null>(null);

  useEffect(() => {
    if (autoCalculate && calculationRequest.items.length > 0) {
      calculateRates();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [calculationRequest, autoCalculate]);

  const calculateRates = async () => {
    setLoading(true);
    setError(null);

    try {
      const apiUrl = configManager.getApiUrl();
      const response = await fetch(
        `${apiUrl}/api/v1/delivery/calculate-universal`,
        {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify(calculationRequest),
        }
      );

      const data: CalculationResponse = await response.json();

      if (data.success && data.data) {
        setQuotes(data.data.providers);
        setCalculationResponse(data);

        // Auto-select recommended quote if no quote is selected
        if (!selectedQuoteId && data.data.recommended) {
          onQuoteSelected?.(data.data.recommended);
        }
      } else {
        setError(data.message || 'Не удалось рассчитать стоимость доставки');
      }
    } catch (err) {
      setError('Ошибка при расчете стоимости доставки');
      console.error('Rate calculation error:', err);
    } finally {
      setLoading(false);
    }
  };

  const formatDeliveryTime = (days: number) => {
    if (days === 1) return '1 рабочий день';
    if (days < 5) return `${days} рабочих дня`;
    return `${days} рабочих дней`;
  };

  const getQuoteStatus = (quote: DeliveryQuote) => {
    if (
      calculationResponse?.data?.recommended?.provider_id === quote.provider_id
    ) {
      return { type: 'recommended', label: 'Рекомендуем' };
    }
    if (
      calculationResponse?.data?.cheapest?.provider_id === quote.provider_id
    ) {
      return { type: 'cheapest', label: 'Самый дешевый' };
    }
    if (calculationResponse?.data?.fastest?.provider_id === quote.provider_id) {
      return { type: 'fastest', label: 'Самый быстрый' };
    }
    return null;
  };

  const handleQuoteSelect = (quote: DeliveryQuote) => {
    onQuoteSelected?.(quote);
  };

  if (loading) {
    return (
      <div className={`card bg-base-100 shadow-lg ${className}`}>
        <div className="card-body p-6 text-center">
          <ArrowPathIcon className="w-12 h-12 mx-auto text-primary animate-spin mb-4" />
          <h3 className="text-lg font-semibold mb-2">
            Расчет стоимости доставки
          </h3>
          <p className="text-base-content/60">
            Сравниваем предложения от всех провайдеров...
          </p>
          <div className="mt-4 space-y-2">
            <div className="skeleton h-4 w-3/4 mx-auto"></div>
            <div className="skeleton h-4 w-1/2 mx-auto"></div>
          </div>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className={`card bg-base-100 shadow-lg ${className}`}>
        <div className="card-body p-6">
          <div className="alert alert-error">
            <ExclamationTriangleIcon className="w-5 h-5" />
            <div>
              <div className="font-semibold">Ошибка расчета</div>
              <div className="text-sm">{error}</div>
            </div>
            <button className="btn btn-sm btn-outline" onClick={calculateRates}>
              Повторить
            </button>
          </div>
        </div>
      </div>
    );
  }

  if (quotes.length === 0 && !loading) {
    return (
      <div className={`card bg-base-100 shadow-lg ${className}`}>
        <div className="card-body p-6 text-center">
          <TruckIcon className="w-16 h-16 mx-auto text-base-content/30 mb-4" />
          <h3 className="text-lg font-semibold mb-2">
            Нет доступных вариантов доставки
          </h3>
          <p className="text-base-content/60 mb-4">
            Для указанного маршрута не найдено подходящих вариантов доставки
          </p>
          <button className="btn btn-primary" onClick={calculateRates}>
            <ArrowPathIcon className="w-4 h-4" />
            Пересчитать
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <div className="p-2 bg-primary/10 rounded-lg">
            <TruckIcon className="w-6 h-6 text-primary" />
          </div>
          <div>
            <h3 className="text-lg font-semibold">Выберите способ доставки</h3>
            <p className="text-sm text-base-content/60">
              Найдено {quotes.length} вариантов доставки
            </p>
          </div>
        </div>

        <button
          className="btn btn-sm btn-ghost"
          onClick={calculateRates}
          disabled={loading}
        >
          <ArrowPathIcon className="w-4 h-4" />
          Обновить
        </button>
      </div>

      {/* Quotes List */}
      <div className="space-y-4">
        {quotes.map((quote) => {
          const isSelected = selectedQuoteId === quote.provider_id;
          const status = getQuoteStatus(quote);

          return (
            <div
              key={quote.provider_id}
              className={`
                card cursor-pointer transition-all duration-200 border-2
                ${
                  isSelected
                    ? 'border-primary bg-primary/5 shadow-lg'
                    : 'border-transparent hover:border-primary/30 hover:shadow-md'
                }
              `}
              onClick={() => handleQuoteSelect(quote)}
            >
              <div className="card-body p-4">
                <div className="flex items-start justify-between">
                  <div className="flex items-start gap-4 flex-1">
                    {/* Provider Logo */}
                    <div className="flex-shrink-0">
                      {PROVIDER_LOGOS[quote.provider_code] ? (
                        <img
                          src={PROVIDER_LOGOS[quote.provider_code]}
                          alt={quote.provider_name}
                          className="w-12 h-12 object-contain rounded-lg bg-white p-1"
                        />
                      ) : (
                        <div className="w-12 h-12 bg-base-200 rounded-lg flex items-center justify-center">
                          <TruckIcon className="w-6 h-6 text-base-content/60" />
                        </div>
                      )}
                    </div>

                    {/* Quote Details */}
                    <div className="flex-1 space-y-2">
                      <div className="flex items-center gap-2">
                        <h4 className="font-semibold text-lg">
                          {quote.provider_name}
                        </h4>
                        {status && (
                          <span
                            className={`
                            badge badge-sm text-xs
                            ${status.type === 'recommended' ? 'badge-primary' : ''}
                            ${status.type === 'cheapest' ? 'badge-success' : ''}
                            ${status.type === 'fastest' ? 'badge-secondary' : ''}
                          `}
                          >
                            {status.label}
                          </span>
                        )}
                      </div>

                      <div className="flex items-center gap-4 text-sm text-base-content/70">
                        <div className="flex items-center gap-1">
                          <ClockIcon className="w-4 h-4" />
                          <span>
                            {formatDeliveryTime(quote.estimated_days)}
                          </span>
                        </div>
                        {quote.estimated_delivery_date && (
                          <div className="flex items-center gap-1">
                            <MapPinIcon className="w-4 h-4" />
                            <span>
                              до{' '}
                              {new Date(
                                quote.estimated_delivery_date
                              ).toLocaleDateString()}
                            </span>
                          </div>
                        )}
                      </div>

                      {/* Services */}
                      {quote.services.length > 0 && (
                        <div className="flex flex-wrap gap-2">
                          {quote.services
                            .filter((s) => s.included)
                            .map((service, idx) => (
                              <div
                                key={idx}
                                className="flex items-center gap-1 text-xs text-success"
                              >
                                <ShieldCheckIcon className="w-3 h-3" />
                                <span>{service.name}</span>
                              </div>
                            ))}
                        </div>
                      )}
                    </div>
                  </div>

                  {/* Price and Selection */}
                  <div className="text-right space-y-2">
                    <div className="text-2xl font-bold text-primary">
                      {quote.total_price.toFixed(0)} RSD
                    </div>
                    <div className="text-xs text-base-content/60">
                      доставка {quote.delivery_cost.toFixed(0)} RSD
                      {quote.insurance_cost && quote.insurance_cost > 0 && (
                        <> + страховка {quote.insurance_cost.toFixed(0)} RSD</>
                      )}
                    </div>
                    {isSelected && (
                      <div className="flex items-center gap-1 text-success text-sm">
                        <CheckIcon className="w-4 h-4" />
                        <span>Выбрано</span>
                      </div>
                    )}
                  </div>
                </div>

                {/* Expanded details for selected quote */}
                {isSelected && showComparison && (
                  <div className="mt-4 pt-4 border-t border-base-200 space-y-3">
                    <h5 className="font-medium text-sm">
                      Детализация стоимости:
                    </h5>

                    <div className="grid grid-cols-2 gap-4 text-sm">
                      <div className="space-y-2">
                        <div className="flex justify-between">
                          <span>Базовая стоимость:</span>
                          <span>
                            {quote.cost_breakdown.base_price.toFixed(0)} RSD
                          </span>
                        </div>
                        {quote.cost_breakdown.weight_surcharge &&
                          quote.cost_breakdown.weight_surcharge > 0 && (
                            <div className="flex justify-between">
                              <span>Надбавка за вес:</span>
                              <span>
                                {quote.cost_breakdown.weight_surcharge.toFixed(
                                  0
                                )}{' '}
                                RSD
                              </span>
                            </div>
                          )}
                        {quote.cost_breakdown.fragile_surcharge &&
                          quote.cost_breakdown.fragile_surcharge > 0 && (
                            <div className="flex justify-between">
                              <span>За хрупкость:</span>
                              <span>
                                {quote.cost_breakdown.fragile_surcharge.toFixed(
                                  0
                                )}{' '}
                                RSD
                              </span>
                            </div>
                          )}
                      </div>

                      <div className="space-y-2">
                        {quote.cost_breakdown.oversized_surcharge &&
                          quote.cost_breakdown.oversized_surcharge > 0 && (
                            <div className="flex justify-between">
                              <span>За габариты:</span>
                              <span>
                                {quote.cost_breakdown.oversized_surcharge.toFixed(
                                  0
                                )}{' '}
                                RSD
                              </span>
                            </div>
                          )}
                        {quote.insurance_cost && quote.insurance_cost > 0 && (
                          <div className="flex justify-between">
                            <span>Страхование:</span>
                            <span>{quote.insurance_cost.toFixed(0)} RSD</span>
                          </div>
                        )}
                        {quote.cod_fee && quote.cod_fee > 0 && (
                          <div className="flex justify-between">
                            <span>Комиссия COD:</span>
                            <span>{quote.cod_fee.toFixed(0)} RSD</span>
                          </div>
                        )}
                      </div>
                    </div>
                  </div>
                )}
              </div>
            </div>
          );
        })}
      </div>

      {/* Summary Information */}
      {calculationResponse?.data && (
        <div className="card bg-gradient-to-r from-info/5 to-info/10">
          <div className="card-body p-4">
            <div className="flex items-start gap-3">
              <InformationCircleIcon className="w-5 h-5 text-info flex-shrink-0 mt-0.5" />
              <div className="space-y-2 text-sm">
                <div className="font-semibold">Информация о расчете:</div>
                <div className="space-y-1">
                  <div>
                    • Сравнены предложения от {quotes.length} провайдеров
                  </div>
                  <div>
                    • Разница в цене:{' '}
                    {(
                      Math.max(...quotes.map((q) => q.total_price)) -
                      Math.min(...quotes.map((q) => q.total_price))
                    ).toFixed(0)}{' '}
                    RSD
                  </div>
                  <div>
                    • Разница во времени доставки:{' '}
                    {Math.max(...quotes.map((q) => q.estimated_days)) -
                      Math.min(...quotes.map((q) => q.estimated_days))}{' '}
                    дней
                  </div>
                  {calculationResponse.data.recommended && (
                    <div className="text-primary">
                      • Рекомендуем{' '}
                      <strong>
                        {calculationResponse.data.recommended.provider_name}
                      </strong>{' '}
                      как оптимальный вариант
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
