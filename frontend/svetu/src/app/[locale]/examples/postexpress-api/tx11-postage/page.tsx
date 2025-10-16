'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { apiClient } from '@/services/api-client';
import RequestResponseDisplay from '@/components/postexpress/RequestResponseDisplay';
import {
  CheckCircleIcon,
  CurrencyDollarIcon,
} from '@heroicons/react/24/outline';
import Link from 'next/link';

export default function TX11PostagePage() {
  const t = useTranslations('postexpressTest.tx11');
  const tCommon = useTranslations('postexpressTest.common');

  // Предзаполненные значения для расчёта стоимости
  const [serviceId, setServiceId] = useState(71); // Express pošiljka
  const [fromPostalCode, setFromPostalCode] = useState('11000'); // Beograd
  const [toPostalCode, setToPostalCode] = useState('21000'); // Novi Sad
  const [weight, setWeight] = useState(500); // 500 грамм
  const [services, setServices] = useState('PNA'); // Povratnica

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [requestData, setRequestData] = useState<any | null>(null);
  const [responseData, setResponseData] = useState<any | null>(null);
  const [processingTime, setProcessingTime] = useState<number | undefined>(
    undefined
  );
  const [postageResult, setPostageResult] = useState<any | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(null);
    setPostageResult(null);
    setResponseData(null);

    const startTime = Date.now();

    const payload = {
      service_id: serviceId,
      from_postal_code: fromPostalCode,
      to_postal_code: toPostalCode,
      weight,
      services: services.trim() || undefined,
    };

    // Сохраняем request для отображения
    setRequestData({
      method: 'POST',
      endpoint: '/api/v2/postexpress/test/tx11-calculate-postage',
      headers: {
        'Content-Type': 'application/json',
        Authorization: 'Bearer ***',
      },
      body: payload,
    });

    try {
      const response = await apiClient.post(
        '/postexpress/test/tx11-calculate-postage',
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
        setPostageResult(response.data.data);
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
    setServiceId(71);
    setFromPostalCode('11000');
    setToPostalCode('21000');
    setWeight(500);
    setServices('PNA');
    setPostageResult(null);
    setError(null);
    setRequestData(null);
    setResponseData(null);
    setProcessingTime(undefined);
  };

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-gradient-to-r from-rose-600 to-pink-600 text-white py-8">
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
          <h1 className="text-4xl font-bold mb-2">
            {t('pageTitle')}
          </h1>
          <p className="text-xl opacity-90">
            {t('pageDescription')}
          </p>
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
                <h2 className="card-title text-2xl mb-4">{t('formTitle')}</h2>

                {/* Service ID */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('serviceId')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered"
                    value={serviceId}
                    onChange={(e) => setServiceId(Number(e.target.value))}
                    required
                  />
                  <label className="label">
                    <span className="label-text-alt opacity-70">
                      {t('serviceIdHint')}
                    </span>
                  </label>
                </div>

                {/* From Postal Code */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('fromPostalCode')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={fromPostalCode}
                    onChange={(e) => setFromPostalCode(e.target.value)}
                    placeholder={t('fromPostalCodePlaceholder')}
                    required
                  />
                </div>

                {/* To Postal Code */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('toPostalCode')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={toPostalCode}
                    onChange={(e) => setToPostalCode(e.target.value)}
                    placeholder={t('toPostalCodePlaceholder')}
                    required
                  />
                </div>

                {/* Weight */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('weight')} <span className="text-error">*</span>
                    </span>
                  </label>
                  <input
                    type="number"
                    className="input input-bordered"
                    value={weight}
                    onChange={(e) => setWeight(Number(e.target.value))}
                    min="1"
                    max="30000"
                    step="100"
                    required
                  />
                  <label className="label">
                    <span className="label-text-alt opacity-70">
                      {t('weightHint')}
                    </span>
                  </label>
                </div>

                {/* Additional Services */}
                <div className="form-control">
                  <label className="label">
                    <span className="label-text font-semibold">
                      {t('services')}
                    </span>
                  </label>
                  <input
                    type="text"
                    className="input input-bordered"
                    value={services}
                    onChange={(e) => setServices(e.target.value)}
                    placeholder={t('servicesPlaceholder')}
                  />
                  <label className="label">
                    <span className="label-text-alt opacity-70">
                      {t('servicesHint')}
                    </span>
                  </label>
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
                        {t('calculating')}
                      </>
                    ) : (
                      <>
                        <CurrencyDollarIcon className="w-5 h-5" />
                        {t('calculateButton')}
                      </>
                    )}
                  </button>
                </div>
              </div>
            </form>

            {/* Postage Result */}
            {postageResult && (
              <div className="card bg-base-100 shadow-xl mt-6">
                <div className="card-body">
                  <h2 className="card-title text-lg mb-4">
                    {t('resultTitle')}
                  </h2>

                  {/* Total Amount */}
                  <div className="mb-4">
                    <div className="alert alert-info">
                      <CurrencyDollarIcon className="w-6 h-6" />
                      <div>
                        <div className="font-semibold text-lg">
                          {t('totalCost')}: {postageResult.Iznos} RSD
                        </div>
                        <div className="text-sm opacity-70">
                          {t('weightLabel')}: {postageResult.Masa} г
                        </div>
                      </div>
                    </div>
                  </div>

                  {/* Price Breakdown */}
                  {postageResult.CenovniStavovi &&
                    postageResult.CenovniStavovi.length > 0 && (
                      <div className="overflow-x-auto">
                        <h3 className="font-semibold mb-2">
                          {t('priceBreakdown')}:
                        </h3>
                        <table className="table table-sm table-zebra">
                          <thead>
                            <tr>
                              <th>{t('service')}</th>
                              <th className="text-right">{t('amount')}</th>
                            </tr>
                          </thead>
                          <tbody>
                            {postageResult.CenovniStavovi.map(
                              (item: any, idx: number) => (
                                <tr key={idx}>
                                  <td>{item.Opis || item.Naziv || '-'}</td>
                                  <td className="text-right font-mono">
                                    {item.Iznos}
                                  </td>
                                </tr>
                              )
                            )}
                          </tbody>
                          <tfoot>
                            <tr className="font-semibold">
                              <td>{t('total')}:</td>
                              <td className="text-right font-mono text-lg">
                                {postageResult.Iznos} RSD
                              </td>
                            </tr>
                          </tfoot>
                        </table>
                      </div>
                    )}

                  {/* Additional Info */}
                  <div className="mt-4 text-sm opacity-70">
                    {postageResult.Otkupnina !== undefined && (
                      <div>{t('codAmount')}: {postageResult.Otkupnina} RSD</div>
                    )}
                    {postageResult.Vrednost !== undefined && (
                      <div>
                        {t('insuredValue')}: {postageResult.Vrednost} RSD
                      </div>
                    )}
                  </div>

                  {/* Route Info */}
                  <div className="mt-4">
                    <div className="text-sm opacity-70">
                      <div className="font-semibold mb-2">{t('routeLabel')}:</div>
                      <div className="flex items-center gap-2">
                        <span className="badge badge-outline">
                          {fromPostalCode}
                        </span>
                        <span>→</span>
                        <span className="badge badge-outline">
                          {toPostalCode}
                        </span>
                      </div>
                    </div>
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
