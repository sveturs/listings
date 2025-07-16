'use client';

import { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
// import { useTranslations } from 'next-intl';
import Link from 'next/link';

export default function PaymentSuccessPage() {
  const searchParams = useSearchParams();
  const paymentId = searchParams?.get('id');

  const [orderData, setOrderData] = useState<any>(null);
  const [showDetails, setShowDetails] = useState(false);

  useEffect(() => {
    // Эффект конфетти можно добавить позже через динамический импорт
    // import('canvas-confetti').then(confetti => {
    //   confetti.default({
    //     particleCount: 100,
    //     spread: 70,
    //     origin: { y: 0.6 }
    //   });
    // });

    // Загружаем данные заказа
    if (paymentId) {
      const paymentData = localStorage.getItem(`mock_payment_${paymentId}`);
      if (paymentData) {
        setOrderData(JSON.parse(paymentData));
        // Очищаем после загрузки
        localStorage.removeItem(`mock_payment_${paymentId}`);
      }
    }
  }, [paymentId]);

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-success/10 to-base-100">
      <div className="max-w-md w-full">
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body text-center">
            {/* Success Icon */}
            <div className="mb-4">
              <div className="w-24 h-24 mx-auto bg-success rounded-full flex items-center justify-center animate-bounce">
                <svg
                  className="w-12 h-12 text-success-content"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>
            </div>

            {/* Title */}
            <h1 className="text-3xl font-bold text-base-content mb-2">
              Платеж успешен!
            </h1>

            {/* Subtitle */}
            <p className="text-base-content/70 mb-6">
              Ваш заказ принят и обрабатывается
            </p>

            {/* Order Info */}
            {orderData && (
              <div className="bg-base-200 rounded-lg p-4 mb-6 text-left">
                <div className="space-y-2">
                  <div className="flex justify-between">
                    <span className="text-base-content/70">ID заказа:</span>
                    <span className="font-mono">
                      {paymentId?.slice(0, 8)}...
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">Сумма:</span>
                    <span className="font-semibold">
                      {new Intl.NumberFormat('sr-RS', {
                        style: 'currency',
                        currency: orderData.currency || 'RSD',
                      }).format(orderData.amount)}
                    </span>
                  </div>
                  <div className="flex justify-between">
                    <span className="text-base-content/70">Статус:</span>
                    <span className="badge badge-success">Оплачено</span>
                  </div>
                </div>

                {/* Expandable Details */}
                <button
                  className="btn btn-ghost btn-xs w-full mt-3"
                  onClick={() => setShowDetails(!showDetails)}
                >
                  {showDetails ? 'Скрыть детали' : 'Показать детали'}
                </button>

                {showDetails && (
                  <div className="mt-3 pt-3 border-t border-base-300 text-sm space-y-1">
                    <p>
                      <span className="text-base-content/70">Дата:</span>{' '}
                      {new Date().toLocaleString('sr-RS')}
                    </p>
                    <p>
                      <span className="text-base-content/70">Метод:</span>{' '}
                      Банковская карта
                    </p>
                    <p>
                      <span className="text-base-content/70">
                        ID транзакции:
                      </span>{' '}
                      {paymentId}
                    </p>
                    {orderData.card && (
                      <p>
                        <span className="text-base-content/70">Карта:</span>{' '}
                        **** **** **** {orderData.card.lastFour}
                      </p>
                    )}
                  </div>
                )}
              </div>
            )}

            {/* Escrow Info */}
            <div className="alert alert-info text-left mb-6">
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
                <h3 className="font-bold">Эскроу защита</h3>
                <div className="text-sm">
                  Средства удерживаются в безопасности до получения товара.
                  Продавец получит оплату после подтверждения доставки.
                </div>
              </div>
            </div>

            {/* Actions */}
            <div className="card-actions justify-center gap-2">
              <Link href="/profile/orders" className="btn btn-primary">
                Мои заказы
              </Link>
              <Link href="/" className="btn btn-ghost">
                Продолжить покупки
              </Link>
            </div>

            {/* Receipt */}
            <div className="mt-6 pt-6 border-t border-base-300">
              <p className="text-sm text-base-content/60">
                Чек отправлен на email:{' '}
                {orderData?.buyer_info?.email || 'покупатель'}
              </p>
            </div>
          </div>
        </div>

        {/* What's Next */}
        <div className="mt-4 p-4 bg-base-100/80 backdrop-blur rounded-lg">
          <h3 className="font-semibold mb-2">Что дальше?</h3>
          <ol className="text-sm space-y-1 text-base-content/70">
            <li>1. Продавец уведомлен о заказе</li>
            <li>2. Отслеживайте статус в разделе &quot;Мои заказы&quot;</li>
            <li>3. Подтвердите получение товара</li>
            <li>4. Оставьте отзыв о покупке</li>
          </ol>
        </div>
      </div>
    </div>
  );
}
