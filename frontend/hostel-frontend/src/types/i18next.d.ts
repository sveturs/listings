import 'i18next';

// Enhanced type definitions for i18next to fix TypeScript errors
declare module 'i18next' {
  interface i18n {
    changeLanguage(lng?: string): Promise<TFunction>;
    language: string;
    services: any;
    exists(key: string, options?: any): boolean;
    getFixedT(lng: string, ns?: string | string[]): TFunction;
    use(module: any): i18n;
  }

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