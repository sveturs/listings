'use client';

import { useState, useEffect } from 'react';
import { useTranslations, useLocale } from 'next-intl';
import Link from 'next/link';
import { ArrowLeftIcon } from '@heroicons/react/24/outline';
import { EditProductProvider } from '@/contexts/EditProductContext';
import { useEditProduct } from '@/contexts/EditProductContext';
import EditProductWizard from '@/components/products/EditProductWizard';
import { storefrontProductsService } from '@/services/storefrontProducts';

interface PageProps {
  params: Promise<{
    locale: string;
    slug: string;
    id: string;
  }>;
}

function EditProductPageContent({
  slug,
  productId,
}: {
  slug: string;
  productId: number;
}) {
  const t = useTranslations('storefronts.products');
  const tStorefronts.products.errors = useTranslations('storefronts.products.errors');
  const { state, loadProduct, dispatch } = useEditProduct();
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const fetchProduct = async () => {
      try {
        setIsLoading(true);
        dispatch({ type: 'SET_LOADING', payload: true });

        const product = await storefrontProductsService.getProduct(
          slug,
          productId
        );
        loadProduct(product);
      } catch (error: any) {
        console.error('Error loading product:', error);
        setError(error.message || tStorefronts.products.errors('loadFailed'));
      } finally {
        setIsLoading(false);
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    };

    fetchProduct();
  }, [slug, productId, loadProduct, dispatch, t]);

  if (isLoading) {
    return (
      <div className="container mx-auto px-4 py-8 flex justify-center">
        <div className="text-center">
          <span className="loading loading-spinner loading-lg"></span>
          <p className="mt-4 text-base-content/70">
            {t('loadingProduct')}
          </p>
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="bg-error/10 border border-error rounded-2xl p-6 text-center">
          <h3 className="text-lg font-semibold text-error mb-2">
            {tStorefronts.products.errors('loadFailed')}
          </h3>
          <p className="text-error/80 mb-4">{error}</p>
          <Link
            href={`/storefronts/${slug}/products`}
            className="btn btn-primary"
          >
            {t('backToProducts')}
          </Link>
        </div>
      </div>
    );
  }

  if (!state.originalProduct) {
    return (
      <div className="container mx-auto px-4 py-8">
        <div className="bg-error/10 border border-error rounded-2xl p-6 text-center">
          <h3 className="text-lg font-semibold text-error mb-2">
            {tStorefronts.products.errors('productNotFound')}
          </h3>
          <Link
            href={`/storefronts/${slug}/products`}
            className="btn btn-primary"
          >
            {t('backToProducts')}
          </Link>
        </div>
      </div>
    );
  }

  return <EditProductWizard storefrontSlug={slug} productId={productId} />;
}

export default function EditProductPage({ params }: PageProps) {
  const [slug, setSlug] = useState<string>('');
  const [productId, setProductId] = useState<number>(0);
  const t = useTranslations('storefronts.products');
  const locale = useLocale();

  useEffect(() => {
    params.then((p) => {
      setSlug(p.slug);
      setProductId(parseInt(p.id, 10));
    });
  }, [params]);

  if (!slug || !productId) {
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
            href={`/${locale}/storefronts/${slug}/products`}
            className="inline-flex items-center text-primary hover:underline mb-4"
          >
            <ArrowLeftIcon className="w-4 h-4 mr-2" />
            {t('backToProducts')}
          </Link>
          <h1 className="text-3xl font-bold text-base-content">
            {t('editProduct')}
          </h1>
          <p className="text-base-content/70 mt-2">
            {t('editProductDescription')}
          </p>
        </div>
      </div>

      {/* Wizard */}
      <div className="container mx-auto px-4 py-8">
        <EditProductProvider>
          <EditProductPageContent slug={slug} productId={productId} />
        </EditProductProvider>
      </div>
    </div>
  );
}
