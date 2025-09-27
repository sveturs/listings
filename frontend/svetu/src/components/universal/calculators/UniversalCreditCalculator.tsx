'use client';

import { FC, useState, useEffect, useMemo } from 'react';
import { useTranslations } from 'next-intl';
import {
  FaCalculator,
  FaPercent,
  FaCalendarAlt,
  FaMoneyBillWave,
  FaInfoCircle,
} from 'react-icons/fa';

// Конфигурация для разных категорий
export interface CreditConfig {
  category: string;
  type: 'credit' | 'installment' | 'mortgage' | 'leasing';
  minDownPayment: number; // В процентах от цены
  maxDownPayment: number;
  defaultDownPayment: number;
  minTerm: number; // В месяцах
  maxTerm: number;
  defaultTerm: number;
  minRate: number; // Годовая процентная ставка
  maxRate: number;
  defaultRate: number;
  // Банки и программы
  banks?: {
    id: string;
    name: string;
    logo?: string;
    rates: {
      term: number;
      rate: number;
      conditions?: string;
    }[];
  }[];
  // Дополнительные расходы
  additionalCosts?: {
    id: string;
    label: string;
    type: 'fixed' | 'percentage';
    value: number;
    required?: boolean;
    description?: string;
  }[];
  // Специальные предложения
  specialOffers?: {
    id: string;
    label: string;
    description: string;
    discount?: number;
    conditions?: string;
  }[];
}

// Предустановленные конфигурации
const CREDIT_CONFIGS: Record<string, CreditConfig> = {
  cars: {
    category: 'cars',
    type: 'credit',
    minDownPayment: 0,
    maxDownPayment: 90,
    defaultDownPayment: 20,
    minTerm: 12,
    maxTerm: 84,
    defaultTerm: 60,
    minRate: 5.9,
    maxRate: 24.9,
    defaultRate: 9.9,
    banks: [
      {
        id: 'bank1',
        name: 'AutoBank',
        rates: [
          { term: 12, rate: 5.9 },
          { term: 24, rate: 7.9 },
          { term: 36, rate: 8.9 },
          { term: 60, rate: 9.9 },
          { term: 84, rate: 11.9 },
        ],
      },
      {
        id: 'bank2',
        name: 'CarCredit',
        rates: [
          { term: 12, rate: 6.9 },
          { term: 36, rate: 9.9 },
          { term: 60, rate: 10.9 },
        ],
      },
    ],
    additionalCosts: [
      {
        id: 'insurance',
        label: 'CASCO Insurance',
        type: 'percentage',
        value: 5,
        required: true,
      },
      {
        id: 'registration',
        label: 'Registration',
        type: 'fixed',
        value: 500,
        required: true,
      },
      { id: 'evaluation', label: 'Car Evaluation', type: 'fixed', value: 100 },
    ],
  },
  real_estate: {
    category: 'real_estate',
    type: 'mortgage',
    minDownPayment: 10,
    maxDownPayment: 80,
    defaultDownPayment: 20,
    minTerm: 60,
    maxTerm: 360,
    defaultTerm: 240,
    minRate: 4.5,
    maxRate: 12.5,
    defaultRate: 6.5,
    banks: [
      {
        id: 'mortgage1',
        name: 'HomeBank',
        rates: [
          { term: 60, rate: 4.5 },
          { term: 120, rate: 5.5 },
          { term: 180, rate: 6.0 },
          { term: 240, rate: 6.5 },
          { term: 360, rate: 7.5 },
        ],
      },
    ],
    additionalCosts: [
      {
        id: 'insurance',
        label: 'Property Insurance',
        type: 'percentage',
        value: 0.3,
        required: true,
      },
      {
        id: 'evaluation',
        label: 'Property Evaluation',
        type: 'fixed',
        value: 300,
        required: true,
      },
      { id: 'notary', label: 'Notary Services', type: 'fixed', value: 500 },
      {
        id: 'commission',
        label: 'Bank Commission',
        type: 'percentage',
        value: 1,
      },
    ],
  },
  electronics: {
    category: 'electronics',
    type: 'installment',
    minDownPayment: 0,
    maxDownPayment: 50,
    defaultDownPayment: 0,
    minTerm: 3,
    maxTerm: 24,
    defaultTerm: 12,
    minRate: 0,
    maxRate: 0,
    defaultRate: 0,
    specialOffers: [
      {
        id: 'zero',
        label: '0% Installment',
        description: 'No interest for up to 12 months',
      },
      {
        id: 'cashback',
        label: '5% Cashback',
        description: 'Get 5% back on your purchase',
        discount: 5,
      },
    ],
  },
};

interface UniversalCreditCalculatorProps {
  price: number;
  category: string;
  config?: Partial<CreditConfig>;
  onApply?: (calculation: CreditCalculation) => void;
  className?: string;
}

export interface CreditCalculation {
  price: number;
  downPayment: number;
  loanAmount: number;
  term: number;
  rate: number;
  monthlyPayment: number;
  totalPayment: number;
  totalInterest: number;
  additionalCosts?: number;
  selectedBank?: string;
  paymentSchedule?: {
    month: number;
    payment: number;
    principal: number;
    interest: number;
    balance: number;
  }[];
}

const UniversalCreditCalculator: FC<UniversalCreditCalculatorProps> = ({
  price,
  category,
  config: customConfig,
  onApply,
  className = '',
}) => {
  const t = useTranslations('calculator');

  // Получаем конфигурацию для категории
  const defaultConfig = CREDIT_CONFIGS[category] || CREDIT_CONFIGS.electronics;
  const config = { ...defaultConfig, ...customConfig };

  // Состояние калькулятора
  const [downPaymentPercent, setDownPaymentPercent] = useState(
    config.defaultDownPayment
  );
  const [term, setTerm] = useState(config.defaultTerm);
  const [rate, setRate] = useState(config.defaultRate);
  const [selectedBank, setSelectedBank] = useState<string | null>(null);
  const [selectedCosts, setSelectedCosts] = useState<string[]>([]);
  const [showSchedule, setShowSchedule] = useState(false);

  // Расчет платежей
  const calculation = useMemo((): CreditCalculation => {
    const downPayment = (price * downPaymentPercent) / 100;
    const loanAmount = price - downPayment;

    if (loanAmount <= 0) {
      return {
        price,
        downPayment: price,
        loanAmount: 0,
        term,
        rate,
        monthlyPayment: 0,
        totalPayment: price,
        totalInterest: 0,
      };
    }

    let monthlyPayment = 0;
    let totalPayment = 0;
    let totalInterest = 0;

    if (rate === 0) {
      // Беспроцентная рассрочка
      monthlyPayment = loanAmount / term;
      totalPayment = loanAmount;
      totalInterest = 0;
    } else {
      // Аннуитетный платеж
      const monthlyRate = rate / 100 / 12;
      const factor = Math.pow(1 + monthlyRate, term);
      monthlyPayment = (loanAmount * monthlyRate * factor) / (factor - 1);
      totalPayment = monthlyPayment * term;
      totalInterest = totalPayment - loanAmount;
    }

    // Дополнительные расходы
    let additionalCosts = 0;
    if (config.additionalCosts) {
      config.additionalCosts.forEach((cost) => {
        if (cost.required || selectedCosts.includes(cost.id)) {
          if (cost.type === 'fixed') {
            additionalCosts += cost.value;
          } else {
            additionalCosts += (price * cost.value) / 100;
          }
        }
      });
    }

    // График платежей
    const paymentSchedule = [];
    let balance = loanAmount;
    const monthlyRate = rate / 100 / 12;

    for (let month = 1; month <= term; month++) {
      const interestPayment = balance * monthlyRate;
      const principalPayment = monthlyPayment - interestPayment;
      balance -= principalPayment;

      paymentSchedule.push({
        month,
        payment: monthlyPayment,
        principal: principalPayment,
        interest: interestPayment,
        balance: Math.max(0, balance),
      });
    }

    return {
      price,
      downPayment,
      loanAmount,
      term,
      rate,
      monthlyPayment,
      totalPayment: totalPayment + downPayment,
      totalInterest,
      additionalCosts,
      selectedBank: selectedBank ?? undefined,
      paymentSchedule,
    };
  }, [
    price,
    downPaymentPercent,
    term,
    rate,
    selectedCosts,
    selectedBank,
    config.additionalCosts,
  ]);

  // Обновление ставки при выборе банка
  useEffect(() => {
    if (selectedBank && config.banks) {
      const bank = config.banks.find((b) => b.id === selectedBank);
      if (bank) {
        const rateInfo = bank.rates.find((r) => r.term === term);
        if (rateInfo) {
          setRate(rateInfo.rate);
        }
      }
    }
  }, [selectedBank, term, config.banks]);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'EUR',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0,
    }).format(amount);
  };

  const handleApply = () => {
    if (onApply) {
      onApply(calculation);
    }
  };

  return (
    <div className={`card bg-base-100 shadow-lg ${className}`}>
      <div className="card-body">
        <h3 className="card-title flex items-center gap-2">
          <FaCalculator />
          {config.type === 'mortgage'
            ? t('mortgageCalculator')
            : config.type === 'installment'
              ? t('installmentCalculator')
              : config.type === 'leasing'
                ? t('leasingCalculator')
                : t('creditCalculator')}
        </h3>

        {/* Основные параметры */}
        <div className="space-y-4">
          {/* Первоначальный взнос */}
          <div>
            <label className="label">
              <span className="label-text">{t('downPayment')}</span>
              <span className="label-text-alt">
                {formatCurrency(calculation.downPayment)} ({downPaymentPercent}
                %)
              </span>
            </label>
            <input
              type="range"
              className="range range-primary"
              min={config.minDownPayment}
              max={config.maxDownPayment}
              step={5}
              value={downPaymentPercent}
              onChange={(e) => setDownPaymentPercent(Number(e.target.value))}
            />
            <div className="flex justify-between text-xs mt-1">
              <span>{config.minDownPayment}%</span>
              <span>{config.maxDownPayment}%</span>
            </div>
          </div>

          {/* Срок кредита */}
          <div>
            <label className="label">
              <span className="label-text">{t('loanTerm')}</span>
              <span className="label-text-alt">
                {term} {t('months')}
                {term >= 12 && ` (${Math.floor(term / 12)} ${t('years')})`}
              </span>
            </label>
            <input
              type="range"
              className="range range-primary"
              min={config.minTerm}
              max={config.maxTerm}
              step={config.minTerm < 12 ? 1 : 12}
              value={term}
              onChange={(e) => setTerm(Number(e.target.value))}
            />
            <div className="flex justify-between text-xs mt-1">
              <span>
                {config.minTerm} {t('months')}
              </span>
              <span>
                {config.maxTerm} {t('months')}
              </span>
            </div>
          </div>

          {/* Процентная ставка (если не рассрочка) */}
          {config.type !== 'installment' && (
            <div>
              <label className="label">
                <span className="label-text">{t('interestRate')}</span>
                <span className="label-text-alt">
                  {rate}% {t('perYear')}
                </span>
              </label>
              <input
                type="range"
                className="range range-primary"
                min={config.minRate}
                max={config.maxRate}
                step={0.1}
                value={rate}
                onChange={(e) => setRate(Number(e.target.value))}
                disabled={selectedBank !== null}
              />
              <div className="flex justify-between text-xs mt-1">
                <span>{config.minRate}%</span>
                <span>{config.maxRate}%</span>
              </div>
            </div>
          )}

          {/* Выбор банка */}
          {config.banks && config.banks.length > 0 && (
            <div>
              <label className="label">
                <span className="label-text">{t('selectBank')}</span>
              </label>
              <select
                className="select select-bordered w-full"
                value={selectedBank || ''}
                onChange={(e) => setSelectedBank(e.target.value || null)}
              >
                <option value="">{t('customRate')}</option>
                {config.banks.map((bank) => (
                  <option key={bank.id} value={bank.id}>
                    {bank.name}
                  </option>
                ))}
              </select>
            </div>
          )}

          {/* Дополнительные расходы */}
          {config.additionalCosts && config.additionalCosts.length > 0 && (
            <div>
              <label className="label">
                <span className="label-text">{t('additionalCosts')}</span>
              </label>
              <div className="space-y-2">
                {config.additionalCosts.map((cost) => (
                  <label
                    key={cost.id}
                    className="flex items-start gap-2 cursor-pointer"
                  >
                    <input
                      type="checkbox"
                      className="checkbox checkbox-sm mt-0.5"
                      checked={cost.required || selectedCosts.includes(cost.id)}
                      disabled={cost.required}
                      onChange={(e) => {
                        if (e.target.checked) {
                          setSelectedCosts([...selectedCosts, cost.id]);
                        } else {
                          setSelectedCosts(
                            selectedCosts.filter((id) => id !== cost.id)
                          );
                        }
                      }}
                    />
                    <div className="flex-1">
                      <span className="text-sm">
                        {cost.label}
                        {cost.required && (
                          <span className="text-error ml-1">*</span>
                        )}
                      </span>
                      <span className="text-xs text-base-content/60 block">
                        {cost.type === 'fixed'
                          ? formatCurrency(cost.value)
                          : `${cost.value}% ${t('ofPrice')}`}
                        {cost.description && ` - ${cost.description}`}
                      </span>
                    </div>
                  </label>
                ))}
              </div>
            </div>
          )}

          {/* Специальные предложения */}
          {config.specialOffers && config.specialOffers.length > 0 && (
            <div className="alert alert-info">
              <FaInfoCircle />
              <div>
                <h4 className="font-semibold">{t('specialOffers')}</h4>
                {config.specialOffers.map((offer) => (
                  <div key={offer.id} className="text-sm mt-1">
                    <span className="font-medium">{offer.label}:</span>{' '}
                    {offer.description}
                    {offer.discount && ` (-${offer.discount}%)`}
                  </div>
                ))}
              </div>
            </div>
          )}
        </div>

        {/* Результаты расчета */}
        <div className="divider"></div>

        <div className="space-y-3">
          <div className="flex justify-between">
            <span className="text-base-content/70">{t('loanAmount')}:</span>
            <span className="font-semibold">
              {formatCurrency(calculation.loanAmount)}
            </span>
          </div>

          <div className="flex justify-between text-lg">
            <span className="text-base-content/70">{t('monthlyPayment')}:</span>
            <span className="font-bold text-primary">
              {formatCurrency(calculation.monthlyPayment)}
            </span>
          </div>

          <div className="flex justify-between">
            <span className="text-base-content/70">{t('totalPayment')}:</span>
            <span className="font-semibold">
              {formatCurrency(calculation.totalPayment)}
            </span>
          </div>

          {calculation.totalInterest > 0 && (
            <div className="flex justify-between">
              <span className="text-base-content/70">
                {t('totalInterest')}:
              </span>
              <span className="text-error">
                {formatCurrency(calculation.totalInterest)}
              </span>
            </div>
          )}

          {calculation.additionalCosts && calculation.additionalCosts > 0 && (
            <div className="flex justify-between">
              <span className="text-base-content/70">
                {t('additionalCosts')}:
              </span>
              <span className="text-warning">
                {formatCurrency(calculation.additionalCosts)}
              </span>
            </div>
          )}

          {calculation.additionalCosts && calculation.additionalCosts > 0 && (
            <div className="flex justify-between text-lg pt-2 border-t">
              <span className="font-semibold">{t('totalWithCosts')}:</span>
              <span className="font-bold text-primary">
                {formatCurrency(
                  calculation.totalPayment + calculation.additionalCosts
                )}
              </span>
            </div>
          )}
        </div>

        {/* График платежей */}
        {calculation.paymentSchedule &&
          calculation.paymentSchedule.length > 0 && (
            <div className="mt-4">
              <button
                className="btn btn-ghost btn-sm w-full"
                onClick={() => setShowSchedule(!showSchedule)}
              >
                <FaCalendarAlt className="mr-2" />
                {showSchedule ? t('hideSchedule') : t('showSchedule')}
              </button>

              {showSchedule && (
                <div className="mt-4 max-h-64 overflow-y-auto">
                  <table className="table table-xs">
                    <thead>
                      <tr>
                        <th>{t('month')}</th>
                        <th>{t('payment')}</th>
                        <th>{t('principal')}</th>
                        <th>{t('interest')}</th>
                        <th>{t('balance')}</th>
                      </tr>
                    </thead>
                    <tbody>
                      {calculation.paymentSchedule.slice(0, 12).map((row) => (
                        <tr key={row.month}>
                          <td>{row.month}</td>
                          <td>{formatCurrency(row.payment)}</td>
                          <td>{formatCurrency(row.principal)}</td>
                          <td>{formatCurrency(row.interest)}</td>
                          <td>{formatCurrency(row.balance)}</td>
                        </tr>
                      ))}
                      {calculation.paymentSchedule.length > 12 && (
                        <tr>
                          <td colSpan={5} className="text-center">
                            ...{' '}
                            {t('andMore', {
                              count: calculation.paymentSchedule.length - 12,
                            })}{' '}
                            ...
                          </td>
                        </tr>
                      )}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
          )}

        {/* Кнопка применения */}
        {onApply && (
          <div className="card-actions justify-end mt-4">
            <button className="btn btn-primary" onClick={handleApply}>
              <FaMoneyBillWave className="mr-2" />
              {t('applyForCredit')}
            </button>
          </div>
        )}
      </div>
    </div>
  );
};

export default UniversalCreditCalculator;
