import { Suspense } from 'react';
import MapClient from './MapClient';

export default function MapPage() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen flex items-center justify-center">
          Loading...
        </div>
      }
    >
      <MapClient />
    </Suspense>
  );
}
