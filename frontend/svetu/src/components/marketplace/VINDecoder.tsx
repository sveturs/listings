'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { Search, Info, Car, Loader2 } from 'lucide-react';
import { CarsService } from '@/services/CarsService';
import type { components } from '@/types/generated/api';

type VINDecodeResult =
  components['schemas']['models.VINDecodeResult'];

interface VINDecoderProps {
  onVINDecoded?: (result: VINDecodeResult) => void;
  className?: string;
}

/**
 * VIN Decoder component for decoding Vehicle Identification Numbers
 * Provides a user-friendly interface to decode VIN and display vehicle information
 */
export const VINDecoder: React.FC<VINDecoderProps> = ({
  onVINDecoded,
  className = '',
}) => {
  const t = useTranslations('cars');
  // const tCommon = useTranslations('common');

  const [vin, setVIN] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<VINDecodeResult | null>(null);
  const [showExample, setShowExample] = useState(false);

  // Example VINs for testing
  const exampleVINs = [
    'WVWZZZ1JZ3W386752', // VW Golf
    '1HGBH41JXMN109186', // Honda Civic
    'WBAVB13566PT22180', // BMW 3 Series
  ];

  const handleDecode = async () => {
    if (!vin || vin.length !== 17) {
      setError(t('vinInvalidLength'));
      return;
    }

    setLoading(true);
    setError(null);
    setResult(null);

    const response = await CarsService.decodeVIN(vin);

    setLoading(false);

    if (response.success && response.data) {
      setResult(response.data);
      if (onVINDecoded) {
        onVINDecoded(response.data);
      }
    } else {
      setError(response.error || t('vinDecodeError'));
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === 'Enter' && vin.length === 17) {
      handleDecode();
    }
  };

  const handleExampleClick = (exampleVIN: string) => {
    setVIN(exampleVIN);
    setShowExample(false);
  };

  const formatValue = (value: any): string => {
    if (value === null || value === undefined) return '-';
    if (typeof value === 'boolean') return value ? t('yes') : t('no');
    if (typeof value === 'number') return value.toString();
    return value.toString();
  };

  return (
    <div className={`space-y-4 ${className}`}>
      {/* Input Section */}
      <div className="form-control">
        <label className="label">
          <span className="label-text font-semibold">{t('enterVIN')}</span>
          <button
            type="button"
            onClick={() => setShowExample(!showExample)}
            className="label-text-alt link link-primary"
          >
            {t('showExamples')}
          </button>
        </label>

        <div className="input-group">
          <input
            type="text"
            value={vin}
            onChange={(e) => {
              setVIN(e.target.value.toUpperCase());
              setError(null);
            }}
            onKeyPress={handleKeyPress}
            placeholder="e.g., WVWZZZ1JZ3W386752"
            className={`input input-bordered flex-1 font-mono ${
              error ? 'input-error' : ''
            }`}
            maxLength={17}
            disabled={loading}
          />
          <button
            onClick={handleDecode}
            disabled={loading || vin.length !== 17}
            className="btn btn-primary"
          >
            {loading ? (
              <Loader2 className="w-5 h-5 animate-spin" />
            ) : (
              <Search className="w-5 h-5" />
            )}
            <span className="hidden sm:inline ml-2">{t('decode')}</span>
          </button>
        </div>

        <label className="label">
          <span className="label-text-alt">
            {vin.length}/17 {t('characters')}
          </span>
          {error && <span className="label-text-alt text-error">{error}</span>}
        </label>
      </div>

      {/* Example VINs */}
      {showExample && (
        <div className="alert alert-info">
          <Info className="w-5 h-5" />
          <div className="flex-1">
            <h4 className="font-semibold">{t('exampleVINs')}</h4>
            <div className="mt-2 space-y-1">
              {exampleVINs.map((exampleVIN, index) => (
                <button
                  key={index}
                  onClick={() => handleExampleClick(exampleVIN)}
                  className="block text-left hover:underline font-mono text-sm"
                >
                  {exampleVIN}
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Loading State */}
      {loading && (
        <div className="flex justify-center py-8">
          <div className="loading loading-spinner loading-lg"></div>
        </div>
      )}

      {/* Results */}
      {result && (
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body">
            <h3 className="card-title">
              <Car className="w-5 h-5" />
              {t('vehicleInformation')}
            </h3>

            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 mt-4">
              {/* Basic Information */}
              <div className="space-y-2">
                <h4 className="font-semibold text-sm opacity-70">
                  {t('basicInfo')}
                </h4>

                {result.make_name && (
                  <div className="flex justify-between">
                    <span className="text-sm opacity-70">{t('make')}:</span>
                    <span className="font-medium">{result.make_name}</span>
                  </div>
                )}

                {result.model_name && (
                  <div className="flex justify-between">
                    <span className="text-sm opacity-70">{t('model')}:</span>
                    <span className="font-medium">{result.model_name}</span>
                  </div>
                )}

                {result.year && (
                  <div className="flex justify-between">
                    <span className="text-sm opacity-70">{t('year')}:</span>
                    <span className="font-medium">{result.year}</span>
                  </div>
                )}

                {result.body_type && (
                  <div className="flex justify-between">
                    <span className="text-sm opacity-70">{t('bodyType')}:</span>
                    <span className="font-medium">{result.body_type}</span>
                  </div>
                )}
              </div>

              {/* Technical Information */}
              <div className="space-y-2">
                <h4 className="font-semibold text-sm opacity-70">
                  {t('technicalInfo')}
                </h4>

                {result.engine?.type && (
                  <div className="flex justify-between">
                    <span className="text-sm opacity-70">{t('engine')}:</span>
                    <span className="font-medium">
                      {result.engine.type}
                      {result.engine.displacement &&
                        ` ${result.engine.displacement}L`}
                      {result.engine.power_hp && ` ${result.engine.power_hp}HP`}
                    </span>
                  </div>
                )}

                {result.transmission && (
                  <div className="flex justify-between">
                    <span className="text-sm opacity-70">
                      {t('transmission')}:
                    </span>
                    <span className="font-medium">{result.transmission}</span>
                  </div>
                )}

                {result.drive_type && (
                  <div className="flex justify-between">
                    <span className="text-sm opacity-70">{t('drive')}:</span>
                    <span className="font-medium">{result.drive_type}</span>
                  </div>
                )}

                {result.fuel_type && (
                  <div className="flex justify-between">
                    <span className="text-sm opacity-70">{t('fuelType')}:</span>
                    <span className="font-medium">{result.fuel_type}</span>
                  </div>
                )}
              </div>
            </div>

            {/* Additional Details */}
            {result.raw_data && Object.keys(result.raw_data).length > 0 && (
              <div className="mt-6">
                <h4 className="font-semibold text-sm opacity-70 mb-2">
                  {t('additionalDetails')}
                </h4>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-2">
                  {Object.entries(result.raw_data).map(([key, value]) => (
                    <div key={key} className="flex justify-between">
                      <span className="text-sm opacity-70">{key}:</span>
                      <span className="font-medium text-sm">
                        {formatValue(value)}
                      </span>
                    </div>
                  ))}
                </div>
              </div>
            )}

            {/* Action Buttons */}
            <div className="card-actions justify-end mt-4">
              <button
                onClick={() => {
                  setResult(null);
                  setVIN('');
                }}
                className="btn btn-ghost btn-sm"
              >
                {t('clearResults')}
              </button>
              {onVINDecoded && (
                <button className="btn btn-primary btn-sm">
                  {t('useInListing')}
                </button>
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default VINDecoder;
