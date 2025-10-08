'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import Link from 'next/link';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import { CreateProductProvider } from '@/contexts/CreateProductContext';
import ProductWizard from '@/components/products/ProductWizard';

interface PageProps {
  params: Promise<{
    locale: string;
    slug: string;
  }>;
}

export default function NewProductPage({ params }: PageProps) {
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
          <h1 className="text-3xl font-bold text-base-content">
            {t('addNewProduct')}
          </h1>
          <p className="text-base-content/70 mt-2">{t('wizardDescription')}</p>
        </div>
      </div>

      {/* Wizard */}
      <div className="container mx-auto px-4 py-8">
        <CreateProductProvider>
          <ProductWizard storefrontSlug={slug} />
        </CreateProductProvider>
      </div>
    </div>
  );
}
