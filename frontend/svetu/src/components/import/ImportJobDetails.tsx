'use client';

import React, { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { useTranslations } from 'next-intl';
import {
  setJobDetailsModalOpen,
  fetchJobDetails,
  setErrorsModalOpen,
  cancelImportJob,
  retryImportJob,
} from '@/store/slices/importSlice';
import { IMPORT_STATUS_COLORS, IMPORT_STATUS_ICONS } from '@/types/import';

export default function ImportJobDetails() {
  const dispatch = useAppDispatch();
  const t = useTranslations('storefronts');
  const { isJobDetailsModalOpen, currentJob, isLoading } = useAppSelector(
    (state) => state.import
  );

  useEffect(() => {
    if (isJobDetailsModalOpen && currentJob) {
      dispatch(
        fetchJobDetails({
          storefrontId: currentJob.storefront_id,
          jobId: currentJob.id,
        })
      );
    }
  }, [isJobDetailsModalOpen, currentJob, dispatch]);

  const handleClose = () => {
    dispatch(setJobDetailsModalOpen(false));
  };

  const handleShowErrors = () => {
    dispatch(setErrorsModalOpen(true));
  };

  const handleCancelJob = async () => {
    if (currentJob && confirm(t('actions.confirmCancel'))) {
      await dispatch(
        cancelImportJob({
          storefrontId: currentJob.storefront_id,
          jobId: currentJob.id,
        })
      );
      handleClose();
    }
  };

  const handleRetryJob = async () => {
    if (currentJob) {
      await dispatch(
        retryImportJob({
          storefrontId: currentJob.storefront_id,
          jobId: currentJob.id,
        })
      );
      handleClose();
    }
  };

  if (!isJobDetailsModalOpen || !currentJob) return null;

  const formatDate = (dateString?: string) => {
    if (!dateString) return t('jobs.details.notSet');
    return new Date(dateString).toLocaleString();
  };

  const calculateProgress = () => {
    if (currentJob.total_records === 0) return 0;
    return Math.round(
      (currentJob.processed_records / currentJob.total_records) * 100
    );
  };

  const getProcessingTime = () => {
    if (!currentJob.started_at) return null;

    const startTime = new Date(currentJob.started_at);
    const endTime = currentJob.completed_at
      ? new Date(currentJob.completed_at)
      : new Date();

    const diffMs = endTime.getTime() - startTime.getTime();
    const diffSeconds = Math.floor(diffMs / 1000);
    const diffMinutes = Math.floor(diffSeconds / 60);
    const diffHours = Math.floor(diffMinutes / 60);

    if (diffHours > 0) {
      return `${diffHours}h ${diffMinutes % 60}m ${diffSeconds % 60}s`;
    } else if (diffMinutes > 0) {
      return `${diffMinutes}m ${diffSeconds % 60}s`;
    } else {
      return `${diffSeconds}s`;
    }
  };

  const getStatusIcon = (status: string) => {
    const iconName =
      IMPORT_STATUS_ICONS[status as keyof typeof IMPORT_STATUS_ICONS] ||
      'help-circle';

    switch (iconName) {
      case 'clock':
        return (
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
              d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"
            />
          </svg>
        );
      case 'refresh':
        return (
          <svg
            className="w-5 h-5 animate-spin"
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
            className="w-5 h-5"
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
            className="w-5 h-5"
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
            className="w-5 h-5"
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

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <h2 className="text-xl font-semibold text-gray-900">
            {t('jobs.details.title')}
          </h2>
          <button
            onClick={handleClose}
            className="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <svg
              className="w-6 h-6"
              fill="none"
              stroke="currentColor"
              viewBox="0 0 24 24"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                strokeWidth={2}
                d="M6 18L18 6M6 6l12 12"
              />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Status and Progress */}
          <div className="bg-gray-50 rounded-lg p-4">
            <div className="flex items-center justify-between mb-4">
              <span
                className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${
                  IMPORT_STATUS_COLORS[
                    currentJob.status as keyof typeof IMPORT_STATUS_COLORS
                  ]
                }`}
              >
                {getStatusIcon(currentJob.status)}
                <span className="ml-2">
                  {t(`jobs.status.${currentJob.status}`)}
                </span>
              </span>

              {currentJob.status === 'processing' && (
                <span className="text-sm text-gray-600">
                  {calculateProgress()}% {t('jobs.details.complete')}
                </span>
              )}
            </div>

            {currentJob.status === 'processing' && (
              <div className="w-full bg-gray-200 rounded-full h-3">
                <div
                  className="bg-blue-600 h-3 rounded-full transition-all duration-300"
                  style={{ width: `${calculateProgress()}%` }}
                />
              </div>
            )}
          </div>

          {/* File Information */}
          <div>
            <h3 className="text-lg font-medium text-gray-900 mb-3">
              {t('jobs.details.fileInfo')}
            </h3>
            <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <dt className="text-sm font-medium text-gray-500">
                  {t('jobs.details.fileName')}
                </dt>
                <dd className="mt-1 text-sm text-gray-900">
                  {currentJob.file_name}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">
                  {t('jobs.details.fileType')}
                </dt>
                <dd className="mt-1 text-sm text-gray-900">
                  {currentJob.file_type.toUpperCase()}
                </dd>
              </div>
              {currentJob.file_url && (
                <div className="sm:col-span-2">
                  <dt className="text-sm font-medium text-gray-500">
                    {t('jobs.details.fileUrl')}
                  </dt>
                  <dd className="mt-1 text-sm text-blue-600 break-all">
                    <a
                      href={currentJob.file_url}
                      target="_blank"
                      rel="noopener noreferrer"
                    >
                      {currentJob.file_url}
                    </a>
                  </dd>
                </div>
              )}
            </dl>
          </div>

          {/* Processing Statistics */}
          <div>
            <h3 className="text-lg font-medium text-gray-900 mb-3">
              {t('jobs.details.statistics')}
            </h3>
            <div className="grid grid-cols-2 gap-4 sm:grid-cols-4">
              <div className="bg-blue-50 rounded-lg p-3 text-center">
                <div className="text-2xl font-bold text-blue-600">
                  {currentJob.total_records}
                </div>
                <div className="text-sm text-blue-600">
                  {t('jobs.details.total')}
                </div>
              </div>
              <div className="bg-yellow-50 rounded-lg p-3 text-center">
                <div className="text-2xl font-bold text-yellow-600">
                  {currentJob.processed_records}
                </div>
                <div className="text-sm text-yellow-600">
                  {t('jobs.details.processed')}
                </div>
              </div>
              <div className="bg-green-50 rounded-lg p-3 text-center">
                <div className="text-2xl font-bold text-green-600">
                  {currentJob.successful_records}
                </div>
                <div className="text-sm text-green-600">
                  {t('jobs.details.successful')}
                </div>
              </div>
              <div className="bg-red-50 rounded-lg p-3 text-center">
                <div className="text-2xl font-bold text-red-600">
                  {currentJob.failed_records}
                </div>
                <div className="text-sm text-red-600">
                  {t('jobs.details.failed')}
                </div>
              </div>
            </div>
          </div>

          {/* Timing Information */}
          <div>
            <h3 className="text-lg font-medium text-gray-900 mb-3">
              {t('jobs.details.timing')}
            </h3>
            <dl className="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <dt className="text-sm font-medium text-gray-500">
                  {t('jobs.details.created')}
                </dt>
                <dd className="mt-1 text-sm text-gray-900">
                  {formatDate(currentJob.created_at)}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">
                  {t('jobs.details.started')}
                </dt>
                <dd className="mt-1 text-sm text-gray-900">
                  {formatDate(currentJob.started_at)}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">
                  {t('jobs.details.completed')}
                </dt>
                <dd className="mt-1 text-sm text-gray-900">
                  {formatDate(currentJob.completed_at)}
                </dd>
              </div>
              {getProcessingTime() && (
                <div>
                  <dt className="text-sm font-medium text-gray-500">
                    {t('jobs.details.duration')}
                  </dt>
                  <dd className="mt-1 text-sm text-gray-900">
                    {getProcessingTime()}
                  </dd>
                </div>
              )}
            </dl>
          </div>

          {/* Error Message */}
          {currentJob.status === 'failed' && currentJob.error_message && (
            <div>
              <h3 className="text-lg font-medium text-gray-900 mb-3">
                {t('jobs.details.error')}
              </h3>
              <div className="bg-red-50 border border-red-200 rounded-md p-4">
                <p className="text-sm text-red-800">
                  {currentJob.error_message}
                </p>
              </div>
            </div>
          )}

          {/* Actions */}
          <div className="flex items-center justify-between pt-6 border-t">
            <div className="flex space-x-3">
              {currentJob.failed_records > 0 && (
                <button
                  onClick={handleShowErrors}
                  className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
                >
                  <svg
                    className="w-4 h-4 mr-2"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  {t('actions.viewErrors')}
                </button>
              )}

              {currentJob.status === 'processing' && (
                <button
                  onClick={handleCancelJob}
                  className="inline-flex items-center px-4 py-2 border border-red-300 shadow-sm text-sm font-medium rounded-md text-red-700 bg-white hover:bg-red-50"
                >
                  <svg
                    className="w-4 h-4 mr-2"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M6 18L18 6M6 6l12 12"
                    />
                  </svg>
                  {t('actions.cancel')}
                </button>
              )}

              {currentJob.status === 'failed' && (
                <button
                  onClick={handleRetryJob}
                  disabled={isLoading}
                  className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 disabled:opacity-50"
                >
                  <svg
                    className="w-4 h-4 mr-2"
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
                  {isLoading ? t('actions.retrying') : t('actions.retry')}
                </button>
              )}
            </div>

            <button
              onClick={handleClose}
              className="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
            >
              {t('actions.close')}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}
