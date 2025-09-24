'use client';

import { PublicEnvScript } from 'next-runtime-env';
import { useEffect, useState } from 'react';

export default function EnvScriptWrapper() {
  const [isMounted, setIsMounted] = useState(false);

  useEffect(() => {
    setIsMounted(true);
  }, []);

  // Рендерим только на клиенте
  if (!isMounted) {
    return null;
  }

  return <PublicEnvScript />;
}
