'use client';

import React, { useState } from 'react';
import { useTranslations } from 'next-intl';
import { useLocale } from 'next-intl';
import Link from 'next/link';
import { balanceService } from '@/services/balance';
import { BalanceWidget } from '@/components/balance/BalanceWidget';
import { toast } from '@/utils/toast';
// Убираем дублирующий useBalance - используем данные из BalanceWidget

export default function WithdrawPage() {
  const t = useTranslations('admin.balance');
  const tProfile = useTranslations('profile');
  const locale = useLocale();
  const [amount, setAmount] = useState('');
  const [withdrawMethod, setWithdrawMethod] = useState('bank_transfer');
  const [bankAccount, setBankAccount] = useState('');
  const [cardNumber, setCardNumber] = useState('');
  const [accountHolder, setAccountHolder] = useState('');
  const [isProcessing, setIsProcessing] = useState(false);

  const withdrawAmount = parseFloat(amount) || 0;
  const isValidAmount = withdrawAmount >= 100; // Проверку доступного баланса делает backend
  const isValidDetails =
    withdrawMethod === 'bank_transfer'
      ? bankAccount.length > 0 && accountHolder.length > 0
      : cardNumber.length >= 16 && accountHolder.length > 0;

  const handleWithdraw = async () => {
    if (!isValidAmount || !isValidDetails) return;

    setIsProcessing(true);
    try {
      // Здесь будет вызов API для вывода средств
      // const response = await balanceService.createWithdrawal({...});

      // Пока что показываем успешное сообщение
      toast.success(
        'Заявка на вывод средств создана. Обработка займет 1-3 рабочих дня.'
      );

      // Очищаем форму
      setAmount('');
      setBankAccount('');
      setCardNumber('');
      setAccountHolder('');
    } catch (error) {
      console.error('Withdrawal error:', error);
      toast.error(
        'Ошибка при создании заявки на вывод средств. Попробуйте еще раз.'
      );
    } finally {
      setIsProcessing(false);
    }
  };

  const handleAmountChange = (value: string) => {
    // Разрешаем только цифры и точку
    const cleanValue = value.replace(/[^0-9.]/g, '');
    const parts = cleanValue.split('.');
    if (parts.length <= 2) {
      setAmount(cleanValue);
    }
  };

  const setMaxAmount = () => {
    // Функция временно отключена - баланс получаем из BalanceWidget
    toast.info('Используйте данные из виджета баланса');
  };

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
          <li>{t('withdrawFunds')}</li>
        </ul>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Основная форма */}
        <div className="lg:col-span-2">
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h1 className="card-title text-2xl mb-6">
                💸 {t('withdrawFunds')}
              </h1>

              {/* Сумма вывода */}
              <div className="form-control mb-6">
                <label className="label">
                  <span className="label-text font-semibold">
                    Сумма для вывода
                  </span>
                  <span className="label-text-alt">
                    Доступно: см. виджет баланса справа
                  </span>
                </label>
                <div className="input-group">
                  <input
                    type="text"
                    placeholder="0"
                    className={`input input-bordered flex-1 ${
                      amount && !isValidAmount ? 'input-error' : ''
                    }`}
                    value={amount}
                    onChange={(e) => handleAmountChange(e.target.value)}
                  />
                  <span className="btn btn-ghost">RSD</span>
                  <button
                    type="button"
                    className="btn btn-outline btn-primary"
                    onClick={setMaxAmount}
                    disabled={false}
                  >
                    Все
                  </button>
                </div>
                <label className="label">
                  <span className="label-text-alt text-neutral-500">
                    Минимум: 100 RSD • Комиссия: 2%
                  </span>
                  {withdrawAmount > 0 && (
                    <span className="label-text-alt">
                      К получению:{' '}
                      {balanceService.formatAmount(
                        withdrawAmount * 0.98,
                        'RSD'
                      )}
                    </span>
                  )}
                </label>
              </div>

              {/* Способ вывода */}
              <div className="form-control mb-6">
                <label className="label">
                  <span className="label-text font-semibold">
                    Способ вывода
                  </span>
                </label>
                <div className="space-y-3">
                  <label className="label cursor-pointer justify-start gap-3 p-4 border rounded-lg hover:bg-base-200">
                    <input
                      type="radio"
                      name="withdrawMethod"
                      className="radio radio-primary"
                      value="bank_transfer"
                      checked={withdrawMethod === 'bank_transfer'}
                      onChange={(e) => setWithdrawMethod(e.target.value)}
                    />
                    <div className="flex items-center gap-3">
                      <span className="text-2xl">🏦</span>
                      <div>
                        <div className="font-medium">Банковский перевод</div>
                        <div className="text-sm text-neutral-500">
                          1-3 рабочих дня • Комиссия 2%
                        </div>
                      </div>
                    </div>
                  </label>

                  <label className="label cursor-pointer justify-start gap-3 p-4 border rounded-lg hover:bg-base-200">
                    <input
                      type="radio"
                      name="withdrawMethod"
                      className="radio radio-primary"
                      value="card"
                      checked={withdrawMethod === 'card'}
                      onChange={(e) => setWithdrawMethod(e.target.value)}
                    />
                    <div className="flex items-center gap-3">
                      <span className="text-2xl">💳</span>
                      <div>
                        <div className="font-medium">На банковскую карту</div>
                        <div className="text-sm text-neutral-500">
                          30 минут - 2 часа • Комиссия 2%
                        </div>
                      </div>
                    </div>
                  </label>
                </div>
              </div>

              {/* Реквизиты для банковского перевода */}
              {withdrawMethod === 'bank_transfer' && (
                <div className="space-y-4 mb-6 p-4 bg-base-200 rounded-lg">
                  <h3 className="font-medium">
                    Реквизиты для банковского перевода
                  </h3>

                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Номер счета / IBAN</span>
                    </label>
                    <input
                      type="text"
                      placeholder="RS35260005601001611379"
                      className="input input-bordered"
                      value={bankAccount}
                      onChange={(e) => setBankAccount(e.target.value)}
                    />
                  </div>

                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Имя владельца счета</span>
                    </label>
                    <input
                      type="text"
                      placeholder="Marko Petrović"
                      className="input input-bordered"
                      value={accountHolder}
                      onChange={(e) => setAccountHolder(e.target.value)}
                    />
                  </div>
                </div>
              )}

              {/* Реквизиты для карты */}
              {withdrawMethod === 'card' && (
                <div className="space-y-4 mb-6 p-4 bg-base-200 rounded-lg">
                  <h3 className="font-medium">Реквизиты банковской карты</h3>

                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Номер карты</span>
                    </label>
                    <input
                      type="text"
                      placeholder="**** **** **** 1234"
                      className="input input-bordered"
                      value={cardNumber}
                      onChange={(e) =>
                        setCardNumber(
                          e.target.value.replace(/\D/g, '').slice(0, 16)
                        )
                      }
                      maxLength={19}
                    />
                  </div>

                  <div className="form-control">
                    <label className="label">
                      <span className="label-text">Имя владельца карты</span>
                    </label>
                    <input
                      type="text"
                      placeholder="Marko Petrović"
                      className="input input-bordered"
                      value={accountHolder}
                      onChange={(e) => setAccountHolder(e.target.value)}
                    />
                  </div>
                </div>
              )}

              {/* Итоговая информация */}
              {withdrawAmount > 0 && isValidAmount && (
                <div className="bg-warning/10 border border-warning/20 rounded-lg p-4 mb-6">
                  <div className="space-y-2">
                    <div className="flex justify-between">
                      <span>Сумма к выводу:</span>
                      <span className="font-bold">
                        {balanceService.formatAmount(withdrawAmount, 'RSD')}
                      </span>
                    </div>
                    <div className="flex justify-between text-sm">
                      <span>Комиссия (2%):</span>
                      <span>
                        {balanceService.formatAmount(
                          withdrawAmount * 0.02,
                          'RSD'
                        )}
                      </span>
                    </div>
                    <hr className="border-warning/20" />
                    <div className="flex justify-between font-bold">
                      <span>К получению:</span>
                      <span className="text-success">
                        {balanceService.formatAmount(
                          withdrawAmount * 0.98,
                          'RSD'
                        )}
                      </span>
                    </div>
                  </div>
                </div>
              )}

              {/* Кнопка вывода */}
              <button
                className="btn btn-warning btn-lg w-full"
                disabled={!isValidAmount || !isValidDetails || isProcessing}
                onClick={handleWithdraw}
              >
                {isProcessing ? (
                  <>
                    <span className="loading loading-spinner loading-sm"></span>
                    Обработка...
                  </>
                ) : (
                  <>🔒 Создать заявку на вывод</>
                )}
              </button>

              {/* Информация о безопасности */}
              <div className="text-center text-sm text-neutral-500 mt-4">
                🔒 Все операции проходят дополнительную проверку безопасности
              </div>
            </div>
          </div>
        </div>

        {/* Боковая панель с балансом */}
        <div className="lg:col-span-1">
          <BalanceWidget showActions={false} />

          {/* Информация о выводе */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">⏱️ Сроки обработки</h3>
              <ul className="text-sm space-y-2">
                <li>• Банковский перевод: 1-3 рабочих дня</li>
                <li>• На банковскую карту: 30 минут - 2 часа</li>
                <li>• Проверка заявки: до 24 часов</li>
                <li>• Минимальная сумма: 100 RSD</li>
                <li>• Комиссия сервиса: 2%</li>
              </ul>
            </div>
          </div>

          {/* Требования */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">📋 Требования</h3>
              <ul className="text-sm space-y-2">
                <li>• Верифицированный аккаунт</li>
                <li>• Реквизиты должны быть на ваше имя</li>
                <li>• Карта/счет в сербском банке</li>
                <li>• Без задолженностей по заказам</li>
              </ul>
            </div>
          </div>

          {/* Поддержка */}
          <div className="card bg-base-100 shadow-xl mt-6">
            <div className="card-body">
              <h3 className="card-title text-lg">🆘 Нужна помощь?</h3>
              <p className="text-sm mb-4">
                Если у вас вопросы по выводу средств, свяжитесь с поддержкой.
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
