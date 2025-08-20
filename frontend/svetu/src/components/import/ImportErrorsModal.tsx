'use client';

import React, { useEffect, useState, useCallback } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { useTranslations } from 'next-intl';
import { setErrorsModalOpen } from '@/store/slices/importSlice';
import { ImportApi, downloadFile } from '@/services/importApi';
import type { ImportError } from '@/types/import';

export default function ImportErrorsModal() {
  const dispatch = useAppDispatch();
  const t = useTranslations('storefronts');
  const { isErrorsModalOpen, currentJob } = useAppSelector(
    (state) => state.import
  );

  const [errors, setErrors] = useState<ImportError[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [currentPage, setCurrentPage] = useState(1);
  const [searchTerm, setSearchTerm] = useState('');
  const [selectedField, setSelectedField] = useState('all');

  const errorsPerPage = 10;

  const loadErrors = useCallback(async () => {
    if (!currentJob) return;

    setIsLoading(true);
    try {
      const jobDetails = await ImportApi.getJobDetails(currentJob.id);
      setErrors((jobDetails as any).errors || []);
    } catch (error) {
      console.error('Failed to load import errors:', error);
    } finally {
      setIsLoading(false);
    }
  }, [currentJob]);

  useEffect(() => {
    if (isErrorsModalOpen && currentJob) {
      loadErrors();
    }
  }, [isErrorsModalOpen, currentJob, loadErrors]);

  const handleClose = () => {
    dispatch(setErrorsModalOpen(false));
    setErrors([]);
    setCurrentPage(1);
    setSearchTerm('');
    setSelectedField('all');
  };

  const handleExportErrors = async () => {
    if (!currentJob) return;

    try {
      const blob = await ImportApi.exportResults(currentJob.id);
      downloadFile(blob, `import_errors_${currentJob.id}.csv`);
    } catch (error) {
      console.error('Failed to export errors:', error);
    }
  };

  const filteredErrors = errors.filter((error) => {
    const matchesSearch =
      searchTerm === '' ||
      error.error_message.toLowerCase().includes(searchTerm.toLowerCase()) ||
      error.raw_data.toLowerCase().includes(searchTerm.toLowerCase());

    const matchesField =
      selectedField === 'all' || error.field_name === selectedField;

    return matchesSearch && matchesField;
  });

  const totalPages = Math.ceil(filteredErrors.length / errorsPerPage);
  const startIndex = (currentPage - 1) * errorsPerPage;
  const endIndex = startIndex + errorsPerPage;
  const currentErrors = filteredErrors.slice(startIndex, endIndex);

  const uniqueFields = [...new Set(errors.map((error) => error.field_name))];

  if (!isErrorsModalOpen || !currentJob) return null;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-6xl w-full max-h-[90vh] flex flex-col">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <div>
            <h2 className="text-xl font-semibold text-gray-900">
              {t('errors.title')}
            </h2>
            <p className="text-sm text-gray-600">
              {t('errors.subtitle', { fileName: currentJob.file_name })}
            </p>
          </div>
          <div className="flex items-center space-x-3">
            {errors.length > 0 && (
              <button
                onClick={handleExportErrors}
                className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50"
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
                    d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
                  />
                </svg>
                {t('errors.export')}
              </button>
            )}
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
        </div>

        {/* Filters */}
        <div className="p-6 border-b bg-gray-50">
          <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between space-y-3 sm:space-y-0 sm:space-x-3">
            <div className="flex-1">
              <input
                type="text"
                placeholder={t('errors.searchPlaceholder')}
                value={searchTerm}
                onChange={(e) => setSearchTerm(e.target.value)}
                className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
              />
            </div>
            <div className="flex items-center space-x-3">
              <select
                value={selectedField}
                onChange={(e) => setSelectedField(e.target.value)}
                className="rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
              >
                <option value="all">{t('errors.allFields')}</option>
                {uniqueFields.map((field) => (
                  <option key={field} value={field}>
                    {field}
                  </option>
                ))}
              </select>
              <span className="text-sm text-gray-600">
                {filteredErrors.length} {t('errors.of')} {errors.length}
              </span>
            </div>
          </div>
        </div>

        {/* Content */}
        <div className="flex-1 overflow-y-auto">
          {isLoading ? (
            <div className="flex items-center justify-center h-64">
              <div className="flex items-center space-x-2">
                <svg
                  className="w-5 h-5 animate-spin text-blue-600"
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
                <span className="text-gray-600">{t('errors.loading')}</span>
              </div>
            </div>
          ) : errors.length === 0 ? (
            <div className="text-center py-12">
              <svg
                className="w-12 h-12 text-green-400 mx-auto mb-4"
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
              <h3 className="text-lg font-medium text-gray-900 mb-1">
                {t('errors.noErrors.title')}
              </h3>
              <p className="text-gray-600">
                {t('errors.noErrors.description')}
              </p>
            </div>
          ) : currentErrors.length === 0 ? (
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
                  d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
                />
              </svg>
              <h3 className="text-lg font-medium text-gray-900 mb-1">
                {t('errors.noResults.title')}
              </h3>
              <p className="text-gray-600">
                {t('errors.noResults.description')}
              </p>
            </div>
          ) : (
            <div className="divide-y divide-gray-200">
              {currentErrors.map((error, index) => (
                <div key={`${error.id}-${index}`} className="p-6">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      <div className="flex items-center space-x-2 mb-2">
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
                          {t('errors.line')} {error.line_number}
                        </span>
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                          {error.field_name}
                        </span>
                      </div>

                      <p className="text-sm text-red-800 mb-3">
                        {error.error_message}
                      </p>

                      {error.raw_data && (
                        <div className="bg-gray-50 rounded-md p-3">
                          <p className="text-xs font-medium text-gray-500 mb-1">
                            {t('errors.rawData')}:
                          </p>
                          <p className="text-sm text-gray-800 font-mono break-all">
                            {error.raw_data.length > 200
                              ? `${error.raw_data.substring(0, 200)}...`
                              : error.raw_data}
                          </p>
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Pagination */}
        {totalPages > 1 && (
          <div className="px-6 py-4 border-t bg-gray-50">
            <div className="flex items-center justify-between">
              <div className="text-sm text-gray-700">
                {t('errors.showing')} {startIndex + 1} {t('errors.to')}{' '}
                {Math.min(endIndex, filteredErrors.length)} {t('errors.of')}{' '}
                {filteredErrors.length} {t('errors.entries')}
              </div>

              <div className="flex items-center space-x-2">
                <button
                  onClick={() =>
                    setCurrentPage((prev) => Math.max(prev - 1, 1))
                  }
                  disabled={currentPage === 1}
                  className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <svg
                    className="w-4 h-4 mr-1"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M15 19l-7-7 7-7"
                    />
                  </svg>
                  {t('errors.previous')}
                </button>

                <span className="text-sm text-gray-700">
                  {currentPage} / {totalPages}
                </span>

                <button
                  onClick={() =>
                    setCurrentPage((prev) => Math.min(prev + 1, totalPages))
                  }
                  disabled={currentPage === totalPages}
                  className="inline-flex items-center px-3 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {t('errors.next')}
                  <svg
                    className="w-4 h-4 ml-1"
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
                </button>
              </div>
            </div>
          </div>
        )}

        {/* Footer */}
        <div className="px-6 py-4 border-t bg-gray-50">
          <div className="flex justify-end">
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
