import { useState, useCallback } from 'react';
import { useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import { usePaymentService } from '@/services/payment/paymentServiceFactory';
import type { BuyerInfo } from '@/types/payment';
import * as navigationHelpers from './useAllSecurePaymentHelpers';

// Re-export for convenience
export const { navigateToUrl, getLocationOrigin } = navigationHelpers;

export function useAllSecurePayment() {
  const router = useRouter();
  const locale = useLocale();
  const paymentService = usePaymentService();

  const [paymentProcessing, setPaymentProcessing] = useState(false);
  const [paymentError, setPaymentError] = useState<string | null>(null);

  const createPayment = useCallback(
    async (data: {
      listingId: string;
      amount: number;
      currency: string;
      buyerInfo: BuyerInfo;
    }) => {
      try {
        setPaymentProcessing(true);
        setPaymentError(null);

        console.log('Creating payment with data:', data);

        // Сохраняем данные локально для восстановления после redirect
        if (typeof window !== 'undefined') {
          sessionStorage.setItem(
            'pending_payment',
            JSON.stringify({
              ...data,
              timestamp: Date.now(),
            })
          );
        }

        // Создаем платеж
        const response = await paymentService.createPayment({
          listing_id: data.listingId,
          amount: data.amount,
          currency: data.currency,
          buyer_info: data.buyerInfo,
          return_url: `${navigationHelpers.getLocationOrigin()}/${locale}/payment/process`,
          locale: locale,
        });

        console.log('Payment created:', response);

        // Сохраняем ID платежа
        if (typeof window !== 'undefined') {
          sessionStorage.setItem('payment_id', response.id);
        }

        // Redirect на страницу оплаты (с локалью)
        navigationHelpers.navigateToUrl(`/${locale}${response.redirectUrl}`);

        return response;
      } catch (error) {
        console.error('Payment creation failed:', error);
        const errorMessage =
          error instanceof Error ? error.message : 'Неизвестная ошибка';
        setPaymentError(errorMessage);
        throw error;
      } finally {
        setPaymentProcessing(false);
      }
    },
    [paymentService, locale]
  );

  const checkPaymentStatus = useCallback(
    async (paymentId: string) => {
      try {
        console.log('Checking payment status for:', paymentId);
        const status = await paymentService.getPaymentStatus(paymentId);
        console.log('Payment status:', status);
        return status;
      } catch (error) {
        console.error('Failed to check payment status:', error);
        throw error;
      }
    },
    [paymentService]
  );

  const handlePaymentReturn = useCallback(async () => {
    const paymentId =
      typeof window !== 'undefined'
        ? sessionStorage.getItem('payment_id')
        : null;

    if (!paymentId) {
      router.push(`/${locale}/payment/error?reason=no_payment_id`);
      return;
    }

    try {
      // Проверяем статус
      const status = await checkPaymentStatus(paymentId);

      // Очищаем временные данные
      if (typeof window !== 'undefined') {
        sessionStorage.removeItem('payment_id');
        sessionStorage.removeItem('pending_payment');
      }

      // Redirect в зависимости от статуса
      switch (status.status) {
        case 'captured':
        case 'authorized':
          router.push(`/${locale}/payment/success?id=${paymentId}`);
          break;
        case 'failed':
          router.push(
            `/${locale}/payment/error?id=${paymentId}&reason=${status.error_code || 'declined'}`
          );
          break;
        case 'cancelled':
          router.push(
            `/${locale}/payment/error?id=${paymentId}&reason=cancelled`
          );
          break;
        default:
          // Если статус pending, показываем страницу ожидания
          router.push(`/${locale}/payment/process/${paymentId}`);
      }
    } catch (error) {
      console.error('Error handling payment return:', error);
      router.push(`/${locale}/payment/error?reason=status_check_failed`);
    }
  }, [router, checkPaymentStatus, locale]);

  // Функция для создания тестового платежа (только для разработки)
  const createTestPayment = useCallback(
    async (amount: number = 5000) => {
      return createPayment({
        listingId: 'test_listing_123',
        amount,
        currency: 'RSD',
        buyerInfo: {
          name: 'Test User',
          email: 'test@example.com',
          phone: '+381601234567',
          address: 'Test Address 123, Belgrade',
        },
      });
    },
    [createPayment]
  );

  return {
    createPayment,
    createTestPayment,
    checkPaymentStatus,
    handlePaymentReturn,
    isProcessing: paymentProcessing,
    error: paymentError,
    clearError: () => setPaymentError(null),
  };
}
