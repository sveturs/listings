import { Suspense } from 'react';
import { getTranslations } from 'next-intl/server';
import ListingGrid from '@/components/marketplace/ListingGrid';
import ListingFilters from '@/components/marketplace/ListingFilters';
import Loading from './loading';

export default async function MarketplacePage() {
  const t = await getTranslations('marketplace');
  
  return (
    <div className="flex gap-6">
      <aside className="w-64 flex-shrink-0">
        <ListingFilters />
      </aside>
      
      <main className="flex-1">
        <h1 className="text-3xl font-bold mb-6">{t('title')}</h1>
        <Suspense fallback={<Loading />}>
          <ListingGrid />
        </Suspense>
      </main>
    </div>
  );
}