'use client';

import { useState } from 'react';
import { useForm } from 'react-hook-form';

interface Mock3DSFormProps {
  onSubmit: (code: string) => void;
  amount?: number;
}

export default function Mock3DSForm({ onSubmit, amount }: Mock3DSFormProps) {
  const [isSubmitting, setIsSubmitting] = useState(false);

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm();

  const handleFormSubmit = async (data: any) => {
    setIsSubmitting(true);
    try {
      await onSubmit(data.code);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="space-y-6">
      {/* 3DS Header */}
      <div className="text-center">
        <div className="text-lg font-semibold text-primary mb-2">
          3D Secure Authentication
        </div>
        <p className="text-sm text-base-content/70">
          Ваш банк запросил дополнительную аутентификацию для обеспечения
          безопасности транзакции
        </p>
      </div>

      {/* Bank Simulation */}
      <div className="mockup-browser bg-base-300 border border-base-300">
        <div className="mockup-browser-toolbar">
          <div className="input text-xs">https://secure.bank.rs/3ds-auth</div>
        </div>

        <div className="p-6 bg-base-100">
          <div className="text-center mb-4">
            <div className="text-lg font-bold text-accent">Банк Србије</div>
            <div className="text-sm text-base-content/60">
              3D Secure Verification
            </div>
          </div>

          {amount !== undefined && (
            <div className="alert alert-info mb-4">
              <div className="flex justify-between w-full">
                <span>Сума за плаћање:</span>
                <span className="font-bold">
                  {new Intl.NumberFormat('sr-RS', {
                    style: 'currency',
                    currency: 'RSD',
                    minimumFractionDigits: 2,
                    maximumFractionDigits: 2,
                  }).format(amount)}
                </span>
              </div>
            </div>
          )}

          <form onSubmit={handleSubmit(handleFormSubmit)} className="space-y-4">
            <div className="form-control">
              <label className="label" htmlFor="verification-code">
                <span className="label-text">Унесите код из SMS-а:</span>
              </label>
              <input
                id="verification-code"
                type="text"
                className={`input input-bordered text-center font-mono text-lg ${errors.code ? 'input-error' : ''}`}
                placeholder="123456"
                maxLength={6}
                {...register('code', {
                  required: 'Код је обавезан',
                  minLength: {
                    value: 3,
                    message: 'Код мора имати најмање 3 цифре',
                  },
                })}
              />
              {errors.code && (
                <label className="label">
                  <span className="label-text-alt text-error">
                    {String(errors.code?.message)}
                  </span>
                </label>
              )}
            </div>

            <div className="text-xs text-base-content/60 text-center">
              💡 За тестирање користите код:{' '}
              <span className="font-mono font-bold">123</span>
            </div>

            <div className="flex gap-2">
              <button
                type="submit"
                className={`btn btn-primary flex-1 ${isSubmitting ? 'loading' : ''}`}
                disabled={isSubmitting}
              >
                {isSubmitting ? 'Потврђивање...' : 'Потврди'}
              </button>
              <button
                type="button"
                className="btn btn-ghost"
                onClick={async () => {
                  setIsSubmitting(true);
                  try {
                    await onSubmit('cancel');
                  } finally {
                    setIsSubmitting(false);
                  }
                }}
                disabled={isSubmitting}
              >
                Откажи
              </button>
            </div>
          </form>

          {/* Security Info */}
          <div className="mt-6 pt-4 border-t border-base-300">
            <div className="flex items-center justify-center gap-2 text-xs text-base-content/60">
              <svg className="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fillRule="evenodd"
                  d="M2.166 4.999A11.954 11.954 0 0010 1.944 11.954 11.954 0 0017.834 5c.11.65.166 1.32.166 2.001 0 5.225-3.34 9.67-8 11.317C5.34 16.67 2 12.225 2 7c0-.682.057-1.35.166-2.001zm11.541 3.708a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                  clipRule="evenodd"
                />
              </svg>
              <span>Заштићено 256-битним SSL шифровањем</span>
            </div>
          </div>
        </div>
      </div>

      {/* Help */}
      <div className="alert alert-info">
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
          <h3 className="font-bold">Информације о 3D Secure</h3>
          <div className="text-sm">
            3D Secure је додатни слој безбедности за онлајн плаћања. SMS код је
            послан на број телефона регистрован у вашој банци.
          </div>
        </div>
      </div>
    </div>
  );
}
