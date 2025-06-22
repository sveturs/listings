import { Suspense } from 'react';
import SearchPage from './SearchPage';

export default function Page() {
  return (
    <Suspense
      fallback={
        <div className="min-h-screen bg-base-200 flex items-center justify-center">
          <span className="loading loading-spinner loading-lg text-primary"></span>
        </div>
      }
    >
      <SearchPage />
    </Suspense>
  );
}
