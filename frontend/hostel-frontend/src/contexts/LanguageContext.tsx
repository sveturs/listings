// frontend/hostel-frontend/src/contexts/LanguageContext.tsx
import React, { createContext, useState, useContext, useEffect, ReactNode } from 'react';
import i18n from '../i18n/config';

// Типы для языкового контекста
export interface Language {
  code: string;
  name: string;
  flag: string;
}

export interface LanguageContextType {
  language: string;
  setLanguage: (lang: string) => void;
  supportedLanguages: Language[];
}

interface LanguageProviderProps {
  children: ReactNode;
}

// Создаем контекст с начальным значением null
const LanguageContext = createContext<LanguageContextType | null>(null);

export const LanguageProvider: React.FC<LanguageProviderProps> = ({ children }) => {
  // Определяем поддерживаемые языки
  const supportedLanguages: Language[] = [
    { code: 'en', name: 'English', flag: '🇬🇧' },
    { code: 'sr', name: 'Српски', flag: '🇷🇸' },
    { code: 'ru', name: 'Русский', flag: '🇷🇺' }
  ];

  const [language, setLanguage] = useState<string>(() => {
    // Проверяем сохраненный язык в localStorage
    const savedLang = localStorage.getItem('preferredLanguage');
    if (savedLang && supportedLanguages.some(lang => lang.code === savedLang)) {
      return savedLang;
    }
    
    // Проверяем язык браузера
    const browserLang = navigator.language.split('-')[0];
    if (supportedLanguages.some(lang => lang.code === browserLang)) {
      return browserLang;
    }
    
    // Возвращаем язык по умолчанию
    return 'sr';
  });

  useEffect(() => {
    localStorage.setItem('preferredLanguage', language);
    document.documentElement.lang = language;
    
    // Обновляем i18next
    i18n.changeLanguage(language);
  }, [language]);

  const value: LanguageContextType = {
    language,
    setLanguage,
    supportedLanguages
  };

  return (
    <LanguageContext.Provider value={value}>
      {children}
    </LanguageContext.Provider>
  );
};

export const useLanguage = (): LanguageContextType => {
  const context = useContext(LanguageContext);
  if (!context) {
    throw new Error('useLanguage must be used within a LanguageProvider');
  }
  return context;
};