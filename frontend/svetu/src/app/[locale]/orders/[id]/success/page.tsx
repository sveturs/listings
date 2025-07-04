'use client';

import { use, useState, useEffect } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { apiClient } from '@/services/api-client';
import { balanceService } from '@/services/balance';

interface Props {
  params: Promise<{ id: string }>;
}

export default function OrderSuccessPage({ params }: Props) {
  const { id } = use(params);
  const locale = useLocale();
  const _t = useTranslations();
  const router = useRouter();
  const { isAuthenticated } = useAuth();
  const searchParams = useSearchParams();

  const sessionId = searchParams.get('session_id');
  const [order, setOrder] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isAuthenticated) {
      router.push(`/${locale}/auth/login`);
      return;
    }

    const fetchOrder = async () => {
      try {
        // Получаем детали заказа из marketplace API
        const response = await apiClient.get(
          `/api/v1/marketplace/orders/${id}`
        );
        console.log('Order details response:', response);

        if (response.data?.success && response.data?.data) {
          setOrder(response.data.data);
        } else {
          setError('Заказ не найден');
        }
      } catch (error: any) {
        console.error('Error fetching order:', error);
        if (error.response?.status === 404) {
          setError('Заказ не найден');
        } else {
          setError('Ошибка загрузки заказа');
        }
      } finally {
        setLoading(false);
      }
    };

    fetchOrder();
  }, [id, isAuthenticated, locale, router]);

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center min-h-[400px]">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-2xl">
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body text-center">
            <div className="text-6xl mb-4">❌</div>
            <h1 className="card-title text-2xl justify-center mb-4">Ошибка</h1>
            <p className="text-base-content/70 mb-6">{error}</p>
            <div className="card-actions justify-center">
              <Link href={`/${locale}/marketplace`} className="btn btn-primary">
                Вернуться к маркетплейсу
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body text-center">
          {/* Иконка успеха */}
          <div className="text-6xl mb-4">✅</div>

          <h1 className="card-title text-2xl justify-center mb-4">
            Заказ успешно оплачен!
          </h1>

          <p className="text-base-content/70 mb-6">
            Ваш заказ #{id} был успешно создан и оплачен.
            {sessionId && ` Номер платежа: ${sessionId}`}
          </p>

          {/* Информация о заказе */}
          {order && (
            <div className="bg-base-200 rounded-lg p-4 mb-6 text-left">
              <h3 className="font-semibold mb-3">Детали заказа:</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span>Номер заказа:</span>
                  <span className="font-medium">#{order.id}</span>
                </div>
                <div className="flex justify-between">
                  <span>Товар:</span>
                  <span className="font-medium">
                    {order.listing?.title || 'Загружается...'}
                  </span>
                </div>
                {order.listing?.description && (
                  <div className="flex justify-between">
                    <span>Описание:</span>
                    <span className="font-medium text-right max-w-[200px]">
                      {order.listing.description}
                    </span>
                  </div>
                )}
                <div className="flex justify-between">
                  <span>Сумма:</span>
                  <span className="font-medium">
                    {order.item_price
                      ? balanceService.formatAmount(order.item_price, 'RSD')
                      : 'Не указано'}
                  </span>
                </div>
                <div className="flex justify-between">
                  <span>Статус:</span>
                  <span className="badge badge-success">
                    {order.status === 'paid' ? 'Оплачен' : order.status}
                  </span>
                </div>
                {order.seller && (
                  <div className="flex justify-between">
                    <span>Продавец:</span>
                    <span className="font-medium">
                      {order.seller.name || order.seller.email}
                    </span>
                  </div>
                )}
              </div>
            </div>
          )}

          {/* Следующие шаги */}
          <div className="alert alert-info mb-6">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="stroke-current shrink-0 h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth="2"
                d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
              ></path>
            </svg>
            <div>
              <h4 className="font-semibold">Что дальше?</h4>
              <p className="text-sm mt-1">
                Продавец получил уведомление о вашем заказе и свяжется с вами
                для согласования деталей доставки. Вы можете отслеживать статус
                заказа в разделе &quot;Мои покупки&quot;.
              </p>
            </div>
          </div>

          {/* Действия */}
          <div className="card-actions justify-center space-x-4">
            <Link
              href={`/${locale}/profile/orders`}
              className="btn btn-primary"
            >
              Мои заказы
            </Link>
            <Link href={`/${locale}/marketplace`} className="btn btn-outline">
              Продолжить покупки
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
