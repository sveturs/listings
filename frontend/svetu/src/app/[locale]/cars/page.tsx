import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';
import CarsPageClient from './CarsPageClient';

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'cars' });

  return {
    title: t('pageTitle'),
    description: t('pageDescription'),
  };
}

export default async function CarsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;

  return <CarsPageClient locale={locale} />;
}
