'use client';

import { useState, useEffect } from 'react';
import {
  MagnifyingGlassIcon,
  TruckIcon,
  ClockIcon,
  CheckIcon,
  MapPinIcon,
  DocumentTextIcon,
  ArrowPathIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  PhoneIcon,
  BuildingStorefrontIcon,
  BanknotesIcon,
} from '@heroicons/react/24/outline';
// import { useTranslations } from 'next-intl';

interface TrackingEvent {
  id: number;
  event_type: string;
  event_code: string;
  description: string;
  timestamp: string;
  location: string;
  postal_code?: string;
  note?: string;
  is_current: boolean;
}

interface ShipmentInfo {
  id: number;
  tracking_number: string;
  status: string;
  created_at: string;
  recipient_name: string;
  recipient_phone: string;
  delivery_address?: string;
  office_code?: string;
  office_name?: string;
  cod_amount?: number;
  weight: number;
  declared_value: number;
  estimated_delivery?: string;
  tracking_events: TrackingEvent[];
}

interface Props {
  initialTrackingNumber?: string;
  onTrackingUpdate?: (shipment: ShipmentInfo | null) => void;
  className?: string;
}

export default function PostExpressTracker({
  initialTrackingNumber = '',
  onTrackingUpdate,
  className = '',
}: Props) {
  // const t = useTranslations('delivery');
  const [trackingNumber, setTrackingNumber] = useState(initialTrackingNumber);
  const [shipment, setShipment] = useState<ShipmentInfo | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [autoRefresh, setAutoRefresh] = useState(false);

  useEffect(() => {
    if (initialTrackingNumber) {
      handleTrack();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [initialTrackingNumber]);

  // Автообновление каждые 30 секунд
  useEffect(() => {
    if (autoRefresh && shipment) {
      const interval = setInterval(() => {
        handleTrack(true); // silent refresh
      }, 30000);

      return () => clearInterval(interval);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [autoRefresh, shipment]);

  const handleTrack = async (silent = false) => {
    if (!trackingNumber.trim()) return;

    if (!silent) setLoading(true);
    setError(null);

    try {
      const response = await fetch(
        `/api/v1/postexpress/track/${encodeURIComponent(trackingNumber)}`
      );
      const data = await response.json();

      if (data.success) {
        const shipmentData = data.data;
        setShipment(shipmentData);
        onTrackingUpdate?.(shipmentData);
      } else {
        setError(data.message || 'Отправление не найдено');
        setShipment(null);
        onTrackingUpdate?.(null);
      }
    } catch (err) {
      setError('Ошибка при отслеживании посылки');
      console.error('Tracking error:', err);
    } finally {
      if (!silent) setLoading(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status.toLowerCase()) {
      case 'delivered':
        return CheckIcon;
      case 'in_transit':
        return TruckIcon;
      case 'at_office':
        return BuildingStorefrontIcon;
      case 'pending':
        return ClockIcon;
      default:
        return DocumentTextIcon;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'delivered':
        return 'text-success';
      case 'in_transit':
        return 'text-primary';
      case 'at_office':
        return 'text-warning';
      case 'pending':
        return 'text-base-content/60';
      default:
        return 'text-base-content/60';
    }
  };

  const getStatusText = (status: string) => {
    switch (status.toLowerCase()) {
      case 'delivered':
        return 'Доставлено';
      case 'in_transit':
        return 'В пути';
      case 'at_office':
        return 'В отделении';
      case 'pending':
        return 'Ожидает отправки';
      case 'cancelled':
        return 'Отменено';
      default:
        return status;
    }
  };

  const formatDateTime = (dateTime: string) => {
    return new Date(dateTime).toLocaleString('sr-RS', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <div className={`space-y-6 ${className}`}>
      {/* Поиск */}
      <div className="card bg-base-100 shadow-lg">
        <div className="card-body p-6">
          <h3 className="text-xl font-bold mb-4 flex items-center gap-2">
            <MagnifyingGlassIcon className="w-6 h-6 text-primary" />
            Отслеживание посылки Post Express
          </h3>

          <div className="flex flex-col sm:flex-row gap-3">
            <div className="flex-1">
              <input
                type="text"
                className="input input-bordered focus:input-primary w-full"
                placeholder="Введите номер отслеживания (например: PE123456789RS)"
                value={trackingNumber}
                onChange={(e) => setTrackingNumber(e.target.value)}
                onKeyPress={(e) => e.key === 'Enter' && handleTrack()}
              />
            </div>
            <div className="flex gap-2">
              <button
                className={`btn btn-primary ${loading ? 'loading' : ''}`}
                onClick={() => handleTrack()}
                disabled={loading || !trackingNumber.trim()}
              >
                {!loading && <MagnifyingGlassIcon className="w-5 h-5" />}
                Отследить
              </button>

              {shipment && (
                <button
                  className={`btn ${autoRefresh ? 'btn-success' : 'btn-outline'}`}
                  onClick={() => setAutoRefresh(!autoRefresh)}
                  title="Автообновление каждые 30 секунд"
                >
                  <ArrowPathIcon
                    className={`w-5 h-5 ${autoRefresh ? 'animate-spin' : ''}`}
                  />
                </button>
              )}
            </div>
          </div>

          {error && (
            <div className="alert alert-error mt-4">
              <ExclamationTriangleIcon className="w-5 h-5" />
              <span>{error}</span>
            </div>
          )}
        </div>
      </div>

      {/* Информация о посылке */}
      {shipment && (
        <div className="space-y-6">
          {/* Общая информация */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body p-6">
              <div className="flex flex-col lg:flex-row justify-between items-start gap-4 mb-6">
                <div>
                  <h4 className="text-xl font-bold mb-2">
                    Отправление {shipment.tracking_number}
                  </h4>
                  <div className="flex items-center gap-2">
                    {(() => {
                      const StatusIcon = getStatusIcon(shipment.status);
                      return (
                        <StatusIcon
                          className={`w-6 h-6 ${getStatusColor(shipment.status)}`}
                        />
                      );
                    })()}
                    <span
                      className={`font-semibold text-lg ${getStatusColor(shipment.status)}`}
                    >
                      {getStatusText(shipment.status)}
                    </span>
                  </div>
                </div>

                {autoRefresh && (
                  <div className="badge badge-primary badge-lg">
                    <ArrowPathIcon className="w-4 h-4 mr-1 animate-spin" />
                    Автообновление
                  </div>
                )}
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
                {/* Получатель */}
                <div className="space-y-1">
                  <div className="text-sm text-base-content/60">Получатель</div>
                  <div className="font-medium">{shipment.recipient_name}</div>
                  <div className="text-sm flex items-center gap-1">
                    <PhoneIcon className="w-4 h-4" />
                    {shipment.recipient_phone}
                  </div>
                </div>

                {/* Адрес доставки */}
                <div className="space-y-1">
                  <div className="text-sm text-base-content/60">
                    Адрес доставки
                  </div>
                  <div className="font-medium">
                    {shipment.delivery_address || (
                      <>
                        {shipment.office_name}
                        <div className="text-sm text-base-content/60">
                          Код: {shipment.office_code}
                        </div>
                      </>
                    )}
                  </div>
                </div>

                {/* Параметры */}
                <div className="space-y-1">
                  <div className="text-sm text-base-content/60">Параметры</div>
                  <div className="font-medium">{shipment.weight} кг</div>
                  <div className="text-sm text-base-content/60">
                    Ценность: {shipment.declared_value} RSD
                  </div>
                </div>

                {/* Наложенный платеж */}
                {shipment.cod_amount && (
                  <div className="space-y-1">
                    <div className="text-sm text-base-content/60">
                      Наложенный платеж
                    </div>
                    <div className="font-medium text-warning flex items-center gap-1">
                      <BanknotesIcon className="w-4 h-4" />
                      {shipment.cod_amount} RSD
                    </div>
                  </div>
                )}
              </div>

              {/* Ожидаемая доставка */}
              {shipment.estimated_delivery && (
                <div className="mt-4 p-3 bg-primary/5 rounded-lg">
                  <div className="flex items-center gap-2">
                    <ClockIcon className="w-5 h-5 text-primary" />
                    <span className="font-medium">
                      Ожидаемая доставка:{' '}
                      {formatDateTime(shipment.estimated_delivery)}
                    </span>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Timeline событий */}
          <div className="card bg-base-100 shadow-lg">
            <div className="card-body p-6">
              <h4 className="text-lg font-bold mb-6">История отслеживания</h4>

              <div className="space-y-4">
                {shipment.tracking_events.map((event, index) => (
                  <div
                    key={event.id}
                    className={`
                      flex gap-4 p-4 rounded-lg transition-all
                      ${
                        event.is_current
                          ? 'bg-primary/10 border-l-4 border-primary'
                          : 'bg-base-200/50'
                      }
                    `}
                  >
                    {/* Timeline dot */}
                    <div className="flex flex-col items-center">
                      <div
                        className={`
                          w-4 h-4 rounded-full border-2 flex-shrink-0
                          ${
                            event.is_current
                              ? 'bg-primary border-primary'
                              : 'bg-base-content/20 border-base-content/20'
                          }
                        `}
                      />
                      {index < shipment.tracking_events.length - 1 && (
                        <div className="w-0.5 h-8 bg-base-content/20 mt-2" />
                      )}
                    </div>

                    {/* Content */}
                    <div className="flex-1 min-w-0">
                      <div className="flex flex-col sm:flex-row sm:items-center justify-between gap-2 mb-2">
                        <div className="font-semibold">{event.description}</div>
                        <div className="text-sm text-base-content/60 flex-shrink-0">
                          {formatDateTime(event.timestamp)}
                        </div>
                      </div>

                      <div className="flex items-center gap-4 text-sm text-base-content/70">
                        {event.location && (
                          <div className="flex items-center gap-1">
                            <MapPinIcon className="w-4 h-4" />
                            {event.location}
                            {event.postal_code && ` (${event.postal_code})`}
                          </div>
                        )}

                        {event.event_code && (
                          <div className="badge badge-outline badge-sm">
                            {event.event_code}
                          </div>
                        )}
                      </div>

                      {event.note && (
                        <div className="mt-2 text-sm text-base-content/80 italic">
                          {event.note}
                        </div>
                      )}
                    </div>
                  </div>
                ))}

                {shipment.tracking_events.length === 0 && (
                  <div className="text-center py-8">
                    <InformationCircleIcon className="w-12 h-12 mx-auto text-base-content/30 mb-3" />
                    <div className="text-base-content/60">
                      События отслеживания пока отсутствуют
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Действия */}
          <div className="card bg-gradient-to-r from-base-200/50 to-base-300/30">
            <div className="card-body p-6">
              <h4 className="font-semibold text-lg mb-4">
                Дополнительные действия
              </h4>

              <div className="flex flex-wrap gap-3">
                <button
                  className="btn btn-outline btn-sm"
                  onClick={() => handleTrack()}
                >
                  <ArrowPathIcon className="w-4 h-4" />
                  Обновить статус
                </button>

                {shipment.status === 'at_office' && (
                  <div className="alert alert-info">
                    <InformationCircleIcon className="w-5 h-5" />
                    <div>
                      <strong>Посылка готова к получению!</strong>
                      <br />
                      Отделение: {shipment.office_name}
                      <br />
                      Не забудьте взять документ удостоверяющий личность
                    </div>
                  </div>
                )}

                {shipment.cod_amount && shipment.status !== 'delivered' && (
                  <div className="alert alert-warning">
                    <BanknotesIcon className="w-5 h-5" />
                    <div>
                      <strong>
                        Наложенный платеж: {shipment.cod_amount} RSD
                      </strong>
                      <br />
                      Подготовьте сумму для оплаты при получении
                    </div>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Пустое состояние */}
      {!shipment && !loading && !error && (
        <div className="text-center py-12">
          <TruckIcon className="w-16 h-16 mx-auto text-base-content/30 mb-4" />
          <h3 className="text-lg font-semibold mb-2">Отслеживание посылок</h3>
          <p className="text-base-content/60">
            Введите номер отслеживания для получения актуальной информации о
            доставке
          </p>
        </div>
      )}
    </div>
  );
}
