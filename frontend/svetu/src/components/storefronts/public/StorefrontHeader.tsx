'use client';

import Image from 'next/image';
import { useTranslations } from 'next-intl';
import type { Storefront } from '@/types/storefront';

interface StorefrontHeaderProps {
  storefront: Storefront;
  isOwner: boolean;
  onImageClick: (images: string[], index: number) => void;
}

export default function StorefrontHeader({ 
  storefront, 
  isOwner: _isOwner,
  onImageClick 
}: StorefrontHeaderProps) {
  const t = useTranslations();
  
  const bannerImage = storefront.banner_image_url || '/storefront-banner-default.jpg';
  const logoImage = storefront.logo_url || '/storefront-logo-default.png';
  
  const handleBannerClick = () => {
    if (storefront.banner_image_url) {
      onImageClick([storefront.banner_image_url], 0);
    }
  };
  
  const handleLogoClick = () => {
    if (storefront.logo_url) {
      onImageClick([storefront.logo_url], 0);
    }
  };

  return (
    <div className="relative">
      {/* Banner */}
      <div 
        className="relative h-80 bg-gradient-to-br from-primary/20 to-secondary/20 cursor-pointer group"
        onClick={handleBannerClick}
      >
        <Image
          src={bannerImage}
          alt={storefront.name || 'Storefront banner'}
          fill
          className="object-cover group-hover:scale-105 transition-transform duration-300"
          priority
        />
        <div className="absolute inset-0 bg-gradient-to-t from-black/60 via-transparent to-transparent" />
        
        {/* Verified Badge */}
        {storefront.is_verified && (
          <div className="absolute top-4 right-4 badge badge-success gap-2 backdrop-blur-sm">
            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {t('storefronts.verified')}
          </div>
        )}
      </div>

      {/* Store Info Bar */}
      <div className="bg-base-100 shadow-lg">
        <div className="container mx-auto px-4">
          <div className="flex flex-col md:flex-row items-start md:items-center gap-6 py-6">
            {/* Logo */}
            <div 
              className="relative -mt-16 md:-mt-12 cursor-pointer group"
              onClick={handleLogoClick}
            >
              <div className="avatar">
                <div className="w-32 h-32 rounded-xl ring-4 ring-base-100 overflow-hidden bg-base-200">
                  <Image
                    src={logoImage}
                    alt={`${storefront.name || 'Store'} logo`}
                    width={128}
                    height={128}
                    className="object-cover group-hover:scale-110 transition-transform duration-300"
                  />
                </div>
              </div>
            </div>

            {/* Name and Category */}
            <div className="flex-1">
              <h1 className="text-3xl font-bold mb-1">{storefront.name || 'Storefront'}</h1>
              <p className="text-base-content/60">
                {storefront.business_type === 'retail' && t('storefronts.business_types.retail')}
                {storefront.business_type === 'service' && t('storefronts.business_types.service')}
                {storefront.business_type === 'restaurant' && t('storefronts.business_types.restaurant')}
                {storefront.business_type === 'other' && t('storefronts.business_types.other')}
              </p>
            </div>

            {/* Quick Stats */}
            <div className="flex gap-6 text-center">
              <div>
                <div className="text-2xl font-bold text-primary">
                  {storefront.stats?.average_rating?.toFixed(1) || '0.0'}
                </div>
                <div className="text-sm text-base-content/60">{t('common.rating')}</div>
              </div>
              <div className="divider divider-horizontal mx-0"></div>
              <div>
                <div className="text-2xl font-bold">
                  {storefront.stats?.total_products || 0}
                </div>
                <div className="text-sm text-base-content/60">{t('storefronts.products.title')}</div>
              </div>
              <div className="divider divider-horizontal mx-0"></div>
              <div>
                <div className="text-2xl font-bold">
                  {storefront.stats?.total_reviews || 0}
                </div>
                <div className="text-sm text-base-content/60">{t('common.reviews')}</div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}