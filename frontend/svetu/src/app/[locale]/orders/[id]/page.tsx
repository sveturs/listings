'use client';

import { use, useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import Image from 'next/image';
import { useAuth } from '@/contexts/AuthContext';
import { apiClient } from '@/services/api-client';
import { balanceService } from '@/services/balance';
import { ImageGallery } from '@/components/reviews/ImageGallery';

interface Props {
  params: Promise<{ id: string }>;
}

export default function OrderDetailsPage({ params }: Props) {
  const { id } = use(params);
  const locale = useLocale();
  const t = useTranslations('orders');
  const tCommon = useTranslations('common');
  const router = useRouter();
  const { user, isAuthenticated } = useAuth();
  const [order, setOrder] = useState<any>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [galleryOpen, setGalleryOpen] = useState(false);
  const [selectedImageIndex, setSelectedImageIndex] = useState(0);

  useEffect(() => {
    if (!isAuthenticated) {
      router.push(`/${locale}/auth/login`);
      return;
    }

    const fetchOrder = async () => {
      try {
        const response = await apiClient.get(
          `/api/v1/marketplace/orders/${id}`
        );
        console.log('Order details:', response);

        if (response.data?.success && response.data?.data) {
          setOrder(response.data.data);
        } else {
          setError(t('orderNotFound'));
        }
      } catch (error: any) {
        console.error('Error fetching order:', error);
        if (error.response?.status === 404) {
          setError(t('orderNotFound'));
        } else if (error.response?.status === 403) {
          setError(t('orderAccessDenied'));
        } else {
          setError(t('orderLoadError'));
        }
      } finally {
        setLoading(false);
      }
    };

    fetchOrder();
  }, [id, isAuthenticated, locale, router]);

  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'pending':
        return 'badge-warning';
      case 'paid':
        return 'badge-info';
      case 'shipped':
        return 'badge-primary';
      case 'delivered':
        return 'badge-success';
      case 'completed':
        return 'badge-success';
      case 'cancelled':
        return 'badge-error';
      case 'disputed':
        return 'badge-error';
      case 'refunded':
        return 'badge-ghost';
      default:
        return 'badge-ghost';
    }
  };

  const getStatusLabel = (status: string) => {
    return t(`status.${status}` as any) || status;
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return new Intl.DateTimeFormat(locale, {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
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

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8 max-w-4xl">
        <div className="alert alert-error">
          <span>{error}</span>
          <Link href={`/${locale}/profile/orders`} className="btn btn-sm">
            {tCommon('back')}
          </Link>
        </div>
      </div>
    );
  }

  if (!order) {
    return null;
  }

  const isBuyer = user?.id === order.buyer_id;
  const isSeller = user?.id === order.seller_id;

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-3xl font-bold">
          {t('order')} #{order.id}
        </h1>
        <div className={`badge badge-lg ${getStatusBadgeClass(order.status)}`}>
          {getStatusLabel(order.status)}
        </div>
      </div>

      {/* Main Content */}
      <div className="grid gap-6 lg:grid-cols-3">
        {/* Order Details */}
        <div className="lg:col-span-2 space-y-6">
          {/* Product Info */}
          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <h2 className="card-title">{t('productInfo')}</h2>

              {order.listing && (
                <div className="space-y-4">
                  <div className="flex gap-4">
                    {order.listing.images &&
                      order.listing.images.length > 0 && (
                        <div
                          className="w-32 h-32 flex-shrink-0 cursor-pointer relative overflow-hidden rounded-lg"
                          onClick={() => {
                            setSelectedImageIndex(0);
                            setGalleryOpen(true);
                          }}
                        >
                          <Image
                            src={order.listing.images[0].url}
                            alt={order.listing.title}
                            fill
                            className="object-cover hover:scale-110 transition-transform duration-200"
                            sizes="(max-width: 768px) 100vw, 300px"
                          />
                          {order.listing.images.length > 1 && (
                            <div className="absolute bottom-2 right-2 bg-black/60 text-white text-xs px-2 py-1 rounded">
                              +{order.listing.images.length - 1}
                            </div>
                          )}
                        </div>
                      )}

                    <div className="flex-1">
                      <h3 className="text-xl font-semibold">
                        {order.listing.title}
                      </h3>
                      {order.listing.description && (
                        <p className="text-base-content/70 mt-2">
                          {order.listing.description}
                        </p>
                      )}
                    </div>
                  </div>

                  <div className="divider"></div>

                  <div className="grid grid-cols-2 gap-4 text-sm">
                    <div>
                      <span className="text-base-content/70">
                        {t('productPrice')}:
                      </span>
                      <p className="font-semibold">
                        {balanceService.formatAmount(order.item_price, 'RSD')}
                      </p>
                    </div>
                    <div>
                      <span className="text-base-content/70">
                        {t('orderDate')}:
                      </span>
                      <p className="font-semibold">
                        {formatDate(order.created_at)}
                      </p>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>

          {/* Shipping Info */}
          {(order.status === 'shipped' ||
            order.status === 'delivered' ||
            order.status === 'completed') && (
            <div className="card bg-base-100 shadow">
              <div className="card-body">
                <h2 className="card-title">{t('shippingInfo')}</h2>

                <div className="space-y-2 text-sm">
                  {order.shipping_method && (
                    <div>
                      <span className="text-base-content/70">
                        {t('shippingMethod')}:
                      </span>
                      <p className="font-semibold">{order.shipping_method}</p>
                    </div>
                  )}

                  {order.tracking_number && (
                    <div>
                      <span className="text-base-content/70">
                        {t('tracking.number')}:
                      </span>
                      <p className="font-semibold">{order.tracking_number}</p>
                    </div>
                  )}

                  {order.shipped_at && (
                    <div>
                      <span className="text-base-content/70">
                        {t('shippedDate')}:
                      </span>
                      <p className="font-semibold">
                        {formatDate(order.shipped_at)}
                      </p>
                    </div>
                  )}

                  {order.delivered_at && (
                    <div>
                      <span className="text-base-content/70">
                        {t('deliveredDate')}:
                      </span>
                      <p className="font-semibold">
                        {formatDate(order.delivered_at)}
                      </p>
                    </div>
                  )}
                </div>
              </div>
            </div>
          )}

          {/* Actions */}
          {isBuyer && order.status === 'shipped' && (
            <div className="card bg-base-100 shadow">
              <div className="card-body">
                <h2 className="card-title">{t('actions')}</h2>
                <p className="text-sm text-base-content/70">
                  {t('confirmReceiptMessage')}
                </p>
                <div className="card-actions justify-end">
                  <button className="btn btn-primary">
                    {t('confirmReceipt')}
                  </button>
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Sidebar */}
        <div className="space-y-6">
          {/* Participants */}
          <div className="card bg-base-100 shadow">
            <div className="card-body">
              <h2 className="card-title text-lg">{t('dealParticipants')}</h2>

              <div className="space-y-4">
                {/* Seller */}
                <div>
                  <p className="text-sm text-base-content/70">{t('seller')}:</p>
                  <p className="font-semibold">
                    {order.seller?.name || order.seller?.email || t('unknown')}
                    {isSeller && ` (${t('you')})`}
                  </p>
                </div>

                {/* Buyer */}
                <div>
                  <p className="text-sm text-base-content/70">{t('buyer')}:</p>
                  <p className="font-semibold">
                    {order.buyer?.name || order.buyer?.email || t('unknown')}
                    {isBuyer && ` (${t('you')})`}
                  </p>
                </div>
              </div>
            </div>
          </div>

          {/* Protection Period */}
          {order.protection_expires_at && (
            <div className="card bg-base-100 shadow">
              <div className="card-body">
                <h2 className="card-title text-lg">{t('protectionPeriod')}</h2>
                <p className="text-sm text-base-content/70">
                  {t('protectionActiveUntil')}:
                </p>
                <p className="font-semibold">
                  {formatDate(order.protection_expires_at)}
                </p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Back Button */}
      <div className="mt-8">
        <Link href={`/${locale}/profile/orders`} className="btn btn-outline">
          ← {tCommon('back')}
        </Link>
      </div>

      {/* Image Gallery Modal */}
      {order?.listing?.images && order.listing.images.length > 0 && (
        <ImageGallery
          images={order.listing.images.map((img: any) => img.url)}
          initialIndex={selectedImageIndex}
          isOpen={galleryOpen}
          onClose={() => setGalleryOpen(false)}
        />
      )}
    </div>
  );
}
