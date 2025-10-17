'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import {
  ChevronDownIcon,
  ChevronUpIcon,
  ClipboardDocumentIcon,
  CheckIcon,
} from '@heroicons/react/24/outline';

interface RequestData {
  method: string;
  endpoint: string;
  headers?: Record<string, string>;
  body?: any;
}

interface ResponseData {
  status: number;
  statusText: string;
  headers?: Record<string, string>;
  data: any;
}

interface RequestResponseDisplayProps {
  request: RequestData;
  response: ResponseData | null;
  loading: boolean;
  error: string | null;
  processingTime?: number;
}

export default function RequestResponseDisplay({
  request,
  response,
  loading,
  error,
  processingTime,
}: RequestResponseDisplayProps) {
  const t = useTranslations('postexpressTest');

  const [requestExpanded, setRequestExpanded] = useState(true);
  const [responseExpanded, setResponseExpanded] = useState(true);
  const [requestCopied, setRequestCopied] = useState(false);
  const [responseCopied, setResponseCopied] = useState(false);

  const copyToClipboard = async (
    text: string,
    type: 'request' | 'response'
  ) => {
    try {
      await navigator.clipboard.writeText(text);
      if (type === 'request') {
        setRequestCopied(true);
        setTimeout(() => setRequestCopied(false), 2000);
      } else {
        setResponseCopied(true);
        setTimeout(() => setResponseCopied(false), 2000);
      }
    } catch (err) {
      console.error('Failed to copy:', err);
    }
  };

  const sanitizeHeaders = (headers?: Record<string, string>) => {
    if (!headers) return {};
    const sanitized = { ...headers };
    if (sanitized.Authorization) {
      sanitized.Authorization = 'Bearer ***';
    }
    return sanitized;
  };

  const getStatusBadgeClass = (status: number) => {
    if (status >= 200 && status < 300) return 'badge-success';
    if (status >= 400 && status < 500) return 'badge-warning';
    if (status >= 500) return 'badge-error';
    return 'badge-info';
  };

  return (
    <div className="space-y-4">
      {/* Request Section */}
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <h3 className="font-semibold text-lg">
                {t('request.title', { defaultValue: '📤 Request' })}
              </h3>
              <div className="badge badge-primary">{request.method}</div>
            </div>
            <div className="flex items-center gap-2">
              <button
                onClick={() =>
                  copyToClipboard(
                    JSON.stringify(
                      {
                        method: request.method,
                        endpoint: request.endpoint,
                        headers: sanitizeHeaders(request.headers),
                        body: request.body,
                      },
                      null,
                      2
                    ),
                    'request'
                  )
                }
                className="btn btn-ghost btn-sm"
                title={t('request.copy', { defaultValue: 'Copy request' })}
              >
                {requestCopied ? (
                  <CheckIcon className="w-4 h-4 text-success" />
                ) : (
                  <ClipboardDocumentIcon className="w-4 h-4" />
                )}
              </button>
              <button
                onClick={() => setRequestExpanded(!requestExpanded)}
                className="btn btn-ghost btn-sm"
              >
                {requestExpanded ? (
                  <ChevronUpIcon className="w-4 h-4" />
                ) : (
                  <ChevronDownIcon className="w-4 h-4" />
                )}
              </button>
            </div>
          </div>

          {requestExpanded && (
            <div className="space-y-3 mt-4">
              {/* Endpoint */}
              <div>
                <div className="text-xs text-base-content/60 mb-1">
                  {t('request.endpoint', { defaultValue: 'Endpoint' })}:
                </div>
                <div className="font-mono text-sm bg-base-200 p-2 rounded">
                  {request.endpoint}
                </div>
              </div>

              {/* Headers */}
              {request.headers && (
                <div>
                  <div className="text-xs text-base-content/60 mb-1">
                    {t('request.headers', { defaultValue: 'Headers' })}:
                  </div>
                  <pre className="bg-base-200 p-3 rounded text-xs overflow-x-auto">
                    {JSON.stringify(sanitizeHeaders(request.headers), null, 2)}
                  </pre>
                </div>
              )}

              {/* Body */}
              {request.body && (
                <div>
                  <div className="text-xs text-base-content/60 mb-1">
                    {t('request.body', { defaultValue: 'Body' })}:
                  </div>
                  <pre className="bg-base-200 p-3 rounded text-xs overflow-x-auto max-h-96">
                    {JSON.stringify(request.body, null, 2)}
                  </pre>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Response Section */}
      <div className="card bg-base-100 border border-base-300">
        <div className="card-body p-4">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2">
              <h3 className="font-semibold text-lg">
                {t('response.title', { defaultValue: '📥 Response' })}
              </h3>
              {response && (
                <div
                  className={`badge ${getStatusBadgeClass(response.status)}`}
                >
                  {response.status} {response.statusText}
                </div>
              )}
              {processingTime !== undefined && (
                <div className="badge badge-ghost">{processingTime}ms</div>
              )}
            </div>
            {response && (
              <div className="flex items-center gap-2">
                <button
                  onClick={() =>
                    copyToClipboard(
                      JSON.stringify(
                        {
                          status: response.status,
                          statusText: response.statusText,
                          headers: response.headers,
                          data: response.data,
                        },
                        null,
                        2
                      ),
                      'response'
                    )
                  }
                  className="btn btn-ghost btn-sm"
                  title={t('response.copy', { defaultValue: 'Copy response' })}
                >
                  {responseCopied ? (
                    <CheckIcon className="w-4 h-4 text-success" />
                  ) : (
                    <ClipboardDocumentIcon className="w-4 h-4" />
                  )}
                </button>
                <button
                  onClick={() => setResponseExpanded(!responseExpanded)}
                  className="btn btn-ghost btn-sm"
                >
                  {responseExpanded ? (
                    <ChevronUpIcon className="w-4 h-4" />
                  ) : (
                    <ChevronDownIcon className="w-4 h-4" />
                  )}
                </button>
              </div>
            )}
          </div>

          {loading && (
            <div className="flex items-center justify-center py-8">
              <span className="loading loading-spinner loading-lg"></span>
              <span className="ml-4 text-base-content/60">
                {t('response.loading', {
                  defaultValue: 'Waiting for response...',
                })}
              </span>
            </div>
          )}

          {error && !loading && (
            <div className="alert alert-error mt-4">
              <span>{error}</span>
            </div>
          )}

          {response && !loading && responseExpanded && (
            <div className="space-y-3 mt-4">
              {/* Headers */}
              {response.headers && (
                <div>
                  <div className="text-xs text-base-content/60 mb-1">
                    {t('response.headers', { defaultValue: 'Headers' })}:
                  </div>
                  <pre className="bg-base-200 p-3 rounded text-xs overflow-x-auto">
                    {JSON.stringify(response.headers, null, 2)}
                  </pre>
                </div>
              )}

              {/* Data */}
              <div>
                <div className="text-xs text-base-content/60 mb-1">
                  {t('response.data', { defaultValue: 'Data' })}:
                </div>
                <pre className="bg-base-200 p-3 rounded text-xs overflow-x-auto max-h-96">
                  {JSON.stringify(response.data, null, 2)}
                </pre>
              </div>
            </div>
          )}

          {!response && !loading && !error && (
            <div className="text-center py-8 text-base-content/50">
              {t('response.empty', {
                defaultValue: 'Response will appear here after request',
              })}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
