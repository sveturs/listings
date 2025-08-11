import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';
import TranslationsDashboardDemo from '@/components/admin/translations/TranslationsDashboardDemo';

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
  // Используем демо-версию для тестирования
  return <TranslationsDashboardDemo />;
}
