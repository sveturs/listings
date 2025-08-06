'use client';

import { useTranslations } from 'next-intl';
import {
  CreditCardIcon,
  BanknotesIcon,
  WalletIcon,
} from '@heroicons/react/24/outline';

interface PaymentMethod {
  id: string;
  name: string;
  description: string;
  icon: React.ReactNode;
  fee?: number;
  disabled?: boolean;
}

interface PaymentMethodSelectorProps {
  selectedMethod: string | null;
  onMethodSelect: (methodId: string) => void;
  amount: number;
  currency: string;
}

export default function PaymentMethodSelector({
  selectedMethod,
  onMethodSelect,
  amount,
  currency = 'RSD',
}: PaymentMethodSelectorProps) {
  const _t = useTranslations('payment');
  const tPayment = useTranslations('payment');

  const paymentMethods: PaymentMethod[] = [
    {
      id: 'card',
      name: tPayment('methods.card'),
      description: tPayment('methods.cardDescription'),
      icon: <CreditCardIcon className="w-6 h-6" />,
    },
    {
      id: 'cash_on_delivery',
      name: tPayment('methods.cashOnDelivery'),
      description: tPayment('methods.cashOnDeliveryDescription'),
      icon: <BanknotesIcon className="w-6 h-6" />,
      fee: amount * 0.02, // 2% fee
    },
    {
      id: 'bank_transfer',
      name: tPayment('methods.bankTransfer'),
      description: tPayment('methods.bankTransferDescription'),
      icon: <WalletIcon className="w-6 h-6" />,
    },
  ];

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('sr-RS', {
      style: 'currency',
      currency: currency,
    }).format(value);
  };

  return (
    <div className="space-y-3">
      <h3 className="text-lg font-semibold text-base-content">
        {tPayment('selectMethod')}
      </h3>

      <div className="space-y-2">
        {paymentMethods.map((method) => (
          <div
            key={method.id}
            className={`card border-2 cursor-pointer transition-all ${
              selectedMethod === method.id
                ? 'border-primary bg-primary/5'
                : 'border-base-300 hover:border-primary/50'
            } ${method.disabled ? 'opacity-50 cursor-not-allowed' : ''}`}
            onClick={() => !method.disabled && onMethodSelect(method.id)}
          >
            <div className="card-body p-4">
              <div className="flex items-start gap-4">
                <div
                  className={`${
                    selectedMethod === method.id
                      ? 'text-primary'
                      : 'text-base-content/70'
                  }`}
                >
                  {method.icon}
                </div>

                <div className="flex-1">
                  <h4 className="font-medium text-base-content">
                    {method.name}
                  </h4>
                  <p className="text-sm text-base-content/70 mt-1">
                    {method.description}
                  </p>

                  {method.fee && (
                    <p className="text-sm text-warning mt-2">
                      {tPayment('additionalFee')}: {formatCurrency(method.fee)}
                    </p>
                  )}
                </div>

                <input
                  type="radio"
                  name="payment-method"
                  className="radio radio-primary"
                  checked={selectedMethod === method.id}
                  onChange={() => onMethodSelect(method.id)}
                  disabled={method.disabled}
                />
              </div>
            </div>
          </div>
        ))}
      </div>

      <div className="alert alert-info mt-4">
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
        <span className="text-sm">{tPayment('securePaymentInfo')}</span>
      </div>
    </div>
  );
}
