'use client';

import { useTranslations } from 'next-intl';
import { formatDistanceToNow, format } from 'date-fns';
import { ru, sr } from 'date-fns/locale';
import { useLocale } from 'next-intl';

interface DeliveryData {
  id: number;
  order_id: number;
  tracking_token: string;
  status: string;
  courier_name: string;
  courier_phone: string;
  pickup_address: string;
  delivery_address: string;
  estimated_delivery_time: string;
  courier_location?: {
    latitude: number;
    longitude: number;
    speed: number;
    heading: number;
    updated_at: string;
  };
}

interface DeliveryInfoProps {
  delivery: DeliveryData;
  connectionStatus: 'connecting' | 'connected' | 'disconnected';
  onRequestETA?: () => void;
}

export function DeliveryInfo({
  delivery,
  connectionStatus,
  onRequestETA,
}: DeliveryInfoProps) {
  const t = useTranslations('tracking.info');
  const locale = useLocale();
  const dateLocale = locale === 'ru' ? ru : locale === 'sr' ? sr : undefined;

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending':
        return 'text-warning';
      case 'picked_up':
        return 'text-info';
      case 'in_transit':
        return 'text-primary';
      case 'delivered':
        return 'text-success';
      case 'failed':
        return 'text-error';
      default:
        return 'text-base-content';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending':
        return '⏳';
      case 'picked_up':
        return '📦';
      case 'in_transit':
        return '🚚';
      case 'delivered':
        return '✅';
      case 'failed':
        return '❌';
      default:
        return '📍';
    }
  };

  const formatETA = (dateString: string) => {
    try {
      const date = new Date(dateString);
      const now = new Date();

      if (date > now) {
        return {
          relative: formatDistanceToNow(date, {
            addSuffix: true,
            locale: dateLocale,
          }),
          absolute: format(date, 'HH:mm', { locale: dateLocale }),
        };
      } else {
        return {
          relative: t('eta.overdue'),
          absolute: format(date, 'HH:mm', { locale: dateLocale }),
        };
      }
    } catch {
      return {
        relative: t('eta.unknown'),
        absolute: '--:--',
      };
    }
  };

  const eta = formatETA(delivery.estimated_delivery_time);

  return (
    <div className="space-y-4">
      {/* Статус доставки */}
      <div className="bg-white rounded-lg shadow-lg p-4">
        <h2 className="text-lg font-semibold text-base-content mb-3">
          {t('status.title')}
        </h2>

        <div className="flex items-center gap-3 mb-4">
          <div className="text-2xl">{getStatusIcon(delivery.status)}</div>
          <div>
            <div className={`font-semibold ${getStatusColor(delivery.status)}`}>
              {t(`status.${delivery.status}`)}
            </div>
            <div className="text-sm text-base-content/70">
              {t('orderNumber', { number: delivery.order_id })}
            </div>
          </div>
        </div>

        {/* ETA */}
        <div className="bg-base-200 rounded-lg p-3 mb-3">
          <div className="flex items-center justify-between">
            <div>
              <div className="text-sm text-base-content/70">
                {t('eta.title')}
              </div>
              <div className="font-semibold text-lg">{eta.absolute}</div>
              <div className="text-xs text-base-content/60">{eta.relative}</div>
            </div>
            {onRequestETA && (
              <button
                onClick={onRequestETA}
                className="btn btn-sm btn-ghost"
                title={t('eta.refresh')}
              >
                🔄
              </button>
            )}
          </div>
        </div>

        {/* Live tracking indicator */}
        {delivery.courier_location && (
          <div className="flex items-center gap-2 text-sm">
            <div
              className={`w-2 h-2 rounded-full ${
                connectionStatus === 'connected'
                  ? 'bg-success animate-pulse'
                  : 'bg-error'
              }`}
            ></div>
            <span className="text-base-content/70">
              {connectionStatus === 'connected'
                ? t('liveTracking')
                : t('connectionLost')}
            </span>
          </div>
        )}
      </div>

      {/* Информация о курьере */}
      <div className="bg-white rounded-lg shadow-lg p-4">
        <h2 className="text-lg font-semibold text-base-content mb-3">
          {t('courier.title')}
        </h2>

        <div className="flex items-center gap-3 mb-3">
          <div className="avatar placeholder">
            <div className="bg-primary text-primary-content rounded-full w-12">
              <span className="text-xl">🚴‍♂️</span>
            </div>
          </div>
          <div className="flex-1">
            <div className="font-semibold">{delivery.courier_name}</div>
            <div className="text-sm text-base-content/70">
              {t('courier.role')}
            </div>
          </div>
        </div>

        {/* Контакты */}
        <div className="space-y-2">
          <a
            href={`tel:${delivery.courier_phone}`}
            className="btn btn-sm btn-outline w-full gap-2"
          >
            📞 {t('courier.call')}
          </a>

          {/* Viber чат (если открыто в Viber) */}
          {window.navigator.userAgent.includes('Viber') && (
            <button className="btn btn-sm btn-primary w-full gap-2">
              💬 {t('courier.chat')}
            </button>
          )}
        </div>

        {/* Текущие данные курьера */}
        {delivery.courier_location && (
          <div className="mt-3 pt-3 border-t border-base-300">
            <div className="text-sm text-base-content/70 mb-2">
              {t('courier.currentData')}
            </div>
            <div className="grid grid-cols-2 gap-2 text-xs">
              <div className="bg-base-200 rounded p-2">
                <div className="text-base-content/60">{t('courier.speed')}</div>
                <div className="font-semibold">
                  {delivery.courier_location.speed.toFixed(1)} км/ч
                </div>
              </div>
              <div className="bg-base-200 rounded p-2">
                <div className="text-base-content/60">
                  {t('courier.direction')}
                </div>
                <div className="font-semibold">
                  {delivery.courier_location.heading}°
                </div>
              </div>
            </div>
            <div className="text-xs text-base-content/50 mt-2">
              {t('courier.lastUpdate')}:{' '}
              {formatDistanceToNow(
                new Date(delivery.courier_location.updated_at),
                { addSuffix: true, locale: dateLocale }
              )}
            </div>
          </div>
        )}
      </div>

      {/* Адреса */}
      <div className="bg-white rounded-lg shadow-lg p-4">
        <h2 className="text-lg font-semibold text-base-content mb-3">
          {t('addresses.title')}
        </h2>

        <div className="space-y-3">
          {/* Адрес получения */}
          <div className="flex gap-3">
            <div className="text-lg mt-1">🏢</div>
            <div className="flex-1">
              <div className="text-sm text-base-content/70">
                {t('addresses.pickup')}
              </div>
              <div className="font-medium text-sm">
                {delivery.pickup_address}
              </div>
            </div>
          </div>

          {/* Разделитель */}
          <div className="flex items-center gap-2">
            <div className="h-px bg-base-300 flex-1"></div>
            <div className="text-xs text-base-content/50">→</div>
            <div className="h-px bg-base-300 flex-1"></div>
          </div>

          {/* Адрес доставки */}
          <div className="flex gap-3">
            <div className="text-lg mt-1">🏠</div>
            <div className="flex-1">
              <div className="text-sm text-base-content/70">
                {t('addresses.delivery')}
              </div>
              <div className="font-medium text-sm">
                {delivery.delivery_address}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Дополнительные действия */}
      <div className="bg-white rounded-lg shadow-lg p-4">
        <h2 className="text-lg font-semibold text-base-content mb-3">
          {t('actions.title')}
        </h2>

        <div className="space-y-2">
          <button className="btn btn-sm btn-outline w-full gap-2">
            📱 {t('actions.share')}
          </button>

          <button className="btn btn-sm btn-outline w-full gap-2">
            🔔 {t('actions.notifications')}
          </button>

          <button className="btn btn-sm btn-outline w-full gap-2">
            ℹ️ {t('actions.support')}
          </button>
        </div>
      </div>
    </div>
  );
}
