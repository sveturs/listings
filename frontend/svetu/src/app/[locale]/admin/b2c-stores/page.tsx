import { Suspense } from 'react';
import { setRequestLocale } from 'next-intl/server';
import { getTranslations } from 'next-intl/server';
import AdminB2CStoresTable from './AdminB2CStoresTable';

interface PageProps {
  params: Promise<{ locale: string }>;
}

// Делаем страницу динамической из-за API вызовов
export const dynamic = 'force-dynamic';

export default async function StorefrontsPage({ params }: PageProps) {
  const { locale } = await params;
  setRequestLocale(locale);
  const t = await getTranslations('admin');

  return (
    <div>
      <Suspense
        fallback={
          <div className="flex justify-center items-center h-64">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        }
      >
        <AdminB2CStoresTable />
      </Suspense>
    </div>
  );
}
