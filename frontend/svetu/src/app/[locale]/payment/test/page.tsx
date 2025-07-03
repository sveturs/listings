'use client';

import { useState } from 'react';
import { useAllSecurePayment } from '@/hooks/useAllSecurePayment';
import { paymentConfig } from '@/config/payment';

export default function PaymentTestPage() {
  const { createTestPayment, isProcessing, error, clearError } =
    useAllSecurePayment();
  const [amount, setAmount] = useState(5000);
  const [customAmount, setCustomAmount] = useState('');

  const predefinedAmounts = [
    { label: '50 RSD', value: 5000 },
    { label: '100 RSD', value: 10000 },
    { label: '500 RSD', value: 50000 },
    { label: '1000 RSD', value: 100000 },
  ];

  const handleTestPayment = async () => {
    clearError();
    try {
      const finalAmount = customAmount
        ? parseFloat(customAmount) * 100
        : amount;
      await createTestPayment(finalAmount);
    } catch (err) {
      console.error('Test payment failed:', err);
    }
  };

  const formatAmount = (value: number) => {
    return new Intl.NumberFormat('sr-RS', {
      style: 'currency',
      currency: 'RSD',
    }).format(value / 100);
  };

  return (
    <div className="min-h-screen bg-base-200 py-8">
      <div className="container mx-auto px-4">
        <div className="max-w-2xl mx-auto">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-3xl font-bold mb-2">
              Тестирование платежной системы
            </h1>
            <p className="text-base-content/70">
              Страница для тестирования интеграции с AllSecure (Mock режим)
            </p>
          </div>

          {/* Status Card */}
          <div className="card bg-base-100 shadow-xl mb-6">
            <div className="card-body">
              <h2 className="card-title">Статус системы</h2>
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="stat bg-base-200 rounded-lg p-4">
                  <div className="stat-title">Режим</div>
                  <div className="stat-value text-lg">
                    <span
                      className={`badge ${paymentConfig.mode === 'mock' ? 'badge-warning' : 'badge-success'}`}
                    >
                      {paymentConfig.mode.toUpperCase()}
                    </span>
                  </div>
                </div>
                <div className="stat bg-base-200 rounded-lg p-4">
                  <div className="stat-title">Вероятность успеха</div>
                  <div className="stat-value text-lg">
                    {(paymentConfig.mock?.config.successRate || 0) * 100}%
                  </div>
                </div>
                <div className="stat bg-base-200 rounded-lg p-4">
                  <div className="stat-title">3D Secure</div>
                  <div className="stat-value text-lg">
                    {(paymentConfig.mock?.config.require3DSRate || 0) * 100}%
                  </div>
                </div>
                <div className="stat bg-base-200 rounded-lg p-4">
                  <div className="stat-title">Задержка API</div>
                  <div className="stat-value text-lg">
                    {paymentConfig.mock?.config.apiDelay || 0}ms
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Test Payment Card */}
          <div className="card bg-base-100 shadow-xl mb-6">
            <div className="card-body">
              <h2 className="card-title">Создать тестовый платеж</h2>

              {error && (
                <div className="alert alert-error mb-4">
                  <svg
                    xmlns="http://www.w3.org/2000/svg"
                    className="stroke-current shrink-0 h-6 w-6"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth="2"
                      d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
                    />
                  </svg>
                  <span>{error}</span>
                </div>
              )}

              <div className="form-control mb-4">
                <label className="label">
                  <span className="label-text">Выберите сумму:</span>
                </label>
                <div className="grid grid-cols-2 md:grid-cols-4 gap-2 mb-4">
                  {predefinedAmounts.map((preset) => (
                    <button
                      key={preset.value}
                      className={`btn ${amount === preset.value ? 'btn-primary' : 'btn-outline'}`}
                      onClick={() => {
                        setAmount(preset.value);
                        setCustomAmount('');
                      }}
                    >
                      {preset.label}
                    </button>
                  ))}
                </div>
              </div>

              <div className="form-control mb-4">
                <label className="label">
                  <span className="label-text">
                    Или введите свою сумму (RSD):
                  </span>
                </label>
                <input
                  type="number"
                  className="input input-bordered"
                  placeholder="100.00"
                  value={customAmount}
                  onChange={(e) => {
                    setCustomAmount(e.target.value);
                    setAmount(0);
                  }}
                  min="1"
                  max="100000"
                  step="0.01"
                />
              </div>

              <div className="alert alert-info mb-4">
                <svg
                  xmlns="http://www.w3.org/2000/svg"
                  fill="none"
                  viewBox="0 0 24 24"
                  className="stroke-current shrink-0 w-6 h-6"
                >
                  <path
                    strokeLinecap="round"
                    strokeLinejoin="round"
                    strokeWidth="2"
                    d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"
                  ></path>
                </svg>
                <div>
                  <h3 className="font-bold">
                    К оплате:{' '}
                    {formatAmount(
                      customAmount ? parseFloat(customAmount) * 100 : amount
                    )}
                  </h3>
                  <div className="text-sm">
                    Тестовая транзакция с фиктивными данными
                  </div>
                </div>
              </div>

              <button
                className={`btn btn-primary w-full ${isProcessing ? 'loading' : ''}`}
                onClick={handleTestPayment}
                disabled={isProcessing || (!amount && !customAmount)}
              >
                {isProcessing
                  ? 'Создание платежа...'
                  : 'Создать тестовый платеж'}
              </button>
            </div>
          </div>

          {/* Test Cards Info */}
          <div className="card bg-base-100 shadow-xl">
            <div className="card-body">
              <h2 className="card-title">Тестовые карты</h2>
              <p className="text-base-content/70 mb-4">
                Используйте эти номера карт для тестирования различных
                сценариев:
              </p>

              <div className="overflow-x-auto">
                <table className="table table-compact w-full">
                  <thead>
                    <tr>
                      <th>Номер карты</th>
                      <th>Тип</th>
                      <th>Описание</th>
                    </tr>
                  </thead>
                  <tbody>
                    {paymentConfig.mock?.testCards.map((card, index) => (
                      <tr key={index}>
                        <td className="font-mono">{card.number}</td>
                        <td>
                          <span
                            className={`badge ${
                              card.type === 'success'
                                ? 'badge-success'
                                : card.type === 'declined'
                                  ? 'badge-error'
                                  : card.type === '3ds_required'
                                    ? 'badge-warning'
                                    : 'badge-info'
                            }`}
                          >
                            {card.type}
                          </span>
                        </td>
                        <td>{card.description}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>

              <div className="mt-4 p-4 bg-base-200 rounded-lg">
                <h4 className="font-semibold mb-2">
                  Дополнительные данные для тестирования:
                </h4>
                <ul className="text-sm space-y-1">
                  <li>
                    <strong>Имя на карте:</strong> любое (например, TEST USER)
                  </li>
                  <li>
                    <strong>Срок действия:</strong> любая будущая дата
                    (например, 12/25)
                  </li>
                  <li>
                    <strong>CVV:</strong> любые 3-4 цифры (например, 123)
                  </li>
                  <li>
                    <strong>3D Secure код:</strong> 123 (для успешной
                    аутентификации)
                  </li>
                </ul>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
