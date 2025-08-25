import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';
import ShipmentsListClient from './shipments-client';

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'admin' });

  return {
    title: `${t('logistics.shipments.title')} | ${t('title')}`,
    description: t('logistics.shipments.description'),
  };
}

export default async function ShipmentsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'admin' });

  return (
    <div className="container mx-auto px-4 py-6">
      {/* Header */}
      <div className="mb-6">
        <h1 className="text-3xl font-bold">{t('logistics.shipments.title')}</h1>
        <p className="text-base-content/70 mt-2">
          {t('logistics.shipments.description')}
        </p>
      </div>

      {/* Client Component */}
      <ShipmentsListClient />
    </div>
  );
}
