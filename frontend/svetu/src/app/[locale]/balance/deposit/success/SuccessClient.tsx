'use client';

import React, { useEffect, useState } from 'react';
import { useTranslations } from 'next-intl';
import { useLocale } from 'next-intl';
import Link from 'next/link';
import { useSearchParams } from 'next/navigation';
import { BalanceWidget } from '@/components/balance/BalanceWidget';

export default function SuccessClient() {
  const t = useTranslations('admin');
  const tProfile = useTranslations('profile');
  const locale = useLocale();
  const searchParams = useSearchParams();
  const [transactionId, setTransactionId] = useState<string | null>(null);
  const [amount, setAmount] = useState<string | null>(null);
  const [key, setKey] = useState(0); // Ключ для перерендеринга BalanceWidget
  const [isUpdatingBalance, setIsUpdatingBalance] = useState(true);

  useEffect(() => {
    // Получаем параметры из URL
    const txId = searchParams?.get('transaction_id');
    const amountParam = searchParams?.get('amount');

    if (txId) setTransactionId(txId);
    if (amountParam) setAmount(amountParam);

    // Принудительно обновляем BalanceWidget несколько раз
    // для гарантии получения актуального баланса
    const timers: NodeJS.Timeout[] = [];

    // Обновляем через 500ms, 1.5s, 3s и 5s
    [500, 1500, 3000, 5000].forEach((delay) => {
      const timer = setTimeout(() => {
        setKey((prev) => prev + 1);
      }, delay);
      timers.push(timer);
    });

    // Скрываем индикатор обновления через 6 секунд
    const hideIndicatorTimer = setTimeout(() => {
      setIsUpdatingBalance(false);
    }, 6000);

    return () => {
      timers.forEach((timer) => clearTimeout(timer));
      clearTimeout(hideIndicatorTimer);
    };
  }, [searchParams]);

  return (
    <div className="container mx-auto px-4 py-8 max-w-4xl">
      {/* Хлебные крошки */}
      <div className="breadcrumbs text-sm mb-6">
        <ul>
          <li>
            <Link href={`/${locale}/profile`}>{tProfile('title')}</Link>
          </li>
          <li>
            <Link href={`/${locale}/balance`}>{t('title')}</Link>
          </li>
          <li>
            <Link href={`/${locale}/balance/deposit`}>{t('depositFunds')}</Link>
          </li>
          <li>Успешно</li>
        </ul>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Основное содержимое */}
        <div className="lg:col-span-2">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body text-center">
              {/* Иконка успеха */}
              <div className="text-6xl text-success mb-4">✅</div>

              <h1 className="text-3xl font-bold text-success mb-4">
                Пополнение выполнено успешно!
              </h1>

              {amount && (
                <div className="bg-success/10 rounded-lg p-4 mb-6">
                  <p className="text-lg font-medium mb-2">
                    Баланс пополнен на:
                  </p>
                  <p className="text-3xl font-bold text-success">
                    {amount} RSD
                  </p>
                </div>
              )}

              {transactionId && (
                <div className="bg-base-200 rounded-lg p-4 mb-6">
                  <p className="text-sm text-neutral-600 mb-1">
                    Номер транзакции:
                  </p>
                  <p className="font-mono text-sm break-all">{transactionId}</p>
                </div>
              )}

              <div className="space-y-4">
                <p className="text-lg">
                  Средства уже доступны на вашем балансе и готовы к
                  использованию.
                </p>

                <div className="text-sm text-neutral-500 space-y-1">
                  <p>• Средства зачислены мгновенно</p>
                  <p>• Уведомление отправлено на вашу почту</p>
                  <p>• Транзакция сохранена в истории операций</p>
                </div>
              </div>

              {/* Действия */}
              <div className="card-actions justify-center mt-8 gap-4">
                <Link
                  href={`/${locale}/balance/history`}
                  className="btn btn-outline btn-primary"
                >
                  📋 История транзакций
                </Link>
                <Link href={`/${locale}/profile`} className="btn btn-primary">
                  👤 Перейти в профиль
                </Link>
              </div>
            </div>
          </div>

          {/* Дополнительные действия */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h2 className="card-title">Что дальше?</h2>
              <div className="grid grid-cols-1 sm:grid-cols-2 gap-4 mt-4">
                <Link href={`/${locale}/search`} className="btn btn-outline">
                  🛍️ Покупки в маркетплейсе
                </Link>
                <Link
                  href={`/${locale}/profile/storefronts`}
                  className="btn btn-outline"
                >
                  🏪 Витрины товаров
                </Link>
                <Link
                  href={`/${locale}/balance/deposit`}
                  className="btn btn-outline"
                >
                  💳 Пополнить еще
                </Link>
                <Link
                  href={`/${locale}/balance/withdraw`}
                  className="btn btn-outline"
                >
                  💸 Вывести средства
                </Link>
              </div>
            </div>
          </div>
        </div>

        {/* Боковая панель с обновленным балансом */}
        <div className="lg:col-span-1">
          {isUpdatingBalance && (
            <div className="alert alert-info mb-4">
              <span className="loading loading-spinner loading-sm"></span>
              <span>Обновляем баланс...</span>
            </div>
          )}
          <BalanceWidget key={key} showActions={true} />

          {/* Информация о транзакции */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">📊 Детали операции</h3>
              <div className="space-y-2 text-sm">
                <div className="flex justify-between">
                  <span className="text-neutral-600">Тип:</span>
                  <span className="font-medium">Пополнение баланса</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-neutral-600">Способ оплаты:</span>
                  <span className="font-medium">Банковская карта</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-neutral-600">Комиссия:</span>
                  <span className="font-medium text-success">0 RSD</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-neutral-600">Статус:</span>
                  <span className="badge badge-success">Завершено</span>
                </div>
                <div className="flex justify-between">
                  <span className="text-neutral-600">Время:</span>
                  <span className="font-medium">
                    {new Date().toLocaleString('ru-RU')}
                  </span>
                </div>
              </div>
            </div>
          </div>

          {/* Поддержка */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">🆘 Нужна помощь?</h3>
              <p className="text-sm mb-4">
                Если у вас есть вопросы по транзакции, обратитесь в службу
                поддержки.
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
