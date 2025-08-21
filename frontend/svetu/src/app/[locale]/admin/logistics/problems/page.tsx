import { getTranslations } from 'next-intl/server';
import { Metadata } from 'next';
import ProblemShipmentsClient from './problems-client';

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'admin' });

  return {
    title: `${t('logistics.problems.title')} | ${t('title')}`,
    description: t('logistics.problems.description'),
  };
}

export default async function ProblemsPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'admin' });

  return (
    <div className="space-y-6">
      <div className="border-b pb-4">
        <h1 className="text-3xl font-bold">{t('logistics.problems.title')}</h1>
        <p className="text-gray-600 mt-2">
          {t('logistics.problems.description')}
        </p>
      </div>

      <ProblemShipmentsClient />
    </div>
  );
}
