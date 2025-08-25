'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import {
  FiArrowLeft,
  FiPackage,
  FiTruck,
  FiCheckCircle,
  FiXCircle,
  FiClock,
  FiMapPin,
  FiPhone,
  FiMail,
  FiInfo,
} from 'react-icons/fi';
import { useRouter } from 'next/navigation';
import { apiClientAuth } from '@/lib/api-client-auth';

interface ShipmentDetailsClientProps {
  provider: string;
  id: string;
}

interface ShipmentDetails {
  provider: string;
  id: number;
  tracking_number: string;
  status: string;
  status_text: string;
  cod_amount: number;
  weight_kg: number;
  package_contents?: string;
  reference_number?: string;
  registered_at?: string;
  delivered_at?: string;
  failed_reason?: string;
  created_at: string;
  updated_at: string;
  recipient: {
    name: string;
    address: string;
    city: string;
    zip: string;
    phone: string;
    email?: string;
  };
  sender: {
    name: string;
    address: string;
    city: string;
    zip: string;
    phone: string;
    email?: string;
  };
  status_history?: Array<{
    status: string;
    status_text: string;
    timestamp: string;
    location?: string;
  }>;
}

export default function ShipmentDetailsClient({
  provider,
  id,
}: ShipmentDetailsClientProps) {
  const t = useTranslations('admin');
  const router = useRouter();
  const [shipment, setShipment] = useState<ShipmentDetails | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetchShipmentDetails();
  }, [provider, id]);

  const fetchShipmentDetails = async () => {
    try {
      setLoading(true);
      setError(null);

      const result = await apiClientAuth.get(
        `/admin/logistics/shipments/${provider}/${id}`
      );

      if (!result.success) {
        throw new Error(result.error || 'Failed to fetch shipment details');
      }

      setShipment(result.data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'An error occurred');
    } finally {
      setLoading(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'pending':
        return <FiClock className="w-5 h-5" />;
      case 'in_transit':
        return <FiTruck className="w-5 h-5" />;
      case 'delivered':
        return <FiCheckCircle className="w-5 h-5" />;
      case 'failed':
        return <FiXCircle className="w-5 h-5" />;
      default:
        return <FiPackage className="w-5 h-5" />;
    }
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending':
        return 'badge-warning';
      case 'in_transit':
        return 'badge-info';
      case 'delivered':
        return 'badge-success';
      case 'failed':
        return 'badge-error';
      default:
        return 'badge-ghost';
    }
  };

  const formatDate = (dateString?: string) => {
    if (!dateString) return t('logistics.notAvailable');
    const date = new Date(dateString);
    return new Intl.DateTimeFormat('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    }).format(date);
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  if (error) {
    return (
      <div className="alert alert-error">
        <FiXCircle className="w-5 h-5" />
        <span>{error}</span>
      </div>
    );
  }

  if (!shipment) {
    return (
      <div className="alert alert-warning">
        <FiInfo className="w-5 h-5" />
        <span>{t('logistics.shipmentNotFound')}</span>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      {/* Back button and Actions */}
      <div className="flex justify-between items-center">
        <button onClick={() => router.back()} className="btn btn-ghost gap-2">
          <FiArrowLeft />
          {t('common.back')}
        </button>

        <div className="flex gap-2">
          <button className="btn btn-primary">
            {t('logistics.actions.updateStatus')}
          </button>
          <button className="btn btn-ghost">
            {t('logistics.actions.printLabel')}
          </button>
        </div>
      </div>

      {/* Main Info Card */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="flex justify-between items-start">
            <div>
              <h2 className="card-title text-2xl">
                {t('logistics.trackingNumber')}: {shipment.tracking_number}
              </h2>
              <div className="flex items-center gap-4 mt-2">
                <div
                  className={`badge ${getStatusColor(shipment.status)} gap-2`}
                >
                  {getStatusIcon(shipment.status)}
                  {t(`logistics.status.${shipment.status}`)}
                </div>
                <span className="text-sm text-base-content/70">
                  {t('logistics.provider')}: {shipment.provider}
                </span>
              </div>
            </div>

            {shipment.cod_amount > 0 && (
              <div className="text-right">
                <p className="text-sm text-base-content/70">
                  {t('logistics.codAmount')}
                </p>
                <p className="text-2xl font-bold">{shipment.cod_amount} RSD</p>
              </div>
            )}
          </div>

          {shipment.failed_reason && (
            <div className="alert alert-error mt-4">
              <FiXCircle className="w-5 h-5" />
              <div>
                <h3 className="font-bold">{t('logistics.failedReason')}</h3>
                <p>{shipment.failed_reason}</p>
              </div>
            </div>
          )}
        </div>
      </div>

      {/* Two columns layout */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {/* Sender Info */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title">{t('logistics.sender')}</h3>
            <div className="space-y-3">
              <div className="flex items-start gap-3">
                <FiPackage className="w-5 h-5 mt-1 text-base-content/70" />
                <div>
                  <p className="font-semibold">{shipment.sender.name}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <FiMapPin className="w-5 h-5 mt-1 text-base-content/70" />
                <div>
                  <p>{shipment.sender.address}</p>
                  <p>
                    {shipment.sender.zip} {shipment.sender.city}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <FiPhone className="w-5 h-5 text-base-content/70" />
                <p>{shipment.sender.phone}</p>
              </div>

              {shipment.sender.email && (
                <div className="flex items-center gap-3">
                  <FiMail className="w-5 h-5 text-base-content/70" />
                  <p>{shipment.sender.email}</p>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Recipient Info */}
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title">{t('logistics.recipient')}</h3>
            <div className="space-y-3">
              <div className="flex items-start gap-3">
                <FiPackage className="w-5 h-5 mt-1 text-base-content/70" />
                <div>
                  <p className="font-semibold">{shipment.recipient.name}</p>
                </div>
              </div>

              <div className="flex items-start gap-3">
                <FiMapPin className="w-5 h-5 mt-1 text-base-content/70" />
                <div>
                  <p>{shipment.recipient.address}</p>
                  <p>
                    {shipment.recipient.zip} {shipment.recipient.city}
                  </p>
                </div>
              </div>

              <div className="flex items-center gap-3">
                <FiPhone className="w-5 h-5 text-base-content/70" />
                <p>{shipment.recipient.phone}</p>
              </div>

              {shipment.recipient.email && (
                <div className="flex items-center gap-3">
                  <FiMail className="w-5 h-5 text-base-content/70" />
                  <p>{shipment.recipient.email}</p>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>

      {/* Package Details */}
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h3 className="card-title">{t('logistics.packageDetails')}</h3>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <p className="text-sm text-base-content/70">
                {t('logistics.weight')}
              </p>
              <p className="font-semibold">{shipment.weight_kg} kg</p>
            </div>

            {shipment.package_contents && (
              <div>
                <p className="text-sm text-base-content/70">
                  {t('logistics.contents')}
                </p>
                <p className="font-semibold">{shipment.package_contents}</p>
              </div>
            )}

            {shipment.reference_number && (
              <div>
                <p className="text-sm text-base-content/70">
                  {t('logistics.referenceNumber')}
                </p>
                <p className="font-semibold">{shipment.reference_number}</p>
              </div>
            )}

            <div>
              <p className="text-sm text-base-content/70">
                {t('logistics.createdAt')}
              </p>
              <p className="font-semibold">{formatDate(shipment.created_at)}</p>
            </div>
          </div>
        </div>
      </div>

      {/* Timeline / Status History */}
      {shipment.status_history && shipment.status_history.length > 0 && (
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title">{t('logistics.statusHistory')}</h3>
            <div className="space-y-4">
              {shipment.status_history.map((entry, index) => (
                <div key={index} className="flex gap-4">
                  <div className="flex flex-col items-center">
                    <div
                      className={`w-10 h-10 rounded-full flex items-center justify-center ${
                        index === 0
                          ? 'bg-primary text-primary-content'
                          : 'bg-base-300'
                      }`}
                    >
                      {getStatusIcon(entry.status)}
                    </div>
                    {index < shipment.status_history!.length - 1 && (
                      <div className="w-0.5 h-16 bg-base-300 mt-2"></div>
                    )}
                  </div>

                  <div className="flex-1 pb-8">
                    <div className="flex items-center gap-2">
                      <p className="font-semibold">{entry.status_text}</p>
                      <span className="text-sm text-base-content/70">
                        {formatDate(entry.timestamp)}
                      </span>
                    </div>
                    {entry.location && (
                      <p className="text-sm text-base-content/70 mt-1">
                        <FiMapPin className="inline w-4 h-4 mr-1" />
                        {entry.location}
                      </p>
                    )}
                  </div>
                </div>
              ))}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
