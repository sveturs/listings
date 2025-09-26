'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import api from '@/services/api';
import type { components } from '@/types/generated/api';
import { Car, Search, TrendingUp, Calendar, Filter } from 'lucide-react';
import CarSortingOptions, {
  type CarSortOption,
} from '@/components/marketplace/CarSortingOptions';
import CarQuickFilters from '@/components/marketplace/CarQuickFilters';
import { CarFilters } from '@/components/marketplace/CarFilters';
import CarBrandIcon from '@/components/marketplace/CarBrandIcon';

type CarMake = components['schemas']['backend_internal_domain_models.CarMake'];
type MarketplaceListing =
  components['schemas']['backend_internal_domain_models.MarketplaceListing'];

interface CarsPageClientProps {
  locale: string;
}

export default function CarsPageClient({ locale }: CarsPageClientProps) {
  const t = useTranslations('cars');
  const router = useRouter();
  const [loading, setLoading] = useState(true);
  const [popularMakes, setPopularMakes] = useState<CarMake[]>([]);
  const [latestListings, setLatestListings] = useState<MarketplaceListing[]>(
    []
  );
  const [searchQuery, setSearchQuery] = useState('');
  const [stats, setStats] = useState({
    totalListings: 0,
    totalMakes: 0,
    totalModels: 0,
  });
  const [sortOption, setSortOption] =
    useState<CarSortOption>('created_at_desc');
  const [selectedQuickFilters, setSelectedQuickFilters] = useState<string[]>(
    []
  );
  const [showFilters, setShowFilters] = useState(false);
  const [activeFilters, setActiveFilters] = useState<Record<string, any>>({});
  const [searchResults, setSearchResults] = useState<MarketplaceListing[]>([]);
  const [isSearching, setIsSearching] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);

      // Загружаем популярные марки
      const makesResponse = await api.get('/api/v1/cars/makes', {
        params: { limit: 12 },
      });
      if (makesResponse.data?.data) {
        setPopularMakes(makesResponse.data.data.slice(0, 12));
      }

      // Загружаем последние автомобильные объявления
      const searchParams = {
        limit: 8,
        offset: 0,
        category_ids: '10101,10102,10103,10104', // Автомобильные категории
        sort: 'created_at_desc',
      };

      const listingsResponse = await api.post(
        '/api/v1/marketplace/search',
        searchParams
      );

      if (listingsResponse.data?.data?.items) {
        setLatestListings(listingsResponse.data.data.items);
        setStats({
          totalListings: listingsResponse.data.data.total || 0,
          totalMakes: makesResponse.data.data?.length || 0,
          totalModels: 3788, // Из БД
        });
      }
    } catch (error) {
      console.error('Error loading cars data:', error);
    } finally {
      setLoading(false);
    }
  };

  const handleSearch = () => {
    if (searchQuery) {
      router.push(
        `/${locale}/search?q=${encodeURIComponent(searchQuery)}&category=10101&context=automotive`
      );
    }
  };

  const handleMakeClick = (makeSlug: string) => {
    router.push(
      `/${locale}/search?category=10101&car_make=${makeSlug}&context=automotive`
    );
  };

  const handleQuickFilterToggle = (
    filterId: string,
    filter: Record<string, any>
  ) => {
    setSelectedQuickFilters((prev) => {
      const newFilters = prev.includes(filterId)
        ? prev.filter((id) => id !== filterId)
        : [...prev, filterId];

      // Apply the filter
      const newActiveFilters = { ...activeFilters };
      if (newFilters.includes(filterId)) {
        Object.assign(newActiveFilters, filter);
      } else {
        // Remove the filter
        Object.keys(filter).forEach((key) => {
          delete newActiveFilters[key];
        });
      }
      setActiveFilters(newActiveFilters);

      return newFilters;
    });
  };

  const searchCars = async () => {
    try {
      setIsSearching(true);

      // Prepare search parameters
      const searchParams: any = {
        limit: 20,
        offset: 0,
        category_ids: '10101,10102,10103,10104',
        sort: sortOption,
      };

      // Add search query if present
      if (searchQuery) {
        searchParams.q = searchQuery;
      }

      // Add filters
      Object.entries(activeFilters).forEach(([key, value]) => {
        if (key === 'priceMax') {
          searchParams.price_max = value;
        } else if (key === 'mileageMax') {
          searchParams['attributes.mileage_max'] = value;
        } else if (key === 'fuelType') {
          searchParams['attributes.fuel_type'] = value;
        } else if (key === 'condition') {
          searchParams.condition = value;
        } else if (key === 'bodyTypes') {
          searchParams['attributes.body_type'] = value.join(',');
        } else if (key === 'make') {
          searchParams['attributes.car_make'] = value;
        } else if (key === 'model') {
          searchParams['attributes.car_model'] = value;
        } else if (key === 'yearFrom') {
          searchParams['attributes.year_min'] = value;
        } else if (key === 'yearTo') {
          searchParams['attributes.year_max'] = value;
        }
      });

      const response = await api.post(
        '/api/v1/marketplace/search',
        searchParams
      );

      if (response.data?.data?.items) {
        setSearchResults(response.data.data.items);
      }
    } catch (error) {
      console.error('Error searching cars:', error);
    } finally {
      setIsSearching(false);
    }
  };

  // Trigger search when filters or sort changes
  useEffect(() => {
    if (
      Object.keys(activeFilters).length > 0 ||
      sortOption !== 'created_at_desc'
    ) {
      searchCars();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeFilters, sortOption]);

  if (loading) {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-100">
      {/* Hero Section */}
      <div className="hero min-h-[400px] bg-gradient-to-br from-primary to-primary-focus text-primary-content">
        <div className="hero-content text-center">
          <div className="max-w-md">
            <h1 className="text-5xl font-bold mb-5 flex items-center justify-center gap-3">
              <Car className="w-12 h-12" />
              {t('heroTitle')}
            </h1>
            <p className="mb-8">{t('heroDescription')}</p>

            {/* Search Bar */}
            <div className="join w-full">
              <input
                type="text"
                placeholder={t('searchPlaceholder')}
                className="input input-bordered join-item flex-1 text-base-content"
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                onKeyDown={(e) => e.key === 'Enter' && handleSearch()}
              />
              <button
                className="btn btn-secondary join-item"
                onClick={handleSearch}
              >
                <Search className="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      </div>

      {/* Statistics */}
      <div className="bg-base-200 py-8">
        <div className="container mx-auto px-4">
          <div className="stats stats-horizontal shadow w-full">
            <div className="stat">
              <div className="stat-figure text-primary">
                <Car className="w-8 h-8" />
              </div>
              <div className="stat-title">{t('stats.listings')}</div>
              <div className="stat-value text-primary">
                {stats.totalListings.toLocaleString()}
              </div>
            </div>
            <div className="stat">
              <div className="stat-figure text-secondary">
                <TrendingUp className="w-8 h-8" />
              </div>
              <div className="stat-title">{t('stats.makes')}</div>
              <div className="stat-value text-secondary">
                {stats.totalMakes}
              </div>
            </div>
            <div className="stat">
              <div className="stat-figure text-accent">
                <Calendar className="w-8 h-8" />
              </div>
              <div className="stat-title">{t('stats.models')}</div>
              <div className="stat-value text-accent">
                {stats.totalModels.toLocaleString()}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Popular Makes */}
      <div className="container mx-auto px-4 py-12">
        <h2 className="text-3xl font-bold mb-8">{t('popularMakes')}</h2>
        <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-6 gap-4">
          {popularMakes.map((make) => (
            <div
              key={make.id}
              className="card bg-base-100 shadow-xl cursor-pointer hover:shadow-2xl transition-shadow"
              onClick={() => make.slug && handleMakeClick(make.slug)}
            >
              <div className="card-body items-center text-center p-4">
                <CarBrandIcon brand={make.name || ''} className="w-16 h-16" />
                <h3 className="card-title text-sm">{make.name || ''}</h3>
              </div>
            </div>
          ))}
        </div>

        <div className="text-center mt-8">
          <Link
            href={`/${locale}/search?category=10101&context=automotive`}
            className="btn btn-primary"
          >
            {t('viewAllMakes')}
          </Link>
        </div>
      </div>

      {/* Categories */}
      <div className="bg-base-200 py-12">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-8">{t('carCategories')}</h2>
          <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
            <Link
              href={`/${locale}/search?category=10101&context=automotive`}
              className="btn btn-outline btn-lg"
            >
              {t('categories.passenger')}
            </Link>
            <Link
              href={`/${locale}/search?category=10102&context=automotive`}
              className="btn btn-outline btn-lg"
            >
              {t('categories.suv')}
            </Link>
            <Link
              href={`/${locale}/search?category=10103&context=automotive`}
              className="btn btn-outline btn-lg"
            >
              {t('categories.commercial')}
            </Link>
            <Link
              href={`/${locale}/search?category=10104&context=automotive`}
              className="btn btn-outline btn-lg"
            >
              {t('categories.motorcycle')}
            </Link>
          </div>
        </div>
      </div>

      {/* Advanced Search Section */}
      <div className="container mx-auto px-4 py-12">
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <div className="flex flex-col lg:flex-row gap-4">
              {/* Left side - Filters */}
              <div
                className={`lg:w-1/4 ${showFilters ? '' : 'hidden lg:block'}`}
              >
                <CarFilters
                  onFiltersChange={setActiveFilters}
                  className="w-full"
                />
              </div>

              {/* Right side - Results */}
              <div className="flex-1">
                {/* Controls Bar */}
                <div className="flex flex-col gap-4 mb-6">
                  {/* Mobile filter toggle */}
                  <button
                    className="btn btn-outline lg:hidden"
                    onClick={() => setShowFilters(!showFilters)}
                  >
                    <Filter className="w-5 h-5" />
                    {showFilters ? t('hideFilters') : t('showFilters')}
                  </button>

                  {/* Quick Filters */}
                  <CarQuickFilters
                    selectedFilters={selectedQuickFilters}
                    onToggleFilter={handleQuickFilterToggle}
                  />

                  {/* Sorting Options */}
                  <div className="flex justify-between items-center">
                    <CarSortingOptions
                      value={sortOption}
                      onChange={setSortOption}
                    />
                    <div className="badge badge-lg">
                      {searchResults.length > 0
                        ? `${searchResults.length} ${t('results')}`
                        : `${stats.totalListings} ${t('totalListings')}`}
                    </div>
                  </div>
                </div>

                {/* Search Results */}
                {isSearching ? (
                  <div className="flex justify-center py-12">
                    <div className="loading loading-spinner loading-lg"></div>
                  </div>
                ) : searchResults.length > 0 ? (
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {searchResults.map((listing) => (
                      <Link
                        key={listing.id}
                        href={`/${locale}/listing/${listing.id || 0}`}
                        className="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow"
                      >
                        <figure className="aspect-[4/3] relative">
                          {listing.images &&
                          listing.images[0] &&
                          listing.images[0].thumbnail_url ? (
                            <Image
                              src={listing.images[0].thumbnail_url}
                              alt={listing.title || ''}
                              fill
                              className="object-cover"
                            />
                          ) : (
                            <div className="w-full h-full bg-base-200 flex items-center justify-center">
                              <Car className="w-12 h-12 text-base-content/30" />
                            </div>
                          )}
                        </figure>
                        <div className="card-body p-4">
                          <h3 className="card-title text-base line-clamp-1">
                            {listing.title}
                          </h3>
                          {listing.price && (
                            <p className="text-lg font-bold text-primary">
                              €{listing.price.toLocaleString()}
                            </p>
                          )}
                          <p className="text-sm text-base-content/60">
                            {listing.city || listing.country}
                          </p>
                        </div>
                      </Link>
                    ))}
                  </div>
                ) : (
                  <div className="alert">
                    <p>{t('noResultsFound')}</p>
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Latest Listings */}
      {latestListings.length > 0 && (
        <div className="container mx-auto px-4 py-12">
          <h2 className="text-3xl font-bold mb-8">{t('latestListings')}</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
            {latestListings.map((listing) => (
              <Link
                key={listing.id}
                href={`/${locale}/listing/${listing.id}`}
                className="card bg-base-100 shadow-xl hover:shadow-2xl transition-shadow"
              >
                <figure className="aspect-[4/3] relative">
                  {listing.images &&
                  listing.images[0] &&
                  listing.images[0].thumbnail_url ? (
                    <Image
                      src={listing.images[0].thumbnail_url}
                      alt={listing.title || ''}
                      fill
                      className="object-cover"
                    />
                  ) : (
                    <div className="w-full h-full bg-base-200 flex items-center justify-center">
                      <Car className="w-12 h-12 text-base-content/30" />
                    </div>
                  )}
                </figure>
                <div className="card-body p-4">
                  <h3 className="card-title text-base line-clamp-1">
                    {listing.title}
                  </h3>
                  {listing.price && (
                    <p className="text-lg font-bold text-primary">
                      €{listing.price.toLocaleString()}
                    </p>
                  )}
                  <p className="text-sm text-base-content/60">
                    {listing.city || listing.country}
                  </p>
                </div>
              </Link>
            ))}
          </div>

          <div className="text-center mt-8">
            <Link
              href={`/${locale}/search?categories=10101,10102,10103,10104&context=automotive`}
              className="btn btn-primary btn-lg"
            >
              {t('viewAllListings')}
            </Link>
          </div>
        </div>
      )}

      {/* Quick Links */}
      <div className="bg-base-200 py-12">
        <div className="container mx-auto px-4">
          <h2 className="text-3xl font-bold mb-8">{t('quickLinks')}</h2>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">{t('sellYourCar')}</h3>
                <p>{t('sellYourCarDescription')}</p>
                <div className="card-actions justify-end">
                  <Link
                    href={`/${locale}/create-listing-choice`}
                    className="btn btn-primary"
                  >
                    {t('createListing')}
                  </Link>
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">{t('priceAnalysis')}</h3>
                <p>{t('priceAnalysisDescription')}</p>
                <div className="card-actions justify-end">
                  <button className="btn btn-secondary" disabled>
                    {t('comingSoon')}
                  </button>
                </div>
              </div>
            </div>

            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title">{t('vinDecoder')}</h3>
                <p>{t('vinDecoderDescription')}</p>
                <div className="card-actions justify-end">
                  <Link
                    href={`/${locale}/cars/vin-decoder`}
                    className="btn btn-accent"
                  >
                    {t('decode')}
                  </Link>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
