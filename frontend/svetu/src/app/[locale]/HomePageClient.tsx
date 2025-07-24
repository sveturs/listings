'use client';

import { PageTransition } from '@/components/ui/PageTransition';
import HomePage from '@/components/marketplace/HomePage';
import { Link } from '@/i18n/routing';
import { SearchBar } from '@/components/SearchBar';
import { BentoGrid } from '@/components/ui/BentoGrid';

interface HomePageClientProps {
  title: string;
  description: string;
  createListingText: string;
  initialData: any;
  locale: string;
  error: Error | null;
  paymentsEnabled: boolean;
}

export default function HomePageClient({
  title,
  description,
  createListingText,
  initialData,
  locale,
  error,
  paymentsEnabled,
}: HomePageClientProps) {
  return (
    <PageTransition mode="fade">
      <div className="min-h-screen">
        {/* Hero секция */}
        <div className="bg-gradient-to-b from-base-200/50 to-base-100 py-12 lg:py-16 mb-8">
          <div className="container mx-auto px-4">
            <h1 className="text-4xl lg:text-5xl font-bold text-center mb-4">
              {title}
            </h1>
            <p className="text-center text-base-content/70 text-lg max-w-2xl mx-auto mb-8">
              {description}
            </p>

            {/* Search Bar с поддержкой fuzzy search */}
            <div className="max-w-3xl mx-auto">
              <SearchBar variant="hero" showTrending={true} />
            </div>
          </div>
        </div>

        <div className="container mx-auto px-4">
          {/* BentoGrid секция */}
          <div className="mb-12">
            <h2 className="text-2xl font-bold text-center mb-8">
              Популярные категории и рекомендации
            </h2>
            <BentoGrid
              categories={[
                { id: 'electronics', name: 'Электроника', count: 1243 },
                { id: 'fashion', name: 'Одежда и обувь', count: 856 },
                { id: 'home', name: 'Дом и сад', count: 642 },
                { id: 'auto', name: 'Автотовары', count: 521 },
                { id: 'books', name: 'Книги', count: 387 },
              ]}
              featuredListing={{
                id: '12345',
                title: 'iPhone 15 Pro в отличном состоянии',
                price: '89,000 ₽',
                image: '/api/placeholder/300/200',
                category: 'Телефоны',
              }}
              stats={{
                totalListings: 15420,
                activeUsers: 2840,
                successfulDeals: 8932,
              }}
            />
          </div>

          <HomePage
            initialData={initialData}
            locale={locale}
            error={error}
            paymentsEnabled={paymentsEnabled}
          />

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
