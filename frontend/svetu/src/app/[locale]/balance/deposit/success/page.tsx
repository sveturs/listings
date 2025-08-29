import { Suspense } from 'react';
import SuccessClient from './SuccessClient';

export default async function SuccessPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          Loading...
        </div>
      }
    >
      <SuccessClient />
    </Suspense>
  );
}
