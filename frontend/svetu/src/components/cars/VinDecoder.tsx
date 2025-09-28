'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useAppDispatch } from '@/store/hooks';

interface VINDecodeResult {
  valid: boolean;
  vin: string;
  year?: number;
  make_name?: string;
  model_name?: string;
  trim?: string;
  body_type?: string;
  fuel_type?: string;
  transmission?: string;
  drive_type?: string;
  engine?: {
    displacement?: number;
    cylinders?: number;
    fuel_injection?: string;
  };
  manufacturer?: {
    name?: string;
    country?: string;
  };
}

interface VinDecoderProps {
  onDecode?: (result: VINDecodeResult) => void;
  onAutoFill?: (data: any) => void;
  showAutoFill?: boolean;
}

export default function VinDecoder({
  onDecode,
  onAutoFill,
  showAutoFill = true,
}: VinDecoderProps) {
  const t = useTranslations('cars');
  const _dispatch = useAppDispatch();
  const [vin, setVin] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<VINDecodeResult | null>(null);

  const validateVIN = (value: string): boolean => {
    // Basic VIN validation: 17 characters, no I, O, Q
    const vinRegex = /^[A-HJ-NPR-Z0-9]{17}$/i;
    return vinRegex.test(value);
  };

  const handleVinChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value.toUpperCase();
    setVin(value);
    setError(null);

    if (value.length === 17 && !validateVIN(value)) {
      setError(t('vinDecoder.invalidFormat'));
    }
  };

  const decodeVIN = async () => {
    if (!validateVIN(vin)) {
      setError(t('vinDecoder.invalidFormat'));
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const response = await fetch(`/api/v1/cars/vin/${vin}/decode`);
      const data = await response.json();

      if (data.success && data.data) {
        setResult(data.data);
        onDecode?.(data.data);
      } else {
        setError(data.message || t('vinDecoder.decodeFailed'));
      }
    } catch {
      setError(t('vinDecoder.networkError'));
    } finally {
      setLoading(false);
    }
  };

  const handleAutoFill = () => {
    if (!result) return;

    const autoFillData = {
      year: result.year,
      make: result.make_name,
      model: result.model_name,
      trim: result.trim,
      body_type: result.body_type,
      fuel_type: result.fuel_type,
      transmission: result.transmission,
      drive_type: result.drive_type,
      engine_displacement: result.engine?.displacement,
      engine_cylinders: result.engine?.cylinders,
    };

    onAutoFill?.(autoFillData);
  };

  const getTransmissionLabel = (transmission?: string) => {
    if (!transmission) return '-';
    const key = `filters.transmission.${transmission.toLowerCase()}`;
    return t(key);
  };

  const getFuelTypeLabel = (fuelType?: string) => {
    if (!fuelType) return '-';
    const key = `filters.fuel_type.${fuelType.toLowerCase()}`;
    return t(key);
  };

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body">
        <h3 className="card-title">
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
              d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"
            />
          </svg>
          {t('vinDecoder.title')}
        </h3>

        <div className="form-control">
          <label className="label">
            <span className="label-text">{t('vinDecoder.enterVin')}</span>
            <span className="label-text-alt">{t('vinDecoder.vinFormat')}</span>
          </label>

          <div className="input-group">
            <input
              type="text"
              placeholder="WBABA91060AL04320"
              className={`input input-bordered w-full ${error ? 'input-error' : ''}`}
              value={vin}
              onChange={handleVinChange}
              maxLength={17}
              disabled={loading}
            />
            <button
              className={`btn btn-primary ${loading ? 'loading' : ''}`}
              onClick={decodeVIN}
              disabled={vin.length !== 17 || loading || !!error}
            >
              {loading ? t('common.loading') : t('vinDecoder.decode')}
            </button>
          </div>

          {error && (
            <label className="label">
              <span className="label-text-alt text-error">{error}</span>
            </label>
          )}

          <div className="text-sm text-base-content/60 mt-2">
            {t('vinDecoder.hint')}
          </div>
        </div>

        {result && (
          <div className="mt-6 space-y-4">
            <div className="divider">{t('vinDecoder.decodedInfo')}</div>

            <div className="grid grid-cols-2 gap-4">
              <div>
                <div className="text-sm text-base-content/60">
                  {t('vinDecoder.year')}
                </div>
                <div className="font-semibold">{result.year || '-'}</div>
              </div>

              <div>
                <div className="text-sm text-base-content/60">
                  {t('vinDecoder.make')}
                </div>
                <div className="font-semibold">{result.make_name || '-'}</div>
              </div>

              <div>
                <div className="text-sm text-base-content/60">
                  {t('vinDecoder.model')}
                </div>
                <div className="font-semibold">{result.model_name || '-'}</div>
              </div>

              <div>
                <div className="text-sm text-base-content/60">
                  {t('vinDecoder.trim')}
                </div>
                <div className="font-semibold">{result.trim || '-'}</div>
              </div>

              <div>
                <div className="text-sm text-base-content/60">
                  {t('vinDecoder.bodyType')}
                </div>
                <div className="font-semibold">{result.body_type || '-'}</div>
              </div>

              <div>
                <div className="text-sm text-base-content/60">
                  {t('vinDecoder.fuelType')}
                </div>
                <div className="font-semibold">
                  {getFuelTypeLabel(result.fuel_type)}
                </div>
              </div>

              <div>
                <div className="text-sm text-base-content/60">
                  {t('vinDecoder.transmission')}
                </div>
                <div className="font-semibold">
                  {getTransmissionLabel(result.transmission)}
                </div>
              </div>

              <div>
                <div className="text-sm text-base-content/60">
                  {t('vinDecoder.driveType')}
                </div>
                <div className="font-semibold">{result.drive_type || '-'}</div>
              </div>
            </div>

            {result.engine && (
              <>
                <div className="divider">{t('vinDecoder.engineInfo')}</div>
                <div className="grid grid-cols-2 gap-4">
                  {result.engine.displacement && (
                    <div>
                      <div className="text-sm text-base-content/60">
                        {t('vinDecoder.displacement')}
                      </div>
                      <div className="font-semibold">
                        {result.engine.displacement} L
                      </div>
                    </div>
                  )}
                  {result.engine.cylinders && (
                    <div>
                      <div className="text-sm text-base-content/60">
                        {t('vinDecoder.cylinders')}
                      </div>
                      <div className="font-semibold">
                        {result.engine.cylinders}
                      </div>
                    </div>
                  )}
                </div>
              </>
            )}

            {result.manufacturer && (
              <>
                <div className="divider">
                  {t('vinDecoder.manufacturerInfo')}
                </div>
                <div className="grid grid-cols-2 gap-4">
                  {result.manufacturer.name && (
                    <div>
                      <div className="text-sm text-base-content/60">
                        {t('vinDecoder.manufacturer')}
                      </div>
                      <div className="font-semibold">
                        {result.manufacturer.name}
                      </div>
                    </div>
                  )}
                  {result.manufacturer.country && (
                    <div>
                      <div className="text-sm text-base-content/60">
                        {t('vinDecoder.country')}
                      </div>
                      <div className="font-semibold">
                        {result.manufacturer.country}
                      </div>
                    </div>
                  )}
                </div>
              </>
            )}

            {showAutoFill && onAutoFill && (
              <div className="card-actions justify-end mt-6">
                <button className="btn btn-primary" onClick={handleAutoFill}>
                  <svg
                    className="w-5 h-5 mr-2"
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
                  {t('vinDecoder.autoFill')}
                </button>
              </div>
            )}
          </div>
        )}
      </div>
    </div>
  );
}
