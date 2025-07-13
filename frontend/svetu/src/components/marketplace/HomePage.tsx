'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import dynamic from 'next/dynamic';
import MarketplaceList from './MarketplaceList';
import { UnifiedSearchItem } from '@/services/unifiedSearch';
import { RadiusSearchResult } from '@/components/GIS/types/gis';

// Динамический импорт карты для избежания SSR проблем
const MarketplaceMapWithRadiusSearch = dynamic(
  () => import('./MarketplaceMapWithRadiusSearch'),
  {
    ssr: false,
    loading: () => (
      <div className="w-full h-96 bg-base-200 rounded-lg flex items-center justify-center mb-8">
        <div className="loading loading-spinner loading-lg"></div>
      </div>
    ),
  }
);

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

  const handleListingSelect = (listing: RadiusSearchResult) => {
    setSelectedListing(listing);
    // Можно добавить логику для показа детальной информации
    console.log('Selected listing:', listing);
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

      {/* Переключатель между списком и картой */}
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-2xl font-bold text-base-content">
          {t('latestListings')}
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
        <div className="mb-8">
          <MarketplaceMapWithRadiusSearch
            className="rounded-lg overflow-hidden border border-base-300"
            onListingSelect={handleListingSelect}
            initialViewState={{
              // Центр на Белград, Сербия
              latitude: 44.8176,
              longitude: 20.4649,
              zoom: 11,
            }}
          />

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
        <MarketplaceList initialData={initialData} locale={locale} />
      )}
    </>
  );
}
