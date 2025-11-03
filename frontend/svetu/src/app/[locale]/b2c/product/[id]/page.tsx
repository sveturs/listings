'use client';

import { useEffect } from 'react';
import { useParams, useRouter } from 'next/navigation';
import config from '@/config';

export default function B2CProductRedirect() {
  const params = useParams();
  const router = useRouter();
  const locale = params?.locale as string;
  const productId = params?.id as string;

  useEffect(() => {
    // Загружаем информацию о товаре чтобы получить slug витрины
    const apiUrl = config.getApiUrl();
    fetch(`${apiUrl}/api/v1/storefronts/products/${productId}`, {
      credentials: 'include',
      headers: {
        'Content-Type': 'application/json',
      },
    })
      .then((res) => res.json())
      .then((result) => {
        const data = result.data || result;
        if (data && data.storefront) {
          // Редирект на правильный URL со slug
          router.replace(
            `/${locale}/b2c/${data.storefront.slug}/products/${productId}`
          );
        } else {
          // Fallback - используем дефолтный slug
          router.replace(`/${locale}/b2c/store/products/${productId}`);
        }
      })
      .catch((err) => {
        console.error('Error loading product:', err);
        // В случае ошибки редиректим на главную страницу витрин
        router.replace(`/${locale}/b2c`);
      });
  }, [productId, locale, router]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <div className="loading loading-spinner loading-lg"></div>
        <p className="mt-4">Loading product...</p>
      </div>
    </div>
  );
}
