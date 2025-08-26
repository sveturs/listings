'use client';

import { useState } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useTranslations, useLocale } from 'next-intl';
import Link from 'next/link';
import { CheckCircleIcon, XCircleIcon } from '@heroicons/react/24/solid';

export default function PaymentMockPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const locale = useLocale();
  const t = useTranslations('subscription');
  const _tCommon = useTranslations('common');

  const [processing, setProcessing] = useState(false);
  const [status, setStatus] = useState<'pending' | 'success' | 'failed'>(
    'pending'
  );

  const paymentIntent = searchParams.get('payment_intent');
  const amount = searchParams.get('amount');
  const planCode = searchParams.get('plan');
  const cycle = searchParams.get('cycle');

  const handlePaymentConfirm = async (success: boolean) => {
    setProcessing(true);

    // Simulate payment processing
    await new Promise((resolve) => setTimeout(resolve, 1500));

    if (success) {
      setStatus('success');
      // Redirect to success page after a moment
      setTimeout(() => {
        router.push(`/${locale}/subscription/success?plan=${planCode}`);
      }, 2000);
    } else {
      setStatus('failed');
    }

    setProcessing(false);
  };

  if (status === 'success') {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="card bg-base-100 shadow-xl max-w-md w-full mx-4">
          <div className="card-body text-center">
            <CheckCircleIcon className="w-20 h-20 text-success mx-auto mb-4" />
            <h2 className="card-title text-2xl justify-center mb-2">
              {t('paymentMock.success.title')}
            </h2>
            <p className="text-base-content/70">
              {t('paymentMock.success.description')}
            </p>
          </div>
        </div>
      </div>
    );
  }

  if (status === 'failed') {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="card bg-base-100 shadow-xl max-w-md w-full mx-4">
          <div className="card-body text-center">
            <XCircleIcon className="w-20 h-20 text-error mx-auto mb-4" />
            <h2 className="card-title text-2xl justify-center mb-2">
              {t('paymentMock.failed.title')}
            </h2>
            <p className="text-base-content/70 mb-6">
              {t('paymentMock.failed.description')}
            </p>
            <div className="card-actions justify-center">
              <Link href={`/${locale}/pricing`} className="btn btn-primary">
                {t('paymentMock.failed.backToPricing')}
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200 flex items-center justify-center">
      <div className="card bg-base-100 shadow-xl max-w-lg w-full mx-4">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">{t('paymentMock.title')}</h2>

          <div className="alert alert-info mb-4">
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
            <span>{t('paymentMock.testNotice')}</span>
          </div>

          <div className="space-y-3 mb-6">
            <div className="flex justify-between">
              <span className="font-medium">{t('paymentMock.plan')}:</span>
              <span className="capitalize">
                {planCode ? t(`plans.${planCode}`) : '-'}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">
                {t('paymentMock.billingCycle')}:
              </span>
              <span className="capitalize">
                {cycle ? t(`paymentMock.${cycle}`) : '-'}
              </span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">{t('paymentMock.amount')}:</span>
              <span className="text-xl font-bold">€{amount}</span>
            </div>
            <div className="flex justify-between">
              <span className="font-medium">{t('paymentMock.paymentId')}:</span>
              <span className="text-xs font-mono">
                {paymentIntent?.slice(0, 20)}...
              </span>
            </div>
          </div>

          <div className="divider">{t('paymentMock.testCardDetails')}</div>

          <div className="space-y-3 mb-6">
            <input
              type="text"
              placeholder={t('paymentMock.cardNumber')}
              className="input input-bordered w-full"
              defaultValue="4242 4242 4242 4242"
              disabled
            />
            <div className="flex gap-3">
              <input
                type="text"
                placeholder={t('paymentMock.expiryDate')}
                className="input input-bordered flex-1"
                defaultValue="12/25"
                disabled
              />
              <input
                type="text"
                placeholder={t('paymentMock.cvc')}
                className="input input-bordered flex-1"
                defaultValue="123"
                disabled
              />
            </div>
          </div>

          <div className="card-actions justify-end gap-3">
            <button
              className="btn btn-ghost"
              onClick={() => handlePaymentConfirm(false)}
              disabled={processing}
            >
              {t('paymentMock.simulateFailure')}
            </button>
            <button
              className={`btn btn-primary ${processing ? 'loading' : ''}`}
              onClick={() => handlePaymentConfirm(true)}
              disabled={processing}
            >
              {processing
                ? t('paymentMock.processing')
                : t('paymentMock.completePayment')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
