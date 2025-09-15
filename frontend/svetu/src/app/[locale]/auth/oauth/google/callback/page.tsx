import { Suspense } from 'react';
import OAuthProcessor from './OAuthProcessor';

interface PageProps {
  searchParams: Promise<{ [key: string]: string | string[] | undefined }>;
}

export default async function GoogleCallbackPage({ searchParams }: PageProps) {
  // Extract token from searchParams (server-side) - for logging only
  const params = await searchParams;
  const token = params.token as string | undefined;
  const error = params.error as string | undefined;
  const code = params.code as string | undefined;
  const state = params.state as string | undefined;

  console.log('[GoogleCallbackPage Server] Params received:', {
    tokenPresent: !!token,
    tokenLength: token?.length,
    error,
    codePresent: !!code,
    statePresent: !!state,
  });

  return (
    <Suspense
      fallback={
        <div className="flex items-center justify-center min-h-screen bg-base-200">
          <div className="loading loading-spinner loading-lg text-primary"></div>
        </div>
      }
    >
      <OAuthProcessor />
    </Suspense>
  );
}
