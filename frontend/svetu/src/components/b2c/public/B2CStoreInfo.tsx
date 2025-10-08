'use client';

import { useTranslations, useLocale } from 'next-intl';
import { formatDistanceToNow } from 'date-fns';
import { ru, enUS } from 'date-fns/locale';
import type { B2CStore } from '@/types/b2c';

interface StorefrontInfoProps {
  storefront: B2CStore;
}

export default function StorefrontInfo({ storefront }: StorefrontInfoProps) {
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const locale = useLocale();
  const dateLocale = locale === 'ru' ? ru : enUS;

  const formatHours = (hours: any[]) => {
    if (!hours || !Array.isArray(hours) || hours.length === 0) return null;
    
    const today = new Date().toLocaleDateString('en-US', { weekday: 'short' }).toLowerCase().slice(0, 3);
    const todayHours = hours.find(h => {
      return h && typeof h.day_of_week === 'string' && h.day_of_week.toLowerCase().startsWith(today);
    });
    
    if (!todayHours) return null;
    
    if (todayHours.is_closed) {
      return <span className="text-error">{t('closed')}</span>;
    }
    
    return (
      <span className="text-success">
        {t('open')} • {todayHours.open_time || ''} - {todayHours.close_time || ''}
      </span>
    );
  };

  return (
    <div className="card bg-base-200 shadow-xl">
      <div className="card-body">
        <h3 className="card-title">{t('storeInfo')}</h3>
        
        <div className="space-y-4">
          {/* Business Hours */}
          {storefront.hours && storefront.hours.length > 0 && (
            <div>
              <div className="flex items-center gap-2 mb-2">
                <svg className="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span className="font-semibold">{t('businessHours')}</span>
              </div>
              <div className="text-sm">{formatHours(storefront.hours)}</div>
            </div>
          )}

          {/* Location */}
          {storefront.location && (
            <div>
              <div className="flex items-center gap-2 mb-2">
                <svg className="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                <span className="font-semibold">{tCommon('location')}</span>
              </div>
              <p className="text-sm text-base-content/80">
                {/* TODO: Добавить поддержку локализации когда backend будет поддерживать переводы для витрин */}
                {storefront.location.full_address}
              </p>
            </div>
          )}

          {/* Phone */}
          {storefront.phone && (
            <div>
              <div className="flex items-center gap-2 mb-2">
                <svg className="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                </svg>
                <span className="font-semibold">{tCommon('phone')}</span>
              </div>
              <a href={`tel:${storefront.phone}`} className="text-sm link link-primary">
                {storefront.phone}
              </a>
            </div>
          )}

          {/* Email */}
          {storefront.email && (
            <div>
              <div className="flex items-center gap-2 mb-2">
                <svg className="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                <span className="font-semibold">{tCommon('email')}</span>
              </div>
              <a href={`mailto:${storefront.email}`} className="text-sm link link-primary">
                {storefront.email}
              </a>
            </div>
          )}

          {/* Website */}
          {storefront.website_url && (
            <div>
              <div className="flex items-center gap-2 mb-2">
                <svg className="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                </svg>
                <span className="font-semibold">{t('website')}</span>
              </div>
              <a 
                href={storefront.website_url} 
                target="_blank" 
                rel="noopener noreferrer" 
                className="text-sm link link-primary"
              >
                {storefront.website_url}
              </a>
            </div>
          )}

          {/* Social Links */}
          {(storefront.social_links?.facebook || storefront.social_links?.instagram) && (
            <div>
              <div className="flex items-center gap-2 mb-2">
                <svg className="w-5 h-5 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
                <span className="font-semibold">{t('socialMedia')}</span>
              </div>
              <div className="flex gap-2">
                {storefront.social_links.facebook && (
                  <a 
                    href={storefront.social_links.facebook} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className="btn btn-circle btn-sm"
                  >
                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
                    </svg>
                  </a>
                )}
                {storefront.social_links.instagram && (
                  <a 
                    href={storefront.social_links.instagram} 
                    target="_blank" 
                    rel="noopener noreferrer"
                    className="btn btn-circle btn-sm"
                  >
                    <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 2.163c3.204 0 3.584.012 4.85.07 3.252.148 4.771 1.691 4.919 4.919.058 1.265.069 1.645.069 4.849 0 3.205-.012 3.584-.069 4.849-.149 3.225-1.664 4.771-4.919 4.919-1.266.058-1.644.07-4.85.07-3.204 0-3.584-.012-4.849-.07-3.26-.149-4.771-1.699-4.919-4.92-.058-1.265-.07-1.644-.07-4.849 0-3.204.013-3.583.07-4.849.149-3.227 1.664-4.771 4.919-4.919 1.266-.057 1.645-.069 4.849-.069zm0-2.163c-3.259 0-3.667.014-4.947.072-4.358.2-6.78 2.618-6.98 6.98-.059 1.281-.073 1.689-.073 4.948 0 3.259.014 3.668.072 4.948.2 4.358 2.618 6.78 6.98 6.98 1.281.058 1.689.072 4.948.072 3.259 0 3.668-.014 4.948-.072 4.354-.2 6.782-2.618 6.979-6.98.059-1.28.073-1.689.073-4.948 0-3.259-.014-3.667-.072-4.947-.196-4.354-2.617-6.78-6.979-6.98-1.281-.059-1.69-.073-4.949-.073zM5.838 12a6.162 6.162 0 1112.324 0 6.162 6.162 0 01-12.324 0zM12 16a4 4 0 110-8 4 4 0 010 8zm4.965-10.405a1.44 1.44 0 112.881.001 1.44 1.44 0 01-2.881-.001z"/>
                    </svg>
                  </a>
                )}
              </div>
            </div>
          )}

          {/* Member Since */}
          {storefront.created_at && (
            <div className="pt-4 border-t border-base-300">
              <p className="text-sm text-base-content/60">
                {t('memberSince')} {' '}
                {formatDistanceToNow(new Date(storefront.created_at), {
                  addSuffix: false,
                  locale: dateLocale
                })}
              </p>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}