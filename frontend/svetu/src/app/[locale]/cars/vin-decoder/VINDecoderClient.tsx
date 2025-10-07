'use client';

import { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { ArrowLeft, Car, Info } from 'lucide-react';
import { VINDecoder } from '@/components/marketplace/VINDecoder';
import type { components } from '@/types/generated/api';

type VINDecodeResult =
  components['schemas']['models.VINDecodeResult'];

interface VINDecoderClientProps {
  locale: string;
}

export default function VINDecoderClient({ locale }: VINDecoderClientProps) {
  const t = useTranslations('cars');
  const tCommon = useTranslations('common');
  const router = useRouter();
  const [decodedResult, setDecodedResult] = useState<VINDecodeResult | null>(
    null
  );

  const handleVINDecoded = (result: VINDecodeResult) => {
    setDecodedResult(result);
  };

  const handleUseInListing = () => {
    // Store the decoded data in sessionStorage to use in listing creation
    if (decodedResult) {
      sessionStorage.setItem('vinDecodedData', JSON.stringify(decodedResult));
      router.push(`/${locale}/create-listing-choice`);
    }
  };

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-base-100 shadow-lg">
        <div className="container mx-auto px-4 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <Link href={`/${locale}/cars`} className="btn btn-ghost btn-sm">
                <ArrowLeft className="w-4 h-4" />
                {tCommon('back')}
              </Link>
              <div>
                <h1 className="text-2xl font-bold flex items-center gap-2">
                  <Car className="w-6 h-6" />
                  {t('vinDecoder')}
                </h1>
                <p className="text-sm opacity-70">
                  {t('vinDecoderDescription')}
                </p>
              </div>
            </div>
          </div>
        </div>
      </div>

      {/* Main Content */}
      <div className="container mx-auto px-4 py-8">
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* VIN Decoder */}
          <div className="lg:col-span-2">
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <VINDecoder onVINDecoded={handleVINDecoded} />

                {decodedResult && (
                  <div className="mt-6 flex justify-end gap-2">
                    <button
                      onClick={handleUseInListing}
                      className="btn btn-primary"
                    >
                      {t('useInListing')}
                    </button>
                  </div>
                )}
              </div>
            </div>
          </div>

          {/* Info Sidebar */}
          <div className="space-y-4">
            {/* What is VIN? */}
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title text-lg">
                  <Info className="w-5 h-5" />
                  {t('whatIsVIN')}
                </h3>
                <p className="text-sm opacity-80">{t('vinExplanation')}</p>
              </div>
            </div>

            {/* Where to find VIN? */}
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title text-lg">{t('whereToFindVIN')}</h3>
                <ul className="text-sm space-y-1 opacity-80">
                  <li>• {t('vinLocation1')}</li>
                  <li>• {t('vinLocation2')}</li>
                  <li>• {t('vinLocation3')}</li>
                  <li>• {t('vinLocation4')}</li>
                </ul>
              </div>
            </div>

            {/* Benefits */}
            <div className="card bg-base-100 shadow-xl">
              <div className="card-body">
                <h3 className="card-title text-lg">{t('vinBenefits')}</h3>
                <ul className="text-sm space-y-1 opacity-80">
                  <li>✓ {t('vinBenefit1')}</li>
                  <li>✓ {t('vinBenefit2')}</li>
                  <li>✓ {t('vinBenefit3')}</li>
                  <li>✓ {t('vinBenefit4')}</li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
