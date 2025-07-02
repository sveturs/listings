'use client';

import { useTranslations } from 'next-intl';
import { useCreateListing } from '@/contexts/CreateListingContext';
import { formatDistanceToNow } from 'date-fns';
import { sr, enUS } from 'date-fns/locale';
import { useLocale } from 'next-intl';

export function DraftStatus() {
  const t = useTranslations('create_listing.draft');
  const locale = useLocale();
  const { isSavingDraft, hasUnsavedChanges, lastSavedAt } = useCreateListing();

  // Определяем локаль для date-fns
  const dateLocale = locale === 'sr' ? sr : enUS;

  // Показываем статус сохранения
  if (isSavingDraft) {
    return (
      <div className="flex items-center gap-2 text-sm">
        <span className="loading loading-spinner loading-xs"></span>
        <span className="text-base-content/70">{t('saving')}</span>
      </div>
    );
  }

  // Показываем время последнего сохранения
  if (lastSavedAt && !hasUnsavedChanges) {
    const timeAgo = formatDistanceToNow(lastSavedAt, {
      addSuffix: true,
      locale: dateLocale,
    });

    return (
      <div className="flex items-center gap-2 text-sm">
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-4 w-4 text-success"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M5 13l4 4L19 7"
          />
        </svg>
        <span className="text-base-content/70">
          {t('savedAt', { time: timeAgo })}
        </span>
      </div>
    );
  }

  // Показываем индикатор несохраненных изменений
  if (hasUnsavedChanges) {
    return (
      <div className="flex items-center gap-2 text-sm">
        <div className="w-2 h-2 bg-warning rounded-full animate-pulse"></div>
        <span className="text-base-content/70">{t('unsavedChanges')}</span>
      </div>
    );
  }

  return null;
}

interface DraftIndicatorProps {
  onClick?: () => void;
}

export function DraftIndicator({ onClick }: DraftIndicatorProps) {
  const t = useTranslations('create_listing');
  const { hasUnsavedChanges, isSavingDraft } = useCreateListing();

  return (
    <div className="indicator">
      {hasUnsavedChanges && !isSavingDraft && (
        <span className="indicator-item badge badge-warning badge-xs"></span>
      )}
      <button
        onClick={onClick}
        className="btn btn-ghost btn-sm"
        title={t('draft.viewDrafts')}
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          className="h-5 w-5"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            strokeLinecap="round"
            strokeLinejoin="round"
            strokeWidth={2}
            d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
          />
        </svg>
      </button>
    </div>
  );
}

export function OfflineIndicator() {
  const t = useTranslations('create_listing.draft');
  const isOnline =
    typeof window !== 'undefined' ? window.navigator.onLine : true;

  if (isOnline) return null;

  return (
    <div className="alert alert-warning">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        className="h-6 w-6 shrink-0 stroke-current"
        fill="none"
        viewBox="0 0 24 24"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          strokeWidth={2}
          d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
        />
      </svg>
      <span>{t('offlineMode')}</span>
    </div>
  );
}
