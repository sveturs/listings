'use client';

import { useRouter } from '@/i18n/routing';
import { useEffect } from 'react';

export default function TestEditPage() {
  const router = useRouter();

  useEffect(() => {
    // Redirect to edit page with a test ID
    router.push('/profile/listings/35/edit');
  }, [router]);

  return (
    <div className="container mx-auto px-4 py-8">
      <div className="flex justify-center">
        <span className="loading loading-spinner loading-lg"></span>
      </div>
    </div>
  );
}
