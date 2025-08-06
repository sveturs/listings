import { Suspense } from 'react';
import Header from './Header';

export default function HeaderWrapper() {
  return (
    <Suspense fallback={<div className="navbar bg-base-100 shadow-lg h-16" />}>
      <Header />
    </Suspense>
  );
}
