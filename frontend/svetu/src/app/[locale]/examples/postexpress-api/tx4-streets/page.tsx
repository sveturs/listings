'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import RequestResponseDisplay from '@/components/postexpress/RequestResponseDisplay';
import { MapIcon, CheckCircleIcon } from '@heroicons/react/24/outline';
import Link from 'next/link';

export default function TX4StreetsPage() {
  const t = useTranslations('postexpressTest.tx4');
  const tCommon = useTranslations('postexpressTest.common');

  // Предзаполненные значения
  const [settlementId, setSettlementId] = useState(100001); // Beograd
  const [query, setQuery] = useState('Takovska');
  const [numRecords, setNumRecords] = useState(10);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [requestData, setRequestData] = useState<any | null>(null);
  const [responseData, setResponseData] = useState<any | null>(null);
  const [processingTime, setProcessingTime] = useState<number | undefined>(
    undefined
  );
  const [results, setResults] = useState<any[] | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setResults(null);
    setResponseData(null);

    const startTime = Date.now();

    const payload = {
      settlement_id: settlementId,
      query,
      num_records: numRecords,
    };

    // Сохраняем request для отображения
    setRequestData({
      method: 'POST',
      endpoint: '/api/v2/postexpress/test/tx4-streets',
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ***',
      },
      body: payload,
    });

    try {
      const response = await apiClient.post(
        '/postexpress/test/tx4-streets',
        payload
      );
      const endTime = Date.now();
      setProcessingTime(endTime - startTime);

      // Сохраняем response для отображения
      setResponseData({
        status: 200,
        statusText: 'OK',
        headers: response.headers || {},
        data: response.data,
      });

      if (response.data?.data?.ulice) {
        setResults(response.data.data.ulice);
      }
    } catch (err: any) {
      const endTime = Date.now();
      setProcessingTime(endTime - startTime);

      const errorMessage =
        err.response?.data?.message || err.message || 'Unknown error';
      setError(errorMessage);

      // Сохраняем error response
      setResponseData({
        status: err.response?.status || 500,
        statusText: err.response?.statusText || 'Error',
        headers: err.response?.headers || {},
        data: err.response?.data || { error: errorMessage },
      });
    } finally {
      setLoading(false);
    }
  };

  const handleReset = () => {
    setSettlementId(100001);
    setQuery('Takovska');
    setNumRecords(10);
    setResults(null);
    setError(null);
    setRequestData(null);
    setResponseData(null);
    setProcessingTime(undefined);
  };

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-gradient-to-r from-green-600 to-teal-600 text-white py-8">
        <div className="container mx-auto px-4">
          <Link
            href="/examples/postexpress-api"
            className="btn btn-ghost btn-sm mb-4 text-white hover:bg-white/20"
          >
            <svg
              xmlns="http://www.w3.org/2000/svg"
              fill="none"
              viewBox="0 0 24 24"
              strokeWidth={1.5}
              stroke="currentColor"
              className="w-4 h-4"
            >
              <path
                strokeLinecap="round"
                strokeLinejoin="round"
                d="M10.5 19.5L3 12m0 0l7.5-7.5M3 12h18"
              />
            </svg>
            {t('back') || 'Back'}
          </Link>
          <h1 className="text-4xl font-bold mb-2">{t('title')}</h1>
          <p className="text-xl opacity-90">{t('description')}</p>
          <div className="mt-4 flex gap-2 flex-wrap">
            <div className="badge badge-success badge-lg gap-2">
              <CheckCircleIcon className="w-4 h-4" />
              API Ready
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Form */}
          <div>
            <form
              onSubmit={handleSubmit}
              className="card bg-base-100 shadow-xl"
            >
              <div className="card-body">
                <h2 className="card-title text-2xl mb-4">
                  {t('formTitle') || 'Search Parameters'}
                </h2>

                {/* Settlement ID */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('settlementId')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered"
                    value={settlementId}
                    onChange={(e) => setSettlementId(Number(e.target.value))}
                    placeholder={t('settlementIdPlaceholder')}
                    required
                  />
                  <label className="label">
                    <span className="label-text-alt opacity-70">
                      {t('settlementIdHint') || 'Use TX 3 to get the ID'}
                    </span>
                  </label>
                </div>

                {/* Query */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('query')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder={t('queryPlaceholder')}
                    required
                  />
                  <label className="label">
                    <span className="label-text-alt opacity-70">
                      {t('queryHint') || 'Enter full or partial street name'}
                    </span>
                  </label>
                </div>

                {/* Num Records */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('numRecords')}
                    </span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered"
                    value={numRecords}
                    onChange={(e) => setNumRecords(Number(e.target.value))}
                    min="1"
                    max="100"
                  />
                </div>

                {/* Buttons */}
                <div className="card-actions justify-end mt-6">
                  <button
                    type="button"
                    onClick={handleReset}
                    className="btn btn-ghost"
                  >
                    {t('reset') || 'Reset'}
                  </button>
                  <button
                    type="submit"
                    className="btn btn-primary"
                    disabled={loading}
                  >
                    {loading ? (
                      <>
                        <span className="loading loading-spinner"></span>
                        {t('searching') || 'Searching...'}
                      </>
                    ) : (
                      <>
                        <MapIcon className="w-5 h-5" />
                        {t('searchButton') || 'Find Streets'}
                      </>
                    )}
                  </button>
                </div>
              </div>
            </form>

            {/* Results Summary */}
            {results && (
              <div className="card bg-base-100 shadow-xl mt-6">
                <div className="card-body">
                  <h2 className="card-title text-lg">
                    {t('success', { count: results.length })}
                  </h2>
                  <div className="overflow-x-auto">
                    <table className="table table-sm">
                      <thead>
                        <tr>
                          <th>ID</th>
                          <th>{t('streetName') || 'Street Name'}</th>
                          <th>{t('settlementId')}</th>
                        </tr>
                      </thead>
                      <tbody>
                        {results.map((street: any, idx: number) => (
                          <tr key={idx}>
                            <td className="font-mono">{street.Id}</td>
                            <td>{street.Naziv}</td>
                            <td>{street.IdNaselje}</td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>
                </div>
              </div>
            )}
          </div>

          {/* Request/Response Display */}
          <div>
            {requestData ? (
              <RequestResponseDisplay
                request={requestData}
                response={responseData}
                loading={loading}
                error={error}
                processingTime={processingTime}
              />
            ) : (
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body text-center">
                  <p className="text-base-content/50">
                    {t('resultsPlaceholder')}
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
