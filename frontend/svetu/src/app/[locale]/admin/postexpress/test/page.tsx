'use client';

import { useState, useEffect, useCallback } from 'react';
import { useTranslations } from 'next-intl';
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
  id_rukovanje?: number;
  parcel_locker_code?: string;
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
  const t = useTranslations('admin.postexpressTest');
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
    id_rukovanje: 29,
    parcel_locker_code: '',
  });

  const [config, setConfig] = useState<Config | null>(null);
  const [result, setResult] = useState<TestShipmentResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [selectedScenario, setSelectedScenario] = useState<string | null>(null);

  // ============== NEW: TX 3, 4, 6, 9, 11 states ==============

  // TX 3 - GetSettlements
  const [tx3ModalOpen, setTx3ModalOpen] = useState(false);
  const [tx3Query, setTx3Query] = useState('Beograd');
  const [tx3Settlements, setTx3Settlements] = useState<any[]>([]);
  const [tx3Loading, setTx3Loading] = useState(false);
  const [tx3Error, setTx3Error] = useState<string | null>(null);

  // TX 4 - GetStreets
  const [tx4ModalOpen, setTx4ModalOpen] = useState(false);
  const [tx4SettlementId, setTx4SettlementId] = useState('100001');
  const [tx4Query, setTx4Query] = useState('Takovska');
  const [tx4Streets, setTx4Streets] = useState<any[]>([]);
  const [tx4Loading, setTx4Loading] = useState(false);
  const [tx4Error, setTx4Error] = useState<string | null>(null);

  // TX 6 - ValidateAddress
  const [tx6ModalOpen, setTx6ModalOpen] = useState(false);
  const [tx6Request, setTx6Request] = useState({
    id_naselje: 100001,
    id_ulica: 1186,
    broj: '2',
    postanski_broj: '11000',
  });
  const [tx6Response, setTx6Response] = useState<any>(null);
  const [tx6Loading, setTx6Loading] = useState(false);
  const [tx6Error, setTx6Error] = useState<string | null>(null);

  // TX 9 - CheckServiceAvailability
  const [tx9ModalOpen, setTx9ModalOpen] = useState(false);
  const [tx9Request, setTx9Request] = useState({
    id_rukovanje: 71,
    postanski_broj_odlaska: '11000',
    postanski_broj_dolaska: '21000',
  });
  const [tx9Response, setTx9Response] = useState<any>(null);
  const [tx9Loading, setTx9Loading] = useState(false);
  const [tx9Error, setTx9Error] = useState<string | null>(null);

  // TX 11 - CalculatePostage
  const [tx11ModalOpen, setTx11ModalOpen] = useState(false);
  const [tx11Request, setTx11Request] = useState({
    id_rukovanje: 71,
    postanski_broj_odlaska: '11000',
    postanski_broj_dolaska: '21000',
    masa: 500,
    otkupnina: 0,
    vrednost: 0,
    posebne_usluge: 'PNA',
  });
  const [tx11Response, setTx11Response] = useState<any>(null);
  const [tx11Loading, setTx11Loading] = useState(false);
  const [tx11Error, setTx11Error] = useState<string | null>(null);

  // TX 73 - B2B Manifest (CreateShipmentViaManifest)
  const [tx73ModalOpen, setTx73ModalOpen] = useState(false);
  const [tx73Request, setTx73Request] = useState({
    recipient_name: 'Marko Marković',
    recipient_phone: '+381641234567',
    recipient_email: 'marko@example.com',
    recipient_city: 'Beograd',
    recipient_address: 'Takovska 2',
    recipient_zip: '11000',
    weight: 500,
    content: 'Test paket - TX 73 test',
    cod_amount: 0,
    insured_value: 0,
    services: 'PNA',
  });
  const [tx73Response, setTx73Response] = useState<any>(null);
  const [tx73Loading, setTx73Loading] = useState(false);
  const [tx73Error, setTx73Error] = useState<string | null>(null);

  // Predefined test scenarios (recipient data only, sender will be preserved from config)
  const testScenarios: Record<string, Partial<TestShipmentRequest>> = {
    standard: {
      recipient_name: 'Marko Marković',
      recipient_phone: '+381641234567',
      recipient_email: 'marko@example.com',
      recipient_city: 'Beograd',
      recipient_address: 'Takovska 2',
      recipient_zip: '11000',
      weight: 500,
      content: 'Test paket - standardna dostava',
      cod_amount: 0,
      insured_value: 0,
      services: 'PNA',
      delivery_method: 'K',
      payment_method: 'POF',
    },
    cod: {
      recipient_name: 'Ana Anić',
      recipient_phone: '+381641234568',
      recipient_email: 'ana@example.com',
      recipient_city: 'Beograd',
      recipient_address: 'Kneza Miloša 10',
      recipient_zip: '11000',
      weight: 750,
      content: 'Test paket - sa otkupninom',
      cod_amount: 5000,
      insured_value: 5000,
      services: 'PNA,OTK,VD',
      delivery_method: 'K',
      payment_method: 'N',
    },
    express: {
      recipient_name: 'Petar Petrović',
      recipient_phone: '+381641234569',
      recipient_email: 'petar@example.com',
      recipient_city: 'Novi Sad',
      recipient_address: 'Bulevar oslobođenja 50',
      recipient_zip: '21000',
      weight: 300,
      content: 'Test paket - ekspresna dostava',
      cod_amount: 0,
      insured_value: 1000,
      services: 'PNA,SMS',
      delivery_method: 'K',
      payment_method: 'POF',
      id_rukovanje: 30,
    },
    parcel_locker: {
      recipient_name: 'Jovana Jovanović',
      recipient_phone: '+381647654321',
      recipient_email: 'jovana@example.com',
      recipient_city: 'Beograd',
      recipient_address: 'Trg Republike 5',
      recipient_zip: '11000',
      weight: 500,
      content: 'Test paket - paketomat dostava',
      cod_amount: 0,
      insured_value: 0,
      services: 'PNA',
      delivery_method: 'PAK',
      payment_method: 'POF',
      id_rukovanje: 85,
      parcel_locker_code: 'BG001',
    },
  };

  const loadConfig = useCallback(async () => {
    try {
      const response = await apiClient.get('/postexpress/test/config');
      if (response.data.success && response.data.data) {
        const cfg = response.data.data;
        setConfig(cfg);

        // Заполняем форму дефолтными значениями
        setFormData((prev) => ({
          ...prev,
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
        }));
      }
    } catch (err: any) {
      console.error('Failed to load config:', err);
    }
  }, []);

  useEffect(() => {
    loadConfig();
  }, [loadConfig]);

  const loadScenario = (scenarioType: string) => {
    const scenario = testScenarios[scenarioType];
    if (!scenario) return;

    // Merge scenario data with current form data (preserving sender info from config)
    setFormData((prev) => ({
      ...prev,
      ...scenario,
    }));
    setSelectedScenario(scenarioType);
    setResult(null);
    setError(null);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResult(null);

    try {
      const response = await apiClient.post(
        '/postexpress/test/shipment',
        formData
      );

      if (response.data.success && response.data.data) {
        const shipmentResult = response.data.data;
        setResult(shipmentResult);
      } else {
        setError(response.data.message || 'Failed to create test shipment');
      }
    } catch (err: any) {
      setError(err.response?.data?.message || err.message || 'Unknown error');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement
    >
  ) => {
    const { name, value, type } = e.target;
    setFormData({
      ...formData,
      [name]: type === 'number' ? Number(value) : value,
    });
  };

  // ====================== NEW: TX HANDLERS ======================

  // TX 3 - GetSettlements
  const handleTx3Test = async () => {
    if (!tx3Query.trim()) {
      setTx3Error('Please enter a search query');
      return;
    }

    setTx3Loading(true);
    setTx3Error(null);
    setTx3Settlements([]);

    try {
      const response = await apiClient.post(
        '/postexpress/test/tx3-settlements',
        { query: tx3Query }
      );

      if (response.data.success && response.data.data) {
        const settlements =
          response.data.data.naselja || response.data.data.Naselja || [];
        setTx3Settlements(settlements);
      } else {
        setTx3Error(response.data.message || 'Failed to get settlements');
      }
    } catch (err: any) {
      setTx3Error(err.response?.data?.message || err.message || 'TX 3 failed');
    } finally {
      setTx3Loading(false);
    }
  };

  // TX 4 - GetStreets
  const handleTx4Test = async () => {
    if (!tx4SettlementId.trim()) {
      setTx4Error('Please enter a settlement ID');
      return;
    }

    if (!tx4Query.trim()) {
      setTx4Error('Please enter a search query');
      return;
    }

    setTx4Loading(true);
    setTx4Error(null);
    setTx4Streets([]);

    try {
      const response = await apiClient.post('/postexpress/test/tx4-streets', {
        settlement_id: parseInt(tx4SettlementId),
        query: tx4Query,
      });

      if (response.data.success && response.data.data) {
        const streets =
          response.data.data.ulice || response.data.data.Ulice || [];
        setTx4Streets(streets);
      } else {
        setTx4Error(response.data.message || 'Failed to get streets');
      }
    } catch (err: any) {
      setTx4Error(err.response?.data?.message || err.message || 'TX 4 failed');
    } finally {
      setTx4Loading(false);
    }
  };

  // TX 6 - ValidateAddress
  const handleTx6Test = async () => {
    if (tx6Request.id_naselje <= 0) {
      setTx6Error('Please enter a valid settlement ID');
      return;
    }

    if (!tx6Request.broj || !tx6Request.postanski_broj) {
      setTx6Error('Please fill all required fields');
      return;
    }

    setTx6Loading(true);
    setTx6Error(null);
    setTx6Response(null);

    try {
      // Преобразуем snake_case в PascalCase для WSP API
      const wspRequest = {
        IdNaselje: tx6Request.id_naselje,
        IdUlica: tx6Request.id_ulica,
        Broj: tx6Request.broj,
        PostanskiBroj: tx6Request.postanski_broj,
      };

      const response = await apiClient.post(
        '/postexpress/test/tx6-validate-address',
        wspRequest
      );

      // Обрабатываем ошибки из ApiResponse
      if (response.error) {
        setTx6Error(response.error.message || 'Failed to validate address');
        return;
      }

      // Обрабатываем успешный ответ
      if (response?.data?.success && response?.data?.data) {
        setTx6Response(response.data.data);
      } else {
        setTx6Error('Failed to validate address');
      }
    } catch (err: any) {
      const errorMsg = err.message || 'TX 6 failed';
      setTx6Error(errorMsg);
    } finally {
      setTx6Loading(false);
    }
  };

  // TX 9 - CheckServiceAvailability
  const handleTx9Test = async () => {
    if (tx9Request.id_rukovanje <= 0) {
      setTx9Error('Please enter a valid service ID (id_rukovanje)');
      return;
    }

    if (
      !tx9Request.postanski_broj_odlaska ||
      !tx9Request.postanski_broj_dolaska
    ) {
      setTx9Error('Please fill postal codes');
      return;
    }

    setTx9Loading(true);
    setTx9Error(null);
    setTx9Response(null);

    try {
      // Преобразуем snake_case в PascalCase для WSP API
      const wspRequest = {
        IdRukovanje: tx9Request.id_rukovanje,
        PostanskiBrojOdlaska: tx9Request.postanski_broj_odlaska,
        PostanskiBrojDolaska: tx9Request.postanski_broj_dolaska,
      };

      const response = await apiClient.post(
        '/postexpress/test/tx9-service-availability',
        wspRequest
      );

      // Обрабатываем ошибки из ApiResponse
      if (response.error) {
        setTx9Error(
          response.error.message || 'Failed to check service availability'
        );
        return;
      }

      // Обрабатываем успешный ответ
      if (response?.data?.success && response?.data?.data) {
        setTx9Response(response.data.data);
      } else {
        setTx9Error('Failed to check service availability');
      }
    } catch (err: any) {
      const errorMsg = err.message || 'TX 9 failed';
      setTx9Error(errorMsg);
    } finally {
      setTx9Loading(false);
    }
  };

  // TX 11 - CalculatePostage
  const handleTx11Test = async () => {
    if (tx11Request.id_rukovanje <= 0) {
      setTx11Error('Please enter a valid service ID (id_rukovanje)');
      return;
    }

    if (
      !tx11Request.postanski_broj_odlaska ||
      !tx11Request.postanski_broj_dolaska
    ) {
      setTx11Error('Please fill postal codes');
      return;
    }

    if (tx11Request.masa <= 0) {
      setTx11Error('Please enter valid weight (masa)');
      return;
    }

    setTx11Loading(true);
    setTx11Error(null);
    setTx11Response(null);

    try {
      // Преобразуем snake_case в PascalCase для WSP API
      const wspRequest = {
        IdRukovanje: tx11Request.id_rukovanje,
        PostanskiBrojOdlaska: tx11Request.postanski_broj_odlaska,
        PostanskiBrojDolaska: tx11Request.postanski_broj_dolaska,
        Masa: tx11Request.masa,
        Otkupnina: tx11Request.otkupnina,
        Vrednost: tx11Request.vrednost,
        PosebneUsluge: tx11Request.posebne_usluge,
      };

      const response = await apiClient.post(
        '/postexpress/test/tx11-calculate-postage',
        wspRequest
      );

      // Обрабатываем ошибки из ApiResponse
      if (response.error) {
        setTx11Error(response.error.message || 'Failed to calculate postage');
        return;
      }

      // Обрабатываем успешный ответ
      if (response?.data?.success && response?.data?.data) {
        setTx11Response(response.data.data);
      } else {
        setTx11Error('Failed to calculate postage');
      }
    } catch (err: any) {
      const errorMsg = err.message || 'TX 11 failed';
      setTx11Error(errorMsg);
    } finally {
      setTx11Loading(false);
    }
  };

  // TX 73 - B2B Manifest (CreateShipmentViaManifest)
  const handleTx73Test = async () => {
    if (!tx73Request.recipient_name || !tx73Request.recipient_phone) {
      setTx73Error('Please fill recipient name and phone');
      return;
    }

    if (!tx73Request.recipient_city || !tx73Request.recipient_zip) {
      setTx73Error('Please fill recipient city and postal code');
      return;
    }

    if (tx73Request.weight <= 0) {
      setTx73Error('Please enter valid weight');
      return;
    }

    setTx73Loading(true);
    setTx73Error(null);
    setTx73Response(null);

    try {
      // Merge with sender data from config
      const shipmentData = {
        ...tx73Request,
        sender_name: config?.default_sender.name || 'SVETU d.o.o.',
        sender_phone: config?.default_sender.phone || '+381641234567',
        sender_email: config?.default_sender.email || 'b2b@svetu.rs',
        sender_city: config?.default_sender.city || 'Beograd',
        sender_address: config?.default_sender.address || '',
        sender_zip: config?.default_sender.zip || '11000',
        delivery_method: 'K',
        payment_method: 'POF',
        id_rukovanje: 71, // PE_Danas_za_sutra_isporuka
      };

      const response = await apiClient.post(
        '/postexpress/test/shipment',
        shipmentData
      );

      // Обрабатываем ошибки из ApiResponse
      if (response.error) {
        setTx73Error(response.error.message || 'Failed to create shipment');
        return;
      }

      // Обрабатываем успешный ответ
      if (response?.data?.success && response?.data?.data) {
        setTx73Response(response.data.data);
      } else {
        setTx73Error(response.data.message || 'Failed to create shipment');
      }
    } catch (err: any) {
      const errorMsg = err.message || 'TX 73 failed';
      setTx73Error(errorMsg);
    } finally {
      setTx73Loading(false);
    }
  };

  // ====================== MODAL COMPONENTS ======================
  const Modal = ({
    isOpen,
    onClose,
    title,
    children,
  }: {
    isOpen: boolean;
    onClose: () => void;
    title: string;
    children: React.ReactNode;
  }) => {
    if (!isOpen) return null;

    return (
      <div className="fixed inset-0 z-50 overflow-y-auto">
        <div
          className="fixed inset-0 bg-black bg-opacity-50 transition-opacity"
          onClick={onClose}
        ></div>
        <div className="flex min-h-screen items-center justify-center p-4">
          <div className="relative bg-white rounded-lg shadow-xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
            <div className="sticky top-0 bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between">
              <h3 className="text-xl font-bold text-gray-900">{title}</h3>
              <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600 text-2xl leading-none"
              >
                ×
              </button>
            </div>
            <div className="p-6">{children}</div>
          </div>
        </div>
      </div>
    );
  };

  return (
    <div className="min-h-screen bg-gray-50 py-8 px-4">
      <div className="max-w-7xl mx-auto">
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900">📦 {t('title')}</h1>
          <p className="mt-2 text-gray-600">{t('description')}</p>
        </div>

        {/* WSP API Transaction Tests (NEW SECTION) */}
        <div className="bg-gradient-to-r from-cyan-50 to-green-50 rounded-lg shadow-lg p-6 mb-8">
          <h2 className="text-xl font-bold mb-2 text-gray-800">
            🔧 WSP API Transaction Tests
          </h2>
          <p className="text-gray-600 mb-4">
            ✅ Working: TX3 (settlements), TX4 (streets), TX73 (shipments + COD)
            | ⚠️ Known issues: TX6, TX9, TX11 (Post Express API limitations)
          </p>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {/* TX 3 - GetNaselje ✅ WORKING! */}
            <button
              onClick={() => setTx3ModalOpen(true)}
              className="bg-white p-4 rounded-lg border-2 border-cyan-300 hover:border-cyan-500 transition-all text-left shadow-md hover:shadow-xl ring-2 ring-cyan-200"
            >
              <h3 className="font-semibold text-cyan-700 mb-1">
                📍 TX 3: GetNaselje ✅
              </h3>
              <p className="text-sm text-gray-600">
                Search settlements by name (WORKING!)
              </p>
              <div className="mt-2 text-xs text-gray-500">
                Query: city name → List of settlements (122ms)
              </div>
            </button>

            {/* TX 4 - GetUlica ✅ WORKING! */}
            <button
              onClick={() => setTx4ModalOpen(true)}
              className="bg-white p-4 rounded-lg border-2 border-blue-300 hover:border-blue-500 transition-all text-left shadow-md hover:shadow-xl ring-2 ring-blue-200"
            >
              <h3 className="font-semibold text-blue-700 mb-1">
                🛣️ TX 4: GetUlica ✅
              </h3>
              <p className="text-sm text-gray-600">
                Search streets in settlement (WORKING!)
              </p>
              <div className="mt-2 text-xs text-gray-500">
                SettlementID + Query → List of streets (74ms)
              </div>
            </button>

            {/* TX 6 - ProveraAdrese ⚠️ Known Issue */}
            <button
              onClick={() => setTx6ModalOpen(true)}
              className="bg-white p-4 rounded-lg border-2 border-yellow-200 hover:border-yellow-400 transition-all text-left shadow-md hover:shadow-lg"
            >
              <h3 className="font-semibold text-yellow-700 mb-1">
                ⚠️ TX 6: ProveraAdrese
              </h3>
              <p className="text-sm text-gray-600">
                Validate address (PE API format issue)
              </p>
              <div className="mt-2 text-xs text-gray-500">
                ⚠️ Post Express requires specific format
              </div>
            </button>

            {/* TX 9 - ProveraDostupnostiUsluge ⚠️ Known Issue */}
            <button
              onClick={() => setTx9ModalOpen(true)}
              className="bg-white p-4 rounded-lg border-2 border-yellow-200 hover:border-yellow-400 transition-all text-left shadow-md hover:shadow-lg"
            >
              <h3 className="font-semibold text-yellow-700 mb-1">
                ⚠️ TX 9: Service Availability
              </h3>
              <p className="text-sm text-gray-600">
                Check service availability (PE API requirements)
              </p>
              <div className="mt-2 text-xs text-gray-500">
                ⚠️ Requires more address data than postal codes
              </div>
            </button>

            {/* TX 11 - PostarinaPosiljke */}
            <button
              onClick={() => setTx11ModalOpen(true)}
              className="bg-white p-4 rounded-lg border-2 border-amber-200 hover:border-amber-400 transition-all text-left shadow-md hover:shadow-lg"
            >
              <h3 className="font-semibold text-amber-700 mb-1">
                💰 TX 11: Postage Calculation ❌
              </h3>
              <p className="text-sm text-gray-600">
                Calculate shipping cost (BROKEN)
              </p>
              <div className="mt-2 text-xs text-gray-500">
                ⚠️ Post Express database bug
              </div>
            </button>

            {/* TX 73 - B2B Manifest ✅ WORKING! */}
            <button
              onClick={() => setTx73ModalOpen(true)}
              className="bg-white p-4 rounded-lg border-2 border-emerald-300 hover:border-emerald-500 transition-all text-left shadow-md hover:shadow-xl ring-2 ring-emerald-200"
            >
              <h3 className="font-semibold text-emerald-700 mb-1">
                🚀 TX 73: B2B Manifest ✅
              </h3>
              <p className="text-sm text-gray-600">
                Create shipment + Get cost (WORKING!)
              </p>
              <div className="mt-2 text-xs text-gray-500">
                ⭐ Recommended for price calculation
              </div>
            </button>
          </div>
        </div>

        {/* Quick Test Scenarios */}
        <div className="bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg shadow-lg p-6 mb-8">
          <h2 className="text-xl font-bold mb-2 text-gray-800">
            {t('quickTests')}
          </h2>
          <p className="text-gray-600 mb-4">{t('quickTestsDesc')}</p>

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            {/* Standard Test */}
            <button
              onClick={() => loadScenario('standard')}
              className={`bg-white p-4 rounded-lg border-2 transition-all text-left ${
                selectedScenario === 'standard'
                  ? 'border-blue-500 shadow-md ring-2 ring-blue-200'
                  : 'border-blue-200 hover:border-blue-400'
              }`}
            >
              <h3 className="font-semibold text-blue-700 mb-1">
                📦 {t('standardTest')}
              </h3>
              <p className="text-sm text-gray-600">{t('standardTestDesc')}</p>
              <div className="mt-2 text-xs text-gray-500">
                • 500g • Beograd → Beograd
              </div>
            </button>

            {/* COD Test */}
            <button
              onClick={() => loadScenario('cod')}
              className={`bg-white p-4 rounded-lg border-2 transition-all text-left ${
                selectedScenario === 'cod'
                  ? 'border-green-500 shadow-md ring-2 ring-green-200'
                  : 'border-green-200 hover:border-green-400'
              }`}
            >
              <h3 className="font-semibold text-green-700 mb-1">
                💰 {t('codTest')}
              </h3>
              <p className="text-sm text-gray-600">{t('codTestDesc')}</p>
              <div className="mt-2 text-xs text-gray-500">
                • 750g • 5000 RSD COD • Insurance
              </div>
            </button>

            {/* Express Test */}
            <button
              onClick={() => loadScenario('express')}
              className={`bg-white p-4 rounded-lg border-2 transition-all text-left ${
                selectedScenario === 'express'
                  ? 'border-purple-500 shadow-md ring-2 ring-purple-200'
                  : 'border-purple-200 hover:border-purple-400'
              }`}
            >
              <h3 className="font-semibold text-purple-700 mb-1">
                ⚡ {t('expressTest')}
              </h3>
              <p className="text-sm text-gray-600">{t('expressTestDesc')}</p>
              <div className="mt-2 text-xs text-gray-500">
                • 300g • Beograd → Novi Sad • SMS
              </div>
            </button>

            {/* Parcel Locker Test */}
            <button
              onClick={() => loadScenario('parcel_locker')}
              className={`bg-white p-4 rounded-lg border-2 transition-all text-left ${
                selectedScenario === 'parcel_locker'
                  ? 'border-orange-500 shadow-md ring-2 ring-orange-200'
                  : 'border-orange-200 hover:border-orange-400'
              }`}
            >
              <h3 className="font-semibold text-orange-700 mb-1">
                📬 Paketomat
              </h3>
              <p className="text-sm text-gray-600">Isporuka na paketomatu</p>
              <div className="mt-2 text-xs text-gray-500">
                • 500g • IdRukovanje: 85 • BG001
              </div>
            </button>
          </div>
        </div>

        {/* Working Post Express WSP API Features */}
        <div className="bg-gradient-to-r from-green-50 to-blue-50 rounded-lg shadow-lg p-6 mb-8">
          <h2 className="text-xl font-bold mb-2 text-gray-800">
            ✅ Working Post Express WSP API Feature
          </h2>
          <p className="text-gray-600 mb-4">
            Successfully tested with real Post Express API - only TX 73 (B2B
            Manifest) is working!
          </p>
          <p className="text-sm text-amber-700 bg-amber-50 border border-amber-200 rounded p-3 mb-4">
            ⚠️ <strong>Important:</strong> TX 63 (Tracking) returns error:
            &quot;Kretanja još uvek nisu implementirana za izabranu
            uslugu!&quot; - Tracking is not yet implemented for B2B service by
            Post Express.
          </p>

          <div className="max-w-md mx-auto">
            {/* Manifest Creation Test (Transaction 73) - ONLY WORKING ONE */}
            <div className="bg-white p-6 rounded-lg border-2 border-green-300 hover:border-green-500 transition-all shadow-md">
              <h3 className="font-semibold text-green-700 mb-2 text-lg">
                📦 B2B Manifest (TX 73) ✅
              </h3>
              <p className="text-sm text-gray-600 mb-4">
                CreateShipmentViaManifest - The ONLY working Transaction ID for
                creating shipments via B2B API
              </p>
              <button
                onClick={() =>
                  document
                    .getElementById('shipment-form')
                    ?.scrollIntoView({ behavior: 'smooth' })
                }
                className="w-full bg-green-600 hover:bg-green-700 text-white py-3 rounded-lg text-sm font-medium"
              >
                Go to Shipment Form Below ↓
              </button>
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 lg:grid-cols-2 gap-8">
          {/* Форма */}
          <div id="shipment-form" className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-bold mb-6">{t('shipmentParams')}</h2>

            <form onSubmit={handleSubmit} className="space-y-6">
              {/* Получатель */}
              <div>
                <h3 className="text-lg font-semibold mb-4 text-blue-700">
                  {t('recipient')}
                </h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      {t('fullName')} {t('required')}
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
                        {t('phone')} {t('required')}
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
                        {t('email')}
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
                        {t('city')} {t('required')}
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
                        {t('zip')} {t('required')}
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
                      {t('address')}
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
                <h3 className="text-lg font-semibold mb-4 text-green-700">
                  {t('sender')}
                </h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      {t('companyName')} {t('required')}
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
                        {t('phone')} {t('required')}
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
                        {t('email')}
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
                        {t('city')} {t('required')}
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
                        {t('zip')} {t('required')}
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
                      {t('address')}
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
                <h3 className="text-lg font-semibold mb-4 text-purple-700">
                  {t('packageParams')}
                </h3>
                <div className="space-y-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-1">
                      {t('content')} {t('required')}
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
                        {t('weight')} {t('required')}
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
                        {t('codAmount')}
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
                        {t('insuredValue')}
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
                        {t('services')}
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
                        {t('deliveryMethod')}
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
                        {t('paymentMethod')}
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

                  <div className="grid grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        IdRukovanje
                      </label>
                      <select
                        name="id_rukovanje"
                        value={formData.id_rukovanje || 29}
                        onChange={handleChange}
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                      >
                        <option value={29}>29 - PE_Danas_za_sutra_12</option>
                        <option value={30}>30 - PE_Danas_za_danas</option>
                        <option value={55}>55 - PE_Danas_za_odmah</option>
                        <option value={58}>58 - PE_Danas_za_sutra_19</option>
                        <option value={59}>59 - PE_Danas_za_odmah_Bg</option>
                        <option value={71}>
                          71 - PE_Danas_za_sutra_isporuka
                        </option>
                        <option value={85}>85 - Paketomat</option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-1">
                        Parcel Locker Code
                      </label>
                      <input
                        type="text"
                        name="parcel_locker_code"
                        value={formData.parcel_locker_code || ''}
                        onChange={handleChange}
                        placeholder="BG001 (for IdRukovanje: 85)"
                        className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500 focus:border-transparent"
                      />
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
                    <svg
                      className="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                      xmlns="http://www.w3.org/2000/svg"
                      fill="none"
                      viewBox="0 0 24 24"
                    >
                      <circle
                        className="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        strokeWidth="4"
                      ></circle>
                      <path
                        className="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      ></path>
                    </svg>
                    {t('creating')}
                  </>
                ) : (
                  <>📦 {t('createButton')}</>
                )}
              </button>
            </form>
          </div>

          {/* Результаты */}
          <div className="bg-white rounded-lg shadow-lg p-6">
            <h2 className="text-xl font-bold mb-6">{t('results')}</h2>

            {!result && !error && (
              <div className="text-center py-12 text-gray-400">
                <svg
                  className="mx-auto h-24 w-24 mb-4"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={1}
                    d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4"
                  />
                </svg>
                <p className="text-lg">{t('fillFormPrompt')}</p>
              </div>
            )}

            {error && (
              <div className="bg-red-50 border-l-4 border-red-500 p-4 mb-4">
                <div className="flex items-start">
                  <svg
                    className="h-6 w-6 text-red-500 mr-3 flex-shrink-0"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <div>
                    <h3 className="text-red-800 font-semibold">{t('error')}</h3>
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
                    <svg
                      className="h-6 w-6 text-green-500 mr-3 flex-shrink-0"
                      fill="none"
                      viewBox="0 0 24 24"
                      stroke="currentColor"
                    >
                      <path
                        strokeLinecap="round"
                        strokeLinejoin="round"
                        strokeWidth={2}
                        d="M5 13l4 4L19 7"
                      />
                    </svg>
                    <div>
                      <h3 className="text-green-800 font-semibold">
                        {t('successTitle')}
                      </h3>
                      <p className="text-green-700 text-sm mt-1">
                        {t('processingTime')}: {result.processing_time_ms}ms
                      </p>
                    </div>
                  </div>
                </div>

                {/* Tracking info */}
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-blue-50 p-4 rounded-lg">
                    <p className="text-sm text-blue-600 font-medium mb-1">
                      {t('trackingNumber')}
                    </p>
                    <p className="text-2xl font-bold text-blue-900">
                      {result.tracking_number}
                    </p>
                  </div>

                  <div className="bg-purple-50 p-4 rounded-lg">
                    <p className="text-sm text-purple-600 font-medium mb-1">
                      {t('cost')}
                    </p>
                    <p className="text-2xl font-bold text-purple-900">
                      {result.cost} RSD
                    </p>
                  </div>
                </div>

                {/* IDs */}
                <div className="grid grid-cols-3 gap-4">
                  <div className="bg-gray-50 p-3 rounded">
                    <p className="text-xs text-gray-600 mb-1">
                      {t('manifestId')}
                    </p>
                    <p className="font-mono text-sm font-semibold">
                      {result.manifest_id}
                    </p>
                  </div>

                  <div className="bg-gray-50 p-3 rounded">
                    <p className="text-xs text-gray-600 mb-1">
                      {t('shipmentId')}
                    </p>
                    <p className="font-mono text-sm font-semibold">
                      {result.shipment_id}
                    </p>
                  </div>

                  <div className="bg-gray-50 p-3 rounded">
                    <p className="text-xs text-gray-600 mb-1">
                      {t('externalId')}
                    </p>
                    <p className="font-mono text-xs font-semibold">
                      {result.external_id}
                    </p>
                  </div>
                </div>

                {/* Request data */}
                <div>
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">
                    {t('requestData')}
                  </h3>
                  <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg text-xs overflow-auto max-h-64">
                    {JSON.stringify(result.request_data, null, 2)}
                  </pre>
                </div>

                {/* Response data */}
                <div>
                  <h3 className="text-sm font-semibold text-gray-700 mb-2">
                    {t('responseData')}
                  </h3>
                  <pre className="bg-gray-900 text-gray-100 p-4 rounded-lg text-xs overflow-auto max-h-64">
                    {JSON.stringify(result.response_data, null, 2)}
                  </pre>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* ====================== TX MODALS ====================== */}

        {/* ====================== TX 3 MODAL ====================== */}
        <Modal
          isOpen={tx3ModalOpen}
          onClose={() => setTx3ModalOpen(false)}
          title="TX 3: GetNaselje (Search Settlements)"
        >
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Search Query
              </label>
              <input
                type="text"
                value={tx3Query}
                onChange={(e) => setTx3Query(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleTx3Test()}
                placeholder="Beograd, Niš, Novi Sad"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <button
              onClick={handleTx3Test}
              disabled={tx3Loading}
              className="w-full bg-cyan-600 hover:bg-cyan-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg"
            >
              {tx3Loading ? 'Testing...' : '🔍 Test TX 3'}
            </button>

            {tx3Error && (
              <div className="bg-red-50 border-l-4 border-red-500 p-4 text-red-700">
                {tx3Error}
              </div>
            )}

            {tx3Settlements.length > 0 && (
              <div>
                <h4 className="font-semibold mb-2 text-gray-800">
                  Settlements ({tx3Settlements.length})
                </h4>
                <div className="space-y-2 max-h-96 overflow-y-auto">
                  {tx3Settlements.map((s, i) => (
                    <div
                      key={i}
                      className="border border-gray-200 rounded p-3 hover:bg-gray-50"
                    >
                      <div className="grid grid-cols-2 gap-2 text-sm">
                        <div>
                          <span className="font-medium">ID:</span>{' '}
                          {s.IdNaselje || s.id_naselje}
                        </div>
                        <div>
                          <span className="font-medium">Name:</span>{' '}
                          {s.Naziv || s.naziv}
                        </div>
                        <div>
                          <span className="font-medium">Postal Code:</span>{' '}
                          {s.PostanskiBroj || s.postanski_broj}
                        </div>
                        <div>
                          <span className="font-medium">Region:</span>{' '}
                          {s.NazivOkruga || s.naziv_okruga}
                        </div>
                      </div>
                      <button
                        onClick={() => {
                          setTx4SettlementId(
                            String(s.IdNaselje || s.id_naselje)
                          );
                          setTx3ModalOpen(false);
                          setTx4ModalOpen(true);
                        }}
                        className="mt-2 text-xs bg-blue-100 hover:bg-blue-200 text-blue-700 px-3 py-1 rounded"
                      >
                        Use in TX 4 →
                      </button>
                      <button
                        onClick={() => {
                          setTx6Request((prev) => ({
                            ...prev,
                            id_naselje: s.IdNaselje || s.id_naselje,
                            postanski_broj: s.PostanskiBroj || s.postanski_broj,
                          }));
                          setTx3ModalOpen(false);
                          setTx6ModalOpen(true);
                        }}
                        className="mt-2 ml-2 text-xs bg-green-100 hover:bg-green-200 text-green-700 px-3 py-1 rounded"
                      >
                        Use in TX 6 →
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {tx3Settlements.length > 0 && (
              <details className="mt-4">
                <summary className="cursor-pointer text-sm font-medium text-gray-700">
                  Raw Response JSON
                </summary>
                <pre className="mt-2 bg-gray-900 text-gray-100 p-4 rounded text-xs overflow-auto max-h-64">
                  {JSON.stringify(tx3Settlements, null, 2)}
                </pre>
              </details>
            )}
          </div>
        </Modal>

        {/* ====================== TX 4 MODAL ====================== */}
        <Modal
          isOpen={tx4ModalOpen}
          onClose={() => setTx4ModalOpen(false)}
          title="TX 4: GetUlica (Search Streets)"
        >
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Settlement ID (from TX 3)
              </label>
              <input
                type="number"
                value={tx4SettlementId}
                onChange={(e) => setTx4SettlementId(e.target.value)}
                placeholder="100001 (Beograd)"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Street Search Query
              </label>
              <input
                type="text"
                value={tx4Query}
                onChange={(e) => setTx4Query(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleTx4Test()}
                placeholder="Takovska, Kneza Miloša"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500"
              />
            </div>

            <button
              onClick={handleTx4Test}
              disabled={tx4Loading}
              className="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg"
            >
              {tx4Loading ? 'Testing...' : '🔍 Test TX 4'}
            </button>

            {tx4Error && (
              <div className="bg-red-50 border-l-4 border-red-500 p-4 text-red-700">
                {tx4Error}
              </div>
            )}

            {tx4Streets.length > 0 && (
              <div>
                <h4 className="font-semibold mb-2 text-gray-800">
                  Streets ({tx4Streets.length})
                </h4>
                <div className="space-y-2 max-h-96 overflow-y-auto">
                  {tx4Streets.map((street, i) => (
                    <div
                      key={i}
                      className="border border-gray-200 rounded p-3 hover:bg-gray-50"
                    >
                      <div className="grid grid-cols-2 gap-2 text-sm">
                        <div>
                          <span className="font-medium">ID:</span>{' '}
                          {street.IdUlica || street.id_ulica}
                        </div>
                        <div>
                          <span className="font-medium">Name:</span>{' '}
                          {street.Naziv || street.naziv}
                        </div>
                      </div>
                      <button
                        onClick={() => {
                          setTx6Request((prev) => ({
                            ...prev,
                            id_ulica: street.IdUlica || street.id_ulica,
                          }));
                          setTx4ModalOpen(false);
                          setTx6ModalOpen(true);
                        }}
                        className="mt-2 text-xs bg-green-100 hover:bg-green-200 text-green-700 px-3 py-1 rounded"
                      >
                        Use in TX 6 →
                      </button>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {tx4Streets.length > 0 && (
              <details className="mt-4">
                <summary className="cursor-pointer text-sm font-medium text-gray-700">
                  Raw Response JSON
                </summary>
                <pre className="mt-2 bg-gray-900 text-gray-100 p-4 rounded text-xs overflow-auto max-h-64">
                  {JSON.stringify(tx4Streets, null, 2)}
                </pre>
              </details>
            )}
          </div>
        </Modal>

        {/* ====================== TX 6 MODAL ====================== */}
        <Modal
          isOpen={tx6ModalOpen}
          onClose={() => setTx6ModalOpen(false)}
          title="TX 6: ProveraAdrese (Validate Address)"
        >
          <div className="space-y-4">
            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Settlement ID (from TX 3)
                </label>
                <input
                  type="number"
                  value={tx6Request.id_naselje}
                  onChange={(e) =>
                    setTx6Request({
                      ...tx6Request,
                      id_naselje: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="100001 (Beograd)"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Street ID (from TX 4, optional)
                </label>
                <input
                  type="number"
                  value={tx6Request.id_ulica}
                  onChange={(e) =>
                    setTx6Request({
                      ...tx6Request,
                      id_ulica: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="0"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                />
              </div>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  House Number *
                </label>
                <input
                  type="text"
                  value={tx6Request.broj}
                  onChange={(e) =>
                    setTx6Request({ ...tx6Request, broj: e.target.value })
                  }
                  placeholder="2"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Postal Code *
                </label>
                <input
                  type="text"
                  value={tx6Request.postanski_broj}
                  onChange={(e) =>
                    setTx6Request({
                      ...tx6Request,
                      postanski_broj: e.target.value,
                    })
                  }
                  placeholder="11000"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-green-500"
                />
              </div>
            </div>

            <button
              onClick={handleTx6Test}
              disabled={tx6Loading}
              className="w-full bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg"
            >
              {tx6Loading ? 'Testing...' : '✅ Test TX 6'}
            </button>

            {tx6Error && (
              <div className="bg-red-50 border-l-4 border-red-500 p-4 text-red-700">
                {tx6Error}
              </div>
            )}

            {tx6Response && (
              <div>
                <h4 className="font-semibold mb-2 text-gray-800">
                  Validation Result
                </h4>
                <div className="border border-gray-200 rounded p-4 bg-gray-50">
                  <div className="grid grid-cols-2 gap-3 text-sm">
                    <div>
                      <span className="font-medium">Address Exists:</span>{' '}
                      <span
                        className={
                          tx6Response.postoji_adresa ||
                          tx6Response.PostojiAdresa
                            ? 'text-green-600 font-bold'
                            : 'text-red-600 font-bold'
                        }
                      >
                        {tx6Response.postoji_adresa || tx6Response.PostojiAdresa
                          ? '✅ YES'
                          : '❌ NO'}
                      </span>
                    </div>
                    <div>
                      <span className="font-medium">Result Code:</span>{' '}
                      {tx6Response.rezultat || tx6Response.Rezultat}
                    </div>
                    {(tx6Response.poruka || tx6Response.Poruka) && (
                      <div className="col-span-2">
                        <span className="font-medium">Message:</span>{' '}
                        {tx6Response.poruka || tx6Response.Poruka}
                      </div>
                    )}
                  </div>
                </div>

                <details className="mt-4">
                  <summary className="cursor-pointer text-sm font-medium text-gray-700">
                    Raw Response JSON
                  </summary>
                  <pre className="mt-2 bg-gray-900 text-gray-100 p-4 rounded text-xs overflow-auto max-h-64">
                    {JSON.stringify(tx6Response, null, 2)}
                  </pre>
                </details>
              </div>
            )}
          </div>
        </Modal>

        {/* ====================== TX 9 MODAL ====================== */}
        <Modal
          isOpen={tx9ModalOpen}
          onClose={() => setTx9ModalOpen(false)}
          title="TX 9: ProveraDostupnostiUsluge (Check Service Availability)"
        >
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Service ID (IdRukovanje) *
              </label>
              <select
                value={tx9Request.id_rukovanje}
                onChange={(e) =>
                  setTx9Request({
                    ...tx9Request,
                    id_rukovanje: parseInt(e.target.value),
                  })
                }
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500"
              >
                <option value={29}>29 - PE_Danas_za_sutra_12</option>
                <option value={30}>30 - PE_Danas_za_danas</option>
                <option value={55}>55 - PE_Danas_za_odmah</option>
                <option value={58}>58 - PE_Danas_za_sutra_19</option>
                <option value={59}>59 - PE_Danas_za_odmah_Bg</option>
                <option value={71}>71 - PE_Danas_za_sutra_isporuka</option>
                <option value={85}>85 - Paketomat</option>
              </select>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  From Postal Code *
                </label>
                <input
                  type="text"
                  value={tx9Request.postanski_broj_odlaska}
                  onChange={(e) =>
                    setTx9Request({
                      ...tx9Request,
                      postanski_broj_odlaska: e.target.value,
                    })
                  }
                  placeholder="11000"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  To Postal Code *
                </label>
                <input
                  type="text"
                  value={tx9Request.postanski_broj_dolaska}
                  onChange={(e) =>
                    setTx9Request({
                      ...tx9Request,
                      postanski_broj_dolaska: e.target.value,
                    })
                  }
                  placeholder="21000"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-purple-500"
                />
              </div>
            </div>

            <button
              onClick={handleTx9Test}
              disabled={tx9Loading}
              className="w-full bg-purple-600 hover:bg-purple-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg"
            >
              {tx9Loading ? 'Testing...' : '📦 Test TX 9'}
            </button>

            {tx9Error && (
              <div className="bg-red-50 border-l-4 border-red-500 p-4 text-red-700">
                {tx9Error}
              </div>
            )}

            {tx9Response && (
              <div>
                <h4 className="font-semibold mb-2 text-gray-800">
                  Service Availability
                </h4>
                <div className="border border-gray-200 rounded p-4 bg-gray-50">
                  <div className="grid grid-cols-2 gap-3 text-sm">
                    <div>
                      <span className="font-medium">Available:</span>{' '}
                      <span
                        className={
                          tx9Response.dostupna || tx9Response.Dostupna
                            ? 'text-green-600 font-bold'
                            : 'text-red-600 font-bold'
                        }
                      >
                        {tx9Response.dostupna || tx9Response.Dostupna
                          ? '✅ YES'
                          : '❌ NO'}
                      </span>
                    </div>
                    <div>
                      <span className="font-medium">Expected Days:</span>{' '}
                      {tx9Response.ocekivano_dana || tx9Response.OcekivanoDana}
                    </div>
                    <div>
                      <span className="font-medium">Result Code:</span>{' '}
                      {tx9Response.rezultat || tx9Response.Rezultat}
                    </div>
                    {(tx9Response.poruka || tx9Response.Poruka) && (
                      <div className="col-span-2">
                        <span className="font-medium">Message:</span>{' '}
                        {tx9Response.poruka || tx9Response.Poruka}
                      </div>
                    )}
                  </div>
                </div>

                <details className="mt-4">
                  <summary className="cursor-pointer text-sm font-medium text-gray-700">
                    Raw Response JSON
                  </summary>
                  <pre className="mt-2 bg-gray-900 text-gray-100 p-4 rounded text-xs overflow-auto max-h-64">
                    {JSON.stringify(tx9Response, null, 2)}
                  </pre>
                </details>
              </div>
            )}
          </div>
        </Modal>

        {/* ====================== TX 11 MODAL ====================== */}
        <Modal
          isOpen={tx11ModalOpen}
          onClose={() => setTx11ModalOpen(false)}
          title="TX 11: PostarinaPosiljke (Calculate Postage)"
        >
          <div className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Service ID (IdRukovanje) *
              </label>
              <select
                value={tx11Request.id_rukovanje}
                onChange={(e) =>
                  setTx11Request({
                    ...tx11Request,
                    id_rukovanje: parseInt(e.target.value),
                  })
                }
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500"
              >
                <option value={29}>29 - PE_Danas_za_sutra_12</option>
                <option value={30}>30 - PE_Danas_za_danas</option>
                <option value={55}>55 - PE_Danas_za_odmah</option>
                <option value={58}>58 - PE_Danas_za_sutra_19</option>
                <option value={59}>59 - PE_Danas_za_odmah_Bg</option>
                <option value={71}>71 - PE_Danas_za_sutra_isporuka</option>
                <option value={85}>85 - Paketomat</option>
              </select>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  From Postal Code *
                </label>
                <input
                  type="text"
                  value={tx11Request.postanski_broj_odlaska}
                  onChange={(e) =>
                    setTx11Request({
                      ...tx11Request,
                      postanski_broj_odlaska: e.target.value,
                    })
                  }
                  placeholder="11000"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  To Postal Code *
                </label>
                <input
                  type="text"
                  value={tx11Request.postanski_broj_dolaska}
                  onChange={(e) =>
                    setTx11Request({
                      ...tx11Request,
                      postanski_broj_dolaska: e.target.value,
                    })
                  }
                  placeholder="21000"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500"
                />
              </div>
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Weight (grams) *
                </label>
                <input
                  type="number"
                  value={tx11Request.masa}
                  onChange={(e) =>
                    setTx11Request({
                      ...tx11Request,
                      masa: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="500"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  COD Amount (para)
                </label>
                <input
                  type="number"
                  value={tx11Request.otkupnina}
                  onChange={(e) =>
                    setTx11Request({
                      ...tx11Request,
                      otkupnina: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="0"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Insured Value (para)
                </label>
                <input
                  type="number"
                  value={tx11Request.vrednost}
                  onChange={(e) =>
                    setTx11Request({
                      ...tx11Request,
                      vrednost: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="0"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Special Services (comma-separated codes)
              </label>
              <input
                type="text"
                value={tx11Request.posebne_usluge}
                onChange={(e) =>
                  setTx11Request({
                    ...tx11Request,
                    posebne_usluge: e.target.value,
                  })
                }
                placeholder="PNA,SMS"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500"
              />
              <p className="text-xs text-gray-500 mt-1">
                Examples: PNA (notification), SMS, OTK (COD), VD (declared
                value)
              </p>
            </div>

            <button
              onClick={handleTx11Test}
              disabled={tx11Loading}
              className="w-full bg-amber-600 hover:bg-amber-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg"
            >
              {tx11Loading ? 'Testing...' : '💰 Test TX 11'}
            </button>

            {tx11Error && (
              <div className="bg-red-50 border-l-4 border-red-500 p-4 text-red-700">
                {tx11Error}
              </div>
            )}

            {tx11Response && (
              <div>
                <h4 className="font-semibold mb-2 text-gray-800">
                  Postage Calculation
                </h4>
                <div className="border border-gray-200 rounded p-4 bg-gradient-to-br from-amber-50 to-yellow-50">
                  <div className="text-center mb-4">
                    <div className="text-3xl font-bold text-amber-700">
                      {(
                        (tx11Response.postarina || tx11Response.Postarina) / 100
                      ).toFixed(2)}{' '}
                      RSD
                    </div>
                    <div className="text-sm text-gray-600">
                      ({tx11Response.postarina || tx11Response.Postarina} para)
                    </div>
                  </div>

                  <div className="grid grid-cols-2 gap-3 text-sm border-t pt-3">
                    <div>
                      <span className="font-medium">Result Code:</span>{' '}
                      {tx11Response.rezultat || tx11Response.Rezultat}
                    </div>
                    {(tx11Response.poruka || tx11Response.Poruka) && (
                      <div className="col-span-2">
                        <span className="font-medium">Message:</span>{' '}
                        {tx11Response.poruka || tx11Response.Poruka}
                      </div>
                    )}
                  </div>
                </div>

                <details className="mt-4">
                  <summary className="cursor-pointer text-sm font-medium text-gray-700">
                    Raw Response JSON
                  </summary>
                  <pre className="mt-2 bg-gray-900 text-gray-100 p-4 rounded text-xs overflow-auto max-h-64">
                    {JSON.stringify(tx11Response, null, 2)}
                  </pre>
                </details>
              </div>
            )}
          </div>
        </Modal>

        {/* ====================== TX 73 MODAL ====================== */}
        <Modal
          isOpen={tx73ModalOpen}
          onClose={() => setTx73ModalOpen(false)}
          title="TX 73: B2B Manifest (Create Shipment + Get Cost)"
        >
          <div className="space-y-4">
            <div className="bg-emerald-50 border border-emerald-200 rounded p-3 mb-4">
              <p className="text-sm text-emerald-700">
                ⭐ <strong>Recommended for postage calculation!</strong> TX 73
                creates a real shipment and returns the actual cost.
              </p>
            </div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Recipient Name *
                </label>
                <input
                  type="text"
                  value={tx73Request.recipient_name}
                  onChange={(e) =>
                    setTx73Request({
                      ...tx73Request,
                      recipient_name: e.target.value,
                    })
                  }
                  placeholder="Marko Marković"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Recipient Phone *
                </label>
                <input
                  type="text"
                  value={tx73Request.recipient_phone}
                  onChange={(e) =>
                    setTx73Request({
                      ...tx73Request,
                      recipient_phone: e.target.value,
                    })
                  }
                  placeholder="+381641234567"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                />
              </div>
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div className="col-span-2">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Recipient City *
                </label>
                <input
                  type="text"
                  value={tx73Request.recipient_city}
                  onChange={(e) =>
                    setTx73Request({
                      ...tx73Request,
                      recipient_city: e.target.value,
                    })
                  }
                  placeholder="Beograd"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Postal Code *
                </label>
                <input
                  type="text"
                  value={tx73Request.recipient_zip}
                  onChange={(e) =>
                    setTx73Request({
                      ...tx73Request,
                      recipient_zip: e.target.value,
                    })
                  }
                  placeholder="11000"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Recipient Address
              </label>
              <input
                type="text"
                value={tx73Request.recipient_address}
                onChange={(e) =>
                  setTx73Request({
                    ...tx73Request,
                    recipient_address: e.target.value,
                  })
                }
                placeholder="Takovska 2"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
              />
            </div>

            <div className="grid grid-cols-3 gap-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Weight (grams) *
                </label>
                <input
                  type="number"
                  value={tx73Request.weight}
                  onChange={(e) =>
                    setTx73Request({
                      ...tx73Request,
                      weight: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="500"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  COD Amount (para)
                </label>
                <input
                  type="number"
                  value={tx73Request.cod_amount}
                  onChange={(e) =>
                    setTx73Request({
                      ...tx73Request,
                      cod_amount: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="0"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                />
              </div>

              <div>
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Insured Value (para)
                </label>
                <input
                  type="number"
                  value={tx73Request.insured_value}
                  onChange={(e) =>
                    setTx73Request({
                      ...tx73Request,
                      insured_value: parseInt(e.target.value) || 0,
                    })
                  }
                  placeholder="0"
                  className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
                />
              </div>
            </div>

            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Package Content *
              </label>
              <input
                type="text"
                value={tx73Request.content}
                onChange={(e) =>
                  setTx73Request({
                    ...tx73Request,
                    content: e.target.value,
                  })
                }
                placeholder="Test paket - TX 73 test"
                className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-emerald-500"
              />
            </div>

            <button
              onClick={handleTx73Test}
              disabled={tx73Loading}
              className="w-full bg-emerald-600 hover:bg-emerald-700 disabled:bg-gray-400 text-white font-semibold py-2 px-4 rounded-lg"
            >
              {tx73Loading
                ? 'Creating Shipment...'
                : '🚀 Create Shipment & Get Cost'}
            </button>

            {tx73Error && (
              <div className="bg-red-50 border-l-4 border-red-500 p-4 text-red-700">
                {tx73Error}
              </div>
            )}

            {tx73Response && tx73Response.success && (
              <div>
                <h4 className="font-semibold mb-2 text-gray-800">
                  ✅ Shipment Created Successfully!
                </h4>

                {/* Cost Display - MOST IMPORTANT! */}
                <div className="bg-gradient-to-br from-emerald-50 to-green-50 border-2 border-emerald-300 rounded-lg p-4 mb-4">
                  <div className="text-center">
                    <p className="text-sm text-gray-600 mb-1">Shipping Cost</p>
                    <div className="text-4xl font-bold text-emerald-700">
                      {tx73Response.cost} RSD
                    </div>
                    <p className="text-xs text-gray-500 mt-1">
                      ({tx73Response.cost * 100} para)
                    </p>
                  </div>
                </div>

                {/* Tracking & IDs */}
                <div className="grid grid-cols-2 gap-4 mb-4">
                  <div className="bg-blue-50 p-3 rounded">
                    <p className="text-xs text-blue-600 mb-1">
                      Tracking Number
                    </p>
                    <p className="font-mono text-sm font-semibold text-blue-900">
                      {tx73Response.tracking_number}
                    </p>
                  </div>

                  <div className="bg-purple-50 p-3 rounded">
                    <p className="text-xs text-purple-600 mb-1">Manifest ID</p>
                    <p className="font-mono text-sm font-semibold text-purple-900">
                      {tx73Response.manifest_id}
                    </p>
                  </div>
                </div>

                <div className="grid grid-cols-2 gap-4 mb-4">
                  <div className="bg-gray-50 p-3 rounded">
                    <p className="text-xs text-gray-600 mb-1">Shipment ID</p>
                    <p className="font-mono text-sm font-semibold">
                      {tx73Response.shipment_id}
                    </p>
                  </div>

                  <div className="bg-gray-50 p-3 rounded">
                    <p className="text-xs text-gray-600 mb-1">
                      Processing Time
                    </p>
                    <p className="font-mono text-sm font-semibold">
                      {tx73Response.processing_time_ms}ms
                    </p>
                  </div>
                </div>

                <details className="mt-4">
                  <summary className="cursor-pointer text-sm font-medium text-gray-700">
                    Full Response JSON
                  </summary>
                  <pre className="mt-2 bg-gray-900 text-gray-100 p-4 rounded text-xs overflow-auto max-h-64">
                    {JSON.stringify(tx73Response, null, 2)}
                  </pre>
                </details>
              </div>
            )}
          </div>
        </Modal>
      </div>
    </div>
  );
}
