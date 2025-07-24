'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useLocale } from 'next-intl';
import Link from 'next/link';
import { balanceService } from '@/services/balance';
import { BalanceWidget } from '@/components/balance/BalanceWidget';
import { toast } from '@/utils/toast';

const PRESET_AMOUNTS = [500, 1000, 2000, 5000, 10000];

export default function DepositPage() {
  const t = useTranslations();
  const locale = useLocale();
  const [selectedAmount, setSelectedAmount] = useState<number | null>(null);
  const [customAmount, setCustomAmount] = useState('');
  const [paymentMethod, setPaymentMethod] = useState('allsecure');
  const [isProcessing, setIsProcessing] = useState(false);

  const handleAmountSelect = (amount: number) => {
    setSelectedAmount(amount);
    setCustomAmount('');
  };

  const handleCustomAmountChange = (value: string) => {
    setCustomAmount(value);
    setSelectedAmount(null);
  };

  const getSelectedAmount = (): number | null => {
    if (selectedAmount) return selectedAmount;
    if (customAmount) {
      const amount = parseFloat(customAmount);
      return isNaN(amount) || amount <= 0 ? null : amount;
    }
    return null;
  };

  const handleDeposit = async () => {
    const amount = getSelectedAmount();
    if (!amount) return;

    setIsProcessing(true);
    try {
      const response = await balanceService.createDeposit({
        amount,
        currency: 'RSD',
        payment_method: paymentMethod,
        return_url: `${window.location.origin}/${locale}/balance/deposit/success`,
      });

      console.log('Deposit response:', response);

      if (response?.success && response.data?.payment_url) {
        // Перенаправляем на страницу оплаты
        window.location.href = response.data.payment_url;
      } else {
        console.error('Invalid response structure:', response);
        toast.error(t('admin.balance.error.deposit_failed'));
      }
    } catch (error) {
      console.error('Deposit error:', error);
      toast.error(t('admin.balance.error.deposit_failed'));
    } finally {
      setIsProcessing(false);
    }
  };

  const amount = getSelectedAmount();
  const isValidAmount = amount && amount >= 100 && amount <= 100000;

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Хлебные крошки */}
      <div className="breadcrumbs text-sm mb-6">
        <ul>
          <li>
            <Link href={`/${locale}/profile`}>{t('profile.title')}</Link>
          </li>
          <li>
            <Link href={`/${locale}/balance`}>{t('admin.balance.title')}</Link>
          </li>
          <li>{t('admin.balance.depositFunds')}</li>
        </ul>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Основная форма */}
        <div className="lg:col-span-2">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h1 className="card-title text-2xl mb-6">
                💳 {t('admin.balance.depositFunds')}
              </h1>

              {/* Выбор суммы */}
              <div className="form-control mb-6">
                <label className="label">
                  <span className="label-text font-semibold">
                    {t('admin.balance.selectAmount')}
                  </span>
                </label>

                {/* Предустановленные суммы */}
                <div className="grid grid-cols-3 sm:grid-cols-5 gap-3 mb-4">
                  {PRESET_AMOUNTS.map((amount) => (
                    <button
                      key={amount}
                      type="button"
                      className={`btn ${
                        selectedAmount === amount
                          ? 'btn-primary'
                          : 'btn-outline btn-primary'
                      }`}
                      onClick={() => handleAmountSelect(amount)}
                    >
                      {balanceService.formatAmount(amount, 'RSD')}
                    </button>
                  ))}
                </div>

                {/* Кастомная сумма */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">
                      {t('admin.balance.customAmount')}
                    </span>
                  </label>
                  <div className="input-group">
                    <input
                      type="number"
                      placeholder="0"
                      className={`input input-bordered flex-1 ${
                        customAmount && !getSelectedAmount()
                          ? 'input-error'
                          : ''
                      }`}
                      value={customAmount}
                      onChange={(e) => handleCustomAmountChange(e.target.value)}
                      min="100"
                      max="100000"
                    />
                    <span className="btn btn-ghost">RSD</span>
                  </div>
                  <label className="label">
                    <span className="label-text-alt text-neutral-500">
                      {t('admin.balance.limits')}
                    </span>
                  </label>
                </div>
              </div>

              {/* Способ оплаты */}
              <div className="form-control mb-6">
                <label className="label">
                  <span className="label-text font-semibold">
                    {t('admin.balance.paymentMethod')}
                  </span>
                </label>
                <div className="space-y-3">
                  <label className="label cursor-pointer justify-start gap-3">
                    <input
                      type="radio"
                      name="paymentMethod"
                      className="radio radio-primary"
                      value="allsecure"
                      checked={paymentMethod === 'allsecure'}
                      onChange={(e) => setPaymentMethod(e.target.value)}
                    />
                    <div className="flex items-center gap-3">
                      <div className="flex items-center gap-2">
                        <span className="text-2xl">💳</span>
                        <div>
                          <div className="font-medium">
                            {t('admin.balance.bankCard')}
                          </div>
                          <div className="text-sm text-neutral-500">
                            {t('admin.balance.supportedCards')}
                          </div>
                        </div>
                      </div>
                    </div>
                  </label>
                </div>
              </div>

              {/* Итоговая информация */}
              {amount && isValidAmount && (
                <div className="bg-base-200 rounded-lg p-4 mb-6">
                  <div className="flex justify-between items-center">
                    <span className="font-medium">
                      {t('admin.balance.amountToDeposit')}:
                    </span>
                    <span className="text-xl font-bold text-primary">
                      {balanceService.formatAmount(amount, 'RSD')}
                    </span>
                  </div>
                  <div className="text-sm text-neutral-500 mt-2">
                    {t('admin.balance.commission')}: 0 RSD •{' '}
                    {t('admin.balance.toCredit')}:{' '}
                    {balanceService.formatAmount(amount, 'RSD')}
                  </div>
                </div>
              )}

              {/* Кнопка оплаты */}
              <button
                className="btn btn-primary btn-lg w-full"
                disabled={!isValidAmount || isProcessing}
                onClick={handleDeposit}
              >
                {isProcessing ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    {t('common.loading')}
                  </>
                ) : (
                  <>🔒 {t('admin.balance.continueToPayment')}</>
                )}
              </button>

              {/* Информация о безопасности */}
              <div className="text-center text-sm text-neutral-500 mt-4">
                🔒 {t('admin.balance.securePayment')}
              </div>
            </div>
          </div>
        </div>

        {/* Боковая панель с балансом */}
        <div className="lg:col-span-1">
          <BalanceWidget showActions={false} />

          {/* Информация */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">
                ℹ️ {t('admin.balance.information')}
              </h3>
              <ul className="text-sm space-y-2">
                <li>• {t('admin.balance.instantCredit')}</li>
                <li>• {t('admin.balance.minAmount')}: 100 RSD</li>
                <li>• {t('admin.balance.maxAmount')}: 100,000 RSD</li>
                <li>• {t('admin.balance.noCommission')}: 0%</li>
                <li>
                  • {t('admin.balance.supportedCurrencies')}: RSD, EUR, USD
                </li>
              </ul>
            </div>
          </div>

          {/* Поддержка */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">
                🆘 {t('admin.balance.needHelp')}
              </h3>
              <p className="text-sm">{t('admin.balance.supportMessage')}</p>
              <Link
                href={`/${locale}/contacts`}
                className="btn btn-sm btn-outline"
              >
                {t('admin.balance.contactSupport')}
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
