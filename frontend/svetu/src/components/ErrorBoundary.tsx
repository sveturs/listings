'use client';

import React, { Component, ErrorInfo, ReactNode, useMemo } from 'react';
import { useTranslations } from 'next-intl';

interface ErrorMessages {
  title: string;
  description: string;
  details: string;
  reload: string;
}

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
  messages?: ErrorMessages;
  name?: string;
}

interface State {
  hasError: boolean;
  error?: Error;
}

class ErrorBoundaryClass extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    const boundaryName = this.props.name || 'ErrorBoundary';
    console.error(`[${boundaryName}] caught an error:`, error, errorInfo);
  }

  render() {
    if (this.state.hasError) {
      const defaultMessages: ErrorMessages = {
        title: 'Что-то пошло не так',
        description:
          'Произошла ошибка при загрузке. Пожалуйста, попробуйте обновить страницу.',
        details: 'Подробности ошибки',
        reload: 'Обновить страницу',
      };

      return (
        this.props.fallback || (
          <DefaultErrorFallback
            error={this.state.error}
            messages={this.props.messages || defaultMessages}
          />
        )
      );
    }

    return this.props.children;
  }
}

function DefaultErrorFallback({
  error,
  messages,
}: {
  error?: Error;
  messages: ErrorMessages;
}) {
  const handleReload = () => {
    window.location.reload();
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-base-200">
      <div className="card w-96 bg-base-100 shadow-xl">
        <div className="card-body text-center">
          <h2 className="card-title justify-center text-error">
            {messages.title}
          </h2>
          <p className="text-base-content/70">{messages.description}</p>
          {error && (
            <details className="text-xs text-left mt-2" role="group">
              <summary
                className="cursor-pointer text-base-content/50"
                aria-controls="error-details"
              >
                {messages.details}
              </summary>
              <div id="error-details" role="region">
                <pre className="mt-2 p-2 bg-base-200 rounded text-error">
                  {error.message}
                </pre>
              </div>
            </details>
          )}
          <div className="card-actions justify-center mt-4">
            <button className="btn btn-primary" onClick={handleReload}>
              {messages.reload}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}

export function AuthErrorBoundary({ children }: { children: ReactNode }) {
  const t = useTranslations('common');

  const messages: ErrorMessages = useMemo(
    () => ({
      title: t('errors.authError.title'),
      description: t('errors.authError.description'),
      details: t('errors.authError.details'),
      reload: t('errors.authError.reload'),
    }),
    [t]
  );

  return (
    <ErrorBoundaryClass messages={messages}>{children}</ErrorBoundaryClass>
  );
}

// Экспортируем класс для использования в других местах
export default ErrorBoundaryClass;
