import React from 'react';

declare global {
  namespace JSX {
    interface IntrinsicElements {
      'model-viewer': React.DetailedHTMLProps<
        React.HTMLAttributes<HTMLElement> & {
          src?: string;
          alt?: string;
          ar?: boolean;
          'ar-modes'?: string;
          'camera-controls'?: boolean;
          'auto-rotate'?: boolean;
          'shadow-intensity'?: string;
          'environment-image'?: string;
          'skybox-image'?: string;
          poster?: string;
          loading?: 'auto' | 'lazy' | 'eager';
          reveal?: 'auto' | 'manual';
          'ios-src'?: string;
          style?: React.CSSProperties;
          slot?: string;
        },
        HTMLElement
      >;
    }
  }
}

export {};
