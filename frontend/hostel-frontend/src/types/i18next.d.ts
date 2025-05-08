import 'i18next';

// declare module 'i18next' to fix some type errors
declare module 'i18next' {
  interface i18n {
    changeLanguage(lng?: string): Promise<TFunction>;
    language: string;
  }
}