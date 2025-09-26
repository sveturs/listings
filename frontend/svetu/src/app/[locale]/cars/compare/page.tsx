import ComparisonTable from '@/components/cars/ComparisonTable';
import { use } from 'react';

type PageProps = {
  params: Promise<{
    locale: string;
  }>;
};

export default function ComparisonPage({ params }: PageProps) {
  const { locale } = use(params);

  return <ComparisonTable locale={locale} />;
}
