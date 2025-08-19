import { Suspense } from 'react';
import ChatClient from './ChatClient';
import ErrorBoundaryClass from '@/components/ErrorBoundary';

export default function ChatPage() {
  return (
    <ErrorBoundaryClass name="ChatPage">
      <Suspense
        fallback={
          <div className="min-h-screen flex items-center justify-center">
            Loading...
          </div>
        }
      >
        <ChatClient />
      </Suspense>
    </ErrorBoundaryClass>
  );
}
