'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import RequestResponseDisplay from '@/components/postexpress/RequestResponseDisplay';
import {
  CheckCircleIcon,
  XCircleIcon,
  ExclamationTriangleIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';

export default function TX6ValidatePage() {
  const t = useTranslations('postexpressTest.tx6');
  const tCommon = useTranslations('postexpressTest.common');

  // Предзаполненные значения для валидации адреса
  const [settlementId, setSettlementId] = useState(100001); // Beograd
  const [streetId, setStreetId] = useState(1186); // Takovska
  const [houseNumber, setHouseNumber] = useState('2');
  const [postalCode, setPostalCode] = useState('11000');

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [requestData, setRequestData] = useState<any | null>(null);
  const [responseData, setResponseData] = useState<any | null>(null);
  const [processingTime, setProcessingTime] = useState<number | undefined>(
    undefined
  );
  const [validationResult, setValidationResult] = useState<any | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setValidationResult(null);
    setResponseData(null);

    const startTime = Date.now();

    const payload = {
      settlement_id: settlementId,
      street_id: streetId,
      house_number: houseNumber,
      postal_code: postalCode,
    };

    // Сохраняем request для отображения
    setRequestData({
      method: 'POST',
      endpoint: '/api/v2/postexpress/test/tx6-validate-address',
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ***',
      },
      body: payload,
    });

    try {
      const response = await apiClient.post(
        '/postexpress/test/tx6-validate-address',
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

      if (response.data?.data) {
        setValidationResult(response.data.data);
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
    setStreetId(1186);
    setHouseNumber('2');
    setPostalCode('11000');
    setValidationResult(null);
    setError(null);
    setRequestData(null);
    setResponseData(null);
    setProcessingTime(undefined);
  };

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-gradient-to-r from-orange-600 to-amber-600 text-white py-8">
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
            {t('back')}
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

      {/* Test Data Warning */}
      <div className="container mx-auto px-4 pt-8">
        <div className="alert alert-warning shadow-lg mb-6">
          <ExclamationTriangleIcon className="w-6 h-6" />
          <div>
            <h3 className="font-bold">{t('testDataWarning.title')}</h3>
            <div className="text-sm">{t('testDataWarning.description')}</div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 pb-8">
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Form */}
          <div>
            <form
              onSubmit={handleSubmit}
              className="card bg-base-100 shadow-xl"
            >
              <div className="card-body">
                <h2 className="card-title text-2xl mb-4">{t('formTitle')}</h2>

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
                    required
                  />
                  <label className="label">
                    <span className="label-text-alt opacity-70">
                      {t('settlementIdHint')}
                    </span>
                  </label>
                </div>

                {/* Street ID */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('streetId')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered"
                    value={streetId}
                    onChange={(e) => setStreetId(Number(e.target.value))}
                    required
                  />
                  <label className="label">
                    <span className="label-text-alt opacity-70">
                      {t('streetIdHint')}
                    </span>
                  </label>
                </div>

                {/* House Number */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('houseNumber')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={houseNumber}
                    onChange={(e) => setHouseNumber(e.target.value)}
                    placeholder={t('houseNumberPlaceholder')}
                    required
                  />
                </div>

                {/* Postal Code */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('postalCode')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={postalCode}
                    onChange={(e) => setPostalCode(e.target.value)}
                    placeholder={t('postalCodePlaceholder')}
                    required
                  />
                </div>

                {/* Buttons */}
                <div className="card-actions justify-end mt-6">
                  <button
                    type="button"
                    onClick={handleReset}
                    className="btn btn-ghost"
                  >
                    {t('reset')}
                  </button>
                  <button
                    type="submit"
                    className="btn btn-primary"
                    disabled={loading}
                  >
                    {loading ? (
                      <>
                        <span className="loading loading-spinner"></span>
                        {t('checking')}
                      </>
                    ) : (
                      <>
                        <CheckCircleIcon className="w-5 h-5" />
                        {t('checkAddress')}
                      </>
                    )}
                  </button>
                </div>
              </div>
            </form>

            {/* Error Display */}
            {error && !loading && (
              <div className="card bg-base-100 shadow-xl mt-6">
                <div className="card-body">
                  <div className="alert alert-error">
                    <XCircleIcon className="w-6 h-6" />
                    <div>
                      <h3 className="font-bold">{t('errorTitle')}</h3>
                      <div className="text-sm">{error}</div>
                    </div>
                  </div>
                </div>
              </div>
            )}

            {/* Validation Result */}
            {validationResult && (
              <div className="card bg-base-100 shadow-xl mt-6">
                <div className="card-body">
                  <h2 className="card-title text-lg mb-4">
                    {t('resultTitle')}
                  </h2>

                  {/* Status Badge */}
                  <div className="mb-4">
                    {validationResult.PostojiAdresa ? (
                      <div className="alert alert-success">
                        <CheckCircleIcon className="w-6 h-6" />
                        <span className="font-semibold">
                          {t('addressExists')}
                        </span>
                      </div>
                    ) : (
                      <div className="alert alert-error">
                        <XCircleIcon className="w-6 h-6" />
                        <span className="font-semibold">
                          {t('addressNotFound')}
                        </span>
                      </div>
                    )}
                  </div>

                  {/* Address Details */}
                  {validationResult.PostojiAdresa && (
                    <div className="overflow-x-auto">
                      <table className="table table-sm">
                        <tbody>
                          <tr>
                            <td className="font-semibold">
                              {t('settlementId')}
                            </td>
                            <td className="font-mono">
                              {validationResult.IdNaselje}
                            </td>
                          </tr>
                          <tr>
                            <td className="font-semibold">
                              {t('settlementName')}
                            </td>
                            <td>{validationResult.NazivNaselja}</td>
                          </tr>
                          <tr>
                            <td className="font-semibold">{t('streetId')}</td>
                            <td className="font-mono">
                              {validationResult.IdUlica}
                            </td>
                          </tr>
                          <tr>
                            <td className="font-semibold">{t('streetName')}</td>
                            <td>{validationResult.NazivUlice}</td>
                          </tr>
                          <tr>
                            <td className="font-semibold">
                              {t('houseNumber')}
                            </td>
                            <td>{validationResult.Broj}</td>
                          </tr>
                          <tr>
                            <td className="font-semibold">{t('postalCode')}</td>
                            <td className="font-mono">
                              {validationResult.PostanskiBroj}
                            </td>
                          </tr>
                        </tbody>
                      </table>
                    </div>
                  )}
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
