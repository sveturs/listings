'use client';

import { useTranslations } from 'next-intl';
import { useCreateStorefrontContext } from '@/contexts/CreateStorefrontContext';

interface PreviewPublishStepProps {
  onBack: () => void;
  onComplete: () => void;
}

const dayNames = [
  'sunday',
  'monday',
  'tuesday',
  'wednesday',
  'thursday',
  'friday',
  'saturday',
];

export default function PreviewPublishStep({
  onBack,
  onComplete,
}: PreviewPublishStepProps) {
  const t = useTranslations('create_storefront');
  const tCommon = useTranslations('common');
  const tPermissions = useTranslations('permissions');
  const { formData, isSubmitting, submitStorefront } =
    useCreateStorefrontContext();

  const handlePublish = async () => {
    const result = await submitStorefront();
    if (result.success) {
      onComplete();
    }
  };

  return (
    <div className="max-w-4xl mx-auto">
      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <h2 className="card-title text-2xl mb-4">
            {t('preview.title')}
          </h2>
          <p className="text-base-content/70 mb-6">
            {t('preview.subtitle')}
          </p>

          {/* Preview sections */}
          <div className="space-y-6">
            {/* Basic Info */}
            <div className="card bg-base-200">
              <div className="card-body">
                <h3 className="card-title text-lg">
                  {t('basic_info')}
                </h3>
                <div className="space-y-2">
                  <p>
                    <strong>{t('basic_info.name')}:</strong>{' '}
                    {formData.name}
                  </p>
                  <p>
                    <strong>{t('basic_info.slug')}:</strong>{' '}
                    svetu.rs/{formData.slug}
                  </p>
                  <p>
                    <strong>
                      {t('basic_info.description')}:
                    </strong>{' '}
                    {formData.description}
                  </p>
                  <p>
                    <strong>
                      {t('basic_info.business_type')}:
                    </strong>{' '}
                    {t(
                      `create_storefront.business_types.${formData.businessType}`
                    )}
                  </p>
                </div>
              </div>
            </div>

            {/* Business Details */}
            <div className="card bg-base-200">
              <div className="card-body">
                <h3 className="card-title text-lg">
                  {t('business_details')}
                </h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                  {formData.registrationNumber && (
                    <p>
                      <strong>
                        {t(
                          'create_storefront.business_details.registration_number'
                        )}
                        :
                      </strong>{' '}
                      {formData.registrationNumber}
                    </p>
                  )}
                  {formData.taxNumber && (
                    <p>
                      <strong>
                        {t('business_details.tax_number')}:
                      </strong>{' '}
                      {formData.taxNumber}
                    </p>
                  )}
                  {formData.vatNumber && (
                    <p>
                      <strong>
                        {t('business_details.vat_number')}:
                      </strong>{' '}
                      {formData.vatNumber}
                    </p>
                  )}
                  {formData.phone && (
                    <p>
                      <strong>
                        {t('business_details.phone')}:
                      </strong>{' '}
                      {formData.phone}
                    </p>
                  )}
                  {formData.email && (
                    <p>
                      <strong>
                        {t('business_details.email')}:
                      </strong>{' '}
                      {formData.email}
                    </p>
                  )}
                  {formData.website && (
                    <p>
                      <strong>
                        {t('business_details.website')}:
                      </strong>{' '}
                      {formData.website}
                    </p>
                  )}
                </div>
              </div>
            </div>

            {/* Location */}
            <div className="card bg-base-200">
              <div className="card-body">
                <h3 className="card-title text-lg">
                  {t('location')}
                </h3>
                <div className="space-y-2">
                  <p>
                    <strong>{t('location.address')}:</strong>{' '}
                    {formData.address}
                  </p>
                  <p>
                    <strong>{t('location.city')}:</strong>{' '}
                    {formData.city}
                  </p>
                  <p>
                    <strong>
                      {t('location.postal_code')}:
                    </strong>{' '}
                    {formData.postalCode}
                  </p>
                  <p>
                    <strong>{t('location.country')}:</strong>{' '}
                    {t(`countries.${formData.country}`)}
                  </p>
                  {formData.latitude && formData.longitude && (
                    <p>
                      <strong>
                        {t('location.coordinates')}:
                      </strong>{' '}
                      {formData.latitude.toFixed(6)},{' '}
                      {formData.longitude.toFixed(6)}
                    </p>
                  )}
                </div>
              </div>
            </div>

            {/* Business Hours */}
            <div className="card bg-base-200">
              <div className="card-body">
                <h3 className="card-title text-lg">
                  {t('business_hours')}
                </h3>
                <div className="space-y-1">
                  {formData.businessHours.map((hours) => (
                    <p key={hours.dayOfWeek}>
                      <strong>
                        {t(`common.days.${dayNames[hours.dayOfWeek]}`)}:
                      </strong>{' '}
                      {hours.isClosed
                        ? tCommon('closed')
                        : `${hours.openTime} - ${hours.closeTime}`}
                    </p>
                  ))}
                </div>
              </div>
            </div>

            {/* Payment & Delivery */}
            <div className="card bg-base-200">
              <div className="card-body">
                <h3 className="card-title text-lg">
                  {t('payment_delivery')}
                </h3>

                {formData.paymentMethods &&
                  formData.paymentMethods.length > 0 && (
                    <div className="mb-4">
                      <h4 className="font-semibold mb-2">
                        {t(
                          'create_storefront.payment_delivery.payment_methods'
                        )}
                        :
                      </h4>
                      <div className="flex flex-wrap gap-2">
                        {formData.paymentMethods.map((method) => (
                          <span key={method} className="badge badge-primary">
                            {t(`payment_methods.${method}`)}
                          </span>
                        ))}
                      </div>
                    </div>
                  )}

                {formData.deliveryOptions &&
                  formData.deliveryOptions.length > 0 && (
                    <div>
                      <h4 className="font-semibold mb-2">
                        {t(
                          'create_storefront.payment_delivery.delivery_options'
                        )}
                        :
                      </h4>
                      <div className="space-y-1">
                        {formData.deliveryOptions.map((option, index) => (
                          <p key={index}>
                            {t(`delivery_providers.${option.providerName}`)} -{' '}
                            {option.deliveryTimeMinutes} min,{' '}
                            {option.deliveryCostRSD} RSD
                            {option.freeDeliveryThresholdRSD &&
                              ` (${tCommon('free_above')} ${option.freeDeliveryThresholdRSD} RSD)`}
                          </p>
                        ))}
                      </div>
                    </div>
                  )}
              </div>
            </div>

            {/* Staff */}
            {formData.staff && formData.staff.length > 0 && (
              <div className="card bg-base-200">
                <div className="card-body">
                  <h3 className="card-title text-lg">
                    {t('staff_setup')}
                  </h3>
                  <div className="space-y-2">
                    {formData.staff.map((member, index) => (
                      <div key={index}>
                        <p>
                          <strong>{member.email}</strong> -{' '}
                          {t(`roles.${member.role}`)}
                        </p>
                        <div className="flex gap-2 text-sm text-base-content/70">
                          {member.canManageProducts && (
                            <span>• {t('products')}</span>
                          )}
                          {member.canManageOrders && (
                            <span>• {t('orders')}</span>
                          )}
                          {member.canManageSettings && (
                            <span>• {t('settings')}</span>
                          )}
                        </div>
                      </div>
                    ))}
                  </div>
                </div>
              </div>
            )}
          </div>

          <div className="card-actions justify-between mt-8">
            <button
              className="btn btn-ghost"
              onClick={onBack}
              disabled={isSubmitting}
            >
              {tCommon('back')}
            </button>
            <button
              className={`btn btn-primary ${isSubmitting ? 'loading' : ''}`}
              onClick={handlePublish}
              disabled={isSubmitting}
            >
              {isSubmitting ? tCommon('publishing') : tCommon('publish')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
