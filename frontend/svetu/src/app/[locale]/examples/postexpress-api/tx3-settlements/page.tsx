'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import RequestResponseDisplay from '@/components/postexpress/RequestResponseDisplay';
import {
  MagnifyingGlassIcon,
  CheckCircleIcon,
} from '@heroicons/react/24/outline';

export default function TX3SettlementsPage() {
  const t = useTranslations('postexpressTest');

  // Предзаполненные значения
  const [query, setQuery] = useState('Beograd');
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
      query,
      num_records: numRecords,
    };

    // Сохраняем request для отображения
    setRequestData({
      method: 'POST',
      endpoint: '/api/v2/postexpress/test/tx3-settlements',
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ***',
      },
      body: payload,
    });

    try {
      const response = await apiClient.post(
        '/postexpress/test/tx3-settlements',
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

      if (response.data?.data?.naselja) {
        setResults(response.data.data.naselja);
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
    setQuery('Beograd');
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
      <div className="bg-gradient-to-r from-blue-600 to-cyan-600 text-white py-8">
        <div className="container mx-auto px-4">
          <h1 className="text-4xl font-bold mb-2">{t('tx3.title')}</h1>
          <p className="text-xl opacity-90">{t('tx3.description')}</p>
          <div className="mt-4 flex gap-2 flex-wrap">
            <div className="badge badge-success badge-lg gap-2">
              <CheckCircleIcon className="w-4 h-4" />
              {t('badges.apiReady')}
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
                <h2 className="card-title text-2xl mb-4">{t('form.title')}</h2>

                {/* Query */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('tx3.query')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={query}
                    onChange={(e) => setQuery(e.target.value)}
                    placeholder={t('tx3.queryPlaceholder')}
                    required
                  />
                  <label className="label">
                    <span className="label-text-alt opacity-70">
                      {t('common.required')}
                    </span>
                  </label>
                </div>

                {/* Num Records */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('tx3.numRecords')}
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
                    {t('form.reset')}
                  </button>
                  <button
                    type="submit"
                    className="btn btn-primary"
                    disabled={loading}
                  >
                    {loading ? (
                      <>
                        <span className="loading loading-spinner"></span>
                        {t('form.submitting')}
                      </>
                    ) : (
                      <>
                        <MagnifyingGlassIcon className="w-5 h-5" />
                        {t('form.submit')}
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
                    {t('tx3.success', { count: results.length })}
                  </h2>
                  <div className="overflow-x-auto">
                    <table className="table table-sm">
                      <thead>
                        <tr>
                          <th>ID</th>
                          <th>{t('tx3.query')}</th>
                          <th>Postal Code</th>
                        </tr>
                      </thead>
                      <tbody>
                        {results.map((settlement: any, idx: number) => (
                          <tr key={idx}>
                            <td className="font-mono">{settlement.Id}</td>
                            <td>{settlement.Naziv}</td>
                            <td>{settlement.PostanskiBroj || '-'}</td>
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
                  <p className="text-base-content/50">{t('response.empty')}</p>
                </div>
              </div>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}
