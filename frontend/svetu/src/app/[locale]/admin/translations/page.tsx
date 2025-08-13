import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';
import EnhancedTranslationsDashboard from '@/components/admin/translations/EnhancedTranslationsDashboard';

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

export default function TranslationsAdminPage() {
  // Используем расширенную версию с полной функциональностью
  return <EnhancedTranslationsDashboard />;
}
