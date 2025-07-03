'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useListingDrafts } from '@/hooks/useListingDraft';
import { formatDistanceToNow } from 'date-fns';
import { sr, enUS } from 'date-fns/locale';
import { useLocale } from 'next-intl';
import { useRouter } from '@/i18n/routing';
import { draftService } from '@/services/draftService';
import { useAuth } from '@/contexts/AuthContext';
import { toast } from '@/utils/toast';

interface DraftsModalProps {
  isOpen: boolean;
  onClose: () => void;
}

export function DraftsModal({ isOpen, onClose }: DraftsModalProps) {
  const t = useTranslations('create_listing.draft');
  const locale = useLocale();
  const router = useRouter();
  const { user } = useAuth();
  const { drafts, isLoading, refreshDrafts } = useListingDrafts();
  const [deletingId, setDeletingId] = useState<string | null>(null);

  // Определяем локаль для date-fns
  const dateLocale = locale === 'sr' ? sr : enUS;

  const handleOpenDraft = (draftId: string) => {
    router.push(`/create-listing?draft=${draftId}`);
    onClose();
  };

  const handleDeleteDraft = async (draftId: string) => {
    if (!user?.id) return;

    setDeletingId(draftId);
    try {
      draftService.deleteDraft(draftId, user.id);
      refreshDrafts();
      toast.success(t('deleted'));
    } catch {
      toast.error(t('deleteError'));
    } finally {
      setDeletingId(null);
    }
  };

  const handleExportDraft = (draftId: string) => {
    if (!user?.id) return;

    const draft = draftService.getDraft(draftId, user.id);
    if (draft) {
      const json = draftService.exportDraft(draft);
      const blob = new Blob([json], { type: 'application/json' });
      const url = URL.createObjectURL(blob);
      const a = document.createElement('a');
      a.href = url;
      a.download = `draft_${draftId}_${new Date().toISOString()}.json`;
      a.click();
      URL.revokeObjectURL(url);
      toast.success(t('exported'));
    }
  };

  if (!isOpen) return null;

  return (
    <div className="modal modal-open">
      <div className="modal-box max-w-4xl">
        <h3 className="font-bold text-lg mb-4">{t('myDrafts')}</h3>

        {isLoading ? (
          <div className="flex justify-center py-8">
            <span className="loading loading-spinner loading-lg"></span>
          </div>
        ) : drafts.length === 0 ? (
          <div className="text-center py-8">
            <svg
              className="mx-auto h-12 w-12 text-base-content/30 mb-4"
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
            <p className="text-base-content/70">{t('noDrafts')}</p>
          </div>
        ) : (
          <div className="space-y-4">
            {drafts.map((draftMeta) => {
              const updatedAgo = formatDistanceToNow(
                new Date(draftMeta.updatedAt),
                { addSuffix: true, locale: dateLocale }
              );
              const expiresIn = formatDistanceToNow(
                new Date(draftMeta.expiresAt),
                { addSuffix: false, locale: dateLocale }
              );

              return (
                <div
                  key={draftMeta.id}
                  className="card bg-base-200 shadow-sm hover:shadow-md transition-shadow"
                >
                  <div className="card-body">
                    <div className="flex justify-between items-start">
                      <div className="flex-1">
                        <h4 className="font-semibold">
                          {draftMeta.title || t('untitled')}
                        </h4>
                        {draftMeta.category && (
                          <p className="text-sm text-base-content/70">
                            {draftMeta.category.name}
                          </p>
                        )}
                        <div className="flex items-center gap-4 mt-2 text-xs text-base-content/50">
                          <span>{t('updatedAgo', { time: updatedAgo })}</span>
                          <span>•</span>
                          <span>{t('expiresIn', { time: expiresIn })}</span>
                          {draftMeta.isComplete && (
                            <>
                              <span>•</span>
                              <span className="text-success">
                                {t('complete')}
                              </span>
                            </>
                          )}
                        </div>
                      </div>

                      <div className="flex gap-2">
                        <button
                          onClick={() => handleOpenDraft(draftMeta.id)}
                          className="btn btn-primary btn-sm"
                        >
                          {t('open')}
                        </button>
                        <div className="dropdown dropdown-end">
                          <div
                            tabIndex={0}
                            role="button"
                            className="btn btn-ghost btn-sm"
                          >
                            <svg
                              xmlns="http://www.w3.org/2000/svg"
                              className="h-4 w-4"
                              fill="none"
                              viewBox="0 0 24 24"
                              stroke="currentColor"
                            >
                              <path
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                strokeWidth={2}
                                d="M12 5v.01M12 12v.01M12 19v.01"
                              />
                            </svg>
                          </div>
                          <ul
                            tabIndex={0}
                            className="dropdown-content menu bg-base-100 rounded-box z-[1] w-52 p-2 shadow"
                          >
                            <li>
                              <button
                                onClick={() => handleExportDraft(draftMeta.id)}
                              >
                                {t('export')}
                              </button>
                            </li>
                            <li>
                              <button
                                onClick={() => handleDeleteDraft(draftMeta.id)}
                                className="text-error"
                                disabled={deletingId === draftMeta.id}
                              >
                                {deletingId === draftMeta.id ? (
                                  <span className="loading loading-spinner loading-xs"></span>
                                ) : (
                                  t('delete')
                                )}
                              </button>
                            </li>
                          </ul>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              );
            })}
          </div>
        )}

        <div className="modal-action">
          <button onClick={onClose} className="btn btn-ghost">
            {t('close')}
          </button>
        </div>
      </div>

      <div className="modal-backdrop" onClick={onClose}>
        <button className="cursor-default">close</button>
      </div>
    </div>
  );
}
