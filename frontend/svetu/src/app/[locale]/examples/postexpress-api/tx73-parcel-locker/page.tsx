'use client';

import { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import {
  BuildingStorefrontIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon,
  MapPinIcon,
  QrCodeIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';

interface TestConfig {
  api_available: boolean;
  test_mode: boolean;
  default_sender: {
    name: string;
    phone: string;
    email: string;
    city: string;
    address: string;
    zip: string;
  };
  default_recipient: {
    name: string;
    phone: string;
    email: string;
    city: string;
    address: string;
    zip: string;
  };
  id_rukovanje_options: Array<{
    id: number;
    name: string;
    description: string;
  }>;
}

interface ShipmentResult {
  success: boolean;
  tracking_number?: string;
  manifest_id?: number;
  shipment_id?: number;
  external_id?: string;
  cost?: number;
  errors?: string[];
  created_at?: string;
  processing_time_ms?: number;
  request_data?: any;
  response_data?: any;
}

export default function TX73ParcelLockerPage() {
  const t = useTranslations('postexpressTest.tx73');

  const [config, setConfig] = useState<TestConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState<ShipmentResult | null>(null);

  // Form state
  const [recipientName, setRecipientName] = useState('');
  const [recipientPhone, setRecipientPhone] = useState('');
  const [recipientEmail, setRecipientEmail] = useState('');

  const [senderName, setSenderName] = useState('');
  const [senderPhone, setSenderPhone] = useState('');
  const [senderEmail, setSenderEmail] = useState('');
  const [senderCity, setSenderCity] = useState('');
  const [senderAddress, setSenderAddress] = useState('');
  const [senderZip, setSenderZip] = useState('');

  const [weight, setWeight] = useState(500);
  const [content, setContent] = useState('');
  const [insuredValue, setInsuredValue] = useState(0);
  const [idRukovanje, setIdRukovanje] = useState(85); // Default for parcel locker
  const [parcelLockerCode, setParcelLockerCode] = useState('');

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const response = await apiClient.get('/postexpress/test/config');
      if (response.data?.data) {
        setConfig(response.data.data);

        // Set default values
        const sender = response.data.data.default_sender;
        setSenderName(sender.name);
        setSenderPhone(sender.phone);
        setSenderEmail(sender.email);
        setSenderCity(sender.city);
        setSenderAddress(sender.address);
        setSenderZip(sender.zip);

        const recipient = response.data.data.default_recipient;
        setRecipientName(recipient.name);
        setRecipientPhone(recipient.phone);
        setRecipientEmail(recipient.email);

        setContent('Test paket - Parcel Locker (паккетомат)');
        setParcelLockerCode('PAK001'); // Default parcel locker code
      }
      setLoading(false);
    } catch (error) {
      console.error('Failed to load config:', error);
      setLoading(false);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setSubmitting(true);
    setResult(null);

    try {
      const payload = {
        recipient_name: recipientName,
        recipient_phone: recipientPhone,
        recipient_email: recipientEmail,

        sender_name: senderName,
        sender_phone: senderPhone,
        sender_email: senderEmail,
        sender_city: senderCity,
        sender_address: senderAddress,
        sender_zip: senderZip,

        weight,
        content,
        insured_value: insuredValue,

        delivery_type: 'parcel_locker',
        id_rukovanje: idRukovanje,
        parcel_locker_code: parcelLockerCode,
      };

      const response = await apiClient.post(
        '/postexpress/test/shipment',
        payload
      );

      if (response.data?.data) {
        setResult(response.data.data);
      }
    } catch (error: any) {
      setResult({
        success: false,
        errors: [
          error.response?.data?.message || error.message || 'Unknown error',
        ],
      });
    } finally {
      setSubmitting(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="loading loading-spinner loading-lg"></div>
        <p className="ml-4">{t('loadingConfig')}</p>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-gradient-to-r from-purple-600 to-indigo-700 text-white py-8">
        <div className="container mx-auto px-4">
          <Link
            href="/examples/postexpress-api"
            className="text-sm hover:underline mb-2 inline-block opacity-80"
          >
            ← {t('back')}
          </Link>
          <h1 className="text-4xl font-bold mb-2 flex items-center gap-3">
            <BuildingStorefrontIcon className="w-10 h-10" />
            {t('titleParcelLocker')}
          </h1>
          <p className="text-xl opacity-90">{t('descriptionParcelLocker')}</p>
          <div className="mt-4 flex gap-2 flex-wrap">
            <div className="badge badge-success badge-lg gap-2">
              <CheckCircleIcon className="w-4 h-4" />
              API Ready
            </div>
            <div className="badge badge-warning badge-lg gap-2">
              <ClockIcon className="w-4 h-4" />
              Test Mode
            </div>
            <div className="badge badge-info badge-lg gap-2">
              <QrCodeIcon className="w-4 h-4" />
              Parcel Locker
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
          {/* Form */}
          <div className="lg:col-span-2">
            <form
              onSubmit={handleSubmit}
              className="card bg-base-100 shadow-xl"
            >
              <div className="card-body">
                <h2 className="card-title text-2xl mb-4">{t('formTitle')}</h2>

                {/* Recipient Section */}
                <div className="mb-6">
                  <h3 className="font-semibold text-lg mb-3 flex items-center gap-2">
                    <MapPinIcon className="w-5 h-5" />
                    {t('recipient.title')}
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          {t('recipient.name')}
                        </span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={recipientName}
                        onChange={(e) => setRecipientName(e.target.value)}
                        placeholder={t('recipient.namePlaceholder')}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          {t('recipient.phone')}
                        </span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={recipientPhone}
                        onChange={(e) => setRecipientPhone(e.target.value)}
                        placeholder={t('recipient.phonePlaceholder')}
                        required
                      />
                    </div>
                    <div className="form-control md:col-span-2">
                      <label className="label">
                        <span className="label-text">
                          {t('recipient.email')}
                        </span>
                      </label>
                      <input
                        type="email"
                        className="input input-bordered"
                        value={recipientEmail}
                        onChange={(e) => setRecipientEmail(e.target.value)}
                        placeholder={t('recipient.emailPlaceholder')}
                        required
                      />
                    </div>
                  </div>
                </div>

                <div className="divider"></div>

                {/* Parcel Locker Section */}
                <div className="mb-6">
                  <h3 className="font-semibold text-lg mb-3 flex items-center gap-2">
                    <BuildingStorefrontIcon className="w-5 h-5 text-info" />
                    {t('deliveryTypes.parcelLocker.name')}
                  </h3>
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text font-bold text-info">
                        {t('shipment.parcelLockerCode')}
                      </span>
                    </label>
                    <input
                      type="text"
                      className="input input-bordered input-info font-mono text-lg"
                      value={parcelLockerCode}
                      onChange={(e) => setParcelLockerCode(e.target.value)}
                      placeholder={t('shipment.parcelLockerCodePlaceholder')}
                      required
                    />
                  </div>
                </div>

                <div className="divider"></div>

                {/* Sender Section */}
                <div className="mb-6">
                  <h3 className="font-semibold text-lg mb-3">
                    {t('sender.title')}
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">{t('sender.name')}</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={senderName}
                        onChange={(e) => setSenderName(e.target.value)}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">{t('sender.phone')}</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={senderPhone}
                        onChange={(e) => setSenderPhone(e.target.value)}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">{t('sender.email')}</span>
                      </label>
                      <input
                        type="email"
                        className="input input-bordered"
                        value={senderEmail}
                        onChange={(e) => setSenderEmail(e.target.value)}
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">{t('sender.city')}</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={senderCity}
                        onChange={(e) => setSenderCity(e.target.value)}
                        required
                      />
                    </div>
                    <div className="form-control md:col-span-2">
                      <label className="label">
                        <span className="label-text">
                          {t('sender.address')}
                        </span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={senderAddress}
                        onChange={(e) => setSenderAddress(e.target.value)}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">{t('sender.zip')}</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={senderZip}
                        onChange={(e) => setSenderZip(e.target.value)}
                        required
                      />
                    </div>
                  </div>
                </div>

                <div className="divider"></div>

                {/* Package Details */}
                <div className="mb-6">
                  <h3 className="font-semibold text-lg mb-3">
                    {t('shipment.title')}
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          {t('shipment.weight')}
                        </span>
                      </label>
                      <input
                        type="number"
                        className="input input-bordered"
                        value={weight}
                        onChange={(e) => setWeight(Number(e.target.value))}
                        placeholder={t('shipment.weightPlaceholder')}
                        min="100"
                        max="20000"
                        step="1"
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          {t('shipment.content')}
                        </span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                        placeholder={t('shipment.contentPlaceholder')}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          {t('shipment.insuredValue')}
                        </span>
                      </label>
                      <input
                        type="number"
                        className="input input-bordered"
                        value={insuredValue}
                        onChange={(e) =>
                          setInsuredValue(Number(e.target.value))
                        }
                        placeholder={t('shipment.insuredValuePlaceholder')}
                        min="0"
                        step="1"
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          {t('shipment.idRukovanje')}
                        </span>
                      </label>
                      <select
                        className="select select-bordered"
                        value={idRukovanje}
                        onChange={(e) => setIdRukovanje(Number(e.target.value))}
                      >
                        {config?.id_rukovanje_options.map((option) => (
                          <option key={option.id} value={option.id}>
                            {option.id} - {option.description}
                          </option>
                        ))}
                      </select>
                    </div>
                  </div>
                </div>

                {/* Parcel Locker Info Alert */}
                <div className="alert alert-info mb-6">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                    className="stroke-current shrink-0 w-6 h-6"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    ></path>
                  </svg>
                  <div>
                    <h3 className="font-bold">
                      {t('deliveryTypes.parcelLocker.name')}
                    </h3>
                    <div className="text-xs">
                      {t('deliveryTypes.parcelLocker.description')}
                    </div>
                  </div>
                </div>

                {/* Submit Button */}
                <div className="card-actions justify-end">
                  <button
                    type="submit"
                    className="btn btn-info btn-lg"
                    disabled={submitting}
                  >
                    {submitting ? (
                      <>
                        <span className="loading loading-spinner"></span>
                        {t('creating')}
                      </>
                    ) : (
                      <>
                        <BuildingStorefrontIcon className="w-5 h-5" />
                        {t('createButton')}
                      </>
                    )}
                  </button>
                </div>
              </div>
            </form>
          </div>

          {/* Result Panel */}
          <div className="lg:col-span-1">
            <div className="card bg-base-100 shadow-xl sticky top-4">
              <div className="card-body">
                <h2 className="card-title">{t('result.title')}</h2>

                {!result && (
                  <div className="text-center py-8 text-base-content/50">
                    <BuildingStorefrontIcon className="w-16 h-16 mx-auto mb-4 opacity-30" />
                    <p>{t('result.waitingForResult')}</p>
                  </div>
                )}

                {result && (
                  <div className="space-y-4">
                    {result.success ? (
                      <>
                        <div className="alert alert-success">
                          <CheckCircleIcon className="w-6 h-6" />
                          <span>{t('result.success')}</span>
                        </div>

                        <div className="space-y-2">
                          <div className="flex justify-between">
                            <span className="text-sm">
                              {t('result.trackingNumber')}:
                            </span>
                            <span className="font-mono font-bold">
                              {result.tracking_number}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">
                              {t('result.manifestId')}:
                            </span>
                            <span className="font-bold">
                              {result.manifest_id}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">
                              {t('result.shipmentId')}:
                            </span>
                            <span className="font-bold">
                              {result.shipment_id}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">
                              {t('result.externalId')}:
                            </span>
                            <span className="font-mono text-xs">
                              {result.external_id}
                            </span>
                          </div>
                          <div className="flex justify-between items-center bg-info/10 p-2 rounded">
                            <span className="text-sm font-bold">
                              Parcel Locker:
                            </span>
                            <span className="font-mono text-info text-lg font-bold">
                              {parcelLockerCode}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">Delivery Cost:</span>
                            <span className="font-bold text-warning">
                              {result.cost} RSD
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">
                              {t('result.processingTime')}:
                            </span>
                            <span className="text-xs">
                              {result.processing_time_ms}ms
                            </span>
                          </div>
                        </div>

                        {result.request_data && (
                          <div className="collapse collapse-arrow bg-base-200">
                            <input type="checkbox" />
                            <div className="collapse-title font-medium text-sm">
                              Request Data
                            </div>
                            <div className="collapse-content">
                              <pre className="text-xs overflow-auto">
                                {JSON.stringify(result.request_data, null, 2)}
                              </pre>
                            </div>
                          </div>
                        )}

                        {result.response_data && (
                          <div className="collapse collapse-arrow bg-base-200">
                            <input type="checkbox" />
                            <div className="collapse-title font-medium text-sm">
                              Response Data
                            </div>
                            <div className="collapse-content">
                              <pre className="text-xs overflow-auto">
                                {JSON.stringify(result.response_data, null, 2)}
                              </pre>
                            </div>
                          </div>
                        )}
                      </>
                    ) : (
                      <>
                        <div className="alert alert-error">
                          <XCircleIcon className="w-6 h-6" />
                          <span>{t('result.failed')}</span>
                        </div>

                        {result.errors && result.errors.length > 0 && (
                          <div className="space-y-2">
                            {result.errors.map((error, index) => (
                              <div key={index} className="text-sm text-error">
                                • {error}
                              </div>
                            ))}
                          </div>
                        )}
                      </>
                    )}
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
