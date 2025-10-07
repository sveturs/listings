import { Metadata } from 'next';
import { getTranslations } from 'next-intl/server';
import CategoryProposalsList from '@/components/admin/CategoryProposalsList';

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}): Promise<Metadata> {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'admin.categoryProposals' });

  return {
    title: t('title'),
    description: t('description'),
  };
}

export default function CategoryProposalsPage() {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-8">
      <div className="container mx-auto px-4">
        <CategoryProposalsList />
      </div>
    </div>
  );
}
