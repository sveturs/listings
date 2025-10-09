import { notFound } from 'next/navigation';
import CarDetailClient from './CarDetailClient';
import type { components } from '@/types/generated/api';

type C2CListing = components['schemas']['models.MarketplaceListing'];

interface CarDetailPageProps {
  params: Promise<{
    locale: string;
    id: string;
  }>;
}

async function getCarDetails(id: string): Promise<C2CListing | null> {
  try {
    // Server Component - используем прямой fetch к backend
    const backendUrl =
      process.env.BACKEND_INTERNAL_URL || 'http://localhost:3000';
    const response = await fetch(`${backendUrl}/api/v1/c2c/listings/${id}`, {
      cache: 'no-store', // Всегда получаем свежие данные
    });

    if (!response.ok) {
      return null;
    }

    const data = await response.json();
    if (data?.data) {
      return data.data;
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
