'use client';

import { useAuth } from '@/contexts/AuthContext';
import { useRouter } from '@/i18n/routing';
import { useEffect, useState } from 'react';

interface AdminGuardProps {
  children: React.ReactNode;
  loading?: React.ReactNode;
}

export default function AdminGuard({ children, loading }: AdminGuardProps) {
  // TEMPORARY: Allow access for testing purposes
  // TODO: Restore authorization check after testing
  return <>{children}</>;
}
