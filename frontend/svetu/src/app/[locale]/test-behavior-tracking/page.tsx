/**
 * Тестовая страница для проверки работы поведенческого трекинга
 */

import { Suspense } from 'react';
import { ExampleBehaviorTrackingComponent } from '@/hooks/__tests__/useBehaviorTracking.example';

export default function TestBehaviorTrackingPage() {
  return (
    <div className="container mx-auto py-8">
      <Suspense
        fallback={
          <div className="loading loading-spinner loading-lg mx-auto"></div>
        }
      >
        <ExampleBehaviorTrackingComponent userId="test-user-123" />
      </Suspense>
    </div>
  );
}
