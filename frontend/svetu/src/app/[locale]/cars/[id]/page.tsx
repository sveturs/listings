import { notFound } from 'next/navigation';
import api from '@/services/api';
import CarDetailClient from './CarDetailClient';
import type { components } from '@/types/generated/api';

type MarketplaceListing =
  components['schemas']['backend_internal_domain_models.MarketplaceListing'];

interface CarDetailPageProps {
  params: Promise<{
    locale: string;
    id: string;
  }>;
}

async function getCarDetails(id: string): Promise<MarketplaceListing | null> {
  try {
    const response = await api.get(`/api/v1/marketplace/listings/${id}`);
    if (response.data?.data) {
      return response.data.data;
    }
    return null;
  } catch (error) {
    console.error('Error fetching car details:', error);
    return null;
  }
}

export default async function CarDetailPage({ params }: CarDetailPageProps) {
  const { id, locale } = await params;
  const car = await getCarDetails(id);

  if (!car) {
    notFound();
  }

  // Проверяем, что это автомобиль
  if (![1003, 1301, 1303].includes(car.category_id || 0)) {
    notFound();
  }

  return <CarDetailClient car={car} locale={locale} />;
}
