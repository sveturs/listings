'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import api from '@/services/api';
import type { components } from '@/types/generated/api';
import { Car, TrendingUp, Calendar, Filter } from 'lucide-react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState, AppDispatch } from '@/store';
import { addToCompare, removeFromCompare } from '@/store/slices/compareSlice';
import CarSortingOptions, {
  type CarSortOption,
} from '@/components/marketplace/CarSortingOptions';
import CarQuickFilters from '@/components/marketplace/CarQuickFilters';
import { CarFilters } from '@/components/marketplace/CarFilters';
import CarBrandIcon from '@/components/marketplace/CarBrandIcon';
import AutocompleteSearch from '@/components/cars/AutocompleteSearch';
import Breadcrumbs from '@/components/cars/Breadcrumbs';
import { CarListingCardEnhanced } from '@/components/cars/CarListingCardEnhanced';
import CarFiltersDrawer from '@/components/cars/CarFiltersDrawer';
import CarQuickViewModal from '@/components/cars/CarQuickViewModal';
import ComparisonBar from '@/components/cars/ComparisonBar';

type CarMake = components['schemas']['backend_internal_domain_models.CarMake'];
type MarketplaceListing =
  components['schemas']['backend_internal_domain_models.MarketplaceListing'];

interface CarsPageClientProps {
  locale: string;
}

export default function CarsPageClient({ locale }: CarsPageClientProps) {
  const t = useTranslations('cars');
  const router = useRouter();
  const dispatch = useDispatch<AppDispatch>();
  const compareItems = useSelector((state: RootState) => state.compare.items);
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
  const [quickViewListing, setQuickViewListing] =
    useState<MarketplaceListing | null>(null);
  const [favoriteIds, setFavoriteIds] = useState<Set<number>>(new Set());

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      setLoading(true);

      // Загружаем статистику автомобилей
      const statsResponse = await api.get('/api/v1/marketplace/cars/stats');
      if (statsResponse.data?.data) {
        setStats({
          totalListings: statsResponse.data.data.totalListings || 0,
          totalMakes: statsResponse.data.data.totalMakes || 0,
          totalModels: statsResponse.data.data.totalModels || 0,
        });
      }

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
        category_ids: '1003,1301,1303', // Автомобильные категории: Automobili, Lični automobili, Auto delovi
        sort: 'created_at_desc',
      };

      const listingsResponse = await api.post(
        '/api/v1/marketplace/search',
        searchParams
      );

      // API возвращает массив напрямую в data, не в data.items
      if (listingsResponse.data?.data) {
        const listings = Array.isArray(listingsResponse.data.data)
          ? listingsResponse.data.data
          : listingsResponse.data.data.items || [];
        setLatestListings(listings);
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
        `/${locale}/search?q=${encodeURIComponent(searchQuery)}&category=1301&context=automotive`
      );
    }
  };

  const handleMakeClick = (makeSlug: string) => {
    router.push(
      `/${locale}/search?category=1301&car_make=${makeSlug}&context=automotive`
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
        category_ids: '1003,1301,1303',
        sort: sortOption,
      };

      // Add search query if present
      if (searchQuery) {
        searchParams.q = searchQuery;
      }

      // Add filters
      Object.entries(activeFilters).forEach(([key, value]) => {
        if (key === 'price_min') {
          searchParams.price_min = value;
        } else if (key === 'price_max') {
          searchParams.price_max = value;
        } else if (key === 'car_mileage_max') {
          searchParams['attributes.mileage_max'] = value;
        } else if (key === 'car_fuel_type') {
          searchParams['attributes.fuel_type'] = value;
        } else if (key === 'car_transmission') {
          searchParams['attributes.transmission'] = value;
        } else if (key === 'condition') {
          searchParams.condition = value;
        } else if (key === 'car_body_type') {
          searchParams['attributes.body_type'] = Array.isArray(value)
            ? value.join(',')
            : value;
        } else if (key === 'car_make') {
          searchParams['attributes.car_make'] = value;
        } else if (key === 'car_model') {
          searchParams['attributes.car_model'] = value;
        } else if (key === 'car_year_from') {
          searchParams['attributes.year_min'] = value;
        } else if (key === 'car_year_to') {
          searchParams['attributes.year_max'] = value;
        }
      });

      const response = await api.post(
        '/api/v1/marketplace/search',
        searchParams
      );

      // API возвращает массив напрямую в data, не в data.items
      if (response.data?.data) {
        const results = Array.isArray(response.data.data)
          ? response.data.data
          : response.data.data.items || [];
        setSearchResults(results);
      }
    } catch (error) {
      console.error('Error searching cars:', error);
    } finally {
      setIsSearching(false);
    }
  };

  // Trigger search when filters or sort changes, or on initial load
  useEffect(() => {
    // Всегда выполняем поиск - либо с фильтрами, либо без них для показа всех автомобилей
    if (!loading) {
      searchCars();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [activeFilters, sortOption, loading]);

  // Handle favorite toggle
  const handleFavorite = (listingId: number) => {
    setFavoriteIds((prev) => {
      const newSet = new Set(prev);
      if (newSet.has(listingId)) {
        newSet.delete(listingId);
      } else {
        newSet.add(listingId);
      }
      return newSet;
    });
  };

  // Handle compare toggle
  const handleCompare = (listing: MarketplaceListing) => {
    const isComparing = compareItems.some((item) => item.id === listing.id);

    if (isComparing) {
      dispatch(removeFromCompare(listing.id || 0));
    } else {
      // Get car attributes
      const carAttrs: any = listing.attributes || {};

      dispatch(
        addToCompare({
          id: listing.id || 0,
          title: listing.title || '',
          price: listing.price || 0,
          year: carAttrs.year || new Date().getFullYear(),
          make: carAttrs.make || '',
          model: carAttrs.model || '',
          mileage: carAttrs.mileage,
          fuelType: carAttrs.fuel_type,
          transmission: carAttrs.transmission,
          engineSize: carAttrs.engine_size,
          power: carAttrs.power,
          bodyType: carAttrs.body_type,
          color: carAttrs.color,
          location: listing.city || listing.country,
          imageUrl: listing.images?.[0]?.thumbnail_url,
          vin: carAttrs.vin,
          driveType: carAttrs.drive_type,
          doors: carAttrs.doors,
          seats: carAttrs.seats,
          condition: carAttrs.condition,
          previousOwners: carAttrs.previous_owners,
          warranty: carAttrs.warranty,
          firstRegistration: carAttrs.first_registration,
          technicalInspection: carAttrs.technical_inspection,
          features: carAttrs.features,
        })
      );
    }
  };

  // Handle share
  const handleShare = (listingId: number) => {
    if (navigator.share) {
      navigator.share({
        title: t('shareTitle'),
        url: `${window.location.origin}/${locale}/listing/${listingId}`,
      });
    } else {
      // Fallback: copy to clipboard
      navigator.clipboard.writeText(
        `${window.location.origin}/${locale}/listing/${listingId}`
      );
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-base-100 flex items-center justify-center">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    );
  }

  // Prepare active filters for Breadcrumbs
  const breadcrumbFilters = Object.entries(activeFilters)
    .filter(([_, value]) => value && value !== '')
    .map(([key, value]) => ({
      key,
      value: String(value),
      label: t(`filters.${key}`),
    }));

  const handleRemoveFilter = (filterKey: string) => {
    const newFilters = { ...activeFilters };
    delete newFilters[filterKey];
    setActiveFilters(newFilters);
  };

  const handleAutoSearch = (query: string, filters?: any) => {
    if (filters) {
      setActiveFilters((prev) => ({ ...prev, ...filters }));
    } else if (query) {
      setSearchQuery(query);
      handleSearch();
    }
  };

  return (
    <div className="min-h-screen bg-base-100">
      {/* Breadcrumbs */}
      <div className="container mx-auto px-4 py-4">
        <Breadcrumbs
          activeFilters={breadcrumbFilters}
          onRemoveFilter={handleRemoveFilter}
        />
      </div>

      {/* Hero Section */}
      <div className="hero min-h-[400px] bg-gradient-to-br from-primary to-primary-focus text-primary-content">
        <div className="hero-content text-center">
          <div className="max-w-lg w-full">
            <h1 className="text-5xl font-bold mb-5 flex items-center justify-center gap-3">
              <Car className="w-12 h-12" />
              {t('heroTitle')}
            </h1>
            <p className="mb-8">{t('heroDescription')}</p>

            {/* Advanced Search Bar */}
            <AutocompleteSearch
              onSearch={handleAutoSearch}
              placeholder={t('searchPlaceholder')}
              className="w-full"
            />
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
            href={`/${locale}/search?category=1301&context=automotive`}
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
              href={`/${locale}/search?category=1301&context=automotive`}
              className="btn btn-outline btn-lg"
            >
              {t('categories.passenger')}
            </Link>
            <Link
              href={`/${locale}/search?category=10174&context=automotive`}
              className="btn btn-outline btn-lg"
            >
              {t('categories.suv')}
            </Link>
            <Link
              href={`/${locale}/search?category=1303&context=automotive`}
              className="btn btn-outline btn-lg"
            >
              {t('categories.commercial')}
            </Link>
            <Link
              href={`/${locale}/search?category=1302&context=automotive`}
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
                        ? `${
                            searchResults.filter((listing) =>
                              [1003, 1301, 1303].includes(
                                listing.category_id || 0
                              )
                            ).length
                          } ${t('results')}`
                        : `${stats.totalListings} ${t('totalListings')}`}
                    </div>
                  </div>
                </div>

                {/* Search Results */}
                {isSearching ? (
                  <div className="flex justify-center py-12">
                    <div className="loading loading-spinner loading-lg"></div>
                  </div>
                ) : searchResults.filter((listing) =>
                    [1003, 1301, 1303].includes(listing.category_id || 0)
                  ).length > 0 ? (
                  <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                    {searchResults
                      .filter((listing) =>
                        // Показываем только автомобили (категории 1003, 1301, 1303)
                        [1003, 1301, 1303].includes(listing.category_id || 0)
                      )
                      .map((listing) => (
                        <CarListingCardEnhanced
                          key={listing.id}
                          listing={listing}
                          locale={locale}
                          onFavorite={() => handleFavorite(listing.id || 0)}
                          onShare={() => handleShare(listing.id || 0)}
                          onCompare={() => handleCompare(listing)}
                          onQuickView={setQuickViewListing}
                          isFavorited={favoriteIds.has(listing.id || 0)}
                          isComparing={compareItems.some(
                            (item) => item.id === listing.id
                          )}
                        />
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
              <CarListingCardEnhanced
                key={listing.id}
                listing={listing}
                locale={locale}
                onFavorite={() => handleFavorite(listing.id || 0)}
                onShare={() => handleShare(listing.id || 0)}
                onCompare={() => handleCompare(listing)}
                onQuickView={setQuickViewListing}
                isFavorited={favoriteIds.has(listing.id || 0)}
                isComparing={compareItems.some(
                  (item) => item.id === listing.id
                )}
              />
            ))}
          </div>

          <div className="text-center mt-8">
            <Link
              href={`/${locale}/search?categories=1003,1301,1303&context=automotive`}
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

      {/* Quick View Modal */}
      {quickViewListing && (
        <CarQuickViewModal
          isOpen={true}
          listing={quickViewListing}
          locale={locale}
          onClose={() => setQuickViewListing(null)}
          onPrevious={() => {
            const currentIndex = searchResults.findIndex(
              (l) => l.id === quickViewListing.id
            );
            if (currentIndex > 0) {
              setQuickViewListing(searchResults[currentIndex - 1]);
            }
          }}
          onNext={() => {
            const currentIndex = searchResults.findIndex(
              (l) => l.id === quickViewListing.id
            );
            if (currentIndex < searchResults.length - 1) {
              setQuickViewListing(searchResults[currentIndex + 1]);
            }
          }}
          hasPrevious={
            searchResults.findIndex((l) => l.id === quickViewListing.id) > 0
          }
          hasNext={
            searchResults.findIndex((l) => l.id === quickViewListing.id) <
            searchResults.length - 1
          }
        />
      )}

      {/* Mobile Filters Drawer */}
      <CarFiltersDrawer
        isOpen={showFilters}
        onClose={() => setShowFilters(false)}
        onApply={(filters) => {
          setActiveFilters(filters);
          setShowFilters(false);
        }}
        onReset={() => setActiveFilters({})}
        currentFilters={activeFilters}
        activeFilterCount={Object.keys(activeFilters).length}
      />

      {/* Comparison Bar */}
      <ComparisonBar />
    </div>
  );
}
