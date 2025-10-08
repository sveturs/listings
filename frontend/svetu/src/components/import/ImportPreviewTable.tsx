'use client';

import React from 'react';
import type { ImportPreviewResponse, ImportPreviewRow } from '@/types/import';

interface ImportPreviewTableProps {
  previewData: ImportPreviewResponse;
  onClose?: () => void;
}

export default function ImportPreviewTable({
  previewData,
  onClose,
}: ImportPreviewTableProps) {
  const getValidationStatusBadge = (row: ImportPreviewRow) => {
    if (row.is_valid) {
      return (
        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
          <svg className="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path
              fillRule="evenodd"
              d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
              clipRule="evenodd"
            />
          </svg>
          Valid
        </span>
      );
    }

    return (
      <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800">
        <svg className="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
          <path
            fillRule="evenodd"
            d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
            clipRule="evenodd"
          />
        </svg>
        Invalid
      </span>
    );
  };

  const renderCellValue = (value: any): string => {
    if (value === null || value === undefined) return '-';
    if (typeof value === 'object') return JSON.stringify(value);
    return String(value);
  };

  return (
    <div className="space-y-4">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h3 className="text-lg font-medium text-gray-900">Import Preview</h3>
          <p className="text-sm text-gray-600">
            Showing {previewData.preview_rows.length} of{' '}
            {previewData.total_rows} rows
          </p>
        </div>
        {onClose && (
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-600"
          >
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path
                fillRule="evenodd"
                d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                clipRule="evenodd"
              />
            </svg>
          </button>
        )}
      </div>

      {/* Validation Summary */}
      <div
        className={`rounded-lg p-4 ${
          previewData.validation_ok
            ? 'bg-green-50 border border-green-200'
            : 'bg-yellow-50 border border-yellow-200'
        }`}
      >
        <div className="flex items-center">
          {previewData.validation_ok ? (
            <svg
              className="w-5 h-5 text-green-500 mr-2"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                clipRule="evenodd"
              />
            </svg>
          ) : (
            <svg
              className="w-5 h-5 text-yellow-500 mr-2"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fillRule="evenodd"
                d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z"
                clipRule="evenodd"
              />
            </svg>
          )}
          <div>
            <p
              className={`text-sm font-medium ${
                previewData.validation_ok ? 'text-green-800' : 'text-yellow-800'
              }`}
            >
              {previewData.validation_ok
                ? 'All preview rows are valid'
                : 'Some rows have validation errors'}
            </p>
            {previewData.error_summary && (
              <p
                className={`text-sm ${
                  previewData.validation_ok
                    ? 'text-green-700'
                    : 'text-yellow-700'
                }`}
              >
                {previewData.error_summary}
              </p>
            )}
          </div>
        </div>
      </div>

      {/* Table */}
      <div className="overflow-x-auto border border-gray-200 rounded-lg">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Line
              </th>
              <th className="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              {previewData.headers &&
                previewData.headers.map((header) => (
                  <th
                    key={header}
                    className="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    {header}
                  </th>
                ))}
              {!previewData.headers &&
                previewData.preview_rows.length > 0 &&
                Object.keys(previewData.preview_rows[0].data).map((key) => (
                  <th
                    key={key}
                    className="px-3 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                  >
                    {key}
                  </th>
                ))}
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {previewData.preview_rows.map((row) => (
              <React.Fragment key={row.line_number}>
                <tr
                  className={`${
                    !row.is_valid ? 'bg-red-50' : ''
                  } hover:bg-gray-50`}
                >
                  <td className="px-3 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {row.line_number}
                  </td>
                  <td className="px-3 py-4 whitespace-nowrap text-sm">
                    {getValidationStatusBadge(row)}
                  </td>
                  {previewData.headers
                    ? previewData.headers.map((header) => (
                        <td
                          key={header}
                          className="px-3 py-4 whitespace-nowrap text-sm text-gray-900"
                        >
                          {renderCellValue(row.data[header])}
                        </td>
                      ))
                    : Object.entries(row.data).map(([key, value]) => (
                        <td
                          key={key}
                          className="px-3 py-4 whitespace-nowrap text-sm text-gray-900"
                        >
                          {renderCellValue(value)}
                        </td>
                      ))}
                </tr>
                {/* Error row - показываем только если есть ошибки */}
                {row.errors && row.errors.length > 0 && (
                  <tr className="bg-red-50">
                    <td
                      colSpan={
                        2 +
                        (previewData.headers?.length ||
                          Object.keys(row.data).length)
                      }
                      className="px-3 py-3"
                    >
                      <div className="space-y-1">
                        {row.errors.map((error, errorIndex) => (
                          <div
                            key={errorIndex}
                            className="flex items-start text-sm"
                          >
                            <svg
                              className="w-4 h-4 text-red-500 mr-2 mt-0.5 flex-shrink-0"
                              fill="currentColor"
                              viewBox="0 0 20 20"
                            >
                              <path
                                fillRule="evenodd"
                                d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                                clipRule="evenodd"
                              />
                            </svg>
                            <div>
                              <span className="font-medium text-red-800">
                                {error.field}:
                              </span>{' '}
                              <span className="text-red-700">
                                {error.message}
                              </span>
                              {error.value !== undefined &&
                                error.value !== null && (
                                  <span className="text-red-600 ml-1">
                                    (value: {renderCellValue(error.value)})
                                  </span>
                                )}
                            </div>
                          </div>
                        ))}
                      </div>
                    </td>
                  </tr>
                )}
              </React.Fragment>
            ))}
          </tbody>
        </table>
      </div>

      {/* File Info */}
      <div className="bg-gray-50 rounded-lg p-4">
        <div className="grid grid-cols-3 gap-4 text-sm">
          <div>
            <p className="text-gray-500">File Type</p>
            <p className="font-medium text-gray-900 uppercase">
              {previewData.file_type}
            </p>
          </div>
          <div>
            <p className="text-gray-500">Total Rows</p>
            <p className="font-medium text-gray-900">
              {previewData.total_rows}
            </p>
          </div>
          <div>
            <p className="text-gray-500">Preview Rows</p>
            <p className="font-medium text-gray-900">
              {previewData.preview_rows.length}
            </p>
          </div>
        </div>
      </div>
    </div>
  );
}
