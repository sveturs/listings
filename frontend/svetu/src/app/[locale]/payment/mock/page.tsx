'use client';

import { useState, useEffect } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useLocale } from 'next-intl';
import MockCardForm from './components/MockCardForm';
import Mock3DSForm from './components/Mock3DSForm';
import { paymentConfig } from '@/config/payment';
import { usePaymentService } from '@/services/payment/paymentServiceFactory';
import { config } from '@/config';

export default function MockPaymentPage() {
  const router = useRouter();
  const locale = useLocale();
  const searchParams = useSearchParams();
  const paymentService = usePaymentService();

  const paymentId = searchParams.get('id');
  const sessionId = searchParams.get('session_id');
  const amount = parseFloat(searchParams.get('amount') || '0');
  const currency = searchParams.get('currency') || '';
  const orderId = searchParams.get('order_id');
  const requires3DS = searchParams.get('require3ds') === 'true';

  const [step, setStep] = useState<'card' | '3ds' | 'processing'>('card');
  const [paymentData, setPaymentData] = useState<any>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    // Поддерживаем оба формата: старый с id и новый с session_id
    if (!paymentId && !sessionId) {
      router.push(`/${locale}/payment/error?reason=no_payment_id`);
      return;
    }

    if (sessionId) {
      // Новый формат с session_id из MockPaymentService
      setPaymentData({
        id: sessionId,
        amount: amount,
        currency: currency.toUpperCase(),
        order_id: orderId ? parseInt(orderId) : null,
        type: orderId ? 'order' : 'balance_deposit',
      });
    } else if (paymentId) {
      // Старый формат - загружаем данные платежа из localStorage
      const data = localStorage.getItem(`mock_payment_${paymentId}`);
      if (data) {
        setPaymentData(JSON.parse(data));
      } else {
        setError('Данные платежа не найдены');
      }
    }

    // Если требуется 3DS, переходим к форме карты
    if (requires3DS) {
      setStep('card');
    }
  }, [
    paymentId,
    sessionId,
    amount,
    currency,
    orderId,
    requires3DS,
    router,
    locale,
  ]);

  const handleCardSubmit = async (_cardData: any) => {
    setStep('processing');
    setError(null);

    try {
      // Симулируем обработку
      await new Promise((resolve) => setTimeout(resolve, 2000));

      if (requires3DS) {
        setStep('3ds');
      } else {
        // Сразу завершаем платеж
        completePayment();
      }
    } catch {
      setError('Ошибка при обработке карты');
      setStep('card');
    }
  };

  const handle3DSSubmit = async (code: string) => {
    if (code === 'cancel') {
      router.push(`/${locale}/payment/error?id=${paymentId}&reason=cancelled`);
      return;
    }

    setStep('processing');
    setError(null);

    try {
      // Симулируем 3DS проверку
      await new Promise((resolve) => setTimeout(resolve, 1500));

      if (paymentService.handle3DSecure) {
        const success = await paymentService.handle3DSecure(paymentId!, code);
        if (success) {
          completePayment();
        } else {
          router.push(
            `/${locale}/payment/error?id=${paymentId}&reason=3ds_failed`
          );
        }
      } else {
        completePayment();
      }
    } catch {
      setError('Ошибка при 3D Secure аутентификации');
      setStep('3ds');
    }
  };

  const completePayment = async () => {
    const success = Math.random() < paymentConfig.mock!.config.successRate;

    if (sessionId) {
      // Новый формат для MockPaymentService
      if (success) {
        // Для пополнения баланса вызываем API для завершения платежа
        if (!orderId) {
          try {
            const response = await fetch(
              `${config.api.url}/balance/mock/complete`,
              {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json',
                  Authorization: `Bearer ${sessionStorage.getItem('svetu_access_token') || localStorage.getItem('access_token')}`,
                },
                body: JSON.stringify({
                  session_id: sessionId,
                  amount: amount,
                }),
              }
            );

            if (!response.ok) {
              console.error('Failed to complete mock payment');
            }
          } catch (error) {
            console.error('Error completing mock payment:', error);
          }
        }

        if (orderId) {
          // Для заказов - подтверждаем оплату через API
          try {
            const response = await fetch(
              `${config.api.url}/marketplace/orders/${orderId}/confirm-payment`,
              {
                method: 'POST',
                headers: {
                  'Content-Type': 'application/json',
                  Authorization: `Bearer ${sessionStorage.getItem('svetu_access_token') || localStorage.getItem('access_token')}`,
                },
                body: JSON.stringify({
                  session_id: sessionId,
                }),
              }
            );

            if (!response.ok) {
              console.error('Failed to confirm order payment');
            }
          } catch (error) {
            console.error('Error confirming order payment:', error);
          }

          router.push(
            `/${locale}/orders/${orderId}/success?session_id=${sessionId}`
          );
        } else {
          // Для пополнения баланса
          router.push(
            `/${locale}/balance/deposit/success?session_id=${sessionId}&amount=${amount}`
          );
        }
      } else {
        // Неуспешный платеж
        if (orderId) {
          router.push(
            `/${locale}/orders/${orderId}/payment-failed?session_id=${sessionId}`
          );
        } else {
          router.push(
            `/${locale}/balance/deposit?error=payment_failed&session_id=${sessionId}`
          );
        }
      }
    } else {
      // Старый формат
      const returnUrl = success
        ? `/${locale}/payment/success?id=${paymentId}`
        : `/${locale}/payment/error?id=${paymentId}&reason=declined`;
      router.push(returnUrl);
    }
  };

  if (error && !paymentData) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center p-4">
        <div className="card bg-base-100 shadow-xl max-w-md w-full">
          <div className="card-body text-center">
            <h2 className="card-title text-error">Ошибка</h2>
            <p>{error}</p>
            <div className="card-actions justify-center">
              <button className="btn btn-primary" onClick={() => router.back()}>
                Назад
              </button>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200 flex items-center justify-center p-4">
      <div className="w-full max-w-md">
        <div className="mockup-browser bg-base-100 shadow-2xl">
          <div className="mockup-browser-toolbar">
            <div className="input text-xs">
              https://sandbox.allsecure.rs/secure-pay
            </div>
          </div>

          <div className="p-6">
            {/* Logo */}
            <div className="text-center mb-6">
              <div className="text-2xl font-bold text-primary">AllSecure</div>
              <div className="text-sm text-base-content/60">
                Test Payment Gateway
              </div>
            </div>

            {/* Error Display */}
            {error && (
              <div className="alert alert-error mb-4">
                <span>{error}</span>
              </div>
            )}

            {/* Payment Amount */}
            {paymentData && (
              <div className="alert alert-info mb-6">
                <div className="flex justify-between w-full">
                  <span>Сумма к оплате:</span>
                  <span className="font-bold">
                    {new Intl.NumberFormat('sr-RS', {
                      style: 'currency',
                      currency: paymentData.currency,
                    }).format(paymentData.amount)}
                  </span>
                </div>
              </div>
            )}

            {/* Step Content */}
            {step === 'card' && (
              <>
                {sessionId && (
                  // Быстрые кнопки для нового формата (MockPaymentService)
                  <div className="space-y-3 mb-6">
                    <div className="text-center">
                      <p className="text-sm font-medium mb-3">
                        Быстрое тестирование:
                      </p>
                    </div>
                    <button
                      className="btn btn-success w-full"
                      onClick={async () => {
                        setStep('processing');

                        // Для пополнения баланса вызываем API для завершения платежа
                        if (!orderId) {
                          try {
                            const token =
                              sessionStorage.getItem('svetu_access_token') ||
                              localStorage.getItem('access_token');
                            console.log(
                              'Mock payment - token:',
                              token ? 'exists' : 'missing'
                            );
                            console.log(
                              'Mock payment - session_id:',
                              sessionId
                            );
                            console.log('Mock payment - amount:', amount);

                            const response = await fetch(
                              `/api/v1/balance/mock/complete`,
                              {
                                method: 'POST',
                                headers: {
                                  'Content-Type': 'application/json',
                                  Authorization: `Bearer ${token}`,
                                },
                                body: JSON.stringify({
                                  session_id: sessionId,
                                  amount: amount,
                                }),
                              }
                            );

                            if (!response.ok) {
                              console.error(
                                'Failed to complete mock payment:',
                                response.status
                              );
                              setError('Ошибка при завершении платежа');
                              setStep('card');
                              return;
                            }

                            const result = await response.json();
                            console.log(
                              'Mock payment completed successfully:',
                              result
                            );
                            console.log(
                              'New balance:',
                              result.data?.new_balance
                            );
                          } catch (error) {
                            console.error(
                              'Error completing mock payment:',
                              error
                            );
                            setError('Ошибка при обработке платежа');
                            setStep('card');
                            return;
                          }
                        }

                        // Для заказов - подтверждаем оплату через API
                        if (orderId) {
                          try {
                            const token =
                              sessionStorage.getItem('svetu_access_token') ||
                              localStorage.getItem('access_token');

                            const response = await fetch(
                              `/api/v1/marketplace/orders/${orderId}/confirm-payment`,
                              {
                                method: 'POST',
                                headers: {
                                  'Content-Type': 'application/json',
                                  Authorization: `Bearer ${token}`,
                                },
                                body: JSON.stringify({
                                  session_id: sessionId,
                                }),
                              }
                            );

                            if (!response.ok) {
                              console.error(
                                'Failed to confirm order payment:',
                                response.status
                              );
                              setError(
                                'Ошибка при подтверждении оплаты заказа'
                              );
                              setStep('card');
                              return;
                            }

                            const result = await response.json();
                            console.log(
                              'Order payment confirmed successfully:',
                              result
                            );
                          } catch (error) {
                            console.error(
                              'Error confirming order payment:',
                              error
                            );
                            setError(
                              'Ошибка при обработке подтверждения оплаты'
                            );
                            setStep('card');
                            return;
                          }
                        }

                        // Небольшая задержка для визуального эффекта
                        setTimeout(() => {
                          if (orderId) {
                            router.push(
                              `/${locale}/orders/${orderId}/success?session_id=${sessionId}`
                            );
                          } else {
                            router.push(
                              `/${locale}/balance/deposit/success?session_id=${sessionId}&amount=${amount}`
                            );
                          }
                        }, 1500);
                      }}
                    >
                      ✅ Имитировать успешный платеж
                    </button>
                    <button
                      className="btn btn-error w-full"
                      onClick={() => {
                        setStep('processing');
                        setTimeout(() => {
                          if (orderId) {
                            router.push(
                              `/${locale}/orders/${orderId}/payment-failed?session_id=${sessionId}`
                            );
                          } else {
                            router.push(
                              `/${locale}/balance/deposit?error=payment_failed&session_id=${sessionId}`
                            );
                          }
                        }, 1500);
                      }}
                    >
                      ❌ Имитировать неуспешный платеж
                    </button>
                    <div className="divider text-xs">или</div>
                  </div>
                )}
                <MockCardForm
                  onSubmit={handleCardSubmit}
                  testCards={paymentConfig.mock!.testCards}
                />
              </>
            )}

            {step === '3ds' && (
              <Mock3DSForm
                onSubmit={handle3DSSubmit}
                amount={paymentData?.amount}
              />
            )}

            {step === 'processing' && (
              <div className="text-center py-8">
                <span className="loading loading-spinner loading-lg text-primary"></span>
                <p className="mt-4">Обработка платежа...</p>
                <p className="text-sm text-base-content/60 mt-2">
                  Пожалуйста, не закрывайте это окно
                </p>
              </div>
            )}

            {/* Security Badge */}
            <div className="mt-6 pt-6 border-t border-base-300">
              <div className="flex items-center justify-center gap-2 text-sm text-base-content/60">
                <svg
                  className="w-4 h-4"
                  fill="currentColor"
                  viewBox="0 0 20 20"
                >
                  <path
                    fillRule="evenodd"
                    d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                    clipRule="evenodd"
                  />
                </svg>
                <span>Secure Test Environment</span>
              </div>
            </div>
          </div>
        </div>

        {/* Debug Panel */}
        {paymentConfig.mock!.debugMode && (
          <div className="mt-4 p-4 bg-base-100 rounded-lg shadow">
            <h3 className="font-bold text-sm mb-2">Debug Info</h3>
            <pre className="text-xs overflow-auto">
              {JSON.stringify(
                {
                  paymentId,
                  requires3DS,
                  step,
                  paymentData: paymentData
                    ? {
                        ...paymentData,
                        buyer_info: paymentData.buyer_info ? '***' : undefined,
                      }
                    : null,
                },
                null,
                2
              )}
            </pre>
          </div>
        )}
      </div>
    </div>
  );
}
