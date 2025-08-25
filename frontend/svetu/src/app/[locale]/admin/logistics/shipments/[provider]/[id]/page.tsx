import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';
import ShipmentDetailsClient from './shipment-details-client';

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string; provider: string; id: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'admin' });

  return {
    title: `${t('logistics.shipmentDetails.title')} | ${t('title')}`,
    description: t('logistics.shipmentDetails.description'),
  };
}

export default async function ShipmentDetailsPage({
  params,
}: {
  params: Promise<{ locale: string; provider: string; id: string }>;
}) {
  const { locale, provider, id } = await params;
  const t = await getTranslations({ locale, namespace: 'admin' });

  return (
    <div className="container mx-auto px-4 py-6">
      {/* Breadcrumbs */}
      <div className="breadcrumbs text-sm mb-4">
        <ul>
          <li>
            <a href={`/${locale}/admin/logistics`}>{t('logistics.title')}</a>
          </li>
          <li>
            <a href={`/${locale}/admin/logistics/shipments`}>
              {t('logistics.shipments.title')}
            </a>
          </li>
          <li>{t('logistics.shipmentDetails.title')}</li>
        </ul>
      </div>

      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold">
          {t('logistics.shipmentDetails.title')}
        </h1>
        <p className="text-base-content/70 mt-2">
          {t('logistics.shipmentDetails.description')}
        </p>
      </div>

      {/* Client Component */}
      <ShipmentDetailsClient provider={provider} id={id} />
    </div>
  );
}
