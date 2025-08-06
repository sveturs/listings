import { Suspense } from 'react';
import PaymentErrorClient from './PaymentErrorClient';

export default function PaymentErrorPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          Loading...
        </div>
      }
    >
      <PaymentErrorClient />
    </Suspense>
  );
}
