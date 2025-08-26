'use client';

import { useEffect } from 'react';
import { useSearchParams, useRouter } from 'next/navigation';
import { useTranslations } from 'next-intl';
import Link from 'next/link';
import {
  CheckCircleIcon,
  SparklesIcon,
  RocketLaunchIcon,
} from '@heroicons/react/24/solid';
import confetti from 'canvas-confetti';

export default function SubscriptionSuccessPage() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const t = useTranslations('subscription');
  const tCommon = useTranslations('common');

  const planCode = searchParams.get('plan');

  useEffect(() => {
    // Trigger confetti animation on mount
    const duration = 3 * 1000;
    const animationEnd = Date.now() + duration;
    const defaults = { startVelocity: 30, spread: 360, ticks: 60, zIndex: 0 };

    function randomInRange(min: number, max: number) {
      return Math.random() * (max - min) + min;
    }

    const interval: any = setInterval(function () {
      const timeLeft = animationEnd - Date.now();

      if (timeLeft <= 0) {
        return clearInterval(interval);
      }

      const particleCount = 50 * (timeLeft / duration);
      confetti(
        Object.assign({}, defaults, {
          particleCount,
          origin: { x: randomInRange(0.1, 0.3), y: Math.random() - 0.2 },
        })
      );
      confetti(
        Object.assign({}, defaults, {
          particleCount,
          origin: { x: randomInRange(0.7, 0.9), y: Math.random() - 0.2 },
        })
      );
    }, 250);

    return () => clearInterval(interval);
  }, []);

  const getPlanName = () => {
    if (!planCode) {
      return t('plans.starter'); // Default to starter if no plan specified
    }
    
    const planKey = `plans.${planCode}` as const;
    return t(planKey);
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-primary/10 via-base-200 to-secondary/10 flex items-center justify-center">
      <div className="card bg-base-100 shadow-2xl max-w-2xl w-full mx-4">
        <div className="card-body text-center p-8">
          <div className="relative">
            <CheckCircleIcon className="w-24 h-24 text-success mx-auto mb-6" />
            <SparklesIcon className="w-8 h-8 text-warning absolute top-0 right-1/4 animate-pulse" />
            <SparklesIcon className="w-6 h-6 text-secondary absolute bottom-2 left-1/4 animate-pulse delay-150" />
          </div>

          <h1 className="text-4xl font-bold mb-4 bg-gradient-to-r from-primary to-secondary bg-clip-text text-transparent">
            {t('success.welcomeTitle', { plan: getPlanName() })}
          </h1>

          <p className="text-lg text-base-content/70 mb-8">
            {t('success.description', { plan: getPlanName() })}
          </p>

          <div className="bg-base-200 rounded-lg p-6 mb-8">
            <h3 className="text-lg font-semibold mb-4">{t('success.whatsNext')}</h3>
            <div className="space-y-3 text-left max-w-md mx-auto">
              <div className="flex items-start gap-3">
                <div className="badge badge-primary badge-sm mt-1">1</div>
                <div>
                  <p className="font-medium">{t('success.step1.title')}</p>
                  <p className="text-sm text-base-content/60">
                    {t('success.step1.description')}
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="badge badge-primary badge-sm mt-1">2</div>
                <div>
                  <p className="font-medium">{t('success.step2.title')}</p>
                  <p className="text-sm text-base-content/60">
                    {t('success.step2.description')}
                  </p>
                </div>
              </div>
              <div className="flex items-start gap-3">
                <div className="badge badge-primary badge-sm mt-1">3</div>
                <div>
                  <p className="font-medium">{t('success.step3.title')}</p>
                  <p className="text-sm text-base-content/60">
                    {t('success.step3.description')}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <div className="card-actions justify-center gap-4">
            <Link
              href="/create-storefront"
              className="btn btn-primary btn-lg gap-2"
            >
              <RocketLaunchIcon className="w-5 h-5" />
              {t('success.createStorefront')}
            </Link>
            <Link href="/dashboard" className="btn btn-ghost btn-lg">
              {t('success.goToDashboard')}
            </Link>
          </div>
        </div>
      </div>
    </div>
  );
}
