'use client';

import React from 'react';
import { useLocale, useTranslations } from 'next-intl';
import Link from 'next/link';
import { motion } from 'framer-motion';
import { PageTransition } from '@/components/ui/PageTransition';

export default function CheckoutSuccessPage() {
  const t = useTranslations('checkout');
  const tOrders = useTranslations('orders');
  const locale = useLocale();

  return (
    <PageTransition>
      <div className="min-h-screen bg-base-100 pt-24">
        <div className="container mx-auto px-4 py-16">
          <div className="max-w-md mx-auto text-center">
            <motion.div
              initial={{ scale: 0 }}
              animate={{ scale: 1 }}
              transition={{ duration: 0.5, type: 'spring' }}
              className="mb-8"
            >
              <div className="inline-flex items-center justify-center w-32 h-32 bg-success/10 rounded-full">
                <svg
                  className="w-16 h-16 text-success"
                  fill="none"
                  viewBox="0 0 24 24"
                  stroke="currentColor"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth={2}
                    d="M5 13l4 4L19 7"
                  />
                </svg>
              </div>
            </motion.div>

            <h1 className="text-3xl font-bold mb-4">{t('success.title')}</h1>
            <p className="text-lg text-base-content/70 mb-8">
              {t('success.description')}
            </p>

            <div className="card bg-base-200 mb-8">
              <div className="card-body">
                <h2 className="card-title justify-center">
                  {t('success.whatNext')}
                </h2>
                <div className="text-left space-y-3 mt-4">
                  <div className="flex gap-3">
                    <span className="text-2xl">📧</span>
                    <div>
                      <p className="font-medium">
                        {t('success.emailConfirmation')}
                      </p>
                      <p className="text-sm text-base-content/60">
                        {t('success.emailConfirmationDesc')}
                      </p>
                    </div>
                  </div>
                  <div className="flex gap-3">
                    <span className="text-2xl">📦</span>
                    <div>
                      <p className="font-medium">{t('success.shipping')}</p>
                      <p className="text-sm text-base-content/60">
                        {t('success.shippingDesc')}
                      </p>
                    </div>
                  </div>
                  <div className="flex gap-3">
                    <span className="text-2xl">💬</span>
                    <div>
                      <p className="font-medium">{t('success.contact')}</p>
                      <p className="text-sm text-base-content/60">
                        {t('success.contactDesc')}
                      </p>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <div className="flex flex-col sm:flex-row gap-4">
              <Link
                href={`/${locale}/orders`}
                className="btn btn-primary flex-1"
              >
                {tOrders('myOrders')}
              </Link>
              <Link href={`/${locale}`} className="btn btn-outline flex-1">
                {t('continueShopping')}
              </Link>
            </div>
          </div>
        </div>
      </div>
    </PageTransition>
  );
}
