// frontend/hostel-frontend/src/contexts/LanguageContext.js
import React, { createContext, useState, useContext, useEffect } from 'react';
import i18n from '../i18n/config';
const LanguageContext = createContext();




export const LanguageProvider = ({ children }) => {
    // Перемещаем определение внутрь компонента
    const supportedLanguages = [
        { code: 'en', name: 'English', flag: '🇬🇧' },
        { code: 'sr', name: 'Српски', flag: '🇷🇸' },
        { code: 'ru', name: 'Русский', flag: '🇷🇺' }
    ];

    const [language, setLanguage] = useState(() => {
        
        const savedLang = localStorage.getItem('preferredLanguage');
        if (savedLang && supportedLanguages.some(lang => lang.code === savedLang)) {
            return savedLang;
        }
        
        const browserLang = navigator.language.split('-')[0];
        if (supportedLanguages.some(lang => lang.code === browserLang)) {
            return browserLang;
        }
        
        return 'sr';
    });

    useEffect(() => {
        localStorage.setItem('preferredLanguage', language);
        document.documentElement.lang = language;
        
        // Обновляем i18next
        i18n.changeLanguage(language);
    }, [language]);

    const value = {
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

export const useLanguage = () => {
    const context = useContext(LanguageContext);
    if (!context) {
        throw new Error('useLanguage must be used within a LanguageProvider');
    }
    return context;
};