'use client';

import React, { useEffect, useState } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { useTranslations } from 'next-intl';
import {
  fetchImportJobs,
  fetchJobStatus,
  setCurrentJob,
  setJobDetailsModalOpen,
  setErrorsModalOpen,
  cancelImportJob,
  retryImportJob,
} from '@/store/slices/importSlice';
import { IMPORT_STATUS_COLORS, IMPORT_STATUS_ICONS } from '@/types/import';
import type { ImportJob } from '@/types/import';

interface ImportJobsListProps {
  storefrontId: number;
  autoRefresh?: boolean;
  refreshInterval?: number;
}

export default function ImportJobsList({
  storefrontId,
  autoRefresh = true,
  refreshInterval = 5000,
}: ImportJobsListProps) {
  const dispatch = useAppDispatch();
  const t = useTranslations('storefronts');
  const { jobs, isLoading } = useAppSelector((state) => state.import);

  const [selectedStatus, setSelectedStatus] = useState<string>('all');

  useEffect(() => {
    dispatch(fetchImportJobs({ storefrontId }));
  }, [dispatch, storefrontId]);

  // Auto-refresh for pending/processing jobs
  useEffect(() => {
    if (!autoRefresh) return;

    const hasActiveJobs = jobs.some(
      (job) => job.status === 'pending' || job.status === 'processing'
    );

    if (!hasActiveJobs) return;

    const interval = setInterval(() => {
      jobs.forEach((job) => {
        if (job.status === 'pending' || job.status === 'processing') {
          dispatch(fetchJobStatus({ storefrontId, jobId: job.id }));
        }
      });
    }, refreshInterval);

    return () => clearInterval(interval);
  }, [jobs, autoRefresh, refreshInterval, dispatch, storefrontId]);

  const handleJobClick = (job: ImportJob) => {
    dispatch(setCurrentJob(job));
    dispatch(setJobDetailsModalOpen(true));
  };

  const handleShowErrors = (job: ImportJob) => {
    dispatch(setCurrentJob(job));
    dispatch(setErrorsModalOpen(true));
  };

  const handleCancelJob = async (jobId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    if (confirm(t('actions.confirmCancel'))) {
      await dispatch(cancelImportJob({ storefrontId, jobId }));
    }
  };

  const handleRetryJob = async (jobId: number, e: React.MouseEvent) => {
    e.stopPropagation();
    await dispatch(retryImportJob({ storefrontId, jobId }));
  };

  const filteredJobs =
    selectedStatus === 'all'
      ? jobs
      : jobs.filter((job) => job.status === selectedStatus);

  const getStatusIcon = (status: string) => {
    const iconName =
      IMPORT_STATUS_ICONS[status as keyof typeof IMPORT_STATUS_ICONS] ||
      'help-circle';

    switch (iconName) {
      case 'clock':
        return (
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
        );
      case 'refresh':
        return (
          <svg
            className="w-4 h-4 animate-spin"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              strokeLinecap="round"
              strokeLinejoin="round"
              strokeWidth={2}
              d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
            />
          </svg>
        );
      case 'check-circle':
        return (
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
              d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        );
      case 'x-circle':
        return (
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
              d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        );
      default:
        return (
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
              d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        );
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleString();
  };

  const calculateProgress = (job: ImportJob) => {
    if (job.total_records === 0) return 0;
    return Math.round((job.processed_records / job.total_records) * 100);
  };

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <h3 className="text-lg font-medium text-gray-900">{t('jobs.title')}</h3>

        {/* Status Filter */}
        <div className="flex items-center space-x-2">
          <select
            value={selectedStatus}
            onChange={(e) => setSelectedStatus(e.target.value)}
            className="text-sm border-gray-300 rounded-md shadow-sm focus:border-blue-500 focus:ring-blue-500"
          >
            <option value="all">{t('jobs.filters.all')}</option>
            <option value="pending">{t('jobs.filters.pending')}</option>
            <option value="processing">{t('jobs.filters.processing')}</option>
            <option value="completed">{t('jobs.filters.completed')}</option>
            <option value="failed">{t('jobs.filters.failed')}</option>
          </select>

          <button
            onClick={() => dispatch(fetchImportJobs({ storefrontId }))}
            disabled={isLoading}
            className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
          >
            <svg
              className={`w-4 h-4 mr-2 ${isLoading ? 'animate-spin' : ''}`}
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"
              />
            </svg>
            {t('actions.refresh')}
          </button>
        </div>
      </div>

      {/* Jobs List */}
      {filteredJobs.length === 0 ? (
        <div className="text-center py-12">
          <svg
            className="w-12 h-12 text-gray-400 mx-auto mb-4"
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
          <h3 className="text-lg font-medium text-gray-900 mb-1">
            {t('jobs.empty.title')}
          </h3>
          <p className="text-gray-600">{t('jobs.empty.description')}</p>
        </div>
      ) : (
        <div className="bg-white shadow overflow-hidden sm:rounded-md">
          <ul className="divide-y divide-gray-200">
            {filteredJobs.map((job) => (
              <li key={job.id}>
                <div
                  onClick={() => handleJobClick(job)}
                  className="px-4 py-4 hover:bg-gray-50 cursor-pointer"
                >
                  <div className="flex items-center justify-between">
                    <div className="flex items-center space-x-3">
                      {/* Status */}
                      <span
                        className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          IMPORT_STATUS_COLORS[
                            job.status as keyof typeof IMPORT_STATUS_COLORS
                          ]
                        }`}
                      >
                        {getStatusIcon(job.status)}
                        <span className="ml-1">
                          {t(`jobs.status.${job.status}`)}
                        </span>
                      </span>

                      {/* File Info */}
                      <div>
                        <div className="flex items-center space-x-2">
                          <p className="text-sm font-medium text-gray-900">
                            {job.file_name}
                          </p>
                          <span className="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-gray-100 text-gray-800">
                            {job.file_type.toUpperCase()}
                          </span>
                        </div>
                        <p className="text-xs text-gray-500">
                          {formatDate(job.created_at)}
                        </p>
                      </div>
                    </div>

                    {/* Stats */}
                    <div className="flex items-center space-x-6">
                      {job.status === 'processing' && (
                        <div className="flex items-center space-x-2">
                          <div className="w-24 bg-gray-200 rounded-full h-2">
                            <div
                              className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                              style={{ width: `${calculateProgress(job)}%` }}
                            />
                          </div>
                          <span className="text-xs text-gray-500">
                            {calculateProgress(job)}%
                          </span>
                        </div>
                      )}

                      <div className="text-right">
                        <div className="flex items-center space-x-4 text-sm text-gray-500">
                          {job.total_records > 0 && (
                            <>
                              <span className="text-green-600">
                                ✓ {job.successful_records}
                              </span>
                              {job.failed_records > 0 && (
                                <span className="text-red-600">
                                  ✗ {job.failed_records}
                                </span>
                              )}
                              <span>/ {job.total_records}</span>
                            </>
                          )}
                        </div>
                      </div>

                      {/* Actions */}
                      <div className="flex items-center space-x-2">
                        <button
                          onClick={() => {
                            dispatch(setCurrentJob(job));
                            dispatch(setJobDetailsModalOpen(true));
                          }}
                          className="text-blue-600 hover:text-blue-800 text-sm"
                        >
                          {t('actions.viewDetails')}
                        </button>

                        {job.failed_records > 0 && (
                          <button
                            onClick={(_e) => handleShowErrors(job)}
                            className="text-red-600 hover:text-red-800 text-sm"
                          >
                            {t('actions.viewErrors')}
                          </button>
                        )}

                        {job.status === 'processing' && (
                          <button
                            onClick={(e) => handleCancelJob(job.id, e)}
                            className="text-red-600 hover:text-red-800 text-sm"
                          >
                            {t('actions.cancel')}
                          </button>
                        )}

                        {job.status === 'failed' && (
                          <button
                            onClick={(e) => handleRetryJob(job.id, e)}
                            className="text-blue-600 hover:text-blue-800 text-sm"
                          >
                            {t('actions.retry')}
                          </button>
                        )}

                        <svg
                          className="w-5 h-5 text-gray-400"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24"
                        >
                          <path
                            strokeLinecap="round"
                            strokeLinejoin="round"
                            strokeWidth={2}
                            d="M9 5l7 7-7 7"
                          />
                        </svg>
                      </div>
                    </div>
                  </div>

                  {/* Error Message */}
                  {job.status === 'failed' && job.error_message && (
                    <div className="mt-2 pl-3">
                      <p className="text-sm text-red-600">
                        {job.error_message}
                      </p>
                    </div>
                  )}
                </div>
              </li>
            ))}
          </ul>
        </div>
      )}
    </div>
  );
}
