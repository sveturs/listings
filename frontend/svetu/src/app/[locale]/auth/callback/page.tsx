import { Suspense } from 'react';
import CallbackClient from './CallbackClient';

export default async function CallbackPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          Loading...
        </div>
      }
    >
      <CallbackClient />
    </Suspense>
  );
}
