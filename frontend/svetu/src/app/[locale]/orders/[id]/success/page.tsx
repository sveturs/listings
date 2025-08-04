import { Suspense } from 'react';
import SuccessClient from './SuccessClient';

interface Props {
  params: Promise<{ id: string }>;
}

export default function SuccessPage({ params }: Props) {
  return (
    <Suspense fallback={<div className="min-h-screen flex items-center justify-center">Loading...</div>}>
      <SuccessClient params={params} />
    </Suspense>
  );
}