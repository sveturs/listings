import { Suspense } from 'react';
import PaymentFailedClient from './PaymentFailedClient';

interface Props {
  params: Promise<{ id: string }>;
}

export default function PaymentFailedPage({ params }: Props) {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          Loading...
        </div>
      }
    >
      <PaymentFailedClient params={params} />
    </Suspense>
  );
}
