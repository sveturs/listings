'use client';

import React, { useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { useTranslations } from 'next-intl';
import {
  setImportModalOpen,
  fetchImportJobs,
  fetchImportFormats,
} from '@/store/slices/importSlice';
import ImportWizard from './ImportWizard';
import ImportJobsList from './ImportJobsList';
import ImportJobDetails from './ImportJobDetails';
import ImportErrorsModal from './ImportErrorsModal';

interface ImportManagerProps {
  storefrontId: number;
  storefrontSlug?: string;
}

export default function ImportManager({
  storefrontId,
  storefrontSlug,
}: ImportManagerProps) {
  const dispatch = useAppDispatch();
  const t = useTranslations('storefronts');
  const { jobs, formats } = useAppSelector((state) => state.import);

  useEffect(() => {
    dispatch(fetchImportJobs({ storefrontId }));
    dispatch(fetchImportFormats());
  }, [dispatch, storefrontId]);

  const handleStartImport = () => {
    dispatch(setImportModalOpen(true));
  };

  const handleImportSuccess = (_jobId: number) => {
    // Refresh jobs list after successful import
    dispatch(fetchImportJobs({ storefrontId }));
  };

  const getJobsStats = () => {
    const pending = jobs.filter((job) => job.status === 'pending').length;
    const processing = jobs.filter((job) => job.status === 'processing').length;
    const completed = jobs.filter((job) => job.status === 'completed').length;
    const failed = jobs.filter((job) => job.status === 'failed').length;

    return { pending, processing, completed, failed };
  };

  const stats = getJobsStats();

  return (
    <div className="space-y-8">
      {/* Header */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <div className="flex items-center justify-between">
            <div>
              <h2 className="text-lg font-medium text-gray-900">
                {t('manager.title')}
              </h2>
              <p className="mt-1 text-sm text-gray-600">
                {t('manager.description')}
              </p>
            </div>
            <button
              onClick={handleStartImport}
              className="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
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
                  d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
                />
              </svg>
              {t('manager.startImport')}
            </button>
          </div>
        </div>
      </div>

      {/* Statistics Cards */}
      {jobs.length > 0 && (
        <div className="grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <svg
                    className="h-6 w-6 text-yellow-400"
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
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">
                      {t('stats.pending')}
                    </dt>
                    <dd className="text-lg font-medium text-gray-900">
                      {stats.pending}
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <svg
                    className="h-6 w-6 text-blue-400 animate-spin"
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
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">
                      {t('stats.processing')}
                    </dt>
                    <dd className="text-lg font-medium text-gray-900">
                      {stats.processing}
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <svg
                    className="h-6 w-6 text-green-400"
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
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">
                      {t('stats.completed')}
                    </dt>
                    <dd className="text-lg font-medium text-gray-900">
                      {stats.completed}
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>

          <div className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
              <div className="flex items-center">
                <div className="flex-shrink-0">
                  <svg
                    className="h-6 w-6 text-red-400"
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
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">
                      {t('stats.failed')}
                    </dt>
                    <dd className="text-lg font-medium text-gray-900">
                      {stats.failed}
                    </dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Supported Formats Info */}
      {formats && (
        <div className="bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium text-gray-900 mb-4">
              {t('formats.title')}
            </h3>
            <div className="grid grid-cols-1 gap-6 sm:grid-cols-3">
              {Object.entries(formats.supported_formats).map(
                ([format, info]) => (
                  <div
                    key={format}
                    className="border border-gray-200 rounded-lg p-4"
                  >
                    <div className="flex items-center mb-2">
                      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                        {format.toUpperCase()}
                      </span>
                    </div>
                    <p className="text-sm text-gray-600 mb-2">
                      {info.description}
                    </p>
                    <div className="text-xs text-gray-500">
                      <p>
                        {t('formats.extensions')}:{' '}
                        {info.file_extensions.join(', ')}
                      </p>
                      {info.encoding && (
                        <p>
                          {t('formats.encoding')}: {info.encoding}
                        </p>
                      )}
                    </div>
                  </div>
                )
              )}
            </div>

            <div className="mt-4 text-sm text-gray-600">
              <p>
                {t('formats.maxFileSize')}: {formats.max_file_size}
              </p>
              <p>
                {t('formats.maxProducts')}:{' '}
                {formats.max_products_per_import.toLocaleString()}
              </p>
            </div>
          </div>
        </div>
      )}

      {/* Import Jobs List */}
      <div className="bg-white shadow rounded-lg">
        <div className="px-4 py-5 sm:p-6">
          <ImportJobsList
            storefrontId={storefrontId}
            autoRefresh={true}
            refreshInterval={5000}
          />
        </div>
      </div>

      {/* Modals */}
      <ImportWizard
        storefrontId={storefrontId}
        storefrontSlug={storefrontSlug}
        onSuccess={handleImportSuccess}
      />
      <ImportJobDetails />
      <ImportErrorsModal />
    </div>
  );
}
