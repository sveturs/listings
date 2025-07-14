'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
// import { useTranslations } from 'next-intl';
import { useAllSecurePayment } from '@/hooks/useAllSecurePayment';
import PaymentProcessing from '@/components/payment/PaymentProcessing';

export default function PaymentProcessPage() {
  const params = useParams();
  const id = params?.id as string;
  const router = useRouter();
  const { checkPaymentStatus } = useAllSecurePayment();

  const [status, setStatus] = useState<string>('checking');
  const [error, setError] = useState<string | null>(null);
  const [attempts, setAttempts] = useState(0);

  const MAX_ATTEMPTS = 20; // 20 * 3 sec = 1 minute max
  const POLL_INTERVAL = 3000;

  useEffect(() => {
    const intervalId = setInterval(() => {}, 1000);
    let mounted = true;

    const checkStatus = async () => {
      try {
        const result = await checkPaymentStatus(id as string);

        if (!mounted) return;

        setStatus(result.status);

        // Если платеж завершен, прекращаем polling
        if (
          ['captured', 'authorized', 'failed', 'cancelled'].includes(
            result.status
          )
        ) {
          clearInterval(newIntervalId);

          // Redirect на соответствующую страницу
          setTimeout(() => {
            if (
              result.status === 'captured' ||
              result.status === 'authorized'
            ) {
              router.push(`/payment/success?id=${id}`);
            } else {
              router.push(
                `/payment/error?id=${id}&reason=${result.error_code || 'unknown'}`
              );
            }
          }, 1500);
        }

        setAttempts((prev) => prev + 1);
      } catch (err) {
        if (!mounted) return;
        console.error('Error checking payment status:', err);
        setError(
          err instanceof Error ? err.message : 'Ошибка проверки статуса'
        );
        clearInterval(newIntervalId);
      }
    };

    // Первая проверка сразу
    checkStatus();

    // Затем каждые 3 секунды
    clearInterval(intervalId);
    const newIntervalId = setInterval(() => {
      if (attempts >= MAX_ATTEMPTS) {
        clearInterval(newIntervalId);
        setError('Превышено время ожидания ответа');
        return;
      }

      checkStatus();
    }, POLL_INTERVAL);

    return () => {
      mounted = false;
      clearInterval(newIntervalId);
    };
  }, [id, attempts, checkPaymentStatus, router]);

  if (error) {
    return (
      <div className="min-h-screen flex items-center justify-center p-4">
        <div className="card bg-base-100 shadow-xl max-w-md w-full">
          <div className="card-body text-center">
            <div className="text-error mb-4">
              <svg
                className="w-16 h-16 mx-auto"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                />
              </svg>
            </div>
            <h2 className="card-title justify-center">Ошибка обработки</h2>
            <p className="text-base-content/70">{error}</p>
            <div className="card-actions justify-center mt-4">
              <button
                className="btn btn-primary"
                onClick={() => router.push('/checkout')}
              >
                Попробовать снова
              </button>
              <button
                className="btn btn-ghost"
                onClick={() => router.push('/')}
              >
                На главную
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4">
      <div className="max-w-md w-full">
        <PaymentProcessing
          status={status}
          attempts={attempts}
          maxAttempts={MAX_ATTEMPTS}
        />

        {/* Debug info */}
        {process.env.NODE_ENV === 'development' && (
          <div className="mt-4 p-4 bg-base-200 rounded-lg text-sm">
            <p>Payment ID: {id}</p>
            <p>Status: {status}</p>
            <p>
              Attempts: {attempts}/{MAX_ATTEMPTS}
            </p>
            <p>Next check in: {POLL_INTERVAL / 1000}s</p>
          </div>
        )}
      </div>
    </div>
  );
}
