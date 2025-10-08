'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import Link from 'next/link';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import { CreateAIProductProvider } from '@/contexts/CreateAIProductContext';
import AIProductWizard from '@/components/b2c/ai/AIProductWizard';

interface PageProps {
  params: Promise<{
    locale: string;
    slug: string;
  }>;
}

export default function NewAIProductPage({ params }: PageProps) {
  const [slug, setSlug] = useState<string>('');
  const t = useTranslations('storefronts');
  const tCommon = useTranslations('common');
  const locale = useLocale();

  useEffect(() => {
    params.then((p) => setSlug(p.slug));
  }, [params]);

  if (!slug) {
    return (
      <div className="container mx-auto px-4 py-8 flex justify-center">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-base-200">
      {/* Header */}
      <div className="bg-base-100 shadow-sm border-b">
        <div className="container mx-auto px-4 py-6">
          <Link
            href={`/${locale}/b2c/${slug}/products`}
            className="inline-flex items-center text-primary hover:underline mb-4"
          >
            <ArrowLeftIcon className="w-4 h-4 mr-2" />
            {tCommon('back')}
          </Link>
          <div className="flex items-center gap-3">
            <div className="p-2 bg-primary/10 rounded-lg">
              <svg
                className="w-6 h-6 text-primary"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  strokeWidth={2}
                  d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"
                />
              </svg>
            </div>
            <div>
              <h1 className="text-3xl font-bold text-base-content">
                {t('addNewProductAI') || 'Create Product with AI'}
              </h1>
              <p className="text-base-content/70 mt-1">
                {t('aiWizardDescription') ||
                  'Upload photos and let AI do the rest'}
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* Wizard */}
      <div className="container mx-auto px-4 py-8">
        <CreateAIProductProvider>
          <AIProductWizard storefrontSlug={slug} />
        </CreateAIProductProvider>
      </div>
    </div>
  );
}
