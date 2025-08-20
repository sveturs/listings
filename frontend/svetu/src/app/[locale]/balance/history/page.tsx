'use client';

import React, { useState, useEffect } from 'react';
import { useTranslations } from 'next-intl';
import { useLocale } from 'next-intl';
import Link from 'next/link';
import { balanceService } from '@/services/balance';
import { BalanceWidget } from '@/components/balance/BalanceWidget';

interface Transaction {
  id: string;
  type: 'deposit' | 'withdrawal' | 'payment' | 'refund';
  amount: number;
  currency: string;
  status: 'pending' | 'completed' | 'failed' | 'cancelled';
  description: string;
  created_at: string;
  payment_method?: string;
}

export default function TransactionHistoryPage() {
  const t = useTranslations('balance');
  const tProfile = useTranslations('profile');
  const locale = useLocale();
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<string>('all');
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const itemsPerPage = 20;

  useEffect(() => {
    loadTransactions();
  }, [currentPage, filter]); // eslint-disable-line react-hooks/exhaustive-deps

  const loadTransactions = async () => {
    try {
      setLoading(true);
      setError(null);

      const params = {
        limit: itemsPerPage,
        offset: (currentPage - 1) * itemsPerPage,
        type: filter === 'all' ? undefined : filter,
      };

      const response = await balanceService.getTransactionHistory(params);

      if (response.success && response.data) {
        setTransactions(response.data.transactions || []);
        setTotalPages(Math.ceil((response.data.total || 0) / itemsPerPage));
      } else {
        setTransactions([]);
      }
    } catch (err) {
      console.error('Error loading transactions:', err);
      setError('Ошибка загрузки истории транзакций');
      setTransactions([]);
    } finally {
      setLoading(false);
    }
  };

  const getTransactionIcon = (type: string) => {
    switch (type) {
      case 'deposit':
        return '⬇️';
      case 'withdrawal':
        return '⬆️';
      case 'payment':
        return '💳';
      case 'refund':
        return '↩️';
      default:
        return '💰';
    }
  };

  const getTransactionTypeText = (type: string) => {
    switch (type) {
      case 'deposit':
        return 'Пополнение';
      case 'withdrawal':
        return 'Вывод средств';
      case 'payment':
        return 'Платёж';
      case 'refund':
        return 'Возврат';
      default:
        return 'Операция';
    }
  };

  const getStatusBadge = (status: string) => {
    const badges = {
      pending: 'badge-warning',
      completed: 'badge-success',
      failed: 'badge-error',
      cancelled: 'badge-neutral',
    };

    const texts = {
      pending: 'В обработке',
      completed: 'Завершено',
      failed: 'Ошибка',
      cancelled: 'Отменено',
    };

    return (
      <span
        className={`badge badge-sm ${badges[status as keyof typeof badges] || 'badge-neutral'}`}
      >
        {texts[status as keyof typeof texts] || status}
      </span>
    );
  };

  const formatDate = (dateString: string) => {
    const date = new Date(dateString);
    return date.toLocaleDateString('ru-RU', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  const filteredTransactions = transactions.filter(
    (tx) => filter === 'all' || tx.type === filter
  );

  return (
    <div className="container mx-auto px-4 py-8 max-w-6xl">
      {/* Хлебные крошки */}
      <div className="breadcrumbs text-sm mb-6">
        <ul>
          <li>
            <Link href={`/${locale}/profile`}>{tProfile('title')}</Link>
          </li>
          <li>
            <Link href={`/${locale}/balance`}>{t('title')}</Link>
          </li>
          <li>{t('transactionHistory')}</li>
        </ul>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
        {/* Основной контент */}
        <div className="lg:col-span-3">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
                <h1 className="card-title text-2xl">
                  📋 {t('transactionHistory')}
                </h1>

                {/* Фильтры */}
                <div className="flex gap-2">
                  <select
                    className="select select-bordered select-sm"
                    value={filter}
                    onChange={(e) => {
                      setFilter(e.target.value);
                      setCurrentPage(1);
                    }}
                  >
                    <option value="all">Все операции</option>
                    <option value="deposit">Пополнения</option>
                    <option value="withdrawal">Выводы</option>
                    <option value="payment">Платежи</option>
                    <option value="refund">Возвраты</option>
                  </select>
                  <button
                    className="btn btn-sm btn-ghost"
                    onClick={loadTransactions}
                    disabled={loading}
                  >
                    🔄
                  </button>
                </div>
              </div>

              {/* Загрузка */}
              {loading && (
                <div className="flex justify-center items-center py-8">
                  <span className="loading loading-spinner loading-lg"></span>
                </div>
              )}

              {/* Ошибка */}
              {error && (
                <div className="alert alert-error">
                  <span>{error}</span>
                  <button
                    className="btn btn-sm btn-ghost"
                    onClick={loadTransactions}
                  >
                    Повторить
                  </button>
                </div>
              )}

              {/* Пустое состояние */}
              {!loading && !error && filteredTransactions.length === 0 && (
                <div className="text-center py-12">
                  <div className="text-6xl mb-4">📊</div>
                  <h3 className="text-xl font-semibold mb-2">
                    История операций пуста
                  </h3>
                  <p className="text-neutral-500 mb-6">
                    {filter === 'all'
                      ? 'У вас пока нет транзакций'
                      : `Нет операций типа "${getTransactionTypeText(filter)}"`}
                  </p>
                  <div className="flex gap-4 justify-center">
                    <Link
                      href={`/${locale}/balance/deposit`}
                      className="btn btn-primary"
                    >
                      💳 Пополнить баланс
                    </Link>
                    <Link
                      href={`/${locale}/marketplace`}
                      className="btn btn-outline"
                    >
                      🛍️ Перейти в маркетплейс
                    </Link>
                  </div>
                </div>
              )}

              {/* Список транзакций */}
              {!loading && !error && filteredTransactions.length > 0 && (
                <>
                  <div className="overflow-x-auto">
                    <table className="table table-zebra">
                      <thead>
                        <tr>
                          <th>Тип</th>
                          <th>Описание</th>
                          <th>Сумма</th>
                          <th>Статус</th>
                          <th>Дата</th>
                        </tr>
                      </thead>
                      <tbody>
                        {filteredTransactions.map((transaction) => (
                          <tr key={transaction.id} className="hover">
                            <td>
                              <div className="flex items-center gap-3">
                                <span className="text-2xl">
                                  {getTransactionIcon(transaction.type)}
                                </span>
                                <div>
                                  <div className="font-medium">
                                    {getTransactionTypeText(transaction.type)}
                                  </div>
                                  {transaction.payment_method && (
                                    <div className="text-sm text-neutral-500">
                                      {transaction.payment_method}
                                    </div>
                                  )}
                                </div>
                              </div>
                            </td>
                            <td>
                              <div className="max-w-xs">
                                <div
                                  className="truncate"
                                  title={transaction.description}
                                >
                                  {transaction.description}
                                </div>
                                <div className="text-xs text-neutral-500">
                                  ID: {transaction.id.slice(0, 8)}...
                                </div>
                              </div>
                            </td>
                            <td>
                              <div
                                className={`font-bold ${
                                  transaction.type === 'deposit' ||
                                  transaction.type === 'refund'
                                    ? 'text-success'
                                    : 'text-error'
                                }`}
                              >
                                {transaction.type === 'deposit' ||
                                transaction.type === 'refund'
                                  ? '+'
                                  : '-'}
                                {balanceService.formatAmount(
                                  Math.abs(transaction.amount),
                                  transaction.currency
                                )}
                              </div>
                            </td>
                            <td>{getStatusBadge(transaction.status)}</td>
                            <td>
                              <div className="text-sm">
                                {formatDate(transaction.created_at)}
                              </div>
                            </td>
                          </tr>
                        ))}
                      </tbody>
                    </table>
                  </div>

                  {/* Пагинация */}
                  {totalPages > 1 && (
                    <div className="flex justify-center mt-6">
                      <div className="join">
                        <button
                          className="join-item btn btn-sm"
                          disabled={currentPage === 1}
                          onClick={() => setCurrentPage(currentPage - 1)}
                        >
                          «
                        </button>

                        {Array.from(
                          { length: Math.min(5, totalPages) },
                          (_, i) => {
                            const page =
                              Math.max(
                                1,
                                Math.min(totalPages - 4, currentPage - 2)
                              ) + i;
                            return (
                              <button
                                key={page}
                                className={`join-item btn btn-sm ${
                                  page === currentPage ? 'btn-active' : ''
                                }`}
                                onClick={() => setCurrentPage(page)}
                              >
                                {page}
                              </button>
                            );
                          }
                        )}

                        <button
                          className="join-item btn btn-sm"
                          disabled={currentPage === totalPages}
                          onClick={() => setCurrentPage(currentPage + 1)}
                        >
                          »
                        </button>
                      </div>
                    </div>
                  )}
                </>
              )}
            </div>
          </div>
        </div>

        {/* Боковая панель */}
        <div className="lg:col-span-1">
          <BalanceWidget showActions={true} compact={true} />

          {/* Статистика */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">📊 Статистика</h3>
              <div className="space-y-3 text-sm">
                <div className="flex justify-between">
                  <span>Всего операций:</span>
                  <span className="font-bold">{transactions.length}</span>
                </div>
                <div className="flex justify-between">
                  <span>Завершенных:</span>
                  <span className="font-bold text-success">
                    {
                      transactions.filter((t) => t.status === 'completed')
                        .length
                    }
                  </span>
                </div>
                <div className="flex justify-between">
                  <span>В обработке:</span>
                  <span className="font-bold text-warning">
                    {transactions.filter((t) => t.status === 'pending').length}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Быстрые действия */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">⚡ Быстрые действия</h3>
              <div className="space-y-2">
                <Link
                  href={`/${locale}/balance/deposit`}
                  className="btn btn-primary btn-sm w-full"
                >
                  💳 Пополнить баланс
                </Link>
                <Link
                  href={`/${locale}/balance/withdraw`}
                  className="btn btn-outline btn-sm w-full"
                >
                  💸 Вывести средства
                </Link>
              </div>
            </div>
          </div>

          {/* Поддержка */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">🆘 Поддержка</h3>
              <p className="text-sm mb-4">
                Вопросы по транзакциям? Мы поможем разобраться.
              </p>
              <Link
                href={`/${locale}/contacts`}
                className="btn btn-sm btn-outline"
              >
                Связаться с поддержкой
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
