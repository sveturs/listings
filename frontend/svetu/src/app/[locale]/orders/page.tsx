'use client';

import { redirect } from 'next/navigation';
import { useParams } from 'next/navigation';

export default function OrdersPage() {
  const params = useParams();
  const locale = params.locale as string;

  redirect(`/${locale}/profile/orders`);
}
