import { getTranslations } from 'next-intl/server';
import FavoritesClient from './FavoritesClient';

export async function generateMetadata({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'favorites' });
  return {
    title: t('pageTitle'),
    description: t('pageDescription'),
  };
}

export default async function FavoritesPage({
  params,
}: {
  params: Promise<{ locale: string }>;
}) {
  const { locale } = await params;
  const t = await getTranslations({ locale, namespace: 'favorites' });

  return (
    <FavoritesClient
      locale={locale}
      translations={{
        title: t('title'),
        description: t('description'),
        emptyTitle: t('emptyTitle'),
        emptyDescription: t('emptyDescription'),
        browseListings: t('browseListings'),
        removeFromFavorites: t('removeFromFavorites'),
        viewDetails: t('viewDetails'),
        contactSeller: t('contactSeller'),
        addToCart: t('addToCart'),
        loading: t('loading'),
      }}
    />
  );
}
