'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/services/api-client';
import {
  TruckIcon,
  CurrencyDollarIcon,
  MapPinIcon,
  BuildingStorefrontIcon,
  CheckCircleIcon,
  XCircleIcon,
  ClockIcon,
} from '@heroicons/react/24/outline';

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
  delivery_types: Array<{
    code: string;
    name: string;
    description: string;
  }>;
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

export default function PostExpressTestPage() {
  const [config, setConfig] = useState<TestConfig | null>(null);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [result, setResult] = useState<ShipmentResult | null>(null);

  // Form state
  const [recipientName, setRecipientName] = useState('');
  const [recipientPhone, setRecipientPhone] = useState('');
  const [recipientEmail, setRecipientEmail] = useState('');
  const [recipientCity, setRecipientCity] = useState('');
  const [recipientAddress, setRecipientAddress] = useState('');
  const [recipientZip, setRecipientZip] = useState('');

  const [senderName, setSenderName] = useState('');
  const [senderPhone, setSenderPhone] = useState('');
  const [senderEmail, setSenderEmail] = useState('');
  const [senderCity, setSenderCity] = useState('');
  const [senderAddress, setSenderAddress] = useState('');
  const [senderZip, setSenderZip] = useState('');

  const [weight, setWeight] = useState(500);
  const [content, setContent] = useState('');
  const [codAmount, setCodAmount] = useState(0);
  const [insuredValue, setInsuredValue] = useState(0);

  const [deliveryType, setDeliveryType] = useState('standard');
  const [idRukovanje, setIdRukovanje] = useState(29);
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
        setRecipientCity(recipient.city);
        setRecipientAddress(recipient.address);
        setRecipientZip(recipient.zip);

        setContent('Test paket za SVETU');
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
        recipient_city: recipientCity,
        recipient_address: recipientAddress,
        recipient_zip: recipientZip,

        sender_name: senderName,
        sender_phone: senderPhone,
        sender_email: senderEmail,
        sender_city: senderCity,
        sender_address: senderAddress,
        sender_zip: senderZip,

        weight,
        content,
        cod_amount: codAmount,
        insured_value: insuredValue,

        delivery_type: deliveryType,
        id_rukovanje: idRukovanje,
        parcel_locker_code:
          deliveryType === 'parcel_locker' ? parcelLockerCode : undefined,
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

  const handleDeliveryTypeChange = (type: string) => {
    setDeliveryType(type);

    // Adjust IdRukovanje based on delivery type
    if (type === 'parcel_locker') {
      setIdRukovanje(85);
    } else if (type === 'cod') {
      setIdRukovanje(29); // or any other suitable option
      if (codAmount === 0) {
        setCodAmount(5000); // Set default COD amount
      }
    } else {
      setIdRukovanje(29);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-gradient-to-r from-blue-600 to-indigo-600 text-white py-8">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl font-bold mb-2">
            Post Express - Тестирование
          </h1>
          <p className="text-xl opacity-90">
            Создание тестовых откупных пошильок и доставки в паккетоматы
          </p>
          <div className="mt-4 flex gap-2 flex-wrap">
            <div className="badge badge-success badge-lg gap-2">
              <CheckCircleIcon className="w-4 h-4" />
              API Ready
            </div>
            <div className="badge badge-warning badge-lg gap-2">
              <ClockIcon className="w-4 h-4" />
              Test Mode
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
                <h2 className="card-title text-2xl mb-4">
                  Создать тестовую пошильку
                </h2>

                {/* Delivery Type Selection */}
                <div className="form-control mb-6">
                  <label className="label">
                    <span className="label-text font-semibold">
                      Тип доставки
                    </span>
                  </label>
                  <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
                    {config?.delivery_types.map((type) => (
                      <label
                        key={type.code}
                        className={`card bordered cursor-pointer transition-all ${
                          deliveryType === type.code
                            ? 'ring-2 ring-primary bg-primary/10'
                            : 'hover:bg-base-200'
                        }`}
                      >
                        <input
                          type="radio"
                          name="delivery_type"
                          className="hidden"
                          checked={deliveryType === type.code}
                          onChange={() => handleDeliveryTypeChange(type.code)}
                        />
                        <div className="card-body p-4 text-center">
                          {type.code === 'standard' && (
                            <TruckIcon className="w-8 h-8 mx-auto mb-2 text-primary" />
                          )}
                          {type.code === 'cod' && (
                            <CurrencyDollarIcon className="w-8 h-8 mx-auto mb-2 text-success" />
                          )}
                          {type.code === 'parcel_locker' && (
                            <BuildingStorefrontIcon className="w-8 h-8 mx-auto mb-2 text-info" />
                          )}
                          <h3 className="font-bold text-sm">{type.name}</h3>
                          <p className="text-xs text-base-content/70">
                            {type.description}
                          </p>
                        </div>
                      </label>
                    ))}
                  </div>
                </div>

                {/* Recipient Section */}
                <div className="mb-6">
                  <h3 className="font-semibold text-lg mb-3">Получатель</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Имя</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={recipientName}
                        onChange={(e) => setRecipientName(e.target.value)}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Телефон</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={recipientPhone}
                        onChange={(e) => setRecipientPhone(e.target.value)}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Email</span>
                      </label>
                      <input
                        type="email"
                        className="input input-bordered"
                        value={recipientEmail}
                        onChange={(e) => setRecipientEmail(e.target.value)}
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Город</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={recipientCity}
                        onChange={(e) => setRecipientCity(e.target.value)}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Адрес</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={recipientAddress}
                        onChange={(e) => setRecipientAddress(e.target.value)}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Индекс</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={recipientZip}
                        onChange={(e) => setRecipientZip(e.target.value)}
                        required
                      />
                    </div>
                  </div>
                </div>

                <div className="divider"></div>

                {/* Sender Section */}
                <div className="mb-6">
                  <h3 className="font-semibold text-lg mb-3">Отправитель</h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Имя/Компания</span>
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
                        <span className="label-text">Телефон</span>
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
                        <span className="label-text">Email</span>
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
                        <span className="label-text">Город</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={senderCity}
                        onChange={(e) => setSenderCity(e.target.value)}
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Адрес</span>
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
                        <span className="label-text">Индекс</span>
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
                    Параметры посылки
                  </h3>
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Вес (грамм)</span>
                      </label>
                      <input
                        type="number"
                        className="input input-bordered"
                        value={weight}
                        onChange={(e) => setWeight(Number(e.target.value))}
                        min="100"
                        max="20000"
                        step="100"
                        required
                      />
                    </div>
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Содержимое</span>
                      </label>
                      <input
                        type="text"
                        className="input input-bordered"
                        value={content}
                        onChange={(e) => setContent(e.target.value)}
                        required
                      />
                    </div>
                    {deliveryType === 'cod' && (
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Сумма откупа (RSD)</span>
                        </label>
                        <input
                          type="number"
                          className="input input-bordered input-success"
                          value={codAmount}
                          onChange={(e) => setCodAmount(Number(e.target.value))}
                          min="0"
                          step="100"
                        />
                      </div>
                    )}
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          Объявленная ценность (RSD)
                        </span>
                      </label>
                      <input
                        type="number"
                        className="input input-bordered"
                        value={insuredValue}
                        onChange={(e) =>
                          setInsuredValue(Number(e.target.value))
                        }
                        min="0"
                        step="100"
                      />
                    </div>
                    {deliveryType === 'parcel_locker' && (
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Код паккетомата</span>
                        </label>
                        <input
                          type="text"
                          className="input input-bordered input-info"
                          value={parcelLockerCode}
                          onChange={(e) => setParcelLockerCode(e.target.value)}
                          placeholder="Например: PAK001"
                        />
                      </div>
                    )}
                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">ID Rukovanje</span>
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

                {/* Submit Button */}
                <div className="card-actions justify-end">
                  <button
                    type="submit"
                    className="btn btn-primary btn-lg"
                    disabled={submitting}
                  >
                    {submitting ? (
                      <>
                        <span className="loading loading-spinner"></span>
                        Создание...
                      </>
                    ) : (
                      <>
                        <CheckCircleIcon className="w-5 h-5" />
                        Создать пошильку
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
                <h2 className="card-title">Результат</h2>

                {!result && (
                  <div className="text-center py-8 text-base-content/50">
                    <MapPinIcon className="w-16 h-16 mx-auto mb-4 opacity-30" />
                    <p>Результат появится здесь после создания пошильки</p>
                  </div>
                )}

                {result && (
                  <div className="space-y-4">
                    {result.success ? (
                      <>
                        <div className="alert alert-success">
                          <CheckCircleIcon className="w-6 h-6" />
                          <span>Пошилька успешно создана!</span>
                        </div>

                        <div className="space-y-2">
                          <div className="flex justify-between">
                            <span className="text-sm">Tracking Number:</span>
                            <span className="font-mono font-bold">
                              {result.tracking_number}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">Manifest ID:</span>
                            <span className="font-bold">
                              {result.manifest_id}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">Shipment ID:</span>
                            <span className="font-bold">
                              {result.shipment_id}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">External ID:</span>
                            <span className="font-mono text-xs">
                              {result.external_id}
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">Cost:</span>
                            <span className="font-bold text-success">
                              {result.cost} RSD
                            </span>
                          </div>
                          <div className="flex justify-between">
                            <span className="text-sm">Processing Time:</span>
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
                          <span>Ошибка при создании пошильки</span>
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
