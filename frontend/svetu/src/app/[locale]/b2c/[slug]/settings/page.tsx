'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTranslations, useLocale } from 'next-intl';
import Image from 'next/image';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchStorefrontBySlug,
  updateStorefront,
} from '@/store/slices/b2cStoreSlice';
import { B2CStoreUpdateDTO } from '@/types/b2c';
import {
  BuildingStorefrontIcon,
  MapPinIcon,
  ClockIcon,
  CreditCardIcon,
  TruckIcon,
  PhotoIcon,
  ArrowLeftIcon,
  CheckCircleIcon,
  ExclamationCircleIcon,
  CogIcon,
  ShieldCheckIcon,
  ChartBarIcon,
  BellIcon,
  DocumentTextIcon,
  KeyIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';
import { toast } from '@/utils/toast';
import { storefrontApi } from '@/services/b2cStoreApi';
import LocationPicker from '@/components/GIS/LocationPicker';

interface DeliveryProvider {
  id: string;
  name: string;
  description: string;
  icon: string;
  enabled: boolean;
  settings?: any;
}

interface PaymentMethod {
  id: string;
  name: string;
  description: string;
  icon: string;
  enabled: boolean;
  settings?: {
    account_number?: string;
    api_key?: string;
    webhook_url?: string;
    min_amount?: number;
    max_amount?: number;
    commission?: number;
  };
}

export default function StorefrontSettingsPage() {
  const params = useParams();
  const router = useRouter();
  const dispatch = useAppDispatch();
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const slug = params?.slug as string;

  const { currentStorefront, isLoading, error } = useAppSelector(
    (state) => state.b2cStores
  );

  const [formData, setFormData] = useState<B2CStoreUpdateDTO>({
    name: '',
    description: '',
    phone: '',
    email: '',
    website: '',
    location: {
      country: '',
      city: '',
      full_address: '',
      postal_code: '',
    },
    settings: {},
  });

  const [activeTab, setActiveTab] = useState<string>('general');

  const [businessHours, setBusinessHours] = useState({
    monday: { open: true, from: '09:00', to: '18:00' },
    tuesday: { open: true, from: '09:00', to: '18:00' },
    wednesday: { open: true, from: '09:00', to: '18:00' },
    thursday: { open: true, from: '09:00', to: '18:00' },
    friday: { open: true, from: '09:00', to: '18:00' },
    saturday: { open: true, from: '10:00', to: '16:00' },
    sunday: { open: false, from: '10:00', to: '16:00' },
  });

  const [deliveryProviders, setDeliveryProviders] = useState<
    DeliveryProvider[]
  >([
    {
      id: 'pickup',
      name: 'Самовывоз',
      description: 'Покупатели могут забрать товар самостоятельно',
      icon: '🏪',
      enabled: true,
      settings: {},
    },
    {
      id: 'local_delivery',
      name: 'Локальная доставка',
      description: 'Доставка курьером в пределах города',
      icon: '🚲',
      enabled: true,
      settings: {
        base_rate: 500,
        free_shipping_threshold: 5000,
        estimated_days: 1,
      },
    },
    {
      id: 'post_express',
      name: 'Post Express',
      description:
        'Национальная почта Сербии - доставка по всей стране с полной интеграцией',
      icon: '📮',
      enabled: false,
      settings: {
        api_username: '',
        api_password: '',
        api_endpoint:
          'https://onlinepostexpress.rs/WSPWebApi/api/app/transakcija',
        sender_name: '',
        sender_address: '',
        sender_city: '',
        sender_postal_code: '',
        sender_phone: '',
        sender_email: '',
        test_mode: false,
        auto_print_labels: false,
        auto_track_shipments: true,
        notify_on_pickup: true,
        notify_on_delivery: true,
        notify_on_failed_delivery: true,
        enable_cod: true,
        enable_insurance: true,
        enable_express: false,
        enable_saturday_delivery: false,
        weight_tiers: {
          '0-2kg': 340,
          '2-5kg': 450,
          '5-10kg': 580,
          '10-20kg': 790,
        },
        cod_fee: 45,
        insurance_base: 15000,
        insurance_rate: 0.01,
        free_shipping_threshold: 5000,
        free_warehouse_pickup_threshold: 2000,
        estimated_days: 1,
        tracking_url: 'https://postexpress.rs/track/',
      },
    },
    {
      id: 'dhl',
      name: 'DHL Express',
      description: 'Международная экспресс-доставка',
      icon: '✈️',
      enabled: false,
      settings: {
        api_key: '',
        base_rate: 2500,
        estimated_days: 5,
        tracking_url: 'https://www.dhl.com/track/',
      },
    },
    {
      id: 'aks',
      name: 'AKS',
      description: 'Доставка по Сербии через AKS',
      icon: '🚛',
      enabled: false,
      settings: {
        api_key: '',
        base_rate: 650,
        free_shipping_threshold: 8000,
        estimated_days: 2,
      },
    },
  ]);

  const [paymentMethods, setPaymentMethods] = useState<PaymentMethod[]>([
    {
      id: 'cash',
      name: 'Наличные',
      description: 'Оплата наличными при получении',
      icon: '💵',
      enabled: true,
      settings: {},
    },
    {
      id: 'card',
      name: 'Банковская карта',
      description: 'Оплата картой онлайн или при получении',
      icon: '💳',
      enabled: true,
      settings: {
        commission: 2.5,
      },
    },
    {
      id: 'bank_transfer',
      name: 'Банковский перевод',
      description: 'Перевод на расчетный счет',
      icon: '🏦',
      enabled: false,
      settings: {
        account_number: '',
      },
    },
    {
      id: 'paypal',
      name: 'PayPal',
      description: 'Оплата через PayPal',
      icon: '💙',
      enabled: false,
      settings: {
        api_key: '',
        webhook_url: '',
        commission: 3.4,
      },
    },
    {
      id: 'crypto',
      name: 'Криптовалюта',
      description: 'Оплата в Bitcoin, Ethereum и других криптовалютах',
      icon: '₿',
      enabled: false,
      settings: {
        api_key: '',
        min_amount: 10,
        commission: 1,
      },
    },
  ]);

  const [notificationSettings, setNotificationSettings] = useState({
    email_notifications: true,
    telegram_notifications: false,
    new_orders: true,
    new_messages: true,
    new_reviews: true,
    low_stock: true,
    daily_summary: false,
    weekly_report: true,
  });

  const [seoSettings, setSeoSettings] = useState({
    meta_title: '',
    meta_description: '',
    meta_keywords: '',
    og_image: '',
    google_analytics: '',
    facebook_pixel: '',
    yandex_metrika: '',
  });

  const [securitySettings, setSecuritySettings] = useState({
    two_factor_auth: false,
    ip_whitelist: false,
    allowed_ips: '',
    api_access: false,
    api_key: '',
    webhook_secret: '',
  });

  const [businessInfo, setBusinessInfo] = useState({
    business_type: 'individual',
    registration_number: '',
    tax_number: '',
    vat_number: '',
    legal_name: '',
    legal_address: '',
  });

  const [logoPreview, setLogoPreview] = useState<string | null>(null);
  const [bannerPreview, setBannerPreview] = useState<string | null>(null);
  const [uploadingLogo, setUploadingLogo] = useState(false);
  const [uploadingBanner, setUploadingBanner] = useState(false);

  const [showLocationPicker, setShowLocationPicker] = useState(false);
  const [coordinates, setCoordinates] = useState<{
    lat: number;
    lng: number;
  } | null>(null);

  useEffect(() => {
    if (slug) {
      dispatch(fetchStorefrontBySlug(slug));
    }
  }, [dispatch, slug]);

  useEffect(() => {
    if (currentStorefront) {
      setFormData({
        name: currentStorefront.name || '',
        description: currentStorefront.description || '',
        phone: currentStorefront.phone || '',
        email: currentStorefront.email || '',
        website: currentStorefront.website || '',
        location: (currentStorefront as any).location || {
          country: '',
          city: '',
          full_address: '',
          postal_code: '',
        },
        settings: currentStorefront.settings || {},
      });

      // Загружаем настройки из витрины
      const settings = (currentStorefront.settings as any) || {};

      if (settings.business_hours) {
        setBusinessHours(settings.business_hours);
      }

      if (settings.delivery_providers) {
        setDeliveryProviders(settings.delivery_providers);
      }

      if (settings.payment_methods) {
        setPaymentMethods(settings.payment_methods);
      }

      if (settings.notifications) {
        setNotificationSettings(settings.notifications);
      }

      if (settings.seo) {
        setSeoSettings(settings.seo);
      }

      if (settings.security) {
        setSecuritySettings(settings.security);
      }

      if (settings.business_info) {
        setBusinessInfo(settings.business_info);
      }

      // Загружаем координаты если есть
      if (currentStorefront.latitude && currentStorefront.longitude) {
        setCoordinates({
          lat: currentStorefront.latitude,
          lng: currentStorefront.longitude,
        });
      }
    }
  }, [currentStorefront]);

  const handleInputChange = (
    e: React.ChangeEvent<
      HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement
    >
  ) => {
    const { name, value } = e.target;
    if (name.includes('.')) {
      const [parent, child] = name.split('.');
      setFormData((prev) => ({
        ...prev,
        [parent]: {
          ...(prev as any)[parent],
          [child]: value,
        },
      }));
    } else {
      setFormData((prev) => ({ ...prev, [name]: value }));
    }
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!currentStorefront) return;

    const updateData: B2CStoreUpdateDTO = {
      ...formData,
      settings: {
        ...formData.settings,
        business_hours: businessHours,
        delivery_providers: deliveryProviders,
        payment_methods: paymentMethods,
        notifications: notificationSettings,
        seo: seoSettings,
        security: securitySettings,
        business_info: businessInfo,
        coordinates: coordinates
          ? { lat: coordinates.lat, lng: coordinates.lng }
          : undefined,
      },
    };

    try {
      await dispatch(
        updateStorefront({
          id: currentStorefront.id!,
          data: updateData,
        })
      ).unwrap();

      toast.success(t('updateSuccess'));
      router.push(`/${locale}/b2c/${slug}/dashboard`);
    } catch {
      toast.error(t('updateError'));
    }
  };

  const handleBusinessHoursChange = (
    day: string,
    field: string,
    value: any
  ) => {
    setBusinessHours((prev) => ({
      ...prev,
      [day]: {
        ...prev[day as keyof typeof prev],
        [field]: value,
      },
    }));
  };

  const toggleDeliveryProvider = (providerId: string) => {
    setDeliveryProviders((prev) =>
      prev.map((provider) =>
        provider.id === providerId
          ? { ...provider, enabled: !provider.enabled }
          : provider
      )
    );
  };

  const updateDeliveryProviderSetting = (
    providerId: string,
    key: string,
    value: any
  ) => {
    setDeliveryProviders((prev) =>
      prev.map((provider) =>
        provider.id === providerId
          ? {
              ...provider,
              settings: {
                ...provider.settings,
                [key]: value,
              },
            }
          : provider
      )
    );
  };

  const togglePaymentMethod = (methodId: string) => {
    setPaymentMethods((prev) =>
      prev.map((method) =>
        method.id === methodId
          ? { ...method, enabled: !method.enabled }
          : method
      )
    );
  };

  const updatePaymentMethodSetting = (
    methodId: string,
    key: string,
    value: any
  ) => {
    setPaymentMethods((prev) =>
      prev.map((method) =>
        method.id === methodId
          ? {
              ...method,
              settings: {
                ...method.settings,
                [key]: value,
              },
            }
          : method
      )
    );
  };

  if (isLoading) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center">
          <span className="loading loading-spinner loading-lg text-primary"></span>
          <p className="mt-4 text-base-content/60">{tCommon('loading')}</p>
        </div>
      </div>
    );
  }

  if (error || !currentStorefront) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="card bg-base-100 shadow-xl max-w-md w-full">
          <div className="card-body text-center">
            <ExclamationCircleIcon className="w-16 h-16 mx-auto text-error mb-4" />
            <h2 className="card-title justify-center text-2xl">
              {error || 'Ошибка загрузки'}
            </h2>
            <div className="card-actions justify-center mt-6">
              <Link
                href={`/${locale}/b2c/${slug}/dashboard`}
                className="btn btn-primary"
              >
                <ArrowLeftIcon className="w-5 h-5" />
                {tCommon('back')}
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

  const tabs = [
    { id: 'general', label: 'Основные', icon: BuildingStorefrontIcon },
    { id: 'location', label: 'Местоположение', icon: MapPinIcon },
    { id: 'hours', label: 'Часы работы', icon: ClockIcon },
    { id: 'delivery', label: 'Доставка', icon: TruckIcon },
    { id: 'payment', label: 'Оплата', icon: CreditCardIcon },
    { id: 'media', label: 'Медиа', icon: PhotoIcon },
    { id: 'notifications', label: 'Уведомления', icon: BellIcon },
    { id: 'seo', label: 'SEO', icon: ChartBarIcon },
    { id: 'security', label: 'Безопасность', icon: ShieldCheckIcon },
    { id: 'business', label: 'Реквизиты', icon: DocumentTextIcon },
  ];

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-base-100 shadow-sm">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Link
                href={`/${locale}/b2c/${slug}/dashboard`}
                className="btn btn-ghost btn-circle"
              >
                <ArrowLeftIcon className="w-5 h-5" />
              </Link>
              <div>
                <h1 className="text-2xl font-bold flex items-center gap-2">
                  <CogIcon className="w-6 h-6" />
                  Настройки витрины
                </h1>
                <p className="text-base-content/60">{currentStorefront.name}</p>
              </div>
            </div>
            <button
              onClick={handleSubmit}
              className="btn btn-primary"
              disabled={isLoading}
            >
              {isLoading && <span className="loading loading-spinner"></span>}
              <CheckCircleIcon className="w-5 h-5" />
              {tCommon('save')}
            </button>
          </div>
        </div>
      </div>

      {/* Content */}
      <div className="container mx-auto px-4 py-6">
        <div className="flex flex-col lg:flex-row gap-6">
          {/* Sidebar */}
          <div className="lg:w-64">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body p-4">
                <ul className="menu">
                  {tabs.map((tab) => (
                    <li key={tab.id}>
                      <button
                        className={activeTab === tab.id ? 'active' : ''}
                        onClick={() => setActiveTab(tab.id)}
                      >
                        <tab.icon className="w-5 h-5" />
                        {tab.label}
                      </button>
                    </li>
                  ))}
                </ul>
              </div>
            </div>
          </div>

          {/* Main Content */}
          <div className="flex-1">
            <form
              onSubmit={handleSubmit}
              className="card bg-base-100 shadow-xl"
            >
              <div className="card-body">
                {/* General Tab */}
                {activeTab === 'general' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">
                      Основная информация
                    </h3>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Название витрины</span>
                      </label>
                      <input
                        type="text"
                        name="name"
                        value={formData.name}
                        onChange={handleInputChange}
                        className="input input-bordered"
                        required
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Описание</span>
                      </label>
                      <textarea
                        name="description"
                        value={formData.description}
                        onChange={handleInputChange}
                        className="textarea textarea-bordered h-32"
                        placeholder="Расскажите о вашем магазине..."
                      />
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Телефон</span>
                        </label>
                        <input
                          type="tel"
                          name="phone"
                          value={formData.phone}
                          onChange={handleInputChange}
                          className="input input-bordered"
                        />
                      </div>

                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Email</span>
                        </label>
                        <input
                          type="email"
                          name="email"
                          value={formData.email}
                          onChange={handleInputChange}
                          className="input input-bordered"
                        />
                      </div>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Веб-сайт</span>
                      </label>
                      <input
                        type="url"
                        name="website"
                        value={formData.website}
                        onChange={handleInputChange}
                        className="input input-bordered"
                        placeholder="https://example.com"
                      />
                    </div>

                    <div className="divider">Социальные сети</div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Instagram</span>
                        </label>
                        <input
                          type="text"
                          placeholder="@username"
                          className="input input-bordered"
                        />
                      </div>

                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Facebook</span>
                        </label>
                        <input
                          type="text"
                          placeholder="facebook.com/page"
                          className="input input-bordered"
                        />
                      </div>

                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Telegram</span>
                        </label>
                        <input
                          type="text"
                          placeholder="@channel"
                          className="input input-bordered"
                        />
                      </div>

                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">WhatsApp</span>
                        </label>
                        <input
                          type="text"
                          placeholder="+381..."
                          className="input input-bordered"
                        />
                      </div>
                    </div>
                  </div>
                )}

                {/* Location Tab */}
                {activeTab === 'location' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">
                      Адрес и местоположение
                    </h3>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Страна</span>
                        </label>
                        <select
                          name="location.country"
                          value={formData.location?.country || ''}
                          onChange={handleInputChange}
                          className="select select-bordered"
                        >
                          <option value="">Выберите страну</option>
                          <option value="RS">Сербия</option>
                          <option value="ME">Черногория</option>
                          <option value="BA">Босния и Герцеговина</option>
                          <option value="HR">Хорватия</option>
                        </select>
                      </div>

                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Город</span>
                        </label>
                        <input
                          type="text"
                          name="location.city"
                          value={formData.location?.city || ''}
                          onChange={handleInputChange}
                          className="input input-bordered"
                        />
                      </div>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Полный адрес</span>
                      </label>
                      <input
                        type="text"
                        name="location.full_address"
                        value={formData.location?.full_address || ''}
                        onChange={handleInputChange}
                        className="input input-bordered"
                        placeholder="Улица, дом, квартира"
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Почтовый индекс</span>
                      </label>
                      <input
                        type="text"
                        name="location.postal_code"
                        value={formData.location?.postal_code || ''}
                        onChange={handleInputChange}
                        className="input input-bordered"
                      />
                    </div>

                    <div className="divider">Координаты на карте</div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          Точное местоположение для покупателей
                        </span>
                      </label>
                      {coordinates ? (
                        <div className="alert alert-success">
                          <CheckCircleIcon className="w-5 h-5" />
                          <div>
                            <p>Координаты установлены</p>
                            <p className="text-sm opacity-80">
                              {coordinates.lat.toFixed(6)},{' '}
                              {coordinates.lng.toFixed(6)}
                            </p>
                          </div>
                          <button
                            type="button"
                            onClick={() => setShowLocationPicker(true)}
                            className="btn btn-sm"
                          >
                            Изменить
                          </button>
                        </div>
                      ) : (
                        <button
                          type="button"
                          onClick={() => setShowLocationPicker(true)}
                          className="btn btn-outline"
                        >
                          <MapPinIcon className="w-5 h-5" />
                          Указать на карте
                        </button>
                      )}
                    </div>

                    {showLocationPicker && (
                      <LocationPicker
                        value={
                          coordinates
                            ? {
                                address: formData.location?.full_address || '',
                                latitude: coordinates.lat,
                                longitude: coordinates.lng,
                                city: formData.location?.city || '',
                                country: formData.location?.country || '',
                                region: '',
                                confidence: 0.8,
                              }
                            : undefined
                        }
                        onChange={(location) => {
                          setCoordinates({
                            lat: location.latitude,
                            lng: location.longitude,
                          });
                          setFormData((prev) => ({
                            ...prev,
                            location: {
                              ...prev.location,
                              full_address: location.address,
                              city: location.city || prev.location?.city || '',
                              country:
                                location.country ||
                                prev.location?.country ||
                                '',
                            },
                          }));
                          setShowLocationPicker(false);
                        }}
                      />
                    )}
                  </div>
                )}

                {/* Business Hours Tab */}
                {activeTab === 'hours' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">Часы работы</h3>

                    <div className="space-y-4">
                      {Object.entries(businessHours).map(([day, hours]) => (
                        <div
                          key={day}
                          className="flex items-center gap-4 p-4 bg-base-200 rounded-lg"
                        >
                          <div className="form-control w-32">
                            <label className="label cursor-pointer">
                              <input
                                type="checkbox"
                                checked={hours.open}
                                onChange={(e) =>
                                  handleBusinessHoursChange(
                                    day,
                                    'open',
                                    e.target.checked
                                  )
                                }
                                className="checkbox checkbox-primary"
                              />
                              <span className="label-text capitalize ml-2 font-medium">
                                {day === 'monday' && 'Понедельник'}
                                {day === 'tuesday' && 'Вторник'}
                                {day === 'wednesday' && 'Среда'}
                                {day === 'thursday' && 'Четверг'}
                                {day === 'friday' && 'Пятница'}
                                {day === 'saturday' && 'Суббота'}
                                {day === 'sunday' && 'Воскресенье'}
                              </span>
                            </label>
                          </div>

                          <div className="flex gap-2 items-center flex-1">
                            <input
                              type="time"
                              value={hours.from}
                              onChange={(e) =>
                                handleBusinessHoursChange(
                                  day,
                                  'from',
                                  e.target.value
                                )
                              }
                              disabled={!hours.open}
                              className="input input-bordered"
                            />
                            <span className="text-base-content/60">—</span>
                            <input
                              type="time"
                              value={hours.to}
                              onChange={(e) =>
                                handleBusinessHoursChange(
                                  day,
                                  'to',
                                  e.target.value
                                )
                              }
                              disabled={!hours.open}
                              className="input input-bordered"
                            />
                            {!hours.open && (
                              <span className="text-error font-medium">
                                Выходной
                              </span>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>

                    <div className="alert alert-info">
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
                      <span>
                        Часы работы отображаются покупателям и влияют на
                        доступность для заказов
                      </span>
                    </div>
                  </div>
                )}

                {/* Delivery Tab */}
                {activeTab === 'delivery' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">Способы доставки</h3>

                    <div className="space-y-4">
                      {deliveryProviders.map((provider) => (
                        <div
                          key={provider.id}
                          className={`card ${
                            provider.enabled
                              ? 'bg-primary/5 border-primary'
                              : 'bg-base-200'
                          } border-2`}
                        >
                          <div className="card-body">
                            <div className="flex items-start justify-between">
                              <div className="flex items-start gap-4">
                                <div className="text-3xl">{provider.icon}</div>
                                <div>
                                  <h4 className="font-semibold text-lg flex items-center gap-2">
                                    {provider.name}
                                    {provider.id === 'post_express' && (
                                      <span className="badge badge-success badge-sm">
                                        NEW
                                      </span>
                                    )}
                                  </h4>
                                  <p className="text-sm text-base-content/70">
                                    {provider.description}
                                  </p>
                                </div>
                              </div>
                              <input
                                type="checkbox"
                                checked={provider.enabled}
                                onChange={() =>
                                  toggleDeliveryProvider(provider.id)
                                }
                                className="checkbox checkbox-primary"
                              />
                            </div>

                            {provider.enabled && provider.settings && (
                              <div className="mt-4 pl-14 space-y-3">
                                {/* Post Express специфические настройки */}
                                {provider.id === 'post_express' ? (
                                  <>
                                    {/* Учетные данные API */}
                                    <div className="divider text-sm">
                                      🔐 API интеграция WSP
                                    </div>
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                      <div className="form-control">
                                        <label className="label">
                                          <span className="label-text">
                                            Username WSP API
                                          </span>
                                        </label>
                                        <input
                                          type="text"
                                          value={
                                            provider.settings.api_username || ''
                                          }
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'api_username',
                                              e.target.value
                                            )
                                          }
                                          className="input input-bordered input-sm"
                                          placeholder="Ваш username..."
                                        />
                                      </div>
                                      <div className="form-control">
                                        <label className="label">
                                          <span className="label-text">
                                            Password WSP API
                                          </span>
                                        </label>
                                        <input
                                          type="password"
                                          value={
                                            provider.settings.api_password || ''
                                          }
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'api_password',
                                              e.target.value
                                            )
                                          }
                                          className="input input-bordered input-sm"
                                          placeholder="Ваш пароль..."
                                        />
                                      </div>
                                    </div>

                                    {/* Данные отправителя */}
                                    <div className="divider text-sm">
                                      📍 Данные отправителя
                                    </div>
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                      <div className="form-control">
                                        <label className="label">
                                          <span className="label-text">
                                            Название/ФИО
                                          </span>
                                        </label>
                                        <input
                                          type="text"
                                          value={
                                            provider.settings.sender_name || ''
                                          }
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'sender_name',
                                              e.target.value
                                            )
                                          }
                                          className="input input-bordered input-sm"
                                        />
                                      </div>
                                      <div className="form-control">
                                        <label className="label">
                                          <span className="label-text">
                                            Телефон
                                          </span>
                                        </label>
                                        <input
                                          type="tel"
                                          value={
                                            provider.settings.sender_phone || ''
                                          }
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'sender_phone',
                                              e.target.value
                                            )
                                          }
                                          className="input input-bordered input-sm"
                                        />
                                      </div>
                                      <div className="form-control">
                                        <label className="label">
                                          <span className="label-text">
                                            Email
                                          </span>
                                        </label>
                                        <input
                                          type="email"
                                          value={
                                            provider.settings.sender_email || ''
                                          }
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'sender_email',
                                              e.target.value
                                            )
                                          }
                                          className="input input-bordered input-sm"
                                        />
                                      </div>
                                      <div className="form-control">
                                        <label className="label">
                                          <span className="label-text">
                                            Город
                                          </span>
                                        </label>
                                        <input
                                          type="text"
                                          value={
                                            provider.settings.sender_city || ''
                                          }
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'sender_city',
                                              e.target.value
                                            )
                                          }
                                          className="input input-bordered input-sm"
                                        />
                                      </div>
                                    </div>
                                    <div className="grid grid-cols-1 md:grid-cols-3 gap-3">
                                      <div className="form-control md:col-span-2">
                                        <label className="label">
                                          <span className="label-text">
                                            Адрес отправки
                                          </span>
                                        </label>
                                        <input
                                          type="text"
                                          value={
                                            provider.settings.sender_address ||
                                            ''
                                          }
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'sender_address',
                                              e.target.value
                                            )
                                          }
                                          className="input input-bordered input-sm"
                                        />
                                      </div>
                                      <div className="form-control">
                                        <label className="label">
                                          <span className="label-text">
                                            Почтовый индекс
                                          </span>
                                        </label>
                                        <input
                                          type="text"
                                          value={
                                            provider.settings
                                              .sender_postal_code || ''
                                          }
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'sender_postal_code',
                                              e.target.value
                                            )
                                          }
                                          className="input input-bordered input-sm"
                                        />
                                      </div>
                                    </div>

                                    {/* Тарифы */}
                                    <div className="divider text-sm">
                                      💰 Тарифы по весу (RSD без НДС)
                                    </div>
                                    <div className="grid grid-cols-2 md:grid-cols-4 gap-3">
                                      <div className="stat bg-base-200 rounded-lg p-3">
                                        <div className="stat-title text-xs">
                                          0-2 кг
                                        </div>
                                        <div className="stat-value text-lg">
                                          340 RSD
                                        </div>
                                      </div>
                                      <div className="stat bg-base-200 rounded-lg p-3">
                                        <div className="stat-title text-xs">
                                          2-5 кг
                                        </div>
                                        <div className="stat-value text-lg">
                                          450 RSD
                                        </div>
                                      </div>
                                      <div className="stat bg-base-200 rounded-lg p-3">
                                        <div className="stat-title text-xs">
                                          5-10 кг
                                        </div>
                                        <div className="stat-value text-lg">
                                          580 RSD
                                        </div>
                                      </div>
                                      <div className="stat bg-base-200 rounded-lg p-3">
                                        <div className="stat-title text-xs">
                                          10-20 кг
                                        </div>
                                        <div className="stat-value text-lg">
                                          790 RSD
                                        </div>
                                      </div>
                                    </div>

                                    {/* Дополнительные услуги */}
                                    <div className="divider text-sm">
                                      🎁 Дополнительные услуги
                                    </div>
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                      <div className="card bg-base-200">
                                        <div className="card-body p-4">
                                          <div className="form-control">
                                            <label className="label cursor-pointer">
                                              <span className="label-text">
                                                <strong>
                                                  💵 Наложенный платеж (COD)
                                                </strong>
                                                <br />
                                                <span className="text-xs opacity-70">
                                                  Комиссия: 45 RSD за транзакцию
                                                </span>
                                              </span>
                                              <input
                                                type="checkbox"
                                                checked={
                                                  provider.settings
                                                    .enable_cod !== false
                                                }
                                                onChange={(e) =>
                                                  updateDeliveryProviderSetting(
                                                    provider.id,
                                                    'enable_cod',
                                                    e.target.checked
                                                  )
                                                }
                                                className="checkbox checkbox-primary"
                                              />
                                            </label>
                                          </div>
                                        </div>
                                      </div>
                                      <div className="card bg-base-200">
                                        <div className="card-body p-4">
                                          <div className="form-control">
                                            <label className="label cursor-pointer">
                                              <span className="label-text">
                                                <strong>🛡️ Страхование</strong>
                                                <br />
                                                <span className="text-xs opacity-70">
                                                  Бесплатно до 15,000 RSD
                                                </span>
                                              </span>
                                              <input
                                                type="checkbox"
                                                checked={
                                                  provider.settings
                                                    .enable_insurance !== false
                                                }
                                                onChange={(e) =>
                                                  updateDeliveryProviderSetting(
                                                    provider.id,
                                                    'enable_insurance',
                                                    e.target.checked
                                                  )
                                                }
                                                className="checkbox checkbox-primary"
                                              />
                                            </label>
                                          </div>
                                        </div>
                                      </div>
                                      <div className="card bg-base-200">
                                        <div className="card-body p-4">
                                          <div className="form-control">
                                            <label className="label cursor-pointer">
                                              <span className="label-text">
                                                <strong>
                                                  ⚡ Express доставка
                                                </strong>
                                                <br />
                                                <span className="text-xs opacity-70">
                                                  В тот же день
                                                  (Белград/Нови-Сад)
                                                </span>
                                              </span>
                                              <input
                                                type="checkbox"
                                                checked={
                                                  provider.settings
                                                    .enable_express || false
                                                }
                                                onChange={(e) =>
                                                  updateDeliveryProviderSetting(
                                                    provider.id,
                                                    'enable_express',
                                                    e.target.checked
                                                  )
                                                }
                                                className="checkbox checkbox-primary"
                                              />
                                            </label>
                                          </div>
                                        </div>
                                      </div>
                                      <div className="card bg-base-200">
                                        <div className="card-body p-4">
                                          <div className="form-control">
                                            <label className="label cursor-pointer">
                                              <span className="label-text">
                                                <strong>
                                                  📅 Субботняя доставка
                                                </strong>
                                                <br />
                                                <span className="text-xs opacity-70">
                                                  Доставка по субботам
                                                </span>
                                              </span>
                                              <input
                                                type="checkbox"
                                                checked={
                                                  provider.settings
                                                    .enable_saturday_delivery ||
                                                  false
                                                }
                                                onChange={(e) =>
                                                  updateDeliveryProviderSetting(
                                                    provider.id,
                                                    'enable_saturday_delivery',
                                                    e.target.checked
                                                  )
                                                }
                                                className="checkbox checkbox-primary"
                                              />
                                            </label>
                                          </div>
                                        </div>
                                      </div>
                                    </div>

                                    {/* Автоматизация */}
                                    <div className="divider text-sm">
                                      🤖 Автоматизация
                                    </div>
                                    <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                      <div className="form-control">
                                        <label className="label cursor-pointer">
                                          <span className="label-text">
                                            🖨️ Автоматическая печать этикеток
                                          </span>
                                          <input
                                            type="checkbox"
                                            checked={
                                              provider.settings
                                                .auto_print_labels || false
                                            }
                                            onChange={(e) =>
                                              updateDeliveryProviderSetting(
                                                provider.id,
                                                'auto_print_labels',
                                                e.target.checked
                                              )
                                            }
                                            className="checkbox checkbox-primary"
                                          />
                                        </label>
                                      </div>
                                      <div className="form-control">
                                        <label className="label cursor-pointer">
                                          <span className="label-text">
                                            📍 Автоматическое отслеживание
                                          </span>
                                          <input
                                            type="checkbox"
                                            checked={
                                              provider.settings
                                                .auto_track_shipments !== false
                                            }
                                            onChange={(e) =>
                                              updateDeliveryProviderSetting(
                                                provider.id,
                                                'auto_track_shipments',
                                                e.target.checked
                                              )
                                            }
                                            className="checkbox checkbox-primary"
                                          />
                                        </label>
                                      </div>
                                      <div className="form-control">
                                        <label className="label cursor-pointer">
                                          <span className="label-text">
                                            📤 Уведомления о заборе
                                          </span>
                                          <input
                                            type="checkbox"
                                            checked={
                                              provider.settings
                                                .notify_on_pickup !== false
                                            }
                                            onChange={(e) =>
                                              updateDeliveryProviderSetting(
                                                provider.id,
                                                'notify_on_pickup',
                                                e.target.checked
                                              )
                                            }
                                            className="checkbox checkbox-primary"
                                          />
                                        </label>
                                      </div>
                                      <div className="form-control">
                                        <label className="label cursor-pointer">
                                          <span className="label-text">
                                            📥 Уведомления о доставке
                                          </span>
                                          <input
                                            type="checkbox"
                                            checked={
                                              provider.settings
                                                .notify_on_delivery !== false
                                            }
                                            onChange={(e) =>
                                              updateDeliveryProviderSetting(
                                                provider.id,
                                                'notify_on_delivery',
                                                e.target.checked
                                              )
                                            }
                                            className="checkbox checkbox-primary"
                                          />
                                        </label>
                                      </div>
                                    </div>

                                    {/* Режим тестирования */}
                                    <div className="alert alert-warning">
                                      <svg
                                        xmlns="http://www.w3.org/2000/svg"
                                        className="stroke-current shrink-0 h-6 w-6"
                                        fill="none"
                                        viewBox="0 0 24 24"
                                      >
                                        <path
                                          strokeLinecap="round"
                                          strokeLinejoin="round"
                                          strokeWidth="2"
                                          d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                                        />
                                      </svg>
                                      <div>
                                        <div className="form-control">
                                          <label className="label cursor-pointer">
                                            <span className="label-text">
                                              <strong>Тестовый режим</strong> -
                                              использовать тестовое API
                                            </span>
                                            <input
                                              type="checkbox"
                                              checked={
                                                provider.settings.test_mode ||
                                                false
                                              }
                                              onChange={(e) =>
                                                updateDeliveryProviderSetting(
                                                  provider.id,
                                                  'test_mode',
                                                  e.target.checked
                                                )
                                              }
                                              className="checkbox checkbox-warning"
                                            />
                                          </label>
                                        </div>
                                      </div>
                                    </div>

                                    {/* Статистика интеграции */}
                                    <div className="stats shadow w-full">
                                      <div className="stat">
                                        <div className="stat-figure text-primary">
                                          <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            className="inline-block w-8 h-8 stroke-current"
                                          >
                                            <path
                                              strokeLinecap="round"
                                              strokeLinejoin="round"
                                              strokeWidth="2"
                                              d="M5 8h14M5 8a2 2 0 110-4h14a2 2 0 110 4M5 8v10a2 2 0 002 2h10a2 2 0 002-2V8m-9 4h4"
                                            ></path>
                                          </svg>
                                        </div>
                                        <div className="stat-title">
                                          Покрытие
                                        </div>
                                        <div className="stat-value text-primary">
                                          5000+
                                        </div>
                                        <div className="stat-desc">
                                          населенных пунктов
                                        </div>
                                      </div>

                                      <div className="stat">
                                        <div className="stat-figure text-secondary">
                                          <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            className="inline-block w-8 h-8 stroke-current"
                                          >
                                            <path
                                              strokeLinecap="round"
                                              strokeLinejoin="round"
                                              strokeWidth="2"
                                              d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"
                                            ></path>
                                          </svg>
                                        </div>
                                        <div className="stat-title">
                                          Почтовые отделения
                                        </div>
                                        <div className="stat-value text-secondary">
                                          180+
                                        </div>
                                        <div className="stat-desc">
                                          по всей Сербии
                                        </div>
                                      </div>

                                      <div className="stat">
                                        <div className="stat-figure text-accent">
                                          <svg
                                            xmlns="http://www.w3.org/2000/svg"
                                            fill="none"
                                            viewBox="0 0 24 24"
                                            className="inline-block w-8 h-8 stroke-current"
                                          >
                                            <path
                                              strokeLinecap="round"
                                              strokeLinejoin="round"
                                              strokeWidth="2"
                                              d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
                                            ></path>
                                          </svg>
                                        </div>
                                        <div className="stat-title">
                                          Срок доставки
                                        </div>
                                        <div className="stat-value text-accent">
                                          1-2
                                        </div>
                                        <div className="stat-desc">
                                          рабочих дня
                                        </div>
                                      </div>
                                    </div>

                                    {/* Информационный блок */}
                                    <div className="alert alert-info">
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
                                      <div className="text-sm">
                                        <p className="font-semibold mb-2">
                                          ✅ Post Express - полная интеграция
                                          включает:
                                        </p>
                                        <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                                          <ul className="space-y-1">
                                            <li>📦 WSP API интеграция v2.0</li>
                                            <li>
                                              🚚 &ldquo;Danas za sutra&rdquo;
                                              доставка
                                            </li>
                                            <li>🏢 180+ почтовых отделений</li>
                                            <li>📍 5000+ населенных пунктов</li>
                                          </ul>
                                          <ul className="space-y-1">
                                            <li>
                                              🖨️ Автоматическая печать этикеток
                                            </li>
                                            <li>📱 SMS/Email уведомления</li>
                                            <li>💵 Наложенный платеж (COD)</li>
                                            <li>📊 Полное отслеживание</li>
                                          </ul>
                                        </div>
                                        <p className="mt-2 text-xs opacity-80">
                                          Для начала работы получите учетные
                                          данные на
                                          <a
                                            href="https://onlinepostexpress.rs/registracija"
                                            target="_blank"
                                            rel="noopener noreferrer"
                                            className="link link-primary ml-1"
                                          >
                                            onlinepostexpress.rs
                                          </a>
                                        </p>
                                      </div>
                                    </div>
                                  </>
                                ) : (
                                  <>
                                    {/* Стандартные настройки для других провайдеров */}
                                    {provider.settings.api_key !==
                                      undefined && (
                                      <div className="form-control">
                                        <label className="label">
                                          <span className="label-text">
                                            API ключ
                                          </span>
                                        </label>
                                        <input
                                          type="password"
                                          value={provider.settings.api_key}
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'api_key',
                                              e.target.value
                                            )
                                          }
                                          className="input input-bordered input-sm"
                                          placeholder="Введите API ключ..."
                                        />
                                      </div>
                                    )}

                                    {provider.settings.base_rate !==
                                      undefined && (
                                      <div className="grid grid-cols-1 md:grid-cols-2 gap-3">
                                        <div className="form-control">
                                          <label className="label">
                                            <span className="label-text">
                                              Базовая стоимость (RSD)
                                            </span>
                                          </label>
                                          <input
                                            type="number"
                                            value={provider.settings.base_rate}
                                            onChange={(e) =>
                                              updateDeliveryProviderSetting(
                                                provider.id,
                                                'base_rate',
                                                Number(e.target.value)
                                              )
                                            }
                                            className="input input-bordered input-sm"
                                          />
                                        </div>

                                        {provider.settings
                                          .free_shipping_threshold !==
                                          undefined && (
                                          <div className="form-control">
                                            <label className="label">
                                              <span className="label-text">
                                                Бесплатно от (RSD)
                                              </span>
                                            </label>
                                            <input
                                              type="number"
                                              value={
                                                provider.settings
                                                  .free_shipping_threshold
                                              }
                                              onChange={(e) =>
                                                updateDeliveryProviderSetting(
                                                  provider.id,
                                                  'free_shipping_threshold',
                                                  Number(e.target.value)
                                                )
                                              }
                                              className="input input-bordered input-sm"
                                            />
                                          </div>
                                        )}
                                      </div>
                                    )}

                                    {provider.settings.estimated_days !==
                                      undefined && (
                                      <div className="form-control">
                                        <label className="label">
                                          <span className="label-text">
                                            Срок доставки (дней)
                                          </span>
                                        </label>
                                        <input
                                          type="number"
                                          value={
                                            provider.settings.estimated_days
                                          }
                                          onChange={(e) =>
                                            updateDeliveryProviderSetting(
                                              provider.id,
                                              'estimated_days',
                                              Number(e.target.value)
                                            )
                                          }
                                          className="input input-bordered input-sm w-32"
                                        />
                                      </div>
                                    )}
                                  </>
                                )}
                              </div>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>

                    <div className="alert alert-success">
                      <svg
                        xmlns="http://www.w3.org/2000/svg"
                        className="stroke-current shrink-0 h-6 w-6"
                        fill="none"
                        viewBox="0 0 24 24"
                      >
                        <path
                          strokeLinecap="round"
                          strokeLinejoin="round"
                          strokeWidth="2"
                          d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                      <div>
                        <h4 className="font-bold">
                          🚀 Post Express полностью интегрирован!
                        </h4>
                        <p className="text-sm mt-1">
                          Национальная почта Сербии готова к работе. Получите
                          учетные данные WSP API для начала использования.
                        </p>
                      </div>
                    </div>
                  </div>
                )}

                {/* Payment Methods Tab */}
                {activeTab === 'payment' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">Способы оплаты</h3>

                    <div className="space-y-4">
                      {paymentMethods.map((method) => (
                        <div
                          key={method.id}
                          className={`card ${
                            method.enabled
                              ? 'bg-primary/5 border-primary'
                              : 'bg-base-200'
                          } border-2`}
                        >
                          <div className="card-body">
                            <div className="flex items-start justify-between">
                              <div className="flex items-start gap-4">
                                <div className="text-3xl">{method.icon}</div>
                                <div>
                                  <h4 className="font-semibold text-lg">
                                    {method.name}
                                  </h4>
                                  <p className="text-sm text-base-content/70">
                                    {method.description}
                                  </p>
                                </div>
                              </div>
                              <input
                                type="checkbox"
                                checked={method.enabled}
                                onChange={() => togglePaymentMethod(method.id)}
                                className="checkbox checkbox-primary"
                              />
                            </div>

                            {method.enabled && method.settings && (
                              <div className="mt-4 pl-14 space-y-3">
                                {method.settings.account_number !==
                                  undefined && (
                                  <div className="form-control">
                                    <label className="label">
                                      <span className="label-text">
                                        Номер счета
                                      </span>
                                    </label>
                                    <input
                                      type="text"
                                      value={method.settings.account_number}
                                      onChange={(e) =>
                                        updatePaymentMethodSetting(
                                          method.id,
                                          'account_number',
                                          e.target.value
                                        )
                                      }
                                      className="input input-bordered input-sm"
                                      placeholder="Введите номер счета..."
                                    />
                                  </div>
                                )}

                                {method.settings.api_key !== undefined && (
                                  <div className="form-control">
                                    <label className="label">
                                      <span className="label-text">
                                        API ключ
                                      </span>
                                    </label>
                                    <input
                                      type="password"
                                      value={method.settings.api_key}
                                      onChange={(e) =>
                                        updatePaymentMethodSetting(
                                          method.id,
                                          'api_key',
                                          e.target.value
                                        )
                                      }
                                      className="input input-bordered input-sm"
                                      placeholder="Введите API ключ..."
                                    />
                                  </div>
                                )}

                                {method.settings.commission !== undefined && (
                                  <div className="form-control">
                                    <label className="label">
                                      <span className="label-text">
                                        Комиссия (%)
                                      </span>
                                    </label>
                                    <input
                                      type="number"
                                      step="0.1"
                                      value={method.settings.commission}
                                      onChange={(e) =>
                                        updatePaymentMethodSetting(
                                          method.id,
                                          'commission',
                                          Number(e.target.value)
                                        )
                                      }
                                      className="input input-bordered input-sm w-32"
                                      disabled
                                    />
                                  </div>
                                )}
                              </div>
                            )}
                          </div>
                        </div>
                      ))}
                    </div>
                  </div>
                )}

                {/* Media Tab */}
                {activeTab === 'media' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">Логотип и баннер</h3>

                    <div>
                      <h4 className="font-medium mb-4">Логотип</h4>
                      <div className="flex items-center gap-6">
                        <div className="avatar">
                          <div className="w-32 rounded-xl bg-base-200">
                            {logoPreview || currentStorefront.logo_url ? (
                              <Image
                                src={
                                  logoPreview ||
                                  currentStorefront.logo_url ||
                                  ''
                                }
                                alt="Logo"
                                fill
                                className="object-cover"
                              />
                            ) : (
                              <div className="w-full h-full flex items-center justify-center">
                                <PhotoIcon className="w-12 h-12 text-base-content/20" />
                              </div>
                            )}
                          </div>
                        </div>
                        <div>
                          <input
                            type="file"
                            id="logo-upload"
                            accept="image/jpeg,image/png,image/webp"
                            onChange={async (e) => {
                              const file = e.target.files?.[0];
                              if (!file) return;

                              if (file.size > 5 * 1024 * 1024) {
                                toast.error('Файл слишком большой (макс. 5MB)');
                                e.target.value = '';
                                return;
                              }

                              const reader = new FileReader();
                              reader.onload = (e) => {
                                setLogoPreview(e.target?.result as string);
                              };
                              reader.readAsDataURL(file);

                              try {
                                setUploadingLogo(true);
                                await storefrontApi.uploadLogo(
                                  currentStorefront.id!,
                                  file
                                );
                                toast.success('Логотип загружен');
                                dispatch(fetchStorefrontBySlug(slug));
                              } catch {
                                toast.error('Ошибка загрузки логотипа');
                                setLogoPreview(null);
                              } finally {
                                setUploadingLogo(false);
                                e.target.value = '';
                              }
                            }}
                            className="hidden"
                          />
                          <div className="flex gap-2">
                            <label
                              htmlFor="logo-upload"
                              className={`btn btn-primary btn-sm ${uploadingLogo ? 'loading' : ''}`}
                            >
                              {uploadingLogo
                                ? 'Загрузка...'
                                : 'Загрузить логотип'}
                            </label>
                            {(logoPreview || currentStorefront.logo_url) && (
                              <button
                                type="button"
                                onClick={() => setLogoPreview(null)}
                                className="btn btn-ghost btn-sm text-error"
                              >
                                Удалить
                              </button>
                            )}
                          </div>
                          <p className="text-sm text-base-content/60 mt-2">
                            PNG, JPG или WebP до 5MB. Рекомендуется 512x512px
                          </p>
                        </div>
                      </div>
                    </div>

                    <div className="divider"></div>

                    <div>
                      <h4 className="font-medium mb-4">Баннер</h4>
                      <div className="space-y-4">
                        <div className="aspect-[3/1] bg-base-200 rounded-xl overflow-hidden relative max-h-64">
                          {bannerPreview || currentStorefront.banner_url ? (
                            <Image
                              src={
                                bannerPreview ||
                                currentStorefront.banner_url ||
                                ''
                              }
                              alt="Banner"
                              fill
                              className="object-cover"
                            />
                          ) : (
                            <div className="w-full h-full flex items-center justify-center">
                              <PhotoIcon className="w-16 h-16 text-base-content/20" />
                            </div>
                          )}
                        </div>
                        <input
                          type="file"
                          id="banner-upload"
                          accept="image/jpeg,image/png,image/webp"
                          onChange={async (e) => {
                            const file = e.target.files?.[0];
                            if (!file) return;

                            if (file.size > 10 * 1024 * 1024) {
                              toast.error('Файл слишком большой (макс. 10MB)');
                              e.target.value = '';
                              return;
                            }

                            const reader = new FileReader();
                            reader.onload = (e) => {
                              setBannerPreview(e.target?.result as string);
                            };
                            reader.readAsDataURL(file);

                            try {
                              setUploadingBanner(true);
                              await storefrontApi.uploadBanner(
                                currentStorefront.id!,
                                file
                              );
                              toast.success('Баннер загружен');
                              dispatch(fetchStorefrontBySlug(slug));
                            } catch {
                              toast.error('Ошибка загрузки баннера');
                              setBannerPreview(null);
                            } finally {
                              setUploadingBanner(false);
                              e.target.value = '';
                            }
                          }}
                          className="hidden"
                        />
                        <div className="flex gap-2">
                          <label
                            htmlFor="banner-upload"
                            className={`btn btn-primary btn-sm ${uploadingBanner ? 'loading' : ''}`}
                          >
                            {uploadingBanner
                              ? 'Загрузка...'
                              : 'Загрузить баннер'}
                          </label>
                          {(bannerPreview || currentStorefront.banner_url) && (
                            <button
                              type="button"
                              onClick={() => setBannerPreview(null)}
                              className="btn btn-ghost btn-sm text-error"
                            >
                              Удалить
                            </button>
                          )}
                        </div>
                        <p className="text-sm text-base-content/60">
                          PNG, JPG или WebP до 10MB. Рекомендуется 1920x640px
                        </p>
                      </div>
                    </div>
                  </div>
                )}

                {/* Notifications Tab */}
                {activeTab === 'notifications' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">Уведомления</h3>

                    <div className="space-y-4">
                      <div className="form-control">
                        <label className="label cursor-pointer">
                          <span className="label-text">Email уведомления</span>
                          <input
                            type="checkbox"
                            checked={notificationSettings.email_notifications}
                            onChange={(e) =>
                              setNotificationSettings((prev) => ({
                                ...prev,
                                email_notifications: e.target.checked,
                              }))
                            }
                            className="checkbox checkbox-primary"
                          />
                        </label>
                      </div>

                      <div className="form-control">
                        <label className="label cursor-pointer">
                          <span className="label-text">
                            Telegram уведомления
                          </span>
                          <input
                            type="checkbox"
                            checked={
                              notificationSettings.telegram_notifications
                            }
                            onChange={(e) =>
                              setNotificationSettings((prev) => ({
                                ...prev,
                                telegram_notifications: e.target.checked,
                              }))
                            }
                            className="checkbox checkbox-primary"
                          />
                        </label>
                      </div>

                      <div className="divider">События для уведомлений</div>

                      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                        <div className="form-control">
                          <label className="label cursor-pointer">
                            <span className="label-text">Новые заказы</span>
                            <input
                              type="checkbox"
                              checked={notificationSettings.new_orders}
                              onChange={(e) =>
                                setNotificationSettings((prev) => ({
                                  ...prev,
                                  new_orders: e.target.checked,
                                }))
                              }
                              className="checkbox checkbox-primary"
                            />
                          </label>
                        </div>

                        <div className="form-control">
                          <label className="label cursor-pointer">
                            <span className="label-text">Новые сообщения</span>
                            <input
                              type="checkbox"
                              checked={notificationSettings.new_messages}
                              onChange={(e) =>
                                setNotificationSettings((prev) => ({
                                  ...prev,
                                  new_messages: e.target.checked,
                                }))
                              }
                              className="checkbox checkbox-primary"
                            />
                          </label>
                        </div>

                        <div className="form-control">
                          <label className="label cursor-pointer">
                            <span className="label-text">Новые отзывы</span>
                            <input
                              type="checkbox"
                              checked={notificationSettings.new_reviews}
                              onChange={(e) =>
                                setNotificationSettings((prev) => ({
                                  ...prev,
                                  new_reviews: e.target.checked,
                                }))
                              }
                              className="checkbox checkbox-primary"
                            />
                          </label>
                        </div>

                        <div className="form-control">
                          <label className="label cursor-pointer">
                            <span className="label-text">
                              Низкий запас товара
                            </span>
                            <input
                              type="checkbox"
                              checked={notificationSettings.low_stock}
                              onChange={(e) =>
                                setNotificationSettings((prev) => ({
                                  ...prev,
                                  low_stock: e.target.checked,
                                }))
                              }
                              className="checkbox checkbox-primary"
                            />
                          </label>
                        </div>

                        <div className="form-control">
                          <label className="label cursor-pointer">
                            <span className="label-text">
                              Ежедневная сводка
                            </span>
                            <input
                              type="checkbox"
                              checked={notificationSettings.daily_summary}
                              onChange={(e) =>
                                setNotificationSettings((prev) => ({
                                  ...prev,
                                  daily_summary: e.target.checked,
                                }))
                              }
                              className="checkbox checkbox-primary"
                            />
                          </label>
                        </div>

                        <div className="form-control">
                          <label className="label cursor-pointer">
                            <span className="label-text">
                              Еженедельный отчет
                            </span>
                            <input
                              type="checkbox"
                              checked={notificationSettings.weekly_report}
                              onChange={(e) =>
                                setNotificationSettings((prev) => ({
                                  ...prev,
                                  weekly_report: e.target.checked,
                                }))
                              }
                              className="checkbox checkbox-primary"
                            />
                          </label>
                        </div>
                      </div>
                    </div>
                  </div>
                )}

                {/* SEO Tab */}
                {activeTab === 'seo' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">SEO и аналитика</h3>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Мета-заголовок</span>
                      </label>
                      <input
                        type="text"
                        value={seoSettings.meta_title}
                        onChange={(e) =>
                          setSeoSettings((prev) => ({
                            ...prev,
                            meta_title: e.target.value,
                          }))
                        }
                        className="input input-bordered"
                        placeholder={currentStorefront.name}
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Мета-описание</span>
                      </label>
                      <textarea
                        value={seoSettings.meta_description}
                        onChange={(e) =>
                          setSeoSettings((prev) => ({
                            ...prev,
                            meta_description: e.target.value,
                          }))
                        }
                        className="textarea textarea-bordered h-24"
                        placeholder="Краткое описание для поисковых систем..."
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Ключевые слова</span>
                      </label>
                      <input
                        type="text"
                        value={seoSettings.meta_keywords}
                        onChange={(e) =>
                          setSeoSettings((prev) => ({
                            ...prev,
                            meta_keywords: e.target.value,
                          }))
                        }
                        className="input input-bordered"
                        placeholder="магазин, товары, доставка..."
                      />
                    </div>

                    <div className="divider">Аналитика</div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Google Analytics ID</span>
                      </label>
                      <input
                        type="text"
                        value={seoSettings.google_analytics}
                        onChange={(e) =>
                          setSeoSettings((prev) => ({
                            ...prev,
                            google_analytics: e.target.value,
                          }))
                        }
                        className="input input-bordered"
                        placeholder="G-XXXXXXXXXX"
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Facebook Pixel ID</span>
                      </label>
                      <input
                        type="text"
                        value={seoSettings.facebook_pixel}
                        onChange={(e) =>
                          setSeoSettings((prev) => ({
                            ...prev,
                            facebook_pixel: e.target.value,
                          }))
                        }
                        className="input input-bordered"
                        placeholder="XXXXXXXXXXXXXXX"
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Яндекс.Метрика ID</span>
                      </label>
                      <input
                        type="text"
                        value={seoSettings.yandex_metrika}
                        onChange={(e) =>
                          setSeoSettings((prev) => ({
                            ...prev,
                            yandex_metrika: e.target.value,
                          }))
                        }
                        className="input input-bordered"
                        placeholder="XXXXXXXX"
                      />
                    </div>
                  </div>
                )}

                {/* Security Tab */}
                {activeTab === 'security' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">Безопасность</h3>

                    <div className="form-control">
                      <label className="label cursor-pointer">
                        <span className="label-text">
                          <strong>Двухфакторная аутентификация</strong>
                          <p className="text-sm text-base-content/70">
                            Дополнительная защита вашего аккаунта
                          </p>
                        </span>
                        <input
                          type="checkbox"
                          checked={securitySettings.two_factor_auth}
                          onChange={(e) =>
                            setSecuritySettings((prev) => ({
                              ...prev,
                              two_factor_auth: e.target.checked,
                            }))
                          }
                          className="checkbox checkbox-primary"
                        />
                      </label>
                    </div>

                    <div className="divider">Ограничение доступа по IP</div>

                    <div className="form-control">
                      <label className="label cursor-pointer">
                        <span className="label-text">
                          <strong>Белый список IP-адресов</strong>
                          <p className="text-sm text-base-content/70">
                            Разрешить доступ только с указанных IP
                          </p>
                        </span>
                        <input
                          type="checkbox"
                          checked={securitySettings.ip_whitelist}
                          onChange={(e) =>
                            setSecuritySettings((prev) => ({
                              ...prev,
                              ip_whitelist: e.target.checked,
                            }))
                          }
                          className="checkbox checkbox-primary"
                        />
                      </label>
                    </div>

                    {securitySettings.ip_whitelist && (
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">
                            Разрешенные IP-адреса
                          </span>
                        </label>
                        <textarea
                          value={securitySettings.allowed_ips}
                          onChange={(e) =>
                            setSecuritySettings((prev) => ({
                              ...prev,
                              allowed_ips: e.target.value,
                            }))
                          }
                          className="textarea textarea-bordered h-24"
                          placeholder="192.168.1.1&#10;10.0.0.1"
                        />
                        <label className="label">
                          <span className="label-text-alt">
                            По одному IP-адресу на строку
                          </span>
                        </label>
                      </div>
                    )}

                    <div className="divider">API доступ</div>

                    <div className="form-control">
                      <label className="label cursor-pointer">
                        <span className="label-text">
                          <strong>Разрешить API доступ</strong>
                          <p className="text-sm text-base-content/70">
                            Позволяет интегрироваться с внешними системами
                          </p>
                        </span>
                        <input
                          type="checkbox"
                          checked={securitySettings.api_access}
                          onChange={(e) =>
                            setSecuritySettings((prev) => ({
                              ...prev,
                              api_access: e.target.checked,
                            }))
                          }
                          className="checkbox checkbox-primary"
                        />
                      </label>
                    </div>

                    {securitySettings.api_access && (
                      <>
                        <div className="form-control">
                          <label className="label">
                            <span className="label-text">API ключ</span>
                          </label>
                          <div className="input-group">
                            <input
                              type="password"
                              value={securitySettings.api_key}
                              readOnly
                              className="input input-bordered flex-1"
                            />
                            <button
                              type="button"
                              onClick={() => {
                                const newKey =
                                  'sk_' +
                                  Math.random().toString(36).substring(2, 15);
                                setSecuritySettings((prev) => ({
                                  ...prev,
                                  api_key: newKey,
                                }));
                              }}
                              className="btn btn-square"
                            >
                              <KeyIcon className="w-5 h-5" />
                            </button>
                          </div>
                        </div>

                        <div className="form-control">
                          <label className="label">
                            <span className="label-text">Webhook Secret</span>
                          </label>
                          <input
                            type="password"
                            value={securitySettings.webhook_secret}
                            onChange={(e) =>
                              setSecuritySettings((prev) => ({
                                ...prev,
                                webhook_secret: e.target.value,
                              }))
                            }
                            className="input input-bordered"
                          />
                        </div>
                      </>
                    )}
                  </div>
                )}

                {/* Business Info Tab */}
                {activeTab === 'business' && (
                  <div className="space-y-6">
                    <h3 className="text-lg font-semibold">
                      Реквизиты компании
                    </h3>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Тип бизнеса</span>
                      </label>
                      <select
                        value={businessInfo.business_type}
                        onChange={(e) =>
                          setBusinessInfo((prev) => ({
                            ...prev,
                            business_type: e.target.value,
                          }))
                        }
                        className="select select-bordered"
                      >
                        <option value="individual">
                          Индивидуальный предприниматель
                        </option>
                        <option value="company">Компания (ООО, АО)</option>
                        <option value="self_employed">Самозанятый</option>
                      </select>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Юридическое название</span>
                      </label>
                      <input
                        type="text"
                        value={businessInfo.legal_name}
                        onChange={(e) =>
                          setBusinessInfo((prev) => ({
                            ...prev,
                            legal_name: e.target.value,
                          }))
                        }
                        className="input input-bordered"
                      />
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">Юридический адрес</span>
                      </label>
                      <input
                        type="text"
                        value={businessInfo.legal_address}
                        onChange={(e) =>
                          setBusinessInfo((prev) => ({
                            ...prev,
                            legal_address: e.target.value,
                          }))
                        }
                        className="input input-bordered"
                      />
                    </div>

                    <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">
                            Регистрационный номер
                          </span>
                        </label>
                        <input
                          type="text"
                          value={businessInfo.registration_number}
                          onChange={(e) =>
                            setBusinessInfo((prev) => ({
                              ...prev,
                              registration_number: e.target.value,
                            }))
                          }
                          className="input input-bordered"
                        />
                      </div>

                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">ИНН</span>
                        </label>
                        <input
                          type="text"
                          value={businessInfo.tax_number}
                          onChange={(e) =>
                            setBusinessInfo((prev) => ({
                              ...prev,
                              tax_number: e.target.value,
                            }))
                          }
                          className="input input-bordered"
                        />
                      </div>
                    </div>

                    <div className="form-control">
                      <label className="label">
                        <span className="label-text">
                          VAT номер (если применимо)
                        </span>
                      </label>
                      <input
                        type="text"
                        value={businessInfo.vat_number}
                        onChange={(e) =>
                          setBusinessInfo((prev) => ({
                            ...prev,
                            vat_number: e.target.value,
                          }))
                        }
                        className="input input-bordered"
                        placeholder="RS123456789"
                      />
                    </div>

                    <div className="alert alert-info">
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
                      <span>
                        Эти данные используются для выставления счетов и
                        налоговой отчетности
                      </span>
                    </div>
                  </div>
                )}
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
  );
}
