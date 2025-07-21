'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import MarketplaceList from './MarketplaceList';
import CategorySidebar from '@/components/categories/CategorySidebar';
import { UnifiedSearchItem } from '@/services/unifiedSearch';
import { RadiusSearchResult } from '@/components/GIS/types/gis';

// Использовем основную карту вместо дублирующего компонента
import { Link } from '@/i18n/routing';

interface HomePageProps {
  initialData: {
    items: UnifiedSearchItem[];
    total: number;
    page: number;
    limit: number;
    has_more: boolean;
  } | null;
  locale: string;
  error?: Error | null;
  paymentsEnabled?: boolean;
}

export default function HomePage({
  initialData,
  locale,
  error,
  paymentsEnabled = false,
}: HomePageProps) {
  const t = useTranslations('home');
  const [showMap, setShowMap] = useState(false);
  const [selectedListing, setSelectedListing] =
    useState<RadiusSearchResult | null>(null);
  const [productTypes, setProductTypes] = useState<
    ('marketplace' | 'storefront')[]
  >(['marketplace', 'storefront']);
  const [showOnlyMarketplace, setShowOnlyMarketplace] = useState(false);
  const [selectedCategoryId, setSelectedCategoryId] = useState<number | null>(
    null
  );

  const _handleListingSelect = (listing: RadiusSearchResult) => {
    setSelectedListing(listing);
    // Можно добавить логику для показа детальной информации
    console.log('Selected listing:', listing);
  };

  const handleProductTypeFilter = (onlyMarketplace: boolean) => {
    setShowOnlyMarketplace(onlyMarketplace);
    if (onlyMarketplace) {
      setProductTypes(['marketplace']);
    } else {
      setProductTypes(['marketplace', 'storefront']);
    }
  };

  const handleCategorySelect = (categoryId: number | null) => {
    setSelectedCategoryId(categoryId);
  };

  return (
    <>
      {paymentsEnabled && (
        <div className="alert alert-info mb-8 shadow-lg">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            fill="none"
            viewBox="0 0 24 24"
            className="stroke-current shrink-0 w-6 h-6"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            ></path>
          </svg>
          <span>🎉 {t('paymentsNowAvailable')}</span>
        </div>
      )}

      {error && (
        <div className="alert alert-error mb-8 shadow-lg">
          <svg
            xmlns="http://www.w3.org/2000/svg"
            className="stroke-current shrink-0 h-6 w-6"
            fill="none"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth="2"
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{t('errorLoadingData')}</span>
        </div>
      )}

      {/* Основной контент с сайдбаром */}
      <div className="flex gap-6">
        {/* Сайдбар с категориями */}
        <div className="w-80 flex-shrink-0">
          <CategorySidebar
            onCategorySelect={handleCategorySelect}
            selectedCategoryId={selectedCategoryId}
            className="sticky top-4"
          />
        </div>

        {/* Основной контент */}
        <div className="flex-1 min-w-0">
          {/* Фильтр типов товаров */}
          <div className="mb-4">
            <div className="form-control">
              <label className="label cursor-pointer justify-start gap-3">
                <input
                  type="checkbox"
                  className="checkbox checkbox-primary"
                  checked={showOnlyMarketplace}
                  onChange={(e) => handleProductTypeFilter(e.target.checked)}
                />
                <span className="label-text">{t('privateListingsOnly')}</span>
              </label>
              <p className="text-xs text-base-content/60 ml-8">
                {t('privateListingsOnlyDescription')}
              </p>
            </div>
          </div>

          {/* Переключатель между списком и картой */}
          <div className="flex items-center justify-between mb-6">
            <h2 className="text-2xl font-bold text-base-content">
              {t('latestListings')}
              {selectedCategoryId && (
                <span className="text-sm font-normal text-base-content/70 ml-2">
                  (фильтр по категории)
                </span>
              )}
            </h2>

            <div className="flex items-center space-x-2">
              <span className="text-sm text-base-content/70">
                {showMap ? t('mapView') : t('listView')}
              </span>
              <div className="join">
                <button
                  className={`btn btn-sm join-item ${!showMap ? 'btn-primary' : 'btn-ghost'}`}
                  onClick={() => setShowMap(false)}
                  aria-label={t('switchToListView')}
                >
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M4 6h16M4 10h16M4 14h16M4 18h16"
                    />
                  </svg>
                </button>
                <button
                  className={`btn btn-sm join-item ${showMap ? 'btn-primary' : 'btn-ghost'}`}
                  onClick={() => setShowMap(true)}
                  aria-label={t('switchToMapView')}
                >
                  <svg
                    className="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M9 20l-5.447-2.724A1 1 0 013 16.382V5.618a1 1 0 011.447-.894L9 7m0 13l6-3m-6 3V7m6 10l4.553 2.276A1 1 0 0021 18.382V7.618a1 1 0 00-.553-.894L15 4m0 13V4m0 0L9 7"
                    />
                  </svg>
                </button>
              </div>
            </div>
          </div>

          {/* Показываем либо карту, либо список */}
          {showMap ? (
            <div className="mb-8 p-4 bg-base-200 rounded-lg">
              <p className="text-center text-base-content/70">
                {t('map.redirecting')}{' '}
                <Link href="/map" className="link link-primary">
                  {t('map.goToMap')}
                </Link>
              </p>

              {/* Информация о выбранном объявлении */}
              {selectedListing && (
                <div className="mt-4 p-4 bg-base-200 rounded-lg">
                  <h3 className="font-semibold text-lg mb-2">
                    {selectedListing.title}
                  </h3>
                  {selectedListing.description && (
                    <p className="text-base-content/70 mb-2">
                      {selectedListing.description}
                    </p>
                  )}
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-4">
                      {selectedListing.price && (
                        <span className="text-xl font-bold text-primary">
                          {selectedListing.price}{' '}
                          {selectedListing.currency || 'RSD'}
                        </span>
                      )}
                      {selectedListing.category && (
                        <span className="badge badge-outline">
                          {selectedListing.category}
                        </span>
                      )}
                    </div>
                    <span className="text-sm text-base-content/70">
                      {selectedListing.distance.toFixed(1)} км от центра поиска
                    </span>
                  </div>
                </div>
              )}
            </div>
          ) : (
            <MarketplaceList
              initialData={initialData}
              locale={locale}
              productTypes={productTypes}
              selectedCategoryId={selectedCategoryId}
            />
          )}
        </div>
      </div>
    </>
  );
}
