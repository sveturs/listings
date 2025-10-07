'use client';

import { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import {
  fetchCategoryProposals,
  fetchPendingCount,
  approveCategoryProposal,
  rejectCategoryProposal,
  setStatusFilter,
  setPage,
  openApproveModal,
  closeApproveModal,
  openRejectModal,
  closeRejectModal,
} from '@/store/slices/categoryProposalsSlice';
import CategoryProposalCard from './CategoryProposalCard';

export default function CategoryProposalsList() {
  const t = useTranslations('admin.categoryProposals');
  const dispatch = useAppDispatch();

  const {
    proposals,
    total,
    page,
    pageSize,
    totalPages,
    pendingCount,
    statusFilter,
    isLoading,
    isApproving,
    isRejecting,
    error,
    isApproveModalOpen,
    isRejectModalOpen,
    selectedProposalId,
  } = useAppSelector((state) => state.categoryProposals);

  const [rejectReason, setRejectReason] = useState('');
  const [createCategory, setCreateCategory] = useState(false);

  // Fetch data on mount and when filters change
  useEffect(() => {
    const filters = {
      page,
      page_size: pageSize,
      ...(statusFilter !== 'all' && { status: statusFilter }),
    };

    dispatch(fetchCategoryProposals(filters));
    dispatch(fetchPendingCount());
  }, [dispatch, page, pageSize, statusFilter]);

  const handleApprove = (id: number) => {
    dispatch(openApproveModal(id));
  };

  const handleReject = (id: number) => {
    dispatch(openRejectModal(id));
  };

  const confirmApprove = async () => {
    if (!selectedProposalId) return;

    try {
      await dispatch(
        approveCategoryProposal({
          id: selectedProposalId,
          request: { create_category: createCategory },
        })
      ).unwrap();

      // Refresh list
      dispatch(
        fetchCategoryProposals({
          page,
          page_size: pageSize,
          ...(statusFilter !== 'all' && { status: statusFilter }),
        })
      );
      dispatch(fetchPendingCount());
    } catch (err) {
      console.error('Failed to approve proposal:', err);
    }
  };

  const confirmReject = async () => {
    if (!selectedProposalId) return;

    try {
      await dispatch(
        rejectCategoryProposal({
          id: selectedProposalId,
          request: { reason: rejectReason || undefined },
        })
      ).unwrap();

      // Reset reason and refresh list
      setRejectReason('');
      dispatch(
        fetchCategoryProposals({
          page,
          page_size: pageSize,
          ...(statusFilter !== 'all' && { status: statusFilter }),
        })
      );
      dispatch(fetchPendingCount());
    } catch (err) {
      console.error('Failed to reject proposal:', err);
    }
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-2xl font-bold text-gray-900 dark:text-white">
            {t('title')}
          </h1>
          <p className="text-sm text-gray-600 dark:text-gray-400 mt-1">
            {t('description')}
          </p>
        </div>
        {pendingCount > 0 && (
          <div className="flex items-center gap-2 px-4 py-2 bg-yellow-100 dark:bg-yellow-900/30 text-yellow-800 dark:text-yellow-400 rounded-lg">
            <svg
              className="w-5 h-5"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
              />
            </svg>
            <span className="font-semibold">{t('pendingCount', { count: pendingCount })}</span>
          </div>
        )}
      </div>

      {/* Filters */}
      <div className="flex gap-2">
        <button
          onClick={() => dispatch(setStatusFilter('all'))}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            statusFilter === 'all'
              ? 'bg-blue-600 text-white'
              : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
          }`}
        >
          {t('allStatuses')}
        </button>
        <button
          onClick={() => dispatch(setStatusFilter('pending'))}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            statusFilter === 'pending'
              ? 'bg-yellow-600 text-white'
              : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
          }`}
        >
          {t('pending')}
        </button>
        <button
          onClick={() => dispatch(setStatusFilter('approved'))}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            statusFilter === 'approved'
              ? 'bg-green-600 text-white'
              : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
          }`}
        >
          {t('approved')}
        </button>
        <button
          onClick={() => dispatch(setStatusFilter('rejected'))}
          className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
            statusFilter === 'rejected'
              ? 'bg-red-600 text-white'
              : 'bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 hover:bg-gray-300 dark:hover:bg-gray-600'
          }`}
        >
          {t('rejected')}
        </button>
      </div>

      {/* Error message */}
      {error && (
        <div className="p-4 bg-red-100 dark:bg-red-900/30 text-red-800 dark:text-red-400 rounded-lg">
          {error}
        </div>
      )}

      {/* Loading state */}
      {isLoading && (
        <div className="flex justify-center items-center py-12">
          <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
        </div>
      )}

      {/* Proposals list */}
      {!isLoading && proposals.length === 0 && (
        <div className="text-center py-12">
          <svg
            className="w-16 h-16 mx-auto text-gray-400 mb-4"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
            />
          </svg>
          <p className="text-gray-600 dark:text-gray-400">{t('noProposals')}</p>
        </div>
      )}

      {!isLoading && proposals.length > 0 && (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {proposals.map((proposal) => (
            <CategoryProposalCard
              key={proposal.id}
              proposal={proposal}
              onApprove={handleApprove}
              onReject={handleReject}
              isApproving={isApproving}
              isRejecting={isRejecting}
            />
          ))}
        </div>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <div className="flex justify-center items-center gap-2 mt-6">
          <button
            onClick={() => dispatch(setPage(page - 1))}
            disabled={page === 1}
            className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
          >
            Previous
          </button>
          <span className="text-sm text-gray-600 dark:text-gray-400">
            Page {page} of {totalPages} ({total} total)
          </span>
          <button
            onClick={() => dispatch(setPage(page + 1))}
            disabled={page >= totalPages}
            className="px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg disabled:opacity-50 disabled:cursor-not-allowed hover:bg-gray-300 dark:hover:bg-gray-600 transition-colors"
          >
            Next
          </button>
        </div>
      )}

      {/* Approve Modal */}
      {isApproveModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full p-6">
            <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-4">
              {t('approveModal.title')}
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
              {t('approveModal.message')}
            </p>

            <div className="space-y-3 mb-6">
              <label className="flex items-start gap-3 p-3 border-2 border-gray-200 dark:border-gray-700 rounded-lg cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                <input
                  type="radio"
                  checked={createCategory}
                  onChange={() => setCreateCategory(true)}
                  className="mt-1"
                />
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">
                    {t('approveModal.createCategory')}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {t('createCategory')}
                  </p>
                </div>
              </label>

              <label className="flex items-start gap-3 p-3 border-2 border-gray-200 dark:border-gray-700 rounded-lg cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700/50 transition-colors">
                <input
                  type="radio"
                  checked={!createCategory}
                  onChange={() => setCreateCategory(false)}
                  className="mt-1"
                />
                <div>
                  <p className="font-medium text-gray-900 dark:text-white">
                    {t('approveModal.justApprove')}
                  </p>
                  <p className="text-xs text-gray-500 dark:text-gray-400">
                    {t('justMarkApproved')}
                  </p>
                </div>
              </label>
            </div>

            <div className="flex gap-3">
              <button
                onClick={() => {
                  dispatch(closeApproveModal());
                  setCreateCategory(false);
                }}
                disabled={isApproving}
                className="flex-1 px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-50 transition-colors"
              >
                {t('approveModal.cancel')}
              </button>
              <button
                onClick={confirmApprove}
                disabled={isApproving}
                className="flex-1 px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 disabled:opacity-50 transition-colors"
              >
                {isApproving ? t('approving') : t('approveModal.confirm')}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* Reject Modal */}
      {isRejectModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-xl max-w-md w-full p-6">
            <h3 className="text-lg font-bold text-gray-900 dark:text-white mb-4">
              {t('rejectModal.title')}
            </h3>
            <p className="text-sm text-gray-600 dark:text-gray-400 mb-4">
              {t('rejectModal.message')}
            </p>

            <div className="mb-4">
              <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                {t('rejectModal.reasonLabel')}
              </label>
              <textarea
                value={rejectReason}
                onChange={(e) => setRejectReason(e.target.value)}
                placeholder={t('rejectReasonPlaceholder')}
                rows={3}
                className="w-full px-3 py-2 border border-gray-300 dark:border-gray-600 rounded-lg bg-white dark:bg-gray-700 text-gray-900 dark:text-white focus:ring-2 focus:ring-red-500 focus:border-transparent"
              />
            </div>

            <div className="flex gap-3">
              <button
                onClick={() => {
                  dispatch(closeRejectModal());
                  setRejectReason('');
                }}
                disabled={isRejecting}
                className="flex-1 px-4 py-2 bg-gray-200 dark:bg-gray-700 text-gray-700 dark:text-gray-300 rounded-lg hover:bg-gray-300 dark:hover:bg-gray-600 disabled:opacity-50 transition-colors"
              >
                {t('rejectModal.cancel')}
              </button>
              <button
                onClick={confirmReject}
                disabled={isRejecting}
                className="flex-1 px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 disabled:opacity-50 transition-colors"
              >
                {isRejecting ? t('rejecting') : t('rejectModal.confirm')}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
