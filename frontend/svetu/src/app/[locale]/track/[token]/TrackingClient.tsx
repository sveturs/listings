'use client';

import { useEffect, useState } from 'react';
import { useSearchParams } from 'next/navigation';
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
  accuracy: number;
  timestamp: string;
}

interface WebSocketMessage {
  type: string;
  data: LocationUpdate | DeliveryData | any;
}

interface TrackingClientProps {
  token: string;
}

export default function TrackingClient({ token }: TrackingClientProps) {
  const t = useTranslations('tracking');
  const searchParams = useSearchParams();

  // Определяем, открыта ли страница из Viber
  const isViber = searchParams.get('viber') === 'true';
  const isEmbedded = searchParams.get('embedded') === 'true';

  const [delivery, setDelivery] = useState<DeliveryData | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [courierLocation, setCourierLocation] = useState<{
    latitude: number;
    longitude: number;
    speed: number;
    heading: number;
    updatedAt: string;
  } | null>(null);

  // WebSocket connection
  // Используем относительный путь который будет проксирован через Next.js
  const wsProtocol =
    typeof window !== 'undefined' && window.location.protocol === 'https:'
      ? 'wss:'
      : 'ws:';
  const wsHost =
    typeof window !== 'undefined' ? window.location.host : 'localhost:3001';
  const { connectionState, sendMessage } = useWebSocket(
    `${wsProtocol}//${wsHost}/ws/tracking/${token}`,
    {
      onMessage: (event: MessageEvent) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data);
          console.log('WebSocket message:', message);

          if (message.type === 'location_update') {
            const update = message.data as LocationUpdate;
            setCourierLocation({
              latitude: update.latitude,
              longitude: update.longitude,
              speed: update.speed,
              heading: update.heading,
              updatedAt: update.timestamp,
            });
          } else if (message.type === 'delivery_update') {
            setDelivery(message.data as DeliveryData);
          }
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error);
        }
      },
      onError: (error) => {
        console.error('WebSocket error:', error);
      },
      reconnectInterval: 5000,
    }
  );

  const isConnected = connectionState === 'connected';

  // Fetch initial delivery data
  useEffect(() => {
    const fetchDelivery = async () => {
      try {
        const response = await fetch(`/api/v1/tracking/${token}`);

        if (!response.ok) {
          if (response.status === 404) {
            setError(t('notFound.description'));
          } else {
            setError(`Error: ${response.status}`);
          }
          setLoading(false);
          return;
        }

        const result = await response.json();
        // The API returns {data: {...}, success: true}
        setDelivery(result.data || result);

        const deliveryData = result.data || result;
        if (deliveryData.courier_location) {
          setCourierLocation({
            latitude: deliveryData.courier_location.latitude,
            longitude: deliveryData.courier_location.longitude,
            speed: deliveryData.courier_location.speed || 0,
            heading: deliveryData.courier_location.heading || 0,
            updatedAt: deliveryData.courier_location.updated_at,
          });
        }
      } catch (err) {
        console.error('Error fetching delivery:', err);
        setError(t('error.title'));
      } finally {
        setLoading(false);
      }
    };

    fetchDelivery();
  }, [token, t]);

  // Request ETA update
  const handleRequestETA = () => {
    if (isConnected) {
      sendMessage(
        JSON.stringify({
          type: 'request_eta',
          delivery_id: delivery?.id,
        })
      );
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
        <div className="text-center max-w-md mx-auto px-4">
          <div className="text-6xl mb-4">🔍</div>
          <h1 className="text-2xl font-bold text-base-content mb-2">
            {t('notFound.title')}
          </h1>
          <p className="text-base-content/70">{t('notFound.description')}</p>
        </div>
      </div>
    );
  }

  // Если страница открыта из Viber во встроенном режиме, показываем только карту
  if (isEmbedded && isViber) {
    return (
      <div className="h-screen w-full relative">
        <TrackingMap
          delivery={{
            ...delivery,
            courier_location: courierLocation
              ? {
                  latitude: courierLocation.latitude,
                  longitude: courierLocation.longitude,
                  speed: courierLocation.speed,
                  heading: courierLocation.heading,
                }
              : delivery.courier_location,
          }}
          onRequestETA={handleRequestETA}
        />

        {/* Компактная информационная панель */}
        <div className="absolute bottom-0 left-0 right-0 bg-white/95 backdrop-blur-sm p-4 shadow-lg">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-3">
              <div className="text-2xl">
                {delivery.status === 'in_transit' ? '🚴' : '📦'}
              </div>
              <div>
                <div className="font-semibold">
                  {t(`status.${delivery.status || 'undefined'}`)}
                </div>
                <div className="text-sm text-base-content/70">
                  ETA:{' '}
                  {delivery.estimated_delivery_time
                    ? new Date(
                        delivery.estimated_delivery_time
                      ).toLocaleTimeString('ru', {
                        hour: '2-digit',
                        minute: '2-digit',
                      })
                    : '--:--'}
                </div>
              </div>
            </div>
            {courierLocation && (
              <div className="text-sm text-right">
                <div>{courierLocation.speed.toFixed(1)} км/ч</div>
                <div className="text-xs text-base-content/60">
                  {connectionState === 'connected' ? '🟢 Live' : '🔴 Offline'}
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header - скрываем для Viber */}
      {!isViber && (
        <div className="bg-base-100 shadow-sm">
          <div className="container mx-auto px-4 py-4">
            <div className="flex items-center justify-between">
              <h1 className="text-xl font-bold text-base-content">
                {t('title')}
              </h1>
              <div className="badge badge-primary badge-lg">
                {t('trackingCode')}: {token}
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Main Content */}
      <div className="container mx-auto px-4 py-6">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Map - Takes 2 columns on large screens */}
          <div className="lg:col-span-2">
            <TrackingMap
              delivery={{
                ...delivery,
                courier_location: courierLocation
                  ? {
                      latitude: courierLocation.latitude,
                      longitude: courierLocation.longitude,
                      speed: courierLocation.speed,
                      heading: courierLocation.heading,
                    }
                  : delivery.courier_location,
              }}
              onRequestETA={handleRequestETA}
            />
          </div>

          {/* Info Panel - Takes 1 column */}
          <div className="lg:col-span-1">
            <DeliveryInfo
              delivery={{
                ...delivery,
                courier_location: courierLocation
                  ? {
                      latitude: courierLocation.latitude,
                      longitude: courierLocation.longitude,
                      speed: courierLocation.speed,
                      heading: courierLocation.heading,
                      updated_at: courierLocation.updatedAt,
                    }
                  : delivery.courier_location,
              }}
              connectionStatus={
                connectionState === 'error' ? 'disconnected' : connectionState
              }
              onRequestETA={handleRequestETA}
            />
          </div>
        </div>
      </div>

      {/* Footer with Viber CTA - показываем только если не из Viber */}
      {!isViber && (
        <div className="fixed bottom-0 left-0 right-0 bg-base-100 border-t border-base-300 p-4">
          <div className="container mx-auto">
            <div className="flex items-center justify-between">
              <div>
                <p className="font-medium text-base-content">
                  {t('viber.title')}
                </p>
                <p className="text-sm text-base-content/70">
                  {t('viber.description')}
                </p>
              </div>
              <a
                href={`viber://pa?chatURI=svetu_bot&text=track_${token}`}
                className="btn btn-primary"
              >
                <svg
                  className="w-5 h-5 mr-2"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path d="M11.398 0C5.078.254.254 5.078 0 11.398v1.204C.254 18.922 5.078 23.746 11.398 24h1.204c6.32-.254 11.144-5.078 11.398-11.398v-1.204C23.746 5.078 18.922.254 12.602 0h-1.204zm.067 2.923c.59-.014 1.186.007 1.766.063 2.356.229 3.892.893 5.086 2.086 1.272 1.272 1.937 2.943 2.165 5.396.076.915.076 1.843 0 2.758-.229 2.453-.893 4.124-2.165 5.396-1.272 1.272-2.943 1.937-5.396 2.165-.457.038-.915.057-1.372.057-.458 0-.916-.019-1.373-.057-2.453-.228-4.124-.893-5.396-2.165-1.272-1.272-1.937-2.943-2.165-5.396-.076-.915-.076-1.843 0-2.758.228-2.453.893-4.124 2.165-5.396 1.194-1.193 2.73-1.857 5.086-2.086.58-.056 1.175-.077 1.765-.063h.002zm.102 2.137c-.495-.012-.997.005-1.485.051-1.908.19-3.16.686-4.05 1.576-.889.889-1.385 2.141-1.575 4.05-.057.68-.057 1.385 0 2.065.19 1.908.686 3.16 1.575 4.05.89.889 2.142 1.385 4.05 1.575.68.057 1.385.057 2.065 0 1.908-.19 3.16-.686 4.05-1.575.889-.89 1.385-2.142 1.575-4.05.057-.68.057-1.385 0-2.065-.19-1.909-.686-3.161-1.575-4.05-.89-.89-2.142-1.386-4.05-1.576-.488-.046-.99-.063-1.485-.051h-.095zm-.034 2.743c1.718 0 3.436.915 4.351 2.299.153.23.077.537-.153.69-.077.039-.154.077-.23.077-.154 0-.307-.077-.422-.23-.689-1.07-2.03-1.762-3.392-1.762-1.439 0-2.78.69-3.469 1.839-.115.152-.268.23-.422.23-.076 0-.153 0-.23-.039-.23-.153-.307-.46-.153-.69.843-1.461 2.414-2.414 4.12-2.414zm.038 1.991c.574 0 1.073.23 1.438.613.153.191.153.46 0 .651-.192.153-.46.153-.652 0-.191-.23-.498-.345-.805-.345-.345 0-.652.153-.843.422-.154.191-.46.23-.652.038-.191-.153-.23-.46-.077-.652.384-.574 1.035-.881 1.744-.804l-.153.077zm2.011 1.225c.191 0 .383.077.498.268.384.575.575 1.264.575 1.953 0 .69-.191 1.379-.575 1.954-.115.19-.307.268-.498.268-.116 0-.231-.039-.346-.115-.268-.192-.345-.575-.154-.843.269-.422.422-.881.422-1.341 0-.46-.153-.919-.422-1.34-.191-.269-.114-.652.154-.844.115-.076.23-.115.346-.115v.155z" />
                </svg>
                {t('viber.open')}
              </a>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
