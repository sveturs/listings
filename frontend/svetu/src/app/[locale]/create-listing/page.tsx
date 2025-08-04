import { Suspense } from 'react';
import CreateListingClient from './CreateListingClient';

export default function CreateListingPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center">Loading...</div>}>
      <CreateListingClient />
    </Suspense>
  );
}