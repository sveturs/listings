import { notFound } from 'next/navigation';
import { getTranslations, setRequestLocale } from 'next-intl/server';
import { routing } from '@/i18n/routing';
import TrackingClient from './TrackingClient';

interface PageProps {
  params: Promise<{
    locale: string;
    token: string;
  }>;
}

export async function generateMetadata({ params }: PageProps) {
  const { locale, token } = await params;
  const t = await getTranslations({ locale, namespace: 'tracking' });

  return {
    title: `${t('title')} - ${token}`,
    description: t('viber.description'),
  };
}

export default async function TrackingPage({ params }: PageProps) {
  const { locale, token } = await params;

  // Ensure that the incoming locale is valid
  if (!routing.locales.includes(locale as any)) {
    notFound();
  }

  // Set the locale for this request
  setRequestLocale(locale);

  return <TrackingClient token={token} />;
}

// Enable static rendering for known locales
export function generateStaticParams() {
  return routing.locales.map((locale) => ({ locale }));
}