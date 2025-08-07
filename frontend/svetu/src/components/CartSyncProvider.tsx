'use client';

import { useCartSync } from '@/hooks/useCartSync';

/**
 * Компонент-провайдер для синхронизации корзин пользователя
 * Должен быть внутри AuthProvider и ReduxProvider
 */
export function CartSyncProvider({ children }: { children: React.ReactNode }) {
  useCartSync();
  return <>{children}</>;
}
