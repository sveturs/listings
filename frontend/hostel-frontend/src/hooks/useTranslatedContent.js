// frontend/hostel-frontend/src/hooks/useTranslatedContent.js
import { useMemo } from 'react';
import { useLanguage } from '../contexts/LanguageContext';
export const useTranslatedContent = (content) => {
    const { language } = useLanguage();

    return useMemo(() => {
        if (!content || !content.translations || !content.original_language) {
            return content;
        }

        // Если текущий язык совпадает с оригинальным - возвращаем контент без изменений
        if (language === content.original_language) {
            return content;
        }

        // Создаем копию контента
        const translatedContent = { ...content };

        // Если есть перевод для текущего языка, применяем его
        if (content.translations[language]) {
            Object.entries(content.translations[language]).forEach(([field, text]) => {
                translatedContent[field] = text;
            });
        }

        return translatedContent;
    }, [content, language]);
};