'use client';

import { useFeature } from '@/hooks/useConfig';

interface PaymentButtonProps {
  amount: number;
  onPayment: () => void;
  className?: string;
}

export function PaymentButton({
  amount,
  onPayment,
  className = '',
}: PaymentButtonProps) {
  const paymentsEnabled = useFeature('enablePayments');

  // Не рендерим если payments отключены
  if (!paymentsEnabled) {
    return null;
  }

  return (
    <button onClick={onPayment} className={`btn btn-primary ${className}`}>
      Pay ${amount.toFixed(2)}
    </button>
  );
}
