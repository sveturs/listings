'use client';

import { use, useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import { useAuth } from '@/contexts/AuthContext';
import { apiClient } from '@/services/api-client';
import SafeImage from '@/components/SafeImage';
import { balanceService } from '@/services/balance';
import { useBalance } from '@/hooks/useBalance';
import { toast } from 'react-hot-toast';

interface Props {
  params: Promise<{ id: string }>;
}

export default function BuyPage({ params }: Props) {
  const { id } = use(params);
  const locale = useLocale();
  const t = useTranslations('marketplace');
  const tHome = useTranslations('marketplace');
  const tCommon = useTranslations('common');
  const router = useRouter();
  const { user, isAuthenticated } = useAuth();
  const { balance, refreshBalance: _refreshBalance } = useBalance();

  const [listing, setListing] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [paymentMethod, setPaymentMethod] = useState('balance');
  const [message, setMessage] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);

  // Загружаем информацию о товаре
  useEffect(() => {
    const fetchListing = async () => {
      try {
        const response = await apiClient.get(
          `/api/v1/marketplace/listings/${id}`
        );
        if (response.data?.data) {
          setListing(response.data.data);
        } else {
          setListing(response.data);
        }
      } catch (error) {
        console.error('Error fetching listing:', error);
        toast.error(t('errorLoadingListing'));
        router.push(`/${locale}/marketplace/${id}`);
      } finally {
        setLoading(false);
      }
    };

    if (isAuthenticated) {
      fetchListing();
    } else {
      router.push(
        `/${locale}/auth/login?redirect=${encodeURIComponent(window.location.pathname)}`
      );
    }
  }, [id, locale, router, isAuthenticated, t]);

  // Проверяем доступность покупки
  useEffect(() => {
    if (listing && user) {
      if (listing.user_id === user.id) {
        toast.error(t('cannotBuyOwnListing'));
        router.push(`/${locale}/marketplace/${id}`);
      }
    }
  }, [listing, user, id, locale, router, t]);

  const handlePurchase = async () => {
    if (!listing || isProcessing) return;

    setIsProcessing(true);

    try {
      // Проверяем баланс
      if (paymentMethod === 'balance') {
        const currentBalance = balance?.balance || 0;
        if (currentBalance < listing.price) {
          toast.error(t('insufficientBalance'));
          // Предлагаем пополнить баланс
          const shouldDeposit = confirm(t('depositPrompt'));
          if (shouldDeposit) {
            router.push(`/${locale}/balance/deposit`);
          }
          setIsProcessing(false);
          return;
        }
      }

      // Создаем заказ
      const response = await apiClient.post(
        '/api/v1/marketplace/orders/create',
        {
          listing_id: listing.id,
          message: message || undefined,
          payment_method: paymentMethod === 'balance' ? 'card' : paymentMethod,
        }
      );

      if (response.data?.success && response.data?.data?.payment_url) {
        // Перенаправляем на страницу оплаты с учетом локали
        const paymentUrl = response.data.data.payment_url;
        // Если URL начинается с /, добавляем локаль
        if (paymentUrl.startsWith('/')) {
          router.push(`/${locale}${paymentUrl}`);
        } else {
          router.push(paymentUrl);
        }
      } else {
        throw new Error('Invalid response');
      }
    } catch (error: any) {
      console.error('Purchase error:', error);
      const errorMessage = error.response?.data?.error || t('purchaseError');
      toast.error(errorMessage);
    } finally {
      setIsProcessing(false);
    }
  };

  if (loading) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="flex justify-center items-center min-h-[400px]">
          <span className="loading loading-spinner loading-lg"></span>
        </div>
      </div>
    );
  }

  if (!listing) {
    return null;
  }

  const platformFee = listing.price * 0.05; // 5% комиссия
  const sellerReceives = listing.price - platformFee;
  const availableBalance = balance?.balance || 0;
  const hasEnoughBalance = availableBalance >= listing.price;

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Хлебные крошки */}
      <div className="breadcrumbs text-sm mb-6">
        <ul>
          <li>
            <Link href={`/${locale}`}>{tHome('title')}</Link>
          </li>
          <li>
            <Link href={`/${locale}/marketplace/${listing.id}`}>
              {listing.title}
            </Link>
          </li>
          <li>{t('buy')}</li>
        </ul>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
        {/* Левая колонка - детали заказа */}
        <div className="lg:col-span-2">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h1 className="card-title text-2xl mb-6">
                {t('confirmPurchase')}
              </h1>

              {/* Информация о товаре */}
              <div className="border rounded-lg p-4 mb-6">
                <div className="flex gap-4">
                  <figure className="relative w-24 h-24 flex-shrink-0 bg-base-200 rounded-lg overflow-hidden">
                    <SafeImage
                      src={listing.images?.[0]?.public_url}
                      alt={listing.title}
                      fill
                      className="object-cover"
                    />
                  </figure>
                  <div className="flex-1">
                    <h3 className="font-semibold text-lg">{listing.title}</h3>
                    <p className="text-sm text-base-content/70 mb-2">
                      {t('seller')}: {listing.user?.name || 'Unknown'}
                    </p>
                    <p className="text-xl font-bold text-primary">
                      {balanceService.formatAmount(listing.price, 'RSD')}
                    </p>
                  </div>
                </div>
              </div>

              {/* Способ оплаты */}
              <div className="form-control mb-6">
                <label className="label">
                  <span className="label-text font-semibold">
                    {t('paymentMethod')}
                  </span>
                </label>
                <div className="space-y-3">
                  {/* Оплата с баланса */}
                  <label className="label cursor-pointer justify-start gap-3 p-4 border rounded-lg hover:bg-base-200">
                    <input
                      type="radio"
                      name="paymentMethod"
                      className="radio radio-primary"
                      value="balance"
                      checked={paymentMethod === 'balance'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                    />
                    <div className="flex-1">
                      <div className="flex items-center gap-2">
                        <span className="text-lg">💰</span>
                        <span className="font-medium">
                          {t('payFromBalance')}
                        </span>
                      </div>
                      <div className="text-sm text-base-content/70 mt-1">
                        {t('availableBalance')}:{' '}
                        <span
                          className={
                            hasEnoughBalance ? 'text-success' : 'text-error'
                          }
                        >
                          {balanceService.formatAmount(availableBalance, 'RSD')}
                        </span>
                      </div>
                    </div>
                  </label>

                  {/* Банковская карта */}
                  <label className="label cursor-pointer justify-start gap-3 p-4 border rounded-lg hover:bg-base-200">
                    <input
                      type="radio"
                      name="paymentMethod"
                      className="radio radio-primary"
                      value="card"
                      checked={paymentMethod === 'card'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                    />
                    <div className="flex items-center gap-2">
                      <span className="text-lg">💳</span>
                      <span className="font-medium">{t('bankCard')}</span>
                    </div>
                  </label>
                </div>

                {paymentMethod === 'balance' && !hasEnoughBalance && (
                  <div className="alert alert-warning mt-3">
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
                        d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                      />
                    </svg>
                    <div>
                      <p>{t('insufficientBalanceMessage')}</p>
                      <Link
                        href={`/${locale}/balance/deposit`}
                        className="link link-primary"
                      >
                        {t('depositNow')}
                      </Link>
                    </div>
                  </div>
                )}
              </div>

              {/* Сообщение продавцу */}
              <div className="form-control mb-6">
                <label className="label">
                  <span className="label-text">{t('messageToSeller')}</span>
                  <span className="label-text-alt">{tCommon('optional')}</span>
                </label>
                <textarea
                  className="textarea textarea-bordered h-24"
                  placeholder={t('messagePlaceholder')}
                  value={message}
                  onChange={(e) => setMessage(e.target.value)}
                  maxLength={500}
                />
              </div>

              {/* Информация о защите */}
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
                  <h4 className="font-semibold">🛡️ {t('buyerProtection')}</h4>
                  <ul className="text-sm mt-2 space-y-1">
                    <li>• {t('protectionPoint1')}</li>
                    <li>• {t('protectionPoint2')}</li>
                    <li>• {t('protectionPoint3')}</li>
                  </ul>
                </div>
              </div>

              {/* Кнопка покупки */}
              <button
                className="btn btn-primary btn-lg w-full"
                onClick={handlePurchase}
                disabled={
                  isProcessing ||
                  (paymentMethod === 'balance' && !hasEnoughBalance)
                }
              >
                {isProcessing ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    {tCommon('processing')}
                  </>
                ) : (
                  <>
                    🔒 {t('confirmAndPay')}{' '}
                    {balanceService.formatAmount(listing.price, 'RSD')}
                  </>
                )}
              </button>
            </div>
          </div>
        </div>

        {/* Правая колонка - сводка */}
        <div className="lg:col-span-1">
          {/* Детали платежа */}
          <div className="card bg-base-100 shadow-xl mb-6">
            <div className="card-body">
              <h3 className="card-title text-lg mb-4">{t('paymentDetails')}</h3>

              <div className="space-y-3">
                <div className="flex justify-between">
                  <span>{t('itemPrice')}</span>
                  <span className="font-medium">
                    {balanceService.formatAmount(listing.price, 'RSD')}
                  </span>
                </div>

                <div className="flex justify-between text-sm text-base-content/70">
                  <span>{t('platformFee')} (5%)</span>
                  <span>
                    -{balanceService.formatAmount(platformFee, 'RSD')}
                  </span>
                </div>

                <div className="divider my-2"></div>

                <div className="flex justify-between">
                  <span>{t('sellerReceives')}</span>
                  <span className="font-medium">
                    {balanceService.formatAmount(sellerReceives, 'RSD')}
                  </span>
                </div>

                <div className="divider my-2"></div>

                <div className="flex justify-between text-lg font-bold">
                  <span>{t('youPay')}</span>
                  <span className="text-primary">
                    {balanceService.formatAmount(listing.price, 'RSD')}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Условия */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h3 className="card-title text-lg mb-4">
                {t('termsAndConditions')}
              </h3>

              <div className="text-sm space-y-2">
                <p>{t('termsText1')}</p>
                <p>{t('termsText2')}</p>
                <p className="text-base-content/70 mt-4">
                  {t('termsText3')}{' '}
                  <Link href={`/${locale}/terms`} className="link link-primary">
                    {t('termsLink')}
                  </Link>
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
