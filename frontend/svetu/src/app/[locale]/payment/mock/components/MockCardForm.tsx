'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import { z } from 'zod';
// import { useTranslations } from 'next-intl';

const cardSchema = z.object({
  cardNumber: z.string().regex(/^\d{16}$/, 'Invalid card number'),
  cardHolder: z.string().min(3, 'Cardholder name required'),
  expiryMonth: z.string().regex(/^(0[1-9]|1[0-2])$/, 'Invalid month'),
  expiryYear: z.string().regex(/^\d{2}$/, 'Invalid year'),
  cvv: z.string().regex(/^\d{3,4}$/, 'Invalid CVV'),
});

interface TestCard {
  number: string;
  type: string;
  description: string;
}

interface MockCardFormProps {
  onSubmit: (data: any) => void;
  testCards: TestCard[];
}

export default function MockCardForm({
  onSubmit,
  testCards,
}: MockCardFormProps) {
  const [showTestCards, setShowTestCards] = useState(true);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(cardSchema),
  });

  const fillTestCard = (card: TestCard) => {
    setValue('cardNumber', card.number);
    setValue('cardHolder', 'TEST USER');
    setValue('expiryMonth', '12');
    setValue('expiryYear', '25');
    setValue('cvv', '123');
  };

  const handleFormSubmit = async (data: any) => {
    setIsSubmitting(true);
    try {
      await onSubmit(data);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
      {/* Test Cards Helper */}
      {showTestCards && (
        <div className="alert alert-warning">
          <div className="w-full">
            <div className="flex justify-between items-center mb-2">
              <span className="font-semibold">Тестовые карты</span>
              <button
                type="button"
                className="btn btn-ghost btn-xs"
                onClick={() => setShowTestCards(false)}
              >
                ✕
              </button>
            </div>
            <div className="space-y-1">
              {testCards.map((card) => (
                <button
                  key={card.number}
                  type="button"
                  className="btn btn-xs btn-outline w-full justify-start"
                  onClick={() => fillTestCard(card)}
                >
                  <span className="font-mono">{card.number}</span>
                  <span className="ml-auto text-xs opacity-70">
                    {card.description}
                  </span>
                </button>
              ))}
            </div>
          </div>
        </div>
      )}

      {/* Card Number */}
      <div className="form-control">
        <label className="label">
          <span className="label-text">Номер карты</span>
        </label>
        <input
          type="text"
          className={`input input-bordered font-mono ${errors.cardNumber ? 'input-error' : ''}`}
          placeholder="1234 5678 9012 3456"
          maxLength={16}
          {...register('cardNumber')}
        />
        {errors.cardNumber && (
          <label className="label">
            <span className="label-text-alt text-error">
              {errors.cardNumber.message}
            </span>
          </label>
        )}
      </div>

      {/* Cardholder Name */}
      <div className="form-control">
        <label className="label">
          <span className="label-text">Имя на карте</span>
        </label>
        <input
          type="text"
          className={`input input-bordered ${errors.cardHolder ? 'input-error' : ''}`}
          placeholder="JOHN DOE"
          style={{ textTransform: 'uppercase' }}
          {...register('cardHolder')}
        />
        {errors.cardHolder && (
          <label className="label">
            <span className="label-text-alt text-error">
              {errors.cardHolder.message}
            </span>
          </label>
        )}
      </div>

      {/* Expiry and CVV */}
      <div className="grid grid-cols-3 gap-4">
        <div className="form-control">
          <label className="label">
            <span className="label-text">Месяц</span>
          </label>
          <input
            type="text"
            className={`input input-bordered ${errors.expiryMonth ? 'input-error' : ''}`}
            placeholder="MM"
            maxLength={2}
            {...register('expiryMonth')}
          />
          {errors.expiryMonth && (
            <label className="label">
              <span className="label-text-alt text-error">Неверный месяц</span>
            </label>
          )}
        </div>

        <div className="form-control">
          <label className="label">
            <span className="label-text">Год</span>
          </label>
          <input
            type="text"
            className={`input input-bordered ${errors.expiryYear ? 'input-error' : ''}`}
            placeholder="YY"
            maxLength={2}
            {...register('expiryYear')}
          />
          {errors.expiryYear && (
            <label className="label">
              <span className="label-text-alt text-error">Неверный год</span>
            </label>
          )}
        </div>

        <div className="form-control">
          <label className="label">
            <span className="label-text">CVV</span>
          </label>
          <input
            type="text"
            className={`input input-bordered ${errors.cvv ? 'input-error' : ''}`}
            placeholder="123"
            maxLength={4}
            {...register('cvv')}
          />
          {errors.cvv && (
            <label className="label">
              <span className="label-text-alt text-error">Неверный CVV</span>
            </label>
          )}
        </div>
      </div>

      {/* Submit Button */}
      <button
        type="submit"
        className={`btn btn-primary w-full ${isSubmitting ? 'loading' : ''}`}
        disabled={isSubmitting}
      >
        {isSubmitting ? 'Обработка...' : 'Оплатить'}
      </button>

      {/* Card Brands */}
      <div className="flex justify-center gap-2 opacity-50 mt-4">
        <div className="badge badge-ghost">VISA</div>
        <div className="badge badge-ghost">MasterCard</div>
        <div className="badge badge-ghost">Maestro</div>
      </div>
    </form>
  );
}
