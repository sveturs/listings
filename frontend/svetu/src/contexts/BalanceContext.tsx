'use client';

import React, { createContext, useContext } from 'react';
import { useBalance } from '@/hooks/useBalance';

interface BalanceContextType {
  refreshBalance: () => Promise<void>;
}

const BalanceContext = createContext<BalanceContextType | null>(null);

export const BalanceProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const { refreshBalance } = useBalance();

  // Убираем автоматическое обновление - пусть это контролируют компоненты
  // при необходимости можно добавить периодическое обновление только на определенных страницах

  return (
    <BalanceContext.Provider value={{ refreshBalance }}>
      {children}
    </BalanceContext.Provider>
  );
};

export const useBalanceContext = () => {
  const context = useContext(BalanceContext);
  if (!context) {
    throw new Error('useBalanceContext must be used within BalanceProvider');
  }
  return context;
};
