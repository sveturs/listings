// Enhanced type definitions for i18next to fix TypeScript errors
declare module 'i18next' {
  // Suppress TypeScript checking for the CustomTypeOptions interface
  interface CustomTypeOptions {
    defaultNS: 'common';
    resources: {
      common: {
        buttons: Record<string, string>;
        common: Record<string, string>;
      };
      gis: {
        layers: Record<string, string>;
        filters: Record<string, string>;
        categories: Record<string, string>;
        search: Record<string, string>;
      };
      marketplace: {
        store: {
          import: Record<string, string>;
          categoryMapping: Record<string, string>;
        };
      };
    };
  }

  // i18n interface
  interface i18n {
    changeLanguage(lng?: string): Promise<TFunction>;
    language: string;
    services: any;
    exists(key: string, options?: any): boolean;
    getFixedT(lng: string, ns?: string | string[]): TFunction;
    use(module: any): i18n;
  }

  // Translation function
  interface TFunction {
    (key: string | string[], options?: any): string;
    (key: string | string[], defaultValue: string, options?: any): string;
  }
}

// Additional declarations for react-i18next
declare module 'react-i18next' {
  import { TFunction, i18n } from 'i18next';

  export function useTranslation(ns?: string | string[], options?: any): {
    t: TFunction;
    i18n: i18n;
    ready: boolean;
  };

  export interface WithTranslation {
    t: TFunction;
    i18n: i18n;
    tReady: boolean;
  }

  export function withTranslation(ns?: string | string[], options?: any):
    <P extends WithTranslation>(component: React.ComponentType<P>) => React.ComponentType<P>;
}