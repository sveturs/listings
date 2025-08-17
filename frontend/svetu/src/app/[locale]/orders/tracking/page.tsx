'use client';

import { useState, useEffect, Suspense } from 'react';
import { useTranslations } from 'next-intl';
import { useSearchParams } from 'next/navigation';
import { PostExpressTracker } from '@/components/delivery/postexpress';
import { PageTransition } from '@/components/ui/PageTransition';
import {
  TruckIcon,
  ClipboardDocumentListIcon,
  MagnifyingGlassIcon,
  MapPinIcon,
} from '@heroicons/react/24/outline';

interface Order {
  id: number;
  order_number: string;
  status: string;
  created_at: string;
  total_amount: number;
  shipping_method: string;
  delivery_provider?: string;
  tracking_number?: string;
  customer_name: string;
  shipping_address: {
    street: string;
    city: string;
    postal_code: string;
    country: string;
  };
  items: Array<{
    id: number;
    product_name: string;
    variant_name?: string;
    quantity: number;
    price: number;
    image_url?: string;
  }>;
}

function OrderTrackingContent() {
  const t = useTranslations('orders');
  const tCommon = useTranslations('common');
  const searchParams = useSearchParams();
  const [searchQuery, setSearchQuery] = useState('');
  const [searchType, setSearchType] = useState<'order' | 'tracking'>('order');
  const [orders, setOrders] = useState<Order[]>([]);
  const [selectedOrder, setSelectedOrder] = useState<Order | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // Get initial tracking number from URL params
  useEffect(() => {
    const trackingNumber = searchParams.get('tracking');
    const orderId = searchParams.get('orderId');

    if (trackingNumber) {
      setSearchQuery(trackingNumber);
      setSearchType('tracking');
    } else if (orderId) {
      setSearchQuery(orderId);
      setSearchType('order');
      handleSearch(orderId, 'order');
    }
  }, [searchParams]);

  const handleSearch = async (
    query: string = searchQuery,
    type: string = searchType
  ) => {
    if (!query.trim()) return;

    setLoading(true);
    setError(null);

    try {
      let endpoint = '';
      if (type === 'order') {
        endpoint = `/api/v1/orders/search?q=${encodeURIComponent(query)}`;
      } else {
        endpoint = `/api/v1/orders/by-tracking?tracking=${encodeURIComponent(query)}`;
      }

      const response = await fetch(endpoint);
      const data = await response.json();

      if (data.success) {
        if (type === 'order') {
          setOrders(data.data || []);
          setSelectedOrder(null);
        } else {
          // For tracking search, we expect a single order
          const order = data.data;
          if (order) {
            setSelectedOrder(order);
            setOrders([order]);
          } else {
            setOrders([]);
            setSelectedOrder(null);
          }
        }
      } else {
        setError(data.message || 'Поиск не дал результатов');
        setOrders([]);
        setSelectedOrder(null);
      }
    } catch (err) {
      setError('Ошибка при поиске заказов');
      console.error('Search error:', err);
    } finally {
      setLoading(false);
    }
  };

  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'text-warning';
      case 'processing':
        return 'text-info';
      case 'shipped':
        return 'text-primary';
      case 'delivered':
        return 'text-success';
      case 'cancelled':
        return 'text-error';
      default:
        return 'text-base-content';
    }
  };

  const getStatusText = (status: string) => {
    switch (status.toLowerCase()) {
      case 'pending':
        return 'Ожидает обработки';
      case 'processing':
        return 'В обработке';
      case 'shipped':
        return 'Отправлен';
      case 'delivered':
        return 'Доставлен';
      case 'cancelled':
        return 'Отменен';
      default:
        return status;
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  return (
    <PageTransition>
      <div className="min-h-screen bg-base-100 pt-24">
        <div className="container mx-auto px-4 py-8">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold mb-4 flex items-center justify-center gap-3">
              <TruckIcon className="w-8 h-8 text-primary" />
              Отслеживание заказов
            </h1>
            <p className="text-base-content/70 max-w-2xl mx-auto">
              Найдите свой заказ по номеру заказа или трек-номеру Post Express и
              отследите статус доставки в реальном времени
            </p>
          </div>

          {/* Search Form */}
          <div className="card bg-base-100 shadow-lg mb-8">
            <div className="card-body p-6">
              <div className="flex flex-col sm:flex-row gap-4">
                {/* Search Type Toggle */}
                <div className="flex-shrink-0">
                  <div className="tabs tabs-boxed">
                    <button
                      className={`tab ${searchType === 'order' ? 'tab-active' : ''}`}
                      onClick={() => setSearchType('order')}
                    >
                      <ClipboardDocumentListIcon className="w-4 h-4 mr-2" />
                      По номеру заказа
                    </button>
                    <button
                      className={`tab ${searchType === 'tracking' ? 'tab-active' : ''}`}
                      onClick={() => setSearchType('tracking')}
                    >
                      <TruckIcon className="w-4 h-4 mr-2" />
                      По трек-номеру
                    </button>
                  </div>
                </div>

                {/* Search Input */}
                <div className="flex-1">
                  <div className="relative">
                    <MagnifyingGlassIcon className="absolute left-3 top-1/2 -translate-y-1/2 w-5 h-5 text-base-content/40" />
                    <input
                      type="text"
                      className="input input-bordered w-full pl-11"
                      placeholder={
                        searchType === 'order'
                          ? 'Введите номер заказа (например: ORD-2024-001)'
                          : 'Введите трек-номер Post Express (например: PE123456789RS)'
                      }
                      value={searchQuery}
                      onChange={(e) => setSearchQuery(e.target.value)}
                      onKeyPress={(e) => e.key === 'Enter' && handleSearch()}
                    />
                  </div>
                </div>

                {/* Search Button */}
                <div className="flex-shrink-0">
                  <button
                    className={`btn btn-primary ${loading ? 'loading' : ''}`}
                    onClick={() => handleSearch()}
                    disabled={loading || !searchQuery.trim()}
                  >
                    {!loading && <MagnifyingGlassIcon className="w-5 h-5" />}
                    Найти
                  </button>
                </div>
              </div>

              {error && (
                <div className="alert alert-error mt-4">
                  <span>{error}</span>
                </div>
              )}
            </div>
          </div>

          {/* Search Results */}
          {orders.length > 0 && !selectedOrder && searchType === 'order' && (
            <div className="card bg-base-100 shadow-lg mb-8">
              <div className="card-body p-6">
                <h3 className="font-semibold text-lg mb-4">Найденные заказы</h3>

                <div className="space-y-4">
                  {orders.map((order) => (
                    <div
                      key={order.id}
                      className="card bg-base-200 cursor-pointer hover:shadow-md transition-shadow"
                      onClick={() => setSelectedOrder(order)}
                    >
                      <div className="card-body p-4">
                        <div className="flex flex-col sm:flex-row justify-between items-start gap-4">
                          <div className="flex-1">
                            <div className="flex items-center gap-3 mb-2">
                              <h4 className="font-semibold text-lg">
                                Заказ {order.order_number}
                              </h4>
                              <span
                                className={`badge ${getStatusColor(order.status)}`}
                              >
                                {getStatusText(order.status)}
                              </span>
                            </div>

                            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 text-sm">
                              <div>
                                <div className="text-base-content/60">
                                  Дата заказа:
                                </div>
                                <div className="font-medium">
                                  {formatDate(order.created_at)}
                                </div>
                              </div>
                              <div>
                                <div className="text-base-content/60">
                                  Сумма:
                                </div>
                                <div className="font-medium">
                                  {order.total_amount.toFixed(2)} RSD
                                </div>
                              </div>
                              <div>
                                <div className="text-base-content/60">
                                  Получатель:
                                </div>
                                <div className="font-medium">
                                  {order.customer_name}
                                </div>
                              </div>
                              {order.tracking_number && (
                                <div>
                                  <div className="text-base-content/60">
                                    Трек-номер:
                                  </div>
                                  <div className="font-medium font-mono text-primary">
                                    {order.tracking_number}
                                  </div>
                                </div>
                              )}
                            </div>
                          </div>

                          <div className="flex-shrink-0">
                            <button className="btn btn-primary btn-sm">
                              Подробности
                            </button>
                          </div>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}

          {/* Order Details and Tracking */}
          {selectedOrder && (
            <div className="space-y-8">
              {/* Order Information */}
              <div className="card bg-base-100 shadow-lg">
                <div className="card-body p-6">
                  <div className="flex flex-col sm:flex-row justify-between items-start gap-4 mb-6">
                    <div>
                      <h2 className="text-2xl font-bold mb-2">
                        Заказ {selectedOrder.order_number}
                      </h2>
                      <div className="flex items-center gap-2">
                        <span
                          className={`badge badge-lg ${getStatusColor(selectedOrder.status)}`}
                        >
                          {getStatusText(selectedOrder.status)}
                        </span>
                        {selectedOrder.delivery_provider === 'post_express' && (
                          <span className="badge badge-primary badge-outline">
                            Post Express
                          </span>
                        )}
                      </div>
                    </div>

                    <button
                      className="btn btn-outline btn-sm"
                      onClick={() => setSelectedOrder(null)}
                    >
                      Назад к поиску
                    </button>
                  </div>

                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
                    {/* Customer Info */}
                    <div>
                      <h4 className="font-semibold mb-3">Получатель</h4>
                      <div className="space-y-1 text-sm">
                        <div>{selectedOrder.customer_name}</div>
                        <div className="flex items-start gap-2">
                          <MapPinIcon className="w-4 h-4 text-base-content/60 mt-0.5 flex-shrink-0" />
                          <div>
                            <div>{selectedOrder.shipping_address.street}</div>
                            <div>
                              {selectedOrder.shipping_address.city},{' '}
                              {selectedOrder.shipping_address.postal_code}
                            </div>
                            <div>{selectedOrder.shipping_address.country}</div>
                          </div>
                        </div>
                      </div>
                    </div>

                    {/* Order Info */}
                    <div>
                      <h4 className="font-semibold mb-3">
                        Информация о заказе
                      </h4>
                      <div className="space-y-2 text-sm">
                        <div className="flex justify-between">
                          <span className="text-base-content/70">
                            Дата заказа:
                          </span>
                          <span>{formatDate(selectedOrder.created_at)}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/70">Сумма:</span>
                          <span className="font-medium">
                            {selectedOrder.total_amount.toFixed(2)} RSD
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/70">Товаров:</span>
                          <span>{selectedOrder.items.length}</span>
                        </div>
                      </div>
                    </div>

                    {/* Tracking Info */}
                    {selectedOrder.tracking_number && (
                      <div>
                        <h4 className="font-semibold mb-3">Отслеживание</h4>
                        <div className="space-y-2 text-sm">
                          <div>
                            <div className="text-base-content/70 text-xs">
                              Трек-номер Post Express:
                            </div>
                            <div className="font-mono text-primary font-medium break-all">
                              {selectedOrder.tracking_number}
                            </div>
                          </div>
                          <div>
                            <div className="text-base-content/70 text-xs">
                              Способ доставки:
                            </div>
                            <div>{selectedOrder.shipping_method}</div>
                          </div>
                        </div>
                      </div>
                    )}
                  </div>

                  {/* Order Items */}
                  <div className="mt-6 pt-6 border-t">
                    <h4 className="font-semibold mb-4">Товары в заказе</h4>
                    <div className="space-y-3">
                      {selectedOrder.items.map((item) => (
                        <div
                          key={item.id}
                          className="flex items-center gap-4 p-3 bg-base-200 rounded-lg"
                        >
                          {item.image_url && (
                            <div className="w-12 h-12 rounded-lg overflow-hidden bg-base-300 flex-shrink-0">
                              <img
                                src={item.image_url}
                                alt={item.product_name}
                                className="w-full h-full object-cover"
                              />
                            </div>
                          )}
                          <div className="flex-1">
                            <div className="font-medium">
                              {item.product_name}
                            </div>
                            {item.variant_name && (
                              <div className="text-sm text-base-content/60">
                                {item.variant_name}
                              </div>
                            )}
                            <div className="text-sm">
                              {item.quantity} × {item.price.toFixed(2)} RSD
                            </div>
                          </div>
                          <div className="text-right font-medium">
                            {(item.quantity * item.price).toFixed(2)} RSD
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              </div>

              {/* Post Express Tracking */}
              {selectedOrder.tracking_number &&
                selectedOrder.delivery_provider === 'post_express' && (
                  <PostExpressTracker
                    initialTrackingNumber={selectedOrder.tracking_number}
                    onTrackingUpdate={(shipment) => {
                      console.log('Tracking updated:', shipment);
                    }}
                  />
                )}
            </div>
          )}

          {/* Tracking for direct tracking number search */}
          {searchType === 'tracking' &&
            searchQuery &&
            !orders.length &&
            !loading &&
            !error && (
              <PostExpressTracker
                initialTrackingNumber={searchQuery}
                onTrackingUpdate={(shipment) => {
                  console.log('Tracking updated:', shipment);
                }}
              />
            )}

          {/* Empty State */}
          {!orders.length && !loading && !error && !searchQuery && (
            <div className="text-center py-12">
              <TruckIcon className="w-16 h-16 mx-auto text-base-content/30 mb-4" />
              <h3 className="text-lg font-semibold mb-2">
                Отслеживание заказов
              </h3>
              <p className="text-base-content/60 max-w-md mx-auto">
                Введите номер заказа или трек-номер Post Express для получения
                информации о статусе доставки
              </p>
            </div>
          )}

          {/* No Results */}
          {!orders.length && !loading && searchQuery && !error && (
            <div className="text-center py-12">
              <MagnifyingGlassIcon className="w-16 h-16 mx-auto text-base-content/30 mb-4" />
              <h3 className="text-lg font-semibold mb-2">Заказы не найдены</h3>
              <p className="text-base-content/60">
                По запросу "{searchQuery}" ничего не найдено. Проверьте
                правильность номера.
              </p>
            </div>
          )}
        </div>
      </div>
    </PageTransition>
  );
}

export default function OrderTrackingPage() {
  return (
    <Suspense
      fallback={
        <PageTransition>
          <div className="min-h-screen bg-base-100 pt-24">
            <div className="container mx-auto px-4 py-8">
              <div className="flex justify-center items-center">
                <span className="loading loading-spinner loading-lg"></span>
              </div>
            </div>
          </div>
        </PageTransition>
      }
    >
      <OrderTrackingContent />
    </Suspense>
  );
}
