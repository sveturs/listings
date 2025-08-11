import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'admin' });

  return {
    title: t('translations.pageTitle'),
    description: t('translations.pageDescription'),
  };
}

export default function TranslationsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return <div className="min-h-screen bg-base-200">{children}</div>;
}
