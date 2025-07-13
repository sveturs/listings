'use client';

import dynamic from 'next/dynamic';

const DistrictMapSearch = dynamic(
  () => import('@/components/search/DistrictMapSearch'),
  { ssr: false }
);

export default function DistrictSearchPage() {
  return (
    <main className="min-h-screen">
      <DistrictMapSearch />
    </main>
  );
}
