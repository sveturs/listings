'use client';

import { useTranslations } from 'next-intl';
import { CategoryProposal } from '@/types/categoryProposals';
import { formatDistanceToNow } from 'date-fns';
import { enUS, ru, sr } from 'date-fns/locale';
import { useLocale } from 'next-intl';

interface CategoryProposalCardProps {
  proposal: CategoryProposal;
  onApprove: (id: number) => void;
  onReject: (id: number) => void;
  onView?: (id: number) => void;
  isApproving?: boolean;
  isRejecting?: boolean;
}

const getDateFnsLocale = (locale: string) => {
  switch (locale) {
    case 'ru':
      return ru;
    case 'sr':
      return sr;
    default:
      return enUS;
  }
};

export default function CategoryProposalCard({
  proposal,
  onApprove,
  onReject,
  onView,
  isApproving,
  isRejecting,
}: CategoryProposalCardProps) {
  const t = useTranslations('admin.categoryProposals');
  const locale = useLocale();

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800 dark:bg-yellow-900/30 dark:text-yellow-400';
      case 'approved':
        return 'bg-green-100 text-green-800 dark:bg-green-900/30 dark:text-green-400';
      case 'rejected':
        return 'bg-red-100 text-red-800 dark:bg-red-900/30 dark:text-red-400';
      default:
        return 'bg-gray-100 text-gray-800 dark:bg-gray-800 dark:text-gray-300';
    }
  };

  const formatDate = (dateString: string) => {
    try {
      return formatDistanceToNow(new Date(dateString), {
        addSuffix: true,
        locale: getDateFnsLocale(locale),
      });
    } catch {
      return dateString;
    }
  };

  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700 p-6 hover:shadow-md transition-shadow">
      {/* Header */}
      <div className="flex items-start justify-between mb-4">
        <div className="flex-1">
          <h3 className="text-lg font-semibold text-gray-900 dark:text-white mb-1">
            {proposal.name}
          </h3>
          {proposal.external_category_source && (
            <p className="text-sm text-gray-500 dark:text-gray-400">
              <span className="font-medium">{t('externalSource')}:</span>{' '}
              {proposal.external_category_source}
            </p>
          )}
        </div>
        <span
          className={`px-3 py-1 rounded-full text-xs font-medium ${getStatusColor(
            proposal.status
          )}`}
        >
          {t(`status.${proposal.status}`)}
        </span>
      </div>

      {/* Description */}
      {proposal.description && (
        <p className="text-sm text-gray-600 dark:text-gray-300 mb-4">
          {proposal.description}
        </p>
      )}

      {/* AI Reasoning */}
      {proposal.reasoning && (
        <div className="mb-4 p-3 bg-blue-50 dark:bg-blue-900/20 rounded-lg">
          <p className="text-xs font-semibold text-blue-900 dark:text-blue-400 mb-1">
            {t('reasoning')}
          </p>
          <p className="text-sm text-blue-800 dark:text-blue-300">
            {proposal.reasoning}
          </p>
        </div>
      )}

      {/* Translations */}
      <div className="mb-4">
        <p className="text-xs font-semibold text-gray-700 dark:text-gray-300 mb-2">
          {t('translations')}
        </p>
        <div className="flex flex-wrap gap-2">
          {Object.entries(proposal.name_translations).map(([lang, name]) => {
            if (!name) return null;
            return (
              <span
                key={lang}
                className="px-2 py-1 bg-gray-100 dark:bg-gray-700 rounded text-xs text-gray-700 dark:text-gray-300"
              >
                <span className="font-semibold uppercase">{lang}:</span> {name}
              </span>
            );
          })}
        </div>
      </div>

      {/* Meta info */}
      <div className="flex flex-wrap gap-4 mb-4 text-sm text-gray-600 dark:text-gray-400">
        {proposal.expected_products > 0 && (
          <div className="flex items-center gap-1">
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
                d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4"
              />
            </svg>
            <span>{t('products', { count: proposal.expected_products })}</span>
          </div>
        )}
        <div className="flex items-center gap-1">
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
              d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
          <span>{formatDate(proposal.created_at)}</span>
        </div>
      </div>

      {/* Tags */}
      {proposal.tags && proposal.tags.length > 0 && (
        <div className="mb-4">
          <div className="flex flex-wrap gap-2">
            {proposal.tags.map((tag, index) => (
              <span
                key={index}
                className="px-2 py-1 bg-purple-100 dark:bg-purple-900/30 text-purple-700 dark:text-purple-400 rounded-full text-xs"
              >
                #{tag}
              </span>
            ))}
          </div>
        </div>
      )}

      {/* Actions */}
      {proposal.status === 'pending' && (
        <div className="flex gap-2 pt-4 border-t border-gray-200 dark:border-gray-700">
          <button
            onClick={() => onApprove(proposal.id)}
            disabled={isApproving || isRejecting}
            className="flex-1 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm font-medium"
          >
            {isApproving ? t('approving') : t('approve')}
          </button>
          <button
            onClick={() => onReject(proposal.id)}
            disabled={isApproving || isRejecting}
            className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors text-sm font-medium"
          >
            {isRejecting ? t('rejecting') : t('reject')}
          </button>
          {onView && (
            <button
              onClick={() => onView(proposal.id)}
              className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors text-sm font-medium"
            >
              {t('view')}
            </button>
          )}
        </div>
      )}

      {/* Reviewed info for approved/rejected */}
      {(proposal.status === 'approved' || proposal.status === 'rejected') &&
        proposal.reviewed_at && (
          <div className="pt-4 border-t border-gray-200 dark:border-gray-700">
            <p className="text-xs text-gray-500 dark:text-gray-400">
              {t('reviewedAt')}: {formatDate(proposal.reviewed_at)}
              {proposal.reviewed_by_user_id && (
                <span className="ml-2">
                  {t('reviewedBy')}: #{proposal.reviewed_by_user_id}
                </span>
              )}
            </p>
          </div>
        )}
    </div>
  );
}
