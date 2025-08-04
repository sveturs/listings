import { Suspense } from 'react';
import MockClient from './MockClient';

export default function MockPage() {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center">Loading...</div>}>
      <MockClient />
    </Suspense>
  );
}