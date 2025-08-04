'use client';

import { useTranslations } from 'next-intl';

export default function UsersPage() {
  const t = useTranslations('admin');

  return (
    <div>
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-3xl font-bold">{t('users.title')}</h1>
      </div>

      <div className="card bg-base-100 shadow-xl">
        <div className="card-body">
          <div className="flex flex-col items-center justify-center py-12">
            <svg
              xmlns="http://www.w3.org/2000/svg"
              className="h-24 w-24 text-base-300 mb-4"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={1}
                d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"
              />
            </svg>
            <p className="text-lg text-base-content/60">{t('users.comingSoon')}</p>
            <p className="text-sm text-base-content/40 mt-2">
              {t('users.description')}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
