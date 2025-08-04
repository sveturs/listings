'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { balanceService } from '@/services/balance';
import { useBalance } from '@/hooks/useBalance';
import Link from 'next/link';
import { useLocale } from 'next-intl';

interface BalanceWidgetProps {
  showActions?: boolean;
  compact?: boolean;
}

export const BalanceWidget: React.FC<BalanceWidgetProps> = ({
  showActions = true,
  compact = false,
}) => {
  const t = useTranslations('admin.balance');
  const locale = useLocale();
  const { balance, loading, error } = useBalance();
  const [selectedCurrency, setSelectedCurrency] = useState('RSD');

  // Автоматическое обновление баланса уже происходит в хуке useBalance

  if (loading && !balance) {
    return (
      <div className="flex justify-center items-center p-4">
        <span className="loading loading-spinner loading-md"></span>
      </div>
    );
  }

  if (error && error !== 'Authentication required') {
    return (
      <div className="alert alert-error">
        <span>
          {t('errorLoading')}: {error}
        </span>
      </div>
    );
  }

  // Показать заглушку для неавторизованных пользователей
  if (!balance && error === 'Authentication required') {
    return (
      <div className="card bg-gradient-to-r from-base-300 to-base-200 shadow-lg">
        <div className="card-body text-center">
          <h2 className="card-title justify-center">
            🔒 {t('title')}
          </h2>
          <p className="text-sm opacity-70">
            Войдите в систему для просмотра баланса
          </p>
          <div className="card-actions justify-center mt-4">
            <Link
              href={`/${locale}/auth/login`}
              className="btn btn-primary btn-sm"
            >
              Войти
            </Link>
          </div>
        </div>
      </div>
    );
  }

  if (!balance) {
    return null;
  }

  // Получаем баланс для выбранной валюты
  const displayBalance = balance.balance || 0;
  const frozenBalance = balance.frozen_balance || 0;
  const availableBalance = displayBalance - frozenBalance;

  return (
    <div
      className={`card bg-gradient-to-r from-primary to-secondary text-primary-content shadow-lg ${compact ? 'card-compact' : ''}`}
    >
      <div className="card-body">
        <div className="flex justify-between items-start">
          <div className="flex-1">
            <p className="text-sm opacity-80">
              {t('availableBalance')}
            </p>
            <h2 className="text-3xl font-bold">
              {balanceService.formatAmount(availableBalance, selectedCurrency)}
            </h2>
            {frozenBalance > 0 && (
              <p className="text-sm opacity-80 mt-1">
                {t('frozen')}:{' '}
                {balanceService.formatAmount(frozenBalance, selectedCurrency)}
              </p>
            )}
          </div>

          <select
            className="select select-sm bg-primary-content text-primary"
            value={selectedCurrency}
            onChange={(e) => setSelectedCurrency(e.target.value)}
          >
            <option value="RSD">RSD</option>
            <option value="EUR">EUR</option>
            <option value="USD">USD</option>
          </select>
        </div>

        {showActions && (
          <div className="card-actions justify-center mt-4">
            <Link
              href={`/${locale}/balance/deposit`}
              className="btn btn-primary btn-outline btn-sm"
            >
              💳 {t('deposit')}
            </Link>
            <Link
              href={`/${locale}/balance/withdraw`}
              className="btn btn-primary btn-outline btn-sm"
            >
              💸 {t('withdraw')}
            </Link>
          </div>
        )}
      </div>
    </div>
  );
};
