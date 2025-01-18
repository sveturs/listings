// frontend/hostel-frontend/src/contexts/LanguageContext.js
import React, { createContext, useState, useContext, useEffect } from 'react';

const LanguageContext = createContext();

const SUPPORTED_LANGUAGES = [
    { code: 'en', name: 'English', flag: '🇬🇧' },
    { code: 'sr', name: 'Српски', flag: '🇷🇸' },
    { code: 'ru', name: 'Русский', flag: '🇷🇺' }
];

export const LanguageProvider = ({ children }) => {
    const [language, setLanguage] = useState(() => {
        // Пытаемся получить язык из localStorage
        const savedLang = localStorage.getItem('preferredLanguage');
        if (savedLang && SUPPORTED_LANGUAGES.some(lang => lang.code === savedLang)) {
            return savedLang;
        }

        // Или определяем язык браузера
        const browserLang = navigator.language.split('-')[0];
        if (SUPPORTED_LANGUAGES.some(lang => lang.code === browserLang)) {
            return browserLang;
        }

        return 'en'; // По умолчанию английский
    });

    useEffect(() => {
        localStorage.setItem('preferredLanguage', language);
        document.documentElement.lang = language;
    }, [language]);

    const value = {
        language,
        setLanguage,
        supportedLanguages: SUPPORTED_LANGUAGES
    };

    return (
        <LanguageContext.Provider value={value}>
            {children}
        </LanguageContext.Provider>
    );
};

export const useLanguage = () => {
    const context = useContext(LanguageContext);
    if (!context) {
        throw new Error('useLanguage must be used within a LanguageProvider');
    }
    return context;
};