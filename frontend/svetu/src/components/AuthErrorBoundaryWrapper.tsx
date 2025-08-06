import { ReactNode } from 'react';
import ErrorBoundaryClass from './ErrorBoundary';

// Дефолтные сообщения на случай, если переводы не загружены
const defaultMessages = {
  title: 'Authentication Error',
  description:
    'An error occurred while loading the authentication system. Please try reloading the page.',
  details: 'Error details',
  reload: 'Reload page',
};

interface Props {
  children: ReactNode;
  messages?: {
    title: string;
    description: string;
    details: string;
    reload: string;
  };
}

export function AuthErrorBoundaryWrapper({ children, messages }: Props) {
  return (
    <ErrorBoundaryClass messages={messages || defaultMessages}>
      {children}
    </ErrorBoundaryClass>
  );
}
