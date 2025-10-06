'use client';

import { useState, useEffect } from 'react';
import { apiClient } from '@/services/api-client';

interface TestShipmentRequest {
  recipient_name: string;
  recipient_phone: string;
  recipient_email: string;
  recipient_city: string;
  recipient_address: string;
  recipient_zip: string;
  sender_name: string;
  sender_phone: string;
  sender_email: string;
  sender_city: string;
  sender_address: string;
  sender_zip: string;
  weight: number;
  content: string;
  cod_amount: number;
  insured_value: number;
  services: string;
  delivery_method: string;
  payment_method: string;
}

interface TestShipmentResponse {
  success: boolean;
  tracking_number?: string;
  manifest_id?: number;
  shipment_id?: number;
  external_id?: string;
  cost?: number;
  errors?: string[];
  request_data?: any;
  response_data?: any;
  created_at?: string;
  processing_time_ms?: number;
}

interface Config {
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
  delivery_methods: Array<{ code: string; name: string }>;
  payment_methods: Array<{ code: string; name: string }>;
  services: Array<{ code: string; name: string }>;
}

export default function PostExpressTestPage() {
  const [formData, setFormData] = useState<TestShipmentRequest>({
    recipient_name: '',
    recipient_phone: '',
    recipient_email: '',
    recipient_city: '',
    recipient_address: '',
    recipient_zip: '',
    sender_name: '',
    sender_phone: '',
    sender_email: '',
    sender_city: '',
    sender_address: '',
    sender_zip: '',
    weight: 500,
    content: 'Test paket za SVETU',
    cod_amount: 0,
    insured_value: 0,
    services: 'PNA',
    delivery_method: 'K',
    payment_method: 'POF',
  });

  const [config, setConfig] = useState<Config | null>(null);
  const [result, setResult] = useState<TestShipmentResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    loadConfig();
  }, []);

  const loadConfig = async () => {
    try {
      const response = await apiClient.get('/postexpress/test/config');
      if (response.data.success && response.data.data) {
        const cfg = response.data.data;
        setConfig(cfg);

        // Заполняем форму дефолтными значениями
        setFormData({
          ...formData,
          sender_name: cfg.default_sender.name,
          sender_phone: cfg.default_sender.phone,
          sender_email: cfg.default_sender.email,
          sender_city: cfg.default_sender.city,
          sender_address: cfg.default_sender.address,
          sender_zip: cfg.default_sender.zip,
          recipient_name: cfg.default_recipient.name,
          recipient_phone: cfg.default_recipient.phone,
          recipient_email: cfg.default_recipient.email,
          recipient_city: cfg.default_recipient.city,
          recipient_address: cfg.default_recipient.address,
          recipient_zip: cfg.default_recipient.zip,
        });
      }
    } catch (err: any) {
      console.error('Failed to load config:', err);
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await apiClient.post('/postexpress/test/shipment', formData);

      if (response.data.success && response.data.data) {
        setResult(response.data.data);
      } else {
        setError(response.data.message || 'Failed to create test shipment');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || err.message || 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>) => {
    const { name, value, type } = e.target;
    setFormData({
      ...formData,
      [name]: type === 'number' ? Number(value) : value,
    });
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">📦 Post Express - Visual Testing</h1>
          <p className="mt-2 text-gray-600">
            Полное визуальное тестирование Post Express WSP API интеграции
          </p>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Форма */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-bold mb-6">Параметры отправления</h2>

            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Получатель */}
              <div>
                <h3 className="text-lg font-semibold mb-4 text-blue-700">Получатель</h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      ФИО *
                    </label>
                    <input
                      type="text"
                      name="recipient_name"
                      value={formData.recipient_name}
                      onChange={handleChange}
                      required
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      placeholder="Petar Petrović"
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Телефон *
                      </label>
                      <input
                        type="text"
                        name="recipient_phone"
                        value={formData.recipient_phone}
                        onChange={handleChange}
                        required
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        placeholder="0641234567"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Email
                      </label>
                      <input
                        type="email"
                        name="recipient_email"
                        value={formData.recipient_email}
                        onChange={handleChange}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        placeholder="petar@example.com"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-4">
                    <div className="col-span-2">
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Город *
                      </label>
                      <input
                        type="text"
                        name="recipient_city"
                        value={formData.recipient_city}
                        onChange={handleChange}
                        required
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        placeholder="Beograd"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Индекс *
                      </label>
                      <input
                        type="text"
                        name="recipient_zip"
                        value={formData.recipient_zip}
                        onChange={handleChange}
                        required
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                        placeholder="11000"
                      />
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Адрес
                    </label>
                    <input
                      type="text"
                      name="recipient_address"
                      value={formData.recipient_address}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
                      placeholder="Takovska 2"
                    />
                  </div>
                </div>
              </div>

              {/* Отправитель */}
              <div>
                <h3 className="text-lg font-semibold mb-4 text-green-700">Отправитель</h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Название *
                    </label>
                    <input
                      type="text"
                      name="sender_name"
                      value={formData.sender_name}
                      onChange={handleChange}
                      required
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
                      placeholder="Sve Tu d.o.o."
                    />
                  </div>

                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Телефон *
                      </label>
                      <input
                        type="text"
                        name="sender_phone"
                        value={formData.sender_phone}
                        onChange={handleChange}
                        required
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
                        placeholder="0641234567"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Email
                      </label>
                      <input
                        type="email"
                        name="sender_email"
                        value={formData.sender_email}
                        onChange={handleChange}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
                        placeholder="b2b@svetu.rs"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-4">
                    <div className="col-span-2">
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Город *
                      </label>
                      <input
                        type="text"
                        name="sender_city"
                        value={formData.sender_city}
                        onChange={handleChange}
                        required
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
                        placeholder="Beograd"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Индекс *
                      </label>
                      <input
                        type="text"
                        name="sender_zip"
                        value={formData.sender_zip}
                        onChange={handleChange}
                        required
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
                        placeholder="11000"
                      />
                    </div>
                  </div>

                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Адрес
                    </label>
                    <input
                      type="text"
                      name="sender_address"
                      value={formData.sender_address}
                      onChange={handleChange}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500 focus:border-transparent"
                      placeholder="Bulevar kralja Aleksandra 73"
                    />
                  </div>
                </div>
              </div>

              {/* Параметры посылки */}
              <div>
                <h3 className="text-lg font-semibold mb-4 text-purple-700">Параметры посылки</h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      Содержимое *
                    </label>
                    <textarea
                      name="content"
                      value={formData.content}
                      onChange={handleChange}
                      required
                      rows={2}
                      className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                      placeholder="Test paket za SVETU"
                    />
                  </div>

                  <div className="grid grid-cols-3 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Вес (г) *
                      </label>
                      <input
                        type="number"
                        name="weight"
                        value={formData.weight}
                        onChange={handleChange}
                        required
                        min="1"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                        placeholder="500"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Наложенный платеж (RSD)
                      </label>
                      <input
                        type="number"
                        name="cod_amount"
                        value={formData.cod_amount}
                        onChange={handleChange}
                        min="0"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                        placeholder="0"
                      />
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Ценность (RSD)
                      </label>
                      <input
                        type="number"
                        name="insured_value"
                        value={formData.insured_value}
                        onChange={handleChange}
                        min="0"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                        placeholder="0"
                      />
                    </div>
                  </div>

                  <div className="grid grid-cols-3 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Услуги
                      </label>
                      <select
                        name="services"
                        value={formData.services}
                        onChange={handleChange}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                      >
                        {config?.services.map((s) => (
                          <option key={s.code} value={s.code}>
                            {s.name}
                          </option>
                        ))}
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Способ доставки
                      </label>
                      <select
                        name="delivery_method"
                        value={formData.delivery_method}
                        onChange={handleChange}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                      >
                        {config?.delivery_methods.map((m) => (
                          <option key={m.code} value={m.code}>
                            {m.name}
                          </option>
                        ))}
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Способ оплаты
                      </label>
                      <select
                        name="payment_method"
                        value={formData.payment_method}
                        onChange={handleChange}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                      >
                        {config?.payment_methods.map((m) => (
                          <option key={m.code} value={m.code}>
                            {m.name}
                          </option>
                        ))}
                      </select>
                    </div>
                  </div>
                </div>
              </div>

              {/* Кнопка */}
              <button
                type="submit"
                disabled={loading}
                className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-3 px-6 rounded-lg transition-colors duration-200 flex items-center justify-center"
              >
                {loading ? (
                  <>
                    <svg className="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                      <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                      <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                    </svg>
                    Создание...
                  </>
                ) : (
                  <>
                    📦 Создать тестовое отправление
                  </>
                )}
              </button>
            </form>
          </div>

          {/* Результаты */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-bold mb-6">Результаты</h2>

            {!result && !error && (
              <div className="text-center py-12 text-gray-400">
                <svg className="mx-auto h-24 w-24 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                </svg>
                <p className="text-lg">Заполните форму и создайте тестовое отправление</p>
              </div>
            )}

            {error && (
              <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-4">
                <div className="flex items-start">
                  <svg className="h-6 w-6 text-red-500 mr-3 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  <div>
                    <h3 className="text-red-800 font-semibold">Ошибка</h3>
                    <p className="text-red-700 mt-1">{error}</p>
                  </div>
                </div>
              </div>
            )}

            {result && result.success && (
              <div className="space-y-4">
                {/* Success banner */}
                <div className="bg-green-50 border-l-4 border-green-500 p-4">
                  <div className="flex items-start">
                    <svg className="h-6 w-6 text-green-500 mr-3 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <div>
                      <h3 className="text-green-800 font-semibold">Отправление создано успешно!</h3>
                      <p className="text-green-700 text-sm mt-1">Время обработки: {result.processing_time_ms}ms</p>
                    </div>
                  </div>
                </div>

                {/* Tracking info */}
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-blue-50 p-4 rounded-lg">
                    <p className="text-sm text-blue-600 font-medium mb-1">Трек-номер</p>
                    <p className="text-2xl font-bold text-blue-900">{result.tracking_number}</p>
                  </div>

                  <div className="bg-purple-50 p-4 rounded-lg">
                    <p className="text-sm text-purple-600 font-medium mb-1">Стоимость</p>
                    <p className="text-2xl font-bold text-purple-900">{result.cost} RSD</p>
                  </div>
                </div>

                {/* IDs */}
                <div className="grid grid-cols-3 gap-4">
                  <div className="bg-gray-50 p-3 rounded">
                    <p className="text-xs text-gray-600 mb-1">Manifest ID</p>
                    <p className="font-mono text-sm font-semibold">{result.manifest_id}</p>
                  </div>

                  <div className="bg-gray-50 p-3 rounded">
                    <p className="text-xs text-gray-600 mb-1">Shipment ID</p>
                    <p className="font-mono text-sm font-semibold">{result.shipment_id}</p>
                  </div>

                  <div className="bg-gray-50 p-3 rounded">
                    <p className="text-xs text-gray-600 mb-1">External ID</p>
                    <p className="font-mono text-xs font-semibold">{result.external_id}</p>
                  </div>
                </div>

                {/* Request data */}
                <div>
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">Отправленные данные</h3>
                  <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg text-xs overflow-auto max-h-64">
                    {JSON.stringify(result.request_data, null, 2)}
                  </pre>
                </div>

                {/* Response data */}
                <div>
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">Ответ API</h3>
                  <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg text-xs overflow-auto max-h-64">
                    {JSON.stringify(result.response_data, null, 2)}
                  </pre>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
