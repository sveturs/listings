'use client';

import { useSearchParams, useRouter } from 'next/navigation';
// import { useTranslations } from 'next-intl';
import Link from 'next/link';

export default function PaymentErrorPage() {
  const router = useRouter();
  const searchParams = useSearchParams();

  const paymentId = searchParams.get('id');
  const reason = searchParams.get('reason') || 'unknown';

  const getErrorInfo = () => {
    const errorMap: Record<
      string,
      {
        title: string;
        description: string;
        suggestions: string[];
        icon: string;
      }
    > = {
      declined: {
        icon: '❌',
        title: 'Платеж отклонен',
        description:
          'Ваш банк отклонил транзакцию. Это может произойти по различным причинам.',
        suggestions: [
          'Проверьте правильность данных карты',
          'Убедитесь в наличии достаточных средств',
          'Обратитесь в ваш банк для уточнения',
        ],
      },
      insufficient_funds: {
        icon: '💳',
        title: 'Недостаточно средств',
        description: 'На карте недостаточно средств для совершения покупки.',
        suggestions: ['Проверьте баланс карты', 'Попробуйте другую карту'],
      },
      timeout: {
        icon: '⏰',
        title: 'Время ожидания истекло',
        description: 'Транзакция прервана из-за превышения времени ожидания.',
        suggestions: ['Проверьте интернет соединение', 'Попробуйте еще раз'],
      },
      cancelled: {
        icon: '🚫',
        title: 'Платеж отменен',
        description: 'Вы отменили процесс оплаты.',
        suggestions: [],
      },
      '3ds_failed': {
        icon: '🔐',
        title: '3D Secure неуспешно',
        description: 'Не удалось пройти аутентификацию 3D Secure.',
        suggestions: [
          'Проверьте правильность введенного кода',
          'Убедитесь что SMS дошло на правильный номер',
        ],
      },
      no_payment_id: {
        icon: '❓',
        title: 'Ошибка обработки',
        description: 'Не удалось найти данные платежа.',
        suggestions: ['Попробуйте начать процесс оплаты заново'],
      },
      unknown: {
        icon: '⚠️',
        title: 'Неизвестная ошибка',
        description: 'Произошла неожиданная ошибка при обработке платежа.',
        suggestions: ['Попробуйте еще раз', 'Обратитесь в службу поддержки'],
      },
    };

    return errorMap[reason] || errorMap.unknown;
  };

  const errorInfo = getErrorInfo();

  return (
    <div className="min-h-screen flex items-center justify-center p-4 bg-gradient-to-br from-error/10 to-base-100">
      <div className="max-w-md w-full">
        <div className="card bg-base-100 shadow-xl">
          <div className="card-body text-center">
            {/* Error Icon */}
            <div className="mb-4">
              <div className="text-6xl mb-2">{errorInfo.icon}</div>
            </div>

            {/* Title */}
            <h1 className="text-2xl font-bold text-base-content mb-2">
              {errorInfo.title}
            </h1>

            {/* Description */}
            <p className="text-base-content/70 mb-6">{errorInfo.description}</p>

            {/* Error Details */}
            {paymentId && (
              <div className="bg-base-200 rounded-lg p-4 mb-6 text-sm">
                <p className="text-base-content/60">
                  Номер для справки:{' '}
                  <span className="font-mono">{paymentId}</span>
                </p>
                <p className="text-base-content/60 mt-1">
                  Код ошибки:{' '}
                  <span className="font-mono">{reason.toUpperCase()}</span>
                </p>
              </div>
            )}

            {/* Suggestions */}
            {errorInfo.suggestions.length > 0 && (
              <div className="text-left mb-6">
                <h3 className="font-semibold mb-2">Что можно сделать:</h3>
                <ul className="space-y-1">
                  {errorInfo.suggestions.map((suggestion, index) => (
                    <li key={index} className="flex items-start gap-2">
                      <span className="text-primary">•</span>
                      <span className="text-sm text-base-content/70">
                        {suggestion}
                      </span>
                    </li>
                  ))}
                </ul>
              </div>
            )}

            {/* Actions */}
            <div className="card-actions justify-center gap-2">
              <button className="btn btn-primary" onClick={() => router.back()}>
                Попробовать снова
              </button>
              <Link href="/" className="btn btn-ghost">
                На главную
              </Link>
            </div>

            {/* Support */}
            <div className="mt-6 pt-6 border-t border-base-300">
              <p className="text-sm text-base-content/60">
                Нужна помощь?{' '}
                <Link href="/support" className="link link-primary">
                  Свяжитесь с поддержкой
                </Link>
              </p>
            </div>
          </div>
        </div>

        {/* Additional Help */}
        <div className="mt-4 p-4 bg-base-100/80 backdrop-blur rounded-lg">
          <h3 className="font-semibold mb-2">Часто задаваемые вопросы</h3>
          <div className="space-y-2 text-sm text-base-content/70">
            <details className="collapse collapse-arrow">
              <summary className="collapse-title text-sm font-medium">
                Почему мой платеж отклонен?
              </summary>
              <div className="collapse-content">
                <p>
                  Платеж может быть отклонен банком по соображениям
                  безопасности, недостатке средств или технических проблемах.
                </p>
              </div>
            </details>

            <details className="collapse collapse-arrow">
              <summary className="collapse-title text-sm font-medium">
                Когда будут возвращены средства?
              </summary>
              <div className="collapse-content">
                <p>
                  Если средства были списаны, они вернутся на карту в течение
                  3-7 рабочих дней.
                </p>
              </div>
            </details>
          </div>
        </div>
      </div>
    </div>
  );
}
