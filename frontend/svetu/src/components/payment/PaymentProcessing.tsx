// import { useTranslations } from 'next-intl';

interface PaymentProcessingProps {
  status: string;
  attempts: number;
  maxAttempts: number;
}

export default function PaymentProcessing({
  status,
  attempts,
  maxAttempts,
}: PaymentProcessingProps) {
  const getStatusInfo = () => {
    switch (status) {
      case 'checking':
        return {
          icon: '🔍',
          title: 'Проверка платежа',
          description: 'Получаем информацию о статусе транзакции',
          color: 'text-info',
        };
      case 'pending':
        return {
          icon: '⏳',
          title: 'Ожидание обработки',
          description: 'Платеж находится в очереди на обработку',
          color: 'text-warning',
        };
      case 'processing':
        return {
          icon: '⚡',
          title: 'Обработка платежа',
          description: 'Банк обрабатывает ваш платеж',
          color: 'text-primary',
        };
      case 'authorized':
        return {
          icon: '✅',
          title: 'Платеж авторизован',
          description: 'Средства заблокированы, завершаем транзакцию',
          color: 'text-success',
        };
      default:
        return {
          icon: '❓',
          title: status.charAt(0).toUpperCase() + status.slice(1),
          description: 'Обрабатываем статус платежа',
          color: 'text-base-content',
        };
    }
  };

  const info = getStatusInfo();
  const progress = Math.min((attempts / maxAttempts) * 100, 100);

  return (
    <div className="card bg-base-100 shadow-xl">
      <div className="card-body text-center">
        {/* Animated Icon */}
        <div className={`text-6xl mb-4 ${info.color} animate-pulse`}>
          {info.icon}
        </div>

        {/* Title */}
        <h2 className="card-title justify-center text-2xl">{info.title}</h2>

        {/* Description */}
        {info.description && (
          <p className="text-base-content/70 mt-2">{info.description}</p>
        )}

        {/* Progress Bar */}
        <div className="w-full mt-6">
          <div className="flex justify-between text-sm text-base-content/60 mb-2">
            <span>Проверяем статус</span>
            <span>
              {attempts}/{maxAttempts}
            </span>
          </div>
          <progress
            className="progress progress-primary w-full"
            value={progress}
            max="100"
          />
        </div>

        {/* Loading Animation */}
        <div className="flex justify-center gap-1 mt-6">
          <div
            className="w-2 h-2 bg-primary rounded-full animate-bounce"
            style={{ animationDelay: '0ms' }}
          ></div>
          <div
            className="w-2 h-2 bg-primary rounded-full animate-bounce"
            style={{ animationDelay: '150ms' }}
          ></div>
          <div
            className="w-2 h-2 bg-primary rounded-full animate-bounce"
            style={{ animationDelay: '300ms' }}
          ></div>
        </div>

        {/* Info Message */}
        <div className="alert alert-info mt-6">
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
          <span className="text-sm">
            Пожалуйста, не закрывайте это окно во время обработки платежа
          </span>
        </div>

        {/* Payment Steps */}
        <div className="mt-6">
          <h3 className="font-semibold mb-3">Этапы обработки:</h3>
          <ul className="steps steps-vertical lg:steps-horizontal w-full">
            <li
              className={`step ${['checking', 'pending', 'processing', 'authorized'].includes(status) ? 'step-primary' : ''}`}
            >
              Проверка
            </li>
            <li
              className={`step ${['pending', 'processing', 'authorized'].includes(status) ? 'step-primary' : ''}`}
            >
              Обработка
            </li>
            <li
              className={`step ${['processing', 'authorized'].includes(status) ? 'step-primary' : ''}`}
            >
              Авторизация
            </li>
            <li
              className={`step ${status === 'authorized' ? 'step-primary' : ''}`}
            >
              Завершение
            </li>
          </ul>
        </div>
      </div>
    </div>
  );
}
