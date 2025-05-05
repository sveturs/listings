// frontend/hostel-frontend/src/i18n/config.js
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import Backend from 'i18next-http-backend';
import LanguageDetector from 'i18next-browser-languagedetector';

// Импортируем все файлы переводов
import commonRu from '../locales/ru/common.json';
import commonEn from '../locales/en/common.json';
import commonSr from '../locales/sr/common.json';
import marketplaceRu from '../locales/ru/marketplace.json';
import marketplaceEn from '../locales/en/marketplace.json';
import marketplaceSr from '../locales/sr/marketplace.json';
import gisRu from '../locales/ru/gis.json';
import gisEn from '../locales/en/gis.json';
import gisSr from '../locales/sr/gis.json';

const resources = {
  ru: {
    common: commonRu,
    marketplace: marketplaceRu,
    gis: gisRu
  },
  en: {
    common: commonEn,
    marketplace: marketplaceEn,
    gis: gisEn
  },
  sr: {
    common: commonSr,
    marketplace: marketplaceSr,
    gis: gisSr
  }
};

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources,
    fallbackLng: 'sr',
    supportedLngs: ['en', 'ru', 'sr'],
    
    ns: ['common', 'marketplace', 'gis'],
    defaultNS: 'common',
    lng: 'sr', 
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