'use client';

import { PageTransition } from '@/components/ui/PageTransition';
import { Link } from '@/i18n/routing';
import { BentoGrid } from '@/components/ui/BentoGrid';
import { BentoGridCategories } from '@/components/ui/BentoGridCategories';
import { BentoGridListings } from '@/components/ui/BentoGridListings';
import { useEffect, useState, useCallback } from 'react';

interface HomePageData {
  categories: Array<{
    id: string;
    name: string;
    count: number;
  }>;
  featuredListing?: {
    id: string;
    title: string;
    price: string;
    image: string;
    category: string;
  };
  stats: {
    totalListings: number;
    activeUsers: number;
    successfulDeals: number;
  };
  popularSearches: any[];
  nearbyListings: Array<{
    id: string;
    latitude: number;
    longitude: number;
    price: number;
  }>;
  error: Error | null;
}

interface HomePageClientProps {
  title: string;
  description: string;
  createListingText: string;
  initialData: any;
  homePageData: HomePageData | null;
  locale: string;
  error: Error | null;
  paymentsEnabled: boolean;
}

export default function HomePageClient({
  createListingText,
  initialData,
  homePageData,
  locale,
  error,
  paymentsEnabled,
}: HomePageClientProps) {
  const [bentoData, setBentoData] = useState<HomePageData | null>(homePageData);
  const [isLoading, setIsLoading] = useState(!homePageData);
  const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(null);
  const [filters, setFilters] = useState<Record<string, any>>({});

  // Стабильные коллбэки для предотвращения лишних рендеров
  const handleCategorySelect = useCallback((categoryId: number | null) => {
    setSelectedCategoryId(categoryId);
  }, []);

  const handleFiltersChange = useCallback((newFilters: Record<string, any>) => {
    setFilters(newFilters);
  }, []);

  useEffect(() => {
    // Если данные не были загружены на сервере (dev mode), загружаем на клиенте
    if (!homePageData && typeof window !== 'undefined') {
      console.log('[HomePageClient] No SSR data, loading on client...');
      import('./actions').then(({ getHomePageData }) => {
        getHomePageData(locale).then((data) => {
          console.log('[HomePageClient] Client data loaded:', data);
          setBentoData(data);
          setIsLoading(false);
        });
      });
    } else if (homePageData) {
      console.log('[HomePageClient] Using SSR data:', homePageData);
    }
  }, [homePageData, locale]);
  return (
    <PageTransition mode="fade">
      <div className="min-h-screen">
        <div className="container mx-auto px-4 pt-8">
          {/* BentoGrid секция */}
          <div className="mb-12">
            {isLoading ? (
              // Skeleton loader для BentoGrid
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 p-4">
                {Array.from({ length: 8 }).map((_, i) => (
                  <div
                    key={i}
                    className={`skeleton h-40 ${i < 2 ? 'col-span-1 md:col-span-2 lg:col-span-2 row-span-2' : ''}`}
                  ></div>
                ))}
              </div>
            ) : (
              <BentoGrid
                categories={bentoData?.categories || []}
                featuredListing={bentoData?.featuredListing}
                stats={
                  bentoData?.stats || {
                    totalListings: 0,
                    activeUsers: 0,
                    successfulDeals: 0,
                  }
                }
                nearbyListings={bentoData?.nearbyListings || []}
              />
            )}
          </div>

          {/* Новая секция BentoGrid с категориями и объявлениями */}
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 p-4">
            <BentoGridCategories 
              onCategorySelect={handleCategorySelect}
              selectedCategoryId={selectedCategoryId}
              filters={filters}
              onFiltersChange={handleFiltersChange}
            />
            <BentoGridListings
              locale={locale}
              selectedCategoryId={selectedCategoryId}
              filters={filters}
            />
          </div>

          {/* Плавающая кнопка создания объявления */}
          <Link
            href="/create-listing"
            className="fixed bottom-6 right-6 btn btn-primary btn-circle btn-lg shadow-xl hover:shadow-2xl hover:scale-110 transition-all duration-200 z-50"
            title={createListingText}
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-6 w-6"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 4v16m8-8H4"
              />
            </svg>
          </Link>
        </div>
      </div>
    </PageTransition>
  );
}
