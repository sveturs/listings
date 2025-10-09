'use client';

import React, { useState, useEffect } from 'react';
import { useAppDispatch, useAppSelector } from '@/store/hooks';
import { useTranslations } from 'next-intl';
import {
  fetchImportFormats,
  setImportModalOpen,
  setSelectedFiles,
  setImportUrl,
  setSelectedFileType,
  setUpdateMode,
  setCategoryMappingMode,
  importFromFile,
  importFromUrl,
  previewImportFile,
  downloadCsvTemplate,
  resetForm,
  clearError,
  clearPreview,
} from '@/store/slices/importSlice';
import ImportPreviewTable from './ImportPreviewTable';
import ImportAnalysisWizard from './ImportAnalysisWizard';
import {
  validateFileType,
  validateFileSize,
  formatFileSize,
  getFileTypeFromExtension,
} from '@/services/importApi';
import { IMPORT_FILE_CONFIG } from '@/types/import';

interface ImportWizardProps {
  storefrontId: number;
  storefrontSlug?: string;
  onSuccess?: (jobId: number) => void;
  onClose?: () => void;
}

export default function ImportWizard({
  storefrontId,
  storefrontSlug,
  onSuccess,
  onClose,
}: ImportWizardProps) {
  const dispatch = useAppDispatch();
  const t = useTranslations('storefronts');
  const {
    isImportModalOpen,
    selectedFiles,
    importUrl,
    selectedFileType,
    updateMode,
    categoryMappingMode,
    isLoading,
    isUploading,
    uploadProgress,
    error,
    formats,
    previewData,
    isPreviewLoading,
    previewError,
  } = useAppSelector((state) => state.import);

  const [activeTab, setActiveTab] = useState<'file' | 'url'>('file');
  const [dragActive, setDragActive] = useState(false);
  const [showPreview, setShowPreview] = useState(false);
  type ImportMode = 'classic' | 'enhanced';
  const [importMode, setImportMode] = useState<ImportMode>('classic');
  const [previewLimit, setPreviewLimit] = useState<number>(100); // Default: show 100 rows instead of 10

  useEffect(() => {
    if (isImportModalOpen && !formats) {
      dispatch(fetchImportFormats());
    }
  }, [isImportModalOpen, formats, dispatch]);

  const handleDrag = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    if (e.type === 'dragenter' || e.type === 'dragover') {
      setDragActive(true);
    } else if (e.type === 'dragleave') {
      setDragActive(false);
    }
  };

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    e.stopPropagation();
    setDragActive(false);

    const files = Array.from(e.dataTransfer.files);
    handleFileSelection(files);
  };

  const handleFileInput = (e: React.ChangeEvent<HTMLInputElement>) => {
    const files = Array.from(e.target.files || []);
    handleFileSelection(files);
  };

  const handleFileSelection = (files: File[]) => {
    if (files.length === 0) return;

    const file = files[0]; // Поддерживаем только один файл

    // Validate file type
    if (!validateFileType(file, IMPORT_FILE_CONFIG.allowedTypes)) {
      dispatch(clearError());
      dispatch(setSelectedFiles([]));
      alert(t('errors.invalidFileType'));
      return;
    }

    // Validate file size
    if (!validateFileSize(file, IMPORT_FILE_CONFIG.maxFileSize)) {
      dispatch(clearError());
      dispatch(setSelectedFiles([]));
      alert(
        t('errors.fileTooLarge', {
          size: formatFileSize(IMPORT_FILE_CONFIG.maxFileSize),
        })
      );
      return;
    }

    // Auto-detect file type
    const detectedType = getFileTypeFromExtension(file.name);
    if (detectedType) {
      dispatch(setSelectedFileType(detectedType));
    }

    dispatch(setSelectedFiles([file]));
    dispatch(clearError());
  };

  const handlePreviewFile = async () => {
    if (selectedFiles.length === 0 || !selectedFileType) return;

    try {
      await dispatch(
        previewImportFile({
          storefrontId: storefrontSlug ? undefined : storefrontId,
          storefrontSlug,
          file: selectedFiles[0],
          fileType: selectedFileType,
          previewLimit: previewLimit,
        })
      ).unwrap();

      setShowPreview(true);
    } catch (error) {
      console.error('Preview failed:', error);
    }
  };

  const handleImportFromFile = async () => {
    if (selectedFiles.length === 0 || !selectedFileType) return;

    try {
      const result = await dispatch(
        importFromFile({
          storefrontId,
          storefrontSlug,
          file: selectedFiles[0],
          options: {
            file_type: selectedFileType,
            update_mode: updateMode,
            category_mapping_mode: categoryMappingMode,
          },
        })
      ).unwrap();

      onSuccess?.(result.id);
      handleClose();
    } catch (error) {
      console.error('Import failed:', error);
    }
  };

  const handleImportFromUrl = async () => {
    if (!importUrl || !selectedFileType) return;

    try {
      const result = await dispatch(
        importFromUrl({
          storefrontId,
          storefrontSlug,
          request: {
            file_url: importUrl,
            file_type: selectedFileType,
            update_mode: updateMode,
            category_mapping_mode: categoryMappingMode,
          },
        })
      ).unwrap();

      onSuccess?.(result.id);
      handleClose();
    } catch (error) {
      console.error('Import from URL failed:', error);
    }
  };

  const handleClose = () => {
    dispatch(resetForm());
    dispatch(clearPreview());
    setShowPreview(false);
    setImportMode('classic'); // Reset to classic mode on close
    dispatch(setImportModalOpen(false));
    onClose?.();
  };

  const handleBackFromPreview = () => {
    setShowPreview(false);
    dispatch(clearPreview());
  };

  const handleDownloadTemplate = () => {
    dispatch(downloadCsvTemplate());
  };

  if (!isImportModalOpen) return null;

  // If enhanced mode is selected, render ImportAnalysisWizard instead
  if (importMode === 'enhanced') {
    return (
      <ImportAnalysisWizard
        storefrontId={storefrontId}
        storefrontSlug={storefrontSlug}
        onClose={handleClose}
        onSuccess={onSuccess}
        onSwitchToClassic={() => setImportMode('classic')}
      />
    );
  }

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between p-6 border-b">
          <div className="flex items-center space-x-4">
            <h2 className="text-2xl font-semibold text-gray-900">
              {t('title')}
            </h2>
            {/* Import Mode Toggle */}
            <div className="flex items-center space-x-2 text-sm">
              <button
                onClick={() => setImportMode('classic' as ImportMode)}
                className={
                  importMode === 'classic'
                    ? 'px-3 py-1 rounded-md font-medium transition-colors bg-blue-100 text-blue-700'
                    : 'px-3 py-1 rounded-md font-medium transition-colors text-gray-600 hover:bg-gray-100'
                }
              >
                {t('importMode.classic')}
              </button>
              <button
                onClick={() => setImportMode('enhanced' as ImportMode)}
                className={
                  (importMode as string) === 'enhanced'
                    ? 'px-3 py-1 rounded-md font-medium transition-colors bg-purple-100 text-purple-700'
                    : 'px-3 py-1 rounded-md font-medium transition-colors text-gray-600 hover:bg-gray-100'
                }
              >
                {t('importMode.enhanced')} ✨
              </button>
            </div>
          </div>
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
        <div className="p-6">
          {!showPreview ? (
            <>
              {/* Tabs */}
              <div className="flex space-x-1 bg-gray-100 rounded-lg p-1 mb-6">
                <button
                  onClick={() => setActiveTab('file')}
                  className={`flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors ${
                    activeTab === 'file'
                      ? 'bg-white text-blue-600 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  {t('tabs.uploadFile')}
                </button>
                <button
                  onClick={() => setActiveTab('url')}
                  className={`flex-1 py-2 px-4 rounded-md text-sm font-medium transition-colors ${
                    activeTab === 'url'
                      ? 'bg-white text-blue-600 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'
                  }`}
                >
                  {t('tabs.importFromUrl')}
                </button>
              </div>

              {/* File Upload Tab */}
              {activeTab === 'file' && (
                <div className="space-y-6">
                  {/* Drag & Drop Area */}
                  <div
                    className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
                      dragActive
                        ? 'border-blue-500 bg-blue-50'
                        : selectedFiles.length > 0
                          ? 'border-green-500 bg-green-50'
                          : 'border-gray-300 hover:border-gray-400'
                    }`}
                    onDragEnter={handleDrag}
                    onDragLeave={handleDrag}
                    onDragOver={handleDrag}
                    onDrop={handleDrop}
                  >
                    {selectedFiles.length > 0 ? (
                      <div className="space-y-2">
                        <div className="flex items-center justify-center">
                          <svg
                            className="w-8 h-8 text-green-500"
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
                        <p className="text-sm font-medium text-gray-900">
                          {selectedFiles[0].name}
                        </p>
                        <p className="text-xs text-gray-500">
                          {formatFileSize(selectedFiles[0].size)}
                        </p>
                        <button
                          onClick={() => dispatch(setSelectedFiles([]))}
                          className="text-sm text-red-600 hover:text-red-800"
                        >
                          {t('actions.removeFile')}
                        </button>
                      </div>
                    ) : (
                      <div className="space-y-2">
                        <div className="flex items-center justify-center">
                          <svg
                            className="w-8 h-8 text-gray-400"
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
                        </div>
                        <p className="text-sm text-gray-600">
                          {t('dragDrop.title')}
                        </p>
                        <p className="text-xs text-gray-500">
                          {t('dragDrop.subtitle')}
                        </p>
                        <input
                          type="file"
                          accept=".csv,.xml,.zip"
                          onChange={handleFileInput}
                          className="hidden"
                          id="file-upload"
                        />
                        <label
                          htmlFor="file-upload"
                          className="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md text-white bg-blue-600 hover:bg-blue-700 cursor-pointer"
                        >
                          {t('actions.selectFile')}
                        </label>
                      </div>
                    )}
                  </div>

                  {/* Download Template */}
                  <div className="bg-gray-50 rounded-lg p-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <h4 className="text-sm font-medium text-gray-900">
                          {t('template.title')}
                        </h4>
                        <p className="text-sm text-gray-600">
                          {t('template.description')}
                        </p>
                      </div>
                      <button
                        onClick={handleDownloadTemplate}
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
                        {t('actions.downloadTemplate')}
                      </button>
                    </div>
                  </div>
                </div>
              )}

              {/* URL Import Tab */}
              {activeTab === 'url' && (
                <div className="space-y-6">
                  <div>
                    <label
                      htmlFor="import-url"
                      className="block text-sm font-medium text-gray-700 mb-2"
                    >
                      {t('url.label')}
                    </label>
                    <input
                      type="url"
                      id="import-url"
                      value={importUrl}
                      onChange={(e) => dispatch(setImportUrl(e.target.value))}
                      placeholder="https://example.com/products.csv"
                      className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                    />
                    <p className="mt-2 text-sm text-gray-500">
                      {t('url.help')}
                    </p>
                  </div>
                </div>
              )}

              {/* File Type Selection */}
              {(selectedFiles.length > 0 || importUrl) && (
                <div className="space-y-4 mt-6">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">
                      {t('fileType.label')}
                    </label>
                    <select
                      value={selectedFileType}
                      onChange={(e) =>
                        dispatch(setSelectedFileType(e.target.value as any))
                      }
                      className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                    >
                      <option value="">{t('fileType.select')}</option>
                      <option value="csv">CSV</option>
                      <option value="xml">XML</option>
                      <option value="zip">ZIP</option>
                    </select>
                  </div>

                  {/* Import Options */}
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        {t('options.updateMode.label')}
                      </label>
                      <select
                        value={updateMode}
                        onChange={(e) =>
                          dispatch(setUpdateMode(e.target.value as any))
                        }
                        className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                      >
                        <option value="upsert">
                          {t('options.updateMode.upsert')}
                        </option>
                        <option value="create_only">
                          {t('options.updateMode.createOnly')}
                        </option>
                        <option value="update_only">
                          {t('options.updateMode.updateOnly')}
                        </option>
                      </select>
                    </div>

                    <div>
                      <label className="block text-sm font-medium text-gray-700 mb-2">
                        {t('options.categoryMapping.label')}
                      </label>
                      <select
                        value={categoryMappingMode}
                        onChange={(e) =>
                          dispatch(
                            setCategoryMappingMode(e.target.value as any)
                          )
                        }
                        className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm"
                      >
                        <option value="auto">
                          {t('options.categoryMapping.auto')}
                        </option>
                        <option value="manual">
                          {t('options.categoryMapping.manual')}
                        </option>
                        <option value="skip">
                          {t('options.categoryMapping.skip')}
                        </option>
                      </select>
                    </div>
                  </div>
                </div>
              )}

              {/* Upload Progress */}
              {isUploading && uploadProgress && (
                <div className="mt-6">
                  <div className="flex items-center justify-between text-sm text-gray-600 mb-2">
                    <span>{t('progress.uploading')}</span>
                    <span>{uploadProgress.percentage}%</span>
                  </div>
                  <div className="w-full bg-gray-200 rounded-full h-2">
                    <div
                      className="bg-blue-600 h-2 rounded-full transition-all duration-300"
                      style={{ width: `${uploadProgress.percentage}%` }}
                    />
                  </div>
                </div>
              )}

              {/* Error Display */}
              {error && (
                <div className="mt-6 p-4 bg-red-50 border border-red-200 rounded-md">
                  <div className="flex">
                    <svg
                      className="w-5 h-5 text-red-400"
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
                    <div className="ml-3">
                      <p className="text-sm text-red-800">{error}</p>
                    </div>
                  </div>
                </div>
              )}

              {/* Actions */}
              <div className="flex items-center justify-between mt-8 pt-6 border-t">
                <button
                  onClick={handleClose}
                  className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  {t('actions.cancel')}
                </button>

                <div className="flex items-center space-x-3">
                  {activeTab === 'file' &&
                    selectedFiles.length > 0 &&
                    selectedFileType && (
                      <>
                        {/* Preview rows selector */}
                        <div className="flex items-center space-x-2">
                          <label
                            htmlFor="preview-limit"
                            className="text-sm text-gray-600"
                          >
                            Preview rows:
                          </label>
                          <select
                            id="preview-limit"
                            value={previewLimit}
                            onChange={(e) =>
                              setPreviewLimit(Number(e.target.value))
                            }
                            className="px-3 py-2 text-sm border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
                          >
                            <option value={10}>10</option>
                            <option value={25}>25</option>
                            <option value={50}>50</option>
                            <option value={100}>100</option>
                          </select>
                        </div>

                        <button
                          onClick={handlePreviewFile}
                          disabled={isPreviewLoading}
                          className="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
                        >
                          {isPreviewLoading ? (
                            <>
                              <svg
                                className="animate-spin -ml-1 mr-2 h-4 w-4 inline-block"
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                              >
                                <circle
                                  className="opacity-25"
                                  cx="12"
                                  cy="12"
                                  r="10"
                                  stroke="currentColor"
                                  strokeWidth="4"
                                ></circle>
                                <path
                                  className="opacity-75"
                                  fill="currentColor"
                                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                                ></path>
                              </svg>
                              Loading Preview...
                            </>
                          ) : (
                            'Preview Import'
                          )}
                        </button>
                      </>
                    )}

                  <button
                    onClick={
                      activeTab === 'file'
                        ? handleImportFromFile
                        : handleImportFromUrl
                    }
                    disabled={
                      isLoading ||
                      isUploading ||
                      !selectedFileType ||
                      (activeTab === 'file'
                        ? selectedFiles.length === 0
                        : !importUrl)
                    }
                    className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 disabled:opacity-50"
                  >
                    {isLoading || isUploading
                      ? t('actions.importing')
                      : t('actions.import')}
                  </button>
                </div>
              </div>
            </>
          ) : (
            /* Preview Results */
            <div className="space-y-6">
              {previewError ? (
                <div className="bg-red-50 border border-red-200 rounded-lg p-4">
                  <div className="flex items-center">
                    <svg
                      className="w-5 h-5 text-red-500 mr-2"
                      fill="currentColor"
                      viewBox="0 0 20 20"
                    >
                      <path
                        fillRule="evenodd"
                        d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                        clipRule="evenodd"
                      />
                    </svg>
                    <p className="text-sm text-red-800">{previewError}</p>
                  </div>
                </div>
              ) : previewData ? (
                <ImportPreviewTable previewData={previewData} />
              ) : null}

              <div className="flex justify-between pt-4 border-t">
                <button
                  onClick={handleBackFromPreview}
                  className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50"
                >
                  {t('actions.back')}
                </button>
                <button
                  onClick={handleImportFromFile}
                  disabled={isLoading || isUploading}
                  className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  {isLoading || isUploading ? (
                    <>
                      <svg
                        className="animate-spin -ml-1 mr-2 h-4 w-4 inline-block"
                        xmlns="http://www.w3.org/2000/svg"
                        fill="none"
                        viewBox="0 0 24 24"
                      >
                        <circle
                          className="opacity-25"
                          cx="12"
                          cy="12"
                          r="10"
                          stroke="currentColor"
                          strokeWidth="4"
                        ></circle>
                        <path
                          className="opacity-75"
                          fill="currentColor"
                          d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                        ></path>
                      </svg>
                      Importing...
                    </>
                  ) : (
                    'Proceed with Import'
                  )}
                </button>
              </div>
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
