'use client';

import { use, useState, useEffect } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';

interface Props {
  params: Promise<{ id: string }>;
}

export default function OrderPaymentFailedPage({ params }: Props) {
  const { id } = use(params);
  const locale = useLocale();
  const t = useTranslations();
  const router = useRouter();
  const { user, isAuthenticated } = useAuth();
  const searchParams = useSearchParams();
  
  const sessionId = searchParams.get('session_id');
  const reason = searchParams.get('reason') || 'unknown';

  useEffect(() => {
    if (!isAuthenticated) {
      router.push(`/${locale}/auth/login`);
      return;
    }
  }, [isAuthenticated, locale, router]);

  const getErrorMessage = (reason: string) => {
    switch (reason) {
      case 'cancelled':
        return 'Платеж был отменен пользователем';
      case 'declined':
        return 'Платеж был отклонен банком';
      case '3ds_failed':
        return '3D Secure аутентификация не прошла';
      case 'insufficient_funds':
        return 'Недостаточно средств на карте';
      case 'expired_card':
        return 'Карта просрочена';
      default:
        return 'Произошла ошибка при обработке платежа';
    }
  };

  const getSuggestion = (reason: string) => {
    switch (reason) {
      case 'cancelled':
        return 'Вы можете попробовать оплатить заказ снова или выбрать другой способ оплаты';
      case 'declined':
        return 'Проверьте данные карты и попробуйте еще раз, или обратитесь в ваш банк';
      case '3ds_failed':
        return 'Убедитесь, что SMS с кодом подтверждения доходят на ваш телефон';
      case 'insufficient_funds':
        return 'Пополните счет карты или используйте другую карту';
      case 'expired_card':
        return 'Используйте действующую карту';
      default:
        return 'Попробуйте еще раз или свяжитесь с поддержкой';
    }
  };

  return (
    <div className="container mx-auto px-4 py-8 max-w-2xl">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body text-center">
          {/* Иконка ошибки */}
          <div className="text-6xl mb-4">❌</div>
          
          <h1 className="card-title text-2xl justify-center mb-4">
            Оплата не прошла
          </h1>
          
          <p className="text-base-content/70 mb-4">
            К сожалению, при оплате заказа #{id} произошла ошибка.
          </p>

          {/* Детали ошибки */}
          <div className="alert alert-error mb-6">
            <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
            </svg>
            <div>
              <h4 className="font-semibold">Причина:</h4>
              <p>{getErrorMessage(reason)}</p>
            </div>
          </div>

          {/* Рекомендации */}
          <div className="alert alert-info mb-6">
            <svg xmlns="http://www.w3.org/2000/svg" className="stroke-current shrink-0 h-6 w-6" fill="none" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
            <div>
              <h4 className="font-semibold">Что делать?</h4>
              <p className="text-sm mt-1">{getSuggestion(reason)}</p>
            </div>
          </div>

          {sessionId && (
            <p className="text-sm text-base-content/60 mb-6">
              Номер попытки платежа: {sessionId}
            </p>
          )}

          {/* Действия */}
          <div className="card-actions justify-center space-x-4">
            <Link 
              href={`/${locale}/marketplace/${id.split('_')[0] || '27'}/buy`} 
              className="btn btn-primary"
            >
              Попробовать снова
            </Link>
            <Link href={`/${locale}/marketplace`} className="btn btn-outline">
              Вернуться к маркетплейсу
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}