'use client';

import { useTranslations } from 'next-intl';
import {
  ShieldCheckIcon,
  ClockIcon,
  CheckCircleIcon,
  XCircleIcon,
  InformationCircleIcon,
} from '@heroicons/react/24/outline';

interface EscrowStatusProps {
  status: 'held' | 'released' | 'refunded' | 'pending';
  amount: number;
  currency?: string;
  releaseDate?: Date;
  sellerTier?: 'new' | 'verified' | 'gold';
}

export default function EscrowStatus({
  status,
  amount,
  currency = 'RSD',
  releaseDate,
  sellerTier = 'new',
}: EscrowStatusProps) {
  const t = useTranslations('payment.escrow');

  const formatCurrency = (value: number) => {
    return new Intl.NumberFormat('sr-RS', {
      style: 'currency',
      currency: currency,
    }).format(value);
  };

  const formatDate = (date: Date) => {
    return new Intl.DateTimeFormat('sr-RS', {
      dateStyle: 'medium',
      timeStyle: 'short',
    }).format(date);
  };

  const getEscrowDays = () => {
    switch (sellerTier) {
      case 'gold':
        return 1;
      case 'verified':
        return 3;
      default:
        return 7;
    }
  };

  const statusConfig = {
    held: {
      icon: <ShieldCheckIcon className="w-6 h-6" />,
      color: 'info',
      label: t('held'),
      description: t('payment.escrow.heldDescription', {
        days: getEscrowDays(),
      }),
    },
    released: {
      icon: <CheckCircleIcon className="w-6 h-6" />,
      color: 'success',
      label: t('released'),
      description: t('releasedDescription'),
    },
    refunded: {
      icon: <XCircleIcon className="w-6 h-6" />,
      color: 'error',
      label: t('refunded'),
      description: t('refundedDescription'),
    },
    pending: {
      icon: <ClockIcon className="w-6 h-6" />,
      color: 'warning',
      label: t('pending'),
      description: t('pendingDescription'),
    },
  };

  const config = statusConfig[status];

  return (
    <div className={`alert alert-${config.color}`}>
      <div className="flex flex-col gap-2 w-full">
        <div className="flex items-center gap-3">
          {config.icon}
          <div className="flex-1">
            <h4 className="font-semibold">{config.label}</h4>
            <p className="text-sm opacity-90">{config.description}</p>
          </div>
          <div className="text-right">
            <p className="font-bold text-lg">{formatCurrency(amount)}</p>
            {releaseDate && status === 'held' && (
              <p className="text-xs opacity-75">
                {t('releaseDate')}: {formatDate(releaseDate)}
              </p>
            )}
          </div>
        </div>

        {status === 'held' && (
          <div className="mt-2 pt-2 border-t border-current/20">
            <div className="flex items-start gap-2">
              <InformationCircleIcon className="w-5 h-5 shrink-0 mt-0.5" />
              <div className="text-sm space-y-1">
                <p>{t('buyerProtection')}</p>
                <ul className="list-disc list-inside ml-2 space-y-1">
                  <li>{t('confirmReceipt')}</li>
                  <li>{t('disputeOption')}</li>
                  <li>
                    {t('payment.escrow.autoRelease', { days: getEscrowDays() })}
                  </li>
                </ul>
              </div>
            </div>
          </div>
        )}

        {status === 'held' && sellerTier && (
          <div className="flex items-center gap-2 text-sm">
            <span className="opacity-75">
              {t('sellerTier')}:
            </span>
            <span
              className={`badge badge-sm ${
                sellerTier === 'gold'
                  ? 'badge-warning'
                  : sellerTier === 'verified'
                    ? 'badge-primary'
                    : 'badge-neutral'
              }`}
            >
              {t(`payment.escrow.tier.${sellerTier}`)}
            </span>
          </div>
        )}
      </div>
    </div>
  );
}
