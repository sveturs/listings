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

  // Type for i18next init options
  interface InitOptions {
    resources?: Record<string, Record<string, Record<string, any>>>;
    fallbackLng?: string | string[] | false;
    supportedLngs?: string[];
    ns?: string | string[];
    defaultNS?: string;
    lng?: string;
    interpolation?: {
      escapeValue?: boolean;
      [key: string]: any;
    };
    backend?: {
      loadPath?: string;
      [key: string]: any;
    };
    detection?: {
      order?: string[];
      caches?: string[];
      [key: string]: any;
    };
    react?: {
      useSuspense?: boolean;
      [key: string]: any;
    };
    [key: string]: any;
  }

  // i18n interface
  interface i18n {
    changeLanguage(lng?: string): Promise<TFunction>;
    language: string;
    services: any;
    exists(key: string, options?: any): boolean;
    getFixedT(lng: string, ns?: string | string[]): TFunction;
    use(module: any): i18n;
    init(options: InitOptions): i18n;
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
  
  export const initReactI18next: any;

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