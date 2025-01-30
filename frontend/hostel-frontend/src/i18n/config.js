// frontend/hostel-frontend/src/i18n/config.js
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import Backend from 'i18next-http-backend';
import LanguageDetector from 'i18next-browser-languagedetector';
import commonRu from '../locales/ru/common.json';
import commonEn from '../locales/en/common.json';
import commonSr from '../locales/sr/common.json';

const resources = {
  ru: {
    common: commonRu
  },
  en: {
    common: commonEn
  },
  sr: {
    common: commonSr
  }
};

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    debug: process.env.NODE_ENV === 'development',
    fallbackLng: 'en',
    supportedLngs: ['en', 'ru', 'sr'],
    
    ns: ['common', 'forms', 'navigation', 'marketplace', 'errors'],
    defaultNS: 'common',
    
    interpolation: {
      escapeValue: false,
    },

    backend: {
      loadPath: '/locales/{{lng}}/{{ns}}.json',
    },

    detection: {
      order: ['localStorage', 'navigator'],
      caches: ['localStorage'],
    },

    react: {
      useSuspense: true,
    },
  });

export default i18n;