'use client';

import React, { Component, ReactNode } from 'react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

export class ChatErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false };
  }

  static getDerivedStateFromError(error: Error): State {
    // Update state so the next render will show the fallback UI
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: React.ErrorInfo) {
    // Log error to console in development
    if (process.env.NODE_ENV === 'development') {
      console.error('Chat Error Boundary caught:', error, errorInfo);
    }

    // TODO: Send error to monitoring service (e.g., Sentry)
    // logErrorToService(error, errorInfo);
  }

  handleReset = () => {
    this.setState({ hasError: false, error: undefined });
  };

  render() {
    if (this.state.hasError) {
      // Custom fallback UI with DaisyUI styling
      if (this.props.fallback) {
        return this.props.fallback;
      }

      return (
        <div className="flex items-center justify-center min-h-[400px] p-4">
          <div className="card bg-base-100 shadow-xl max-w-md w-full">
            <div className="card-body text-center">
              {/* Error icon */}
              <div className="flex justify-center mb-4">
                <div className="w-20 h-20 rounded-full bg-error/10 flex items-center justify-center">
                  <svg
                    className="w-10 h-10 text-error"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke="currentColor"
                  >
                    <path
                      strokeLinecap="round"
                      strokeLinejoin="round"
                      strokeWidth={2}
                      d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"
                    />
                  </svg>
                </div>
              </div>

              <h2 className="card-title justify-center text-error">
                Что-то пошло не так
              </h2>

              <p className="text-base-content/70 mb-4">
                Произошла ошибка при загрузке чата. Попробуйте обновить страницу
                или вернуться позже.
              </p>

              {/* Error details in development */}
              {process.env.NODE_ENV === 'development' && this.state.error && (
                <div className="collapse collapse-arrow bg-base-200 text-left mb-4">
                  <input type="checkbox" />
                  <div className="collapse-title text-sm font-medium">
                    Детали ошибки (dev only)
                  </div>
                  <div className="collapse-content">
                    <pre className="text-xs overflow-auto p-2 bg-base-300 rounded">
                      {this.state.error.message}
                      {'\n'}
                      {this.state.error.stack}
                    </pre>
                  </div>
                </div>
              )}

              <div className="card-actions justify-center">
                <button onClick={this.handleReset} className="btn btn-primary">
                  Попробовать снова
                </button>
                <button
                  onClick={() => window.location.reload()}
                  className="btn btn-ghost"
                >
                  Обновить страницу
                </button>
              </div>
            </div>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ChatErrorBoundary;
