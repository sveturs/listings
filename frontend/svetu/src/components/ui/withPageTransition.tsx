import { ComponentType } from 'react';
import { PageTransition } from './PageTransition';

export type TransitionMode = 'fade' | 'slide' | 'scale' | 'slideUp';

interface WithPageTransitionOptions {
  mode?: TransitionMode;
  duration?: number;
}

export function withPageTransition<P extends object>(
  Component: ComponentType<P>,
  options: WithPageTransitionOptions = {}
) {
  const WrappedComponent = (props: P) => {
    return (
      <PageTransition mode={options.mode} duration={options.duration}>
        <Component {...props} />
      </PageTransition>
    );
  };

  WrappedComponent.displayName = `withPageTransition(${Component.displayName || Component.name})`;

  return WrappedComponent;
}
