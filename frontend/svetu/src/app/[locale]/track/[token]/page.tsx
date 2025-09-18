'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import { TrackingMap } from '@/components/tracking/TrackingMap';
import { DeliveryInfo } from '@/components/tracking/DeliveryInfo';
import { useWebSocket } from '@/hooks/useWebSocket';
import { useTranslations } from 'next-intl';

interface DeliveryData {
  id: number;
  order_id: number;
  tracking_token: string;
  status: string;
  courier_name: string;
  courier_phone: string;
  pickup_address: string;
  delivery_address: string;
  pickup_latitude: number;
  pickup_longitude: number;
  delivery_latitude: number;
  delivery_longitude: number;
  estimated_delivery_time: string;
  courier_location?: {
    latitude: number;
    longitude: number;
    speed: number;
    heading: number;
    updated_at: string;
  };
  nearby_items?: Array<{
    id: number;
    title: string;
    price: number;
    currency: string;
    location: {
      latitude: number;
      longitude: number;
    };
    distance_meters: number;
    store_name?: string;
    store_logo?: string;
  }>;
}

interface LocationUpdate {
  delivery_id: number;
  latitude: number;
  longitude: number;
  speed: number;
  heading: number;
  timestamp: string;
  eta: string;
  distance_meters: number;
}

interface WebSocketMessage {
  type: string;
  data: LocationUpdate | DeliveryData | any;
}

export default function TrackingPage() {
  const params = useParams();
  const token = params.token as string;
  const t = useTranslations('tracking');

  const [delivery, setDelivery] = useState<DeliveryData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [connectionStatus, setConnectionStatus] = useState<
    'connecting' | 'connected' | 'disconnected'
  >('connecting');

  // WebSocket подключение
  const { lastMessage, connectionState, sendMessage } = useWebSocket(
    `${process.env.NEXT_PUBLIC_WEBSOCKET_URL?.replace('http', 'ws')}/tracking?token=${token}`,
    {
      onOpen: () => setConnectionStatus('connected'),
      onClose: () => setConnectionStatus('disconnected'),
      onError: (error) => {
        console.error('WebSocket error:', error);
        setConnectionStatus('disconnected');
      },
      shouldReconnect: (_closeEvent) => true,
    }
  );

  // Загрузка первоначальных данных о доставке
  useEffect(() => {
    const fetchDeliveryData = async () => {
      try {
        const response = await fetch(
          `${process.env.NEXT_PUBLIC_API_URL}/api/v1/tracking/${token}`
        );

        if (!response.ok) {
          throw new Error(
            response.status === 404
              ? 'Доставка не найдена'
              : 'Ошибка загрузки данных'
          );
        }

        const data = await response.json();
        setDelivery(data.data);
      } catch (err) {
        setError(err instanceof Error ? err.message : 'Неизвестная ошибка');
      } finally {
        setLoading(false);
      }
    };

    if (token) {
      fetchDeliveryData();
    }
  }, [token]);

  // Обработка WebSocket сообщений
  useEffect(() => {
    if (lastMessage) {
      try {
        const message: WebSocketMessage = JSON.parse(lastMessage.data);

        switch (message.type) {
          case 'location_update':
            const locationUpdate = message.data as LocationUpdate;
            setDelivery((prev) => {
              if (!prev || prev.id !== locationUpdate.delivery_id) return prev;

              return {
                ...prev,
                courier_location: {
                  latitude: locationUpdate.latitude,
                  longitude: locationUpdate.longitude,
                  speed: locationUpdate.speed,
                  heading: locationUpdate.heading,
                  updated_at: locationUpdate.timestamp,
                },
                estimated_delivery_time: locationUpdate.eta,
              };
            });
            break;

          case 'delivery_info':
            const deliveryInfo = message.data as DeliveryData;
            setDelivery((prev) =>
              prev ? { ...prev, ...deliveryInfo } : deliveryInfo
            );
            break;

          case 'nearby_items':
            const nearbyItems = message.data;
            setDelivery((prev) =>
              prev ? { ...prev, nearby_items: nearbyItems } : prev
            );
            break;

          case 'eta_update':
            const etaUpdate = message.data;
            setDelivery((prev) =>
              prev
                ? {
                    ...prev,
                    estimated_delivery_time: etaUpdate.eta,
                  }
                : prev
            );
            break;

          default:
            console.log('Unknown message type:', message.type);
        }
      } catch (error) {
        console.error('Error parsing WebSocket message:', error);
      }
    }
  }, [lastMessage]);

  // Запрос обновленного ETA
  const requestETA = () => {
    if (connectionState === 'connected') {
      sendMessage(JSON.stringify({ type: 'request_eta' }));
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center">
          <div className="loading loading-spinner loading-lg text-primary"></div>
          <p className="mt-4 text-base-content">{t('loading')}</p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center max-w-md mx-auto px-4">
          <div className="text-6xl mb-4">📦</div>
          <h1 className="text-2xl font-bold text-base-content mb-2">
            {t('error.title')}
          </h1>
          <p className="text-base-content/70 mb-6">{error}</p>
          <button
            className="btn btn-primary"
            onClick={() => window.location.reload()}
          >
            {t('error.retry')}
          </button>
        </div>
      </div>
    );
  }

  if (!delivery) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center">
          <div className="text-6xl mb-4">🔍</div>
          <h1 className="text-2xl font-bold text-base-content">
            {t('notFound.title')}
          </h1>
          <p className="text-base-content/70">{t('notFound.description')}</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-white shadow-sm">
        <div className="container mx-auto px-4 py-4">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-xl font-bold text-base-content">
                {t('title')} #{delivery.order_id}
              </h1>
              <p className="text-sm text-base-content/70">
                {t('trackingCode')}: {delivery.tracking_token}
              </p>
            </div>

            {/* Connection Status */}
            <div className="flex items-center gap-2">
              <div
                className={`w-2 h-2 rounded-full ${
                  connectionStatus === 'connected'
                    ? 'bg-success'
                    : connectionStatus === 'connecting'
                      ? 'bg-warning'
                      : 'bg-error'
                }`}
              ></div>
              <span className="text-sm text-base-content/70">
                {connectionStatus === 'connected'
                  ? t('status.connected')
                  : connectionStatus === 'connecting'
                    ? t('status.connecting')
                    : t('status.disconnected')}
              </span>
            </div>
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Карта */}
          <div className="lg:col-span-2">
            <div className="bg-white rounded-lg shadow-lg overflow-hidden">
              <TrackingMap delivery={delivery} onRequestETA={requestETA} />
            </div>
          </div>

          {/* Информация о доставке */}
          <div className="lg:col-span-1">
            <DeliveryInfo
              delivery={delivery}
              connectionStatus={connectionStatus}
              onRequestETA={requestETA}
            />
          </div>
        </div>

        {/* Viber Bot Link (если открыто не в Viber) */}
        {!window.navigator.userAgent.includes('Viber') && (
          <div className="mt-6 bg-gradient-to-r from-purple-500 to-purple-600 rounded-lg p-4 text-white">
            <div className="flex items-center gap-3">
              <div className="text-2xl">💬</div>
              <div className="flex-1">
                <h3 className="font-semibold">{t('viber.title')}</h3>
                <p className="text-sm opacity-90">{t('viber.description')}</p>
              </div>
              <a
                href={`viber://pa?chatURI=SveTuBot&context=track_${delivery.tracking_token}`}
                className="btn btn-sm bg-white text-purple-600 hover:bg-gray-100 border-none"
              >
                {t('viber.open')}
              </a>
            </div>
          </div>
        )}
      </div>
    </div>
  );
}
