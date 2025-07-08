'use client';

import { useTranslations } from 'next-intl';

export default function SearchAnalytics() {
  const t = useTranslations('admin.search.analytics');

  return (
    <div className="space-y-6">
      <div className="card bg-base-100 shadow-md">
        <div className="card-body">
          <h3 className="card-title">{t('title')}</h3>
          <p className="text-base-content/60">{t('description')}</p>
          <div className="alert alert-info mt-4">
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
            <span>Компонент общей аналитики поиска будет реализован позже</span>
          </div>
        </div>
      </div>
    </div>
  );
}
