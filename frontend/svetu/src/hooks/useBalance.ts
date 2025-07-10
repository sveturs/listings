import { useState, useEffect, useCallback } from 'react';
import { balanceService } from '@/services/balance';
import { toast } from 'react-hot-toast';
import type { components } from '@/types/generated/api';

type UserBalance =
  components['schemas']['backend_internal_domain_models.UserBalance'];
type Transaction =
  components['schemas']['backend_internal_domain_models.BalanceTransaction'];

interface UseBalanceReturn {
  balance: UserBalance | null;
  transactions: Transaction[];
  loading: boolean;
  error: string | null;
  refreshBalance: () => Promise<void>;
  refreshTransactions: () => Promise<void>;
  depositFunds: (data: {
    amount: number;
    currency: string;
    payment_method: string;
  }) => Promise<void>;
}

export const useBalance = (): UseBalanceReturn => {
  const [balance, setBalance] = useState<UserBalance | null>(null);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  // Получить баланс
  const fetchBalance = async () => {
    try {
      setError(null);
      const data = await balanceService.getBalance();
      setBalance(data);
    } catch (err: any) {
      // Если ошибка авторизации, просто устанавливаем пустые данные
      if (err.response?.status === 401) {
        setBalance(null);
        setError('Authentication required');
        return;
      }
      const message = err.response?.data?.error || 'Failed to fetch balance';
      setError(message);
      console.error('Error fetching balance:', err);
    }
  };

  // Получить историю транзакций
  const fetchTransactions = async () => {
    try {
      const data = await balanceService.getTransactionHistory({ limit: 10 });
      setTransactions(data?.transactions || []);
    } catch (err: any) {
      // Если ошибка авторизации, просто устанавливаем пустой массив
      if (err.response?.status === 401) {
        setTransactions([]);
        return;
      }
      console.error('Error fetching transactions:', err);
      setTransactions([]);
    }
  };

  // Пополнить баланс
  const depositFunds = useCallback(
    async (data: {
      amount: number;
      currency: string;
      payment_method: string;
    }) => {
      try {
        setLoading(true);
        setError(null);

        const response = await balanceService.createDeposit({
          ...data,
          return_url: `${window.location.origin}/balance/deposit/success`,
        });

        // Перенаправление на платёжный шлюз
        if ((response as any).redirect_url) {
          window.location.href = (response as any).redirect_url;
        } else {
          toast.error('No redirect URL received from payment gateway');
        }
      } catch (err: any) {
        const message = err.response?.data?.error || 'Failed to create deposit';
        setError(message);
        toast.error(message);
        throw err;
      } finally {
        setLoading(false);
      }
    },
    []
  );

  // Обновить баланс
  const refreshBalance = useCallback(async () => {
    setLoading(true);
    await fetchBalance();
    setLoading(false);
  }, []);

  // Обновить транзакции
  const refreshTransactions = useCallback(async () => {
    await fetchTransactions();
  }, []);

  // Загрузка при монтировании
  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      await Promise.all([fetchBalance(), fetchTransactions()]);
      setLoading(false);
    };

    loadData();
  }, []); // Убираем зависимости чтобы избежать бесконечных циклов

  return {
    balance,
    transactions,
    loading,
    error,
    refreshBalance,
    refreshTransactions,
    depositFunds,
  };
};
