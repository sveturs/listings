'use client';

import React, { useState } from 'react';
import { SveTuLogoStatic } from '@/components/logos/SveTuLogoStatic';
import { AnimatedSection } from '@/components/ui/AnimatedSection';

const EscrowPayment = () => {
  const [currentStep, setCurrentStep] = useState(1);
  const [paymentMethod, setPaymentMethod] = useState<
    'card' | 'bank' | 'crypto'
  >('card');
  const [dealStatus, setDealStatus] = useState<
    'pending' | 'paid' | 'shipped' | 'delivered' | 'completed'
  >('pending');

  const product = {
    title: 'iPhone 14 Pro Max 256GB',
    price: 899,
    seller: 'TechStore',
    buyer: 'Иван Петров',
    image:
      '/api/minio/download?fileName=listings/0a47e66f-d8da-459f-a2ba-8e2b85ae0163/38ad29e6-7b07-4bfc-9db2-d965cb6b966f.jpg',
  };

  const escrowSteps = [
    { id: 1, title: 'Оформление', icon: '📝', status: 'completed' },
    {
      id: 2,
      title: 'Оплата',
      icon: '💳',
      status: currentStep >= 2 ? 'completed' : 'pending',
    },
    {
      id: 3,
      title: 'Доставка',
      icon: '📦',
      status: currentStep >= 3 ? 'completed' : 'pending',
    },
    {
      id: 4,
      title: 'Получение',
      icon: '✅',
      status: currentStep >= 4 ? 'completed' : 'pending',
    },
    {
      id: 5,
      title: 'Завершение',
      icon: '🎉',
      status: currentStep >= 5 ? 'completed' : 'pending',
    },
  ];

  const handlePayment = () => {
    setCurrentStep(2);
    setDealStatus('paid');
    setTimeout(() => {
      setCurrentStep(3);
      setDealStatus('shipped');
    }, 2000);
  };

  const handleDeliveryConfirm = () => {
    setCurrentStep(4);
    setDealStatus('delivered');
  };

  const handleDealComplete = () => {
    setCurrentStep(5);
    setDealStatus('completed');
  };

  return (
    <div className="min-h-screen bg-gradient-to-br from-base-100 to-base-200">
      {/* Header */}
      <div className="navbar bg-base-100 shadow-lg">
        <div className="navbar-start">
          <SveTuLogoStatic variant="gradient" width={120} height={40} />
        </div>
        <div className="navbar-center">
          <h1 className="text-xl font-bold">🔒 Эскроу-защита платежей</h1>
        </div>
        <div className="navbar-end">
          <div className="badge badge-success gap-2">
            <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path
                fillRule="evenodd"
                d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                clipRule="evenodd"
              />
            </svg>
            Безопасная сделка
          </div>
        </div>
      </div>

      <div className="container mx-auto px-4 py-8 max-w-6xl">
        {/* Progress Steps */}
        <AnimatedSection animation="fadeIn">
          <div className="mb-8">
            <ul className="steps steps-horizontal w-full">
              {escrowSteps.map((step, idx) => (
                <li
                  key={step.id}
                  className={`step ${step.status === 'completed' ? 'step-primary' : ''}`}
                  data-content={step.status === 'completed' ? '✓' : step.icon}
                >
                  <span className="text-xs font-semibold">{step.title}</span>
                </li>
              ))}
            </ul>
          </div>
        </AnimatedSection>

        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8">
          {/* Main Content */}
          <div className="lg:col-span-2 space-y-6">
            {/* Deal Info */}
            <AnimatedSection animation="slideLeft">
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h2 className="card-title mb-4">Информация о сделке</h2>
                  <div className="flex gap-4">
                    <img
                      src={product.image}
                      alt={product.title}
                      className="w-24 h-24 rounded-lg object-cover"
                    />
                    <div className="flex-1">
                      <h3 className="font-bold text-lg">{product.title}</h3>
                      <div className="mt-2 space-y-1 text-sm">
                        <div className="flex justify-between">
                          <span className="text-base-content/60">
                            Продавец:
                          </span>
                          <span className="font-semibold">
                            {product.seller}
                          </span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">
                            Покупатель:
                          </span>
                          <span className="font-semibold">{product.buyer}</span>
                        </div>
                        <div className="flex justify-between">
                          <span className="text-base-content/60">Сумма:</span>
                          <span className="text-2xl font-bold text-primary">
                            €{product.price}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </AnimatedSection>

            {/* Payment Section */}
            {currentStep === 1 && dealStatus === 'pending' && (
              <AnimatedSection animation="slideUp">
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="card-title mb-4">Выберите способ оплаты</h3>

                    <div className="space-y-3">
                      <label
                        className={`card cursor-pointer ${paymentMethod === 'card' ? 'ring-2 ring-primary' : ''}`}
                      >
                        <div className="card-body p-4">
                          <div className="flex items-center gap-3">
                            <input
                              type="radio"
                              name="payment"
                              className="radio radio-primary"
                              checked={paymentMethod === 'card'}
                              onChange={() => setPaymentMethod('card')}
                            />
                            <div className="text-2xl">💳</div>
                            <div className="flex-1">
                              <div className="font-semibold">
                                Банковская карта
                              </div>
                              <div className="text-sm text-base-content/60">
                                Visa, Mastercard, Мир
                              </div>
                            </div>
                            <div className="badge badge-success">Быстро</div>
                          </div>
                        </div>
                      </label>

                      <label
                        className={`card cursor-pointer ${paymentMethod === 'bank' ? 'ring-2 ring-primary' : ''}`}
                      >
                        <div className="card-body p-4">
                          <div className="flex items-center gap-3">
                            <input
                              type="radio"
                              name="payment"
                              className="radio radio-primary"
                              checked={paymentMethod === 'bank'}
                              onChange={() => setPaymentMethod('bank')}
                            />
                            <div className="text-2xl">🏦</div>
                            <div className="flex-1">
                              <div className="font-semibold">
                                Банковский перевод
                              </div>
                              <div className="text-sm text-base-content/60">
                                SWIFT, SEPA
                              </div>
                            </div>
                          </div>
                        </div>
                      </label>

                      <label
                        className={`card cursor-pointer ${paymentMethod === 'crypto' ? 'ring-2 ring-primary' : ''}`}
                      >
                        <div className="card-body p-4">
                          <div className="flex items-center gap-3">
                            <input
                              type="radio"
                              name="payment"
                              className="radio radio-primary"
                              checked={paymentMethod === 'crypto'}
                              onChange={() => setPaymentMethod('crypto')}
                            />
                            <div className="text-2xl">₿</div>
                            <div className="flex-1">
                              <div className="font-semibold">Криптовалюта</div>
                              <div className="text-sm text-base-content/60">
                                Bitcoin, Ethereum, USDT
                              </div>
                            </div>
                            <div className="badge badge-info">Анонимно</div>
                          </div>
                        </div>
                      </label>
                    </div>

                    <div className="divider"></div>

                    <button
                      className="btn btn-primary btn-block"
                      onClick={handlePayment}
                    >
                      Оплатить €{product.price}
                    </button>
                  </div>
                </div>
              </AnimatedSection>
            )}

            {/* Payment Processing */}
            {currentStep === 2 && dealStatus === 'paid' && (
              <AnimatedSection animation="fadeIn">
                <div className="card bg-success/10 border border-success">
                  <div className="card-body text-center">
                    <div className="text-6xl mb-4 animate-bounce">✅</div>
                    <h3 className="text-2xl font-bold text-success">
                      Оплата получена!
                    </h3>
                    <p className="mt-2">
                      Средства находятся на эскроу-счете до подтверждения
                      получения товара
                    </p>
                    <div className="loading loading-dots loading-lg text-success mt-4"></div>
                    <p className="text-sm text-base-content/60 mt-2">
                      Ожидание отправки товара продавцом...
                    </p>
                  </div>
                </div>
              </AnimatedSection>
            )}

            {/* Shipping */}
            {currentStep === 3 && dealStatus === 'shipped' && (
              <AnimatedSection animation="slideUp">
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="card-title mb-4">📦 Товар отправлен</h3>
                    <div className="bg-info/10 rounded-lg p-4 mb-4">
                      <div className="flex items-center gap-3">
                        <svg
                          className="w-6 h-6 text-info"
                          fill="currentColor"
                          viewBox="0 0 20 20"
                        >
                          <path d="M8 16.5a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0zM15 16.5a1.5 1.5 0 11-3 0 1.5 1.5 0 013 0z" />
                          <path d="M3 4a1 1 0 00-1 1v10a1 1 0 001 1h1.05a2.5 2.5 0 014.9 0H10a1 1 0 001-1V5a1 1 0 00-1-1H3zM14 7h4.05C19.166 7 20 7.834 20 8.95V13h-2a2.5 2.5 0 00-4.9 0H12V7h2z" />
                        </svg>
                        <div>
                          <p className="font-semibold">
                            Трек-номер: RS123456789
                          </p>
                          <p className="text-sm text-base-content/60">
                            Курьерская служба: DHL Express
                          </p>
                        </div>
                      </div>
                    </div>

                    <div className="space-y-3">
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-success flex items-center justify-center text-white">
                          ✓
                        </div>
                        <div className="flex-1">
                          <p className="font-semibold">Передано курьеру</p>
                          <p className="text-sm text-base-content/60">
                            Сегодня, 10:30
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-primary flex items-center justify-center text-white animate-pulse">
                          •••
                        </div>
                        <div className="flex-1">
                          <p className="font-semibold">В пути</p>
                          <p className="text-sm text-base-content/60">
                            Ожидаемая доставка: завтра
                          </p>
                        </div>
                      </div>
                      <div className="flex items-center gap-3">
                        <div className="w-10 h-10 rounded-full bg-base-300 flex items-center justify-center"></div>
                        <div className="flex-1 opacity-50">
                          <p className="font-semibold">Доставлено</p>
                          <p className="text-sm text-base-content/60">
                            Ожидается
                          </p>
                        </div>
                      </div>
                    </div>

                    <button
                      className="btn btn-primary btn-block mt-4"
                      onClick={handleDeliveryConfirm}
                    >
                      Подтвердить получение
                    </button>
                    <p className="text-xs text-center text-base-content/60 mt-2">
                      После подтверждения средства будут переведены продавцу
                    </p>
                  </div>
                </div>
              </AnimatedSection>
            )}

            {/* Confirmation */}
            {currentStep === 4 && dealStatus === 'delivered' && (
              <AnimatedSection animation="slideUp">
                <div className="card bg-base-100 shadow-xl">
                  <div className="card-body">
                    <h3 className="card-title mb-4">Подтверждение получения</h3>
                    <div className="space-y-4">
                      <div className="form-control">
                        <label className="label cursor-pointer">
                          <span className="label-text">
                            Товар получен в полной комплектации
                          </span>
                          <input
                            type="checkbox"
                            className="checkbox checkbox-primary"
                            defaultChecked
                          />
                        </label>
                      </div>
                      <div className="form-control">
                        <label className="label cursor-pointer">
                          <span className="label-text">
                            Товар соответствует описанию
                          </span>
                          <input
                            type="checkbox"
                            className="checkbox checkbox-primary"
                            defaultChecked
                          />
                        </label>
                      </div>
                      <div className="form-control">
                        <label className="label cursor-pointer">
                          <span className="label-text">
                            Претензий к продавцу нет
                          </span>
                          <input
                            type="checkbox"
                            className="checkbox checkbox-primary"
                            defaultChecked
                          />
                        </label>
                      </div>

                      <div className="divider"></div>

                      <div className="form-control">
                        <label className="label">
                          <span className="label-text">Оставьте отзыв</span>
                        </label>
                        <div className="rating rating-lg">
                          {[1, 2, 3, 4, 5].map((star) => (
                            <input
                              key={star}
                              type="radio"
                              name="rating"
                              className="mask mask-star-2 bg-orange-400"
                              defaultChecked={star === 5}
                            />
                          ))}
                        </div>
                        <textarea
                          className="textarea textarea-bordered mt-2"
                          placeholder="Ваш отзыв о товаре и продавце..."
                          rows={3}
                        ></textarea>
                      </div>

                      <button
                        className="btn btn-success btn-block"
                        onClick={handleDealComplete}
                      >
                        Завершить сделку
                      </button>
                    </div>
                  </div>
                </div>
              </AnimatedSection>
            )}

            {/* Completed */}
            {currentStep === 5 && dealStatus === 'completed' && (
              <AnimatedSection animation="zoomIn">
                <div className="card bg-gradient-to-r from-success/10 to-primary/10 border-2 border-success">
                  <div className="card-body text-center">
                    <div className="text-8xl mb-4 animate-bounce">🎉</div>
                    <h3 className="text-3xl font-bold text-success">
                      Сделка завершена!
                    </h3>
                    <p className="mt-2 text-lg">
                      Спасибо за использование эскроу-защиты
                    </p>
                    <div className="mt-6 space-y-2">
                      <p className="text-sm">✅ Средства переведены продавцу</p>
                      <p className="text-sm">✅ Отзыв опубликован</p>
                      <p className="text-sm">✅ +10 баллов к вашему рейтингу</p>
                    </div>
                    <button className="btn btn-primary mt-6">
                      Перейти к покупкам
                    </button>
                  </div>
                </div>
              </AnimatedSection>
            )}
          </div>

          {/* Sidebar */}
          <div className="space-y-6">
            {/* Security Info */}
            <AnimatedSection animation="slideRight">
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title mb-4">🛡️ Как это работает?</h3>
                  <div className="space-y-4">
                    <div className="flex items-start gap-3">
                      <div className="text-2xl">1️⃣</div>
                      <div>
                        <h4 className="font-semibold">Безопасная оплата</h4>
                        <p className="text-sm text-base-content/60">
                          Деньги хранятся на защищенном счете
                        </p>
                      </div>
                    </div>
                    <div className="flex items-start gap-3">
                      <div className="text-2xl">2️⃣</div>
                      <div>
                        <h4 className="font-semibold">Контроль доставки</h4>
                        <p className="text-sm text-base-content/60">
                          Отслеживание на всех этапах
                        </p>
                      </div>
                    </div>
                    <div className="flex items-start gap-3">
                      <div className="text-2xl">3️⃣</div>
                      <div>
                        <h4 className="font-semibold">Гарантия возврата</h4>
                        <p className="text-sm text-base-content/60">
                          100% возврат при проблемах
                        </p>
                      </div>
                    </div>
                    <div className="flex items-start gap-3">
                      <div className="text-2xl">4️⃣</div>
                      <div>
                        <h4 className="font-semibold">Арбитраж споров</h4>
                        <p className="text-sm text-base-content/60">
                          Независимое решение конфликтов
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </AnimatedSection>

            {/* Benefits */}
            <AnimatedSection animation="slideRight" delay={0.2}>
              <div className="card bg-gradient-to-r from-primary/10 to-secondary/10">
                <div className="card-body">
                  <h3 className="card-title mb-4">✨ Преимущества</h3>
                  <ul className="space-y-2">
                    <li className="flex items-center gap-2">
                      <svg
                        className="w-5 h-5 text-success"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span className="text-sm">Комиссия всего 2%</span>
                    </li>
                    <li className="flex items-center gap-2">
                      <svg
                        className="w-5 h-5 text-success"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span className="text-sm">Страхование до €10,000</span>
                    </li>
                    <li className="flex items-center gap-2">
                      <svg
                        className="w-5 h-5 text-success"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span className="text-sm">24/7 поддержка</span>
                    </li>
                    <li className="flex items-center gap-2">
                      <svg
                        className="w-5 h-5 text-success"
                        fill="currentColor"
                        viewBox="0 0 20 20"
                      >
                        <path
                          fillRule="evenodd"
                          d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                          clipRule="evenodd"
                        />
                      </svg>
                      <span className="text-sm">Юридическая защита</span>
                    </li>
                  </ul>
                </div>
              </div>
            </AnimatedSection>

            {/* FAQ */}
            <AnimatedSection animation="slideRight" delay={0.3}>
              <div className="card bg-base-100 shadow-xl">
                <div className="card-body">
                  <h3 className="card-title mb-4">❓ Частые вопросы</h3>
                  <div className="join join-vertical w-full">
                    <div className="collapse collapse-arrow join-item border border-base-300">
                      <input type="radio" name="faq" defaultChecked />
                      <div className="collapse-title text-sm font-medium">
                        Что если товар не соответствует?
                      </div>
                      <div className="collapse-content text-sm">
                        <p>
                          Вы можете открыть спор и получить полный возврат
                          средств
                        </p>
                      </div>
                    </div>
                    <div className="collapse collapse-arrow join-item border border-base-300">
                      <input type="radio" name="faq" />
                      <div className="collapse-title text-sm font-medium">
                        Сколько времени на проверку?
                      </div>
                      <div className="collapse-content text-sm">
                        <p>
                          У вас есть 3 дня после получения для проверки товара
                        </p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </AnimatedSection>
          </div>
        </div>
      </div>
    </div>
  );
};

export default EscrowPayment;
