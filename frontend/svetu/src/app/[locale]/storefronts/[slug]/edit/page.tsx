'use client';

import { useEffect, useState } from 'react';
import { useParams, useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchStorefrontBySlug,
  updateStorefront,
} from '@/store/slices/storefrontSlice';
import { StorefrontUpdateDTO } from '@/types/storefront';
import {
  BuildingStorefrontIcon,
  MapPinIcon,
  PhoneIcon,
  EnvelopeIcon,
  GlobeAltIcon,
  ClockIcon,
  CreditCardIcon,
  TruckIcon,
  PhotoIcon,
  ArrowLeftIcon,
  CheckCircleIcon,
  ExclamationCircleIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';
import { toast } from '@/utils/toast';

export default function EditStorefrontPage() {
  const params = useParams();
  const router = useRouter();
  const dispatch = useAppDispatch();
  const t = useTranslations();
  const slug = params.slug as string;

  const { currentStorefront, isLoading, error } = useAppSelector(
    (state) => state.storefronts
  );

  const [formData, setFormData] = useState<StorefrontUpdateDTO>({
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

  const [activeTab, setActiveTab] = useState<
    'basic' | 'location' | 'hours' | 'payment' | 'delivery' | 'media'
  >('basic');

  const [businessHours, setBusinessHours] = useState({
    monday: { open: true, from: '09:00', to: '18:00' },
    tuesday: { open: true, from: '09:00', to: '18:00' },
    wednesday: { open: true, from: '09:00', to: '18:00' },
    thursday: { open: true, from: '09:00', to: '18:00' },
    friday: { open: true, from: '09:00', to: '18:00' },
    saturday: { open: true, from: '10:00', to: '16:00' },
    sunday: { open: false, from: '10:00', to: '16:00' },
  });

  const [paymentMethods, setPaymentMethods] = useState({
    cash: true,
    card: true,
    bank_transfer: false,
    crypto: false,
  });

  const [deliveryOptions, setDeliveryOptions] = useState({
    pickup: true,
    local_delivery: true,
    shipping: false,
  });

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

      // Загружаем часы работы из настроек
      if ((currentStorefront as any).settings?.business_hours) {
        setBusinessHours((currentStorefront as any).settings.business_hours);
      }

      // Загружаем методы оплаты
      if ((currentStorefront as any).settings?.payment_methods) {
        setPaymentMethods((currentStorefront as any).settings.payment_methods);
      }

      // Загружаем способы доставки
      if ((currentStorefront as any).settings?.delivery_options) {
        setDeliveryOptions(
          (currentStorefront as any).settings.delivery_options
        );
      }
    }
  }, [currentStorefront]);

  const handleInputChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
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

    const updateData: StorefrontUpdateDTO = {
      ...formData,
      settings: {
        ...formData.settings,
        business_hours: businessHours,
        payment_methods: paymentMethods,
        delivery_options: deliveryOptions,
      },
    };

    try {
      await dispatch(
        updateStorefront({
          id: currentStorefront.id!,
          data: updateData,
        })
      ).unwrap();

      toast.success(t('storefronts.updateSuccess'));
      router.push(`/storefronts/${slug}/dashboard`);
    } catch {
      toast.error(t('storefronts.updateError'));
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

  if (isLoading) {
    return (
      <div className="min-h-screen bg-base-200 flex items-center justify-center">
        <div className="text-center">
          <span className="loading loading-spinner loading-lg text-primary"></span>
          <p className="mt-4 text-base-content/60">{t('common.loading')}</p>
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
              {t('storefronts.notFound')}
            </h2>
            <p className="text-base-content/70">
              {error || t('storefronts.loadError')}
            </p>
            <div className="card-actions justify-center mt-6">
              <Link href="/profile/storefronts" className="btn btn-primary">
                <ArrowLeftIcon className="w-5 h-5" />
                {t('common.back')}
              </Link>
            </div>
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-base-100 shadow-sm">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Link
                href={`/storefronts/${slug}/dashboard`}
                className="btn btn-ghost btn-circle"
              >
                <ArrowLeftIcon className="w-5 h-5" />
              </Link>
              <div>
                <h1 className="text-2xl font-bold">
                  {t('storefronts.editStorefront')}
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
              {t('common.save')}
            </button>
          </div>
        </div>
      </div>

      {/* Tabs */}
      <div className="container mx-auto px-4 py-6">
        <div className="tabs tabs-boxed bg-base-100 mb-6">
          <button
            className={`tab ${activeTab === 'basic' ? 'tab-active' : ''}`}
            onClick={() => setActiveTab('basic')}
          >
            <BuildingStorefrontIcon className="w-4 h-4 mr-2" />
            {t('storefronts.basicInfo')}
          </button>
          <button
            className={`tab ${activeTab === 'location' ? 'tab-active' : ''}`}
            onClick={() => setActiveTab('location')}
          >
            <MapPinIcon className="w-4 h-4 mr-2" />
            {t('storefronts.location')}
          </button>
          <button
            className={`tab ${activeTab === 'hours' ? 'tab-active' : ''}`}
            onClick={() => setActiveTab('hours')}
          >
            <ClockIcon className="w-4 h-4 mr-2" />
            {t('storefronts.businessHours')}
          </button>
          <button
            className={`tab ${activeTab === 'payment' ? 'tab-active' : ''}`}
            onClick={() => setActiveTab('payment')}
          >
            <CreditCardIcon className="w-4 h-4 mr-2" />
            {t('storefronts.payment')}
          </button>
          <button
            className={`tab ${activeTab === 'delivery' ? 'tab-active' : ''}`}
            onClick={() => setActiveTab('delivery')}
          >
            <TruckIcon className="w-4 h-4 mr-2" />
            {t('storefronts.delivery')}
          </button>
          <button
            className={`tab ${activeTab === 'media' ? 'tab-active' : ''}`}
            onClick={() => setActiveTab('media')}
          >
            <PhotoIcon className="w-4 h-4 mr-2" />
            {t('storefronts.media')}
          </button>
        </div>

        {/* Form Content */}
        <form onSubmit={handleSubmit} className="card bg-base-100 shadow-xl">
          <div className="card-body">
            {/* Basic Info Tab */}
            {activeTab === 'basic' && (
              <div className="space-y-6">
                <div className="form-control">
                  <label className="label">
                    <span className="label-text">
                      {t('storefronts.storeName')}
                    </span>
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
                    <span className="label-text">
                      {t('storefronts.description')}
                    </span>
                  </label>
                  <textarea
                    name="description"
                    value={formData.description}
                    onChange={handleInputChange}
                    className="textarea textarea-bordered h-32"
                    placeholder={t('storefronts.descriptionPlaceholder')}
                  />
                </div>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">
                        {t('storefronts.phone')}
                      </span>
                    </label>
                    <div className="input-group">
                      <span className="bg-base-200">
                        <PhoneIcon className="w-5 h-5" />
                      </span>
                      <input
                        type="tel"
                        name="phone"
                        value={formData.phone}
                        onChange={handleInputChange}
                        className="input input-bordered flex-1"
                      />
                    </div>
                  </div>

                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">
                        {t('storefronts.email')}
                      </span>
                    </label>
                    <div className="input-group">
                      <span className="bg-base-200">
                        <EnvelopeIcon className="w-5 h-5" />
                      </span>
                      <input
                        type="email"
                        name="email"
                        value={formData.email}
                        onChange={handleInputChange}
                        className="input input-bordered flex-1"
                      />
                    </div>
                  </div>
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">
                      {t('storefronts.website')}
                    </span>
                  </label>
                  <div className="input-group">
                    <span className="bg-base-200">
                      <GlobeAltIcon className="w-5 h-5" />
                    </span>
                    <input
                      type="url"
                      name="website"
                      value={formData.website}
                      onChange={handleInputChange}
                      className="input input-bordered flex-1"
                      placeholder="https://example.com"
                    />
                  </div>
                </div>
              </div>
            )}

            {/* Location Tab */}
            {activeTab === 'location' && (
              <div className="space-y-6">
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">
                        {t('storefronts.country')}
                      </span>
                    </label>
                    <input
                      type="text"
                      name="location.country"
                      value={formData.location?.country || ''}
                      onChange={handleInputChange}
                      className="input input-bordered"
                    />
                  </div>

                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">
                        {t('storefronts.city')}
                      </span>
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
                    <span className="label-text">
                      {t('storefronts.address')}
                    </span>
                  </label>
                  <input
                    type="text"
                    name="location.full_address"
                    value={formData.location?.full_address || ''}
                    onChange={handleInputChange}
                    className="input input-bordered"
                  />
                </div>

                <div className="form-control">
                  <label className="label">
                    <span className="label-text">
                      {t('storefronts.postalCode')}
                    </span>
                  </label>
                  <input
                    type="text"
                    name="location.postal_code"
                    value={formData.location?.postal_code || ''}
                    onChange={handleInputChange}
                    className="input input-bordered"
                  />
                </div>
              </div>
            )}

            {/* Business Hours Tab */}
            {activeTab === 'hours' && (
              <div className="space-y-4">
                {Object.entries(businessHours).map(([day, hours]) => (
                  <div key={day} className="flex items-center gap-4">
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
                        <span className="label-text capitalize ml-2">
                          {t(`days.${day}`)}
                        </span>
                      </label>
                    </div>

                    <div className="flex gap-2 items-center flex-1">
                      <input
                        type="time"
                        value={hours.from}
                        onChange={(e) =>
                          handleBusinessHoursChange(day, 'from', e.target.value)
                        }
                        disabled={!hours.open}
                        className="input input-bordered input-sm"
                      />
                      <span className="text-base-content/60">-</span>
                      <input
                        type="time"
                        value={hours.to}
                        onChange={(e) =>
                          handleBusinessHoursChange(day, 'to', e.target.value)
                        }
                        disabled={!hours.open}
                        className="input input-bordered input-sm"
                      />
                    </div>
                  </div>
                ))}
              </div>
            )}

            {/* Payment Methods Tab */}
            {activeTab === 'payment' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold mb-4">
                  {t('storefronts.acceptedPaymentMethods')}
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  {Object.entries(paymentMethods).map(([method, enabled]) => (
                    <div key={method} className="form-control">
                      <label className="label cursor-pointer">
                        <span className="label-text">
                          {t(`storefronts.payment.${method}`)}
                        </span>
                        <input
                          type="checkbox"
                          checked={enabled}
                          onChange={(e) =>
                            setPaymentMethods((prev) => ({
                              ...prev,
                              [method]: e.target.checked,
                            }))
                          }
                          className="checkbox checkbox-primary"
                        />
                      </label>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Delivery Options Tab */}
            {activeTab === 'delivery' && (
              <div className="space-y-4">
                <h3 className="text-lg font-semibold mb-4">
                  {t('storefronts.deliveryOptions')}
                </h3>
                <div className="space-y-4">
                  {Object.entries(deliveryOptions).map(([option, enabled]) => (
                    <div key={option} className="form-control">
                      <label className="label cursor-pointer justify-start">
                        <input
                          type="checkbox"
                          checked={enabled}
                          onChange={(e) =>
                            setDeliveryOptions((prev) => ({
                              ...prev,
                              [option]: e.target.checked,
                            }))
                          }
                          className="checkbox checkbox-primary mr-3"
                        />
                        <div>
                          <span className="label-text font-medium">
                            {t(`storefronts.delivery.${option}`)}
                          </span>
                          <p className="text-sm text-base-content/60">
                            {t(`storefronts.delivery.${option}Description`)}
                          </p>
                        </div>
                      </label>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Media Tab */}
            {activeTab === 'media' && (
              <div className="space-y-6">
                <div>
                  <h3 className="text-lg font-semibold mb-4">
                    {t('storefronts.logo')}
                  </h3>
                  <div className="flex items-center gap-6">
                    <div className="avatar">
                      <div className="w-32 rounded-xl bg-base-200">
                        {currentStorefront.logo_url ? (
                          <img src={currentStorefront.logo_url} alt="Logo" />
                        ) : (
                          <div className="w-full h-full flex items-center justify-center">
                            <PhotoIcon className="w-12 h-12 text-base-content/20" />
                          </div>
                        )}
                      </div>
                    </div>
                    <div>
                      <button type="button" className="btn btn-primary btn-sm">
                        {t('storefronts.uploadLogo')}
                      </button>
                      <p className="text-sm text-base-content/60 mt-2">
                        {t('storefronts.logoRequirements')}
                      </p>
                    </div>
                  </div>
                </div>

                <div className="divider"></div>

                <div>
                  <h3 className="text-lg font-semibold mb-4">
                    {t('storefronts.banner')}
                  </h3>
                  <div className="space-y-4">
                    <div className="aspect-[3/1] bg-base-200 rounded-xl overflow-hidden">
                      {currentStorefront.banner_url ? (
                        <img
                          src={currentStorefront.banner_url}
                          alt="Banner"
                          className="w-full h-full object-cover"
                        />
                      ) : (
                        <div className="w-full h-full flex items-center justify-center">
                          <PhotoIcon className="w-16 h-16 text-base-content/20" />
                        </div>
                      )}
                    </div>
                    <button type="button" className="btn btn-primary btn-sm">
                      {t('storefronts.uploadBanner')}
                    </button>
                    <p className="text-sm text-base-content/60">
                      {t('storefronts.bannerRequirements')}
                    </p>
                  </div>
                </div>
              </div>
            )}
          </div>
        </form>
      </div>
    </div>
  );
}
